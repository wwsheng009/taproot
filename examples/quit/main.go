package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/tui/app"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs/quit"
	"github.com/wwsheng009/taproot/internal/tui/page"
	"github.com/wwsheng009/taproot/internal/tui/util"
)

const (
	pageHome page.PageID = "home"
)

func main() {
	application := app.NewApp()
	homePage := NewHomePage()
	application.RegisterPage(pageHome, homePage)
	application.SetPage(pageHome)

	p := tea.NewProgram(application, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

// HomePage is the home page
type HomePage struct {
	hasChanges bool
}

func NewHomePage() HomePage {
	return HomePage{hasChanges: false}
}

func (h HomePage) Init() tea.Cmd {
	return nil
}

func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if h.hasChanges {
				// Show quit confirmation dialog
				dialog := quit.New(h.hasChanges)
				return h, func() tea.Msg {
					return dialogs.OpenDialogMsg{Model: dialog}
				}
			}
			return h, tea.Quit
		case "ctrl+c":
			return h, tea.Quit
		case "s":
			// Simulate making changes
			h.hasChanges = true
			return h, util.ReportInfo("Changes made! Press q to see confirmation")
		}
	case util.InfoMsg:
		// Clear changes after a save (simulation)
		if msg.Type == util.InfoTypeSuccess {
			h.hasChanges = false
		}
	}

	return h, nil
}

func (h HomePage) View() string {
	t := lipgloss.NewStyle()
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Quit Dialog Demo")

	var b strings.Builder
	b.WriteString(title + "\n\n")

	status := "No unsaved changes"
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	if h.hasChanges {
		status = "⚠️  You have unsaved changes!"
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	}

	b.WriteString(statusStyle.Render(status))
	b.WriteString("\n\n")
	b.WriteString("Commands:\n")
	b.WriteString("  s    - Simulate making changes\n")
	b.WriteString("  q    - Quit (with confirmation if changes)\n")
	b.WriteString("  ctrl+c - Force quit\n")

	return b.String()
}
