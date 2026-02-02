package image

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/wwsheng009/taproot/ui/styles"
)

// Test that renderBlocks output size changes with zoom
func TestRenderBlocksSizeChanges(t *testing.T) {
	// Create a test image (100x100 pixels)
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

	// Initial scaling and render
	img.scaleImage()
	view1 := img.renderBlocks()

	// Count the number of lines in the view
	lines1 := countLines(view1)
	t.Logf("Initial view lines: %d (scale=%.2f)", lines1, img.scale)

	// Zoom in
	img.SetScale(2.0)
	view2 := img.renderBlocks()

	lines2 := countLines(view2)
	t.Logf("After zoom to 2.0x view lines: %d (scale=%.2f)", lines2, img.scale)

	if lines2 <= lines1 {
		t.Errorf("Expected more lines after zoom (got %d, was %d)", lines2, lines1)
	}

	// Zoom in more
	img.SetScale(4.0)
	view3 := img.renderBlocks()

	lines3 := countLines(view3)
	t.Logf("After zoom to 4.0x view lines: %d (scale=%.2f)", lines3, img.scale)

	if lines3 <= lines2 {
		t.Errorf("Expected even more lines after more zoom (got %d, was %d)", lines3, lines2)
	}

	// Check that the actual content size (not including borders) increased
	contentSize1 := getRenderedContentSize(view1)
	contentSize2 := getRenderedContentSize(view2)
	contentSize3 := getRenderedContentSize(view3)

	t.Logf("Content sizes: %d -> %d -> %d", contentSize1, contentSize2, contentSize3)

	if contentSize2 <= contentSize1 {
		t.Errorf("Expected content size to increase (got %d, was %d)", contentSize2, contentSize1)
	}
}

func countLines(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	return count
}

func getRenderedContentSize(s string) int {
	// Count actual color blocks (escape sequences with "m ")
	count := 0
	for i := 0; i < len(s)-2; i++ {
		if s[i] == 'm' && s[i+1] == ' ' {
			count++
		}
	}
	return count
}

// Test cache invalidation more thoroughly
func TestCacheInvalidationDetailed(t *testing.T) {
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

	// Initial render
	img.scaleImage()
	view1 := img.renderBlocks()
	lines1 := countLines(view1)

	t.Logf("Step 1: scale=%.2f, cacheValid=%v, lines=%d", img.scale, img.cacheValid, lines1)

	// Zoom in - this should invalidate cache
	img.SetScale(2.0)
	t.Logf("Step 2: scale=%.2f, cacheValid=%v", img.scale, img.cacheValid)

	if img.cacheValid {
		t.Error("Cache should be invalid after SetScale")
	}

	view2 := img.renderBlocks()
	lines2 := countLines(view2)

	t.Logf("Step 3: after render, scale=%.2f, cacheValid=%v, lines=%d", img.scale, img.cacheValid, lines2)

	if lines2 == lines1 {
		t.Errorf("Lines should be different: %d vs %d", lines2, lines1)
	}

	// Call renderBlocks again - should use cache
	view3 := img.renderBlocks()

	if view2 != view3 {
		t.Error("Cached view should be identical")
	}

	// Zoom again
	img.SetScale(3.0)
	t.Logf("Step 4: scale=%.2f, cacheValid=%v", img.scale, img.cacheValid)

	if img.cacheValid {
		t.Error("Cache should be invalid after another SetScale")
	}

	view4 := img.renderBlocks()
	lines4 := countLines(view4)

	t.Logf("Step 5: after render, scale=%.2f, lines=%d", img.scale, lines4)

	fmt.Printf("\n=== Render Size Debug ===\n")
	fmt.Printf("Scale 1.0: %d lines\n", lines1)
	fmt.Printf("Scale 2.0: %d lines\n", lines2)
	fmt.Printf("Scale 3.0: %d lines\n", lines4)
	fmt.Printf("========================\n")
}
