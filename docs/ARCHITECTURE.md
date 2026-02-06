# Taproot TUI Framework - Architecture

## Overview

Taproot is a production-ready TUI (Terminal User Interface) framework for Go, built on top of Bubbletea with a v2.0 engine-agnostic architecture. It provides comprehensive components and utilities for building feature-rich terminal applications.

**Module**: `github.com/wwsheng009/taproot`
**Go Version**: 1.24.2
**Current Version**: 2.0.0

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         Application Layer                        │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                    AppModel (tui/app/)                      │ │
│  │  - Page Management (navigation, history)                   │ │
│  │  - Dialog Stack (overlay, focus management)                │ │
│  │  - Status Bar (InfoMsg with TTL)                           │ │
│  │  - Global Key Bindings (Ctrl+C, G, P, M, Esc)             │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                       v2.0 Engine-Agnostic Layer                 │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                  render.Model Interface                     │ │
│  │  - Init() Cmd                                              │ │
│  │  - Update(msg any) (Model, Cmd)                            │ │
│  │  - View() string                                           │ │
│  └────────────────────────────────────────────────────────────┘ │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐ │
│  │   Lists      │  │   Dialogs    │  │      Forms           │ │
│  │ (ui/list/)   │  │ (ui/dialog/) │  │   (ui/forms/)        │ │
│  └──────────────┘  └──────────────┘  └──────────────────────┘ │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐ │
│  │  Messages    │  │    Status    │  │     Progress         │ │
│  │(ui/messages/) │  │(ui/status/)  │  │  (ui/progress/)      │ │
│  └──────────────┘  └──────────────┘  └──────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                          Engine Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐ │
│  │  Bubbletea   │  │  Ultraviolet  │  │      Direct          │ │
│  │  (production)│  │   (beta)     │  │   (testing)          │ │
│  └──────────────┘  └──────────────┘  └──────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Package Structure

```
taproot/
├── layout/              # Core interfaces
│   ├── layout.go        # Focusable, Sizeable, Positional, Help
│   └── layout_test.go
│
├── ui/                  # UI components (v2.0 engine-agnostic)
│   ├── styles/          # Theme system with HCL gradients
│   ├── list/            # Virtualized lists (filtering, grouping)
│   ├── render/          # Engine abstraction layer
│   ├── dialog/          # Dialog system
│   ├── forms/           # Form components
│   ├── components/      # UI components (messages, status, progress)
│   │   ├── messages/    # Chat messages, diagnostics, todos
│   │   ├── status/       # Service status (LSP, MCP)
│   │   ├── progress/    # Progress bars, spinners
│   │   ├── attachments/  # File attachments
│   │   ├── pills/       # Status pills
│   │   └── sidebar/     # Navigation sidebar
│   ├── layout/          # Layout utilities
│   └── tools/           # Tools (clipboard, shell, watcher)
│       ├── clipboard/
│       ├── shell/
│       └── watcher/
│
├── tui/                 # Framework components (Bubbletea-based)
│   ├── app/             # Application framework
│   ├── page/            # Page management
│   ├── util/            # Utilities (Model, InfoMsg)
│   ├── keys.go          # Global key bindings
│   └── components/      # High-level components
│
└── examples/            # 50+ example programs
```

## Core Framework (layout/)

The foundation of Taproot is a set of composable interfaces that all UI components implement.

### Interfaces

| Interface | Purpose | Methods |
|-----------|---------|---------|
| `Focusable` | Keyboard focus management | `Focus()`, `Blur()`, `Focused()` |
| `Sizeable` | Dimension management | `Size()`, `SetSize(w, h)` |
| `Positional` | Position management | `Position()`, `SetPosition(x, y)` |
| `Help` | Help text provider | `Help() []string` |

### Benefits

- **Composability**: Mix and match components through interfaces
- **Type Safety**: Compile-time guarantees
- **Flexibility**: Components work independently

## Application Framework (tui/app/)

The AppModel provides a complete application framework with:

### Features

- **Page Management**: Navigate between different screens
- **Dialog Stack**: Overlay multiple dialogs
- **Status Bar**: Display info messages with auto-clear
- **Key Bindings**: Global shortcuts (Ctrl+C, G, P, M, Esc)

### Global Key Bindings

| Key | Action |
|-----|--------|
| `Ctrl+C` | Quit application |
| `Ctrl+G` | Toggle help |
| `Ctrl+M` | Open model selector |
| `Ctrl+P` | Open command palette |
| `Ctrl+S` | Open session selector |
| `Esc` | Close dialog / Go back |

## Theme System (ui/styles/)

### Features

- **HCL Color Space**: Smooth color blending
- **Gradient Text**: Foreground/background gradients
- **Semantic Colors**: 20+ predefined color categories
- **Markdown Styles**: Glamour-based markdown
- **Syntax Highlighting**: Chroma theme integration

### Color Categories

```go
type Styles struct {
    // Brand
    Primary, Secondary, Tertiary, Accent lipgloss.Color

    // Backgrounds
    BgBase, BgBaseLighter, BgSubtle, BgOverlay lipgloss.Color

    // Foregrounds
    FgBase, FgMuted, FgHalfMuted, FgSubtle, FgSelected lipgloss.Color

    // Status
    Success, Error, Warning, Info lipgloss.Color

    // ... and more
}
```

### Usage

```go
s := styles.DefaultStyles()

// Use style presets
text := s.S().Base.Foreground(s.Primary).Render("Hello")

// Apply gradients
gradient := styles.ApplyForegroundGrad(&s, "Gradient", s.Primary, s.Secondary)
```

## v2.0 Engine-Agnostic Architecture

### Render Abstraction (ui/render/)

The v2.0 architecture introduces an engine abstraction layer:

```go
// Core interface - same across all engines
type Model interface {
    Init() Cmd
    Update(msg any) (Model, Cmd)
    View() string
}

// Command abstraction
type Cmd interface {
    Execute() error
}

// Built-in commands
func None() Cmd
func Quit() Cmd
func Batch(cmds ...Cmd) Cmd
func Tick(interval time.Duration, fn func(Msg)) Cmd
```

### Engine Types

| Engine | Type | Status | Use Case |
|--------|------|--------|----------|
| Bubbletea | `EngineBubbletea` | Production | Complex TUIs |
| Ultraviolet | `EngineUltraviolet` | Beta | High-performance UIs |
| Direct | `EngineDirect` | Testing | Unit tests, CI/CD |

### Message Types

```go
type KeyMsg struct {
    Key      string  // "up", "down", "enter", "q", etc.
    Runes    []rune
    Alt      bool
    Ctrl     bool
}

type WindowSizeMsg struct {
    Width  int
    Height int
}

type TickMsg struct {
    Time time.Time
}

type QuitMsg struct{}
```

## Component Systems

### Lists (ui/list/)

Virtualized list components with:

- **Viewport**: Windowed scrolling for large datasets
- **Filter**: Real-time search with match highlighting
- **Selection**: Single/multiple selection modes
- **Groups**: Expandable grouped items
- **Key Bindings**: Configurable keyboard navigation

### Dialogs (ui/dialog/)

Modal dialog system with:

- **Dialog Types**: Info, Confirm, Input, Select
- **Overlay Manager**: Stack-based dialog management
- **Callbacks**: Result handling via callbacks
- **Auto-centering**: Automatic positioning

### Forms (ui/forms/)

Form input components:

- **TextInput**: Single-line with validation
- **TextArea**: Multi-line with word wrap
- **Select**: Dropdown selection
- **Checkbox**: Boolean toggle
- **Radio**: Single selection from options
- **Form Container**: Focus traversal and validation

### Messages (ui/components/messages/)

Chat message display:

- **UserMessage**: User input with code blocks
- **AssistantMessage**: AI responses with markdown
- **ToolMessage**: Tool calls and results
- **FetchMessage**: Network requests (4 types)
- **DiagnosticMessage**: LSP diagnostics
- **TodoMessage**: Task lists with progress

### Status (ui/components/status/)

Service status display:

- **ServiceCmp**: Single service status
- **DiagnosticStatusCmp**: Diagnostic summary
- **LSPList**: Multiple LSP services
- **MCPList**: MCP services with tool counts

## Tools (ui/tools/)

### Clipboard (ui/tools/clipboard/)

Cross-platform clipboard:

- **OSC 52**: Terminal clipboard (write-only)
- **Native**: OS clipboard (full support)
- **Auto-detection**: Automatic provider selection
- **History**: Clipboard history with persistence

### Shell (ui/tools/shell/)

Command execution:

- **Synchronous/Async**: Execution modes
- **Progress Callbacks**: Real-time output
- **Piping**: Command chaining
- **Timeout**: Cancellation support
- **Cross-platform**: Shell detection

### Watcher (ui/tools/watcher/)

File system monitoring:

- **Event Types**: Create, Write, Remove, Rename, Chmod
- **Filtering**: Patterns, extensions, file size
- **Debouncing**: Event storm prevention
- **Batching**: Event aggregation

## Data Flow

### Message Flow (Bubbletea)

```
User Input → tea.KeyMsg
                  ↓
            Model.Update()
                  ↓
            ┌─────────────┐
            │  New State  │
            └─────────────┘
                  ↓
            tea.Cmd (optional)
                  ↓
            Side Effects
                  ↓
            tea.Msg (result)
                  ↓
            Model.Update()
                  ↓
            Model.View()
                  ↓
            Render Output
```

### Dialog Flow

```
User Action → OpenDialogMsg
                  ↓
            ┌─────────────────┐
            │ Dialog Overlay  │
            │   (Stack)       │
            └─────────────────┘
                  ↓
            Dialog.SetFocused(true)
                  ↓
            User Interaction
                  ↓
            DialogResultMsg
                  ↓
            Callback(result, data)
                  ↓
            CloseDialogMsg
```

## Design Patterns

### 1. Interface-Based Design

Components implement small, focused interfaces:

```go
type MyComponent struct {
    // state
}

func (c *MyComponent) Focus()  { c.focused = true }
func (c *MyComponent) Blur()   { c.focused = false }
func (c *MyComponent) Focused() bool { return c.focused }

// Now MyComponent can be used anywhere Focusable is required
```

### 2. Elm Architecture (Model-View-Update)

All components follow the Elm pattern:

```go
type Model struct {
    // state
}

func (m Model) Init() tea.Cmd {
    // Initialize
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Handle events, return new state
}

func (m Model) View() string {
    // Render UI
}
```

### 3. Component Composition

```go
type AppModel struct {
    list    list.BaseList
    dialog  dialog.Overlay
    form    forms.Form
    status  status.StatusBar
}
```

### 4. Dependency Injection

```go
type CommandDialog struct {
    executeCommand func(cmd string) tea.Cmd
    getCompletions func(input string) []Completion
}
```

## Performance Considerations

### Virtualization

- **Lists**: Only render visible items
- **Viewport**: Windowed scrolling reduces render cost

### Caching

- **Style Objects**: Reuse lipgloss.Style
- **Render Output**: Cache static content
- **Sixel Images**: Protocol-level caching

### Optimization

- **strings.Builder**: Efficient string concatenation
- **Object Pools**: Buffer reuse
- **Debouncing**: Reduce event processing

## Platform Support

| Platform | Status | Notes |
|----------|--------|-------|
| Linux | ✅ Full | All features supported |
| macOS | ✅ Full | All features supported |
| Windows | ✅ Full | Some terminal limitations |

## Dependencies

### Core

- `github.com/charmbracelet/bubbletea` v1.3.10 - Elm architecture
- `github.com/charmbracelet/lipgloss` v1.1.1 - Styling

### v2.0 Engine

- `github.com/charmbracelet/ultraviolet` - High-performance rendering

### UI Components

- `github.com/charmbracelet/bubbles` v0.21.0 - Interactive components
- `github.com/charmbracelet/glamour` v0.10.0 - Markdown rendering
- `github.com/alecthomas/chroma/v2` v2.23.1 - Syntax highlighting

### Utilities

- `github.com/lucasb-eyer/go-colorful` v1.3.0 - Color blending
- `mvdan.cc/sh/v3` v3.12.0 - Shell parsing

## Best Practices

1. **Use Interfaces**: Compose through `Focusable`, `Sizeable`, etc.
2. **Theme System**: Always use `styles.DefaultStyles()` for colors
3. **Handle Resizing**: Implement `tea.WindowSizeMsg`
4. **Return Nil Commands**: Not every update needs a command
5. **Optimize Views**: Use `strings.Builder` for complex views
6. **Write Tests**: Use `DirectEngine` for unit tests

## Migration Guide

See [MIGRATION_V2.md](MIGRATION_V2.md) for details on migrating from v1.0 to v2.0.

### Key Changes

```go
// v1.0 (Bubbletea-only)
type Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (Model, tea.Cmd)
    View() string
}

// v2.0 (engine-agnostic)
type Model interface {
    Init() render.Cmd
    Update(any) (Model, render.Cmd)
    View() string
}
```

---

**Version**: 2.0.0
**Last Updated**: 2026-02-06
