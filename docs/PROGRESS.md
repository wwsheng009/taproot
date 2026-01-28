# Taproot TUI Framework - Migration Progress

## Overview

Taproot 是从 Crush CLI 提取的 TUI 框架，提供可复用的终端 UI 组件。

**进度**: Phase 1, 2, 3 完成 (约 70%)

---

## 已完成组件 ✅

### Phase 1: 框架基础 (100%)

| 组件 | 文件 | 功能 | 代码行数 |
|------|------|------|----------|
| **布局接口** | `internal/layout/layout.go` | Focusable, Sizeable, Positional, Help | ~40 |
| **工具函数** | `internal/tui/util/util.go` | Model, InfoMsg, ExecShell | ~70 |
| **快捷键** | `internal/tui/keys.go` | KeyMap, DefaultKeyMap | ~30 |
| **主题系统** | `internal/tui/styles/` | Theme, Manager, 颜色混合 | ~350 |
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

### 示例程序

| 示例 | 文件 | 功能 |
|------|------|------|
| **demo** | `examples/demo/main.go` | 简单计数器 |
| **list** | `examples/list/main.go` | 可选择列表 |
| **app** | `examples/app/main.go` | 页面/对话框演示 |

### 文档

| 文档 | 路径 | 内容 |
|------|------|------|
| **架构分析** | `docs/ARCHITECTURE.md` | Crush TUI 架构完整分析 |
| **迁移计划** | `docs/MIGRATION_PLAN.md` | 5阶段迁移路线图 |
| **替代方案** | `docs/ALTERNATIVES.md` | 技术选型分析 |
| **任务清单** | `docs/TASKS.md` | 详细待办事项 |
| **开发指南** | `AGENTS.md` | Agent 工作指南 |

---

## 总体统计

```
已完成代码: ~3,105 行
完成阶段: Phase 1 + Phase 2 + Phase 3 (70%)
组件数量: 17 个核心组件
示例程序: 8 个
文档页数: 5 个
```

---

## 下一步计划

### 短期 (Phase 3 继续)

1. **自动完成组件** (completions/) - 高价值
   - 模糊搜索
   - 键盘导航
   - 预估: 5-6h

2. **虚拟化列表** (exp/list/) - 极高价值
   - 大数据支持
   - 过滤/分组
   - 预估: 10-12h

3. **Diff查看器** (exp/diffview/) - 高价值
   - 统一/分屏diff
   - 语法高亮
   - 预估: 8-10h

### 中期 (Phase 4)

4. **文件选择器** (dialogs/filepicker/)
5. **退出确认** (dialogs/quit/)
6. **命令面板基础** (dialogs/commands/)
7. **模型选择基础** (dialogs/models/)

### 长期 (Phase 5)

8. **图片渲染** (image/)
9. **消息渲染** (messages/)

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
├── github.com/charmbracelet/x/ansi (v0.11.4)
├── github.com/charmbracelet/x/exp/charmtone
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
?   	github.com/yourorg/taproot/examples/app	[no test files]
?   	github.com/yourorg/taproot/examples/demo	[no test files]
?   	github.com/yourorg/taproot/examples/list	[no test files]
ok  	github.com/yourorg/taproot/internal/layout	(cached)
?   	github.com/yourorg/taproot/internal/tui/*	[no test files]
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
    "github.com/yourorg/taproot/internal/tui/app"
    "github.com/yourorg/taproot/internal/tui/page"
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
3. **Markdown**: 依赖 glamour,需要额外集成
4. **编辑器**: 复杂度太高,建议使用 bubbles/textarea

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

## 路线图更新

```
2024-01-28: Phase 1 + Phase 2 完成 ✅
2024-01-28: Phase 3.4 (Logo) 完成 ✅
2024-01-28: Phase 3.2 (简化列表) 完成 ✅
2024-01-28: Phase 3.1 (自动完成) 完成 ✅
2024-01-28: Phase 3.2 (过滤列表) 完成 ✅
2024-01-28: Phase 3.2 (分组列表) 完成 ✅
2024-01-28: Phase 3.3 (Diff查看器) 完成 ✅
2024-01-29: Phase 4.4 (命令面板) - 计划中
2024-01-30: Phase 4.5 (模型选择) - 计划中
```

---

**最后更新**: 2024-01-28
**当前版本**: 0.4.0
**下一步**: Phase 4 对话框系统 (命令面板, 模型选择)
