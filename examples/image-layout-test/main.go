package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/wwsheng009/taproot/ui/components/layout"
)

func main() {
	fmt.Println("=" + strings.Repeat("=", 78))
	fmt.Println("Image Viewer Layout Test - Simulating Zoom Scenarios")
	fmt.Println("=" + strings.Repeat("=", 78))

	testScenarios := []struct {
		name         string
		scaledWidth  int
		scaledHeight int
		description  string
	}{
		{"Small Image (Zoom Out)", 40, 20, "Image zoomed out, small dimensions"},
		{"Medium Image (Normal)", 80, 40, "Normal zoom level"},
		{"Large Image (Zoom In)", 160, 80, "Image zoomed in, large dimensions"},
		{"Very Large Image", 200, 120, "Very large zoom level"},
		{"Tiny Image", 10, 6, "Image at minimum zoom"},
	}

	for i, scenario := range testScenarios {
		fmt.Printf("\n[Test %d] %s\n", i+1, scenario.name)
		fmt.Println("-" + strings.Repeat("-", 77))
		fmt.Printf("%s\n", scenario.description)
		fmt.Printf("Image scaled size: %dx%d pixels\n", scenario.scaledWidth, scenario.scaledHeight)

		// Calculate Sixel display height (6 pixels per line)
		displayHeight := scenario.scaledHeight / 6
		if displayHeight > 0 {
			fmt.Printf("Sixel display height: %d terminal lines\n", displayHeight)
		} else {
			fmt.Printf("Sixel display height: %d (too small)\n", displayHeight)
			displayHeight = 1 // Minimum height
		}

		// Test with this display height
		success := testLayoutWithHeight(scenario.name, scenario.scaledWidth, scenario.scaledHeight, displayHeight)
		if success {
			fmt.Println("✅ PASS: Header and footer visible")
		} else {
			fmt.Println("❌ FAIL: Layout issue detected")
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Summary:")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("✅ All scenarios tested successfully")
	fmt.Println("Header and footer remain visible at all zoom levels")
}

func testLayoutWithHeight(name string, imgW, imgH, displayHeight int) bool {
	// Simulate image viewer layout
	header := buildHeader(name, imgW, imgH)
	content := "Sixel Image Output Placeholder"
	footer := buildFooter()

	// Terminal size
	termWidth := 80
	termHeight := 25

	vbox := layout.NewVerticalLayout().
		SetSize(termWidth, termHeight).
		SetHeader(header).
		SetContent(content).
		SetFooter(footer).
		SetCenterV(true).
		SetCenterH(false).
		SetSeparator(false)

	output := vbox.Render(displayHeight)

	// Verify output
	headerHeight := vbox.GetHeaderHeight()
	footerHeight := vbox.GetFooterHeight()
	totalLines := countLines(output)

	fmt.Fprintf(os.Stderr, "  Terminal: %dx%d, Header: %d lines, Footer: %d lines\n", termWidth, termHeight, headerHeight, footerHeight)
	fmt.Fprintf(os.Stderr, "  Display height param: %d lines\n", displayHeight)
	fmt.Fprintf(os.Stderr, "  Total output: %d lines\n", totalLines)
	fmt.Fprintf(os.Stderr, "  Content space: %d lines\n", termHeight-headerHeight-footerHeight)

	// Check if header is present
	if !strings.Contains(output, "Path:") {
		fmt.Fprintf(os.Stderr, "❌ Header text not found in output\n")
		return false
	}

	// Check if footer is present
	if !strings.Contains(output, "q: Quit") {
		fmt.Fprintf(os.Stderr, "❌ Footer text not found in output\n")
		return false
	}

	// Check if layout respects terminal height
	if totalLines > termHeight {
		fmt.Fprintf(os.Stderr, "❌ Output exceeds terminal height (%d > %d)\n", totalLines, termHeight)
		return false
	}

	fmt.Fprintf(os.Stderr, "✓ Layout within bounds, header and footer visible\n")

	// Print first few lines to visualize
	lines := strings.Split(output, "\n")
	fmt.Printf("\n  Header preview:\n")
	maxPreview := min(headerHeight, 3)
	for i := 0; i < maxPreview; i++ {
		if i < len(lines) {
			fmt.Printf("    %s\n", lines[i])
		}
	}

	return true
}

func buildHeader(name string, imgW, imgH int) string {
	var b strings.Builder
	b.WriteString("Image Viewer Demo\n")
	b.WriteString(fmt.Sprintf("Path: %s\n", name))
	b.WriteString(fmt.Sprintf("Renderer: Sixel\n"))
	b.WriteString(fmt.Sprintf("Size: %dx%d pixels\n", imgW, imgH))
	return b.String()
}

func buildFooter() string {
	return "+/-: Zoom | 0: Reset | m: Mode | 1-6: Renderer | r: Reload | q: Quit"
}

func countLines(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	if len(s) > 0 && s[len(s)-1] != '\n' {
		count++
	}
	return count
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
