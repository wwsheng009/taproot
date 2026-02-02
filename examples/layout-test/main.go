package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/wwsheng009/taproot/ui/components/layout"
)

func main() {
	fmt.Println("=== Test 1: Normal Layout ===")
	testNormalLayout()

	fmt.Println("\n=== Test 2: Content Overflow ===")
	testContentOverflow()

	fmt.Println("\n=== Test 3: Sixel Display Height ===")
	testSixelDisplayHeight()

	fmt.Println("\n=== Test 4: Empty Sections ===")
	testEmptySections()

	fmt.Println("\n=== Test 5: Large Header ===")
	testLargeHeader()

	fmt.Println("\n✅ All manual tests completed")
}

func testNormalLayout() {
	header := "Image Viewer Demo\nPath: test.png\nRenderer: Auto\nTerminal: 80x24"
	content := generateASCIIArt(10) // 10 lines
	footer := "+/-: Zoom | 0: Reset | q: Quit  [Zoom: Fit 100%]"

	vbox := layout.NewVerticalLayout().
		SetSize(80, 24).
		SetHeader(header).
		SetContent(content).
		SetFooter(footer).
		SetCenterV(true).
		SetCenterH(false).
		SetSeparator(false)

	output := vbox.Render(0)
	fmt.Println(output)

	lines := countLines(output)
	headerLines := vbox.GetHeaderHeight()
	footerLines := vbox.GetFooterHeight()

	fmt.Fprintf(os.Stderr, "Stats:\n")
	fmt.Fprintf(os.Stderr, "  Total lines: %d\n", lines)
	fmt.Fprintf(os.Stderr, "  Header height: %d\n", headerLines)
	fmt.Fprintf(os.Stderr, "  Footer height: %d\n", footerLines)
	fmt.Fprintf(os.Stderr, "  Content height: %d\n", lines-headerLines-footerLines)

	// Verify header is present
	if !strings.Contains(output, "Image Viewer Demo") {
		fmt.Fprintf(os.Stderr, "❌ Header not found\n")
	} else {
		fmt.Fprintf(os.Stderr, "✓ Header present\n")
	}

	// Verify footer is present
	if !strings.Contains(output, "q: Quit") {
		fmt.Fprintf(os.Stderr, "❌ Footer not found\n")
	} else {
		fmt.Fprintf(os.Stderr, "✓ Footer present\n")
	}
}

func testContentOverflow() {
	header := "Header"
	content := generateASCIIArt(50) // Large content
	footer := "Footer"

	vbox := layout.NewVerticalLayout().
		SetSize(80, 20).
		SetHeader(header).
		SetContent(content).
		SetFooter(footer).
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	output := vbox.Render(0)
	lines := countLines(output)

	fmt.Println(output)
	fmt.Fprintf(os.Stderr, "Content that exceeds available space should be clipped\n")
	fmt.Fprintf(os.Stderr, "Content lines: 50, Terminal height: 20\n")
	fmt.Fprintf(os.Stderr, "Actual output lines: %d\n", lines)

	if lines <= 20 {
		fmt.Fprintf(os.Stderr, "✓ Content properly clipped to fit terminal\n")
	} else {
		fmt.Fprintf(os.Stderr, "❌ Content exceeds terminal height\n")
	}

	// Verify header and footer are still present
	if !strings.Contains(output, header) {
		fmt.Fprintf(os.Stderr, "❌ Header not found after clipping\n")
	} else {
		fmt.Fprintf(os.Stderr, "✓ Header preserved\n")
	}

	if !strings.Contains(output, footer) {
		fmt.Fprintf(os.Stderr, "❌ Footer not found after clipping\n")
	} else {
		fmt.Fprintf(os.Stderr, "✓ Footer preserved\n")
	}
}

func testSixelDisplayHeight() {
	header := "Header\nLine2\nLine3"
	content := "Sixel Image Content" // Will be simulated with displayHeight
	footer := "Footer"

	vbox := layout.NewVerticalLayout().
		SetSize(80, 30).
		SetHeader(header).
		SetContent(content).
		SetFooter(footer).
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	// Simulate Sixel image taking 10 lines
	displayHeight := 10
	output := vbox.Render(displayHeight)

	fmt.Println(output)
	fmt.Fprintf(os.Stderr, "Sixel display height: %d lines\n", displayHeight)
	fmt.Fprintf(os.Stderr, "Terminal height: 30\n")
	fmt.Fprintf(os.Stderr, "Header height: %d\n", vbox.GetHeaderHeight())
	fmt.Fprintf(os.Stderr, "Footer height: %d\n", vbox.GetFooterHeight())
	fmt.Fprintf(os.Stderr, "Expected content space: %d\n", 30-vbox.GetHeaderHeight()-vbox.GetFooterHeight())

	lines := countLines(output)
	if lines <= 30 {
		fmt.Fprintf(os.Stderr, "✓ Layout respects terminal height\n")
	} else {
		fmt.Fprintf(os.Stderr, "❌ Layout exceeds terminal height\n")
	}
}

func testEmptySections() {
	vbox := layout.NewVerticalLayout().
		SetSize(80, 20).
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	output := vbox.Render(0)
	fmt.Println(output)
	fmt.Fprintf(os.Stderr, "Empty layout should not panic\n")
	fmt.Fprintf(os.Stderr, "✓ Empty layout handled\n")
}

func testLargeHeader() {
	header := generateASCIIArt(15) // 15 lines
	content := "Content"
	footer := generateASCIIArt(5) // 5 lines

	vbox := layout.NewVerticalLayout().
		SetSize(80, 25).
		SetHeader(header).
		SetContent(content).
		SetFooter(footer).
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	output := vbox.Render(0)

	fmt.Println(output)
	fmt.Fprintf(os.Stderr, "Header: 15 lines, Footer: 5 lines, Total: 25 lines\n")
	fmt.Fprintf(os.Stderr, "Available for content: %d lines\n", 25-15-5)
	fmt.Fprintf(os.Stderr, "Content should be very limited or clipped\n")

	lines := countLines(output)
	if lines <= 25 {
		fmt.Fprintf(os.Stderr, "✓ Total output fits in terminal\n")
	} else {
		fmt.Fprintf(os.Stderr, "❌ Total output exceeds terminal\n")
	}

	if strings.Contains(output, header) && strings.Contains(output, footer) {
		fmt.Fprintf(os.Stderr, "✓ Both header and footer present\n")
	}
}

func generateASCIIArt(lines int) string {
	var b strings.Builder
	width := 60
	for i := 0; i < lines; i++ {
		for j := 0; j < width; j++ {
			if (i+j)%2 == 0 {
				b.WriteString("/")
			} else {
				b.WriteString("\\")
			}
		}
		if i < lines-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
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
