package image

import (
	"strings"
	"testing"

	"github.com/wwsheng009/taproot/ui/styles"
)

func TestNew(t *testing.T) {
	img := New("test.png")
	if img == nil {
		t.Fatal("New() returned nil")
	}

	if img.Path() != "test.png" {
		t.Errorf("Expected path 'test.png', got '%s'", img.Path())
	}

	// Newly created Images are not loaded until Reload is called
	if img.IsLoaded() {
		t.Error("Expected IsLoaded to be false after New()")
	}
}

func TestPath(t *testing.T) {
	img := New("image.png")
	path := img.Path()
	if path != "image.png" {
		t.Errorf("Expected 'image.png', got '%s'", path)
	}
}

func TestSetPath(t *testing.T) {
	img := New("old.png")
	oldPath := img.Path()

	cmd := img.SetPath("new.png")
	if cmd == nil {
		t.Error("SetPath should return a command")
	}

	newPath := img.Path()
	if newPath == oldPath {
		t.Error("Path should have changed after SetPath")
	}

	if newPath != "new.png" {
		t.Errorf("Expected 'new.png', got '%s'", newPath)
	}
}

func TestSize(t *testing.T) {
	img := New("test.png")
	w, h := img.Size()

	if w != 0 {
		t.Errorf("Expected width 0, got %d", w)
	}
	if h != 0 {
		t.Errorf("Expected height 0, got %d", h)
	}

	img.SetSize(80, 24)
	w, h = img.Size()
	if w != 80 {
		t.Errorf("Expected width 80, got %d", w)
	}
	if h != 24 {
		t.Errorf("Expected height 24, got %d", h)
	}
}

func TestSetSize(t *testing.T) {
	img := New("test.png")
	img.SetSize(100, 50)

	w, h := img.Size()
	if w != 100 {
		t.Errorf("Expected width 100, got %d", w)
	}
	if h != 50 {
		t.Errorf("Expected height 50, got %d", h)
	}
}

func TestSetRenderer(t *testing.T) {
	img := New("test.png")

	cmd := img.SetRenderer(RendererKitty)
	if cmd == nil {
		t.Error("SetRenderer should return a command")
	}
}

func TestRendererTypes(t *testing.T) {
	img := New("test.png")

	tests := []struct {
		name     string
		renderer RendererType
	}{
		{"Auto", RendererAuto},
		{"Kitty", RendererKitty},
		{"iTerm2", RendereriTerm2},
		{"Blocks", RendererBlocks},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := img.SetRenderer(tt.renderer)
			if cmd == nil {
				t.Error("SetRenderer should always return a command")
			}
		})
	}
}

func TestReload(t *testing.T) {
	img := New("test.png")

	cmd := img.Reload()
	if cmd == nil {
		t.Error("Reload should return a command")
	}
}

func TestInit(t *testing.T) {
	img := New("test.png")
	cmd := img.Init()

	if cmd != nil {
		t.Error("Init should return nil (no-op)")
	}
}

func TestView(t *testing.T) {
	img := New("test.png")
	img.SetSize(60, 20)

	view := img.View()
	if view == "" {
		t.Error("View should return a non-empty string")
	}

	// Check for basic elements in the view
	if !strings.Contains(view, "Image") && !strings.Contains(view, "Loading") && !strings.Contains(view, "⚠️") {
		t.Error("View should contain 'Image', 'Loading', or error indicator")
	}
}

func TestViewWithSize(t *testing.T) {
	tests := []struct {
		name  string
		width int
		height int
	}{
		{"Small", 20, 10},
		{"Medium", 60, 20},
		{"Large", 120, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := New("test.png")
			img.SetSize(tt.width, tt.height)

			view := img.View()
			if view == "" {
				t.Error("View should return a non-empty string")
			}
		})
	}
}

func TestRendererTypeString(t *testing.T) {
	types := []struct {
		name     string
		renderer RendererType
	}{
		{"RendererAuto", RendererAuto},
		{"RendererKitty", RendererKitty},
		{"RendereriTerm2", RendereriTerm2},
		{"RendererBlocks", RendererBlocks},
	}

	for i, tt := range types {
		if int(tt.renderer) != i {
			t.Errorf("%s has unexpected value: %d, expected %d", tt.name, tt.renderer, i)
		}
	}
}

func TestError(t *testing.T) {
	// For a non-existent file, should eventually show error
	img := New("non-existent-file-that-does-not-exist.png")

	// Initial state may not have error yet
	// The error is set after trying to load
	_ = img.Reload()
}

func TestStylesIntegration(t *testing.T) {
	s := styles.DefaultStyles()
	// Check styles are valid by calling a method on them
	_ = s.Base.Render("test")

	img := New("test.png")
	img.SetSize(60, 20)
	view := img.View()

	if view == "" {
		t.Error("View should work with DefaultStyles")
	}
}

func TestViewReturnsValidString(t *testing.T) {
	img := New("test.png")
	img.SetSize(80, 24)

	view := img.View()

	// Should be a valid string
	if view == "" {
		t.Fatal("View returned empty string")
	}

	// Should not contain null bytes
	if strings.Contains(view, "\x00") {
		t.Error("View should not contain null bytes")
	}

	// Should have reasonable length (not too long for terminal)
	if len(view) > 10000 {
		t.Errorf("View is too long: %d characters", len(view))
	}
}

func TestSetSizeZero(t *testing.T) {
	img := New("test.png")
	img.SetSize(0, 0)

	w, h := img.Size()
	if w != 0 || h != 0 {
		t.Errorf("SetSize(0, 0) should set size to 0,0, got %d,%d", w, h)
	}

	view := img.View()
	// Should still return valid output even with zero size
	if view == "" {
		t.Error("View should return something even with zero size")
	}
}

func TestSetSizeNegative(t *testing.T) {
	img := New("test.png")
	img.SetSize(-10, -20)

	w, h := img.Size()
	// Should accept the negative values (component handles them)
	if w != -10 || h != -20 {
		t.Errorf("SetSize(-10, -20) should store values, got %d,%d", w, h)
	}
}

func TestMultipleSetOperations(t *testing.T) {
	img := New("test.png")

	// Multiple operations should work without panic
	img.SetPath("a.png")
	img.SetSize(50, 30)
	img.SetSize(80, 40)
	_ = img.SetRenderer(RendererKitty)

	view := img.View()
	if view == "" {
		t.Error("View should work after multiple Set operations")
	}
}

func TestConcurrentAccess(t *testing.T) {
	img := New("test.png")
	img.SetSize(60, 20)

	// Simulate concurrent access (basic test)
	done := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			_ = img.View()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			img.SetSize(60+i, 20+i)
		}
		done <- true
	}()

	<-done
	<-done

	// Should not panic
	view := img.View()
	if view == "" {
		t.Error("View should work after concurrent access")
	}
}

func TestNilStyles(t *testing.T) {
	img := New("test.png")
	// Even without setting styles, should work with defaults
	view := img.View()
	if view == "" {
		t.Error("View should work without explicit styles")
	}
}

func BenchmarkView(b *testing.B) {
	img := New("test.png")
	img.SetSize(80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = img.View()
	}
}

func BenchmarkSetSize(b *testing.B) {
	img := New("test.png")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img.SetSize(80, 24)
	}
}
