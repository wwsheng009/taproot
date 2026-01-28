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

Terminal image display with multiple protocols.

```go
type RendererType int

const (
    RendererAuto RendererType = iota
    RendererKitty
    RendereriTerm2
    RendererBlocks
)

type Image struct {
    // Internal state
}

func New(path string) *Image
func (img *Image) SetRenderer(renderer RendererType) tea.Cmd
func (img *Image) Reload() tea.Cmd
func (img *Image) SetSize(w, h int)
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

### Theme Structure (`internal/tui/styles`)

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

**Version**: 0.9.0

**Last Updated**: 2024-01-28
