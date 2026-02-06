package main

import (
	"fmt"
	"os"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <image-path>")
		os.Exit(1)
	}

	path := os.Args[1]

	// Test decoder
	dec := decoder.NewDecoder()
	fmt.Printf("Decoding: %s\n", path)

	data, err := dec.DecodeFile(path)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success!\n")
	fmt.Printf("  Size: %dx%d\n", data.Width, data.Height)
	fmt.Printf("  Format: %s\n", data.Format)
	fmt.Printf("  Memory: ~%d KB\n", data.Bytes()/1024)

	// Test pixel access
	r, g, b, a := data.GetPixelColor(0, 0)
	fmt.Printf("  Pixel(0,0): R=%d G=%d B=%d A=%d\n", r, g, b, a)

	// Test RGBA conversion
	rgba := data.GetRGBA()
	fmt.Printf("  RGBA: %dx%d stride:%d\n", rgba.Bounds().Dx(), rgba.Bounds().Dy(), rgba.Stride)
}
