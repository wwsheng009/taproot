# Taproot Development Guide

## Project Overview

Taproot is a TUI (Terminal User Interface) framework for Go built on top of [Bubbletea](https://github.com/charmbracelet/bubbletea). It provides reusable interfaces, components, and utilities extracted from the Crush CLI project for building terminal applications.

**Module**: `github.com/yourorg/taproot`
**Go Version**: 1.24.2

## Project Structure

```
taproot/
├── internal/
│   ├── layout/           # Core TUI component interfaces
│   │   ├── layout.go     # Focusable, Sizeable, Positional, Help interfaces
│   │   └── layout_test.go
│   ├── ui/
│   │   ├── styles/       # Theme system with gradient support
│   │   │   ├── styles.go     # Theme manager, color blending, gradients
│   │   │   ├── markdown.go   # Markdown rendering with glamour
│   │   │   ├── chroma.go     # Chroma syntax highlighting theme
│   │   │   ├── palette.go    # Charmtone color palette constants
│   │   │   └── icons.go      # Icon definitions
│   │   ├── list/         # Engine-agnostic list components (v2.0)
│   │   │   ├── types.go      # Core interfaces (Item, Filterable, Selectable, etc.)
│   │   │   ├── item.go       # Basic item implementations
│   │   │   ├── list.go       # Base list with key bindings
│   │   │   ├── viewport.go   # Virtualized scroll management
│   │   │   ├── filter.go     # Filtering with match highlighting
│   │   │   ├── selection.go  # Selection state management
│   │   │   ├── group.go      # Grouped/expanding list support
│   │   │   └── list_test.go  # Tests
│   │   └── render/       # Rendering engine abstraction (v2.0)
│   │       ├── types.go      # Model, Msg, Cmd interfaces
│   │       ├── engine.go     # Engine registry and factory
│   │       ├── direct.go     # DirectEngine for testing
│   │       └── render_test.go # Tests
│   └── tui/              # Extracted TUI framework from Crush
│       ├── keys.go       # Global key bindings
│       ├── util/         # TUI utilities (Model, InfoMsg, Cursor, etc.)
│       │   └── util.go
│       ├── highlight/    # Syntax highlighting support
│       │   └── highlight.go  # Syntax highlight function using chroma
│       ├── anim/         # Animated spinner component
│       │   └── anim.go
│       └── components/
│           └── core/     # Core UI components
│               ├── core.go       # Section, Title, Status, Button helpers
│               └── status/       # Status bar component
│                   └── status.go
├── examples/             # Example applications demonstrating usage
│   ├── demo/            # Basic counter demo
│   └── list/            # Selectable list example
├── docs/                # Documentation
├── go.mod
└── go.sum
```

## Essential Commands

### Build and Run

```bash
# Run an example
go run examples/demo/main.go
go run examples/list/main.go

# Build an example
go build -o bin/demo.exe examples/demo/main.go

# Build all packages
go build ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests in specific package
go test ./internal/layout/
go test ./internal/tui/...

# Run tests with verbose output
go test -v ./...

# Run tests and open coverage in browser
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Dependency Management

```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify

# Add a new dependency
go get github.com/charmbracelet/lipgloss@latest
```

## Code Organization and Patterns

### Framework Architecture

The TUI framework (`internal/tui/`) is extracted from Crush with these components:

#### **Utilities (`util/`)**
- `Model`: Standard Bubbletea model interface (Init, Update, View)
- `InfoMsg`: Status/info messages with types (Info, Success, Warn, Error, Update)
- `ClearStatusMsg`: Message to clear status
- `ReportError()`, `ReportInfo()`, `ReportWarn()`: Command generators
- `ExecShell()`: Execute shell commands with proper TTY handling

#### **Key Bindings (`keys.go`)**
- `KeyMap`: Global key binding structure
- `DefaultKeyMap()`: Returns default keymap (quit: ctrl+c, help: ctrl+g, etc.)

#### **Styles System (`internal/ui/styles/`)**
- `Styles`: Complete style definition with colors for all UI elements
- `DefaultStyles()`: Returns the default styles
- `ApplyForegroundGrad()`: Render text with gradient
- `ChromaTheme()`: Method to get Chroma syntax highlighting theme entries
- `Markdown`, `PlainMarkdown`: Markdown style configurations

**Charmtone Palette** (`colors.go`):
Predefined color constants for the charmtone color scheme:
- **Neutrals**: Smoke, Charcoal, Oyster, Salt
- **Reds**: Coral, Sriracha, Bengal, Cherry
- **Yellows**: Zest, Butter, Citron, Cumin
- **Greens**: Guac, Julep, Guppy, Bok
- **Blues**: Malibu, Zinc, Squid
- **Purples**: Charple, Pony, Mauve, Hazy, Cheeky

#### **Animation (`anim/`)**
- `Anim`: Bubbletea component for animated spinner with gradients
- `Settings`: Configuration for size, label, colors, cycling
- Supports staggered entrance, color cycling, ellipsis animation

#### **Syntax Highlighting (`highlight/`)**
- `SyntaxHighlight()`: Perform syntax highlighting on source code
- Uses `alecthomas/chroma` for tokenization and formatting
- Supports auto-detection of lexer from filename or content
- Integrates with theme system for styled output

#### **Core Components (`components/core/`)**
- `Section()`, `SectionWithInfo()`: Render section headers with lines
- `Title()`: Render gradient title
- `Status()`: Render status with icon, title, description
- `SelectableButton()`, `SelectableButtons()`: Interactive buttons
- `SelectableButtonsVertical()`: Vertical button layout

#### **Status Bar (`components/core/status/`)**
- `StatusCmp`: Status bar component interface
- `NewStatusCmp()`: Create status bar with help and message display
- Handles InfoMsg, shows errors/warnings/success with colors

### Interface-Based Design

Taproot is built around small, composable interfaces defined in `internal/layout/`:

- **Focusable**: Components that can receive/lose focus
  - `Focus()`, `Blur()`, `Focused() bool`
- **Sizeable**: Components with dimensions
  - `Size() (width, height int)`, `SetSize(width, height int)`
- **Positional**: Components with x,y coordinates
  - `Position() (x, y int)`, `SetPosition(x, y int)`
- **Help**: Components that provide help text
  - `Help() []string`

### Bubbletea Architecture

The framework follows the Elm architecture (Model-View-Update) used by Bubbletea:

```go
// Every component implements these three methods:
type Model interface {
    Init() tea.Cmd           // Initialization
    Update(msg tea.Msg) (Model, tea.Cmd)  // Handle events
    View() string            // Render UI
}
```

### Component Implementation Pattern

Components typically:
1. Implement relevant interfaces from `layout` package or `util.Model`
2. Use a struct to hold state
3. Handle `tea.KeyMsg` for keyboard input
4. Use `strings.Builder` for efficient view rendering
5. Return `tea.Quit` command to exit
6. Use `util.InfoMsg` for status reporting

### Theme System

The theme system uses `lipgloss.Color` for all colors:

```go
// Get default styles
s := styles.DefaultStyles()

// Access colors
primary := s.Primary
mutedText := s.Base.Foreground(s.FgMuted).Render("text")

// Apply gradients
gradientText := styles.ApplyForegroundGrad(&s, "Hello", s.Primary, s.Secondary)
```

Predefined color categories:
- **Brand**: Primary, Secondary, Tertiary, Accent
- **Backgrounds**: BgBase, BgBaseLighter, BgSubtle, BgOverlay
- **Foregrounds**: FgBase, FgMuted, FgHalfMuted, FgSubtle, FgSelected
- **Borders**: Border, BorderFocus
- **Status**: Success, Error, Warning, Info

## Naming Conventions

### Files and Packages
- Package names are lowercase: `layout`, `util`, `styles`, `anim`
- File names match the primary type: `layout.go`, `anim.go`, `theme.go`
- Test files use `_test.go` suffix: `layout_test.go`
- Key bindings separated into `keys.go` files

### Code Style
- **Interfaces**: Named with `-able` suffix for capabilities: `Focusable`, `Sizeable`, `Positional`, `Help`
- **Methods**: PascalCase for exported, camelCase for internal
- **State**: Structs hold component state, typically named after the component
- **Functions**: `initialModel()` or `NewComponent()` for constructors
- **Constants**: All caps for exported: `MAX_WIDTH`

### Key Binding Patterns
- **quit**: `ctrl+c`, `q`
- **navigation**: Arrow keys, `j`/`k` (vim-style)
- **selection**: `space`, `enter`
- **increment/decrement**: `+`/`=` and `-`/`_`
- **help**: `ctrl+g`
- **commands**: `ctrl+p`

## Testing Approach

### Test Structure
- Tests use standard Go testing package
- Mock components implement all interfaces for comprehensive testing
- Table-driven tests for multiple scenarios
- Interface compliance tests using type assertions

### Mock Component Pattern
```go
type mockComponent struct {
    // implements all interface fields
}

func (m *mockComponent) Focus()        { /* ... */ }
func (m *mockComponent) Blur()         { /* ... */ }
// ... implement all methods
```

### Running Tests
- Use `go test ./...` to run all tests
- Focus on interface contracts - ensure implementations satisfy all methods
- Test edge cases: empty inputs, boundary conditions, focus transitions

## Important Gotchas

### Bubbletea Specifics
1. **Model updates return new model state**: Always return the modified model, even if no changes
2. **Commands can be nil**: Not every update needs to return a command
3. **Msg type switching**: Always use type switch for `tea.Msg`: `switch msg := msg.(type) {`
4. **View is called every frame**: Optimize view rendering, use `strings.Builder` for concatenation

### Platform Considerations
- **Windows**: Some terminal features may be limited
- **Cross-platform shell**: The `ExecShell()` function in util uses `mvdan.cc/sh/v3` for proper shell parsing

### Dependency Notes
- Uses `charmbracelet/bubbletea` v1.3.10 as core framework
- `charmbracelet/bubbles` for interactive components (help, textinput, etc.)
- `charmbracelet/lipgloss` for styling
- `charmbracelet/glamour` for markdown rendering
- `charmbracelet/x/ansi` for ANSI string manipulation
- `alecthomas/chroma/v2` for syntax highlighting
- All dependencies are properly vendored in go.sum

### Framework-Specific Patterns
1. **Message passing**: Use `util.InfoMsg` for status messages, handle in parent components
2. **Theme access**: Use `styles.DefaultStyles()` to get the styles struct
3. **Color gradients**: `ApplyForegroundGrad()` and `ApplyBoldForegroundGrad()` for gradient text
4. **Animation**: `anim.Anim` requires `tea.Tick` commands for frame updates
5. **Help system**: Use `help.Model` from bubbles with custom `KeyMap` implementations

### Import Paths
- The framework uses standard public packages (github.com)
- Not the `charm.land` proxy used in Crush development
- All imports should use `github.com/charmbracelet/*` paths

## Development Workflow

1. **Add new interface**: Define in `internal/layout/layout.go`
2. **Implement interface**: Create component struct and methods
3. **Add styling**: Use `styles.DefaultStyles()` for colors
4. **Write tests**: Add comprehensive tests in `*_test.go`
5. **Create example**: Add demo in `examples/` if applicable
6. **Run tests**: `go test ./...`
7. **Tidy dependencies**: `go mod tidy`

## Examples

The `examples/` directory demonstrates framework usage:

- **demo/main.go**: Simple counter with keyboard controls
- **list/main.go**: Selectable list with cursor navigation

Run examples with `go run examples/<name>/main.go`

## Framework Components Reference

### Utilities (`internal/tui/util/`)
- `Model`: Base interface for all Bubbletea models
- `InfoMsg`: Status message with type (Info/Success/Warn/Error/Update)
- `ReportError(err) tea.Cmd`: Report an error to status
- `ReportInfo(msg) tea.Cmd`: Report info to status
- `ReportWarn(msg) tea.Cmd`: Report warning to status
- `ExecShell()`: Execute shell commands

### Styles (`internal/ui/styles/`)
- `DefaultStyles() Styles`: Get default styles
- `ApplyForegroundGrad(s, text, c1, c2) string`: Render gradient text
- `ApplyBoldForegroundGrad(s, text, c1, c2) string`: Render bold gradient text
- `Styles` struct: Contains all styles (Base, Primary, Markdown, etc.)

### Highlight (`internal/tui/highlight/`)
- `SyntaxHighlight(s, source, fileName, bg) (string, error)`: Highlight source code
- Auto-detects lexer from filename or content
- Uses theme colors for syntax highlighting

### Animation (`internal/tui/anim/`)
- `New(opts Settings) *Anim`: Create animated spinner
- `Settings`: Size, Label, LabelColor, GradColorA/B, CycleColors
- Implements `util.Model` for Bubbletea integration

### Core Components (`internal/tui/components/core/`)
- `Section(text, width) string`: Render section header with line
- `Title(title, width) string`: Render gradient title
- `Status(opts StatusOpts, width) string`: Render status line
- `SelectableButton(opts ButtonOpts) string`: Render one button
- `SelectableButtons(opts []ButtonOpts, spacing) string`: Render button row

### Status Bar (`internal/tui/components/core/status/`)
- `NewStatusCmp() StatusCmp`: Create status bar component
- Displays help by default, shows InfoMsg when received
- Auto-clears messages after TTL (default 5 seconds)

---

## v2.0.0 Architecture (New)

The v2.0.0 release adds engine-agnostic UI components for dual-engine support (Bubbletea + Ultraviolet).

### List Components (`internal/ui/list/`)

**Core Interfaces** (`types.go`):
- `Item`: Base interface for list items (requires `ID() string`)
- `FilterableItem`: Items that can be filtered (`FilterValue() string`)
- `Sizeable`: Components with dimensions (`Size()`, `SetSize()`)
- `Focusable`: Components that can receive focus (`Focus()`, `Blur()`, `Focused()`)
- `Selectable`: Items with selection state (`Selected()`, `SetSelected()`)
- `Toggleable`: Items with expandable state (`Expanded()`, `SetExpanded()`)

**Item Implementations** (`item.go`):
- `ListItem`: Basic filterable item with ID, title, description
- `SectionItem`: Section header for grouped lists
- `SelectableItem`: Item with selection tracking
- `ExpandableItem`: Item with expandable state

**Viewport** (`viewport.go`):
- `NewViewport(visible, total)`: Create viewport
- `MoveUp()`, `MoveDown()`, `PageUp()`, `PageDown()`: Navigation
- `MoveToTop()`, `MoveToBottom()`: Jump navigation
- `Range()`: Get visible range (start, end)
- `ScrollIndicator()`: Get scroll indicator string

**Selection Manager** (`selection.go`):
- `SelectionModeNone`, `SelectionModeSingle`, `SelectionModeMultiple`
- `Select()`, `Deselect()`, `Toggle()`: Selection operations
- `SelectAll()`, `DeselectAll()`, `InvertSelection()`: Bulk operations
- `SelectedIDs()`, `GetSelected()`: Get selected items

**Filter** (`filter.go`):
- `SetQuery()`, `Clear()`: Filter management
- `Apply(items)`: Apply filter to items
- `Highlight(text, before, after)`: Highlight matches with ANSI
- `SetCaseSensitive()`: Case sensitivity toggle

**Group Manager** (`group.go`):
- `Group`: Section with title, items, expanded state
- `GroupManager`: Manages groups with flattened view
- `ExpandAll()`, `CollapseAll()`: Bulk operations
- `ToggleCurrentGroup()`: Toggle group at cursor

**Base List** (`list.go`):
- `Action`: Keyboard action constants
- `KeyMap`: Keyboard shortcut configuration
- `DefaultKeyMap()`: Returns default key bindings
- `BaseList`: Shared state for all list types

### Render Engine (`internal/ui/render/`)

**Core Types** (`types.go`):
- `Model`: Elm architecture interface (`Init()`, `Update()`, `View()`)
- `Cmd`: Command interface for side effects
- `Msg`: Message interface (events)
- `KeyMsg`: Keyboard event with modifiers
- `WindowSizeMsg`, `TickMsg`, `QuitMsg`: Standard messages

**Engine** (`engine.go`):
- `Engine`: Rendering engine interface
- `EngineType`: Bubbletea, Ultraviolet, Direct
- `RegisterEngine()`: Register engine factory
- `CreateEngine()`: Create engine by type
- `EngineConfig`: Configuration for engine initialization

**Direct Engine** (`direct.go`):
- `DirectEngine`: Simple engine for testing
- Writes to buffer (no terminal output)
- Useful for unit tests and CI/CD

### Using v2.0.0 Components

```go
import "github.com/yourorg/taproot/internal/ui/list"

// Create items
items := []list.FilterableItem{
    list.NewListItem("1", "Apple", "Red fruit"),
    list.NewListItem("2", "Banana", "Yellow fruit"),
}

// Create filter
filter := list.NewFilter()
filtered := filter.Apply(items)

// Create viewport
viewport := list.NewViewport(10, len(filtered))

// Create selection manager
selMgr := list.NewSelectionManager(list.SelectionModeMultiple)
selMgr.SelectAll(items)

// Use key map
keyMap := list.DefaultKeyMap()
action := keyMap.MatchAction("k")  // Returns ActionMoveUp
```

---

## Dialog Components (`internal/ui/dialog/`)

**Core Interfaces** (`types.go`):
- `Dialog`: Base interface for all dialogs (`ID()`, `Title()`, `SetSize()`)
- `ActionResult`: Dialog result constants (`ActionConfirm`, `ActionCancel`, `ActionDismiss`)
- `Callback`: Callback function type for dialog results

**UI Components**:
- `Button`: Single button with label and selection state
- `ButtonGroup`: Group of buttons with selection management
- `InputField`: Text input with cursor, max length, hidden (password) support
- `SelectList`: Dropdown selection list with viewport

**Dialog Types** (`dialogs.go`, `input.go`):
- `InfoDialog`: Simple informational message dialog
- `ConfirmDialog`: Yes/No confirmation dialog with callback
- `InputDialog`: Text input prompt with validation
- `SelectListDialog`: Single-select from list of items

**Overlay Manager** (`overlay.go`):
- `Overlay`: Stack manager for multiple dialogs
- `Push()`, `Pop()`, `Peek()`: Dialog stack operations
- `ActiveDialog()`: Get topmost dialog
- `FindByID()`, `RemoveByID()`: Dialog lookup operations
- `CalculateBounds()`: Center dialog on screen

**Messages**:
- `OpenDialogMsg`: Open a new dialog
- `CloseDialogMsg`: Close active dialog
- `CloseAllDialogsMsg`: Close all dialogs
- `DialogActionMsg`: Dialog action notification
- `DialogResultMsg`: Final dialog result

**Utility Functions**:
- `DefaultWidth()`, `DefaultHeight()`: Default dialog dimensions
- `MaxWidth()`, `MaxHeight()`: Maximum dimensions based on screen size

### Using Dialog Components

```go
import "github.com/yourorg/taproot/internal/ui/dialog"

// Info dialog
info := dialog.NewInfoDialog("Success", "Operation completed successfully!")
info.SetCallback(func(result dialog.ActionResult, data any) {
    // Handle acknowledgment
})

// Confirm dialog with callback
confirm := dialog.NewConfirmDialog(
    "Delete File",
    "Are you sure you want to delete this file?",
    func(result dialog.ActionResult, data any) {
        if result == dialog.ActionConfirm {
            // Execute delete
        }
    },
)

// Input dialog
input := dialog.NewInputDialog(
    "Enter Name",
    "Name:",
    func(value string) {
        fmt.Printf("User entered: %s\n", value)
    },
)
input.SetPlaceholder("John Doe")
input.SetMaxLength(50)

// Select list dialog
items := []string{"Option 1", "Option 2", "Option 3"}
select := dialog.NewSelectListDialog(
    "Choose an Option",
    items,
    func(index int, value string) {
        fmt.Printf("Selected: %s\n", value)
    },
)

// Using overlay manager
overlay := dialog.NewOverlay()
overlay.Push(confirm)
overlay.Pop()
```

---

## Related Projects

This framework is extracted from the Crush CLI tool (E:/projects/ai/crush), which contains extensive TUI implementations including:
- Virtualized lists (`internal/tui/exp/list/`)
- Diff viewers (`internal/tui/exp/diffview/`)
- Dialog systems (`internal/tui/components/dialogs/`)
- Text editors with completions (`internal/tui/components/chat/editor/`)
- Session management

Refer to Crush's TUI codebase for advanced patterns and implementations.
