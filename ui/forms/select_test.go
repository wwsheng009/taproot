package forms

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSelect_Init(t *testing.T) {
	opts := []string{"A", "B", "C"}
	s := NewSelect("Test", opts)

	if s.Value() != "" {
		t.Errorf("Expected empty value, got %q", s.Value())
	}
	if s.SelectedIndex() != -1 {
		t.Errorf("Expected -1 index, got %d", s.SelectedIndex())
	}
	if s.expanded {
		t.Error("Expected collapsed by default")
	}
}

func TestSelect_Selection(t *testing.T) {
	opts := []string{"A", "B", "C"}
	s := NewSelect("Test", opts)

	s.SetSelectedIndex(1)
	if s.Value() != "B" {
		t.Errorf("Expected 'B', got %q", s.Value())
	}

	s.SetValue("C")
	if s.SelectedIndex() != 2 {
		t.Errorf("Expected index 2, got %d", s.SelectedIndex())
	}

	s.SetValue("Invalid")
	if s.SelectedIndex() != -1 {
		t.Errorf("Expected index -1 for invalid value, got %d", s.SelectedIndex())
	}
}

func TestSelect_Interaction(t *testing.T) {
	opts := []string{"Option 1", "Option 2"}
	s := NewSelect("Test", opts)
	s.Focus()

	// 1. Enter to expand
	s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !s.expanded {
		t.Error("Expected expanded after Enter")
	}

	// 2. Down to highlight next (default highlighted is 0)
	s.Update(tea.KeyMsg{Type: tea.KeyDown})
	if s.highlighted != 1 {
		t.Errorf("Expected highlighted 1, got %d", s.highlighted)
	}

	// 3. Enter to select
	s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if s.expanded {
		t.Error("Expected collapsed after selection")
	}
	if s.Value() != "Option 2" {
		t.Errorf("Expected 'Option 2' selected, got %q", s.Value())
	}

	// 4. Blur should collapse if open
	s.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Open again
	if !s.expanded {
		t.Error("Failed to reopen")
	}
	s.Blur()
	if s.expanded {
		t.Error("Expected collapsed after Blur")
	}
}
