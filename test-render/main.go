package main

import (
	"fmt"
	"os"

	"github.com/wwsheng009/taproot/ui/components/image"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <image-path>")
		os.Exit(1)
	}

	path := os.Args[1]

	// Create image component
	img := image.New(path)

	// Set size
	img.SetSize(80, 40)

	// Check if loaded
	if !img.IsLoaded() {
		fmt.Printf("ERROR: Image not loaded: %s\n", img.Error())
		os.Exit(1)
	}

	fmt.Printf("Image loaded successfully!\n")
	w, h := img.GetImageDimensions()
	fmt.Printf("Original size: %dx%d\n", w, h)
	fmt.Printf("Renderer: %s\n", img.GetRenderer())

	// Get rendered output
	view := img.View()
	fmt.Printf("\n=== RENDERED OUTPUT ===\n")
	fmt.Printf("%s\n", view)
	fmt.Printf("=== END ===\n")

	fmt.Printf("\nRendered length: %d characters\n", len(view))
}
