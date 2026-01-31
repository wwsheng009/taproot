package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/tui/page"
	"github.com/wwsheng009/taproot/tui/util"
)

// MockModel is a simple model for testing
type MockModel struct {
	util.Model
	initCalled bool
}

func (m *MockModel) Init() tea.Cmd {
	m.initCalled = true
	return nil
}

func (m *MockModel) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	return m, nil
}

func (m *MockModel) View() string {
	return "mock view"
}

func TestAppModelImmutability(t *testing.T) {
	app := NewApp()

	// Store original state
	originalQuitting := app.quitting

	// Apply update that modifies state
	updatedModel, _ := app.Update(tea.QuitMsg{})

	// Convert to concrete type
	updated, ok := updatedModel.(AppModel)
	if !ok {
		t.Fatal("Updated model is not AppModel")
	}

	// Original should be unchanged
	if app.quitting != originalQuitting {
		t.Errorf("Original model was modified. Expected quitting=%v, got %v",
			originalQuitting, app.quitting)
	}

	// Updated should have new state
	if updated.quitting != true {
		t.Errorf("Updated model quitting not set. Expected true, got %v",
			updated.quitting)
	}
}

func TestAppModelSetPage(t *testing.T) {
	app := NewApp()

	// Register a page
	mockModel := &MockModel{}
	app.RegisterPage("test", mockModel)

	// Store original state
	originalPage := app.currentPage
	originalStackLen := len(app.pageStack)

	// Set new page
	newApp := app.SetPage("test")

	// Original should be unchanged
	if app.currentPage != originalPage {
		t.Errorf("Original currentPage was modified")
	}
	if len(app.pageStack) != originalStackLen {
		t.Errorf("Original pageStack was modified")
	}

	// New model should have new page
	if newApp.currentPage != "test" {
		t.Errorf("New page not set. Expected 'test', got '%s'", newApp.currentPage)
	}
	if len(newApp.pageStack) != originalStackLen {
		t.Errorf("pageStack should not change on first SetPage")
	}
}

func TestAppModelPageChangeMsg(t *testing.T) {
	app := NewApp()

	// Register pages
	mockModel1 := &MockModel{}
	mockModel2 := &MockModel{}
	app.RegisterPage("page1", mockModel1)
	app.RegisterPage("page2", mockModel2)

	// Set initial page
	app.currentPage = "page1"

	// Apply page change
	updatedModel, _ := app.Update(page.PageChangeMsg{ID: "page2"})

	// Convert to concrete type
	updated, ok := updatedModel.(AppModel)
	if !ok {
		t.Fatal("Updated model is not AppModel")
	}

	// Original should be unchanged
	if app.currentPage != "page1" {
		t.Errorf("Original currentPage was modified")
	}

	// Updated should have new page
	if updated.currentPage != "page2" {
		t.Errorf("Page not changed. Expected 'page2', got '%s'", updated.currentPage)
	}

	// Page stack should have old page
	if len(updated.pageStack) != 1 {
		t.Errorf("Page stack should have 1 page, got %d", len(updated.pageStack))
	}
	if updated.pageStack[0] != "page1" {
		t.Errorf("Old page not in stack. Expected 'page1', got '%s'", updated.pageStack[0])
	}
}

func TestAppModelPageBackMsg(t *testing.T) {
	app := NewApp()

	// Register pages
	mockModel1 := &MockModel{}
	mockModel2 := &MockModel{}
	app.RegisterPage("page1", mockModel1)
	app.RegisterPage("page2", mockModel2)

	// Setup initial state
	app.currentPage = "page2"
	app.pageStack = []page.PageID{"page1"}

	// Apply page back
	updatedModel, _ := app.Update(page.PageBackMsg{})

	// Convert to concrete type
	updated, ok := updatedModel.(AppModel)
	if !ok {
		t.Fatal("Updated model is not AppModel")
	}

	// Original should be unchanged
	if app.currentPage != "page2" {
		t.Errorf("Original currentPage was modified")
	}

	// Updated should have gone back to page1
	if updated.currentPage != "page1" {
		t.Errorf("Page not gone back. Expected 'page1', got '%s'", updated.currentPage)
	}

	// Page stack should be empty
	if len(updated.pageStack) != 0 {
		t.Errorf("Page stack should be empty, got %d", len(updated.pageStack))
	}
}

func TestAppModelWindowSizeMsg(t *testing.T) {
	app := NewApp()

	// Store original size
	originalWidth := app.width
	originalHeight := app.height

	// Apply window size change
	updatedModel, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Convert to concrete type
	updated, ok := updatedModel.(AppModel)
	if !ok {
		t.Fatal("Updated model is not AppModel")
	}

	// Original should be unchanged
	if app.width != originalWidth || app.height != originalHeight {
		t.Errorf("Original size was modified")
	}

	// Updated should have new size
	if updated.width != 120 || updated.height != 40 {
		t.Errorf("Size not updated. Expected (120, 40), got (%d, %d)",
			updated.width, updated.height)
	}
}

func TestAppModelLifecycle(t *testing.T) {
	app := NewApp()

	// Register pages
	mockModel1 := &MockModel{}
	mockModel2 := &MockModel{}
	app.RegisterPage("page1", mockModel1)
	app.RegisterPage("page2", mockModel2)

	// Set page1
	app.currentPage = "page1"
	app.lifecycleMgr.Register("page1")

	// Store original context
	ctx1, exists1 := app.GetPageContext("page1")
	if !exists1 {
		t.Error("Context for page1 should exist")
	}

	// Change to page2
	updatedModel, _ := app.Update(page.PageChangeMsg{ID: "page2"})
	updated, _ := updatedModel.(AppModel)

	// Original context should be cancelled
	select {
	case <-ctx1.Done():
		// Expected
	default:
		t.Error("Old page context should be cancelled")
	}

	// New context should exist
	ctx2, exists2 := updated.GetPageContext("page2")
	if !exists2 {
		t.Error("Context for page2 should exist")
	}
	if ctx2 == nil {
		t.Error("Context for page2 should not be nil")
	}
}
