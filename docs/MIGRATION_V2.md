# Taproot v1 → v2 Migration Guide

## Overview

This guide helps you migrate Taproot applications from v1.0 (Bubbletea-only) to v2.0 (engine-agnostic design).

**What Changed in v2.0:**
- Engine-agnostic component design
- Multiple rendering engine support (Bubbletea, Ultraviolet, Direct)
- New `render` package with abstracted types
- Refactored component interfaces
- Enhanced list and progress components

**Migration Difficulty:** 🟡 Moderate
Most changes are straightforward type signature updates. Applications that use standard patterns should migrate easily.

---

## Quick Reference

| Concept | v1.0 | v2.0 |
|---------|------|------|
| Model interface | `tea.Model` | `render.Model` |
| Command type | `tea.Cmd` | `render.Cmd` |
| Message type | `tea.Msg` | `any` (or `render.KeyMsg` for keys) |
| Quit command | `tea.Quit` | `render.Quit()` |
| Tick command | `tea.Tick(...)` | `render.Tick(...)` |
| Key message | `tea.KeyMsg` | `render.KeyMsg` |
| WindowSize message | `tea.WindowSizeMsg` | `render.WindowSizeMsg` |

---

## Step-by-Step Migration

### Step 1: Update Import Paths

#### Old Imports (v1.0)

```go
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/list"
)
```

#### New Imports (v2.0)

```go
import (
    "github.com/wwsheng009/taproot/ui/render"
    "github.com/wwsheng009/taproot/ui/list"
    "github.com/wwsheng009/taproot/ui/dialog"
    "github.com/wwsheng009/taproot/ui/components/status"
    // Bubbletea still needed for running
    tea "github.com/charmbracelet/bubbletea"
)
```

**Key Points:**
- All Taproot components are under `github.com/wwsheng009/taproot/ui/`
- Bubbletea import still required (unless using Ultraviolet)

---

### Step 2: Update Type Signatures

#### Old Model (v1.0)

```go
type Model struct {
    items       []list.Item
    list        list.Model
    ready       bool
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // Handle messages
    return m, nil
}

func (m Model) View() string {
    // Render view
    return "Hello"
}
```

#### New Model (v2.0)

```go
type Model struct {
    items       []list.FilterableItem  // Changed type
    list        list.BaseList         // Changed type
    ready       bool
}

func (m Model) Init() render.Cmd {          // Changed return type
    return nil
}

func (m Model) Update(msg any) (render.Model, render.Cmd) {  // Changed signature
    // Handle messages
    return m, nil
}

func (m Model) View() string {
    // Render view
    return "Hello"
}
```

**Key Changes:**
1. `tea.Model` → `render.Model`
2. `tea.Cmd` → `render.Cmd`
3. `tea.Msg` → `any` (more flexible)
4. Return type for `Update()`: `(Model, tea.Cmd)` → `(render.Model, render.Cmd)`

---

### Step 3: Update Commands

#### Old Commands (v1.0)

```go
// Quit
return m, tea.Quit

// Tick
return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
    return TickMsg(t)
})

// Batch
return m, tea.Batch(cmd1, cmd2)
```

#### New Commands (v2.0)

```go
// Quit
return m, render.Quit()

// Tick
return m, render.Tick(time.Second, func(t time.Time) render.Msg {
    return TickMsg(t)
})

// Batch
return m, render.Batch(cmd1, cmd2)
```

**Change Summary:**
- `tea.Quit` → `render.Quit()` (function call)
- `tea.Tick(...)` → `render.Tick(...)`
- `tea.Batch(...)` → `render.Batch(...)`

---

### Step 4: Update Messages

#### Key Messages

```go
// v1.0
case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyUp:
        // Handle up arrow

// v2.0
case render.KeyMsg:
    switch msg.Key {
    case "up":
        // Handle up arrow
```

**Important:** v2.0 uses string-based key mapping instead of enum types. See the [Key Mapping](#key-mapping-reference) section below.

#### Window Size Messages

```go
// v1.0
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height

// v2.0
case render.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
```

---

### Step 5: Update Component Initialization

#### List Components

```go
// v1.0 (using bubbles/list)
delegate := list.NewDefaultDelegate()
l := list.New(items, delegate, 0, 0)
l.Title = "My List"

// v2.0 (using ui/list)
keyMap := list.DefaultKeyMap()
l := list.NewBaseList()
l.SetKeyMap(keyMap)
l.SetItems(items)
```

#### Dialog Components

```go
// v1.0 - Custom dialog implementation

// v2.0 - Pre-built dialogs
dialog := dialog.NewInfoDialog(
    "Success",
    "Operation completed successfully!",
)
dialog.SetCallback(func(result dialog.ActionResult, data any) {
    // Handle result
})
```

---

### Step 6: Update Engine Creation

#### Standard Bubbletea Engine (v1.0 and v2.0)

```go
// This doesn't change for basic apps
p := tea.NewProgram(initialModel, tea.WithAltScreen())
if _, err := p.Run(); err != nil {
    panic(err)
}
```

#### Engine-Agnostic Creation (v2.0)

```go
// New v2.0 approach (engines abstracted)
config := render.DefaultConfig()
config.EnableAltScreen = true

engine, err := render.CreateEngine(render.EngineBubbletea, config)
if err != nil {
    panic(err)
}

engine.Start(initialModel)
```

**Benefits:**
- Easy to switch between engines
- Consistent interface across engines
- Testing with direct engine

---

## Key Mapping Reference

### v1.0 Key Types (Bubbletea)

```go
msg.Type == tea.KeyUp        // Enum-based
msg.Type == tea.KeyCtrlC
msg.Type == tea.KeyEnter
msg.Rune == 'a'              // Printable character
```

### v2.0 Key Strings (render)

```go
msg.Key == "up"              // String-based
msg.Key == "ctrl+c"
msg.Key == "enter"
msg.Key == "a"               // Direct character
```

### Common Key Mappings

| v1.0 | v2.0 | Description |
|------|------|-------------|
| `tea.KeyQuit`, `tea.KeyCtrlC` | `"q"`, `"ctrl+c"` | Quit |
| `tea.KeyUp` | `"up"` | Up arrow |
| `tea.KeyDown` | `"down"` | Down arrow |
| `tea.KeyLeft` | `"left"` | Left arrow |
| `tea.KeyRight` | `"right"` | Right arrow |
| `tea.KeyEnter` | `"enter"` | Enter |
| `tea.KeySpace` | `" "` | Space |
| `tea.KeyBackspace` | `"backspace"` | Backspace |
| `tea.KeyDelete` | `"delete"` | Delete |
| `tea.KeyTab` | `"tab"` | Tab |
| `msg.Rune == 'a'` | `msg.Key == "a"` | Character 'a' |

**Example: Key Handler Migration**

```go
// v1.0
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyUp:
            m.cursor--
        case tea.KeyDown:
            m.cursor++
        case tea.KeyEnter:
            return m, m.submit()
        }
        if msg.Rune == 'q' {
            return m, tea.Quit
        }
    }
    return m, nil
}

// v2.0
func (m Model) Update(msg any) (render.Model, render.Cmd) {
    switch key := msg.(type) {
    case render.KeyMsg:
        switch key.Key {
        case "up":
            m.cursor--
        case "down":
            m.cursor++
        case "enter":
            return m, m.submit()
        case "q":
            return m, render.Quit()
        }
    }
    return m, render.None()
}
```

---

## Component Migration Guide

### List Component Migration

#### v1.0 (bubbles/list)

```go
package main

import tea "github.com/charmbracelet/bubbletea"
import "github.com/charmbracelet/bubbles/list"

type item string
func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type Model struct {
    list list.Model
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}
```

#### v2.0 (ui/list)

```go
package main

import (
    "github.com/wwsheng009/taproot/ui/list"
    "github.com/wwsheng009/taproot/ui/render"
)

type MyItem struct {
    list.ListItem
}

type Model struct {
    list *list.BaseList
}

func (m Model) Init() render.Cmd {
    return nil
}

func (m Model) Update(msg any) (render.Model, render.Cmd) {
    var cmd render.Cmd
    m.list = m.list.Update(msg.(*list.BaseList))
    return m, cmd
}
```

**Key Differences:**
1. Custom item types vs. built-in `list.Item`
2. `list.Model` vs. pointer `*list.BaseList`
3. Direct `Update()` call vs. bubbletea-style pattern

---

### Dialog Component Migration

#### v1.0 (Custom Implementation)

```go
// You had to implement your own dialogs
type DialogModel struct {
    title string
    message string
    visible bool
}
```

#### v2.0 (Pre-built Dialogs)

```go
import "github.com/wwsheng009/taproot/ui/dialog"

// Ready-to-use dialog types
infoDialog := dialog.NewInfoDialog("Title", "Message")
confirmDialog := dialog.NewConfirmDialog("Confirm", "Are you sure?", callback)
inputDialog := dialog.NewInputDialog("Input", "Enter:", callback)
selectDialog := dialog.NewSelectListDialog("Choose", options, callback)
```

---

### Progress Component Migration

#### v1.0 (Manual Implementation)

```go
// You had to implement progress bars manually
type ProgressModel struct {
    percent float64
    maxWidth int
}
```

#### v2.0 (Pre-built Components)

```go
import "github.com/wwsheng009/taproot/ui/components/progress"

// Ready-to-use components
progressBar := progress.NewProgressBar()
progressBar.SetCurrent(75)
progressBar.SetTotal(100)

spinner := progress.NewSpinner(progress.SpinnerDots)
spinner.SetFPS(30)
spinner.Start()
```

---

## Common Migration Patterns

### Pattern 1: Simple Counter

```go
// v1.0
type CounterModel struct {
    count int
}

func (m CounterModel) Update(msg tea.Msg) (CounterModel, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyUp {
        m.count++
    }
    return m, nil
}

// v2.0
type CounterModel struct {
    count int
}

func (m CounterModel) Update(msg any) (render.Model, render.Cmd) {
    if key, ok := msg.(render.KeyMsg); ok && key.Key == "up" {
        m.count++
    }
    return m, render.None()
}
```

### Pattern 2: File Browser List

```go
// v1.0
type FileModel struct {
    files []file.Item
    list  list.Model
}

func (m FileModel) Update(msg tea.Msg) (FileModel, tea.Cmd) {
    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

// v2.0
type FileModel struct {
    files []list.FilterableItem
    list  *list.BaseList
}

func (m FileModel) Update(msg any) (render.Model, render.Cmd) {
    m.list = m.list.Update(msg)
    return m, render.None()
}
```

### Pattern 3: Animation with Tick

```go
// v1.0
type AnimationModel struct {
    frame int
}

func (m AnimationModel) Init() tea.Cmd {
    return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

func (m AnimationModel) Update(msg tea.Msg) (AnimationModel, tea.Cmd) {
    switch msg.(type) {
    case TickMsg:
        m.frame++
        if m.frame >= len(frames) {
            m.frame = 0
        }
        return m, m.tickCmd()
    }
    return m, nil
}

func (m AnimationModel) tickCmd() tea.Cmd {
    return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

// v2.0
type AnimationModel struct {
    frame int
}

func (m AnimationModel) Init() render.Cmd {
    return render.Tick(time.Millisecond*100, func(t time.Time) render.Msg {
        return TickMsg(t)
    })
}

func (m AnimationModel) Update(msg any) (render.Model, render.Cmd) {
    switch msg.(type) {
    case TickMsg:
        m.frame++
        if m.frame >= len(frames) {
            m.frame = 0
        }
        return m, m.tickCmd()
    }
    return m, render.None()
}

func (m AnimationModel) tickCmd() render.Cmd {
    return render.Tick(time.Millisecond*100, func(t time.Time) render.Msg {
        return TickMsg(t)
    })
}
```

---

## Testing Migration

### Unit Tests

```go
// v1.0 - Testing with mocked tea.Model
func TestModel(t *testing.T) {
    model := initialModel()
    msg := tea.KeyMsg{Type: tea.KeyUp}
    newModel, _ := model.Update(msg)
    // assertions...
}

// v2.0 - Testing with render.Model
func TestModel(t *testing.T) {
    model := initialModel()
    msg := render.KeyMsg{Key: "up"}
    newModel, _ := model.Update(msg)
    // assertions...

    // Or use DirectEngine for full integration tests
    engine, _ := render.CreateEngine(render.EngineDirect, render.DefaultConfig())
    engine.Start(model)
}
```

---

## Breaking Changes Summary

### Required Changes

1. **Import paths**: Update all Taproot imports to `github.com/wwsheng009/taproot/ui/`
2. **Type signatures**: Update Model, Cmd, Msg types
3. **Commands**: Use `render.Quit()`, `render.Tick()`, etc.
4. **Keys**: Switch from enum-based to string-based keys
5. **List components**: Use new `list.BaseList` with `FilterableItem`

### Optional Changes

1. **Engine abstraction**: Use `render.CreateEngine()` instead of direct `tea.NewProgram()`
2. **Dialogs**: Replace custom dialogs with pre-built `ui/dialog` components
3. **Progress**: Use new progress components instead of manual implementation

### No Changes Required

- **Views**: The `View()` string rendering remains the same
- **Styles**: Lipgloss usage is unchanged
- **Utilities**: Helper functions remain compatible

---

## Migration Checklist

Use this checklist to track your migration progress:

### Preparation
- [ ] Read this migration guide
- [ ] Check your current dependencies
- [ ] Create a feature branch for migration
- [ ] Run tests to ensure baseline passing

### Code Changes
- [ ] Update import statements
- [ ] Update Model type signatures
- [ ] Update Command types
- [ ] Update Message handling ( especially keys)
- [ ] Update component initialization
- [ ] Update engine creation (if desired)

### Testing
- [ ] Fix any compilation errors
- [ ] Update key handling tests
- [ ] Test all interactive features
- [ ] Verify quit functionality
- [ ] Test with different terminal sizes

### Documentation
- [ ] Update README with new imports
- [ ] Update any custom documentation
- [ ] Note version requirement in go.mod

### Deployment
- [ ] Merge migration branch
- [ ] Update CI/CD pipelines if needed
- [ ] Deploy to staging for final verification
- [ ] Merge to main and release

---

## Troubleshooting

### Issue: "undefined: tea.Model"

**Cause:** Using v1.0 type signature with v2.0 imports.

**Solution:** Update type to `render.Model`:

```go
// Wrong
func (m Model) Init() tea.Cmd

// Correct
func (m Model) Init() render.Cmd
```

### Issue: "cannot use tea.Quit (type tea.Cmd) as type render.Cmd"

**Cause:** Using v1.0 command with v2.0 type.

**Solution:** Use `render.Quit()`:

```go
// Wrong
return m, tea.Quit

// Correct
return m, render.Quit()
```

### Issue: Keys not responding

**Cause:** v2.0 uses string-based keys, not enums.

**Solution:** Update key handling:

```go
// Wrong
case tea.KeyUp:

// Correct
case "up":
```

### Issue: List component errors

**Cause:** Using v1.0 bubbles/list with v2.0 types.

**Solution:** Use ui/list components:

```go
// Wrong
import "github.com/charmbracelet/bubbles/list"

// Correct
import "github.com/wwsheng009/taproot/ui/list"
```

---

## Additional Resources

### Documentation
- [MIGRATION_PLAN.md](./MIGRATION_PLAN.md) - Original 5-phase migration strategy
- [ULTRAVIOLET.md](./ULTRAVIOLET.md) - Ultraviolet引擎指南
- [API.md](./API.md) - v2.0 API reference
- [EXAMPLES_V2.md](./EXAMPLES_V2.md) - v2.0 examples

### Examples
- `examples/demo/main.go` - Basic counter (migrated)
- `examples/list/main.go` - List component (v2.0)
- `examples/ultraviolet/main.go` - Ultraviolet example

### Code Reference
- `ui/render/types.go` - Core type definitions
- `ui/render/adapter_tea.go` - Bubbletea adapter reference
- `ui/list/types.go` - List component interfaces

---

## Getting Help

If you encounter issues during migration:

1. Check this guide's [Troubleshooting](#troubleshooting) section
2. Review [EXAMPLES_V2.md](./EXAMPLES_V2.md) for working examples
3. Check the [API Reference](./API.md) for component documentation
4. Review existing Taproot components for patterns

---

**Last Updated:** Phase 10 completion (February 2026)

**Migration Status:** 🟡 Beta - Framework stable, documentation in progress
