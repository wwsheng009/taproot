# Image Buffer Test

使用 Buffer 渲染系统进行图像和图形渲染的测试示例。

## 功能

本示例展示了 Buffer 系统的图形渲染能力，包括：

- **Test 1**: ASCII 艺术图像渲染
- **Test 2**: 彩色方块渲染
- **Test 3**: 渐变色渲染
- **Test 4**: 宽字符渲染（中文、日文、韩文）
- **Test 5**: 文本自动换行测试
- **Test 6**: 性能基准测试

## 运行

```bash
go run main.go
```

或编译后运行：

```bash
go build -o image-buffer-test.exe main.go
./image-buffer-test.exe
```

## 性能测试结果

示例程序包含性能基准测试，评估 Buffer 系统的渲染性能：

- **FillRect**: ~10,000 FPS
- **WriteString**: ~9,800 FPS
- **Render**: ~39,000 FPS

所有测试均远超 60 FPS 的实时渲染要求。

## Buffer 渲染示例

```go
// 创建缓冲区
buf := buffer.NewBuffer(width, height)

// 填充背景
buf.FillRect(buffer.Rect{X: 0, Y: 0, Width: width, Height: height}, ' ', buffer.Style{})

// 写入文本（支持宽字符）
buf.WriteString(buffer.Point{X: 2, Y: 5}, "中文测试", buffer.Style{Foreground: "#86", Bold: true})

// 写入自动换行文本
buf.WriteStringWrapped(buffer.Point{X: 2, Y: 10}, 60, longText, buffer.Style{})

// 渲染输出
output := buf.Render()
fmt.Println(output)
```

## 技术特点

1. **纯 Buffer 系统**：不依赖其他 UI 组件，直接使用 Buffer API
2. **宽字符支持**：正确处理中文、日文、韩文等 CJK 字符
3. **ANSI 样式**：支持颜色、粗体、斜体等样式
4. **高性能**：使用对象池和样式缓存优化
5. **灵活性**：适合快速原型和图形渲染测试

## 相关示例

- `examples/image/` - 使用 Image 组件的完整图像查看器
- `examples/optimization-test/` - Buffer 性能优化测试
- `examples/unicode-demo/` - Unicode 宽字符综合测试
