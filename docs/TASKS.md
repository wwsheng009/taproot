# Taproot TUI - 详细任务清单

## Phase 2: 应用框架 🚧

### 2.1 页面系统 (page/)

- [ ] 创建 `internal/tui/page/` 目录
- [ ] 迁移 `page.go`
  - [ ] 定义 `PageID` 类型
  - [ ] 定义 `PageChangeMsg` 消息
  - [ ] 定义 `PageCloseMsg` 消息
  - [ ] 定义 `PageBackMsg` 消息（支持返回）
- [ ] 实现页面管理器
  - [ ] 页面注册机制
  - [ ] 页面切换逻辑
  - [ ] 页面栈管理（支持前进/后退）
  - [ ] 页面生命周期（Init/Update/View）
- [ ] 编写测试
  - [ ] 测试页面切换
  - [ ] 测试页面栈
  - [ ] 测试页面生命周期
- [ ] 创建示例程序

**预估**: 2-3 小时

---

### 2.2 对话框管理器 (dialogs/)

- [ ] 创建 `internal/tui/components/dialogs/` 目录
- [ ] 迁移 `dialogs.go`
  - [ ] 定义 `DialogID` 类型
  - [ ] 定义 `DialogModel` 接口
    - [ ] `Init() tea.Cmd`
    - [ ] `Update(msg tea.Msg) (Model, tea.Cmd)`
    - [ ] `View() string`
    - [ ] `Position() (int, int)`
    - [ ] `ID() DialogID`
  - [ ] 定义 `OpenDialogMsg` 消息
  - [ ] 定义 `CloseDialogMsg` 消息
  - [ ] 定义 `DialogCmp` 接口
  - [ ] 实现对话框堆栈
    - [ ] Push 对话框
    - [ ] Pop 对话框
    - [ ] 获取活动对话框
  - [ ] 实现键盘导航
    - [ ] ESC 关闭
    - [ ] Tab 切换
  - [ ] 实现层级渲染
    - [ ] 使用 `lipgloss.Layer`
    - [ ] 处理窗口大小变化
- [ ] 创建基础对话框示例
- [ ] 编写测试
  - [ ] 测试对话框打开/关闭
  - [ ] 测试对话框堆栈
  - [ ] 测试键盘导航
- [ ] 创建示例程序

**预估**: 4-5 小时

---

### 2.3 应用主循环 (app/)

- [ ] 创建 `internal/tui/app/` 目录
- [ ] 创建 `app.go`
  - [ ] 定义 `AppModel` 结构
    - [ ] 页面管理
    - [ ] 对话框管理
    - [ ] 全局状态
    - [ ] 窗口尺寸
  - [ ] 实现初始化逻辑
  - [ ] 实现 Update 方法
    - [ ] 路由消息到页面/对话框
    - [ ] 处理全局快捷键
    - [ ] 处理窗口大小变化
  - [ ] 实现 View 方法
    - [ ] 渲染当前页面
    - [ ] 渲染对话框层
    - [ ] 渲染状态栏
  - [ ] 实现页面切换逻辑
  - [ ] 实现对话框集成
  - [ ] 实现全局快捷键
    - [ ] `ctrl+c`: 退出
    - [ ] `ctrl+g`: 切换帮助
    - [ ] `ESC`: 关闭对话框
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 3-4 小时

---

## Phase 3: 通用组件 ⏳

### 3.1 自动完成组件 (completions/)

- [ ] 创建 `internal/tui/components/completions/` 目录
- [ ] 迁移 `completions.go`
  - [ ] 定义 `CompletionItem` 结构
  - [ ] 定义 `Completions` 模型
  - [ ] 实现模糊匹配算法
  - [ ] 实现键盘导航
    - [ ] 上下箭头
    - [ ] Enter 确认
    - [ ] ESC 取消
  - [ ] 实现高亮显示
  - [ ] 实现多选支持
- [ ] 迁移 `keys.go`
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 5-6 小时

---

### 3.2 虚拟化列表 (exp/list/)

- [ ] 创建 `internal/tui/exp/list/` 目录
- [ ] 迁移 `items.go`
  - [ ] 定义 `Item` 接口
  - [ ] 定义 `StandardItem` 实现
- [ ] 迁移 `list.go`
  - [ ] 定义 `ListModel` 结构
  - [ ] 实现虚拟化渲染
  - [ ] 实现滚动逻辑
  - [ ] 实现键盘导航
  - [ ] 实现选中状态管理
- [ ] 迁移 `keys.go`
- [ ] 迁移 `filterable.go`
  - [ ] 实现过滤功能
  - [ ] 实现搜索高亮
- [ ] 迁移 `filterable_group.go`
  - [ ] 实现分组过滤
- [ ] 迁移 `grouped.go`
  - [ ] 实现分组显示
  - [ ] 实现分组折叠
- [ ] 迁移 `list_test.go`
- [ ] 迁移 `filterable_test.go`
- [ ] 创建示例程序

**预估**: 10-12 小时

---

### 3.3 Diff查看器 (exp/diffview/)

- [ ] 创建 `internal/tui/exp/diffview/` 目录
- [ ] 迁移 `diffview.go`
  - [ ] 定义 `DiffView` 结构
  - [ ] 实现 unified diff 视图
  - [ ] 实现分屏 diff 视图
  - [ ] 实现滚动同步
- [ ] 迁移 `split.go`
  - [ ] 实现分屏布局
  - [ ] 实现光标同步
- [ ] 迁移 `style.go`
  - [ ] 定义颜色样式
  - [ ] 实现行号样式
  - [ ] 实现代码样式
- [ ] 迁移 `chroma.go`
  - [ ] 集成 Chroma 语法高亮
  - [ ] 定义 Chroma 主题
- [ ] 迁移 `util.go`
  - [ ] 实现解析函数
  - [ ] 实现辅助函数
- [ ] 迁移测试文件
  - [ ] `diffview_test.go`
  - [ ] `udiff_test.go`
  - [ ] `util_test.go`
- [ ] 创建示例程序

**预估**: 8-10 小时

---

### 3.4 Logo渲染 (logo/)

- [ ] 创建 `internal/tui/components/logo/` 目录
- [ ] 迁移 `logo.go`
  - [ ] 定义 `Logo` 结构
  - [ ] 实现 ASCII logo 渲染
  - [ ] 实现颜色渐变
- [ ] 迁移 `rand.go`
  - [ ] 实现随机 logo 生成
- [ ] 创建示例程序

**预估**: 2 小时

---

### 3.5 文件列表 (files/)

- [ ] 创建 `internal/tui/components/files/` 目录
- [ ] 迁移 `files.go`
  - [ ] 定义 `FileItem` 结构
  - [ ] 定义 `FilesModel` 结构
  - [ ] 实现目录遍历
  - [ ] 实现文件图标
  - [ ] 实现排序功能
  - [ ] 实现过滤功能
  - [ ] 实现隐藏文件显示
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 3-4 小时

---

### 3.6 语法高亮 (highlight/)

- [ ] 创建 `internal/tui/components/highlight/` 目录
- [ ] 迁移 `highlight.go`
  - [ ] 定义 `Highlighter` 接口
  - [ ] 实现 Chroma 集成
  - [ ] 实现主题映射
  - [ ] 实现语言检测
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 2-3 小时

---

## Phase 4: 对话框系统 ⏳

### 4.1 文件选择器 (dialogs/filepicker/)

- [ ] 创建 `internal/tui/components/dialogs/filepicker/` 目录
- [ ] 迁移 `filepicker.go`
  - [ ] 定义 `FilePicker` 结构
  - [ ] 实现目录浏览
  - [ ] 实现文件过滤
  - [ ] 实现键盘导航
  - [ ] 实现对话框接口
- [ ] 迁移 `keys.go`
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 5-6 小时

---

### 4.2 退出确认 (dialogs/quit/)

- [ ] 创建 `internal/tui/components/dialogs/quit/` 目录
- [ ] 迁移 `quit.go`
  - [ ] 定义 `QuitDialog` 结构
  - [ ] 实现 "未保存更改" 检查
  - [ ] 实现确认逻辑
- [ ] 迁移 `keys.go`
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 2 小时

---

### 4.3 推理显示 (dialogs/reasoning/)

- [ ] 创建 `internal/tui/components/dialogs/reasoning/` 目录
- [ ] 迁移 `reasoning.go`
  - [ ] 定义 `ReasoningDialog` 结构
  - [ ] 实现可折叠内容
  - [ ] 实现 Markdown 渲染
  - [ ] 实现流式更新
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 2-3 小时

---

### 4.4 基础命令面板 (dialogs/commands/)

- [ ] 创建 `internal/tui/components/dialogs/commands/` 目录
- [ ] 定义接口
  - [ ] `CommandProvider` 接口
  - [ ] `Command` 结构
  - [ ] `ArgDef` 结构
- [ ] 迁移 `commands.go`
  - [ ] 实现命令列表显示
  - [ ] 实现模糊搜索
  - [ ] 实现参数输入
  - [ ] 解耦执行逻辑（使用回调）
- [ ] 迁移 `arguments.go`
  - [ ] 实现参数输入界面
- [ ] 迁移 `keys.go`
- [ ] 编写测试
- [ ] 创建示例程序（带回调）

**预估**: 7-8 小时

---

### 4.5 基础模型选择 (dialogs/models/)

- [ ] 创建 `internal/tui/components/dialogs/models/` 目录
- [ ] 定义接口
  - [ ] `ModelProvider` 接口
  - [ ] `ConfigProvider` 接口
  - [ ] `Model` 结构
- [ ] 迁移 `models.go`
  - [ ] 实现模型列表显示
  - [ ] 实现搜索过滤
  - [ ] 实现最近使用
  - [ ] 解耦业务逻辑
- [ ] 迁移 `list.go`
  - [ ] 实现模型列表组件
- [ ] 迁移 `apikey.go`
  - [ ] 实现 API 密钥输入
- [ ] 迁移 `keys.go`
- [ ] 编写测试
- [ ] 创建示例程序（带 mock Provider）

**预估**: 5-6 小时

---

### 4.6 基础会话切换 (dialogs/sessions/)

- [ ] 创建 `internal/tui/components/dialogs/sessions/` 目录
- [ ] 定义接口
  - [ ] `SessionProvider` 接口
  - [ ] `Session` 结构
- [ ] 迁移 `sessions.go`
  - [ ] 实现会话列表显示
  - [ ] 实现搜索功能
  - [ ] 实现新建会话
  - [ ] 实现删除会话
  - [ ] 解耦业务逻辑
- [ ] 迁移 `keys.go`
- [ ] 编写测试
- [ ] 创建示例程序（带 mock Provider）

**预估**: 5-6 小时

---

## Phase 5: 高级组件 ⏳

### 5.1 图片渲染 (image/)

- [ ] 创建 `internal/tui/components/image/` 目录
- [ ] 迁移 `image.go`
  - [ ] 定义 `Image` 结构
  - [ ] 实现 kitty 协议
  - [ ] 实现 iterm2 协议
  - [ ] 实现自动检测
- [ ] 迁移 `load.go`
  - [ ] 实现图片加载
  - [ ] 实现缓存
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 7-8 小时

---

### 5.2 消息渲染 (messages/)

- [ ] 创建 `internal/tui/components/messages/` 目录
- [ ] 迁移 `messages.go`
  - [ ] 定义 `MessagesModel` 结构
  - [ ] 实现消息列表
  - [ ] 实现滚动逻辑
- [ ] 迁移 `renderer.go`
  - [ ] 集成 Glamour
  - [ ] 实现 Markdown 渲染
  - [ ] 实现代码块渲染
- [ ] 迁移 `tool.go`
  - [ ] 实现工具调用显示
  - [ ] 实现工具结果显示
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 9-10 小时

---

### 5.3 文本编辑器 (editor/)

**注意**: 此组件极其复杂，建议作为独立项目

- [ ] 创建 `internal/tui/components/editor/` 目录
- [ ] 迁移 `editor.go`
  - [ ] 定义 `Editor` 结构
  - [ ] 实现多行编辑
  - [ ] 实现语法高亮
  - [ ] 实现自动缩进
- [ ] 迁移 `clipboard.go`
  - [ ] 实现剪贴板接口
- [ ] 迁移 `clipboard_supported.go`
  - [ ] 实现各平台支持
- [ ] 迁移 `clipboard_not_supported.go`
  - [ ] 实现无剪贴板支持
- [ ] 迁移 `keys.go`
- [ ] 编写测试
- [ ] 创建示例程序

**预估**: 18-20 小时（作为独立项目）

---

## 文档任务 📚

- [x] 创建 `docs/ARCHITECTURE.md`
- [x] 创建 `docs/MIGRATION_PLAN.md`
- [x] 创建 `docs/ALTERNATIVES.md`
- [x] 创建 `docs/TASKS.md` (本文档)
- [ ] 创建 `docs/API.md` - API 文档
- [ ] 创建 `docs/EXAMPLES.md` - 示例集合
- [ ] 创建 `docs/CONTRIBUTING.md` - 贡献指南
- [ ] 创建 `docs/CHANGELOG.md` - 变更日志
- [ ] 更新 `AGENTS.md`
- [ ] 创建 README.md

---

## 测试任务 🧪

- [ ] 为每个组件添加单元测试
- [ ] 为关键组件添加集成测试
- [ ] 添加性能基准测试
- [ ] 添加跨平台测试
- [ ] 设置 CI/CD

---

## 发布任务 🚀

- [ ] 设置版本号策略
- [ ] 创建 GitHub Releases
- [ ] 发布到 GoPackages
- [ ] 创建示例应用
- [ ] 写博客介绍

---

## 优先级说明

- **P0**: 必须完成，框架核心功能
- **P1**: 强烈推荐，常用功能
- **P2**: 推荐完成，增强功能
- **P3**: 可选，特殊场景

---

## 当前进度

```
Phase 1: ████████████████████ 100% (已完成)
Phase 2: ░░░░░░░░░░░░░░░░░░░░   0%
Phase 3: ░░░░░░░░░░░░░░░░░░░░   0%
Phase 4: ░░░░░░░░░░░░░░░░░░░░   0%
Phase 5: ░░░░░░░░░░░░░░░░░░░░   0%
文档:    ████████████░░░░░░░░  50%
```

---

## 下一步行动 (按顺序)

1. ⏳ **Phase 2.1**: 实现页面系统 (2-3h)
2. ⏳ **Phase 2.2**: 实现对话框管理器 (4-5h)
3. ⏳ **Phase 2.3**: 实现应用主循环 (3-4h)
4. ⏳ **Phase 3.2**: 实现虚拟化列表 (10-12h) - 高价值
5. ⏳ **Phase 3.3**: 实现 Diff 查看器 (8-10h) - 高价值

**完成这些后，Taproot 将拥有完整的 TUI 框架基础！**
