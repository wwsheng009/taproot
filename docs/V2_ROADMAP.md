# Taproot v2.0.0 - 迁移计划

基于对 Crush 项目的详细审查，以下是未迁移但有价值的组件清单。

## Phase 6: 基础设施完善 (P0)

### 6.1 全局按键系统
**估计时间**: 1-2 小时

- [ ] 迁移 `internal/tui/keys.go`
  - [x] 定义 KeyMap 结构
  - [ ] 实现默认按键绑定
  - [ ] 添加按键帮助文本
  - [ ] 编写测试
- [ ] 集成到应用框架

**价值**: 统一按键管理，所有组件共享

### 6.2 通用布局组件
**估计时间**: 2-3 小时

- [ ] 迁移 `internal/tui/components/core/layout/layout.go`
  - [ ] 定义布局接口和类型
  - [ ] 实现布局管理器
  - [ ] 支持动态布局调整
  - [ ] 编写测试
- [ ] 创建示例程序

**价值**: 可复用的布局系统，支持复杂界面

### 6.3 对话框按键绑定
**估计时间**: 1-2 小时

- [ ] 迁移各对话框的 keys.go
  - [ ] `dialogs/commands/keys.go`
  - [ ] `dialogs/models/keys.go`
  - [ ] `dialogs/sessions/keys.go`
  - [ ] `dialogs/filepicker/keys.go`
  - [ ] `dialogs/quit/keys.go`
  - [ ] `dialogs/reasoning/keys.go`
- [ ] 统一按键模式
- [ ] 添加测试

**价值**: 完整的按键系统，用户体验一致性

## Phase 7: 列表组件完善 (P1)

### 7.1 列表键盘导航
**估计时间**: 1-2 小时

- [ ] 迁移 `exp/list/keys.go`
  - [ ] 定义列表按键绑定
  - [ ] 实现导航快捷键
  - [ ] 集成到列表组件
  - [ ] 编写测试

### 7.2 可过滤分组列表
**估计时间**: 2-3 小时

- [ ] 迁移 `exp/list/filterable_group.go`
  - [ ] 实现过滤+分组组合功能
  - [ ] 性能优化
  - [ ] 编写测试
  - [ ] 创建示例程序

### 7.3 列表测试套件
**估计时间**: 2-3 小时

- [ ] 迁移 `exp/list/list_test.go`
- [ ] 迁移 `exp/list/filterable_test.go`
- [ ] 添加边界测试
- [ ] 添加性能测试

**价值**: 保证组件质量和稳定性

## Phase 8: Diff查看器完善 (P1)

### 8.1 Diff 工具函数
**估计时间**: 1-2 小时

- [ ] 迁移 `exp/diffview/util.go`
  - [ ] 实现 unified diff 解析
  - [ ] 实现 git diff 解析
  - [ ] 错误处理
  - [ ] 编写测试

### 8.2 Diff 样式系统
**估计时间**: 1-2 小时

- [ ] 迁移 `exp/diffview/style.go`
  - [ ] 定义 diff 颜色主题
  - [ ] 集成主题系统
  - [ ] 编写测试

### 8.3 分屏 Diff 视图
**估计时间**: 3-4 小时

- [ ] 迁移 `exp/diffview/split.go`
  - [ ] 实现左右分屏布局
  - [ ] 同步滚动
  - [ ] 光标同步
  - [ ] 编写测试
  - [ ] 创建示例程序

### 8.4 Diff 测试套件
**估计时间**: 2-3 小时

- [ ] 迁移 `exp/diffview/*_test.go`
- [ ] 添加测试数据
- [ ] 边界测试

## Phase 9: 高级消息组件 (P2)

### 9.1 任务列表组件
**估计时间**: 3-4 小时

- [ ] 从 `components/chat/todos/todos.go` 提取
- [ ] 解耦 session.Todo 依赖
- [ ] 泛化为通用任务列表
- [ ] 添加 CRUD 操作
- [ ] 编写测试
- [ ] 创建示例程序

**价值**: 通用任务管理组件

### 9.2 消息渲染器
**估计时间**: 2-3 小时

- [ ] 迁移 `components/chat/messages/renderer.go`
  - [ ] 解耦 message.Message 类型
  - [ ] 支持多种消息格式
  - [ ] 集成 Markdown 渲染
  - [ ] 编写测试

### 9.3 工具调用显示
**估计时间**: 2-3 小时

- [ ] 迁移 `components/chat/messages/tool.go`
  - [ ] 解耦 tool 类型
  - [ ] 支持展开/折叠
  - [ ] 支持异步状态更新
  - [ ] 编写测试
  - [ ] 创建示例程序

## Phase 10: Shell 执行工具 (P2)

### 10.1 Shell 命令执行
**估计时间**: 1-2 小时

- [ ] 重写 `util/shell.go`
  - [ ] 移除 uiutil 依赖
  - [ ] 直接使用 mvdan.cc/sh/v3
  - [ ] 实现 ExecShell 函数
  - [ ] 编写测试

**价值**: 独立的 shell 执行能力

## Phase 11: 示例程序补充 (P2)

### 11.1 缺失示例
**估计时间**: 4-6 小时

- [ ] `examples/keys/` - 按键绑定演示
- [ ] `examples/layout/` - 布局组件演示
- [ ] `examples/todos/` - 任务列表演示
- [ ] `examples/diffsplit/` - 分屏 diff 演示
- [ ] `examples/filterablegroup/` - 过滤分组列表演示

## Phase 12: 质量保证 (P1)

### 12.1 测试覆盖提升
**估计时间**: 6-8 小时

- [ ] 为所有核心组件添加测试
- [ ] 目标覆盖率 >70%
- [ ] 添加集成测试
- [ ] 添加基准测试

### 12.2 文档完善
**估计时间**: 2-3 小时

- [ ] 更新 AGENTS.md
- [ ] 添加 API 使用示例
- [ ] 添加故障排查指南
- [ ] 更新 README 特性列表

## 优先级路线图

### 快速胜利 (1-2 天)
```
Week 1:
- Day 1-2: Phase 6.1 (keys.go)
- Day 3-4: Phase 6.2 (layout.go)
- Day 5: Phase 6.3 (对话框按键)
```

### 核心功能 (1 周)
```
Week 2:
- Day 1-2: Phase 7.1-7.2 (列表完善)
- Day 3-4: Phase 8.1-8.3 (Diff 完善)
- Day 5: Phase 10.1 (Shell 执行)
```

### 高级特性 (1 周)
```
Week 3:
- Day 1-2: Phase 9.1-9.3 (消息组件)
- Day 3-4: Phase 11.1 (示例补充)
- Day 5: Phase 12.1-12.2 (质量保证)
```

## 总体时间估算

| Phase | 内容 | 时间 |
|-------|------|------|
| Phase 6 | 基础设施 | 4-7 小时 |
| Phase 7 | 列表完善 | 5-8 小时 |
| Phase 8 | Diff 完善 | 7-11 小时 |
| Phase 9 | 高级消息 | 7-10 小时 |
| Phase 10 | Shell 工具 | 1-2 小时 |
| Phase 11 | 示例补充 | 4-6 小时 |
| Phase 12 | 质量保证 | 8-11 小时 |
| **总计** | **36-55 小时** (约 1-2 周全职工作) |

## 里程碑

- **v1.1.0**: Phase 6 完成（基础设施）
- **v1.2.0**: Phase 7-8 完成（列表和 Diff）
- **v1.3.0**: Phase 9-10 完成（高级功能）
- **v2.0.0**: Phase 11-12 完成（完整框架）

## 开始顺序建议

1. ✅ **keys.go** - 最简单，无依赖
2. ✅ **layout.go** - 独立组件
3. ✅ **对话框 keys.go** - 批量迁移
4. ✅ **list keys.go** - 列表增强
5. ✅ **diffview util/style** - 工具函数

要开始从 Phase 6.1 (keys.go) 开始吗？
