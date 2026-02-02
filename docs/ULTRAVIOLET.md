# Ultraviolet Rendering Engine Guide

## Overview

Ultraviolet is a high-performance terminal rendering engine that provides an alternative to Bubbletea for Taproot applications. It offers direct terminal control with potentially better performance for certain use cases.

**Status**: 🟡 Experimental/Alpha

The Ultraviolet engine integration is currently in active development and should be considered experimental for production use.

## What is Ultraviolet?

Ultraviolet (`github.com/charmbracelet/ultraviolet`) is a modern TUI rendering engine designed for:

- **Performance**: Lower overhead than Bubbletea for simple applications
- **Direct Control**: More direct terminal manipulation
- **Alternative Approach**: Different event loop architecture

### Key Characteristics

| Feature | Bubbletea | Ultraviolet |
|---------|-----------|-------------|
| Architecture | Elm-style (Model-View-Update) | Direct event loop |
| Complexity | Higher (mature, feature-rich) | Lower (simpler, focused) |
| Community | Large and active | Smaller, growing |
| Ecosystem | `bubbles` component library | Limited ecosystem |
| Learning Curve | Steeper | Gentler for direct control |
| Best For | Complex, rich TUI apps | Simple, performance-critical apps |

## Architecture

### Engine Abstraction Layer

Taproot v2.0 provides an engine-agnostic interface through the `ui/render` package:

```go
// Engine interface
type Engine interface {
    Type() EngineType
    Start(model Model) error
    Stop() error
    Send(msg Msg) error
    Resize(width, height int) error
    Running() bool
}
```

This allows components to work transparently with either Bubbletea or Ultraviolet.

### Implementation Files

```
ui/render/
├── types.go          # Core interfaces (Model, Msg, Cmd, Engine)
├── engine.go         # Engine registry and factory
├── commands.go       # Command implementations (Quit, Batch, None)
├── direct.go         # DirectEngine for testing
├── adapter_tea.go    # Bubbletea engine adapter ✅ Production
└── adapter_uv.go     # Ultraviolet engine adapter 🟡 Experimental
```

## Using Ultraviolet

### Basic Setup

```go
package main

import (
    "github.com/wwsheng009/taproot/ui/render"
)

// Your model implements render.Model
type CounterModel struct {
    count int
}

func (m *CounterModel) Init() render.Cmd {
    return nil
}

func (m *CounterModel) Update(msg any) (render.Model, render.Cmd) {
    if key, ok := msg.(render.KeyMsg); ok {
        if key.Key == "q" {
            return m, render.Quit()
        }
        // Handle other keys...
    }
    return m, render.None()
}

func (m *CounterModel) View() string {
    return fmt.Sprintf("Count: %d\nPress q to quit", m.count)
}

func main() {
    // Create Ultraviolet engine
    config := render.DefaultConfig()
    config.EnableAltScreen = true

    engine, err := render.CreateEngine(render.EngineUltraviolet, config)
    if err != nil {
        panic(err)
    }

    model := &CounterModel{count: 0}

    // Start the engine (blocks until quit)
    engine.Start(model)
}
```

### Engine Configuration

```go
config := render.DefaultConfig()

// Enable mouse support
config.EnableMouse = true

// Enable alternate screen mode (clears terminal)
config.EnableAltScreen = true

// Show/hide cursor
config.EnableCursor = false

// Custom IO (nil uses stdin/stdout)
config.Input = someReader
config.Output = someWriter
```

### Creating Different Engines

```go
// Bubbletea engine (default, stable)
bubbleteaEngine, _ := render.CreateEngine(
    render.EngineBubbletea,
    render.DefaultConfig(),
)

// Ultraviolet engine (experimental)
ultravioletEngine, _ := render.CreateEngine(
    render.EngineUltraviolet,
    render.DefaultConfig(),
)

// Direct engine (for testing)
directEngine, _ := render.CreateEngine(
    render.EngineDirect,
    render.DefaultConfig(),
)
```

## Message Types

Ultraviolet adapter provides the same message types as the render abstraction:

### KeyMsg

Received when a key is pressed:

```go
func (m *Model) Update(msg any) (render.Model, render.Cmd) {
    if key, ok := msg.(render.KeyMsg); ok {
        switch key.Key {
        case "q", "ctrl+c":
            return m, render.Quit()
        case "up", "down", "left", "right":
            // Arrow keys
        }
    }
    return m, render.None()
}
```

### WindowSizeMsg

Received when terminal is resized:

```go
func (m *Model) Update(msg any) (render.Model, render.Cmd) {
    if size, ok := msg.(render.WindowSizeMsg); ok {
        m.width = size.Width
        m.height = size.Height
    }
    return m, render.None()
}
```

### TickMsg

For animation/timing:

```go
func (m *Model) Init() render.Cmd {
    return render.Tick(time.Second, func(time.Time) render.Msg {
        return TickMsg(time.Now())
    })
}

func (m *Model) Update(msg any) (render.Model, render.Cmd) {
    if _, ok := msg.(TickMsg); ok {
        m.count++
        return m, m.tickCmd()
    }
    return m, render.None()
}
```

## Commands

The render package provides engine-agnostic commands:

```go
// No operation
render.None()

// Quit the application
render.Quit()

// Combine multiple commands
render.Batch(
    cmd1,
    cmd2,
)

// Custom command
type CustomMsg struct{ data string }

render.Tick(time.Second, func(t time.Time) render.Msg {
    return CustomMsg{data: "tick"}
})
```

## Component Compatibility

All Taproot v2.0 components implement `render.Model` and work with any engine:

```go
// List components
import "github.com/wwsheng009/taproot/ui/list"

myList := list.NewBaseList()
// Works with both Bubbletea and Ultraviolet!

// Dialog components
import "github.com/wwsheng009/taproot/ui/dialog"

myDialog := dialog.NewInfoDialog("Title", "Message")
// Works with both Bubbletea and Ultraviolet!

// Message components
import "github.com/wwsheng009/taproot/ui/components/messages"

myMessage := messages.NewAssistantMessage("id", "Content")
// Works with both Bubbletea and Ultraviolet!
```

## Current Limitations

### Ultraviolet Adapter Limitations

The Ultraviolet adapter (`adapter_uv.go`) is currently experimental:

1. **Blocking Start**: The `Start()` method currently uses blocking semantics that need refinement
2. **Message Injection**: The `Send()` method doesn't fully support async message injection
3. **Graceful Shutdown**: Uses `os.Exit(0)` instead of proper cleanup
4. **Key Event Mapping**: Key events are simplified (uses string representation)
5. **No Thread Safety**: Not safe for concurrent message sending

### Missing Features

Compared to Bubbletea, Ultraviolet integration lacks:

- ⚠️ Advanced mouse events (click, drag, scroll)
- ⚠️ Focus management API
- ⚠️ Paste event handling
- ⚠️ Bracketed paste mode
- ⚠️ Custom signal handling
- ⚠️ Alternative screen buffer save/restore

## When to Use Ultraviolet

### Good Candidates For Ultraviolet

- ✅ Simple counter/display applications
- ✅ Performance-critical dashboards
- ✅ Learning and experimentation
- ✅ Applications that don't need complex interaction
- ✅ Testing/rendering validation

### Use Bubbletea Instead

- ❌ Complex multi-page applications
- ❌ Heavy user interaction (forms, navigation)
- ❌ Applications needing mouse support
- ❌ Production environments requiring stability
- ❌ Applications with sophisticated state management

## Example: Ultraviolet Counter

See `examples/ultraviolet/main.go` for a complete working example:

```bash
go run examples/ultraviolet/main.go
```

Features demonstrated:
- Counter increment/decrement
- Progress bar visualization
- Pause/resume toggle
- Reset functionality
- Quit on `q` or `ctrl+c`

## Migration from Bubbletea to Ultraviolet

### Step 1: Change Imports

```go
// Before
import "github.com/charmbracelet/bubbletea"

// After
import "github.com/wwsheng009/taproot/ui/render"
```

### Step 2: Change Type Signatures

```go
// Before
type Model struct{}
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string

// After
type Model struct{}
func (m Model) Init() render.Cmd
func (m Model) Update(msg any) (render.Model, render.Cmd)
func (m Model) View() string
```

### Step 3: Update Commands

```go
// Before
return tea.Quit

// After
return render.Quit()
```

### Step 4: Update Messages

```go
// Before
case tea.KeyMsg:

// After
case render.KeyMsg:
```

### Step 5: Change Engine Creation

```go
// Before
p := tea.NewProgram(initialModel, tea.WithAltScreen())
if _, err := p.Run(); err != nil {
    panic(err)
}

// After
engine, _ := render.CreateEngine(render.EngineBubbletea, render.DefaultConfig())
engine.Start(initialModel)
```

## Development Status

### Current Implementation

| Component | Status | Notes |
|-----------|--------|-------|
| Core interfaces | ✅ Stable | `types.go` - production ready |
| Engine registry | ✅ Stable | `engine.go` - production ready |
| Commands | ✅ Stable | `commands.go` - production ready |
| Bubbletea adapter | ✅ Stable | `adapter_tea.go` - production ready |
| Direct engine | ✅ Stable | `direct.go` - testing only |
| Ultraviolet adapter | 🟡 Experimental | `adapter_uv.go` - WIP |

### Known Issues

1. adapter_uv.go:102 - Blocking implementation needs channel-based cleanup
2. adapter_uv.go:148 - Uses os.Exit instead of proper shutdown
3. adapter_uv.go:66-70 - Key event mapping is simplified
4. adapter_uv.go:153-160 - Send() method doesn't support async

### TODO Items

- [ ] Implement proper blocking with done channel
- [ ] Add async message injection via channel
- [ ] Implement graceful shutdown without os.Exit
- [ ] Improve key event mapping
- [ ] Add mouse event support
- [ ] Add comprehensive tests for Ultraviolet adapter
- [ ] Document Ultraviolet-specific features
- [ ] Create migration guide from Bubbletea to Ultraviolet

## Contributing

To improve the Ultraviolet adapter:

1. Read `adapter_uv.go` and understand the current implementation
2. Check `adapter_tea.go` for reference on expected behavior
3. Consult Ultraviolet documentation: https://github.com/charmbracelet/ultraviolet
4. Add tests in `ui/render/render_test.go`
5. Update this documentation

## Resources

### Official Documentation

- [Ultraviolet GitHub](https://github.com/charmbracelet/ultraviolet)
- [Bubbletea documentation](https://github.com/charmbracelet/bubbletea)
- [Charm tutorials](https://charm.sh/)

### Taproot Resources

- `ui/render/types.go` - Type definitions
- `ui/render/engine.go` - Engine factory and registry
- `ui/render/adapter_tea.go` - Bubbletea reference implementation
- `examples/ultraviolet/main.go` - Working example
- `examples/demo/main.go` - Bubbletea comparison example

## FAQ

### Q: Should I use Ultraviolet or Bubbletea?

**A**: Use Bubbletea by default. Ultraviolet is experimental. Consider Ultraviolet only for:
- Simple, performance-critical applications
- Learning terminal rendering
- Contributing to the Taproot project

### Q: Can I switch engines at runtime?

**A**: No, the engine is selected at startup via `CreateEngine()`. However, your components work with both engines automatically.

### Q: Are all Bubbletea features available in Ultraviolet?

**A**: No. Ultraviolet is a simpler engine. Many Bubbletea features (complex mouse handling, alternate screen management, etc.) are not yet available in the Ultraviolet adapter.

### Q: Will my Bubbletea code work with Ultraviolet?

**A**: After updating imports and type signatures to use `render.Model` instead of `tea.Model`, yes! The engine-agnostic design ensures compatibility.

### Q: How do I know which engine is being used?

**A**: Check the `Engine.Type()` method:

```go
engine, _ := render.CreateEngine(render.EngineUltraviolet, config)
if engine.Type() == render.EngineUltraviolet {
    println("Using Ultraviolet!")
}
```

## Troubleshooting

### Issue: Ultraviolet engine doesn't start

**Solution**: Ensure you have the Ultraviolet dependency:
```bash
go get github.com/charmbracelet/ultraviolet
```

### Issue: Keys not recognized correctly

**Solution**: Ultraviolet adapter uses simplified key mapping. Check adapter_uv.go:66-70 for current implementation. Consider using Bubbletea for complex key combinations.

### Issue: Application doesn't exit cleanly

**Solution**: Ultraviolet adapter currently uses `os.Exit(0)`. This is a known limitation (adapter_uv.go:148).

### Issue: Can't run both engines in tests

**Solution**: Use `EngineDirect` for unit tests:
```go
engine, _ := render.CreateEngine(render.EngineDirect, render.DefaultConfig())
```

---

**Last Updated**: Phase 10 completion (February 2026)

**Maintained by**: Taproot Development Team
