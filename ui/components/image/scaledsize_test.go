package image

import (
	"fmt"
	"testing"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

// TestScaledSizeAlgorithm tests the ScaledSize calculation
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
			name:            "Fit at 2.0x zoom",
			displayW:        80,
			displayH:        40,
			zoomLevel:       2.0,
			zoomMode:        ZoomFit,
			expectedMinW:    70,
			expectedMaxW:    90,
			expectedMinH:    70,
			expectedMaxH:    90,
			description:     "2x zoom should double base size: 40x40 -> 80x80",
		},
		{
			name:            "Fit at 0.5x zoom",
			displayW:        80,
			displayH:        40,
			zoomLevel:       0.5,
			zoomMode:        ZoomFit,
			expectedMinW:    15,
			expectedMaxW:    25,
			expectedMinH:    15,
			expectedMaxH:    25,
			description:     "0.5x zoom should halve base size: 40x40 -> 20x20",
		},
		{
			name:            "Stretch at 2.0x zoom",
			displayW:        80,
			displayH:        40,
			zoomLevel:       2.0,
			zoomMode:        ZoomStretch,
			expectedMinW:    155,
			expectedMaxW:    165,
			expectedMinH:    75,
			expectedMaxH:    85,
			description:     "Stretch mode: 2x display size = 160x80",
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
			name:            "Original at 0.5x zoom",
			displayW:        80,
			displayH:        40,
			zoomLevel:       0.5,
			zoomMode:        ZoomOriginal,
			expectedMinW:    45,
			expectedMaxW:    55,
			expectedMinH:    45,
			expectedMaxH:    55,
			description:     "Should be half of original (50x50)",
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
			t.Logf("  Result: %dx%d", w, h)

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
		fmt.Printf("%s: %dx%d (%s)\n", tt.name, w, h, tt.description)
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

	t.Logf("Wide image (200x100) in 80x40 display at 1.0x zoom: %dx%d", w, h)

	// Should fit within 80x40, maintaining aspect ratio (2:1)
	if w > 80 || h > 40 {
		t.Errorf("Scaled size %dx%d exceeds display bounds 80x40", w, h)
	}

	// Check aspect ratio is approximately 2:1
	aspectRatio := float64(w) / float64(h)
	if aspectRatio < 1.9 || aspectRatio > 2.1 {
		t.Errorf("Aspect ratio %.2f not close to 2.0 for wide image", aspectRatio)
	}

	// Test zoom to 2x
	img.zoomLevel = 2.0
	w2, h2 := img.ScaledSize()
	t.Logf("Wide image (200x100) in 80x40 display at 2.0x zoom: %dx%d", w2, h2)

	// Should be 2x the base size
	if w2 != w*2 || h2 != h*2 {
		t.Errorf("2x zoom should double size: expected %dx%d, got %dx%d", w*2, h*2, w2, h2)
	}
}
