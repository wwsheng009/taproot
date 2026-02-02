package buffer

import (
	"strings"
)

// Point represents a coordinate in the buffer
type Point struct {
	X int
	Y int
}

// Size represents dimensions
type Size struct {
	Width  int
	Height int
}

// Rect represents a rectangular area
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Style represents cell styling
type Style struct {
	Foreground string
	Background string
	Bold       bool
	Italic     bool
	Underline  bool
	Reverse    bool
}

// Cell represents a single character cell in the buffer
// For wide characters (2 columns), the first cell has Width=2 and IsContinuation=false,
// and the second cell has Width=0 or 1 and IsContinuation=true
type Cell struct {
	Char           rune
	Width          int
	Style          Style
	IsContinuation bool // true if this is the second cell of a wide character
}

// Buffer represents a 2D grid of cells
type Buffer struct {
	width  int
	height int
	cells  [][]Cell
}

// NewBuffer creates a new buffer with specified dimensions
func NewBuffer(width, height int) *Buffer {
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	b := &Buffer{
		width:  width,
		height: height,
		cells:  make([][]Cell, height),
	}

	// Initialize all cells
	for y := 0; y < height; y++ {
		b.cells[y] = make([]Cell, width)
		for x := 0; x < width; x++ {
			b.cells[y][x] = Cell{Char: ' ', Width: 1}
		}
	}

	return b
}

// Size returns the buffer dimensions
func (b *Buffer) Size() Size {
	return Size{Width: b.width, Height: b.height}
}

// Width returns the buffer width
func (b *Buffer) Width() int {
	return b.width
}

// Height returns the buffer height
func (b *Buffer) Height() int {
	return b.height
}

// Valid checks if a point is within buffer bounds
func (b *Buffer) Valid(p Point) bool {
	return p.X >= 0 && p.X < b.width && p.Y >= 0 && p.Y < b.height
}

// SetCell sets a cell at the given position
func (b *Buffer) SetCell(p Point, cell Cell) bool {
	if !b.Valid(p) {
		return false
	}
	b.clearCellAt(p.X, p.Y)
	b.cells[p.Y][p.X] = cell
	return true
}

// clearCellAt clears a cell and handles wide character continuation
// This prevents "ghost characters" when overwriting wide characters
func (b *Buffer) clearCellAt(x, y int) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return
	}

	cell := b.cells[y][x]

	// If this is a continuation cell, we need to clear the head (left cell) too
	if cell.IsContinuation && x > 0 {
		head := b.cells[y][x-1]
		if head.Width == 2 {
			b.cells[y][x-1] = Cell{Char: ' ', Width: 1}
		}
	}

	// If this is a wide character head, we need to clear the continuation (right cell) too
	if cell.Width == 2 && x+1 < b.width {
		b.cells[y][x+1] = Cell{Char: ' ', Width: 1, IsContinuation: false}
	}

	// Clear the current cell
	b.cells[y][x] = Cell{Char: ' ', Width: 1, IsContinuation: false}
}

// FillRect fills a rectangular area with a character
func (b *Buffer) FillRect(r Rect, _char rune, style Style) {
	// Clamp rect to buffer bounds
	if r.X < 0 {
		r.X = 0
	}
	if r.Y < 0 {
		r.Y = 0
	}
	if r.X+r.Width > b.width {
		r.Width = b.width - r.X
	}
	if r.Y+r.Height > b.height {
		r.Height = b.height - r.Y
	}

	// Fill the area
	for y := r.Y; y < r.Y+r.Height; y++ {
		for x := r.X; x < r.X+r.Width; x++ {
			b.clearCellAt(x, y)
			b.cells[y][x] = Cell{
				Char:           _char,
				Width:          1,
				Style:          style,
				IsContinuation: false,
			}
		}
	}
}

// WriteString writes a string at the given position
// Returns the number of columns used
func (b *Buffer) WriteString(p Point, text string, style Style) int {
	if !b.Valid(p) {
		return 0
	}

	x := p.X
	y := p.Y
	colsUsed := 0
	textRunes := []rune(text)

	for _, r := range textRunes {
		if x >= b.width {
			break
		}

		width := 1
		if isWideChar(r) {
			width = 2
			// Check if there's enough space for wide char
			if x+1 >= b.width {
				break
			}
		}

		// Clear the cell first to handle continuation cells
		b.clearCellAt(x, y)

		// Set the head cell
		b.cells[y][x] = Cell{
			Char:           r,
			Width:          width,
			Style:          style,
			IsContinuation: false,
		}

		// If wide character, set the continuation cell
		if width == 2 {
			b.clearCellAt(x+1, y)
			b.cells[y][x+1] = Cell{
				Char:           0,
				Width:          0,
				Style:          style,
				IsContinuation: true,
			}
		}

		x += width
		colsUsed += width
	}

	return colsUsed
}

// WriteStringWrapped writes a string with word wrapping
// Returns the number of lines used
func (b *Buffer) WriteStringWrapped(p Point, maxWidth int, text string, style Style) int {
	if !b.Valid(p) {
		return 0
	}

	if maxWidth <= 0 {
		maxWidth = b.width - p.X
	}

	x := p.X
	y := p.Y
	linesUsed := 1
	wordStart := 0
	wordCols := 0

	for i, r := range text {
		width := 1
		if isWideChar(r) {
			width = 2
		}

		if r == ' ' || r == '\t' || r == '\n' {
			// Write the word
			if x+wordCols > p.X+maxWidth {
				// Word doesn't fit, move to next line
				x = p.X
				y++
				linesUsed++
				if y >= b.height {
					break
				}
			}

			// Write word content
			wordRunes := []rune(text[wordStart:i])
			for _, r := range wordRunes {
				if x >= b.width {
					break
				}
				width := cellWidthForRune(r)
				b.clearCellAt(x, y)
				b.cells[y][x] = Cell{
					Char:           r,
					Width:          width,
					Style:          style,
					IsContinuation: false,
				}
				// If wide character, set the continuation cell
				if width == 2 {
					if x+1 < b.width {
						b.clearCellAt(x+1, y)
						b.cells[y][x+1] = Cell{
							Char:           0,
							Width:          0,
							Style:          style,
							IsContinuation: true,
						}
					}
				}
				x += width
			}

			if r == '\n' {
				// Move to next line
				x = p.X
				y++
				linesUsed++
				if y >= b.height {
					break
				}
			} else {
				// Write space
				if x < b.width {
					b.cells[y][x] = Cell{
						Char:  ' ',
						Width: 1,
						Style: style,
					}
					x++
				}
			}

			wordStart = i + 1
			wordCols = 0
		} else {
			wordCols += width
		}
	}

	// Write remaining word
	if wordStart < len(text) && y < b.height {
		if x+wordCols > p.X+maxWidth {
			x = p.X
			y++
			linesUsed++
		}

		wordRunes := []rune(text[wordStart:])
		for _, r := range wordRunes {
			if x >= b.width {
				break
			}
			width := cellWidthForRune(r)
			b.clearCellAt(x, y)
			b.cells[y][x] = Cell{
				Char:           r,
				Width:          width,
				Style:          style,
				IsContinuation: false,
			}
			// If wide character, set the continuation cell
			if width == 2 {
				if x+1 < b.width {
					b.clearCellAt(x+1, y)
					b.cells[y][x+1] = Cell{
						Char:           0,
						Width:          0,
						Style:          style,
						IsContinuation: true,
					}
				}
			}
			x += width
		}
	}

	return linesUsed
}

// WriteBuffer writes another buffer's content into this buffer
// Returns true if successful
func (b *Buffer) WriteBuffer(p Point, other *Buffer) bool {
	if other == nil {
		return false
	}

	for y := 0; y < other.height; y++ {
		for x := 0; x < other.width; x++ {
			targetX := p.X + x
			targetY := p.Y + y

			if targetX < b.width && targetY < b.height {
				b.cells[targetY][targetX] = other.cells[y][x]
			}
		}
	}

	return true
}

// Render converts the buffer to a string for display
func (b *Buffer) Render() string {
	output := GetStringBuilder()
	defer PutStringBuilder(output)

	output.Grow(b.width * b.height * 2) // Pre-allocate memory

	for y := 0; y < b.height; y++ {
		b.renderLineToBuilder(y, output)
		if y < b.height-1 {
			output.WriteString("\n")
		}
	}

	return output.String()
}

// renderLineToBuilder renders a single line directly to strings.Builder (optimized)
func (b *Buffer) renderLineToBuilder(y int, output *strings.Builder) {
	x := 0
	var lastStyleStr string

	for x < b.width {
		cell := b.cells[y][x]

		// Skip continuation cells (second part of wide characters)
		if cell.IsContinuation {
			x++
			continue
		}

		// Skip empty cells
		if cell.Char == ' ' && cell.Style.Foreground == "" {
			x++
			continue
		}

		// Use cached style string
		styleStr := globalStyleCache.Get(cell.Style)

		// Only write reset code if style changed
		if styleStr != lastStyleStr {
			if lastStyleStr != "" {
				output.WriteString("\x1b[0m")
			}
			if styleStr != "" {
				output.WriteString(styleStr)
			}
			lastStyleStr = styleStr
		}

		output.WriteRune(cell.Char)
		x += cell.Width
	}

	// Reset style at end of line
	if lastStyleStr != "" {
		output.WriteString("\x1b[0m")
	}
}

// isWideChar checks if a rune should be displayed with double width
func isWideChar(r rune) bool {
	// Simple heuristic: CJK characters are wide
	return r >= 0x1100 &&
		(r <= 0x115F ||
			r == 0x2329 ||
			r == 0x232A ||
			(r >= 0x2E80 && r <= 0xA4CF && r != 0x303F) ||
			(r >= 0xAC00 && r <= 0xD7A3) ||
			(r >= 0xF900 && r <= 0xFAFF) ||
			(r >= 0xFE10 && r <= 0xFE19) ||
			(r >= 0xFE30 && r <= 0xFE6F) ||
			(r >= 0xFF00 && r <= 0xFF60) ||
			(r >= 0xFFE0 && r <= 0xFFE6) ||
			(r >= 0x20000 && r <= 0x2FFFD) ||
			(r >= 0x30000 && r <= 0x3FFFD))
}

// cellWidthForRune returns the display width of a rune
func cellWidthForRune(r rune) int {
	if isWideChar(r) {
		return 2
	}
	return 1
}
