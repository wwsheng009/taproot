package main

import (
	"fmt"

	"github.com/wwsheng009/taproot/internal/tui/exp/diffview"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("          Taproot Diff Viewer Demo")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Example 1: Simple unified diff
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("1. Unified Diff View")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	before := `package main

import "fmt"

func hello() {
	fmt.Println("Hello, World!")
}

func main() {
	hello()
}`

	after := `package main

import "fmt"

func hello() {
	fmt.Println("Hello, Taproot!")
}

func goodbye() {
	fmt.Println("Goodbye!")
}

func main() {
	hello()
	goodbye()
}`

	dv := diffview.New()
	dv.Before(before)
	dv.After(after)
	dv.SetSize(80, 15)

	fmt.Println(dv.Render())
	fmt.Println()

	// Example 2: Split view style info
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("2. Diff Style Information")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	lightStyle := diffview.DefaultLightStyle()
	fmt.Printf("Light Theme: SyntaxHighlight=%v, SyntaxTheme=%s\n",
		lightStyle.SyntaxHighlight, lightStyle.SyntaxTheme)

	darkStyle := diffview.DefaultDarkStyle()
	fmt.Printf("Dark Theme: SyntaxHighlight=%v, SyntaxTheme=%s\n",
		darkStyle.SyntaxHighlight, darkStyle.SyntaxTheme)

	fmt.Println()

	// Example 3: Scrolling controls
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("3. Scrolling Features")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	dv2 := diffview.New()
	dv2.Before(before)
	dv2.After(after)
	dv2.SetSize(60, 5)

	fmt.Printf("Can scroll up: %v\n", dv2.CanScrollUp())
	fmt.Printf("Can scroll down: %v\n", dv2.CanScrollDown())

	dv2.ScrollDown()
	fmt.Printf("After ScrollDown - can scroll down: %v\n", dv2.CanScrollDown())

	dv2.ScrollToTop()
	fmt.Printf("After ScrollToTop - can scroll up: %v\n", dv2.CanScrollUp())

	fmt.Println()

	// Example 4: Layout options
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("4. Layout Options")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	fmt.Println("LayoutUnified: Traditional unified diff format")
	fmt.Println("LayoutSplit: Side-by-side comparison (framework ready)")

	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("                    Demo Complete")
	fmt.Println("═══════════════════════════════════════════════════════════════")
}
