package status

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// DiagnosticStatusCmp is a component that displays diagnostic information
// (errors, warnings, info, hints) typically from LSP/MCP services.
type DiagnosticStatusCmp struct {
	summary      DiagnosticSummary
	source       string
	sourceCounts map[string]DiagnosticSummary
	focused      bool
	compact      bool
	showHints    bool
	maxWidth     int
	initialized  bool
}

// NewDiagnosticStatus creates a new diagnostic status component.
func NewDiagnosticStatus(source string) *DiagnosticStatusCmp {
	return &DiagnosticStatusCmp{
		summary:      DiagnosticSummary{},
		source:       source,
		sourceCounts: make(map[string]DiagnosticSummary),
		focused:      false,
		compact:      false,
		showHints:    true,
		maxWidth:     0,
		initialized:  false,
	}
}

// Init initializes the component.
// Implements render.Model interface.
func (d *DiagnosticStatusCmp) Init() render.Cmd {
	d.initialized = true
	return nil
}

// Update handles incoming messages and returns updated state and commands.
// Implements render.Model interface.
func (d *DiagnosticStatusCmp) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		d.focused = true
	case *render.BlurMsg:
		d.focused = false
	}
	return d, render.None()
}

// View returns the string representation for rendering.
// Implements render.Model interface.
func (d *DiagnosticStatusCmp) View() string {
	sty := styles.DefaultStyles()

	if d.compact {
		return d.renderCompact(sty)
	}
	return d.renderExpanded(sty)
}

// renderCompact returns the compact view showing only error and warning counts.
func (d *DiagnosticStatusCmp) renderCompact(sty styles.Styles) string {
	var result string

	// Render error count
	if d.summary.Error > 0 {
		result += sty.LSP.ErrorDiagnostic.Foreground(sty.Error).
			Render(string(styles.ErrorIcon) + " ")
		result += sty.LSP.ErrorDiagnostic.Render(fmt.Sprintf("%d", d.summary.Error))
	}

	// Render warning count
	if d.summary.Warning > 0 {
		if d.summary.Error > 0 {
			result += " "
		}
		result += sty.LSP.WarningDiagnostic.Foreground(sty.Warning).
			Render(string(styles.WarningIcon) + " ")
		result += sty.LSP.WarningDiagnostic.Render(fmt.Sprintf("%d", d.summary.Warning))
	}

	// Render info count (only if no errors/warnings)
	if d.summary.Error == 0 && d.summary.Warning == 0 && d.summary.Information > 0 {
		result += sty.LSP.InfoDiagnostic.Foreground(sty.Info).
			Render(string(styles.InfoIcon) + " ")
		result += sty.LSP.InfoDiagnostic.Render(fmt.Sprintf("%d", d.summary.Information))
	}

	if d.maxWidth > 0 && lipgloss.Width(result) > d.maxWidth {
		result = lipgloss.NewStyle().MaxWidth(d.maxWidth).Render(result)
	}

	return result
}

// renderExpanded returns the expanded view showing all diagnostic counts.
func (d *DiagnosticStatusCmp) renderExpanded(sty styles.Styles) string {
	var result string

	// Add source label if focused
	if d.focused && d.source != "" {
		sourceLabel := sty.Base.Foreground(sty.FgMuted).Render(d.source + ":")
		result += sourceLabel + " "
	}

	// Render each severity type
	hasErrors := d.summary.Error > 0
	hasWarnings := d.summary.Warning > 0
	hasInfo := d.summary.Information > 0
	hasHints := d.showHints && d.summary.Hint > 0

	if hasErrors {
		result += sty.LSP.ErrorDiagnostic.Foreground(sty.Error).
			Render(string(styles.ErrorIcon)) + " "
		result += sty.LSP.ErrorDiagnostic.Render(fmt.Sprintf("%d", d.summary.Error))
	}

	if hasWarnings {
		if hasErrors {
			result += " "
		}
		result += sty.LSP.WarningDiagnostic.Foreground(sty.Warning).
			Render(string(styles.WarningIcon)) + " "
		result += sty.LSP.WarningDiagnostic.Render(fmt.Sprintf("%d", d.summary.Warning))
	}

	if hasInfo {
		if hasErrors || hasWarnings {
			result += " "
		}
		result += sty.LSP.InfoDiagnostic.Foreground(sty.Info).
			Render(string(styles.InfoIcon)) + " "
		result += sty.LSP.InfoDiagnostic.Render(fmt.Sprintf("%d", d.summary.Information))
	}

	if hasHints {
		if hasErrors || hasWarnings || hasInfo {
			result += " "
		}
		result += sty.LSP.HintDiagnostic.Foreground(sty.FgSubtle).
			Render(string(styles.HintIcon)) + " "
		result += sty.LSP.HintDiagnostic.Render(fmt.Sprintf("%d", d.summary.Hint))
	}

	// Show "No diagnostics" if empty
	if d.summary.Total() == 0 {
		result = sty.Base.Foreground(sty.FgMuted).Render("No diagnostics")
	}

	if d.maxWidth > 0 && lipgloss.Width(result) > d.maxWidth {
		result = lipgloss.NewStyle().MaxWidth(d.maxWidth).Render(result)
	}

	return result
}

// Summary returns the current diagnostic summary.
func (d *DiagnosticStatusCmp) Summary() DiagnosticSummary {
	return d.summary
}

// SetSummary sets the diagnostic summary.
func (d *DiagnosticStatusCmp) SetSummary(summary DiagnosticSummary) {
	d.summary = summary
}

// AddDiagnostic adds a diagnostic to the summary.
func (d *DiagnosticStatusCmp) AddDiagnostic(severity DiagnosticSeverity) {
	d.summary.Add(severity)
}

// Clear clears all diagnostics.
func (d *DiagnosticStatusCmp) Clear() {
	d.summary.Clear()
	d.sourceCounts = make(map[string]DiagnosticSummary)
}

// Source returns the source name.
func (d *DiagnosticStatusCmp) Source() string {
	return d.source
}

// SetSource sets the source name.
func (d *DiagnosticStatusCmp) SetSource(source string) {
	d.source = source
}

// GetSourceSummary returns diagnostic summary for a specific source.
func (d *DiagnosticStatusCmp) GetSourceSummary(source string) DiagnosticSummary {
	return d.sourceCounts[source]
}

// AddSourceDiagnostic adds a diagnostic from a specific source.
func (d *DiagnosticStatusCmp) AddSourceDiagnostic(source string, severity DiagnosticSeverity) {
	summary := d.sourceCounts[source]
	summary.Add(severity)
	d.sourceCounts[source] = summary
}

// HasProblems returns true if there are any errors or warnings.
func (d *DiagnosticStatusCmp) HasProblems() bool {
	return d.summary.HasProblems()
}

// HasErrors returns true if there are any errors.
func (d *DiagnosticStatusCmp) HasErrors() bool {
	return d.summary.HasErrors()
}

// Total returns the total number of diagnostics.
func (d *DiagnosticStatusCmp) Total() int {
	return d.summary.Total()
}

// Focus focuses the component.
func (d *DiagnosticStatusCmp) Focus() {
	d.focused = true
}

// Blur unfocuses the component.
func (d *DiagnosticStatusCmp) Blur() {
	d.focused = false
}

// Focused returns whether the component is focused.
func (d *DiagnosticStatusCmp) Focused() bool {
	return d.focused
}

// SetCompact sets whether to use compact mode.
func (d *DiagnosticStatusCmp) SetCompact(compact bool) {
	d.compact = compact
}

// Compact returns whether compact mode is enabled.
func (d *DiagnosticStatusCmp) Compact() bool {
	return d.compact
}

// SetShowHints sets whether to show hints in expanded mode.
func (d *DiagnosticStatusCmp) SetShowHints(show bool) {
	d.showHints = show
}

// ShowHints returns whether hints are shown.
func (d *DiagnosticStatusCmp) ShowHints() bool {
	return d.showHints
}

// SetMaxWidth sets the maximum width for rendering.
func (d *DiagnosticStatusCmp) SetMaxWidth(width int) {
	d.maxWidth = width
}

// MaxWidth returns the maximum width.
func (d *DiagnosticStatusCmp) MaxWidth() int {
	return d.maxWidth
}

// ErrorCount returns the number of errors.
func (d *DiagnosticStatusCmp) ErrorCount() int {
	return d.summary.Error
}

// WarningCount returns the number of warnings.
func (d *DiagnosticStatusCmp) WarningCount() int {
	return d.summary.Warning
}

// InfoCount returns the number of info diagnostics.
func (d *DiagnosticStatusCmp) InfoCount() int {
	return d.summary.Information
}

// HintCount returns the number of hints.
func (d *DiagnosticStatusCmp) HintCount() int {
	return d.summary.Hint
}
