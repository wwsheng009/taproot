package messages

import (
	"testing"
	"time"

	"github.com/wwsheng009/taproot/ui/render"
)

func TestUserMessage(t *testing.T) {
	t.Run("NewAndID", func(t *testing.T) {
		msg := NewUserMessage("test-id", "Hello, world!")
		if msg.ID() != "test-id" {
			t.Errorf("Expected ID 'test-id', got '%s'", msg.ID())
		}
		if msg.Content() != "Hello, world!" {
			t.Errorf("Expected content 'Hello, world!', got '%s'", msg.Content())
		}
		if msg.Role() != RoleUser {
			t.Errorf("Expected role RoleUser, got %v", msg.Role())
		}
	})

	t.Run("FocusAndBlur", func(t *testing.T) {
		msg := NewUserMessage("test-id", "Test content")
		if msg.Focused() {
			t.Error("Message should not be focused initially")
		}

		msg.Focus()
		if !msg.Focused() {
			t.Error("Message should be focused after Focus()")
		}

		msg.Blur()
		if msg.Focused() {
			t.Error("Message should not be focused after Blur()")
		}
	})

	t.Run("Attachments", func(t *testing.T) {
		msg := NewUserMessage("test-id", "Test content")
		attachment := Attachment{
			ID:   "att-1",
			Type: "file",
			Name: "test.txt",
		}

		msg.AddAttachment(attachment)
		if len(msg.Attachments()) != 1 {
			t.Errorf("Expected 1 attachment, got %d", len(msg.Attachments()))
		}

		msg.RemoveAttachment(0)
		if len(msg.Attachments()) != 0 {
			t.Errorf("Expected 0 attachments after removal, got %d", len(msg.Attachments()))
		}
	})

	t.Run("WordCount", func(t *testing.T) {
		msg := NewUserMessage("test-id", "Hello world, this is a test")
		if msg.WordCount() != 6 {
			t.Errorf("Expected 6 words, got %d", msg.WordCount())
		}
	})

	t.Run("CharacterCount", func(t *testing.T) {
		msg := NewUserMessage("test-id", "Hello")
		if msg.CharacterCount() != 5 {
			t.Errorf("Expected 5 characters, got %d", msg.CharacterCount())
		}
	})

	t.Run("Timestamp", func(t *testing.T) {
		msg := NewUserMessage("test-id", "Test")
		if msg.Timestamp().IsZero() {
			t.Error("Timestamp should not be zero")
		}

		ts := time.Now().Add(-time.Hour)
		msg.SetTimestamp(ts)
		if !msg.Timestamp().Equal(ts) {
			t.Error("Timestamp should match the set value")
		}
	})

	t.Run("InitAndView", func(t *testing.T) {
		msg := NewUserMessage("test-id", "Test content")
		if err := msg.Init(); err != nil {
			t.Errorf("Init() should not return error, got %v", err)
		}

		view := msg.View()
		if view == "" {
			t.Error("View() should not return empty string")
		}
	})

	t.Run("Update", func(t *testing.T) {
		msg := NewUserMessage("test-id", "Test content")

		newMsg, cmd := msg.Update(&render.FocusGainMsg{})
		if newMsg == nil {
			t.Error("Update() should return a model")
		}
		if cmd != nil {
			t.Error("Update() should return nil command")
		}

		if !msg.Focused() {
			t.Error("Message should be focused after FocusGainMsg")
		}
	})
}

func TestAssistantMessage(t *testing.T) {
	t.Run("NewAndID", func(t *testing.T) {
		msg := NewAssistantMessage("test-id", "Assistant response")
		if msg.ID() != "test-id" {
			t.Errorf("Expected ID 'test-id', got '%s'", msg.ID())
		}
		if msg.Content() != "Assistant response" {
			t.Errorf("Expected content 'Assistant response', got '%s'", msg.Content())
		}
		if msg.Role() != RoleAssistant {
			t.Errorf("Expected role RoleAssistant, got %v", msg.Role())
		}
	})

	t.Run("ThinkingContent", func(t *testing.T) {
		msg := NewAssistantMessage("test-id", "Response")
		thinking := "Let me think about this..."
		msg.SetThinking(thinking)

		if msg.Thinking() != thinking {
			t.Errorf("Expected thinking '%s', got '%s'", thinking, msg.Thinking())
		}
	})

	t.Run("ThinkingExpansion", func(t *testing.T) {
		msg := NewAssistantMessage("test-id", "Response")
		thinking := "Thinking content"
		msg.SetThinking(thinking)

		if msg.ThinkingExpanded() {
			t.Error("Thinking should not be expanded initially")
		}

		msg.SetThinkingExpanded(true)
		if !msg.ThinkingExpanded() {
			t.Error("Thinking should be expanded after SetThinkingExpanded(true)")
		}

		msg.ToggleExpanded()
		if msg.ThinkingExpanded() {
			t.Error("Thinking should be collapsed after ToggleExpanded()")
		}
	})

	t.Run("FinishedState", func(t *testing.T) {
		msg := NewAssistantMessage("test-id", "Response")

		if !msg.Finished() {
			t.Error("Message should be finished by default")
		}

		msg.SetFinished(false)
		if msg.Finished() {
			t.Error("Message should not be finished after SetFinished(false)")
		}

		msg.SetFinishReason("error")
		if msg.FinishReason() != "error" {
			t.Errorf("Expected finish reason 'error', got '%s'", msg.FinishReason())
		}
	})

	t.Run("SetContent", func(t *testing.T) {
		msg := NewAssistantMessage("test-id", "Initial")
		newContent := "Updated content"
		msg.SetContent(newContent)

		if msg.Content() != newContent {
			t.Errorf("Expected content '%s', got '%s'", newContent, msg.Content())
		}
	})
}

func TestToolMessage(t *testing.T) {
	t.Run("NewAndID", func(t *testing.T) {
		msg := NewToolMessage("test-id", "search_files")
		if msg.ID() != "test-id" {
			t.Errorf("Expected ID 'test-id', got '%s'", msg.ID())
		}
		if msg.Name() != "search_files" {
			t.Errorf("Expected name 'search_files', got '%s'", msg.Name())
		}
		if msg.Role() != RoleTool {
			t.Errorf("Expected role RoleTool, got %v", msg.Role())
		}
	})

	t.Run("ToolCalls", func(t *testing.T) {
		msg := NewToolMessage("test-id", "search_files")

		call := ToolCall{
			ID:       "call-1",
			Name:     "search",
			Arguments: map[string]any{"query": "test"},
			Status:   ToolStatusPending,
			Timestamp: time.Now(),
		}

		msg.AddCall(call)
		if msg.CallCount() != 1 {
			t.Errorf("Expected 1 call, got %d", msg.CallCount())
		}

		msg.UpdateCallStatus("call-1", ToolStatusCompleted)
		if msg.CompletedCount() != 1 {
			t.Errorf("Expected 1 completed call, got %d", msg.CompletedCount())
		}

		msg.UpdateCallResult("call-1", "Result text", "")
		if msg.calls[0].Result != "Result text" {
			t.Errorf("Expected result 'Result text', got '%s'", msg.calls[0].Result)
		}
	})

	t.Run("ToolCallsFromConstructor", func(t *testing.T) {
		calls := []ToolCall{
			{
				ID:       "call-1",
				Name:     "read_file",
				Status:   ToolStatusCompleted,
				Timestamp: time.Now(),
			},
		}
		msg := NewToolMessageFromCalls("test-id", "Tool", calls)

		if msg.CallCount() != 1 {
			t.Errorf("Expected 1 call, got %d", msg.CallCount())
		}
	})

	t.Run("Expansion", func(t *testing.T) {
		msg := NewToolMessage("test-id", "Tool")
		call := ToolCall{
			ID:       "call-1",
			Name:     "test",
			Status:   ToolStatusPending,
			Timestamp: time.Now(),
		}
		msg.AddCall(call)

		if msg.Expanded() {
			t.Error("Message should not be expanded initially")
		}

		msg.ToggleExpanded()
		if !msg.Expanded() {
			t.Error("Message should be expanded after ToggleExpanded()")
		}
	})
}

func TestDiagnosticMessage(t *testing.T) {
	t.Run("NewAndID", func(t *testing.T) {
		msg := NewDiagnosticMessage("test-id", "Build Errors")
		if msg.ID() != "test-id" {
			t.Errorf("Expected ID 'test-id', got '%s'", msg.ID())
		}
		if msg.Title() != "Build Errors" {
			t.Errorf("Expected title 'Build Errors', got '%s'", msg.Title())
		}
	})

	t.Run("Diagnostics", func(t *testing.T) {
		msg := NewDiagnosticMessage("test-id", "Test")

		diag := Diagnostic{
			Severity:  SeverityError,
			Message:   "undefined variable",
			File:      "test.go",
			Line:      10,
			Timestamp: time.Now(),
		}

		msg.AddDiagnostic(diag)
		if msg.DiagnosticCount() != 1 {
			t.Errorf("Expected 1 diagnostic, got %d", msg.DiagnosticCount())
		}

		if msg.ErrorCount() != 1 {
			t.Errorf("Expected 1 error, got %d", msg.ErrorCount())
		}
	})

	t.Run("Severities", func(t *testing.T) {
		diagnostics := []Diagnostic{
			{Severity: SeverityError, Message: "Error 1"},
			{Severity: SeverityError, Message: "Error 2"},
			{Severity: SeverityWarning, Message: "Warning 1"},
			{Severity: SeverityInfo, Message: "Info 1"},
			{Severity: SeverityHint, Message: "Hint 1"},
		}

		msg := NewDiagnosticMessageFromDiagnostics("test-id", "All", diagnostics)

		if msg.ErrorCount() != 2 {
			t.Errorf("Expected 2 errors, got %d", msg.ErrorCount())
		}

		if msg.WarningCount() != 1 {
			t.Errorf("Expected 1 warning, got %d", msg.WarningCount())
		}

		if msg.InfoCount() != 1 {
			t.Errorf("Expected 1 info, got %d", msg.InfoCount())
		}

		if msg.HintCount() != 1 {
			t.Errorf("Expected 1 hint, got %d", msg.HintCount())
		}
	})
}

func TestTodoMessage(t *testing.T) {
	t.Run("NewAndID", func(t *testing.T) {
		msg := NewTodoMessage("test-id", "Project Tasks")
		if msg.ID() != "test-id" {
			t.Errorf("Expected ID 'test-id', got '%s'", msg.ID())
		}
		if msg.Title() != "Project Tasks" {
			t.Errorf("Expected title 'Project Tasks', got '%s'", msg.Title())
		}
	})

	t.Run("Todos", func(t *testing.T) {
		msg := NewTodoMessage("test-id", "Tasks")

		todo := Todo{
			ID:          "todo-1",
			Description: "Write tests",
			Status:      TodoStatusPending,
			Timestamp:   time.Now(),
		}

		msg.AddTodo(todo)
		if msg.TodoCount() != 1 {
			t.Errorf("Expected 1 todo, got %d", msg.TodoCount())
		}

		if msg.PendingCount() != 1 {
			t.Errorf("Expected 1 pending todo, got %d", msg.PendingCount())
		}
	})

	t.Run("TodoStatusUpdates", func(t *testing.T) {
		msg := NewTodoMessage("test-id", "Tasks")
		todo := Todo{
			ID:          "todo-1",
			Description: "Task 1",
			Status:      TodoStatusPending,
			Progress:    0.0,
			Timestamp:   time.Now(),
		}
		msg.AddTodo(todo)

		msg.UpdateTodoStatus("todo-1", TodoStatusInProgress)
		if msg.InProgressCount() != 1 {
			t.Errorf("Expected 1 in-progress todo, got %d", msg.InProgressCount())
		}

		msg.UpdateTodoProgress("todo-1", 1.0)
		if msg.CompletedCount() != 1 {
			t.Errorf("Expected 1 completed todo, got %d", msg.CompletedCount())
		}

		if msg.todos[0].Progress != 1.0 {
			t.Errorf("Expected progress 1.0, got %f", msg.todos[0].Progress)
		}
	})

	t.Run("OverallProgress", func(t *testing.T) {
		todos := []Todo{
			{ID: "1", Status: TodoStatusCompleted, Progress: 1.0},
			{ID: "2", Status: TodoStatusCompleted, Progress: 1.0},
			{ID: "3", Status: TodoStatusPending, Progress: 0.0},
		}
		msg := NewTodoMessageFromTodos("test-id", "All", todos)

		progress := msg.OverallProgress()
		expected := 2.0 / 3.0

		if progress != expected {
			t.Errorf("Expected progress %f, got %f", expected, progress)
		}
	})
}

func TestMessageInterfaces(t *testing.T) {
	t.Run("AllMessagesImplementModel", func(t *testing.T) {
		messages := []struct {
			name    string
			msg     render.Model
			id      string
		}{
			{"UserMessage", NewUserMessage("u1", "test"), "u1"},
			{"AssistantMessage", NewAssistantMessage("a1", "test"), "a1"},
			{"ToolMessage", NewToolMessage("t1", "tool"), "t1"},
			{"DiagnosticMessage", NewDiagnosticMessage("d1", "test"), "d1"},
			{"TodoMessage", NewTodoMessage("todo1", "test"), "todo1"},
		}

		for _, tc := range messages {
			t.Run(tc.name, func(t *testing.T) {
				if err := tc.msg.Init(); err != nil {
					t.Errorf("%s: Init() should not return error, got %v", tc.name, err)
				}

				view := tc.msg.View()
				if view == "" {
					t.Errorf("%s: View() should not return empty string", tc.name)
				}
			})
		}
	})

	t.Run("FocusableInterface", func(t *testing.T) {
		messages := []Focusable{
			NewUserMessage("u1", "test"),
			NewAssistantMessage("a1", "test"),
			NewToolMessage("t1", "tool"),
			NewDiagnosticMessage("d1", "test"),
			NewTodoMessage("todo1", "test"),
		}

		for i, msg := range messages {
			msg.Focus()
			if !msg.Focused() {
				t.Errorf("Message %d should be focused after Focus()", i)
			}

			msg.Blur()
			if msg.Focused() {
				t.Errorf("Message %d should not be focused after Blur()", i)
			}
		}
	})

	t.Run("IdentifiableInterface", func(t *testing.T) {
		idables := []Identifiable{
			NewUserMessage("u1", "test"),
			NewAssistantMessage("a1", "test"),
			NewToolMessage("t1", "tool"),
			NewDiagnosticMessage("d1", "test"),
			NewTodoMessage("todo1", "test"),
		}

		expectedIDs := []string{"u1", "a1", "t1", "d1", "todo1"}

		for i, idable := range idables {
			if idable.ID() != expectedIDs[i] {
				t.Errorf("Expected ID '%s', got '%s'", expectedIDs[i], idable.ID())
			}
		}
	})

	t.Run("ExpandableInterface", func(t *testing.T) {
		expandables := []Expandable{
			NewToolMessage("t1", "tool"),
			NewDiagnosticMessage("d1", "test"),
			NewTodoMessage("todo1", "test"),
		}

		for i, exp := range expandables {
			
			if exp.Expanded() {
				t.Errorf("Expandable %d should not be expanded initially", i)
			}

			exp.SetExpanded(true)
			if !exp.Expanded() {
				t.Errorf("Expandable %d should be expanded after SetExpanded(true)", i)
			}

			exp.ToggleExpanded()
			if exp.Expanded() {
				t.Errorf("Expandable %d should be collapsed after ToggleExpanded()", i)
			}
		}
	})
}

func TestEnums(t *testing.T) {
	t.Run("MessageRoleString", func(t *testing.T) {
		tests := []struct {
			role    MessageRole
			str     string
		}{
			{RoleUser, "user"},
			{RoleAssistant, "assistant"},
			{RoleSystem, "system"},
			{RoleTool, "tool"},
		}

		for _, tc := range tests {
			if tc.role.String() != tc.str {
				t.Errorf("Expected %v.String() = %s, got %s", tc.role, tc.str, tc.role.String())
			}
		}
	})

	t.Run("ToolStatus", func(t *testing.T) {
		if ToolStatusPending.String() != "pending" {
			t.Errorf("ToolStatusPending.String() should return 'pending'")
		}
		if ToolStatusError.String() != "error" {
			t.Errorf("ToolStatusError.String() should return 'error'")
		}
	})

	t.Run("DiagnosticSeverity", func(t *testing.T) {
		if SeverityError.String() != "error" {
			t.Errorf("SeverityError.String() should return 'error'")
		}
		if SeverityWarning.String() != "warning" {
			t.Errorf("SeverityWarning.String() should return 'warning'")
		}
	})

	t.Run("TodoStatus", func(t *testing.T) {
		if TodoStatusPending.String() != "pending" {
			t.Errorf("TodoStatusPending.String() should return 'pending'")
		}
		if TodoStatusCompleted.String() != "completed" {
			t.Errorf("TodoStatusCompleted.String() should return 'completed'")
		}
	})
}

func TestMessageConfig(t *testing.T) {
	t.Run("DefaultMessageConfig", func(t *testing.T) {
		config := DefaultMessageConfig()

		if config.MaxWidth != 80 {
			t.Errorf("Expected MaxWidth 80, got %d", config.MaxWidth)
		}

		if !config.ShowTimestamp {
			t.Error("Expected ShowTimestamp to be true")
		}

		if !config.EnableMarkdown {
			t.Error("Expected EnableMarkdown to be true")
		}
	})
}
