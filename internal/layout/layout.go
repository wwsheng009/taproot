// Package layout provides core interfaces for TUI components.
package layout

// Focusable represents a component that can receive and manage focus.
type Focusable interface {
	// Focus sets the component to receive keyboard input.
	Focus()

	// Blur removes focus from the component.
	Blur()

	// Focused returns true if the component currently has focus.
	Focused() bool
}

// Sizeable represents a component that can report and set its size.
type Sizeable interface {
	// Size returns the current dimensions of the component.
	Size() (width, height int)

	// SetSize updates the component's dimensions.
	SetSize(width, height int)
}

// Help represents a component that can provide help information.
type Help interface {
	// Help returns help text for the component.
	Help() []string
}

// Positional represents a component that can report and set its position.
type Positional interface {
	// Position returns the component's x, y coordinates.
	Position() (x, y int)

	// SetPosition updates the component's x, y coordinates.
	SetPosition(x, y int)
}
