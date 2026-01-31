package forms

import (
	"strings"
	"time"

	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// BlinkMsg is sent to blink the cursor.
type BlinkMsg struct{}

// BlinkCmd returns a command to blink the cursor.
func BlinkCmd() render.Cmd {
	return func() render.Msg {
		time.Sleep(time.Millisecond * 500)
		return BlinkMsg{}
	}
}

// TextInput is a text input component.
type TextInput struct {
	value       string
	placeholder string
	focused     bool
	hidden      bool // For passwords
	cursor      int
	width       int
	maxLength   int
	
	// Validation
	validators []Validator
	err        error

	// Cursor blinking
	blink      bool
	
	// Styles
	styles styles.Styles
}

// NewTextInput creates a new text input.
func NewTextInput(placeholder string) *TextInput {
	return &TextInput{
		placeholder: placeholder,
		width:       40,
		blink:       true,
		styles:      styles.DefaultStyles(),
	}
}

// Init implements render.Model.
func (t *TextInput) Init() error {
	return nil
}

// Update implements render.Model.
func (t *TextInput) Update(msg any) (render.Model, render.Cmd) {
	if !t.focused {
		return t, nil
	}

	var cmd render.Cmd

	switch msg := msg.(type) {
	case BlinkMsg:
		t.blink = !t.blink
		cmd = BlinkCmd()
		
	case render.KeyMsg:
		switch msg.String() {
		case "backspace", "ctrl+h":
			t.Delete()
		case "delete":
			t.DeleteForward()
		case "ctrl+u": // Clear line
			t.Clear()
		case "left", "ctrl+b":
			t.MoveCursorLeft()
		case "right", "ctrl+f":
			t.MoveCursorRight()
		case "home", "ctrl+a":
			t.MoveCursorToStart()
		case "end", "ctrl+e":
			t.MoveCursorToEnd()
		default:
			// Regular character input
			// Simple check for printable characters
			if len(msg.String()) == 1 {
				r := []rune(msg.String())[0]
				if r >= 32 { // Skip control chars
					t.Insert(r)
				}
			}
		}
		
		// Validate on input
		t.Validate()
	}

	return t, cmd
}

// View implements render.Model.
func (t *TextInput) View() string {
	var b strings.Builder

	// Value or Placeholder
	val := t.value
	if t.hidden {
		val = strings.Repeat("•", len(val))
	}
	
	if val == "" && t.placeholder != "" && !t.focused {
		return t.styles.Base.Foreground(t.styles.FgMuted).Render(t.placeholder)
	}

	// Render with cursor if focused
	if t.focused {
		cursorStyle := t.styles.Base
		if t.blink {
			cursorStyle = cursorStyle.Reverse(true)
		}
		
		// Split by cursor position
		before := val[:t.cursor]
		
		var charAtCursor string
		if t.cursor < len(val) {
			charAtCursor = string(val[t.cursor])
		} else {
			charAtCursor = " "
		}
		
		after := ""
		if t.cursor+1 < len(val) {
			after = val[t.cursor+1:]
		}
		
		b.WriteString(before)
		b.WriteString(cursorStyle.Render(charAtCursor))
		b.WriteString(after)
	} else {
		b.WriteString(val)
	}
	
	return b.String()
}

// Focus focuses the input and starts blinking.
func (t *TextInput) Focus() render.Cmd {
	t.focused = true
	t.blink = true
	return BlinkCmd()
}

// Blur blurs the input.
func (t *TextInput) Blur() {
	t.focused = false
}

// Value returns the current value.
func (t *TextInput) Value() string {
	return t.value
}

// SetValue sets the value.
func (t *TextInput) SetValue(val string) {
	t.value = val
	t.setCursor(len(val))
	t.Validate()
}

// Validate runs validators and updates error state.
func (t *TextInput) Validate() error {
	t.err = nil
	for _, v := range t.validators {
		if err := v(t.value); err != nil {
			t.err = err
			return err
		}
	}
	return nil
}

// Error returns the current validation error.
func (t *TextInput) Error() error {
	return t.err
}

// AddValidator adds a validator.
func (t *TextInput) AddValidator(v Validator) {
	t.validators = append(t.validators, v)
}

// ... Cursor movement and editing methods ported from dialog/types.go ...

// Insert inserts a rune at the cursor position.
func (t *TextInput) Insert(r rune) {
	if t.maxLength > 0 && len(t.value) >= t.maxLength {
		return
	}
	
	before := t.value[:t.cursor]
	after := t.value[t.cursor:]
	t.value = before + string(r) + after
	t.cursor++
}

// Delete deletes the rune before the cursor.
func (t *TextInput) Delete() {
	if t.cursor > 0 {
		before := t.value[:t.cursor-1]
		after := t.value[t.cursor:]
		t.value = before + after
		t.cursor--
	}
}

// DeleteForward deletes the rune at the cursor.
func (t *TextInput) DeleteForward() {
	if t.cursor < len(t.value) {
		before := t.value[:t.cursor]
		after := t.value[t.cursor+1:]
		t.value = before + after
	}
}

// MoveCursorLeft moves the cursor left.
func (t *TextInput) MoveCursorLeft() {
	if t.cursor > 0 {
		t.cursor--
	}
}

// MoveCursorRight moves the cursor right.
func (t *TextInput) MoveCursorRight() {
	if t.cursor < len(t.value) {
		t.cursor++
	}
}

// MoveCursorToStart moves cursor to start.
func (t *TextInput) MoveCursorToStart() {
	t.cursor = 0
}

// MoveCursorToEnd moves cursor to end.
func (t *TextInput) MoveCursorToEnd() {
	t.cursor = len(t.value)
}

// Clear clears the input value.
func (t *TextInput) Clear() {
	t.value = ""
	t.cursor = 0
}

// setCursor sets the cursor position safely.
func (t *TextInput) setCursor(pos int) {
	if pos < 0 {
		pos = 0
	} else if pos > len(t.value) {
		pos = len(t.value)
	}
	t.cursor = pos
}

// SetPlaceholder sets the placeholder.
func (t *TextInput) SetPlaceholder(p string) {
	t.placeholder = p
}

// SetHidden sets the hidden state.
func (t *TextInput) SetHidden(h bool) {
	t.hidden = h
}

// SetMaxLength sets the max length.
func (t *TextInput) SetMaxLength(l int) {
	t.maxLength = l
}

// SetWidth sets the width.
func (t *TextInput) SetWidth(w int) {
	t.width = w
}
