package dialogs

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/tui/util"
)

// MockDialog is a simple dialog for testing
type MockDialog struct {
	id     DialogID
	closed bool
}

func (m *MockDialog) Init() tea.Cmd {
	return nil
}

func (m *MockDialog) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	return m, nil
}

func (m *MockDialog) View() string {
	return "mock dialog"
}

func (m *MockDialog) Position() (int, int) {
	return 0, 0
}

func (m *MockDialog) ID() DialogID {
	return m.id
}

func TestDialogCmpImmutability(t *testing.T) {
	d := NewDialogCmp()

	// Store original dialog count
	originalCount := len(d.Dialogs())

	// Add a dialog
	mockDialog := &MockDialog{id: "test1"}
	openMsg := OpenDialogMsg{Model: mockDialog}
	updatedModel, _ := d.Update(openMsg)

	// Convert back to DialogCmp interface
	updated, ok := updatedModel.(DialogCmp)
	if !ok {
		t.Fatal("Updated model is not DialogCmp")
	}

	// Original should be unchanged
	if len(d.Dialogs()) != originalCount {
		t.Errorf("Original dialog count was modified. Expected %d, got %d",
			originalCount, len(d.Dialogs()))
	}

	// Updated should have dialog
	if len(updated.Dialogs()) != originalCount+1 {
		t.Errorf("Updated dialog count incorrect. Expected %d, got %d",
			originalCount+1, len(updated.Dialogs()))
	}
}

func TestDialogCmpOpenDialog(t *testing.T) {
	d := NewDialogCmp()

	mockDialog := &MockDialog{id: "test1"}
	openMsg := OpenDialogMsg{Model: mockDialog}

	// Open dialog
	updatedModel, _ := d.Update(openMsg)

	// Convert back to DialogCmp interface
	updated, ok := updatedModel.(DialogCmp)
	if !ok {
		t.Fatal("Updated model is not DialogCmp")
	}

	// Check dialog was added
	dialogs := updated.Dialogs()
	if len(dialogs) != 1 {
		t.Errorf("Expected 1 dialog, got %d", len(dialogs))
	}

	// Check active dialog
	activeID := updated.ActiveDialogID()
	if activeID != "test1" {
		t.Errorf("Expected active dialog ID 'test1', got '%s'", activeID)
	}
}

func TestDialogCmpCloseDialog(t *testing.T) {
	d := NewDialogCmp()

	// Open dialog
	mockDialog := &MockDialog{id: "test1"}
	openMsg := OpenDialogMsg{Model: mockDialog}
	dUpdated, _ := d.Update(openMsg)

	// Convert back to DialogCmp
	d, _ = dUpdated.(DialogCmp)

	// Store count
	countBefore := len(d.Dialogs())

	// Close dialog
	updatedModel, _ := d.Update(CloseDialogMsg{})

	// Convert back to DialogCmp
	updated, _ := updatedModel.(DialogCmp)

	// Check dialog was removed
	if len(updated.Dialogs()) != countBefore-1 {
		t.Errorf("Expected %d dialogs, got %d", countBefore-1, len(updated.Dialogs()))
	}

	// Original should be unchanged
	if len(d.Dialogs()) != countBefore {
		t.Errorf("Original dialog count was modified")
	}
}

func TestDialogCmpMultipleDialogs(t *testing.T) {
	d := NewDialogCmp()

	// Open multiple dialogs
	var model util.Model
	model, _ = d.Update(OpenDialogMsg{Model: &MockDialog{id: "dialog1"}})
	d, _ = model.(DialogCmp)

	model, _ = d.Update(OpenDialogMsg{Model: &MockDialog{id: "dialog2"}})
	d, _ = model.(DialogCmp)

	model, _ = d.Update(OpenDialogMsg{Model: &MockDialog{id: "dialog3"}})
	d, _ = model.(DialogCmp)

	// Check count
	if len(d.Dialogs()) != 3 {
		t.Errorf("Expected 3 dialogs, got %d", len(d.Dialogs()))
	}

	// Check active dialog (should be the last one)
	activeID := d.ActiveDialogID()
	if activeID != "dialog3" {
		t.Errorf("Expected active dialog 'dialog3', got '%s'", activeID)
	}

	// Close should remove last dialog
	updatedModel, _ := d.Update(CloseDialogMsg{})
	updated, _ := updatedModel.(DialogCmp)
	if len(updated.Dialogs()) != 2 {
		t.Errorf("Expected 2 dialogs after close, got %d", len(updated.Dialogs()))
	}

	activeID = updated.ActiveDialogID()
	if activeID != "dialog2" {
		t.Errorf("Expected active dialog 'dialog2', got '%s'", activeID)
	}
}

func TestDialogCmpCloseAllDialogs(t *testing.T) {
	d := NewDialogCmp()

	// Open multiple dialogs
	var model util.Model
	model, _ = d.Update(OpenDialogMsg{Model: &MockDialog{id: "dialog1"}})
	d, _ = model.(DialogCmp)

	model, _ = d.Update(OpenDialogMsg{Model: &MockDialog{id: "dialog2"}})
	d, _ = model.(DialogCmp)

	model, _ = d.Update(OpenDialogMsg{Model: &MockDialog{id: "dialog3"}})
	d, _ = model.(DialogCmp)

	// Store original
	original := d

	// Close all
	updatedModel, _ := d.Update(CloseDialogMsg{})
	updated, _ := updatedModel.(DialogCmp)
	updatedModel, _ = updated.Update(CloseDialogMsg{})
	updated, _ = updatedModel.(DialogCmp)
	updatedModel, _ = updated.Update(CloseDialogMsg{})
	updated, _ = updatedModel.(DialogCmp)

	// Check no dialogs
	if len(updated.Dialogs()) != 0 {
		t.Errorf("Expected 0 dialogs, got %d", len(updated.Dialogs()))
	}

	// Original should be unchanged
	if len(original.Dialogs()) != 3 {
		t.Errorf("Original dialog count was modified")
	}
}

func TestDialogCmpHasDialogs(t *testing.T) {
	d := NewDialogCmp()

	// Initially no dialogs
	if d.HasDialogs() {
		t.Error("Should not have dialogs initially")
	}

	// Add dialog
	model, _ := d.Update(OpenDialogMsg{Model: &MockDialog{id: "test1"}})
	d, _ = model.(DialogCmp)

	// Now should have dialogs
	if !d.HasDialogs() {
		t.Error("Should have dialogs after adding")
	}
}
