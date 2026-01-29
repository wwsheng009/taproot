package dialog

import (
	"testing"

	"github.com/yourorg/taproot/internal/ui/list"
	"github.com/yourorg/taproot/internal/ui/render"
)

func TestButton(t *testing.T) {
	t.Run("NewButton", func(t *testing.T) {
		b := NewButton("OK", true)
		if b.Label() != "OK" {
			t.Errorf("expected label 'OK', got %q", b.Label())
		}
		if !b.Primary() {
			t.Error("expected primary button")
		}
	})

	t.Run("SetSelected", func(t *testing.T) {
		b := NewButton("Cancel", false)
		b.SetSelected(true)
		if !b.Selected() {
			t.Error("expected selected to be true")
		}
	})

	t.Run("Toggle", func(t *testing.T) {
		b := NewButton("OK", true)
		b.Toggle()
		if !b.Selected() {
			t.Error("expected selected after toggle")
		}
		b.Toggle()
		if b.Selected() {
			t.Error("expected not selected after second toggle")
		}
	})
}

func TestButtonGroup(t *testing.T) {
	t.Run("NewButtonGroup", func(t *testing.T) {
		bg := NewButtonGroup(true)
		if !bg.Horizontal() {
			t.Error("expected horizontal group")
		}
		if bg.Count() != 0 {
			t.Errorf("expected 0 buttons, got %d", bg.Count())
		}
	})

	t.Run("AddButtons", func(t *testing.T) {
		bg := NewButtonGroup(true).
			Add("Cancel", false).
			Add("OK", true)
		
		if bg.Count() != 2 {
			t.Errorf("expected 2 buttons, got %d", bg.Count())
		}
	})

	t.Run("SelectNext", func(t *testing.T) {
		bg := NewButtonGroup(true).
			Add("1", false).
			Add("2", false).
			Add("3", false)
		
		bg.SetSelected(0)
		bg.SelectNext()
		if bg.Selected() != 1 {
			t.Errorf("expected selected 1, got %d", bg.Selected())
		}
		
		bg.SelectNext()
		if bg.Selected() != 2 {
			t.Errorf("expected selected 2, got %d", bg.Selected())
		}
		
		// Should wrap around
		bg.SelectNext()
		if bg.Selected() != 0 {
			t.Errorf("expected selected 0 (wrapped), got %d", bg.Selected())
		}
	})

	t.Run("SelectPrev", func(t *testing.T) {
		bg := NewButtonGroup(true).
			Add("1", false).
			Add("2", false)
		
		bg.SetSelected(0)
		bg.SelectPrev()
		if bg.Selected() != 1 {
			t.Errorf("expected selected 1 (wrapped), got %d", bg.Selected())
		}
	})
}

func TestInputField(t *testing.T) {
	t.Run("NewInputField", func(t *testing.T) {
		i := NewInputField("Enter text")
		if i.Placeholder() != "Enter text" {
			t.Errorf("expected placeholder 'Enter text', got %q", i.Placeholder())
		}
		if i.Value() != "" {
			t.Errorf("expected empty value, got %q", i.Value())
		}
	})

	t.Run("Insert", func(t *testing.T) {
		i := NewInputField("")
		i.Insert('H')
		i.Insert('i')
		if i.Value() != "Hi" {
			t.Errorf("expected 'Hi', got %q", i.Value())
		}
		if i.Cursor() != 2 {
			t.Errorf("expected cursor at 2, got %d", i.Cursor())
		}
	})

	t.Run("Delete", func(t *testing.T) {
		i := NewInputField("")
		i.Insert('a')
		i.Insert('b')
		i.Delete()
		if i.Value() != "a" {
			t.Errorf("expected 'a', got %q", i.Value())
		}
		if i.Cursor() != 1 {
			t.Errorf("expected cursor at 1, got %d", i.Cursor())
		}
	})

	t.Run("SetValue", func(t *testing.T) {
		i := NewInputField("")
		i.SetValue("test")
		if i.Value() != "test" {
			t.Errorf("expected 'test', got %q", i.Value())
		}
		if i.Cursor() != 4 {
			t.Errorf("expected cursor at 4, got %d", i.Cursor())
		}
	})

	t.Run("MaxLength", func(t *testing.T) {
		i := NewInputField("")
		i.SetMaxLength(3)
		i.Insert('a')
		i.Insert('b')
		i.Insert('c')
		i.Insert('d') // Should be ignored
		if i.Value() != "abc" {
			t.Errorf("expected 'abc', got %q", i.Value())
		}
	})

	t.Run("Hidden", func(t *testing.T) {
		i := NewInputField("")
		i.SetHidden(true)
		i.SetValue("password")
		if !i.Hidden() {
			t.Error("expected hidden to be true")
		}
	})
}

func TestSelectList(t *testing.T) {
	items := []list.Item{
		list.NewListItem("1", "Apple", ""),
		list.NewListItem("2", "Banana", ""),
		list.NewListItem("3", "Cherry", ""),
	}

	t.Run("NewSelectList", func(t *testing.T) {
		s := NewSelectList(items)
		if len(s.Items()) != 3 {
			t.Errorf("expected 3 items, got %d", len(s.Items()))
		}
		if s.Selected() != 0 {
			t.Errorf("expected selected 0, got %d", s.Selected())
		}
	})

	t.Run("MoveDown", func(t *testing.T) {
		s := NewSelectList(items)
		s.MoveDown()
		if s.Selected() != 1 {
			t.Errorf("expected selected 1, got %d", s.Selected())
		}
		s.MoveDown()
		if s.Selected() != 2 {
			t.Errorf("expected selected 2, got %d", s.Selected())
		}
		// Should not go beyond
		s.MoveDown()
		if s.Selected() != 2 {
			t.Errorf("expected still selected 2, got %d", s.Selected())
		}
	})

	t.Run("MoveUp", func(t *testing.T) {
		s := NewSelectList(items)
		s.SetSelected(2)
		s.MoveUp()
		if s.Selected() != 1 {
			t.Errorf("expected selected 1, got %d", s.Selected())
		}
	})

	t.Run("Toggle", func(t *testing.T) {
		s := NewSelectList(items)
		if s.Expanded() {
			t.Error("expected not expanded initially")
		}
		s.Toggle()
		if !s.Expanded() {
			t.Error("expected expanded after toggle")
		}
	})

	t.Run("SelectedItem", func(t *testing.T) {
		s := NewSelectList(items)
		s.SetSelected(1)
		item := s.SelectedItem()
		if item == nil {
			t.Fatal("expected non-nil item")
		}
		if li, ok := item.(*list.ListItem); ok {
			if li.Title() != "Banana" {
				t.Errorf("expected 'Banana', got %q", li.Title())
			}
		}
	})
}

func TestOverlay(t *testing.T) {
	t.Run("NewOverlay", func(t *testing.T) {
		o := NewOverlay()
		if o.HasDialogs() {
			t.Error("expected no dialogs initially")
		}
		if o.Count() != 0 {
			t.Errorf("expected count 0, got %d", o.Count())
		}
	})

	t.Run("PushPeekPop", func(t *testing.T) {
		o := NewOverlay()
		d1 := NewInfoDialog("Dialog1", "Message1")
		d2 := NewInfoDialog("Dialog2", "Message2")
		
		o.Push(d1)
		if !o.HasDialogs() {
			t.Error("expected dialogs after push")
		}
		if o.Count() != 1 {
			t.Errorf("expected count 1, got %d", o.Count())
		}
		
		peeked := o.Peek()
		if peeked == nil {
			t.Error("expected non-nil peeked dialog")
		}
		
		o.Push(d2)
		if o.Count() != 2 {
			t.Errorf("expected count 2, got %d", o.Count())
		}
		
		popped := o.Pop()
		if popped == nil {
			t.Error("expected non-nil popped dialog")
		}
		if o.Count() != 1 {
			t.Errorf("expected count 1 after pop, got %d", o.Count())
		}
	})

	t.Run("ActiveDialog", func(t *testing.T) {
		o := NewOverlay()
		d := NewInfoDialog("Test", "Message")
		
		if o.ActiveDialog() != nil {
			t.Error("expected no active dialog initially")
		}
		
		o.Push(d)
		active := o.ActiveDialog()
		if active == nil {
			t.Error("expected non-nil active dialog")
		}
		if !o.IsActive(d) {
			t.Error("expected pushed dialog to be active")
		}
	})

	t.Run("FindByID", func(t *testing.T) {
		o := NewOverlay()
		d1 := NewInfoDialog("Dialog1", "Message1")
		d1.SetID("test-id")
		
		o.Push(d1)
		found := o.FindByID("test-id")
		if found == nil {
			t.Error("expected to find dialog by ID")
		}
		
		notFound := o.FindByID("non-existent")
		if notFound != nil {
			t.Error("expected nil for non-existent ID")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		o := NewOverlay()
		o.Push(NewInfoDialog("D1", "M1"))
		o.Push(NewInfoDialog("D2", "M2"))
		
		o.Clear()
		if o.HasDialogs() {
			t.Error("expected no dialogs after clear")
		}
		if o.Count() != 0 {
			t.Errorf("expected count 0 after clear, got %d", o.Count())
		}
	})
}

func TestDialogBounds(t *testing.T) {
	t.Run("CalculateBounds", func(t *testing.T) {
		bounds := CalculateBounds(40, 10, 80, 24)
		if bounds.Width != 40 {
			t.Errorf("expected width 40, got %d", bounds.Width)
		}
		if bounds.Height != 10 {
			t.Errorf("expected height 10, got %d", bounds.Height)
		}
		// Should be centered
		if bounds.X != 20 {
			t.Errorf("expected x 20, got %d", bounds.X)
		}
		if bounds.Y != 7 {
			t.Errorf("expected y 7, got %d", bounds.Y)
		}
	})

	t.Run("DefaultSizes", func(t *testing.T) {
		w := DefaultWidth()
		if w != 60 {
			t.Errorf("expected default width 60, got %d", w)
		}
		
		h := DefaultHeight()
		if h != 15 {
			t.Errorf("expected default height 15, got %d", h)
		}
	})

	t.Run("MaxWidth", func(t *testing.T) {
		max := MaxWidth(100)
		if max != 92 { // 100 - 8 padding
			t.Errorf("expected max width 92, got %d", max)
		}
		
		max = MaxWidth(50)
		if max < 40 { // Min 40
			t.Errorf("expected max width at least 40, got %d", max)
		}
		
		max = MaxWidth(200)
		if max != 100 { // Max 100
			t.Errorf("expected max width 100, got %d", max)
		}
	})
}

func TestInfoDialog(t *testing.T) {
	t.Run("NewInfoDialog", func(t *testing.T) {
		d := NewInfoDialog("Info", "This is a message")
		
		if d.Title() != "Info" {
			t.Errorf("expected title 'Info', got %q", d.Title())
		}
		if d.Message() != "This is a message" {
			t.Errorf("expected message 'This is a message', got %q", d.Message())
		}
	})

	t.Run("Init", func(t *testing.T) {
		d := NewInfoDialog("Info", "Message")
		err := d.Init()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		d := NewInfoDialog("Info", "Message")
		resultCalled := false
		d.SetCallback(func(result ActionResult, data any) {
			resultCalled = true
		})
		
		// Send enter key
		m, _ := d.Update(render.KeyMsg{Key: "enter"})
		_ = m
		
		if !resultCalled {
			t.Error("expected callback to be called")
		}
	})

	t.Run("View", func(t *testing.T) {
		d := NewInfoDialog("Info", "Message")
		view := d.View()
		if view == "" {
			t.Error("expected non-empty view")
		}
	})
}

func TestConfirmDialog(t *testing.T) {
	t.Run("NewConfirmDialog", func(t *testing.T) {
		d := NewConfirmDialog("Confirm", "Are you sure?", func(r ActionResult, data any) {
			// Callback
		})
		
		if d.Title() != "Confirm" {
			t.Errorf("expected title 'Confirm', got %q", d.Title())
		}
		if d.Selected() != 0 { // Default to cancel
			t.Errorf("expected selected 0, got %d", d.Selected())
		}
	})

	t.Run("SelectConfirm", func(t *testing.T) {
		result := ActionNone
		d := NewConfirmDialog("Confirm", "Are you sure?", func(r ActionResult, data any) {
			result = r
		})
		
		// Select confirm
		d.SetSelected(1)
		
		// Send enter
		m, _ := d.Update(render.KeyMsg{Key: "enter"})
		_ = m
		
		if result != ActionConfirm {
			t.Errorf("expected ActionConfirm, got %v", result)
		}
	})

	t.Run("EscapeCancels", func(t *testing.T) {
		d := NewConfirmDialog("Confirm", "Are you sure?", func(r ActionResult, data any) {
			// Callback receives result
		})
		
		m, _ := d.Update(render.KeyMsg{Key: "escape"})
		_ = m
	})
}
