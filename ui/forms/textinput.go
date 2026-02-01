package forms

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)



// TextInput represents a text input field.
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
	blink            bool
	blinkCtx         int // Context ID for blink command
	cursorStyle      lipgloss.Style
	placeHolderStyle lipgloss.Style
	textStyle        lipgloss.Style
	errorStyle       lipgloss.Style
	prompt           string

	// Border
	showBorder   bool
	focusedStyle lipgloss.Style
	blurredStyle lipgloss.Style
}

// NewTextInput creates a new text input.
func NewTextInput(placeholder string) *TextInput {
	s := styles.DefaultStyles()
	return &TextInput{
		placeholder:      placeholder,
		width:            40,
		maxLength:        0, // 0 means no limit
		blink:            true,
		cursorStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Reverse(true),
		placeHolderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		textStyle:        lipgloss.NewStyle(),
		errorStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		focusedStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(s.BorderColor).
			Padding(0, 1),
		blurredStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(s.Border).
			Padding(0, 1),
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

	// Handle both tea.KeyMsg and render.KeyMsg
	var keyStr string
	if k, ok := msg.(tea.KeyMsg); ok {
		keyStr = k.String()
	} else if k, ok := msg.(render.KeyMsg); ok {
		keyStr = k.String()
	}

	switch keyStr {
		case "up", "down", "enter", "tab", "shift+tab":
			// These keys should be handled by parent (for navigation, etc.)
			return t, nil
		case "backspace", "ctrl+h":
			t.deleteBefore()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "delete", "ctrl+d":
			t.deleteAfter()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "left", "ctrl+b":
			t.moveLeft()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "right", "ctrl+f":
			t.moveRight()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "home", "ctrl+a":
			t.cursor = 0
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "end", "ctrl+e":
			t.cursor = utf8.RuneCountInString(t.value)
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "ctrl+k":
			t.value = t.value[:t.cursor]
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "ctrl+u":
			t.value = t.value[t.cursor:]
			t.cursor = 0
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		default:
			// Insert single printable characters only
			// Check: single character string AND not a control character
			if len(keyStr) == 1 {
				r := []rune(keyStr)[0]
				// Only accept printable characters (space and above)
				if r >= 32 && r <= 126 {
					t.insert(r)
					t.blink = true
					t.blinkCtx = NextBlinkID()
					return t, BlinkCmd(t.blinkCtx)
				}
			}
		}

	// Handle BlinkMsg
	if msg, ok := msg.(BlinkMsg); ok {
		if t.focused && msg.id == t.blinkCtx {
			t.blink = !t.blink
			return t, BlinkCmd(t.blinkCtx)
		}
	}

	return t, cmd
}

// View implements render.Model.
func (t *TextInput) View() string {
	var b strings.Builder

	// Render prompt if set
	if t.prompt != "" {
		b.WriteString(t.prompt)
		b.WriteString(" ")
	}

	// Get rune count for accurate measurement
	valLength := utf8.RuneCountInString(t.value)
	var content string

	// Handle empty state with placeholder
	if valLength == 0 && t.placeholder != "" && !t.focused {
		// Truncate placeholder to width in visual characters
		ph := truncateToWidth(t.placeholder, t.width)
		content = t.placeHolderStyle.Render(ph)
		// For placeholder, we also need to account for visual length for padding
		// However, lipgloss style might add ANSI codes, so we use the raw string length for padding calc
		// assuming placeholder style doesn't change width (e.g. bold is fine)
		// But wait, we should use the raw string length for padding calculation
		// If we use currentVisualLen based on 'ph', it's correct.
		// Let's refine the flow below.
	}

	// Calculate visual length for padding purposes
	var currentVisualLen int

	if valLength == 0 && t.placeholder != "" && !t.focused {
		ph := truncateToWidth(t.placeholder, t.width)
		content = t.placeHolderStyle.Render(ph)
		currentVisualLen = utf8.RuneCountInString(ph)
	} else {
		// Generate display content using runes (not bytes)
		var displayRunes []rune
		if t.hidden {
			// Use rune count for bullets
			displayRunes = make([]rune, valLength)
			for i := range displayRunes {
				displayRunes[i] = '•'
			}
		} else {
			displayRunes = []rune(t.value)
		}

		// Ensure cursor stays within bounds
		if t.cursor > valLength {
			t.cursor = valLength
		}

		// Simple scrolling based on rune count
		displayRunes, scrollStart := scrollRunes(displayRunes, t.cursor, t.width)

		// Calculate cursor offset relative to scrolled view
		cursorOffset := t.cursor - scrollStart

		// Render with cursor
		currentVisualLen = len(displayRunes)
		var sb strings.Builder
		if cursorOffset >= len(displayRunes) {
			// Cursor at end
			sb.WriteString(t.textStyle.Render(string(displayRunes)))
			if t.focused && t.blink {
				sb.WriteString(t.cursorStyle.Render(" "))
				currentVisualLen++
			}
		} else {
			// Cursor in middle
			before := displayRunes[:cursorOffset]
			cursorChar := displayRunes[cursorOffset]
			after := displayRunes[cursorOffset+1:]

			sb.WriteString(t.textStyle.Render(string(before)))
			if t.focused && t.blink {
				sb.WriteString(t.cursorStyle.Render(string(cursorChar)))
			} else {
				sb.WriteString(t.textStyle.Render(string(cursorChar)))
			}
			sb.WriteString(t.textStyle.Render(string(after)))
		}
		content = sb.String()
	}

	b.WriteString(content)

	// Pad with spaces to match width if showBorder is enabled
	// This ensures the border width is consistent regardless of content length
	if t.showBorder {
		if padding := t.width - currentVisualLen; padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
	}

	// Render error if present
	if t.err != nil {
		b.WriteString("\n")
		b.WriteString(t.errorStyle.Render(t.err.Error()))
	}

	result := b.String()

	if t.showBorder {
		if t.focused {
			return t.focusedStyle.Render(result)
		}
		return t.blurredStyle.Render(result)
	}

	return result
}

// scrollRunes scrolls rune array to keep cursor visible
// Returns the scrolled runes and the start position in the original array
// Strategy: Keep cursor as close to the right as possible, only scroll when needed
func scrollRunes(runes []rune, cursor, width int) ([]rune, int) {
	// If content fits, no scrolling needed
	if len(runes) <= width {
		return runes, 0
	}

	// Calculate scroll start position
	// Keep cursor visible, preferably near the right edge for natural typing feel
	start := 0

	// Only scroll left when cursor is near or beyond the right edge
	// Leave some padding on the right for better UX
	paddingRight := 2
	if cursor >= width-paddingRight {
		start = cursor - width + paddingRight
	}

	// Ensure start is within bounds
	if start < 0 {
		start = 0
	}

	// Calculate end position
	end := start + width
	if end > len(runes) {
		end = len(runes)
		// Adjust start if we're at the end of the text
		start = end - width
		if start < 0 {
			start = 0
		}
	}

	return runes[start:end], start
}

// truncateToWidth truncates a string to fit within maxWidth in visual width
func truncateToWidth(s string, maxWidth int) string {
	// For simplicity, assume 1 byte = 1 visual width for ASCII
	// This works for most use cases
	if len(s) <= maxWidth {
		return s
	}
	return s[:maxWidth]
}

// Value returns the current value.
func (t *TextInput) Value() string {
	return t.value
}

// SetValue sets the value.
func (t *TextInput) SetValue(v string) {
	t.value = v
	cursorLen := utf8.RuneCountInString(v)
	if t.cursor > cursorLen {
		t.cursor = cursorLen
	}
}

// Focus focuses the input.
func (t *TextInput) Focus() render.Cmd {
	t.focused = true
	t.blink = true
	t.blinkCtx = NextBlinkID()
	return BlinkCmd(t.blinkCtx)
}

// Blur blurs the input.
func (t *TextInput) Blur() {
	t.focused = false
}

// Focused returns true if focused.
func (t *TextInput) Focused() bool {
	return t.focused
}

// Validate validates the input.
func (t *TextInput) Validate() error {
	for _, v := range t.validators {
		if err := v(t.value); err != nil {
			t.err = err
			return err
		}
	}
	t.err = nil
	return nil
}

// Error returns the last validation error.
func (t *TextInput) Error() error {
	return t.err
}

// AddValidator adds a validator.
func (t *TextInput) AddValidator(v Validator) {
	t.validators = append(t.validators, v)
}

// SetValidators sets the validators.
func (t *TextInput) SetValidators(v ...Validator) {
	t.validators = v
}

// SetPrompt sets the prompt.
func (t *TextInput) SetPrompt(p string) {
	t.prompt = p
}

// SetHidden sets hidden mode.
func (t *TextInput) SetHidden(h bool) {
	t.hidden = h
}

// SetWidth sets the width.
func (t *TextInput) SetWidth(w int) {
	t.width = w
}

// SetMaxLength sets max length.
func (t *TextInput) SetMaxLength(l int) {
	t.maxLength = l
}

// SetShowBorder sets whether to show the border.
func (t *TextInput) SetShowBorder(show bool) {
	t.showBorder = show
}

// Helpers

func (t *TextInput) insert(r rune) {
	if t.maxLength > 0 && utf8.RuneCountInString(t.value) >= t.maxLength {
		return
	}

	left := t.value[:t.cursor]
	right := t.value[t.cursor:]
	t.value = left + string(r) + right
	t.cursor++
}

func (t *TextInput) deleteBefore() {
	if t.cursor > 0 {
		runes := []rune(t.value)
		left := string(runes[:t.cursor-1])
		right := string(runes[t.cursor:])
		t.value = left + right
		t.cursor--
	}
}

func (t *TextInput) deleteAfter() {
	if t.cursor < utf8.RuneCountInString(t.value) {
		runes := []rune(t.value)
		left := string(runes[:t.cursor])
		right := string(runes[t.cursor+1:])
		t.value = left + right
	}
}

func (t *TextInput) moveLeft() {
	if t.cursor > 0 {
		t.cursor--
	}
}

func (t *TextInput) moveRight() {
	if t.cursor < utf8.RuneCountInString(t.value) {
		t.cursor++
	}
}
