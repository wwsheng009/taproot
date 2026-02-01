package reasoning

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/components/dialogs"
	"github.com/wwsheng009/taproot/ui/styles"
	"github.com/wwsheng009/taproot/tui/util"
)

const (
	ID dialogs.DialogID = "reasoning"
)

// ReasoningDialog displays collapsible reasoning/thought content
type ReasoningDialog struct {
	styles     *styles.Styles
	width      int
	height     int
	expanded   bool
	content    string
	visibleLines int
	scroll     int
}

// New creates a new reasoning dialog
func New(content string) *ReasoningDialog {
	s := styles.DefaultStyles()
	return &ReasoningDialog{
		styles:       &s,
		expanded:     false,
		content:      content,
		visibleLines: 5, // Collapsed height
		scroll:       0,
	}
}

// SetContent updates the reasoning content (for streaming updates)
func (d *ReasoningDialog) SetContent(content string) tea.Cmd {
	d.content = content
	return nil
}

func (d *ReasoningDialog) Init() tea.Cmd {
	return nil
}

func (d *ReasoningDialog) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ReasoningUpdateMsg:
		// Handle streaming content updates
		d.content = msg.Content
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			// Toggle expand/collapse
			d.expanded = !d.expanded
			if d.expanded {
				d.visibleLines = d.height - 8 // Account for borders and padding
			} else {
				d.visibleLines = 5
				d.scroll = 0
			}
		case "esc":
			return d, func() tea.Msg { return dialogs.CloseDialogMsg{} }
		case "up", "k":
			if d.expanded && d.scroll > 0 {
				d.scroll--
			}
		case "down", "j":
			if d.expanded && d.canScrollDown() {
				d.scroll++
			}
		}
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		if d.expanded {
			d.visibleLines = d.height - 8
		}
	}

	return d, nil
}

func (d *ReasoningDialog) canScrollDown() bool {
	lines := strings.Split(d.content, "\n")
	return d.scroll+d.visibleLines < len(lines)
}

func (d *ReasoningDialog) View() string {
	s := d.styles

	// Dialog box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Border).
		Padding(1, 2)

	// Calculate content width
	contentWidth := max(d.width-8, 20)

	var sb strings.Builder

	// Header with expand indicator
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(s.Primary).Render("Reasoning"))
	if d.expanded {
		sb.WriteString(" " + lipgloss.NewStyle().Foreground(s.FgMuted).Render("[-]"))
	} else {
		sb.WriteString(" " + lipgloss.NewStyle().Foreground(s.FgMuted).Render("[+]"))
	}

	// Content with scrolling
	lines := strings.Split(d.content, "\n")
	totalLines := len(lines)

	// Show ellipsis if scrollable
	if d.expanded && totalLines > d.visibleLines {
		// We have scrolling
	}

	// Render visible lines
	end := min(d.scroll+d.visibleLines, totalLines)

	if d.scroll < totalLines {
		for i := d.scroll; i < end; i++ {
			if i < len(lines) {
				line := lines[i]
				// Truncate long lines
				if len(line) > contentWidth {
					line = line[:contentWidth-3] + "..."
				}
				sb.WriteString("\n" + lipgloss.NewStyle().Foreground(s.FgBase).Render(line))
			}
		}
	}

	// Scroll indicators
	if d.expanded {
		if d.scroll > 0 {
			// Show up arrow indicator
		}
		if d.canScrollDown() {
			// Show down arrow indicator
		}
	} else if totalLines > d.visibleLines {
		// Show "more" indicator in collapsed mode
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Italic(true).Render("..."))
	}

	// Hints footer
	sb.WriteString("\n\n")
	hints := lipgloss.NewStyle().Foreground(s.FgMuted).Render("Enter: Toggle | Esc: Close")
	if d.expanded {
		hints = lipgloss.NewStyle().Foreground(s.FgMuted).Render("Enter: Collapse | ↑↓: Scroll | Esc: Close")
	}
	sb.WriteString(hints)

	// Apply box style and center
	dialogWidth := min(contentWidth+6, d.width-4)

	rendered := boxStyle.Width(dialogWidth).Render(sb.String())

	// Center horizontally
	return lipgloss.NewStyle().
		Width(d.width).
		Align(lipgloss.Center).
		Render(rendered)
}

func (d *ReasoningDialog) Position() (int, int) {
	return 0, 0
}

func (d *ReasoningDialog) ID() dialogs.DialogID {
	return ID
}

// ReasoningUpdateMsg is sent when reasoning content is updated
type ReasoningUpdateMsg struct {
	Content string
}

// UpdateReasoning returns a command that updates the reasoning content
func UpdateReasoning(content string) tea.Cmd {
	return func() tea.Msg {
		return ReasoningUpdateMsg{Content: content}
	}
}
