# Taproot v2.0 Examples

This document provides an overview of all example applications included with Taproot v2.0, organized by functionality and difficulty level.

## Quick Reference

| Category | Examples | Description |
|----------|----------|-------------|
| **Getting Started** | 3 | Basic concepts and simple demos |
| **Lists & Navigation** | 8 | List components, filtering, grouping |
| **Dialogs & Forms** | 5 | User input and dialogs |
| **Advanced UI** | 12 | Messages, diff, status, progress |
| **Layout & Styling** | 6 | Layout systems and visual components |
| **Engine & Architecture** | 3 | Engine types and rendering |

**Total Examples:** 37 interactive programs

---

## Getting Started ⭐

These examples introduce basic Taproot concepts and are recommended for beginners.

### demo - Simple Counter

**Location:** `examples/demo/main.go`

**Difficulty:** ★☆☆☆☆

**Description:** The most basic example - a simple counter with increment/decrement controls.

**Key Concepts:**
- Basic Model-View-Update architecture
- render.Model interface
- Simple key handling
- State management

**Run:**
```bash
go run examples/demo/main.go
```

**Key Bindings:**
- `↑` / `+` - Increment counter
- `↓` / `-` - Decrement counter
- `q` / `Ctrl+C` - Quit

---

### list - Selectable List

**Location:** `examples/list/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Introduces basic list component with navigation and selection.

**Key Concepts:**
- List component basics
- Item interfaces
- Cursor management
- Selection state

**Run:**
```bash
go run examples/list/main.go
```

**Key Bindings:**
- `↑` / `k` - Move cursor up
- `↓` / `j` - Move cursor down
- `q` / `Ctrl+C` - Quit

---

### ultraviolet - Ultraviolet Engine Demo

**Location:** `examples/ultraviolet/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Demonstrates the Ultraviolet rendering engine with a counter and progress bar.

**Key Concepts:**
- Multiple rendering engines
- Engine abstraction
- Cross-engine compatibility

**Run:**
```bash
go run examples/ultraviolet/main.go
```

**Key Bindings:**
- `↑` / `→` / `+` / `=` - Increment counter
- `↓` / `←` / `-` / `_` - Decrement counter
- `Space` / `Enter` - Toggle pause
- `r` - Reset counter
- `q` / `Ctrl+C` - Quit

---

## Lists & Navigation ⭐⭐

Examples demonstrating advanced list features, filtering, grouping, and file navigation.

### ui-list - Engine-Agnostic List

**Location:** `examples/ui-list/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Shows the engine-agnostic list component with virtual scrolling and selection.

**Key Concepts:**
- `list.FilterableItem` interface
- `list.Viewport` for virtualized scrolling
- `list.SelectionManager` for multi-selection
- Navigation operations (MoveUp, MoveDown, PageUp, PageDown)

**Run:**
```bash
go run examples/ui-list/main.go
```

**Key Bindings:**
- `↑` / `k` - Move up
- `↓` / `j` - Move down
- `Space` / `Enter` - Toggle selection
- `g` - Jump to top
- `G` - Jump to bottom
- `Ctrl+U` - Page up
- `Ctrl+D` - Page down
- `q` / `Ctrl+C` - Quit

---

### filterablelist - Filtered List

**Location:** `examples/filterablelist/main.go`

**Difficulty:** ★★★☆☆

**Description:** Interactive list with real-time search and filtering.

**Key Concepts:**
- `list.Filter` for search functionality
- Real-time filtering
- Filter highlighting
- Filter query management

**Run:**
```bash
go run examples/filterablelist/main.go
```

**Key Bindings:**
- `/` - Enter filter mode
- `j` / `↓` - Move down
- `k` / `↑` - Move up
- `Enter` - Apply filter
- `Esc` - Clear filter
- `q` / `Ctrl+C` - Quit

---

### groupedlist - Grouped List

**Location:** `examples/groupedlist/main.go`

**Difficulty:** ★★★☆☆

**Description:** List with expandable/collapsible groups (like tree structure).

**Key Concepts:**
- `list.Group` for grouped items
- `list.GroupManager` for group state
- Expand/collapse operations
- Group navigation

**Run:**
```bash
go run examples/groupedlist/main.go
```

**Key Bindings:**
- `j` / `↓` - Move down
- `k` / `↑` - Move up
- `Enter` / `Space` - Toggle group expansion
- `E` - Expand all groups
- `W` - Collapse all groups
- `q` / `Ctrl+C` - Quit

---

### ui-filtergroup - Filter & Group Combined

**Location:** `examples/ui-filtergroup/main.go`

**Difficulty:** ★★★★☆

**Description:** Advanced list combining filtering and grouping features.

**Key Concepts:**
- Combined filter and group functionality
- Dynamic group rebuilding
- Real-time filter mode
- Expand/collapse with filtering

**Run:**
```bash
go run examples/ui-filtergroup/main.go
```

**Key Bindings:**
- `/` - Enter filter mode
- `j` / `↓` - Move down
- `k` / `↑` - Move up
- `Enter` / `Space` - Toggle group expansion
- `E` - Expand all groups
- `W` - Collapse all groups
- `q` / `Ctrl+C` - Quit

---

### files - File List

**Location:** `examples/files/main.go`

**Difficulty:** ★★★☆☆

**Description:** Interactive file browser with list component.

**Key Concepts:**
- File system navigation
- Directory traversal
- File type detection
- Item metadata

**Run:**
```bash
go run examples/files/main.go
```

**Key Bindings:**
- `↑` / `k` - Move up
- `↓` / `j` - Move down
- `Enter` - Enter directory
- `Esc` - Go to parent directory
- `q` / `Ctrl+C` - Quit

---

### treefiles - Tree File Browser

**Location:** `examples/treefiles/main.go`

**Difficulty:** ★★★★☆

**Description:** Hierarchical file tree with expand/collapse.

**Key Concepts:**
- Tree data structure
- Nested navigation
- Expand/collapse logic
- Depth tracking

**Run:**
```bash
go run examples/treefiles/main.go
```

**Key Bindings:**
- `↑` / `k` - Move up
- `↓` / `j` - Move down
- `Enter` - Enter directory
- `Right` / `l` - Expand directory
- `Left` / `h` - Collapse directory / Go to parent
- `q` / `Ctrl+C` - Quit

---

### attachments - Attachment List (Phase 10)

**Location:** `examples/attachments/main.go`

**Difficulty:** ★★★☆☆

**Description:** File attachment management with type detection and filtering.

**Key Concepts:**
- `AttachmentList` component
- File type detection (image, video, audio, document, archive)
- MIME type auto-detection
- Size formatting (KB, MB, GB)
- Filter by name/type
- Selection management

**Run:**
```bash
go run examples/attachments/main.go
```

**Key Bindings:**
- `↑` / `↓` - Navigate attachments
- `a` - Add new attachment
- `r` - Remove current attachment
- `c` - Toggle compact mode
- `p` - Toggle preview mode
- `s` - Toggle size display
- `Esc` - Exit

**Sample Attachments:**
- PDF documents
- Images (PNG, JPG)
- Videos (MP4, AVI)
- Audio (MP3, WAV)
- Archives (ZIP, TAR)
- Text files

---

### file-browser - Advanced File Browser

**Location:** `examples/file-browser/`

**Difficulty:** ★★★★☆

**Description:** Full-featured file browser with filtering and browsing.

**Key Concepts:**
- Advanced file operations
- Filtered file listing
- Navigation history
- File metadata display

**Run:**
```bash
cd examples/file-browser
go run main.go
```

---

## Dialogs & Forms ⭐⭐

Examples demonstrating user input collection through dialogs and forms.

### ui-dialogs - Dialog System

**Location:** `examples/ui-dialogs/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Demonstrates the dialog system with various dialog types.

**Key Concepts:**
- `dialog.Overlay` for dialog stack
- Info, Confirm, Input, SelectList dialogs
- Dialog callbacks and results
- Dialog lifecycle

**Run:**
```bash
go run examples/ui-dialogs/main.go
```

**Key Bindings:**
- `1` - Show info dialog
- `2` - Show confirm dialog
- `3` - Show input dialog
- `4` - Show select list dialog
- `q` / `Ctrl+C` - Quit

---

### forms - Form Components

**Location:** `examples/forms/main.go`

**Difficulty:** ★★★☆☆

**Description:** Form system with multiple input types and validation.

**Key Concepts:**
- `forms.Form` container
- `forms.TextInput`, `forms.TextArea`, `forms.Select`
- `forms.Checkbox`, `forms.RadioGroup`
- Input validation
- Focus traversal (Tab/Shift+Tab)

**Run:**
```bash
go run examples/forms/main.go
```

**Key Bindings:**
- `Tab` / `Shift+Tab` - Navigate between inputs
- `Enter` - Submit form
- `q` / `Ctrl+C` - Quit

**Form Fields:**
- Text input (name)
- Email input (with validation)
- Multi-line text area
- Dropdown select
- Checkbox
- Radio group

---

### app - Multi-Page Application

**Location:** `examples/app/main.go`

**Difficulty:** ★★★★☆

**Description:** Complete app with page system, dialogs, and keyboard shortcuts.

**Key Concepts:**
- Page management system
- Dialog integration
- Keyboard command palette
- Status bar
- Multiple UI components working together

**Run:**
```bash
go run examples/app/main.go
```

**Key Bindings:**
- `Ctrl+P` - Open command palette
- `q` / `Ctrl+C` - Quit

---

### completions - Auto-Completion System

**Location:** `examples/completions/main.go`

**Difficulty:** ★★★☆☆

**Description:** Auto-completion dropdown with fuzzy matching.

**Key Concepts:**
- Completion component
- Fuzzy search
- Keyboard navigation
- Selection highlighting

**Run:**
```bash
go run examples/completions/main.go
```

**Key Bindings:**
- Type to filter completions
- `↑` / `↓` - Navigate options
- `Enter` - Select completion
- `Esc` - Close
- `q` - Quit

---

### commands - Command Palette

**Location:** `examples/commands/main.go`

**Difficulty:** ★★★☆☆

**Description:** Command palette similar to VS Code/IntelliJ.

**Key Concepts:**
- Command registry
- Fuzzy command search
- Command execution
- Keyboard shortcuts

**Run:**
```bash
go run examples/commands/main.go
```

**Key Bindings:**
- `Ctrl+P` - Open command palette
- Type to filter commands
- `Enter` - Execute command
- `Esc` - Close

---

## Advanced UI ⭐⭐⭐

Examples showcasing advanced UI components like messaging, status, progress, and diff viewing.

### messages - Message Components

**Location:** `examples/messages/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Various message types (user, assistant, tool, fetch) with markdown and syntax highlighting.

**Key Concepts:**
- `messages.UserMessage`
- `messages.AssistantMessage`
- `messages.ToolMessage`
- `messages.FetchMessage`
- Markdown rendering
- Syntax highlighting
- Expand/collapse

**Run:**
```bash
go run examples/messages/main.go
```

**Key Bindings:**
- `↑` / `↓` - Navigate messages
- `Enter` - Toggle expansion
- `q` / `Ctrl+C` - Quit

**Message Types:**
- User messages with code blocks
- Assistant messages with markdown
- Tool call messages with results
- Fetch messages (basic, web-fetch, web-search, agentic)

---

### status-demo - Status Indicators

**Location:** `examples/status-demo/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Status bar and service status components.

**Key Concepts:**
- `status.ServiceCmp`
- `status.DiagnosticStatusCmp`
- Service states (starting, ready, error)
- Diagnostic counts

**Run:**
```bash
go run examples/status-demo/main.go
```

**Key Bindings:**
- `1` - Set state to starting
- `2` - Set state to ready
- `3` - Set state to error
- `q` / `Ctrl+C` - Quit

---

### status-list-demo - Service List

**Location:** `examples/status-list-demo/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Multiple services display (LSP and MCP services).

**Key Concepts:**
- `status.LSPList`
- `status.MCPList`
- Multiple service tracking
- Statistics (online count, total tools, etc.)
- Truncation support

**Run:**
```bash
go run examples/status-list-demo/main.go
```

**Key Bindings:**
- `q` / `Ctrl+C` - Quit

**Services Displayed:**
- LSP services (gopls, rust-analyzer, clangd)
- MCP services (filesystem, git, database)

---

### progress - Progress Bars & Animations (Phase 10)

**Location:** `examples/progress/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Progress bars and spinner animations.

**Key Concepts:**
- `progress.ProgressBar` with percentage display
- `progress.Spinner` with 4 animation types
- Real-time updates
- FPS control
- State management (start, stop, reset)

**Run:**
```bash
go run examples/progress/main.go
```

**Key Bindings:**
- `+` / `-` - Increment/decrement progress (bar 1)
- `1` / `2` / `3` - Select progress bar
- `r` - Reset progress
- `t` - Cycle spinner types (Dots → Line → Arrow → Moon)
- `s` - Start/stop spinner
- `r` - Reset spinner
- `q` / `Ctrl+C` - Quit

**Spinner Types:**
- Dots: ⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏
- Line: / - \ |
- Arrow: ← ↑ → ↓
- Moon: 🌑 🌒 🌓 🌔 🌕 🌖 🌗 🌘

---

### pills - Pills Status List (Phase 10)

**Location:** `examples/pills/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Status badge pills with 7 preset statuses and expand/collapse.

**Key Concepts:**
- `pills.PillList` component
- 7 pill statuses (Pending, InProgress, Completed, Error, Warning, Info, Neutral)
- Status icons (☐ ⟳ ✓ × ⚠ ℹ •)
- Badge counts
- Inline mode
- Expand/collapse
- Batch operations

**Run:**
```bash
go run examples/pills/main.go
```

**Key Bindings:**
- `1-6` - Toggle specific pills (Tasks, In Progress, Completed, Errors, Warnings, Info)
- `a` - Expand all
- `x` - Collapse all
- `n` - Add new pill
- `r` - Remove current pill
- `i` - Toggle inline mode
- `q` / `Ctrl+C` - Quit

**Sample Pills:**
- Tasks (5 pending)
- In Progress (2 active)
- Completed (12 done)
- Errors (1 failure)
- Warnings (3 cautions)
- Info (4 notes)

---

### diffview - Diff Viewer

**Location:** `examples/diffview/main.go`

**Difficulty:** ★★★★☆

**Description:** Side-by-side diff viewer with syntax highlighting.

**Key Concepts:**
- Diff parsing (unified/line diff formats)
- Side-by-side comparison
- Syntax highlighting for code diffs
- Line navigation
- Expand/collapse hunks

**Run:**
```bash
go run examples/diffview/main.go
```

**Key Bindings:**
- `↑` / `↓` - Navigate diffs
- `q` / `Ctrl+C` - Quit

---

### diff-demo - Diff Demo

**Location:** `examples/diff-demo/main.go`

**Difficulty:** ★★★★☆

**Description:** Another diff viewer example with different UI style.

**Run:**
```bash
go run examples/diff-demo/main.go
```

---

### reasoning - Reasoning Display

**Location:** `examples/reasoning/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Display AI reasoning steps with expand/collapse.

**Key Concepts:**
- Thought chain visualization
- Nested reasoning steps
- Step navigation
- Expand/collapse reasoning

**Run:**
```bash
go run examples/reasoning/main.go
```

**Key Bindings:**
- `↑` / `↓` - Navigate reasoning steps
- `Enter` - Toggle expansion
- `q` / `Ctrl+C` - Quit

---

### models - Model Selection

**Location:** `examples/models/main.go`

**Difficulty:** ★★★☆☆

**Description:** Dialog for selecting AI models with metadata.

**Key Concepts:**
- Model metadata display
- Search/filter models
- Model selection dialog
- Model compatibility info

**Run:**
```bash
go run examples/models/main.go
```

**Key Bindings:**
- `Ctrl+M` - Open model selection
- Type to filter models
- `Enter` - Select model
- `Esc` - Close

---

### quit - Quit Confirmation

**Location:** `examples/quit/main.go`

**Difficulty:** ★☆☆☆☆

**Description:** Simple quit confirmation dialog.

**Key Concepts:**
- Confirm dialog
- Dialog callback
- Clean application exit

**Run:**
```bash
go run examples/quit/main.go
```

**Key Bindings:**
- `q` / `Ctrl+C` - Show quit confirmation

---

### sessions - Session Management

**Location:** `examples/sessions/main.go`

**Difficulty:** ★★★☆☆

**Description:** Session switching and management UI.

**Key Concepts:**
- Session list
- Session switching dialog
- Session metadata
- Active session indicator

**Run:**
```bash
go run examples/sessions/main.go
```

**Key Bindings:**
- `Ctrl+S` - Open session switcher
- Type to filter sessions
- `Enter` - Switch session

---

### filepicker - File Selection

**Location:** `examples/filepicker/main.go`

**Difficulty:** ★★★☆☆

**Description:** File picker dialog for selecting files.

**Key Concepts:**
- File system navigation in dialog
- File type filtering
- File preview
- Selection confirmation

**Run:**
```bash
go run examples/filepicker/main.go
```

**Key Bindings:**
- `Ctrl+O` - Open file picker
- Navigate directories
- `Enter` - Select file

---

## Layout & Styling ⭐⭐

Examples demonstrating layout systems, visual components, and styling.

### layout-demo - Flex Layout

**Location:** `examples/layout-demo/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Flex-based layout system with dynamic sizing.

**Key Concepts:**
- Flex container
- Flex items with grow/shrink
- Direction (row/column)
- Alignment

**Run:**
```bash
go run examples/layout-demo/main.go
```

**Key Bindings:**
- `q` / `Ctrl+C` - Quit

---

### header-demo - Header Component

**Location:** `examples/header-demo/main.go`

**Difficulty:** ★☆☆☆☆

**Description:** Application header with title and status.

**Key Concepts:**
- Header component
- Title display
- Status indicators
- Breadcrumbs

**Run:**
```bash
go run examples/header-demo/main.go
```

**Key Bindings:**
- `q` / `Ctrl+C` - Quit

---

### sidebar-demo - Sidebar Navigation

**Location:** `examples/sidebar-demo/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Sidebar with navigation items and sections.

**Key Concepts:**
- Sidebar component
- Navigation items
- Active state
- Collapsible sections

**Run:**
```bash
go run examples/sidebar-demo/main.go
```

**Key Bindings:**
- `↑` / `↓` - Navigate sidebar
- `Enter` - Select item
- `q` / `Ctrl+C` - Quit

---

### image - Image Rendering

**Location:** `examples/image/main.go`

**Difficulty:** ★★★☆☆

**Description:** Terminal image rendering with scaling.

**Key Concepts:**
- Image rendering in terminal
- Scaling and resizing
- ANSI art conversion
- Cache optimization

**Run:**
```bash
go run examples/image/main.go
```

**Key Bindings:**
- `+` / `-` - Scale image
- `q` / `Ctrl+C` - Quit

---

### markdown-demo - Markdown Rendering

**Location:** `examples/markdown-demo/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Markdown document rendering with styles.

**Key Concepts:**
- Markdown parsing
- Styled rendering
- Code block formatting
- Syntax highlighting

**Run:**
```bash
go run examples/markdown-demo/main.go
```

---

### autocomplete - Autocomplete Demo

**Location:** `examples/autocomplete/demo.go`

**Difficulty:** ★★☆☆☆

**Description:** Autocomplete dropdown component.

**Key Concepts:**
- Autocomplete popup
- Fuzzy matching
- Keyboard navigation
- Item selection

**Run:**
```bash
go run examples/autocomplete/demo.go
```

---

### notify - Notifications

**Location:** `examples/notify/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Toast notification system.

**Key Concepts:**
- Notification queue
- Fade in/out animations
- Multiple notification types (success, error, warning, info)
- Auto-dismiss

**Run:**
```bash
go run examples/notify/main.go
```

**Key Bindings:**
- `q` / `Ctrl+C` - Quit

---

### tasks-demo - Task List

**Location:** `examples/tasks-demo/main.go`

**Difficulty:** ★★☆☆☆

**Description:** Todo/task management list.

**Key Concepts:**
- Task items
- Task completion
- Priority indicators
- Group by status

**Run:**
```bash
go run examples/tasks-demo/main.go
```

**Key Bindings:**
- `↑` / `↓` - Navigate tasks
- `Space` - Toggle completion
- `q` / `Ctrl+C` - Quit

---

## Engine & Architecture ⭐⭐⭐

Examples demonstrating rendering engines and architectural concepts.

### dual-engine - Dual Engine Demo

**Location:** `examples/dual-engine/main.go`

**Difficulty:** ★★★★☆

**Description:** Demonstrates switching between engines at runtime.

**Key Concepts:**
- Engine switching
- Model persistence across engines
- Engine comparison
- Performance metrics

**Run:**
```bash
go run examples/dual-engine/main.go
```

**Key Bindings:**
- `1` - Switch to Bubbletea engine
- `2` - Switch to Ultraviolet engine
- `q` / `Ctrl+C` - Quit

---

## Running Examples

### Prerequisites

```bash
# Clone the repository
git clone https://github.com/wwsheng009/taproot.git
cd taproot

# Install dependencies
go mod download
```

### Running Individual Examples

```bash
# Basic example
go run examples/demo/main.go

# With custom flags (if the example supports)
go run examples/demo/main.go --theme dark
```

### Building Examples

```bash
# Build for current platform
go build -o demo.exe examples/demo/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o demo-linux examples/demo/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o demo-mac examples/demo/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./ui/list/...
go test ./ui/components/attachments/...
go test ./ui/dialog/...

# Run tests with coverage
go test -cover ./...
```

---

## Contributing Examples

To contribute a new example:

1. Create a new directory under `examples/`
2. Name it descriptively (lowercase, underscores for spaces)
3. Add a `main.go` file with a complete working example
4. Add documentation (this file, comments in code)
5. Update the examples table in `PROGRESS.md`
6. Follow the example structure:
   - Clear model implementation
   - Commented key bindings
   - Usage instructions
   - Dependencies listed

---

## Example Structure Template

```go
package main

import (
    "fmt"
    "github.com/wwsheng009/taproot/ui/render"
    "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
    // Fields...
}

// Init initializes the model
func (m Model) Init() render.Cmd {
    return nil
}

// Update handles incoming messages
func (m Model) Update(msg any) (render.Model, render.Cmd) {
    switch msg := msg.(type) {
    case render.KeyMsg:
        switch msg.Key {
        case "q":
            return m, render.Quit()
        }
    }
    return m, render.None()
}

// View returns the UI representation
func (m Model) View() string {
    return "Hello, Taproot!"
}

func main() {
    p := tea.NewProgram(Model{}, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        panic(err)
    }
}
```

---

## Best Practices for Examples

1. **Keep it simple**: Focus on one concept per example
2. **Add comments**: Explain what each part does
3. **Document key bindings**: List all keyboard shortcuts
4. **Handle errors**: Provide proper error messages
5. **Use render.Model**: Follow v2.0 conventions
6. **Add description**: README or top-level comment explaining the example
7. **Test thoroughly**: Ensure the example runs without errors
8. **Keep dependencies minimal**: Only import what's needed

---

## Common Patterns

### Pattern 1: Simple Counter

See: `examples/demo/main.go`

### Pattern 2: List with Selection

See: `examples/list/main.go`, `examples/ui-list/main.go`

### Pattern 3: Filtered List

See: `examples/filterablelist/main.go`, `examples/ui-filtergroup/main.go`

### Pattern 4: File Browser

See: `examples/files/main.go`, `examples/treefiles/main.go`

### Pattern 5: Dialog System

See: `examples/quit/main.go`, `examples/ui-dialogs/main.go`

### Pattern 6: Form with Validation

See: `examples/forms/main.go`

### Pattern 7: Progress Display

See: `examples/progress/main.go`

### Pattern 8: Page Navigation

See: `examples/app/main.go`

---

## Troubleshooting

### Issue: Example won't run

**Solution:** Ensure you're in the correct directory:
```bash
cd taproot  # Make sure you're in the project root
go run examples/demo/main.go
```

### Issue: Import errors

**Solution:** Update dependencies:
```bash
go mod tidy
go mod download
```

### Issue: Terminal display issues

**Solution:** Try with different terminal settings:
```bash
# Run without alternate screen
go run examples/demo/main.go

# Or set terminal to UTF-8
export LANG=en_US.UTF-8
```

### Issue: Keys not responding

**Solution:** Check that your terminal supports the key codes. Some keys may need different mappings in different terminals.

---

**Last Updated:** Phase 10 completion (February 2026)

**Total Examples:** 37 interactive programs
