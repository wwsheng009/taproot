package main

import (
	"fmt"
	"os"

	"github.com/wwsheng009/taproot/ui/components/image"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug_zoom.go <image-path>")
		os.Exit(1)
	}

	imgPath := os.Args[1]
	img := image.New(imgPath)

	// Simulate terminal size
	img.SetSize(80, 40)

	// Get initial scale
	scale1 := img.GetScale()
	w1, h1 := img.ScaledSize()
	origW, origH := img.GetImageDimensions()
	fmt.Printf("=== Initial State ===\n")
	fmt.Printf("Zoom Level: %.2f\n", scale1)
	fmt.Printf("Scaled Size: %dx%d\n", w1, h1)
	fmt.Printf("Original Size: %dx%d\n", origW, origH)

	// Test ZoomIn
	fmt.Printf("\n=== After ZoomIn() ===\n")
	img.ZoomIn()
	scale2 := img.GetScale()
	w2, h2 := img.ScaledSize()
	fmt.Printf("Zoom Level: %.2f\n", scale2)
	fmt.Printf("Scaled Size: %dx%d\n", w2, h2)

	if w2 == w1 && h2 == h1 {
		fmt.Printf("❌ ERROR: Size did not change after ZoomIn!\n")
	} else {
		fmt.Printf("✅ Size changed: %dx%d -> %dx%d\n", w1, h1, w2, h2)
		ratioW := float64(w2) / float64(w1)
		ratioH := float64(h2) / float64(h1)
		fmt.Printf("   Ratio: %.2fx (width), %.2fx (height)\n", ratioW, ratioH)
	}

	// Test SetScale(2.0)
	fmt.Printf("\n=== After SetScale(2.0) ===\n")
	img.SetScale(2.0)
	scale3 := img.GetScale()
	w3, h3 := img.ScaledSize()
	fmt.Printf("Zoom Level: %.2f\n", scale3)
	fmt.Printf("Scaled Size: %dx%d\n", w3, h3)

	if w3 != w1*2 || h3 != h1*2 {
		fmt.Printf("❌ ERROR: Expected 2x size, got %dx%d (expected %dx%d)\n", w3, h3, w1*2, h1*2)
	} else {
		fmt.Printf("✅ Correctly doubled to 2x size\n")
	}

	// Test View() output
	fmt.Printf("\n=== Testing View() Output ===\n")
	view1 := img.View()
	lines1 := countLines(view1)
	fmt.Printf("View output at zoom=%.2f: %d lines\n", img.GetScale(), lines1)

	img.SetScale(1.0)
	view2 := img.View()
	lines2 := countLines(view2)
	fmt.Printf("View output at zoom=%.2f: %d lines\n", img.GetScale(), lines2)

	if lines1 != lines2 {
		fmt.Printf("✅ View output changes with zoom (%d vs %d lines)\n", lines1, lines2)
	} else {
		fmt.Printf("❌ ERROR: View output does not change with zoom!\n")
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
