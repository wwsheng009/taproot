// Interactive filter and group example using Bubbletea
package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/list"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// InteractiveFilterGroupModel demonstrates filter and group with Bubbletea
type InteractiveFilterGroupModel struct {
	allItems   []list.FilterableItem
	filtered   []list.FilterableItem
	groups     []*list.Group
	groupMgr   *list.GroupManager
	filter     *list.Filter
	cursor     int
	styles     *styles.Styles
	quitting   bool
	filterMode bool
	query      string
}

func NewInteractiveFilterGroupModel() *InteractiveFilterGroupModel {
	s := styles.DefaultStyles()

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
		list.NewGroup("🍎 Apples", []list.Item{
			list.NewListItem("f1", "Fuji Apple", "Red apple from Japan"),
			list.NewListItem("f2", "Gala Apple", "Sweet red apple"),
			list.NewListItem("f3", "Honeycrisp", "Crisp sweet apple"),
		}),
		list.NewGroup("🍊 Citrus", []list.Item{
			list.NewListItem("c1", "Navel Orange", "Seedless orange"),
			list.NewListItem("c2", "Blood Orange", "Red-fleshed orange"),
			list.NewListItem("c3", "Mandarin", "Small easy-peel citrus"),
		}),
		list.NewGroup("🌴 Tropical", []list.Item{
			list.NewListItem("b1", "Cavendish", "Common yellow banana"),
			list.NewListItem("b2", "Plantain", "Cooking banana"),
			list.NewListItem("g1", "Thompson", "Green seedless grape"),
			list.NewListItem("g2", "Concord", "Purple grape"),
		}),
	}

	return &InteractiveFilterGroupModel{
		allItems:   items,
		filtered:   items,
		groups:     groups,
		groupMgr:   list.NewGroupManager(),
		filter:     list.NewFilter(),
		cursor:     0,
		styles:     &s,
		quitting:   false,
		filterMode: false,
		query:      "",
	}
}

// Init implements tea.Model
func (m *InteractiveFilterGroupModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *InteractiveFilterGroupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.quitting {
			return m, tea.Quit
		}

		// Handle filter mode input
		if m.filterMode {
			switch msg.String() {
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
				m.quitting = true
				return m, tea.Quit
			default:
				// Regular character input
				if len(msg.String()) == 1 && msg.String() != "/" {
					m.query += msg.String()
					m.applyFilter()
				}
			}
			return m, nil
		}

		// Normal mode
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

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
		}

	case tea.WindowSizeMsg:
		// Window resize - can adjust viewport here if needed
	}

	return m, nil
}

// View implements tea.Model
func (m *InteractiveFilterGroupModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var b strings.Builder

	// Styles
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235")).
		Padding(0, 2).
		MarginBottom(1)

	groupStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("228")).
		Bold(true)

	itemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	filterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Background(lipgloss.Color("235"))

	// Header
	b.WriteString(headerStyle.Render("╔════════════════════════════════════════════════════════════════╗"))
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("║              Filter & Group List Demo (v2.0)                    ║"))
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("╠════════════════════════════════════════════════════════════════╣"))
	b.WriteString("\n")

	if m.filterMode {
		filterText := fmt.Sprintf("║  Filter: %s_%-51s║", m.query, "")
		b.WriteString(filterStyle.Render(filterText))
		b.WriteString("\n")
	} else {
		filterInfo := ""
		if m.query != "" {
			filterInfo = fmt.Sprintf(" (filtered: %d/%d)", len(m.filtered), len(m.allItems))
		}
		helpText := fmt.Sprintf("║  /filter j/k nav enter=toggle E=expand W=collapse%s%-20s║", filterInfo, "")
		b.WriteString(headerStyle.Render(helpText))
		b.WriteString("\n")
	}

	b.WriteString(headerStyle.Render("╚════════════════════════════════════════════════════════════════╝"))
	b.WriteString("\n")

	// Grouped view
	m.renderGroups(&b, cursorStyle, groupStyle, itemStyle)

	// Filter info footer
	if m.query != "" && !m.filterMode {
		footerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			MarginTop(1)
		footerText := fmt.Sprintf("Filter: \"%s\" | Showing %d of %d items",
			m.query, len(m.filtered), len(m.allItems))
		b.WriteString(footerStyle.Render(footerText))
	}

	return b.String()
}

func (m *InteractiveFilterGroupModel) renderGroups(
	b *strings.Builder,
	cursorStyle, groupStyle, itemStyle lipgloss.Style,
) {
	count := m.groupMgr.Count()
	maxVisible := 12
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
			expandIcon := "▼"
			if !group.Expanded() {
				expandIcon = "▶"
			}
			line := fmt.Sprintf("%s [%s] %s (%d items)", cursor, expandIcon, group.Title(), group.ItemCount())

			if i == m.cursor {
				b.WriteString(cursorStyle.Render("  "+line) + "\n")
			} else {
				b.WriteString(groupStyle.Render("  "+line) + "\n")
			}
		} else {
			// Render item
			group := m.groupMgr.Groups()[groupIdx]
			item := group.Items()[itemIdx]
			if li, ok := item.(*list.ListItem); ok {
				line := fmt.Sprintf("%s    • %s — %s", cursor, li.Title(), li.Desc())

				if i == m.cursor {
					b.WriteString(cursorStyle.Render("  "+line) + "\n")
				} else {
					b.WriteString(itemStyle.Render("  "+line) + "\n")
				}
			}
		}
	}

	// Scroll indicator
	if count > maxVisible {
		footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
		scrollText := fmt.Sprintf("  Showing %d-%d of %d (scroll for more)", start+1, end, count)
		b.WriteString("\n" + footerStyle.Render(scrollText))
	}
}

func (m *InteractiveFilterGroupModel) moveDown() {
	if m.cursor < m.groupMgr.Count()-1 {
		m.cursor++
	}
}

func (m *InteractiveFilterGroupModel) moveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *InteractiveFilterGroupModel) toggleCurrent() {
	if m.groupMgr.IsAtGroup() {
		m.groupMgr.ToggleCurrentGroup()
	}
}

func (m *InteractiveFilterGroupModel) applyFilter() {
	if m.query == "" {
		m.filter.Clear()
		m.filtered = m.allItems
		m.groupMgr.SetGroups(m.groups)
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

func (m *InteractiveFilterGroupModel) clearFilter() {
	m.query = ""
	m.filter.Clear()
	m.filtered = m.allItems
	m.groupMgr.SetGroups(m.groups)
	m.cursor = 0
}

func main() {
	p := tea.NewProgram(
		NewInteractiveFilterGroupModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
