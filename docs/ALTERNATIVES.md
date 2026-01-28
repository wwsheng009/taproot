# Taproot TUI - 替代方案与技术选型

## 概述

本文档分析在 TUI 框架开发中可能遇到的技术选择，以及现有的替代方案。

---

## 核心 TUI 框架对比

### Bubbletea (已选定 ✅)

**优点**:
- Elm 架构，易于理解和测试
- 活跃的社区支持
- Charmbracelet 生态系统完整
- 成熟的生产应用验证 (Crush, Glow, etc.)

**缺点**:
- 性能相比直接渲染略低
- 学习曲线对新手较陡

**结论**: 已作为 Taproot 基础框架，不建议更换

### 其他框架 (对比参考)

| 框架 | 语言 | 架构 | 活跃度 | 适用场景 |
|------|------|------|--------|----------|
| **termui** | Go | 直接渲染 | 低 | 简单仪表板 |
| **tview** | Go | 直接渲染 | 中 | 复杂布局 |
| **gocui** | Go | 直接渲染 | 低 | 直接键盘控制 |
| **rich** | Python | 声明式 | 高 | Python 应用 |
| **textual** | Python | 声明式 | 高 | Python 复杂应用 |

---

## 组件替代方案

### 1. 虚拟化列表

#### 选项 A: 迁移 Crush 实现 ✅ 推荐

**优点**:
- 已经过生产验证
- 功能完整（过滤、分组、搜索）
- 与 Taproot 设计一致

**缺点**:
- 迁移工作量较大 (12h)

#### 选项 B: 使用 bubbles/list

**优点**:
- 官方组件
- 维护良好

**缺点**:
- 功能相对简单
- 不支持虚拟化（大数据性能问题）

#### 选项 C: 自己实现

**优点**:
- 完全控制
- 针对性优化

**缺点**:
- 开发时间长
- 可能引入 bug

**结论**: 推荐迁移 Crush 实现

---

### 2. 对话框系统

#### 选项 A: 迁移 Crush 实现 ✅ 推荐

**优点**:
- 支持层级管理
- 使用 lipgloss.Layer
- 键盘导航完整

**缺点**:
- 需要解耦业务逻辑

#### 选项 B: 使用 bubbles/viewport

**优点**:
- 官方组件
- 简单场景够用

**缺点**:
- 无层级管理
- 无对话框栈

#### 选项 C: 使用 bubblezone

**优点**:
- 区域管理
- 事件委托

**缺点**:
- 学习成本
- 不专门针对对话框

**结论**: 推荐迁移 Crush 实现

---

### 3. 文本编辑器

#### 选项 A: 迁移 Crush 实现

**优点**:
- 功能完整（自动补全、剪贴板）
- 生产验证

**缺点**:
- 极其复杂 (20h+)
- 跨平台剪贴板支持困难

#### 选项 B: 使用 bubbles/textarea

**优点**:
- 官方组件
- 稳定可靠

**缺点**:
- 功能有限
- 无自动补全

#### 选项 C: 使用 micro/vim/nano

**优点**:
- 功能强大
- 用户熟悉

**缺点**:
- 用户体验不连贯
- 难以集成到 TUI

**结论**: Phase 5 可选，推荐使用 bubbles/textarea

---

### 4. Markdown 渲染

#### 选项 A: Glamour ✅ 推荐

**优点**:
- Charmbracelet 官方
- 样式可定制
- 支持代码高亮

**缺点**:
- 性能开销较大

#### 选项 B: goldmark + 自定义渲染

**优点**:
- 性能更好
- 灵活控制

**缺点**:
- 开发工作量大

#### 选项 C: 纯文本显示

**优点**:
- 简单
- 性能最好

**缺点**:
- 用户体验差

**结论**: 推荐使用 Glamour

---

### 5. 语法高亮

#### 选项 A: Chroma ✅ 推荐

**优点**:
- 支持 200+ 语言
- 颜色主题丰富
- 性能良好

**缺点**:
- 依赖包较大

#### 选项 B: pygments

**优点**:
- 功能强大

**缺点**:
- Python 依赖
- 集成困难

#### 选项 C: 无高亮

**优点**:
- 简单

**缺点**:
- 可读性差

**结论**: 推荐使用 Chroma

---

### 6. 图片渲染

#### 选项 A: 迁移 Crush 实现

**优点**:
- 支持多种协议 (kitty, iterm2)
- 性能优化

**缺点**:
- 终端兼容性问题
- 维护成本高

#### 选项 B: 使用 viu/termimage

**优点**:
- 独立工具
- 社区维护

**缺点**:
- 外部依赖
- 集成复杂

#### 选项 C: ASCII 艺术

**优点**:
- 通用兼容
- 无外部依赖

**缺点**:
- 质量较低

**结论**: Phase 5 可选，根据需求决定

---

### 7. 剪贴板支持

#### 选项 A: 迁移 Crush 实现

**优点**:
- 跨平台
- 功能完整

**缺点**:
- 平台特定代码
- 维护成本

#### 选项 B: clipboard

**优点**:
- 独立库
- 跨平台

**缺点**:
- 依赖 CGo
- 编译复杂

#### 选项 C: 无剪贴板支持

**优点**:
- 简单

**缺点**:
- 功能受限

**结论**: 如需要编辑器，使用 `clipboard` 库

---

## 样式与主题

### 选项 A: Lipgloss ✅ 已选定

**优点**:
- Charmbracelet 官方
- 功能强大
- 性能良好

**结论**: 已作为 Taproot 基础

### 选项 B: Termenv

**优点**:
- 支持更多终端特性
- 颜色查询

**缺点**:
- 与 lipgloss 功能重叠

**结论**: 可用于高级特性（如终端颜色查询）

---

## 性能优化方案

### 渲染优化

#### 策略 1: 缓存预渲染

```go
type Component struct {
    cachedView string
    dirty      bool
}

func (c *Component) View() string {
    if !c.dirty {
        return c.cachedView
    }
    c.cachedView = c.render()
    c.dirty = false
    return c.cachedView
}
```

#### 策略 2: 使用 strings.Builder

```go
func (c *Component) View() string {
    var b strings.Builder
    b.Grow(1024) // 预分配
    b.WriteString(header)
    b.WriteString(body)
    b.WriteString(footer)
    return b.String()
}
```

#### 策略 3: 虚拟化渲染

```go
// 只渲染可见区域
type VirtualList struct {
    items    []Item
    visible  int
    offset   int
}

func (v *VirtualList) View() string {
    start := v.offset
    end := min(start+v.visible, len(v.items))
    // 只渲染 start:end
}
```

### 事件处理优化

#### 策略 1: 节流

```go
type Throttle struct {
    lastTime time.Time
    interval time.Duration
}

func (t *Throttle) ShouldProcess() bool {
    now := time.Now()
    if now.Sub(t.lastTime) < t.interval {
        return false
    }
    t.lastTime = now
    return true
}
```

#### 策略 2: 批处理

```go
func batch(cmds ...tea.Cmd) tea.Cmd {
    return tea.Batch(cmds...)
}
```

---

## 测试替代方案

### 单元测试

#### 选项 A: 标准库 testing ✅ 推荐

**优点**:
- Go 原生
- 简单直接

**结论**: 已使用

#### 选项 B: testify

**优点**:
- 断言更丰富
- Mock 支持

**缺点**:
- 额外依赖

**结论**: 可选，复杂场景推荐

### Golden File 测试

```go
func TestComponentView(t *testing.T) {
    m := NewComponent()
    got := m.View()

    golden := "testdata/view.golden"
    if *update {
        os.WriteFile(golden, []byte(got), 0644)
    }

    want, _ := os.ReadFile(golden)
    if diff := cmp.Diff(string(want), got); diff != "" {
        t.Errorf("View() mismatch (-want +got):\n%s", diff)
    }
}
```

---

## 持续集成

### GitHub Actions 示例

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go: ['1.22', '1.23', '1.24']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - run: go test ./...
      - run: go test -race ./...
```

---

## 发布策略

### 选项 A: 单体发布 ✅ 推荐

**优点**:
- 使用简单
- 依赖清晰

**缺点**:
- 包体积大

**适用**: 大多数场景

### 选项 B: 模块化发布

```
github.com/yourorg/taproot
github.com/yourorg/taproot/dialogs
github.com/yourorg/taproot/list
github.com/yourorg/taproot/diffview
```

**优点**:
- 按需引入
- 包体积小

**缺点**:
- 版本管理复杂
- 依赖混乱

**适用**: 大型项目

### 选项 C: 插件系统

**优点**:
- 极度灵活
- 第三方扩展

**缺点**:
- 复杂度高
- 稳定性风险

**适用**: 特殊需求

**结论**: 推荐单体发布，必要时考虑模块化

---

## 文档方案

### 选项 A: Godoc ✅ 推荐

**优点**:
- Go 标准
- 自动生成

**结论**: 已使用

### 选项 B: pkg.go.dev

**优点**:
- 官方托管
- 版本浏览

**结论**: 配合 Godoc 使用

### 选项 C: 自建文档

**优点**:
- 可定制
- 丰富内容

**缺点**:
- 维护成本

**工具推荐**:
- Docusaurus
- Hugo
- MkDocs

---

## 总结

### 推荐技术栈

| 类别 | 选择 | 理由 |
|------|------|------|
| **TUI 框架** | Bubbletea | Elm 架构，生态完善 |
| **样式** | Lipgloss | 官方支持，功能强大 |
| **列表** | 迁移 Crush | 功能完整，生产验证 |
| **对话框** | 迁移 Crush | 层级管理，键盘导航 |
| **编辑器** | bubbles/textarea | 简单够用 |
| **Markdown** | Glamour | 官方支持 |
| **语法高亮** | Chroma | 语言支持多 |
| **图片** | 可选 | 兼容性问题 |
| **测试** | testing + testify | 标准库 + 增强 |
| **CI** | GitHub Actions | 免费易用 |

### 技术债务

1. **剪贴板支持**: 跨平台复杂，考虑使用独立库
2. **图片渲染**: 终端兼容性差，建议作为可选功能
3. **文本编辑器**: 过于复杂，建议使用简化版本

### 未来考虑

1. **WebAssembly**: 终端 TUI 在浏览器中运行
2. **远程渲染**: 分离渲染和逻辑
3. **插件系统**: 允许第三方扩展
