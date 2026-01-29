package list

// Item represents a list item with a unique identifier.
type Item interface {
	// ID returns a unique identifier for this item.
	ID() string
}

// FilterableItem represents an item that can be filtered by a search query.
type FilterableItem interface {
	Item
	// FilterValue returns the string value used for filtering.
	FilterValue() string
}

// Filterable represents an item that supports filter state management.
type Filterable interface {
	// SetFilter sets the filter query and match indexes for highlighting.
	SetFilter(filter string, matchIndexes []int)
	// FilterValue returns the value to filter against.
	FilterValue() string
	// MatchIndexes returns the current match indexes.
	MatchIndexes() []int
	// ClearFilter clears the filter state.
	ClearFilter()
}

// Sizeable represents a component with configurable dimensions.
type Sizeable interface {
	// Size returns the current dimensions as (width, height).
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
	// Focused returns true if the component currently has focus.
	Focused() bool
}

// HasMatchIndexes represents an item that can show match highlights.
type HasMatchIndexes interface {
	// MatchIndexes sets the match indexes for highlighting.
	MatchIndexes([]int)
}

// Indexable represents an item with a positional index.
type Indexable interface {
	SetIndex(int)
}

// ItemSection represents a section header in a list.
type ItemSection interface {
	Item
	Sizeable
	Indexable
	// SetInfo sets additional info text for the section.
	SetInfo(info string)
	// Title returns the section title.
	Title() string
}

// Selectable represents an item that can be selected/deselected.
type Selectable interface {
	Item
	// Selected returns true if the item is selected.
	Selected() bool
	// SetSelected sets the selection state.
	SetSelected(bool)
}

// Toggleable represents an item with expandable/collapsible state.
type Toggleable interface {
	Item
	// Expanded returns true if the item is expanded.
	Expanded() bool
	// SetExpanded sets the expanded state.
	SetExpanded(bool)
}
