package forms

import (
	"testing"

	"github.com/wwsheng009/taproot/ui/render"
)

func TestRadioGroup_Init(t *testing.T) {
	rg := NewRadioGroup("Colors", []string{"Red", "Green", "Blue"})
	if rg.Init() != nil {
		t.Error("Init should return nil")
	}
}

func TestRadioGroup_Selection(t *testing.T) {
	opts := []string{"Red", "Green", "Blue"}
	rg := NewRadioGroup("Colors", opts)

	// Default selection
	if rg.Value() != "Red" {
		t.Errorf("Expected default selection 'Red', got '%s'", rg.Value())
	}

	// Set value
	rg.SetValue("Green")
	if rg.Value() != "Green" {
		t.Errorf("Expected 'Green', got '%s'", rg.Value())
	}

	// Set invalid value -> default to first
	rg.SetValue("Yellow")
	if rg.Value() != "Red" {
		t.Errorf("Expected fallback to 'Red' for invalid value, got '%s'", rg.Value())
	}

	// Set index
	rg.SetSelectedIndex(2)
	if rg.Value() != "Blue" {
		t.Errorf("Expected 'Blue', got '%s'", rg.Value())
	}

	// Set invalid index
	rg.SetSelectedIndex(99)
	if rg.Value() != "Blue" {
		t.Error("Should ignore invalid index")
	}
}

func TestRadioGroup_Navigation(t *testing.T) {
	rg := NewRadioGroup("Colors", []string{"Red", "Green", "Blue"})
	rg.Focus()

	// Down (j)
	rg.Update(render.KeyMsg{Key: "j"})
	if rg.Value() != "Green" {
		t.Errorf("Down should select Green, got %s", rg.Value())
	}

	// Down (down arrow)
	rg.Update(render.KeyMsg{Key: "down"})
	if rg.Value() != "Blue" {
		t.Errorf("Down should select Blue, got %s", rg.Value())
	}

	// Wrap around (down)
	rg.Update(render.KeyMsg{Key: "down"})
	if rg.Value() != "Red" {
		t.Errorf("Wrap around should select Red, got %s", rg.Value())
	}

	// Up (k)
	rg.Update(render.KeyMsg{Key: "k"})
	if rg.Value() != "Blue" {
		t.Errorf("Up should select Blue (wrap), got %s", rg.Value())
	}
}

func TestRadioGroup_View(t *testing.T) {
	rg := NewRadioGroup("Colors", []string{"Red", "Green"})
	view := rg.View()
	
	if len(view) == 0 {
		t.Error("View should not be empty")
	}
}
