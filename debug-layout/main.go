package main

import (
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/components/layout"
)

func main() {
	// Test the problematic case: small display height (3-6 lines)
	fmt.Println("=== Debugging Small Display Height Issue ===")

	testCases := []struct {
		displayHeight int
	}{
		{1},
		{3},
		{6},
		{10},
		{15},
		{20},
	}

	for _, tc := range testCases {
		fmt.Printf("Test with displayHeight=%d\n", tc.displayHeight)
		fmt.Println(strings.Repeat("-", 60))

		header := buildHeader("Test", 100, tc.displayHeight*6)
		content := "Sixel Image Content"
		footer := buildFooter()

		vbox := layout.NewVerticalLayout().
			SetSize(80, 25).
			SetHeader(header).
			SetContent(content).
			SetFooter(footer).
			SetCenterV(true).
			SetCenterH(false).
			SetSeparator(false)

		output := vbox.Render(tc.displayHeight)
		totalLines := countLines(output)
		headerHeight := vbox.GetHeaderHeight()
		footerHeight := vbox.GetFooterHeight()

		// Count actual newlines in output
		newlineCount := strings.Count(output, "\n")

		fmt.Printf("  Display height param: %d lines\n", tc.displayHeight)
		fmt.Printf("  Header lines: %d\n", headerHeight)
		fmt.Printf("  Footer lines: %d\n", footerHeight)
		lines := strings.Split(output, "\n")
		if len(lines) > 7 {
			fmt.Printf("  Content lines in output: %s\n", lines[7])
		}
		fmt.Printf("  Total newlines in output: %d\n", newlineCount)
		fmt.Printf("  Total lines counted: %d\n", totalLines)
		fmt.Printf("  Expected total: %d (header) + %d (content) + %d (footer) = %d\n",
			headerHeight, tc.displayHeight, footerHeight, headerHeight+tc.displayHeight+footerHeight)

		if totalLines > 25 {
			fmt.Printf("  ❌ EXCEEDS terminal height (%d > 25)\n", totalLines)

			// Print first 30 lines to see what's being output
			lines := strings.Split(output, "\n")
			fmt.Printf("\n  First %d lines of output:\n", min(30, len(lines)))
			for i := 0; i < min(30, len(lines)); i++ {
				if i < len(lines) {
					fmt.Printf("    [%2d] %s\n", i, lines[i])
				}
			}
		} else {
			fmt.Printf("  ✅ Within bounds (%d lines)\n", totalLines)
		}
		fmt.Println()
	}
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
