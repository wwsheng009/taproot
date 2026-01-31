package render

// Renderer is the interface for rendering UI components.
// Different rendering engines (Bubbletea, Ultraviolet, etc.) implement this.
type Renderer interface {
	// Render returns the rendered string representation.
	Render() string
}

// Model is the core interface for interactive components.
// It follows the Elm architecture (Model-View-Update).
type Model interface {
	// Init performs initial setup.
	Init() error
	// Update handles incoming messages and returns commands.
	Update(msg any) (Model, Cmd)
	// View returns the string representation for rendering.
	View() string
}

// Cmd represents a side effect command (async operations, etc.).
type Cmd interface{}

// Command is a helper type for simple commands.
type Command func() error

// Execute runs the command.
func (c Command) Execute() error {
	return c()
}

// None returns a nil command (no operation).
func None() Cmd {
	return nil
}

// Batch combines multiple commands into one.
func Batch(cmds ...Cmd) Cmd {
	if len(cmds) == 0 {
		return None()
	}
	return func() error {
		for _, cmd := range cmds {
			if cmd != nil {
				if c, ok := cmd.(Command); ok {
					if err := c.Execute(); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

// Quit returns a command that will quit the application.
func Quit() Cmd {
	return quitCmd{}
}

type quitCmd struct{}

func (q quitCmd) Execute() error {
	return nil
}

func (q quitCmd) IsQuit() bool {
	return true
}

// IsQuit checks if a command is a quit command.
func IsQuit(cmd Cmd) bool {
	if cmd == nil {
		return false
	}
	if q, ok := cmd.(interface{ IsQuit() bool }); ok {
		return q.IsQuit()
	}
	return false
}

// Msg represents any message that can be sent to a Model.
// This could be key events, mouse events, timer events, etc.
type Msg interface{}

// KeyMsg represents a keyboard input message.
type KeyMsg struct {
	Key    string
	Alt    bool
	Ctrl   bool
	Shift  bool
	Type   KeyType
}

// KeyType represents the type of key event.
type KeyType int

const (
	// KeyPress is a standard key press.
	KeyPress KeyType = iota
	// KeyRelease is a key release event.
	KeyRelease
)

// String returns the key string including modifiers.
func (k KeyMsg) String() string {
	prefix := ""
	if k.Alt {
		prefix += "alt+"
	}
	if k.Ctrl {
		prefix += "ctrl+"
	}
	if k.Shift {
		prefix += "shift+"
	}
	return prefix + k.Key
}

// IsMouse returns true if this is a mouse event.
func (k KeyMsg) IsMouse() bool {
	return false
}

// ResizeMsg represents a terminal resize event.
type ResizeMsg struct {
	Width  int
	Height int
}

// TickMsg represents a timer tick event.
type TickMsg struct {
	Time interface{} // time.Time, but using interface{} to avoid import
}

// FocusMsg represents a focus change event.
type FocusMsg struct {
	Focused bool
}

// BlurMsg is sent when a component loses focus.
type BlurMsg struct{}

// FocusMsg is sent when a component receives focus.
type FocusGainMsg struct{}

// ErrorMsg represents an error event.
type ErrorMsg struct {
	Error error
}

// QuitMsg is sent to request application exit.
type QuitMsg struct{}

// WindowSizeMsg is sent when the terminal window is resized.
type WindowSizeMsg struct {
	Width  int
	Height int
}

// CustomMsg allows custom message types.
type CustomMsg struct {
	Data any
	Type string
}
