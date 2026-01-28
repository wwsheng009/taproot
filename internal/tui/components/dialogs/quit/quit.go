package quit

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/tui/components/dialogs"
	"github.com/yourorg/taproot/internal/tui/styles"
	"github.com/yourorg/taproot/internal/tui/util"
)

const (
	ID dialogs.DialogID = "quit"
)

type QuitDialog struct {
	id         dialogs.DialogID
	width      int
	height     int
	hasChanges bool
	selected   int // 0: Cancel, 1: Quit
}

func New(hasChanges bool) *QuitDialog {
	return &QuitDialog{
		id:         ID,
		hasChanges: hasChanges,
		selected:   0,
	}
}

func (m *QuitDialog) Init() tea.Cmd {
	return nil
}

func (m *QuitDialog) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.selected = 0
		case "right", "l":
			m.selected = 1
		case "enter":
			if m.selected == 1 {
				return m, tea.Quit
			}
			return m, func() tea.Msg { return dialogs.CloseDialogMsg{} }
		case "esc":
			return m, func() tea.Msg { return dialogs.CloseDialogMsg{} }
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *QuitDialog) View() string {
	theme := styles.CurrentTheme()
	
	// Box style
	s := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2).
		Align(lipgloss.Center)

	var sb strings.Builder
	
	// Title
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("Confirm Quit"))
	sb.WriteString("\n\n")
	
	// Message
	if m.hasChanges {
		sb.WriteString(lipgloss.NewStyle().Foreground(theme.Warning).Render("⚠️  You have unsaved changes!"))
		sb.WriteString("\n")
		sb.WriteString("Are you sure you want to quit?")
	} else {
		sb.WriteString("Are you sure you want to quit?")
	}
	sb.WriteString("\n\n")
	
	// Buttons
	btnStyle := lipgloss.NewStyle().Padding(0, 2).Foreground(theme.FgBase)
	selectedStyle := lipgloss.NewStyle().Padding(0, 2).Background(theme.Primary).Foreground(theme.FgSelected).Bold(true)
	
	cancelBtn := "Cancel"
	quitBtn := "Quit"
	
	if m.selected == 0 {
		cancelBtn = selectedStyle.Render(cancelBtn)
		quitBtn = btnStyle.Render(quitBtn)
	} else {
		cancelBtn = btnStyle.Render(cancelBtn)
		quitBtn = selectedStyle.Render(quitBtn)
	}
	
	sb.WriteString(cancelBtn + "   " + quitBtn)
	
	return s.Render(sb.String())
}

func (m *QuitDialog) Position() (int, int) {
	return 0, 0
}

func (m *QuitDialog) ID() dialogs.DialogID {
	return m.id
}
