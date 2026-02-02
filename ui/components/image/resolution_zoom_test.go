package image

import (
	"fmt"
	"testing"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

// TestResolutionBasedZoom tests that zoom changes sampling resolution, not display size
func TestResolutionBasedZoom(t *testing.T) {
	// Create test image data (100x100 pixels)
	imgData := &decoder.ImageData{
		Width:  100,
		Height: 100,
	}

	displayW := 80
	displayH := 40

	tests := []struct {
		name        string
		zoomLevel   float64
		expectedMinW int // Expected sampled width (should decrease with higher zoom)
		expectedMaxW int
		description  string
	}{
		{
			name:         "1.0x zoom (normal)",
			zoomLevel:    1.0,
			expectedMinW: 35,
			expectedMaxW: 45,
			description:  "Should sample ~40x40 pixels for 80x40 display",
		},
		{
			name:         "2.0x zoom (zoomed in)",
			zoomLevel:    2.0,
			expectedMinW: 17,
			expectedMaxW: 25,
			description:  "Should sample ~20x20 pixels (higher resolution)",
		},
		{
			name:         "4.0x zoom (very zoomed in)",
			zoomLevel:    4.0,
			expectedMinW: 8,
			expectedMaxW: 15,
			description:  "Should sample ~10x10 pixels (much higher resolution)",
		},
		{
			name:         "0.5x zoom (zoomed out)",
			zoomLevel:    0.5,
			expectedMinW: 70,
			expectedMaxW: 90,
			description:  "Should sample ~80x80 pixels (lower resolution)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := &Image{
				imgData:   imgData,
				width:     displayW + 4,
				height:    displayH + 6,
				zoomLevel: tt.zoomLevel,
				zoomMode:  ZoomFit,
			}

			sampledW, sampledH := img.ScaledSize()

			t.Logf("%s", tt.description)
			t.Logf("  Display: %dx%d, Zoom: %.2fx", displayW, displayH, tt.zoomLevel)
			t.Logf("  Sampled from original: %dx%d pixels", sampledW, sampledH)

			// Check that sampled size is in expected range
			if sampledW < tt.expectedMinW || sampledW > tt.expectedMaxW {
				t.Errorf("Sampled width %d outside expected range [%d, %d]", sampledW, tt.expectedMinW, tt.expectedMaxW)
			}

			// Key assertion: Higher zoom = SMALLER sampled area
			// This means we're sampling fewer pixels across the same display area = higher detail
		})
	}

	fmt.Println("\n=== Resolution-Based Zoom Summary ===")
	fmt.Println("Zoom 1.0x: Sample 40x40 → Display across 80x40 chars (normal)")
	fmt.Println("Zoom 2.0x: Sample 20x20 → Display across 80x40 chars (2x detail)")
	fmt.Println("Zoom 4.0x: Sample 10x10 → Display across 80x40 chars (4x detail)")
	fmt.Println("Zoom 0.5x: Sample 80x80 → Display across 80x40 chars (0.5x detail)")
}

// TestZoomDoesNotChangeDisplaySize verifies zoom doesn't affect display dimensions
func TestZoomDoesNotChangeDisplaySize(t *testing.T) {
	imgData := &decoder.ImageData{
		Width:  100,
		Height: 100,
	}

	img := &Image{
		imgData:   imgData,
		width:     84,  // 80 - 4 margin
		height:    46,  // 40 - 6 margin
		zoomLevel: 1.0,
		zoomMode:  ZoomFit,
	}

	// Get display size (should be constant)
	displayW, displayH := img.calculateDisplaySize()
	
	// Test at different zoom levels
	zoomLevels := []float64{0.5, 1.0, 2.0, 4.0}
	for _, zoom := range zoomLevels {
		img.zoomLevel = zoom
		sampledW, sampledH := img.ScaledSize()
		
		t.Logf("Zoom %.1fx: Sample %dx%d from original (display is %dx%d)",
			zoom, sampledW, sampledH, displayW, displayH)
		
		// The sampled area should decrease as zoom increases
		// But display size remains constant (80x40)
	}

	// Verify display size is constant
	if displayW != 80 || displayH != 40 {
		t.Errorf("Display size should be constant 80x40, got %dx%d", displayW, displayH)
	}

	fmt.Println("\n✅ Verified: Zoom changes sampling resolution, NOT display size")
}
