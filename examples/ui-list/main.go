// Simple list example using the engine-agnostic list package
package main

import (
	"fmt"
	"strings"

	"github.com/yourorg/taproot/internal/ui/list"
	"github.com/yourorg/taproot/internal/ui/render"
)

// SimpleListModel demonstrates using the list package with a custom model
type SimpleListModel struct {
	items    []list.Item
	cursor   int
	selected map[string]struct{}
	viewport *list.Viewport
	selMgr   *list.SelectionManager
	width    int
	height   int
	quit     bool
}

func NewSimpleListModel() *SimpleListModel {
	items := []list.Item{
		list.NewListItem("1", "Apple", "A red fruit"),
		list.NewListItem("2", "Banana", "A yellow fruit"),
		list.NewListItem("3", "Cherry", "A small red fruit"),
		list.NewListItem("4", "Date", "A sweet brown fruit"),
		list.NewListItem("5", "Elderberry", "A dark purple fruit"),
		list.NewListItem("6", "Fig", "A soft sweet fruit"),
		list.NewListItem("7", "Grape", "A small juicy fruit"),
		list.NewListItem("8", "Honeydew", "A sweet melon"),
		list.NewListItem("9", "Kiwi", "A fuzzy green fruit"),
		list.NewListItem("10", "Lemon", "A sour citrus fruit"),
	}

	return &SimpleListModel{
		items:    items,
		cursor:   0,
		selected: make(map[string]struct{}),
		viewport: list.NewViewport(5, len(items)),
		selMgr:   list.NewSelectionManager(list.SelectionModeMultiple),
		width:    60,
		height:   10,
		quit:     false,
	}
}

// Init initializes the model
func (m *SimpleListModel) Init() error {
	m.viewport.SetTotal(len(m.items))
	return nil
}

// Update handles incoming messages
func (m *SimpleListModel) Update(msg any) (render.Model, render.Cmd) {
	switch msg := msg.(type) {
	case render.KeyMsg:
		// Convert key string to action
		key := msg.String()
		action := m.matchAction(key)

		switch action {
		case list.ActionMoveUp:
			m.moveUp()
		case list.ActionMoveDown:
			m.moveDown()
		case list.ActionToggleSelection:
			m.toggleSelection()
		case list.ActionMoveToTop:
			m.moveToTop()
		case list.ActionMoveToBottom:
			m.moveToBottom()
		case list.ActionQuit:
			m.quit = true
		}

	case render.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.SetVisible(m.height - 4) // Leave room for header/footer
	}

	return m, nil
}

// View renders the model
func (m *SimpleListModel) View() string {
	if m.quit {
		return "Goodbye!"
	}

	var b strings.Builder

	// Header
	b.WriteString("╔════════════════════════════════════════════════════════╗\n")
	b.WriteString("║              Simple List Demo (v2.0)                   ║\n")
	b.WriteString("╠════════════════════════════════════════════════════════╣\n")
	b.WriteString("║  Keys: ↑/k=up  ↓/j=down  space=select  g=top  G=bottom ║\n")
	b.WriteString("╚════════════════════════════════════════════════════════╝\n\n")

	// Items
	start, end := m.viewport.Range()
	for i := start; i < end; i++ {
		if i >= len(m.items) {
			break
		}

		item := m.items[i]
		cursor := " "
		if i == m.cursor {
			cursor = "→"
		}

		checked := " "
		if _, ok := m.selected[item.ID()]; ok {
			checked = "✓"
		}

		// Get title from ListItem
		if li, ok := item.(*list.ListItem); ok {
			line := fmt.Sprintf("%s [%s] %s", cursor, checked, li.Title())
			if li.Desc() != "" {
				line += fmt.Sprintf(" - %s", li.Desc())
			}
			b.WriteString(line + "\n")
		}
	}

	// Footer
	b.WriteString(fmt.Sprintf("\n  Selected: %d | Showing %d-%d of %d | Scroll: %s",
		len(m.selected),
		start+1, end,
		len(m.items),
		m.viewport.ScrollIndicator()))

	return b.String()
}

func (m *SimpleListModel) matchAction(key string) list.Action {
	// Simple key matching
	switch key {
	case "up", "k":
		return list.ActionMoveUp
	case "down", "j":
		return list.ActionMoveDown
	case " ":
		return list.ActionToggleSelection
	case "g":
		return list.ActionMoveToTop
	case "G":
		return list.ActionMoveToBottom
	case "q", "ctrl+c":
		return list.ActionQuit
	default:
		return list.ActionNone
	}
}

func (m *SimpleListModel) moveUp() {
	if m.cursor > 0 {
		m.cursor--
		m.viewport.MoveUp()
	}
}

func (m *SimpleListModel) moveDown() {
	if m.cursor < len(m.items)-1 {
		m.cursor++
		m.viewport.MoveDown()
	}
}

func (m *SimpleListModel) toggleSelection() {
	id := m.items[m.cursor].ID()
	if _, ok := m.selected[id]; ok {
		delete(m.selected, id)
	} else {
		m.selected[id] = struct{}{}
	}
}

func (m *SimpleListModel) moveToTop() {
	m.cursor = 0
	m.viewport.MoveToTop()
}

func (m *SimpleListModel) moveToBottom() {
	m.cursor = len(m.items) - 1
	m.viewport.MoveToBottom()
}

func main() {
	fmt.Println("=== Engine-Agnostic List Demo ===")
	fmt.Println()
	fmt.Println("This demo uses the new internal/ui/list package")
	fmt.Println("which is independent of any rendering engine.")
	fmt.Println()
	fmt.Println("The components demonstrated:")
	fmt.Println("  - list.Item: Basic item interface")
	fmt.Println("  - list.Viewport: Virtualized scrolling")
	fmt.Println("  - list.SelectionManager: Selection state")
	fmt.Println("  - render.Model: Elm architecture")
	fmt.Println()

	// Create and run the model using DirectEngine
	model := NewSimpleListModel()
	engine := render.NewDirectEngine(nil)

	if err := engine.Start(model); err != nil {
		fmt.Printf("Error starting engine: %v\n", err)
		return
	}
	defer engine.Stop()

	// Simulate some user input
	demoInputs := []render.Msg{
		render.WindowSizeMsg{Width: 60, Height: 15},
		render.KeyMsg{Key: "j"},
		render.KeyMsg{Key: "j"},
		render.KeyMsg{Key: " "},
		render.KeyMsg{Key: "j"},
		render.KeyMsg{Key: "j"},
		render.KeyMsg{Key: " "},
		render.KeyMsg{Key: "G"},
		render.KeyMsg{Key: "q"},
	}

	fmt.Println("--- Initial State ---")
	if d, ok := engine.(*render.DirectEngine); ok {
		fmt.Println(d.Output())
		fmt.Println()
	}

	for _, input := range demoInputs {
		if err := engine.Send(input); err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			continue
		}

		// Show state after each input
		if d, ok := engine.(*render.DirectEngine); ok {
			if km, ok := input.(render.KeyMsg); ok {
				fmt.Printf("--- After key: %s ---\n", km.String())
			} else {
				fmt.Println("--- After input ---")
			}
			fmt.Println(d.Output())
			fmt.Println()
		}
	}

	fmt.Println("--- Final State ---")
	if d, ok := engine.(*render.DirectEngine); ok {
		fmt.Println(d.Output())
	}

	// Show selected items
	if d, ok := engine.(*render.DirectEngine); ok {
		finalModel := d.Model()
		fmt.Println("\n=== Selected Items ===")
		if m, ok := finalModel.(*SimpleListModel); ok {
			for _, item := range m.items {
				if _, ok := m.selected[item.ID()]; ok {
					if li, ok := item.(*list.ListItem); ok {
						fmt.Printf("  - %s: %s\n", li.ID(), li.Title())
					}
				}
			}
		}
	}
}
