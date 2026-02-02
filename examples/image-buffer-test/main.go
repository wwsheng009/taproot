package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render/buffer"
)

func main() {
	// Create a buffer with terminal size
	width, height := 80, 24
	screenBuffer := buffer.NewBuffer(width, height)

	// Load image path
	imgPath := "placeholder.png"
	if len(os.Args) > 1 {
		imgPath = os.Args[1]
	}

	fmt.Println("============================================================================")
	fmt.Println("                   Image Buffer Rendering Test")
	fmt.Println("============================================================================")
	fmt.Println()

	// Test 1: Render ASCII art to buffer
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[Test 1] ASCII Art to Buffer")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	screenBuffer.FillRect(buffer.Rect{X: 0, Y: 0, Width: width, Height: height}, ' ', buffer.Style{})

	// Draw a simple ASCII art image
	asciiArt := `
    .--.
   |__  |
.--.'  '.
|    '   |
|  __|   |
|_  \.___|
   \    |
    \   |
     \__|
`
	screenBuffer.WriteString(buffer.Point{X: 5, Y: 2}, imgPath, buffer.Style{Foreground: "#86", Bold: true})
	screenBuffer.WriteString(buffer.Point{X: 5, Y: 4}, asciiArt, buffer.Style{Foreground: "#196"})

	fmt.Println(screenBuffer.Render())
	fmt.Println()

	// Test 2: Render color blocks to buffer
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[Test 2] Color Blocks Test")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	screenBuffer.FillRect(buffer.Rect{X: 0, Y: 0, Width: width, Height: height}, ' ', buffer.Style{})

	// Draw color blocks
	screenBuffer.WriteString(buffer.Point{X: 2, Y: 2}, "Red Block:", buffer.Style{Bold: true})
	screenBuffer.FillRect(buffer.Rect{X: 12, Y: 2, Width: 10, Height: 3}, '#', buffer.Style{Foreground: "#ff0000"})

	screenBuffer.WriteString(buffer.Point{X: 2, Y: 6}, "Green Block:", buffer.Style{Bold: true})
	screenBuffer.FillRect(buffer.Rect{X: 14, Y: 6, Width: 10, Height: 3}, '#', buffer.Style{Foreground: "#00ff00"})

	screenBuffer.WriteString(buffer.Point{X: 2, Y: 10}, "Blue Block:", buffer.Style{Bold: true})
	screenBuffer.FillRect(buffer.Rect{X: 13, Y: 10, Width: 10, Height: 3}, '#', buffer.Style{Foreground: "#0000ff"})

	screenBuffer.WriteString(buffer.Point{X: 2, Y: 14}, "Yellow Block:", buffer.Style{Bold: true})
	screenBuffer.FillRect(buffer.Rect{X: 15, Y: 14, Width: 10, Height: 3}, '#', buffer.Style{Foreground: "#ffff00"})

	screenBuffer.WriteString(buffer.Point{X: 2, Y: 18}, "Purple Block:", buffer.Style{Bold: true})
	screenBuffer.FillRect(buffer.Rect{X: 16, Y: 18, Width: 10, Height: 3}, '#', buffer.Style{Foreground: "#ff00ff"})

	fmt.Println(screenBuffer.Render())
	fmt.Println()

	// Test 3: Render gradient to buffer
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[Test 3] Gradient Test")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	screenBuffer.FillRect(buffer.Rect{X: 0, Y: 0, Width: width, Height: height}, ' ', buffer.Style{})

	// Draw gradient from red to blue
	for x := 0; x < 40; x++ {
		r := 255 - (x * 255 / 40)
		g := 0
		b := x * 255 / 40
		color := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		screenBuffer.FillRect(buffer.Rect{X: x, Y: 2, Width: 1, Height: 5}, '#', buffer.Style{Foreground: color})
	}

	screenBuffer.WriteString(buffer.Point{X: 2, Y: 8}, "Red -> Blue Gradient", buffer.Style{Bold: true})

	// Draw gradient from green to yellow
	for x := 0; x < 40; x++ {
		r := x * 255 / 40
		g := 255
		b := 0
		color := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		screenBuffer.FillRect(buffer.Rect{X: x, Y: 10, Width: 1, Height: 5}, '#', buffer.Style{Foreground: color})
	}

	screenBuffer.WriteString(buffer.Point{X: 2, Y: 16}, "Green -> Yellow Gradient", buffer.Style{Bold: true})

	fmt.Println(screenBuffer.Render())
	fmt.Println()

	// Test 4: Wide character rendering with buffer
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[Test 4] Wide Characters with Buffer")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	screenBuffer.FillRect(buffer.Rect{X: 0, Y: 0, Width: width, Height: height}, ' ', buffer.Style{})

	texts := []struct {
		text  string
		style buffer.Style
	}{
		{"中文测试: 你好世界", buffer.Style{Foreground: "#196", Bold: true}},
		{"日文测试: こんにちは", buffer.Style{Foreground: "#208", Bold: true}},
		{"韩文测试: 안녕하세요", buffer.Style{Foreground: "#226", Bold: true}},
		{"Mixed: Hello 世界 こんにちは 안녕하세요", buffer.Style{Foreground: "#246", Bold: true}},
	}

	for i, t := range texts {
		screenBuffer.WriteString(buffer.Point{X: 2, Y: 2 + i*2}, t.text, t.style)
	}

	fmt.Println(screenBuffer.Render())
	fmt.Println()

	// Test 5: Line wrapping test
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[Test 5] Text Wrapping with Wide Characters")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	screenBuffer.FillRect(buffer.Rect{X: 0, Y: 0, Width: width, Height: height}, ' ', buffer.Style{})

	longText := "这是一个很长的中文文本，测试Buffer的自动换行功能是否能够正确处理宽字符。Buffer的WriteStringWrapped方法应该能够根据maxWidth参数自动将文本分行显示，并且每个汉字占用2个列宽。The quick brown fox jumps over the lazy dog."

	screenBuffer.WriteString(buffer.Point{X: 2, Y: 2}, "长文本自动换行测试:", buffer.Style{Foreground: "#86", Bold: true})
	screenBuffer.WriteStringWrapped(buffer.Point{X: 2, Y: 4}, 60, longText, buffer.Style{Foreground: "#246"})

	fmt.Println(screenBuffer.Render())
	fmt.Println()

	// Test 6: Performance benchmark
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[Test 6] Performance Benchmark")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	iterations := 1000

	// Benchmark 1: Simple fill
	testBuffer := buffer.NewBuffer(width, height)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		testBuffer.FillRect(buffer.Rect{X: 0, Y: 0, Width: width, Height: height}, ' ', buffer.Style{})
	}
	elapsed := time.Since(start)
	avgNs := elapsed.Nanoseconds() / int64(iterations)
	fps := int64(1000000000 / avgNs)
	fmt.Printf("测试 1: FillRect\n")
	fmt.Printf("  迭代次数: %d\n", iterations)
	fmt.Printf("  总耗时: %v\n", elapsed)
	fmt.Printf("  平均耗时: %d ns (%.3f μs)\n", avgNs, float64(avgNs)/1000)
	fmt.Printf("  理论 FPS: %d\n\n", fps)

	// Benchmark 2: WriteString
	start = time.Now()
	testText := "Buffer rendering test with 中文 characters 世界"
	for i := 0; i < iterations; i++ {
		testBuffer.FillRect(buffer.Rect{X: 0, Y: 0, Width: width, Height: height}, ' ', buffer.Style{})
		testBuffer.WriteString(buffer.Point{X: 2, Y: 5}, testText, buffer.Style{Foreground: "#86"})
	}
	elapsed = time.Since(start)
	avgNs = elapsed.Nanoseconds() / int64(iterations)
	fps = int64(1000000000 / avgNs)
	fmt.Printf("测试 2: WriteString with wide chars\n")
	fmt.Printf("  迭代次数: %d\n", iterations)
	fmt.Printf("  总耗时: %v\n", elapsed)
	fmt.Printf("  平均耗时: %d ns (%.3f μs)\n", avgNs, float64(avgNs)/1000)
	fmt.Printf("  理论 FPS: %d\n\n", fps)

	// Benchmark 3: Render
	start = time.Now()
	for i := 0; i < iterations; i++ {
		testBuffer.Render()
	}
	elapsed = time.Since(start)
	avgNs = elapsed.Nanoseconds() / int64(iterations)
	fps = int64(1000000000 / avgNs)
	fmt.Printf("测试 3: Render\n")
	fmt.Printf("  迭代次数: %d\n", iterations)
	fmt.Printf("  总耗时: %v\n", elapsed)
	fmt.Printf("  平均耗时: %d ns (%.3f μs)\n", avgNs, float64(avgNs)/1000)
	fmt.Printf("  理论 FPS: %d\n", fps)

	// Test FPS requirements
	fmt.Println("\n60 FPS 要求: 每帧 <= 16,666,667 ns (16.67ms)")
	fmt.Println("30 FPS 要求: 每帧 <= 33,333,333 ns (33.33ms)")

	// Summary
	fmt.Println()
	fmt.Println("============================================================================")
	fmt.Println(fmt.Sprintf("                               %s 结果", lipgloss.NewStyle().Foreground(lipgloss.Color("#46f")).Bold(true).Render("所有测试完成")))
	fmt.Println("============================================================================")
	fmt.Println()
	fmt.Println("测试说明:")
	fmt.Println("  • Test 1: 使用Buffer渲染ASCII艺术图像")
	fmt.Println("  • Test 2: 使用Buffer渲染彩色方块")
	fmt.Println("  • Test 3: 使用Buffer渲染渐变色")
	fmt.Println("  • Test 4: 宽字符渲染测试")
	fmt.Println("  • Test 5: 文本自动换行测试")
	fmt.Println("  • Test 6: 性能基准测试")
	fmt.Println()
	fmt.Println("所有测试均使用Buffer系统完成，未依赖其他UI组件。")
}
