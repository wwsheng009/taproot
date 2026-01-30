package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/render"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// DiagnosticMessage represents a diagnostic message containing error/warning/info diagnostics.
//
// This component supports:
// - Multiple diagnostics with different severity levels
// - File path and location information
// - Diagnostic messages and error codes
// - Collapsible diagnostic details
// - Focus state with different styling
// - Caching for performance
type DiagnosticMessage struct {
	id         string
	title      string
	diagnostics []Diagnostic
	timestamp  time.Time
	focused    bool
	expanded   bool
	config     *MessageConfig
	maxWidth   int

	// Caching fields
	cachedRender string
	cachedWidth  int
	cachedHeight int
	cacheValid   bool
}

// NewDiagnosticMessage creates a new DiagnosticMessage component.
func NewDiagnosticMessage(id, title string) *DiagnosticMessage {
	return &DiagnosticMessage{
		id:           id,
		title:        title,
		diagnostics:  []Diagnostic{},
		timestamp:    time.Now(),
		focused:      false,
		expanded:     false,
		config:       &MessageConfig{},
	}
}

// NewDiagnosticMessageFromDiagnostics creates a new DiagnosticMessage from existing diagnostics.
func NewDiagnosticMessageFromDiagnostics(id, title string, diagnostics []Diagnostic) *DiagnosticMessage {
	return &DiagnosticMessage{
		id:           id,
		title:        title,
		diagnostics:  diagnostics,
		timestamp:    time.Now(),
		focused:      false,
		expanded:     false,
		config:       &MessageConfig{},
	}
}

// Init initializes the component. Implements render.Model.
func (m *DiagnosticMessage) Init() error {
	return nil
}

// Update handles incoming messages. Implements render.Model.
func (m *DiagnosticMessage) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		m.focused = true
	case *render.BlurMsg:
		m.focused = false
	}

	// Invalidate cache on any update
	m.cacheValid = false

	return m, render.None()
}

// View renders the diagnostic message. Implements render.Model.
func (m *DiagnosticMessage) View() string {
	width := m.config.MaxWidth
	// Apply max width limit
	if m.maxWidth > 0 && width > m.maxWidth {
		width = m.maxWidth
	}

	return m.render(width)
}

// ID returns the message ID. Implements Identifiable interface.
func (m *DiagnosticMessage) ID() string {
	return m.id
}

// Role returns RoleTool (diagnostics are system/tool messages). Implements Message interface.
func (m *DiagnosticMessage) Role() MessageRole {
	return RoleTool
}

// Content returns a summary of diagnostics. Implements Message interface.
func (m *DiagnosticMessage) Content() string {
	return fmt.Sprintf("%s (%d diagnostics)", m.title, len(m.diagnostics))
}

// SetContent sets the title. Implements Message interface.
func (m *DiagnosticMessage) SetContent(content string) {
	m.title = content
	m.cacheValid = false
}

// Timestamp returns the message timestamp. Implements Message interface.
func (m *DiagnosticMessage) Timestamp() time.Time {
	return m.timestamp
}

// SetTimestamp sets the message timestamp.
func (m *DiagnosticMessage) SetTimestamp(ts time.Time) {
	m.timestamp = ts
	m.cacheValid = false
}

// Title returns the diagnostic message title.
func (m *DiagnosticMessage) Title() string {
	return m.title
}

// SetTitle sets the title.
func (m *DiagnosticMessage) SetTitle(title string) {
	m.title = title
	m.cacheValid = false
}

// Diagnostics returns the diagnostics.
func (m *DiagnosticMessage) Diagnostics() []Diagnostic {
	return m.diagnostics
}

// SetDiagnostics sets the diagnostics.
func (m *DiagnosticMessage) SetDiagnostics(diagnostics []Diagnostic) {
	m.diagnostics = diagnostics
	m.cacheValid = false
}

// AddDiagnostic adds a diagnostic.
func (m *DiagnosticMessage) AddDiagnostic(diagnostic Diagnostic) {
	m.diagnostics = append(m.diagnostics, diagnostic)
	m.cacheValid = false
}

// Expanded returns true if the message is expanded.
func (m *DiagnosticMessage) Expanded() bool {
	return m.expanded
}

// SetExpanded sets the expansion state.
func (m *DiagnosticMessage) SetExpanded(expanded bool) {
	m.expanded = expanded
	m.cacheValid = false
}

// ToggleExpanded toggles the expansion state.
func (m *DiagnosticMessage) ToggleExpanded() {
	m.expanded = !m.expanded
	m.cacheValid = false
}

// Focus focuses the component. Implements Focusable interface.
func (m *DiagnosticMessage) Focus() {
	m.focused = true
	m.cacheValid = false
}

// Blur blurs the component. Implements Focusable interface.
func (m *DiagnosticMessage) Blur() {
	m.focused = false
	m.cacheValid = false
}

// Focused returns true if the component is focused. Implements Focusable interface.
func (m *DiagnosticMessage) Focused() bool {
	return m.focused
}

// SetMaxWidth sets the maximum width for rendering.
func (m *DiagnosticMessage) SetMaxWidth(width int) {
	m.maxWidth = width
	m.cacheValid = false
}

// SetConfig sets the message configuration.
func (m *DiagnosticMessage) SetConfig(config *MessageConfig) {
	m.config = config
	m.cacheValid = false
}

// render renders the diagnostic message with the given width.
func (m *DiagnosticMessage) render(width int) string {
	// Check cache first
	if m.cacheValid && m.cachedWidth == width && m.cachedRender != "" {
		return m.cachedRender
	}

	sty := styles.DefaultStyles()

	var builder strings.Builder

	// Render header
	header := m.renderHeader(&sty)
	builder.WriteString(header)

	// Render diagnostics if expanded
	if m.expanded && len(m.diagnostics) > 0 {
		for i, diag := range m.diagnostics {
			if i > 0 {
				builder.WriteString("\n\n")
			}
			builder.WriteString(m.renderDiagnostic(&sty, diag, width))
		}
	}

	// Add expand hint if collapsed and there are diagnostics
	if !m.expanded && len(m.diagnostics) > 0 {
		builder.WriteString("\n")
		expandHint := sty.Base.Foreground(sty.FgMuted).Italic(true).Render(
			fmt.Sprintf("[%d diagnostic(s) - press enter to expand]", len(m.diagnostics)),
		)
		builder.WriteString(expandHint)
	}

	result := builder.String()

	// Apply focus styling
	if m.focused {
		result = lipgloss.NewStyle().Foreground(sty.Warning).Render(result)
	} else {
		result = lipgloss.NewStyle().Foreground(sty.FgMuted).Render(result)
	}

	// Cache the result
	m.cachedRender = result
	m.cachedWidth = width
	m.cachedHeight = lipgloss.Height(result)
	m.cacheValid = true

	return result
}

// renderHeader renders the diagnostic message header.
func (m *DiagnosticMessage) renderHeader(sty *styles.Styles) string {
	severityIcon := "ℹ️"
	severityCount := 0

	// Determine the overall severity (worst one)
	for _, diag := range m.diagnostics {
		severityCount++
		if diag.Severity == SeverityError {
			severityIcon = "❌"
			break
		} else if diag.Severity == SeverityWarning {
			severityIcon = "⚠️"
		}
	}

	// If we have diagnostics, count them
	headerText := fmt.Sprintf("%s %s", severityIcon, m.title)
	if severityCount > 0 {
		headerText = fmt.Sprintf("%s %s (%d)", severityIcon, m.title, severityCount)
	}

	headerStyle := sty.Base.Bold(true)
	if severityCount > 0 {
		// Use color based on highest severity
		hasError := false
		hasWarning := false
		for _, diag := range m.diagnostics {
			if diag.Severity == SeverityError {
				hasError = true
				break
			} else if diag.Severity == SeverityWarning {
				hasWarning = true
			}
		}

		if hasError {
			headerStyle = headerStyle.Foreground(sty.Error)
		} else if hasWarning {
			headerStyle = headerStyle.Foreground(sty.Warning)
		} else {
			headerStyle = headerStyle.Foreground(sty.Info)
		}
	}

	return headerStyle.Render(headerText)
}

// renderDiagnostic renders a single diagnostic.
func (m *DiagnosticMessage) renderDiagnostic(sty *styles.Styles, diag Diagnostic, width int) string {
	var parts []string

	// Diagnostic header with severity, file, and location
	severityIcon := getSeverityIcon(diag.Severity)
	location := ""
	if diag.File != "" {
		location = diag.File
		if diag.Line > 0 {
			location += fmt.Sprintf(":%d", diag.Line)
			if diag.Column > 0 {
				location += fmt.Sprintf(":%d", diag.Column)
			}
		}
	}

	headerText := fmt.Sprintf("  %s", severityIcon)
	if location != "" {
		headerText += fmt.Sprintf(" %s", location)
	} else {
		headerText += fmt.Sprintf(" %s", diag.Severity.String())
	}

	headerStyle := sty.Base.Bold(true)
	switch diag.Severity {
	case SeverityError:
		headerStyle = headerStyle.Foreground(sty.Error)
	case SeverityWarning:
		headerStyle = headerStyle.Foreground(sty.Warning)
	case SeverityInfo:
		headerStyle = headerStyle.Foreground(sty.Info)
	case SeverityHint:
		headerStyle = headerStyle.Foreground(sty.FgMuted)
	}

	parts = append(parts, headerStyle.Render(headerText))

	// Diagnostic message
	if diag.Message != "" {
		messageStyle := sty.Base.Foreground(sty.FgBase)
		parts = append(parts, messageStyle.Render("  "+diag.Message))
	}

	// Error code if present
	if diag.Code != "" {
		codeStyle := sty.Base.Foreground(sty.FgMuted).Italic(true)
		parts = append(parts, codeStyle.Render(fmt.Sprintf("  [%s]", diag.Code)))
	}

	// Note hint for hints
	if diag.Severity == SeverityHint {
		hintStyle := sty.Base.Foreground(sty.FgMuted).Italic(true)
		parts = append(parts, hintStyle.Render("  💡 Hint"))
	}

	return strings.Join(parts, "\n")
}

// getSeverityIcon returns the icon for a diagnostic severity.
func getSeverityIcon(severity DiagnosticSeverity) string {
	switch severity {
	case SeverityError:
		return "❌"
	case SeverityWarning:
		return "⚠️"
	case SeverityInfo:
		return "ℹ️"
	case SeverityHint:
		return "💡"
	default:
		return "❓"
	}
}

// handleFocusGained is a helper for handling focus gain.
func (m *DiagnosticMessage) handleFocusGained() {
	m.focused = true
	m.cacheValid = false
}

// HandleClick handles click events for expanding/collapsing.
// Implements ClickHandler interface.
func (m *DiagnosticMessage) HandleClick(line, col int) render.Cmd {
	if len(m.diagnostics) > 0 {
		m.ToggleExpanded()
		return render.None()
	}
	return render.None()
}

// HandleKeyEvent handles keyboard events.
// Implements KeyEventHandler interface.
func (m *DiagnosticMessage) HandleKeyEvent(msg any) (bool, render.Cmd) {
	keyMsg, ok := msg.(*render.KeyMsg)
	if !ok {
		return false, nil
	}

	switch keyMsg.String() {
	case " ", "enter":
		// Toggle expansion if there are diagnostics
		if len(m.diagnostics) > 0 {
			m.ToggleExpanded()
			return true, render.None()
		}
	}

	return false, nil
}

// DiagnosticCount returns the total number of diagnostics.
func (m *DiagnosticMessage) DiagnosticCount() int {
	return len(m.diagnostics)
}

// ErrorCount returns the number of error diagnostics.
func (m *DiagnosticMessage) ErrorCount() int {
	count := 0
	for _, diag := range m.diagnostics {
		if diag.Severity == SeverityError {
			count++
		}
	}
	return count
}

// WarningCount returns the number of warning diagnostics.
func (m *DiagnosticMessage) WarningCount() int {
	count := 0
	for _, diag := range m.diagnostics {
		if diag.Severity == SeverityWarning {
			count++
		}
	}
	return count
}

// InfoCount returns the number of info diagnostics.
func (m *DiagnosticMessage) InfoCount() int {
	count := 0
	for _, diag := range m.diagnostics {
		if diag.Severity == SeverityInfo {
			count++
		}
	}
	return count
}

// HintCount returns the number of hint diagnostics.
func (m *DiagnosticMessage) HintCount() int {
	count := 0
	for _, diag := range m.diagnostics {
		if diag.Severity == SeverityHint {
			count++
		}
	}
	return count
}
