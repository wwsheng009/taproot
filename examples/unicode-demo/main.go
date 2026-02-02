package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/wwsheng009/taproot/ui/render/buffer"
)

func main() {
	fmt.Println("============================================================================")
	fmt.Println("                   Buffer 宽字符支持综合测试")
	fmt.Println("============================================================================")
	fmt.Println()

	// 测试 1: 基本宽字符渲染
	testBasicWideChars()

	// 测试 2: 中英文混合
	testMixedText()

	// 测试 3: 多语言支持
	testMultipleLanguages()

	// 测试 4: 样式应用
	testStyledText()

	// 测试 5: 边界情况
	testEdgeCases()

	// 测试 6: 布局管理器
	testLayoutManager()

	// 测试 7: 性能测试
	testPerformance()

	// 测试 8: 实际场景
	testRealWorldScenario()

	fmt.Println("============================================================================")
	fmt.Println("                           ✅ 所有测试完成！")
	fmt.Println("============================================================================")
}

func testBasicWideChars() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[测试 1] 基本宽字符渲染")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	buf := buffer.NewBuffer(50, 10)
	style := buffer.Style{Foreground: "202"}

	// 中文
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "中文测试: 你好世界", style)
	buf.WriteString(buffer.Point{X: 0, Y: 1}, "每个汉字占用 2 个列宽", style)

	// 日文
	buf.WriteString(buffer.Point{X: 0, Y: 2}, "日文测试: こんにちは", style)
	buf.WriteString(buffer.Point{X: 0, Y: 3}, "日本語: 平仮名と漢字", style)

	// 韩文
	buf.WriteString(buffer.Point{X: 0, Y: 4}, "韩文测试: 안녕하세요", style)
	buf.WriteString(buffer.Point{X: 0, Y: 5}, "한글: 글자수 계산", style)

	// 宽度验证
	cols := buf.WriteString(buffer.Point{X: 0, Y: 7}, "测试宽度: 你好abc", buffer.Style{Foreground: "201"})
	buf.WriteString(buffer.Point{X: 0, Y: 8}, fmt.Sprintf("占用列数: %d (2汉字+3字母 = 4+3 = 7)", cols), buffer.Style{})

	fmt.Println(buf.Render())
	fmt.Println()
}

func testMixedText() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[测试 2] 中英文混合文本")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	buf := buffer.NewBuffer(60, 8)
	style := buffer.Style{Foreground: "202"}

	// 混合文本示例
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "Hello, 世界! 这是一个中英文混合的例子。", style)
	buf.WriteString(buffer.Point{X: 0, Y: 1}, "Python编程语言支持多种编码格式。", style)
	buf.WriteString(buffer.Point{X: 0, Y: 2}, "Go语言: fmt.Println(\"Hello, 世界\")", style)
	buf.WriteString(buffer.Point{X: 0, Y: 3}, "Terminal User Interface (TUI) with CJK support", style)
	buf.WriteString(buffer.Point{X: 0, Y: 4}, "终端用户界面开发 (Terminal UI Development)", style)

	// 统计
	buf.WriteString(buffer.Point{X: 0, Y: 6}, "✓ 英文: 1列/字符", buffer.Style{Foreground: "32"})
	buf.WriteString(buffer.Point{X: 0, Y: 7}, "✓ 中文: 2列/字符", buffer.Style{Foreground: "32"})

	fmt.Println(buf.Render())
	fmt.Println()
}

func testMultipleLanguages() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[测试 3] 多语言支持")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	buf := buffer.NewBuffer(70, 12)
	style := buffer.Style{Foreground: "202"}

	// 各种语言
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "中文: 你好，世界！欢迎使用 Taproot TUI 框架", style)
	buf.WriteString(buffer.Point{X: 0, Y: 1}, "日文: こんにちは、世界。Taproot TUI フレームワークへようこそ", style)
	buf.WriteString(buffer.Point{X: 0, Y: 2}, "韩文: 안녕하세요, 세계! Taproot TUI 프레임워크에 오신 것을 환영합니다", style)
	buf.WriteString(buffer.Point{X: 0, Y: 3}, "繁體中文: 你好，世界！歡迎使用 Taproot TUI 框架", style)
	buf.WriteString(buffer.Point{X: 0, Y: 4}, "English: Hello, World! Welcome to Taproot TUI Framework", style)
	buf.WriteString(buffer.Point{X: 0, Y: 5}, "Español: ¡Hola, Mundo! Bienvenido al marco TUI de Taproot", style)
	buf.WriteString(buffer.Point{X: 0, Y: 6}, "Français: Bonjour, le Monde! Bienvenue dans le framework TUI Taproot", style)
	buf.WriteString(buffer.Point{X: 0, Y: 7}, "Deutsch: Hallo, Welt! Willkommen beim Taproot TUI-Framework", style)
	buf.WriteString(buffer.Point{X: 0, Y: 8}, "Русский: Привет, мир! Добро пожаловать в фреймворк Taproot TUI", style)
	buf.WriteString(buffer.Point{X: 0, Y: 9}, "العربية: مرحبا بالعالم! مرحبا بكم في إطار TUI Taproot", style)

	// 说明
	buf.WriteString(buffer.Point{X: 0, Y: 11}, "✓ 支持多种语言字符，正确计算列宽", buffer.Style{Foreground: "32"})

	fmt.Println(buf.Render())
	fmt.Println()
}

func testStyledText() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[测试 4] 样式应用")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	buf := buffer.NewBuffer(70, 10)

	// 标题 - 粗体红色
	titleStyle := buffer.Style{Foreground: "196", Bold: true}
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "Taproot TUI 框架 - 终端用户界面开发工具集", titleStyle)

	// 副标题 - 斜体蓝色
	subtitleStyle := buffer.Style{Foreground: "33", Italic: true}
	buf.WriteString(buffer.Point{X: 0, Y: 2}, "基于 Bubbletea 的高性能 TUI 框架", subtitleStyle)

	// 功能列表
	bodyStyle := buffer.Style{Foreground: "246"}
	buf.WriteString(buffer.Point{X: 2, Y: 4}, "• Buffer 渲染系统", bodyStyle)
	buf.WriteString(buffer.Point{X: 2, Y: 5}, "• 宽字符支持（CJK）", bodyStyle)
	buf.WriteString(buffer.Point{X: 2, Y: 6}, "• 组件化设计", bodyStyle)
	buf.WriteString(buffer.Point{X: 2, Y: 7}, "• 高性能优化", bodyStyle)

	// 强调文本
	emphasisStyle := buffer.Style{Foreground: "202", Underline: true}
	buf.WriteString(buffer.Point{X: 0, Y: 9}, "关键特性: 准确的高度/宽度计算", emphasisStyle)

	fmt.Println(buf.Render())
	fmt.Println()
}

func testEdgeCases() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[测试 5] 边界情况处理")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	buf := buffer.NewBuffer(40, 12)
	style := buffer.Style{Foreground: "202"}

	// 测试 1: 行尾宽字符
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "行尾测试: abcdefghijklmnopqrst", style)
	buf.WriteString(buffer.Point{X: 0, Y: 1}, "行尾宽字一二三四五六七八九", style)

	// 使用分隔线
	borderStyle := buffer.Style{Foreground: "240"}
	for i := 0; i < 40; i++ {
		buf.SetCell(buffer.Point{X: i, Y: 2}, buffer.Cell{Char: '─', Width: 1, Style: borderStyle})
		buf.SetCell(buffer.Point{X: i, Y: 5}, buffer.Cell{Char: '─', Width: 1, Style: borderStyle})
		buf.SetCell(buffer.Point{X: i, Y: 8}, buffer.Cell{Char: '─', Width: 1, Style: borderStyle})
	}

	// 测试 2: 换行处理
	wrappedStyle := buffer.Style{Foreground: "201"}
	buf.WriteString(buffer.Point{X: 0, Y: 3}, "自动换行测试:", wrappedStyle)
	buf.WriteStringWrapped(buffer.Point{X: 0, Y: 4}, 40, "这是一个很长的文本，测试自动换行功能是否正确处理宽字符，确保中文和英文都能正确换行。", style)

	// 测试 3: 混合换行
	buf.WriteString(buffer.Point{X: 0, Y: 6}, "手动换行测试:", wrappedStyle)
	buf.WriteString(buffer.Point{X: 0, Y: 7}, "第一行\n第二行\n第三行", style)

	// 测试 4: 非常长的文本
	buf.WriteString(buffer.Point{X: 0, Y: 9}, "超长文本截断:", wrappedStyle)
	buf.WriteString(buffer.Point{X: 0, Y: 10}, "这是一段非常长的文本，应该被截断处理，因为缓冲区的宽度是有限的，需要正确处理边界情况", style)

	fmt.Println(buf.Render())
	fmt.Println()
}

func testLayoutManager() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[测试 6] 布局管理器与宽字符")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	lm := buffer.NewLayoutManager(60, 15)

	// Header
	header := buffer.NewTextComponent(
		"╔══════════════════════════════════╗\n"+
			"║   Taproot TUI 框架演示程序          ║\n"+
			"║   支持宽字符的布局管理器              ║\n"+
			"╚══════════════════════════════════╝",
		buffer.Style{Foreground: "208", Bold: true},
	).SetCenterH(true)

	// Content
	content := buffer.NewTextComponent(
		"这是一个使用布局管理器的示例。\n\n"+
			"支持中文、日文、韩文等多种语言字符。\n\n"+
			"Header: 固定 5 行高度\nContent: 中间区域\nFooter: 固定 1 行高度",
		buffer.Style{Foreground: "244"},
	)

	// Footer
	footer := buffer.NewTextComponent(
		"按 Ctrl+C 退出 | 按 R 重新渲染",
		buffer.Style{Foreground: "202", Italic: true},
	).SetCenterH(true)

	lm.AddComponent("header", header)
	lm.AddComponent("content", content)
	lm.AddComponent("footer", footer)

	lm.CalculateLayout()
	fmt.Println(lm.Render())
	fmt.Println()
}

func testPerformance() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[测试 7] 性能测试")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	width := 80
	height := 24
	iterations := 1000

	// 纯英文
	fmt.Println("测试 1: 纯英文文本")
	testPerf(fmt.Sprintf("Hello, this is a test string for benchmarking the buffer rendering performance. "), width, height, iterations)

	// 纯中文
	fmt.Println("\n测试 2: 纯中文文本")
	testPerf("这是一个用于测试缓冲区渲染性能的字符串，包含大量中文字符。", width, height, iterations)

	// 混合文本
	fmt.Println("\n测试 3: 中英文混合")
	testPerf("This is a mixed text with some Chinese characters 这是一个混合文本包含一些中文字符。", width, height, iterations)

	// 复杂样式
	fmt.Println("\n测试 4: 复杂样式文本")
	testStyledPerf(width, height, iterations)
}

func testPerf(text string, width, height, iterations int) {
	wrappedText := strings.Repeat(text, 10)

	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		buf := buffer.GetBuffer(width, height)
		style := buffer.Style{Foreground: "202"}
		buf.WriteStringWrapped(buffer.Point{X: 0, Y: 0}, width, wrappedText, style)
		_ = buf.Render()
		buffer.PutBuffer(buf)
	}
	elapsed := time.Since(startTime)

	avgTime := elapsed.Nanoseconds() / int64(iterations)
	fps := float64(time.Second.Nanoseconds()) / float64(avgTime)

	fmt.Printf("  渲染次数: %d\n", iterations)
	fmt.Printf("  总耗时: %v\n", elapsed)
	fmt.Printf("  平均耗时: %d ns (%.3f μs)\n", avgTime, float64(avgTime)/1000.0)
	fmt.Printf("  理论 FPS: %.0f\n", fps)
}

func testStyledPerf(width, height, iterations int) {
	styles := []buffer.Style{
		{Foreground: "196", Bold: true},
		{Foreground: "202", Italic: true},
		{Foreground: "226", Underline: true},
		{Foreground: "38", Bold: true, Underline: true},
	}

	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		buf := buffer.GetBuffer(width, height)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				style := styles[(x+y)%4]
				buf.FillRect(buffer.Rect{
					X:      x,
					Y:      y,
					Width:  1,
					Height: 1,
				}, 'A', style)
			}
		}

		_ = buf.Render()
		buffer.PutBuffer(buf)
	}
	elapsed := time.Since(startTime)

	avgTime := elapsed.Nanoseconds() / int64(iterations)
	fps := float64(time.Second.Nanoseconds()) / float64(avgTime)

	fmt.Printf("  渲染次数: %d\n", iterations)
	fmt.Printf("  总耗时: %v\n", elapsed)
	fmt.Printf("  平均耗时: %d ns (%.3f μs)\n", avgTime, float64(avgTime)/1000.0)
	fmt.Printf("  理论 FPS: %.0f\n", fps)
}

func testRealWorldScenario() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("[测试 8] 真实场景: 一个简单的应用程序界面")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	lm := buffer.NewLayoutManager(60, 16)

	// 应用标题
	appTitle := buffer.NewTextComponent(
		"  🚀 代码编辑器 - Taproot TUI  ",
		buffer.Style{Foreground: "210", Bold: true},
	)

	// 菜单栏
	menuBar := buffer.NewTextComponent(
		" File  Edit  View  Tools  Help  ",
		buffer.Style{Foreground: "244", Bold: true},
	)

	// 状态栏
	statusBar := buffer.NewTextComponent(
		"就绪 | Ln 1, Col 1 | UTF-8 | 🇨🇳 中文支持 | 60 FPS",
		buffer.Style{Foreground: "33", Background: "234"},
	)

	// 内容区 - 模拟代码编辑
	codeContent := buffer.NewTextComponent(
		"  1 │ package main\n"+
			"  2 │ \n"+
			"  3 │ import \"fmt\"\n"+
			"  4 │ \n"+
			"  5 │ func main() {\n"+
			"  6 │     fmt.Println(\"你好，世界！\")\n"+
			"  7 │     fmt.Println(\"Hello, World!\")\n"+
			"  8 │ }\n"+
			"  9 │ \n"+
			" 10 │ // 支持中文注释",
		buffer.Style{},
	)

	// 侧边栏
	sideBar := buffer.NewTextComponent(
		"📁 项目\n"+
			"\n"+
			"  ├─ main.go\n"+
			"  ├─ util.go\n"+
			"  └─ README.md\n"+
			"\n"+
			"🔍 搜索\n"+
			"  输入关键词...",
		buffer.Style{Foreground: "244"},
	)

	lm.AddComponent("appTitle", appTitle)
	lm.AddComponent("menuBar", menuBar)
	lm.AddComponent("statusBar", statusBar)
	lm.AddComponent("codeContent", codeContent)
	lm.AddComponent("sideBar", sideBar)

	lm.CalculateLayout()
	fmt.Println(lm.Render())
	fmt.Println()
}
