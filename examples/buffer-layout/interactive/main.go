// Buffer Interactive Example - Full Bubbletea app using buffer layout
//
// This example demonstrates:
// - Integration with Bubbletea framework
// - Interactive components with keyboard navigation
// - Dynamic content updates
// - Real-time layout recalculation
//
// Usage: go run main.go

package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render/buffer"
)

// Model represents the application state
type Model struct {
	width      int
	height     int
	quitting   bool
	cursor     int
	items      []string
	layoutMgr  *buffer.LayoutManager
	ticks      int
}

// ItemComponent renders a single menu item
type ItemComponent struct {
	label  string
	cursor bool
}

func (i *ItemComponent) Render(buf *buffer.Buffer, rect buffer.Rect) {
	prefix := "  "
	suffix := ""
	if i.cursor {
		prefix = "▶ "
		suffix = " ◀"
	}

	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, prefix+i.label+suffix, buffer.Style{
		Foreground: "#86",
		Bold:       i.cursor,
	})
}

func (i *ItemComponent) MinSize() (int, int)     { return 20, 1 }
func (i *ItemComponent) PreferredSize() (int, int) { return 30, 1 }

// StatsComponent shows statistics
type StatsComponent struct {
	ticks int
	width int
}

func (s *StatsComponent) Render(buf *buffer.Buffer, rect buffer.Rect) {
	stats := []string{
		fmt.Sprintf("Width: %d", s.width),
		fmt.Sprintf("Height: %d", rect.Height),
		fmt.Sprintf("Ticks: %d", s.ticks),
		fmt.Sprintf("Items: %d", 5),
	}

	for i, stat := range stats {
		buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + i}, stat, buffer.Style{Foreground: "#244"})
	}
}

func (s *StatsComponent) MinSize() (int, int)     { return 15, 4 }
func (s *StatsComponent) PreferredSize() (int, int) { return 20, 6 }

// InfoComponent shows information panel
type InfoComponent struct {
	cursor int
	items  []string
}

func (i *InfoComponent) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw border
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, "┌", buffer.Style{Foreground: "#86"})
	for x := 1; x < rect.Width-1; x++ {
		buf.WriteString(buffer.Point{X: rect.X + x, Y: rect.Y}, "─", buffer.Style{Foreground: "#86"})
	}
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y}, "┐", buffer.Style{Foreground: "#86"})

	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + rect.Height - 1}, "└", buffer.Style{Foreground: "#86"})
	for x := 1; x < rect.Width-1; x++ {
		buf.WriteString(buffer.Point{X: rect.X + x, Y: rect.Y + rect.Height - 1}, "─", buffer.Style{Foreground: "#86"})
	}
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y + rect.Height - 1}, "┘", buffer.Style{Foreground: "#86"})

	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + 1}, "│", buffer.Style{Foreground: "#86"})
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y + 1}, "│", buffer.Style{Foreground: "#86"})

	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + rect.Height - 2}, "│", buffer.Style{Foreground: "#86"})
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y + rect.Height - 2}, "│", buffer.Style{Foreground: "#86"})

	// Draw side borders
	for y := 2; y < rect.Height-2; y++ {
		buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + y}, "│", buffer.Style{Foreground: "#86"})
		buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y + y}, "│", buffer.Style{Foreground: "#86"})
	}

	// Draw title
	title := " Selected Item "
	buf.WriteString(buffer.Point{X: rect.X + 2, Y: rect.Y}, title, buffer.Style{Foreground: "#235", Background: "#86", Bold: true})

	// Draw selected item info
	if i.cursor < len(i.items) {
		item := i.items[i.cursor]
		lines := []string{
			"",
			"Name: " + item,
			"",
			"Description:",
			"  This is a sample item",
			"  demonstrating the buffer",
			"  layout system with",
			"  Bubbletea integration.",
		}

		y := rect.Y + 2
		for _, line := range lines {
			if y < rect.Y+rect.Height-2 {
				buf.WriteString(buffer.Point{X: rect.X + 2, Y: y}, line, buffer.Style{Foreground: "#250"})
				y++
			}
		}
	}
}

func (i *InfoComponent) MinSize() (int, int)     { return 25, 10 }
func (i *InfoComponent) PreferredSize() (int, int) { return 35, 15 }

func NewModel() Model {
	items := []string{
		"Dashboard",
		"Projects",
		"Settings",
		"About",
		"Quit",
	}

	return Model{
		width:    80,
		height:   24,
		cursor:   0,
		items:    items,
		layoutMgr: buffer.NewLayoutManager(80, 24),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type TickMsg time.Time

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.cursor == len(m.items)-1 {
				m.quitting = true
				return m, tea.Quit
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.cursor == len(m.items)-1 {
				m.quitting = true
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layoutMgr.SetSize(msg.Width, msg.Height)
	case TickMsg:
		m.ticks++
		return m, m.Init()
	}
	return m, m.Init()
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	// Create buffer
	buf := buffer.NewBuffer(m.width, m.height)

	// Draw header
	header := buffer.NewTextComponent(
		"  Taproot Buffer Layout + Bubbletea  ",
		buffer.Style{Bold: true, Foreground: "#15", Background: "#32"},
	)
	header.Render(buf, buffer.Rect{X: 0, Y: 0, Width: m.width, Height: 1})

	// Draw separator
	for x := 0; x < m.width; x++ {
		buf.WriteString(buffer.Point{X: x, Y: 1}, "─", buffer.Style{Foreground: "#244"})
	}

	// Menu items
	menuY := 3
	for i, item := range m.items {
		itemComp := &ItemComponent{
			label:  item,
			cursor: i == m.cursor,
		}
		itemComp.Render(buf, buffer.Rect{X: 5, Y: menuY + i, Width: 30, Height: 1})
	}

	// Stats panel
	statsComp := &StatsComponent{
		ticks: m.ticks,
		width: m.width,
	}
	statsComp.Render(buf, buffer.Rect{X: 5, Y: menuY + len(m.items) + 2, Width: 20, Height: 5})

	// Info panel
	infoComp := &InfoComponent{
		cursor: m.cursor,
		items:  m.items,
	}
	infoComp.Render(buf, buffer.Rect{X: m.width/2 + 2, Y: 3, Width: m.width/2 - 5, Height: m.height - 5})

	// Footer
	footerText := " ↑↓: Navigate | Enter/Space: Select | Q: Quit "
	footerX := (m.width - len(footerText)) / 2
	buf.WriteString(buffer.Point{X: footerX, Y: m.height - 1}, footerText, buffer.Style{Foreground: "#244"})

	return buf.Render()
}

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
