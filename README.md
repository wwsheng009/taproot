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

  <em>Built on top of <a href="https://github.com/charmbracelet/bubbletea">Bubbletea</a></em>
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
    <img src="https://img.shields.io/badge/v1.0.0-green.svg" alt="Version: 1.0.0" />
  </a>
</p>

---

Taproot provides reusable, composable components and utilities for building terminal applications in Go. Extracted from production use, it offers a solid foundation for TUI development without the boilerplate.

## ✨ Features

- **🎨 Theme System** - Dynamic themes with HCL color space blending and gradients
- **📦 Component Library** - Pre-built components (dialogs, lists, forms, etc.)
- **🔧 Easy Composable** - Interface-based design for maximum flexibility
- **📱 Responsive Layout** - Automatic size management and positioning
- **🎯 Type Safe** - Full type safety with compile-time guarantees
- **📝 Markdown Rendering** - Glamour-based markdown with syntax highlighting
- **🎨 Syntax Highlighting** - Chroma-powered code highlighting
- **🚀 Zero Dependencies** - Only depends on Bubbletea ecosystem

## 🚀 Quick Start

```bash
go get github.com/wwsheng009/taproot
```

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/wwsheng009/taproot/tui/app"
    "github.com/wwsheng009/taproot/tui/components/dialogs"
    "github.com/wwsheng009/taproot/tui/components/dialogs/commands"
    "github.com/wwsheng009/taproot/tui/page"
    "github.com/wwsheng009/taproot/tui/util"
)

func main() {
    // Create application
    application := app.NewApp()
    
    // Register pages
    application.RegisterPage("home", HomePage{})
    application.SetPage("home")
    
    // Run
    p := tea.NewProgram(application, tea.WithAltScreen())
    p.Run()
}

type HomePage struct{}

func (h HomePage) Init() tea.Cmd { return nil }
func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) { return h, nil }
func (h HomePage) View() string { return "Hello, Taproot!" }
```

## 📦 Components

### Core Framework
| Component | Description |
|-----------|-------------|
| **Layout** | Interfaces for composable components |
| **Theme** | Dynamic theming with gradients |
| **App** | Page management and dialog system |
| **Status Bar** | Info messages with TTL |

### UI Components
| Component | Description |
|-----------|-------------|
| **Commands** | Command palette with fuzzy search |
| **Models** | Model selection dialog |
| **Sessions** | Session management |
| **Messages** | Chat message display |
| **Lists** | Virtualized lists with filtering |
| **DiffView** | Unified diff viewer |
| **FilePicker** | File browser dialog |
| **Quit** | Unsaved changes confirmation |
| **Reasoning** | Collapsible reasoning display |
| **Image** | Terminal image rendering |

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

Run any example:

```bash
# Basic counter
go run examples/demo/main.go

# Command palette
go run examples/commands/main.go

# Model selection
go run examples/models/main.go

# Session management
go run examples/sessions/main.go

# Messages display
go run examples/messages/main.go

# Dialog system
go run examples/app/main.go
```

See the [`examples/`](examples/) directory for more examples.

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

## 📖 Documentation

- [Architecture](docs/ARCHITECTURE.md) - Detailed architecture analysis
- [Migration Plan](docs/MIGRATION_PLAN.md) - Development roadmap
- [Tasks](docs/TASKS.md) - Detailed task list
- [Alternatives](docs/ALTERNATIVES.md) - Technology choices

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./layout/
go test ./tui/util/
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
│   ├── list/       # Virtualized list components
│   ├── dialog/     # Dialog system
│   ├── layout/     # Layout utilities
│   ├── render/     # Rendering engine abstraction
│   └── components/ # UI components (files, messages, etc.)
├── tui/             # Framework-level components
│   ├── app/        # Application framework
│   ├── page/       # Page management
│   ├── anim/       # Animations
│   ├── util/       # Utilities
│   ├── components/ # High-level components
│   └── exp/        # Experimental features
├── examples/        # Example programs
├── docs/           # Documentation
└── go.mod
```

### Code Style

- Package names: lowercase
- Interfaces: `-able` suffix (Focusable, Sizeable)
- Functions: PascalCase (exported), camelCase (internal)
- Always use `styles.CurrentTheme()` for colors

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

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
Phase 5: ██████████░░░░░░░░░░  60% ✅ Advanced Components
```

**Current Version**: 1.0.0

**Components**: 38 core components, 15 examples

**Test Coverage**: 21 tests passing

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built on top of amazing projects:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - The Elm architecture for Go
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for nice terminal layouts
- [Charmbracelet Bubbles](https://github.com/charmbracelet/bubbles) - TUI components for Bubbletea

## 📮 Contact

For questions, suggestions, or contributions, please open an issue on GitHub.

---

**Taproot** - Deep roots, beautiful interfaces 🌳
