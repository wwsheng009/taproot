package forms

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render"
)

// Form manages a collection of inputs, handling focus traversal and validation.
type Form struct {
	Inputs       []Input
	focusedIndex int
	width        int
}

// NewForm creates a new form with the given inputs.
func NewForm(inputs ...Input) *Form {
	f := &Form{
		Inputs:       inputs,
		focusedIndex: 0,
	}
	
	// Ensure only the first input is focused initially if inputs exist
	if len(f.Inputs) > 0 {
		for i, input := range f.Inputs {
			if i == 0 {
				input.Focus()
			} else {
				input.Blur()
			}
		}
	}
	
	return f
}

// Init implements render.Model.
func (f *Form) Init() error {
	// Initialize all inputs
	for _, input := range f.Inputs {
		if err := input.Init(); err != nil {
			return err
		}
	}
	
	// Return blink command for the initially focused input if applicable
	if len(f.Inputs) > 0 {
		// We can't easily extract the Init cmd from here without calling Init again.
		// But Focus() returns a Cmd.
		// However, NewForm acts as a constructor. Init is called by the runtime.
		// Let's rely on the inputs' Init.
	}
	return nil
}

// Update implements render.Model.
func (f *Form) Update(msg any) (render.Model, render.Cmd) {
	if len(f.Inputs) == 0 {
		return f, nil
	}

	var cmds []render.Cmd

	// Handle global form navigation (Tab/Shift+Tab)
	var keyStr string
	var isKeyMsg bool
	
	if k, ok := msg.(tea.KeyMsg); ok {
		keyStr = k.String()
		isKeyMsg = true
	} else if k, ok := msg.(render.KeyMsg); ok {
		keyStr = k.String()
		isKeyMsg = true
	}

	if isKeyMsg {
		switch keyStr {
		case "tab":
			cmds = append(cmds, f.FocusNext())
			return f, render.Batch(cmds...)
		case "shift+tab":
			cmds = append(cmds, f.FocusPrev())
			return f, render.Batch(cmds...)
		case "enter":
			// Special handling for Enter
			// If it's a TextArea or expanded Select, let it handle Enter.
			// Otherwise, move to next field.
			current := f.Inputs[f.focusedIndex]
			
			shouldPassEnter := false
			switch v := current.(type) {
			case *TextArea:
				shouldPassEnter = true
			case *Select:
				if v.expanded {
					shouldPassEnter = true
				}
			}
			
			if !shouldPassEnter {
				cmds = append(cmds, f.FocusNext())
				return f, render.Batch(cmds...)
			}
		}
	}

	// Pass message to currently focused input
	updatedModel, cmd := f.Inputs[f.focusedIndex].Update(msg)
	if updatedInput, ok := updatedModel.(Input); ok {
		f.Inputs[f.focusedIndex] = updatedInput
	}
	cmds = append(cmds, cmd)

	return f, render.Batch(cmds...)
}

// View implements render.Model.
// It renders inputs vertically with spacing.
func (f *Form) View() string {
	var b strings.Builder
	for i, input := range f.Inputs {
		b.WriteString(input.View())
		if i < len(f.Inputs)-1 {
			b.WriteString("\n\n") // Gap between inputs
		}
	}
	return b.String()
}

// FocusNext moves focus to the next input.
func (f *Form) FocusNext() render.Cmd {
	if len(f.Inputs) == 0 {
		return nil
	}
	
	f.Inputs[f.focusedIndex].Blur()
	f.focusedIndex = (f.focusedIndex + 1) % len(f.Inputs)
	return f.Inputs[f.focusedIndex].Focus()
}

// FocusPrev moves focus to the previous input.
func (f *Form) FocusPrev() render.Cmd {
	if len(f.Inputs) == 0 {
		return nil
	}
	
	f.Inputs[f.focusedIndex].Blur()
	f.focusedIndex--
	if f.focusedIndex < 0 {
		f.focusedIndex = len(f.Inputs) - 1
	}
	return f.Inputs[f.focusedIndex].Focus()
}

// Validate validates all inputs and returns the first error found, or nil.
func (f *Form) Validate() error {
	for _, input := range f.Inputs {
		if err := input.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// FocusedIndex returns the index of the currently focused input.
func (f *Form) FocusedIndex() int {
	return f.focusedIndex
}

// SetFocusedIndex sets the focused input index.
func (f *Form) SetFocusedIndex(i int) render.Cmd {
	if i < 0 || i >= len(f.Inputs) {
		return nil
	}
	
	f.Inputs[f.focusedIndex].Blur()
	f.focusedIndex = i
	return f.Inputs[f.focusedIndex].Focus()
}
