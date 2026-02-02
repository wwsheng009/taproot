# Image Buffer Example

使用 Buffer Layout 系统构建的图像查看器示例。

## 特性

本示例展示了如何使用 Buffer 的 LayoutManager 来构建复杂的 UI：

- **Buffer Layout 系统**：使用 `buffer.LayoutManager` 进行组件布局
- **组件化设计**：Header、Footer、Content 三个独立组件
- **Renderable 接口**：实现 `Render()`、`MinSize()`、`PreferredSize()` 方法
- **动态布局**：根据窗口大小自动调整布局
- **图像渲染器模拟**：支持多种渲染模式（Kitty、iTerm2、Sixel、Blocks、ASCII）

## 运行

```bash
go run main.go [image-path]
```

或编译后运行：

```bash
go build -o image-buffer.exe main.go
./image-buffer.exe path/to/image.png
```

## 键盘快捷键

| 按键 | 功能 |
|------|------|
| `1-6` | 切换渲染器（Auto/Kitty/iTerm2/Blocks/Sixel/ASCII） |
| `+/-` | 放大/缩小 |
| `0` | 重置缩放 |
| `m` | 切换缩放模式（Fit/Fill） |
| `r` | 重新加载图像 |
| `s` | 切换图像路径 |
| `q` / `Ctrl+C` | 退出 |

## 架构说明

### LayoutManager

```go
layoutMgr := buffer.NewLayoutManager(width, height)

// 添加组件
layoutMgr.AddComponent("header", header)
layoutMgr.AddComponent("footer", footer)
layoutMgr.AddComponent("content", content)

// 计算布局（支持图像模式）
layoutMgr.ImageLayout(contentHeightHint)

// 渲染
output := layoutMgr.Render()
```

### Renderable 接口

每个组件必须实现 `buffer.Renderable` 接口：

```go
type Renderable interface {
    Render(buf *Buffer, rect Rect)           // 渲染组件
    MinSize() (int, int)                     // 最小尺寸
    PreferredSize() (int, int)               // 首选尺寸
}
```

### 组件示例

```go
type Header struct {
    imgPath     string
    renderer    RendererType
    imageLoaded bool
    // ...
}

func (h *Header) Render(buf *buffer.Buffer, rect buffer.Rect) {
    // 渲染标题和图像信息
    buf.WriteString(buffer.Point{X: 2, Y: 0}, "Title", buffer.Style{...})
}

func (h *Header) MinSize() (int, int) {
    return 60, 5
}

func (h *Header) PreferredSize() (int, int) {
    return 80, 6
}
```

## 与原始 Image 示例的对比

| 特性 | 原始 Image | Buffer Layout 版本 |
|------|------------|-------------------|
| 布局系统 | lipgloss 垂直布局 | Buffer LayoutManager |
| 组件接口 | 自定义 View() | buffer.Renderable |
| 内存管理 | 字符串拼接 | Buffer 对象池 |
| 宽字符支持 | 依赖 lipgloss | Buffer 内置支持 |
| 性能优化 | 无 | 对象池 + 样式缓存 |

## 相关示例

- `examples/image/` - 使用 lipgloss 布局的原始图像查看器
- `examples/image-buffer-test/` - Buffer 图形渲染测试
- `examples/layout-test/` - LayoutManager 使用示例
