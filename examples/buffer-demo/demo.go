// Buffer-based rendering demo for Taproot
//
// This demo demonstrates the key advantages of buffer-based rendering
// over string-based layout calculations:
//
// 1. ACCURATE DIMENSIONS: Buffer grid guarantees exact cell counts
// 2. LAYOUT INDEPENDENCE: Layout planned before rendering content
// 3. NO STRING CALCULATION HELL: No need to count newlines or estimate widths
// 4. COMPONENT ISOLATION: Each component renders to its own buffer
// 5. WIDE CHARACTER SUPPORT: CJK characters handled correctly
//
// Usage: go run main.go [width] [height]
// Example: go run main.go 80 30

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wwsheng009/taproot/ui/render/buffer"
)

func main() {
	// Parse terminal size from command line or use defaults
	width := 80
	height := 30

	if len(os.Args) >= 2 {
		if w, err := strconv.Atoi(os.Args[1]); err == nil && w > 0 {
			width = w
		}
	}
	if len(os.Args) >= 3 {
		if h, err := strconv.Atoi(os.Args[2]); err == nil && h > 0 {
			height = h
		}
	}

	fmt.Println(repeat("=", width))
	fmt.Println("Buffer-Based Rendering Demo - Taproot")
	fmt.Printf("Terminal Size: %d x %d\n", width, height)
	fmt.Println(repeat("=", width))
	fmt.Println()

	// Demo 1: Basic three-pane layout (header, content, footer)
	demoBasicLayout(width, height)

	fmt.Println("\n" + repeat("=", width))
	fmt.Println()

	// Demo 2: Image viewer layout with display height hint
	demoImageLayout(width, height)

	fmt.Println("\n" + repeat("=", width))
	fmt.Println()

	// Demo 3: Small content centered in large space
	demoCentering(width, height)

	fmt.Println("\n" + repeat("=", width))
	fmt.Println("Key Advantages of Buffer-Based Rendering:")
	fmt.Println()
	fmt.Println("1. ACCURATE DIMENSIONS:")
	fmt.Println("   - Buffer size is exact: 80 columns × 30 rows = 2400 cells")
	fmt.Println("   - Each cell is exactly 1 display column")
	fmt.Println("   - No estimation or approximation needed")
	fmt.Println()
	fmt.Println("2. LAYOUT INDEPENDENCE:")
	fmt.Println("   - Calculate layout (CalculateLayout()) before rendering")
	fmt.Println("   - Layout rectangles are exact: x,y,width,height")
	fmt.Println("   - Content height doesn't affect layout calculation")
	fmt.Println()
	fmt.Println("3. NO STRING CALCULATION HELL:")
	fmt.Println("   - String-based: height = countNewlines(string) // WRONG!")
	fmt.Println("   - Buffer-based: height = buffer.Height() // CORRECT!")
	fmt.Println()
	fmt.Println("4. COMPONENT ISOLATION:")
	fmt.Println("   - Each component has its own buffer")
	fmt.Println("   - Render sub-buffer, then write to main buffer")
	fmt.Println("   - Components don't interfere with each other")
	fmt.Println()
	fmt.Println("5. ACCURATE SIXEL SUPPORT:")
	fmt.Println("   - Sixel display height vs actual lines: NO GUESSING")
	fmt.Println("   - Use exact displayHeight hint from image data")
	fmt.Println("   - Layout system handles vertical positioning")
	fmt.Println(repeat("=", width))
}

func demoBasicLayout(width, height int) {
	fmt.Println("Demo 1: Basic Three-Pane Layout")
	fmt.Println("(Header + Content + Footer)")
	fmt.Println()

	// Create layout manager
	lm := buffer.NewLayoutManager(width, height)

	// Calculate layout
	lm.CalculateLayout()

	// Create components
	header := buffer.NewTextComponent(
		"HEADER\nThis is the header area\nFixed at 5 lines tall",
		buffer.Style{
			Bold:       true,
			Foreground: "202", // Orange
			Background: "235", // Dark gray
		},
	).SetCenterH(true)

	content := buffer.NewTextComponent(
		"CONTENT\n\nThis is the content area.\nIt takes all remaining space\nbetween header and footer.\n\n• Point 1\n• Point 2\n• Point 3\n\nThe buffer system calculates\nexact dimensions without\nstring manipulation!",
		buffer.Style{Foreground: "250"},
	).SetWrap(true)

	footer := buffer.NewTextComponent(
		"Press Ctrl+C to exit",
		buffer.Style{
			Bold:       true,
			Foreground: "244", // Gray
		},
	).SetCenterH(true)

	// Add components
	lm.AddComponent("header", header)
	lm.AddComponent("content", content)
	lm.AddComponent("footer", footer)

	// Render
	output := lm.Render()
	fmt.Println(output)

	// Show layout info
	fmt.Println("\nLayout Information:")
	fmt.Printf("  Header: y=0, height=5 (%d cols × 5 rows)\n", width)
	fmt.Printf("  Content: y=5, height=%d (%d cols × %d rows)\n", height-6, width, height-6)
	fmt.Printf("  Footer: y=%d, height=1 (%d cols × 1 row)\n", height-1, width)
	fmt.Printf("  Total: %d cells\n", width*height)
}

func demoImageLayout(width, height int) {
	fmt.Println("Demo 2: Image Viewer Layout")
	fmt.Println("With Sixel display height hint")
	fmt.Println()

	// Simulate different image display heights
	testHeights := []int{5, 10, 15, 20}

	for i, hint := range testHeights {
		fmt.Printf("\n> Test %d: Content height hint = %d lines\n", i+1, hint)

		lm := buffer.NewLayoutManager(width, height)
		lm.ImageLayout(hint)

		// Create image placeholder
		image := buffer.NewImageComponent(width-20, hint*3)

		header := buffer.NewTextComponent(
			fmt.Sprintf("Image Viewer Test %d", i+1),
			buffer.Style{
				Bold:       true,
				Foreground: "81", // Cyan
			},
		).SetCenterH(true)

		footer := buffer.NewTextComponent(
			fmt.Sprintf("Hint: %d lines | Total: %d lines", hint, height),
			buffer.Style{Foreground: "244"},
		).SetCenterH(true)

		lm.AddComponent("header", header)
		lm.AddComponent("content", image)
		lm.AddComponent("footer", footer)

		output := lm.Render()
		fmt.Println(output)
	}
}

func demoCentering(width, height int) {
	fmt.Println("Demo 3: Vertical Centering")
	fmt.Println("Small content in large space")
	fmt.Println()

	lm := buffer.NewLayoutManager(width, height)
	lm.ImageLayout(5) // Content is only 5 lines tall

	smallContent := buffer.NewTextComponent(
		"Small Content\n\nThis text is only 5 lines tall\nbut we have a 30-line terminal.\n\nThe buffer system centers it perfectly!",
		buffer.Style{
			Bold:       true,
			Foreground: "120", // Green
		},
	).SetCenterV(true).SetCenterH(true)

	header := buffer.NewTextComponent(
		"Vertical Centering Demo",
		buffer.Style{
			Bold:       true,
			Foreground: "226", // Yellow
		},
	).SetCenterH(true)

	footer := buffer.NewTextComponent(
		"Watch how the content is centered vertically",
		buffer.Style{Foreground: "244"},
	).SetCenterH(true)

	lm.AddComponent("header", header)
	lm.AddComponent("content", smallContent)
	lm.AddComponent("footer", footer)

	output := lm.Render()
	fmt.Println(output)

	fmt.Println("\nVertical Centering Calculation:")
	fmt.Printf("  Total space: %d lines\n", height)
	fmt.Printf("  Header: 5 lines\n")
	fmt.Printf("  Footer: 1 line\n")
	fmt.Printf("  Available content space: %d lines\n", height-6)
	fmt.Printf("  Content height: 5 lines\n")
	fmt.Printf("  Padding top: %d lines\n", (height-6-5)/2)
	fmt.Printf("  Padding bottom: %d lines\n", (height-6-5)-(height-6-5)/2)
}

// Simple string repeat utility
func repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	result := s
	for i := 1; i < count; i++ {
		result += s
	}
	return result
}
