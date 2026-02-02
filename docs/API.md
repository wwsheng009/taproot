# Taproot API Documentation

Complete API reference for Taproot TUI framework components.

## Table of Contents

- [Core Framework](#core-framework)
- [Application Layer](#application-layer)
- [Components](#components)
- [Utilities](#utilities)
- [Theme System](#theme-system)

---

## Core Framework

### Layout Interfaces (`internal/layout`)

#### Focusable

Components that can receive and lose focus.

```go
type Focusable interface {
    Focus()
    Blur()
    Focused() bool
}
```

**Methods:**
- `Focus()` - Sets the component to focused state
- `Blur()` - Removes focus from the component
- `Focused() bool` - Returns true if component is focused

#### Sizeable

Components with dimensions.

```go
type Sizeable interface {
    Size() (width, height int)
    SetSize(width, height int)
}
```

#### Positional

Components with x, y coordinates.

```go
type Positional interface {
    Position() (x, y int)
    SetPosition(x, y int)
}
```

#### Help

Components that provide help text.

```go
type Help interface {
    Help() []string
}
```

---

## Application Layer

### App Model (`internal/tui/app`)

Main application model that manages pages and dialogs.

```go
type AppModel struct {
    // Contains page management, dialog stack, global state
}

func NewApp() *AppModel
func (a *AppModel) RegisterPage(id page.PageID, model util.Model)
func (a *AppModel) SetPage(id page.PageID) tea.Cmd
func (a *AppModel) Init() tea.Cmd
func (a *AppModel) Update(msg tea.Msg) (util.Model, tea.Cmd)
func (a *AppModel) View() string
```

**Usage:**
```go
app := app.NewApp()
app.RegisterPage("home", HomePage{})
app.SetPage("home")

p := tea.NewProgram(app, tea.WithAltScreen())
p.Run()
```

### Page System (`internal/tui/page`)

Page identifiers and messages.

```go
type PageID string

type PageChangeMsg struct {
    ID PageID
}

type PageCloseMsg struct{}
type PageBackMsg struct{}
```

**Usage:**
```go
// Change to a different page
return func() tea.Msg {
    return page.PageChangeMsg{ID: "settings"}
}

// Go back to previous page
return func() tea.Msg {
    return page.PageBackMsg{}
}
```

---

## Components

### Dialog System (`internal/tui/components/dialogs`)

#### Dialog Model Interface

```go
type DialogModel interface {
    util.Model
    Position() (int, int)
    ID() DialogID
}

type DialogID string
```

#### Dialog Messages

```go
type OpenDialogMsg struct {
    Model DialogModel
}

type CloseDialogMsg struct{}
```

**Opening a Dialog:**
```go
dialog := quit.New(hasChanges)
return func() tea.Msg {
    return dialogs.OpenDialogMsg{Model: dialog}
}
```

**Closing a Dialog:**
```go
return func() tea.Msg {
    return dialogs.CloseDialogMsg{}
}
```

### Command Palette (`internal/tui/components/dialogs/commands`)

Fuzzy-searchable command palette with argument collection.

```go
type Command struct {
    ID          string
    Title       string
    Description string
    Args        []ArgDef
    Callback    func(args map[string]string) tea.Cmd
}

type CommandProvider interface {
    Commands() []Command
}

func NewCommandsDialog(provider CommandProvider) DialogModel
```

**Usage:**
```go
provider := &MyCommandProvider{}
dialog := commands.NewCommandsDialog(provider)
return func() tea.Msg {
    return dialogs.OpenDialogMsg{Model: dialog}
}
```

### Model Selection (`internal/tui/components/dialogs/models`)

Model selection dialog with search and recent models.

```go
type Model struct {
    ID          string
    Name        string
    Provider    string
    ContextSize int
    Description string
}

type ModelProvider interface {
    Models() []Model
    RecentModels() []Model
    SetModel(modelID string) tea.Cmd
}

func NewModelsDialog(provider ModelProvider) DialogModel
```

### Session Management (`internal/tui/components/dialogs/sessions`)

Session management with CRUD operations.

```go
type Session struct {
    ID        string
    Title     string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type SessionProvider interface {
    Sessions() []Session
    CreateSession(title string) tea.Cmd
    DeleteSession(id string) tea.Cmd
    SelectSession(id string) tea.Cmd
}

func NewSessionsDialog(provider SessionProvider) DialogModel
```

### Messages Display (`internal/tui/components/messages`)

Chat messages list with scrolling and tool call display.

```go
type Message struct {
    Role    string // "user", "assistant", "system", "tool"
    Content string
    ToolUse *ToolUse
}

type ToolUse struct {
    Name      string
    Arguments string
    Result    string
}

type MessagesModel struct {
    // Internal state
}

func New() *MessagesModel
func (m *MessagesModel) SetWidth(w int)
func (m *MessagesModel) SetHeight(h int)
func (m *MessagesModel) AddMessage(msg Message) tea.Cmd
func (m *MessagesModel) Clear()
```

**Usage:**
```go
messages := messages.New()
messages.SetWidth(80)
messages.SetHeight(20)

messages.AddMessage(messages.Message{
    Role:    "user",
    Content: "Hello!",
})
```

### Lists (`internal/tui/exp/list`)

Virtualized list components with filtering and grouping.

#### Simple List

```go
type SimpleList struct {
    // Internal state
}

func NewSimpleList() *SimpleList
func (l *SimpleList) SetItems(items []string)
func (l *SimpleList) Cursor() int
func (l *SimpleList) SetCursor(cursor int)
```

#### Filterable List

```go
type FilterableList struct {
    // Internal state
}

func NewFilterableList() *FilterableList
func (l *FilterableList) SetFilter(filter string)
func (l *FilterableList) FilteredItems() []string
```

#### Grouped List

```go
type GroupedList struct {
    // Internal state
}

type GroupedItem struct {
    Group   string
    Content string
}

func NewGroupedList() *GroupedList
func (l *GroupedList) SetGroups(groups map[string][]string)
func (l *GroupedList) ToggleGroup(group string)
```

### Diff Viewer (`internal/tui/exp/diffview`)

Unified diff viewer with scrolling.

```go
type DiffView struct {
    // Internal state
}

func NewDiffView() *DiffView
func (d *DiffView) SetDiff(unifiedDiff string)
func (d *DiffView) SetScroll(scroll int)
```

### Image Renderer (`internal/tui/components/image`)

Enhanced terminal image display with **multi-protocol support** and **automatic capability detection**.

**Features:**
- **Sixel Protocol**: High-quality rendering for Sixel-compatible terminals (XTerm, WezTerm, Mintty)
- **Kitty/iTerm2**: Native protocol support for modern terminals
- **Unicode Blocks**: Universal fallback with ANSI colors
- **ASCII Art**: Full compatibility fallback
- **Auto Detection**: Automatically queries terminal capabilities using DA query
- **Tmux Support**: Automatic passthrough sequence wrapping for tmux sessions
- **Zoom Control**: Zoom in/out with aspect ratio preservation
- **Smart Caching**: Protocol-level caching for optimal performance

```go
type RendererType int

const (
    RendererAuto RendererType = iota // Auto-detect best renderer
    RendererSixel                    // Sixel protocol (high-quality)
    RendererKitty                    // Kitty graphics protocol
    RendereriTerm2                   // iTerm2 inline images
    RendererBlocks                   // Unicode blocks with ANSI colors
    RendererASCII                    // ASCII art fallback
)

type ZoomMode int

const (
    ZoomFitScreen ZoomMode = iota // Fit within space (maintain aspect ratio)
    ZoomFitWidth                  // Fit width, allow height overflow
    ZoomFitHeight                 // Fit height, allow width overflow
    ZoomFill                      // Fill space (may distort)
)

type Image struct {
    // Internal state
    width   int
    height  int
    path    string
    renderer RendererType
    loaded  bool
    img     image.Image
    scaled  image.Image
    scale   float64      // Zoom scale factor (1.0 = fit to screen)
    zoomMode ZoomMode     // Current zoom mode
    // Caches
    cachedView   string
    sixelCache   string
    asciiCache   string
}

// Image loading
func New(path string) *Image
func (img *Image) SetPath(path string) tea.Cmd
func (img *Image) Reload() tea.Cmd
func (img *Image) IsLoaded() bool
func (img *Image) Error() string

// Rendering
func (img *Image) SetRenderer(renderer RendererType) tea.Cmd
func (img *Image) SetSize(w, h int)
func (img *Image) View() string

// Zoom control
func (img *Image) ZoomIn() tea.Cmd
func (img *Image) ZoomOut() tea.Cmd
func (img *Image) ResetZoom() tea.Cmd
func (img *Image) GetScale() float64
func (img *Image) SetZoomMode(mode ZoomMode) tea.Cmd
func (img *Image) CycleZoomMode() tea.Cmd
func (img *Image) GetZoomModeName() string

// Query methods
func (img *Image) Path() string
func (img *Image) Size() (int, int)
func (img *Image) ScaledSize() (int, int)
```

**Usage Example:**

```go
import "github.com/wwsheng009/taproot/tui/components/image"

// Create image with auto-detection
img := image.New("/path/to/image.png")

// Set size and enable auto-scaling
img.SetSize(80, 40)

// Manually set renderer (optional)
img.SetRenderer(image.RendererSixel)

// Zoom in by 25%
cmd := img.ZoomIn()

// Cycle zoom modes: FitScreen -> FitWidth -> FitHeight -> Fill
cmd = img.CycleZoomMode()
```

---

## Utilities

### Model Interface (`internal/tui/util`)

Base interface for all Bubbletea models.

```go
type Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (Model, tea.Cmd)
    View() string
}
```

### Info Messages (`internal/tui/util`)

Status and notification messages.

```go
type InfoType int

const (
    InfoTypeInfo InfoType = iota
    InfoTypeSuccess
    InfoTypeWarn
    InfoTypeError
    InfoTypeUpdate
)

type InfoMsg struct {
    Type InfoType
    Msg  string
}
```

**Creating Messages:**
```go
// Convenience functions
util.NewInfoMsg("Info message")
util.NewSuccessMsg("Success!")
util.NewWarnMsg("Warning")
util.NewErrorMsg(err)

// Reporting commands
util.ReportInfo("Info")
util.ReportSuccess("Success")
util.ReportWarn("Warning")
util.ReportError(err)
```

### Commands

```go
func CmdHandler(msg tea.Msg) tea.Cmd
```

---

## Theme System

### Theme Structure (`internal/ui/styles`)

```go
type Theme struct {
    // Brand colors
    Primary   lipgloss.Color
    Secondary lipgloss.Color
    Tertiary  lipgloss.Color
    Accent    lipgloss.Color
    
    // Backgrounds
    BgBase       lipgloss.Color
    BgBaseLighter lipgloss.Color
    BgSubtle     lipgloss.Color
    BgOverlay    lipgloss.Color
    
    // Foregrounds
    FgBase      lipgloss.Color
    FgMuted     lipgloss.Color
    FgHalfMuted lipgloss.Color
    FgSubtle    lipgloss.Color
    FgSelected  lipgloss.Color
    
    // Borders
    Border      lipgloss.Color
    BorderFocus lipgloss.Color
    
    // Status
    Success lipgloss.Color
    Error   lipgloss.Color
    Warning lipgloss.Color
    Info    lipgloss.Color
    
    // Style presets
    S *StyleSet
}

type StyleSet struct {
    Base      lipgloss.Style
    Primary   lipgloss.Style
    Secondary lipgloss.Style
    // ... more presets
}
```

### Theme Functions

```go
func CurrentTheme() *Theme
func ApplyForegroundGrad(text string, c1, c2 lipgloss.Color) string
func ApplyBoldForegroundGrad(text string, c1, c2 lipgloss.Color) string
func blendColors(c1, c2 lipgloss.Color, t float64) lipgloss.Color
```

**Usage:**
```go
t := styles.CurrentTheme()

// Use style presets
text := t.S().Base.Foreground(t.Primary).Render("Hello")

// Apply gradients
gradient := styles.ApplyForegroundGrad("Gradient Text", t.Primary, t.Secondary)
```

### Animation (`internal/tui/anim`)

Animated spinner component.

```go
type Anim struct {
    // Internal state
}

type Settings struct {
    Label       string
    LabelColor  lipgloss.Color
    GradColorA  lipgloss.Color
    GradColorB  lipgloss.Color
    CycleColors bool
}

func New(opts Settings) *Anim
```

---

## Keyboard Shortcuts

### Default Key Bindings

| Key | Action |
|-----|--------|
| `ctrl+c` | Quit application |
| `ctrl+g` | Toggle help |
| `ESC` | Close dialog / Go back |
| `ctrl+p` | Open command palette |
| `ctrl+m` | Open model selector |
| `ctrl+s` | Open session selector |

### Navigation

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `h` / `←` | Move left |
| `l` / `→` | Move right |
| `g` | Jump to top |
| `G` | Jump to bottom |

---

## Best Practices

1. **Always use interfaces** - Compose through interfaces like `Focusable`, `Sizeable`
2. **Use the theme system** - Never hardcode colors
3. **Handle window resizing** - Implement `tea.WindowSizeMsg`
4. **Return nil commands** - Not every update needs to return a command
5. **Test your models** - Write tests for model logic
6. **Keep views simple** - Use `strings.Builder` for complex views

---

## Migration Guide

See [MIGRATION_PLAN.md](MIGRATION_PLAN.md) for detailed migration information.

---

## v2.0 Component Reference

The v2.0 architecture introduces engine-agnostic components that work with multiple rendering engines (Bubbletea, Ultraviolet, Direct).

### Render Abstraction (`ui/render`)

#### Model Interface

Core interface that all v2.0 components implement.

```go
type Model interface {
    Init() Cmd
    Update(msg any) (Model, Cmd)
    View() string
}
```

#### Command Interface

Side effect abstraction.

```go
type Cmd interface {
    Execute() error
}

// Built-in commands
func None() Cmd                           // No operation
func Quit() Cmd                           // Exit application
func Batch(cmds ...Cmd) Cmd               // Combine multiple commands
func Tick(interval time.Duration, fn func(time.Time) Msg) Cmd
```

#### Messages

```go
type KeyMsg struct {
    Key      string  // "up", "down", "enter", "q", etc.
    Runes    []rune  // Raw runes for typed characters
    Alt      bool    // Alt modifier
    Ctrl     bool    // Ctrl modifier
}

type WindowSizeMsg struct {
    Width  int
    Height int
}

type TickMsg struct {
    Time time.Time
}

type QuitMsg struct{}
```

#### Engine Management

```go
type EngineType int

const (
    EngineBubbletea EngineType = iota
    EngineUltraviolet
    EngineDirect
)

type Engine interface {
    Type() EngineType
    Start(model Model) error
    Stop() error
    Send(msg Msg) error
    Resize(width, height int) error
    Running() bool
}

// Create an engine
engine, err := render.CreateEngine(render.EngineBubbletea, render.DefaultConfig())
engine.Start(initialModel)
```

---

### List Components (`ui/list`)

#### Core Interfaces

```go
// Item: Base interface for all list items
type Item interface {
    ID() string
}

// FilterableItem: Items that can be filtered
type FilterableItem interface {
    Item
    FilterValue() string
}

// Selectable: Items with selection state
type Selectable interface {
    Selected() bool
    SetSelected(bool)
}

// Toggleable: Items with expand/collapse state
type Toggleable interface {
    Expanded() bool
    SetExpanded(bool)
}
```

#### ListItem

Basic filterable item implementation.

```go
type ListItem struct {
    ID          string
    Title       string
    Description string
}

func NewListItem(id, title, description string) ListItem
func (l ListItem) FilterValue() string

// Usage
item := list.NewListItem("1", "Apple", "Red fruit")
```

#### Viewport

Virtualized scrolling manager.

```go
type Viewport struct {
    visible  int
    total    int
    position int
}

func NewViewport(visible, total int) Viewport
func (v *Viewport) MoveUp()
func (v *Viewport) MoveDown()
func (v *Viewport) PageUp()
func (v *Viewport) PageDown()
func (v *Viewport) MoveToTop()
func (v *Viewport) MoveToBottom()
func (v *Viewport) Range() (start, end int)
func (v *Viewport) ScrollIndicator() string

// Usage
viewport := list.NewViewport(10, 100)
viewport.MoveUp()
visibleItems := items[viewport.Range()]
```

#### Filter

Search functionality with highlighting.

```go
type Filter struct {
    query          string
    caseSensitive  bool
}

func NewFilter() Filter
func (f *Filter) SetQuery(query string)
func (f *Filter) Clear()
func (f *Filter) Apply(items []FilterableItem) []FilterableItem
func (f *Filter) SetCaseSensitive(sensitive bool)
func (f *Filter) Highlight(text, before, after string) string

// Usage
filter := list.NewFilter()
filter.SetQuery("apple")
filtered := filter.Apply(items)
```

#### Selection Manager

Multi-selection state management.

```go
type SelectionMode int

const (
    SelectionModeNone SelectionMode = iota
    SelectionModeSingle
    SelectionModeMultiple
)

type SelectionManager struct {
    mode  SelectionMode
    items map[string]struct{}
}

func NewSelectionManager(mode SelectionMode) SelectionManager
func (s *SelectionManager) Select(item FilterableItem)
func (s *SelectionManager) Deselect(item FilterableItem)
func (s *SelectionManager) Toggle(item FilterableItem)
func (s *SelectionManager) SelectAll(items []FilterableItem)
func (s *SelectionManager) DeselectAll()
func (s *SelectionManager) InvertSelection(items []FilterableItem)
func (s *SelectionManager) SelectedIDs() []string
func (s *SelectionManager) GetSelected(items []FilterableItem) []FilterableItem

// Usage
selMgr := list.NewSelectionManager(list.SelectionModeMultiple)
selMgr.SelectAll(items)
selected := selMgr.GetSelected(items)
```

#### Base List

Complete list component with key bindings.

```go
type BaseList struct {
    items       []FilterableItem
    viewport    Viewport
    filter      Filter
    selection   SelectionManager
    cursor      int
    keyMap      KeyMap
}

type Action int

const (
    ActionMoveUp Action = iota
    ActionMoveDown
    ActionPageUp
    ActionPageDown
    ActionMoveToTop
    ActionMoveToBottom
    ActionToggleSelection
    ActionSelectAll
    ActionSelect
    ActionGoBack
    // ... more actions
)

type KeyMap struct {
    Up     []string
    Down   []string
    // ... more key bindings
}

func DefaultKeyMap() KeyMap
func (k *KeyMap) MatchAction(key string) Action

func NewBaseList() *BaseList
func (l *BaseList) SetItems(items []FilterableItem)
func (l *BaseList) SetKeyMap(keyMap KeyMap)
func (l *BaseList) SetFilter(query string)
func (l *BaseList) ClearFilter()
func (l *BaseList) Update(msg any) *BaseList
func (l *BaseList) View() string

// Usage
list := list.NewBaseList()
list.SetKeyMap(list.DefaultKeyMap())
list.SetItems(items)
newList := list.Update(render.KeyMsg{Key: "down"})
```

---

### Dialog Components (`ui/dialog`)

#### Dialog Types

```go
type DialogResult int

const (
    ActionConfirm DialogResult = iota
    ActionCancel
    ActionDismiss
)

type Callback func(result DialogResult, data any)

// Info Dialog
func NewInfoDialog(title, message string) *InfoDialog

// Confirm Dialog
func NewConfirmDialog(title, message string, callback Callback) *ConfirmDialog

// Input Dialog
func NewInputDialog(prompt, placeholder string, callback func(string)) *InputDialog

// Select List Dialog
func NewSelectListDialog(title string, options []string, callback func(int, string)) *SelectListDialog
```

#### Dialog Overlay

Stack management for multiple dialogs.

```go
type Overlay struct {
    dialogs []Dialog
}

func NewOverlay() Overlay
func (o *Overlay) Push(dialog Dialog)
func (o *Overlay) Pop() Dialog
func (o *Overlay) Peek() Dialog
func (o *Overlay) Count() int
func (o *Overlay) Clear()

// Usage
overlay := dialog.NewOverlay()
overlay.Push(infoDialog)
active := overlay.Peek()
```

---

### Form Components (`ui/forms`)

#### Input Types

```go
// Text Input
type TextInput struct {
    id        string
    label     string
    value     string
    required  bool
    validation func(string) error
}

func NewTextInput(id, label string) TextInput
func (t TextInput) Value() string
func (t *TextInput) SetValue(value string)
func (t *TextInput) Validate() error

// Text Area
type TextArea struct {
    // Multi-line text input
}

// Select Dropdown
type Select struct {
    options []string
    value   string
}

// Checkbox
type Checkbox struct {
    checked bool
}

// Radio Group
type RadioGroup struct {
    options []string
    value   string
}
```

#### Form Container

```go
type Form struct {
    inputs    []Input
    focus     int
}

func NewForm(inputs ...Input) Form
func (f Form) Update(msg any) (Form, render.Cmd)
func (f Form) View() string
func (f Form) Validate() error

// Usage
nameInput := forms.NewTextInput("name", "Name")
nameInput.SetRequired(true)

emailInput := forms.NewTextInput("email", "Email")
emailInput.SetValidation(validateEmail)

form := forms.NewForm(nameInput, emailInput)
```

---

### Message Components (`ui/components/messages`)

#### Message Types

```go
// User Message
type UserMessage struct {
    id               string
    content          string
    codeBlocks       []CodeBlock
    fileAttachments  []string
    filesAdded       int
    filesRemoved     int
    copyMode         bool
}

func NewUserMessage(id, content string) UserMessage
func (m *UserMessage) AddCodeBlock(language, code string)
func (m *UserMessage) SetFileAttachments(files []string)

// Assistant Message
type AssistantMessage struct {
    id           string
    content      string
    inputTokens  int
    outputTokens int
}

func NewAssistantMessage(id, content string) AssistantMessage

// Tool Message
type ToolMessage struct {
    id        string
    toolType  string
    toolName  string
    arguments string
    result    string
    error     string
}

func NewToolMessage(id, toolType string, arguments []string, result string) ToolMessage

// Fetch Message
type FetchType int

const (
    FetchTypeBasic FetchType = iota
    FetchTypeWebFetch
    FetchTypeWebSearch
    FetchTypeAgentic
)

type FetchMessage struct {
    id         string
    fetchType  FetchType
    request    *FetchRequest
    result     *FetchResult
    loaded     bool
}

type FetchRequest struct {
    URL    string
    Params map[string]string
    Query  string
}

// Diagnostic Message
type DiagnosticSeverity int

const (
    DiagnosticError DiagnosticSeverity = iota
    DiagnosticWarning
    DiagnosticInfo
    DiagnosticHint
)

type DiagnosticMessage struct {
    id        string
    source    string
    severity  DiagnosticSeverity
    message   string
    code      string
    line      int
    column    int
}

// Todo Message
type TodoStatus int

const (
    StatusPending TodoStatus = iota
    StatusInProgress
    StatusCompleted
)

type TodoMessage struct {
    todos        []Todo
    active       bool
}

type Todo struct {
    id          string
    label       string
    status      TodoStatus
    description string
    priority    int
}

func NewTodoMessage(id, label string, status TodoStatus) TodoMessage
func (m *TodoMessage) AddTodo(id, label string, status TodoStatus)
func (m *TodoMessage) ToggleExpanded(id string)
```

---

### Status Components (`ui/components/status`)

#### Service Status

```go
type State int

const (
    StateDisabled State = iota
    StateStarting
    StateReady
    StateError
)

type ServiceCmp struct {
    name       string
    state      State
    errorCount int
}

func NewService(name, language string) ServiceCmp
func (s *ServiceCmp) SetStatus(state State)
func (s *ServiceCmp) SetErrorCount(count int)
func (s *ServiceCmp) SetCompact(compact bool)
func (s *ServiceCmp) View() string
```

#### LSP Service List

```go
type LSPServiceInfo struct {
    Name         string
    Language     string
    State        State
    Diagnostics  DiagnosticSummary
}

type LSPList struct {
    services []LSPServiceInfo
    width    int
}

func NewLSPList() LSPList
func (l *LSPList) AddService(service LSPServiceInfo)
func (l *LSPList) OnlineCount() int
func (l *LSPList) TotalErrors() int
func (l *LSPList) TotalWarnings() int
func (l *LSPList) TotalDiagnostics() int
func (l *LSPList) View() string
```

#### MCP Service List

```go
type ToolCounts struct {
    Tools   int
    Prompts int
}

type MCPServiceInfo struct {
    Name       string
    State      State
    ToolCounts ToolCounts
}

type MCPList struct {
    services []MCPServiceInfo
    width    int
}

func NewMCPList() MCPList
func (m *MCPList) AddService(service MCPServiceInfo)
func (m *MCPList) ConnectedCount() int
func (m *MCPList) TotalTools() int
func (m *MCPList) TotalPrompts() int
func (m *MCPList) View() string
```

---

### Progress Components (`ui/components/progress`)

#### Progress Bar

```go
type ProgressBar struct {
    label             string
    current           int64
    total             int64
    width             int
    showLabel         bool
    showLabelInline   bool
    color             lipgloss.Color
}

func NewProgressBar() ProgressBar
func (p *ProgressBar) SetCurrent(current int64)
func (p *ProgressBar) SetTotal(total int64)
func (p *ProgressBar) SetLabel(label string)
func (p *ProgressBar) SetWidth(width int)
func (p *ProgressBar) ShowLabel(show bool)
func (p *ProgressBar) Completed() bool
func (p *ProgressBar) Percent() float64
func (p ProgressBar) View() string

// Usage
bar := progress.NewProgressBar()
bar.SetLabel("Downloading")
bar.SetTotal(100)
bar.SetCurrent(75)
view := bar.View() // "█████████░░░░░░░░░░░░ Downloading 75/100 (75%)"
```

#### Spinner

```go
type SpinnerType int

const (
    SpinnerDots SpinnerType = iota
    SpinnerLine
    SpinnerArrow
    SpinnerMoon
)

type Spinner struct {
    type        SpinnerType
    label       string
    current     int
    fps         int
    state       StateType
}

func NewSpinner(typ SpinnerType) Spinner
func (s *Spinner) SetType(typ SpinnerType)
func (s *Spinner) SetFPS(fps int)
func (s *Spinner) SetLabel(label string)
func (s *Spinner) Start()
func (s *Spinner) Stop()
func (s *Spinner) Reset()
func (s *Spinner) Running() bool
func (s Spinner) View() string

// Usage
spinner := progress.NewSpinner(progress.SpinnerDots)
spinner.SetFPS(30)
spinner.Start()
```

---

### Attachments Component (`ui/components/attachments`)

```go
type AttachmentType int

const (
    FileType AttachmentType = iota
    ImageType
    VideoType
    AudioType
    DocumentType
    ArchiveType
)

type Attachment struct {
    ID        string
    Name      string
    Path      string
    Type      AttachmentType
    Size      int64
    SizeStr   string
    Extension string
    MIMEType  string
    Added     time.Time
    Modified  time.Time
}

type AttachmentList struct {
    attachments []Attachment
    filter      string
    cursor      int
    selected    []string
    config      AttachmentConfig
    cache       map[string]string
}

type AttachmentConfig struct {
    Compact      bool
    ShowPreview  bool
    ShowSize     bool
    MaxWidth     int
}

func NewAttachmentList() AttachmentList
func (a *AttachmentList) Add(att Attachment)
func (a *AttachmentList) Remove(id string)
func (a *AttachmentList) SetFilter(filter string)
func (a *AttachmentList) Select(id string)
func (a *AttachmentList) Statistics() AttachmentStatistics

// Usage
list := attachments.NewAttachmentList()
list.SetCompact(true)
list.Add(attachments.Attachment{
    ID:   "att-1",
    Name: "document.pdf",
    Path: "/path/to/document.pdf",
    Size: 2400000,
})
```

---

### Pills Component (`ui/components/pills`)

```go
type PillStatus int

const (
    StatusPending PillStatus = iota
    StatusInProgress
    StatusCompleted
    StatusError
    StatusWarning
    StatusInfo
    StatusNeutral
)

type Pill struct {
    ID       string
    Label    string
    Status   PillStatus
    Count    int
    Expanded bool
}

type PillList struct {
    pills   []Pill
    cursor  int
    expanded map[string]bool
    config  PillConfig
    cache   map[string]string
}

type PillConfig struct {
    ShowIcons bool
    Inline    bool
}

func NewPillList() PillList
func (p *PillList) Add(pill Pill)
func (p *PillList) Remove(id string)
func (p *PillList) ExpandAll()
func (p *PillList) CollapseAll()
func (p *PillList) ToggleExpanded(id string)
func (p *PillList) SetInlineMode(inline bool)

// Usage
pills := pills.NewPillList()
pills.Add(pills.Pill{
    ID:     "pill-1",
    Label:  "Tasks",
    Status: pills.StatusPending,
    Count:  5,
})
```

---

## v2.0 API Migration Notes

### Type Changes

```go
// v1.0 (Bubbletea-only)
type Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (Model, tea.Cmd)
    View() string
}

// v2.0 (engine-agnostic)
type Model interface {
    Init() render.Cmd
    Update(any) (Model, render.Cmd)
    View() string
}
```

### Message Handling

```go
// v1.0
case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyUp:
        // Handle up

// v2.0
case render.KeyMsg:
    switch key.Key {
    case "up":
        // Handle up
```

### Commands

```go
// v1.0
tea.Quit
tea.Tick(duration, fn)
tea.Batch(cmd1, cmd2)

// v2.0
render.Quit()
render.Tick(duration, fn)
render.Batch(cmd1, cmd2)
```

---

**Version**: 2.0.0

**Last Updated**: Phase 10 completion (February 2026)
