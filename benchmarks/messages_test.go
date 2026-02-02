package benchmarks

import (
	"fmt"
	"testing"

	"github.com/wwsheng009/taproot/ui/components/messages"
)

// BenchmarkUserMessageCreation measures performance of creating user messages
func BenchmarkUserMessageCreation(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = messages.NewUserMessage("id-1", "This is a test message")
	}
}

// BenchmarkAssistantMessageCreation measures performance of creating assistant messages
func BenchmarkAssistantMessageCreation(b *testing.B) {
	content := "# Heading\n\nThis is **markdown** content.\n\n- Item 1\n- Item 2\n\n```go\nfunc hello() {\n    println(\"Hello\")\n}\n```"
	for n := 0; n < b.N; n++ {
		_ = messages.NewAssistantMessage("id-1", content)
	}
}

// BenchmarkMessageView_Small measures rendering performance for small messages
func BenchmarkMessageView_Small(b *testing.B) {
	msg := messages.NewUserMessage("id-1", "Short message")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = msg.View()
	}
}

// BenchmarkMessageView_Medium measures rendering performance for medium messages
func BenchmarkMessageView_Medium(b *testing.B) {
	content := "This is a medium-sized message that contains multiple lines of text " +
		"to simulate typical chat messages. It has enough content to require " +
		"some processing but is still relatively small.\n\n" +
		"Second paragraph here with more content."
	msg := messages.NewUserMessage("id-1", content)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = msg.View()
	}
}

// BenchmarkMessageView_Large measures rendering performance for large messages
func BenchmarkMessageView_Large(b *testing.B) {
	content := generateLargeContent(50) // 50 lines
	msg := messages.NewUserMessage("id-1", content)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = msg.View()
	}
}

// BenchmarkAssistantMessageView measures assistant message rendering with markdown
func BenchmarkAssistantMessageView(b *testing.B) {
	content := "# Markdown Test\n\n" +
		"This is a **bold** and *italic* text.\n\n" +
		"## Code Blocks\n\n" +
		"```go\n" +
		"func main() {\n" +
		"    println(\"Hello, World!\")\n" +
		"}\n" +
		"```\n\n" +
		"## Lists\n\n" +
		"- First item\n" +
		"- Second item\n" +
		"  - Nested item\n\n" +
		"## Tables\n\n" +
		"| Column 1 | Column 2 |\n" +
		"|----------|----------|\n" +
		"| Data 1   | Data 2   |\n" +
		"| Data 3   | Data 4   |\n"
	msg := messages.NewAssistantMessage("id-1", content)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = msg.View()
	}
}

// BenchmarkTodoMessageOperations measures todo list operations
func BenchmarkTodoMessageOperations(b *testing.B) {
	msg := messages.NewTodoMessage("id-1", "Task List")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		msg.AddTodo(messages.Todo{
			ID:          "task-1",
			Description: "Task description",
			Status:      messages.TodoStatusPending,
		})
		msg.SetExpanded(true)
		_ = msg.View()
		msg.SetTodos([]messages.Todo{})
	}
}

// BenchmarkMessageExpandCollapse measures expand/collapse performance
func BenchmarkMessageExpandCollapse(b *testing.B) {
	msg := messages.NewAssistantMessage("id-1", generateLargeContent(20))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		msg.ToggleExpanded()
		msg.ToggleExpanded()
	}
}

// BenchmarkMultipleMessages measures rendering multiple messages
func BenchmarkMultipleMessages(b *testing.B) {
	messages := createTestMessages(10)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, msg := range messages {
			_ = msg.View()
		}
	}
}

// Helper function to generate large content
func generateLargeContent(lines int) string {
	content := ""
	for i := 0; i < lines; i++ {
		content += "Line " + generateRandomText(50) + "\n"
	}
	return content
}

// Helper function to generate random text
func generateRandomText(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz "
	result := ""
	for i := 0; i < length; i++ {
		result += string(chars[i%len(chars)])
	}
	return result
}

// Helper function to create test messages
func createTestMessages(count int) []messages.MessageItem {
	msgs := make([]messages.MessageItem, count*2)
	for i := 0; i < count; i++ {
		// User message
		msgs[i*2] = messages.NewUserMessage(
			fmt.Sprintf("user-%d", i),
			fmt.Sprintf("User message %d with some content", i),
		)
		// Assistant message
		msgs[i*2+1] = messages.NewAssistantMessage(
			fmt.Sprintf("assistant-%d", i),
			fmt.Sprintf("# Response %d\n\nThis is the assistant's response.", i),
		)
	}
	return msgs
}
