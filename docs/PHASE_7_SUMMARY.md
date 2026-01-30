# Phase 7.1-7.3: Core Component Library (v2.0)

## Overview

This document summarizes the implementation of Phase 7.1-7.3: Core Component Library for Taproot v2.0. These components are designed to be engine-agnostic, implementing the `render.Model` interface for compatibility with both Bubbletea and Ultraviolet rendering engines.

---

## Phase 7.1: Header Component

### Implementation

**File:** `internal/ui/components/header/header_v2.go`

**Type:** `HeaderV2`

**Features:**
- Implements `render.Model` interface for engine-agnostic rendering
- Displays brand name, title, and gradient text
- Shows token usage progress bar with visual indicator (╱)
- Displays working directory with intelligent truncation (max 4 path components)
- Shows error/warning counts with colored icons
- Supports keyboard shortcuts (Ctrl+D to toggle details state)
- Handles window resize events
- Manual ANSI-safe text truncation to prevent multi-line rendering
- Proper padding and width management

**State Management:**
```go
type HeaderV2 struct {
    width        int
    height       int
    brand        string
    title        string
    sessionTitle string
    workingDir   string
    tokenUsed    int
    tokenMax     int
    cost         float64
    errorCount   int
    warnCount    int
    detailsOpen  bool
    initialized  bool
}
```

**Key Methods:**
- `Init() error` - Initialize the component
- `Update(msg any) (Model, Cmd)` - Handle events
- `View() string` - Render the header
- `SetSize(width, height int)` - Set dimensions
- `SetBrand(brand, title string)` - Set brand text
- `SetWorkingDirectory(cwd string)` - Set CWD
- `SetTokenUsage(used, max int, cost float64)` - Set token stats
- `SetErrorCount(count int)` - Set error count
- `SetWarnCount(count int)` - Set warning count
- `SetDetailsOpen(open bool)` - Toggle details state

**Messages Supported:**
- `*render.KeyMsg` - Keyboard input
- `*render.WindowSizeMsg` - Window resize
- `*render.FocusMsg` - Focus events
- `*render.CustomMsg` - Custom messages

---

## Phase 7.2: Button/Label/Text Basic Components

### Implementation

**Directory:** `internal/ui/components/basic/`

### Button Component

**File:** `basic/button.go`

**Type:** `Button`

**Interfaces Implemented:**
- `render.Model` - For rendering and event handling
- `Clickable` - For click event handling

**Features:**
- Four states: Normal, Focused, Pressed, Disabled
- Keyboard shortcuts: Enter or Space to click
- Click handler callbacks
- Focus/blur support
- Configurable styled appearance
- ID-based identification

**State Management:**
```go
type Button struct {
    id           string
    label        string
    state        ButtonState
    clickHandler func()
    focused      bool
    disabled     bool
    style        *ButtonStyle
    width        int
    initialized  bool
}
```

**Key Methods:**
- `Init() error` - Initialize
- `Update(msg any) (Model, Cmd)` - Handle events
- `View() string` - Render button text
- `SetClickHandler(handler func())` - Set callback
- `Focus() / Blur()` - Manage focus
- `SetDisabled(disabled bool)` - Enable/disable

**Messages Supported:**
- `*render.KeyMsg` - Enter/Space to click, navigation keys
- `*render.FocusGainMsg` - Focus gained
- `*render.BlurMsg` - Focus lost

### Label Component

**File:** `basic/label.go`

**Type:** `Label`

**Features:**
- Simple text label display
- Configurable styling with `lipgloss.Style`
- Width and alignment support (Left/Center/Right)
- Focus tracking

**Key Methods:**
- `SetText(text string)` - Set label text
- `SetStyle(style lipgloss.Style)` - Set styling
- `SetWidth(width int)` - Set width
- `SetAlign(align lipgloss.Position)` - Set alignment

### Text Component

**File:** `basic/text.go`

**Type:** `Text`

**Features:**
- Multiline text support
- Word wrapping to specified width
- Height truncation (max lines)
- Dynamic content updates
- Configurable styling

**Key Methods:**
- `SetContent(content string)` - Set text content
- `AppendText(text string)` - Append text
- `ClearText()` - Clear all content
- `SetWidth(width int)` - Enable/disable wrapping
- `SetHeight(height int)` - Set max height limit
- `SetWrapText(wrap bool)` - Toggle wrapping

### Types

**File:** `basic/types.go`

**Definitions:**
- `Clickable` interface - Click behavior
- `ButtonState` enum - Button states
- `ButtonStyle` struct - Style configuration with `DefaultButtonStyle()`

---

## Phase 7.3: ProgressBar/Spinner Components

### Implementation

**Directory:** `internal/ui/components/progress/`

### ProgressBar Component

**File:** `progress/progressbar.go`

**Type:** `ProgressBar`

**Features:**
- Visual progress bar with full/empty indicator characters
- Percentage display (e.g., "50%")
- Optional label prefix
- Configurable styling (full bar vs empty bar colors)
- Current/total value management
- Increment and add operations
- Completion detection
- Reset functionality

**Visual Style:**
- Full bar: "█" (configurable color)
- Empty bar: "░" (configurable color)
- Percent: "50%" (same color as full bar)

**Key Methods:**
- `SetCurrent(current float64)` - Set current progress
- `Add(delta float64)` - Add to progress
- `Increment()` - Add 1
- `SetTotal(total float64)` - Set total
- `SetLabel(label string)` - Set optional label
- `Percent() float64` - Get 0-100 percentage
- `Completed() bool` - Check if complete
- `Reset()` - Reset to 0
- `SetStyle(*ProgressBarStyle)` - Set styling

**Style Configuration:**
```go
type ProgressBarStyle struct {
    FullBarStyle  lipgloss.Style
    EmptyBarStyle lipgloss.Style
    ShowPercent   bool
    ShowLabel     bool
    Width         int
}
```

### Spinner Component

**File:** `progress/spinner.go`

**Type:** `Spinner`

**Features:**
- Animated spinner with multiple frame types
- Multiple spinner styles: Dots, Line, Arrow, Moon
- Configurable FPS (frames per second)
- Optional label prefix
- Start/stop/reset control
- State tracking across frames
- Tick message handling for animation

**Spinner Types:**
1. **Dots:** `⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏` (default)
2. **Line:** `- \ | /`
3. **Arrow:** `← ↖ ↑ ↗ → ↘ ↓ ↙`
4. **Moon:** `🌑 🌒 🌓 🌔 🌕 🌖 🌗 🌘`

**Key Methods:**
- `SetType(SpinnerType)` - Set spinner style
- `SetLabel(label string)` - Set label
- `SetColor(color lipgloss.Color)` - Set color
- `SetFPS(fps int)` - Set animation speed
- `Stop()` - Stop animation
- `Reset()` - Reset to initial state
- `Running() bool` - Check if running

**Messages Supported:**
- `*TickMsg` - Internal tick messages
- `*render.TickMsg` - Generic render engine ticks

---

## Implementation Details

### Error Handling

All components follow these error handling patterns:
- `Init()` returns `error` but typically returns `nil`
- `Update()` returns `(Model, Cmd)` to allow for command chaining
- State is protected from invalid values (clamping, bounds checking)

### State Management

All components implement proper state management:
- Immutable-like updates return new model references
- State mutations happen during message handling
- `initialized` flag tracks initialization status

### Rendering Philosophy

1. **Width Management:** All components respect specified widths and truncate or wrap appropriately
2. **ANSI Safety:** Text truncation preserves ANSI escape sequences for proper styling
3. **Single Line:** Components avoid unintended multi-line rendering
4. **Style Composure:** Uses `lipgloss` for consistent styling with theme support

### Testing

Each component includes comprehensive test coverage:
- Construction/initialization tests
- Interface compliance tests
- State mutation tests
- Rendering/output tests
- Event handling tests
- Edge case tests

---

## Files Created

### Phase 7.1 - Header
- `internal/ui/components/header/header_v2.go` (318 lines) - Main implementation
- `internal/ui/components/header/header_v2_test.go` (166 lines) - Tests

### Phase 7.2 - Basic Components
- `internal/ui/components/basic/types.go` (42 lines) - Type definitions
- `internal/ui/components/basic/button.go` (191 lines) - Button component
- `internal/ui/components/basic/label.go` (107 lines) - Label component
- `internal/ui/components/basic/text.go` (182 lines) - Text component
- `internal/ui/components/basic/basic_test.go` (224 lines) - Tests

### Phase 7.3 - Progress Components
- `internal/ui/components/progress/progressbar.go` (205 lines) - ProgressBar component
- `internal/ui/components/progress/spinner.go` (184 lines) - Spinner component
- `internal/ui/components/progress/progress_test.go` (364 lines) - Tests

**Total Lines:** ~1,883 lines of code (including tests)

---

## Test Results

All tests pass:

```bash
# Header
ok  	github.com/wwsheng009/taproot/internal/ui/components/header	0.467s

# Basic Components
ok  	github.com/wwsheng009/taproot/internal/ui/components/basic	0.221s

# Progress Components
ok  	github.com/wwsheng009/taproot/internal/ui/components/progress	0.366s
```

---

## Design Decisions

### Engine Agnosticism
- Components implement `render.Model` instead of `tea.Model`
- Use `render` package types (KeyMsg, TickMsg, etc.)
- Compatible with both Bubbletea and Ultraviolet via adapters

### Separation of Concerns
- Each component is self-contained in its own file
- Shared types in `types.go` file
- Tests in `*_test.go` files alongside implementation

### Performance Considerations
- Minimal state recalculation
- Cached style configurations
- Efficient string building with `strings.Builder`
- Truncation algorithms preserve ANSI sequences

---

## Next Steps

Potential enhancements for future phases:
- **Phase 7.4:** Input fields (text input, password)
- **Phase 7.5:** Select/dropdown components
- **Phase 7.6:** Table/grid components
- **Phase 7.7:** Tabs component
- **Phase 7.8:** Menu/context menu
- **Phase 7.9:** Toggle/checkbox components
- **Phase 7.10:** Slider component

---

## Acknowledgments

The component designs are inspired by:
- Bubble Tea's component patterns
- Lips gloss styling system
- Original Taproot v1.0 component implementations
- Modern TUI best practices

---
