package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Model holds the application state
type model struct {
	items   []string
	cursor  int
	selected map[int]struct{}
}

// InitialModel creates the initial model
func initialModel() model {
	return model{
		items:   []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"},
		selected: make(map[int]struct{}),
	}
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ", "enter":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	var b strings.Builder

	b.WriteString("Taproot TUI Framework - List Example\n")
	b.WriteString("Use arrow keys or j/k to move, space to select, q to quit\n\n")

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, item))
	}

	b.WriteString("\n")
	if len(m.selected) > 0 {
		b.WriteString(fmt.Sprintf("Selected: %d items", len(m.selected)))
	} else {
		b.WriteString("No items selected")
	}

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
