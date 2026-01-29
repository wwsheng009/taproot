package list

// SelectionMode determines how selection works in a list.
type SelectionMode int

const (
	// SelectionNone disables selection.
	SelectionModeNone SelectionMode = iota
	// SelectionSingle allows only one item to be selected at a time.
	SelectionModeSingle
	// SelectionMultiple allows multiple items to be selected.
	SelectionModeMultiple
)

// SelectionManager manages item selection state.
type SelectionManager struct {
	mode      SelectionMode
	selected  map[string]struct{}
	cursorIdx int
}

// NewSelectionManager creates a new selection manager.
func NewSelectionManager(mode SelectionMode) *SelectionManager {
	return &SelectionManager{
		mode:     mode,
		selected: make(map[string]struct{}),
	}
}

// Mode returns the current selection mode.
func (sm *SelectionManager) Mode() SelectionMode {
	return sm.mode
}

// SetMode changes the selection mode, clearing existing selections if needed.
func (sm *SelectionManager) SetMode(mode SelectionMode) {
	if mode != sm.mode {
		sm.mode = mode
		sm.Clear()
	}
}

// SetCursor updates the cursor index.
func (sm *SelectionManager) SetCursor(idx int) {
	sm.cursorIdx = idx
}

// Cursor returns the current cursor index.
func (sm *SelectionManager) Cursor() int {
	return sm.cursorIdx
}

// Select adds an item to the selection.
func (sm *SelectionManager) Select(id string) {
	if sm.mode == SelectionModeNone {
		return
	}
	if sm.mode == SelectionModeSingle {
		sm.selected = make(map[string]struct{})
	}
	sm.selected[id] = struct{}{}
}

// Deselect removes an item from the selection.
func (sm *SelectionManager) Deselect(id string) {
	delete(sm.selected, id)
}

// Toggle toggles the selection state of an item.
func (sm *SelectionManager) Toggle(id string) {
	if sm.IsSelected(id) {
		sm.Deselect(id)
	} else {
		sm.Select(id)
	}
}

// IsSelected returns true if the item is selected.
func (sm *SelectionManager) IsSelected(id string) bool {
	_, ok := sm.selected[id]
	return ok
}

// SelectOnly selects only the specified item, deselecting all others.
func (sm *SelectionManager) SelectOnly(id string) {
	sm.selected = make(map[string]struct{})
	sm.selected[id] = struct{}{}
}

// Clear clears all selections.
func (sm *SelectionManager) Clear() {
	sm.selected = make(map[string]struct{})
}

// Count returns the number of selected items.
func (sm *SelectionManager) Count() int {
	return len(sm.selected)
}

// SelectedIDs returns a list of selected item IDs.
func (sm *SelectionManager) SelectedIDs() []string {
	if len(sm.selected) == 0 {
		return nil
	}
	ids := make([]string, 0, len(sm.selected))
	for id := range sm.selected {
		ids = append(ids, id)
	}
	return ids
}

// SetSelectedIDs sets the selected items from a list of IDs.
func (sm *SelectionManager) SetSelectedIDs(ids []string) {
	sm.selected = make(map[string]struct{}, len(ids))
	for _, id := range ids {
		sm.selected[id] = struct{}{}
	}
}

// AllSelected returns true if all items would be selected.
func (sm *SelectionManager) AllSelected(total int) bool {
	return total > 0 && len(sm.selected) == total
}

// HasSelection returns true if any items are selected.
func (sm *SelectionManager) HasSelection() bool {
	return len(sm.selected) > 0
}

// SelectRange selects all items between two indices (inclusive).
// The items slice is used to look up IDs by index.
func (sm *SelectionManager) SelectRange(items []Item, start, end int) {
	if sm.mode == SelectionModeNone {
		return
	}
	if sm.mode == SelectionModeSingle {
		sm.Clear()
		if start >= 0 && start < len(items) {
			sm.Select(items[start].ID())
		}
		return
	}

	// Ensure start <= end
	if start > end {
		start, end = end, start
	}

	// Clamp to valid range
	start = max(0, start)
	end = min(len(items)-1, end)

	for i := start; i <= end; i++ {
		sm.Select(items[i].ID())
	}
}

// SelectVisible selects all visible items in the viewport.
func (sm *SelectionManager) SelectVisible(items []Item, start, end int) {
	sm.SelectRange(items, start, end)
}

// InvertSelection inverts the selection state of all items.
func (sm *SelectionManager) InvertSelection(items []Item) {
	newSelected := make(map[string]struct{})
	for _, item := range items {
		if !sm.IsSelected(item.ID()) {
			newSelected[item.ID()] = struct{}{}
		}
	}
	sm.selected = newSelected
}

// SelectAll selects all items.
func (sm *SelectionManager) SelectAll(items []Item) {
	if sm.mode == SelectionModeNone {
		return
	}
	if sm.mode == SelectionModeSingle && len(items) > 0 {
		sm.SelectOnly(items[0].ID())
		return
	}
	for _, item := range items {
		sm.Select(item.ID())
	}
}

// GetSelected retrieves the selected items from a slice.
func (sm *SelectionManager) GetSelected(items []Item) []Item {
	var result []Item
	for _, item := range items {
		if sm.IsSelected(item.ID()) {
			result = append(result, item)
		}
	}
	return result
}
