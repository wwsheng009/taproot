# UI Package Examples

This directory contains examples demonstrating the new v2.0.0 engine-agnostic UI components in `internal/ui/list/` and `internal/ui/render/`.

## Layout System Comparisons

### file-browser-layout

**Using Taproot Core Layout System** (`ui/layout/`)

A file browser demonstrating the core layout system with area-based layout calculations and lipgloss rendering.

- **Layout primitives**: `SplitHorizontal`, `SplitVertical`, `NewArea`
- **Constraint system**: `Fixed`, `Percent`, `Ratio`, `Grow`
- **Rendering**: lipgloss styles with layout delegation
- **Features**: Resizable panels, file preview, command palette, search mode

```bash
go run examples/file-browser-layout/main.go
```

**Key bindings:**
- `↑/k`, `↓/j` - Navigate
- `[`, `]` - Resize panels
- `Tab` - Toggle panels
- `!` - Command mode
- `/` - Search mode
- `u` - Parent directory

### file-browser-buffer

**Using Buffer Layout System** (`ui/render/buffer/`)

A file browser demonstrating the buffer layout system with cell-based rendering and native wide character support.

- **Buffer API**: `NewBuffer`, `SetCell`, `WriteString`, `WriteBuffer`
- **Cell-based rendering**: Exact dimension calculations
- **Native CJK support**: Built-in wide character handling
- **High performance**: ~0.15ms per frame
- **Features**: Same as file-browser-layout

```bash
go run examples/file-browser-buffer/main.go
```

**Comparison:**
| Aspect | Core Layout | Buffer Layout |
|--------|-------------|---------------|
| Layout | Area-based | Cell grid |
| Rendering | lipgloss | Direct buffer ops |
| Wide Characters | Handled by lipgloss | Native |
| Performance | Good | Excellent (~0.15ms) |
| Learning Curve | Easier (based on image.Rectangle) | More concepts |

---

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

## Ultraviolet & Dual-Engine Examples (v2.0.0)

### ultraviolet

Demonstrates the new **Ultraviolet rendering engine** interface:
- High-performance direct terminal drawing
- Same model architecture as Bubbletea, different engine
- Counter example with progress bar and pause functionality

```bash
go run examples/ultraviolet/main.go
```

**Key bindings:**
- `↑` / `→` / `+` / `=` : Increment (if not paused)
- `↓` / `←` / `-` / `_` : Decrement (if not paused)
- `Space` / `Enter` : Toggle pause
- `r` : Reset counter
- `q` / `Ctrl+C` : Quit

**Features demonstrated:**
- Using `render.EngineUltraviolet` for rendering
- Model implements `render.Model` interface
- Progress bar with Unicode block characters
- Pause/resume functionality

### dual-engine

Shows the **same model running on different engines**:
- Demonstrates engine-agnostic model architecture
- Performance comparison between Bubbletea and Ultraviolet
- History graph visualization
- Command-line argument to switch engines

```bash
# Run with Bubbletea (default)
go run examples/dual-engine/main.go

# Run with Ultraviolet
go run examples/dual-engine/main.go -engine=ultraviolet
```

**Key bindings:**
- `↑` / `→` / `+` / `=` : Increment counter
- `↓` / `←` / `-` / `_` : Decrement counter
- `Space` : Save to history
- `q` / `Ctrl+C` : Quit

**Features demonstrated:**
- Same `DualEngineModel` works with both engines
- Engine type display (`bubbletea` vs `ultraviolet`)
- Real-time FPS estimation
- History graph with ASCII visualization
- Command-line flag parsing

## Engine Comparison

| Feature | Bubbletea | Ultraviolet |
|---------|-----------|-------------|
| Maturity | ✅ Production-ready | 🔄 Beta |
| Performance | Good | ⚡ Excellent |
| API | Event-driven ( tea.Cmd ) | Direct drawing ( uv.Screen ) |
| Use Cases | Complex TUIs | High-performance UIs |
| Integration | Built-in | Adapter required |

Both engines use the **same model interface** (`render.Model`), making it easy to switch between them.

## Auto-Complete Component (v2.0.0)

### autocomplete/demo.go

Demonstrates the **engine-agnostic auto-completion system**:
- Interactive text input with inline completion suggestions
- Dynamic filtering as you type
- Multiple provider types (String, File, Command)
- ASCII popup box for completion display

```bash
go run examples/autocomplete/demo.go
```

**Key bindings:**
- `/` - Toggle completion popup
- `1` / `2` - Switch between provider types
- Type to filter completions in real-time
- `↑` / `↓` - Navigate completions
- `Enter` - Select current completion
- `Esc` - Close popup, cancel completion
- `q` / `Ctrl+C` - Quit

**Features demonstrated:**
- Engine-agnostic `AutoCompletion` component
- Three built-in provider types:
  - **StringProvider**: Simple string-based completions (e.g., fruit names)
  - **FileProvider**: File system completions with depth control
  - **CommandProvider**: Command completions with action handlers
- Real-time filtering with match highlighting
- Popup positioning and dimension calculation
- Cursor navigation (circular wrapping)

**Provider Type 1 - String Provider:**
- Simple list of strings
- Best for: keywords, options, fixed values
- Example data: Apple, Banana, Cherry, Grape, etc.

**Provider Type 2 - Command Provider:**
- Completions with associated actions
- Best for: commands, functions with callbacks
- Example data: help, clear, quit (with handlers)

**Built-in Providers:**
The `completions` package includes three ready-to-use providers:

```go
// StringProvider - Simple string completions
provider := completions.NewStringProvider([]string{
    "Apple", "Banana", "Cherry",
})

// FileProvider - File system completions
provider := completions.NewFileProvider(
    ".",
    false,  // caseSensitive
    "/",    // ignorePathSeparator
    3,      // maxDepth
    nil,    // ignoreDirs
)

// CommandProvider - Commands with handlers
handler := func(value string) error {
    fmt.Printf("Executing: %s\n", value)
    return nil
}

cmd1 := completions.NewCommandItem("help", "Show help", handler)
cmd2 := completions.NewCommandItem("clear", "Clear screen", handler)

provider := completions.NewCommandProvider([]*completions.CommandItem{cmd1, cmd2})
```

**Custom Providers:**
Create custom providers by implementing the `Provider` interface:

```go
type MyProvider struct {
    items []CompletionItem
}

func (p *MyProvider) GetCompletions() ([]CompletionItem, error) {
    return p.items, nil
}
```
