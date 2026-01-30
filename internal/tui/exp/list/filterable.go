package list

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// FilterableList is a list with filtering support
type FilterableList struct {
	styles      *styles.Styles
	allItems    []ListItem
	filtered    []ListItem
	cursor      int
	selected    map[string]struct{}
	width       int
	height      int
	visible     int
	offset      int
	focused     bool
	query       string
	filtering   bool
}

// ListItem represents an item in the filterable list
type ListItem struct {
	ID    string
	Title string
	Desc  string
}

// NewFilterableList creates a new filterable list
func NewFilterableList(items []ListItem) *FilterableList {
	s := styles.DefaultStyles()
	return &FilterableList{
		styles:    &s,
		allItems:  items,
		filtered:  items,
		selected:  make(map[string]struct{}),
		visible:   10,
		focused:   true,
		filtering: false,
	}
}

func (l *FilterableList) Init() tea.Cmd {
	return nil
}

func (l *FilterableList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !l.focused {
		return l, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If in filtering mode, capture all input except navigation
		if l.filtering {
			switch msg.String() {
			case "enter":
				l.filtering = false
				return l, nil
			case "esc":
				l.query = ""
				l.applyFilter()
				l.filtering = false
				return l, nil
			case "backspace", "ctrl+h":
				if len(l.query) > 0 {
					l.query = l.query[:len(l.query)-1]
					l.applyFilter()
				}
				return l, nil
			case "ctrl+c", "q":
				return l, tea.Quit
			default:
				// Regular character input
				if len(msg.String()) == 1 {
					l.query += msg.String()
					l.applyFilter()
				}
				return l, nil
			}
		}

		// Normal mode
		switch msg.String() {
		case "/":
			// Enter filtering mode
			l.filtering = true
			l.query = ""
			return l, nil
		case "up", "k":
			l.moveUp()
		case "down", "j":
			l.moveDown()
		case " ", "enter":
			l.toggleSelected()
		case "g":
			l.cursor = 0
			l.offset = 0
		case "G":
			l.cursor = len(l.filtered) - 1
			l.offset = max(0, len(l.filtered)-l.visible)
		}
	}

	return l, nil
}

func (l *FilterableList) moveUp() {
	if l.cursor > 0 {
		l.cursor--
		if l.cursor < l.offset {
			l.offset--
		}
	}
}

func (l *FilterableList) moveDown() {
	if l.cursor < len(l.filtered)-1 {
		l.cursor++
		if l.cursor >= l.offset+l.visible {
			l.offset++
		}
	}
}

func (l *FilterableList) toggleSelected() {
	if len(l.filtered) == 0 {
		return
	}
	item := l.filtered[l.cursor]
	if _, ok := l.selected[item.ID]; ok {
		delete(l.selected, item.ID)
	} else {
		l.selected[item.ID] = struct{}{}
	}
}

func (l *FilterableList) applyFilter() {
	if l.query == "" {
		l.filtered = l.allItems
	} else {
		query := strings.ToLower(l.query)
		l.filtered = []ListItem{}
		for _, item := range l.allItems {
			if strings.Contains(strings.ToLower(item.Title), query) ||
				strings.Contains(strings.ToLower(item.Desc), query) {
				l.filtered = append(l.filtered, item)
			}
		}
	}
	// Reset cursor
	l.cursor = 0
	l.offset = 0
}

func (l *FilterableList) View() string {
	s := l.styles

	var b strings.Builder

	// Header
	headerStyle := s.Base.Bold(true).Foreground(s.Primary)
	if l.filtering {
		b.WriteString(headerStyle.Render("Filter: " + l.query + "_"))
	} else {
		filteredCount := ""
		if l.query != "" {
			filteredCount = fmt.Sprintf(" (%d/%d)", len(l.filtered), len(l.allItems))
		}
		b.WriteString(headerStyle.Render(fmt.Sprintf("Items%s - Press / to filter", filteredCount)))
	}
	b.WriteString("\n\n")

	// Calculate visible range
	start := l.offset
	end := min(start+l.visible, len(l.filtered))

	for i := start; i < end; i++ {
		item := l.filtered[i]
		cursor := " "
		if l.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := l.selected[item.ID]; ok {
			checked = "x"
		}

		itemStyle := s.Base
		if l.cursor == i && l.focused {
			itemStyle = s.TextSelection
		}

		// Highlight matching text
		title := item.Title
		if l.query != "" {
			title = l.highlightMatches(item.Title, l.query)
		}

		line := fmt.Sprintf("%s [%s] %s", cursor, checked, title)
		if item.Desc != "" {
			line += fmt.Sprintf(": %s", item.Desc)
		}

		b.WriteString(itemStyle.Render(line) + "\n")
	}

	// Footer
	footer := fmt.Sprintf("Selected: %d items", len(l.selected))
	if len(l.filtered) > l.visible {
		footer += fmt.Sprintf(" | Showing %d-%d of %d", start+1, end, len(l.filtered))
	}
	b.WriteString("\n" + s.Base.Foreground(s.FgMuted).Render(footer))

	return b.String()
}

func (l *FilterableList) highlightMatches(text, query string) string {
	if query == "" {
		return text
	}

	s := l.styles
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	result := ""
	i := 0
	for i < len(text) {
		matchStart := strings.Index(lowerText[i:], lowerQuery)
		if matchStart == -1 {
			result += text[i:]
			break
		}
		matchStart += i
		result += text[i:matchStart]
		result += s.Base.Foreground(s.Primary).Bold(true).Render(text[matchStart:matchStart+len(query)])
		i = matchStart + len(query)
		lowerText = strings.ToLower(text[i:])
	}
	return result
}

func (l *FilterableList) Size() (width, height int) {
	return l.width, l.height
}

func (l *FilterableList) SetSize(width, height int) {
	l.width = width
	l.height = height
	l.visible = max(1, height-4) // Leave space for header/footer
}

func (l *FilterableList) Focus() {
	l.focused = true
}

func (l *FilterableList) Blur() {
	l.focused = false
}

func (l *FilterableList) Focused() bool {
	return l.focused
}

// SelectedIDs returns the IDs of selected items
func (l *FilterableList) SelectedIDs() []string {
	var result []string
	for id := range l.selected {
		result = append(result, id)
	}
	return result
}

// SelectedItems returns the full selected items
func (l *FilterableList) SelectedItems() []ListItem {
	var result []ListItem
	for _, item := range l.allItems {
		if _, ok := l.selected[item.ID]; ok {
			result = append(result, item)
		}
	}
	return result
}

// SetItems sets a new list of items
func (l *FilterableList) SetItems(items []ListItem) {
	l.allItems = items
	l.applyFilter()
}

// Query returns the current filter query
func (l *FilterableList) Query() string {
	return l.query
}

// Filtering returns true if in filtering mode
func (l *FilterableList) Filtering() bool {
	return l.filtering
}
