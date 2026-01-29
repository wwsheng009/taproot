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
- [ ] 迁移 `internal/ui/dialog/`
  - [ ] Dialog 接口定义
  - [ ] Overlay 管理器
  - [ ] Action 消息系统
- [ ] 从 TUI 对话框中提取通用部分
  - [ ] 按钮组件
  - [ ] 输入组件
  - [ ] 选择组件
- [ ] 创建通用对话框
  - [ ] InfoDialog
  - [ ] ConfirmDialog
  - [ ] InputDialog
  - [ ] SelectListDialog

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
- [ ] 合并 TUI 和 UI 的自动完成
  - [ ] 触发字符 (@)
  - [ ] 弹窗定位算法
  - [ ] 模糊匹配
  - [ ] 键盘导航
- [ ] 数据提供者接口
  - [ ] FileProvider
  - [ ] CommandProvider
  - [ ] CustomProvider

**源文件**:
```
E:/projects/ai/crush/internal/ui/completions/completions.go
E:/projects/ai/crush/internal/tui/components/completions/completions.go
```

---

## Phase 7: 核心组件库 (3-4周)

### 7.1 文件列表组件 ⭐⭐⭐

**目标**: 创建通用的文件列表显示组件

**任务**:
- [ ] 创建 `internal/ui/components/files/`
- [ ] 迁移 `internal/tui/components/files/` 逻辑
- [ ] 实现 FileItem 结构
- [ ] 文件图标系统
- [ ] 排序功能
- [ ] 过滤功能
- [ ] 添加/删除行数显示

**源文件**:
```
E:/projects/ai/crush/internal/tui/components/files/files.go
E:/projects/ai/crush/internal/ui/chat/sidebar/ (参考)
```

---

### 7.2 状态显示组件 ⭐⭐

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

### 7.3 Diff 查看器完善 ⭐⭐

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
- [ ] 创建 `internal/ui/layout/`
- [ ] Flexbox 布局
- [ ] Grid 布局
- [ ] 响应式断点
- [ ] 自适应大小

**参考**:
```
E:/projects/ai/crush/internal/ui/model/ui.go (generateLayout)
E:/projects/ai/crush/internal/ui/common/elements.go
```

---

### 9.2 侧边栏组件 ⭐⭐

**目标**: 通用侧边栏组件

**任务**:
- [ ] 创建 `internal/ui/components/sidebar/`
- [ ] 迁移 `internal/tui/components/chat/sidebar/`
- [ ] 多面板支持
- [ ] 折叠/展开
- [ ] 紧凑模式

**源文件**:
```
E:/projects/ai/crush/internal/tui/components/chat/sidebar/sidebar.go
```

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
| 8.1 | 消息渲染框架 | 2周 | ⭐⭐⭐ |

### P1 - 强烈推荐 (常用功能)

| Phase | 组件 | 预估时间 | 价值 |
|-------|------|----------|------|
| 7.2 | 状态显示组件 | 1周 | ⭐⭐ |
| 9.1 | 通用布局组件 | 1周 | ⭐⭐⭐ |
| 9.2 | 侧边栏组件 | 1周 | ⭐⭐ |
| 12.1 | 文档完善 | 持续 | ⭐⭐⭐ |

### P2 - 推荐完成 (增强功能)

| Phase | 组件 | 预估时间 | 价值 |
|-------|------|----------|------|
| 6.3 | 自动完成组件 | 1周 | ⭐⭐ |
| 7.3 | Diff 查看器完善 | 1周 | ⭐⭐ |
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

**文档版本**: v2.0.0
**创建日期**: 2025-01-29
**最后更新**: 2025-01-29

---

## 📝 更新日志

### 2025-01-29
- ✅ Phase 6.1 部分完成: `internal/ui/styles/` 已创建并迁移
  - 主题系统 (theme.go)
  - Markdown 渲染器 (markdown.go)
  - Chroma 语法高亮 (chroma.go)
  - Charmtone 颜色调色板 (palette.go, charmtone.go)
  - 图标系统 (icons.go)
- ✅ 移除旧的 `internal/tui/styles/` 包
- ✅ 所有组件已更新为使用注入式样式
