package forms

import (
	"strings"

	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// TextArea is a multi-line text input component.
type TextArea struct {
	lines       []string
	placeholder string
	focused     bool
	cursorRow   int
	cursorCol   int
	width       int
	height      int
	offset      int // Vertical scroll offset (first visible line)
	
	validators []Validator
	err        error
	blink      bool
	styles     styles.Styles
}

// NewTextArea creates a new text area.
func NewTextArea(placeholder string) *TextArea {
	return &TextArea{
		lines:       []string{""},
		placeholder: placeholder,
		width:       40,
		height:      5,
		blink:       true,
		styles:      styles.DefaultStyles(),
	}
}

// Init implements render.Model.
func (t *TextArea) Init() error {
	return nil
}

// Update implements render.Model.
func (t *TextArea) Update(msg any) (render.Model, render.Cmd) {
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
		case "enter":
			t.InsertNewline()
		case "backspace", "ctrl+h":
			t.Delete()
		case "up":
			t.MoveUp()
		case "down":
			t.MoveDown()
		case "left":
			t.MoveLeft()
		case "right":
			t.MoveRight()
		default:
			if len(msg.String()) == 1 {
				r := []rune(msg.String())[0]
				if r >= 32 {
					t.Insert(r)
				}
			}
		}
	}
	return t, cmd
}

// View implements render.Model.
func (t *TextArea) View() string {
	var b strings.Builder
	
	// Ensure visible lines
	visibleLines := t.height
	if len(t.lines) < visibleLines {
		visibleLines = len(t.lines)
	}
	
	// Adjust offset if cursor is out of view
	if t.cursorRow < t.offset {
		t.offset = t.cursorRow
	} else if t.cursorRow >= t.offset+t.height {
		t.offset = t.cursorRow - t.height + 1
	}
	
	endLine := t.offset + t.height
	if endLine > len(t.lines) {
		endLine = len(t.lines)
	}
	
	for i := t.offset; i < endLine; i++ {
		line := t.lines[i]
		
		if t.focused && i == t.cursorRow {
			cursorStyle := t.styles.Base
			if t.blink {
				cursorStyle = cursorStyle.Reverse(true)
			}
			
			// Render line with cursor
			if t.cursorCol >= len(line) {
				b.WriteString(line)
				b.WriteString(cursorStyle.Render(" "))
			} else {
				b.WriteString(line[:t.cursorCol])
				b.WriteString(cursorStyle.Render(string(line[t.cursorCol])))
				b.WriteString(line[t.cursorCol+1:])
			}
		} else {
			b.WriteString(line)
		}
		
		if i < endLine-1 {
			b.WriteString("\n")
		}
	}
	
	return b.String()
}

// Value returns the full text.
func (t *TextArea) Value() string {
	return strings.Join(t.lines, "\n")
}

// SetValue sets the value.
func (t *TextArea) SetValue(val string) {
	t.lines = strings.Split(val, "\n")
	if len(t.lines) == 0 {
		t.lines = []string{""}
	}
	t.cursorRow = len(t.lines) - 1
	t.cursorCol = len(t.lines[len(t.lines)-1])
}

// Focus focuses the text area.
func (t *TextArea) Focus() render.Cmd {
	t.focused = true
	t.blink = true
	return BlinkCmd()
}

// Blur blurs the text area.
func (t *TextArea) Blur() {
	t.focused = false
}

// Insert inserts a rune.
func (t *TextArea) Insert(r rune) {
	line := t.lines[t.cursorRow]
	before := line[:t.cursorCol]
	after := line[t.cursorCol:]
	t.lines[t.cursorRow] = before + string(r) + after
	t.cursorCol++
}

// InsertNewline inserts a newline.
func (t *TextArea) InsertNewline() {
	line := t.lines[t.cursorRow]
	before := line[:t.cursorCol]
	after := line[t.cursorCol:]
	
	t.lines[t.cursorRow] = before
	// Insert new line after current
	t.lines = append(t.lines[:t.cursorRow+1], append([]string{after}, t.lines[t.cursorRow+1:]...)...)
	
	t.cursorRow++
	t.cursorCol = 0
}

// Delete deletes the character before cursor.
func (t *TextArea) Delete() {
	if t.cursorCol > 0 {
		line := t.lines[t.cursorRow]
		before := line[:t.cursorCol-1]
		after := line[t.cursorCol:]
		t.lines[t.cursorRow] = before + after
		t.cursorCol--
	} else if t.cursorRow > 0 {
		// Merge with previous line
		prevLine := t.lines[t.cursorRow-1]
		currLine := t.lines[t.cursorRow]
		
		newCol := len(prevLine)
		t.lines[t.cursorRow-1] = prevLine + currLine
		
		// Remove current line
		t.lines = append(t.lines[:t.cursorRow], t.lines[t.cursorRow+1:]...)
		
		t.cursorRow--
		t.cursorCol = newCol
	}
}

// MoveUp moves cursor up.
func (t *TextArea) MoveUp() {
	if t.cursorRow > 0 {
		t.cursorRow--
		if t.cursorCol > len(t.lines[t.cursorRow]) {
			t.cursorCol = len(t.lines[t.cursorRow])
		}
	}
}

// MoveDown moves cursor down.
func (t *TextArea) MoveDown() {
	if t.cursorRow < len(t.lines)-1 {
		t.cursorRow++
		if t.cursorCol > len(t.lines[t.cursorRow]) {
			t.cursorCol = len(t.lines[t.cursorRow])
		}
	}
}

// MoveLeft moves cursor left.
func (t *TextArea) MoveLeft() {
	if t.cursorCol > 0 {
		t.cursorCol--
	} else if t.cursorRow > 0 {
		t.cursorRow--
		t.cursorCol = len(t.lines[t.cursorRow])
	}
}

// MoveRight moves cursor right.
func (t *TextArea) MoveRight() {
	if t.cursorCol < len(t.lines[t.cursorRow]) {
		t.cursorCol++
	} else if t.cursorRow < len(t.lines)-1 {
		t.cursorRow++
		t.cursorCol = 0
	}
}

// SetWidth sets the width.
func (t *TextArea) SetWidth(w int) {
	t.width = w
}

// SetHeight sets the height.
func (t *TextArea) SetHeight(h int) {
	t.height = h
}
