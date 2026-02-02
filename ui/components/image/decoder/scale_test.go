package decoder

import (
	"image"
	"image/color"
	"testing"
)

// TestScaleWithOptions tests the ScaleWithOptions method
func TestScaleWithOptions(t *testing.T) {
	// Create a simple 10x10 test image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	imgData := &ImageData{
		Width:  10,
		Height: 10,
		Img:    img,
	}

	tests := []struct {
		name          string
		maxWidth      int
		maxHeight     int
		allowUpscale  bool
		expectedW     int
		expectedH     int
		description   string
	}{
		{
			name:         "Downscale only",
			maxWidth:     5,
			maxHeight:    5,
			allowUpscale: false,
			expectedW:    5,
			expectedH:    5,
			description:  "Should scale down to 5x5",
		},
		{
			name:         "Upscale denied",
			maxWidth:     20,
			maxHeight:    20,
			allowUpscale: false,
			expectedW:    10,
			expectedH:    10,
			description:  "Should NOT upscale (stays 10x10)",
		},
		{
			name:         "Upscale allowed 2x",
			maxWidth:     20,
			maxHeight:    20,
			allowUpscale: true,
			expectedW:    20,
			expectedH:    20,
			description:  "Should upscale to 20x20",
		},
		{
			name:         "Upscale allowed 4x",
			maxWidth:     40,
			maxHeight:    40,
			allowUpscale: true,
			expectedW:    40,
			expectedH:    40,
			description:  "Should upscale to 40x40",
		},
		{
			name:         "Maintain aspect ratio when upscaling",
			maxWidth:     100,
			maxHeight:    50,
			allowUpscale: true,
			expectedW:    50,
			expectedH:    50,
			description:  "Limited by height, should upscale to 50x50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := imgData.ScaleWithOptions(tt.maxWidth, tt.maxHeight, tt.allowUpscale)

			t.Logf("%s", tt.description)
			t.Logf("  Request: %dx%d, allowUpscale: %v", tt.maxWidth, tt.maxHeight, tt.allowUpscale)
			t.Logf("  Result: %dx%d (expected %dx%d)", w, h, tt.expectedW, tt.expectedH)

			if w != tt.expectedW || h != tt.expectedH {
				t.Errorf("Expected %dx%d, got %dx%d", tt.expectedW, tt.expectedH, w, h)
			}
		})
	}
}

// TestScaleBackwardCompatibility tests that old Scale() behavior is preserved
func TestScaleBackwardCompatibility(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	imgData := &ImageData{
		Width:  10,
		Height: 10,
		Img:    img,
	}

	// Test downscaling (should work as before)
	w, h := imgData.Scale(5, 5)
	if w != 5 || h != 5 {
		t.Errorf("Downscaling failed: expected 5x5, got %dx%d", w, h)
	}

	// Test upscaling (should be denied by default)
	w, h = imgData.Scale(20, 20)
	if w != 10 || h != 10 {
		t.Errorf("Default Scale() should not upscale: expected 10x10, got %dx%d", w, h)
	}

	t.Logf("✅ Backward compatibility preserved")
}
