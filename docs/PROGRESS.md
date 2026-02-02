# Taproot TUI Framework - Migration Progress

## Overview

Taproot 是从 Crush CLI 提取的 TUI 框架，提供可复用的终端 UI 组件。

**进度**: Phase 1-4 完成, Phase 5 60% (约 97%)

---

## 已完成组件 ✅

### Phase 1: 框架基础 (100%)

| 组件 | 文件 | 功能 | 代码行数 |
|------|------|------|----------|
| **布局接口** | `internal/layout/layout.go` | Focusable, Sizeable, Positional, Help | ~40 |
| **工具函数** | `internal/tui/util/util.go` | Model, InfoMsg, ExecShell | ~70 |
| **快捷键** | `internal/tui/keys.go` | KeyMap, DefaultKeyMap | ~30 |
| **主题系统** | `internal/ui/styles/` | Theme, Manager, 颜色混合 | ~350 |
| **动画组件** | `internal/tui/anim/` | 渐变加载动画 | ~250 |
| **核心UI** | `internal/tui/components/core/` | Section, Title, Button | ~180 |
| **状态栏** | `internal/tui/components/core/status/` | 状态栏组件 | ~100 |

**小计**: ~1,020 行

### Phase 2: 应用框架 (100%)

| 组件 | 文件 | 功能 | 代码行数 |
|------|------|------|----------|
| **页面系统** | `internal/tui/page/page.go` | PageID, PageChangeMsg | ~15 |
| **对话框管理** | `internal/tui/components/dialogs/dialogs.go` | DialogCmp, 堆栈管理 | ~140 |
| **应用主循环** | `internal/tui/app/app.go` | AppModel, 页面/对话框集成 | ~150 |

**小计**: ~305 行

### Phase 3: 通用组件 (部分完成)

| 组件 | 文件 | 功能 | 代码行数 |
|------|------|------|----------|
| **Logo渲染** | `internal/tui/components/logo/` | ASCII logo, 渐变 | ~280 |

**小计**: ~280 行

### Phase 4: 对话框系统 (66%)

|| 组件 | 文件 | 功能 | 代码行数 |
||------|------|------|----------|
|| **命令面板** | `internal/tui/components/dialogs/commands/` | 命令列表, 参数输入 | ~330 |
|| **模型选择** | `internal/tui/components/dialogs/models/` | 模型列表, 最近使用 | ~260 |
|| **文件选择** | `internal/tui/components/dialogs/filepicker/` | 目录浏览, 文件过滤 | ~280 |
|| **退出确认** | `internal/tui/components/dialogs/quit/` | 未保存检查 | ~110 |

**小计**: ~980 行

**待完成** (34%):
- 推理显示 (reasoning/)
- 会话切换 (sessions/)

---

### Phase 8: 消息系统 (100%)

|| 组件 | 文件 | 功能 | 代码行数 |
||------|------|------|----------|
|| **消息接口** | `ui/components/messages/types.go` | Message, MessageItem, Focusable, Expandable | ~250 |
|| **助手消息** | `ui/components/messages/assistant.go` | Markdown渲染, Token统计, 可展开 | ~200 |
|| **用户消息** | `ui/components/messages/user.go` | 代码块, 文件附件, 复制模式 | ~250 |
|| **工具消息** | `ui/components/messages/tools.go` | 工具调用详情, 状态跟踪 | ~300 |
|| **Fetch消息** | `ui/components/messages/fetch.go` | Agentic fetch, 嵌套消息, 树形渲染 | ~730 |
|| **诊断消息** | `ui/components/messages/diagnostics.go` | 诊断汇总, 代码高亮, 可展开 | ~200 |
|| **TODO消息** | `ui/components/messages/todos.go` | TODO列表, 进度条, 状态图标 | ~540 |
|| **Markdown** | `ui/styles/markdown.go` | 表格, 任务列表, 链接, 图片渲染 | ~400 |

**小计**: ~3,040 行 (组件) + ~400 行 (Markdown) = ~3,440 行

**测试覆盖**:
- `messages_test.go`: ~570 lines, 60+ tests
- `markdown_test.go`: ~250 lines, 10+ tests
- 所有测试通过 ✅

| 示例 | 文件 | 功能 |
|------|------|------|
| **demo** | `examples/demo/main.go` | 简单计数器 |
| **list** | `examples/list/main.go` | 可选择列表 |
| **app** | `examples/app/main.go` | 页面/对话框演示 |
| **completions** | `examples/completions/main.go` | 自动完成演示 |
| **commands** | `examples/commands/main.go` | 命令面板演示 |
| **models** | `examples/models/main.go` | 模型选择演示 |
| **filepicker** | `examples/filepicker/main.go` | 文件选择器演示 |
| **quit** | `examples/quit/main.go` | 退出确认演示 |
| **reasoning** | `examples/reasoning/main.go` | 推理显示演示 |
| **sessions** | `examples/sessions/main.go` | 会话切换演示 |
| **diffview** | `examples/diffview/main.go` | Diff 查看器演示 |
| **filterablelist** | `examples/filterablelist/main.go` | 过滤列表演示 |
| **groupedlist** | `examples/groupedlist/main.go` | 分组列表演示 |

| **attachments** | `examples/attachments/main.go` | 附件列表演示 |
| **pills** | `examples/pills/main.go` | 胶囊状态列表演示 |
| **progress** | `examples/progress/main.go` | 进度条和动画演示 |

---

## Phase 10: 高级功能 ✅

### Phase 10.1: 附件系统

| 组件 | 文件 | 功能描述 |
|------|------|----------|
| **AttachmentList** | `ui/components/attachments/attachments.go` | 文件附件列表组件，支持文件类型检测、过滤、选择、统计 |
| **Attachment Types** | `ui/components/attachments/types.go` | 定义附件类型（文件/图片/视频/音频/文档/归档）、MIME类型检测、大小格式化 |

**测试文件**: `ui/components/attachments/attachments_test.go` (409 行，18+ 测试)

**核心功能**:
- 6种附件类型自动检测（50+文件扩展名）
- MIME类型自动识别
- 文件大小格式化（KB/MB/GB/TB）
- 过滤功能支持
- 多选/单选支持
- 展开/折叠支持
- 统计功能（总数、过滤数、选择数）
- 渲染缓存优化
- 引擎无关设计（实现render.Model）

### Phase 10.2: 胶囊系统

| 组件 | 文件 | 功能描述 |
|------|------|----------|
| **PillList** | `ui/components/pills/pills.go` | 状态胶囊列表组件，支持7种预设状态、展开/折叠、批量操作 |

**测试文件**: `ui/components/pills/pills_test.go` (452 行，23+ 测试)

**核心功能**:
- 7种预设状态（待处理/进行中/完成/错误/警告/信息/中性）
- 状态图标（☐ ⟳ ✓ × ⚠ ℹ •）
- 行内模式支持
- 展开/折叠支持
- 批量展开/折叠
- 徽章计数支持
- 自定义模式（用户定义图标）
- 渲染缓存优化

### Phase 10.3: 进度条和动画

| 组件 | 文件 | 功能描述 |
|------|------|----------|
| **ProgressBar** | `ui/components/progress/progressbar.go` | 进度条组件，支持百分比显示、标签、自定义颜色 |
| **Spinner** | `ui/components/progress/spinner.go` | 动画旋转器，支持4种动画类型、可配置FPS、状态管理 |

**测试文件**: `ui/components/progress/progress_test.go` (419 行，30+ 测试)

**核心功能**:
- 进度条可视化：████████░░░░░░░ 8/10 (80%)
- 百分比显示
- 标签支持（分离或内联）
- 自定义颜色
- 边界检查（限制在0-total范围）
- 4种预设动画类型（Dots/Line/Arrow/Moon）
- 可配置FPS
- 状态管理（Start/Stop/Reset）
- 双Tick消息支持（tea.Tick和render.Tick）


### 文档

| 文档 | 路径 | 内容 |
|------|------|------|
| **架构分析** | `docs/ARCHITECTURE.md` | Crush TUI 架构完整分析 |
| **迁移计划** | `docs/MIGRATION_PLAN.md` | 5阶段迁移路线图 |
| **替代方案** | `docs/ALTERNATIVES.md` | 技术选型分析 |
| **任务清单** | `docs/TASKS.md` | 详细待办事项 |
| **开发指南** | `AGENTS.md` | Agent 工作指南 |
| **V2 路线图** | `docs/V2_ROADMAP.md` | v2.0 完整开发路线图 |
| **Phase 7 摘要** | `docs/PHASE_7_SUMMARY.md` | Phase 7.1-7.3 完成报告 |
| **Phase 8 摘要** | `docs/PHASE_8_SUMMARY.md` | Phase 8 消息系统完成报告 |
| **Phase 10 摘要** | `docs/PHASE_10_SUMMARY.md` | Phase 10 高级功能完成报告 |

---

## 总体统计

```
已完成代码: ~13,350 行
完成阶段: Phase 1 + Phase 2 + Phase 3 + Phase 4 + Phase 5 + Phase 8 + Phase 10 (100%)
组件数量: 55+ 核心组件
示例程序: 16 个
文档页数: 8 个
```

---

### 下一步计划

#### 下一步计划

#### 已完成 ✅

1. ✅ **图片渲染** (image/) - 已完成
2. ✅ **消息渲染** (messages/) - 已完成
3. ✅ **README.md** - 已完成
4. ✅ **API.md** - 已完成
5. ✅ **Markdown 渲染** (styles/markdown.go) - 已完成
6. ✅ **Chroma 语法高亮** (styles/chroma.go, highlight/) - 已完成
7. ✅ **Charmtone 调色板** (styles/palette.go) - 已完成

#### 可选任务

1. **编写更多测试** - 提高代码覆盖率
2. **代码质量优化** - 修复剩余 diagnostics
3. **EXAMPLES.md** - 示例集合文档
4. **CONTRIBUTING.md** - 贡献指南
5. **CHANGELOG.md** - 变更日志
6. **发布准备** - 版本号, 发布说明

---

## 技术亮点

### 已实现特性

✅ **主题系统**
- 动态主题切换
- HCL 色彩空间混合
- 渐变文本渲染
- 20+ 预定义颜色

✅ **对话框管理**
- 对话框堆栈
- 键盘导航 (ESC关闭)
- 位置管理

✅ **页面系统**
- 页面切换
- 页面栈 (支持返回)
- 生命周期管理

✅ **状态栏**
- InfoMsg 类型 (Info/Success/Warn/Error)
- TTL 自动清除
- Help 集成

✅ **动画**
- 渐变色彩
- 错位入场
- 省略号动画

---

## 依赖关系

```
外部依赖:
├── github.com/charmbracelet/bubbletea (v1.3.10)
├── github.com/charmbracelet/bubbles (v0.21.0)
├── github.com/charmbracelet/lipgloss (v1.1.x)
├── github.com/charmbracelet/glamour (v0.8.0) ✅ 新增
├── github.com/alecthomas/chroma/v2 (v2.23.1) ✅ 新增
├── github.com/charmbracelet/x/ansi (v0.11.4)
├── github.com/lucasb-eyer/go-colorful (v1.3.0)
└── mvdan.cc/sh/v3 (v3.12.0)

无内部依赖 - 完全解耦! ✅
```

---

## 与 Crush 对比

| 特性 | Crush | Taproot | 状态 |
|------|-------|---------|------|
| TUI 框架 | ✅ | ✅ | 已迁移 |
| 主题系统 | ✅ | ✅ | 已迁移 |
| 动画 | ✅ | ✅ | 已迁移 (简化) |
| 状态栏 | ✅ | ✅ | 已迁移 |
| 对话框管理 | ✅ | ✅ | 已迁移 (无 Layer) |
| 页面系统 | ✅ | ✅ | 已迁移 |
| Logo | ✅ | ✅ | 已迁移 (改为 Taproot) |
| 自动完成 | ✅ | ⏳ | 待迁移 |
| 虚拟化列表 | ✅ | ⏳ | 待迁移 |
| Diff查看器 | ✅ | ⏳ | 待迁移 |
| 文件选择器 | ✅ | ⏳ | 待迁移 |
| 编辑器 | ✅ | ❌ | 复杂度太高,暂不迁移 |
| 聊天组件 | ✅ | ❌ | 业务耦合,不适合框架 |

---

## 测试状态

```bash
$ go test ./...
?   	github.com/wwsheng009/taproot/examples/app	[no test files]
?   	github.com/wwsheng009/taproot/examples/demo	[no test files]
?   	github.com/wwsheng009/taproot/examples/list	[no test files]
ok  	github.com/wwsheng009/taproot/internal/layout	(cached)
?   	github.com/wwsheng009/taproot/internal/tui/*	[no test files]
```

✅ 所有包编译通过
⏳ 测试覆盖率待提高

---

## 使用示例

### 简单页面应用

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/wwsheng009/taproot/internal/tui/app"
    "github.com/wwsheng009/taproot/internal/tui/page"
)

func main() {
    application := app.NewApp()
    application.RegisterPage("home", HomePage{})
    application.SetPage("home")
    
    p := tea.NewProgram(application, tea.WithAltScreen())
    p.Run()
}

type HomePage struct{}

func (h HomePage) Init() tea.Cmd { return nil }
func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) { return h, nil }
func (h HomePage) View() string { return "Hello, Taproot!" }
```

### 使用对话框

```go
// 打开对话框
return func() tea.Msg {
    return dialogs.OpenDialogMsg{Model: MyDialog{}}
}

// 关闭对话框
return func() tea.Msg {
    return dialogs.CloseDialogMsg{}
}
```

### 使用主题

```go
t := styles.CurrentTheme()
text := t.S().Base.Foreground(t.Primary).Render("Hello")
gradient := styles.ApplyForegroundGrad("Text", t.Primary, t.Secondary)
```

---

## 性能特点

- **零拷贝**: 使用 `strings.Builder` 优化字符串拼接
- **缓存**: 主题对象单例,动画帧预渲染
- **虚拟化**: (待实现) 列表组件支持大数据

---

## 已知限制

1. **lipgloss.Layer**: 公共版本可能不支持层级渲染,已简化实现
2. **剪贴板**: 未迁移,需要平台特定代码
3. **Markdown**: 已集成 glamour,提供主题化渲染 ✅
4. **语法高亮**: 已集成 chroma,支持自动语言检测 ✅
5. **编辑器**: 复杂度太高,建议使用 bubbles/textarea

---

## 贡献指南

### 添加新组件

1. 在 `internal/tui/components/` 创建目录
2. 实现 `util.Model` 接口
3. 添加测试
4. 创建示例程序
5. 更新文档

### 代码规范

- 包名: 小写
- 接口: `-able` 后缀 (Focusable, Sizeable)
- 函数: PascalCase (导出), camelCase (内部)
- 样式: 使用 `styles.CurrentTheme()`

---

## v2.0.0 开发进度 (新架构)

### Phase 6: 双引擎基础 (50% 完成 ✅)

| 组件 | 文件 | 功能 | 代码行数 | 状态 |
|------|------|------|----------|------|
| **渲染类型** | `internal/ui/render/types.go` | Model, Msg, Cmd, KeyMsg | 140 | ✅ |
| **引擎注册** | `internal/ui/render/engine.go` | Engine, Factory, Registry | 108 | ✅ |
| **Direct 引擎** | `internal/ui/render/direct.go` | 测试用直接引擎 | 249 | ✅ |
| **Bubbletea 适配器** | `internal/ui/render/adapter_tea.go` | Bubbletea 集成 | 172 | ✅ |
| **Ultraviolet 适配器** | `internal/ui/render/adapter_uv.go` | Ultraviolet 集成 | 163 | ✅ |
| **渲染测试** | `internal/ui/render/render_test.go` | 单元测试 | 303 | ✅ |
| **UV 示例** | `examples/ultraviolet/main.go` | Ultraviolet 演示 | 120 | ✅ |
| **双引擎示例** | `examples/dual-engine/main.go` | 引擎对比 | 170 | ✅ |

**小计**: ~1,425 行

**已完成**:
- ✅ 引擎抽象层 (`Engine` 接口)
- ✅ 引擎工厂模式 (注册 + 创建)
- ✅ Direct 引擎 (用于测试)
- ✅ Bubbletea 适配器 (无缝集成)
- ✅ Ultraviolet 适配器 (高性能渲染)
- ✅ 示例程序 (ultraviolet, dual-engine)
- ✅ 文档更新 (UI_EXAMPLES.md)

**下一步**:
- Phase 6.2: 增强对话框系统 ✅ (已在 `internal/ui/dialog/`)
- Phase 6.3: 自动完成组件 ✅ (已完成)
- Phase 7: 核心组件库

---

### Phase 6.2: 对话框系统 (100% 完成 ✅)

| 组件 | 文件 | 功能 | 代码行数 | 状态 |
|------|------|------|----------|------|
| **对话框类型** | `internal/ui/dialog/types.go` | Dialog 接口, Action | 80 | ✅ |
| **UI 组件** | `internal/ui/dialog/*.go` | Button, Input, SelectList | 200+ | ✅ |
| **对话框实现** | `internal/ui/dialog/dialogs.go` | Info, Confirm, Input, Select | 250+ | ✅ |
| **输入对话框** | `internal/ui/dialog/input.go` | 文本输入组件 | 150+ | ✅ |
| **覆盖层管理** | `internal/ui/dialog/overlay.go` | 对话框堆栈 | 200+ | ✅ |
| **对话框测试** | `internal/ui/dialog/dialog_test.go` | 单元测试 | 100+ | ✅ |
| **对话框示例** | `examples/ui-dialogs/` | 交互演示 | - | ✅ |

**小计**: ~1,000 行

### Phase 6.3: 自动完成组件 (100% 完成 ✅)

| 组件 | 文件 | 功能 | 代码行数 | 状态 |
|------|------|------|----------|------|
| **自动完成核心** | `internal/ui/completions/completions.go` | AutoCompletion 核心逻辑 | 230 | ✅ |
| **数据提供者** | `internal/ui/completions/providers.go` | String/File/Command Provider | 200+ | ✅ |
| **单元测试** | `internal/ui/completions/completions_test.go` | 完整测试覆盖 | 330 | ✅ |
| **示例程序** | `examples/autocomplete/demo.go` | 交互式演示 | 265 | ✅ |

**小计**: ~1,025 行

**已完成**:
- ✅ 引擎无关的自动完成组件
- ✅ Provider 接口 (StringProvider, FileProvider, CommandProvider)
- ✅ 过滤和虚拟化滚动
- ✅ 键盘导航 (上下、选择、关闭)
- ✅ 维度自适应
- ✅ 完整单元测试 (所有测试通过)
- ✅ 交互式示例程序

**特性**:
- 模糊搜索 (case-insensitive)
- 动态宽高计算
- 虚拟化视图 (支持大量项目)
- 可配置的最小/最大尺寸
- 提供者模式便于扩展

---

## 路线图更新

```
2024-01-28: Phase 1 + Phase 2 完成 ✅
2024-01-28: Phase 3 完成 (Logo, Lists, Completions, DiffView) ✅
2024-01-28: Phase 4 完成 (所有对话框组件) ✅
2024-01-28: Phase 5.2 完成 (Messages) ✅
2024-01-28: Phase 5.1 完成 (Image) ✅
2024-01-28: README.md 完成 ✅
2024-01-28: API.md 完成 ✅
2024-01-28: Markdown 渲染完成 ✅
2024-01-28: Chroma 语法高亮完成 ✅
2024-01-28: **Taproot TUI Framework v1.0.0 发布就绪！** 🎉
2024-01-29: Phase 6.1 完成 (双引擎基础) ✅
2024-01-29: Phase 6.2 完成 (对话框系统) ✅
2024-01-29: Phase 6.3 完成 (自动完成组件) ✅
2024-01-29: **Taproot v2.0.0-alpha1 就绪！** 🚀
2024-02-02: Phase 7.1 完成 (文件列表组件) ✅
2024-02-02: Phase 7.2 完成 (树文件组件) ✅
2024-02-02: Phase 7.3 完成 (状态显示组件) ✅
2024-02-02: Phase 7.4 完成 (Diff 查看器完善) ✅
```

---

### Phase 7.4: Diff 查看器完善 (100% 完成 ✅)

|| 组件 | 文件 | 功能 | 代码行数 | 状态 |
||------|------|------|----------|------|
|| **DiffView 核心** | `tui/exp/diffview/diffview.go` | Unified/Split 分屏视图 | 687 | ✅ |
|| **Split Hunk 转换** | `tui/exp/diffview/split.go` | 分屏 diff 转换逻辑 | 71 | ✅ |
|| **样式定义** | `tui/exp/diffview/style.go` | 主题样式配置 | 101 | ✅ |
|| **测试套件** | `tui/exp/diffview/diffview_test.go` | 26 个测试用例 | 683 | ✅ |

**小计**: ~1,542 行

**已完成**:
- ✅ Unified 和 Split 两种布局模式
- ✅ Split view 支持语法高亮 (之前回退到 unified，现在修复)
- ✅ 同步滚动 (垂直和水平)
- ✅ 行号显示
- ✅ 响应式列宽
- ✅ 完整的测试覆盖 (26 个测试全部通过)

**新增功能**:
- Split view 现在支持语法高亮（不再回退到 unified view）
- 水平滚动在 split view 中正常工作
- Split view specific tests (rendering, alignment, horizontal scrolling)

**测试详情**:
- `TestSplitViewRendering` - 验证 split view 基本渲染
- `TestSplitViewWithSyntaxHighlighting` - 验证 split view 语法高亮
- `TestSplitViewScrolling` - 验证 split view 垂直滚动
- `TestSplitHorizontalScrolling` - 验证 split view 水平滚动

---

### Phase 7: 核心组件库 (部分完成)

|| 组件 | 进度 | 状态 |
||------|------|------|
|| Phase 7.1: 文件列表组件 | 100% | ✅ 完成 |
|| Phase 7.2: 树文件组件 | 100% | ✅ 完成 |
|| Phase 7.3: 状态显示组件 | 100% | ✅ 完成 |
|| Phase 7.4: Diff 查看器完善 | 100% | ✅ 完成 |

---


### v2.0.0 总体进度

| Phase | 描述 | 进度 | 状态 |
|-------|------|------|------|
| Phase 6.1 | 双引擎基础 | 100% | ✅ 完成 |
| Phase 6.2 | 对话框系统 | 100% | ✅ 完成 |
| Phase 6.3 | 自动完成组件 | 100% | ✅ 完成 |
| Phase 6 总计 | 双引擎基础架构 | 100% | ✅ 完成 |
| Phase 7 | 核心组件库 | 0% | ⏳ 待开始 |

**总代码行数**: ~1,450 行

---

**最后更新**: 2024-01-29
**当前版本**: 2.0.0-alpha1
**状态**: Phase 6 完成，准备进入 Phase 7 🚀
