# Taproot v2.0.0 迁移路线图

## 📋 概述

Taproot v2.0.0 是一次重大架构升级，旨在整合 Crush 项目中两个互补的 UI 系统：
- **TUI 系统** (`internal/tui/`) - 基于 Bubbletea，成熟稳定
- **UI 系统** (`internal/ui/`) - 基于 Ultraviolet，性能更优

## 🎯 v2.0.0 目标

1. **双引擎支持**: 同时支持 Bubbletea 和 Ultraviolet 渲染
2. **组件库完善**: 覆盖常见 TUI 组件需求
3. **完全解耦**: 无业务逻辑依赖的通用框架
4. **性能优化**: 利用 Ultraviolet 的直接绘制能力
5. **向后兼容**: 保持 v1.0.0 API 不变

---

## Phase 6: 双引擎基础 (2-3周)

### 6.1 Ultraviolet 集成 ⭐⭐⭐

**目标**: 为 Taproot 添加 Ultraviolet 渲染引擎支持

**任务**:
- [x] 创建 `internal/ui/` 目录 ✅
- [x] 迁移 `internal/ui/list/list.go` ✅
  - [x] Item 接口定义
  - [x] 虚拟化渲染
  - [x] 滚动逻辑
  - [x] 选择管理
  - [x] 过滤支持 (filterable.go)
  - [x] 分组支持 (grouped.go)
- [x] 迁移 `internal/ui/styles/` ✅
  - [x] 布局工具函数
  - [x] Markdown 渲染器
  - [x] Chroma 语法高亮
  - [x] Charmtone 颜色调色板
  - [x] 主题系统
- [x] 创建 `internal/ui/render/` ✅
  - [x] 渲染引擎抽象层
  - [x] DirectEngine 实现
  - [x] 引擎注册系统
  - [ ] Ultraviolet 适配器
  - [ ] Bubbletea 适配器

**源文件**:
```
E:/projects/ai/crush/internal/ui/list/list.go
E:/projects/ai/crush/internal/ui/list/filterable.go
E:/projects/ai/crush/internal/ui/list/grouped.go
E:/projects/ai/crush/internal/ui/common/common.go
E:/projects/ai/crush/internal/ui/common/markdown.go
E:/projects/ai/crush/internal/ui/common/elements.go
```

**预期成果**:
```go
// 使用示例
package main

import (
    "github.com/yourorg/taproot/internal/ui"
    "github.com/yourorg/taproot/internal/ui/list"
    uv "github.com/charmbracelet/ultraviolet"
)

type Model struct {
    list *list.List
}

func (m *Model) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
    m.list.Draw(scr, area)
    return nil
}
```

---

### 6.2 增强对话框系统 ⭐⭐⭐

**目标**: 整合 TUI 和 UI 的对话框框架

**任务**:
- [x] 迁移 `internal/ui/dialog/` ✅
  - [x] Dialog 接口定义
  - [x] Overlay 管理器
  - [x] Action 消息系统
- [x] 从 TUI 对话框中提取通用部分 ✅
  - [x] 按钮组件
  - [x] 输入组件
  - [x] 选择组件
- [x] 创建通用对话框 ✅
  - [x] InfoDialog
  - [x] ConfirmDialog
  - [x] InputDialog
  - [x] SelectListDialog

**源文件**:
```
E:/projects/ai/crush/internal/ui/dialog/dialog.go
E:/projects/ai/crush/internal/ui/dialog/commands.go
E:/projects/ai/crush/internal/ui/dialog/arguments.go
E:/projects/ai/crush/internal/tui/components/dialogs/dialogs.go
```

**预期成果**:
```go
// 对话框接口
type Dialog interface {
    ID() string
    Init() tea.Cmd
    Update(msg tea.Msg) (Dialog, tea.Cmd)
    Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor
}

// 使用示例
dialog := dialogs.NewConfirmDialog(
    "Confirm Action",
    "Are you sure?",
    func(confirmed bool) tea.Cmd {
        if confirmed {
            return executeAction()
        }
        return nil
    },
)
```

---

### 6.3 自动完成组件 ⭐⭐

**目标**: 整合两个系统的自动完成实现

**任务**:
- [x] 合并 TUI 和 UI 的自动完成 ✅
  - [x] 触发字符 (/) ✅
  - [x] 弹窗定位算法 ✅
  - [x] 模糊匹配 ✅
  - [x] 键盘导航 ✅
- [x] 数据提供者接口 ✅
  - [x] FileProvider ✅
  - [x] CommandProvider ✅
  - [x] StringProvider ✅
  - [x] CustomProvider ✅
- [x] 测试用例 ✅
  - [x] 单元测试 (completions_test.go) ✅
  - [x] Provider 测试 ✅
  - [x] 过滤测试 ✅
  - [x] 导航测试 ✅
- [x] 示例程序 ✅
  - [x] examples/autocomplete/demo.go ✅

**源文件**:
```
E:/projects/ai/crush/internal/ui/completions/completions.go
E:/projects/ai/crush/internal/tui/components/completions/completions.go
E:/projects/ai/Taproot/internal/ui/completions/completions.go
E:/projects/ai/Taproot/internal/ui/completions/providers.go
E:/projects/ai/Taproot/internal/ui/completions/completions_test.go
E:/projects/ai/Taproot/examples/autocomplete/demo.go
```

**预期成果**:
```go
// 自动完成组件接口
type AutoCompletion struct {
    provider   Provider
    visible    bool
    cursor     int
    // ...
}

// 提供者接口
type Provider interface {
    GetCompletions() ([]CompletionItem, error)
}

// 提供者类型
type StringProvider struct { /* ... */ }
type FileProvider struct { /* ... */ }
type CommandProvider struct { /* ... */ }

// 使用示例
provider := completions.NewStringProvider([]string{"Apple", "Banana"})
auto := completions.NewAutoCompletion(provider, triggerChar)
auto.SetQuery("Ap")  // 过滤
auto.MoveNext()       // 导航
selected, ok := auto.Select()  // 选择
```

---

## Phase 7: 核心组件库 (3-4周)

### 7.1 文件列表组件 ⭐⭐⭐

**目标**: 创建通用的文件列表显示组件

**任务**:
- [x] 创建 `internal/ui/components/files/` ✅
- [x] 迁移 `internal/tui/components/files/` 逻辑 ✅
- [x] 实现 FileItem 结构 ✅
- [x] 文件图标系统 ✅
- [x] 排序功能 ✅
- [x] 过滤功能 ✅
- [x] 添加/删除行数显示 ✅
- [x] 创建示例程序 ✅ (`examples/files/main.go`, 250 lines, 演示所有功能)

**源文件**:
```
E:/projects/ai/crush/internal/tui/components/files/files.go
E:/projects/ai/crush/internal/ui/chat/sidebar/ (参考)
```

---

### 7.2 树文件组件 ⭐⭐⭐

**目标**: 创建树形文件浏览器组件，支持目录展开/折叠

**任务**:
- [x] 添加 `internal/ui/components/treefiles/` ✅
- [x] 文件节点结构 (`FileNode`, `FileTree`) - 480+ lines ✅
- [x] 树形可视化图标系统 - 190+ lines ✅
- [x] 展开/折叠功能 ✅
- [x] 扁平化显示 (`Flatten()` 方法) ✅
- [x] 排序支持 (名称/大小/时间/类型) ✅
- [x] 隐藏文件切换 ✅
- [x] 最大深度限制 ✅
- [x] 树统计信息 (`Stats()`) ✅
- [x] 综合测试套件 - 490+ lines, 20+ tests ✅
- [x] 交互式树形浏览器演示 (`examples/treefiles/main.go`, 340 lines) ✅

---

### 7.3 状态显示组件 ⭐⭐

**目标**: LSP 和 MCP 状态显示

**任务**:
- [ ] 创建 `internal/ui/components/status/`
- [ ] 迁移 `internal/tui/components/lsp/`
- [ ] 迁移 `internal/tui/components/mcp/`
- [ ] 诊断计数显示
- [ ] 工具数量显示
- [ ] 状态图标

**源文件**:
```
E:/projects/ai/crush/internal/tui/components/lsp/lsp.go
E:/projects/ai/crush/internal/tui/components/mcp/mcp.go
```

---

### 7.4 Diff 查看器完善 ⭐⭐

**目标**: 实现分屏 diff 视图

**任务**:
- [ ] 迁移 `internal/tui/exp/diffview/split.go`
- [ ] 实现分屏布局
- [ ] 同步滚动
- [ ] 语法高亮集成
- [ ] 测试用例

**源文件**:
```
E:/projects/ai/crush/internal/tui/exp/diffview/split.go
E:/projects/ai/crush/internal/tui/exp/diffview/style.go
```

---

## Phase 8: 消息系统 (3-4周)

### 8.1 消息渲染框架 ⭐⭐⭐

**目标**: 创建解耦的消息渲染系统

**任务**:
- [ ] 创建 `internal/ui/components/messages/`
- [ ] 迁移 `internal/ui/chat/` 组件
  - [ ] messages.go - 基础消息
  - [ ] assistant.go - 助手消息
  - [ ] user.go - 用户消息
  - [ ] tools.go - 工具调用
  - [ ] fetch.go - Agentic fetch
  - [ ] diagnostics.go - 诊断信息
  - [ ] todos.go - 任务列表
- [ ] 解耦 message.Message 依赖
  - [ ] 定义通用接口
  - [ ] 适配器模式

**源文件**:
```
E:/projects/ai/crush/internal/ui/chat/*.go
```

---

### 8.2 Markdown 渲染增强 ⭐⭐

**目标**: 更强大的 Markdown 渲染

**任务**:
- [x] 增强 `internal/ui/styles/styles.go` (incorporating chroma/markdown logic) ✅
- [x] 代码块语法高亮 ✅
- [ ] 表格渲染
- [ ] 任务列表
- [ ] 链接处理
- [ ] 图片引用处理

---

### 8.3 任务列表组件 ⭐⭐

**目标**: TODO/Tasks 列表显示

**任务**:
- [ ] 迁移 `internal/ui/chat/todos.go`
- [ ] 任务状态图标
- [ ] 进度条
- [ ] 展开/折叠
- [ ] 动画效果

---

## Phase 9: 布局系统 (2-3周)

### 9.1 通用布局组件 ⭐⭐⭐

**目标**: 创建响应式布局系统

**任务**:
- [x] 创建 `internal/ui/layout/` ✅
- [x] 核心类型和约束 (area.go - Position, Area, Constraints) ✅
  - [x] Fixed 约束
  - [x] Percent 约束
  - [x] Ratio 约束
  - [x] Grow 约束
  - [x] MinSize/MaxSize 约束
- [x] Split 布局 (split.go) ✅
  - [x] SplitVertical - 垂直分割
  - [x] SplitHorizontal - 水平分割
  - [x] CenterRect - 居中矩形
  - [x] TopLeftRect, TopCenterRect, TopRightRect - 顶部位置
  - [x] LeftCenterRect, RightCenterRect - 中部位置
  - [x] BottomLeftRect, BottomCenterRect, BottomRightRect - 底部位置
  - [x] Pad - 统一内边距
  - [x] Inset - 四边独立内边距
- [x] Flexbox 布局 (flex.go) ✅
  - [x] FlexChild - 子元素配置 (支持 grow/shrink)
  - [x] RowLayout - 水平弹性布局
  - [x] ColumnLayout - 垂直弹性布局
  - [x] FlexRow/FlexColumn - 便捷函数
  - [x] 支持固定尺寸、比例分配、自动扩展
- [x] Grid 布局 (grid.go) ✅
  - [x] GridConfig - 网格配置 (行/列/间距)
  - [x] GridLayout - 创建均匀网格
  - [x] GetCell/GetRow/GetColumn - 访问网格单元
  - [x] SpanCell - 跨行跨列
  - [x] FixedGrid - 固定单元格大小
  - [x] UniformGrid - 均匀分布
- [x] 综合测试套件 (layout_test.go, 500+ lines, 30+ tests) ✅
- [x] 交互式布局演示 (examples/layout-demo/main.go, 390+ lines) ✅
  - [x] 8种布局类型演示
  - [x] 实时调整大小
  - [x] 详细信息查看

**参考**:
```
E:/projects/ai/crush/internal/ui/model/ui.go (generateLayout)
E:/projects/ai/crush/internal/ui/common/elements.go
github.com/charmbracelet/ultraviolet (layout primitives)
```

**文件结构** (6 files, ~1500 lines):
```
internal/ui/layout/
├── area.go        (180 lines) - 核心类型和约束
├── split.go       (200 lines) - 分割布局和定位
├── flex.go        (220 lines) - 弹性布局系统
├── grid.go        (230 lines) - 网格布局
└── layout_test.go (500+ lines) - 全面测试

examples/layout-demo/
└── main.go        (390+ lines) - 交互式演示
```

---

### 9.2 侧边栏组件 ⭐⭐

**目标**: 通用侧边栏组件

**任务**:
- [x] 创建 `internal/ui/components/sidebar/` ✅
- [x] 核心类型定义 (types.go) ✅
  - [x] Sidebar 接口
  - [x] ModelInfo - 模型信息 (名称、提供者、推理能力、上下文窗口)
  - [x] SessionInfo - 会话信息 (ID、标题、Token 使用、成本、工作目录)
  - [x] FileInfo - 文件信息 (路径、增删统计)
  - [x] LSPService - LSP 服务状态
  - [x] MCPService - MCP 服务状态
  - [x] Config - 配置选项 (宽度、高度、Logo、模式)
- [x] 主组件实现 (sidebar.go) ✅
  - [x] Logo 显示 (响应式，小屏幕用简化版)
  - [x] Session 标题和信息
  - [x] 工作目录显示
  - [x] 当前模型信息 (带推理状态)
  - [x] Token 使用和成本显示 (百分比、格式化、警告)
  - [x] 文件修改列表 (增删统计)
  - [x] LSP 服务列表 (状态、错误计数)
  - [x] MCP 服务列表 (状态)
  - [x] 紧凑模式支持
  - [x] 响应式布局 (垂直/水平)
- [x] 多面板支持 (✅ - 文件/LSP/MCP 三个面板)
- [x] 紧凑模式 (✅ - 移除 Logo、减少 Padding)
- [x] 综合测试套件 ✅
- [x] 交互式演示 (examples/sidebar-demo) ✅

**源文件**:
```
E:/projects/ai/crush/internal/tui/components/chat/sidebar/sidebar.go
```

**文件结构** (3 files, ~750 lines):
```
internal/ui/components/sidebar/
├── types.go       (140 lines) - 核心接口和类型定义
├── sidebar.go     (550 lines) - 主组件实现
└── sidebar_test.go (420 lines) - 测试套件

examples/sidebar-demo/
└── main.go        (240 lines) - 交互式演示
```

**特性**:
- Engine-agnostic 设计，可配合多种渲染引擎使用
- 响应式布局，根据屏幕尺寸自动调整显示内容
- 支持自定义 Logo 提供者
- 文件路径自动截断
- Token 使用百分比显示 (超过 80% 显示警告)
- 支持多个 LSP 和 MCP 服务
- 可配置显示数量限制

---

### 9.3 头部组件 ⭐⭐

**目标**: 通用头部组件

**任务**:
- [ ] 创建 `internal/ui/components/header/`
- [ ] 迁移 `internal/tui/components/chat/header/`
- [ ] 标题显示
- [ ] 信息区域
- [ ] 操作按钮

**源文件**:
```
E:/projects/ai/crush/internal/tui/components/chat/header/header.go
```

---

## Phase 10: 高级功能 (3-4周)

### 10.1 附件系统 ⭐⭐

**目标**: 文件附件管理

**任务**:
- [ ] 创建 `internal/ui/components/attachments/`
- [ ] 迁移 `internal/ui/attachments/`
- [ ] 图片附件预览
- [ ] 文件附件显示
- [ ] 删除模式
- [ ] 拖拽支持

**源文件**:
```
E:/projects/ai/crush/internal/ui/attachments/attachments.go
```

---

### 10.2 Pills 系统 ⭐⭐

**目标**: TODO/Queue 胶囊显示

**任务**:
- [ ] 创建 `internal/ui/components/pills/`
- [ ] 迁移 `internal/tui/page/chat/pills.go`
- [ ] TODO 胶囊
- [ ] Queue 胶囊
- [ ] 展开/折叠
- [ ] 动画效果

**源文件**:
```
E:/projects/ai/crush/internal/tui/page/chat/pills.go
E:/projects/ai/crush/internal/ui/model/pills.go
```

---

### 10.3 进度条和动画 ⭐

**目标**: 统一的动画系统

**任务**:
- [ ] 增强 `internal/tui/anim/`
- [ ] 进度条组件
- [ ] 加载动画
- [ ] 过渡动画
- [ ] 性能优化

---

## Phase 11: 工具和实用程序 (2-3周)

### 11.1 Shell 执行工具 ⭐⭐

**目标**: 跨平台 shell 命令执行

**任务**:
- [ ] 完善 `internal/tui/util/shell.go`
- [ ] 命令构建器
- [ ] 输出捕获
- [ ] 异步执行
- [ ] 进度回调

---

### 11.2 文件监控 ⭐

**目标**: 文件变化监控

**任务**:
- [ ] 创建 `internal/ui/watch/`
- [ ] fsnotify 集成
- [ ] 事件过滤
- [ ] 防抖动
- [ ] 批量更新

---

### 11.3 剪贴板支持 ⭐

**目标**: 跨平台剪贴板操作

**任务**:
- [ ] 创建 `internal/ui/clipboard/`
- [ ] OSC 52 支持
- [ ] 原生剪贴板
- [ ] 图片支持
- [ ] 历史记录

---

## Phase 12: 文档和示例 (2-3周)

### 12.1 文档完善 ⭐⭐⭐

**任务**:
- [ ] 更新 `docs/ARCHITECTURE.md`
- [ ] 添加 `docs/ULTRAVIOLET.md` - Ultraviolet 指南
- [ ] 添加 `docs/MIGRATION_V2.md` - v1→v2 迁移指南
- [ ] 更新 `docs/API.md`
- [ ] 添加 `docs/EXAMPLES_V2.md`
- [ ] 更新 `AGENTS.md`

---

### 12.2 示例程序 ⭐⭐⭐

**任务**:
- [ ] `examples/ultraviolet/` - UV 引擎演示
- [ ] `examples/dual-engine/` - 双引擎对比
- [ ] `examples/file-browser/` - 文件浏览器
- [ ] `examples/dashboard/` - 仪表板
- [ ] `examples/chat-ui/` - 聊天界面
- [ ] `examples/complete-app/` - 完整应用

---

### 12.3 性能基准 ⭐

**任务**:
- [ ] 创建 `benchmarks/`
- [ ] 列表性能测试
- [ ] 渲染性能测试
- [ ] 内存使用测试
- [ ] 性能优化建议

---

## 🎯 优先级矩阵

### P0 - 必须完成 (核心价值)

| Phase | 组件 | 预估时间 | 价值 |
|-------|------|----------|------|
| 6.1 | Ultraviolet 集成 | 1周 | ⭐⭐⭐ |
| 6.2 | 增强对话框系统 | 1周 | ⭐⭐⭐ |
| 7.1 | 文件列表组件 | 1周 | ⭐⭐⭐ |
| 7.2 | 树文件组件 | 1周 | ⭐⭐⭐ |
| 8.1 | 消息渲染框架 | 2周 | ⭐⭐⭐ |

### P1 - 强烈推荐 (常用功能)

| Phase | 组件 | 预估时间 | 价值 |
|-------|------|----------|------|
| 7.3 | 状态显示组件 | 1周 | ⭐⭐ |
| 9.1 | 通用布局组件 | 1周 | ⭐⭐⭐ |
| 9.2 | 侧边栏组件 | 1周 | ⭐⭐ |
| 12.1 | 文档完善 | 持续 | ⭐⭐⭐ |

### P2 - 推荐完成 (增强功能)

| Phase | 组件 | 预估时间 | 价值 |
|-------|------|----------|------|
| ~~6.3~~ | ~~自动完成组件~~ | ~~1周~~ | ~~⭐⭐~~ |
| 7.4 | Diff 查看器完善 | 1周 | ⭐⭐ |
| 8.2 | Markdown 渲染增强 | 1周 | ⭐⭐ |
| 10.1 | 附件系统 | 1周 | ⭐⭐ |

### P3 - 可选 (特殊场景)

| Phase | 组件 | 预估时间 | 价值 |
|-------|------|----------|------|
| 8.3 | 任务列表组件 | 1周 | ⭐ |
| 9.3 | 头部组件 | 3天 | ⭐ |
| 10.2 | Pills 系统 | 1周 | ⭐ |
| 10.3 | 进度条和动画 | 1周 | ⭐ |

---

## 📊 时间线

```
Month 1:  Phase 6 (双引擎基础)
Month 2:  Phase 7 (核心组件库)
Month 3:  Phase 8 (消息系统) + Phase 9 (布局系统)
Month 4:  Phase 10 (高级功能)
Month 5:  Phase 11 (工具和实用程序)
Month 6:  Phase 12 (文档和示例) + 发布准备
```

**总预估**: 4-6 个月达到生产就绪

---

## 🚀 发布里程碑

### v2.0.0-alpha1 (Month 2)
- ✅ Ultraviolet 集成
- ✅ 增强对话框系统
- ✅ 基础示例

### v2.0.0-beta1 (Month 4)
- ✅ 核心组件库
- ✅ 消息渲染框架
- ✅ 布局系统

### v2.0.0-rc1 (Month 5)
- ✅ 高级功能
- ✅ 完整文档
- ✅ 性能优化

### v2.0.0 (Month 6)
- ✅ 生产就绪
- ✅ 稳定 API
- ✅ 社区反馈

---

## 🔍 成功指标

### 技术指标
- [ ] 渲染性能 < 16ms (60fps)
- [ ] 内存占用 < 50MB (空载)
- [ ] 测试覆盖率 > 70%
- [ ] API 稳定性 > 6个月

### 生态指标
- [ ] 10+ 个示例程序
- [ ] 5+ 个外部项目使用
- [ ] 活跃的社区讨论
- [ ] 完整的文档覆盖

---

## 📚 参考资料

### Crush 项目源码
```
E:/projects/ai/crush/internal/ui/
E:/projects/ai/crush/internal/tui/
```

### 关键文档
- [Ultraviolet 文档](https://github.com/charmbracelet/ultraviolet)
- [Bubbletea 文档](https://github.com/charmbracelet/bubbletea)
- [Crush UI 分析报告](docs/CRUSH_UI_ANALYSIS.md)

---

---
**文档版本**: v2.0.0
**创建日期**: 2025-01-29
**最后更新**: 2025-01-29

---

## 📝 更新日志

### 2025-01-29
- ✅ Phase 6.1 完成: Dual engine foundation complete
  - 渲染引擎抽象层 (`render/`)
  - DirectEngine 测试引擎
  - Bubbletea 适配器 (`adapter_tea.go`)
  - Ultraviolet 适配器 (`adapter_uv.go`)
  - Ultraviolet 演示程序 (`examples/ultraviolet/main.go`)
  - 双引擎对比演示 (`examples/dual-engine/main.go`)
- ✅ Phase 6.2 完成: Dialog system integrated
  - Engine-agnostic dialog framework
  - InfoDialog, ConfirmDialog, InputDialog, SelectListDialog
  - Overlay manager for dialog stacking
- ✅ Phase 6.3 完成: Auto-complete component created
  - Engine-agnostic `AutoCompletion` component (`completions.go`, 230 lines)
  - Three built-in providers: StringProvider, FileProvider, CommandProvider (`providers.go`, 200+ lines)
  - Comprehensive test suite (`completions_test.go`, 330 lines, 7 test suites, 28 subtests)
  - Interactive demo (`examples/autocomplete/demo.go`, 265 lines)
  - Real-time filtering with match highlighting
  - ASCII popup box for completions
- ✅ 主题系统 (theme.go)
  - Markdown 渲染器 (markdown.go)
  - Chroma 语法高亮 (chroma.go)
  - Charmtone 颜色调色板 (palette.go, charmtone.go)
  - 图标系统 (icons.go)
- ✅ 移除旧的 `internal/tui/styles/` 包
- ✅ 所有组件已更新为使用注入式样式
- ✅ Phase 7.1 完成: File list component
  - FileList manager with sorting/filtering (`files.go`, 290 lines)
  - FileItem interface and FileInfo implementation (`types.go`, 206 lines)
  - Icon system with file type detection (`icon.go`)
  - Flexible sorting by name/size/time/extension (`sort.go`)
  - Pattern filtering with wildcard support (`filter.go`)
  - Comprehensive test suite (`files_test.go`, 200+ lines)
  - Interactive file browser demo (`examples/files/main.go`, 250 lines)
- ✅ Phase 7.2 完成: Tree file component
  - FileNode and FileTree structures (`tree.go`, 480+ lines)
  - Tree visualization with expand/collapse support
  - Tree icon system (│ └ ├ ─) (`icons.go`, 190+ lines)
  - File type icons for 50+ formats
  - Flatten() method for visible node traversal
  - ExpandAll() / CollapseAll() bulk operations
  - Sorting by name/size/time/type
  - Hidden file toggle support
  - Max depth limiting for large trees
  - Tree statistics (nodes, files, dirs, size)
  - Comprehensive test suite (`tree_test.go`, 490+ lines, 20+ tests)
  - Interactive tree browser demo (`examples/treefiles/main.go`, 340 lines)
- ✅ Phase 9.1 完成: Layout system
  - Core types and constraints (`area.go`, 180 lines)
    - Fixed, Percent, Ratio, Grow constraints
    - MinSize/MaxSize constraints
    - Area methods (TopLeft, BottomRight, Dx, Dy, Empty, Intersect, Union)
  - Split layout primitives (`split.go`, 200 lines)
    - SplitVertical / SplitHorizontal
    - Absolute positioning: CenterRect, TopLeftRect, TopCenterRect, etc.
    - Padding utilities: Pad, Inset
  - Flexbox layout system (`flex.go`, 220 lines)
    - FlexChild with grow/shrink support
    - RowLayout / ColumnLayout
    - FlexRow / FlexColumn convenience functions
    - Support for fixed sizes, ratios, and auto-expand
  - Grid layout (`grid.go`, 230 lines)
    - GridConfig with rows/cols/gaps
    - GetCell / GetRow / GetColumn utilities
    - SpanCell for cross-row/column spans
    - FixedGrid and UniformGrid
  - Comprehensive test suite (`layout_test.go`, 500+ lines, 30+ tests)
  - Interactive layout demo (`examples/layout-demo/main.go`, 390+ lines)
    - 8 layout type demonstrations
    - Real-time size adjustment
    - Detailed area information view
- ✅ Phase 9.2 完成: Sidebar component
  - Core types and interfaces (`types.go`, 140 lines)
    - Sidebar interface with layout.Sizeable
    - ModelInfo, SessionInfo, FileInfo, LSPService, MCPService
    - Config with width, height, logo, mode, limits
  - Main implementation (`sidebar.go`, 550 lines)
    - Logo display (responsive)
    - Session title and info
    - Working directory
    - Model info with reasoning status
    - Token usage (percentage, warnings > 80%)
    - Modified files with diff stats
    - LSP/MCP services list
    - Compact mode support
    - Responsive layout (vertical/horizontal)
  - Comprehensive test suite (`sidebar_test.go`, 420 lines, 20+ tests)
  - Interactive demo (`examples/sidebar-demo/main.go`, 240 lines)
    - Toggle compact mode, add/remove files, update session, reload data
    - Real-time resize support


