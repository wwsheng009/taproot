package status

import (
	"time"

	"github.com/wwsheng009/taproot/ui/render"
)

// State represents the connection state of a service.
type State int

const (
	// StateDisabled means the service is disabled or inactive.
	StateDisabled State = iota

	// StateStarting means the service is starting up.
	StateStarting

	// StateReady means the service is connected and ready.
	StateReady

	// StateError means the service encountered an error.
	StateError
)

// String returns the string representation of the state.
func (s State) String() string {
	switch s {
	case StateDisabled:
		return "disabled"
	case StateStarting:
		return "starting"
	case StateReady:
		return "ready"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

// Icon returns the icon for the state.
func (s State) Icon() string {
	switch s {
	case StateDisabled:
		return "○" // Offline/inactive
	case StateStarting:
		return "⟳" // Busy/spinner
	case StateReady:
		return "●" // Online/ready
	case StateError:
		return "×" // Error
	default:
		return "?"
	}
}

// DiagnosticCounts represents diagnostic counts by severity.
type DiagnosticCounts struct {
	Error       int
	Warning     int
	Information int
	Hint        int
}

// Total returns the total number of diagnostics.
func (d DiagnosticCounts) Total() int {
	return d.Error + d.Warning + d.Information + d.Hint
}

// HasAny returns true if there are any diagnostics.
func (d DiagnosticCounts) HasAny() bool {
	return d.Total() > 0
}

// HasErrors returns true if there are any errors.
func (d *DiagnosticCounts) HasErrors() bool {
	return d.Error > 0
}

// HasWarnings returns true if there are any warnings.
func (d *DiagnosticCounts) HasWarnings() bool {
	return d.Warning > 0
}

// HasProblems returns true if there are any errors or warnings.
func (d *DiagnosticCounts) HasProblems() bool {
	return d.Error > 0 || d.Warning > 0
}

// Add adds a diagnostic to the summary.
func (d *DiagnosticCounts) Add(severity DiagnosticSeverity) {
	switch severity {
	case DiagnosticSeverityError:
		d.Error++
	case DiagnosticSeverityWarning:
		d.Warning++
	case DiagnosticSeverityInfo:
		d.Information++
	case DiagnosticSeverityHint:
		d.Hint++
	}
}

// Clear clears all diagnostics.
func (d *DiagnosticCounts) Clear() {
	d.Error = 0
	d.Warning = 0
	d.Information = 0
	d.Hint = 0
}

// DiagnosticSeverity represents the severity level of a diagnostic.
type DiagnosticSeverity int

const (
	DiagnosticSeverityError DiagnosticSeverity = iota
	DiagnosticSeverityWarning
	DiagnosticSeverityInfo
	DiagnosticSeverityHint
)

// String returns the string representation of the diagnostic severity.
func (d DiagnosticSeverity) String() string {
	switch d {
	case DiagnosticSeverityError:
		return "error"
	case DiagnosticSeverityWarning:
		return "warning"
	case DiagnosticSeverityInfo:
		return "info"
	case DiagnosticSeverityHint:
		return "hint"
	default:
		return "unknown"
	}
}

// ToolCounts represents tool and prompt counts for MCP services.
type ToolCounts struct {
	Tools   int
	Prompts int
}

// Total returns the total number of tools and prompts.
func (t ToolCounts) Total() int {
	return t.Tools + t.Prompts
}

// HasAny returns true if there are any tools or prompts.
func (t ToolCounts) HasAny() bool {
	return t.Total() > 0
}

// DiagnosticSummary is an alias for DiagnosticCounts for compatibility.
type DiagnosticSummary = DiagnosticCounts

// LSPService represents an LSP (Language Server Protocol) service.
type LSPService struct {
	// Name is the service name (e.g., "gopls", "rust-analyzer").
	Name string

	// Language is the programming language (e.g., "go", "python", "rust").
	Language string

	// State is the current connection state.
	State State

	// Error is the error message if State is StateError.
	Error string

	// Diagnostics contains diagnostic counts.
	Diagnostics DiagnosticCounts

	// ConnectedAt is the time when the service was connected.
	ConnectedAt time.Time
}

// MCPService represents an MCP (Model Context Protocol) service.
type MCPService struct {
	// Name is the service name (e.g., "filesystem", "git").
	Name string

	// State is the current connection state.
	State State

	// Error is the error message if State is StateError.
	Error string

	// Counts contains tool and prompt counts.
	Counts ToolCounts

	// ConnectedAt is the time when the service was connected.
	ConnectedAt time.Time
}

// Status represents a status display component.
type Status interface {
	render.Model

	// SetWidth sets the display width.
	SetWidth(width int)

	// SetLSPServices sets the LSP services to display.
	SetLSPServices(services []LSPService)

	// SetMCPServices sets the MCP services to display.
	SetMCPServices(services []MCPService)

	// Services returns the current service information.
	Services() ([]LSPService, []MCPService)
}

// Config represents configuration options for status display.
type Config struct {
	// Width is the display width.
	Width int

	// MaxLSPs is the maximum number of LSP services to display.
	MaxLSPs int

	// MaxMCPs is the maximum number of MCP services to display.
	MaxMCPs int

	// ShowIcons enables/disables status icons.
	ShowIcons bool

	// ShowErrorCounts enables/disables diagnostic error counts for LSP.
	ShowErrorCounts bool

	// ShowToolCounts enables/disables tool counts for MCP.
	ShowToolCounts bool
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Width:           40,
		MaxLSPs:         5,
		MaxMCPs:         5,
		ShowIcons:       true,
		ShowErrorCounts: true,
		ShowToolCounts:  true,
	}
}

// Service defines the interface for a service status component.
type Service interface {
	render.Model
	// ID returns the service identifier.
	ID() string
	// Name returns the service name.
	Name() string
	// Status returns the current service status.
	Status() ServiceStatus
	// SetStatus sets the service status.
	SetStatus(status ServiceStatus)
	// IsOnline returns whether the service is online.
	IsOnline() bool
	// ErrorCount returns the error count.
	ErrorCount() int
	// SetErrorCount sets the error count.
	SetErrorCount(count int)
}


