package header

import (
	"fmt"
)

// Header represents the application header component.
type Header interface {
	Sizeable
	SetBrand(brand, title string)
	SetSessionTitle(title string)
	SetWorkingDirectory(cwd string)
	SetTokenUsage(used, max int, cost float64)
	SetErrorCount(count int)
	SetDetailsOpen(open bool)
	ShowingDetails() bool
}

// Sizeable defines the sizing interface.
type Sizeable interface {
	Size() (width, height int)
	SetSize(width, height int)
}

// Config holds the header configuration.
type Config struct {
	Width   int
	Height  int
	Brand   string
	Title   string
	Compact bool
}

// SessionInfo holds session-related data.
type SessionInfo struct {
	Title       string
	WorkingDir  string
	TokenUsed   int
	TokenMax    int
	Cost        float64
}

// DiagnosticInfo holds diagnostic-related data.
type DiagnosticInfo struct {
	ErrorCount int
	WarnCount  int
}

// FormatTokenUsage formats token usage percentage and cost.
// Returns formatted string like "50% (1.2K) $0.75"
func FormatTokenUsage(used, max int, cost float64) string {
	if max <= 0 {
		return fmt.Sprintf("%d $%.2f", used, cost)
	}
	percentage := int(float64(used) / float64(max) * 100)

	// Format used tokens with K suffix for large numbers
	usedStr := fmt.Sprintf("%d", used)
	if used >= 1000 {
		usedStr = fmt.Sprintf("%.1fK", float64(used)/1000)
	}

	return fmt.Sprintf("%d%% (%s) $%.2f", percentage, usedStr, cost)
}

// FormatErrorMessage formats the error count for display.
func FormatErrorMessage(count int) string {
	if count == 0 {
		return ""
	}
	return fmt.Sprintf("%d", count)
}
