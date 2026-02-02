package main

import (
	"fmt"

	"github.com/wwsheng009/taproot/ui/render/buffer"
)

func main() {
	fmt.Println("=== 宽字符测试 ===")
	fmt.Println()

	// 测试 1: 纯中文
	fmt.Println("[测试 1] 纯中文文本")
	buf := buffer.NewBuffer(40, 3)
	style := buffer.Style{Foreground: "202"}
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "你好，世界！", style)
	buf.WriteString(buffer.Point{X: 0, Y: 1}, "欢迎使用 Taproot TUI 框架", style)
	fmt.Println(buf.Render())
	fmt.Println()

	// 测试 2: 混合文本
	fmt.Println("[测试 2] 中英文混合")
	buf = buffer.NewBuffer(50, 4)
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "Hello, 世界!", style)
	buf.WriteString(buffer.Point{X: 0, Y: 1}, "这是中英文混合的例子", style)
	buf.WriteString(buffer.Point{X: 0, Y: 2}, "Mixed text: 你好world", style)
	fmt.Println(buf.Render())
	fmt.Println()

	// 测试 3: 日文
	fmt.Println("[测试 3] 日文文本")
	buf = buffer.NewBuffer(40, 2)
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "こんにちは、世界", style)
	buf.WriteString(buffer.Point{X: 0, Y: 1}, "日本語のテスト", style)
	fmt.Println(buf.Render())
	fmt.Println()

	// 测试 4: 带样式的文本
	fmt.Println("[测试 4] 带样式的标题")
	buf = buffer.NewBuffer(50, 3)
	titleStyle := buffer.Style{Foreground: "201", Bold: true}
	buf.WriteString(buffer.Point{X: 0, Y: 0}, "终端用户界面开发工具集", titleStyle)
	bodyStyle := buffer.Style{Foreground: "244"}
	buf.WriteString(buffer.Point{X: 0, Y: 2}, "支持多语言字符的缓冲区渲染系统", bodyStyle)
	fmt.Println(buf.Render())
	fmt.Println()

	fmt.Println("✅ 所有测试完成！")
}
