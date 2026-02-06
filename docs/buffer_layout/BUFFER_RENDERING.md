# Buffer-Based Rendering System - Taproot

## Overview

Buffer-based rendering is a new TUI rendering approach for Taproot that provides **accurate layout calculations** by working with actual cell coordinates instead of string estimation.

## Problem with String-Based Layout

Traditional TUI frameworks (like Bubbletea) render to strings and then estimate dimensions:

```go
// STRING-BASED APPROACH (PROBLEMATIC)
header := "Header\nLine 2\nLine 3\nLine 4\nLine 5"
content := sixelImage + "\n"
footer := "Footer"
output := header + "\n" + content + footer

// PROBLEM: Can't accurately calculate:
// - Sixel display height vs actual lines (6 pixels = ~1 line)
// - Word wrapping line count
// - Centering padding
// - Component dimensions
```

**Result**: Unreliable layout calculations, guesswork with image heights, and broken layouts on resize.

## Buffer-Based Solution

```go
// BUFFER-BASED APPROACH (RELIABLE)
// Pass 1: Calculate Layout
lm := NewLayoutManager(width, height)
lm.ImageLayout(displayHeight)  // Plan where everything goes

// Pass 2: Render to Buffer
mainBuf := NewBuffer(width, height)
componentBuf := NewBuffer(rect.Width, rect.Height)
component.Render(componentBuf, rect)  // Component knows its exact size
mainBuf.WriteBuffer(point, componentBuf)  // Copy to correct position

// Pass 3: Output String
output := mainBuf.Render()
```

## Key Advantages

### 1. **ACCURATE DIMENSIONS** ✅
```go
// Buffer grid guarantees exact cell counts
buf := NewBuffer(80, 30)  // → Exactly 80×30 = 2400 cells
width := buf.Width()      // → 80 (exact, not estimated)
height := buf.Height()    // → 30 (exact, not estimated)
```

### 2. **LAYOUT INDEPENDENCE** ✅
```go
// Calculate layout BEFORE rendering content
lm.CalculateLayout()

headerRect := lm.layouts["header"]
// → Rect{X: 0, Y: 0, Width: 80, Height: 5} (exact!)
// Content height doesn't affect layout calculation
```

### 3. **NO STRING CALCULATION HELL** ✅
```go
// String-based:
lines = strings.Count(output, "\n") + 1  // ❌ WRONG for Sixel images

// Buffer-based:
height := buf.Height()  // ✅ ACCURATE!
```

### 4. **COMPONENT ISOLATION** ✅
```go
// Each component has its own buffer
componentBuf := NewBuffer(componentWidth, componentHeight)
component.Render(componentBuf, rect)

// Write to main buffer at exact position
mainBuf.WriteBuffer(Point{X: 0, Y: 5}, componentBuf)
// → Components don't interfere with each other
```

### 5. **ACCURATE SIXEL SUPPORT** ✅
```go
// Sixel display height: NO GUESSING!
// Image: 300×600 pixels
// Terminal: 80×30 cells
// Display: 10 columns × 100 lines (Sixel)
// LayoutManager handles vertical positioning:

lm.ImageLayout(10)  // Exact display height in lines!
// → Content centered perfectly between header and footer
```

## Architecture

### Core Types

```go
// Point represents a coordinate in the buffer
type Point struct {
    X int
    Y int
}

// Rect represents a rectangular area
type Rect struct {
    X      int
    Y      int
    Width  int
    Height int
}

// Cell represents a single character cell
type Cell struct {
    Char  rune
    Width int  // 1 for normal, 2 for CJK characters
    Style Style
}

// Buffer represents a 2D grid of cells
type Buffer struct {
    width  int
    height int
    cells  [][]Cell
}
```

### Core Methods

```go
// Buffer operations
NewBuffer(width, height int) *Buffer
SetCell(p Point, cell Cell) bool
FillRect(r Rect, char rune, style Style)
WriteString(p Point, text string, style Style) int
WriteStringWrapped(p Point, maxWidth int, text string, style Style) int
WriteBuffer(p Point, other *Buffer) bool
Render() string  // Convert to output string
```

### Components

```go
// All components implement Renderable interface
type Renderable interface {
    Render(buf *Buffer, rect Rect)
    MinSize() (int, int)
    PreferredSize() (int, int)
}

// Built-in components
TextComponent   // Text with wrap/center
ImageComponent  // Image placeholder with borders
FillComponent   // Fill area with character/style
```

### Layout Manager

```go
// LayoutManager orchestrates component layout
lm := NewLayoutManager(width, height)

// Calculate layout
lm.CalculateLayout()           // Standard header/content/footer
lm.ImageLayout(displayHeight)  // For image viewers

// Add components
lm.AddComponent("header", header)
lm.AddComponent("content", content)
lm.AddComponent("footer", footer)

// Render
output := lm.Render()
```

## Test Results

### Unit Tests (100% Pass ✅)

```
=== RUN   TestNewBuffer
--- PASS: TestNewBuffer (0.00s)
=== RUN   TestSize
--- PASS: TestSize (0.00s)
=== RUN   TestWidthHeight
--- PASS: TestWidthHeight (0.00s)
=== RUN   TestValid
--- PASS: TestValid (0.00s)
=== RUN   TestSetCell
--- PASS: TestSetCell (0.00s)
=== RUN   TestFillRect
--- PASS: TestFillRect (0.00s)
=== RUN   TestWriteString
--- PASS: TestWriteString (0.00s)
=== RUN   TestWriteStringWideChars
--- PASS: TestWriteStringWideChars (0.00s)
=== RUN   TestWriteStringWrapped
--- PASS: TestWriteStringWrapped (0.00s)
=== RUN   TestWriteBuffer
--- PASS: TestWriteBuffer (0.00s)
=== RUN   TestRender
--- PASS: TestRender (0.00s)
=== RUN   TestTextComponent
--- PASS: TestTextComponent (0.00s)
=== RUN   TestImageComponent
--- PASS: TestImageComponent (0.00s)
=== RUN   TestFillComponent
--- PASS: TestFillComponent (0.00s)
=== RUN   TestLayoutManager
--- PASS: TestLayoutManager (0.00s)
=== RUN   TestLayoutManagerImageLayout
--- PASS: TestLayoutManagerImageLayout (0.00s)
=== RUN   TestLayoutManagerRender
--- PASS: TestLayoutManagerRender (0.00s)
=== RUN   TestLayoutManagerSetSize
--- PASS: TestLayoutManagerSetSize (0.00s)

PASS
```

### Performance Benchmarks

```
BenchmarkFillRect-16                10000    102,538 ns/op    0 B/op    0 allocs/op
BenchmarkWriteString-16           1,524,304      794 ns/op    0 B/op    0 allocs/op
BenchmarkWriteStringWrapped-16     502,207    2,455 ns/op    0 B/op    0 allocs/op
BenchmarkWriteBuffer-16           649,096    1,702 ns/op    0 B/op    0 allocs/op
BenchmarkRender-16                 67,072   16,970 ns/op  904 B/op   32 allocs/op
BenchmarkTextComponentRender-16   973,038    1,257 ns/op   16 B/op    1 alloc/op
BenchmarkLayoutCalculate-16      3,843,465      300 ns/op   96 B/op    3 allocs/op
BenchmarkLayoutRender-16            8,883   150,900 ns/op  235KB/op  79 allocs/op
```

**Interpretation**:
- **WriteString**: 794 ns (extremely fast)
- **Render entire buffer**: 16,970 ns (~0.017ms)
- **Full layout calculation**: 300 ns (negligible)
- **Full layout + render**: 150,900 ns (~0.15ms)

**Result**: Fast enough for 60fps TUI applications (16.6ms budget per frame)

### Wide Character Support ✅

```go
// CJK characters handled correctly (2 columns)
buf := NewBuffer(20, 5)
cols := buf.WriteString(Point{0, 0}, "你好世界", Style{})
// → cols = 8 (4 chars × 2 columns each)
```

### Memory Efficiency ✅

- **Zero allocations** for basic operations (FillRect, WriteString, WriteBuffer)
- **Small allocations** for rendering (32-79 allocs per full frame)
- **Pre-allocated buffer grid** - no resizing during render

## Comparison: String vs Buffer

| Feature | String-Based | Buffer-Based |
|---------|--------------|--------------|
| **Dimension Accuracy** | ❌ Estimated | ✅ Exact |
| **Layout Calculation** | ❌ After render | ✅ Before render |
| **Sixel Image Support** | ❌ Guesswork | ✅ Exact height |
| **Component Isolation** | ❌ Shared string | ✅ Isolated buffers |
| **Wide Chars** | ❌ Complex | ✅ Native support |
| **Performance** | Fast | Fast (-0.13ms overhead) |
| **Debugging** | Hard | Easy (inspect buffer) |

## Use Cases

### 1. Image Viewer

```go
// Perfect vertical centering
lm := NewLayoutManager(80, 30)
lm.ImageLayout(10)  // Sixel display height in lines
// → Content exactly centered at line 14-24
```

### 2. Forms with Dynamic Height

```go
// Form expands/contracts based on content
form := NewForm(...)
height := form.MinHeight()  // Exact calculation!
lm.SetHeight(headerHeight + height + footerHeight)
```

### 3. Complex Layouts

```go
// Multi-column, multi-row layouts
lm := NewLayoutManager(width, height)
lm.CreateGrid(rows, cols)
// → Each cell has exact dimensions
```

### 4. Responsive Design

```go
// Resize handler
func (m Model) Resize(w, h int) {
    lm := buffer.NewLayoutManager(w, h)
    lm.ImageLayout(m.displayHeight)
    m.buffer = lm.GetBuffer()
}
```

## Examples

### Basic Three-Pane Layout

```bash
go run examples/buffer-demo/demo.go 80 30
```

Demonstrates: Header, content, footer with exact layout calculation.

### Image Viewer with Real Image

```bash
cd examples/image-viewer
go run main.go
```

Demonstrates: Sixel image with buffer-based layout system.

### Unit Tests

```bash
go test ./ui/render/buffer/... -v
```

### Benchmarks

```bash
go test ./ui/render/buffer/... -bench=. -benchmem
```

## File Structure

```
ui/render/buffer/
├── buffer.go       # Core buffer implementation
├── components.go   # Renderable components (Text, Image, Fill)
├── layout.go       # LayoutManager for orchestrating layouts
└── buffer_test.go  # Comprehensive unit tests

examples/
├── buffer-demo/
│   └── demo.go     # Basic layout demo
└── image-viewer/
    ├── main.go     # Image viewer with actual image
    └── test-image.jpeg  # Test image
```

## API Reference

### Creating Buffers

```go
// New buffer
buf := NewBuffer(width, height)

// Get dimensions
w := buf.Width()
h := buf.Height()
size := buf.Size()  // Size{Width, Height}
```

### Writing to Buffers

```go
// Write styled text
buf.WriteString(Point{X: 10, Y: 5}, "Hello", Style{Foreground: "196", Bold: true})

// Write with wrapping
buf.WriteStringWrapped(Point{X: 0, Y: 0}, 40, "long text...", Style{})

// Fill rectangle
buf.FillRect(Rect{X: 0, Y: 0, Width: 80, Height: 5}, '░', Style{Foreground: "245"})

// Set single cell
buf.SetCell(Point{X: 10, Y: 5}, Cell{Char: 'X', Style: Style{Bold: true}})
```

### Sub-Buffers

```go
// Create component sub-buffer
compBuf := NewBuffer(rect.Width, rect.Height)

// Component renders to its own buffer
component.Render(compBuf, Rect{X: 0, Y: 0, Width: rect.Width, Height: rect.Height})

// Write to main buffer
mainBuf.WriteBuffer(Point{X: rect.X, Y: rect.Y}, compBuf)
```

### Rendering

```go
// Convert buffer to output string
output := buf.Render()

// Output has consistent newline structure:
// - line 1\n
// - line 2\n
// ...
// - line N (no trailing newline)
```

### Layout Manager

```go
// Create
lm := NewLayoutManager(width, height)

// Set size
lm.SetSize(newWidth, newHeight)

// Calculate layout
lm.CalculateLayout()          // Standard header/content/footer
lm.ImageLayout(displayHeight) // For image viewers

// Add components
lm.AddComponent("header", header)
lm.AddComponent("content", content)
lm.AddComponent("footer", footer)

// Render
output := lm.Render()

// Get buffer for debugging
buf := lm.GetBuffer()
```

## Best Practices

1. **Use LayoutManager for complex layouts** - Handles positioning automatically
2. **Prefer component isolation** - Each component to its own buffer
3. **Use exact display heights** - Don't estimate for Sixel images
4. **Cache buffer when possible** - Reuse buffers instead of recreating
5. **Test with various terminal sizes** - Ensure responsive design

## Future Enhancements

- [ ] Diff/patch comparison buffers
- [ ] Scissor clipping regions
- [ ] Blending/overlay support
- [ ] Animation frame buffers
- [ ] Tearing-free rendering (double buffering)
- [ ] Hardware acceleration hints

## Conclusion

Buffer-based rendering provides **reliable, accurate layouts** for Taproot TUI applications. By working with cell coordinates instead of string estimation, it eliminates the "height calculation hell" that plagues traditional TUI frameworks, especially when dealing with graphics (Sixel) and dynamic content.

**Performance impact**: ~0.15ms per frame (negligible for 60fps applications)
**Developer productivity**: Faster, more reliable UI development
**User experience**: Consistent, accurate layouts across all scenarios
