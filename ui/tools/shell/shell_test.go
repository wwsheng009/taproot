// Package shell provides shell command execution utilities.
package shell

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestProcessState_String(t *testing.T) {
	tests := []struct {
		state    ProcessState
		expected string
	}{
		{StateIdle, "idle"},
		{StateRunning, "running"},
		{StateCompleted, "completed"},
		{StateFailed, "failed"},
		{StateCancelled, "cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("ProcessState.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "simple error",
			err:  &Error{Code: "TEST", Message: "test message"},
			want: "test message",
		},
		{
			name: "error with cause",
			err:  &Error{Code: "TEST", Message: "test message", Cause: TestError{msg: "inner error"}},
			want: "test message: inner error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	innerErr := TestError{msg: "inner error"}
	err := &Error{Message: "outer error", Cause: innerErr}

	if got := err.Unwrap(); got != innerErr {
		t.Errorf("Error.Unwrap() = %v, want %v", got, innerErr)
	}

	errNoCause := &Error{Message: "no cause error"}
	if got := errNoCause.Unwrap(); got != nil {
		t.Errorf("Error.Unwrap() = %v, want nil", got)
	}
}

func TestDefaultCommandOptions(t *testing.T) {
	opts := DefaultCommandOptions()

	if opts.Type != CommandShell {
		t.Errorf("CommandOptions.Type = %v, want %v", opts.Type, CommandShell)
	}

	if opts.Direction != OutputBoth {
		t.Errorf("CommandOptions.Direction = %v, want %v", opts.Direction, OutputBoth)
	}

	if opts.Timeout != 0 {
		t.Errorf("CommandOptions.Timeout = %v, want 0", opts.Timeout)
	}
}

func TestCommandBuilder_Command(t *testing.T) {
	builder := NewCommandBuilder()
	builder.Command("echo", "hello", "world")

	cmd, args, opts, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if cmd != "echo" {
		t.Errorf("Command = %v, want echo", cmd)
	}

	if len(args) != 2 {
		t.Fatalf("Args len = %v, want 2", len(args))
	}

	if args[0] != "hello" || args[1] != "world" {
		t.Errorf("Args = %v, want [hello world]", args)
	}

	if opts.Type != CommandDirect {
		t.Errorf("Type = %v, want %v", opts.Type, CommandDirect)
	}
}

func TestCommandBuilder_ShellCommand(t *testing.T) {
	builder := NewCommandBuilder()
	builder.ShellCommand("echo 'hello world'")

	cmd, args, opts, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if cmd != "echo 'hello world'" {
		t.Errorf("Command = %v, want 'echo 'hello world''", cmd)
	}

	if args != nil {
		t.Errorf("Args = %v, want nil (shell command)", args)
	}

	if opts.Type != CommandShell {
		t.Errorf("Type = %v, want %v", opts.Type, CommandShell)
	}
}

func TestCommandBuilder_Build(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*CommandBuilder)
		wantErr    bool
		errMessage string
	}{
		{
			name: "empty command",
			setup: func(b *CommandBuilder) {
				// Don't add any command
			},
			wantErr:    true,
			errMessage: ErrEmptyCommand.Error(),
		},
		{
			name: "valid direct command",
			setup: func(b *CommandBuilder) {
				b.Command("echo", "test")
			},
			wantErr: false,
		},
		{
			name: "valid shell command",
			setup: func(b *CommandBuilder) {
				b.ShellCommand("echo test")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewCommandBuilder()
			tt.setup(builder)

			_, _, _, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		 if err != nil && tt.errMessage != "" {
				if !strings.Contains(err.Error(), tt.errMessage) {
					t.Errorf("Build() error = %v, want to contain %v", err.Error(), tt.errMessage)
				}
			}
		})
	}
}

func TestCommandBuilder_FluentInterface(t *testing.T) {
	cmd, args, opts, err := NewCommandBuilder().
		Command("ls", "-la").
		SetWorkingDir("/tmp").
		SetTimeout(5 * time.Second).
		SetEnv("TEST=value").
		SetOutputDirection(OutputStdout).
		SetQuiet(true).
		Build()

	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if cmd != "ls" {
		t.Errorf("Command = %v, want ls", cmd)
	}

	if args[0] != "-la" {
		t.Errorf("Args[0] = %v, want -la", args[0])
	}

	if opts.WorkingDir != "/tmp" {
		t.Errorf("WorkingDir = %v, want /tmp", opts.WorkingDir)
	}

	if opts.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want 5s", opts.Timeout)
	}

	if opts.Direction != OutputStdout {
		t.Errorf("Direction = %v, want %v", opts.Direction, OutputStdout)
	}

	if !opts.Quiet {
		t.Errorf("Quiet = %v, want true", opts.Quiet)
	}

	if len(opts.Env) != 1 || opts.Env[0] != "TEST=value" {
		t.Errorf("Env = %v, want [TEST=value]", opts.Env)
	}
}

func TestNewExecutor(t *testing.T) {
	exec := NewExecutor()

	if exec == nil {
		t.Fatal("NewExecutor() returned nil")
	}

	if exec.runningProcesses == nil {
		t.Fatal("NewExecutor() has nil processes map")
	}
}

func TestExecutor_Quick(t *testing.T) {
	// Find a reliable echo-like command for the platform
	cmd, args := "echo", []string{"hello world"}
	if runtime.GOOS == "windows" {
		cmd, args = "cmd.exe", []string{"/c", "echo hello world"}
	}

	// Test with a simple command that should work
	output, err := Quick(cmd, args...)
	if err != nil {
		t.Fatalf("Quick() error = %v", err)
	}

	if !strings.Contains(output, "hello world") {
		t.Errorf("Quick() output = %v, want to contain 'hello world'", output)
	}
}

func TestExecutor_Quiet(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to command availability")
	}
	// Test quiet execution
	output, err := Quiet("echo", "quiet test")
	if err != nil {
		t.Fatalf("Quiet() error = %v", err)
	}

	// Remove trailing newline
	output = strings.TrimSpace(output)
	if output != "quiet test" {
		t.Errorf("Quiet() output = %v, want 'quiet test'", output)
	}
}

func TestExecutor_ReadOutput(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to command availability")
	}
	stdout, stderr, err := ReadOutput("sh", "-c", "echo stdout line; echo stderr line >&2")
	if err != nil {
		t.Fatalf("ReadOutput() error = %v", err)
	}

	if !strings.Contains(stdout, "stdout line") {
		t.Errorf("stdout = %v, want to contain 'stdout line'", stdout)
	}

	if !strings.Contains(stderr, "stderr line") {
		t.Errorf("stderr = %v, want to contain 'stderr line'", stderr)
	}
}

func TestExecutor_Execute(t *testing.T) {
	exec := NewExecutor()

	tests := []struct {
		name    string
		command string
		args    []string
		wantOK  bool
		skipWin bool
	}{
		{
			name:    "successful command",
			command: "echo",
			args:    []string{"test"},
			wantOK:  true,
			skipWin: true, // echo not in path on Windows
		},
		{
			name:    "command with args",
			command: "echo",
			args:    []string{"-n", "no newline"},
			wantOK:  true,
			skipWin: true, // echo not in path on Windows
		},
		{
			name:    "failing command",
			command: "false",
			args:    nil,
			wantOK:  false, // false command exits with non-zero
			skipWin: true, // false not in path on Windows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipWin && runtime.GOOS == "windows" {
				t.Skip("Skipping on Windows due to command availability")
			}

			result, err := exec.Execute(tt.command, tt.args...)

			if tt.wantOK && err != nil {
				t.Errorf("Execute() error = %v, want no error", err)
			}

			if !tt.wantOK && err == nil && result.ExitCode == 0 {
				t.Errorf("Execute() expected error, got nil with exit code 0")
			}
		})
	}
}

func TestExecutor_ExecuteWithTimeout(t *testing.T) {
	exec := NewExecutor()

	// Test command that times out - use a long timeout on Windows
	timeout := 100 * time.Millisecond
	if runtime.GOOS == "windows" {
		timeout = 2 * time.Second // Windows ping command is slower
	}

	cmd, args := "sleep", []string{"5"}
	if runtime.GOOS == "windows" {
		cmd, args = "ping", []string{"-n", "6", "127.0.0.1"}
	}

	opts := CommandOptions{
		Type:      CommandDirect,
		Direction: OutputNone,
		Timeout:   timeout,
	}

	result, err := exec.ExecuteWithOptions(cmd, args, opts)

	if err == nil {
		t.Fatal("ExecuteWithTimeout() expected timeout error, got nil")
	}

	if !result.TimedOut {
		t.Errorf("ExecuteWithTimeout() TimedOut = %v, want true", result.TimedOut)
	}
}

func TestExecutor_ExecuteWithCancellation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to different process behavior")
	}

	exec := NewExecutor()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	opts := CommandOptions{
		Type:      CommandDirect,
		Direction: OutputNone,
		Context:   ctx,
	}

	result, err := exec.ExecuteWithOptions("sleep", []string{"5"}, opts)

	if err == nil {
		t.Fatal("ExecuteWithCancellation() expected error, got nil")
	}

	if !result.Cancelled {
		t.Errorf("ExecuteWithCancellation() Cancelled = %v, want true", result.Cancelled)
	}
}

func TestExecutor_ExecuteWithPathCapture(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to command availability")
	}

	exec := NewExecutor()

	tests := []struct {
		name      string
		direction CommandDirection
		wantStdout bool
		wantStderr bool
	}{
		{
			name:      "capture stdout only",
			direction: OutputStdout,
			wantStdout: true,
			wantStderr: false,
		},
		{
			name:      "capture stderr only",
			direction: OutputStderr,
			wantStdout: false,
			wantStderr: true,
		},
		{
			name:      "capture both",
			direction: OutputBoth,
			wantStdout: true,
			wantStderr: true,
		},
		{
			name:      "capture none",
			direction: OutputNone,
			wantStdout: false,
			wantStderr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := CommandOptions{
				Type:      CommandDirect,
				Direction: tt.direction,
			}

			result, err := exec.ExecuteWithOptions("sh", []string{"-c", "echo stdout; echo stderr >&2"}, opts)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			hasStdout := result.Stdout != ""
			hasStderr := result.Stderr != ""

			if hasStdout != tt.wantStdout {
				t.Errorf("Execute() hasStdout = %v, want %v", hasStdout, tt.wantStdout)
			}

			if hasStderr != tt.wantStderr {
				t.Errorf("Execute() hasStderr = %v, want %v", hasStderr, tt.wantStderr)
			}
		})
	}
}

func TestExecutor_ExecuteWithProgress(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to command availability")
	}

	exec := NewExecutor()

	progressUpdates := make([]ProgressUpdate, 0)
	opts := CommandOptions{
		Type:      CommandDirect,
		Direction: OutputStdout,
		Progress: func(update ProgressUpdate) {
			progressUpdates = append(progressUpdates, update)
		},
	}

	result, err := exec.ExecuteWithOptions("echo", []string{"line1", "line2"}, opts)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(progressUpdates) == 0 {
		t.Fatal("Execute() progress callback not called")
	}

	// Check that we received stdout
	hasStdout := false
	for _, update := range progressUpdates {
		if update.Stdout != "" {
			hasStdout = true
			break
		}
	}

	if !hasStdout {
		t.Error("Execute() did not receive stdout in progress updates")
	}

	if result.Stdout == "" {
		t.Error("Execute() result.Stdout was empty")
	}
}

func TestShellCommandType(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to command availability")
	}

	result, err := Shell("echo 'shell test'")
	if err != nil {
		t.Fatalf("Shell() error = %v", err)
	}

	if !strings.Contains(result.Stdout, "shell test") {
		t.Errorf("Shell() output = %v, want to contain 'shell test'", result.Stdout)
	}
}

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		wantExec string
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "simple command",
			command:  "echo hello",
			wantExec: "echo",
			wantArgs: []string{"hello"},
			wantErr:  false,
		},
		{
			name:     "command with quotes",
			command:  `echo "hello world"`,
			wantExec: "echo",
			wantArgs: []string{"hello world"},
			wantErr:  false,
		},
		{
			name:     "empty command",
			command:  "",
			wantExec: "",
			wantArgs: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec, args, err := SplitCommand(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("SplitCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if exec != tt.wantExec {
					t.Errorf("SplitCommand() exec = %v, want %v", exec, tt.wantExec)
				}

				if len(args) != len(tt.wantArgs) {
					t.Errorf("SplitCommand() args len = %v, want %v", len(args), len(tt.wantArgs))
				}

				for i, wantArg := range tt.wantArgs {
					if args[i] != wantArg {
						t.Errorf("SplitCommand() args[%d] = %v, want %v", i, args[i], wantArg)
					}
				}
			}
		})
	}
}

func TestWhich(t *testing.T) {
	// Test with a command that should exist
	path, err := Which("go")
	if err != nil {
		t.Fatalf("Which() error = %v", err)
	}

	if path == "" {
		t.Error("Which() returned empty path")
	}

	// Test with a command that doesn't exist
	_, err = Which("nonexistent-command-xyz")
	if err == nil {
		t.Error("Which() expected error for nonexistent command, got nil")
	}
}

func TestGetDefaultShell(t *testing.T) {
	shell := GetDefaultShell()

	if shell == "" {
		t.Error("GetDefaultShell() returned empty string")
	}

	// On Windows, should return either PowerShell or cmd.exe
	// On Unix, should return /bin/sh or similar
	t.Logf("Default shell: %s", shell)
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "expand tilde",
			input:    "~/test",
			// On Windows, ~ expands to user directory path
			// On Unix, it expands to /home/user/test
			contains: "/test", // This will fail on Windows, but we check for expansion anyway
		},
		{
			name:     "no tilde",
			input:    "/absolute/path",
			// On Windows, this becomes \absolute\path
			// On Unix, it stays /absolute/path
			contains: "absolute",
		},
		{
			name:     "relative path",
			input:    "./relative",
			contains: "relative", // After Clean()
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Logf("ExpandPath() = %v, want to contain %v", result, tt.contains)
				// Don't fail on Windows for path separator differences
				if runtime.GOOS != "windows" || tt.name == "relative path" {
					t.Errorf("ExpandPath() = %v, want to contain %v", result, tt.contains)
				}
			}
		})
	}
}

func TestExecutor_AsyncExecution(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to command availability")
	}

	exec := NewExecutor()

	process, err := exec.ExecuteAsync("sleep", "1")
	if err != nil {
		t.Fatalf("ExecuteAsync() error = %v", err)
	}

	if process.State() != StateRunning {
		t.Errorf("Process.State() = %v, want StateRunning", process.State())
	}

	// Wait for completion
	time.Sleep(2 * time.Second)

	result := process.Result()
	if result.Error != nil && !result.Cancelled && !result.TimedOut {
		t.Errorf("Process.Result().Error = %v, want nil (or cancelled/timedout)", result.Error)
	}

	if process.State() != StateCompleted && process.State() != StateCancelled {
		t.Errorf("Process.State() after completion = %v, want StateCompleted or StateCancelled", process.State())
	}
}

func TestExecutor_CancelProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to different process behavior")
	}

	exec := NewExecutor()

	process, err := exec.ExecuteAsync("sleep", "10")
	if err != nil {
		t.Fatalf("ExecuteAsync() error = %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	err = exec.Cancel(1)
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}

	// Wait for cancellation to propagate
	time.Sleep(500 * time.Millisecond)

	result := process.Result()
	if !result.Cancelled {
		t.Errorf("Process.Result().Cancelled = %v, want true", result.Cancelled)
	}

	if process.State() != StateCancelled {
		t.Errorf("Process.State() = %v, want StateCancelled", process.State())
	}
}

func TestExecutor_CancelMissingProcess(t *testing.T) {
	exec := NewExecutor()

	err := exec.Cancel(999)
	if err == nil {
		t.Error("Cancel() expected error for nonexistent process, got nil")
	}

	if err != ErrProcessNotFound {
		t.Errorf("Cancel() error = %v, want ErrProcessNotFound", err)
	}
}

func TestPipe(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows due to command availability")
	}

	tests := []struct {
		name    string
		cmds    [][]string
		wantOut string
		wantErr bool
	}{
		{
			name:    "single command",
			cmds:    [][]string{{"echo", "hello"}},
			wantOut: "hello",
			wantErr: false,
		},
		{
			name:    "two commands piped",
			cmds:    [][]string{{"echo", "hello world"}, {"grep", "hello"}},
			wantOut: "hello world",
			wantErr: false,
		},
		{
			name:    "empty command list",
			cmds:    [][]string{},
			wantOut: "",
			wantErr: true,
		},
		{
			name:    "empty command in list",
			cmds:    [][]string{{"echo", "test"}, {}},
			wantOut: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Pipe(tt.cmds...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Pipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !strings.Contains(result.Stdout, tt.wantOut) {
				t.Errorf("Pipe() stdout = %v, want to contain %v", result.Stdout, tt.wantOut)
			}
		})
	}
}

// TestError is a simple error for testing Unwrap
type TestError struct {
	msg string
}

func (e TestError) Error() string {
	return e.msg
}

func BenchmarkExecutor_Quick(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Quick("echo", "benchmark")
	}
}

func BenchmarkExecutor_ExecuteWithOutput(b *testing.B) {
	exec := NewExecutor()
	for i := 0; i < b.N; i++ {
		_, _ = exec.Execute("echo", "benchmark")
	}
}
