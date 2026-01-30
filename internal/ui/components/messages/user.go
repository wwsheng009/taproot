package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/render"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// UserMessage represents a user message in the chat UI.
type UserMessage struct {
	id              string
	content         string
	timestamp       time.Time
	attachments     []Attachment
	focused         bool
	config          MessageConfig
	maxWidth        int
	cachedRender    string
	cachedWidth     int
	cachedHeight    int
	cacheValid      bool
	initialized     bool
}

// NewUserMessage creates a new user message.
func NewUserMessage(id, content string) *UserMessage {
	return &UserMessage{
		id:           id,
		content:      content,
		timestamp:    time.Now(),
		attachments:  []Attachment{},
		focused:      false,
		config:       DefaultMessageConfig(),
		maxWidth:     80,
		cacheValid:   false,
		initialized:  false,
	}
}

// Init initializes the component.
// Implements render.Model interface.
func (m *UserMessage) Init() error {
	m.initialized = true
	return nil
}

// Update handles incoming messages and returns updated state and commands.
// Implements render.Model interface.
func (m *UserMessage) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		m.focused = true
		m.cacheValid = false
	case *render.BlurMsg:
		m.focused = false
		m.cacheValid = false
	}
	return m, render.None()
}

// View returns the string representation for rendering.
// Implements render.Model interface.
func (m *UserMessage) View() string {
	sty := styles.DefaultStyles()

	// Use configured max width or component max width
	width := m.config.MaxWidth
	if m.maxWidth > 0 && m.maxWidth < width {
		width = m.maxWidth
	}

	content, _ := m.renderContent(width, sty)
	// Apply user-focused or user-blurred style
	if m.focused {
		return lipgloss.NewStyle().Foreground(sty.Primary).Bold(true).Render(content)
	}
	return sty.Base.Foreground(sty.FgBase).Render(content)
}

// ID returns the unique identifier.
// Implements Identifiable interface.
func (m *UserMessage) ID() string {
	return m.id
}

// Role returns the message role.
// Implements Message interface.
func (m *UserMessage) Role() MessageRole {
	return RoleUser
}

// Content returns the message content.
// Implements Message interface.
func (m *UserMessage) Content() string {
	return m.content
}

// Timestamp returns the message timestamp.
// Implements Message interface.
func (m *UserMessage) Timestamp() time.Time {
	return m.timestamp
}

// SetContent sets the message content.
// Implements Message interface.
func (m *UserMessage) SetContent(content string) {
	m.content = content
	m.cacheValid = false
}

// Focus focuses the component.
// Implements Focusable interface.
func (m *UserMessage) Focus() {
	m.focused = true
	m.cacheValid = false
}

// Blur unfocuses the component.
// Implements Focusable interface.
func (m *UserMessage) Blur() {
	m.focused = false
	m.cacheValid = false
}

// Focused returns whether the component is focused.
// Implements Focusable interface.
func (m *UserMessage) Focused() bool {
	return m.focused
}

// Attachments returns the attachments.
func (m *UserMessage) Attachments() []Attachment {
	return m.attachments
}

// SetAttachments sets the attachments.
func (m *UserMessage) SetAttachments(attachments []Attachment) {
	m.attachments = attachments
	m.cacheValid = false
}

// AddAttachment adds an attachment.
func (m *UserMessage) AddAttachment(attachment Attachment) {
	m.attachments = append(m.attachments, attachment)
	m.cacheValid = false
}

// SetTimestamp sets the message timestamp.
func (m *UserMessage) SetTimestamp(timestamp time.Time) {
	m.timestamp = timestamp
}

// Config returns the current configuration.
func (m *UserMessage) Config() MessageConfig {
	return m.config
}

// SetConfig sets the message configuration.
func (m *UserMessage) SetConfig(config MessageConfig) {
	m.config = config
	m.cacheValid = false
}

// MaxWidth returns the maximum width.
func (m *UserMessage) MaxWidth() int {
	return m.maxWidth
}

// SetMaxWidth sets the maximum width.
func (m *UserMessage) SetMaxWidth(width int) {
	m.maxWidth = width
	m.cacheValid = false
}

// invalidateCache marks the cached render as invalid.
func (m *UserMessage) invalidateCache() {
	m.cacheValid = false
}

// renderContent renders the message content.
func (m *UserMessage) renderContent(width int, sty styles.Styles) (string, int) {
	// Check cache
	if m.cacheValid && m.cachedWidth == width {
		return m.cachedRender, m.cachedHeight
	}

	var b strings.Builder

	// Render attachments if present
	if len(m.attachments) > 0 && !m.config.CompactMode {
		attachments := m.renderAttachments(sty)
		if attachments != "" {
			b.WriteString(attachments)
			b.WriteString("\n\n")
		}
	}

	// Render message content
	if m.content != "" {
		content := m.renderContentText(m.content, width, sty)
		b.WriteString(content)
	}

	// Render timestamp if enabled
	if m.config.ShowTimestamp && !m.config.CompactMode {
		b.WriteString("\n")
		timestamp := m.renderTimestamp(sty)
		b.WriteString(timestamp)
	}

	content := b.String()
	height := lipgloss.Height(content)

	// Update cache
	m.cachedRender = content
	m.cachedWidth = width
	m.cachedHeight = height
	m.cacheValid = true

	return content, height
}

// renderAttachments renders the attachments.
func (m *UserMessage) renderAttachments(sty styles.Styles) string {
	var b strings.Builder

	for i, attachment := range m.attachments {
		if i > 0 {
			b.WriteString(" ")
		}

		// Render attachment based on type
		switch attachment.Type {
		case "file":
			icon := "📄"
			name := sty.Base.Foreground(sty.FgBase).Render(attachment.Name)
			format := sty.Base.Foreground(sty.FgMuted).Render(fmt.Sprintf("(%d B)", attachment.Size))
			b.WriteString(fmt.Sprintf("%s %s %s", icon, name, format))
		case "image":
			icon := "🖼️"
			name := sty.Base.Foreground(sty.FgBase).Render(attachment.Name)
			format := sty.Base.Foreground(sty.FgMuted).Render(fmt.Sprintf("%dx%d", attachment.Size, attachment.Size))
			b.WriteString(fmt.Sprintf("%s %s %s", icon, name, format))
		default:
			icon := "📎"
			name := sty.Base.Foreground(sty.FgBase).Render(attachment.Name)
			b.WriteString(fmt.Sprintf("%s %s", icon, name))
		}
	}

	return b.String()
}

// renderContentText renders the message text content.
func (m *UserMessage) renderContentText(text string, width int, sty styles.Styles) string {
	if m.config.EnableMarkdown {
		// Use Markdown rendering - simplified approach
		// Note: In a full implementation, use glamour with the Markdown config
		// For now, just render with basic styling
		return sty.Base.Width(width).Render(text)
	}

	// Simple text rendering
	wrapped := sty.Base.Width(width).Render(text)
	return wrapped
}

// renderTimestamp renders the timestamp.
func (m *UserMessage) renderTimestamp(sty styles.Styles) string {
	relativeTime := m.formatRelativeTime(m.timestamp)
	timestamp := sty.Base.Foreground(sty.FgMuted).Render(relativeTime)
	return timestamp
}

// formatRelativeTime formats the timestamp as a relative time string.
func (m *UserMessage) formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "now"
	}

	if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", minutes)
	}

	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}

	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}

	return t.Format("Jan 2")
}

// HasAttachments returns whether the message has attachments.
func (m *UserMessage) HasAttachments() bool {
	return len(m.attachments) > 0
}

// AttachmentCount returns the number of attachments.
func (m *UserMessage) AttachmentCount() int {
	return len(m.attachments)
}

// GetAttachment returns an attachment by index.
func (m *UserMessage) GetAttachment(index int) (Attachment, bool) {
	if index < 0 || index >= len(m.attachments) {
		return Attachment{}, false
	}
	return m.attachments[index], true
}

// RemoveAttachment removes an attachment by index.
func (m *UserMessage) RemoveAttachment(index int) bool {
	if index < 0 || index >= len(m.attachments) {
		return false
	}
	m.attachments = append(m.attachments[:index], m.attachments[index+1:]...)
	m.cacheValid = false
	return true
}

// ClearAttachments removes all attachments.
func (m *UserMessage) ClearAttachments() {
	m.attachments = []Attachment{}
	m.cacheValid = false
}

// WordCount returns the word count of the content.
func (m *UserMessage) WordCount() int {
	words := strings.Fields(m.content)
	return len(words)
}

// CharacterCount returns the character count of the content.
func (m *UserMessage) CharacterCount() int {
	return len(m.content)
}
