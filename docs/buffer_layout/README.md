# Buffer Layout 文档

本目录包含 Taproot 框架中 Buffer-Based Rendering System 和 Layout System 的核心文档。

## 概述

Buffer-Based Rendering 是 Taproot 的核心技术，通过二维字符网格提供精确的布局计算，解决了传统 TUI 框架中基于字符串的布局计算问题。

## 核心优势

| 特性 | 字符串布局 | Buffer布局 |
|------|----------|-----------|
| 维度计算 | 估算（不可靠） | **精确（网格保证）** |
| 布局时机 | 渲染后 | **渲染前** |
| Sixel图像 | 猜测高度 | **精确高度** |
| 组件隔离 | 共享字符串 | **独立缓冲区** |
| 宽字符 | 复杂处理 | **原生支持** |

## 文档列表

### 📚 核心文档

#### [BUFFER_RENDERING.md](./BUFFER_RENDERING.md)
Buffer-Based Rendering 系统完整文档

**包含内容**：
- 系统概述与问题分析
- 核心数据结构（Cell, Buffer, Point, Rect）
- Buffer 操作（SetCell, FillRect, WriteString, WriteBuffer）
- LayoutManager 使用
- 组件系统（Renderable接口, TextComponent, ImageComponent）
- 性能基准测试结果
- 使用示例和最佳实践

**适用读者**：所有开发者，必读文档

---

#### [IMAGE_ZOOM_BUFFER_ARCHITECTURE.md](./IMAGE_ZOOM_BUFFER_ARCHITECTURE.md)
图像缩放与 Buffer Layer 架构详解（中文）

**包含内容**：
- 图像缩放系统设计
- Zoom Modes（Fit/Fill/Stretch/Original）
- Zoom Level 缩放控制
- 分辨率驱动的缩放算法
- Buffer Layer 布局系统
- 宽字符处理
- 渲染优化技术
- 多渲染器支持（Kitty, iTerm2, Sixel, Blocks, ASCII）
- 键盘快捷键
- 性能特性

**适用读者**：图像查看器开发者、深入理解系统者

---

#### [LAYOUT_FIX.md](./LAYOUT_FIX.md)
布局系统修复文档

**包含内容**：
- 布局系统历史问题
- 修复方案说明
- 测试验证结果

**适用读者**：维护者和高级开发者

---

### 📖 示例文档

#### [BUFFER_EXAMPLES.md](./BUFFER_EXAMPLES.md)
Buffer Layout 示例 - 图像查看器

**包含内容**：
- 使用 LayoutManager 构建复杂 UI
- 组件化设计（Header、Footer、Content）
- Renderable 接口实现
- 动态布局支持
- 多种渲染模式演示
- 键盘快捷键
- 架构说明

**运行示例**：
```bash
cd examples/image-buffer
go run main.go [image-path]
```

---

#### [BUFFER_TEST_EXAMPLES.md](./BUFFER_TEST_EXAMPLES.md)
Buffer 测试示例

**包含内容**：
- Buffer 核心功能测试
- 宽字符支持测试
- 性能测试示例

**运行示例**：
```bash
cd examples/image-buffer-test
go run main.go
```

---

## 快速入门

### 1. 理解基本概念

```go
// 创建 Buffer
buf := buffer.NewBuffer(80, 30)

// 设置单元格
buf.SetCell(buffer.Point{X:10, Y:5}, buffer.Cell{
    Char:  'A',
    Width: 1,
    Style: buffer.Style{Foreground: "red"},
})

// 渲染为字符串
output := buf.Render()
```

### 2. 使用 LayoutManager

```go
// 创建布局管理器
lm := buffer.NewLayoutManager(width, height)

// 计算布局
lm.CalculateLayout()

// 添加组件
lm.AddComponent("header", header)
lm.AddComponent("content", content)
lm.AddComponent("footer", footer)

// 渲染
output := lm.Render()
```

### 3. 创建自定义组件

```go
type MyComponent struct {
    content string
    style   buffer.Style
}

// 实现 Renderable 接口
func (c *MyComponent) Render(buf *buffer.Buffer, rect buffer.Rect) {
    buf.WriteString(
        buffer.Point{X: rect.X, Y: rect.Y},
        c.content,
        c.style,
    )
}

func (c *MyComponent) MinSize() (int, int) {
    return 10, 1
}

func (c *MyComponent) PreferredSize() (int, int) {
    return 80, 1
}
```

## 技术架构

### 核心组件

```
buffer/
├── buffer.go       # 核心 Buffer 实现
├── layout.go       # LayoutManager 布局管理
├── components.go   # Renderable 组件
├── cache.go        # 样式缓存
└── pool.go         # Buffer 池
```

### Layout 组件

```
layout/
├── area.go         # Area, Constraint, Fixed, Percent, Grow
├── flex.go         # Flex布局（RowLayout, ColumnLayout）
├── grid.go         # Grid布局
├── split.go        # 分割布局（SplitVertical, SplitHorizontal）
└── layout_test.go  # 布局测试
```

## 性能数据

| 操作 | 时间 | 分配 |
|------|------|------|
| FillRect | 102,538 ns/op | 0 B/op |
| WriteString | 794 ns/op | 0 B/op |
| WriteStringWrapped | 2,455 ns/op | 0 B/op |
| WriteBuffer | 1,702 ns/op | 0 B/op |
| Render | 16,970 ns/op | 904 B/op |
| LayoutCalculate | 300 ns/op | 96 B/op |
| LayoutRender | 150,900 ns/op | 235KB/op |

**结论**：性能足够 60fps TUI 应用（16.6ms/帧预算）

## 相关资源

### 代码位置

| 功能 | 文件位置 |
|------|---------|
| Buffer 核心 | `ui/render/buffer/buffer.go` |
| Layout 管理器 | `ui/render/buffer/layout.go` |
| 组件系统 | `ui/render/buffer/components.go` |
| Buffer 池/缓存 | `ui/render/buffer/cache.go`, `pool.go` |
| 布局约束 | `ui/layout/area.go` |
| Flex 布局 | `ui/layout/flex.go` |
| Grid 布局 | `ui/layout/grid.go` |
| 分割布局 | `ui/layout/split.go` |
| 垂直布局 | `ui/components/layout/vbox.go` |

### 示例代码

| 示例 | 位置 |
|------|------|
| 基础布局示例 | `examples/buffer-demo/` |
| 图像查看器（Buffer） | `examples/image-buffer/` |
| 图像测试 | `examples/image-buffer-test/` |

### 其他文档

- **整体架构**: `docs/ARCHITECTURE.md`
- **API 参考**: `docs/API.md`
- **v2.0 迁移**: `docs/MIGRATION_V2.md`
- **字符支持**: `CHARACTER_SUPPORT.md`
- **性能分析**: `PERFORMANCE_ANALYSIS.md`

## 最佳实践

### ✅ 推荐做法

1. **使用 LayoutManager 处理复杂布局**
   ```go
   lm := buffer.NewLayoutManager(width, height)
   lm.ImageLayout(displayHeight)
   ```

2. **组件独立渲染到子 Buffer**
   ```go
   compBuf := buffer.NewBuffer(rect.Width, rect.Height)
   component.Render(compBuf, Rect{0,0, rect.Width, rect.Height})
   mainBuf.WriteBuffer(Point{rect.X, rect.Y}, compBuf)
   ```

3. **对 Sixel 使用精确 displayHeight**
   ```go
   actualDisplayHeight := getSixelDisplayHeight(imageData)
   lm.ImageLayout(actualDisplayHeight)
   ```

4. **复用 Buffer 池**
   ```go
   buf := buffer.GetBuffer(width, height)
   defer buffer.PutBuffer(buf)
   ```

### ❌ 避免做法

1. **不要在每次渲染都创建新 Buffer**（性能问题）
2. **不要用字符串操作计算高度**（不准确）
3. **不要忽略宽字符处理**（显示错误）
4. **不要重复计算样式**（使用缓存）

## 常见问题

### Q: 为什么使用 Buffer 而不是直接字符串？

**A**: Buffer 提供精确的布局计算，特别是在处理 Sixel 图像和复杂布局时，避免了字符串操作的"高度计算地狱"。

### Q: Buffer 会影响性能吗？

**A**: 不会。基准测试显示，完整的布局+渲染仅需 ~0.15ms，远低于 60fps 的 16.6ms 预算。

### Q: 如何支持动画？

**A**: 使用 `tea.Tick` 命令定期更新组件状态，组件会自动重新渲染。

### Q: 如何处理窗口大小调整？

**A**: 在 Update 中处理 `WindowSizeMsg`，重新计算布局：
   ```go
   func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
       switch msg := msg.(type) {
       case tea.WindowSizeMsg:
           lm := buffer.NewLayoutManager(msg.Width, msg.Height)
           lm.ImageLayout(m.displayHeight)
           // ...
       }
   }
   ```

## 贡献

如果您想改进 Buffer Layout 系统：

1. 集成测试：添加到 `ui/render/buffer/buffer_test.go`
2. 性能测试：添加到 `benchmarks/`
3. 文档更新：更新本目录下的相关文档
4. 示例代码：添加到 `examples/buffer-*/`

## 版本历史

| 版本 | 日期 | 主要变更 |
|------|------|---------|
| v1.0 | 2025-01-15 | 初始版本 |
| v1.1 | 2025-01-20 | 添加 Buffer Layer |
| v1.2 | 2025-01-25 | 添加 Zoom Modes |
| v1.3 | 2025-02-01 | 添加 Sixel 渲染器 |
| v1.4 | 2025-02-03 | 修复工具栏定位问题 |

---

**相关联系列**:
- [v2.0 架构](../V2_ROADMAP.md)
- [API 文档](../API.md)
- [项目 README](../../README.md)

**维护者**: Taproot Team
**最后更新**: 2025-02-03
