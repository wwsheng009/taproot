package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDialogFooterVisibility(t *testing.T) {
	// 1. Initialize
	m := NewDialogsModel()
	m.Init()

	// 2. Set Size (80x24)
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Check initial view
	view := m.View()
	lines := strings.Split(view, "\n")
	if len(lines) != 24 {
		t.Errorf("Expected 24 lines, got %d", len(lines))
	}
	if !strings.Contains(view, "Navigate") {
		t.Error("Footer missing in initial view")
	}

	// 3. Open Dialog (Space)
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	
	// Verify dialog is open
	if !m.overlay.HasDialogs() {
		t.Fatal("Dialog failed to open")
	}

	// View should still be 24 lines (handled by Overlay.Render -> eventually, but Overlay might handle it differently)
	// Actually Overlay.Render now adapts to background size.
	// Since background is 24 lines (padded), Overlay.Render will be at least 24 lines.
	viewDialog := m.View()
	if len(strings.Split(viewDialog, "\n")) < 24 {
		t.Errorf("Expected at least 24 lines with dialog, got %d", len(strings.Split(viewDialog, "\n")))
	}

	// 4. Close Dialog (Enter)
	// We need to send Enter to the DIALOG. 
	// The model forwards messages to updateDialog.
	// InfoDialog closes on Enter.
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify dialog is closed
	if m.overlay.HasDialogs() {
		// Maybe it needs an update cycle?
		// updateDialog calls m.overlay.Pop() immediately if View() == "".
		// InfoDialog sets quitting=true on Enter, so next View() is "".
		// We need to ensure the update loop processes it.
		// Let's send another update or check state.
		// The updateDialog function returns the updated model.
		// Wait, m.updateDialog returns (m, cmd). It modifies m.overlay in place.
		// Let's check if it actually popped.
		// InfoDialog.Update sets quitting=true.
		// Then updateDialog checks updatedDialog.View() == "".
		// InfoDialog.View() returns "" if quitting.
		// So it should pop.
		t.Fatal("Dialog failed to close")
	}

	// 5. Verify Footer is present
	finalView := m.View()
	if !strings.Contains(finalView, "Navigate") {
		t.Error("Footer missing after closing dialog")
	}
	
	finalLines := strings.Split(finalView, "\n")
	if len(finalLines) != 24 {
		t.Errorf("Expected 24 lines after closing dialog, got %d", len(finalLines))
	}
	
	// Check for "ghosting" or garbage?
	// Hard to check programmatically without a reference, but the line count is the key metric for screen clearing.
}
