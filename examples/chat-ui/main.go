package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/components/messages"
	"github.com/wwsheng009/taproot/ui/render"
)

// Model - Chat UI model
type Model struct {
	messages        []messages.MessageItem
	selectedMessage int // Index of currently selected message (-1 = none)
	input           textinput.Model
	width           int
	height          int
	quitting        bool
}

// NewModel - Create new chat UI model
func NewModel() Model {
	// Add initial welcome message
	welcomeMsg := messages.NewAssistantMessage(
		"msg-0",
		"# Welcome to Taproot Chat UI!\n\nThis is an interactive chat interface demonstrating the **v2.0 messages** component.\n\n## Demo Commands\n\n- `help` - Show help message\n- `demo` - Show message types\n- `clear` - Clear all messages\n\nType a message and press Enter to chat!",
	)
	welcomeMsg.SetShowTimestamp(false)

	// Initialize and focus text input
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()

	return Model{
		messages: []messages.MessageItem{welcomeMsg},
		input:    ti,
		width:    80,
		height:   24,
		quitting: false,
	}
}

// Init - Initialize model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update - Update model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle input
		m.input, cmd = m.input.Update(msg)

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyUp:
			// Navigate to previous message
			if len(m.messages) > 0 {
				if m.selectedMessage <= 0 {
					m.selectedMessage = len(m.messages) - 1
				} else {
					m.selectedMessage--
				}
			}

		case tea.KeyDown:
			// Navigate to next message
			if len(m.messages) > 0 {
				if m.selectedMessage >= len(m.messages)-1 || m.selectedMessage < 0 {
					m.selectedMessage = 0
				} else {
					m.selectedMessage++
				}
			}

		case tea.KeyEnter:
			// Check if text input has content
			inputText := strings.TrimSpace(m.input.Value())
			if inputText != "" {
				// Send the message
				m = m.handleInput(inputText)
				m.input.Reset()
				m.selectedMessage = len(m.messages) - 1 // Select the new message
				return m, cmd
			}
			// If input is empty, toggle selected message
			if m.selectedMessage >= 0 && m.selectedMessage < len(m.messages) {
				if exp, ok := m.messages[m.selectedMessage].(interface{ ToggleExpanded() }); ok {
					exp.ToggleExpanded()
					// Save the updated message back to the slice
					m.messages[m.selectedMessage] = exp.(messages.MessageItem)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update all messages to handle their own state
	for i := range m.messages {
		var updatedMsg render.Model
		updatedMsg, _ = m.messages[i].Update(msg)
		m.messages[i] = updatedMsg.(messages.MessageItem)
	}

	return m, cmd
}

// handleInput - Process user input
func (m Model) handleInput(input string) Model {
	// Handle special commands
	switch strings.ToLower(input) {
	case "help":
		return m.showHelp()
	case "clear":
		m.messages = []messages.MessageItem{}
		clearMsg := messages.NewAssistantMessage(
			"msg-cleared",
			"Chat cleared! 🧹",
		)
		clearMsg.SetShowTimestamp(false)
		m.messages = append(m.messages, clearMsg)
		return m
	case "demo":
		return m.showDemo()
	}

	// Create user message
	userMsg := messages.NewUserMessage(
		fmt.Sprintf("msg-%d-user", len(m.messages)),
		input,
	)
	userMsg.SetShowTimestamp(false)
	m.messages = append(m.messages, userMsg)

	// Generate AI response
	response := m.generateResponse(input)
	assistantMsg := messages.NewAssistantMessage(
		fmt.Sprintf("msg-%d", len(m.messages)),
		response,
	)
	assistantMsg.SetShowTimestamp(false)
	m.messages = append(m.messages, assistantMsg)

	return m
}

// showHelp - Show help message
func (m Model) showHelp() Model {
	helpMsg := messages.NewAssistantMessage(
		"msg-help",
		"# Help\n\n"+
			"## Commands\n\n"+
			"- `help` - Show this help\n"+
			"- `demo` - Show message types demo\n"+
			"- `clear` - Clear messages\n\n"+
			"## Message Types\n\n"+
			"1. **UserMessage** - Messages from you\n"+
			"2. **AssistantMessage** - AI responses (Markdown)\n"+
			"3. **ToolMessage** - Tool execution\n"+
			"4. **FetchMessage** - Web searches\n"+
			"5. **DiagnosticMessage** - LSP errors\n"+
			"6. **TodoMessage** - Task lists\n\n"+
			"## Features\n\n"+
			"- Markdown support with syntax highlighting\n"+
			"- Multiple message types\n"+
			"- Responsive layout\n"+
			"- Clean, modern UI",
	)
	m.messages = append(m.messages, helpMsg)
	return m
}

// showDemo - Show demo messages
func (m Model) showDemo() Model {
	// User Message
	userMsg := messages.NewUserMessage(
		"msg-demo-user",
		"Here's my code:\n\n```go\nfunc hello() {\n    println(\"Hello!\")\n}\n```",
	)
	m.messages = append(m.messages, userMsg)

	// Assistant Message
	assistantMsg := messages.NewAssistantMessage(
		"msg-demo-assistant",
		"Great code! Here's some **bold** and *italic* text.\n\n## Lists\n\n- Item 1\n- Item 2\n\n## Code\n\n```python\nprint('Python too!')\n```",
	)
	m.messages = append(m.messages, assistantMsg)

	// Todo Message
	todoMsg := messages.NewTodoMessage(
		"msg-demo-todo",
		"My Tasks",
	)
	todoMsg.AddTodo(messages.Todo{
		ID:          "task-1",
		Description: "Learn Taproot",
		Status:      messages.TodoStatusCompleted,
	})
	todoMsg.AddTodo(messages.Todo{
		ID:          "task-2",
		Description: "Build an app",
		Status:      messages.TodoStatusInProgress,
		Progress:    0.5,
	})
	todoMsg.AddTodo(messages.Todo{
		ID:          "task-3",
		Description: "Ship to production",
		Status:      messages.TodoStatusPending,
	})
	m.messages = append(m.messages, todoMsg)

	// Diagnostic Message
	diagMsg := messages.NewDiagnosticMessage(
		"msg-demo-diag",
		"compiler",
	)
	m.messages = append(m.messages, diagMsg)

	// Fetch Message
	fetchMsg := messages.NewFetchMessage(
		"msg-demo-fetch",
		messages.FetchTypeWebSearch,
	)
	m.messages = append(m.messages, fetchMsg)

	// Explanation
	explainMsg := messages.NewAssistantMessage(
		"msg-demo-explain",
		"That's all 6 message types! Each serves a different purpose in a chat/assistant interface.",
	)
	m.messages = append(m.messages, explainMsg)

	return m
}

// generateResponse - Generate AI response (simulated)
func (m Model) generateResponse(input string) string {
	inputLower := strings.ToLower(input)

	switch {
	case strings.Contains(inputLower, "hello") || strings.Contains(inputLower, "hi"):
		return "Hi there! 👋 Welcome to the Taproot Chat UI demo!\n\n" +
			"This demo showcases the v2.0 messages component. Try typing `demo` to see all message types!"

	case strings.Contains(inputLower, "demo") || strings.Contains(inputLower, "example"):
		return "Sure! Type `demo` to see all message types in action.\n\n" +
			"You'll see examples of:\n" +
			"- User and assistant messages\n" +
			"- Todo lists with progress\n" +
			"- Diagnostic messages\n" +
			"- Fetch messages"

	case strings.Contains(inputLower, "message") || strings.Contains(inputLower, "type"):
		return "The messages component supports 6 types:\n\n" +
			"1. **UserMessage** - Your messages\n" +
			"2. **AssistantMessage** - AI responses\n" +
			"3. **ToolMessage** - Tool calls\n" +
			"4. **FetchMessage** - Web searches\n" +
			"5. **DiagnosticMessage** - Errors\n" +
			"6. **TodoMessage** - Tasks\n\n" +
			"Type `demo` to see them all!"

	case strings.Contains(inputLower, "markdown") || strings.Contains(inputLower, "code"):
		return "The assistant messages support full **Markdown** formatting:\n\n" +
			"## Headers\n\n" +
			"**Bold** and *italic* text\n\n" +
			"- Lists\n- Items\n\n" +
			"```go\n// Code blocks with syntax highlighting\nfunc hello() {\n    println(\"Hello!\")\n}\n```\n\n" +
			"Plus links, quotes, tables, and more!"

	case strings.Contains(inputLower, "todo") || strings.Contains(inputLower, "task"):
		return "Todo messages help track tasks:\n\n" +
			"- **Pending**: Not started\n" +
			"- **In Progress**: Working on it\n" +
			"- **Completed**: Done!\n\n" +
			"Each todo can have:\n" +
			"- Description\n" +
			"- Status\n" +
			"- Progress (0.0 to 1.0)\n\n" +
			"Type `demo` to see an example!"

	default:
		return fmt.Sprintf("Thanks for writing: %q\n\n", input) +
			"This is a demo interface. Type `help` to see commands or `demo` to explore message types!"
	}
}

// View - Render chat UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Header
	header := styleHeader.Width(m.width).Render("💬 Taproot Chat UI")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Messages area
	messagesHeight := m.height - 5 // header + input + footer
	messagesBox := styleMessagesBox.
		Height(messagesHeight).
		Width(m.width).
		Render(m.renderMessages())

	b.WriteString(messagesBox)
	b.WriteString("\n")

	// Input area
	inputBox := styleInputBox.Width(m.width).Render(
		fmt.Sprintf("%s %s",
			stylePrompt.Render(">"),
			m.input.View(),
		),
	)

	b.WriteString(inputBox)
	b.WriteString("\n")

	// Footer
	footer := styleFooter.Width(m.width).Render(
		"[Enter] Send  [Ctrl+C] Quit  Type `help` for commands",
	)
	b.WriteString(footer)

	return b.String()
}

// renderMessages - Render messages with chat-style layout
func (m Model) renderMessages() string {
	if len(m.messages) == 0 {
		return styleEmpty.Render("No messages yet. Type a message to start!")
	}

	var b strings.Builder
	availableWidth := m.width - 6 // Account for padding and borders

	// Show messages with chat-style layout
	for i, msg := range m.messages {
		// Determine message type and styling
		var messageStyle lipgloss.Style
		var isSelected bool
		var isUserMessage bool

		// Determine if this message is selected
		if i == m.selectedMessage {
			isSelected = true
		}

		// Get message content
		messageView := msg.View()

		// Determine message type and apply appropriate style
		switch m := msg.(type) {
		case *messages.UserMessage:
			// User message - right aligned
			isUserMessage = true
			if isSelected {
				messageStyle = styleUserMessageSelected
			} else {
				messageStyle = styleUserMessage
			}
			_ = m // Avoid unused variable warning

		case *messages.AssistantMessage:
			// Assistant message - left aligned
			isUserMessage = false
			if isSelected {
				messageStyle = styleAssistantMessageSelected
			} else {
				messageStyle = styleAssistantMessage
			}
			_ = m // Avoid unused variable warning

		default:
			// Other message types - left aligned
			isUserMessage = false
			if isSelected {
				messageStyle = styleOtherMessageSelected
			} else {
				messageStyle = styleOtherMessage
			}
		}

		// Apply width constraint before rendering
		// Reserve space for margins and selection indicator
		// Use 70% of available width for messages (leaves 30% for spacing)
		maxMsgWidth := int(float64(availableWidth) * 0.7)
		if maxMsgWidth < 30 {
			maxMsgWidth = 30 // Minimum width for readability
		}
		messageStyle = messageStyle.Width(maxMsgWidth)

		// Apply style and render message
		styledMessage := messageStyle.Render(messageView)
		msgWidth := lipgloss.Width(styledMessage)

		// Position the message based on type
		if isUserMessage {
			// User message - right aligned
			var indicator string
			if isSelected {
				indicator = "▶"
			} else {
				indicator = " "
			}

			// Calculate left padding to push message to the right
			leftPadding := availableWidth - msgWidth
			if leftPadding < 1 {
				leftPadding = 1 // Minimum padding to separate from indicator
			}

			// Split and render with indicator at position 0
			msgLines := strings.Split(styledMessage, "\n")
			for lineIdx, line := range msgLines {
				if lineIdx > 0 {
					b.WriteString("\n")
					b.WriteString(" ") // Use space for subsequent lines to maintain alignment
					paddingStr := strings.Repeat(" ", leftPadding)
					b.WriteString(paddingStr + line)
				} else {
					b.WriteString(indicator)
					paddingStr := strings.Repeat(" ", leftPadding)
					b.WriteString(paddingStr + line)
				}
			}
		} else {
			// Assistant/other message - left aligned
			var indicator string
			if isSelected {
				indicator = "▶"
			} else {
				indicator = " "
			}

			// Split and render with indicator at position 0
			msgLines := strings.Split(styledMessage, "\n")
			for lineIdx, line := range msgLines {
				if lineIdx > 0 {
					b.WriteString("\n")
					b.WriteString(" ") // Space on subsequent lines to maintain alignment
					b.WriteString(line)
				} else {
					b.WriteString(indicator)
					b.WriteString(line)
				}
			}
		}

		// Add spacing between messages
		if i < len(m.messages)-1 {
			b.WriteString("\n\n")
		}
	}

	return b.String()
}

// Styles
var (
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("#1e1e2e")).
			Foreground(lipgloss.Color("#cba6f7"))

	styleMessagesBox = lipgloss.NewStyle().
				Padding(1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#45475a"))

	styleInputBox = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("#313244")).
			Foreground(lipgloss.Color("#cdd6f4"))

	stylePrompt = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#89b4fa"))

	styleFooter = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("#1e1e2e")).
			Foreground(lipgloss.Color("#6c7086"))

	styleEmpty = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086"))

	styleSelectedIndicator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cba6f7")).
				Bold(true)

	// User message styles (right-aligned with double border)
	styleUserMessage = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#89b4fa")).
			Foreground(lipgloss.Color("#cdd6f4")).
			Padding(0, 1)

	styleUserMessageSelected = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#b4befe")).
			Foreground(lipgloss.Color("#cdd6f4")).
			Padding(0, 1)

	// Assistant message styles (left-aligned with border)
	styleAssistantMessage = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#45475a")).
			Foreground(lipgloss.Color("#cdd6f4")).
			Padding(0, 1)

	styleAssistantMessageSelected = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#585b70")).
			Foreground(lipgloss.Color("#cdd6f4")).
			Padding(0, 1)

	// Other message types (left-aligned with neutral border)
	styleOtherMessage = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#313244")).
			Foreground(lipgloss.Color("#a6adc8")).
			Padding(0, 1).
			MaxWidth(60)

	styleOtherMessageSelected = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#45475a")).
			Foreground(lipgloss.Color("#cdd6f4")).
			Padding(0, 1).
			MaxWidth(60)

)

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
