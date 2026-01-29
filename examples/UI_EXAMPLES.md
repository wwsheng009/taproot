# UI Package Examples

This directory contains examples demonstrating the new v2.0.0 engine-agnostic UI components in `internal/ui/list/` and `internal/ui/render/`.

## Examples

### ui-list

A simple list example showing:
- `list.Item` interface implementation
- `list.Viewport` for virtualized scrolling
- `list.SelectionManager` for selection state
- `render.Model` Elm architecture
- Keyboard navigation (j/k/down/up, space, g, G, q)

```bash
go run examples/ui-list/main.go
```

**Features demonstrated:**
- Creating custom items with `list.NewListItem()`
- Setting up viewport with visible/total item count
- Managing selection state
- Handling keyboard input via `render.KeyMsg`
- Rendering views with custom formatting

### ui-filtergroup

An advanced example showing filtering and grouping:
- `list.Filter` for search/filter functionality
- `list.Group` for grouped lists
- `list.GroupManager` for group state management
- Filter mode with real-time matching
- Expand/collapse groups

```bash
go run examples/ui-filtergroup/main.go
```

**Features demonstrated:**
- Creating filterable items
- Grouping items by category
- Real-time filtering with query input
- Expanding/collapsing groups
- Dynamic group rebuild on filter

## Architecture

These components are **engine-agnostic** - they don't depend on Bubbletea or any specific rendering engine. This allows them to work with:

1. **Bubbletea** (current TUI system)
2. **Ultraviolet** (future high-performance renderer)
3. **DirectEngine** (testing, as shown in examples)

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

// Model: Elm architecture
type Model interface {
    Init() error
    Update(msg any) (Model, Cmd)
    View() string
}
```

## Comparison with Old Examples

| Old (TUI-based) | New (UI-based) |
|-----------------|----------------|
| `examples/list/` | `examples/ui-list/` |
| `examples/filterablelist/` | `examples/ui-filtergroup/` |
| `examples/groupedlist/` | (merged into ui-filtergroup) |

The old examples use `internal/tui/exp/list/` which depends on Bubbletea.
The new examples use `internal/ui/list/` which is engine-agnostic.

## Migration Path

For v2.0.0:

1. **v1.x (current)**: Continue using `internal/tui/exp/list/` and Bubbletea
2. **v2.0.0**: New code can use `internal/ui/list/` for engine flexibility
3. **Future**: Bubbletea adapter will be added to `internal/ui/render/`

## Next Steps

See `docs/V2_ROADMAP.md` for the complete v2.0.0 roadmap including:
- Phase 6.2: Enhanced dialog system
- Phase 7: Core components (file list, status display, diff viewer)
- Phase 8: Message system
- Phase 9: Layout system
