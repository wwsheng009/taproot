package main

import (
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/internal/ui/components/messages"
)

func main() {
	// Create sample messages
	var msgItems []interface {
		Init() error
		View() string
	}

	// 1. User Message
	userMsg := messages.NewUserMessage("msg-1", "Can you help me understand the Taproot framework?")
	userMsg.AddAttachment(messages.Attachment{
		ID:   "att-1",
		Type: "file",
		Name: "README.md",
		Size: 2048,
	})
	msgItems = append(msgItems, userMsg)

	// 2. Assistant Message with thinking
	assistantMsg := messages.NewAssistantMessage("msg-2", "Of course! Taproot provides engine-agnostic UI components.")
	assistantMsg.SetThinking(`Let me think about how to explain this...

Taproot is built on top of several key concepts:
1. Engine-agnostic rendering (Bubbletea, Ultraviolet)
2. Reusable component interfaces
3. Clean separation of concerns

The framework allows you to write UI components once
and use them with different rendering engines.`)
	msgItems = append(msgItems, assistantMsg)

	// 3. Tool Message
	toolMsg := messages.NewToolMessage("msg-3", "read_file")
	toolMsg.AddCall(messages.ToolCall{
		ID:        "call-1",
		Name:      "read_file",
		Arguments: map[string]any{"path": "/tmp/file.txt"},
		Status:    messages.ToolStatusCompleted,
		Result:    "File content successfully read",
	})
	msgItems = append(msgItems, toolMsg)

	// 4. Diagnostic Message
	diagMsg := messages.NewDiagnosticMessage("msg-4", "Build Issues")
	diagMsg.AddDiagnostic(messages.Diagnostic{
		Severity:  messages.SeverityError,
		Message:   "undefined variable: taproot",
		File:      "main.go",
		Line:      42,
		Column:    10,
		Code:      "E123",
	})
	diagMsg.SetExpanded(true)
	msgItems = append(msgItems, diagMsg)

	// 5. Todo Message
	todoMsg := messages.NewTodoMessage("msg-5", "Project Tasks")
	todoMsg.AddTodo(messages.Todo{
		ID:          "todo-1",
		Description: "Implement UserMessage component",
		Status:      messages.TodoStatusCompleted,
		Progress:    1.0,
		Tags:        []string{"component", "v2.0"},
	})
	todoMsg.AddTodo(messages.Todo{
		ID:          "todo-2",
		Description: "Implement AssistantMessage component",
		Status:      messages.TodoStatusCompleted,
		Progress:    1.0,
		Tags:        []string{"component", "v2.0"},
	})
	todoMsg.AddTodo(messages.Todo{
		ID:          "todo-3",
		Description: "Implement ToolMessage component",
		Status:      messages.TodoStatusCompleted,
		Progress:    1.0,
		Tags:        []string{"component", "v2.0"},
	})
	todoMsg.AddTodo(messages.Todo{
		ID:          "todo-4",
		Description: "Write comprehensive tests",
		Status:      messages.TodoStatusInProgress,
		Progress:    0.75,
		Tags:        []string{"testing"},
	})
	todoMsg.AddTodo(messages.Todo{
		ID:          "todo-5",
		Description: "Create interactive demo",
		Status:      messages.TodoStatusPending,
		Progress:    0.0,
		Tags:        []string{"demo"},
	})
	msgItems = append(msgItems, todoMsg)

	// Configure messages
	maxWidth := 80
	for _, msg := range msgItems {
		if m, ok := msg.(interface{ SetMaxWidth(int) }); ok {
			m.SetMaxWidth(maxWidth)
		}
		if m, ok := msg.(interface{ SetConfig(*messages.MessageConfig) }); ok {
			m.SetConfig(&messages.MessageConfig{
				MaxWidth:           maxWidth,
				ShowTimestamp:      true,
				CompactMode:        false,
				EnableMarkdown:     true,
				SyntaxHighlighting: true,
			})
		}
	}

	// Render each message
	var output strings.Builder

	output.WriteString(strings.Repeat("=", 80) + "\n")
	output.WriteString("Taproot Message Components Demo\n")
	output.WriteString(strings.Repeat("=", 80) + "\n\n")

	messageNames := []string{
		"User Message",
		"Assistant Message",
		"Tool Message",
		"Diagnostic Message",
		"Todo Message",
	}

	for i, msg := range msgItems {
		// Initialize the message
		msg.Init()

		// Render message header
		output.WriteString(strings.Repeat("-", 80) + "\n")
		output.WriteString("Message Type: " + messageNames[i] + "\n")
		output.WriteString(strings.Repeat("-", 80) + "\n")

		// Render the message
		output.WriteString(msg.View())
		output.WriteString("\n")

		if i < len(msgItems)-1 {
			output.WriteString("\n")
		}
	}

	output.WriteString(strings.Repeat("=", 80) + "\n")

	// Print the output
	fmt.Println(output.String())
}
