# Taproot Development Guide

## Project Overview

Taproot is a TUI (Terminal User Interface) framework for Go built on top of [Bubbletea](https://github.com/charmbracelet/bubbletea). It provides reusable interfaces, components, and utilities extracted from the Crush CLI project for building terminal applications.

**Module**: `github.com/wwsheng009/taproot`
**Go Version**: 1.24.2

## Project Structure

```
taproot/
├── layout/           # Core TUI component interfaces
│   ├── layout.go     # Focusable, Sizeable, Positional, Help interfaces
│   └── layout_test.go
├── ui/
│   ├── styles/       # Theme system with gradient support
│   │   ├── styles.go     # Theme manager, color blending, gradients
│   │   ├── markdown.go   # Markdown rendering with glamour
│   │   ├── chroma.go     # Chroma syntax highlighting theme
│   │   ├── palette.go    # Charmtone color palette constants
│   │   └── icons.go      # Icon definitions
│   ├── list/         # Engine-agnostic list components (v2.0)
│   │   ├── types.go      # Core interfaces (Item, Filterable, Selectable, etc.)
│   │   ├── item.go       # Basic item implementations
│   │   ├── list.go       # Base list with key bindings
│   │   ├── viewport.go   # Virtualized scroll management
│   │   ├── filter.go     # Filtering with match highlighting
│   │   ├── selection.go  # Selection state management
│   │   ├── group.go      # Grouped/expanding list support
│   │   └── list_test.go  # Tests
│   ├── render/       # Rendering engine abstraction (v2.0)
│   │   ├── types.go      # Model, Msg, Cmd interfaces
│   │   ├── engine.go     # Engine registry and factory
│   │   ├── direct.go     # DirectEngine for testing
│   │   └── render_test.go # Tests
│   └── components/   # UI components
│       ├── basic/       # Basic widgets (button, label, text)
│       ├── completions/ # Auto-completion system
│       ├── dialogs/     # Dialog system
│       ├── files/       # File list with filtering
│       ├── forms/       # Form components
│       ├── header/      # Header component
│       ├── layout/      # Layout system (flex, grid, split)
│       ├── list/        # Engine-agnostic list components
│       ├── messages/    # Message components (user, assistant, tool, fetch, diagnostic, todo)
│       ├── progress/    # Progress bars and spinners
│       ├── render/      # Rendering engine abstraction
│       ├── sidebar/     # Sidebar navigation
│       ├── status/      # Status indicators (LSP, MCP, diagnostic)
│       └── styles/      # Theme system with gradient support
├── tui/              # TUI framework
│   ├── keys.go       # Global key bindings
│   ├── util/         # TUI utilities (Model, InfoMsg, Cursor, etc.)
│   │   └── util.go
│   ├── highlight/    # Syntax highlighting support
│   │   └── highlight.go  # Syntax highlight function using chroma
│   ├── anim/         # Animated spinner component
│   │   └── anim.go
│   ├── components/   # Framework-level components
│   │   ├── core/       # Core UI components
│   │   │   ├── core.go       # Section, Title, Status, Button helpers
│   │   │   └── status/       # Status bar component
│   │   │       └── status.go
│   │   ├── completions/ # Auto-completion system
│   │   ├── dialogs/     # Dialog system (commands, file picker, etc.)
│   │   ├── logo/        # Logo component
│   │   ├── messages/    # Message components
│   │   └── image/       # Image rendering
│   ├── exp/          # Experimental components
│   │   ├── diffview/    # Diff viewer
│   │   └── list/        # Experimental list implementations
│   ├── app/          # Application framework
│   │   └── app.go
│   ├── lifecycle/    # Lifecycle management
│   │   └── lifecycle.go
│   └── page/         # Page management
│       └── page.go
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
go test ./layout/
go test ./tui/util/

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
5. For v1.0.0 (Bubbletea only): Return `tea.Quit` command to exit
6. For v2.0.0 (engine-agnostic): Return `render.Quit()` command to exit
7. Use `util.InfoMsg` for status reporting

**Note**: When using v2.0.0 engine-agnostic components with Bubbletea, use `render.Quit()` instead of `tea.Quit()`. The adapter will handle conversion automatically.

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

**Commands**:
- `None() Cmd`: Return no operation command
- `Batch(cmds ...Cmd) Cmd`: Combine multiple commands
- `Quit() Cmd`: Return quit command to exit application
- `IsQuit(cmd Cmd) bool`: Check if command is a quit command

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

**Engine Adapters**:
- `adapter_tea.go`: Bubbletea engine adapter
  - Converts `render.Model` to `tea.Model`
  - Handles `tea.Quit` from `render.Quit()` command
- `adapter_uv.go`: Ultraviolet engine adapter
  - Implements direct terminal drawing
  - Handles quit via `render.Quit()` command or `QuitMsg`

### Using v2.0.0 Components

```go
import "github.com/wwsheng009/taproot/internal/ui/list"

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
import "github.com/wwsheng009/taproot/internal/ui/dialog"

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

## Form Components (`ui/forms/`)

**Core Interfaces** (`types.go`):
- `Input`: Interface for all form inputs (`Value()`, `SetValue()`, `Focus()`, `Validate()`)

**Container** (`form.go`):
- `Form`: Manages a collection of inputs
- `NewForm(inputs ...Input)`: Create a new form
- Handles focus traversal (Tab/Shift+Tab) automatically
- Delegated event handling to child components

**Input Components**:
- `TextInput`: Single-line text input with validation and masking
- `TextArea`: Multi-line text area with word wrap
- `Select`: Dropdown selection
- `Checkbox`: Boolean toggle
- `RadioGroup`: Single selection from options

### Using Forms

```go
import "github.com/wwsheng009/taproot/ui/forms"

// Create inputs
nameInput := forms.NewTextInput("Name", "Enter your name")
nameInput.SetRequired(true)

emailInput := forms.NewTextInput("Email", "Enter email")
emailInput.SetValidation(func(s string) error {
    if !strings.Contains(s, "@") {
        return errors.New("invalid email")
    }
    return nil
})

// Create form
form := forms.NewForm(nameInput, emailInput)

// In your Update loop
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Form handles navigation and input updates
    newForm, cmd := m.form.Update(msg)
    m.form = newForm.(*forms.Form)
    
    // Check for submission
    if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
        // Check if form handles Enter (e.g. multiline), otherwise submit
        if m.form.Validate() == nil {
             // Process form data
             name := nameInput.Value()
        }
    }
    return m, cmd
}
```

---

## Message Components (`ui/components/messages/`)

**Core Interfaces** (`types.go`):
- `MessageItem`: Base interface for all message components (`ID() string`, implements `render.Model`)
- `Identifiable`: ID-based identification
- `Expandable`: Expand/collapse support (`Expanded()`, `SetExpanded()`, `ToggleExpanded()`)
- `MessageConfig`: Rendering options (expanded state, max width, etc.)

**Message Types**:

1. **AssistantMessage** (`assistant.go`, 200+ lines):
   - Markdown rendering with syntax highlighting
   - Token usage display (input/output/total)
   - Expandable content sections
   - Configurable width and expansion state

2. **UserMessage** (`user.go`, 250+ lines):
   - Plain text rendering
   - Code block support with language detection
   - File attachments display (paths, add/remove counts)
   - Copy mode toggle

3. **ToolMessage** (`tools.go`, 300+ lines):
   - Tool call details (type, name)
   - Arguments formatting and display
   - Result rendering (truncated for long content)
   - Error state indication with error codes

4. **FetchMessage** (`fetch.go`, 730+ lines):
   - Four fetch types: `FetchTypeBasic`, `FetchTypeWebFetch`, `FetchTypeWebSearch`, `FetchTypeAgentic`
   - **FetchTypeBasic**: Fast fetch for simple URL content retrieval
   - **FetchTypeWebFetch**: Fetch with URL parameter for webpage analysis
   - **FetchTypeWebSearch**: Web search with query term extraction
   - **FetchTypeAgentic**: Multi-step search + fetch with tree-structured nested message rendering
   - Request and result structures with status tracking
   - Nested message support for agentic fetch (store []MessageItem)
   - Tree-structure rendering with visual connectors (├─, └─, │)
   - Collapsible/expandable UI with auto-collapse nested messages
   - Error handling with error codes and messages
   - Loading states and completion status
   - Support for large content saved to file (SavedPath)
   - Cache-based rendering optimization

5. **DiagnosticMessage** (`diagnostics.go`, 200+ lines):
   - Diagnostic source and severity (Error, Warning, Info, Hint)
   - Code snippet with line/column highlighting
   - Expandable full message display
   - Multiple diagnostics per message support
   - Color-coded severity indicators

6. **TodoMessage** (`todos.go`, 450+ lines):
   - Todo list with status icons (✗ pending, ✓ completed, ⟳ in-progress)
   - Progress tracking (active/total counts)
   - Expandable todo items with descriptions
   - Inactive/active state support for visibility
   - Priority indicators
   - Automatic status validation

**Testing** (`messages_test.go`, 570+ lines):
- All 6 message types tested with 10+ test cases each
- Interface compliance tests (Focusable, Identifiable, Expandable)
- State management tests (expand/collapse, focus)
- Rendering tests with various configurations
- Edge case handling (empty content, large content, errors)

### Using Message Components

```go
import "github.com/wwsheng009/taproot/ui/components/messages"

// Assistant message with markdown
assistant := messages.NewAssistantMessage(
    "msg-1",
    "This is **markdown** with `code` blocks",
)
assistant.SetInputTokens(100)
assistant.SetOutputTokens(200)
assistant.SetExpanded(true)

// User message with attachments
user := messages.NewUserMessage(
    "msg-2",
    "Here's a code snippet",
)
user.AddCodeBlock("go", "func main() { println(\"Hello\") }")
user.SetFileAttachments([]string{"file1.go", "file2.go"})
user.SetFilesAdded(2)
user.SetFilesRemoved(0)

// Tool message
tool := messages.NewToolMessage(
    "msg-3",
    "bash",
    []string{"run", "build"},
    `Build started...`,
)
tool.SetResult(`Build completed in 1.2s
Binary: /build/app.exe`)

// Fetch message - Basic fetch (fast fetch for simple URL content)
basicFetch := messages.NewFetchMessage(
    "msg-4",
    messages.FetchTypeBasic,
    "Fast fetch example",
)
basicFetch.SetRequest(&messages.FetchRequest{
    URL: "https://example.com/data.json",
})
basicFetch.SetResult(&messages.FetchResult{
    Success: true,
    Content: `{"data": "value"}`,
})

// Fetch message - WebFetch (fetch with URL parameter)
webFetch := messages.NewFetchMessage(
    "msg-5",
    messages.FetchTypeWebFetch,
    "Web fetch example",
)
webFetch.SetRequest(&messages.FetchRequest{
    URL: "https://example.com/article",
    Params: map[string]string{"format": "markdown"},
})
webFetch.SetResult(&messages.FetchResult{
    Success: true,
    Summary: "Article fetched successfully",
    Content: "# Article Title\n\nContent here...",
})

// Fetch message - WebSearch (search with query)
searchFetch := messages.NewFetchMessage(
    "msg-6",
    messages.FetchTypeWebSearch,
    "Web search example",
)
searchFetch.SetRequest(&messages.FetchRequest{
    Query: "golang tutorial",
})
searchFetch.SetResult(&messages.FetchResult{
    Success: true,
    Summary: "Found 5 results",
    Content: "[websearch] 5 results found",
})

// Fetch message - Agentic (multi-step search + fetch with nested messages)
agenticFetch := messages.NewFetchMessage(
    "msg-7",
    messages.FetchTypeAgentic,
    "Agentic fetch example",
)
agenticFetch.SetRequest(&messages.FetchRequest{
    Query: "What is Bubbletea?",
    URL:   "https://charm.sh/blog/bubbletea/",
})

// Add nested tool messages for agentic fetch
webSearchMsg := messages.NewToolMessage(
    "tool-1",
    "web_search",
    []string{"Bubbletea TUI framework"},
    "",
)
webSearchMsg.SetResult("[web_search] About Bubbletea...")

webFetchMsg := messages.NewToolMessage(
    "tool-2",
    "web_fetch",
    []string{"https://charm.sh/blog/bubbletea/"},
    "",
)
webFetchMsg.SetResult(`An Elegant Terminal UI Framework for Go...`)

// Add nested messages as MessageItem interface
agenticFetch.AddNestedMessage(webSearchMsg)
agenticFetch.AddNestedMessage(webFetchMsg)

agenticFetch.SetResult(&messages.FetchResult{
    Success: true,
    Summary: "Bubbletea is a Go framework for building terminal UIs",
    Content: "Answer based on search and fetch results...",
})

// Set message as loaded
agenticFetch.SetLoaded(true)

// Diagnostic message
diag := messages.NewDiagnosticMessage(
    "msg-8",
    "compiler",
    messages.SeverityError,
    "invalid syntax",
    42,
    5,
)
diag.SetMessage("unexpected token in expression")
diag.SetCode(`func main() {
    x = 5 +  // error here
    println(x)
}`)

// Todo message
todo := messages.NewTodoMessage(
    "msg-9",
    "task-1",
    "Fix memory leak",
    messages.StatusPending,
)
todo.SetProgress(1, 5)
todo.SetDescription("Track down the issue in the cache layer")

// Add more todos
todo.AddTodo("task-2", "Write tests", messages.StatusInProgress)
todo.AddTodo("task-3", "Update docs", messages.StatusCompleted)

// Set inactive state (hide from UI)
todo.SetActive(false)

// Using in Update loop
func (m Model) Update(msg any) (render.Model, render.Cmd) {
    switch msg := msg.(type) {
    case render.KeyMsg:
        // Handle keyboard input for expandable messages
        if msg.Type == render.KeyEnter {
            if assistant.Focused() {
                assistant.ToggleExpanded()
            }
        }
    }

    // Update child components
    var cmd render.Cmd
    m.assistant, cmd = m.assistant.Update(msg)
    return m, cmd
}
```

---

## Status Components (`ui/components/status/`)

**Core Types** (`types.go`):
- `State`: Service connection state enum
  - `StateDisabled`: Service is disabled or inactive (○)
  - `StateStarting`: Service is starting up (⟳)
  - `StateReady`: Service is connected and ready (●)
  - `StateError`: Service encountered an error (×)
- `DiagnosticCounts`: Diagnostic summary by severity
  - Error, Warning, Information, Hint counts
  - Total(), HasAny(), HasErrors(), HasWarnings(), HasProblems()
  - Add(severity), Clear()
- `ToolCounts`: Tool and prompt counts for MCP services
  - Tools, Prompts counts
  - Total(), HasAny()
- `DiagnosticSeverity`: Severity levels (Error, Warning, Info, Hint)
- `Service`: Interface for service status components
- `LSPService`/`MCPService`: Service information structures

**Service Components**:

1. **ServiceCmp** (`service.go`, 200+ lines):
   - Single service status display component
   - Implements `render.Model` and `Service` interface
   - Status icons (●, ×, ○, ⟳) with colors
   - Error count display
   - Focus/blur support
   - Compact mode (icon + name only)
   - Configurable max width

2. **DiagnosticStatusCmp** (`diagnostic.go`, 285+ lines):
   - Diagnostic summary display component
   - Compact mode: errors and warnings only
   - Expanded mode: all severities with icons
   - Icons: × (error), ⚠ (warning), ⓘ (info), ∵ (hint)
   - Source-specific diagnostic tracking
   - Focus/blur, compact mode, show hints toggle

3. **LSPList** (`lsp.go`, 280+ lines):
   - Multiple LSP services display
   - Diagnostic counts per service
   - Truncation with "…and X more" indicator
   - Statistics: OnlineCount(), TotalErrors(), TotalWarnings(), TotalDiagnostics()
   - Configurable width, max items, title display

4. **MCPList** (`mcp.go`, 270+ lines):
   - Multiple MCP services display
   - Tool and prompt counts (singular/plural labels)
   - Truncation support
   - Statistics: ConnectedCount(), TotalTools(), TotalPrompts()
   - Configurable width, max items, title display

**Backward Compatibility**:
- `ServiceStatus` is an alias for `State`
- `ServiceStatusOffline`, `ServiceStatusStarting`, etc. constants map to new states

### Using Status Components

```go
import "github.com/wwsheng009/taproot/ui/components/status"

// Create a single service status component
service := status.NewService("gopls", "Go LSP")
service.SetStatus(status.StateReady)
service.SetErrorCount(3)
service.SetCompact(true)
view := service.View()

// Create diagnostic status component
diag := status.NewDiagnosticStatus("workspace")
diag.AddDiagnostic(status.DiagnosticSeverityError)
diag.AddDiagnostic(status.DiagnosticSeverityWarning)
diag.SetCompact(true)
diagView := diag.View()

// Create LSP list
lspList := status.NewLSPList()
lspList.SetWidth(50)
lspList.SetMaxItems(5)
lspList.SetShowTitle(true)

// Add LSP services
lspList.AddService(status.LSPServiceInfo{
    Name:     "gopls",
    Language: "go",
    State:    status.StateReady,
    Diagnostics: status.DiagnosticSummary{
        Error:   0,
        Warning: 2,
        Hint:    5,
    },
})
lspList.AddService(status.LSPServiceInfo{
    Name:     "rust-analyzer",
    Language: "rust",
    State:    status.StateStarting,
})
lspList.AddService(status.LSPServiceInfo{
    Name:     "clangd",
    Language: "c++",
    State:    status.StateError,
    Error:    "failed to start",
    Diagnostics: status.DiagnosticSummary{
        Error: 3,
    },
})

// Get statistics
onlineCount := lspList.OnlineCount()       // Number of ready services
totalErrors := lspList.TotalErrors()       // Total error count
totalWarnings := lspList.TotalWarnings()   // Total warning count

// Create MCP list
mcpList := status.NewMCPList()
mcpList.SetWidth(50)
mcpList.SetMaxItems(5)

// Add MCP services
mcpList.AddService(status.MCPServiceInfo{
    Name:  "filesystem",
    State: status.StateReady,
    ToolCounts: status.ToolCounts{
        Tools:   5,
        Prompts: 0,
    },
})
mcpList.AddService(status.MCPServiceInfo{
    Name:  "git",
    State: status.StateReady,
    ToolCounts: status.ToolCounts{
        Tools:   3,
        Prompts: 1,
    },
})
mcpList.AddService(status.MCPServiceInfo{
    Name:  "database",
    State: status.StateStarting,
})

// Get statistics
connectedCount := mcpList.ConnectedCount()  // Number of ready services
totalTools := mcpList.TotalTools()          // Total tool count
totalPrompts := mcpList.TotalPrompts()      // Total prompt count

// Using in Update loop
func (m Model) Update(msg any) (render.Model, render.Cmd) {
    // Update components
    m.lspList.Update(msg)
    m.mcpList.Update(msg)
    return m, render.None()
}

// Render view
func (m Model) View() string {
    var b strings.Builder
    b.WriteString(m.lspList.View())
    b.WriteString("\n\n")
    b.WriteString(m.mcpList.View())
    return b.String()
}
```

---

## Tools and Utilities (Phase 11)

### Shell Execution Tool (`ui/tools/shell/`)

**Purpose**: Cross-platform shell command execution with advanced features.

**Core Files**:
- `types.go` (267 lines): Command types, options, builders
- `executor.go` (566 lines): Execution engine with async support
- `shell_test.go` (780+ lines): Comprehensive test coverage
- `examples/shell/main.go` (500+ lines): Interactive demo

**Core Types**:

```go
// CommandType - Execution type
type CommandType int
const (
    CommandShell  // Shell command (sh/bash/cmd.exe)
    CommandDirect // Direct execution
)

// CommandOptions - Configuration
type CommandOptions struct {
    Type            CommandType
    Direction       CommandDirection
    WorkingDir      string
    Timeout         time.Duration
    Progress        func(ProgressUpdate)
    Env             []string
}

// CommandResult - Execution result
type CommandResult struct {
    ExitCode    int
    Stdout      string
    Stderr      string
    Duration    time.Duration
    Error       error
    Cancelled   bool
    TimedOut    bool
}

// CommandBuilder - Fluent interface
func NewCommandBuilder() *CommandBuilder {
    return &CommandBuilder{}
}

func (b *CommandBuilder) Command(cmd string, args ...string) *CommandBuilder
func (b *CommandBuilder) ShellCommand(cmd string) *CommandBuilder
func (b *CommandBuilder) SetWorkingDir(dir string) *CommandBuilder
func (b *CommandBuilder) SetTimeout(timeout time.Duration) *CommandBuilder
func (b *CommandBuilder) SetProgress(callback func(ProgressUpdate)) *CommandBuilder
func (b *CommandBuilder) Build() (string, []string, CommandOptions, error)
```

**Usage Examples**:

```go
import "github.com/wwsheng009/taproot/ui/tools/shell"

// Basic execution
executor := shell.NewExecutor()
result, err := executor.Execute("echo", []string{"hello"})
fmt.Printf("Output: %s\n", result.Stdout)

// Command builder
cmd, args, opts, err := shell.NewCommandBuilder().
    ShellCommand("ls -la /tmp").
    SetTimeout(5 * time.Second).
    SetWorkingDir("/home/user").
    SetProgress(func(update shell.ProgressUpdate) {
        fmt.Printf("[%s] %s", update.Stdout, update.Stderr)
    }).
    Build()

result, err := executor.ExecuteWithOptions(cmd, args, opts)

// Async execution with process tracking
process, err := executor.ExecuteAsync("long-running", []string{"task"}, opts)

// Monitor progress
for process.IsRunning() {
    fmt.Printf("PID: %d, State: %s\n", process.PID(), process.State())
    time.Sleep(100 * time.Millisecond)
}

// Cancel async process
err = executor.Cancel(process.ID())

// Pipe commands
result, err = executor.Pipe(
    []string{"cat", "file.txt"},
    []string{"grep", "pattern"},
)
```

**Features**:
- Synchronous and async execution
- Real-time progress callbacks
- Timeout and cancellation
- Working directory and environment variables
- Command piping
- Cross-platform (Windows/Unix)
- Process management (list, cancel, status)

**Path Utilities**:
```go
// Find executable in PATH
which := executor.Which("python")

// Get default shell
shellPath := shell.GetDefaultShell() // "cmd.exe" on Windows, "/bin/sh" on Unix

// Expand tilde and variables
expanded := shell.ExpandPath("~/Documents")
```

---

### File Watcher (`ui/tools/watcher/`)

**Purpose**: File system event monitoring with filtering and batching.

**Core Files**:
- `types.go` (286 lines): Event types, filters, configuration
- `watcher.go` (800+ lines): fsnotify integration and event processing
- `examples/watcher/main.go`: Interactive demo

**Core Types**:

```go
// EventType - File system event type
type EventType int
const (
    EventCreate EventType
    EventWrite
    EventRemove
    EventRename
    EventChmod
)

// Event - File system event
type Event struct {
    Type      EventType
    Path      string
    OldPath   string // for rename events
    Timestamp time.Time
}

// Filter - Event filtering configuration
type Filter struct {
    IncludePatterns   []string // glob patterns
    ExcludePatterns   []string
    Directories       bool
    HiddenFiles       bool
    EventTypes        []EventType
    MinSize, MaxSize  int64
    Extensions        []string
    Recursive         bool
}

// DebounceConfig - Debounce configuration
type DebounceConfig struct {
    Enabled      bool
    Delay        time.Duration
    MaxWait      time.Duration
    MergeEvents  bool
    MergeWindow  time.Duration
}

// Watcher - Main watcher struct
type Watcher struct {
    fsnotify    *fsnotify.Watcher
    filter      Filter
    debounce    DebounceConfig
    batch       BatchConfig
    handler     func([]Event)
    errorHandler func(error)
}
```

**Usage Examples**:

```go
import "github.com/wwsheng009/taproot/ui/tools/watcher"

// Basic watching
w, err := watcher.NewWatcher(
    filter.Filter{
        IncludePatterns: []string{"*.go"},
        Recursive:       true,
    },
    func(events []watcher.Event) {
        for _, e := range events {
            fmt.Printf("%s: %s\n", e.Type.String(), e.Path)
        }
    },
    func(err error) {
        log.Printf("Watcher error: %v", err)
    },
)

// Add paths
w.AddRecursive("/path/to/project")
w.Add("/tmp/single-file.txt")

// Start watching
w.Start()

// Wait for events
w.Wait()

// Stop watching
w.Stop()
```

**Advanced Filtering**:
```go
// Filter by extension
filter := watcher.Filter{
    Extensions: []string{".go", ".md", ".txt"},
    Recursive:  true,
}

// Filter by event types
filter.EventTypes = []watcher.EventType{
    watcher.EventWrite,
    watcher.EventCreate,
}

// Exclude patterns
filter.ExcludePatterns = []string{
    "node_modules/*",
    ".git/*",
    "*.tmp",
}

// Size filtering
filter.MinSize = 100  // minimum 100 bytes
filter.MaxSize = 1024 * 1024 // maximum 1MB
```

**Debouncing and Batching**:
```go
// Enable debouncing to reduce event storms
w.SetDebounceConfig(watcher.DebounceConfig{
    Enabled:     true,
    Delay:       100 * time.Millisecond,
    MaxWait:     500 * time.Millisecond,
    MergeEvents: true,
    MergeWindow: 50 * time.Millisecond,
})

// Enable batching
w.SetBatchConfig(watcher.BatchConfig{
    Enabled:    true,
    MaxSize:    10,
    MaxWait:    200 * time.Millisecond,
    MinSize:    1,
})
```

**Convenience Functions**:
```go
// Quick watch with defaults
err = watcher.Watch("/path", handler, errorHandler)

// Watch directory recursively
err = watcher.WatchRecursive("/path", handler, errorHandler)

// Watch specific files
err = watcher.WatchFiles([]string{"/path1", "/path2"}, handler, errorHandler)
```

**Statistics**:
```go
// Get statistics
stats := w.GetStats()
fmt.Printf("Events: %d (dropped: %d, debounced: %d, batched: %d)\n",
    stats.TotalEvents, stats.Dropped, stats.Debounced, stats.Batched)
fmt.Printf("Watching: %d files, %d directories\n",
    stats.WatchedFiles, stats.WatchedDirectories)
```

---

### Clipboard Support (`ui/tools/clipboard/`)

**Purpose**: Cross-platform clipboard operations with OSC 52 and native support.

**Core Files**:
- `types.go` (370+ lines): Clipboard types, interfaces, data structures
- `osc52.go` (270+ lines): OSC 52 terminal clipboard
- `native.go` (410+ lines): Native OS clipboard provider
- `windows.go` (170+ lines): Windows API clipboard implementation
- `manager.go` (530+ lines): Unified clipboard manager with history
- `clipboard_test.go` (430+ lines): Comprehensive tests
- `examples/clipboard/main.go` (330+ lines): Interactive demo

**Core Types**:

```go
// ClipboardType - Clipboard implementation type
type ClipboardType int
const (
    ClipboardOSC52    // Terminal clipboard (write-only)
    ClipboardNative   // OS clipboard (full support)
    ClipboardPlatform // Auto-detect
)

// ClipboardData - Clipboard data container
type ClipboardData struct {
    Format    Format
    Text      string
    Bytes     []byte
    Timestamp time.Time
}

// Format - Data format
type Format string
const (
    FormatText       = "text/plain"
    FormatHTML       = "text/html"
    FormatRTF        = "text/rtf"
    FormatImagePNG   = "image/png"
    FormatImageJPEG  = "image/jpeg"
    FormatImageGIF   = "image/gif"
)

// HistoryConfig - History configuration
type HistoryConfig struct {
    MaxItems    int
    Expiration  time.Duration
    PersistPath string
    Deduplicate bool
}

// Manager - Unified clipboard manager
type Manager struct {
    primary Provider
    history *HistoryManager
    config  ManagerConfig
}
```

**Usage Examples**:

```go
import "github.com/wwsheng009/taproot/ui/tools/clipboard"

// Create manager with default config
mgr := clipboard.NewDefaultManager()

// Simple copy/paste
err := mgr.Copy("Hello, clipboard!")
text, _ := mgr.Paste()
fmt.Println("Pasted:", text)

// Copy with custom data
data := clipboard.NewClipboardData(clipboard.FormatText, "Custom text")
err = mgr.Write(data)

// Read with format
data, err := mgr.Read(clipboard.FormatText)
fmt.Println(data.Text)
```

**History Management**:
```go
// Enable history
config := clipboard.DefaultManagerConfig()
config.HistoryConfig.MaxItems = 100
config.HistoryConfig.Deduplicate = true
mgr := clipboard.NewManager(config)

// View history
history := mgr.History()
for i, data := range history {
    fmt.Printf("[%d] %s\n", i, data.Text)
}

// Restore from history
err = mgr.RestoreFromHistory(0)

// Clear history
mgr.ClearHistory()

// Remove specific entry
mgr.RemoveFromHistory(5)
```

**OSC 52 Terminal Clipboard**:
```go
// Create OSC 52 provider
config := clipboard.DefaultOSC52Config()
config.MaxSize = 100 * 1024 // 100KB limit
config.EncodeBase64 = true

provider := clipboard.NewOSC52Provider(config)

// Check terminal support
if provider.Available() {
    fmt.Println("Terminal supports OSC 52")
    ctx := context.Background()
    data := clipboard.NewClipboardData(clipboard.FormatText, "Text via OSC 52")
    err := provider.Write(ctx, data)
}

// Encode/decode utilities
encoded := clipboard.EncodeOSC52("Hello")
decoded, _ := clipboard.DecodeOSC52(encoded)

// Parse OSC 52 sequence
selection, data, err := clipboard.ParseOSC52Sequence("\x1b]52;c;SGVsbG8=\x1b\\")
```

**Native Clipboard**:
```go
// Platform-specific clipboard
config := clipboard.DefaultNativeConfig()
config.Timeout = 5 * time.Second
config.RetryCount = 3

provider := clipboard.NewNativeProvider(config)

// Check read/write support
fmt.Printf("Read supported: %v\n", provider.IsReadSupported())
fmt.Printf("Write supported: %v\n", provider.IsWriteSupported())
fmt.Printf("Platform: %s\n", provider.GetPlatformName())

// Copy and paste
if provider.IsWriteSupported() {
    ctx := context.Background()
    data := clipboard.NewClipboardData(clipboard.FormatText, "Native clipboard")
    err := provider.Write(ctx, data)
}

if provider.IsReadSupported() {
    ctx := context.Background()
    data, err := provider.Read(ctx)
    if err == nil {
        fmt.Println("Clipboard:", data.Text)
    }
}
```

**Platform Support**:
- **Windows**: Windows API (full read/write support)
- **Linux**: xclip or xsel (full read/write support)
- **macOS**: pbcopy/pbpaste (full read/write support)
- **OSC 52**: Terminal clipboard (write-only, works on all platforms)

**Provider Switching**:
```go
// Switch between providers
err := mgr.SetProvider(clipboard.ClipboardOSC52)
err = mgr.SetProvider(clipboard.ClipboardNative)
err = mgr.SetProvider(clipboard.ClipboardPlatform) // auto-detect

// Get current provider
provider := mgr.GetProvider()
fmt.Printf("Active provider: %s\n", mgr.Type().String())

// Get platform information
info := mgr.GetPlatformInfo()
fmt.Printf("OS: %s\n", info["OS"])
fmt.Printf("OSC52 Available: %v\n", info["OSC52Available"])
fmt.Printf("Native Available: %v\n", info["NativeAvailable"])
```

**Platform Detection**:
```go
// Terminal detection
terminalName := clipboard.GetTerminalName()
terminalVersion := clipboard.GetTerminalVersion()
fmt.Printf("Terminal: %s %s\n", terminalName, terminalVersion)

// Check OSC 52 support
if clipboard.IsOSC52Supported() {
    fmt.Println("OSC 52 is supported")
}

// Tool availability (Linux/macOS)
tools := clipboard.CheckToolAvailability()
fmt.Printf("xclip: %v\n", tools["xclip"])
fmt.Printf("pbcopy: %v\n", tools["pbcopy"])
```

**Common Patterns**:

```go
// Auto-fallback copy
err := mgr.TryCopy("fallback text")
// Tries OSC 52 first, then native clipboard

// Auto-fallback read
text, err := mgr.TryRead()
// Tries native clipboard first (OSC 52 doesn't support read)

// Batch clipboard operations
for _, item := range items {
    mgr.Copy(item)
    // Automatic history tracking
}

// Clear clipboard on exit
defer mgr.Clear()

// Persist history
config.HistoryConfig.PersistPath = "/path/to/clipboard.json"
err = mgr.PersistHistory()
err = mgr.LoadHistory()
```

---

## v2.0 Development Best Practices

### Component Architecture Patterns

#### 1. Implement render.Model for All Components

All v2.0 components should implement the `render.Model` interface for engine compatibility:

```go
import "github.com/wwsheng009/taproot/ui/render"

type MyComponent struct {
    // State fields
}

func (c *MyComponent) Init() render.Cmd {
    return nil
}

func (c *MyComponent) Update(msg any) (render.Model, render.Cmd) {
    switch msg := msg.(type) {
    case render.KeyMsg:
        // Handle keyboard input
    case render.WindowSizeMsg:
        // Handle resize
    }
    return c, render.None()
}

func (c *MyComponent) View() string {
    return "Component view"
}
```

#### 2. Use Type Assertions for Messages

Since v2.0 uses `any` for messages, use type assertions:

```go
func (m Model) Update(msg any) (render.Model, render.Cmd) {
    switch msg := msg.(type) {
    case render.KeyMsg:
        // Keyboard input
    case render.WindowSizeMsg:
        // Terminal resize
    case render.TickMsg:
        // Timer tick
    case CustomMsg:
        // Custom message type
    default:
        // Unknown message - ignore
    }
    return m, render.None()
}
```

#### 3. String-Based Key Handling

v2.0 uses string-based keys instead of Bubbletea enums:

```go
// v2.0 key handling
case render.KeyMsg:
    switch key.Key {
    case "up", "k":       // Up arrow or vim k
    case "down", "j":     // Down arrow or vim j
    case "enter", " ":    // Enter or space
    case "q", "ctrl+c":   // Quit
    }
```

### State Management Patterns

#### 1. Immutable State Updates

Return new model state instead of modifying:

```go
// Good: Return updated model
func (m Model) Increment() Model {
    m.count++
    return m
}

// Avoid: Modify and return pointer
func (m *Model) Increment() *Model {
    m.count++
    return m
}
```

#### 2. Cache Management

Implement render caching for performance:

```go
type MyComponent struct {
    cache       map[string]string
    lastConfig  ComponentConfig
}

func (c *MyComponent) View() string {
    cacheKey := fmt.Sprintf("%v", c.lastConfig)

    if cached, ok := c.cache[cacheKey]; ok {
        return cached
    }

    view := c.renderView()
    c.cache[cacheKey] = view
    return view
}

func (c *MyComponent) configChanged(newConfig ComponentConfig) bool {
    return !reflect.DeepEqual(c.lastConfig, newConfig)
}
```

#### 3. Boundary Checking

Always validate state boundaries:

```go
func (p *ProgressBar) SetCurrent(current int64) {
    if p.total <= 0 {
        p.current = 0
        return
    }
    if current < 0 {
        current = 0
    }
    if current > p.total {
        current = p.total
    }
    p.current = current
}
```

### Testing Strategies

#### 1. Use DirectEngine for Tests

The DirectEngine is perfect for unit tests:

```go
func TestMyComponent(t *testing.T) {
    component := NewMyComponent()

    // Use DirectEngine for testing (no terminal)
    engine, _ := render.CreateEngine(render.EngineDirect, render.DefaultConfig())
    engine.Start(component)

    // Test update
    newComponent := component.Update(render.KeyMsg{Key: "up"})
    
    // Verify state change
    assert.Equal(t, expectedValue, newComponent.SomeField())
}
```

#### 2. Test Interface Compliance

Ensure components implement required interfaces:

```go
func TestComponentImplementsRenderModel(t *testing.T) {
    var _ render.Model = (*MyComponent)(nil)
}

func TestComponentFocusable(t *testing.T) {
    c := NewMyComponent()
    _, ok := interface{}(c).(Focusable)
    assert.True(t, ok, "Component should be Focusable")
}
```

#### 3. Test Edge Cases

Test boundary conditions:

```go
func TestProgressBarBoundaryConditions(t *testing.T) {
    bar := NewProgressBar()

    // Test negative values
    bar.SetCurrent(-10)
    assert.Equal(t, int64(0), bar.Current())

    // Test overflow
    bar.SetTotal(100)
    bar.SetCurrent(200)
    assert.Equal(t, int64(100), bar.Current())

    // Test zero total
    bar.SetTotal(0)
    assert.Equal(t, 0.0, bar.Percent())
}
```

### Command Patterns

#### 1. Quit Command

Always use `render.Quit()` for clean exit:

```go
case render.KeyMsg:
    switch key.Key {
    case "q", "ctrl+c":
        return m, render.Quit()
    }
```

#### 2. Batch Commands

Combine multiple commands:

```go
return m, render.Batch(
    cmd1,
    cmd2,
    render.Tick(time.Second, tickCallback),
)
```

#### 3. Custom Commands

Define custom message types:

```go
type CustomMsg struct {
    Data string
}

func MyCommand(data string) render.Cmd {
    return func() render.Msg {
        return CustomMsg{Data: data}
    }
}
```

### Error Handling

#### 1. Report Errors via InfoMsg

Use the status system for errors:

```go
import "github.com/wwsheng009/taproot/tui/util"

if err != nil {
    return m, util.ReportError(fmt.Errorf("operation failed: %w", err))
}
```

#### 2. Return Commands for Errors

Commands can handle errors asynchronously:

```go
func doWork() render.Cmd {
    return func() render.Msg {
        result, err := someOperation()
        if err != nil {
            return ErrorMsg{Err: err}
        }
        return SuccessMsg{Result: result}
    }
}
```

### Performance Optimization

#### 1. Efficient View Rendering

Use `strings.Builder` for concatenation:

```go
func (m *Model) View() string {
    var b strings.Builder
    b.WriteString("Header\n")
    b.WriteString(m.renderList())
    b.WriteString("Footer\n")
    return b.String()
}
```

#### 2. Lazy Evaluation

Calculate view only when needed:

```go
func (m *Model) View() string {
    if !m.visible {
        return ""
    }

    if m.cachedView != "" && !m.dirty {
        return m.cachedView
    }

    view := m.renderView()
    m.cachedView = view
    m.dirty = false
    return view
}
```

#### 3. Minimize Allocations

Reuse large buffers:

```go
type Model struct {
    viewBuffer strings.Builder
}

func (m *Model) View() string {
    m.viewBuffer.Reset()
    // Build view...
    return m.viewBuffer.String()
}
```

### Common Components Patterns

#### 1. List Component Pattern

```go
type MyListComponent struct {
    items    []list.FilterableItem
    viewport *list.Viewport
    filter   *list.Filter
    cursor   int
}

func (m *MyListComponent) Update(msg any) (render.Model, render.Cmd) {
    switch key := msg.(type) {
    case render.KeyMsg:
        action := list.DefaultKeyMap().MatchAction(key.Key)
        switch action {
        case list.ActionMoveUp:
            m.viewport.MoveUp()
        case list.ActionMoveDown:
            m.viewport.MoveDown()
        case list.ActionFilter:
            m.filter.SetQuery("") // Start filtering
        }
    }
    return m, render.None()
}
```

#### 2. Dialog Pattern

```go
type DialogModel struct {
    overlay   *dialog.Overlay
    callback  dialog.Callback
}

func NewDialog(title, message string, callback dialog.Callback) *DialogModel {
    infoDialog := dialog.NewInfoDialog(title, message)
    infoDialog.SetCallback(callback)

    overlay := dialog.NewOverlay()
    overlay.Push(infoDialog)

    return &DialogModel{
        overlay:  overlay,
        callback: callback,
    }
}
```

#### 3. Form Pattern

```go
type FormModel struct {
    form     *forms.Form
    submitted bool
}

func NewForm() *FormModel {
    nameInput := forms.NewTextInput("Name", "Enter your name")
    nameInput.SetRequired(true)

    emailInput := forms.NewTextInput("Email", "Enter email")
    emailInput.SetValidation(validateEmail)

    return &FormModel{
        form: forms.NewForm(nameInput, emailInput),
    }
}
```

### Engine Considerations

#### 1. Cross-Engine Compatibility

Write code that works with all engines:

```go
// Don't use engine-specific types
func (m Model) Update(msg tea.Msg) // ❌ Bubbletea-specific

// Use generic message type
func (m Model) Update(msg any) // ✅ Engine-agnostic
```

#### 2. Tick Message Support

Support both tick types:

```go
func (s *Spinner) Update(msg any) (render.Model, render.Cmd) {
    switch msg := msg.(type) {
    case render.TickMsg:
        // render.Tick() message
    case tea.TickMsg: // tea.Tick from adapter
        // Handle if using Bubbletea
    }
    return s, render.NextTick()
}
```

#### 3. Engine Configuration

Provide flexible configuration:

```go
type ComponentConfig struct {
    Width       int
    Height      int
    Theme       *styles.Styles
    EngineMode  EngineMode // "bubbletea" or "ultraviolet"
}
```

### Documentation Best Practices

#### 1. Document Component Purpose

```go
// ProgressBar displays progress with visual bar and percentage.
// It supports automatic boundary checking (0-100%) and can show
// labels and percentages in various modes.
type ProgressBar struct {
    current int64
    total   int64
    width   int
}
```

#### 2. Document Key Bindings

```go
// Key Bindings:
//   up/k     : Move up
//   down/j   : Move down
//   enter    : Select item
//   q/ctrl+c : Quit
```

#### 3. Provide Usage Examples

```go
// Example:
//   bar := progress.NewProgressBar()
//   bar.SetTotal(100)
//   bar.SetCurrent(75)
//   view := bar.View()
//   // Output: ████████░░░░░░░ 75/100 (75%)
```

### Debugging Tips

#### 1. Logging State Changes

```go
func (m Model) Update(msg any) (render.Model, render.Cmd) {
    debug.Printf("Received message: %T", msg)
    debug.Printf("Current state: %+v", m)
    // ...
}
```

#### 2. Visualization Mode

Add debug view mode:

```go
type Model struct {
    debugMode bool
}

func (m *Model) View() string {
    if m.debugMode {
        return fmt.Sprintf("DEBUG:\nState: %+v\nMessages: %d",
            m, m.messageCount)
    }
    return m.normalView()
}
```

#### 3. Step-Through Execution

Use tick delay for debugging:

```go
func (m Model) Init() render.Cmd {
    if debugMode {
        return render.Tick(time.Second, tickCallback) // Slow tick for debugging
    }
    return nil
}
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
