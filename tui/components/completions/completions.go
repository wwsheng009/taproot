package completions

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/util"
	"github.com/wwsheng009/taproot/ui/styles"
)

const maxCompletionsHeight = 10

// CompletionItem represents a single completion item
type CompletionItem struct {
	id    string
	title string
	value any
}

// NewCompletionItem creates a new completion item
func NewCompletionItem(id, title string, value any) CompletionItem {
	return CompletionItem{
		id:    id,
		title: title,
		value: value,
	}
}

// ID returns the unique identifier for the item
func (c CompletionItem) ID() string {
	return c.id
}

// Title returns the display title
func (c CompletionItem) Title() string {
	return c.title
}

// Value returns the underlying value
func (c CompletionItem) Value() any {
	return c.value
}

// Messages

// OpenCompletionsMsg is sent to open completions with items
type OpenCompletionsMsg struct {
	Completions []CompletionItem
	X           int
	Y           int
	MaxResults  int
}

// FilterCompletionsMsg is sent to filter completions
type FilterCompletionsMsg struct {
	Query  string
	Reopen bool
	X      int
	Y      int
}

// CloseCompletionsMsg is sent to close completions
type CloseCompletionsMsg struct{}

// CompletionsClosedMsg is sent when completions are closed
type CompletionsClosedMsg struct{}

// SelectCompletionMsg is sent when a completion is selected
type SelectCompletionMsg struct {
	Value  any
	Insert bool
}

// CompletionsCmp represents the completions component
type CompletionsCmp interface {
	util.Model
	Open() bool
	Query() string
	Position() (int, int)
	Width() int
	Height() int
}

type completionsCmp struct {
	styles          *styles.Styles
	wWidth, wHeight int
	width, height   int
	x, y            int
	open            bool
	query           string

	items         []CompletionItem
	filteredItems []CompletionItem
	selectedIndex int
	maxResults    int

	normalStyle  lipgloss.Style
	focusedStyle lipgloss.Style
}

// New creates a new completions component
func New() CompletionsCmp {
	s := styles.DefaultStyles()
	return &completionsCmp{
		styles:        &s,
		width:         0,
		height:        maxCompletionsHeight,
		open:          false,
		query:         "",
		items:         []CompletionItem{},
		selectedIndex: -1,
		maxResults:    0,
		normalStyle:   s.Base,
		focusedStyle:  s.TextSelection,
	}
}

func (c *completionsCmp) Init() tea.Cmd {
	return nil
}

func (c *completionsCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.wWidth, c.wHeight = msg.Width, msg.Height
		return c, nil

	case OpenCompletionsMsg:
		c.open = true
		c.items = msg.Completions
		c.filteredItems = msg.Completions
		c.x, c.y = msg.X, msg.Y
		c.maxResults = msg.MaxResults
		c.query = ""
		c.selectedIndex = -1
		if len(c.filteredItems) > 0 {
			c.selectedIndex = 0
		}
		return c, nil

	case FilterCompletionsMsg:
		if !c.open && !msg.Reopen {
			return c, nil
		}
		c.query = strings.ToLower(msg.Query)
		c.x, c.y = msg.X, msg.Y
		c.filterItems()
		if c.open {
			c.selectedIndex = 0
			if len(c.filteredItems) > 0 {
				c.selectedIndex = 0
			} else {
				c.selectedIndex = -1
			}
		}
		return c, nil

	case CloseCompletionsMsg:
		c.open = false
		c.items = []CompletionItem{}
		c.filteredItems = []CompletionItem{}
		c.selectedIndex = -1
		return c, util.CmdHandler(CompletionsClosedMsg{})

	case tea.KeyMsg:
		if !c.open {
			return c, nil
		}

		switch msg.String() {
		case "up", "ctrl+p":
			return c.selectPrev()
		case "down", "ctrl+n":
			return c.selectNext()
		case "enter", "tab", "ctrl+y":
			return c.selectCurrent(false)
		case "esc":
			c.open = false
			return c, util.CmdHandler(CompletionsClosedMsg{})
		}
	}

	return c, nil
}

func (c *completionsCmp) filterItems() {
	if c.query == "" {
		c.filteredItems = c.items
		return
	}

	c.filteredItems = []CompletionItem{}
	for _, item := range c.items {
		if strings.Contains(strings.ToLower(item.title), c.query) {
			c.filteredItems = append(c.filteredItems, item)
		}
	}
}

func (c *completionsCmp) selectPrev() (util.Model, tea.Cmd) {
	if len(c.filteredItems) == 0 {
		return c, nil
	}
	c.selectedIndex--
	if c.selectedIndex < 0 {
		c.selectedIndex = len(c.filteredItems) - 1
	}
	return c, nil
}

func (c *completionsCmp) selectNext() (util.Model, tea.Cmd) {
	if len(c.filteredItems) == 0 {
		return c, nil
	}
	c.selectedIndex++
	if c.selectedIndex >= len(c.filteredItems) {
		c.selectedIndex = 0
	}
	return c, nil
}

func (c *completionsCmp) selectCurrent(insert bool) (util.Model, tea.Cmd) {
	if c.selectedIndex < 0 || c.selectedIndex >= len(c.filteredItems) {
		return c, nil
	}
	item := c.filteredItems[c.selectedIndex]
	c.open = false
	c.filteredItems = []CompletionItem{}
	c.items = []CompletionItem{}
	c.selectedIndex = -1
	return c, util.CmdHandler(SelectCompletionMsg{
		Value:  item.value,
		Insert: insert,
	})
}

func (c *completionsCmp) View() string {
	if !c.open || len(c.filteredItems) == 0 {
		return ""
	}

	s := c.styles

	// Calculate visible items
	visibleCount := len(c.filteredItems)
	if c.maxResults > 0 {
		visibleCount = min(visibleCount, c.maxResults)
	}
	visibleCount = min(visibleCount, c.height)

	// Calculate offset to keep selected item visible
	offset := 0
	if c.selectedIndex >= visibleCount {
		offset = c.selectedIndex - visibleCount + 1
	}

	var b strings.Builder
	b.WriteString(s.Base.Foreground(s.FgMuted).Render("┌─ Completions ─────────┐") + "\n")

	for i := 0; i < visibleCount; i++ {
		idx := offset + i
		if idx >= len(c.filteredItems) {
			break
		}

		item := c.filteredItems[idx]
		prefix := " "
		if idx == c.selectedIndex {
			prefix = ">"
		}

		itemStyle := c.normalStyle
		if idx == c.selectedIndex {
			itemStyle = c.focusedStyle
		}

		// Highlight matching characters
		title := item.title
		if c.query != "" {
			title = c.highlightMatches(item.title, c.query)
		}

		line := lipgloss.NewStyle().Width(c.width).Render(itemStyle.Render(prefix + " " + title))
		b.WriteString("│" + line + "│\n")
	}

	b.WriteString(s.Base.Foreground(s.FgMuted).Render("└─────────────────────┘"))

	return b.String()
}

func (c *completionsCmp) highlightMatches(text, query string) string {
	if query == "" {
		return text
	}

	s := c.styles
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	// Guard against potential unicode length mismatches
	if len(lowerText) != len(text) {
		return text
	}

	var sb strings.Builder
	sb.Grow(len(text) + len(text)/2) // Pre-allocate with some buffer

	i := 0
	for i < len(text) {
		if i <= len(text)-len(query) && lowerText[i:i+len(query)] == lowerQuery {
			sb.WriteString(s.Base.Foreground(s.Primary).Bold(true).Render(text[i : i+len(query)]))
			i += len(query)
		} else {
			sb.WriteByte(text[i])
			i++
		}
	}
	return sb.String()
}

func (c *completionsCmp) Open() bool {
	return c.open
}

func (c *completionsCmp) Query() string {
	return c.query
}

func (c *completionsCmp) Position() (int, int) {
	return c.x, c.y
}

func (c *completionsCmp) Width() int {
	return c.width
}

func (c *completionsCmp) Height() int {
	return c.height
}
