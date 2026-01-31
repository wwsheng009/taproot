package messages

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMessagesModelImmutability(t *testing.T) {
	// Test that Update returns new instance
	m := New()
	// Add some content so we can scroll
	m.messages = []Message{
		{Role: "user", Content: "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6"},
	}
	m.height = 3

	// Store original scroll (need to access private field for testing)
	originalScroll := m.scroll

	// Apply update
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	// Convert to concrete type
	updated, ok := updatedModel.(*MessagesModel)
	if !ok {
		t.Fatal("Updated model is not *MessagesModel")
	}

	// Original should be unchanged
	if m.scroll != originalScroll {
		t.Errorf("Original model was modified. Expected scroll=%d, got %d",
			originalScroll, m.scroll)
	}

	// Updated should have new scroll
	if updated.scroll != originalScroll+1 {
		t.Errorf("Updated model scroll not incremented. Expected %d, got %d",
			originalScroll+1, updated.scroll)
	}
}

func TestMessagesModelAddMessage(t *testing.T) {
	m := New()
	originalCount := len(m.messages)

	// Add message
	msg := Message{
		Role:    "user",
		Content: "Hello",
	}
	newModel := m.AddMessage(msg)

	// Original should be unchanged
	if len(m.messages) != originalCount {
		t.Errorf("Original model was modified. Expected %d messages, got %d",
			originalCount, len(m.messages))
	}

	// New model should have message
	if len(newModel.messages) != originalCount+1 {
		t.Errorf("New model does not have new message. Expected %d messages, got %d",
			originalCount+1, len(newModel.messages))
	}
}

func TestMessagesModelClear(t *testing.T) {
	m := New()
	m.messages = []Message{
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi"},
	}

	// Clear messages
	newModel := m.Clear()

	// Original should be unchanged
	if len(m.messages) != 2 {
		t.Errorf("Original model was modified. Expected 2 messages, got %d",
			len(m.messages))
	}

	// New model should be empty
	if len(newModel.messages) != 0 {
		t.Errorf("New model not cleared. Expected 0 messages, got %d",
			len(newModel.messages))
	}
}

func TestMessagesModelSetSize(t *testing.T) {
	m := New()
	originalWidth := m.width
	originalHeight := m.height

	// Set new size
	newModel := m.SetWidth(100).SetHeight(50)

	// Original should be unchanged
	if m.width != originalWidth || m.height != originalHeight {
		t.Errorf("Original model was modified")
	}

	// New model should have new size
	if newModel.width != 100 || newModel.height != 50 {
		t.Errorf("New model size not set correctly. Expected (100, 50), got (%d, %d)",
			newModel.width, newModel.height)
	}
}

func TestMessagesModelScroll(t *testing.T) {
	m := New()
	m.messages = []Message{
		{Role: "user", Content: "Line 1"},
		{Role: "user", Content: "Line 2"},
		{Role: "user", Content: "Line 3"},
	}
	m.SetHeight(5)

	// Scroll down
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	updated, _ := updatedModel.(*MessagesModel)

	if updated.scroll != 1 {
		t.Errorf("Scroll not incremented. Expected 1, got %d", updated.scroll)
	}

	// Original should be unchanged
	if m.scroll != 0 {
		t.Errorf("Original model was modified. Expected scroll=0, got %d", m.scroll)
	}

	// Scroll up
	updatedModel, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	updated, _ = updatedModel.(*MessagesModel)

	if updated.scroll != 0 {
		t.Errorf("Scroll not decremented. Expected 0, got %d", updated.scroll)
	}
}
