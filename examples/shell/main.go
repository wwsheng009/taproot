// Shell execution utilities demo
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/wwsheng009/taproot/ui/tools/shell"
)

type Model struct {
	executor     *shell.Executor
	running      bool
	currentCmd   string
	lastOutput   string
	lastStderr   string
	lastDuration time.Duration
	progress     []string
}

func NewModel() Model {
	return Model{
		executor: shell.NewExecutor(),
		running:  false,
		progress: make([]string, 0),
	}
}

func (m Model) Run() {
	m.ShowHelp()

	for {
		m.ShowStatus()
		fmt.Print("\nEnter command (or h for help, q to quit): ")
		var input string
		fmt.Scanln(&input)

		switch input {
		case "q", "quit":
			fmt.Println("\nGoodbye!")
			return
		case "h", "help":
			m.ShowHelp()
		case "1":
			m.RunBasicCommand()
		case "2":
			m.RunCommandBuilder()
		case "3":
			m.RunQuickCommand()
		case "4":
			m.RunShellCommand()
		case "5":
			m.RunWithTimeout()
		case "6":
			m.RunWithProgress()
		case "7":
			m.RunAsyncCommand()
		case "8":
			m.RunPipedCommands()
		case "9":
			m.ShowPathUtilities()
		default:
			fmt.Printf("\nUnknown command: %s\n", input)
		}

		fmt.Println("\n" + strings.Repeat("─", 60))
		m.lastOutput = ""
		m.lastStderr = ""
		m.progress = make([]string, 0)
	}
}

func (m Model) ShowHelp() {
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("     Shell Execution Utilities Demo")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println("\nCommands:")
	fmt.Println("  1. Basic command execution")
	fmt.Println("  2. Command builder demo")
	fmt.Println("  3. Quick synchronous command")
	fmt.Println("  4. Shell command execution")
	fmt.Println("  5. Command with timeout")
	fmt.Println("  6. Command with progress callback")
	fmt.Println("  7. Async command execution")
	fmt.Println("  8. Piped commands")
	fmt.Println("  9. Path utilities")
	fmt.Println("  h. Show this help")
	fmt.Println("  q. Quit")
	fmt.Println(strings.Repeat("═", 60))
}

func (m Model) ShowStatus() {
	fmt.Println("\n┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Shell Executor Status                                      │")
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Printf("│ Running Processes: %-40d │\n", m.executor.CountRunning())
	fmt.Printf("│ Platform:          %-40s │\n", runtime.GOOS)
	fmt.Printf("│ Default Shell:     %-40s │\n", shell.GetDefaultShell())
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
}

func (m *Model) RunBasicCommand() {
	fmt.Println("\n📝 Running basic command...")

	cmd := "go"
	args := []string{"version"}
	if runtime.GOOS == "windows" {
		cmd = "powershell.exe"
		args = []string{"-Command", "Get-Host"}
	}

	m.currentCmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	startTime := time.Now()

	result, err := m.executor.Execute(cmd, args...)
	m.lastDuration = time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		m.lastOutput = err.Error()
		return
	}

	m.lastOutput = result.Stdout
	m.ShowResult(result)
}

func (m *Model) RunCommandBuilder() {
	fmt.Println("\n📝 Using command builder...")

	cmd, args := "go", []string{"list", "-m"}
	if runtime.GOOS == "windows" {
		cmd, args = "powershell.exe", []string{"-Command", "Get-Date"}
	}

	builder := shell.NewCommandBuilder().
		Command(cmd, args...).
		SetWorkingDir("").
		SetTimeout(30 * time.Second).
		SetOutputDirection(shell.OutputStdout).
		SetQuiet(false)

	m.currentCmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	startTime := time.Now()

	 _, _, opts, err := builder.Build()
	if err != nil {
		fmt.Printf("❌ Build error: %v\n", err)
		return
	}

	result, err := m.executor.ExecuteWithOptions(cmd, args, opts)
	m.lastDuration = time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		m.lastOutput = err.Error()
		return
	}

	m.lastOutput = result.Stdout
	m.ShowResult(result)
}

func (m *Model) RunQuickCommand() {
	fmt.Println("\n📝 Running quick command...")

	cmd, args := "go", []string{"env", "GOVERSION"}
	if runtime.GOOS == "windows" {
		cmd, args = "powershell.exe", []string{"-Command", "$env:OS"}
	}

	m.currentCmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))

	startTime := time.Now()
	output, err := shell.Quick(cmd, args...)
	m.lastDuration = time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		m.lastOutput = err.Error()
		return
	}

	m.lastOutput = strings.TrimSpace(output)
	fmt.Printf("✅ Output:\n%s\n", m.lastOutput)
}

func (m *Model) RunShellCommand() {
	fmt.Println("\n📝 Running shell command...")
	m.ShowWarning("Shell command execution")

	if runtime.GOOS == "windows" {
		result, err := shell.Shell("powershell -Command 'Write-Host \"Hello from PowerShell\"'")
		m.currentCmd = "powershell -Command ..."
		m.lastDuration = result.Duration

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			m.lastOutput = err.Error()
			return
		}

		m.lastOutput = result.Stdout
		m.ShowResult(result)
	} else {
		result, err := shell.Shell("echo 'Hello from Shell'")
		m.currentCmd = "echo 'Hello from Shell'"
		m.lastDuration = result.Duration

		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			m.lastOutput = err.Error()
			return
		}

		m.lastOutput = result.Stdout
		m.ShowResult(result)
	}
}

func (m *Model) RunWithTimeout() {
	fmt.Println("\n📝 Running command with timeout...")

	timeout := 100 * time.Millisecond
	cmd, args := "sleep", []string{"5"}
	if runtime.GOOS == "windows" {
		timeout = 2 * time.Second
		cmd, args = "ping", []string{"-n", "6", "127.0.0.1"}
	}

	opts := shell.CommandOptions{
		Type:      shell.CommandDirect,
		Direction: shell.OutputNone,
		Timeout:   timeout,
	}

	m.currentCmd = fmt.Sprintf("%s %s (timeout: %s)", cmd, strings.Join(args, " "), timeout)
	startTime := time.Now()

	result, err := m.executor.ExecuteWithOptions(cmd, args, opts)
	m.lastDuration = time.Since(startTime)

	if result.TimedOut {
		fmt.Printf("⏱️  Command timed out after %s as expected\n", timeout)
		fmt.Printf("   Duration: %s\n", result.Duration)
		return
	}

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Command completed in %s\n", result.Duration)
}

func (m *Model) RunWithProgress() {
	fmt.Println("\n📝 Running command with progress callback...")
	m.ShowWarning("Progress tracking")

	cmd, args := "echo", []string{"test"}
	if runtime.GOOS == "windows" {
		cmd, args = "cmd.exe", []string{"/c", "echo test"}
	}

	m.progress = make([]string, 0)
	opts := shell.CommandOptions{
		Type:      shell.CommandDirect,
		Direction: shell.OutputStdout,
		Progress: func(update shell.ProgressUpdate) {
			m.progress = append(m.progress, fmt.Sprintf("[%s] State: %s",
				update.Timestamp.Format("15:04:05.000"),
				update.State.String()))
			if update.Stdout != "" {
				fmt.Printf("   STDOUT: %s", update.Stdout)
			}
			if update.Stderr != "" {
				fmt.Printf("   STDERR: %s", update.Stderr)
			}
		},
	}

	m.currentCmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	startTime := time.Now()

	result, err := m.executor.ExecuteWithOptions(cmd, args, opts)
	m.lastDuration = time.Since(startTime)

	if err != nil {
		fmt.Printf("\n❌ Error: %v\n", err)
		m.lastOutput = err.Error()
		return
	}

	fmt.Printf("\n✅ Command completed\n")
	fmt.Printf("   Total progress updates: %d\n", len(m.progress))
	m.lastOutput = result.Stdout
}

func (m *Model) RunAsyncCommand() {
	fmt.Println("\n📝 Running async command...")
	m.ShowWarning("Async execution")

	if runtime.GOOS == "windows" {
		fmt.Println("⚠️  Skipping on Windows (different process behavior)")
		return
	}

	cmd, args := "sleep", []string{"2"}
	m.currentCmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))

	startTime := time.Now()
	process, err := m.executor.ExecuteAsync(cmd, args...)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("🔄 Process started (PID: 1)\n")
	fmt.Printf("   State: %s\n", process.State())

	// Wait a bit and check status
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("   After 500ms: State = %s, Elapsed = %v\n", process.State(), process.Elapsed())

	// Wait for completion
	time.Sleep(2 * time.Second)
	m.lastDuration = time.Since(startTime)

	result := process.Result()
	fmt.Printf("\n✅ Process completed\n")
	fmt.Printf("   Final state: %s\n", process.State())
	fmt.Printf("   Exit code: %d\n", result.ExitCode)
	fmt.Printf("   Duration: %s\n", result.Duration)

	if result.Error != nil {
		fmt.Printf("   Error: %v\n", result.Error)
	}
}

func (m *Model) RunPipedCommands() {
	fmt.Println("\n📝 Running piped commands...")

	if runtime.GOOS == "windows" {
		fmt.Println("⚠️  Skipping on Windows (command availability)")
		return
	}

	cmds := [][]string{
		{"echo", "hello world"},
		{"grep", "hello"},
	}

	m.currentCmd = ""
	for _, cmd := range cmds {
		m.currentCmd += fmt.Sprintf("%s ", strings.Join(cmd, " "))
		m.currentCmd += "| "
	}
	m.currentCmd = strings.TrimSuffix(m.currentCmd, "| ")

	startTime := time.Now()
	result, err := shell.Pipe(cmds...)
	m.lastDuration = time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		m.lastOutput = err.Error()
		return
	}

	m.lastOutput = result.Stdout
	m.ShowResult(result)
}

func (m *Model) ShowPathUtilities() {
	fmt.Println("\n📝 Path utilities...")

	// Find executable
	if path, err := shell.Which("go"); err == nil {
		fmt.Printf("� Found 'go' at: %s\n", path)
	} else {
		fmt.Printf("❌ Could not find 'go': %v\n", err)
	}

	// Get default shell
	shellPath := shell.GetDefaultShell()
	fmt.Printf("� Default shell: %s\n", shellPath)

	// Split command
	cmd, args, err := shell.SplitCommand("echo hello world")
	if err == nil {
		fmt.Printf("� Split command: cmd='%s' args=%v\n", cmd, args)
	}

	// Expand path
	if runtime.GOOS != "windows" {
		fmt.Printf("� Expanded path: %s -> %s\n", "~/test", shell.ExpandPath("~/test"))
	}

	// Join path
	joined := shell.JoinPath("home", "user", "documents")
	fmt.Printf("� Joined path: %s\n", joined)
}

func (m Model) ShowWarning(msg string) {
	fmt.Printf("\n⚠️  %s\n", msg)
}

func (m Model) ShowResult(result shell.CommandResult) {
	fmt.Printf("\n✅ Command completed successfully\n")
	fmt.Printf("   Exit code: %d\n", result.ExitCode)
	fmt.Printf("   Duration: %v\n", result.Duration)
	fmt.Printf("   Stdout bytes: %d\n", len(result.Stdout))
	fmt.Printf("   Stderr bytes: %d\n", len(result.Stderr))

	if result.Stdout != "" {
		fmt.Printf("\n📄 STDOUT:\n")
		maxLines := 10
		lines := strings.Split(result.Stdout, "\n")
		for i, line := range lines {
			if i >= maxLines {
				fmt.Printf("   ... (%d more lines)\n", len(lines)-maxLines)
				break
			}
			fmt.Printf("   %s\n", line)
		}
	}

	if result.Stderr != "" {
		fmt.Printf("\n⚠️  STDERR:\n%s\n", result.Stderr)
	}
}

func main() {
	model := NewModel()
	model.Run()
}
