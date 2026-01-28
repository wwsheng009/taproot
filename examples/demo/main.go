package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type model struct {
	count int
}

// InitialModel returns the initial model
func initialModel() model {
	return model{count: 0}
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "right", "+", "=":
			m.count++
		case "down", "left", "-", "_":
			if m.count > 0 {
				m.count--
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	var b strings.Builder

	b.WriteString("Taproot TUI Framework - Basic Demo\n\n")
	b.WriteString("Press arrow keys or +/- to change counter\n")
	b.WriteString("Press q or ctrl+c to quit\n\n")
	b.WriteString(fmt.Sprintf("Count: %d\n", m.count))
	b.WriteString(strings.Repeat("█", m.count))

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
