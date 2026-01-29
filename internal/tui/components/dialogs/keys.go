package dialogs

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines keyboard bindings for components.
type KeyMap struct {
	Close key.Binding
}

// DefaultKeyMap returns the default key bindings for dialogs.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Close: key.NewBinding(
			key.WithKeys("esc", "alt+esc"),
			key.WithHelp("esc", "close"),
		),
	}
}

// KeyBindings returns all key bindings as a slice.
func (k KeyMap) KeyBindings() []key.Binding {
	return []key.Binding{
		k.Close,
	}
}

// FullHelp implements help.KeyMap and returns help grouped by category.
func (k KeyMap) FullHelp() [][]key.Binding {
	m := [][]key.Binding{}
	slice := k.KeyBindings()
	for i := 0; i < len(slice); i += 4 {
		end := min(i+4, len(slice))
		m = append(m, slice[i:end])
	}
	return m
}

// ShortHelp implements help.KeyMap and returns essential key bindings.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Close,
	}
}
