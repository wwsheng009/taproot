package image

import "os"

// env returns the value of an environment variable
func env(key string) string {
	return os.Getenv(key)
}

// termEnv returns the TERM environment variable
func termEnv() string {
	return env("TERM")
}

// hasEnv checks if an environment variable is set and non-empty
func hasEnv(key string) bool {
	return env(key) != ""
}

// envContains checks if an environment variable contains a substring
func envContains(key, substr string) bool {
	val := env(key)
	return val != "" && contains(val, substr)
}

// contains is a simple string contains helper
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (
		s[:len(substr)] == substr ||
		s[len(s)-len(substr):] == substr ||
		containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// PlatformInfo contains information about the terminal/platform
type PlatformInfo struct {
	Term         string
	TermProgram  string
	SupportsKitty bool
	SupportsITerm2 bool
	SupportsSixel bool
	IsWindows    bool
	ColorDepth   int // 1, 4, 8, 24
}

// GetPlatformInfo detects platform and terminal capabilities
func GetPlatformInfo() *PlatformInfo {
	info := &PlatformInfo{
		Term:        termEnv(),
		TermProgram: env("TERM_PROGRAM"),
		ColorDepth:  detectColorDepth(),
	}

	info.SupportsKitty = DetectKitty()
	info.SupportsITerm2 = DetectITerm2()
	info.SupportsSixel = DetectSixel()
	info.IsWindows = detectWindows()

	return info
}

// detectColorDepth detects the terminal color depth
func detectColorDepth() int {
	// Check for truecolor (24-bit)
	if env("COLORTERM") == "truecolor" || env("COLORTERM") == "24bit" {
		return 24
	}

	term := termEnv()

	// Check for 256-color terminals
	if contains(term, "256color") || contains(term, "xterm-256") {
		return 8
	}

	// Check for 16-color terminals
	if contains(term, "color") || contains(term, "ansi") {
		return 4
	}

	// Default to monochrome
	return 1
}

// detectWindows detects if running on Windows
func detectWindows() bool {
	return env("OS") == "Windows_NT" || env("TERM") == "cygwin" || env("TERM") == "msys"
}
