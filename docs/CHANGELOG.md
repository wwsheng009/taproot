# Changelog

All notable changes to Taproot TUI Framework will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-01-28

### Added

#### Phase 1: Framework Foundation
- Core layout interfaces (Focusable, Sizeable, Positional, Help)
- TUI utilities (Model, InfoMsg, ExecShell)
- Global key bindings system (KeyMap, DefaultKeyMap)
- Theme system with HCL color space blending
- Animated spinner component with gradients
- Core UI components (Section, Title, Status, Button)
- Status bar component with TTL-based message clearing

#### Phase 2: Application Framework
- Page management system with page stack (navigation support)
- Dialog manager with stack-based rendering
- Application main loop (AppModel)
- Global shortcuts (ctrl+c, ctrl+g, ESC)

#### Phase 3: Common Components
- Logo rendering with ASCII art and gradients
- Auto-complete component with fuzzy matching
- Virtualized list with scrolling and filtering
- Grouped list with expand/collapse
- Diff viewer with unified diff view
- Syntax highlighting with Chroma integration
- **Markdown rendering with Glamour**
- **Charmtone color palette (24 predefined colors)**

#### Phase 4: Dialog System
- Command Palette dialog (ctrl+p)
- Model Selection dialog (ctrl+m)
- File Picker dialog
- Quit Confirmation dialog
- Reasoning Display dialog
- Session Switcher dialog

#### Phase 5: Advanced Components
- Messages component with tool call display
- Image rendering with kitty/iterm2 protocol support
- **Icons definition set**

### Dependencies

- github.com/charmbracelet/bubbletea v1.3.10
- github.com/charmbracelet/bubbles v0.21.0
- github.com/charmbracelet/lipgloss v1.1.1
- github.com/charmbracelet/glamour v0.10.0
- github.com/alecthomas/chroma/v2 v2.23.1
- github.com/charmbracelet/x/ansi v0.11.4
- github.com/lucasb-eyer/go-colorful v1.3.0
- mvdan.cc/sh/v3 v3.12.0

### Documentation

- ARCHITECTURE.md - Complete framework architecture analysis
- MIGRATION_PLAN.md - 5-phase migration roadmap
- ALTERNATIVES.md - Technology selection analysis
- TASKS.md - Detailed task checklist
- AGENTS.md - Agent development guide
- README.md - Project overview and quick start
- API.md - API documentation

### Testing

- Unit tests for layout, page, dialogs, util packages
- Tests for styles (43.9% coverage)
- Tests for highlight (87.0% coverage)
- Tests for util (55.0% coverage)

## [0.9.0] - 2025-01-28

### Added

- Complete Phase 1-4 framework implementation
- All dialog components
- Page and dialog management systems
- Theme system with gradient support

## [0.1.0] - 2025-01-28

### Added

- Initial Taproot TUI framework extraction
- Core interfaces and utilities
- Basic theme system

---

## Version Policy

- **Major version (X.0.0)**: Breaking changes, major feature additions
- **Minor version (0.X.0)**: New features, backward compatible
- **Patch version (0.0.X)**: Bug fixes, minor improvements
