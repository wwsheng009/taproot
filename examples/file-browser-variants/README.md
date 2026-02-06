# File Browser Examples - Layout System Comparison

This directory contains two versions of the same file browser application, demonstrating different layout approaches in Taproot.

## Examples

### 1. file-browser-layout/

**Using Taproot Core Layout System**

The core layout system (`ui/layout/`) provides flexible, constraint-based layout primitives based on `image.Rectangle`. It focuses on area calculations while delegating rendering to the caller (typically lipgloss).

**Key Features:**

- **Area-based layout**: Uses `layout.Area` (based on `image.Rectangle`) to define regions
- **Constraint system**: Flexible sizing with `Fixed`, `Percent`, `Ratio`, `Grow` constraints
- **Layout primitives**:
  - `SplitHorizontal/Vertical` - Divide areas
  - `RowLayout/ColumnLayout` - Flex layouts
  - `GridLayout` - Grid-based positioning
  - `CenterRect`, `TopLeftRect` - Positioning helpers
- **Rendering delegation**: Layout returns areas, rendering is handled separately via lipgloss

**Layout Pattern:**
```go
// Calculate areas
m.mainArea = layout.NewArea(0, 2, m.width, totalHeight)
m.fileListArea, m.rightPanelArea = layout.SplitHorizontal(
    m.mainArea,
    layout.Fixed(fileListWidth),
)

// Render with lipgloss
footer := lipgloss.NewStyle().Padding(0, 1).Render(text)
```

**Best for:**
- Applications needing flexible, dynamic layouts
- Component-based UI construction
- Projects already using lipgloss for styling
- Complex panel management (resizable, foldable, etc.)

---

### 2. file-browser-buffer/

**Using Buffer Layout System**

The buffer layout system (`ui/render/buffer/`) provides a 2D cell grid-based rendering system with exact dimension calculations and native wide character support.

**Key Features:**

- **Cell-based rendering**: 2D grid of cells with exact positioning
- **Native wide character support**: Built-in CJK character handling
- **Component isolation**: Sub-buffers for independent rendering
- **LayoutManager**: Orchestrates multi-component layouts
- **Buffer API**:
  - `SetCell`, `WriteString`, `FillRect` - Cell operations
  - `WriteBuffer` - Buffer composition
  - `Render` - Final ANSI output

**Layout Pattern:**
```go
// Create buffer and layout manager
buf := buffer.NewBuffer(width, height)
lm := buffer.NewLayoutManager(width, height)

// Calculate regions and render
mainBuf := buffer.NewBuffer(m.width, m.height)
mainBuf.WriteBuffer(buffer.Point{X: 0, Y: 0}, headerBuf)
mainBuf.WriteBuffer(buffer.Point{X: fileListWidth, Y: 2}, contentBuf)

// Render final output
output := mainBuf.Render()
```

**Best for:**
- Applications needing precise pixel control
- Complex text rendering with CJK characters
- UIs with frequent layout changes
- Performance-critical rendering (~0.15ms/frame)

---

## Comparison

| Aspect | Core Layout | Buffer Layout |
|--------|-------------|---------------|
| **Based on** | `image.Rectangle` areas | 2D cell grid |
| **API Style** | Layout calculation + lipgloss rendering | Direct buffer operations |
| **Wide Characters** | Handled by lipgloss | Native support |
| **Styling** | lipglos styles | `buffer.Style` struct |
| **Layout Changes** | Recalculate areas | Re-render buffers |
| **Debugging** | Area dimensions visible | Cell-level inspection |
| **Performance** | Depends on lipgloss | Very high (~0.15ms/frame) |
| **Learning Curve** | Easier for existing lipgloss users | Concepts similar to canvas/pixel UI |

---

## Running the Examples

### Taproot Layout Version:
```bash
go run examples/file-browser-layout/main.go
```

### Buffer Layout Version:
```bash
go run examples/file-browser-buffer/main.go
```

---

## Key Differences in Implementation

### Layout Calculation

**Core Layout:**
```go
func (m *Model) recalculateLayout() {
    // Using Area-based layout
    m.mainArea = layout.NewArea(0, 2, m.width, totalHeight)
    m.fileListArea, m.rightPanelArea = layout.SplitHorizontal(
        m.mainArea,
        layout.Fixed(fileListWidth),
    )
}
```

**Buffer Layout:**
```go
func (m *Model) recalculateLayout() {
    // Direct rectangle calculation
    m.mainRect = buffer.Rect{X: 0, Y: 0, Width: m.width, Height: mainHeight}
    m.fileListRect = buffer.Rect{X: 0, Y: 0, Width: fileListWidth, Height: mainHeight}
    m.contentRect = buffer.Rect{X: fileListWidth, Y: 0, Width: m.width - fileListWidth, Height: mainHeight}
}
```

### Rendering

**Core Layout:**
```go
func (m Model) View() string {
    // Calculate areas
    // ...

    // Render with lipgloss
    header := styleHeader.Render(text)
    panel := panelStyle.Width(width).Render(content)

    // Join and return
    return lipgloss.JoinHorizontal(lipgloss.Left, header, panel)
}
```

**Buffer Layout:**
```go
func (m Model) View() string {
    // Create main buffer
    mainBuf := buffer.NewBuffer(m.width, m.height)

    // Render sub-buffers
    headerBuf, _ := header.Render()
    contentBuf, _ := content.Render()

    // Compose
    mainBuf.WriteBuffer(buffer.Point{X: 0, Y: 0}, headerBuf)
    mainBuf.WriteBuffer(buffer.Point{X: fileListWidth, Y: 2}, contentBuf)

    // Render final
    return mainBuf.Render()
}
```

---

## Which One to Use?

### Choose Core Layout if:
- You're already familiar with lipgloss
- Your UI has complex, dynamic component hierarchies
- You want flexible panel sizes and constraints
- Performance is not the primary concern

### Choose Buffer Layout if:
- You need precise control over every cell
- Your app renders CJK/wide characters
- Layout changes frequently (animations, real-time updates)
- Performance is critical
- You prefer a canvas-like rendering approach

---

## Common Features

Both examples include:

- **File browser** with directory navigation
- **Resizable panels** with `[` / `]` keys
- **File preview** with markdown rendering
- **Command palette** with shell execution
- **Search mode** with real-time filtering
- **Keyboard navigation** (hjkl, arrows, PageUp/Down)
- **Focus management** between panels
- **Todo list integration**
- **LSP/MCP status indicators**

---

## Status and Roadmap

Both versions are fully functional. Future enhancements may include:

- [ ] More layout examples (grid, tabs, split view)
- [ ] Performance benchmarks
- [ ] Hybrid approach discussion
- [ ] Video tutorials
- [ ] Interactive layout playground

For questions or contributions, see the main Taproot repository.
