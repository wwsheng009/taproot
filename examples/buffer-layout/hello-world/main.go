// Buffer Hello World - Simplest buffer layout example
//
// This example demonstrates the core concepts of buffer-based rendering:
// 1. Create a LayoutManager with terminal dimensions
// 2. Add components (header, content, footer)
// 3. Calculate layout to determine component positions
// 4. Render to get the final output
//
// Usage: go run main.go

package main

import (
	"fmt"

	"github.com/wwsheng009/taproot/ui/render/buffer"
)

func main() {
	// Terminal dimensions (can be adjusted)
	width := 80
	height := 24

	// Create a layout manager - this manages the 2D grid for rendering
	lm := buffer.NewLayoutManager(width, height)

	// Create components using the built-in TextComponent
	header := buffer.NewTextComponent(
		"Hello, Buffer World!",
		buffer.Style{
			Bold:       true,
			Foreground: "#86", // Cyan
			Background: "#235", // Dark blue
		},
	).SetCenterH(true) // Center horizontally

	content := buffer.NewTextComponent(
		"Welcome to Taproot's Buffer Layout System\n\n"+
			"Key Concepts:\n"+
			"• Buffer is a 2D grid (width × height cells)\n"+
			"• Layout is calculated BEFORE rendering\n"+
			"• Each component gets its own position rectangle\n"+
			"• Rendering writes to exact cell positions\n\n"+
			"This example shows the simplest possible layout:\n"+
			"  1. Header at the top\n"+
			"  2. Content fills the middle\n"+
			"  3. Footer at the bottom",
		buffer.Style{Foreground: "#250"},
	).SetWrap(true) // Wrap text to fit width

	footer := buffer.NewTextComponent(
		"Press Ctrl+C to exit | Visit: github.com/wwsheng009/taproot",
		buffer.Style{Foreground: "#244", Bold: false},
	).SetCenterH(true)

	// Add components to the layout manager
	// The order matters: header is added first, then content, then footer
	lm.AddComponent("header", header)
	lm.AddComponent("content", content)
	lm.AddComponent("footer", footer)

	// Calculate layout - this determines where each component goes
	// Default layout: header (5 lines) + content (remaining) + footer (1 line)
	lm.CalculateLayout()

	// Render the complete layout to a string
	output := lm.Render()

	// Print the result
	fmt.Print(output)

	// Show layout information
	fmt.Printf("\n\nLayout Information:\n")
	fmt.Printf("  Terminal: %d cols × %d rows = %d cells\n", width, height, width*height)
	fmt.Printf("  Components: 3 (header, content, footer)\n")
	fmt.Printf("  Header: row 0-4 (5 lines)\n")
	fmt.Printf("  Content: row 5-22 (18 lines)\n")
	fmt.Printf("  Footer: row 23 (1 line)\n")
	fmt.Printf("\nKey Insight: The buffer system knows EXACT dimensions!\n")
	fmt.Printf("No counting newlines, no estimating heights - just pixel-perfect layout.\n")
}
