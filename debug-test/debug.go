package main

import (
	"fmt"
	"os"

	"github.com/wwsheng009/taproot/ui/components/image"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug.go <image-path>")
		os.Exit(1)
	}

	imgPath := os.Args[1]
	img := image.New(imgPath)
	
	// Set terminal size
	img.SetSize(80, 40)

	origW, origH := img.GetImageDimensions()
	fmt.Printf("=== Original Image ===\n")
	fmt.Printf("Dimensions: %dx%d pixels\n", origW, origH)

	fmt.Printf("\n=== Testing Zoom Levels ===\n")

	for _, zoom := range []float64{0.5, 1.0, 1.5, 2.0, 4.0} {
		img.SetScale(zoom)
		w, h := img.ScaledSize()
		fmt.Printf("Zoom %.1fx: ScaledSize returns %dx%d\n", zoom, w, h)
		
		// View the image to see actual output
		view := img.View()
		lineCount := 0
		for _, c := range view {
			if c == '\n' {
				lineCount++
			}
		}
		fmt.Printf("  View output: %d lines\n", lineCount)
	}
}
