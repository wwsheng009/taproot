package main

import (
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/completions"
	"github.com/wwsheng009/taproot/ui/render"
)

// AutoCompleteModel demonstrates auto-completion with different providers
type AutoCompleteModel struct {
	providerType    string
	currentQuery    string
	lastSelection   string
	showHelp        bool
	autocomplete    *completions.AutoCompletion
	stringProvider  *completions.StringProvider
	commandProvider *completions.CommandProvider
}

// NewAutoCompleteModel creates a new auto-complete model
func NewAutoCompleteModel() *AutoCompleteModel {
	// Create string provider with fruits
	fruits := []string{
		"Apple", "Apricot", "Avocado",
		"Banana", "Blackberry", "Blueberry",
		"Cherry", "Cantaloupe", "Coconut",
		"Date", "Dragonfruit", "Durian",
		"Elderberry", "Fig", "Grape",
		"Grapefruit", "Guava", "Honeydew",
		"Kiwi", "Kumquat", "Lemon",
		"Lime", "Mango", "Melon",
		"Nectarine", "Orange", "Papaya",
		"Peach", "Pear", "Persimmon",
		"Pineapple", "Plum", "Pomegranate",
		"Raspberry", "Strawberry", "Watermelon",
	}
	stringProvider := completions.NewStringProviderFromStrings(fruits)

	// Create command provider with commands
	commands := []completions.CommandItem{
		{Name: "help", Description: "Show help", Handler: func(args ...string) any { return "Showing help" }},
		{Name: "quit", Description: "Exit application", Handler: func(args ...string) any { return "Quitting" }},
		{Name: "clear", Description: "Clear screen", Handler: func(args ...string) any { return "Cleared" }},
		{Name: "status", Description: "Show status", Handler: func(args ...string) any { return "Status: OK" }},
		{Name: "ls", Description: "List files", Handler: func(args ...string) any { return "Files listed" }},
		{Name: "cd", Description: "Change directory", Handler: func(args ...string) any { return "Directory changed" }},
	}
	commandProvider := completions.NewCommandProvider(commands)

	// Create auto-complete with string provider
	autocomplete := completions.NewAutoCompletion(stringProvider, 3, 10, 60)

	return &AutoCompleteModel{
		providerType:    "fruit",
		currentQuery:    "",
		lastSelection:   "",
		showHelp:        true,
		autocomplete:    autocomplete,
		stringProvider:  stringProvider,
		commandProvider: commandProvider,
	}
}

// Init initializes the model
func (m *AutoCompleteModel) Init() render.Cmd {
	return nil
}

// Update handles incoming messages
func (m *AutoCompleteModel) Update(msg any) (render.Model, render.Cmd) {
	if msg == nil {
		return m, render.None()
	}

	switch msg := msg.(type) {
	case render.KeyMsg:
		switch msg.Key {
		case "q", "ctrl+c":
			// Quit
			return m, render.Command(func() error {
				return nil
			})
		case "?":
			m.showHelp = !m.showHelp
		case "1":
			// Switch to fruit provider
			m.providerType = "fruit"
			m.autocomplete.Close()
			m.autocomplete = completions.NewAutoCompletion(m.stringProvider, 3, 10, 60)
			m.currentQuery = ""
		case "2":
			// Switch to command provider
			m.providerType = "command"
			m.autocomplete.Close()
			m.autocomplete = completions.NewAutoCompletion(m.commandProvider, 3, 10, 60)
			m.currentQuery = ""
		case "/":
			// Toggle autocomplete
			if m.autocomplete.IsOpen() {
				m.autocomplete.Close()
			} else {
				m.autocomplete.Open()
				m.autocomplete.SetQuery(m.currentQuery)
			}
		case "enter":
			// Select current
			if m.autocomplete.IsOpen() {
				selected := m.autocomplete.Selected()
				if selected != nil {
					m.lastSelection = selected.Display()
					m.autocomplete.Close()
				}
			}
		case "esc":
			// Close autocomplete
			m.autocomplete.Close()
		case "backspace":
			if len(m.currentQuery) > 0 {
				m.currentQuery = m.currentQuery[:len(m.currentQuery)-1]
				if m.autocomplete.IsOpen() {
					m.autocomplete.SetQuery(m.currentQuery)
				}
			}
		default:
			// Typing - add to query
			if len(msg.Key) == 1 && msg.Key[0] >= 32 && msg.Key[0] < 127 {
				m.currentQuery += msg.Key
				if !m.autocomplete.IsOpen() {
					m.autocomplete.Open()
				}
				m.autocomplete.SetQuery(m.currentQuery)
			}
		}
	}

	return m, render.None()
}

// View returns the string representation for rendering
func (m *AutoCompleteModel) View() string {
	var b strings.Builder

	b.WriteString("Taproot v2.0.0 - Auto-Completion Demo\n")
	b.WriteString("======================================\n\n")

	if m.showHelp {
		b.WriteString("Keys:\n")
		b.WriteString("  1          : Switch to fruit completions\n")
		b.WriteString("  2          : Switch to command completions\n")
		b.WriteString("  /          : Toggle auto-complete popup\n")
		b.WriteString("  Type       : Filter completions\n")
		b.WriteString("  ↑/↓ or j/k : Navigate results\n")
		b.WriteString("  Enter      : Select current\n")
		b.WriteString("  Esc        : Close popup\n")
		b.WriteString("  Backspace  : Delete last char\n")
		b.WriteString("  ?          : Toggle this help\n")
		b.WriteString("  q/ctrl+c   : Quit\n\n")
	}

	b.WriteString("──────────────────────────────────────\n\n")

	// Show provider info
	b.WriteString(fmt.Sprintf("Provider: %s\n\n", m.providerType))

	// Show input line
	b.WriteString("Search: ")
	if len(m.currentQuery) == 0 {
		b.WriteString("█")
	} else {
		b.WriteString(m.currentQuery + "█")
	}
	b.WriteString("\n\n")

	// Show autocomplete popup
	if m.autocomplete.IsOpen() && m.autocomplete.HasItems() {
		items := m.autocomplete.Items()
		cursor := m.autocomplete.Cursor()
		start, end := m.autocomplete.VisibleRange()
		width, _ := m.autocomplete.Size()

		// Render popup box
		b.WriteString("┌─ Completions ")
		for i := 0; i < width-15; i++ {
			b.WriteString("─")
		}
		b.WriteString("┐\n")

		for i := start; i < end && i < len(items); i++ {
			item := items[i]
			prefix := " "
			if i == cursor {
				prefix = ">"
			}

			// Highlight matches
			display := item.Display()
			if len(m.currentQuery) > 0 {
				// Simple bold effect for matches
				display = strings.ReplaceAll(display, m.currentQuery, fmt.Sprintf("[%s]", m.currentQuery))
			}

			// Truncate if too long
			if len(display) > width-4 {
				display = display[:width-7] + "..."
			}

			b.WriteString(fmt.Sprintf("│ %s %-*s │\n", prefix, width-5, display))
		}

		b.WriteString("└")
		for i := 0; i < width-2; i++ {
			b.WriteString("─")
		}
		b.WriteString("┘\n")

		// Show stats
		total := m.autocomplete.ItemCount()
		b.WriteString(fmt.Sprintf("Showing %d of %d results\n\n", end-start, total))
	}

	// Show last selection
	if m.lastSelection != "" {
		b.WriteString(fmt.Sprintf("Last selected: %s\n", m.lastSelection))
	}

	return b.String()
}

func main() {
	model := NewAutoCompleteModel()

	config := render.DefaultConfig()
	config.EnableAltScreen = true

	engine, err := render.CreateEngine(render.EngineBubbletea, config)
	if err != nil {
		fmt.Printf("Failed to create engine: %v\n", err)
		return
	}

	fmt.Println("Starting auto-completion demo...")
	fmt.Println("Press ? for help")

	if err := engine.Start(model); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
