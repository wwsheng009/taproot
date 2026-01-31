package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/components/completions"
)

// Model holds the application state
type model struct {
	width       int
	height      int
	input       string
	cursor      int
	completions completions.CompletionsCmp
}

// InitialModel creates the initial model
func initialModel() model {
	return model{
		input:       "",
		cursor:      0,
		completions: completions.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// If completions are open, forward keys to it first
		if m.completions.Open() {
			switch msg.String() {
			case "up", "down", "ctrl+p", "ctrl+n", "enter", "tab", "ctrl+y", "esc":
				newComps, cmd := m.completions.Update(msg)
				m.completions = newComps.(completions.CompletionsCmp)
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			// Trigger completions on enter
			return m, m.openCompletions()

		case "ctrl+space":
			// Toggle completions with ctrl+space
			if m.completions.Open() {
				return m, func() tea.Msg {
					return completions.CloseCompletionsMsg{}
				}
			}
			return m, m.openCompletions()

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				if m.cursor > len(m.input) {
					m.cursor = len(m.input)
				}
				// Auto-filter completions if open
				if m.completions.Open() {
					return m, m.filterCompletions()
				}
			}

		case "ctrl+h": // Alternative to backspace
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				if m.cursor > len(m.input) {
					m.cursor = len(m.input)
				}
				if m.completions.Open() {
					return m, m.filterCompletions()
				}
			}

		default:
			// Regular character input
			if len(msg.String()) == 1 {
				m.input += msg.String()
				m.cursor++
				// Auto-filter completions if open
				if m.completions.Open() {
					return m, m.filterCompletions()
				}
			}
		}

	case completions.SelectCompletionMsg:
		// Handle selected completion
		if val, ok := msg.Value.(string); ok {
			m.input = val
			m.cursor = len(m.input)
		}
		m.completions = completions.New()
		return m, nil

	case completions.CompletionsClosedMsg:
		m.completions = completions.New()
		return m, nil
	}

	// Forward other messages to completions
	newComps, cmd := m.completions.Update(msg)
	m.completions = newComps.(completions.CompletionsCmp)
	return m, cmd
}

func (m model) openCompletions() tea.Cmd {
	items := getCompletions(m.input)
	return func() tea.Msg {
		return completions.OpenCompletionsMsg{
			Completions: items,
			X:           0,
			Y:           3,
			MaxResults:  10,
		}
	}
}

func (m model) filterCompletions() tea.Cmd {
	return func() tea.Msg {
		return completions.FilterCompletionsMsg{
			Query:  m.input,
			Reopen: true,
			X:      0,
			Y:      3,
		}
	}
}

func (m model) View() string {
	t := lipgloss.NewStyle()

	var b strings.Builder

	// Title
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Completions Demo")
	b.WriteString(title + "\n\n")

	// Instructions
	b.WriteString("Type to filter commands, press Enter or Ctrl+Space to show completions\n\n")

	// Input box
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("240")).
		Padding(0, 1)

	inputLabel := "> "
	inputDisplay := inputLabel + m.input
	if m.cursor < len(m.input) {
		inputDisplay += "_"
	} else {
		inputDisplay += " "
	}

	b.WriteString(inputStyle.Width(min(60, m.width)).Render(inputDisplay))
	b.WriteString("\n\n")

	// Show completions if open
	if m.completions.Open() {
		b.WriteString(m.completions.View())
	} else {
		b.WriteString("Press Enter to see available commands\n")
	}

	// Footer
	b.WriteString("\n" + t.Foreground(lipgloss.Color("244")).Render("Press q or ctrl+c to quit"))

	return b.String()
}

// getCompletions returns available completion items based on input
func getCompletions(input string) []completions.CompletionItem {
	commands := []struct {
		id    string
		title string
		value string
	}{
		{"help", "help - Show help information", "help "},
		{"status", "status - Show current status", "status "},
		{"config", "config - Configuration commands", "config "},
		{"config-set", "config set <key> <value> - Set configuration", "config set "},
		{"config-get", "config get <key> - Get configuration", "config get "},
		{"build", "build - Build the project", "build "},
		{"build-verbose", "build --verbose - Build with verbose output", "build --verbose"},
		{"run", "run - Run the application", "run "},
		{"test", "test - Run tests", "test "},
		{"test-cover", "test --cover - Run tests with coverage", "test --cover"},
		{"clean", "clean - Clean build artifacts", "clean "},
		{"install", "install - Install dependencies", "install "},
		{"update", "update - Update dependencies", "update "},
	}

	items := make([]completions.CompletionItem, 0, len(commands))
	for _, cmd := range commands {
		items = append(items, completions.NewCompletionItem(cmd.id, cmd.title, cmd.value))
	}

	return items
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
