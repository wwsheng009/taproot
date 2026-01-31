# Contributing to Taproot

感谢你对 Taproot 的兴趣！我们欢迎所有形式的贡献。

## 🚀 快速开始

### 开发环境设置

1. **Fork 并克隆仓库**
   ```bash
   git clone https://github.com/YOUR_USERNAME/taproot.git
   cd taproot
   git remote add upstream https://github.com/wwsheng009/taproot.git
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **运行测试**
   ```bash
   go test ./...
   go build ./...
   ```

## 📋 开发流程

### 1. 分支策略

- `main` - 主分支，稳定版本
- `develop` - 开发分支（如果需要）
- `feature/*` - 功能分支
- `fix/*` - 修复分支

### 2. 提交规范

使用语义化提交消息：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type):**
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

**示例：**
```
feat(dialog): 添加文件选择对话框

实现基于文件系统的文件浏览功能，支持：
- 目录导航
- 文件过滤
- 键盘快捷键

Closes #123
```

### 3. 代码规范

#### Go 代码风格
- 遵循 [Effective Go](https://golang.org/doc/effective_go) 指南
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量

#### 命名约定
- **包名**: 小写单数 (`layout`, `styles`)
- **接口**: `-able` 后缀 (`Focusable`, `Sizeable`)
- **常量**: 大驼峰或全大写 (`MaxWidth` 或 `MAX_WIDTH`)
- **私有**: 小驼峰 (`internalState`)

#### 组件设计原则
- 实现相关接口 (`Focusable`, `Sizeable` 等)
- 提供 `New()` 构造函数
- 实现 `Init()`, `Update()`, `View()` 方法
- 使用 `styles.DefaultStyles()` 获取主题

### 4. 测试要求

所有新功能必须包含测试：

```go
func TestNewComponent(t *testing.T) {
    comp := NewComponent()
    
    if comp == nil {
        t.Error("NewComponent() should return non-nil")
    }
    
    // 测试默认值
    if comp.Width() != 80 {
        t.Errorf("Expected default width 80, got %d", comp.Width())
    }
}
```

运行测试：
```bash
go test ./...
go test -v ./path/to/package
go test -cover ./...
```

### 5. 文档要求

更新相关文档：
- **新组件**: 更新 README.md 的组件列表
- **破坏性变更**: 更新 CHANGELOG.md
- **新功能**: 添加示例到 `examples/`

## 📦 Pull Request 流程

### 1. 创建 PR 前

- [ ] 确保所有测试通过
- [ ] 添加必要的测试
- [ ] 更新文档
- [ ] 代码已格式化 (`gofmt`)
- [ ] 提交消息符合规范

### 2. 提交 PR

```bash
git checkout -b feature/amazing-feature
# 进行开发...
git add .
git commit -m "feat(component): 添加 amazing 功能"
git push origin feature/amazing-feature
```

然后在 GitHub 上创建 Pull Request。

### 3. PR 标题格式

```
<type>: <subject>

类型：
- feat: 新功能
- fix: 修复
- docs: 文档
- refactor: 重构
- test: 测试
- chore: 构建
```

### 4. PR 描述模板

```markdown
## 变更说明
简要描述这个 PR 做了什么。

## 变更类型
- [ ] 新功能
- [ ] Bug 修复
- [ ] 破坏性变更
- [ ] 文档更新

## 测试
描述如何测试这些变更。

## 截图
如果适用，添加截图或 GIF。

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 自我审查代码
- [ ] 添加了必要的测试
- [ ] 所有测试通过
- [ ] 更新了文档
```

## 🐛 报告 Bug

### Bug 报告模板

```markdown
**问题描述**
清晰简洁地描述 bug 是什么。

**复现步骤**
1. 运行 '...'
2. 点击 '....'
3. 滚动到 '....'
4. 看到错误

**预期行为**
描述你期望发生什么。

**实际行为**
描述实际发生了什么。

**环境信息**
- OS: [e.g. Windows 11, macOS 14.0]
- Go 版本: [e.g. 1.24.2]
- Taproot 版本: [e.g. v1.0.0]

**截图**
如果适用，添加截图帮助说明问题。

**附加信息**
其他相关信息、日志、错误消息等。
```

## 💡 功能建议

### 功能建议模板

```markdown
**功能描述**
清晰简洁地描述你希望添加的功能。

**问题背景**
这个功能解决什么问题？为什么需要它？

**建议的解决方案**
详细描述你希望如何实现这个功能。

**替代方案**
描述你考虑过的其他替代方案。

**附加信息**
其他相关信息、示例代码等。
```

## 🎨 开发工具

### 必需工具
- Go 1.24+
- Git

### 推荐工具
- [gofmt](https://golang.org/pkg/cmd/gofmt/) - 代码格式化
- [golint](https://github.com/golang/lint) - 代码检查
- [goimports](https://github.com/golang/tools/tree/master/cmd/goimports) - 导入管理

### 调试技巧
```go
// 使用 t.Log() 输出调试信息
func TestExample(t *testing.T) {
    t.Log("Debug info:", someValue)
}

// 使用 tea.Model 的 String() 方法
func (m Model) View() string {
    // 返回可读的状态信息
}
```

## 📝 代码审查准则

### 我们会关注
- **代码质量**: 清晰、可读、可维护
- **测试覆盖**: 新功能必须有测试
- **文档**: 公共 API 必须有文档
- **一致性**: 遵循项目现有风格
- **性能**: 避免不必要的内存分配

### 不会接受
- 破坏现有 API 的变更（除非必要且经过讨论）
- 没有测试的新功能
- 不符合项目规范的代码
- 仅格式化的 PR

## 🌟 社区规范

### 行为准则
- 尊重不同观点
- 建设性反馈
- 关注问题而非个人

### 沟通渠道
- **GitHub Issues**: Bug 报告、功能建议
- **Pull Requests**: 代码贡献
- **Discussions**: 一般讨论

## 📜 许可证

通过贡献代码，你同意你的贡献将在 [MIT License](LICENSE) 下发布。

## 🙏 致谢

感谢所有贡献者！你的贡献让 Taproot 变得更好。

---

有任何问题？欢迎创建 Issue 或 Discussion！
