package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/render"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// ToolMessage represents a tool call message that displays tool execution status.
//
// This component supports:
// - Multiple tool calls in a single message
// - Tool status tracking (Pending, Running, Completed, Error, Canceled)
// - Collapsible tool call details
// - Error display for failed tool calls
// - Focus state with different styling
// - Caching for performance
type ToolMessage struct {
	id         string
	name       string
	calls      []ToolCall
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

// NewToolMessage creates a new ToolMessage component with a single tool call.
func NewToolMessage(id, toolName string) *ToolMessage {
	return &ToolMessage{
		id:        id,
		name:      toolName,
		calls:     []ToolCall{},
		timestamp: time.Now(),
		focused:   false,
		expanded:  false,
		config:    &MessageConfig{},
	}
}

// NewToolMessageFromCalls creates a new ToolMessage from existing tool calls.
func NewToolMessageFromCalls(id, toolName string, calls []ToolCall) *ToolMessage {
	return &ToolMessage{
		id:        id,
		name:      toolName,
		calls:     calls,
		timestamp: time.Now(),
		focused:   false,
		expanded:  false,
		config:    &MessageConfig{},
	}
}

// Init initializes the component. Implements render.Model.
func (m *ToolMessage) Init() error {
	return nil
}

// Update handles incoming messages. Implements render.Model.
func (m *ToolMessage) Update(msg any) (render.Model, render.Cmd) {
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

// View renders the tool message. Implements render.Model.
func (m *ToolMessage) View() string {
	width := m.config.MaxWidth
	// Apply max width limit
	if m.maxWidth > 0 && width > m.maxWidth {
		width = m.maxWidth
	}

	return m.render(width)
}

// ID returns the message ID. Implements Identifiable interface.
func (m *ToolMessage) ID() string {
	return m.id
}

// Role returns RoleTool. Implements Message interface.
func (m *ToolMessage) Role() MessageRole {
	return RoleTool
}

// Content returns a summary of tool calls. Implements Message interface.
func (m *ToolMessage) Content() string {
	return fmt.Sprintf("%s (%d calls)", m.name, len(m.calls))
}

// SetContent sets the tool name. Implements Message interface.
func (m *ToolMessage) SetContent(content string) {
	m.name = content
	m.cacheValid = false
}

// Timestamp returns the message timestamp. Implements Message interface.
func (m *ToolMessage) Timestamp() time.Time {
	return m.timestamp
}

// SetTimestamp sets the message timestamp.
func (m *ToolMessage) SetTimestamp(ts time.Time) {
	m.timestamp = ts
	m.cacheValid = false
}

// Name returns the tool name.
func (m *ToolMessage) Name() string {
	return m.name
}

// SetName sets the tool name.
func (m *ToolMessage) SetName(name string) {
	m.name = name
	m.cacheValid = false
}

// Calls returns the tool calls.
func (m *ToolMessage) Calls() []ToolCall {
	return m.calls
}

// SetCalls sets the tool calls.
func (m *ToolMessage) SetCalls(calls []ToolCall) {
	m.calls = calls
	m.cacheValid = false
}

// AddCall adds a tool call.
func (m *ToolMessage) AddCall(call ToolCall) {
	m.calls = append(m.calls, call)
	m.cacheValid = false
}

// UpdateCallStatus updates the status of a tool call by ID.
func (m *ToolMessage) UpdateCallStatus(callID string, status ToolStatus) {
	for i := range m.calls {
		if m.calls[i].ID == callID {
			m.calls[i].Status = status
			m.calls[i].Timestamp = time.Now()
			m.cacheValid = false
			return
		}
	}
}

// UpdateCallResult updates the result of a tool call by ID.
func (m *ToolMessage) UpdateCallResult(callID string, result string, err string) {
	for i := range m.calls {
		if m.calls[i].ID == callID {
			m.calls[i].Result = result
			m.calls[i].Error = err
			m.calls[i].Timestamp = time.Now()
			if err != "" {
				m.calls[i].Status = ToolStatusError
			} else {
				m.calls[i].Status = ToolStatusCompleted
			}
			m.cacheValid = false
			return
		}
	}
}

// Expanded returns true if the message is expanded.
func (m *ToolMessage) Expanded() bool {
	return m.expanded
}

// SetExpanded sets the expansion state.
func (m *ToolMessage) SetExpanded(expanded bool) {
	m.expanded = expanded
	m.cacheValid = false
}

// ToggleExpanded toggles the expansion state.
func (m *ToolMessage) ToggleExpanded() {
	m.expanded = !m.expanded
	m.cacheValid = false
}

// Focus focuses the component. Implements Focusable interface.
func (m *ToolMessage) Focus() {
	m.focused = true
	m.cacheValid = false
}

// Blur blurs the component. Implements Focusable interface.
func (m *ToolMessage) Blur() {
	m.focused = false
	m.cacheValid = false
}

// Focused returns true if the component is focused. Implements Focusable interface.
func (m *ToolMessage) Focused() bool {
	return m.focused
}

// SetMaxWidth sets the maximum width for rendering.
func (m *ToolMessage) SetMaxWidth(width int) {
	m.maxWidth = width
	m.cacheValid = false
}

// SetConfig sets the message configuration.
func (m *ToolMessage) SetConfig(config *MessageConfig) {
	m.config = config
	m.cacheValid = false
}

// render renders the tool message with the given width.
func (m *ToolMessage) render(width int) string {
	// Check cache first
	if m.cacheValid && m.cachedWidth == width && m.cachedRender != "" {
		return m.cachedRender
	}

	sty := styles.DefaultStyles()

	var builder strings.Builder

	// Render tool header
	header := m.renderHeader(&sty)
	builder.WriteString(header)

	// Render tool calls if expanded or there's content
	if m.expanded || len(m.calls) > 0 {
		for i, call := range m.calls {
			if i > 0 {
				builder.WriteString("\n\n")
			}
			builder.WriteString(m.renderCall(&sty, call, width))
		}
	}

	// Add expand hint if collapsed and there are calls
	if !m.expanded && len(m.calls) > 0 {
		builder.WriteString("\n")
		expandHint := sty.Base.Foreground(sty.FgMuted).Italic(true).Render(
			fmt.Sprintf("[%d tool call(s) - press enter to expand]", len(m.calls)),
		)
		builder.WriteString(expandHint)
	}

	result := builder.String()

	// Apply focus styling
	if m.focused {
		result = lipgloss.NewStyle().Foreground(sty.Error).Render(result)
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

// renderHeader renders the tool message header.
func (m *ToolMessage) renderHeader(sty *styles.Styles) string {
	toolIcon := "🔧"
	headerText := fmt.Sprintf("%s %s", toolIcon, m.name)

	headerStyle := sty.Base.
		Bold(true).
		Foreground(sty.Primary)

	// Add status icon if there are calls
	if len(m.calls) > 0 {
		statusIcon := m.getOverallStatusIcon()
		headerText = fmt.Sprintf("%s %s", statusIcon, headerText)
	}

	return headerStyle.Render(headerText)
}

// renderCall renders a single tool call.
func (m *ToolMessage) renderCall(sty *styles.Styles, call ToolCall, width int) string {
	var parts []string

	// Tool name and status
	statusIcon := getStatusIcon(call.Status)
	callHeader := fmt.Sprintf("  %s %s", statusIcon, call.Name)

	if call.Arguments != nil && len(call.Arguments) > 0 {
		args := m.formatArguments(call.Arguments)
		callHeader += fmt.Sprintf(" %s", args)
	}

	headerStyle := sty.Base.Bold(true)
	if call.Status == ToolStatusError {
		headerStyle = headerStyle.Foreground(sty.Error)
	} else if call.Status == ToolStatusCompleted {
		headerStyle = headerStyle.Foreground(lipgloss.Color("#4ade80"))
	} else if call.Status == ToolStatusRunning {
		headerStyle = headerStyle.Foreground(lipgloss.Color("#fbbf24"))
	}

	parts = append(parts, headerStyle.Render(callHeader))

	// Render result or error if completed
	if call.Status == ToolStatusCompleted || call.Status == ToolStatusError {
		if m.expanded {
			if call.Error != "" {
				// Render error
				errorStyle := sty.Base.Foreground(sty.Error)
				parts = append(parts, "")
				parts = append(parts, errorStyle.Render(call.Error))
			} else if call.Result != "" {
				// Render result (basic text wrapping)
				resultLines := m.wrapText(call.Result, width-4)
				resultStyle := sty.Base.Foreground(sty.FgMuted)
				for _, line := range resultLines {
					parts = append(parts, resultStyle.Render("  "+line))
				}
			}
		}
	}

	return strings.Join(parts, "\n")
}

// wrapText wraps text to fit within the given width.
func (m *ToolMessage) wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		testLine := currentLine.String()
		if testLine == "" {
			testLine = word
		} else {
			testLine = testLine + " " + word
		}

		if lipgloss.Width(testLine) <= width {
			currentLine.WriteString(word + " ")
		} else {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
			}
			currentLine.Reset()
			currentLine.WriteString(word + " ")
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// formatArguments formats tool arguments for display.
func (m *ToolMessage) formatArguments(args map[string]any) string {
	if args == nil || len(args) == 0 {
		return ""
	}

	var argStrings []string
	for key, value := range args {
		argStrings = append(argStrings, fmt.Sprintf("%s=%v", key, value))
	}

	return fmt.Sprintf("(%s)", strings.Join(argStrings, ", "))
}

// getOverallStatusIcon returns an icon representing the overall status of all tool calls.
func (m *ToolMessage) getOverallStatusIcon() string {
	if len(m.calls) == 0 {
		return "⏸"
	}

	// Check for errors
	for _, call := range m.calls {
		if call.Status == ToolStatusError {
			return "❌"
		}
	}

	// Check if any are still running
	for _, call := range m.calls {
		if call.Status == ToolStatusRunning {
			return "⚡"
		}
	}

	// Check if any are pending
	for _, call := range m.calls {
		if call.Status == ToolStatusPending {
			return "⏳"
		}
	}

	// All completed
	return "✅"
}

// getStatusIcon returns the icon for a tool status.
func getStatusIcon(status ToolStatus) string {
	switch status {
	case ToolStatusPending:
		return "⏳"
	case ToolStatusRunning:
		return "⚡"
	case ToolStatusCompleted:
		return "✅"
	case ToolStatusError:
		return "❌"
	case ToolStatusCanceled:
		return "⏹"
	default:
		return "❓"
	}
}

// HandleClick handles click events for expanding/collapsing.
// Implements ClickHandler interface.
func (m *ToolMessage) HandleClick(line, col int) render.Cmd {
	if len(m.calls) > 0 {
		m.ToggleExpanded()
		return render.None()
	}
	return render.None()
}

// HandleKeyEvent handles keyboard events.
// Implements KeyEventHandler interface.
func (m *ToolMessage) HandleKeyEvent(msg any) (bool, render.Cmd) {
	keyMsg, ok := msg.(*render.KeyMsg)
	if !ok {
		return false, nil
	}

	switch keyMsg.String() {
	case " ", "enter":
		// Toggle expansion if there are tool calls
		if len(m.calls) > 0 {
			m.ToggleExpanded()
			return true, render.None()
		}
	}

	return false, nil
}

// CallCount returns the number of tool calls.
func (m *ToolMessage) CallCount() int {
	return len(m.calls)
}

// CompletedCount returns the number of completed tool calls.
func (m *ToolMessage) CompletedCount() int {
	count := 0
	for _, call := range m.calls {
		if call.Status == ToolStatusCompleted {
			count++
		}
	}
	return count
}

// FailedCount returns the number of failed tool calls.
func (m *ToolMessage) FailedCount() int {
	count := 0
	for _, call := range m.calls {
		if call.Status == ToolStatusError {
			count++
		}
	}
	return count
}
