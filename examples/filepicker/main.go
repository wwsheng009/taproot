package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/tui/app"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs/filepicker"
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
	selectedFile string
}

func NewHomePage() *HomePage {
	return &HomePage{}
}

func (h *HomePage) Init() tea.Cmd {
	return nil
}

func (h *HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+o":
			// Open file picker
			fp := filepicker.New(".", func(path string) tea.Cmd {
				h.selectedFile = path
				return func() tea.Msg {
					return util.ReportInfo(fmt.Sprintf("Selected: %s", path))
				}
			})
			return h, func() tea.Msg {
				return dialogs.OpenDialogMsg{Model: fp}
			}
		case "ctrl+c", "q":
			return h, tea.Quit
		}
	}

	return h, nil
}

func (h *HomePage) View() string {
	t := lipgloss.NewStyle()
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot FilePicker Demo")

	var b strings.Builder
	b.WriteString(title + "\n\n")
	b.WriteString("Press ctrl+o to open the file picker\n\n")
	
	if h.selectedFile != "" {
		b.WriteString(fmt.Sprintf("Selected File: %s\n\n", h.selectedFile))
	} else {
		b.WriteString("No file selected\n\n")
	}

	b.WriteString("Press q or ctrl+c to quit")

	return b.String()
}
