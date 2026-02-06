// Package shell provides shell command execution utilities.
package shell

import (
	"context"
	"io"
	"sync"
	"time"
)

// CommandType represents the type of shell command.
type CommandType int

const (
	// CommandShell executes through a shell (sh/bash/cmd.exe)
	CommandShell CommandType = iota
	// CommandDirect executes the command directly (no shell)
	CommandDirect
)

// CommandDirection defines which output streams to capture.
type CommandDirection int

const (
	// OutputNone discards all output
	OutputNone CommandDirection = iota
	// OutputStdout captures only stdout
	OutputStdout
	// OutputStderr captures only stderr
	OutputStderr
	// OutputBoth captures both stdout and stderr
	OutputBoth
)

// ProcessState represents the current state of a running process.
type ProcessState int

const (
	// StateIdle is the initial state before execution
	StateIdle ProcessState = iota
	// StateRunning indicates the process is currently running
	StateRunning
	// StateCompleted indicates the process completed successfully
	StateCompleted
	// StateFailed indicates the process failed
	StateFailed
	// StateCancelled indicates the process was cancelled
	StateCancelled
)

// String returns the string representation of the process state.
func (s ProcessState) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateRunning:
		return "running"
	case StateCompleted:
		return "completed"
	case StateFailed:
		return "failed"
	case StateCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// CommandResult contains the result of a command execution.
type CommandResult struct {
	ExitCode   int           // Process exit code
	Stdout     string        // Captured stdout
	Stderr     string        // Captured stderr
	Duration   time.Duration // Execution duration
	Error      error         // Execution error (if any)
	Cancelled  bool          // Whether the command was cancelled
	TimedOut   bool          // Whether the command timed out
}

// ProgressUpdate represents progress updates during command execution.
type ProgressUpdate struct {
	Timestamp time.Time // Update timestamp
	Stdout    string    // Stdout output since last update
	Stderr    string    // Stderr output since last update
	State     ProcessState
}

// ProgressCallback is called during command execution for progress updates.
type ProgressCallback func(update ProgressUpdate)

// CommandOptions configures command execution behavior.
type CommandOptions struct {
	// Type determines how the command is executed
	Type CommandType

	// Direction specifies which streams to capture
	Direction CommandDirection

	// Working directory for the command
	WorkingDir string

	// Environment variables (will be merged with os.Environ())
	Env []string

	// Timeout for execution (0 = no timeout)
	Timeout time.Duration

	// Progress callback for real-time updates
	Progress ProgressCallback

	// Context for cancellation
	Context context.Context

	// Stdin provides input to the command
	Stdin io.Reader

	// Quiet mode suppresses error logging
	Quiet bool

	// RedirectStderrToStdout merges stderr into stdout
	RedirectStderrToStdout bool
}

// DefaultCommandOptions returns default command options.
func DefaultCommandOptions() CommandOptions {
	return CommandOptions{
		Type:                    CommandShell,
		Direction:               OutputBoth,
		Timeout:                 0,
		RedirectStderrToStdout:  false,
	}
}

// CommandBuilder provides a fluent interface for building shell commands.
type CommandBuilder struct {
	args    []string
	options CommandOptions
}

// NewCommandBuilder creates a new command builder with default options.
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		args:    make([]string, 0),
		options: DefaultCommandOptions(),
	}
}

// Command adds the executable and arguments.
func (b *CommandBuilder) Command(cmd string, args ...string) *CommandBuilder {
	b.args = append([]string{cmd}, args...)
	b.options.Type = CommandDirect
	return b
}

// ShellCommand adds a complete shell command string.
func (b *CommandBuilder) ShellCommand(command string) *CommandBuilder {
	b.args = []string{command}
	b.options.Type = CommandShell
	return b
}

// SetCommandType sets the command execution type.
func (b *CommandBuilder) SetCommandType(cmdType CommandType) *CommandBuilder {
	b.options.Type = cmdType
	return b
}

// SetWorkingDir sets the working directory.
func (b *CommandBuilder) SetWorkingDir(dir string) *CommandBuilder {
	b.options.WorkingDir = dir
	return b
}

// SetTimeout sets the execution timeout.
func (b *CommandBuilder) SetTimeout(timeout time.Duration) *CommandBuilder {
	b.options.Timeout = timeout
	return b
}

// SetEnv sets environment variables.
func (b *CommandBuilder) SetEnv(env ...string) *CommandBuilder {
	b.options.Env = env
	return b
}

// SetOutputDirection sets which streams to capture.
func (b *CommandBuilder) SetOutputDirection(dir CommandDirection) *CommandBuilder {
	b.options.Direction = dir
	return b
}

// SetProgressCallback sets the progress callback.
func (b *CommandBuilder) SetProgressCallback(cb ProgressCallback) *CommandBuilder {
	b.options.Progress = cb
	return b
}

// SetContext sets the execution context.
func (b *CommandBuilder) SetContext(ctx context.Context) *CommandBuilder {
	b.options.Context = ctx
	return b
}

// SetStdin sets the stdin input.
func (b *CommandBuilder) SetStdin(r io.Reader) *CommandBuilder {
	b.options.Stdin = r
	return b
}

// SetQuiet sets quiet mode.
func (b *CommandBuilder) SetQuiet(quiet bool) *CommandBuilder {
	b.options.Quiet = quiet
	return b
}

// RedirectStderrToStdout merges stderr into stdout.
func (b *CommandBuilder) RedirectStderrToStdout(merge bool) *CommandBuilder {
	b.options.RedirectStderrToStdout = merge
	return b
}

// Build validates and returns the command arguments and options.
func (b *CommandBuilder) Build() (string, []string, CommandOptions, error) {
	if len(b.args) == 0 {
		return "", nil, CommandOptions{}, ErrEmptyCommand
	}
	switch b.options.Type {
	case CommandDirect:
		if len(b.args) < 1 {
			return "", nil, CommandOptions{}, ErrMissingExecutable
		}
		return b.args[0], b.args[1:], b.options, nil
	case CommandShell:
		return b.args[0], nil, b.options, nil
	default:
		return "", nil, CommandOptions{}, ErrUnknownCommandType
	}
}

// Process represents a running command process.
type Process struct {
	cmd       *ProcessInfo
	startTime time.Time
	state     ProcessState
	stateMu   sync.RWMutex // Protects state field
	cancel    context.CancelFunc
	result    CommandResult
	resultMu  sync.RWMutex // Protects result field
}

// ProcessInfo holds internal process information.
type ProcessInfo struct {
	Command   string
	Args      []string
	Options   CommandOptions
	StdoutBuf string
	StderrBuf string
}

// State returns the current process state.
func (p *Process) State() ProcessState {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state
}

// Result returns the final execution result (only valid after completion).
func (p *Process) Result() CommandResult {
	p.resultMu.RLock()
	defer p.resultMu.RUnlock()
	return p.result
}

// Running returns true if the process is still running.
func (p *Process) Running() bool {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state == StateRunning
}

// Elapsed returns the time elapsed since the process started.
func (p *Process) Elapsed() time.Duration {
	p.stateMu.RLock()
	state := p.state
	p.stateMu.RUnlock()

	if state == StateIdle {
		return 0
	}
	return time.Since(p.startTime)
}

// Command returns the process command information.
func (p *Process) Command() *ProcessInfo {
	return p.cmd
}

// Errors
var (
	ErrEmptyCommand         = &Error{Code: "EMPTY_COMMAND", Message: "command is empty"}
	ErrMissingExecutable    = &Error{Code: "MISSING_EXEC", Message: "missing executable"}
	ErrUnknownCommandType   = &Error{Code: "UNKNOWN_TYPE", Message: "unknown command type"}
	ErrInvalidArgument      = &Error{Code: "INVALID_ARG", Message: "invalid argument"}
	ErrExecutionFailed      = &Error{Code: "EXEC_FAILED", Message: "command execution failed"}
	ErrTimeout              = &Error{Code: "TIMEOUT", Message: "command timed out"}
	ErrCancelled            = &Error{Code: "CANCELLED", Message: "command cancelled"}
	ErrProcessNotFound      = &Error{Code: "NOT_FOUND", Message: "process not found"}
	ErrWorkingDirNotFound   = &Error{Code: "WORKING_DIR_NOT_FOUND", Message: "working directory does not exist"}
	ErrCommandBuilderClosed = &Error{Code: "BUILDER_CLOSED", Message: "command builder is closed"}
)

// Error represents a shell execution error.
type Error struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	return e.Cause
}
