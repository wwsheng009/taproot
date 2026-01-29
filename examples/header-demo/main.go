package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourorg/taproot/internal/ui/components/header"
)

type model struct {
	header        *header.HeaderComponent
	contentHeight int
	errorCount    int
	workingDir    string
	tokenUsed     int
	tokenMax      int
	cost          float64
	detailsOpen   bool
	compactMode   bool
	brand         string
	title         string
}

func initialModel() model {
	h := header.New()
	h.SetSize(100, 2)
	h.SetWorkingDirectory("/projects/ai/Taproot")
	h.SetTokenUsage(64000, 128000, 1.50)
	h.SetErrorCount(3)

	return model{
		header:       h,
		errorCount:   3,
		workingDir:   "/projects/ai/Taproot",
		tokenUsed:    64000,
		tokenMax:     128000,
		cost:         1.50,
		detailsOpen:  false,
		compactMode:  false,
		brand:        "Charm™",
		title:        "CRUSH",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "d":
			// Toggle details
			m.detailsOpen = !m.detailsOpen
			m.header.SetDetailsOpen(m.detailsOpen)
		case "e":
			// Add errors
			m.errorCount++
			m.header.SetErrorCount(m.errorCount)
		case "r":
			// Reset errors
			m.errorCount = 0
			m.header.SetErrorCount(m.errorCount)
		case "t":
			// Update token usage
			m.tokenUsed = 80000
			m.cost = 2.00
			m.header.SetTokenUsage(m.tokenUsed, m.tokenMax, m.cost)
		case "h":
			// Change working directory
			if m.workingDir == "/projects/ai/Taproot" {
				m.workingDir = "/home/user/projects/taproot/examples"
			} else {
				m.workingDir = "/projects/ai/Taproot"
			}
			m.header.SetWorkingDirectory(m.workingDir)
		case "c":
			// Change compact mode
			m.compactMode = true
			m.header.SetCompactMode(m.compactMode)
		case "n":
			// Normal mode
			m.compactMode = false
			m.header.SetCompactMode(m.compactMode)
		case "b":
			// Change brand
			m.brand = "MyBrand™"
			m.title = "APP"
			m.header.SetBrand(m.brand, m.title)
		case "s":
			// Reset to default brand
			m.brand = "Charm™"
			m.title = "CRUSH"
			m.header.SetBrand(m.brand, m.title)
		}

	case tea.WindowSizeMsg:
		// Update header size
		m.header.SetSize(msg.Width, 2)
		m.contentHeight = msg.Height - 2
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Render header
	b.WriteString(m.header.View())
	b.WriteString("\n")

	// Render content area
	content := strings.Builder{}
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("Header Component Demo\n\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("Press keys to interact:\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  q / ctrl+c  - Quit\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  d           - Toggle details\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  e           - Add error\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  r           - Reset errors\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  t           - Update token usage\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  h           - Change working dir\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  c           - Compact mode\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  n           - Normal mode\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  b           - Change brand\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  s           - Reset brand\n\n")

	// Show current state
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("Current State:\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  Details Open: ")
	if m.detailsOpen {
		content.WriteString("Yes\n")
	} else {
		content.WriteString("No\n")
	}
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  Error Count: ")
	content.WriteString(fmt.Sprintf("%d\n", m.errorCount))
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  Working Dir: ")
	content.WriteString(m.workingDir)
	content.WriteString("\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  Token Usage: ")
	content.WriteString(fmt.Sprintf("%d/%d\n", m.tokenUsed, m.tokenMax))
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  Brand: ")
	content.WriteString(m.brand)
	content.WriteString(" ")
	content.WriteString(m.title)
	content.WriteString("\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  Compact Mode: ")
	if m.compactMode {
		content.WriteString("Yes\n")
	} else {
		content.WriteString("No\n")
	}

	// Pad remaining height
	lines := strings.Split(content.String(), "\n")
	for i := len(lines); i < m.contentHeight; i++ {
		b.WriteString("\n")
	}
	b.WriteString(content.String())

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
