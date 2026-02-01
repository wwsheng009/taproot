package forms

import (
	"testing"
	"unicode/utf8"
)

// TestTextInput_PasswordScrolling tests that password bullets display correctly
// when the text is longer than the display width and requires scrolling.
func TestTextInput_PasswordScrolling(t *testing.T) {
	input := NewTextInput("password")
	input.SetHidden(true)
	input.SetWidth(10) // Small width to force scrolling

	// Test: Text longer than width
	longText := "12345678901234567890" // 20 chars
	input.SetValue(longText)
	input.Focus()

	// Verify value was set correctly
	if input.Value() != longText {
		t.Fatalf("Expected value %q, got %q", longText, input.Value())
	}

	// Test view output
	view := input.View()

	// Count bullets - should be at most the width (10)
	bulletCount := utf8.RuneCountInString(view)
	for _, r := range view {
		if r != '•' && r != ' ' {
			// Allow bullets and cursor space, nothing else
			t.Errorf("View contains unexpected character: %c (U+%04X)", r, r)
		}
	}

	// We should see at most width bullets (plus maybe cursor space)
	// Since we have 20 chars but width=10, we should see max 10 bullets
	expectedMaxBullets := input.width
	if bulletCount > expectedMaxBullets+1 { // +1 for potential cursor space
		t.Errorf("Too many bullets in view: expected max %d, got %d\nView: %q",
			expectedMaxBullets+1, bulletCount, view)
	}

	// Test that no matter where cursor is, we don't get excessive bullets
	// Note: Can't directly set cursor from test (unexported)
	// But we can test different text lengths

	// Short text (no scrolling)
	input2 := NewTextInput("password")
	input2.SetHidden(true)
	input2.SetWidth(20)
	input2.SetValue("abc")
	view2 := input2.View()
	bulletCount2 := 0
	for _, r := range view2 {
		if r == '•' {
			bulletCount2++
		}
	}
	if bulletCount2 != 3 {
		t.Errorf("Short text: expected 3 bullets, got %d. View: %q", bulletCount2, view2)
	}

	// Medium text (still under width)
	input3 := NewTextInput("password")
	input3.SetHidden(true)
	input3.SetWidth(20)
	input3.SetValue("1234567890123456789") // 19 chars
	view3 := input3.View()
	bulletCount3 := 0
	for _, r := range view3 {
		if r == '•' {
			bulletCount3++
		}
	}
	if bulletCount3 != 19 {
		t.Errorf("Medium text: expected 19 bullets, got %d. View: %q", bulletCount3, view3)
	}

	// Long text (triggers scrolling)
	input4 := NewTextInput("password")
	input4.SetHidden(true)
	input4.SetWidth(20)
	input4.SetValue("123456789012345678901234567890") // 30 chars
	view4 := input4.View()
	bulletCount4 := 0
	for _, r := range view4 {
		if r == '•' {
			bulletCount4++
		}
	}
	// Should display at most width (20) bullets due to scrolling
	if bulletCount4 > 21 { // 20 + 1 for potential cursor space
		t.Errorf("Long text: expected max 21 bullets, got %d. View length: %d. View: %q",
			bulletCount4, len(view4), view4)
	}
}
