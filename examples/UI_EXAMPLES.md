# UI Package Examples

This directory contains examples demonstrating the new v2.0.0 engine-agnostic UI components in `internal/ui/list/` and `internal/ui/render/`.

## Interactive Examples (Bubbletea)

These examples use **Bubbletea** for real interactive TUI applications.

### ui-list

A simple interactive list example showing:
- `list.Item` interface implementation
- `list.Viewport` for virtualized scrolling
- `list.SelectionManager` for selection state
- Keyboard navigation (j/k/down/up, space, g, G, q)
- Real-time styling with lipgloss

```bash
go run examples/ui-list/main.go
```

**Key bindings:**
- `↑` / `k` - Move up
- `↓` / `j` - Move down
- `Space` / `Enter` - Toggle selection
- `g` - Jump to top
- `G` - Jump to bottom
- `Ctrl+U` - Page up
- `Ctrl+D` - Page down
- `q` / `Ctrl+C` - Quit

**Features demonstrated:**
- Creating custom items with `list.NewListItem()`
- Setting up viewport with visible/total item count
- Managing selection state with `map[string]struct{}`
- Handling keyboard input via `tea.KeyMsg`
- Rendering views with lipgloss styling

### ui-filtergroup

An interactive example showing filtering and grouping:
- `list.Filter` for search/filter functionality
- `list.Group` for grouped lists with expand/collapse
- `list.GroupManager` for group state management
- Real-time filter mode with live matching
- Expand/collapse groups with keyboard

```bash
go run examples/ui-filtergroup/main.go
```

**Key bindings:**
- `/` - Enter filter mode
- `j` / `↓` - Move down
- `k` / `↑` - Move up
- `Enter` / `Space` - Toggle group expansion or navigate
- `E` - Expand all groups
- `W` - Collapse all groups
- `q` / `Ctrl+C` - Quit

**Filter mode:**
- Type to filter items in real-time
- `Enter` - Apply filter
- `Esc` - Clear filter

**Features demonstrated:**
- Creating filterable items
- Grouping items by category (Apples, Citrus, Tropical)
- Real-time filtering with query input
- Expanding/collapsing groups
- Dynamic group rebuild on filter

## Architecture

These components are **engine-agnostic** - they don't depend directly on Bubbletea. The examples use Bubbletea for the event loop, but the core list components (`internal/ui/list/`) are independent and can work with:

1. **Bubbletea** (current TUI system - used in examples)
2. **Ultraviolet** (future high-performance renderer)
3. **Any custom engine** implementing `render.Model`

### Key Interfaces

```go
// Item: Basic list item
type Item interface {
    ID() string
}

// FilterableItem: Item that can be filtered
type FilterableItem interface {
    Item
    FilterValue() string
}

// Viewport: Virtualized scrolling
type Viewport struct {
    // MoveUp(), MoveDown(), PageUp(), PageDown()
    // Range() returns visible start/end
}

// Filter: Search functionality
type Filter struct {
    // SetQuery(), Apply(), Highlight()
}
```

## Component Integration

The examples show how to wrap engine-agnostic components for use with Bubbletea:

```go
type MyModel struct {
    // Engine-agnostic components
    viewport *list.Viewport
    selMgr   *list.SelectionManager
    filter   *list.Filter
    groupMgr *list.GroupManager
    
    // Bubbletea-specific state
    quitting bool
}

func (m *MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Convert tea.KeyMsg to actions
    // Update engine-agnostic components
    // Return new model state
}
```

## Comparison with Old Examples

| Old (TUI-based) | New (UI-based) |
|-----------------|----------------|
| `examples/list/` | `examples/ui-list/` |
| `examples/filterablelist/` | `examples/ui-filtergroup/` |
| `examples/groupedlist/` | (merged into ui-filtergroup) |

The old examples use `internal/tui/exp/list/` which has Bubbletea dependencies.
The new examples use `internal/ui/list/` which is engine-agnostic.

## Running the Examples

```bash
# Simple list with selection
go run examples/ui-list/main.go

# Filter and grouped list
go run examples/ui-filtergroup/main.go
```

Both examples run in **alternate screen mode** for a proper TUI experience.
