# Taproot

> A composable TUI (Terminal User Interface) framework for Go, built on top of [Bubbletea](https://github.com/charmbracelet/bubbletea).

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-FF6C6C?style=flat)](LICENSE)

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
go get github.com/yourorg/taproot
```

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourorg/taproot/internal/tui/app"
    "github.com/yourorg/taproot/internal/tui/components/dialogs"
    "github.com/yourorg/taproot/internal/tui/components/dialogs/commands"
    "github.com/yourorg/taproot/internal/tui/page"
    "github.com/yourorg/taproot/internal/tui/util"
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
import "github.com/yourorg/taproot/internal/ui/styles"

// Get current theme
t := styles.CurrentTheme()

// Use theme colors
text := t.S().Base.Foreground(t.Primary).Render("Hello")

// Apply gradients
gradient := styles.ApplyForegroundGrad("Gradient Text", t.Primary, t.Secondary)
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
go test ./internal/layout/
go test ./internal/tui/util/
```

## 🛠️ Development

### Prerequisites

- Go 1.24 or later
- Bubbletea v1.3.10+
- Lipgloss v1.1.x

### Project Structure

```
taproot/
├── internal/
│   ├── layout/          # Core interfaces
│   ├── ui/
│   │   ├── styles/     # Theme system
│   ├── tui/
│   │   ├── app/        # Application framework
│   │   ├── page/       # Page system
│   │   ├── anim/       # Animations
│   │   ├── util/       # Utilities
│   │   ├── components/ # UI components
│   │   └── exp/        # Experimental components
│   └── ...
├── examples/            # Example programs
├── docs/               # Documentation
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

**Current Version**: 0.9.0

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
