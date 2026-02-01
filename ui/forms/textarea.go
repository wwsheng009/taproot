package forms

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// TextArea is a multi-line text input component.
type TextArea struct {
	lines       []string
	label       string // Label displayed above the text area
	placeholder string
	focused     bool
	cursorRow   int
	cursorCol   int
	width       int
	height      int
	offset      int // Vertical scroll offset (first visible line)
	wrap        bool // Enable automatic word wrap

	validators []Validator
	err        error
	blink      bool
	blinkCtx   int
	styles     styles.Styles

	// Border
	showBorder   bool
	focusedStyle lipgloss.Style
	blurredStyle lipgloss.Style
	labelStyle   lipgloss.Style
}

// NewTextArea creates a new text area.
func NewTextArea(placeholder string) *TextArea {
	s := styles.DefaultStyles()
	return &TextArea{
		lines:       []string{""},
		placeholder: placeholder,
		width:       40,
		height:      5,
		wrap:        true, // Enable word wrap by default
		blink:       true,
		styles:      s,
		focusedStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(s.BorderColor).
			Padding(0, 1),
		blurredStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(s.Border).
			Padding(0, 1),
		labelStyle: lipgloss.NewStyle().
			Foreground(s.FgBase).
			Bold(true),
	}
}

// Init implements render.Model.
func (t *TextArea) Init() render.Cmd {
	return nil
}

// Update implements render.Model.
func (t *TextArea) Update(msg any) (render.Model, render.Cmd) {
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

	if keyStr != "" {
		switch keyStr {
		case "tab", "shift+tab", "ctrl+c":
			// These keys should be handled by parent
			return t, nil
		case "enter":
			t.InsertNewline()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "backspace", "ctrl+h":
			t.Delete()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "up":
			t.MoveUp()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "down":
			t.MoveDown()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "left":
			t.MoveLeft()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "right":
			t.MoveRight()
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "home", "ctrl+a":
			// Move to start of current line
			t.cursorCol = 0
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		case "end", "ctrl+e":
			// Move to end of current line
			if t.cursorRow < len(t.lines) {
				t.cursorCol = len(t.lines[t.cursorRow])
			}
			t.blink = true
			t.blinkCtx = NextBlinkID()
			return t, BlinkCmd(t.blinkCtx)
		default:
			// Insert single printable characters only
			if len(keyStr) == 1 {
				r := []rune(keyStr)[0]
				// Only accept printable characters (space and above)
				if r >= 32 && r <= 126 {
					t.Insert(r)
					t.blink = true
					t.blinkCtx = NextBlinkID()
					return t, BlinkCmd(t.blinkCtx)
				}
			}
		}
	}

	// Handle BlinkMsg
	if msg, ok := msg.(BlinkMsg); ok {
		if t.focused && msg.id == t.blinkCtx {
			t.blink = !t.blink
			cmd = BlinkCmd(t.blinkCtx)
		}
	}

	return t, cmd
}

// View implements render.Model.
func (t *TextArea) View() string {
	var b strings.Builder

	// Render Label
	if t.label != "" {
		b.WriteString(t.labelStyle.Render(t.label))
		b.WriteString("\n")
	}

	// Handle empty state with placeholder
	if len(t.lines) == 1 && t.lines[0] == "" && t.placeholder != "" && !t.focused {
		// Render placeholder
		ph := t.placeholder
		ph = ph[:min(len(ph), t.width)]
		// Pad to width
		if len(ph) < t.width {
			ph += strings.Repeat(" ", t.width-utf8.RuneCountInString(ph))
		}
		
		// Use a gray color for placeholder to match TextInput
		placeholderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		b.WriteString(placeholderStyle.Render(ph))
		b.WriteString("\n")

		// Fill remaining lines to match height
		emptyLine := strings.Repeat(" ", t.width)
		for i := 1; i < t.height; i++ {
			b.WriteString(emptyLine)
			if i < t.height-1 {
				b.WriteString("\n")
			}
		}
	} else {
		// Get display lines and visual cursor position
		displayLines, vRow, vCol := t.getDisplayInfo()

		// Adjust vertical offset if cursor is out of view
		if vRow != -1 {
			if vRow < t.offset {
				t.offset = vRow
			} else if vRow >= t.offset+t.height {
				t.offset = vRow - t.height + 1
			}
		}

		// Calculate which display lines to show
		startLine := t.offset
		endLine := min(startLine+t.height, len(displayLines))


		for i := startLine; i < endLine; i++ {
			line := displayLines[i]

			// Pad line to full width for consistent rendering
			if utf8.RuneCountInString(line) < t.width {
				line += strings.Repeat(" ", t.width-utf8.RuneCountInString(line))
			}

			// Render cursor if focused and on cursor row
			if t.focused && i == vRow {
				cursorStyle := t.styles.Base
				if t.blink {
					cursorStyle = cursorStyle.Reverse(true)
				}

				// Render line with cursor
				if vCol >= utf8.RuneCountInString(line) {
					b.WriteString(line)
					b.WriteString(cursorStyle.Render(" "))
				} else {
					lineRunes := []rune(line)
					before := string(lineRunes[:vCol])
					cursorChar := string(lineRunes[vCol])
					after := string(lineRunes[vCol+1:])
					b.WriteString(before)
					b.WriteString(cursorStyle.Render(cursorChar))
					b.WriteString(after)
				}
			} else {
				b.WriteString(line)
			}

			if i < endLine-1 {
				b.WriteString("\n")
			}
		}
		
		// If we haven't filled up to t.height lines (e.g. content is shorter than height), fill with empty lines
		// This ensures consistent height for the TextArea
		linesRendered := endLine - startLine
		if linesRendered < t.height {
			if linesRendered > 0 {
				b.WriteString("\n")
			}
			emptyLine := strings.Repeat(" ", t.width)
			for i := linesRendered; i < t.height; i++ {
				b.WriteString(emptyLine)
				if i < t.height-1 {
					b.WriteString("\n")
				}
			}
		}
	}

	// Render error if present
	if t.err != nil {
		b.WriteString("\n")
		b.WriteString(t.styles.Base.Foreground(t.styles.Error).Render(t.err.Error()))
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

// getDisplayInfo returns lines for display and the visual cursor coordinates
func (t *TextArea) getDisplayInfo() ([]string, int, int) {
	if !t.wrap {
		return t.lines, t.cursorRow, t.cursorCol
	}

	var wrapped []string
	vRow, vCol := -1, -1

	for r, line := range t.lines {
		runes := []rune(line)

		// Handle empty line
		if len(runes) == 0 {
			if r == t.cursorRow {
				vRow = len(wrapped)
				vCol = 0
			}
			wrapped = append(wrapped, "")
			continue
		}

		// Break line into chunks of width
		for i := 0; i < len(runes); i += t.width {
			end := min(i+t.width, len(runes))

			// Check if cursor falls in this chunk
			if r == t.cursorRow {
				isLastChunk := (end == len(runes))
				// Cursor belongs to this chunk if it's within [i, end)
				// OR if it's at 'end' and this is the last chunk
				if t.cursorCol >= i && t.cursorCol < end {
					vRow = len(wrapped)
					vCol = t.cursorCol - i
				} else if t.cursorCol == end && isLastChunk {
					vRow = len(wrapped)
					vCol = t.cursorCol - i
				}
			}

			wrapped = append(wrapped, string(runes[i:end]))
		}
	}

	return wrapped, vRow, vCol
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
	if t.cursorRow >= 0 {
		t.cursorCol = len(t.lines[t.cursorRow])
	}
}

// Focus focuses the text area.
func (t *TextArea) Focus() render.Cmd {
	t.focused = true
	t.blink = true
	t.blinkCtx = NextBlinkID()
	return BlinkCmd(t.blinkCtx)
}

// Blur blurs the text area.
func (t *TextArea) Blur() {
	t.focused = false
	t.blink = false
}

// Focused returns true if focused.
func (t *TextArea) Focused() bool {
	return t.focused
}

// Validate validates the input.
func (t *TextArea) Validate() error {
	for _, v := range t.validators {
		if err := v(t.Value()); err != nil {
			t.err = err
			return err
		}
	}
	t.err = nil
	return nil
}

// Error returns the last validation error.
func (t *TextArea) Error() error {
	return t.err
}

// AddValidator adds a validator.
func (t *TextArea) AddValidator(v Validator) {
	t.validators = append(t.validators, v)
}

// SetValidators sets the validators.
func (t *TextArea) SetValidators(v ...Validator) {
	t.validators = v
}

// Insert inserts a rune.
func (t *TextArea) Insert(r rune) {
	if len(t.lines) == 0 {
		t.lines = []string{""}
	}
	line := t.lines[t.cursorRow]
	runes := []rune(line)
	before := string(runes[:t.cursorCol])
	after := string(runes[t.cursorCol:])
	t.lines[t.cursorRow] = before + string(r) + after
	t.cursorCol++
}

// InsertNewline inserts a newline.
func (t *TextArea) InsertNewline() {
	if len(t.lines) == 0 {
		t.lines = []string{""}
	}
	line := t.lines[t.cursorRow]
	runes := []rune(line)
	before := string(runes[:t.cursorCol])
	after := string(runes[t.cursorCol:])

	t.lines[t.cursorRow] = before
	// Insert new line after current
	t.lines = append(t.lines[:t.cursorRow+1], append([]string{after}, t.lines[t.cursorRow+1:]...)...)

	t.cursorRow++
	t.cursorCol = 0
}

// Delete deletes the character before cursor.
func (t *TextArea) Delete() {
	if len(t.lines) == 0 {
		return
	}
	if t.cursorCol > 0 {
		line := t.lines[t.cursorRow]
		runes := []rune(line)
		before := string(runes[:t.cursorCol-1])
		after := string(runes[t.cursorCol:])
		t.lines[t.cursorRow] = before + after
		t.cursorCol--
	} else if t.cursorRow > 0 {
		// Merge with previous line
		prevLine := t.lines[t.cursorRow-1]
		currLine := t.lines[t.cursorRow]

		newCol := utf8.RuneCountInString(prevLine)
		t.lines[t.cursorRow-1] = prevLine + currLine

		// Remove current line
		t.lines = append(t.lines[:t.cursorRow], t.lines[t.cursorRow+1:]...)

		t.cursorRow--
		t.cursorCol = newCol
	}
}

// MoveUp moves cursor up.
func (t *TextArea) MoveUp() {
	_, vRow, vCol := t.getDisplayInfo()
	if vRow > 0 {
		t.cursorRow, t.cursorCol = t.resolveVisualPos(vRow-1, vCol)
	}
}

// MoveDown moves cursor down.
func (t *TextArea) MoveDown() {
	lines, vRow, vCol := t.getDisplayInfo()
	if vRow < len(lines)-1 {
		t.cursorRow, t.cursorCol = t.resolveVisualPos(vRow+1, vCol)
	}
}

// MoveLeft moves cursor left.
func (t *TextArea) MoveLeft() {
	if t.cursorCol > 0 {
		t.cursorCol--
	} else if t.cursorRow > 0 {
		t.cursorRow--
		t.cursorCol = utf8.RuneCountInString(t.lines[t.cursorRow])
	}
}

// MoveRight moves cursor right.
func (t *TextArea) MoveRight() {
	if t.cursorRow < len(t.lines) {
		lineLen := utf8.RuneCountInString(t.lines[t.cursorRow])
		if t.cursorCol < lineLen {
			t.cursorCol++
		} else if t.cursorRow < len(t.lines)-1 {
			t.cursorRow++
			t.cursorCol = 0
		}
	}
}

// resolveVisualPos resolves a visual position to a logical position.
func (t *TextArea) resolveVisualPos(targetVRow, targetVCol int) (int, int) {
	if !t.wrap {
		// Clamp to valid range
		if targetVRow < 0 {
			targetVRow = 0
		}
		if targetVRow >= len(t.lines) {
			targetVRow = len(t.lines) - 1
		}
		// Clamp column
		lineLen := utf8.RuneCountInString(t.lines[targetVRow])
		targetVCol = min(targetVCol, lineLen)
		return targetVRow, targetVCol
	}

	// For wrapped text, we need to find the logical position
	currentVRow := 0

	// Handle negative targetVRow (clamp to top)
	if targetVRow < 0 {
		return 0, 0
	}

	for lRow, line := range t.lines {
		runes := []rune(line)

		// Handle empty line
		if len(runes) == 0 {
			if currentVRow == targetVRow {
				return lRow, 0
			}
			currentVRow++
			continue
		}

		// Iterate chunks
		for i := 0; i < len(runes); i += t.width {
			end := min(i+t.width, len(runes))

			if currentVRow == targetVRow {
				// Found the visual row
				chunkLen := end - i

				// Clamp targetVCol to chunk length
				colInChunk := targetVCol
				colInChunk = min(colInChunk, chunkLen)

				// Logical col = start of chunk + offset
				return lRow, i + colInChunk
			}
			currentVRow++
		}
	}

	// If targetVRow is beyond the end, return the very last position
	lastRow := len(t.lines) - 1
	if lastRow < 0 {
		return 0, 0
	}
	lastCol := utf8.RuneCountInString(t.lines[lastRow])
	return lastRow, lastCol
}

// SetWidth sets the width.
func (t *TextArea) SetWidth(w int) {
	t.width = w
}

// SetHeight sets the height.
func (t *TextArea) SetHeight(h int) {
	t.height = h
}

// SetWrap enables or disables word wrap.
func (t *TextArea) SetWrap(wrap bool) {
	t.wrap = wrap
}

// SetShowBorder sets whether to show the border.
func (t *TextArea) SetShowBorder(show bool) {
	t.showBorder = show
}

// SetLabel sets the label displayed above the text area.
func (t *TextArea) SetLabel(label string) {
	t.label = label
}
