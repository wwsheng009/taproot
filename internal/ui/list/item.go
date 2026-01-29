package list

// ListItem is a basic implementation of FilterableItem.
type ListItem struct {
	id    string
	title string
	desc  string
}

// NewListItem creates a new ListItem.
func NewListItem(id, title, desc string) *ListItem {
	return &ListItem{
		id:    id,
		title: title,
		desc:  desc,
	}
}

// ID returns the item's unique identifier.
func (i *ListItem) ID() string {
	return i.id
}

// Title returns the item's title.
func (i *ListItem) Title() string {
	return i.title
}

// SetTitle sets the item's title.
func (i *ListItem) SetTitle(title string) {
	i.title = title
}

// Desc returns the item's description.
func (i *ListItem) Desc() string {
	return i.desc
}

// SetDesc sets the item's description.
func (i *ListItem) SetDesc(desc string) {
	i.desc = desc
}

// FilterValue returns the combined title and description for filtering.
func (i *ListItem) FilterValue() string {
	return i.title + " " + i.desc
}

// SectionItem implements ItemSection for section headers.
type SectionItem struct {
	id     string
	title  string
	info   string
	width  int
	height int
	index  int
}

// NewSectionItem creates a new section item.
func NewSectionItem(title string) *SectionItem {
	return &SectionItem{
		title:  title,
		info:   "",
		width:  0,
		height: 1,
		index:  -1,
		id:     "section-" + title,
	}
}

// ID returns the section's unique identifier.
func (s *SectionItem) ID() string {
	return s.id
}

// Title returns the section title.
func (s *SectionItem) Title() string {
	return s.title
}

// SetTitle sets the section title.
func (s *SectionItem) SetTitle(title string) {
	s.title = title
}

// SetInfo sets additional info text for the section.
func (s *SectionItem) SetInfo(info string) {
	s.info = info
}

// Info returns the section info.
func (s *SectionItem) Info() string {
	return s.info
}

// Size returns the current dimensions.
func (s *SectionItem) Size() (width, height int) {
	return s.width, s.height
}

// SetSize updates the dimensions.
func (s *SectionItem) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// SetIndex sets the positional index.
func (s *SectionItem) SetIndex(index int) {
	s.index = index
}

// Index returns the positional index.
func (s *SectionItem) Index() int {
	return s.index
}

// IsSectionHeader identifies this as a section header.
func (s *SectionItem) IsSectionHeader() bool {
	return true
}

// SelectableItem is a list item that tracks its selection state.
type SelectableItem struct {
	*ListItem
	selected bool
}

// NewSelectableItem creates a new selectable item.
func NewSelectableItem(id, title, desc string) *SelectableItem {
	return &SelectableItem{
		ListItem: NewListItem(id, title, desc),
		selected: false,
	}
}

// Selected returns the selection state.
func (i *SelectableItem) Selected() bool {
	return i.selected
}

// SetSelected sets the selection state.
func (i *SelectableItem) SetSelected(selected bool) {
	i.selected = selected
}

// Toggle reverses the selection state.
func (i *SelectableItem) Toggle() {
	i.selected = !i.selected
}

// ExpandableItem is a list item with expandable state (for grouped lists).
type ExpandableItem struct {
	*ListItem
	expanded bool
}

// NewExpandableItem creates a new expandable item.
func NewExpandableItem(id, title, desc string) *ExpandableItem {
	return &ExpandableItem{
		ListItem: NewListItem(id, title, desc),
		expanded: false,
	}
}

// Expanded returns the expanded state.
func (i *ExpandableItem) Expanded() bool {
	return i.expanded
}

// SetExpanded sets the expanded state.
func (i *ExpandableItem) SetExpanded(expanded bool) {
	i.expanded = expanded
}

// Toggle reverses the expanded state.
func (i *ExpandableItem) Toggle() {
	i.expanded = !i.expanded
}
