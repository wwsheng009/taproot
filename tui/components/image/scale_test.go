package image

import (
	"testing"
)

func TestSetScale(t *testing.T) {
	img := &Image{
		scale: 1.0,
	}

	// Test setting scale
	cmd := img.SetScale(2.0)
	if cmd != nil {
		t.Errorf("SetScale should return nil command, got %v", cmd)
	}

	if img.scale != 2.0 {
		t.Errorf("Expected scale to be 2.0, got %f", img.scale)
	}

	// Test clamping minimum
	img.SetScale(0.05)
	if img.scale != 0.1 {
		t.Errorf("Expected scale to be clamped to 0.1, got %f", img.scale)
	}

	// Test clamping maximum
	img.SetScale(10.0)
	if img.scale != 5.0 {
		t.Errorf("Expected scale to be clamped to 5.0, got %f", img.scale)
	}
}

func TestGetRenderer(t *testing.T) {
	img := &Image{
		renderer: RendererKitty,
	}

	if img.GetRenderer() != RendererKitty {
		t.Errorf("Expected renderer to be RendererKitty, got %v", img.GetRenderer())
	}
}

func TestGetImageDimensions(t *testing.T) {
	// Test with nil image
	img := &Image{}
	w, h := img.GetImageDimensions()
	if w != 0 || h != 0 {
		t.Errorf("Expected 0x0 for nil image, got %dx%d", w, h)
	}
}

func TestGetRendererDescription(t *testing.T) {
	tests := []struct {
		renderer   RendererType
		contains   string
	}{
		{RendererAuto, "Auto"},
		{RendererSixel, "Sixel"},
		{RendererKitty, "Kitty"},
		{RendereriTerm2, "iTerm2"},
		{RendererBlocks, "Blocks"},
		{RendererASCII, "ASCII"},
	}

	for _, tt := range tests {
		desc := GetRendererDescription(tt.renderer)
		if len(desc) == 0 {
			t.Errorf("GetRendererDescription(%v) returned empty string", tt.renderer)
		}
	}
}
