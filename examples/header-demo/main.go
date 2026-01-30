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

func initialModel() *model {
	h := header.New()
	h.SetSize(100, 2)
	h.SetWorkingDirectory("/projects/ai/Taproot")
	h.SetTokenUsage(0, 128000, 0.00)
	h.SetErrorCount(3)

	return &model{
		header:       h,
		errorCount:   3,
		workingDir:   "/projects/ai/Taproot",
		tokenUsed:    0,
		tokenMax:     128000,
		cost:         0.00,
		detailsOpen:  false,
		compactMode:  false,
		brand:        "Charm™",
		title:        "CRUSH",
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q":
			return m, tea.Quit
		case "t", "T":
			// Update token usage (cycle 0 -> 100% -> 0)
			m.tokenUsed += 10000
			if m.tokenUsed > m.tokenMax {
				m.tokenUsed = 0
			}
			m.cost = float64(m.tokenUsed) / float64(m.tokenMax) * 3.00
			m.header.SetTokenUsage(m.tokenUsed, m.tokenMax, m.cost)
		case "d", "D":
			// Toggle details
			m.detailsOpen = !m.detailsOpen
			m.header.SetDetailsOpen(m.detailsOpen)
		case "e", "E":
			// Add errors
			m.errorCount++
			m.header.SetErrorCount(m.errorCount)
		case "r", "R":
			// Reset errors
			m.errorCount = 0
			m.header.SetErrorCount(m.errorCount)
		case "h", "H":
			// Change working directory
			if m.workingDir == "/projects/ai/Taproot" {
				m.workingDir = "/home/user/projects/taproot/examples"
			} else {
				m.workingDir = "/projects/ai/Taproot"
			}
			m.header.SetWorkingDirectory(m.workingDir)
		case "c", "C":
			// Change compact mode
			m.compactMode = true
			m.header.SetCompactMode(m.compactMode)
		case "n", "N":
			// Normal mode
			m.compactMode = false
			m.header.SetCompactMode(m.compactMode)
		case "b", "B":
			// Change brand
			m.brand = "MyBrand™"
			m.title = "APP"
			m.header.SetBrand(m.brand, m.title)
		case "s", "S":
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

func (m *model) View() string {
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
	content.WriteString("  Q / Ctrl+C  - Quit\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  T           - Cycle token % (0->100->0)\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  D           - Toggle details\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  E           - Add error\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  R           - Reset errors\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  H           - Change working dir\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  C           - Compact mode\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  N           - Normal mode\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  B           - Change brand\n")
	content.WriteString(strings.Repeat(" ", 2))
	content.WriteString("  S           - Reset brand\n")

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
