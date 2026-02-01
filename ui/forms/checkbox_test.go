package forms

import (
	"testing"

	"github.com/wwsheng009/taproot/ui/render"
)

func TestCheckbox_Init(t *testing.T) {
	cb := NewCheckbox("Accept Terms")
	if cb.Init() != nil {
		t.Error("Init should return nil")
	}
}

func TestCheckbox_Value(t *testing.T) {
	cb := NewCheckbox("Test")
	
	if cb.Value() != "false" {
		t.Errorf("Expected 'false', got '%s'", cb.Value())
	}
	if cb.Checked() {
		t.Error("Expected not checked")
	}

	cb.SetChecked(true)
	if cb.Value() != "true" {
		t.Errorf("Expected 'true', got '%s'", cb.Value())
	}

	cb.SetValue("false")
	if cb.Checked() {
		t.Error("Expected not checked after SetValue('false')")
	}
}

func TestCheckbox_Focus(t *testing.T) {
	cb := NewCheckbox("Test")
	
	if cb.Focused() {
		t.Error("Should not be focused initially")
	}

	cb.Focus()
	if !cb.Focused() {
		t.Error("Should be focused after Focus()")
	}

	cb.Blur()
	if cb.Focused() {
		t.Error("Should not be focused after Blur()")
	}
}

func TestCheckbox_Update(t *testing.T) {
	cb := NewCheckbox("Test")
	cb.Focus()

	// Test Space toggle
	cb.Update(render.KeyMsg{Key: " "})
	if !cb.Checked() {
		t.Error("Space should toggle checked state")
	}

	// Test Enter toggle
	cb.Update(render.KeyMsg{Key: "enter"})
	if cb.Checked() {
		t.Error("Enter should toggle checked state")
	}
}

func TestCheckbox_View(t *testing.T) {
	cb := NewCheckbox("Test Checkbox")
	
	// Basic check that it contains label
	view := cb.View()
	if len(view) == 0 {
		t.Error("View should not be empty")
	}
	
	// We can't easily check for exact string due to ANSI codes, 
	// but we can check if it runs without panic
}
