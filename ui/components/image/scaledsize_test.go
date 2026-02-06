package image

import (
	"fmt"
	"testing"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

// TestScaledSizeAlgorithm tests the ScaledSize calculation
// ScaledSize returns the sampling area from the original image, NOT the display size.
// Higher zoom = smaller sampling area (more detail, less pixels sampled)
// zoomLevel 2.0 = sample half the area (2x zoom in)
// zoomLevel 0.5 = sample double the area (2x zoom out)
func TestScaledSizeAlgorithm(t *testing.T) {
	// Create test image data (100x100 pixels)
	imgData := &decoder.ImageData{
		Width:  100,
		Height: 100,
	}

	tests := []struct {
		name            string
		displayW        int
		displayH        int
		zoomLevel       float64
		zoomMode        ZoomMode
		expectedMinW    int
		expectedMaxW    int
		expectedMinH    int
		expectedMaxH    int
		description     string
	}{
		{
			name:            "Fit at 1.0x zoom",
			displayW:        80,
			displayH:        40,
			zoomLevel:       1.0,
			zoomMode:        ZoomFit,
			expectedMinW:    35,
			expectedMaxW:    45,
			expectedMinH:    35,
			expectedMaxH:    45,
			description:     "Square 100x100 image in 80x40: limited by height to ~40x40",
		},
		{
			name:            "Fit at 2.0x zoom (zoom in)",
			displayW:        80,
			displayH:        40,
			zoomLevel:       2.0,
			zoomMode:        ZoomFit,
			expectedMinW:    15,
			expectedMaxW:    25,
			expectedMinH:    15,
			expectedMaxH:    25,
			description:     "2x zoom = sample half: 40x40 -> 20x20",
		},
		{
			name:            "Fit at 0.5x zoom (zoom out)",
			displayW:        80,
			displayH:        40,
			zoomLevel:       0.5,
			zoomMode:        ZoomFit,
			expectedMinW:    70,
			expectedMaxW:    90,
			expectedMinH:    70,
			expectedMaxH:    90,
			description:     "0.5x zoom = sample double: 40x40 -> 80x80",
		},
		{
			name:            "Stretch at 2.0x zoom (zoom in)",
			displayW:        80,
			displayH:        40,
			zoomLevel:       2.0,
			zoomMode:        ZoomStretch,
			expectedMinW:    35,
			expectedMaxW:    45,
			expectedMinH:    15,
			expectedMaxH:    25,
			description:     "Stretch 2x zoom: sample half of display (80x40 -> 40x20)",
		},
		{
			name:            "Original at 1.0x zoom",
			displayW:        80,
			displayH:        40,
			zoomLevel:       1.0,
			zoomMode:        ZoomOriginal,
			expectedMinW:    95,
			expectedMaxW:    105,
			expectedMinH:    95,
			expectedMaxH:    105,
			description:     "Should use original image size (100x100)",
		},
		{
			name:            "Original at 0.5x zoom (zoom out)",
			displayW:        80,
			displayH:        40,
			zoomLevel:       0.5,
			zoomMode:        ZoomOriginal,
			expectedMinW:    95,
			expectedMaxW:    105,
			expectedMinH:    95,
			expectedMaxH:    105,
			description:     "0.5x zoom limited by original size: 100x100 -> 100x100 (can't sample more than original)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := &Image{
				imgData:   imgData,
				width:     tt.displayW + 4,
				height:    tt.displayH + 6,
				zoomLevel: tt.zoomLevel,
				zoomMode:  tt.zoomMode,
			}

			w, h := img.ScaledSize()

			t.Logf("%s", tt.description)
			t.Logf("  Display: %dx%d, Zoom: %.2fx, Mode: %v", tt.displayW, tt.displayH, tt.zoomLevel, tt.zoomMode)
			t.Logf("  Sampling area: %dx%d", w, h)

			// Check reasonable bounds
			if w < tt.expectedMinW || w > tt.expectedMaxW {
				t.Errorf("Width %d outside expected range [%d, %d]", w, tt.expectedMinW, tt.expectedMaxW)
			}
			if h < tt.expectedMinH || h > tt.expectedMaxH {
				t.Errorf("Height %d outside expected range [%d, %d]", h, tt.expectedMinH, tt.expectedMaxH)
			}

			// For square image and Fit mode, width should equal height
			if tt.zoomMode == ZoomFit {
				aspectRatio := float64(w) / float64(h)
				if aspectRatio < 0.95 || aspectRatio > 1.05 {
					t.Errorf("Aspect ratio %.2f not close to 1.0 for square image in Fit mode", aspectRatio)
				}
			}
		})
	}

	fmt.Println("\n=== ScaledSize Algorithm Test Summary ===")
	for _, tt := range tests {
		img := &Image{
			imgData:   imgData,
			width:     tt.displayW + 4,
			height:    tt.displayH + 6,
			zoomLevel: tt.zoomLevel,
			zoomMode:  tt.zoomMode,
		}
		w, h := img.ScaledSize()
		fmt.Printf("%s: Sample %dx%d (%s)\n", tt.name, w, h, tt.description)
	}
}

// TestWideImageScaledSize tests scaling with wide images
func TestWideImageScaledSize(t *testing.T) {
	// Create wide test image (200x100 pixels)
	imgData := &decoder.ImageData{
		Width:  200,
		Height: 100,
	}

	img := &Image{
		imgData:   imgData,
		width:     84,  // 80 - 4 margin
		height:    46,  // 40 - 6 margin
		zoomLevel: 1.0,
		zoomMode:  ZoomFit,
	}

	w, h := img.ScaledSize()

	t.Logf("Wide image (200x100) in 80x40 display at 1.0x zoom: sample %dx%d", w, h)

	// Should fit within display bounds
	if w > 84 || h > 46 {
		t.Errorf("Sampling area %dx%d exceeds display bounds", w, h)
	}

	// Check aspect ratio is approximately 2:1
	aspectRatio := float64(w) / float64(h)
	if aspectRatio < 1.9 || aspectRatio > 2.1 {
		t.Errorf("Aspect ratio %.2f not close to 2.0 for wide image", aspectRatio)
	}

	// Test zoom to 2x (zoom in = sample half)
	img.zoomLevel = 2.0
	w2, h2 := img.ScaledSize()
	t.Logf("Wide image (200x100) in 80x40 display at 2.0x zoom: sample %dx%d", w2, h2)

	// Should be half the base sampling size
	if w2 != w/2 || h2 != h/2 {
		t.Errorf("2x zoom should halve sampling area: expected %dx%d, got %dx%d", w/2, h/2, w2, h2)
	}
}
