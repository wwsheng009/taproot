package sessions

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/components/dialogs"
	"github.com/wwsheng009/taproot/ui/styles"
	"github.com/wwsheng009/taproot/tui/util"
)

const (
	ID dialogs.DialogID = "sessions"
)

// Session represents a chat session
type Session struct {
	ID        string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SessionProvider provides sessions for the dialog
type SessionProvider interface {
	Sessions() []Session
	CreateSession(title string) tea.Cmd
	DeleteSession(id string) tea.Cmd
	SelectSession(id string) tea.Cmd
}

// SessionsDialog displays a list of sessions with search and actions
type SessionsDialog struct {
	styles      *styles.Styles
	provider    SessionProvider
	width       int
	height      int
	filter      string
	cursor      int
	scroll      int
	mode        dialogMode // modeBrowse, modeCreate, modeConfirmDelete
	creatingTitle string
	deleteConfirmID string
}

type dialogMode int

const (
	modeBrowse dialogMode = iota
	modeCreate
	modeConfirmDelete
)

const (
	maxVisibleItems = 8
)

// New creates a new sessions dialog
func New(provider SessionProvider) *SessionsDialog {
	s := styles.DefaultStyles()
	return &SessionsDialog{
		styles:   &s,
		provider: provider,
		mode:     modeBrowse,
	}
}

func (d *SessionsDialog) Init() tea.Cmd {
	return nil
}

func (d *SessionsDialog) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return d.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
	}

	return d, nil
}

func (d *SessionsDialog) handleKeyMsg(msg tea.KeyMsg) (util.Model, tea.Cmd) {
	switch d.mode {
	case modeBrowse:
		return d.handleBrowseMode(msg)
	case modeCreate:
		return d.handleCreateMode(msg)
	case modeConfirmDelete:
		return d.handleDeleteConfirmMode(msg)
	}
	return d, nil
}

func (d *SessionsDialog) handleBrowseMode(msg tea.KeyMsg) (util.Model, tea.Cmd) {
	sessions := d.filteredSessions()

	switch msg.String() {
	case "esc":
		if d.filter != "" {
			d.filter = ""
			d.cursor = 0
			d.scroll = 0
			return d, nil
		}
		return d, func() tea.Msg { return dialogs.CloseDialogMsg{} }

	case "q":
		return d, func() tea.Msg { return dialogs.CloseDialogMsg{} }

	case "/":
		// Start filtering (just type)
		return d, nil

	case "n", "ctrl+n":
		// Create new session
		d.mode = modeCreate
		d.creatingTitle = ""
		return d, nil

	case "d", "ctrl+d":
		// Delete selected session
		if len(sessions) > 0 {
			d.deleteConfirmID = sessions[d.cursor].ID
			d.mode = modeConfirmDelete
		}
		return d, nil

	case "enter":
		// Select session
		if len(sessions) > 0 {
			selected := sessions[d.cursor]
			return d, tea.Sequence(
				d.provider.SelectSession(selected.ID),
				func() tea.Msg { return dialogs.CloseDialogMsg{} },
			)
		}
		return d, nil

	case "up", "k":
		if d.cursor > 0 {
			d.cursor--
			if d.cursor < d.scroll {
				d.scroll--
			}
		}

	case "down", "j":
		if d.cursor < len(sessions)-1 {
			d.cursor++
			if d.cursor >= d.scroll+maxVisibleItems {
				d.scroll++
			}
		}

	default:
		// Filter typing
		if len(msg.String()) == 1 && msg.Type == tea.KeyRunes {
			d.filter += msg.String()
			d.cursor = 0
			d.scroll = 0
		}
	}

	return d, nil
}

func (d *SessionsDialog) handleCreateMode(msg tea.KeyMsg) (util.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		d.mode = modeBrowse
		d.creatingTitle = ""
		return d, nil

	case "enter":
		if d.creatingTitle != "" {
			return d, d.provider.CreateSession(d.creatingTitle)
		}
		return d, nil

	case "backspace":
		if len(d.creatingTitle) > 0 {
			d.creatingTitle = d.creatingTitle[:len(d.creatingTitle)-1]
		}

	default:
		if len(msg.String()) == 1 && msg.Type == tea.KeyRunes {
			d.creatingTitle += msg.String()
		}
	}

	return d, nil
}

func (d *SessionsDialog) handleDeleteConfirmMode(msg tea.KeyMsg) (util.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "n":
		d.mode = modeBrowse
		d.deleteConfirmID = ""
		return d, nil

	case "y":
		if d.deleteConfirmID != "" {
			return d, d.provider.DeleteSession(d.deleteConfirmID)
		}
	}

	return d, nil
}

func (d *SessionsDialog) filteredSessions() []Session {
	allSessions := d.provider.Sessions()

	if d.filter == "" {
		return allSessions
	}

	filter := strings.ToLower(d.filter)
	var result []Session
	for _, s := range allSessions {
		if strings.Contains(strings.ToLower(s.Title), filter) ||
			strings.Contains(strings.ToLower(s.ID), filter) {
			result = append(result, s)
		}
	}
	return result
}

func (d *SessionsDialog) View() string {
	s := d.styles

	// Dialog box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Border).
		Padding(1, 2)

	// Calculate dialog size
	dialogWidth := min(d.width-4, 60)
	dialogHeight := min(d.height-4, 20)

	var sb strings.Builder

	// Header
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(s.Primary).Render("Sessions"))
	sb.WriteString("\n\n")

	// Mode-specific content
	switch d.mode {
	case modeBrowse:
		d.renderBrowseMode(&sb)
	case modeCreate:
		d.renderCreateMode(&sb)
	case modeConfirmDelete:
		d.renderDeleteConfirmMode(&sb)
	}

	// Footer hints
	sb.WriteString("\n")
	switch d.mode {
	case modeBrowse:
		hints := "Enter: Select | n: New | d: Delete | /: Filter | Esc: Close"
		sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render(hints))
	case modeCreate:
		hints := "Enter: Create | Esc: Cancel"
		sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render(hints))
	case modeConfirmDelete:
		hints := "y: Delete | n: Cancel"
		sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render(hints))
	}

	// Apply size and render
	rendered := boxStyle.Width(dialogWidth).MaxHeight(dialogHeight).Render(sb.String())

	return lipgloss.NewStyle().
		Width(d.width).
		Height(d.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(rendered)
}

func (d *SessionsDialog) renderBrowseMode(sb *strings.Builder) {
	s := d.styles
	sessions := d.filteredSessions()

	if d.filter != "" {
		sb.WriteString(lipgloss.NewStyle().Foreground(s.Secondary).Render("Filter: /"+d.filter))
		sb.WriteString("\n\n")
	}

	if len(sessions) == 0 {
		msg := "No sessions found"
		if d.filter != "" {
			msg = "No matching sessions"
		}
		sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Italic(true).Render(msg))
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render("Press n to create a new session"))
		return
	}

	// Calculate visible range
	end := d.scroll + maxVisibleItems
	if end > len(sessions) {
		end = len(sessions)
	}

	for i := d.scroll; i < end; i++ {
		sess := sessions[i]
		isSelected := i == d.cursor

		// Format title
		title := sess.Title
		if title == "" {
			title = "Untitled"
		}
		if len(title) > 30 {
			title = title[:27] + "..."
		}

		// Cursor
		cursor := " "
		if isSelected {
			cursor = ">"
		}

		// Style
		style := lipgloss.NewStyle().Foreground(s.FgBase)
		if isSelected {
			style = style.Foreground(s.Primary).Bold(true)
		}

		sb.WriteString(style.Render(cursor + " " + title))

		// Timestamp
		timeStr := formatTime(sess.UpdatedAt)
		if timeStr != "" {
			padding := strings.Repeat(" ", 35-len(title))
			sb.WriteString(padding + lipgloss.NewStyle().Foreground(s.FgMuted).Render(timeStr))
		}

		sb.WriteString("\n")
	}

	// Scroll indicator
	if len(sessions) > maxVisibleItems {
		scrollInfo := fmt.Sprintf("%d/%d", d.cursor+1, len(sessions))
		sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render(scrollInfo))
	}
}

func (d *SessionsDialog) renderCreateMode(sb *strings.Builder) {
	s := d.styles
	sb.WriteString(lipgloss.NewStyle().Foreground(s.Secondary).Render("Create new session"))
	sb.WriteString("\n\n")
	sb.WriteString("Title: ")
	sb.WriteString(lipgloss.NewStyle().Foreground(s.Primary).Render(d.creatingTitle))
	if time.Now().Unix()%2 == 0 {
		sb.WriteString("▋")
	}
	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render("Enter a title for the new session"))
}

func (d *SessionsDialog) renderDeleteConfirmMode(sb *strings.Builder) {
	s := d.styles
	sb.WriteString(lipgloss.NewStyle().Foreground(s.Error).Bold(true).Render("Delete Session?"))
	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render("This action cannot be undone."))
	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render("Press y to confirm, n to cancel"))
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	diff := time.Since(t)
	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		return fmt.Sprintf("%dm", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh", int(diff.Hours()))
	}
	if diff < 30*24*time.Hour {
		return fmt.Sprintf("%dd", int(diff.Hours()/24))
	}
	return t.Format("Jan 2")
}

func (d *SessionsDialog) Position() (int, int) {
	return 0, 0
}

func (d *SessionsDialog) ID() dialogs.DialogID {
	return ID
}
