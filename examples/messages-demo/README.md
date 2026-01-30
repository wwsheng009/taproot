# Taproot Message Components Demo

This demo showcases the engine-agnostic message components implemented in Phase 8.1 of the Taproot framework roadmap.

## Components

The demo demonstrates the following message types:

1. **UserMessage** - User message content with optional file/image attachments
2. **AssistantMessage** - Assistant responses with expandable thinking/reasoning content
3. **ToolMessage** - Tool call execution status and results
4. **DiagnosticMessage** - Error/warning/info diagnostics with severity levels
5. **TodoMessage** - Task list with progress tracking and status

## Running the Demo

```bash
# Build the demo
go build -o bin/messages-demo.exe examples/messages-demo/main.go

# Run the demo
./bin/messages-demo.exe
```

Or directly:

```bash
go run examples/messages-demo/main.go
```

## Component Features

### UserMessage
- Multiple file/image attachments
- Relative timestamp display
- Word and character count
- Configurable styling

### AssistantMessage
- Collapsible thinking/reasoning box
- Markdown rendering for content
- Finish reason display (canceled, error)
- Focus state with different styles

### ToolMessage
- Multiple tool calls tracking
- Status icons (pending, running, completed, error, canceled)
- Collapsible tool details
- Error message display

### DiagnosticMessage
- Multiple diagnostics with different severity levels
- File path and location display
- Error codes
- Severity-based coloring

### TodoMessage
- Progress tracking with percentage
- Status-based styling (pending, in-progress, completed)
- Tag/label support
- Collapsible todo list

## Architecture

All message components implement the `render.Model` interface, making them engine-agnostic:

```go
type Model interface {
    Init() error
    Update(msg any) (Model, Cmd)
    View() string
}
```

They also implement various optional interfaces:
- `Message` - Base message interface with ID, Role, Content, Timestamp
- `Focusable` - Components that can receive focus
- `Expandable` - Components with expandable/collapsible content
- `Identifiable` - Components with unique IDs
- `ClickHandler` - Components that handle mouse clicks
- `KeyEventHandler` - Components that handle keyboard events

## Configuration

Each message can be configured with `MessageConfig`:

```go
config := &messages.MessageConfig{
    MaxWidth:           80,
    ShowTimestamp:      true,
    CompactMode:        false,
    EnableMarkdown:     true,
    SyntaxHighlighting: true,
}
```

## Usage Example

```go
// Create a user message
userMsg := messages.NewUserMessage("msg-1", "Hello!")
userMsg.AddAttachment(messages.Attachment{
    ID:   "att-1",
    Type: "file",
    Name: "document.txt",
    Size: 1024,
})

// Create an assistant message with thinking
assistantMsg := messages.NewAssistantMessage("msg-2", "I can help!")
assistantMsg.SetThinking("Let me think about this...")
assistantMsg.SetThinkingExpanded(true)

// Render the message
userMsg.Init()
assistantMsg.Init()

fmt.Println(userMsg.View())
fmt.Println(assistantMsg.View())
```

## Testing

Run the tests:

```bash
go test ./internal/ui/components/messages/
```

## Phase 8.1 Status

- ✅ types.go - Core type definitions and interfaces
- ✅ user.go - UserMessage component
- ✅ assistant.go - AssistantMessage component
- ✅ tools.go - ToolMessage component
- ✅ diagnostics.go - DiagnosticMessage component
- ✅ todos.go - TodoMessage component
- ✅ messages_test.go - Comprehensive tests
- ✅ examples/messages-demo/ - Interactive demo
