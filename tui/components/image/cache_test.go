package image

import (
	"image"
	"image/color"
	"testing"

	"github.com/wwsheng009/taproot/ui/styles"
)

// Test that cache invalidation works correctly
func TestCacheInvalidation(t *testing.T) {
	// Create a test image
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			testImg.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	s := styles.DefaultStyles()
	img := &Image{
		img:      testImg,
		width:    80,
		height:   40,
		scale:    1.0,
		zoomMode: ZoomFitScreen,
		loaded:   true,
		styles:   &s,
	}

	// Initial scaling
	img.scaleImage()
	img.View() // This should cache the view

	// Check cache is valid after first render
	if !img.cacheValid {
		t.Error("Expected cache to be valid after first View() call")
	}

	// Get initial cached view
	initialView := img.cachedView

	// Zoom in
	img.ZoomIn()

	// Check cache was invalidated
	if img.cacheValid {
		t.Error("Expected cache to be invalid after ZoomIn()")
	}

	if img.cachedView != "" {
		t.Error("Expected cachedView to be cleared after ZoomIn()")
	}

	// Scale should have been updated
	if img.scaled == nil {
		t.Error("Expected scaled image to be non-nil after ZoomIn()")
	}

	// Call View() again to generate new view
	newView := img.View()

	// Check that new view is different from cached view
	if newView == initialView && initialView != "" {
		t.Error("Expected View() to return different result after zoom")
	}
}
