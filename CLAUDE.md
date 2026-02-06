# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Taproot** is a composable TUI (Terminal User Interface) framework for Go, built on top of Bubbletea. It was extracted from the Crush CLI project and provides reusable interfaces, components, and utilities for building terminal applications.

- **Module**: `github.com/wwsheng009/taproot`
- **Go Version**: 1.24.2
- **Core Dependencies**: Bubbletea v1.3.10, Lipgloss v1.1.1

## Build and Test Commands

```bash
# Run examples
go run examples/demo/main.go
go run examples/commands/main.go
go run examples/models/main.go
go run examples/sessions/main.go
go run examples/messages/main.go
go run examples/app/main.go

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./layout/
go test ./tui/util/

# Build project
go build ./...
```

## Architecture

Taproot follows the Elm Architecture (Model-View-Update) used by Bubbletea:

```
AppModel
├── Page Management    (current page, page switching)
├── Dialog Stack       (commands, models, sessions, file picker, etc.)
└── Status Bar         (InfoMsg with TTL)
```

### Package Structure

| Package | Purpose |
|---------|---------|
| `layout/` | Core interfaces: `Focusable`, `Sizeable`, `Positional`, `Help` |
| `ui/styles/` | Theme system with HCL color blending and gradients |
| `ui/list/` | Engine-agnostic virtualized list components (v2.0) |
| `ui/render/` | Rendering engine abstraction (Bubbletea/Ultraviolet support) |
| `ui/components/` | UI components (dialogs, forms, messages, status, etc.) |
| `tui/app/` | Application framework with page/dialog management |
| `tui/page/` | Page management system |
| `tui/util/` | Utilities (Model interface, InfoMsg, shell execution) |
| `tui/components/` | High-level framework components |
| `tui/exp/` | Experimental features (diff viewer, advanced lists) |

## Key Interfaces

Components implement composable interfaces from `layout/`:

- **Focusable**: `Focus()`, `Blur()`, `Focused() bool`
- **Sizeable**: `Size() (width, height int)`, `SetSize(width, height int)`
- **Positional**: `Position() (x, y int)`, `SetPosition(x, y int)`
- **Help**: `Help() []string`

## Bubbletea Model Pattern

All components implement the Model interface:

```go
type Model interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Model, tea.Cmd)
    View() string
}
```

**Key patterns**:
- Use type switch for message handling: `switch msg := msg.(type) {`
- View is called every frame - optimize with `strings.Builder`
- Return `tea.Quit` to exit (v1.0) or `render.Quit()` for v2.0 engine-agnostic components

## Theme System

Always use `styles.DefaultStyles()` for colors:

```go
s := styles.DefaultStyles()
text := s.Base.Foreground(s.Primary).Render("Hello")
gradient := styles.ApplyForegroundGrad(&s, "Gradient Text", s.Primary, s.Secondary)
```

Color categories: `Primary`, `Secondary`, `BgBase`, `FgMuted`, `Success`, `Error`, `Warning`, `Info`

## Status Reporting

Use `util.InfoMsg` for status messages:

```go
return m, util.ReportError(err)
return m, util.ReportInfo("Operation complete")
return m, util.ReportWarn("Resource low")
```

Message types: `Info`, `Success`, `Warn`, `Error`, `Update`

## Key Binding Conventions

| Action | Keys |
|--------|------|
| Quit | `ctrl+c`, `q` |
| Navigation | Arrow keys, `j`/`k` (vim-style) |
| Selection | `space`, `enter` |
| Help | `ctrl+g` |
| Commands | `ctrl+p` |

## Code Style

- Package names: lowercase (`layout`, `util`, `styles`)
- Interfaces: `-able` suffix (`Focusable`, `Sizeable`)
- Functions: PascalCase (exported), camelCase (internal)
- Constructors: `NewComponent()` or `initialModel()`
- Test files: `_test.go` suffix

## v2.0 Engine-Agnostic Components

The framework supports dual rendering engines (Bubbletea + Ultraviolet) via `ui/render/`:

```go
import "github.com/wwsheng009/taproot/ui/render"

// Components implement render.Model
func (c *Component) Update(msg any) (render.Model, render.Cmd)

// Use render.Quit() instead of tea.Quit()
return m, render.Quit()

// String-based key handling
case render.KeyMsg:
    switch key.Key {
    case "up", "k": ...
    }
```

## Dialog System

Dialogs use a stack-based overlay manager:

```go
overlay := dialog.NewOverlay()
overlay.Push(confirmDialog)
overlay.Pop()

// Dialog types: InfoDialog, ConfirmDialog, InputDialog, SelectListDialog
```

## Common Gotchas

1. **Model updates**: Always return the modified model, even if no changes
2. **Commands can be nil**: Not every update needs to return a command
3. **View optimization**: Pre-render static content, cache where appropriate
4. **Platform differences**: Windows has some terminal limitations
5. **Shell execution**: Use `util.ExecShell()` for proper TTY handling with `mvdan.cc/sh/v3`

## Testing

- Use table-driven tests for multiple scenarios
- Test interface compliance with type assertions
- Mock components should implement all relevant interfaces
- For v2.0 components, use `DirectEngine` for unit tests (no terminal required)

## Dependencies

Core dependencies are properly pinned:
- `charmbracelet/bubbletea` v1.3.10 - Elm architecture
- `charmbracelet/bubbles` v0.21.0 - Interactive components
- `charmbracelet/lipgloss` v1.1.1 - Styling
- `charmbracelet/glamour` v0.10.0 - Markdown rendering
- `alecthomas/chroma/v2` v2.23.1 - Syntax highlighting
- `mvdan.cc/sh/v3` v3.12.0 - Shell parsing
