package diffview

import (
	"testing"
)

func TestNew(t *testing.T) {
	dv := New()

	if dv == nil {
		t.Error("New() should return non-nil DiffView")
	}

	if dv.layout != LayoutUnified {
		t.Errorf("Expected default layout LayoutUnified, got %v", dv.layout)
	}

	if !dv.lineNumbers {
		t.Error("Expected line numbers to be enabled by default")
	}
}

func TestBefore(t *testing.T) {
	dv := New()
	content := "before content"

	result := dv.Before(content)

	if result != dv {
		t.Error("Before() should return the same DiffView instance")
	}

	if dv.before != content {
		t.Errorf("Expected before content '%s', got '%s'", content, dv.before)
	}
}

func TestAfter(t *testing.T) {
	dv := New()
	content := "after content"

	result := dv.After(content)

	if result != dv {
		t.Error("After() should return the same DiffView instance")
	}

	if dv.after != content {
		t.Errorf("Expected after content '%s', got '%s'", content, dv.after)
	}
}

func TestSetLayout(t *testing.T) {
	dv := New()

	dv.SetLayout(LayoutSplit)

	if dv.layout != LayoutSplit {
		t.Errorf("Expected layout LayoutSplit, got %v", dv.layout)
	}
}

func TestSetLineNumbers(t *testing.T) {
	dv := New()

	dv.SetLineNumbers(false)

	if dv.lineNumbers {
		t.Error("Expected line numbers to be disabled")
	}

	dv.SetLineNumbers(true)

	if !dv.lineNumbers {
		t.Error("Expected line numbers to be enabled")
	}
}

func TestSetSize(t *testing.T) {
	dv := New()

	dv.SetSize(100, 50)

	if dv.width != 100 {
		t.Errorf("Expected width 100, got %d", dv.width)
	}

	if dv.height != 50 {
		t.Errorf("Expected height 50, got %d", dv.height)
	}
}

func TestSetXOffset(t *testing.T) {
	dv := New()

	dv.SetXOffset(10)

	if dv.xOffset != 10 {
		t.Errorf("Expected xOffset 10, got %d", dv.xOffset)
	}

	// Test negative offset
	dv.SetXOffset(-5)

	if dv.xOffset != 0 {
		t.Errorf("Expected xOffset to be clamped to 0, got %d", dv.xOffset)
	}
}

func TestSetYOffset(t *testing.T) {
	dv := New()
	dv.SetSize(80, 20)

	// Setup some lines - need more than height to test scrolling
	dv.before = "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\nline11\nline12\nline13\nline14\nline15\nline16\nline17\nline18\nline19\nline20\nline21\nline22\nline23\nline24\nline25"
	dv.after = dv.before
	dv.Compute()

	dv.SetYOffset(5)

	if dv.yOffset != 5 {
		t.Errorf("Expected yOffset 5, got %d", dv.yOffset)
	}

	// Test offset beyond total
	dv.SetYOffset(100)

	if dv.yOffset > dv.totalLines {
		t.Errorf("Expected yOffset to be clamped to totalLines, got %d", dv.yOffset)
	}
}

func TestCanScrollDown(t *testing.T) {
	dv := New()
	dv.SetSize(80, 10)

	// Setup content
	dv.before = "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\nline11"
	dv.after = dv.before
	dv.Compute()

	if !dv.CanScrollDown() {
		t.Error("Expected to be able to scroll down")
	}

	dv.ScrollToBottom()

	if dv.CanScrollDown() {
		t.Error("Expected to NOT be able to scroll down at bottom")
	}
}

func TestCanScrollUp(t *testing.T) {
	dv := New()
	dv.SetSize(80, 5)  // Height smaller than content

	// Setup content - need more lines than height
	dv.before = "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"
	dv.after = dv.before
	dv.Compute()

	if dv.CanScrollUp() {
		t.Error("Expected to NOT be able to scroll up at top")
	}

	dv.ScrollDown()

	if !dv.CanScrollUp() {
		t.Error("Expected to be able to scroll up after scrolling down")
	}
}

func TestScrollDown(t *testing.T) {
	dv := New()
	dv.SetSize(80, 5)  // Height smaller than content
	dv.before = "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"
	dv.after = dv.before
	dv.Compute()

	initialOffset := dv.yOffset
	dv.ScrollDown()

	if dv.yOffset <= initialOffset {
		t.Error("ScrollDown() should increase yOffset")
	}
}

func TestScrollUp(t *testing.T) {
	dv := New()
	dv.SetSize(80, 5)  // Height smaller than content
	dv.before = "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"
	dv.after = dv.before
	dv.Compute()

	dv.ScrollDown()
	initialOffset := dv.yOffset
	dv.ScrollUp()

	if dv.yOffset >= initialOffset {
		t.Error("ScrollUp() should decrease yOffset")
	}
}

func TestScrollLeft(t *testing.T) {
	dv := New()
	dv.SetXOffset(20)

	dv.ScrollLeft()

	if dv.xOffset != 16 {
		t.Errorf("Expected xOffset 16 after ScrollLeft(), got %d", dv.xOffset)
	}
}

func TestScrollRight(t *testing.T) {
	dv := New()

	dv.ScrollRight()

	if dv.xOffset != 4 {
		t.Errorf("Expected xOffset 4 after ScrollRight(), got %d", dv.xOffset)
	}
}

func TestScrollToTop(t *testing.T) {
	dv := New()
	dv.SetSize(80, 10)
	dv.before = "line1\nline2\nline3"
	dv.after = dv.before
	dv.Compute()

	dv.ScrollDown()
	dv.ScrollDown()
	dv.ScrollToTop()

	if dv.yOffset != 0 {
		t.Errorf("Expected yOffset 0 after ScrollToTop(), got %d", dv.yOffset)
	}
}

func TestScrollToBottom(t *testing.T) {
	dv := New()
	dv.SetSize(80, 10)
	dv.before = "line1\nline2\nline3"
	dv.after = dv.before
	dv.Compute()

	dv.ScrollToBottom()

	if dv.yOffset != 0 {
		t.Error("Expected yOffset 0 when content fits in view")
	}
}

func TestRender(t *testing.T) {
	dv := New()
	dv.SetSize(80, 10)
	dv.before = "line1\nline2\nline3"
	dv.after = "line1\nline2 modified\nline3"
	dv.Compute()

	result := dv.Render()

	if result == "" {
		t.Error("Render() should not return empty string")
	}
}

func TestDefaultLightStyle(t *testing.T) {
	style := DefaultLightStyle()

	// Just check that the function returns a valid style
	if style.SyntaxTheme != "github" {
		t.Errorf("Expected syntax theme 'github', got '%s'", style.SyntaxTheme)
	}
}

func TestDefaultDarkStyle(t *testing.T) {
	style := DefaultDarkStyle()

	// Just check that the function returns a valid style
	if style.SyntaxTheme != "monokai" {
		t.Errorf("Expected syntax theme 'monokai', got '%s'", style.SyntaxTheme)
	}

	if !style.SyntaxHighlight {
		t.Error("Syntax highlighting should be enabled by default")
	}
}

func TestSetFilename(t *testing.T) {
	dv := New()
	filename := "main.go"

	dv.SetFilename(filename)

	if dv.filename != filename {
		t.Errorf("Expected filename '%s', got '%s'", filename, dv.filename)
	}
}

func TestSetSyntaxHighlighting(t *testing.T) {
	dv := New()

	// Default should be disabled
	if dv.useSyntaxHighlighting {
		t.Error("Syntax highlighting should be disabled by default in DiffView")
	}

	dv.SetSyntaxHighlighting(true)

	if !dv.useSyntaxHighlighting {
		t.Error("Syntax highlighting should be enabled after SetSyntaxHighlighting(true)")
	}

	dv.SetSyntaxHighlighting(false)

	if dv.useSyntaxHighlighting {
		t.Error("Syntax highlighting should be disabled after SetSyntaxHighlighting(false)")
	}
}

func TestSyntaxHighlightingRendering(t *testing.T) {
	dv := New()
	dv.SetSize(100, 20)

	// Set content and filename
	dv.SetFilename("test.go")
	dv.before = `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	dv.after = `package main

import "fmt"

func main() {
	fmt.Println("Hello, Taproot!")
}
`
	dv.Compute()

	// Render without syntax highlighting
	dv.SetSyntaxHighlighting(false)
	resultNoHighlight := dv.Render()

	if resultNoHighlight == "" {
		t.Error("Render() should not return empty string without syntax highlighting")
	}

	// Render with syntax highlighting
	dv.SetSyntaxHighlighting(true)
	resultWithHighlight := dv.Render()

	if resultWithHighlight == "" {
		t.Error("Render() should not return empty string with syntax highlighting")
	}

	// The highlighted version should be different from non-highlighted
	// (due to ANSI color codes for syntax)
	if resultNoHighlight == resultWithHighlight {
		t.Error("Render() output should differ with syntax highlighting enabled")
	}
}

func TestGetHighlightedLines(t *testing.T) {
	dv := New()
	dv.SetFilename("test.go")
	dv.before = `package main

func main() {}`
	dv.after = `package main

func main() {
	fmt.Println("hello")
}`

	// Without syntax highlighting enabled
	dv.SetSyntaxHighlighting(false)
	beforeLines, afterLines := dv.getHighlightedLines()

	// Before content: line1, empty, func line = 3 lines
	if len(beforeLines) != 3 {
		t.Errorf("Expected 3 before lines, got %d", len(beforeLines))
	}

	// After content: line1, empty, func line, fmt line, closing brace = 5 lines
	if len(afterLines) != 5 {
		t.Errorf("Expected 5 after lines, got %d", len(afterLines))
	}

	// With syntax highlighting enabled
	dv.SetSyntaxHighlighting(true)
	highlightedBefore, highlightedAfter := dv.getHighlightedLines()

	if len(highlightedBefore) != 3 {
		t.Errorf("Expected 3 highlighted before lines, got %d", len(highlightedBefore))
	}

	if len(highlightedAfter) != 5 {
		t.Errorf("Expected 5 highlighted after lines, got %d", len(highlightedAfter))
	}
}

func TestSplitLineContentWithHighlight(t *testing.T) {
	dv := New()
	dv.SetFilename("test.go")

	beforeLines := []string{"line1", "line2", "line3"}
	afterLines := []string{"line1", "line2 modified", "line3"}

	// Test with syntax highlighting disabled (falls back to splitLineContent)
	dv.SetSyntaxHighlighting(false)

	tests := []struct {
		name          string
		line          DiffLine
		expectedLeft  string
		expectedRight string
	}{
		{
			name:          "Deleted line",
			line:          DiffLine{Type: LineDeleted, Content: "line2", LineNum: 2},
			expectedLeft:  "line2",
			expectedRight: "",
		},
		{
			name:          "Added line",
			line:          DiffLine{Type: LineAdded, Content: "line2 modified", LineNum: 2},
			expectedLeft:  "",
			expectedRight: "line2 modified",
		},
		{
			name:          "Context line",
			line:          DiffLine{Type: LineContext, Content: "line3", LineNum: 3},
			expectedLeft:  "line3",
			expectedRight: "line3",
		},
		{
			name:          "Header line",
			line:          DiffLine{Type: LineHeader, Content: "@@ line info @@", LineNum: 0},
			expectedLeft:  "@@ line info @@",
			expectedRight: "@@ line info @@",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeIdx := 0
			afterIdx := 0
			left, right := dv.splitLineContentWithHighlight(tt.line, &beforeIdx, &afterIdx, beforeLines, afterLines)

			if left != tt.expectedLeft {
				t.Errorf("Expected left content '%s', got '%s'", tt.expectedLeft, left)
			}
			if right != tt.expectedRight {
				t.Errorf("Expected right content '%s', got '%s'", tt.expectedRight, right)
			}
		})
	}
}

func TestSplitLineContentWithHighlight_IndexTracking(t *testing.T) {
	dv := New()
	dv.SetFilename("test.go")
	dv.SetSyntaxHighlighting(true)

	beforeLines := []string{"line1", "line2", "line3"}
	afterLines := []string{"line1", "line2 modified", "line3"}

	beforeIdx := 0
	afterIdx := 0

	// Test deleted line - increments before index
	line1 := DiffLine{Type: LineDeleted, Content: "line1"}
	_, _ = dv.splitLineContentWithHighlight(line1, &beforeIdx, &afterIdx, beforeLines, afterLines)
	if beforeIdx != 1 || afterIdx != 0 {
		t.Errorf("Deleted line should increment beforeIdx, got beforeIdx=%d, afterIdx=%d", beforeIdx, afterIdx)
	}

	// Test added line - increments after index
	line2 := DiffLine{Type: LineAdded, Content: "line1"}
	_, _ = dv.splitLineContentWithHighlight(line2, &beforeIdx, &afterIdx, beforeLines, afterLines)
	if beforeIdx != 1 || afterIdx != 1 {
		t.Errorf("Added line should increment afterIdx, got beforeIdx=%d, afterIdx=%d", beforeIdx, afterIdx)
	}

	// Test context line - increments both indices
	line3 := DiffLine{Type: LineContext, Content: "line2 modified"}
	_, _ = dv.splitLineContentWithHighlight(line3, &beforeIdx, &afterIdx, beforeLines, afterLines)
	if beforeIdx != 2 || afterIdx != 2 {
		t.Errorf("Context line should increment both indices, got beforeIdx=%d, afterIdx=%d", beforeIdx, afterIdx)
	}
}

