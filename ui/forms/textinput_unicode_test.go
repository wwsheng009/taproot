package forms

import (
	"testing"
	"unicode/utf8"
)

func TestTextInput_PasswordMasking(t *testing.T) {
	input := NewTextInput("password")
	input.SetHidden(true)
	input.SetValue("test123")

	// Test raw value
	if input.Value() != "test123" {
		t.Errorf("Expected 'test123', got %q", input.Value())
	}

	// Test rune count should be 7
	if utf8.RuneCountInString(input.Value()) != 7 {
		t.Errorf("Expected 7 runes, got %d", utf8.RuneCountInString(input.Value()))
	}

	// Focus and test view
	input.Focus()
	view := input.View()

	// Count bullets in view - should be 7
	bulletCount := 0
	for _, r := range view {
		if r == '•' {
			bulletCount++
		}
	}

	if bulletCount != 7 {
		t.Errorf("Expected 7 bullets, got %d\nView: %q", bulletCount, view)
	}
}

func TestTextInput_RuneCount(t *testing.T) {
	input := NewTextInput("test")

	// ASCII characters
	input.SetValue("hello")
	if utf8.RuneCountInString(input.Value()) != 5 {
		t.Errorf("Expected 5 runes, got %d", utf8.RuneCountInString(input.Value()))
	}

	// Mixed ASCII
	input.SetValue("test123")
	if utf8.RuneCountInString(input.Value()) != 7 {
		t.Errorf("Expected 7 runes, got %d", utf8.RuneCountInString(input.Value()))
	}
}
