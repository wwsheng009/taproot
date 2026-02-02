// Package shell provides shell command execution utilities.
package shell

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"mvdan.cc/sh/v3/shell"
)

// Executor manages shell command execution.
type Executor struct {
	runningProcesses map[int]*Process
	nextPID          int
	mu               sync.RWMutex
}

// NewExecutor creates a new shell executor.
func NewExecutor() *Executor {
	return &Executor{
		runningProcesses: make(map[int]*Process),
		nextPID:          1,
	}
}

// Execute runs a command synchronously and returns the result.
func (e *Executor) Execute(command string, args ...string) (CommandResult, error) {
	return e.ExecuteWithOptions(command, args, DefaultCommandOptions())
}

// ExecuteWithOptions runs a command with custom options.
func (e *Executor) ExecuteWithOptions(command string, args []string, opts CommandOptions) (CommandResult, error) {
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Validate working directory
	if opts.WorkingDir != "" {
		if !e.dirExists(opts.WorkingDir) {
			result := CommandResult{
				Error: fmt.Errorf("%w: %s", ErrWorkingDirNotFound, opts.WorkingDir),
			}
			return result, result.Error
		}
	}

	var cmd *exec.Cmd
	if opts.Type == CommandShell {
		// Parse shell command
		fields, err := shell.Fields(command, nil)
		if err != nil {
			result := CommandResult{
				Error: err,
			}
			return result, result.Error
		}
		if len(fields) == 0 {
			result := CommandResult{
				Error: ErrEmptyCommand,
			}
			return result, result.Error
		}

		cmd = exec.CommandContext(ctx, fields[0], fields[1:]...)
	} else {
		// Direct execution
		if command == "" {
			result := CommandResult{
				Error: ErrEmptyCommand,
			}
			return result, result.Error
		}
		cmd = exec.CommandContext(ctx, command, args...)
	}

	// Set working directory
	cmd.Dir = opts.WorkingDir

	// Set environment
	if len(opts.Env) > 0 {
		cmd.Env = append(os.Environ(), opts.Env...)
	}

	// Set stdin
	if opts.Stdin != nil {
		cmd.Stdin = opts.Stdin
	}

	// Set up output capture
	var stdout, stderr io.ReadCloser
	var stdoutBuf, stderrBuf strings.Builder

	if opts.Direction == OutputStdout || opts.Direction == OutputBoth {
		var err error
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			result := CommandResult{
				Error: fmt.Errorf("failed to create stdout pipe: %w", err),
			}
			return result, result.Error
		}
	}

	if opts.Direction == OutputStderr || opts.Direction == OutputBoth {
		if !opts.RedirectStderrToStdout {
			var err error
			stderr, err = cmd.StderrPipe()
			if err != nil {
				result := CommandResult{
					Error: fmt.Errorf("failed to create stderr pipe: %w", err),
				}
				return result, result.Error
			}
		}
	}

	// Merge stderr to stdout if requested
	if opts.RedirectStderrToStdout {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stdoutBuf
	}

	// Start the command
	startTime := time.Now()

	// Progress update handling
	var stdoutReader, stderrReader *bufio.Reader
	if stdout != nil {
		stdoutReader = bufio.NewReader(stdout)
	}
	if stderr != nil {
		stderrReader = bufio.NewReader(stderr)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		result := CommandResult{
			Error:    fmt.Errorf("failed to start command: %w", err),
			Duration: time.Since(startTime),
		}
		return result, result.Error
	}

	// Capture output in real-time with progress updates
	var wg sync.WaitGroup
	outputMu := sync.Mutex{}

	// Capture stdout
	if stdoutReader != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				line, err := stdoutReader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						if !opts.Quiet {
							slog.Error("Error reading stdout", "error", err)
						}
					}
					break
				}
				outputMu.Lock()
				stdoutBuf.WriteString(line)
				outputMu.Unlock()

				// Send progress update
				if opts.Progress != nil {
					opts.Progress(ProgressUpdate{
						Timestamp: time.Now(),
						Stdout:    line,
						State:     StateRunning,
					})
				}
			}
		}()
	}

	// Capture stderr
	if stderrReader != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				line, err := stderrReader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						if !opts.Quiet {
							slog.Error("Error reading stderr", "error", err)
						}
					}
					break
				}
				outputMu.Lock()
				stderrBuf.WriteString(line)
				outputMu.Unlock()

				// Send progress update
				if opts.Progress != nil {
					opts.Progress(ProgressUpdate{
						Timestamp: time.Now(),
						Stderr:    line,
						State:     StateRunning,
					})
				}
			}
		}()
	}

	// Wait for the command
	err := cmd.Wait()
	wg.Wait()

	duration := time.Since(startTime)

	result := CommandResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		Duration: duration,
	}

	// Check for context errors
	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.Error = ErrTimeout
		if !opts.Quiet {
			slog.Warn("Command timed out", "command", command, "timeout", opts.Timeout)
		}
		return result, result.Error
	}

	if ctx.Err() == context.Canceled {
		result.Cancelled = true
		result.Error = ErrCancelled
		if !opts.Quiet {
			slog.Info("Command cancelled", "command", command)
		}
		return result, result.Error
	}

	// Handle execution error
	if err != nil {
		result.Error = err
		// Get exit code if available
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	} else {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	return result, result.Error
}

// ExecuteAsync starts a command asynchronously and returns a Process.
func (e *Executor) ExecuteAsync(command string, args ...string) (*Process, error) {
	return e.ExecuteAsyncWithOptions(command, args, DefaultCommandOptions())
}

// ExecuteAsyncWithOptions starts a command asynchronously with custom options.
func (e *Executor) ExecuteAsyncWithOptions(command string, args []string, opts CommandOptions) (*Process, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	pid := e.nextPID
	e.nextPID++

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)

	process := &Process{
		cmd: &ProcessInfo{
			Command: command,
			Args:    args,
			Options: opts,
		},
		startTime: time.Now(),
		state:     StateIdle,
		cancel:    cancel,
	}

	e.runningProcesses[pid] = process

	// Start async execution
	go e.runProcess(process, pid)

	return process, nil
}

// runProcess executes the process in a goroutine.
func (e *Executor) runProcess(process *Process, pid int) {
	process.state = StateRunning

	defer func() {
		e.mu.Lock()
		delete(e.runningProcesses, pid)
		e.mu.Unlock()
	}()

	// Execute the command
	result, err := e.ExecuteWithOptions(process.cmd.Command, process.cmd.Args, process.cmd.Options)

	process.result = result

	if err != nil {
		if result.TimedOut {
			process.state = StateFailed
		} else if result.Cancelled {
			process.state = StateCancelled
		} else {
			process.state = StateFailed
		}
	} else {
		process.state = StateCompleted
	}
}

// Cancel cancels a running process by PID.
func (e *Executor) Cancel(pid int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	process, exists := e.runningProcesses[pid]
	if !exists {
		return ErrProcessNotFound
	}

	if process.state != StateRunning {
		return fmt.Errorf("process %d is not running (state: %s)", pid, process.state)
	}

	process.cancel()
	if opts := process.cmd.Options; opts.Progress != nil {
		opts.Progress(ProgressUpdate{
			Timestamp: time.Now(),
			State:     StateCancelled,
		})
	}

	return nil
}

// GetProcess returns a running process by PID.
func (e *Executor) GetProcess(pid int) (*Process, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	process, exists := e.runningProcesses[pid]
	if !exists {
		return nil, ErrProcessNotFound
	}

	return process, nil
}

// ListProcesses returns all running processes.
func (e *Executor) ListProcesses() []*Process {
	e.mu.RLock()
	defer e.mu.RUnlock()

	processes := make([]*Process, 0, len(e.runningProcesses))
	for _, process := range e.runningProcesses {
		processes = append(processes, process)
	}

	return processes
}

// CancelAll cancels all running processes.
func (e *Executor) CancelAll() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, process := range e.runningProcesses {
		if process.state == StateRunning {
			process.cancel()
			if opts := process.cmd.Options; opts.Progress != nil {
				opts.Progress(ProgressUpdate{
					Timestamp: time.Now(),
					State:     StateCancelled,
				})
			}
		}
	}
}

// CountRunning returns the number of running processes.
func (e *Executor) CountRunning() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	count := 0
	for _, process := range e.runningProcesses {
		if process.state == StateRunning {
			count++
		}
	}

	return count
}

// dirExists checks if a directory exists.
func (e *Executor) dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Quiet executes a command quietly (no error logging, output capture).
func Quiet(command string, args ...string) (string, error) {
	result, err := NewExecutor().ExecuteWithOptions(command, args, CommandOptions{
		Type:      CommandDirect,
		Direction: OutputStdout,
		Quiet:     true,
	})
	if err != nil {
		return "", err
	}
	return result.Stdout, nil
}

// Shell executes a command through the system shell (sh/bash/cmd.exe).
func Shell(command string) (CommandResult, error) {
	return NewExecutor().ExecuteWithOptions(command, nil, CommandOptions{
		Type:      CommandShell,
		Direction: OutputBoth,
	})
}

// Quick is a convenience function for quick synchronous execution with output capture.
func Quick(command string, args ...string) (string, error) {
	result, err := NewExecutor().ExecuteWithOptions(command, args, CommandOptions{
		Type:      CommandDirect,
		Direction: OutputStdout,
	})
	if err != nil {
		return "", err
	}
	return result.Stdout, nil
}

// Which finds the path to an executable.
func Which(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

// GetDefaultShell returns the default shell for the current platform.
func GetDefaultShell() string {
	switch runtime.GOOS {
	case "windows":
		// Prefer PowerShell if available, otherwise cmd.exe
		if path, err := Which("powershell"); err == nil {
			return path
		}
		return "cmd.exe"
	case "darwin", "linux", "freebsd", "openbsd", "netbsd":
		// Prefer bash if available, otherwise sh
		if path, err := Which("bash"); err == nil {
			return path
		}
		return "/bin/sh"
	default:
		return "/bin/sh"
	}
}

// JoinPath joins path elements safely for the current platform.
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// ExpandPath expands ~ in paths to the user's home directory.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return filepath.Clean(path)
}

// SplitCommand splits a command into executable and arguments based on shell syntax.
func SplitCommand(command string) (string, []string, error) {
	fields, err := shell.Fields(command, nil)
	if err != nil {
		return "", nil, err
	}

	if len(fields) == 0 {
		return "", nil, ErrEmptyCommand
	}

	return fields[0], fields[1:], nil
}

// ReadOutput reads the complete output from a command synchronously.
func ReadOutput(command string, args ...string) (stdout, stderr string, err error) {
	result, err := NewExecutor().ExecuteWithOptions(command, args, CommandOptions{
		Type:      CommandDirect,
		Direction: OutputBoth,
	})
	return result.Stdout, result.Stderr, err
}

// Pipe executes multiple commands with piped output.
func Pipe(commands ...[]string) (CommandResult, error) {
	if len(commands) == 0 {
		result := CommandResult{
			Error: ErrEmptyCommand,
		}
		return result, result.Error
	}

	if len(commands) == 1 {
		return NewExecutor().ExecuteWithOptions(commands[0][0], commands[0][1:], CommandOptions{
			Type:      CommandDirect,
			Direction: OutputBoth,
		})
	}

	// Create pipes
	prevStdout := new(bytes.Buffer)
	var result CommandResult

	for i, cmd := range commands {
		if len(cmd) == 0 {
			result.Error = ErrEmptyCommand
			return result, result.Error
		}

		execCmd := exec.Command(cmd[0], cmd[1:]...)

		if i == 0 {
			// First command: stdout to pipe
			execCmd.Stdout = prevStdout
		} else if i == len(commands)-1 {
			// Last command: read from pipe, capture stdout/stderr
			execCmd.Stdin = prevStdout
			var outBuf, errBuf bytes.Buffer
			execCmd.Stdout = &outBuf
			execCmd.Stderr = &errBuf

			startTime := time.Now()
			err := execCmd.Run()
			result.Duration = time.Since(startTime)
			result.Stdout = outBuf.String()
			result.Stderr = errBuf.String()

			if err != nil {
				result.Error = err
				if exitErr, ok := err.(*exec.ExitError); ok {
					result.ExitCode = exitErr.ExitCode()
				}
			} else {
				result.ExitCode = execCmd.ProcessState.ExitCode()
			}
		} else {
			// Middle commands: read from pipe, write to next pipe
			newStdout := new(bytes.Buffer)
			execCmd.Stdin = prevStdout
			execCmd.Stdout = newStdout

			if err := execCmd.Run(); err != nil && !result.TimedOut {
				result.Error = err
			}
			prevStdout = newStdout
		}
	}

	return result, result.Error
}
