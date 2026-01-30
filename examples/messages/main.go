package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/tui/components/messages"
	"github.com/wwsheng009/taproot/internal/tui/util"
)

func main() {
	model := NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

type Model struct {
	messages *messages.MessagesModel
	quitting bool
}

func NewModel() Model {
	msgs := messages.New()
	msgs.SetWidth(80)
	msgs.SetHeight(20)

	// Add some sample messages
	msgs.AddMessage(messages.Message{
		Role:    "system",
		Content: "Welcome to Taproot Messages Demo! This is a system message.",
	})

	msgs.AddMessage(messages.Message{
		Role:    "user",
		Content: "Hello! Can you help me understand how the messages component works?",
	})

	msgs.AddMessage(messages.Message{
		Role: "assistant",
		Content: "Of course! The Messages component displays a list of chat messages with:\n\n" +
			"• Different roles (user, assistant, system, tool)\n" +
			"• Automatic scrolling support\n" +
			"• Color-coded headers\n" +
			"• Tool call display\n" +
			"• Responsive layout",
	})

	msgs.AddMessage(messages.Message{
		Role:    "user",
		Content: "Can you show me a tool call example?",
	})

	msgs.AddMessage(messages.Message{
		Role: "assistant",
		Content: "Sure! Here's an example of a tool call:",
		ToolUse: &messages.ToolUse{
			Name:      "get_weather",
			Arguments: `{"location": "San Francisco", "unit": "celsius"}`,
			Result:    `{"temperature": 18, "condition": "partly cloudy", "humidity": 65}`,
		},
	})

	return Model{
		messages: msgs,
		quitting: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			m.messages.Update(msg)
		case "down", "j":
			m.messages.Update(msg)
		case "g", "home":
			m.messages.Update(msg)
		case "G", "end":
			m.messages.Update(msg)
		case "a":
			// Add a new user message
			m.messages.AddMessage(messages.Message{
				Role:    "user",
				Content: "This is a new message added at runtime!",
			})
			return m, nil
		case "c":
			// Clear all messages
			m.messages.Clear()
			return m, util.ReportInfo("Messages cleared")
		}
	case tea.WindowSizeMsg:
		m.messages.SetWidth(msg.Width)
		m.messages.SetHeight(msg.Height - 4) // Leave room for header/footer
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	theme := lipgloss.NewStyle()
	title := theme.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Messages Component Demo")

	var b strings.Builder

	// Header
	b.WriteString(title + "\n\n")

	// Messages view
	b.WriteString(m.messages.View())

	// Footer hints
	b.WriteString("\n")
	hints := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(
		"↑↓: Scroll | a: Add message | c: Clear | g/G: Home/End | q: Quit",
	)
	b.WriteString(hints)

	return b.String()
}
