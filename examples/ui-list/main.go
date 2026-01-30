// Interactive list example using the new engine-agnostic list components with Bubbletea
package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/list"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// InteractiveListModel wraps the engine-agnostic list components for Bubbletea
type InteractiveListModel struct {
	items    []list.Item
	cursor   int
	selected map[string]struct{}
	viewport *list.Viewport
	selMgr   *list.SelectionManager
	styles   *styles.Styles
	quitting bool
}

func NewInteractiveListModel() *InteractiveListModel {
	s := styles.DefaultStyles()
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
		list.NewListItem("11", "Mango", "A tropical sweet fruit"),
		list.NewListItem("12", "Nectarine", "A smooth-skinned peach"),
		list.NewListItem("13", "Orange", "A citrus fruit"),
		list.NewListItem("14", "Papaya", "A tropical fruit"),
		list.NewListItem("15", "Quince", "A fragrant fruit"),
	}

	return &InteractiveListModel{
		items:    items,
		cursor:   0,
		selected: make(map[string]struct{}),
		viewport: list.NewViewport(8, len(items)),
		selMgr:   list.NewSelectionManager(list.SelectionModeMultiple),
		styles:   &s,
		quitting: false,
	}
}

// Init implements tea.Model
func (m *InteractiveListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *InteractiveListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.quitting {
			return m, tea.Quit
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			m.moveUp()

		case "down", "j":
			m.moveDown()

		case " ", "enter":
			m.toggleSelection()

		case "g":
			m.moveToTop()

		case "G":
			m.moveToBottom()

		case "ctrl+u":
			m.viewport.PageUp()

		case "ctrl+d":
			m.viewport.PageDown()
		}

	case tea.WindowSizeMsg:
		// Window resized - update viewport visible count
		height := msg.Height - 8 // Leave room for header/footer
		if height > 0 {
			m.viewport.SetVisible(height)
		}
	}

	return m, nil
}

// View implements tea.Model
func (m *InteractiveListModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var b strings.Builder

	// Header with gradient
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235")).
		Padding(0, 2).
		MarginBottom(1)

	b.WriteString(headerStyle.Render("╔════════════════════════════════════════════════════════╗"))
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("║         Engine-Agnostic List Demo (v2.0)                ║"))
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("╠════════════════════════════════════════════════════════╣"))
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("║  ↑/k up  ↓/j down  space=select  g=top  G=bottom  q=quit ║"))
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("╚════════════════════════════════════════════════════════╝"))
	b.WriteString("\n")

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

		// Style based on cursor position
		itemStyle := m.styles.Base
		if i == m.cursor {
			itemStyle = m.styles.TextSelection.Foreground(lipgloss.Color("226"))
		} else {
			itemStyle = m.styles.Base.Foreground(lipgloss.Color("252"))
		}

		// Render item
		if li, ok := item.(*list.ListItem); ok {
			line := fmt.Sprintf("%s [%s] %s", cursor, checked, li.Title())
			if li.Desc() != "" {
				line += fmt.Sprintf(" — %s", li.Desc())
			}
			b.WriteString(itemStyle.Render(line) + "\n")
		}
	}

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		MarginTop(1)

	selectedCount := len(m.selected)
	scrollInfo := m.viewport.ScrollIndicator()
	footerText := fmt.Sprintf("Selected: %d | Showing %d-%d of %d | Scroll: %s",
		selectedCount, start+1, end, len(m.items), scrollInfo)

	b.WriteString(footerStyle.Render(footerText))

	return b.String()
}

func (m *InteractiveListModel) moveUp() {
	if m.cursor > 0 {
		m.cursor--
		m.viewport.MoveUp()
	}
}

func (m *InteractiveListModel) moveDown() {
	if m.cursor < len(m.items)-1 {
		m.cursor++
		m.viewport.MoveDown()
	}
}

func (m *InteractiveListModel) toggleSelection() {
	id := m.items[m.cursor].ID()
	if _, ok := m.selected[id]; ok {
		delete(m.selected, id)
	} else {
		m.selected[id] = struct{}{}
	}
}

func (m *InteractiveListModel) moveToTop() {
	m.cursor = 0
	m.viewport.MoveToTop()
}

func (m *InteractiveListModel) moveToBottom() {
	m.cursor = len(m.items) - 1
	m.viewport.MoveToBottom()
}

func main() {
	p := tea.NewProgram(
		NewInteractiveListModel(),
		tea.WithAltScreen(),       // Use alternate screen
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
