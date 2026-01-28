# Taproot TUI 迁移计划

## 总体目标

将 Crush CLI 项目中的 TUI 框架完整迁移到 Taproot，创建一个可复用的终端 UI 框架库。

## 迁移路线图

```
Phase 1 (已完成) ████████████████████  框架基础
Phase 2 (进行中) ██████████░░░░░░░░░░  应用框架
Phase 3           ░░░░░░░░░░░░░░░░░░░░  通用组件
Phase 4           ░░░░░░░░░░░░░░░░░░░░  对话框系统
Phase 5           ░░░░░░░░░░░░░░░░░░░░  高级组件
```

---

## Phase 1: 框架基础 ✅ 已完成

### 已完成组件

| 组件 | 状态 | 文件 |
|------|------|------|
| 布局接口 | ✅ | `internal/layout/layout.go` |
| 工具函数 | ✅ | `internal/tui/util/util.go` |
| 快捷键 | ✅ | `internal/tui/keys.go` |
| 主题系统 | ✅ | `internal/tui/styles/` |
| 动画组件 | ✅ | `internal/tui/anim/` |
| 核心UI组件 | ✅ | `internal/tui/components/core/` |
| 状态栏 | ✅ | `internal/tui/components/core/status/` |

**代码量**: ~800 行

---

## Phase 2: 应用框架 🚧 进行中

### 目标组件

| 组件 | 优先级 | 复杂度 | 预估工时 |
|------|--------|--------|----------|
| **页面系统** | P0 | 低 | 2h |
| **对话框管理器** | P0 | 中 | 4h |
| **应用主循环** | P0 | 中 | 4h |

### 2.1 页面系统 (page/)

**文件**: `internal/tui/page/`

**功能**:
- 页面标识符 (PageID)
- 页面切换消息 (PageChangeMsg)
- 页面生命周期管理

**实现步骤**:
1. 创建 `internal/tui/page/` 目录
2. 迁移 `page.go`
3. 实现页面注册和切换逻辑
4. 添加页面栈管理（支持前进/后退）

**代码示例**:
```go
// internal/tui/page/page.go
package page

type PageID string

type PageChangeMsg struct {
    ID PageID
}

type PageCloseMsg struct{}

type PageBackMsg struct{}
```

### 2.2 对话框管理器 (dialogs/)

**文件**: `internal/tui/components/dialogs/dialogs.go`

**功能**:
- 对话框堆栈管理
- 层级渲染 (使用 lipgloss.Layer)
- 键盘导航 (ESC关闭)
- 对话框位置管理

**依赖**:
- `internal/tui/util/`
- `github.com/charmbracelet/lipgloss`

**实现步骤**:
1. 创建对话框接口
2. 实现对话框堆栈
3. 添加 Open/Close 消息处理
4. 实现层级渲染

**代码结构**:
```go
// internal/tui/components/dialogs/dialogs.go
package dialogs

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/yourorg/taproot/internal/tui/util"
)

type DialogID string

type DialogModel interface {
    util.Model
    Position() (int, int)
    ID() DialogID
}

type OpenDialogMsg struct { Model DialogModel }
type CloseDialogMsg struct{}

type DialogCmp interface {
    util.Model
    Dialogs() []DialogModel
    HasDialogs() bool
    GetLayers() []*lipgloss.Layer
    ActiveModel() util.Model
}

func NewDialogCmp() DialogCmp
```

### 2.3 应用主循环 (app/)

**文件**: `internal/tui/app/app.go`

**功能**:
- 页面管理
- 对话框管理
- 全局状态
- 键盘路由
- 窗口大小处理

**实现步骤**:
1. 创建应用模型结构
2. 实现页面切换逻辑
3. 集成对话框管理
4. 添加全局快捷键处理

---

## Phase 3: 通用组件 ⏳ 待开始

### 目标组件

| 组件 | 优先级 | 复杂度 | 预估工时 |
|------|--------|--------|----------|
| **自动完成** | P1 | 中 | 6h |
| **虚拟化列表** | P0 | 高 | 12h |
| **Diff查看器** | P1 | 高 | 10h |
| **Logo渲染** | P2 | 低 | 2h |
| **文件列表** | P1 | 中 | 4h |
| **语法高亮** | P2 | 低 | 3h |

### 3.1 自动完成组件 (completions/)

**源文件**: `internal/tui/components/completions/`

**功能**:
- 自动完成弹窗
- 键盘导航
- 模糊匹配
- 多选支持

**依赖**:
- 无外部依赖

**实现步骤**:
1. 迁移 `completions.go`
2. 迁移 `keys.go`
3. 解耦业务逻辑
4. 添加测试

### 3.2 虚拟化列表 (exp/list/)

**源文件**: `internal/tui/exp/list/` (7个文件)

**功能**:
- 窗口化渲染
- 懒加载
- 过滤功能
- 分组支持
- 键盘导航
- 滚动条

**依赖**:
- `internal/tui/util/`
- `internal/tui/styles/`

**文件列表**:
```
list.go          - 主列表组件
items.go         - 列表项类型定义
filterable.go    - 可过滤列表
filterable_group.go - 分组可过滤列表
grouped.go       - 分组列表
keys.go          - 列表快捷键
list_test.go     - 测试
```

**实现步骤**:
1. 迁移基础类型 (items.go)
2. 迁移核心列表 (list.go)
3. 迁移过滤功能 (filterable.go)
4. 迁移分组功能 (grouped.go)
5. 添加测试

### 3.3 Diff查看器 (exp/diffview/)

**源文件**: `internal/tui/exp/diffview/` (7个文件)

**功能**:
- 统一diff视图
- 分屏diff视图
- 语法高亮 (Chroma集成)
- 自定义样式
- 制表符处理

**依赖**:
- `github.com/alecthomas/chroma/v2`
- `internal/tui/styles/`

**文件列表**:
```
diffview.go   - 主diff查看器
split.go      - 分屏布局
style.go      - 样式定义
chroma.go     - Chroma集成
util.go       - 工具函数
diffview_test.go
udiff_test.go
util_test.go
```

**实现步骤**:
1. 迁移核心组件 (diffview.go)
2. 迁移分屏逻辑 (split.go)
3. 迁移样式系统 (style.go)
4. 集成Chroma (chroma.go)
5. 添加测试

### 3.4 其他组件

**Logo渲染**:
- 迁移 `logo.go` 和 `rand.go`
- 依赖: 无

**文件列表**:
- 迁移 `files.go`
- 实现文件图标
- 实现目录遍历

**语法高亮**:
- 迁移 `highlight.go`
- 集成 Chroma 语法高亮

---

## Phase 4: 对话框系统 ⏳ 待开始

### 目标组件

| 组件 | 优先级 | 复杂度 | 预估工时 |
|------|--------|--------|----------|
| **文件选择器** | P1 | 中 | 6h |
| **退出确认** | P2 | 低 | 2h |
| **推理显示** | P2 | 低 | 3h |
| **基础命令面板** | P1 | 高 | 8h |
| **基础模型选择** | P1 | 中 | 6h |
| **基础会话切换** | P1 | 中 | 6h |

### 4.1 文件选择器 (dialogs/filepicker/)

**功能**:
- 目录浏览
- 文件过滤
- 键盘导航
- 隐藏文件显示

**实现步骤**:
1. 迁移 `filepicker.go`
2. 迁移 `keys.go`
3. 使用标准库 `os` 替代 Crush 的文件系统抽象

### 4.2 退出确认 (dialogs/quit/)

**功能**:
- 简单确认对话框
- "是否有未保存的更改"提示

**实现步骤**:
1. 迁移 `quit.go`
2. 迁移 `keys.go`

### 4.3 推理显示 (dialogs/reasoning/)

**功能**:
- 显示AI推理过程
- 可折叠/展开
- Markdown渲染

**实现步骤**:
1. 迁移 `reasoning.go`
2. 解耦Markdown渲染

### 4.4 基础命令面板 (dialogs/commands/)

**功能**:
- 命令列表
- 模糊搜索
- 参数输入
- 命令历史

**解耦策略**:
- 使用回调函数替代直接执行
- 命令提供者接口

**接口设计**:
```go
type CommandProvider interface {
    Commands() []Command
    Execute(cmd Command, args []string) tea.Cmd
    Complete(input string) []Completion
}

type Command struct {
    ID          string
    Label       string
    Description string
    Args        []ArgDef
}

type ArgDef struct {
    Name        string
    Description string
    Required    bool
    Type        ArgType
}
```

### 4.5 基础模型选择 (dialogs/models/)

**功能**:
- 模型列表
- 搜索过滤
- API密钥输入
- 最近使用

**解耦策略**:
- 模型提供者接口
- 配置提供者接口

**接口设计**:
```go
type ModelProvider interface {
    Models() []Model
    RecentModels() []Model
    SetModel(modelID string) error
}

type Model struct {
    ID          string
    Name        string
    Provider    string
    ContextSize int
}
```

### 4.6 基础会话切换 (dialogs/sessions/)

**功能**:
- 会话列表
- 会话搜索
- 新建会话
- 删除会话

**解耦策略**:
- 会话提供者接口

**接口设计**:
```go
type SessionProvider interface {
    Sessions() ([]Session, error)
    GetSession(id string) (*Session, error)
    CreateSession(name string) (*Session, error)
    DeleteSession(id string) error
}

type Session struct {
    ID        string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

## Phase 5: 高级组件 ⏳ 待开始

### 目标组件

| 组件 | 优先级 | 复杂度 | 预估工时 |
|------|--------|--------|----------|
| **图片渲染** | P2 | 高 | 8h |
| **消息渲染** | P2 | 高 | 10h |
| **文本编辑器** | P3 | 极高 | 20h |

### 5.1 图片渲染 (image/)

**功能**:
- 终端图片显示 (kitty, iterm2)
- 图片缩放
- 图片缓存

**挑战**:
- 终端兼容性
- 性能优化

### 5.2 消息渲染 (messages/)

**功能**:
- Markdown渲染
- 代码块语法高亮
- 工具调用显示
- 流式更新

**依赖**:
- `github.com/charmbracelet/glamour`
- `github.com/alecthomas/chroma`

### 5.3 文本编辑器 (editor/)

**功能**:
- 多行文本输入
- 语法高亮
- 自动补全
- 剪贴板支持
- 撤销/重做

**挑战**:
- 复杂度极高
- 跨平台剪贴板
- 建议作为独立项目

---

## 依赖解耦策略

### 策略1: 接口抽象

```go
// 原始代码 (紧耦合)
type Component struct {
    app *crushApp  // 具体依赖
}

// 解耦后
type Component struct {
    provider DataProvider  // 抽象接口
}

type DataProvider interface {
    GetData() ([]Item, error)
    SaveData(item Item) error
}
```

### 策略2: 回调函数

```go
type Component struct {
    onAction func(id string) tea.Cmd
}

func NewComponent(onAction func(string) tea.Cmd) *Component {
    return &Component{onAction: onAction}
}
```

### 策略3: 消息传递

```go
type RequestMsg struct {
    Query string
    ReplyTo chan tea.Msg
}

type ResponseMsg struct {
    Result []string
}

func (c *Component) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case RequestMsg:
        return c, func() tea.Msg {
            // 执行查询
            result := c.query(msg.Query)
            return ResponseMsg{Result: result}
        }
    case ResponseMsg:
        // 处理响应
    }
}
```

---

## 测试策略

### 单元测试

每个组件都需要测试:
- 模型初始化
- 消息处理
- 视图渲染
- 边界条件

### 集成测试

- 页面切换流程
- 对话框打开/关闭
- 组件交互
- 主题切换

### 基准测试

- 渲染性能
- 大数据量处理
- 内存使用

---

## 时间估算

| Phase | 组件数 | 预估工时 | 完成度 |
|-------|--------|----------|--------|
| Phase 1 | 7 | 16h | ✅ 100% |
| Phase 2 | 3 | 10h | 🚧 0% |
| Phase 3 | 6 | 37h | ⏳ 0% |
| Phase 4 | 6 | 31h | ⏳ 0% |
| Phase 5 | 3 | 38h | ⏳ 0% |
| **总计** | **25** | **132h** | **15%** |

---

## 里程碑

- [x] **M1**: 框架基础完成 (Phase 1)
- [ ] **M2**: 应用框架完成 (Phase 2)
- [ ] **M3**: 通用组件完成 (Phase 3)
- [ ] **M4**: 对话框系统完成 (Phase 4)
- [ ] **M5**: 高级组件完成 (Phase 5)

---

## 风险与应对

### 风险1: 依赖复杂

**应对**: 分阶段迁移，优先迁移低依赖组件

### 风险2: 工作量大

**应对**: 社区贡献，分批发布

### 风险3: API设计不确定

**应对**: 保持向后兼容，使用抽象接口

### 风险4: 性能问题

**应对**: 基准测试，性能优化

---

## 下一步行动

1. ✅ 完成 Phase 2 - 页面系统和对话框管理器
2. ⏳ 实现 Phase 3.1 - 自动完成组件
3. ⏳ 实现 Phase 3.2 - 虚拟化列表
4. ⏳ 编写更多示例程序
