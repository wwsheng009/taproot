package buffer

import (
	"fmt"
	"strings"
)

// Renderable is an interface for components that can render to a buffer
type Renderable interface {
	// Render renders the component to the given buffer within the specified rect
	Render(buf *Buffer, Rect Rect)
	// MinSize returns the minimum size required (width, height)
	MinSize() (int, int)
	// PreferredSize returns the preferred size (width, height)
	PreferredSize() (int, int)
}

// TextComponent renders text content to a buffer
type TextComponent struct {
	content string
	style   Style
	wrap    bool
	centerV bool
	centerH bool
}

// NewTextComponent creates a new text component
func NewTextComponent(content string, style Style) *TextComponent {
	return &TextComponent{
		content: content,
		style:   style,
		wrap:    false,
		centerV: false,
		centerH: false,
	}
}

// Render renders the text to the buffer
func (tc *TextComponent) Render(buf *Buffer, rect Rect) {
	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}

	lines := strings.Split(tc.content, "\n")
	lineCount := len(lines)

	// Calculate vertical padding if centering
	paddingTop := 0
	if tc.centerV && lineCount < rect.Height {
		paddingTop = (rect.Height - lineCount) / 2
	}

	// Render each line
	for i, line := range lines {
		y := rect.Y + paddingTop + i
		if y >= rect.Y+rect.Height {
			break
		}

		if tc.wrap {
			buf.WriteStringWrapped(Point{X: rect.X, Y: y}, rect.Width, line, tc.style)
		} else {
			if tc.centerH {
				lineWidth := stringWidth(line)
				if lineWidth < rect.Width {
					padding := (rect.Width - lineWidth) / 2
					buf.WriteString(Point{X: rect.X + padding, Y: y}, line, tc.style)
				} else {
					buf.WriteString(Point{X: rect.X, Y: y}, truncateString(line, rect.Width), tc.style)
				}
			} else {
				buf.WriteString(Point{X: rect.X, Y: y}, line, tc.style)
			}
		}
	}
}

// MinSize returns the minimum size
func (tc *TextComponent) MinSize() (int, int) {
	lines := strings.Split(tc.content, "\n")
	maxWidth := 0
	for _, line := range lines {
		width := stringWidth(line)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth, len(lines)
}

// PreferredSize returns the preferred size
func (tc *TextComponent) PreferredSize() (int, int) {
	return tc.MinSize()
}

// SetWrap enables/disables word wrapping
func (tc *TextComponent) SetWrap(wrap bool) *TextComponent {
	tc.wrap = wrap
	return tc
}

// SetCenterV enables/disables vertical centering
func (tc *TextComponent) SetCenterV(center bool) *TextComponent {
	tc.centerV = center
	return tc
}

// SetCenterH enables/disables horizontal centering
func (tc *TextComponent) SetCenterH(center bool) *TextComponent {
	tc.centerH = center
	return tc
}

// ImageComponent renders a placeholder for images
type ImageComponent struct {
	width   int
	height  int
	bgChar  rune
	bgStyle Style
}

// NewImageComponent creates a new image component
func NewImageComponent(width, height int) *ImageComponent {
	return &ImageComponent{
		width:   width,
		height:  height,
		bgChar:  '·',
		bgStyle: Style{Foreground: "244"},
	}
}

// Render renders the image placeholder
func (ic *ImageComponent) Render(buf *Buffer, rect Rect) {
	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}

	// Fill with placeholder characters
	buf.FillRect(rect, ic.bgChar, ic.bgStyle)

	// Draw border
	borderStyle := Style{Foreground: "245"}
	borderChars := []rune{'┌', '┐', '└', '┘', '─', '│'}

	// Top-left
	buf.SetCell(Point{X: rect.X, Y: rect.Y}, Cell{Char: borderChars[0], Style: borderStyle})
	// Top-right
	if rect.Width > 1 {
		buf.SetCell(Point{X: rect.X + rect.Width - 1, Y: rect.Y}, Cell{Char: borderChars[1], Style: borderStyle})
	}
	// Bottom-left
	if rect.Height > 1 {
		buf.SetCell(Point{X: rect.X, Y: rect.Y + rect.Height - 1}, Cell{Char: borderChars[2], Style: borderStyle})
	}
	// Bottom-right
	if rect.Width > 1 && rect.Height > 1 {
		buf.SetCell(Point{X: rect.X + rect.Width - 1, Y: rect.Y + rect.Height - 1}, Cell{Char: borderChars[3], Style: borderStyle})
	}

	// Top and bottom borders
	for x := rect.X + 1; x < rect.X+rect.Width-1; x++ {
		buf.SetCell(Point{X: x, Y: rect.Y}, Cell{Char: borderChars[4], Style: borderStyle})
		if rect.Height > 1 {
			buf.SetCell(Point{X: x, Y: rect.Y + rect.Height - 1}, Cell{Char: borderChars[4], Style: borderStyle})
		}
	}

	// Left and right borders
	for y := rect.Y + 1; y < rect.Y+rect.Height-1; y++ {
		buf.SetCell(Point{X: rect.X, Y: y}, Cell{Char: borderChars[5], Style: borderStyle})
		if rect.Width > 1 {
			buf.SetCell(Point{X: rect.X + rect.Width - 1, Y: y}, Cell{Char: borderChars[5], Style: borderStyle})
		}
	}

	// Draw label
	if rect.Width > 20 && rect.Height > 3 {
		label := fmt.Sprintf("%dx%d", ic.width, ic.height)
		labelStyle := Style{Foreground: "86", Bold: true}
		labelX := rect.X + (rect.Width-len(label))/2
		labelY := rect.Y + rect.Height/2
		for i, r := range label {
			buf.SetCell(Point{X: labelX + i, Y: labelY}, Cell{Char: r, Style: labelStyle})
		}
	}
}

// MinSize returns the minimum size
func (ic *ImageComponent) MinSize() (int, int) {
	minW := ic.width
	minH := ic.height
	if minW < 10 {
		minW = 10
	}
	if minH < 5 {
		minH = 5
	}
	return minW, minH
}

// PreferredSize returns the preferred size
func (ic *ImageComponent) PreferredSize() (int, int) {
	return ic.width, ic.height
}

// FillComponent fills an area with a single character/style
type FillComponent struct {
	char  rune
	style Style
}

// NewFillComponent creates a new fill component
func NewFillComponent(char rune, style Style) *FillComponent {
	return &FillComponent{
		char:  char,
		style: style,
	}
}

// Render fills the entire rect with the character
func (fc *FillComponent) Render(buf *Buffer, rect Rect) {
	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}
	buf.FillRect(rect, fc.char, fc.style)
}

// MinSize returns the minimum size
func (fc *FillComponent) MinSize() (int, int) {
	return 1, 1
}

// PreferredSize returns the preferred size
func (fc *FillComponent) PreferredSize() (int, int) {
	return 1, 1
}

// Helper functions

func stringWidth(s string) int {
	width := 0
	for _, r := range s {
		if isWideChar(r) {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

func truncateString(s string, maxWidth int) string {
	r := []rune(s)
	width := 0
	for i, c := range r {
		cw := 1
		if isWideChar(c) {
			cw = 2
		}
		if width+cw > maxWidth {
			return string(r[:i])
		}
		width += cw
	}
	return s
}
