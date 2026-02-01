package forms

import (
	"errors"
	"testing"
)

func TestNewTextInput(t *testing.T) {
	input := NewTextInput("test placeholder")
	if input == nil {
		t.Fatal("returned nil input")
	}
	if input.placeholder != "test placeholder" {
		t.Errorf("expected placeholder 'test placeholder', got %q", input.placeholder)
	}
	if input.Value() != "" {
		t.Errorf("expected empty value, got %q", input.Value())
	}
}

func TestTextInput_Insert(t *testing.T) {
	input := NewTextInput("test")
	input.Focus()

	// Test basic insert
	input.insert('a')
	if input.Value() != "a" {
		t.Errorf("expected 'a', got %q", input.Value())
	}

	// Test insert at position
	input.insert('b')
	input.insert('c')
	input.cursor = 1
	input.insert('x')
	if input.Value() != "axbc" {
		t.Errorf("expected 'axbc', got %q", input.Value())
	}
}

func TestTextInput_DeleteBefore(t *testing.T) {
	input := NewTextInput("test")
	input.Focus()

	// Insert some characters
	for _, r := range "hello" {
		input.insert(r)
	}

	// Test backspace
	input.cursor = 3
	input.deleteBefore()
	if input.Value() != "helo" {
		t.Errorf("expected 'helo', got %q", input.Value())
	}
	if input.cursor != 2 {
		t.Errorf("expected cursor at 2, got %d", input.cursor)
	}
}

func TestTextInput_MaxLength(t *testing.T) {
	input := NewTextInput("test")
	input.SetMaxLength(5)
	input.Focus()

	for _, r := range "abcdefgh" {
		input.insert(r)
	}

	if input.Value() != "abcde" {
		t.Errorf("expected 'abcde', got %q", input.Value())
	}
}

func TestTextInput_Validation(t *testing.T) {
	input := NewTextInput("email")

	// Test Required validator
	input.AddValidator(Required)
	input.SetValue("test@example.com")
	err := input.Validate()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	input.SetValue("")
	err = input.Validate()
	if err == nil {
		t.Error("expected error for empty value")
	}

	// Test Email validator (replaces validators)
	input.SetValidators(Email)
	input.SetValue("test@example.com")
	err = input.Validate()
	if err != nil {
		t.Errorf("expected nil error for valid email, got %v", err)
	}

	input.SetValue("invalid-email")
	err = input.Validate()
	if err == nil {
		t.Error("expected error for invalid email")
	}

	// Test Multiple validators
	input.AddValidator(Required)
	input.SetValue("test@example.com")
	err = input.Validate()
	if err != nil {
		t.Errorf("expected nil error with multiple validators, got %v", err)
	}
}

func TestTextInput_Hidden(t *testing.T) {
	input := NewTextInput("password")
	input.SetHidden(true)
	input.Focus()

	input.SetValue("secretpass")
	if input.Value() != "secretpass" {
		t.Errorf("expected 'secretpass', got %q", input.Value())
	}

	// Hidden input should display dots in View()
	view := input.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestNewNumberInput(t *testing.T) {
	input := NewNumberInput("age")
	if input == nil {
		t.Fatal("returned nil input")
	}
	if input.step != 1 {
		t.Errorf("expected step 1, got %f", input.step)
	}
	if input.min != -1e9 || input.max != 1e9 {
		t.Error("unexpected default range")
	}
}

func TestNumberInput_Increment(t *testing.T) {
	input := NewNumberInput("age")
	input.SetFloatValue(25)
	input.Increment()
	if input.FloatValue() != 26 {
		t.Errorf("expected 26, got %f", input.FloatValue())
	}
}

func TestNumberInput_Decrement(t *testing.T) {
	input := NewNumberInput("age")
	input.SetFloatValue(25)
	input.Decrement()
	if input.FloatValue() != 24 {
		t.Errorf("expected 24, got %f", input.FloatValue())
	}
}

func TestNumberInput_Range(t *testing.T) {
	input := NewNumberInput("age")
	input.SetRange(0, 100)

	input.SetFloatValue(50)
	if input.FloatValue() != 50 {
		t.Errorf("expected 50, got %f", input.FloatValue())
	}

	// Test out of range - should clamp
	input.SetFloatValue(150)
	if input.FloatValue() != 100 {
		t.Errorf("expected 100 (clamped), got %f", input.FloatValue())
	}

	input.SetFloatValue(-10)
	if input.FloatValue() != 0 {
		t.Errorf("expected 0 (clamped), got %f", input.FloatValue())
	}
}

func TestNumberInput_Validation(t *testing.T) {
	input := NewNumberInput("age")
	input.SetRange(18, 100)

	// Valid
	input.SetValue("25")
	err := input.Validate()
	if err != nil {
		t.Errorf("expected nil error for valid age, got %v", err)
	}

	// Below minimum
	input.SetValue("10")
	err = input.Validate()
	if err == nil {
		t.Error("expected error for age below minimum")
	}

	// Above maximum
	input.SetValue("150")
	err = input.Validate()
	if err == nil {
		t.Error("expected error for age above maximum")
	}

	// Invalid number
	input.SetValue("not-a-number")
	err = input.Validate()
	if err == nil {
		t.Error("expected error for non-numeric value")
	}

	// Empty should be allowed by default
	input.SetValue("")
	err = input.Validate()
	if err != nil {
		t.Errorf("expected nil error for empty value, got %v", err)
	}
}

func TestNumberInput_Precision(t *testing.T) {
	input := NewNumberInput("price")
	input.SetPrecision(2)
	input.SetFloatValue(123.4567)

	if input.Value() != "123.46" {
		t.Errorf("expected '123.46', got %q", input.Value())
	}
}

func TestNewTextArea(t *testing.T) {
	area := NewTextArea("bio")
	if area == nil {
		t.Fatal("returned nil area")
	}
	if area.height != 5 {
		t.Errorf("expected height 5, got %d", area.height)
	}
	if len(area.lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(area.lines))
	}
}

func TestTextArea_Insert(t *testing.T) {
	area := NewTextArea("bio")

	// Test basic insert
	area.Insert('H')
	area.Insert('i')
	if area.Value() != "Hi" {
		t.Errorf("expected 'Hi', got %q", area.Value())
	}
}

func TestTextArea_InsertNewline(t *testing.T) {
	area := NewTextArea("bio")
	area.Insert('H')
	area.Insert('i')

	area.InsertNewline()
	area.Insert('B')
	area.Insert('y')
	area.Insert('e')

	if area.Value() != "Hi\nBye" {
		t.Errorf("expected 'Hi\\nBye', got %q", area.Value())
	}
	if len(area.lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(area.lines))
	}
}

func TestTextArea_Delete(t *testing.T) {
	area := NewTextArea("bio")

	// Insert characters
	for _, r := range "Hello" {
		area.Insert(r)
	}

	// Test backspace at end
	area.cursorCol = 5
	area.Delete()
	if area.Value() != "Hell" {
		t.Errorf("expected 'Hell', got %q", area.Value())
	}

	// Test delete in middle
	area.cursorCol = 2
	area.Delete()
	if area.Value() != "Hll" {
		t.Errorf("expected 'Hll', got %q", area.Value())
	}

	// Test delete at start (should do nothing)
	area.cursorCol = 0
	area.cursorRow = 0
	area.Delete()
	if area.Value() != "Hll" {
		t.Errorf("expected 'Hll', got %q", area.Value())
	}
}

func TestTextArea_MultiLineDelete(t *testing.T) {
	area := NewTextArea("bio")

	// Create two lines
	area.Insert('A')
	area.Insert('B')
	area.InsertNewline()
	area.Insert('C')
	area.Insert('D')

	// Cursor at start of second line
	area.cursorCol = 0
	area.cursorRow = 1

	// Delete should merge lines
	area.Delete()

	if area.Value() != "ABCD" {
		t.Errorf("expected 'ABCD', got %q", area.Value())
	}
	if len(area.lines) != 1 {
		t.Errorf("expected 1 line after merge, got %d", len(area.lines))
	}
}

func TestTextArea_SetValue(t *testing.T) {
	area := NewTextArea("bio")

	area.SetValue("Line 1\nLine 2\nLine 3")

	if len(area.lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(area.lines))
	}
	if area.lines[0] != "Line 1" {
		t.Errorf("expected 'Line 1', got %q", area.lines[0])
	}
	if area.lines[1] != "Line 2" {
		t.Errorf("expected 'Line 2', got %q", area.lines[1])
	}
	if area.lines[2] != "Line 3" {
		t.Errorf("expected 'Line 3', got %q", area.lines[2])
	}
}

func TestTextArea_Validation(t *testing.T) {
	area := NewTextArea("bio")

	// Test MinLength validator
	area.AddValidator(MinLength(10))
	area.SetValue("Short")
	err := area.Validate()
	if err == nil {
		t.Error("expected error for short text")
	}

	area.SetValue("This is long enough text")
	err = area.Validate()
	if err != nil {
		t.Errorf("expected nil error for valid length, got %v", err)
	}

	// Test Multiple validators
	area.AddValidator(Required)
	area.SetValue("")
	err = area.Validate()
	if err == nil {
		t.Error("expected error for empty value")
	}
}

// Test Focus/Blur methods
func TestTextInput_Focus(t *testing.T) {
	input := NewTextInput("test")

	input.Focus()
	if !input.Focused() {
		t.Error("expected focused to be true after Focus()")
	}

	input.Blur()
	if input.Focused() {
		t.Error("expected focused to be false after Blur()")
	}
}

func TestNumberInput_Focus(t *testing.T) {
	input := NewNumberInput("age")

	input.Focus()
	if !input.Focused() {
		t.Error("expected focused to be true after Focus()")
	}
}

func TestTextArea_Focus(t *testing.T) {
	area := NewTextArea("bio")

	area.Focus()
	if !area.Focused() {
		t.Error("expected focused to be true after Focus()")
	}
}

// Test custom validators
func Test_CustomValidator(t *testing.T) {
	input := NewTextInput("username")

	// Custom validator: username must be lowercase
	customValidator := func(v string) error {
		for _, r := range v {
			if r >= 'A' && r <= 'Z' {
				return errors.New("username must be lowercase only")
			}
		}
		return nil
	}

	input.AddValidator(customValidator)

	input.SetValue("validusername")
	err := input.Validate()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	input.SetValue("InvalidUsername")
	err = input.Validate()
	if err == nil {
		t.Error("expected error for uppercase letters")
	}
}
