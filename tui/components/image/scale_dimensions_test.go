package image

import (
	"image"
	"image/color"
	"testing"
)

// Test that zoom actually changes the scaled image dimensions
func TestZoomChangesScaledDimensions(t *testing.T) {
	// Create a test image (100x100 pixels)
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			testImg.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	img := &Image{
		img:     testImg,
		width:   80,
		height:  40,
		scale:   1.0,
			zoomMode: ZoomFitScreen,
		loaded:  true,
	}

	// Initial scaling
	img.scaleImage()

	initialW, initialH := img.ScaledSize()
	t.Logf("Initial scaled size: %dx%d (scale=%.2f)", initialW, initialH, img.scale)

	// Zoom in (should increase dimensions)
	img.ZoomIn()
	newW, newH := img.ScaledSize()
	t.Logf("After ZoomIn: %dx%d (scale=%.2f)", newW, newH, img.scale)

	if newW <= initialW || newH <= initialH {
		t.Errorf("Expected dimensions to increase after ZoomIn, got %dx%d (was %dx%d)",
			newW, newH, initialW, initialH)
	}

	// Zoom again
	img.ZoomIn()
	newW2, newH2 := img.ScaledSize()
	t.Logf("After second ZoomIn: %dx%d (scale=%.2f)", newW2, newH2, img.scale)

	if newW2 <= newW || newH2 <= newH {
		t.Errorf("Expected dimensions to increase again, got %dx%d (was %dx%d)",
			newW2, newH2, newW, newH)
	}

	// Set scale to 2.0 directly
	img.SetScale(2.0)
	w2, h2 := img.ScaledSize()
	t.Logf("After SetScale(2.0): %dx%d (scale=%.2f)", w2, h2, img.scale)

	// At 2.0 scale, should be approximately 2x the initial size
	expectedW := initialW * 2
	expectedH := initialH * 2
	if w2 < expectedW-1 || w2 > expectedW+1 {
		t.Errorf("Expected width ~%d at 2.0x scale, got %d", expectedW, w2)
	}
	if h2 < expectedH-1 || h2 > expectedH+1 {
		t.Errorf("Expected height ~%d at 2.0x scale, got %d", expectedH, h2)
	}

	// Set scale to 0.5
	img.SetScale(0.5)
	w05, h05 := img.ScaledSize()
	t.Logf("After SetScale(0.5): %dx%d (scale=%.2f)", w05, h05, img.scale)

	// At 0.5 scale, should be approximately half the initial size
	expectedW05 := initialW / 2
	expectedH05 := initialH / 2
	if w05 < expectedW05-1 || w05 > expectedW05+1 {
		t.Errorf("Expected width ~%d at 0.5x scale, got %d", expectedW05, w05)
	}
	if h05 < expectedH05-1 || h05 > expectedH05+1 {
		t.Errorf("Expected height ~%d at 0.5x scale, got %d", expectedH05, h05)
	}
}
