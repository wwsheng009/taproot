package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/tui/app"
	"github.com/yourorg/taproot/internal/tui/components/dialogs"
	"github.com/yourorg/taproot/internal/tui/components/dialogs/commands"
	"github.com/yourorg/taproot/internal/tui/page"
	"github.com/yourorg/taproot/internal/tui/util"
)

const (
	pageHome page.PageID = "home"
)

// MyCommandProvider implements the CommandProvider interface
type MyCommandProvider struct {
	actionLog []string
}

func (p *MyCommandProvider) Commands() []commands.Command {
	return []commands.Command{
		{
			ID:          "save",
			Title:       "Save",
			Description: "Save current work",
			Callback: func() tea.Cmd {
				p.actionLog = append(p.actionLog, "Saved!")
				return func() tea.Msg {
					return util.ReportInfo("File saved successfully")
				}
			},
		},
		{
			ID:          "load",
			Title:       "Load",
			Description: "Load a file",
			Callback: func() tea.Cmd {
				p.actionLog = append(p.actionLog, "Opened file dialog")
				return func() tea.Msg {
					return util.ReportInfo("File loaded")
				}
			},
		},
		{
			ID:          "settings",
			Title:       "Settings",
			Description: "Open settings",
			Callback: func() tea.Cmd {
				p.actionLog = append(p.actionLog, "Opened settings")
				return func() tea.Msg {
					return util.ReportInfo("Settings opened")
				}
			},
		},
		{
			ID:          "export",
			Title:       "Export",
			Description: "Export data",
			Callback: func() tea.Cmd {
				p.actionLog = append(p.actionLog, "Exported data")
				return util.ReportInfo("Data exported")
			},
		},
		{
			ID:          "help",
			Title:       "Help",
			Description: "Show help information",
			Callback: func() tea.Cmd {
				p.actionLog = append(p.actionLog, "Opened help")
				return util.ReportInfo("Help: Press ctrl+p for commands")
			},
		},
		{
			ID:          "quit",
			Title:       "Quit",
			Description: "Exit the application",
			Callback: func() tea.Cmd {
				return tea.Quit
			},
		},
	}
}

func main() {
	application := app.NewApp()
	provider := &MyCommandProvider{actionLog: []string{}}

	homePage := NewHomePage(provider)
	application.RegisterPage(pageHome, homePage)
	application.SetPage(pageHome)

	p := tea.NewProgram(application, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}

	// Show action log
	fmt.Println("\n=== Action Log ===")
	for _, action := range provider.actionLog {
		fmt.Println("-", action)
	}
}

// HomePage is the home page
type HomePage struct {
	provider *MyCommandProvider
}

func NewHomePage(provider *MyCommandProvider) HomePage {
	return HomePage{provider: provider}
}

func (h HomePage) Init() tea.Cmd {
	return nil
}

func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+p":
			// Open command palette
			dialog := commands.NewCommandsDialog(h.provider)
			return h, func() tea.Msg {
				return dialogs.OpenDialogMsg{Model: dialog}
			}
		case "ctrl+c", "q":
			return h, tea.Quit
		}
	}

	return h, nil
}

func (h HomePage) View() string {
	t := lipgloss.NewStyle()
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Commands Palette Demo")

	var b strings.Builder
	b.WriteString(title + "\n\n")
	b.WriteString("Press ctrl+p to open the command palette\n\n")
	b.WriteString("Available commands:\n")
	b.WriteString("  - Save: Save current work\n")
	b.WriteString("  - Load: Load a file\n")
	b.WriteString("  - Settings: Open settings\n")
	b.WriteString("  - Export: Export data\n")
	b.WriteString("  - Help: Show help\n")
	b.WriteString("  - Quit: Exit the application\n\n")
	b.WriteString("Press q or ctrl+c to quit")

	return b.String()
}
