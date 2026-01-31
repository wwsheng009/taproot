package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

const (
	// maxThinkingHeight defines the maximum number of lines to show when thinking is collapsed.
	maxThinkingHeight = 10

	// thinkingTruncateFormat is the text shown when thinking content is truncated.
	thinkingTruncateFormat = "… (%d lines hidden) [click or space to expand]"
)

// AssistantMessage represents an assistant chat message with thinking/reasoning content.
//
// This component supports:
// - Collapsible thinking/reasoning content box
// - Markdown rendering for main content
// - Focus state with different styling
// - Finish reason display (canceled, error)
// - Relative timestamp
// - Caching for performance
type AssistantMessage struct {
	id                string
	content           string
	thinking          string
	timestamp         time.Time
	focused           bool
	finished          bool
	finishReason      string
	cancelled         bool
	thinkingExpanded  bool
	config            *MessageConfig
	maxWidth          int

	// Caching fields
	cachedRender      string
	cachedWidth       int
	cachedHeight      int
	cacheValid        bool
}

// NewAssistantMessage creates a new AssistantMessage component.
func NewAssistantMessage(id, content string) *AssistantMessage {
	return &AssistantMessage{
		id:               id,
		content:          content,
		timestamp:        time.Now(),
		focused:          false,
		finished:         true,
		thinkingExpanded: false,
		config:           &MessageConfig{},
	}
}

// Init initializes the component. Implements render.Model.
func (m *AssistantMessage) Init() error {
	return nil
}

// Update handles incoming messages. Implements render.Model.
func (m *AssistantMessage) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		m.Focus()
	case *render.BlurMsg:
		m.focused = false
	}

	// Invalidate cache on any update
	m.cacheValid = false

	return m, render.None()
}

// View renders the assistant message. Implements render.Model.
func (m *AssistantMessage) View() string {
	width := m.config.MaxWidth
	// Apply max width limit
	if m.maxWidth > 0 && width > m.maxWidth {
		width = m.maxWidth
	}

	return m.render(width)
}

// ID returns the message ID. Implements Identifiable interface.
func (m *AssistantMessage) ID() string {
	return m.id
}

// Role returns RoleAssistant. Implements Message interface.
func (m *AssistantMessage) Role() MessageRole {
	return RoleAssistant
}

// Content returns the main message content. Implements Message interface.
func (m *AssistantMessage) Content() string {
	return m.content
}

// SetContent sets the main message content. Implements Message interface.
func (m *AssistantMessage) SetContent(content string) {
	m.content = content
	m.cacheValid = false
}

// Timestamp returns the message timestamp. Implements Message interface.
func (m *AssistantMessage) Timestamp() time.Time {
	return m.timestamp
}

// SetTimestamp sets the message timestamp.
func (m *AssistantMessage) SetTimestamp(ts time.Time) {
	m.timestamp = ts
	m.cacheValid = false
}

// Thinking returns the thinking/reasoning content.
func (m *AssistantMessage) Thinking() string {
	return m.thinking
}

// SetThinking sets the thinking/reasoning content.
func (m *AssistantMessage) SetThinking(thinking string) {
	m.thinking = thinking
	m.cacheValid = false
}

// ThinkingExpanded returns true if thinking content is expanded.
func (m *AssistantMessage) ThinkingExpanded() bool {
	return m.thinkingExpanded
}

// SetThinkingExpanded sets the thinking expansion state.
func (m *AssistantMessage) SetThinkingExpanded(expanded bool) {
	m.thinkingExpanded = expanded
	m.cacheValid = false
}

// ToggleExpanded toggles the thinking content expansion state.
func (m *AssistantMessage) ToggleExpanded() {
	m.thinkingExpanded = !m.thinkingExpanded
	m.cacheValid = false
}

// Finished returns true if the message is finished.
func (m *AssistantMessage) Finished() bool {
	return m.finished
}

// SetFinished sets the finished state.
func (m *AssistantMessage) SetFinished(finished bool) {
	m.finished = finished
	m.cacheValid = false
}

// FinishReason returns the finish reason.
func (m *AssistantMessage) FinishReason() string {
	return m.finishReason
}

// SetFinishReason sets the finish reason.
func (m *AssistantMessage) SetFinishReason(reason string) {
	m.finishReason = reason
	m.cacheValid = false
}

// Canceled returns true if the message was canceled.
func (m *AssistantMessage) Canceled() bool {
	return m.cancelled
}

// SetCanceled sets the canceled state.
func (m *AssistantMessage) SetCanceled(cancelled bool) {
	m.cancelled = cancelled
	m.cacheValid = false
}

// Focus focuses the component. Implements Focusable interface.
func (m *AssistantMessage) Focus() {
	m.focused = true
	m.cacheValid = false
}

// Blur blurs the component. Implements Focusable interface.
func (m *AssistantMessage) Blur() {
	m.focused = false
	m.cacheValid = false
}

// Focused returns true if the component is focused. Implements Focusable interface.
func (m *AssistantMessage) Focused() bool {
	return m.focused
}

// SetMaxWidth sets the maximum width for rendering.
func (m *AssistantMessage) SetMaxWidth(width int) {
	m.maxWidth = width
	m.cacheValid = false
}

// SetConfig sets the message configuration.
func (m *AssistantMessage) SetConfig(config *MessageConfig) {
	m.config = config
	m.cacheValid = false
}

// render renders the assistant message with the given width.
func (m *AssistantMessage) render(width int) string {
	// Check cache first
	if m.cacheValid && m.cachedWidth == width && m.cachedRender != "" {
		return m.cachedRender
	}

	sty := styles.DefaultStyles()

	var parts []string

	// Show "Thinking..." text if not finished and no content
	if !m.finished && m.content == "" && m.thinking == "" {
		thinkingText := sty.Base.Foreground(sty.FgMuted).Italic(true).Render("Thinking...")
		parts = append(parts, thinkingText)
	}

	// Render thinking content if present
	thinking := strings.TrimSpace(m.thinking)
	if thinking != "" {
		parts = append(parts, m.renderThinking(width, thinking))
	}

	// Render main content if present
	content := strings.TrimSpace(m.content)
	if content != "" {
		// Add spacer between thinking and content
		if thinking != "" {
			parts = append(parts, "")
		}
		parts = append(parts, m.renderContent(width, content))
	}

	// Render finish reason info if finished
	if m.finished {
		if m.cancelled {
			canceledText := sty.Base.Foreground(sty.FgMuted).Italic(true).Render("Canceled")
			parts = append(parts, canceledText)
		} else if m.finishReason == "error" {
			parts = append(parts, m.renderError(width))
		}
	}

	var result string
	if len(parts) > 0 {
		result = strings.Join(parts, "\n")
	}

	// Apply focus styling
	if m.focused {
		result = lipgloss.NewStyle().Foreground(sty.Secondary).Render(result)
	} else {
		result = lipgloss.NewStyle().Foreground(sty.FgBase).Render(result)
	}

	// Cache the result
	m.cachedRender = result
	m.cachedWidth = width
	m.cachedHeight = lipgloss.Height(result)
	m.cacheValid = true

	return result
}

// renderThinking renders the thinking/reasoning content with truncation and footer.
func (m *AssistantMessage) renderThinking(width int, thinking string) string {
	sty := styles.DefaultStyles()

	// Simplified markdown rendering (no glamour for now)
	rendered := strings.TrimSpace(thinking)
	rendered = strings.ReplaceAll(rendered, "\n\n", "\n")
	rendered = strings.ReplaceAll(rendered, "\n\n\n", "\n\n")

	lines := strings.Split(rendered, "\n")
	totalLines := len(lines)

	// Truncate if collapsed and too many lines
	displayLines := lines
	if !m.thinkingExpanded && totalLines > maxThinkingHeight {
		displayLines = lines[totalLines-maxThinkingHeight:]

		truncationHint := sty.Base.Foreground(sty.FgMuted).Italic(true).Render(
			fmt.Sprintf(thinkingTruncateFormat, totalLines-maxThinkingHeight),
		)
		displayLines = append([]string{truncationHint, ""}, displayLines...)
	}

	// Render thinking box with border
	borderStyle := lipgloss.NewStyle().
		Foreground(sty.Border).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(sty.Border).
		Padding(0, 1).
		Width(width)

	thinkingBox := borderStyle.Render(strings.Join(displayLines, "\n"))

	var footer string
	// Add footer if thinking is done
	if m.finished {
		duration := time.Since(m.timestamp).Round(time.Second)
		if duration > 0 {
			durationText := sty.Base.Foreground(sty.Primary).Bold(true).Render("Thought for ") +
				sty.Base.Foreground(sty.FgMuted).Render(duration.String())
			footer = durationText
		}
	}

	if footer != "" {
		thinkingBox += "\n\n" + footer
	}

	return thinkingBox
}

// renderContent renders the main content with markdown support.
func (m *AssistantMessage) renderContent(width int, content string) string {
	sty := styles.DefaultStyles()

	// Simplified markdown rendering
	rendered := strings.TrimSpace(content)

	// Basic code block detection and formatting
	lines := strings.Split(rendered, "\n")
	var formattedLines []string
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for code block markers
		if trimmed == "```" {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			// Style code block content
			formattedLines = append(formattedLines, sty.Base.Foreground(sty.FgMuted).Render(line))
		} else {
			formattedLines = append(formattedLines, line)
		}
	}

	rendered = strings.Join(formattedLines, "\n")
	rendered = strings.TrimSpace(rendered)

	return rendered
}

// renderError renders an error message.
func (m *AssistantMessage) renderError(width int) string {
	sty := styles.DefaultStyles()

	// Simple error rendering
	errorTag := sty.Base.
		Foreground(lipgloss.Color("204")).
		Bold(true).
		Render("ERROR")

	errorText := sty.Base.
		Bold(true).
		Render("An error occurred")

	return fmt.Sprintf("%s %s", errorTag, errorText)
}

// WordCount returns the number of words in the content.
func (m *AssistantMessage) WordCount() int {
	words := strings.Fields(m.content)
	return len(words)
}

// CharacterCount returns the number of characters in the content.
func (m *AssistantMessage) CharacterCount() int {
	return len(m.content)
}

// GetThinkingBoxHeight returns the height of the thinking box in lines.
// This is useful for click detection.
func (m *AssistantMessage) GetThinkingBoxHeight() int {
	if m.thinking == "" {
	return 0
	}

	width := m.config.MaxWidth
	if m.maxWidth > 0 && width > m.maxWidth {
		width = m.maxWidth
	}

	// Render to calculate height
	thinkingBox := m.renderThinking(width, m.thinking)
	return lipgloss.Height(thinkingBox)
}

// HandleClick handles click events, particularly for expanding/collapsing thinking.
// Implements ClickHandler interface.
func (m *AssistantMessage) HandleClick(line, col int) render.Cmd {
	thinkingHeight := m.GetThinkingBoxHeight()
	if m.thinking != "" && thinkingHeight > 0 && line <= thinkingHeight {
		m.ToggleExpanded()
		return render.None()
	}
	return render.None()
}

// HandleKeyEvent handles keyboard events.
// Implements KeyEventHandler interface.
func (m *AssistantMessage) HandleKeyEvent(msg any) (bool, render.Cmd) {
	keyMsg, ok := msg.(*render.KeyMsg)
	if !ok {
		return false, nil
	}
	
	switch keyMsg.String() {
	case "c", "y":
		// Copy to clipboard would go here
		// For now, just return true to indicate handled
		return true, nil
	case " ", "enter":
		// Toggle thinking expansion if thinking content exists
		if m.thinking != "" {
			m.ToggleExpanded()
			return true, render.None()
		}
	}

	return false, nil
}
