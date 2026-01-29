package completions

import (
	"strings"

	"github.com/yourorg/taproot/internal/ui/list"
)

// Provider is a data source for completion items.
type Provider interface {
	// GetItems returns all available completion items.
	GetItems() []CompletionItem
	// GetFilterValue returns the string to filter on for an item.
	GetFilterValue(item CompletionItem) string
}

// CompletionItem represents a single completion option.
type CompletionItem interface {
	// ID returns a unique identifier for the item.
	ID() string
	// Display returns the text to display.
	Display() string
	// Value returns the underlying value to use when selected.
	Value() any
}

// SimpleCompletionItem is a basic implementation of CompletionItem.
type SimpleCompletionItem struct {
	id      string
	display string
	value   any
}

// NewSimpleCompletionItem creates a new SimpleCompletionItem.
func NewSimpleCompletionItem(id, display string, value any) SimpleCompletionItem {
	return SimpleCompletionItem{
		id:      id,
		display: display,
		value:   value,
	}
}

// ID returns the unique identifier.
func (s SimpleCompletionItem) ID() string {
	return s.id
}

// Display returns the display text.
func (s SimpleCompletionItem) Display() string {
	return s.display
}

// Value returns the underlying value.
func (s SimpleCompletionItem) Value() any {
	return s.value
}

// AutoCompletion is an engine-agnostic auto-completion component.
type AutoCompletion struct {
	provider Provider
	list     *list.BaseList
	filter   *list.Filter
	viewport *list.Viewport
	items    []CompletionItem
	filtered []CompletionItem

	open    bool
	query   string
	cursor  int
	width   int
	height  int

	minHeight int
	maxHeight int
	maxWidth  int
}

// NewAutoCompletion creates a new auto-completion component.
func NewAutoCompletion(provider Provider, minHeight, maxHeight, maxWidth int) *AutoCompletion {
	return &AutoCompletion{
		provider:  provider,
		list:      list.NewBaseList(),
		filter:    list.NewFilter(),
		viewport:  list.NewViewport(0, 0),
		items:     []CompletionItem{},
		filtered:  []CompletionItem{},
		open:      false,
		query:     "",
		cursor:    0,
		width:     0,
		height:    0,
		minHeight: minHeight,
		maxHeight: maxHeight,
		maxWidth:  maxWidth,
	}
}

// Open opens the completion popup.
func (a *AutoCompletion) Open() {
	a.items = a.provider.GetItems()
	a.filtered = a.items
	a.open = true
	a.query = ""
	a.cursor = 0
	a.recalcDimensions()
}

// Close closes the completion popup.
func (a *AutoCompletion) Close() {
	a.open = false
	a.items = []CompletionItem{}
	a.filtered = []CompletionItem{}
	a.query = ""
	a.cursor = 0
}

// IsOpen returns whether the popup is open.
func (a *AutoCompletion) IsOpen() bool {
	return a.open
}

// Query returns the current filter query.
func (a *AutoCompletion) Query() string {
	return a.query
}

// SetQuery filters the completions with the given query.
func (a *AutoCompletion) SetQuery(query string) {
	a.query = strings.ToLower(query)
	a.filtered = []CompletionItem{}

	for _, item := range a.items {
		filterValue := strings.ToLower(a.provider.GetFilterValue(item))
		if strings.Contains(filterValue, a.query) {
			a.filtered = append(a.filtered, item)
		}
	}

	// Reset cursor to first item
	if len(a.filtered) > 0 {
		a.cursor = 0
	} else {
		a.cursor = -1
	}

	a.recalcDimensions()
}

// Cursor returns the currently selected index.
func (a *AutoCompletion) Cursor() int {
	return a.cursor
}

// SetCursor sets the selected index.
func (a *AutoCompletion) SetCursor(index int) {
	if index >= 0 && index < len(a.filtered) {
		a.cursor = index
	}
}

// MoveUp moves selection up one item with wraparound.
func (a *AutoCompletion) MoveUp() {
	if len(a.filtered) == 0 {
		return
	}
	a.cursor--
	if a.cursor < 0 {
		a.cursor = len(a.filtered) - 1
	}
}

// MoveDown moves selection down one item with wraparound.
func (a *AutoCompletion) MoveDown() {
	if len(a.filtered) == 0 {
		return
	}
	a.cursor++
	if a.cursor >= len(a.filtered) {
		a.cursor = 0
	}
}

// Selected returns the currently selected item.
func (a *AutoCompletion) Selected() CompletionItem {
	if a.cursor < 0 || a.cursor >= len(a.filtered) {
		return nil
	}
	return a.filtered[a.cursor]
}

// SelectedValue returns the value of the currently selected item.
func (a *AutoCompletion) SelectedValue() any {
	item := a.Selected()
	if item == nil {
		return nil
	}
	return item.Value()
}

// HasItems returns whether there are visible items.
func (a *AutoCompletion) HasItems() bool {
	return len(a.filtered) > 0
}

// Items returns all filtered items.
func (a *AutoCompletion) Items() []CompletionItem {
	return a.filtered
}

// ItemCount returns the number of visible items.
func (a *AutoCompletion) ItemCount() int {
	return len(a.filtered)
}

// Size returns the current dimensions of the popup.
func (a *AutoCompletion) Size() (width, height int) {
	return a.width, a.height
}

// recalcDimensions recalculates the popup dimensions based on items.
func (a *AutoCompletion) recalcDimensions() {
	if !a.open || len(a.filtered) == 0 {
		a.height = 0
		a.width = 0
		return
	}

	// Calculate height
	a.height = len(a.filtered)
	if a.height > a.maxHeight {
		a.height = a.maxHeight
	}
	if a.height < a.minHeight {
		a.height = a.minHeight
	}

	// Calculate width based on display text
	maxDisplay := 0
	for _, item := range a.filtered {
		display := item.Display()
		if len(display) > maxDisplay {
			maxDisplay = len(display)
		}
	}

	a.width = maxDisplay + 4 // Add padding
	if a.width > a.maxWidth {
		a.width = a.maxWidth
	}

	// Update viewport
	a.viewport = list.NewViewport(a.height, len(a.filtered))
	a.viewport.MoveTo(a.cursor)
}

// ScrollToVisible scrolls the viewport to make the cursor visible.
func (a *AutoCompletion) ScrollToVisible() {
	start, end := a.viewport.Range()
	if a.cursor < start {
		a.viewport.MoveUp()
	} else if a.cursor >= end {
		a.viewport.MoveDown()
	}
}

// VisibleRange returns the start and end indices of visible items.
func (a *AutoCompletion) VisibleRange() (start, end int) {
	return a.viewport.Range()
}
