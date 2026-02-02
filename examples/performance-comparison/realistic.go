package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/wwsheng009/taproot/ui/render/buffer"
)

// Scenario 1: Simple layout (real TUI usage, single frame)
func scenario1_StringBased(width, height int) string {
	// Single frame render - typical TUI usage
	var sb strings.Builder

	// Header
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteByte('\n')
	sb.WriteString("HEADER\n")
	sb.WriteString(strings.Repeat(" ", (width-10)/2))
	sb.WriteString("Line 2\n")
	sb.WriteString(strings.Repeat(" ", (width-10)/2))
	sb.WriteString("Line 3\n")
	sb.WriteString(strings.Repeat(" ", (width-10)/2))
	sb.WriteString("Line 4\n")
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteByte('\n')

	// Content area
	contentHeight := height - 7
	for i := 0; i < contentHeight; i++ {
		sb.WriteString("|")
		sb.WriteString(strings.Repeat(" ", width-2))
		sb.WriteString("|\n")
	}

	// Footer
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteByte('\n')
	sb.WriteString(strings.Repeat(" ", (width-20)/2))
	sb.WriteString("Footer\n")

	return sb.String()
}

func scenario1_BufferBased(width, height int) string {
	// Single frame render with buffer - typical TUI usage
	lm := buffer.NewLayoutManager(width, height)

	header := buffer.NewTextComponent(
		"HEADER\nLine 2\nLine 3\nLine 4",
		buffer.Style{Foreground: "202"},
	).SetCenterH(true)

	content := buffer.NewFillComponent(' ', buffer.Style{})

	footer := buffer.NewTextComponent(
		"Footer",
		buffer.Style{Foreground: "244"},
	).SetCenterH(true)

	lm.AddComponent("header", header)
	lm.AddComponent("content", content)
	lm.AddComponent("footer", footer)

	lm.CalculateLayout()
	return lm.Render()
}

// Scenario 2: Partial update (only content changes)
func scenario2_StringBased(width, height, frame int) string {
	// Only content changes, header/footer same
	var sb strings.Builder

	// Header (same every frame)
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteByte('\n')
	sb.WriteString("HEADER (static)\n")
	sb.WriteString(strings.Repeat(" ", (width-12)/2))
	sb.WriteString("Line 2\n")
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteByte('\n')

	// Content (changes every frame)
	contentHeight := height - 6
	for i := 0; i < contentHeight; i++ {
		sb.WriteString("|")
		content := fmt.Sprintf(" Frame %d Line %d ", frame, i)
		padding := width - 2 - len(content)
		if padding < 0 {
			content = content[:width-2]
			padding = 0
		}
		sb.WriteString(content)
		sb.WriteString(strings.Repeat(" ", padding))
		sb.WriteString("|\n")
	}

	// Footer (same every frame)
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteByte('\n')

	return sb.String()
}

type Scenario2Buffer struct {
	lm     *buffer.LayoutManager
	header buffer.Renderable
	footer buffer.Renderable
}

func newScenario2Buffer(width, height int) *Scenario2Buffer {
	lm := buffer.NewLayoutManager(width, height)

	header := buffer.NewTextComponent(
		"HEADER (static)\nLine 2",
		buffer.Style{Foreground: "202"},
	).SetCenterH(true)

	footer := buffer.NewTextComponent(
		fmt.Sprintf("Footer"),
		buffer.Style{Foreground: "244"},
	).SetCenterH(true)

	lm.AddComponent("header", header)
	lm.AddComponent("footer", footer)
	lm.CalculateLayout()

	return &Scenario2Buffer{lm: lm, header: header, footer: footer}
}

func (s *Scenario2Buffer) render(frame int, width, height int) string {
	// Only content changes - can reuse layout
	content := buffer.NewTextComponent(
		fmt.Sprintf("Frame %d Content", frame),
		buffer.Style{Foreground: "246"},
	).SetCenterV(true).SetCenterH(true)

	s.lm.AddComponent("content", content)
	s.lm.SetSize(width, height)
	s.lm.CalculateLayout()
	return s.lm.Render()
}

// Scenario 3: Complex layout with multiple components
func scenario3_StringBased(width, height int) string {
	var sb strings.Builder

	// Header with border
	sb.WriteString("┌")
	sb.WriteString(strings.Repeat("─", width-2))
	sb.WriteString("┐\n")
	sb.WriteString("│")
	sb.WriteString(strings.Repeat(" ", (width-20)/2))
	sb.WriteString("Complex Layout")
	sb.WriteString(strings.Repeat(" ", width-2-(width-20)/2-13))
	sb.WriteString("│\n")
	sb.WriteString("╞")
	sb.WriteString(strings.Repeat("═", width-2))
	sb.WriteString("╡\n")

	// Two columns
	col1Width := width/2 - 2
	col2Width := width - col1Width - 4

	contentHeight := height - 4
	for i := 0; i < contentHeight; i++ {
		sb.WriteString("│")
		sb.WriteString(strings.Repeat(" ", col1Width))
		sb.WriteString("│")
		sb.WriteString(strings.Repeat(" ", col2Width))
		sb.WriteString("│\n")
	}

	// Footer
	sb.WriteString("└")
	sb.WriteString(strings.Repeat("─", width-2))
	sb.WriteString("┘\n")

	return sb.String()
}

func scenario3_BufferBased(width, height int) string {
	lm := buffer.NewLayoutManager(width, height)

	// Header
	header := buffer.NewTextComponent(
		"Complex Layout",
		buffer.Style{Foreground: "81", Bold: true},
	).SetCenterH(true)

	lm.AddComponent("header", header)
	lm.AddComponent("content", buffer.NewFillComponent(' ', buffer.Style{}))
	lm.AddComponent("footer", buffer.NewTextComponent("", buffer.Style{}))

	lm.CalculateLayout()
	return lm.Render()
}

// Benchmark function
func benchmark(name string, fn func() string, iterations int) (time.Duration, int, int) {
	runtime.GC() // Force GC before benchmark
	start := time.Now()

	var totalSize int
	for i := 0; i < iterations; i++ {
		output := fn()
		totalSize += len(output)
	}

	return time.Since(start), iterations, totalSize / iterations
}

func printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func main() {
	width := 80
	height := 24
	iterations := 1000

	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println("Realistic TUI Performance Comparison")
	fmt.Printf("Terminal: %dx%d, Iterations: %d\n", width, height, iterations)
	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println()

	// Scenario 1: Single frame render
	fmt.Println("SCENARIO 1: Single Frame Render (typical TUI usage)")
	fmt.Println(strings.Repeat("-", 80))

	// Warm up
	scenario1_StringBased(width, height)
	scenario1_BufferBased(width, height)

	// Test string-based
	time1, _, size1 := benchmark("Scenario1", func() string {
		return scenario1_StringBased(width, height)
	}, iterations)
	avg1 := time.Duration(int64(time1) / int64(iterations))

	// Test buffer-based
	time2, _, size2 := benchmark("Scenario1", func() string {
		return scenario1_BufferBased(width, height)
	}, iterations)
	avg2 := time.Duration(int64(time2) / int64(iterations))

	// Results
	fmt.Printf("String-based:  %v total, %v avg, %d bytes/output\n", time1, avg1, size1)
	fmt.Printf("Buffer-based:  %v total, %v avg, %d bytes/output\n", time2, avg2, size2)

	ratio := float64(time2) / float64(time1)
	if ratio > 1.0 {
		fmt.Printf("Buffer-based: %.2fx slower (+%.1f%%)\n", ratio, (ratio-1.0)*100)
	} else {
		fmt.Printf("Buffer-based: %.2fx faster (-%.1f%%)\n", 1.0/ratio, (1.0-ratio)*100)
	}

	// Frame rate analysis
	fpsString := float64(iterations) / time1.Seconds()
	fpsBuffer := float64(iterations) / time2.Seconds()
	fmt.Printf("FPS Potential: String=%d, Buffer=%d\n", int(fpsString), int(fpsBuffer))

	// Scenario 2: Partial update
	fmt.Println()
	fmt.Println("SCENARIO 2: Partial Update (only content changes)")
	fmt.Println(strings.Repeat("-", 80))

	scenario2 := newScenario2Buffer(width, height)

	// Warm up
	for i := 0; i < 10; i++ {
		scenario2_StringBased(width, height, i)
		scenario2.render(i, width, height)
	}

	// Test string-based
	time3, _, size3 := benchmark("Scenario2", func() string {
		static := 0
		return func() string {
			defer func() { static++ }()
			return scenario2_StringBased(width, height, static)
		}()
	}, iterations)
	avg3 := time.Duration(int64(time3) / int64(iterations))

	// Test buffer-based
	time4, _, size4 := benchmark("Scenario2", func() string {
		static := 0
		return func() string {
			defer func() { static++ }()
			return scenario2.render(static, width, height)
		}()
	}, iterations)
	avg4 := time.Duration(int64(time4) / int64(iterations))

	// Results
	fmt.Printf("String-based:  %v total, %v avg, %d bytes/output\n", time3, avg3, size3)
	fmt.Printf("Buffer-based:  %v total (reusing layout), %v avg, %d bytes/output\n",
		time4, avg4, size4)

	ratio2 := float64(time4) / float64(time3)
	if ratio2 > 1.0 {
		fmt.Printf("Buffer-based: %.2fx slower (+%.1f%%)\n", ratio2, (ratio2-1.0)*100)
	} else {
		fmt.Printf("Buffer-based: %.2fx faster (-%.1f%%)\n", 1.0/ratio2, (1.0-ratio2)*100)
	}

	fpsString2 := float64(iterations) / time3.Seconds()
	fpsBuffer2 := float64(iterations) / time4.Seconds()
	fmt.Printf("FPS Potential: String=%d, Buffer=%d\n", int(fpsString2), int(fpsBuffer2))

	// Summary
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println("SUMMARY")
	fmt.Println("=" + strings.Repeat("=", 79))

	fmt.Println()
	fmt.Println("Frame Budget Analysis (60fps = 16.67ms budget):")
	fmt.Printf("  Scenario 1 (full render):\n")
	fmt.Printf("    String-based: %v (%.1f%% of budget) → %d FPS\n", avg1, float64(avg1)/16.67*100, int(fpsString))
	fmt.Printf("    Buffer-based: %v (%.1f%% of budget) → %d FPS\n", avg2, float64(avg2)/16.67*100, int(fpsBuffer))

	fmt.Printf("  Scenario 2 (partial update):\n")
	fmt.Printf("    String-based: %v (%.1f%% of budget) → %d FPS\n", avg3, float64(avg3)/16.67*100, int(fpsString2))
	fmt.Printf("    Buffer-based: %v (%.1f%% of budget) → %d FPS\n", avg4, float64(avg4)/16.67*100, int(fpsBuffer2))

	fmt.Println()
	fmt.Println("CONCLUSION:")
	if float64(avg2)/16.67 < 50.0 {
		fmt.Println("  ✓ Buffer-based rendering is fast enough for 60fps TUI applications")
		fmt.Println("  ✓ Overhead is well within the 16.67ms frame budget")
	} else if float64(avg2)/16.67 < 100.0 {
		fmt.Println("  ✓ Buffer-based rendering is acceptable for most TUI applications")
		fmt.Println("  → May need optimization for very complex layouts")
	} else {
		fmt.Println("  ⚠ Buffer-based rendering may be slow for complex TUI applications")
		fmt.Println("  → Consider: caching, partial updates, or hybrid approach")
	}

	fmt.Println()
	fmt.Println("BENEFITS OVERHEAD TRADE-OFF:")
	fmt.Println("  + Accurate dimension calculations (no estimation)")
	fmt.Println("  + Layout independence (content doesn't affect layout)")
	fmt.Println("  + Component isolation (no interference)")
	fmt.Println("  + Accurate Sixel image support")
	fmt.Println("  + Better debugging (inspect buffer state)")
	fmt.Println("  - Performance overhead (2-5x slower for realistic usage)")
	fmt.Printf("  → For typical TUI apps: %s\n", map[bool]string{
		true:  "RECOMMENDED (benefits > overhead)",
		false: "Consider carefully (benefits vs overhead)",
	}[float64(avg2)/16.67 < 50.0])

	fmt.Println("=" + strings.Repeat("=", 79))
}
