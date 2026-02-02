package buffer

import (
	"testing"
)

func TestNewBuffer(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		wantW    int
		wantH    int
		wantCell int
	}{
		{"default size", 0, 0, 80, 24, 80 * 24},
		{"negative size", -10, -20, 80, 24, 80 * 24},
		{"valid size", 40, 20, 40, 20, 40 * 20},
		{"small size", 5, 3, 5, 3, 5 * 3},
		{"square", 10, 10, 10, 10, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer(tt.width, tt.height)

			if b.width != tt.wantW {
				t.Errorf("NewBuffer() width = %d, want %d", b.width, tt.wantW)
			}
			if b.height != tt.wantH {
				t.Errorf("NewBuffer() height = %d, want %d", b.height, tt.wantH)
			}
			if len(b.cells) != tt.wantH {
				t.Errorf("NewBuffer() cells rows = %d, want %d", len(b.cells), tt.wantH)
			}
			if len(b.cells) > 0 && len(b.cells[0]) != tt.wantW {
				t.Errorf("NewBuffer() cells cols = %d, want %d", len(b.cells[0]), tt.wantW)
			}

			// Check all cells are initialized to space
			for y := 0; y < b.height; y++ {
				for x := 0; x < b.width; x++ {
					cell := b.cells[y][x]
					if cell.Char != ' ' {
						t.Errorf("NewBuffer() cell[%d][%d] = %q, want ' '", y, x, cell.Char)
					}
				}
			}
		})
	}
}

func TestSize(t *testing.T) {
	b := NewBuffer(50, 30)

	w, h := b.Size().Width, b.Size().Height
	if w != 50 {
		t.Errorf("Size() Width = %d, want 50", w)
	}
	if h != 30 {
		t.Errorf("Size() Height = %d, want 30", h)
	}
}

func TestWidthHeight(t *testing.T) {
	b := NewBuffer(100, 50)

	if b.Width() != 100 {
		t.Errorf("Width() = %d, want 100", b.Width())
	}
	if b.Height() != 50 {
		t.Errorf("Height() = %d, want 50", b.Height())
	}
}

func TestValid(t *testing.T) {
	b := NewBuffer(10, 10)

	tests := []struct {
		name string
		p    Point
		want bool
	}{
		{"valid point", Point{X: 5, Y: 5}, true},
		{"origin", Point{X: 0, Y: 0}, true},
		{"edge", Point{X: 9, Y: 9}, true},
		{"negative x", Point{X: -1, Y: 5}, false},
		{"negative y", Point{X: 5, Y: -1}, false},
		{"x too large", Point{X: 10, Y: 5}, false},
		{"y too large", Point{X: 5, Y: 10}, false},
		{"both too large", Point{X: 100, Y: 100}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := b.Valid(tt.p); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetCell(t *testing.T) {
	b := NewBuffer(10, 10)
	style := Style{Foreground: "196", Bold: true}
	cell := Cell{Char: 'X', Width: 1, Style: style}

	// Valid position
	if !b.SetCell(Point{X: 5, Y: 5}, cell) {
		t.Error("SetCell() returned false for valid position")
	}
	if b.cells[5][5].Char != 'X' {
		t.Errorf("SetCell() cell char = %q, want 'X'", b.cells[5][5].Char)
	}

	// Invalid position
	if b.SetCell(Point{X: 20, Y: 20}, cell) {
		t.Error("SetCell() returned true for invalid position")
	}
}

func TestFillRect(t *testing.T) {
	b := NewBuffer(20, 20)
	style := Style{Foreground: "33"}
	runeChar := '█'

	// Normal rect
	b.FillRect(Rect{X: 5, Y: 5, Width: 10, Height: 5}, runeChar, style)

	// Check filled area
	for y := 5; y < 10; y++ {
		for x := 5; x < 15; x++ {
			if b.cells[y][x].Char != runeChar {
				t.Errorf("FillRect() cell[%d][%d] = %q, want %q", y, x, b.cells[y][x].Char, runeChar)
			}
		}
	}

	// Check unfilled area
	if b.cells[4][5].Char != ' ' {
		t.Errorf("FillRect() affected cell outside rect")
	}
	if b.cells[15][5].Char != ' ' {
		t.Errorf("FillRect() affected cell outside rect")
	}

	// Rect crossing bounds
	b = NewBuffer(10, 10)
	b.FillRect(Rect{X: 5, Y: 5, Width: 20, Height: 20}, '●', Style{})

	// Should only fill visible area
	if b.cells[9][9].Char != '●' {
		t.Errorf("FillRect() didn't fill boundary cell")
	}
}

func TestWriteString(t *testing.T) {
	style := Style{Foreground: "82"}

	tests := []struct {
		name      string
		p         Point
		text      string
		wantCols  int
		wantChars []string
	}{
		{
			name:      "simple text",
			p:         Point{X: 2, Y: 2},
			text:      "hello",
			wantCols:  5,
			wantChars: []string{"h", "e", "l", "l", "o"},
		},
		{
			name:      "text at edge",
			p:         Point{X: 15, Y: 2},
			text:      "test",
			wantCols:  4,
			wantChars: []string{"t", "e", "s", "t"},
		},
		{
			name:      "text beyond edge",
			p:         Point{X: 18, Y: 2},
			text:      "toolong",
			wantCols:  2,
			wantChars: []string{"t", "o"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b2 := NewBuffer(20, 5)
			cols := b2.WriteString(tt.p, tt.text, style)

			if cols != tt.wantCols {
				t.Errorf("WriteString() cols = %d, want %d", cols, tt.wantCols)
			}

			for i, wantChar := range tt.wantChars {
				if b2.cells[tt.p.Y][tt.p.X+i].Char != rune(wantChar[0]) {
					t.Errorf("WriteString() cell char = %q, want %q", b2.cells[tt.p.Y][tt.p.X+i].Char, wantChar)
				}
			}
		})
	}
}

func TestWriteStringWideChars(t *testing.T) {
	b := NewBuffer(20, 5)
	style := Style{}

	// English text
	cols := b.WriteString(Point{X: 0, Y: 0}, "hello", style)
	if cols != 5 {
		t.Errorf("WriteString() cols for English = %d, want 5", cols)
	}

	// Chinese text (wide chars)
	b = NewBuffer(20, 5)
	cols = b.WriteString(Point{X: 0, Y: 0}, "你好世界", style)
	if cols != 8 { // Each Chinese char is 2 columns
		t.Errorf("WriteString() cols for Chinese = %d, want 8", cols)
	}

	// Mixed text
	b = NewBuffer(20, 5)
	cols = b.WriteString(Point{X: 0, Y: 0}, "Hello世界", style)
	if cols != 9 { // 5 + 4
		t.Errorf("WriteString() cols for mixed = %d, want 9", cols)
	}
}

func TestWriteStringWrapped(t *testing.T) {
	style := Style{}

	tests := []struct {
		name      string
		text      string
		maxWidth  int
		wantLines int
	}{
		{
			name:      "short text",
			text:      "hello world",
			maxWidth:  20,
			wantLines: 1,
		},
		{
			name:      "text wrapping",
			text:      "this is a very long text that needs wrapping",
			maxWidth:  10,
			wantLines: 5,
		},
		{
			name:      "newlines",
			text:      "line1\nline2\nline3",
			maxWidth:  20,
			wantLines: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b2 := NewBuffer(20, 10)
			lines := b2.WriteStringWrapped(Point{X: 0, Y: 0}, tt.maxWidth, tt.text, style)

			if lines != tt.wantLines {
				t.Errorf("WriteStringWrapped() lines = %d, want %d", lines, tt.wantLines)
			}
		})
	}
}

func TestWriteBuffer(t *testing.T) {
	b1 := NewBuffer(20, 20)
	b2 := NewBuffer(5, 5)

	// Fill b2 with pattern
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			b2.cells[y][x] = Cell{Char: '█', Width: 1}
		}
	}

	// Write b2 into b1
	if !b1.WriteBuffer(Point{X: 5, Y: 5}, b2) {
		t.Error("WriteBuffer() returned false")
	}

	// Check content was written
	if b1.cells[5][5].Char != '█' {
		t.Errorf("WriteBuffer() didn't write content")
	}
	if b1.cells[9][9].Char != '█' {
		t.Errorf("WriteBuffer() didn't write all content")
	}

	// Check boundary
	if b1.cells[4][5].Char != ' ' {
		t.Errorf("WriteBuffer() wrote outside bounds")
	}

	// Nil buffer
	if b1.WriteBuffer(Point{X: 0, Y: 0}, nil) {
		t.Error("WriteBuffer() returned true for nil buffer")
	}
}

func TestRender(t *testing.T) {
	b := NewBuffer(10, 5)

	// Write some content
	style := Style{Foreground: "196"}
	b.WriteString(Point{X: 2, Y: 2}, "hello", style)

	output := b.Render()

	// Check has 5 lines (4 newlines)
	newlineCount := 0
	for _, c := range output {
		if c == '\n' {
			newlineCount++
		}
	}
	if newlineCount != 4 {
		t.Errorf("Render() output has %d newlines, want 4", newlineCount)
	}
}

func TestTextComponent(t *testing.T) {
	content := "Hello\nWorld"
	style := Style{Foreground: "82"}
	tc := NewTextComponent(content, style)

	// Test MinSize
	w, h := tc.MinSize()
	if w != 5 || h != 2 {
		t.Errorf("TextComponent.MinSize() = (%d, %d), want (5, 2)", w, h)
	}

	// Test PreferredSize
	w, h = tc.PreferredSize()
	if w != 5 || h != 2 {
		t.Errorf("TextComponent.PreferredSize() = (%d, %d), want (5, 2)", w, h)
	}

	// Test Render
	buf := NewBuffer(10, 10)
	tc.Render(buf, Rect{X: 2, Y: 2, Width: 10, Height: 10})

	if buf.cells[2][2].Char != 'H' {
		t.Errorf("TextComponent.Render() didn't render correctly")
	}
	if buf.cells[3][2].Char != 'W' {
		t.Errorf("TextComponent.Render() didn't render second line")
	}

	// Test centering
	tc = NewTextComponent("test", style).SetCenterH(true)
	buf = NewBuffer(20, 5)
	tc.Render(buf, Rect{X: 0, Y: 0, Width: 20, Height: 5})

	// Should be centered
	if buf.cells[0][8].Char != 't' {
		t.Errorf("TextComponent.SetCenterH() didn't center text")
	}
}

func TestImageComponent(t *testing.T) {
	ic := NewImageComponent(100, 200)

	// Test MinSize - should be at least 10x5, so returns (100, 200)
	w, h := ic.MinSize()
	if w != 100 || h != 200 {
		t.Errorf("ImageComponent.MinSize() = (%d, %d), want (100, 200)", w, h)
	}

	// Test MinSize with small values
	ic2 := NewImageComponent(3, 2)
	w, h = ic2.MinSize()
	if w != 10 || h != 5 {
		t.Errorf("ImageComponent.MinSize() for small image = (%d, %d), want (10, 5)", w, h)
	}

	// Test PreferredSize
	ic = NewImageComponent(100, 200)
	w, h = ic.PreferredSize()
	if w != 100 || h != 200 {
		t.Errorf("ImageComponent.PreferredSize() = (%d, %d), want (100, 200)", w, h)
	}

	// Test Render
	buf := NewBuffer(20, 15)
	ic.Render(buf, Rect{X: 2, Y: 2, Width: 16, Height: 11})

	// Should have border
	if buf.cells[2][2].Char != '┌' {
		t.Errorf("ImageComponent.Render() didn't draw top-left corner")
	}
	if buf.cells[2][17].Char != '┐' {
		t.Errorf("ImageComponent.Render() didn't draw top-right corner")
	}
	if buf.cells[12][2].Char != '└' {
		t.Errorf("ImageComponent.Render() didn't draw bottom-left corner")
	}
	if buf.cells[12][17].Char != '┘' {
		t.Errorf("ImageComponent.Render() didn't draw bottom-right corner")
	}
}

func TestFillComponent(t *testing.T) {
	style := Style{Foreground: "33"}
	fc := NewFillComponent('█', style)

	buf := NewBuffer(10, 10)
	fc.Render(buf, Rect{X: 2, Y: 2, Width: 6, Height: 6})

	// Should fill the rect
	for y := 2; y < 8; y++ {
		for x := 2; x < 8; x++ {
			if buf.cells[y][x].Char != '█' {
				t.Errorf("FillComponent.Render() didn't fill cell[%d][%d]", y, x)
			}
		}
	}

	// Should not fill outside rect
	if buf.cells[1][2].Char != ' ' {
		t.Errorf("FillComponent.Render() filled outside rect")
	}
	if buf.cells[2][1].Char != ' ' {
		t.Errorf("FillComponent.Render() filled outside rect")
	}
}

func TestLayoutManager(t *testing.T) {
	lm := NewLayoutManager(80, 24)

	// Test initial state
	if lm.width != 80 || lm.height != 24 {
		t.Errorf("NewLayoutManager() size = (%d, %d), want (80, 24)", lm.width, lm.height)
	}

	// Test AddComponent
	tc := NewTextComponent("test", Style{})
	lm.AddComponent("content", tc)

	if len(lm.components) != 1 {
		t.Errorf("AddComponent() didn't add component")
	}

	// Test CalculateLayout
	lm.CalculateLayout()

	if len(lm.layouts) != 3 {
		t.Errorf("CalculateLayout() didn't create 3 layouts")
	}

	// Check header layout
	headerRect := lm.layouts["header"]
	if headerRect == nil {
		t.Error("CalculateLayout() didn't create header layout")
	} else if headerRect.Height != 5 {
		t.Errorf("CalculateLayout() header height = %d, want 5", headerRect.Height)
	}

	// Check footer layout
	footerRect := lm.layouts["footer"]
	if footerRect == nil {
		t.Error("CalculateLayout() didn't create footer layout")
	} else if footerRect.Y != 23 {
		t.Errorf("CalculateLayout() footer Y = %d, want 23", footerRect.Y)
	}
}

func TestLayoutManagerImageLayout(t *testing.T) {
	lm := NewLayoutManager(80, 24)

	// Test with height hint
	lm.ImageLayout(10)

	contentRect := lm.layouts["content"]
	if contentRect == nil {
		t.Error("ImageLayout() didn't create content layout")
	} else if contentRect.Height != 10 {
		t.Errorf("ImageLayout() content height = %d, want 10", contentRect.Height)
	}

	// Content should be vertically centered
	if contentRect.Y < 5 {
		t.Errorf("ImageLayout() content not centered: Y = %d", contentRect.Y)
	}
}

func TestLayoutManagerRender(t *testing.T) {
	lm := NewLayoutManager(80, 24)

	// Add components
	header := NewTextComponent("HEADER", Style{Bold: true}).SetCenterH(true)
	content := NewTextComponent("Content", Style{})
	content.SetWrap(true)
	footer := NewTextComponent("Footer", Style{}).SetCenterH(true)

	lm.AddComponent("header", header)
	lm.AddComponent("content", content)
	lm.AddComponent("footer", footer)

	// Calculate layout
	lm.CalculateLayout()

	// Render
	output := lm.Render()

	if len(output) == 0 {
		t.Error("Render() returned empty string")
	}

	// Should have newlines
	newlineCount := 0
	for _, c := range output {
		if c == '\n' {
			newlineCount++
		}
	}
	if newlineCount != 23 {
		t.Errorf("Render() has %d newlines, want 23", newlineCount)
	}
}

func TestLayoutManagerSetSize(t *testing.T) {
	lm := NewLayoutManager(80, 24)

	lm.SetSize(100, 50)

	if lm.width != 100 || lm.height != 50 {
		t.Errorf("SetSize() size = (%d, %d), want (100, 50)", lm.width, lm.height)
	}

	if lm.mainBuf.Width() != 100 || lm.mainBuf.Height() != 50 {
		t.Errorf("SetSize() didn't recreate buffer")
	}
}

// Performance benchmarks

func BenchmarkFillRect(b *testing.B) {
	style := Style{}
	buf := NewBuffer(100, 100)
	rect := Rect{X: 10, Y: 10, Width: 80, Height: 80}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.FillRect(rect, '█', style)
	}
}

func BenchmarkWriteString(b *testing.B) {
	style := Style{}
	buf := NewBuffer(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.WriteString(Point{X: 0, Y: 0}, "The quick brown fox jumps over the lazy dog", style)
	}
}

func BenchmarkWriteStringWrapped(b *testing.B) {
	style := Style{}
	buf := NewBuffer(80, 24)
	text := "This is a long text that needs to be wrapped across multiple lines for testing purposes"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.WriteStringWrapped(Point{X: 0, Y: 0}, 40, text, style)
	}
}

func BenchmarkWriteBuffer(b *testing.B) {
	source := NewBuffer(20, 10)
	dest := NewBuffer(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dest.WriteBuffer(Point{X: 10, Y: 10}, source)
	}
}

func BenchmarkRender(b *testing.B) {
	buf := NewBuffer(80, 24)
	style := Style{Foreground: "82"}
	buf.WriteString(Point{X: 10, Y: 5}, "Test Content", style)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Render()
	}
}

func BenchmarkTextComponentRender(b *testing.B) {
	tc := NewTextComponent("Test content for rendering performance", Style{})
	tc.SetWrap(true)
	buf := NewBuffer(80, 24)
	rect := Rect{X: 0, Y: 0, Width: 80, Height: 24}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc.Render(buf, rect)
	}
}

func BenchmarkLayoutCalculate(b *testing.B) {
	lm := NewLayoutManager(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lm.CalculateLayout()
	}
}

func BenchmarkLayoutRender(b *testing.B) {
	lm := NewLayoutManager(80, 24)

	header := NewTextComponent("HEADER", Style{Bold: true}).SetCenterH(true)
	content := NewTextComponent("Content", Style{})
	content.SetWrap(true)
	footer := NewTextComponent("Footer", Style{}).SetCenterH(true)

	lm.AddComponent("header", header)
	lm.AddComponent("content", content)
	lm.AddComponent("footer", footer)
	lm.CalculateLayout()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lm.Render()
	}
}
