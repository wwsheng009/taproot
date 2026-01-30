# Taproot Examples

This directory contains example applications demonstrating Taproot TUI Framework features.

## Table of Contents

- [Basic Examples](#basic-examples)
- [Component Examples](#component-examples)
- [Dialog Examples](#dialog-examples)
- [Advanced Examples](#advanced-examples)

## Basic Examples

### demo
**Path**: `examples/demo/main.go`

A simple counter application demonstrating basic Bubbletea patterns with Taproot styling.

```bash
go run examples/demo/main.go
```

Features:
- Keyboard controls (+/- to increment/decrement)
- Theme integration
- Status bar with help

### list
**Path**: `examples/list/main.go`

Demonstrates the virtualized list component with cursor navigation.

```bash
go run examples/list/main.go
```

Features:
- Scrollable list with items
- Cursor movement (arrows, j/k)
- Selected state highlighting

## Component Examples

### completions
**Path**: `examples/completions/main.go`

Shows the auto-complete component with fuzzy matching.

```bash
go run examples/completions/main.go
```

Features:
- Fuzzy search filtering
- Keyboard navigation
- Highlighted matching text

### filterablelist
**Path**: `examples/filterablelist/main.go`

Demonstrates filterable list with search functionality.

```bash
go run examples/filterablelist/main.go
```

Features:
- Real-time filtering
- Search result highlighting
- Large dataset handling

### groupedlist
**Path**: `examples/groupedlist/main.go`

Shows grouped list with expandable/collapsible sections.

```bash
go run examples/groupedlist/main.go
```

Features:
- Hierarchical grouping
- Expand/collapse groups
- Group header styling

### diffview
**Path**: `examples/diffview/main.go`

Displays unified diff view with syntax highlighting.

```bash
go run examples/diffview/main.go
```

Features:
- Unified diff rendering
- Color-coded additions/deletions
- Synchronized scrolling

### messages
**Path**: `examples/messages/main.go`

Demonstrates message rendering with markdown and tool calls.

```bash
go run examples/messages/main.go
```

Features:
- Markdown rendering ( Glamour)
- Tool call display
- Message list with scrolling

### image
**Path**: `examples/image/main.go`

Shows image rendering with terminal graphics protocols.

```bash
go run examples/image/main.go
```

Features:
- Kitty graphics protocol
- iTerm2 graphics protocol
- Auto-detection

## Dialog Examples

### app
**Path**: `examples/app/main.go`

Complete application demonstrating page management and dialogs.

```bash
go run examples/app/main.go
```

Features:
- Multiple pages with navigation
- Dialog stack management
- Global shortcuts
- Status bar integration

### commands
**Path**: `examples/commands/main.go`

Command palette dialog (ctrl+p).

```bash
go run examples/commands/main.go
```

Features:
- Fuzzy command search
- Parameter input
- Command execution callbacks

### models
**Path**: `examples/models/main.go`

Model selection dialog (ctrl+m).

```bash
go run examples/models/main.go
```

Features:
- Model list display
- Recent models
- Mock provider interface

### filepicker
**Path**: `examples/filepicker/main.go`

File picker dialog for directory navigation.

```bash
go run examples/filepicker/main.go
```

Features:
- Directory browsing
- File filtering
- Path completion

### quit
**Path**: `examples/quit/main.go`

Quit confirmation dialog.

```bash
go run examples/quit/main.go
```

Features:
- Unsaved changes check
- Confirmation prompt
- Custom actions

### reasoning
**Path**: `examples/reasoning/main.go`

Reasoning display dialog for AI thought process.

```bash
go run examples/reasoning/main.go
```

Features:
- Collapsible content
- Stream updates
- Markdown rendering

### sessions
**Path**: `examples/sessions/main.go`

Session switcher dialog.

```bash
go run examples/sessions/main.go
```

Features:
- Session list
- New session creation
- Session deletion

## Advanced Examples

### models
**Path**: `examples/models/main.go`

Advanced model management with API key configuration.

```bash
go run examples/models/main.go
```

## Running Examples

All examples can be run directly with `go run`:

```bash
# Run a specific example
go run examples/<name>/main.go

# Or build first
go build -o bin/<name>.exe examples/<name>/main.go
./bin/<name>.exe
```

## Example Structure

Each example follows this structure:

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/wwsheng009/taproot/internal/tui/..."
)

func main() {
    // Create model
    m := NewModel()
    
    // Run program
    p := tea.NewProgram(m, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        panic(err)
    }
}
```

## Tips for Learning

1. **Start with `demo`** - Simplest example, shows basic patterns
2. **Try `app`** - Shows complete application structure
3. **Explore dialogs** - Each demonstrates specific dialog patterns
4. **Check components** - See individual component usage

## Common Patterns

### Using Themes

```go
t := styles.CurrentTheme()
text := t.S().Base.Foreground(t.Primary).Render("Hello")
```

### Status Messages

```go
return util.ReportInfo("Operation completed")
return util.ReportError(errors.New("failed"))
```

### Dialogs

```go
// Open dialog
return func() tea.Msg {
    return dialogs.OpenDialogMsg{Model: MyDialog{}}
}

// Close dialog
return func() tea.Msg {
    return dialogs.CloseDialogMsg{}
}
```

### Pages

```go
// Register page
application.RegisterPage("home", HomePage{})

// Navigate to page
application.SetPage("home")

// Go back
application.GoBack()
```

## Contributing Examples

To add a new example:

1. Create directory in `examples/<name>/`
2. Add `main.go` with example code
3. Update this README.md
4. Ensure example follows patterns
5. Test example runs successfully
