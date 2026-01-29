// Filter and Group example using the engine-agnostic list package
package main

import (
	"fmt"
	"strings"

	"github.com/yourorg/taproot/internal/ui/list"
	"github.com/yourorg/taproot/internal/ui/render"
)

// FilterGroupModel demonstrates filter and group functionality
type FilterGroupModel struct {
	allItems   []list.FilterableItem
	filtered   []list.FilterableItem
	groups     []*list.Group
	groupMgr   *list.GroupManager
	filter     *list.Filter
	cursor     int
	width      int
	height     int
	quit       bool
	filterMode bool
	query      string
}

func NewFilterGroupModel() *FilterGroupModel {
	// Create filterable items
	items := []list.FilterableItem{
		list.NewListItem("f1", "Fuji Apple", "Red apple from Japan"),
		list.NewListItem("f2", "Gala Apple", "Sweet red apple"),
		list.NewListItem("f3", "Honeycrisp", "Crisp sweet apple"),
		list.NewListItem("c1", "Navel Orange", "Seedless orange"),
		list.NewListItem("c2", "Blood Orange", "Red-fleshed orange"),
		list.NewListItem("c3", "Mandarin", "Small easy-peel citrus"),
		list.NewListItem("b1", "Cavendish", "Common yellow banana"),
		list.NewListItem("b2", "Plantain", "Cooking banana"),
		list.NewListItem("g1", "Thompson", "Green seedless grape"),
		list.NewListItem("g2", "Concord", "Purple grape"),
	}

	// Create groups
	groups := []*list.Group{
		list.NewGroup("Apples", []list.Item{
			list.NewListItem("f1", "Fuji Apple", "Red apple from Japan"),
			list.NewListItem("f2", "Gala Apple", "Sweet red apple"),
			list.NewListItem("f3", "Honeycrisp", "Crisp sweet apple"),
		}),
		list.NewGroup("Citrus", []list.Item{
			list.NewListItem("c1", "Navel Orange", "Seedless orange"),
			list.NewListItem("c2", "Blood Orange", "Red-fleshed orange"),
			list.NewListItem("c3", "Mandarin", "Small easy-peel citrus"),
		}),
		list.NewGroup("Tropical", []list.Item{
			list.NewListItem("b1", "Cavendish", "Common yellow banana"),
			list.NewListItem("b2", "Plantain", "Cooking banana"),
			list.NewListItem("g1", "Thompson", "Green seedless grape"),
			list.NewListItem("g2", "Concord", "Purple grape"),
		}),
	}

	return &FilterGroupModel{
		allItems:   items,
		filtered:   items,
		groups:     groups,
		groupMgr:   list.NewGroupManager(),
		filter:     list.NewFilter(),
		cursor:     0,
		width:      70,
		height:     15,
		quit:       false,
		filterMode: false,
		query:      "",
	}
}

// Init initializes the model
func (m *FilterGroupModel) Init() error {
	m.groupMgr.SetGroups(m.groups)
	return nil
}

// Update handles incoming messages
func (m *FilterGroupModel) Update(msg any) (render.Model, render.Cmd) {
	switch msg := msg.(type) {
	case render.KeyMsg:
		key := msg.String()

		// Handle filter mode input
		if m.filterMode {
			switch key {
			case "enter":
				m.filterMode = false
			case "esc":
				m.clearFilter()
				m.filterMode = false
			case "backspace", "ctrl+h":
				if len(m.query) > 0 {
					m.query = m.query[:len(m.query)-1]
					m.applyFilter()
				}
			case "ctrl+c", "q":
				m.quit = true
			default:
				// Regular character input
				if len(key) == 1 && key != "/" {
					m.query += key
					m.applyFilter()
				}
			}
			return m, nil
		}

		// Normal mode
		switch key {
		case "/":
			m.filterMode = true
			m.query = ""
		case "j", "down":
			m.moveDown()
		case "k", "up":
			m.moveUp()
		case "enter", " ":
			m.toggleCurrent()
		case "E":
			m.groupMgr.ExpandAll()
		case "W":
			m.groupMgr.CollapseAll()
		case "q", "ctrl+c":
			m.quit = true
		}

	case render.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the model
func (m *FilterGroupModel) View() string {
	if m.quit {
		return "Goodbye!"
	}

	var b strings.Builder

	// Header
	b.WriteString("╔════════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║           Filter & Group List Demo (v2.0)                      ║\n")
	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	if m.filterMode {
		b.WriteString(fmt.Sprintf("║  Filter: %s_                                                   ║\n", m.query))
	} else {
		filterInfo := ""
		if m.query != "" {
			filterInfo = fmt.Sprintf(" (filtered: %d/%d)", len(m.filtered), len(m.allItems))
		}
		b.WriteString(fmt.Sprintf("║  Keys: /filter  j/k=nav  enter=toggle  E=expand  W=collapse%s  ║\n", filterInfo))
	}

	b.WriteString("╚════════════════════════════════════════════════════════════════╝\n\n")

	// Grouped view
	m.renderGroups(&b)

	// Filter info
	if m.query != "" && !m.filterMode {
		b.WriteString(fmt.Sprintf("\n  Filter: \"%s\" | Showing %d of %d items\n", m.query, len(m.filtered), len(m.allItems)))
	}

	return b.String()
}

func (m *FilterGroupModel) renderGroups(b *strings.Builder) {
	count := m.groupMgr.Count()
	maxVisible := m.height - 6
	start := 0
	end := min(count, start+maxVisible)

	for i := start; i < end; i++ {
		isGroup, groupIdx, itemIdx := m.groupMgr.GetItemAt(i)
		cursor := " "
		if i == m.cursor {
			cursor = "→"
		}

		if isGroup {
			// Render group header
			group := m.groupMgr.Groups()[groupIdx]
			expandIcon := "+"
			if group.Expanded() {
				expandIcon = "-"
			}
			line := fmt.Sprintf("%s [%s] %s (%d items)", cursor, expandIcon, group.Title(), group.ItemCount())
			b.WriteString("  " + line + "\n")
		} else {
			// Render item
			group := m.groupMgr.Groups()[groupIdx]
			item := group.Items()[itemIdx]
			if li, ok := item.(*list.ListItem); ok {
				line := fmt.Sprintf("%s    • %s - %s", cursor, li.Title(), li.Desc())
				b.WriteString("  " + line + "\n")
			}
		}
	}

	// Footer info
	if count > maxVisible {
		b.WriteString(fmt.Sprintf("\n  Showing %d-%d of %d (scroll for more)", start+1, end, count))
	}
}

func (m *FilterGroupModel) moveDown() {
	if m.cursor < m.groupMgr.Count()-1 {
		m.cursor++
	}
}

func (m *FilterGroupModel) moveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *FilterGroupModel) toggleCurrent() {
	if m.groupMgr.IsAtGroup() {
		m.groupMgr.ToggleCurrentGroup()
	}
}

func (m *FilterGroupModel) applyFilter() {
	if m.query == "" {
		m.filter.Clear()
		m.filtered = m.allItems
	} else {
		m.filter.SetQuery(m.query)
		m.filtered = m.filter.Apply(m.allItems)

		// Rebuild groups with filtered items
		newGroups := []*list.Group{}
		for _, g := range m.groups {
			filteredItems := []list.Item{}
			for _, item := range g.Items() {
				if fi, ok := item.(list.FilterableItem); ok {
					if m.filter.HasMatchIn(fi.FilterValue()) {
						filteredItems = append(filteredItems, item)
					}
				}
			}
			if len(filteredItems) > 0 {
				newGroups = append(newGroups, list.NewGroup(g.Title(), filteredItems))
			}
		}
		m.groupMgr.SetGroups(newGroups)
	}
	m.cursor = 0
}

func (m *FilterGroupModel) clearFilter() {
	m.query = ""
	m.filter.Clear()
	m.filtered = m.allItems
	m.groupMgr.SetGroups(m.groups)
	m.cursor = 0
}

func main() {
	fmt.Println("=== Filter & Group Demo ===")
	fmt.Println()
	fmt.Println("This demo demonstrates:")
	fmt.Println("  - list.Filter: Filter items with query")
	fmt.Println("  - list.Group: Grouped list with expand/collapse")
	fmt.Println("  - list.GroupManager: Manage group state")
	fmt.Println()

	// Create and run the model
	model := NewFilterGroupModel()
	engine := render.NewDirectEngine(nil)

	if err := engine.Start(model); err != nil {
		fmt.Printf("Error starting engine: %v\n", err)
		return
	}
	defer engine.Stop()

	// Demo sequence
	demoInputs := []render.Msg{
		render.WindowSizeMsg{Width: 70, Height: 20},
		render.KeyMsg{Key: "j"},
		render.KeyMsg{Key: "j"},
		render.KeyMsg{Key: " "}, // Toggle group
		render.KeyMsg{Key: "k"},
		render.KeyMsg{Key: "/"},
		render.KeyMsg{Key: "a"},
		render.KeyMsg{Key: "p"},
		render.KeyMsg{Key: "l"},
		render.KeyMsg{Key: "e"},
		render.KeyMsg{Key: "enter"},
		render.KeyMsg{Key: "W"},
		render.KeyMsg{Key: "E"},
		render.KeyMsg{Key: "q"},
	}

	for _, input := range demoInputs {
		if err := engine.Send(input); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if d, ok := engine.(*render.DirectEngine); ok {
			if km, ok := input.(render.KeyMsg); ok {
				fmt.Printf("--- Key: %s ---\n", km.String())
			}
			fmt.Println(d.Output())
			fmt.Println()
		}
	}

	fmt.Println("=== Demo Complete ===")
	fmt.Println("\nTry running again and:")
	fmt.Println("  - Press / to enter filter mode, type to filter")
	fmt.Println("  - Press enter on a group to toggle it")
	fmt.Println("  - Press E to expand all groups")
	fmt.Println("  - Press W to collapse all groups")
}
