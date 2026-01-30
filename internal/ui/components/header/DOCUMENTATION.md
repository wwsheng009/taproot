# Header Component 完整技术文档

## 概述 (Overview)

Header Component 是 Taproot TUI 框架中的顶部标题栏组件，显示应用程序品牌、渐变标题、token使用进度条、错误计数、工作目录等信息。

**文件位置**: `internal/ui/components/header/header.go`
**组件类型**: 状态驱动渲染组件 (State-driven rendering)
**设计模式**: 纯函数式渲染 (Pure functional rendering)

---

## 目录 (Table of Contents)

1. [组件架构](#组件架构)
2. [数据结构](#数据结构)
3. [初始化流程详解](#初始化流程详解)
4. [渲染流程详解](#渲染流程详解)
5. [状态更新流程](#状态更新流程)
6. [样式系统](#样式系统)
7. [已知问题与根因分析](#已知问题与根因分析)
8. [调试指南](#调试指南)
9. [API 参考](#api-参考)
10. [性能优化建议](#性能优化建议)

---

## 组件架构

### 设计哲学

Header Component 遵循以下设计原则：

1. **纯数据驱动**: 所有渲染完全依赖于 `HeaderComponent` 结构体的状态字段
2. **不可变渲染**: `View()` 方法是纯函数，不修改组件状态
3. **无内部缓存**: 每次调用 `View()` 都从零开始构建整个渲染字符串
4. **宽度约束优先**: 所有渲染步骤都受到 `width` 字段的约束

### 组件职责

```
HeaderComponent
├── 维护窗口尺寸状态 (width, height)
├── 维护内容状态 (brand, title, cwd, tokens, errors)
├── 渲染品牌和渐变标题
├── 渲染token使用进度条
├── 渲染详情信息 (错误数、百分比、提示)
├── 适应窗口宽度 (自动截断、填充)
└── 确保单行输出 (防换行)
```

### 依赖图

```
HeaderComponent (header.go)
│
├── github.com/charmbracelet/lipgloss
│   ├── NewStyle() - 创建样式对象
│   ├── Foreground() - 设置前景色
│   ├── MaxWidth() - 限制最大宽度
│   ├── MaxHeight() - 限制最大高度
│   ├── Padding() - 设置内边距
│   ├── Width(str) - 计算可视宽度
│   └── Render(str) - 应用样式并渲染
│
├── internal/ui/styles.DefaultStyles()
│   └── 返回预定义的颜色和样式集合
│
└── internal/ui/styles.grad
    ├── ApplyBoldForegroundGrad() - 渐变文本渲染
    └── 依赖:
        ├── github.com/rivo/uniseg (Grapheme聚类)
        └── github.com/lucasb-eyer/go-colorful (颜色混合)
```

---

## 数据结构

### HeaderComponent 结构体

```go
// 文件: internal/ui/components/header/header.go:28-41
type HeaderComponent struct {
    // === 布局状态 ===
    width        int    // 窗口宽度（如 100, 150）
    height       int    // Header高度（通常固定为 1）

    // === 品牌和标题 ===
    brand        string // 品牌名称（如 "Charm™"）
    title        string // 应用标题（如 "CRUSH"）

    // === 详情信息 ===
    sessionTitle string // 会话标题（可选）
    workingDir   string // 当前工作目录路径

    // === Token 使用统计 ===
    tokenUsed    int     // 已使用的token数量
    tokenMax     int     // 最大token数量限制
    cost         float64 // Token成本（用于显示）

    // === 状态指示 ===
    errorCount   int     // 错误计数（0或负数时不显示错误图标）
    detailsOpen  bool    // 详情面板是否打开（显示 "open" 或 "close"）

    // === 显示模式 ===
    compactMode  bool    // 紧凑模式（当前已定义但未实现）
}
```

### 内部常量

```go
// 文件: internal/ui/components/header/header.go:115-121
const (
    gap          = " "       // 元素间间距（一个空格）
    diag         = "╱"       // 进度条斜杠字符
    minDiags     = 3         // 进度条最小宽度（字符数）
    leftPadding  = 1         // 左内边距（空格数）
    rightPadding = 1         // 右内边距（空格数）
)
```

### 接口实现

```go
// 文件: internal/ui/components/header/header.go:15-25
type headerImpl interface {
    Size() (width, height int)              // layout.Sizeable
    SetSize(width, height int)             // layout.Sizeable
    SetBrand(brand, title string)
    SetSessionTitle(title string)
    SetWorkingDirectory(cwd string)
    SetTokenUsage(used, max int, cost float64)
    SetErrorCount(count int)
    SetDetailsOpen(open bool)
    ShowingDetails() bool
}

// HeaderComponent 同时实现了以下接口：
// - headerImpl (私有接口)
// - layout.Sizeable (公有接口: Size() 和 SetSize())
```

---

## 初始化流程详解

### 阶段1: 创建组件实例

**入口**: `New()` 构造函数

```go
// 文件: internal/ui/components/header/header.go:44-51
func New() *HeaderComponent {
    return &HeaderComponent{
        brand:       "Charm™",      // 默认品牌
        title:       "CRUSH",       // 默认标题
        tokenMax:    128000,        // 默认token上限
        compactMode: false,         // 默认非紧凑模式
        // 注意: width, height 默认为 0，需要后续设置
    }
}
```

**初始化状态表**:

| 字段 | 初始值 | 说明 |
|------|--------|------|
| width | 0 | 必须在后续设置 |
| height | 0 | 必须在后续设置 |
| brand | "Charm™" | 可覆盖 |
| title | "CRUSH" | 可覆盖 |
| tokenMax | 128000 | 可覆盖 |
| errorCount | 0 | 默认不显示错误 |
| detailsOpen | false | 默认显示 "open" |
| workingDir | "" | 默认显示 "~" |
| compactMode | false | 已定义但未使用 |

### 阶段2: 设置窗口尺寸

**入口**: `SetSize(width, height int)`

```go
// 文件: internal/ui/components/header/header.go:59-62
func (h *HeaderComponent) SetSize(width, height int) {
    h.width = width   // 窗口总宽度（如 100, 200）
    h.height = height // Header高度（通常为 1）
}
```

**调用时机**:
1. 初始化时: `h.SetSize(100, 1)`
2. 窗口resize时: 在 Bubble Tea 的 `tea.WindowSizeMsg` 处理中调用

**Demo中的实现** (examples/header-demo/main.go:108-112):

```go
case tea.WindowSizeMsg:
    // 更新header尺寸（header为1行高）
    m.header.SetSize(msg.Width, 1)
    m.contentHeight = msg.Height - 1  // 内容区域高度
```

### 阶段3: 配置内容属性

**Setter方法列表**:

```go
// 文件: internal/ui/components/header/header.go

// 设置品牌和标题
func (h *HeaderComponent) SetBrand(brand, title string) {
    h.brand = brand  // 品牌名称
    h.title = title  // 应用标题
}

// 设置会话标题（当前View()中未使用）
func (h *HeaderComponent) SetSessionTitle(title string) {
    h.sessionTitle = title
}

// 设置工作目录
func (h *HeaderComponent) SetWorkingDirectory(cwd string) {
    h.workingDir = cwd
}

// 设置Token使用情况
func (h *HeaderComponent) SetTokenUsage(used, max int, cost float64) {
    h.tokenUsed = used  // 已使用数量
    h.tokenMax = max    // 最大数量
    h.cost = cost       // 成本
}

// 设置错误计数
func (h *HeaderComponent) SetErrorCount(count int) {
    h.errorCount = count
}

// 设置详情面板状态
func (h *HeaderComponent) SetDetailsOpen(open bool) {
    h.detailsOpen = open
}

// 设置紧凑模式（当前View()中未使用）
func (h *HeaderComponent) SetCompactMode(compact bool) {
    h.compactMode = compact
}
```

### 阶段4: 完整初始化示例

**Demo中的完整初始化流程** (examples/header-demo/main.go:25-48):

```go
func initialModel() *model {
    brand := "Charm™"
    title := "CRUSH"

    // Step 1: 创建header实例（New()）
    h := header.New()
    //   -> h.brand = "Charm™"
    //   -> h.title = "CRUSH"
    //   -> h.width = 0
    //   -> h.height = 0

    // Step 2: 设置尺寸
    h.SetSize(100, 1)
    //   -> h.width = 100
    //   -> h.height = 1

    // Step 3: 设置品牌（覆盖默认值）
    h.SetBrand(brand, title)
    //   -> h.brand = "Charm™"
    //   -> h.title = "CRUSH"

    // Step 4: 设置工作目录
    h.SetWorkingDirectory("/projects/ai/Taproot")
    //   -> h.workingDir = "/projects/ai/Taproot"

    // Step 5: 设置Token使用情况
    h.SetTokenUsage(0, 128000, 0.00)
    //   -> h.tokenUsed = 0
    //   -> h.tokenMax = 128000
    //   -> h.cost = 0.00

    // Step 6: 设置错误计数
    h.SetErrorCount(3)
    //   -> h.errorCount = 3

    return &model{
        header:       h,
        errorCount:   3,
        workingDir:   "/projects/ai/Taproot",
        tokenUsed:    0,
        tokenMax:     128000,
        cost:         0.00,
        detailsOpen:  false,
        compactMode:  false,
        brand:        brand,
        title:        title,
    }
}
```

### 初始化流程图

```
应用程序启动
    │
    ├─> header.New()
    │   └─> 返回带默认值的HeaderComponent实例
    │       │
    │       ├─> brand: "Charm™"
    │       ├─> title: "CRUSH"
    │       ├─> tokenMax: 128000
    │       └─> width/height: 0 (待设置)
    │
    ├─> header.SetSize(width, 1)
    │   └─> 更新 h.width 和 h.height
    │
    ├─> header.SetBrand(brand, title)
    │   └─> 更新 h.brand 和 h.title
    │
    ├─> header.SetWorkingDirectory(cwd)
    │   └─> 更新 h.workingDir
    │
    ├─> header.SetTokenUsage(used, max, cost)
    │   └─> 更新 h.tokenUsed, h.tokenMax, h.cost
    │
    └─> header.SetErrorCount(count)
        └─> 更新 h.errorCount
```

---

## 渲染流程详解

### View() 方法总览

**入口**: `h.View() string` (header.go:108-254)

**返回值**: 包含完整ANSI颜色代码的渲染字符串

**核心原则**:
1. 每次都从零开始构建（无缓存）
2. 所有宽度计算基于 `h.width`
3. 多层防护确保单行输出
4. 手动处理ANSI序列截断

### 渲染流程分阶段详解

#### 阶段1: 前置检查和样式准备 (Lines 108-113)

```go
// Line 109-111: 空品牌检查
if h.brand == "" {
    return ""  // 空字符串，不渲染任何内容
}

// Line 113: 获取默认样式
s := styles.DefaultStyles()
```

**styles.DefaultStyles() 解析** (internal/ui/styles/styles.go:456-1050):

```go
func DefaultStyles() Styles {
    // 定义颜色常量
    var (
        primary   = Charple    // 主色调（紫色系）
        secondary = Dolly      // 次要色调（蓝色系）
        tertiary  = Bok        // 第三色调（绿色系）
        fgBase    = Ash        // 基础前景色
        fgMuted   = Squid      // 静音前景色
        fgSubtle  = Oyster     // 微妙前景色
        errorColor= Sriracha   // 错误颜色（红色系）
        // ... 更多颜色定义
    )

    // 创建基础样式
    base := lipgloss.NewStyle().Foreground(fgBase)

    // 构建完整的Styles结构
    return Styles{
        Base:  lipgloss.NewStyle().Foreground(fgBase),
        Muted: lipgloss.NewStyle().Foreground(fgMuted),
        Subtle: lipgloss.NewStyle().Foreground(fgSubtle),

        // 颜色引用
        Primary:   primary,
        Secondary: secondary,
        Tertiary:  tertiary,

        Error: errorColor,
        // ... 更多样式
    }
}
```

**关键点**:
- `s.Base`: 基础样式对象，用于继承和组合样式
- `s.Primary`, `s.Secondary`: 用于渐变文本的颜色端点
- `s.Error`, `s.Muted`, `s.Subtle`: 用于不同元素的样式

#### 阶段2: 渲染品牌和渐变标题 (Lines 115-129)

**常量定义**:
```go
const (
    gap          = " "       // 间距
    diag         = "╱"       // 斜杠
    minDiags     = 3         // 最小斜杠数
    leftPadding  = 1         // 左边距
    rightPadding = 1         // 右边距
)
```

**品牌渲染** (Line 126):
```go
b.WriteString(s.Base.Foreground(s.Secondary).Render(h.brand))
// 示例: h.brand = "Charm™"
// 输出: "\x1b[38;5;245mCharm™\x1b[0m"
```

**标题渐变渲染** (Line 128):
```go
b.WriteString(styles.ApplyBoldForegroundGrad(&s, h.title, s.Secondary, s.Primary))
// 示例: h.title = "CRUSH"
// 输出: 每个字符独立着色，如 "\x1b[38;5;60m\x1b[1mC\x1b[0m\x1b[38;5;58;1mR\x1b[0m..."
```

**ApplyBoldForegroundGrad 完整实现** (internal/ui/styles/grad.go:70-80):

```go
func ApplyBoldForegroundGrad(t *Styles, input string, color1, color2 color.Color) string {
    if input == "" {
        return ""
    }

    var o strings.Builder

    // 调用 ForegroundGrad 获取每个字符的独立着色片段
    clusters := ForegroundGrad(t, input, true, color1, color2)

    // 拼接所有片段
    for _, c := range clusters {
        fmt.Fprint(&o, c)
    }

    return o.String()
}
```

**ForegroundGrad 实现详解** (internal/ui/styles/grad.go:17-43):

```go
func ForegroundGrad(t *Styles, input string, bold bool, color1, color2 color.Color) []string {
    // 空字符串处理
    if input == "" {
        return []string{""}
    }

    // 单字符优化
    if len(input) == 1 {
        style := t.Base.Foreground(colorToLipgloss(color1))
        if bold {
            style.Bold(true)
        }
        return []string{style.Render(input)}
    }

    // 步骤1: 使用uniseg进行grapheme聚类（支持emoji等复合字符）
    var clusters []string
    gr := uniseg.NewGraphemes(input)
    for gr.Next() {
        clusters = append(clusters, string(gr.Runes()))
    }
    // 示例: "👋World" -> ["👋", "W", "o", "r", "l", "d"]

    // 步骤2: 生成颜色渐变色板
    ramp := blendColors(len(clusters), color1, color2)
    // 示例: 6个字符，2个颜色端点 -> 6个中间渐变色
    //       [color1*1.0, color1*0.8, color1*0.6, ..., color2*0.6, color2*0.8, color2*1.0]

    // 步骤3: 为每个字符应用对应的渐变色
    for i, c := range ramp {
        style := t.Base.Foreground(colorToLipgloss(c))
        if bold {
            style.Bold(true)  // 粗体
        }
        clusters[i] = style.Render(clusters[i])
    }

    return clusters
}
```

**blendColors 实现** (internal/ui/styles/grad.go:84-127):

```go
func blendColors(size int, stops ...color.Color) []color.Color {
    // 参数校验
    if len(stops) < 2 {
        return nil
    }

    // 转换为colorful.Color（使用HCL色彩空间以保持在色域内）
    stopsPrime := make([]colorful.Color, len(stops))
    for i, k := range stops {
        stopsPrime[i], _ = colorful.MakeColor(k)
    }

    // 计算分段
    numSegments := len(stopsPrime) - 1  // n个颜色端点 = n-1段
    baseSize := size / numSegments
    remainder := size % numSegments  // 余数分配到前面的段

    segmentSizes := make([]int, numSegments)
    for i := range numSegments {
        segmentSizes[i] = baseSize
        if i < remainder {
            segmentSizes[i]++
        }
    }

    // 为每段生成渐变色
    blended := make([]color.Color, 0, size)
    for i := range numSegments {
        c1 := stopsPrime[i]      // 段起始颜色
        c2 := stopsPrime[i+1]    // 段结束颜色
        segmentSize := segmentSizes[i]

        for j := range segmentSize {
            var t float64
            if segmentSize > 1 {
                t = float64(j) / float64(segmentSize-1)  // 插值参数
            }
            // 使用HCL色彩空间混合
            c := c1.BlendHcl(c2, t)
            blended = append(blended, c)
        }
    }

    return blended
}
```

**渐变文本ANSI结构示例**:

```
输入: "CRUSH" (5个字符)

步骤1: Grapheme聚类
  -> clusters = ["C", "R", "U", "S", "H"]

步骤2: 颜色渐变 (Secondary -> Primary, 5个中间色)
  -> ramp = [color2*1.0, color2*0.75, color2*0.50, color2*0.25, color1*1.0]

步骤3: 每个字符独立着色
  -> clusters = [
       "\x1b[38;5;60m\x1b[1mC\x1b[0m",
       "\x1b[38;5;58m\x1b[1mR\x1b[0m",
       "\x1b[38;5;61m\x1b[1mU\x1b[0m",
       "\x1b[38;5;63m\x1b[1mS\x1b[0m",
       "\x1b[38;5;68m\x1b[1mH\x1b[0m"
     ]

步骤4: 拼接
  -> 输出: "\x1b[38;5;60m\x1b[1mC\x1b[0m\x1b[38;5;58m\x1b[1mR\x1b[0m..."

字节长度:
  - 每个字符约 20-30 字节（包含ANSI序列）
  - "CRUSH" 渲染后约 120字节
  - lipgloss.Width() 仍返回 5（可视字符数）
```

#### 阶段3: 计算宽度分配 (Lines 131-133)

```go
// Line 132: 可用宽度 = 窗口宽度 - 左右边距
availableWidth := h.width - leftPadding - rightPadding
// 示例: h.width = 100, leftPadding = 1, rightPadding = 1
//       -> availableWidth = 98

// Line 133: 进度条宽度 = 可用宽度的25%
progressBarWidth := int(float64(availableWidth) * 0.25)
// 示例: availableWidth = 98
//       -> progressBarWidth = int(98 * 0.25) = int(24.5) = 24
```

**宽度分配图示** (以 width=100 为例):

```
|<-- leftPadding:1 -->|<-- 品牌+标题 -->|<-- gap:1 -->|<-- 进度条:25% -->|<-- gap:1 -->|<-- 详情 -->|<-- rightPadding:1 -->|
|          <--- availableWidth:98 --->                         |                     |
|<------------------------------------ h.width:100 --------------------------------->|
```

#### 阶段4: 渲染进度条 (Lines 135-156)

```go
if progressBarWidth > minDiags {  // 确保有足够空间
    // 步骤1: 计算token使用百分比
    var percentage float64
    if h.tokenUsed > 0 && h.tokenMax > 0 {
        percentage = float64(h.tokenUsed) / float64(h.tokenMax)
    }
    // 示例: h.tokenUsed = 50000, h.tokenMax = 128000
    //       -> percentage = 50000 / 128000 = 0.39 (39%)

    // 步骤2: 计算显示的斜杠数量
    // 即使 0% 也要显示 minDiags 个斜杠
    diagsCount := minDiags + int(float64(progressBarWidth-minDiags)*percentage)
    // 示例: progressBarWidth = 24, minDiags = 3, percentage = 0.39
    //       -> diagsCount = 3 + int((24-3) * 0.39) = 3 + int(8.19) = 11

    // 步骤3: 渲染进度条并填充到固定宽度
    diagsStr := strings.Repeat(diag, diagsCount)
    paddingCount := progressBarWidth - diagsCount
    if paddingCount > 0 {
        diagsStr += strings.Repeat(" ", paddingCount)
    }
    // 示例: diagsCount = 11, progressBarWidth = 24
    //       -> diagsStr = "╱╱╱╱╱╱╱╱╱╱╱" + " " * 13 = "╱╱╱╱╱╱╱╱╱╱╱             "

    // 步骤4: 应用主色调
    b.WriteString(s.Base.Foreground(s.Primary).Render(diagsStr))
    // 输出: "\x1b[38;5;68m╱╱╱╱╱╱╱╱╱╱╱             \x1b[0m"

    // 步骤5: 添加间距
    b.WriteString(gap)
}
```

**进度条可视化示例**:

```
Token使用: 0%   -> ╱╱╱                      (仅minDiags)
Token使用: 39%  -> ╱╱╱╱╱╱╱╱╱╱╱             (填充到25%宽度)
Token使用: 100% -> ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱�╱╱╱╱╱╱╱╱  (完全填充)
```

#### 阶段5: 渲染详情部分 (Lines 158-171)

```go
// 步骤1: 计算已使用宽度
usedWidth := lipgloss.Width(b.String())
// lipgloss.Width() 忽略ANSI序列，只计算可视字符数
// 示例: b.String() = "Charm™" + " " + "CRUSH" + " " + "╱╱╱..." + " "
//       -> usedWidth = 可视字符数 = ~50

// 步骤2: 计算详情可用宽度
detailsAvailWidth := availableWidth - usedWidth
// 示例: availableWidth = 98, usedWidth = 50
//       -> detailsAvailWidth = 48

// 步骤3: 渲染详情
if detailsAvailWidth > minDiags {
    details := h.renderDetails(detailsAvailWidth)  // 见下文详解
    detailsWidth := lipgloss.Width(details)         // 计算可视宽度
    if detailsWidth < detailsAvailWidth {
        // 步骤4: 填充空格以填满剩余宽度
        details += strings.Repeat(" ", detailsAvailWidth-detailsWidth)
    }
    b.WriteString(details)
}
```

**renderDetails() 完整实现** (header.go:256-310):

```go
func (h *HeaderComponent) renderDetails(availWidth int) string {
    s := styles.DefaultStyles()

    var parts []string

    // === 部分1: 错误计数 ===
    if h.errorCount > 0 {
        errorStyle := s.Base.Foreground(s.Error)
        parts = append(parts, errorStyle.Render(fmt.Sprintf("%s%d", styles.ErrorIcon, h.errorCount)))
        // 示例: h.errorCount = 3
        //       -> parts = ["\x1b[38;5;203m×3\x1b[0m"]
    }

    // === 部分2: Token百分比 ===
    var tokenStr string
    if h.tokenMax > 0 {
        percentage := int(float64(h.tokenUsed) / float64(h.tokenMax) * 100)
        tokenStr = fmt.Sprintf("%d%%", percentage)
    } else {
        tokenStr = fmt.Sprintf("%d", h.tokenUsed)
    }
    parts = append(parts, s.Muted.Render(tokenStr))
    // 示例: h.tokenUsed = 50000, h.tokenMax = 128000
    //       -> percentage = 39, tokenStr = "39%"
    //       -> parts = ["...", "\x1b[38;5;245m39%\x1b[0m"]

    // === 部分3: 详情提示 ===
    const keystroke = "ctrl+d"
    if h.detailsOpen {
        parts = append(parts, s.Muted.Render(keystroke)+s.Subtle.Render(" close"))
    } else {
        parts = append(parts, s.Muted.Render(keystroke)+s.Subtle.Render(" open "))
    }
    // 示例: h.detailsOpen = false
    //       -> parts = ["...", "...", "\x1b[38;5;245mctrl+d\x1b[0m\x1b[38;5;251m open \x1b[0m"]

    // === 部分4: 用分隔符连接 ===
    dot := s.Subtle.Render(" • ")
    metadata := strings.Join(parts, dot)
    metadata = dot + metadata  // 在前面也加一个dot
    // 示例: metadata = " • ×3 • 39% • ctrl+d open "

    // === 部分5: 工作目录处理 ===
    cwd := h.workingDir
    if cwd == "" {
        cwd = "~"  // 空路径显示为家目录符号
    }

    // 截断到最多4个组件
    dirs := strings.Split(cwd, string('/'))
    if len(dirs) > 4 {
        cwd = strings.Join(dirs[len(dirs)-4:], "/")
        cwd = "…" + cwd  // 添加省略号前缀
    }
    // 示例: cwd = "/projects/ai/Taproot"
    //       -> dirs = ["", "projects", "ai", "Taproot"]
    //       -> len(dirs) = 4 <= 4, 不截断

    // 示例: cwd = "/Users/john/projects/ai/Taproot"
    //       -> dirs = ["", "Users", "john", "projects", "ai", "Taproot"]
    //       -> len(dirs) = 6 > 4
    //       -> cwd = "…/projects/ai/Taproot"

    // 截断CWD适应剩余空间
    maxCwdWidth := max(0, availWidth-lipgloss.Width(metadata))
    cwd = lipgloss.NewStyle().
        MaxWidth(maxCwdWidth).
        MaxHeight(1).  // 确保单行
        Render(cwd)
    cwd = s.Muted.Render(cwd)
    // 示例: availWidth = 48, width(metadata) = ~20
    //       -> maxCwdWidth = 28
    //       -> cwd = "/projects/ai/Taproot"
    //       -> 如果长度 > 28，则截断为 "/projects/ai/Taproot" 或 "/.../ai/Taproot"

    // === 返回 ===
    return cwd + metadata
}
```

**详情部分结构**:

```
[工作目录] • [错误数] • [百分比] • [快捷键] [提示]
   ↑          ↑          ↑           ↑         ↑
 truncated   ×%        39%        ctrl+d    open
```

#### 阶段6: 手动ANSI截断处理 (Lines 173-228)

**问题背景**:
- 渐变文本产生密集的ANSI序列
- `lipgloss.MaxWidth()` 处理复杂ANSI序列时可能不可靠
- 需要手动截断并保留ANSI序列

**实现代码**:

```go
// === 步骤1: 获取已构建内容 ===
content := b.String()
contentWidth := lipgloss.Width(content)

// === 步骤2: 如果超出宽度则截断 ===
if contentWidth > availableWidth {
    // 手动逐字符截断，保留ANSI序列
    var truncated strings.Builder
    currentWidth := 0
    runes := []rune(content)
    i := 0

    for i < len(runes) && currentWidth < availableWidth {
        r := runes[i]

        // === ANSI转义序列检测 ===
        if r == '\x1b' {  // ESC字符标识ANSI序列开始
            // 找到序列结束（通常以'm'结尾的CSI序列）
            end := i + 1
            for end < len(runes) && runes[end] != 'm' {
                end++
            }

            // 如果找到完整的序列，全部包含
            if end < len(runes) {
                for j := i; j <= end; j++ {
                    truncated.WriteRune(runes[j])
                }
                i = end + 1  // 跳过整个ANSI序列
                continue
            }
        }

        // === 计算字符宽度 ===
        runeWidth := lipgloss.Width(string(r))
        // 检查是否可接受
        if currentWidth+runeWidth > availableWidth {
            break  // 超出限制，停止
        }

        truncated.WriteRune(r)
        currentWidth += runeWidth
        i++
    }

    // === 步骤3: 继续检查宽度 ===
    currentResult := truncated.String()
    if lipgloss.Width(currentResult) > availableWidth {
        // 降级方案: 使用lipgloss.MaxWidth作为最后的保障
        truncatedStyle := lipgloss.NewStyle().
            MaxWidth(availableWidth).
            MaxHeight(1)
        currentResult = truncatedStyle.Render(content)
    }
    content = currentResult
    contentWidth = lipgloss.Width(content)
}

// === 步骤4: 填充空格到完整宽度 ===
if contentWidth < availableWidth {
    content += strings.Repeat(" ", availableWidth-contentWidth)
}
```

**ANSI序列结构**:

```
标准CSI序列: \x1b[参数m

示例:
  \x1b[38;5;245m  -> 置前景色为256色板中第245色
  \x1b[1m          -> 加粗
  \x1b[0m          -> 重置所有属性

嵌套示例（渐变文本）:
  \x1b[38;5;60m\x1b[1mC\x1b[0m
  ├─> 颜色60（灰紫色）
  ├─> 加粗
  ├─> 字符 'C'
  └─> 重置

问题: 如果在颜色序列中间截断，会导致剩余文本呈现错误的颜色
```

**手动漫游图示**:

```
内容: "\x1b[38;5;60m\x1b[1mC\x1b[38;5;58mR\x1b[38;5;61mU\x1b[0m..."

逐字符处理:
  i=0: r='\x1b' -> 检测到ANSI开始
       -> 查找结束，找到 '\x1b[38;5;60m'
       -> 完整写入，不增加currentWidth
       -> i跳到序列结束+1的位置

  i=N: r='C' -> 是普通字符
       -> runeWidth = 1
       -> currentWidth = 1
       -> 写入'C'

  i=N+1: r='\x1b' -> 又是ANSI序列
       -> ...重复此过程
```

#### 阶段7: 应用Padding并最终安全检查 (Lines 235-253)

```go
// === 步骤1: 应用左右padding ===
result := s.Base.Padding(0, rightPadding, 0, leftPadding).Render(content)
// 参数顺序: top, right, bottom, left
// 示例: leftPadding=1, rightPadding=1
//       -> 在content前面加1个空格，后面加1个空格
//       -> 最终宽度 = contentWidth + leftPadding + rightPadding

// === 步骤2: 最终安全检查 ===
if strings.ContainsAny(result, "\n\r") {
    // 理论上不应出现，但作为安全网
    // 移除所有换行符和回车符
    result = strings.ReplaceAll(result, "\n", "")
    result = strings.ReplaceAll(result, "\r", "")

    // 重新填充到正确宽度
    resultWidth := lipgloss.Width(result)
    targetWidth := h.width
    if resultWidth < targetWidth {
        result += strings.Repeat(" ", targetWidth-resultWidth)
    }
}

return result
```

### 渲染输出完整示例

假设:
```
h.width = 100
h.height = 1
h.brand = "Charm™"
h.title = "CRUSH"
h.workingDir = "/projects/ai/Taproot"
h.tokenUsed = 50000
h.tokenMax = 128000
h.errorCount = 3
h.detailsOpen = false
```

**可视化输出**:
```
 Charm™ CRUSH ╱╱╱╱╱╱╱╱╱╱╱              /projects/ai/Taproot • ×3 • 39% • ctrl+d open
↑↑           ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
││           │                                            │                                     │
││           │       进度条 (25%, 约占24列)                 │                                     │
││           └─── 渐变标题 ( Secondary -> Primary)          │                                     │
│└─── 品牌            └─── 详情区域 (cwd, errors, %, hint)                                       │
│                                                                          └── 前提 (ctrl+d open)
└── 左padding                                                                   右padding
```

**宽度分配表**:

| 区域 | 实际宽度 | 占比 |
|------|----------|------|
| 左padding | 1 | 1% |
| 品牌名称 | 6 | 6% |
| 间距 | 1 | 1% |
| 标题 | 5 | 5% |
| 间距 | 1 | 1% |
| 进度条 | 24 | 24% |
| 间距 | 1 | 1% |
| 工作目录 | ~25 | ~25% |
| 分隔符+错误 | ~5 | ~5% |
| 百分比 | ~3 | ~3% |
| 分隔符+提示 | ~22 | ~22% |
| 右padding | 1 | 1% |
| **总计** | **95** | **95%** (剩余~5%为空格填充) |

### View() 流程图

```
View() 调用
    │
    ├─> 前置检查: h.brand == "" ?
    │   └─> 是: 返回 ""
    │   └─> 否: 继续
    │
    ├─> 获取样式: s := styles.DefaultStyles()
    │
    ├─> 构建内容 (使用 strings.Builder)
    │   │
    │   ├─> 步骤1: 渲染品牌和标题
    │   │   ├─> s.Base.Foreground(s.Secondary).Render(h.brand)
    │   │   └─> ApplyBoldForegroundGrad(h.title)
    │   │       └─> ForegroundGrad(UNISEG聚类 + 颜色渐变)
    │   │
    │   ├─> 步骤2: 渲染进度条
    │   │   ├─> 计算percentage = tokenUsed / tokenMax
    │   │   ├─> 计算diags = minDiags + (progressBarWidth-minDiags)*percentage
    │   │   └─> 填充到固定宽度
    │   │
    │   ├─> 步骤3: 渲染详情
    │   │   └─> 调用 renderDetails()
    │   │       ├─> 构建parts列表 (errors, %, ctrl+d)
    │   │       ├─> 用 " • " 连接
    │   │       ├─> 截断工作目录到4个组件
    │   │       └─> 应用MaxWidth和MaxHeight(1)
    │   │
    │   ├─> 步骤4: 手动截断确保宽度
    │   │   ├─> 逐字符遍历，检测ANSI序列
    │   │   ├─> 保留完整ANSI序列
    │   │   ├─> 追踪可视宽度
    │   │   └─> 降级到lipgloss.MaxWidth
    │   │
    │   └─> 步骤5: 填充空格
    │       └─> content += spaces
    │
    ├─> 应用Padding
    │   └─> s.Base.Padding(0, 1, 0, 1).Render(content)
    │
    ├─> 最终安全检查
    │   ├─> 检测换行符
    │   └─> 移除 \n \r
    │
    └─> 返回 result
```

---

## 状态更新流程

### 设计理念

Header Component 采用**无状态更新模式**:

- `Update()` 方法是占位符，不执行任何逻辑
- 所有状态通过 Setter 方法直接修改
- 下次 `View()` 调用时自动反映最新状态

**好处**:
1. 简化代码: 无需复杂的状态管理
2. 可预测: `View()` 输出完全由当前状态决定
3. 易测试: 可直接修改字段验证渲染

### Setter方法工作流程

对于每个 Setter 方法:

```go
// 例如: SetWorkingDirectory
func (h *HeaderComponent) SetWorkingDirectory(cwd string) {
    h.workingDir = cwd  // 直接修改字段
    // 无其他副作用
    // 无缓存清除
    // 无事件触发
}
```

**调用序列** (在 demo 中):

```go
// 1. 用户按 'h' 键
case tea.KeyMsg:
    switch msg.String() {
    case "h":
        // 2. 修改本地状态
        m.workingDir = "/new/path"

        // 3. 更新 header
        m.header.SetWorkingDirectory(m.workingDir)
        //   -> h.workingDir = "/new/path"

        // 4. 下次 Bubble Tea 渲染时
        //    -> 调用 m.header.View()
        //    -> 重新构建所有内容
        //    -> 使用新的 h.workingDir
    }
```

### WindowResize 处理

**完整流程** (examples/header-demo/main.go:108-112):

```go
case tea.WindowSizeMsg:
    // 1. 更新header尺寸
    m.header.SetSize(msg.Width, 1)

    // 2. 更新内容区域高度
    m.contentHeight = msg.Height - 1

    // 3. 下次渲染时
    //    -> m.header.View()
    //    -> 使用新的 h.width 重新计算所有宽度
    //    -> 可能调整详情截断
    //    -> 可能调整进度条宽度（25%计算）
```

**Resize时的内部变化**:

```
旧: h.width = 100
    availableWidth = 98
    progressBarWidth = 24
    detailsAvailWidth = ~45

    ↓ 窗口缩小到 80 ↓

新: h.width = 80
    availableWidth = 78
    progressBarWidth = int(78 * 0.25) = 19
    detailsAvailWidth = ~30

结果:
  - 工作目录被截断更多
  - 进度条变短（斜杠数减少）
  - 可能触发手动截断逻辑
```

### Update() 占位符

```go
// 文件: internal/ui/components/header/header.go:313-316
func (h *HeaderComponent) Update(msg any) (*HeaderComponent, any) {
    // Placeholder - engine-agnostic Update method
    return h, nil
}
```

**为什么Update()是占位符**:

1. **Engine-Agnostic设计**: HeaderComponent设计为可在不同渲染引擎中使用（BubbleTea, Ultraviolet等）
2. **当前实现**: 在BubbleTea中，消息由外层的model处理，通过Setter方法更新header
3. **未来扩展**: 如果需要header内部处理键盘消息，可以在此实现

---

## 样式系统

### DefaultStyles() 详解

**入口**: `styles.DefaultStyles()` (internal/ui/styles/styles.go:456-1050)

**返回值**: 包含所有预定义样式和颜色的 `Styles` 结构

### Styles 结构体 (简化)

```go
type Styles struct {
    // === 颜色 ===
    Primary   lipgloss.Color  // 主色调 (紫色系)
    Secondary lipgloss.Color  // 次要色调 (蓝色系)
    Tertiary  lipgloss.Color  // 第三色调 (绿色系)
    Error     lipgloss.Color  // 错误 (红色系)
    Warning   lipgloss.Color  // 警告 (黄色系)
    Info      lipgloss.Color  // 信息 (蓝色系)

    // === 文本样式 ===
    Base   lipgloss.Style  // 基础样式 (前景色)
    Muted  lipgloss.Style  // 静音样式 (淡文本)
    Subtle lipgloss.Style  // 微妙样式 (更淡)

    // === 预设样式带颜色 ===
    PrimaryStyle   lipgloss.Style
    SecondaryStyle lipgloss.Style

    // ... 更多样式
}
```

### Header 中使用的样式

| 样式引用 | 用途 | 定义位置 |
|----------|------|----------|
| `s.Base.Foreground(s.Secondary)` | 品牌文字颜色 | styles.go:909 |
| `s.Primary` | 渐变的结束色 | styles.go:458 |
| `s.Secondary` | 渐变的开始色 | styles.go:459 |
| `s.Base.Foreground(s.Primary)` | 进度条颜色 | styles.go:909 |
| `s.Base.Foreground(s.Error)` | 错误计数颜色 | styles.go:909 |
| `s.Muted` | Token百分比样式 | styles.go:910 |
| `s.Subtle` | 提示文本样式 | styles.go:912 |

### 图标定义

```go
// 文件: internal/ui/styles/styles.go (示例)
const (
    ToolPending  = "⏳"
    ToolError    = "✕"
    ToolSuccess  = "✓"
    ErrorIcon    = "×"           // Header 中使用
    ArrowRightIcon = "→"
    RadioOn       = "●"
    RadioOff      = "○"
)
```

---

## 已知问题与根因分析

### 问题1: 窗口resize时的残留字符 (Artifacts/Ghosting)

**症状描述**:
- 窗口缩小时，右侧出现之前渲染的旧字符
- 窗口变大时，右侧出现空白或对齐问题
- resize动画过程中可能出现闪烁

**用户反馈** (根据对话摘要):
> "还是不行" - 多次修复后问题仍然存在

#### 根因分析

**原因1: ANSI序列密度过高**

```
渐变文本产生的ANSI序列密度:
  输入: "CRUSH"
  输出: ~120字节（5个字符 × ~24字节/字符）

  每个字符的ANSI结构:
    \x1b[38;5;60m   - 10字节：设置256色调色板颜色
    \x1b[1m          - 4字节：加粗
    X                - 1字节：实际字符
    \x1b[0m          - 4字节：重置
  总计: 19-24字节/字符

问题:
  1. 可视字符少（5）但字节长度大（120）
  2. lipgloss.Width() 返回5，但实际字节长度是120
  3. 截断时以字节或字符为单位，可能破坏ANSI序列
```

**原因2: 手动ANSI序列检测不完整**

```go
// 当前实现 (header.go:191-203)
if r == '\x1b' {
    end := i + 1
    for end < len(runes) && runes[end] != 'm' {
        end++
    }
    // 问题: 只检测以 'm' 结尾的CSI序列
    //      未处理其他类型:
    //      - OSC序列: \x1b]...\x07
    //      - DCS序列: \x1bP...\x1b\
    //      - CSI参数含分号: \x1b[38;2;255;0;0m
}
```

**原因3: lipgloss的内部缓存（推测）**

```
假设（未证实）:
  lipgloss.NewStyle() 可能内部缓存样式对象
  width相同时可能返回缓存结果
  padding操作可能继承之前的宽度信息

可能性:
  - 在某些lipgloss版本中存在（未确认）
  - 需要检查lipgloss源码验证
```

**原因4: width计算时机**

```go
// 当前流程:
content := b.String()  // 已经包含ANSI序列
contentWidth := lipgloss.Width(content)  // 过滤ANSI后计算

if contentWidth > availableWidth {
    // 截断已经着色的内容
    // 问题: width是基于"过滤ANSI后"的宽度
    //      但截断操作在"包含ANSI的"字节流上进行
}
```

#### 已尝试的解决方案（基于历史记录）

**尝试1: 添加MaxWidth和MaxHeight**

```go
// 早期尝试
truncatedStyle := lipgloss.NewStyle().
    MaxWidth(availableWidth).
    Faint(false)
content = truncatedStyle.Render(content)
```

**结果**: 未解决问题

**原因推测**: `MaxWidth`无法可靠处理密集的ANSI序列

---

**尝试2: 添加MaxHeight(1)限制**

```go
truncatedStyle := lipgloss.NewStyle().
    MaxWidth(availableWidth).
    MaxHeight(1).  // 新增
    Faint(false)
```

**结果**: 部分有效，多行问题改善

**原因**: `MaxHeight(1)`强制单行，但不解决宽度截断

---

**尝试3: 手动ANSI序列截断**

```go
// 当前实现 (header.go:187-215)
for i < len(runes) && currentWidth < availableWidth {
    r := runes[i]
    if r == '\x1b' {
        // 保留完整ANSI序列
        end := i + 1
        for end < len(runes) && runes[end] != 'm' {
            end++
        }
        // ... 包含整个序列
    }
    // ... 追踪可视宽度
}
```

**结果**: 测试通过，但用户报告问题仍存在

**可能原因**:
1. 测试环境和终端环境差异
2. 不同终端对ANSI序列的渲染方式不同
3. resize时机（快速连续resize时的问题）

---

**尝试4: renderDetails中添加MaxHeight(1)**

```go
// header.go:303-306
cwd = lipgloss.NewStyle().
    MaxWidth(max(0, availWidth-lipgloss.Width(metadata))).
    MaxHeight(1).  // 新增
    Render(cwd)
```

**结果**: 预防CWD多行问题

---

### 问题2: Header可能占用两行

**症状**:
- 某些情况下，header渲染为两行
- 导致content从第三行开始

**可视化问题**:

```
正确 (1行):
 Charm™ CRUSH ╱... /projects... ×3 39% ctrl+d
 content starts here
 ...

错误 (2行):
 Charm™ CRUSH ╱... /projects... ×3 39% ctrl+d
 (empty or partial line)
 content starts here
 ...
```

**触发条件分析**:

1. **极窄窗口** (width < 30)
   ```
   contentWidth = ~50 (品牌+标题+进度条)
   availableWidth = < 30
   -> 必须截断
   ```

2. **手动截断失败**
   ```go
   // 如果ANSI序列被部分截断
   // 某些终端可能显示为两行
   ```

3. **Padding引入换行**
   ```go
   // 如果content本身包含隐藏的换行符
   // padding操作可能保留或放大问题
   ```

4. **renderDetails换行**
   ```go
   // CWD MaxWidth没有MaxHeight(1)时
   // 可能产生多行
   ```

#### 当前防护机制

```go
// 第一层: renderDetails MaxHeight(1)
cwd = lipgloss.NewStyle().
    MaxWidth(maxCwdWidth).
    MaxHeight(1).  // ← 防护1
    Render(cwd)

// 第二层: 手动截断
// 逐字符检测，保留ANSI

// 第三层: 降级到lipgloss
truncatedStyle := lipgloss.NewStyle().
    MaxWidth(availableWidth).
    MaxHeight(1).  // ← 防护2
    Faint(false)

// 第四层: 最终安全检查
if strings.ContainsAny(result, "\n\r") {
    result = strings.ReplaceAll(result, "\n", "")  // ← 防护3
    result = strings.ReplaceAll(result, "\r", "")
}
```

### 问题3: 不同终端渲染差异

**可能受影响的终端**:
- Windows Terminal
- iTerm2 (macOS)
- VS Code Terminal
- SSH远程连接

**差异来源**:

1. **ANSI兼容性**
   ```
   终端A: 完全支持256色
   终端B: 只支持16色，降级渲染
   终端C: 某些序列不支持，忽略或显示原始编码
   ```

2. **字符宽度计算**
   ```
   终端A: Unicode字体正确渲染，width准确
   终端B: 等宽字体，某些字符宽度不同
   终端C: 双宽度字符（CJK, emoji）支持差异
   ```

3. **窗口resize响应**
   ```
   终端A: resize时清空缓冲区
   终端B: 保留缓冲区，只是滚动
   终端C: resize后需要手动清屏
   ```

### 根本假设（需验证）

**假设A**: lipgloss内部缓存机制

```go
// 验证方法
s1 := lipgloss.NewStyle().MaxWidth(10)
s1.Render("x")
s2 := lipgloss.NewStyle().MaxWidth(10)
s2.Render("y")

// 检查 s1, s2 是否使用相同缓存
// 或 s1.Render() 是否记住上次的参数
```

**假设B**: 样式对象的可变性

```go
// styles.DefaultStyles() 返回值引用同一对象？
s1 := styles.DefaultStyles()
s2 := styles.DefaultStyles()

// s1 和 s2 是同一个对象还是复制？
// 如果是同一个对象，可能存在状态污染
```

**假设C**: 渐变文本生成效率问题

```go
// ForegroundGrad 每次都重新计算颜色
clusters := ForegroundGrad(t, input, true, color1, color2)

// 如果input相同，是否可以缓存结果？
// 毕竟颜色渐变计算开销较大
```

---

## 调试指南

### 调试环境准备

**添加调试日志**:

```go
import "log"

func (h *HeaderComponent) View() string {
    // 调试点1: 入口状态
    log.Printf("[Header Debug] View() called")
    log.Printf("  width=%d, height=%d", h.width, h.height)
    log.Printf("  brand=%q, title=%q", h.brand, h.title)

    // ... 原有渲染逻辑 ...

    // 调试点2: 内容宽度
    content := b.String()
    contentWidth := lipgloss.Width(content)
    log.Printf("  contentWidth=%d, availableWidth=%d", contentWidth, availableWidth)

    // 调试点3: 截断前字符串
    if len(content) > 100 {
        log.Printf("  content (first 100): %q", content[:100])
    } else {
        log.Printf("  content: %q", content)
    }

    // ... 截断逻辑 ...

    // 调试点4: 结果检查
    result := s.Base.Padding(...).Render(content)
    newlineCount := strings.Count(result, "\n")
    if newlineCount > 0 {
        log.Printf("[Header WARNING] Result contains %d newlines!", newlineCount)
    }

    log.Printf("  resultWidth=%d, targetWidth=%d", lipgloss.Width(result), h.width)

    return result
}
```

### 诊断命令

**检查lipgloss版本**:

```bash
go list -m github.com/charmbracelet/lipgloss
# 确保是最新版本
```

**运行带日志的demo**:

```bash
cd examples/header-demo
go build -o demo.exe main.go
./demo.exe 2>&1 | tee debug.log
```

### 单元测试

**多行检测测试**:

```go
func TestHeaderNeverMultiline(t *testing.T) {
    testCases := []struct {
        name  string
        width int
        title string
        cwd   string
    }{
        {"Normal", 100, "CRUSH", "/projects/ai/Taproot"},
        {"Narrow", 60, "CRUSH", "/very/long/path/to/workspace"},
        {"Very Narrow", 40, "LONGTITLE", "/a/b/c/d/e/f"},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            h := New()
            h.SetSize(tc.width, 1)
            h.SetBrand("Charm™", tc.title)
            h.SetWorkingDirectory(tc.cwd)
            h.SetTokenUsage(50000, 128000, 0.00)
            h.SetErrorCount(3)

            render := h.View()
            newlineCount := strings.Count(render, "\n")

            assert.Equal(t, 0, newlineCount,
                "header should never contain newlines, got %d newlines in:\n%s",
                newlineCount, render)
        })
    }
}
```

**宽度压力测试**:

```go
func TestHeaderWidthStress(t *testing.T) {
    h := New()
    h.SetBrand("Charm", "CRUSH")
    h.SetWorkingDirectory("/very/long/path/to/workspace")
    h.SetTokenUsage(75000, 128000, 0.00)
    h.SetErrorCount(5)

    widths := []int{200, 150, 100, 80, 60, 40, 30, 25, 20}

    for _, w := range widths {
        h.SetSize(w, 1)
        render := h.View()
        actualWidth := lipgloss.Width(render)

        assert.LessOrEqual(t, actualWidth, w,
            "width %d should not exceed target width %d",
            actualWidth, w)

        assert.Equal(t, 0, strings.Count(render, "\n"),
            "width %d: header should be single line", w)
    }
}
```

### 终端测试脚本

**创建测试文件** (test_resize.sh):

```bash
#!/bin/bash

echo "Testing header resize behavior..."

# 启动demo并模拟resize
cd examples/header-demo
go build -o demo.exe main.go

# 使用expect或tmux模拟窗口resize
# 或者手动测试并记录观察点
```

### 性能分析

**添加计时**:

```go
import "time"

func TestHeaderPerformance(t *testing.T) {
    h := New()
    h.SetSize(100, 1)
    h.SetBrand("Charm™", "CRUSH")

    iterations := 1000
    start := time.Now()

    for i := 0; i < iterations; i++ {
        h.SetTokenUsage(i*100, 128000, 0.00)
        _ = h.View()
    }

    elapsed := time.Since(start)
    avgTime := elapsed / time.Duration(iterations)

    t.Logf("Average render time: %v per call", avgTime)

    // 目标: < 1ms per View() call
    if avgTime > time.Millisecond {
        t.Logf("WARNING: Render time exceeds 1ms")
    }
}
```

---

## API 参考

### 构造函数

#### `New() *HeaderComponent`

创建新的header组件实例。

```go
func New() *HeaderComponent
```

**返回值**:
- `*HeaderComponent`: 新初始化的header组件

**初始状态**:
```go
&HeaderComponent{
    brand:       "Charm™",
    title:       "CRUSH",
    tokenMax:    128000,
    compactMode: false,
    // width, height = 0
}
```

**示例**:
```go
h := header.New()
h.SetSize(100, 1)
```

---

### 尺寸管理

#### `SetSize(width, height int)`

设置header的尺寸。

```go
func (h *HeaderComponent) SetSize(width, height int)
```

**参数**:
- `width`: 窗口总宽度（列数），必须 > 0
- `height`: Header高度（行数），通常为 1

**示例**:
```go
h.SetSize(100, 1)  // 标准尺寸
h.SetSize(80, 1)   // 窄窗口
h.SetSize(200, 1)  // 宽窗口
```

**注意事项**:
- 必须在调用 `View()` 前设置
- `height` 参数当前未使用（固定为单行）
- 在 `tea.WindowSizeMsg` 处理中调用

---

#### `Size() (width, height int)`

获取当前header尺寸。

```go
func (h *HeaderComponent) Size() (width, height int)
```

**返回值**:
- `width`: 当前宽度
- `height`: 当前高度

**示例**:
```go
w, h := header.Size()
fmt.Printf("Header: %dx%d\n", w, h)
```

---

### 内容设置

#### `SetBrand(brand, title string)`

设置品牌名称和应用标题。

```go
func (h *HeaderComponent) SetBrand(brand, title string)
```

**参数**:
- `brand`: 品牌名称（如 "Charm™"），可包含特殊字符和emoji
- `title`: 应用标题（如 "CRUSH"），将应用渐变色

**示例**:
```go
h.SetBrand("Charm™", "CRUSH")
h.SetBrand("MyApp™", "CLI Tool")
h.SetBrand("🚀", "Rocket")
```

**效果**:
- `brand` 使用次要色调（Secondary）渲染
- `title` 使用渐变色（Secondary -> Primary）粗体渲染

---

#### `SetSessionTitle(title string)`

设置会话标题（可选）。

```go
func (h *HeaderComponent) SetSessionTitle(title string)
```

**参数**:
- `title`: 会话标题

**注意**:
- 当前 `View()` 方法中未使用此字段
- 为未来功能预留

---

#### `SetWorkingDirectory(cwd string)`

设置当前工作目录路径。

```go
func (h *HeaderComponent) SetWorkingDirectory(cwd string)
```

**参数**:
- `cwd`: 工作目录路径（如 "/projects/ai/Taproot"）

**示例**:
```go
h.SetWorkingDirectory("/projects/ai/Taproot")
h.SetWorkingDirectory("~/workspace")
h.SetWorkingDirectory("")  // 显示为 "~"
```

**自动截断规则**:
- 最多显示4个路径组件
- 超过4个时显示最后4个，前面加 "…"
- 根据可用空间进一步截断

**示例**:
```
"/projects/ai/Taproot"           → "/projects/ai/Taproot"
"/a/b/c/d/e/f"                   → "…/c/d/e/f"
"/very/long/path/name/here/now"   → "…/name/here/now"
```

---

#### `SetTokenUsage(used, max int, cost float64)`

设置Token使用统计。

```go
func (h *HeaderComponent) SetTokenUsage(used, max int, cost float64)
```

**参数**:
- `used`: 已使用的token数量
- `max`: 最大token限制（如果 <= 0，只显示数量不显示百分比）
- `cost`: Token成本（当前未使用）

**示例**:
```go
h.SetTokenUsage(50000, 128000, 0.00)  // 39%
h.SetTokenUsage(128000, 128000, 3.00) // 100%
h.SetTokenUsage(0, 128000, 0.00)      // 0%
h.SetTokenUsage(76432, -1, 0.00)      // 无上限，显示"76432"
```

**显示效果**:
- 如果 `max > 0`: 显示百分比（如 "39%"）
- 如果 `max <= 0`: 显示数量（如 "76432"）
- 进度条显示: 即使0%也显示 minDiags 个斜杠

---

#### `SetErrorCount(count int)`

设置错误计数显示。

```go
func (h *HeaderComponent) SetErrorCount(count int)
```

**参数**:
- `count`: 错误数量（0或负数时不显示错误图标）

**示例**:
```go
h.SetErrorCount(3)   // 显示 "×3"
h.SetErrorCount(0)   // 不显示错误
h.SetErrorCount(1)   // 显示 "×1"
h.SetErrorCount(-5)  // 不显示错误
```

**图标**: 使用 `styles.ErrorIcon`（默认 "×"）

---

### 详情面板状态

#### `SetDetailsOpen(open bool)`

设置详情面板打开状态。

```go
func (h *HeaderComponent) SetDetailsOpen(open bool)
```

**参数**:
- `open`: true 显示 "close"，false 显示 "open"

**示例**:
```go
h.SetDetailsOpen(false)  // 显示 "ctrl+d open "
h.SetDetailsOpen(true)   // 显示 "ctrl+d close"
```

**用途**:
- 显示快捷键提示
- 提示用户可用功能

---

#### `ShowingDetails() bool`

检查详情面板是否打开。

```go
func (h *HeaderComponent) ShowingDetails() bool
```

**返回值**:
- `bool`: 当前打开状态

**示例**:
```go
if h.ShowingDetails() {
    // 执行某些操作
}
```

---

### 显示模式

#### `SetCompactMode(compact bool)`

设置紧凑模式（预留）。

```go
func (h *HeaderComponent) SetCompactMode(compact bool)
```

**参数**:
- `compact`: 是否使用紧凑模式

**注意**:
- 字段已定义但 `View()` 中未实现
- 为未来功能预留

---

### 渲染方法

#### `View() string`

渲染header。

```go
func (h *HeaderComponent) View() string
```

**返回值**:
- `string`: 渲染后的header字符串（包含ANSI颜色代码）

**行为**:
- 如果 `h.brand == ""`，返回空字符串
- 始终渲染为单行
- 完整适应 `h.width` 宽度

**调用时机**:
- Bubble Tea 的 `View()` 方法中调用
- 每次屏幕刷新时调用

**示例**:
```go
func (m *model) View() string {
    var b strings.Builder
    b.WriteString(m.header.View())
    b.WriteString("\n")
    b.WriteString(m.content)
    return b.String()
}
```

---

#### `Update(msg any) (*HeaderComponent, any)`

引擎无关的更新方法（占位符）。

```go
func (h *HeaderComponent) Update(msg any) (*HeaderComponent, any)
```

**返回值**:
- `*HeaderComponent`: 更新后的组件（当前是 self）
- `any`: 命令（当前是 nil）

**注意**:
- 当前实现返回 `(h, nil)`
- 预留用于engine-agnostic设计
- 在BubbleTea中由外层model处理消息

---

## 性能优化建议

### 当前性能瓶颈分析

#### 1. 样式对象重复创建

```go
// 当前实现
s := styles.DefaultStyles()  // 每次View()都调用

// 问题:
// - DefaultStyles() 返回完整的Styles结构
// - 可能涉及大量对象创建
```

**优化建议**:

```go
// 选项1: 缓存styles对象
type HeaderComponent struct {
    // ... 现有字段
    styles *styles.Styles  // 缓存
}

func New() *HeaderComponent {
    return &HeaderComponent{
        styles: styles.DefaultStyles(),  // 只调用一次
        // ...
    }
}

// 选项2: 延迟初始化
func (h *HeaderComponent) getStyles() *styles.Styles {
    if h.styles == nil {
        h.styles = styles.DefaultStyles()
    }
    return h.styles
}
```

#### 2. 渐变文本重新计算

```go
// 当前实现
b.WriteString(styles.ApplyBoldForegroundGrad(&s, h.title, s.Secondary, s.Primary))

// 问题:
// - ForegroundGrad 每次都重新计算颜色渐变
// - 如果title不变，可以缓存结果
```

**优化建议**:

```go
type HeaderComponent struct {
    // ... 现有字段
    title     string
    titleGradient string  // 缓存的渐变文本
    titleHash  uint64    // 用于检测变化
}

// 在Setter中更新缓存
func (h *HeaderComponent) SetBrand(brand, title string) {
    h.brand = brand
    h.title = title
    h.titleGradient = ""  // 清除缓存
}

// 在View()中使用缓存
func (h *HeaderComponent) View() string {
    // ...
    if h.titleGradient == "" {
        h.titleGradient = styles.ApplyBoldForegroundGrad(
            &s, h.title, s.Secondary, s.Primary
        )
    }
    b.WriteString(h.titleGradient)
    // ...
}
```

**注意**:
- 需要权衡缓存节省 vs 内存占用
- 大多数情况下重新计算的开销可接受

#### 3. 字符串拼接优化

```go
// 当前实现: 使用strings.Builder（已经优化）
var b strings.Builder

// 这是最佳实践，无需优化
```

#### 4. lipgloss.Width() 重复调用

```go
// 当前实现
contentWidth := lipgloss.Width(content)  // 调用1
// ...
detailsWidth := lipgloss.Width(details)  // 调用2
// ...
if lipgloss.Width(currentResult) > availableWidth {  // 调用3
```

**优化建议**:

```go
// 缓存已知宽度
contentWidth := lipgloss.Width(content)
// 使用 contentWidth 而不是重复计算
```

### 性能基准

**建议目标**:

| 操作 | 目标时间 | 当前估算 |
|------|----------|----------|
| View() 调用 | < 1ms | ~0.5-1ms |
| Set* 方法 | < 10μs | ~1μs |
| 首次渲染 | < 5ms | ~2-5ms |

**测试方法**:

```go
func BenchmarkHeaderView(b *testing.B) {
    h := New()
    h.SetSize(100, 1)
    h.SetBrand("Charm™", "CRUSH")
    h.SetWorkingDirectory("/projects/ai/Taproot")
    h.SetTokenUsage(50000, 128000, 0.00)
    h.SetErrorCount(3)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = h.View()
    }
}
```

---

## 更新日志

### v1.2.0 (2026-01-30)

**修复**:
- 实现手动ANSI序列截断逻辑（lines 173-228）
- 添加多层防护确保单行输出
- 改进resize时内容处理

**文档**:
- 添加完整的渲染流程文档
- 添加样式系统详解
- 添加调试指南

### v1.1.0 (2026-01-30)

**修复**:
- 添加 MaxHeight(1) 限制防止多行
- 改进宽度截断逻辑
- 增强resize响应

### v1.0.0 (2026-01-29)

**初始版本**:
- 完整的header渲染功能
- 支持品牌、标题、进度条、详情
- 支持窗口resize

---

## 附录: ANSI转义序列参考

### CSI (Control Sequence Introducer) 序列

**格式**: `\x1b[` + 参数 + 终止符

**常见类型**:

| 序列 | 功能 | 示例 |
|------|------|------|
| `\x1b[nm` | SGR (Select Graphic Rendition) | `\x1b[38;5;245m` |
| `\x1b[nK` | 擦除行 | `\x1b[2K` (整行) |
| `\x1b[nL` | 插入行 | `\x1b[2L` |

### SGR参数

**颜色设置**:
```
38;5;n     - 256色调色板前景色
48;5;n     - 256色调色板背景色
38;2;r;g;b - RGB前景色
48;2;r;g;b - RGB背景色
0          - 重置所有属性
1          - 粗体
3/4        - 斜体/下划线
```

**示例**:
```
\x1b[38;5;245m      -> 前景色设为256色调色板第245色
\x1b[38;2;255;0;0m  -> 前景色设为纯红色(RGB)
\x1b[1m             -> 加粗
\x1b[0m             -> 重置所有属性
```

### Header中使用的ANSI序列

**品牌文本**:
```
输入: "Charm™"
输出: \x1b[38;5;245mCharm™\x1b[0m
       └── 颜色 ──┘ └─ 文本 ─┘ └─ 重置 ─┘
```

**渐变标题** (每个字符):
```
输入: "CRUSH"

"C" -> \x1b[38;5;60m\x1b[1mC\x1b[0m
"R" -> \x1b[38;5;58m\x1b[1mR\x1b[0m
...
```

**进度条**:
```
输入: "╱╱╱"
输出: \x1b[38;5;68m╱╱╱\x1b[0m
```

---

## 总结

Header Component 是一个功能完备的TUI标题栏组件，具有以下特点:

**优点**:
1. 清晰的数据驱动设计
2. 完整的渐变文本支持
3. 自适应窗口宽度
4. 多层防护确保正确渲染

**挑战**:
1. 复杂的ANSI序列处理
2. 窗口resize时的残留字符问题
3. 不同终端的兼容性

**未来方向**:
1. 改进ANSI序列截断算法
2. 添加可选的渲染缓存
3. 增强紧凑模式实现
4. 添加更多自定义选项（颜色、样式等）

---

**文档版本**: 1.0
**最后更新**: 2026-01-30
**作者**: Crush Assistant
**项目**: Taproot TUI Framework
