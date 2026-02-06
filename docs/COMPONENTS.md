# Taproot Component Reference

Complete reference for all components in the Taproot TUI framework.

## Table of Contents

- [Core Interfaces](#core-interfaces)
- [Dialog System](#dialog-system)
- [Form Components](#form-components)
- [List Components](#list-components)
- [Message Components](#message-components)
- [Status Components](#status-components)
- [Progress Components](#progress-components)
- [Other Components](#other-components)
- [Tools](#tools)

---

## Core Interfaces

### Focusable

Components that can receive and lose keyboard focus.

```go
type Focusable interface {
    Focus()
    Blur()
    Focused() bool
}
```

### Sizeable

Components with dimensions.

```go
type Sizeable interface {
    Size() (width, height int)
    SetSize(width, height int)
}
```

### Positional

Components with x, y coordinates.

```go
type Positional interface {
    Position() (x, y int)
    SetPosition(x, y int)
}
```

### Help

Components that provide help text.

```go
type Help interface {
    Help() []string
}
```

---

## Dialog System

Location: `ui/dialog/`

The dialog system provides modal dialogs for user interaction. All dialogs support:
- Overlay rendering with automatic centering
- Keyboard navigation
- Callback-based result handling
- Stack management for multiple dialogs

### Dialog Types

#### InfoDialog

Simple informational message dialog.

```go
import "github.com/wwsheng009/taproot/ui/dialog"

dialog := dialog.NewInfoDialog("Success", "Operation completed!")
dialog.SetCallback(func(result dialog.ActionResult, data any) {
    // Handle acknowledgment
})
```

#### ConfirmDialog

Yes/No confirmation dialog with callback.

```go
confirm := dialog.NewConfirmDialog(
    "Delete File",
    "Are you sure you want to delete this file?",
    func(result dialog.ActionResult, data any) {
        if result == dialog.ActionConfirm {
            // Execute delete
        }
    },
)
```

#### InputDialog

Text input prompt with validation.

```go
input := dialog.NewInputDialog(
    "Enter Name",
    "Name:",
    func(value string) {
        fmt.Printf("User entered: %s\n", value)
    },
)
input.SetPlaceholder("John Doe")
input.SetMaxLength(50)
```

#### SelectListDialog

Single-select from list of items.

```go
items := []string{"Option 1", "Option 2", "Option 3"}
select := dialog.NewSelectListDialog(
    "Choose an Option",
    items,
    func(index int, value string) {
        fmt.Printf("Selected: %s\n", value)
    },
)
```

### Overlay Manager

Stack management for multiple dialogs.

```go
overlay := dialog.NewOverlay()
overlay.Push(confirmDialog)
overlay.Pop()
active := overlay.Peek()
```

---

## Form Components

Location: `ui/forms/`

Form components provide input handling with validation and focus management.

### TextInput

Single-line text input with validation.

```go
import "github.com/wwsheng009/taproot/ui/forms"

input := forms.NewTextInput("name", "Name")
input.SetRequired(true)
input.SetPlaceholder("Enter your name")

// Custom validation
input.SetValidation(func(s string) error {
    if len(s) < 3 {
        return errors.New("name too short")
    }
    return nil
})

value := input.Value()
```

### TextArea

Multi-line text input with word wrap.

```go
textarea := forms.NewTextArea("description", "Description")
textarea.SetPlaceholder("Enter description...")
textarea.SetMinSize(60, 5)
```

### Select

Dropdown selection with expand/collapse.

```go
select := forms.NewSelect("country", "Country")
select.SetOptions([]string{"USA", "Canada", "UK"})
select.SetValue("USA")
```

### Checkbox

Boolean toggle selection.

```go
checkbox := forms.NewCheckbox("agree", "I agree to terms")
checkbox.SetChecked(true)
if checkbox.Checked() {
    // User agreed
}
```

### Radio

Single selection from options.

```go
radio := forms.NewRadio("plan", "Plan")
radio.SetOptions([]string{"Free", "Pro", "Enterprise"})
radio.SetValue("Pro")
```

### Form Container

Manages multiple inputs with focus traversal.

```go
form := forms.NewForm(nameInput, emailInput, checkbox)

// In Update loop
newForm, cmd := form.Update(msg)
form = newForm.(*forms.Form)

// Validate before submission
if err := form.Validate(); err != nil {
    // Show error
}
```

---

## List Components

Location: `ui/list/`

Virtualized list components for large datasets with filtering and grouping.

### Core Interfaces

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

// Selectable: Item with selection state
type Selectable interface {
    Selected() bool
    SetSelected(bool)
}

// Toggleable: Item with expand/collapse state
type Toggleable interface {
    Expanded() bool
    SetExpanded(bool)
}
```

### ListItem

Basic filterable item implementation.

```go
import "github.com/wwsheng009/taproot/ui/list"

item := list.NewListItem("1", "Apple", "Red fruit")
fmt.Println(item.ID())         // "1"
fmt.Println(item.FilterValue()) // "Apple Red fruit"
```

### Viewport

Virtualized scrolling manager.

```go
viewport := list.NewViewport(10, 100) // visible=10, total=100

viewport.MoveUp()
viewport.MoveDown()
viewport.PageUp()
viewport.PageDown()
viewport.MoveToTop()
viewport.MoveToBottom()

start, end := viewport.Range()
```

### Filter

Search functionality with match highlighting.

```go
filter := list.NewFilter()
filter.SetQuery("apple")
filtered := filter.Apply(items)

// Highlight matches
highlighted := filter.Highlight("apple pineapple", "[", "]")
// Result: "[apple] pineapple"
```

### SelectionManager

Multi-selection state management.

```go
selMgr := list.NewSelectionManager(list.SelectionModeMultiple)

selMgr.Select(item)
selMgr.Deselect(item)
selMgr.Toggle(item)
selMgr.SelectAll(items)
selMgr.DeselectAll()
selMgr.InvertSelection(items)

selectedIDs := selMgr.SelectedIDs()
selected := selMgr.GetSelected(items)
```

### GroupManager

Grouped/expanding list support.

```go
groupMgr := list.NewGroupManager()

group := list.Group{
    Title:    "Fruits",
    Expanded: true,
    Items:    items,
}

groupMgr.AddGroup(group)
groupMgr.ExpandAll()
groupMgr.CollapseAll()
groupMgr.ToggleCurrentGroup(cursor)
```

### BaseList

Complete list component with key bindings.

```go
baseList := list.NewBaseList()
baseList.SetItems(items)
baseList.SetFilter("query")

// Update with key message
action := list.DefaultKeyMap().MatchAction("down")
newList := baseList.HandleAction(action)

view := baseList.View()
```

---

## Message Components

Location: `ui/components/messages/`

Components for displaying chat messages, diagnostics, and todos.

### Message Types

#### AssistantMessage

AI assistant message with markdown rendering.

```go
import "github.com/wwsheng009/taproot/ui/components/messages"

assistant := messages.NewAssistantMessage(
    "msg-1",
    "This is **markdown** with `code` blocks",
)
assistant.SetInputTokens(100)
assistant.SetOutputTokens(200)
assistant.SetExpanded(true)
```

#### UserMessage

User message with code blocks and attachments.

```go
user := messages.NewUserMessage("msg-2", "Here's a code snippet")
user.AddCodeBlock("go", "func main() { println(\"Hello\") }")
user.SetFileAttachments([]string{"file1.go", "file2.go"})
user.SetFilesAdded(2)
```

#### ToolMessage

Tool call/result display.

```go
tool := messages.NewToolMessage(
    "msg-3",
    "bash",
    []string{"run", "build"},
    `Build started...`,
)
tool.SetResult(`Build completed in 1.2s`)
```

#### FetchMessage

Network request display with 4 types.

```go
// Basic fetch
basic := messages.NewFetchMessage("msg-4", messages.FetchTypeBasic, "Fast fetch")
basic.SetRequest(&messages.FetchRequest{URL: "https://example.com"})
basic.SetResult(&messages.FetchResult{Success: true, Content: `{"data": "value"}`})

// Web search
search := messages.NewFetchMessage("msg-5", messages.FetchTypeWebSearch, "Search")
search.SetRequest(&messages.FetchRequest{Query: "golang tutorial"})

// Agentic (multi-step with nested messages)
agentic := messages.NewFetchMessage("msg-6", messages.FetchTypeAgentic, "Agentic")
agentic.AddNestedMessage(toolMessage1)
agentic.AddNestedMessage(toolMessage2)
```

#### DiagnosticMessage

LSP-style diagnostic display.

```go
diag := messages.NewDiagnosticMessage(
    "msg-7",
    "compiler",
    messages.SeverityError,
    "invalid syntax",
    42,
    5,
)
diag.SetMessage("unexpected token in expression")
diag.SetCode(`func main() { x = 5 + }`)
```

#### TodoMessage

Task list with progress tracking.

```go
todo := messages.NewTodoMessage("msg-8", "task-1", "Fix bug", messages.StatusPending)
todo.SetProgress(1, 5)
todo.SetDescription("Track down the issue")

todo.AddTodo("task-2", "Write tests", messages.StatusInProgress)
todo.AddTodo("task-3", "Update docs", messages.StatusCompleted)
```

---

## Status Components

Location: `ui/components/status/`

Service status and diagnostic display components.

### Service States

```go
const (
    StateDisabled State = iota // ○ Disabled/inactive
    StateStarting               // ⟳ Starting up
    StateReady                  // ● Connected and ready
    StateError                  // × Error state
)
```

### ServiceCmp

Single service status display.

```go
import "github.com/wwsheng009/taproot/ui/components/status"

service := status.NewService("gopls", "Go LSP")
service.SetStatus(status.StateReady)
service.SetErrorCount(3)
service.SetCompact(true)

view := service.View()
```

### DiagnosticStatusCmp

Diagnostic summary with severity counts.

```go
diag := status.NewDiagnosticStatus("workspace")
diag.AddDiagnostic(status.DiagnosticSeverityError)
diag.AddDiagnostic(status.DiagnosticSeverityWarning)
diag.SetCompact(true)
```

### LSPList

Multiple LSP services display.

```go
lspList := status.NewLSPList()
lspList.SetWidth(50)
lspList.SetMaxItems(5)

lspList.AddService(status.LSPServiceInfo{
    Name:     "gopls",
    Language: "go",
    State:    status.StateReady,
    Diagnostics: status.DiagnosticSummary{
        Error:   0,
        Warning: 2,
    },
})

onlineCount := lspList.OnlineCount()
totalErrors := lspList.TotalErrors()
```

### MCPList

MCP (Model Context Protocol) services with tool counts.

```go
mcpList := status.NewMCPList()
mcpList.AddService(status.MCPServiceInfo{
    Name:  "filesystem",
    State: status.StateReady,
    ToolCounts: status.ToolCounts{
        Tools:   5,
        Prompts: 0,
    },
})

connectedCount := mcpList.ConnectedCount()
totalTools := mcpList.TotalTools()
```

---

## Progress Components

Location: `ui/components/progress/`

### ProgressBar

Progress bar with percentage display.

```go
import "github.com/wwsheng009/taproot/ui/components/progress"

bar := progress.NewProgressBar()
bar.SetLabel("Downloading")
bar.SetTotal(100)
bar.SetCurrent(75)
bar.SetWidth(40)

completed := bar.Completed()   // false
percent := bar.Percent()       // 75.0

view := bar.View()
// Output: ████████████████████████░░░░░░░░░░ Downloading 75/100 (75%)
```

### Spinner

Animated loading indicator.

```go
spinner := progress.NewSpinner(progress.SpinnerDots)
spinner.SetFPS(30)
spinner.SetLabel("Loading...")
spinner.Start()

view := spinner.View()
spinner.Stop()
```

---

## Other Components

### Attachments

Location: `ui/components/attachments/`

File attachment list with preview.

```go
import "github.com/wwsheng009/taproot/ui/components/attachments"

list := attachments.NewAttachmentList()
list.Add(attachments.Attachment{
    ID:   "att-1",
    Name: "document.pdf",
    Path: "/path/to/document.pdf",
    Size: 2400000,
})

list.SetCompact(true)
list.SetFilter("pdf")
```

### Pills

Location: `ui/components/pills/`

Status/queue pills for metadata display.

```go
import "github.com/wwsheng009/taproot/ui/components/pills"

pills := pills.NewPillList()
pills.Add(pills.Pill{
    ID:     "pill-1",
    Label:  "Tasks",
    Status: pills.StatusPending,
    Count:  5,
})

pills.SetInlineMode(true)
pills.ExpandAll()
```

### Sidebar

Location: `ui/components/sidebar/`

Navigation sidebar with sections.

```go
import "github.com/wwsheng009/taproot/ui/components/sidebar"

sb := sidebar.NewSidebar()
sb.SetWidth(20)
sb.AddSection("Files", items)
sb.AddSection("Edit", items)
```

---

## Tools

### Clipboard

Location: `ui/tools/clipboard/`

Cross-platform clipboard with OSC 52 and native support.

```go
import "github.com/wwsheng009/taproot/ui/tools/clipboard"

mgr := clipboard.NewDefaultManager()

// Copy/paste
mgr.Copy("Hello, clipboard!")
text, _ := mgr.Paste()

// History
history := mgr.History()
mgr.RestoreFromHistory(0)

// Provider switching
mgr.SetProvider(clipboard.ClipboardOSC52)
mgr.SetProvider(clipboard.ClipboardNative)
```

### Shell

Location: `ui/tools/shell/`

Shell command execution with async support.

```go
import "github.com/wwsheng009/taproot/ui/tools/shell"

executor := shell.NewExecutor()

// Basic execution
result, err := executor.Execute("echo", []string{"hello"})

// Command builder
cmd, args, opts, _ := shell.NewCommandBuilder().
    ShellCommand("ls -la").
    SetTimeout(5 * time.Second).
    Build()

result, err = executor.ExecuteWithOptions(cmd, args, opts)

// Async execution
process, _ := executor.ExecuteAsync("long-running", []string{"task"}, opts)
err = executor.Cancel(process.ID())

// Pipe commands
result, err = executor.Pipe(
    []string{"cat", "file.txt"},
    []string{"grep", "pattern"},
)
```

### Watcher

Location: `ui/tools/watcher/`

File system monitoring with debouncing.

```go
import "github.com/wwsheng009/taproot/ui/tools/watcher"

w, _ := watcher.NewWatcher(
    watcher.Filter{
        IncludePatterns: []string{"*.go"},
        Recursive:       true,
    },
    func(events []watcher.Event) {
        for _, e := range events {
            fmt.Printf("%s: %s\n", e.Type.String(), e.Path)
        }
    },
    func(err error) {
        log.Printf("Error: %v", err)
    },
)

w.AddRecursive("/path/to/project")
w.Start()

// Configure debouncing
w.SetDebounceConfig(watcher.DebounceConfig{
    Enabled:     true,
    Delay:       100 * time.Millisecond,
    MaxWait:     500 * time.Millisecond,
    MergeEvents: true,
})
```

---

## Using Components Together

### Example: Complete Form with Dialog

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/wwsheng009/taproot/ui/dialog"
    "github.com/wwsheng009/taproot/ui/forms"
    "github.com/wwsheng009/taproot/ui/styles"
)

type Model struct {
    form forms.Form
    quit bool
}

func NewModel() Model {
    name := forms.NewTextInput("name", "Name")
    name.SetRequired(true)

    email := forms.NewTextInput("email", "Email")
    email.SetValidation(validateEmail)

    agree := forms.NewCheckbox("agree", "I agree")

    return Model{
        form: forms.NewForm(name, email, agree),
    }
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            m.quit = true
            return m, tea.Quit
        case "enter":
            if err := m.form.Validate(); err != nil {
                // Show error dialog
                return m, func() tea.Msg {
                    return dialog.OpenDialogMsg{
                        Model: dialog.NewInfoDialog("Error", err.Error()),
                    }
                }
            }
            // Submit form
        }
    }

    // Update form
    newForm, cmd := m.form.Update(msg)
    m.form = newForm.(forms.Form)
    return m, cmd
}

func (m Model) View() string {
    s := styles.DefaultStyles()
    return s.S().Title.Render("User Registration") + "\n\n" + m.form.View()
}

func validateEmail(s string) error {
    if !strings.Contains(s, "@") {
        return errors.New("invalid email")
    }
    return nil
}
```

---

For more details, see:
- [API Reference](API.md)
- [Architecture](ARCHITECTURE.md)
- [Examples](../examples/)
