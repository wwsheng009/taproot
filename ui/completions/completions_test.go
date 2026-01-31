package completions

import (
	"testing"
)

func TestSimpleCompletionItem(t *testing.T) {
	item := NewSimpleCompletionItem("id", "display", 123)

	t.Run("ID", func(t *testing.T) {
		if item.ID() != "id" {
			t.Errorf("expected \"id\", got %q", item.ID())
		}
	})

	t.Run("Display", func(t *testing.T) {
		if item.Display() != "display" {
			t.Errorf("expected \"display\", got %q", item.Display())
		}
	})

	t.Run("Value", func(t *testing.T) {
		if item.Value() != 123 {
			t.Errorf("expected 123, got %v", item.Value())
		}
	})
}

func TestStringProvider(t *testing.T) {
	items := []CompletionItem{
		NewSimpleCompletionItem("1", "Apple", "a"),
		NewSimpleCompletionItem("2", "Banana", "b"),
		NewSimpleCompletionItem("3", "Cherry", "c"),
	}

	provider := NewStringProvider(items)

	t.Run("GetItems", func(t *testing.T) {
		result := provider.GetItems()
		if len(result) != 3 {
			t.Errorf("expected 3 items, got %d", len(result))
		}
	})

	t.Run("GetFilterValue", func(t *testing.T) {
		result := provider.GetFilterValue(items[0])
		if result != "Apple" {
			t.Errorf("expected \"Apple\", got %q", result)
		}
	})
}

func TestStringProviderFromStrings(t *testing.T) {
	strings := []string{"alpha", "beta", "gamma"}
	provider := NewStringProviderFromStrings(strings)

	t.Run("GetItems", func(t *testing.T) {
		items := provider.GetItems()
		if len(items) != 3 {
			t.Errorf("expected 3 items, got %d", len(items))
		}
	})
}

func TestAutoCompletion(t *testing.T) {
	items := []CompletionItem{
		NewSimpleCompletionItem("1", "Apple", "a"),
		NewSimpleCompletionItem("2", "Banana", "b"),
		NewSimpleCompletionItem("3", "Cherry", "c"),
	}
	provider := NewStringProvider(items)
	autocomplete := NewAutoCompletion(provider, 1, 10, 50)

	t.Run("Initial state", func(t *testing.T) {
		if autocomplete.IsOpen() {
			t.Error("expected closed")
		}
		if autocomplete.Query() != "" {
			t.Errorf("expected empty query, got %q", autocomplete.Query())
		}
	})

	t.Run("Open", func(t *testing.T) {
		autocomplete.Open()
		if !autocomplete.IsOpen() {
			t.Error("expected open")
		}
		if autocomplete.ItemCount() != 3 {
			t.Errorf("expected 3 items, got %d", autocomplete.ItemCount())
		}
	})

	t.Run("SetQuery", func(t *testing.T) {
		autocomplete.SetQuery("ch")
		if autocomplete.Query() != "ch" {
			t.Errorf("expected \"ch\", got %q", autocomplete.Query())
		}
		// Should filter to only "Cherry"
		if autocomplete.ItemCount() != 1 {
			t.Errorf("expected 1 item, got %d", autocomplete.ItemCount())
		}
	})

	t.Run("SetQuery empty", func(t *testing.T) {
		autocomplete.SetQuery("")
		// Should show all items
		if autocomplete.ItemCount() != 3 {
			t.Errorf("expected 3 items, got %d", autocomplete.ItemCount())
		}
	})

	t.Run("Cursor navigation", func(t *testing.T) {
		autocomplete.SetQuery("")
		if autocomplete.Cursor() != 0 {
			t.Errorf("expected cursor at 0, got %d", autocomplete.Cursor())
		}

		autocomplete.MoveDown()
		if autocomplete.Cursor() != 1 {
			t.Errorf("expected cursor at 1, got %d", autocomplete.Cursor())
		}

		autocomplete.MoveDown()
		autocomplete.MoveDown() // Wrap around
		if autocomplete.Cursor() != 0 {
			t.Errorf("expected cursor at 0 after wrap, got %d", autocomplete.Cursor())
		}

		autocomplete.MoveUp() // Should wrap to end
		if autocomplete.Cursor() != 2 {
			t.Errorf("expected cursor at 2 after wrap up, got %d", autocomplete.Cursor())
		}
	})

	t.Run("Selected", func(t *testing.T) {
		autocomplete.SetCursor(1)
		selected := autocomplete.Selected()
		if selected == nil {
			t.Fatal("expected non-nil selected item")
		}
		if selected.ID() != "2" {
			t.Errorf("expected ID \"2\", got %q", selected.ID())
		}
	})

	t.Run("SelectedValue", func(t *testing.T) {
		autocomplete.SetCursor(2)
		value := autocomplete.SelectedValue()
		if value != "c" {
			t.Errorf("expected \"c\", got %v", value)
		}
	})

	t.Run("Close", func(t *testing.T) {
		autocomplete.Close()
		if autocomplete.IsOpen() {
			t.Error("expected closed")
		}
		if autocomplete.ItemCount() != 0 {
			t.Errorf("expected 0 items after close, got %d", autocomplete.ItemCount())
		}
	})
}

func TestAutoCompletionFiltering(t *testing.T) {
	items := []CompletionItem{
		NewSimpleCompletionItem("1", "Apple", "a"),
		NewSimpleCompletionItem("2", "Banana", "b"),
		NewSimpleCompletionItem("3", "Apricot", "c"),
		NewSimpleCompletionItem("4", "Blueberry", "d"),
	}
	provider := NewStringProvider(items)
	autocomplete := NewAutoCompletion(provider, 1, 10, 50)
	autocomplete.Open()

	t.Run("Filter by 'a'", func(t *testing.T) {
		autocomplete.SetQuery("a")
		if autocomplete.ItemCount() != 3 {
			t.Errorf("expected 3 items matching 'a', got %d", autocomplete.ItemCount())
		}
	})

	t.Run("Filter by 'an'", func(t *testing.T) {
		autocomplete.SetQuery("an")
		if autocomplete.ItemCount() != 1 {
			t.Errorf("expected 1 item matching 'an', got %d", autocomplete.ItemCount())
		}
		if autocomplete.Selected().ID() != "2" {
			t.Errorf("expected \"Banana\", got %q", autocomplete.Selected().Display())
		}
	})

	t.Run("No matches", func(t *testing.T) {
		autocomplete.SetQuery("xyz")
		if autocomplete.ItemCount() != 0 {
			t.Errorf("expected 0 items, got %d", autocomplete.ItemCount())
		}
		if autocomplete.Selected() != nil {
			t.Error("expected nil selected when no items")
		}
	})
}

func TestAutoCompletionDimensions(t *testing.T) {
	items := []CompletionItem{
		NewSimpleCompletionItem("1", "Short", "a"),
		NewSimpleCompletionItem("2", "Medium Length", "b"),
		NewSimpleCompletionItem("3", "Very Long Item Name Here", "c"),
	}
	provider := NewStringProvider(items)
	autocomplete := NewAutoCompletion(provider, 1, 10, 50)
	autocomplete.Open()

	t.Run("Width calculation", func(t *testing.T) {
		w, h := autocomplete.Size()
		if w == 0 {
			t.Error("expected non-zero width")
		}
		if h == 0 {
			t.Error("expected non-zero height")
		}
		// Width should accommodate the longest item
		expectedMinWidth := len("Very Long Item Name Here") + 4
		if w < expectedMinWidth {
			t.Errorf("expected width >= %d, got %d", expectedMinWidth, w)
		}
	})

	t.Run("Height limits", func(t *testing.T) {
		// Add more items than maxHeight
		for i := 4; i <= 20; i++ {
			items = append(items, NewSimpleCompletionItem(string(rune('0'+i)), "Item "+string(rune('0'+i)), i))
		}
		provider2 := NewStringProvider(items)
		autocomplete2 := NewAutoCompletion(provider2, 1, 10, 50)
		autocomplete2.Open()

		_, h := autocomplete2.Size()
		if h > 10 {
			t.Errorf("expected height <= 10, got %d", h)
		}
	})
}

func TestVisibleRange(t *testing.T) {
	items := []CompletionItem{}
	for i := 0; i < 20; i++ {
		items = append(items, NewSimpleCompletionItem(string(rune('0'+i)), "Item "+string(rune('0'+i)), i))
	}
	provider := NewStringProvider(items)
	autocomplete := NewAutoCompletion(provider, 5, 10, 50)
	autocomplete.Open()

	t.Run("Initial range", func(t *testing.T) {
		start, end := autocomplete.VisibleRange()
		if start != 0 || end != 10 {
			t.Errorf("expected range [0, 10), got [%d, %d)", start, end)
		}
	})

	t.Run("Scroll down", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			autocomplete.MoveDown()
		}
		autocomplete.ScrollToVisible()
		start, end := autocomplete.VisibleRange()
		// Cursor should still be visible
		if start > autocomplete.Cursor() || end <= autocomplete.Cursor() {
			t.Errorf("cursor %d not in visible range [%d, %d)", autocomplete.Cursor(), start, end)
		}
	})
}
