package status

import "github.com/wwsheng009/taproot/ui/render"

// ServiceStatus represents the status of a service (LSP, MCP, etc.).
type ServiceStatus int

const (
	ServiceStatusOffline ServiceStatus = iota
	ServiceStatusStarting
	ServiceStatusConnecting
	ServiceStatusOnline
	ServiceStatusBusy
	ServiceStatusError
)

// String returns the string representation of the service status.
func (s ServiceStatus) String() string {
	switch s {
	case ServiceStatusOffline:
		return "offline"
	case ServiceStatusStarting:
		return "starting"
	case ServiceStatusConnecting:
		return "connecting"
	case ServiceStatusOnline:
		return "online"
	case ServiceStatusBusy:
		return "busy"
	case ServiceStatusError:
		return "error"
	default:
		return "unknown"
	}
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

// Diagnostic represents a single diagnostic entry.
type Diagnostic struct {
	Severity DiagnosticSeverity
	Message  string
	File     string
	Line     int
	Column   int
}

// DiagnosticSummary represents a summary of diagnostics by severity.
type DiagnosticSummary struct {
	Error   int
	Warning int
	Info    int
	Hint    int
}

// Total returns the total number of diagnostics.
func (d *DiagnosticSummary) Total() int {
	return d.Error + d.Warning + d.Info + d.Hint
}

// HasErrors returns true if there are any errors.
func (d *DiagnosticSummary) HasErrors() bool {
	return d.Error > 0
}

// HasWarnings returns true if there are any warnings.
func (d *DiagnosticSummary) HasWarnings() bool {
	return d.Warning > 0
}

// HasProblems returns true if there are any errors or warnings.
func (d *DiagnosticSummary) HasProblems() bool {
	return d.Error > 0 || d.Warning > 0
}

// Add adds a diagnostic to the summary.
func (d *DiagnosticSummary) Add(severity DiagnosticSeverity) {
	switch severity {
	case DiagnosticSeverityError:
		d.Error++
	case DiagnosticSeverityWarning:
		d.Warning++
	case DiagnosticSeverityInfo:
		d.Info++
	case DiagnosticSeverityHint:
		d.Hint++
	}
}

// Clear clears all diagnostics.
func (d *DiagnosticSummary) Clear() {
	d.Error = 0
	d.Warning = 0
	d.Info = 0
	d.Hint = 0
}
