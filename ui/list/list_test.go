package list

import (
	"testing"
)

func TestViewport(t *testing.T) {
	t.Run("NewViewport", func(t *testing.T) {
		vp := NewViewport(10, 100)
		if vp.Visible() != 10 {
			t.Errorf("expected visible 10, got %d", vp.Visible())
		}
		if vp.Total() != 100 {
			t.Errorf("expected total 100, got %d", vp.Total())
		}
		if vp.Cursor() != 0 {
			t.Errorf("expected cursor 0, got %d", vp.Cursor())
		}
		if vp.Offset() != 0 {
			t.Errorf("expected offset 0, got %d", vp.Offset())
		}
	})

	t.Run("MoveDown", func(t *testing.T) {
		vp := NewViewport(5, 20)
		for i := 0; i < 7; i++ {
			vp.MoveDown()
		}
		if vp.Cursor() != 7 {
			t.Errorf("expected cursor 7, got %d", vp.Cursor())
		}
		// Offset should have adjusted
		if vp.Offset() != 3 {
			t.Errorf("expected offset 3, got %d", vp.Offset())
		}
	})

	t.Run("MoveUp", func(t *testing.T) {
		vp := NewViewport(5, 20)
		vp.SetCursor(10)
		vp.MoveUp()
		if vp.Cursor() != 9 {
			t.Errorf("expected cursor 9, got %d", vp.Cursor())
		}
	})

	t.Run("MoveToTop", func(t *testing.T) {
		vp := NewViewport(5, 20)
		vp.SetCursor(10)
		vp.MoveToTop()
		if vp.Cursor() != 0 {
			t.Errorf("expected cursor 0, got %d", vp.Cursor())
		}
		if vp.Offset() != 0 {
			t.Errorf("expected offset 0, got %d", vp.Offset())
		}
	})

	t.Run("MoveToBottom", func(t *testing.T) {
		vp := NewViewport(5, 20)
		vp.MoveToBottom()
		if vp.Cursor() != 19 {
			t.Errorf("expected cursor 19, got %d", vp.Cursor())
		}
		if vp.Offset() != 15 {
			t.Errorf("expected offset 15, got %d", vp.Offset())
		}
	})

	t.Run("PageDown", func(t *testing.T) {
		vp := NewViewport(5, 20)
		vp.PageDown()
		if vp.Cursor() != 5 {
			t.Errorf("expected cursor 5, got %d", vp.Cursor())
		}
	})

	t.Run("PageUp", func(t *testing.T) {
		vp := NewViewport(5, 20)
		vp.SetCursor(10)
		vp.PageUp()
		if vp.Cursor() != 5 {
			t.Errorf("expected cursor 5, got %d", vp.Cursor())
		}
	})

	t.Run("Range", func(t *testing.T) {
		vp := NewViewport(5, 20)
		start, end := vp.Range()
		if start != 0 {
			t.Errorf("expected start 0, got %d", start)
		}
		if end != 5 {
			t.Errorf("expected end 5, got %d", end)
		}
	})

	t.Run("SetTotal", func(t *testing.T) {
		vp := NewViewport(5, 20)
		vp.SetCursor(10)
		vp.SetTotal(5)
		// Cursor should be clamped
		if vp.Cursor() != 4 {
			t.Errorf("expected cursor 4, got %d", vp.Cursor())
		}
	})
}

func TestSelectionManager(t *testing.T) {
	items := []Item{
		NewListItem("1", "Item 1", ""),
		NewListItem("2", "Item 2", ""),
		NewListItem("3", "Item 3", ""),
	}

	t.Run("NewSelectionManager", func(t *testing.T) {
		sm := NewSelectionManager(SelectionModeMultiple)
		if sm.Mode() != SelectionModeMultiple {
			t.Errorf("expected multiple mode, got %v", sm.Mode())
		}
		if sm.Count() != 0 {
			t.Errorf("expected 0 selections, got %d", sm.Count())
		}
	})

	t.Run("Select", func(t *testing.T) {
		sm := NewSelectionManager(SelectionModeMultiple)
		sm.Select("1")
		if !sm.IsSelected("1") {
			t.Error("expected item 1 to be selected")
		}
		if sm.Count() != 1 {
			t.Errorf("expected count 1, got %d", sm.Count())
		}
	})

	t.Run("Toggle", func(t *testing.T) {
		sm := NewSelectionManager(SelectionModeMultiple)
		sm.Toggle("1")
		if !sm.IsSelected("1") {
			t.Error("expected item 1 to be selected")
		}
		sm.Toggle("1")
		if sm.IsSelected("1") {
			t.Error("expected item 1 to be deselected")
		}
	})

	t.Run("SingleMode", func(t *testing.T) {
		sm := NewSelectionManager(SelectionModeSingle)
		sm.Select("1")
		sm.Select("2")
		if sm.Count() != 1 {
			t.Errorf("expected count 1, got %d", sm.Count())
		}
		if !sm.IsSelected("2") {
			t.Error("expected item 2 to be selected")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		sm := NewSelectionManager(SelectionModeMultiple)
		sm.Select("1")
		sm.Select("2")
		sm.Clear()
		if sm.Count() != 0 {
			t.Errorf("expected count 0, got %d", sm.Count())
		}
	})

	t.Run("SelectAll", func(t *testing.T) {
		sm := NewSelectionManager(SelectionModeMultiple)
		sm.SelectAll(items)
		if sm.Count() != 3 {
			t.Errorf("expected count 3, got %d", sm.Count())
		}
	})
}

func TestFilter(t *testing.T) {
	items := []FilterableItem{
		NewListItem("1", "Apple", "Red fruit"),
		NewListItem("2", "Banana", "Yellow fruit"),
		NewListItem("3", "Cherry", "Red fruit"),
	}

	t.Run("NoFilter", func(t *testing.T) {
		f := NewFilter()
		filtered := f.Apply(items)
		if len(filtered) != 3 {
			t.Errorf("expected 3 items, got %d", len(filtered))
		}
	})

	t.Run("WithQuery", func(t *testing.T) {
		f := NewFilter()
		f.SetQuery("Red")
		filtered := f.Apply(items)
		if len(filtered) != 2 {
			t.Errorf("expected 2 items, got %d", len(filtered))
		}
	})

	t.Run("CaseInsensitive", func(t *testing.T) {
		f := NewFilter()
		f.SetQuery("apple")
		f.SetCaseSensitive(false)
		filtered := f.Apply(items)
		if len(filtered) != 1 {
			t.Errorf("expected 1 item, got %d", len(filtered))
		}
	})

	t.Run("Clear", func(t *testing.T) {
		f := NewFilter()
		f.SetQuery("Red")
		f.Clear()
		if f.Query() != "" {
			t.Errorf("expected empty query, got %q", f.Query())
		}
		if f.Active() {
			t.Error("expected filter to be inactive")
		}
	})
}

func TestGroup(t *testing.T) {
	t.Run("NewGroup", func(t *testing.T) {
		items := []Item{
			NewListItem("1", "Item 1", ""),
			NewListItem("2", "Item 2", ""),
		}
		g := NewGroup("Fruits", items)
		if g.Title() != "Fruits" {
			t.Errorf("expected title 'Fruits', got %q", g.Title())
		}
		if !g.Expanded() {
			t.Error("expected group to be expanded")
		}
		if g.ItemCount() != 2 {
			t.Errorf("expected 2 items, got %d", g.ItemCount())
		}
	})

	t.Run("Toggle", func(t *testing.T) {
		items := []Item{
			NewListItem("1", "Item 1", ""),
		}
		g := NewGroup("Test", items)
		g.Toggle()
		if g.Expanded() {
			t.Error("expected group to be collapsed")
		}
		g.Toggle()
		if !g.Expanded() {
			t.Error("expected group to be expanded")
		}
	})
}

func TestGroupManager(t *testing.T) {
	t.Run("NewGroupManager", func(t *testing.T) {
		gm := NewGroupManager()
		if gm.GroupCount() != 0 {
			t.Errorf("expected 0 groups, got %d", gm.GroupCount())
		}
	})

	t.Run("SetGroups", func(t *testing.T) {
		gm := NewGroupManager()
		groups := []*Group{
			NewGroup("A", []Item{NewListItem("1", "Item 1", "")}),
			NewGroup("B", []Item{NewListItem("2", "Item 2", "")}),
		}
		gm.SetGroups(groups)
		if gm.GroupCount() != 2 {
			t.Errorf("expected 2 groups, got %d", gm.GroupCount())
		}
		// Count includes headers
		if gm.Count() != 4 { // 2 headers + 2 items
			t.Errorf("expected count 4, got %d", gm.Count())
		}
	})

	t.Run("CollapseGroup", func(t *testing.T) {
		gm := NewGroupManager()
		groups := []*Group{
			NewGroup("A", []Item{NewListItem("1", "Item 1", "")}),
		}
		gm.SetGroups(groups)
		gm.CollapseGroup(0)
		if gm.Count() != 1 { // Only header
			t.Errorf("expected count 1, got %d", gm.Count())
		}
	})

	t.Run("ToggleCurrentGroup", func(t *testing.T) {
		gm := NewGroupManager()
		groups := []*Group{
			NewGroup("A", []Item{NewListItem("1", "Item 1", "")}),
		}
		gm.SetGroups(groups)
		gm.ToggleCurrentGroup()
		if gm.Count() != 1 { // Only header
			t.Errorf("expected count 1, got %d", gm.Count())
		}
	})

	t.Run("ExpandAll/CollapseAll", func(t *testing.T) {
		gm := NewGroupManager()
		groups := []*Group{
			NewGroup("A", []Item{NewListItem("1", "Item 1", "")}),
			NewGroup("B", []Item{NewListItem("2", "Item 2", "")}),
		}
		gm.SetGroups(groups)
		gm.CollapseAll()
		if gm.Count() != 2 { // Only headers
			t.Errorf("expected count 2, got %d", gm.Count())
		}
		gm.ExpandAll()
		if gm.Count() != 4 { // Headers + items
			t.Errorf("expected count 4, got %d", gm.Count())
		}
	})
}

func TestAction(t *testing.T) {
	km := DefaultKeyMap()

	t.Run("NavigationKeys", func(t *testing.T) {
		tests := []struct {
			key     string
			action  Action
		}{
			{"up", ActionMoveUp},
			{"k", ActionMoveUp},
			{"down", ActionMoveDown},
			{"j", ActionMoveDown},
			{"g", ActionMoveToTop},
			{"G", ActionMoveToBottom},
		}

		for _, tt := range tests {
			if a := km.MatchAction(tt.key); a != tt.action {
				t.Errorf("key %q: expected %v, got %v", tt.key, tt.action, a)
			}
		}
	})

	t.Run("SelectionKeys", func(t *testing.T) {
		if a := km.MatchAction(" "); a != ActionToggleSelection {
			t.Errorf("space: expected ToggleSelection, got %v", a)
		}
		if a := km.MatchAction("ctrl+a"); a != ActionSelectAll {
			t.Errorf("ctrl+a: expected SelectAll, got %v", a)
		}
	})

	t.Run("FilterKey", func(t *testing.T) {
		if a := km.MatchAction("/"); a != ActionFilter {
			t.Errorf("/: expected Filter, got %v", a)
		}
	})
}
