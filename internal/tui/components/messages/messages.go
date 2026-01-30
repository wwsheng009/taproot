package messages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/tui/util"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// Message represents a single message in the chat
type Message struct {
	Role    string // "user", "assistant", "system", "tool"
	Content string
	ToolUse *ToolUse // Optional, for tool calls
}

// ToolUse represents a tool invocation
type ToolUse struct {
	Name      string
	Arguments string
	Result    string
}

// MessagesModel displays a list of messages with scrolling
type MessagesModel struct {
	width    int
	height   int
	messages []Message
	scroll   int
	cursor   int // For potential selection
	styles   *styles.Styles
}

const (
	maxWidth = 80
)

// New creates a new messages model
func New() *MessagesModel {
	s := styles.DefaultStyles()
	return &MessagesModel{
		messages: []Message{},
		scroll:   0,
		styles:   &s,
	}
}

// SetWidth sets the width for rendering
func (m *MessagesModel) SetWidth(w int) {
	m.width = w
}

// SetHeight sets the height for rendering
func (m *MessagesModel) SetHeight(h int) {
	m.height = h
}

// AddMessage appends a new message
func (m *MessagesModel) AddMessage(msg Message) tea.Cmd {
	m.messages = append(m.messages, msg)
	return m.scrollToBottom()
}

// Clear removes all messages
func (m *MessagesModel) Clear() {
	m.messages = []Message{}
	m.scroll = 0
}

func (m *MessagesModel) scrollToBottom() tea.Cmd {
	// Calculate total content height
	totalHeight := 0
	for _, msg := range m.messages {
		totalHeight += m.messageHeight(msg)
	}

	// Auto-scroll to bottom if content is taller than viewport
	if totalHeight > m.height {
		m.scroll = totalHeight - m.height
	}

	return nil
}

func (m *MessagesModel) Init() tea.Cmd {
	return nil
}

func (m *MessagesModel) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			maxScroll := m.maxScroll()
			if m.scroll < maxScroll {
				m.scroll++
			}
		case "home", "g":
			m.scroll = 0
		case "end", "G":
			m.scroll = m.maxScroll()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Ensure scroll is within bounds
		maxScroll := m.maxScroll()
		if m.scroll > maxScroll {
			m.scroll = maxScroll
		}
	}

	return m, nil
}

func (m *MessagesModel) maxScroll() int {
	totalHeight := 0
	for _, msg := range m.messages {
		totalHeight += m.messageHeight(msg)
	}
	if totalHeight <= m.height {
		return 0
	}
	return totalHeight - m.height
}

func (m *MessagesModel) messageHeight(msg Message) int {
	contentWidth := m.width - 4 // Account for padding
	if contentWidth > maxWidth {
		contentWidth = maxWidth
	}

	lines := 0

	// Header (role)
	lines++

	// Content (rough estimate - wrapped lines)
	contentLines := strings.Count(msg.Content, "\n") + 1
	wrappedWidth := len(msg.Content) / contentWidth
	if wrappedWidth > 0 {
		contentLines += wrappedWidth
	}
	lines += contentLines

	// Tool use if present
	if msg.ToolUse != nil {
		lines += 3 // Tool header
		if msg.ToolUse.Arguments != "" {
			lines += strings.Count(msg.ToolUse.Arguments, "\n") + 1
		}
		if msg.ToolUse.Result != "" {
			lines += 2 // Result header + content
			lines += strings.Count(msg.ToolUse.Result, "\n") + 1
		}
	}

	// Spacing
	lines += 1 // Blank line after each message

	return lines
}

func (m *MessagesModel) View() string {
	s := m.styles

	if len(m.messages) == 0 {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(lipgloss.NewStyle().Foreground(s.FgMuted).Italic(true).Render("No messages yet"))
	}

	var sb strings.Builder
	currentY := 0

	for _, msg := range m.messages {
		msgHeight := m.messageHeight(msg)
		msgView := m.renderMessage(msg, s)
		msgLines := strings.Split(msgView, "\n")

		// Skip if above scroll position
		if currentY+msgHeight <= m.scroll {
			currentY += msgHeight
			continue
		}

		// Render visible lines
		startLine := 0
		if currentY < m.scroll {
			startLine = m.scroll - currentY
		}

		for i := startLine; i < len(msgLines); i++ {
			if currentY+i >= m.scroll+m.height {
				break
			}
			sb.WriteString(msgLines[i] + "\n")
		}

		currentY += msgHeight

		// Stop if we've filled the viewport
		if currentY >= m.scroll+m.height {
			break
		}
	}

	return sb.String()
}

func (m *MessagesModel) renderMessage(msg Message, s *styles.Styles) string {
	contentWidth := m.width - 4
	if contentWidth > maxWidth {
		contentWidth = maxWidth
	}

	var sb strings.Builder

	// Role header with color
	var roleStyle lipgloss.Style
	var roleName string

	switch msg.Role {
	case "user":
		roleStyle = lipgloss.NewStyle().Foreground(s.Primary).Bold(true)
		roleName = "You"
	case "assistant":
		roleStyle = lipgloss.NewStyle().Foreground(s.Secondary).Bold(true)
		roleName = "Assistant"
	case "system":
		roleStyle = lipgloss.NewStyle().Foreground(s.Warning).Bold(true)
		roleName = "System"
	case "tool":
		roleStyle = lipgloss.NewStyle().Foreground(s.Info).Bold(true)
		roleName = "Tool"
	default:
		roleStyle = lipgloss.NewStyle().Foreground(s.FgMuted).Bold(true)
		roleName = msg.Role
	}

	sb.WriteString(roleStyle.Render("╭─ " + roleName))
	sb.WriteString("\n")

	// Content with wrapping
	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Foreground(s.FgBase)

	sb.WriteString(contentStyle.Render(msg.Content))
	sb.WriteString("\n")

	// Tool use if present
	if msg.ToolUse != nil {
		toolStyle := lipgloss.NewStyle().
			Foreground(s.Info).
			Bold(true)

		sb.WriteString(toolStyle.Render("  🔧 " + msg.ToolUse.Name))
		sb.WriteString("\n")

		if msg.ToolUse.Arguments != "" {
			argStyle := lipgloss.NewStyle().
				Width(contentWidth - 2).
				Foreground(s.FgMuted)
			sb.WriteString(argStyle.Render(msg.ToolUse.Arguments))
			sb.WriteString("\n")
		}

		if msg.ToolUse.Result != "" {
			sb.WriteString(lipgloss.NewStyle().Foreground(s.Green).Render("  ✓ Result"))
			sb.WriteString("\n")

			resultStyle := lipgloss.NewStyle().
				Width(contentWidth - 2).
				Foreground(s.FgBase)
			sb.WriteString(resultStyle.Render(msg.ToolUse.Result))
			sb.WriteString("\n")
		}
	}

	// Footer
	sb.WriteString(roleStyle.Render("╰─ " + strings.Repeat("─", contentWidth-2)))
	sb.WriteString("\n")

	return sb.String()
}

// Size returns the current size
func (m *MessagesModel) Size() (int, int) {
	return m.width, m.height
}

// ScrollToBottom returns a command that scrolls to the bottom
func (m *MessagesModel) ScrollToBottom() tea.Cmd {
	return func() tea.Msg {
		m.scrollToBottom()
		return nil
	}
}
