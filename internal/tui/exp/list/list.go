package list

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// SimpleList is a simple list component for demonstration.
type SimpleList struct {
	styles   *styles.Styles
	items    []string
	cursor   int
	selected map[int]struct{}
	width    int
	height   int
	visible  int
	offset   int
	focused  bool
}

// NewSimpleList creates a new simple list.
func NewSimpleList(items []string) *SimpleList {
	s := styles.DefaultStyles()
	return &SimpleList{
		styles:    &s,
		items:     items,
		selected:  make(map[int]struct{}),
		visible:   10,
		focused:   true,
	}
}

func (l *SimpleList) Init() tea.Cmd {
	return nil
}

func (l *SimpleList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !l.focused {
		return l, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if l.cursor > 0 {
				l.cursor--
				if l.cursor < l.offset {
					l.offset--
				}
			}
		case "down", "j":
			if l.cursor < len(l.items)-1 {
				l.cursor++
				if l.cursor >= l.offset+l.visible {
					l.offset++
				}
			}
		case " ", "enter":
			if _, ok := l.selected[l.cursor]; ok {
				delete(l.selected, l.cursor)
			} else {
				l.selected[l.cursor] = struct{}{}
			}
		case "g":
			l.cursor = 0
			l.offset = 0
		case "G":
			l.cursor = len(l.items) - 1
			l.offset = max(0, len(l.items)-l.visible)
		}
	}

	return l, nil
}

func (l *SimpleList) View() string {
	s := l.styles

	var b strings.Builder
	b.WriteString("Simple List\n\n")

	// Calculate visible range
	start := l.offset
	end := min(start+l.visible, len(l.items))

	for i := start; i < end; i++ {
		cursor := " "
		if l.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := l.selected[i]; ok {
			checked = "x"
		}

		itemStyle := s.Base
		if l.cursor == i && l.focused {
			itemStyle = s.TextSelection
		}

		item := itemStyle.Render(fmt.Sprintf("%s [%s] %s", cursor, checked, l.items[i]))
		b.WriteString(item + "\n")
	}

	b.WriteString(fmt.Sprintf("\nSelected: %d items", len(l.selected)))

	return b.String()
}

func (l *SimpleList) Size() (width, height int) {
	return l.width, l.height
}

func (l *SimpleList) SetSize(width, height int) {
	l.width = width
	l.height = height
	l.visible = min(l.height-2, 10) // Leave space for header/footer
}

func (l *SimpleList) Focus() {
	l.focused = true
}

func (l *SimpleList) Blur() {
	l.focused = false
}

func (l *SimpleList) Focused() bool {
	return l.focused
}

func (l *SimpleList) SelectedItems() []string {
	var result []string
	for i := range l.selected {
		result = append(result, l.items[i])
	}
	return result
}
