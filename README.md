<p align="center">
  <pre>
     taproot
        │
        ├─┬─┐
        │ │ │
        │ │ └── view
        │ └──── model
        └────── runtime
  </pre>

  <b>A composable TUI (Terminal User Interface) framework for Go</b>

  <em>Built on top of <a href="https://github.com/charmbracelet/bubbletea">Bubbletea</a> with engine-agnostic v2.0 architecture</em>
</p>

<p align="center">
  <a href="https://github.com/wwsheng009/taproot/actions/workflows/ci.yml">
    <img src="https://github.com/wwsheng009/taproot/actions/workflows/ci.yml/badge.svg" alt="CI" />
  </a>
  <a href="https://goreportcard.com/report/github.com/wwsheng009/taproot">
    <img src="https://goreportcard.com/badge/github.com/wwsheng009/taproot" alt="Go Report Card" />
  </a>
  <a href="https://github.com/wwsheng009/taproot/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License: MIT" />
  </a>
  <a href="https://github.com/wwsheng009/taproot/releases">
    <img src="https://img.shields.io/badge/v2.0.0-green.svg" alt="Version: 2.0.0" />
  </a>
</p>

---

Taproot provides reusable, composable components and utilities for building terminal applications in Go. Extracted from production use, it offers a solid foundation for TUI development without the boilerplate.

## ✨ Features

- **🎨 Theme System** - Dynamic themes with HCL color space blending and gradients
- **📦 Component Library** - 50+ pre-built components (dialogs, lists, forms, messages, etc.)
- **🔧 Easy Composable** - Interface-based design for maximum flexibility
- **📱 Responsive Layout** - Automatic size management and positioning
- **🎯 Type Safe** - Full type safety with compile-time guarantees
- **📝 Markdown Rendering** - Glamour-based markdown with syntax highlighting
- **🎨 Syntax Highlighting** - Chroma-powered code highlighting
- **🚀 Multi-Engine** - v2.0 supports Bubbletea, Ultraviolet, and custom engines

## 🚀 Quick Start

```bash
go get github.com/wwsheng009/taproot@latest
```

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/wwsheng009/taproot/layout"
    "github.com/wwsheng009/taproot/ui/styles"
)

func main() {
    // Get default theme
    s := styles.DefaultStyles()

    // Simple counter model
    model := counter{count: 0}

    p := tea.NewProgram(model, tea.WithAltScreen())
    p.Run()
}

type counter struct {
    count int
}

func (c counter) Init() tea.Cmd { return nil }
func (c counter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok {
        switch key.String() {
        case "q", "ctrl+c":
            return c, tea.Quit
        case "+", "=":
            c.count++
        case "-", "_":
            c.count--
        }
    }
    return c, nil
}
func (c counter) View() string {
    s := styles.DefaultStyles()
    return s.S().Title.Render(fmt.Sprintf("Count: %d", c.count))
}
```

## 📦 Components

### Core Framework
| Component | Description |
|-----------|-------------|
| **Layout** | Interfaces: Focusable, Sizeable, Positional, Help |
| **Theme** | Dynamic theming with HCL color gradients |
| **App** | Page management and dialog system |
| **Status Bar** | Info messages with TTL |

### UI Components (v2.0)
| Component | Description |
|-----------|-------------|
| **Lists** | Virtualized lists with filtering, grouping, selection |
| **Dialogs** | Info, confirm, input, select dialogs with overlay |
| **Forms** | Text input, textarea, select, checkbox, radio |
| **Messages** | Chat messages (user, assistant, tool, diagnostic, todo) |
| **Status** | LSP/MCP service status, diagnostic displays |
| **Progress** | Progress bars and spinners |
| **Attachments** | File attachment list with preview |
| **Pills** | Status/queue pills for metadata |

### Tools
| Tool | Description |
|------|------|
| **Clipboard** | Cross-platform (OSC 52 + native) with history |
| **Shell** | Command execution (sync/async, pipes, timeout) |
| **Watcher** | File system monitoring (debounce, filtering) |

## 🎨 Themes

```go
import "github.com/wwsheng009/taproot/ui/styles"

// Get default theme
s := styles.DefaultStyles()

// Use theme colors
text := s.Base.Foreground(s.Primary).Render("Hello")

// Apply gradients
gradient := styles.ApplyForegroundGrad(&s, "Gradient Text", s.Primary, s.Secondary)
```

## 📚 Examples

### v2.0 Engine-Agnostic Examples

```bash
# Virtualized list with selection
go run examples/ui-list/main.go

# Filterable and grouped list
go run examples/ui-filtergroup/main.go

# Dialog system (v2.0)
go run examples/ui-dialogs/main.go

# Form components
go run examples/forms/main.go

# Auto-complete
go run examples/autocomplete/main.go

# Messages display
go run examples/messages-demo/main.go

# Status components
go run examples/status-demo/main.go
go run examples/status-list-demo/main.go

# Progress indicators
go run examples/progress/main.go

# Attachments and pills
go run examples/attachments/main.go
go run examples/pills/main.go

# Tools: clipboard, shell
go run examples/clipboard/main.go
go run examples/shell/main.go
```

### Layout System Examples

```bash
# File browser with core layout
go run examples/file-browser-layout/main.go

# File browser with buffer layout (better wide char support)
go run examples/file-browser-buffer/main.go

# Layout demo
go run examples/layout-demo/main.go
```

### Dual Engine Examples

```bash
# Ultraviolet engine demo
go run examples/ultraviolet/main.go

# Same model, different engines
go run examples/dual-engine/main.go
go run examples/dual-engine/main.go -engine=ultraviolet
```

### Complete Application Example

```bash
# Full-featured application with multiple pages
go run examples/complete-app/main.go
```

## 🏗️ Architecture

Taproot follows the Elm Architecture (Model-View-Update) used by Bubbletea:

```
┌─────────────────────────────────────┐
│           AppModel                  │
│  ┌─────────────────────────────┐   │
│  │      Page Management        │   │
│  │  ┌──────────┐  ┌──────────┐ │   │
│  │  │  Page 1  │  │  Page 2  │ │   │
│  │  └──────────┘  └──────────┘ │   │
│  └─────────────────────────────┘   │
│  ┌─────────────────────────────┐   │
│  │     Dialog Stack            │   │
│  │  ┌──────────────────────┐   │   │
│  │  │ Commands │ Models ... │   │   │
│  │  └──────────────────────┘   │   │
│  └─────────────────────────────┘   │
│  ┌─────────────────────────────┐   │
│  │      Status Bar             │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
```

### v2.0 Multi-Engine Architecture

```
┌─────────────────────────────────────┐
│         Engine-Agnostic Layer        │
│  ┌─────────────────────────────┐   │
│  │  render.Model Interface     │   │
│  │  - Init() Cmd               │   │
│  │  - Update(msg) (Model, Cmd) │   │
│  │  - View() string            │   │
│  └─────────────────────────────┘   │
│  ┌─────────────────────────────┐   │
│  │  UI Components              │   │
│  │  - lists, dialogs, forms    │   │
│  │  - messages, status         │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
           ↓         ↓         ↓
    ┌─────────┐ ┌────────┐ ┌────────┐
    │Bubbletea│ │Ultraviol│ │ Direct │
    └─────────┘ └────────┘ └────────┘
```

## 📖 Documentation

- [Architecture](docs/ARCHITECTURE.md) - Detailed architecture analysis
- [Components](docs/COMPONENTS.md) - Complete component reference
- [API Reference](docs/API.md) - Full API documentation
- [v2.0 Migration](docs/MIGRATION_V2.md) - Migrating to v2.0
- [UI Examples](examples/UI_EXAMPLES.md) - Example descriptions

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./layout/
go test ./ui/list/
go test ./ui/dialog/
```

## 🛠️ Development

### Prerequisites

- Go 1.24 or later
- Bubbletea v1.3.10+
- Lipgloss v1.1.x

### Project Structure

```
taproot/
├── layout/          # Core interfaces (Focusable, Sizeable, etc.)
├── ui/              # UI components and theming
│   ├── styles/     # Theme system with gradients
│   ├── list/       # Virtualized list components (v2.0)
│   ├── dialog/     # Dialog system (v2.0)
│   ├── forms/      # Form components (v2.0)
│   ├── render/     # Rendering engine abstraction (v2.0)
│   ├── layout/     # Layout utilities
│   ├── components/ # UI components (messages, status, etc.)
│   └── tools/      # Tools (clipboard, shell, watcher)
├── tui/             # Framework-level components
│   ├── app/        # Application framework
│   ├── page/       # Page management
│   ├── util/       # Utilities
│   └── components/ # High-level components
├── examples/        # 50+ example programs
└── docs/           # Documentation
```

### Code Style

- Package names: lowercase
- Interfaces: `-able` suffix (Focusable, Sizeable)
- Functions: PascalCase (exported), camelCase (internal)
- Always use `styles.DefaultStyles()` for colors

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📊 Status

```
Phase 1: ████████████████████ 100% ✅ Core Framework
Phase 2: ████████████████████ 100% ✅ Application Layer
Phase 3: ████████████████████ 100% ✅ UI Components
Phase 4: ████████████████████ 100% ✅ Dialog System
Phase 5: ████████████████████ 100% ✅ v2.0 Engine-Agnostic
```

**Current Version**: 2.0.0

**Components**: 50+ core components, 50+ examples

**Test Coverage**: 30+ tests passing

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built on top of amazing projects:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - The Elm architecture for Go
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for nice terminal layouts
- [Charmbracelet Bubbles](https://github.com/charmbracelet/bubbles) - TUI components for Bubbletea
- [Ultraviolet](https://github.com/charmbracelet/ultraviolet) - High-performance TUI rendering

## 📮 Contact

For questions, suggestions, or contributions, please open an issue on GitHub.

---

**Taproot** - Deep roots, beautiful interfaces 🌳
