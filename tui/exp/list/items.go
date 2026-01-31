package list

// Item represents a list item.
type Item interface {
	ID() string
}

// FilterableItem represents an item that can be filtered.
type FilterableItem interface {
	Item
	FilterValue() string
}

// Filterable represents an item that supports filtering.
type Filterable interface {
	// SetFilter sets the filter query and match indexes.
	SetFilter(filter string, matchIndexes []int)
	// FilterValue returns the value to filter against.
	FilterValue() string
	// MatchIndexes returns the current match indexes.
	MatchIndexes() []int
	// ClearFilter clears the filter.
	ClearFilter()
}

// Sizeable represents a component with size.
type Sizeable interface {
	// Size returns the current dimensions.
	Size() (width, height int)
	// SetSize updates the dimensions.
	SetSize(width, height int)
}

// Focusable represents a component that can receive focus.
type Focusable interface {
	// Focus sets the component to receive keyboard input.
	Focus()
	// Blur removes focus from the component.
	Blur()
	// Focused() returns true if the component currently has focus.
	Focused() bool
}

// HasMatchIndexes represents an item that can show match highlights.
type HasMatchIndexes interface {
	// MatchIndexes sets the match indexes.
	MatchIndexes([]int)
}

// Indexable represents an item with an index.
type Indexable interface {
	SetIndex(int)
}

// ItemSection represents a section header in a list.
type ItemSection interface {
	Item
	Sizeable
	Indexable
	SetInfo(info string)
	Title() string
}

type itemSectionModel struct {
	width int
	title string
	inx   int
	id    string
	info  string
}

// NewItemSection creates a new section item.
func NewItemSection(title string) ItemSection {
	return &itemSectionModel{
		title: title,
		inx:   -1,
		id:    "section-" + title,
	}
}

func (m *itemSectionModel) ID() string {
	return m.id
}

func (m *itemSectionModel) Title() string {
	return m.title
}

func (m *itemSectionModel) SetInfo(info string) {
	m.info = info
}

func (m *itemSectionModel) SetIndex(inx int) {
	m.inx = inx
}

func (m *itemSectionModel) Size() (width, height int) {
	return m.width, 1
}

func (m *itemSectionModel) SetSize(width, height int) {
	m.width = width
}

// IsSectionHeader returns true for section headers.
func (m *itemSectionModel) IsSectionHeader() bool {
	return true
}
