package list

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/styles"
)

// Group represents a group of items with a section header
type Group struct {
	Title    string
	Items    []ListItem
	Expanded bool
}

// GroupedList is a list with group support
type GroupedList struct {
	styles    *styles.Styles
	groups    []Group
	flatItems []flatItem
	cursor    int
	selected  map[string]struct{}
	width     int
	height    int
	visible   int
	offset    int
	focused   bool
}

// flatItem represents either a group header or a list item
type flatItem struct {
	isGroup  bool
	groupIdx int
	itemIdx  int
	group    *Group
	listItem *ListItem
}

// NewGroupedList creates a new grouped list
func NewGroupedList(groups []Group) *GroupedList {
	s := styles.DefaultStyles()
	gl := &GroupedList{
		styles:   &s,
		groups:   groups,
		selected: make(map[string]struct{}),
		visible:  10,
		focused:  true,
	}
	gl.flatten()
	return gl
}

func (l *GroupedList) Init() tea.Cmd {
	return nil
}

func (l *GroupedList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !l.focused {
		return l, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			l.moveUp()
		case "down", "j":
			l.moveDown()
		case " ", "enter":
			l.toggleCurrent()
		case "g":
			l.cursor = 0
			l.offset = 0
		case "G":
			l.cursor = len(l.flatItems) - 1
			l.offset = max(0, len(l.flatItems)-l.visible)
		}
	}

	return l, nil
}

func (l *GroupedList) moveUp() {
	if l.cursor > 0 {
		l.cursor--
		if l.cursor < l.offset {
			l.offset--
		}
	}
}

func (l *GroupedList) moveDown() {
	if l.cursor < len(l.flatItems)-1 {
		l.cursor++
		if l.cursor >= l.offset+l.visible {
			l.offset++
		}
	}
}

func (l *GroupedList) toggleCurrent() {
	if l.cursor >= len(l.flatItems) {
		return
	}

	item := l.flatItems[l.cursor]

	// If it's a group header, toggle expansion
	if item.isGroup {
		group := l.groups[item.groupIdx]
		group.Expanded = !group.Expanded
		l.groups[item.groupIdx] = group
		l.flatten()
		// Adjust cursor if needed
		if l.cursor >= len(l.flatItems) {
			l.cursor = len(l.flatItems) - 1
		}
		return
	}

	// If it's a list item, toggle selection
	listItem := item.listItem
	if _, ok := l.selected[listItem.ID]; ok {
		delete(l.selected, listItem.ID)
	} else {
		l.selected[listItem.ID] = struct{}{}
	}
}

func (l *GroupedList) flatten() {
	l.flatItems = []flatItem{}
	for gi, group := range l.groups {
		// Add group header
		l.flatItems = append(l.flatItems, flatItem{
			isGroup:  true,
			groupIdx: gi,
			group:    &l.groups[gi],
		})

		// Add items if group is expanded
		if group.Expanded {
			for li := range group.Items {
				l.flatItems = append(l.flatItems, flatItem{
					isGroup:  false,
					groupIdx: gi,
					itemIdx:  li,
					group:    &l.groups[gi],
					listItem: &group.Items[li],
				})
			}
		}
	}
}

func (l *GroupedList) View() string {
	s := l.styles

	var b strings.Builder

	// Header
	headerStyle := s.Base.Bold(true).Foreground(s.Primary)
	expandedCount := 0
	totalItems := 0
	for _, g := range l.groups {
		if g.Expanded {
			expandedCount++
			totalItems += len(g.Items)
		}
	}
	b.WriteString(headerStyle.Render(fmt.Sprintf("Groups (%d/%d expanded, %d items)",
		expandedCount, len(l.groups), totalItems)))
	b.WriteString("\n\n")

	// Calculate visible range
	start := l.offset
	end := min(start+l.visible, len(l.flatItems))

	for i := start; i < end; i++ {
		item := l.flatItems[i]
		cursor := " "
		if l.cursor == i {
			cursor = ">"
		}

		itemStyle := s.Base
		if l.cursor == i && l.focused {
			itemStyle = s.TextSelection
		}

		if item.isGroup {
			// Render group header
			expandIcon := "+"
			if item.group.Expanded {
				expandIcon = "-"
			}
			groupLine := fmt.Sprintf("%s [%s] %s (%d items)", cursor, expandIcon,
				item.group.Title, len(item.group.Items))
			b.WriteString(itemStyle.Bold(true).Render(groupLine) + "\n")
		} else {
			// Render list item
			checked := " "
			if _, ok := l.selected[item.listItem.ID]; ok {
				checked = "x"
			}
			indent := "  "
			line := fmt.Sprintf("%s%s [%s] %s", cursor, indent, checked, item.listItem.Title)
			if item.listItem.Desc != "" {
				line += fmt.Sprintf(": %s", item.listItem.Desc)
			}
			b.WriteString(itemStyle.Render(line) + "\n")
		}
	}

	// Footer
	selectedCount := len(l.selected)
	footer := fmt.Sprintf("Selected: %d items", selectedCount)
	if len(l.flatItems) > l.visible {
		footer += fmt.Sprintf(" | Showing %d-%d of %d", start+1, end, len(l.flatItems))
	}
	b.WriteString("\n" + s.Base.Foreground(s.FgMuted).Render(footer))

	return b.String()
}

func (l *GroupedList) Size() (width, height int) {
	return l.width, l.height
}

func (l *GroupedList) SetSize(width, height int) {
	l.width = width
	l.height = height
	l.visible = max(1, height-4) // Leave space for header/footer
}

func (l *GroupedList) Focus() {
	l.focused = true
}

func (l *GroupedList) Blur() {
	l.focused = false
}

func (l *GroupedList) Focused() bool {
	return l.focused
}

// SelectedIDs returns the IDs of selected items
func (l *GroupedList) SelectedIDs() []string {
	var result []string
	for id := range l.selected {
		result = append(result, id)
	}
	return result
}

// SelectedItems returns all selected items
func (l *GroupedList) SelectedItems() []ListItem {
	var result []ListItem
	for _, group := range l.groups {
		for _, item := range group.Items {
			if _, ok := l.selected[item.ID]; ok {
				result = append(result, item)
			}
		}
	}
	return result
}

// SetGroups sets new groups
func (l *GroupedList) SetGroups(groups []Group) {
	l.groups = groups
	// Ensure all groups are expanded by default
	for i := range l.groups {
		if !l.groups[i].Expanded {
			l.groups[i].Expanded = true
		}
	}
	l.flatten()
}

// ToggleGroup toggles a group's expanded state
func (l *GroupedList) ToggleGroup(groupIdx int) {
	if groupIdx >= 0 && groupIdx < len(l.groups) {
		group := l.groups[groupIdx]
		group.Expanded = !group.Expanded
		l.groups[groupIdx] = group
		l.flatten()
	}
}
