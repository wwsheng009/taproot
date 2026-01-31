package messages

import (
	"time"

	"github.com/wwsheng009/taproot/ui/render"
)

// MessageRole represents the role of a message sender.
type MessageRole int

const (
	// RoleUser represents a user message.
	RoleUser MessageRole = iota
	// RoleAssistant represents an assistant message.
	RoleAssistant
	// RoleSystem represents a system message.
	RoleSystem
	// RoleTool represents a tool call or tool result.
	RoleTool
)

// String returns the string representation of the message role.
func (r MessageRole) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleAssistant:
		return "assistant"
	case RoleSystem:
		return "system"
	case RoleTool:
		return "tool"
	default:
		return "unknown"
	}
}

// Message represents a chat message with role, content, and metadata.
// This is a simplified, engine-agnostic interface that doesn't depend on
// the Crush project's internal message.Message type.
type Message interface {
	render.Model

	// ID returns a unique identifier for this message.
	ID() string

	// Role returns the message role (User, Assistant, System, Tool).
	Role() MessageRole

	// Content returns the message content text.
	Content() string

	// Timestamp returns when the message was created/updated.
	Timestamp() time.Time

	// SetContent sets the message content.
	SetContent(content string)
}

// Identifiable is an interface for items that can provide a unique identifier.
type Identifiable interface {
	ID() string
}

// Expandable is an interface for items that can be expanded or collapsed.
type Expandable interface {
	ToggleExpanded()
	Expanded() bool
	SetExpanded(expanded bool)
}

// Focusable is an interface for items that can receive focus.
type Focusable interface {
	Focus()
	Blur()
	Focused() bool
}

// AnimationHandler is an interface for items that handle animation.
type AnimationHandler interface {
	StartAnimation() render.Cmd
	Animate(msg render.TickMsg) render.Cmd
}

// KeyEventHandler is an interface for items that can handle key events.
type KeyEventHandler interface {
	HandleKeyEvent(msg any) (bool, render.Cmd)
}

// ClickHandler is an interface for items that can handle click events.
type ClickHandler interface {
	HandleClick(line, col int) render.Cmd
}

// MessageItem represents a renderable message item.
type MessageItem interface {
	render.Model
	Identifiable
}

// ToolStatus represents the status of a tool call.
type ToolStatus int

const (
	// ToolStatusPending means the tool call is pending execution.
	ToolStatusPending ToolStatus = iota
	// ToolStatusRunning means the tool call is currently executing.
	ToolStatusRunning
	// ToolStatusCompleted means the tool call completed successfully.
	ToolStatusCompleted
	// ToolStatusError means the tool call failed.
	ToolStatusError
	// ToolStatusCanceled means the tool call was canceled.
	ToolStatusCanceled
)

// String returns the string representation of tool status.
func (s ToolStatus) String() string {
	switch s {
	case ToolStatusPending:
		return "pending"
	case ToolStatusRunning:
		return "running"
	case ToolStatusCompleted:
		return "completed"
	case ToolStatusError:
		return "error"
	case ToolStatusCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}

// ToolCall represents a tool invocation.
type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]any
	Status    ToolStatus
	Result    string
	Error     string
	Timestamp time.Time
}

// DiagnosticSeverity represents the severity of a diagnostic message.
type DiagnosticSeverity int

const (
	// SeverityError represents an error diagnostic.
	SeverityError DiagnosticSeverity = iota
	// SeverityWarning represents a warning diagnostic.
	SeverityWarning
	// SeverityInfo represents an informational diagnostic.
	SeverityInfo
	// SeverityHint represents a hint diagnostic.
	SeverityHint
)

// String returns the string representation of diagnostic severity.
func (s DiagnosticSeverity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	case SeverityHint:
		return "hint"
	default:
		return "unknown"
	}
}

// Diagnostic represents a single diagnostic.
type Diagnostic struct {
	Severity  DiagnosticSeverity
	Message   string
	File      string
	Line      int
	Column    int
	Code      string
	Timestamp time.Time
}

// TodoStatus represents the status of a todo item.
type TodoStatus int

const (
	// TodoStatusPending means the todo is not started.
	TodoStatusPending TodoStatus = iota
	// TodoStatusInProgress means the todo is in progress.
	TodoStatusInProgress
	// TodoStatusCompleted means the todo is completed.
	TodoStatusCompleted
)

// String returns the string representation of todo status.
func (s TodoStatus) String() string {
	switch s {
	case TodoStatusPending:
		return "pending"
	case TodoStatusInProgress:
		return "in_progress"
	case TodoStatusCompleted:
		return "completed"
	default:
		return "unknown"
	}
}

// Todo represents a todo item.
type Todo struct {
	ID          string
	Description string
	Status      TodoStatus
	Progress    float64 // 0.0 to 1.0
	Tags        []string
	Timestamp   time.Time
}

// Attachment represents a file or image attachment.
type Attachment struct {
	ID        string
	Type      string // "file" or "image"
	Name      string
	Path      string
	Size      int64
	Thumbnail string // For images
}

// MessageConfig contains configuration for message rendering.
type MessageConfig struct {
	// MaxWidth is the maximum width for message content.
	MaxWidth int

	// ShowTimestamp controls whether to display timestamps.
	ShowTimestamp bool

	// CompactMode enables compact rendering mode.
	CompactMode bool

	// EnableMarkdown enables Markdown rendering.
	EnableMarkdown bool

	// SyntaxHighlighting enables syntax highlighting for code blocks.
	SyntaxHighlighting bool
}

// DefaultMessageConfig returns the default message configuration.
func DefaultMessageConfig() MessageConfig {
	return MessageConfig{
		MaxWidth:           80,
		ShowTimestamp:      true,
		CompactMode:        false,
		EnableMarkdown:     true,
		SyntaxHighlighting: true,
	}
}
