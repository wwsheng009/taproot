package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	fmt.Printf("COLORTERM=%s\n", os.Getenv("COLORTERM"))
	fmt.Printf("TERM=%s\n", os.Getenv("TERM"))
	fmt.Printf("WT_SESSION=%s\n", os.Getenv("WT_SESSION"))
	
	// Test with explicit RGB
	style1 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).  // Red
		Background(lipgloss.Color("46"))    // Green
	
	result1 := style1.Render("▀")
	fmt.Printf("\nWith 256-color codes:\n%s\n", result1)
	fmt.Printf("Raw: %q\n\n", result1)
	
	// Test with hex
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0000")).
		Background(lipgloss.Color("#00ff00"))
	
	result2 := style2.Render("▀")
	fmt.Printf("With hex colors:\n%s\n", result2)
	fmt.Printf("Raw: %q\n", result2)
	
	// Test without styles
	fmt.Printf("\nPlain char: ▀\n")
}
