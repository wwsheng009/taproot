package forms

import (
	"testing"

	"github.com/wwsheng009/taproot/ui/render"
)

func TestNewForm(t *testing.T) {
	in1 := NewTextInput("1")
	in2 := NewTextInput("2")
	
	f := NewForm(in1, in2)
	
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0, got %d", f.FocusedIndex())
	}
	
	if !in1.Focused() {
		t.Error("expected first input to be focused")
	}
	if in2.Focused() {
		t.Error("expected second input to be blurred")
	}
}

func TestForm_Navigation(t *testing.T) {
	in1 := NewTextInput("1")
	in2 := NewTextInput("2")
	in3 := NewTextInput("3")
	
	f := NewForm(in1, in2, in3)
	
	// Test Next (Tab)
	f.Update(render.KeyMsg{Key: "tab"})
	if f.FocusedIndex() != 1 {
		t.Errorf("expected focused index 1 after tab, got %d", f.FocusedIndex())
	}
	if !in2.Focused() {
		t.Error("expected second input to be focused")
	}
	
	// Test Next loop around
	f.Update(render.KeyMsg{Key: "tab"}) // -> 2
	f.Update(render.KeyMsg{Key: "tab"}) // -> 0
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0 after looping, got %d", f.FocusedIndex())
	}
	
	// Test Prev (Shift+Tab)
	f.Update(render.KeyMsg{Key: "shift+tab"})
	if f.FocusedIndex() != 2 {
		t.Errorf("expected focused index 2 after shift+tab, got %d", f.FocusedIndex())
	}
	
	// Test Enter (for TextInput, should act as Next)
	f.Update(render.KeyMsg{Key: "enter"}) // -> 0
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focused index 0 after enter, got %d", f.FocusedIndex())
	}
}

func TestForm_EnterBehavior(t *testing.T) {
	// TextArea should capture Enter
	ta := NewTextArea("bio")
	ti := NewTextInput("name")
	
	f := NewForm(ta, ti)
	
	// Focus TextArea
	if f.FocusedIndex() != 0 {
		t.Error("expected TextArea to be focused initially")
	}
	
	// Send Enter - TextArea should keep focus and receive the newline
	f.Update(render.KeyMsg{Key: "enter"})
	
	if f.FocusedIndex() != 0 {
		t.Error("expected TextArea to retain focus on Enter")
	}
	if ta.Value() != "\n" {
		t.Errorf("expected TextArea to contain newline, got %q", ta.Value())
	}
	
	// Move focus to TextInput
	f.Update(render.KeyMsg{Key: "tab"})
	if f.FocusedIndex() != 1 {
		t.Fatal("failed to switch focus")
	}
	
	// Send Enter - TextInput should yield focus (loop back to start)
	f.Update(render.KeyMsg{Key: "enter"})
	if f.FocusedIndex() != 0 {
		t.Errorf("expected focus to move from TextInput on Enter, got index %d", f.FocusedIndex())
	}
}

func TestForm_Validate(t *testing.T) {
	in1 := NewTextInput("req")
	in1.AddValidator(Required)
	
	in2 := NewTextInput("opt")
	
	f := NewForm(in1, in2)
	
	// Empty required field -> Error
	if err := f.Validate(); err == nil {
		t.Error("expected validation error for empty required field")
	}
	
	// Filled required field -> Success
	in1.SetValue("ok")
	if err := f.Validate(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}
