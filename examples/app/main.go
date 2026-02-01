package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/app"
	"github.com/wwsheng009/taproot/tui/components/dialogs"
	"github.com/wwsheng009/taproot/tui/page"
	"github.com/wwsheng009/taproot/tui/util"
)

const (
	pageHome page.PageID = "home"
	pageMenu page.PageID = "page"
	pageAbout page.PageID = "about"
)

func main() {
	application := app.NewApp()

	// Register pages
	homePage := newHomePage()
	menuPage := newMenuPage()
	aboutPage := newAboutPage()

	application.RegisterPage(pageHome, homePage)
	application.RegisterPage(pageMenu, menuPage)
	application.RegisterPage(pageAbout, aboutPage)

	// Set initial page
	if newApp := application.SetPage(pageHome); newApp != nil {
		application = *newApp
	}

	// Run the application
	p := tea.NewProgram(application, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

// HomePage is the home page
type HomePage struct {
	count int
}

func newHomePage() HomePage {
	return HomePage{}
}

func (h HomePage) Init() tea.Cmd {
	return nil
}

func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			return h, func() tea.Msg {
				return page.PageChangeMsg{ID: pageMenu}
			}
		case "2":
			return h, func() tea.Msg {
				return page.PageChangeMsg{ID: pageAbout}
			}
		case "ctrl+d":
			return h, func() tea.Msg {
				return dialogs.OpenDialogMsg{Model: newDemoDialog()}
			}
		case "up", "right", "+", "=":
			h.count++
		case "down", "left", "-", "_":
			if h.count > 0 {
				h.count--
			}
		}
	}
	return h, nil
}

func (h HomePage) View() string {
	t := lipgloss.NewStyle()
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot TUI Framework - Home")
	
	var b strings.Builder
	b.WriteString(title + "\n\n")
	b.WriteString("Press 1: Go to Menu Page\n")
	b.WriteString("Press 2: Go to About Page\n")
	b.WriteString("Press ctrl+d: Open Demo Dialog\n")
	b.WriteString("Press +/-: Change counter\n\n")
	b.WriteString(fmt.Sprintf("Count: %d\n", h.count))
	b.WriteString(strings.Repeat("█", h.count))

	return b.String()
}

// MenuPage is a menu page
type MenuPage struct {
	cursor int
	items  []string
}

func newMenuPage() MenuPage {
	return MenuPage{
		items: []string{"Option 1", "Option 2", "Option 3", "Back to Home"},
	}
}

func (m MenuPage) Init() tea.Cmd {
	return nil
}

func (m MenuPage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(m.items)-1 {
				return m, func() tea.Msg {
					return page.PageBackMsg{}
				}
			}
		}
	}
	return m, nil
}

func (m MenuPage) View() string {
	t := lipgloss.NewStyle()
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Menu Page")
	
	var b strings.Builder
	b.WriteString(title + "\n\n")

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
	}

	b.WriteString("\nPress j/k or arrow keys to navigate")
	b.WriteString("\nPress enter to select")

	return b.String()
}

// AboutPage is the about page
type AboutPage struct{}

func newAboutPage() AboutPage {
	return AboutPage{}
}

func (a AboutPage) Init() tea.Cmd {
	return nil
}

func (a AboutPage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b", "q", "esc":
			return a, func() tea.Msg {
				return page.PageBackMsg{}
			}
		}
	}
	return a, nil
}

func (a AboutPage) View() string {
	t := lipgloss.NewStyle()
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("About Taproot")
	
	var b strings.Builder
	b.WriteString(title + "\n\n")
	b.WriteString("Taproot is a TUI framework for Go.\n\n")
	b.WriteString("Based on Bubbletea Elm architecture.\n\n")
	b.WriteString("Press b, q, or esc to go back")

	return b.String()
}

// DemoDialog is a simple demo dialog
type DemoDialog struct {
	visible bool
}

func newDemoDialog() *DemoDialog {
	return &DemoDialog{visible: true}
}

func (d *DemoDialog) Init() tea.Cmd {
	return nil
}

func (d *DemoDialog) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		d.visible = false
		return d, func() tea.Msg {
			return dialogs.CloseDialogMsg{}
		}
	}
	return d, nil
}

func (d *DemoDialog) View() string {
	if !d.visible {
		return ""
	}
	
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		Background(lipgloss.Color("62")).
		Padding(1, 2).
		Width(40)
	
	content := "Demo Dialog\n\nPress any key to close"
	return style.Render(content)
}

func (d *DemoDialog) Position() (int, int) {
	return 10, 5
}

func (d *DemoDialog) ID() dialogs.DialogID {
	return "demo"
}
