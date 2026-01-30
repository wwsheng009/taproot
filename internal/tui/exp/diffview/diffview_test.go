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

	// Setup some lines
	dv.before = "line1\nline2\nline3"
	dv.after = "line1\nline2\nline3"
	dv.computeDiff()

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
	dv.computeDiff()

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
	dv.SetSize(80, 10)

	// Setup content
	dv.before = "line1\nline2\nline3"
	dv.after = dv.before
	dv.computeDiff()

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
	dv.SetSize(80, 10)
	dv.before = "line1\nline2\nline3"
	dv.after = dv.before
	dv.computeDiff()

	initialOffset := dv.yOffset
	dv.ScrollDown()

	if dv.yOffset <= initialOffset {
		t.Error("ScrollDown() should increase yOffset")
	}
}

func TestScrollUp(t *testing.T) {
	dv := New()
	dv.SetSize(80, 10)
	dv.before = "line1\nline2\nline3"
	dv.after = dv.before
	dv.computeDiff()

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
	dv.computeDiff()

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
	dv.computeDiff()

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
	dv.computeDiff()

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
