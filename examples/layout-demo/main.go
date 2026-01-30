package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/layout"
	"github.com/wwsheng009/taproot/internal/ui/render"
)

// Styles
var (
	someStyle = lipgloss.NewStyle()
	dimStyle  = lipgloss.NewStyle().Faint(true)
	boldStyle = lipgloss.NewStyle().Bold(true)

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Margin(0, 1)

	labelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	areaInfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("239"))

	colors = []lipgloss.Color{
		lipgloss.Color("196"),
		lipgloss.Color("208"),
		lipgloss.Color("226"),
		lipgloss.Color("46"),
		lipgloss.Color("33"),
		lipgloss.Color("57"),
		lipgloss.Color("213"),
		lipgloss.Color("15"),
	}
)

// LayoutType represents different layout demonstrations
type LayoutType int

const (
	LayoutSplitVertical LayoutType = iota + 1
	LayoutSplitHorizontal
	LayoutFlexRow
	LayoutFlexColumn
	LayoutGrid
	LayoutGridGapped
	LayoutPadding
	LayoutAbsolute
)

func (lt LayoutType) String() string {
	switch lt {
	case LayoutSplitVertical:
		return "Split Vertical"
	case LayoutSplitHorizontal:
		return "Split Horizontal"
	case LayoutFlexRow:
		return "Flex Row"
	case LayoutFlexColumn:
		return "Flex Column"
	case LayoutGrid:
		return "Grid"
	case LayoutGridGapped:
		return "Grid (Gapped)"
	case LayoutPadding:
		return "Padding"
	case LayoutAbsolute:
		return "Absolute Positioning"
	default:
		return "Unknown"
	}
}

type model struct {
	layoutType  LayoutType
	width       int
	height      int
	showDetails bool
}

func (m *model) Init() error {
	return nil
}

func (m *model) Update(msg any) (render.Model, render.Cmd) {
	switch msg := msg.(type) {
	case render.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, render.Quit()
		case "1":
			m.layoutType = LayoutSplitVertical
		case "2":
			m.layoutType = LayoutSplitHorizontal
		case "3":
			m.layoutType = LayoutFlexRow
		case "4":
			m.layoutType = LayoutFlexColumn
		case "5":
			m.layoutType = LayoutGrid
		case "6":
			m.layoutType = LayoutGridGapped
		case "7":
			m.layoutType = LayoutPadding
		case "8":
			m.layoutType = LayoutAbsolute
		case "d":
			m.showDetails = !m.showDetails
		case "right":
			m.width++
		case "left":
			if m.width > 40 {
				m.width--
			}
		case "up":
			if m.height > 20 {
				m.height--
			}
		case "down":
			m.height++
		}
	case render.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *model) View() string {
	var b strings.Builder

	title := labelStyle.Render("Taproot Layout Demo")
	b.WriteString(title)
	b.WriteString("\n\n")

	help := []string{
		fmt.Sprintf("1-8: Change Layout | d: Toggle Details | Arrow keys: Resize | q: Quit"),
		fmt.Sprintf("Current: %s | Size: %dx%d", m.layoutType, m.width, m.height),
	}
	for _, h := range help {
		b.WriteString(dimStyle.Render(h))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	area := layout.NewArea(0, 0, m.width-2, m.height-8)
	layoutView := m.renderLayout(area)
	b.WriteString(layoutView)

	if m.showDetails {
		details := m.renderDetails(area)
		b.WriteString("\n")
		b.WriteString(boxStyle.BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("239")).
			Render(dimStyle.Render(details)))
	}

	return b.String()
}

func (m *model) renderLayout(area layout.Area) string {
	switch m.layoutType {
	case LayoutSplitVertical:
		return m.renderSplitVertical(area)
	case LayoutSplitHorizontal:
		return m.renderSplitHorizontal(area)
	case LayoutFlexRow:
		return m.renderFlexRow(area)
	case LayoutFlexColumn:
		return m.renderFlexColumn(area)
	case LayoutGrid:
		return m.renderGrid(area, false)
	case LayoutGridGapped:
		return m.renderGrid(area, true)
	case LayoutPadding:
		return m.renderPadding(area)
	case LayoutAbsolute:
		return m.renderAbsolute(area)
	default:
		return ""
	}
}

func (m model) renderSplitVertical(area layout.Area) string {
	top, bottom := layout.SplitVertical(area, layout.Percent(50))

	// Border: 2 (left+right), Padding: 2 (left+right), Margin: 2 (left+right)
	tw := top.Dx() - 2
	th := top.Dy() - 2
	bw := bottom.Dx() - 2
	bh := bottom.Dy() - 2

	topStyle := boxStyle.Copy().
		BorderForeground(colors[0]).
		Margin(0).
		Width(tw).
		Height(th)

	bottomStyle := boxStyle.Copy().
		BorderForeground(colors[1]).
		Margin(0).
		Width(bw).
		Height(bh)

	topBox := topStyle.Render(fmt.Sprintf("Top\n%d×%d", tw, th))
	bottomBox := bottomStyle.Render(fmt.Sprintf("Bottom\n%d×%d", bw, bh))

	return lipgloss.JoinVertical(lipgloss.Left, topBox, bottomBox)
}

func (m model) renderSplitHorizontal(area layout.Area) string {
	left, right := layout.SplitHorizontal(area, layout.Percent(50))

	// Border: 2, Padding: 2 (Margin removed for join)
	lw := left.Dx() - 2
	lh := left.Dy() - 2
	rw := right.Dx() - 2
	rh := right.Dy() - 2

	leftStyle := boxStyle.Copy().
		BorderForeground(colors[0]).
		Margin(0).
		Width(lw).
		Height(lh)

	rightStyle := boxStyle.Copy().
		BorderForeground(colors[1]).
		Margin(0).
		Width(rw).
		Height(rh)

	leftBox := leftStyle.Render(fmt.Sprintf("Left\n%d×%d", lw, lh))
	rightBox := rightStyle.Render(fmt.Sprintf("Right\n%d×%d", rw, rh))

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}

func (m model) renderFlexRow(area layout.Area) string {
	children := []layout.FlexChild{
		layout.NewFlexChild(layout.Fixed(20)),
		layout.NewFlexChild(layout.Grow{}).WithGrow(),
		layout.NewFlexChild(layout.Fixed(20)),
	}
	areas := layout.RowLayout(area, children)

	var boxes []string
	for i, a := range areas {
		// Border: 2, Padding: 2
		w := a.Dx() - 2
		h := a.Dy() - 2

		style := boxStyle.Copy().
			BorderForeground(colors[i%len(colors)]).
			Margin(0).
			Width(w).
			Height(h)
		box := style.Render(fmt.Sprintf("Item %d\n%d×%d", i+1, w, h))
		boxes = append(boxes, box)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, boxes...)
}

func (m model) renderFlexColumn(area layout.Area) string {
	children := []layout.FlexChild{
		layout.NewFlexChild(layout.Fixed(5)),
		layout.NewFlexChild(layout.Grow{}).WithGrow(),
		layout.NewFlexChild(layout.Fixed(5)),
	}
	areas := layout.ColumnLayout(area, children)

	var boxes []string
	for i, a := range areas {
		// Border: 2, Padding: 2
		w := a.Dx() - 2
		h := a.Dy() - 2

		style := boxStyle.Copy().
			BorderForeground(colors[i%len(colors)]).
			Margin(0).
			Width(w).
			Height(h)
		box := style.Render(fmt.Sprintf("Item %d\n%d×%d", i+1, w, h))
		boxes = append(boxes, box)
	}

	return lipgloss.JoinVertical(lipgloss.Left, boxes...)
}

func (m model) renderGrid(area layout.Area, gapped bool) string {
	var config layout.GridConfig
	if gapped {
		config = layout.NewGridConfig(2, 3).WithRowGaps(1).WithColGaps(1)
	} else {
		config = layout.NewGridConfig(2, 3)
	}
	cells := layout.GridLayout(area, config)

	var rows [][]string
	for row := 0; row < config.Rows; row++ {
		var rowBoxes []string
		for col := 0; col < config.Cols; col++ {
			cell := layout.GetCell(cells, config, row, col)
			// Border: 2, Padding: 2
			w := cell.Dx() - 2
			h := cell.Dy() - 2

			style := boxStyle.Copy().
				BorderForeground(colors[(row*config.Cols+col)%len(colors)]).
				Margin(0).
				Width(w).
				Height(h)
			box := style.Render(fmt.Sprintf("%d,%d\n%d×%d", row, col, w, h))
			rowBoxes = append(rowBoxes, box)
		}
		rows = append(rows, rowBoxes)
	}

	var rowStrings []string
	for _, row := range rows {
		rowStrings = append(rowStrings, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rowStrings...)
}

func (m model) renderPadding(area layout.Area) string {
	styles := []lipgloss.Style{
		boxStyle.Copy().BorderForeground(colors[0]),
		boxStyle.Copy().BorderForeground(colors[1]),
		boxStyle.Copy().BorderForeground(colors[2]),
		boxStyle.Copy().BorderForeground(colors[3]),
	}

	areas := []layout.Area{
		area,
		layout.Pad(area, 1),
		layout.Pad(area, 2),
		layout.Pad(area, 3),
	}

	lines := []string{"No padding", "Pad(1)", "Pad(2)", "Pad(3)"}

	var rendered []string
	for i, a := range areas {
		// Border: 2, Padding: 2 (show padded content, not container)
		w := a.Dx() - 2
		h := a.Dy() - 2

		style := styles[i].Margin(0).Width(w).Height(h)
		box := style.Render(lines[i])
		rendered = append(rendered, box)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func (m model) renderAbsolute(area layout.Area) string {
	// Box sizes based on 3x3 grid
	boxWidth := area.Dx() / 3
	boxHeight := area.Dy() / 3

	positions := []struct {
		name string
	}{
		{"TopLeft"},
		{"TopCenter"},
		{"TopRight"},
		{"LeftCenter"},
		{"Center"},
		{"RightCenter"},
		{"BottomLeft"},
		{"BottomCenter"},
		{"BottomRight"},
	}

	var boxes []string
	for i, p := range positions {
		// Account for border+padding
		w := boxWidth - 2
		h := boxHeight - 2

		if w < 8 {
			w = 8
		}
		if h < 3 {
			h = 3
		}

		style := boxStyle.Copy().
			BorderForeground(colors[i%len(colors)]).
			Margin(0).
			Width(w).
			Height(h)
		box := style.Render(p.name)
		boxes = append(boxes, box)
	}

	// Create 3x3 grid
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, boxes[0], boxes[1], boxes[2])
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, boxes[3], boxes[4], boxes[5])
	row3 := lipgloss.JoinHorizontal(lipgloss.Top, boxes[6], boxes[7], boxes[8])

	return lipgloss.JoinVertical(lipgloss.Left, row1, row2, row3)
}

func (m model) renderDetails(area layout.Area) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Area Position: (%d, %d)\n", area.TopLeft().X, area.TopLeft().Y))
	b.WriteString(fmt.Sprintf("Area Size: %d×%d\n", area.Dx(), area.Dy()))
	b.WriteString(fmt.Sprintf("Area Bounds: (%d, %d) to (%d, %d)\n",
		area.TopLeft().X, area.TopLeft().Y,
		area.BottomRight().X, area.BottomRight().Y))
	b.WriteString(fmt.Sprintf("Is Empty: %v", area.Empty()))

	return b.String()
}

func main() {
	engine, err := render.CreateEngine(render.EngineBubbletea, render.DefaultConfig())
	if err != nil {
		panic(err)
	}

	m := &model{
		width:       80,
		height:      24,
		showDetails: false,
		layoutType:  LayoutSplitVertical,
	}

	if err := engine.Start(m); err != nil {
		panic(err)
	}
}
