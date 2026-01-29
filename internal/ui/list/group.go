package list

// Group represents a group of items with a header.
type Group struct {
	title    string
	items    []Item
	expanded bool
}

// NewGroup creates a new group.
func NewGroup(title string, items []Item) *Group {
	return &Group{
		title:    title,
		items:    items,
		expanded: true,
	}
}

// Title returns the group title.
func (g *Group) Title() string {
	return g.title
}

// SetTitle sets the group title.
func (g *Group) SetTitle(title string) {
	g.title = title
}

// Items returns the group's items.
func (g *Group) Items() []Item {
	return g.items
}

// SetItems sets the group's items.
func (g *Group) SetItems(items []Item) {
	g.items = items
}

// Expanded returns true if the group is expanded.
func (g *Group) Expanded() bool {
	return g.expanded
}

// SetExpanded sets the expanded state.
func (g *Group) SetExpanded(expanded bool) {
	g.expanded = expanded
}

// Toggle toggles the expanded state.
func (g *Group) Toggle() {
	g.expanded = !g.expanded
}

// ItemCount returns the number of items in the group.
func (g *Group) ItemCount() int {
	return len(g.items)
}

// GroupManager manages a list of groups with flattened view.
type GroupManager struct {
	groups     []*Group
	flatItems  []flatGroupItem
	cursor     int
}

// flatGroupItem represents either a group header or an item.
type flatGroupItem struct {
	isGroup  bool
	groupIdx int
	itemIdx  int
}

// NewGroupManager creates a new group manager.
func NewGroupManager() *GroupManager {
	return &GroupManager{
		groups:     nil,
		flatItems:  nil,
		cursor:     0,
	}
}

// Groups returns the groups.
func (gm *GroupManager) Groups() []*Group {
	return gm.groups
}

// SetGroups sets the groups and rebuilds the flattened view.
func (gm *GroupManager) SetGroups(groups []*Group) {
	gm.groups = groups
	gm.flatten()
	gm.ClampCursor()
}

// AddGroup adds a group to the manager.
func (gm *GroupManager) AddGroup(group *Group) {
	gm.groups = append(gm.groups, group)
	gm.flatten()
}

// Clear removes all groups.
func (gm *GroupManager) Clear() {
	gm.groups = nil
	gm.flatItems = nil
	gm.cursor = 0
}

// Count returns the total number of flattened items (headers + visible items).
func (gm *GroupManager) Count() int {
	return len(gm.flatItems)
}

// VisibleItemCount returns the count of visible (non-header) items.
func (gm *GroupManager) VisibleItemCount() int {
	count := 0
	for _, fi := range gm.flatItems {
		if !fi.isGroup {
			count++
		}
	}
	return count
}

// TotalItemCount returns the total number of items across all groups.
func (gm *GroupManager) TotalItemCount() int {
	count := 0
	for _, g := range gm.groups {
		count += len(g.items)
	}
	return count
}

// GroupCount returns the number of groups.
func (gm *GroupManager) GroupCount() int {
	return len(gm.groups)
}

// ExpandedGroupCount returns the number of expanded groups.
func (gm *GroupManager) ExpandedGroupCount() int {
	count := 0
	for _, g := range gm.groups {
		if g.expanded {
			count++
		}
	}
	return count
}

// flatten rebuilds the flattened item list.
func (gm *GroupManager) flatten() {
	gm.flatItems = nil
	for gi, group := range gm.groups {
		// Add group header
		gm.flatItems = append(gm.flatItems, flatGroupItem{
			isGroup:  true,
			groupIdx: gi,
			itemIdx:  -1,
		})

		// Add items if expanded
		if group.expanded {
			for li := range group.items {
				gm.flatItems = append(gm.flatItems, flatGroupItem{
					isGroup:  false,
					groupIdx: gi,
					itemIdx:  li,
				})
			}
		}
	}
}

// Cursor returns the current cursor position.
func (gm *GroupManager) Cursor() int {
	return gm.cursor
}

// SetCursor sets the cursor position and clamps it.
func (gm *GroupManager) SetCursor(cursor int) {
	gm.cursor = cursor
	gm.ClampCursor()
}

// MoveUp moves the cursor up.
func (gm *GroupManager) MoveUp() {
	if gm.cursor > 0 {
		gm.cursor--
	}
}

// MoveDown moves the cursor down.
func (gm *GroupManager) MoveDown() {
	if gm.cursor < len(gm.flatItems)-1 {
		gm.cursor++
	}
}

// ClampCursor ensures the cursor is within valid bounds.
func (gm *GroupManager) ClampCursor() {
	if len(gm.flatItems) == 0 {
		gm.cursor = 0
		return
	}
	if gm.cursor < 0 {
		gm.cursor = 0
	} else if gm.cursor >= len(gm.flatItems) {
		gm.cursor = len(gm.flatItems) - 1
	}
}

// IsAtGroup returns true if the cursor is on a group header.
func (gm *GroupManager) IsAtGroup() bool {
	if gm.cursor >= len(gm.flatItems) {
		return false
	}
	return gm.flatItems[gm.cursor].isGroup
}

// IsAtItem returns true if the cursor is on an item (not a header).
func (gm *GroupManager) IsAtItem() bool {
	if gm.cursor >= len(gm.flatItems) {
		return false
	}
	return !gm.flatItems[gm.cursor].isGroup
}

// CurrentGroup returns the group at the cursor.
func (gm *GroupManager) CurrentGroup() *Group {
	if gm.cursor >= len(gm.flatItems) {
		return nil
	}
	fi := gm.flatItems[gm.cursor]
	if fi.groupIdx >= 0 && fi.groupIdx < len(gm.groups) {
		return gm.groups[fi.groupIdx]
	}
	return nil
}

// CurrentItem returns the item at the cursor.
func (gm *GroupManager) CurrentItem() Item {
	if gm.cursor >= len(gm.flatItems) {
		return nil
	}
	fi := gm.flatItems[gm.cursor]
	if fi.isGroup {
		return nil
	}
	if fi.groupIdx >= 0 && fi.groupIdx < len(gm.groups) {
		group := gm.groups[fi.groupIdx]
		if fi.itemIdx >= 0 && fi.itemIdx < len(group.items) {
			return group.items[fi.itemIdx]
		}
	}
	return nil
}

// CurrentGroupIndex returns the index of the group at the cursor.
func (gm *GroupManager) CurrentGroupIndex() int {
	if gm.cursor >= len(gm.flatItems) {
		return -1
	}
	return gm.flatItems[gm.cursor].groupIdx
}

// ToggleGroupAt toggles the group at the specified index.
func (gm *GroupManager) ToggleGroupAt(groupIdx int) {
	if groupIdx >= 0 && groupIdx < len(gm.groups) {
		gm.groups[groupIdx].Toggle()
		gm.flatten()
		gm.ClampCursor()
	}
}

// ToggleCurrentGroup toggles the group at the cursor.
func (gm *GroupManager) ToggleCurrentGroup() {
	if !gm.IsAtGroup() {
		return
	}
	groupIdx := gm.CurrentGroupIndex()
	gm.ToggleGroupAt(groupIdx)
}

// ExpandAll expands all groups.
func (gm *GroupManager) ExpandAll() {
	for _, g := range gm.groups {
		g.expanded = true
	}
	gm.flatten()
}

// CollapseAll collapses all groups.
func (gm *GroupManager) CollapseAll() {
	for _, g := range gm.groups {
		g.expanded = false
	}
	gm.flatten()
	gm.ClampCursor()
}

// ExpandGroup expands the group at the specified index.
func (gm *GroupManager) ExpandGroup(groupIdx int) {
	if groupIdx >= 0 && groupIdx < len(gm.groups) {
		gm.groups[groupIdx].SetExpanded(true)
		gm.flatten()
	}
}

// CollapseGroup collapses the group at the specified index.
func (gm *GroupManager) CollapseGroup(groupIdx int) {
	if groupIdx >= 0 && groupIdx < len(gm.groups) {
		gm.groups[groupIdx].SetExpanded(false)
		gm.flatten()
		gm.ClampCursor()
	}
}

// GetItemAt returns the item at the specified flat index.
func (gm *GroupManager) GetItemAt(idx int) (isGroup bool, groupIdx int, itemIdx int) {
	if idx < 0 || idx >= len(gm.flatItems) {
		return false, -1, -1
	}
	fi := gm.flatItems[idx]
	return fi.isGroup, fi.groupIdx, fi.itemIdx
}

// AllItems returns all visible (non-header) items in order.
func (gm *GroupManager) AllItems() []Item {
	var items []Item
	for _, fi := range gm.flatItems {
		if !fi.isGroup && fi.groupIdx >= 0 && fi.groupIdx < len(gm.groups) {
			group := gm.groups[fi.groupIdx]
			if fi.itemIdx >= 0 && fi.itemIdx < len(group.items) {
				items = append(items, group.items[fi.itemIdx])
			}
		}
	}
	return items
}

// FindItemID finds the flat index of an item by its ID.
func (gm *GroupManager) FindItemID(id string) int {
	for i, fi := range gm.flatItems {
		if !fi.isGroup && fi.groupIdx >= 0 && fi.groupIdx < len(gm.groups) {
			group := gm.groups[fi.groupIdx]
			if fi.itemIdx >= 0 && fi.itemIdx < len(group.items) {
				if group.items[fi.itemIdx].ID() == id {
					return i
				}
			}
		}
	}
	return -1
}

// FindGroupTitle finds the flat index of a group by its title.
func (gm *GroupManager) FindGroupTitle(title string) int {
	for i, fi := range gm.flatItems {
		if fi.isGroup && fi.groupIdx >= 0 && fi.groupIdx < len(gm.groups) {
			if gm.groups[fi.groupIdx].title == title {
				return i
			}
		}
	}
	return -1
}
