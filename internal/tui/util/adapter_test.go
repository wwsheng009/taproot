package util

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// mockBubbleteaModel is a simple implementation of util.Model for testing.
type mockBubbleteaModel struct {
	value int
}

func (m *mockBubbleteaModel) Init() tea.Cmd {
	return nil
}

func (m *mockBubbleteaModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case int:
		newModel := *m
		newModel.value = msg
		return &newModel, nil
	}
	return m, nil
}

func (m *mockBubbleteaModel) View() string {
	return "value: " + string(rune(m.value))
}

// TestBubbleteaToRenderModel_Init 验证 Init 方法
func TestBubbleteaToRenderModel_Init(t *testing.T) {
	inner := &mockBubbleteaModel{value: 0}
	adapter := NewBubbleteaToRenderModel(inner)

	err := adapter.Init()
	if err != nil {
		t.Errorf("Init should return nil, got %v", err)
	}
}

// TestBubbleteaToRenderModel_Update 验证 Update 方法
func TestBubbleteaToRenderModel_Update(t *testing.T) {
	inner := &mockBubbleteaModel{value: 0}
	adapter := NewBubbleteaToRenderModel(inner)

	// Update with a tea.Msg (int is a valid tea.Msg)
	newAdapter, cmd := adapter.Update(tea.Msg(42))
	if cmd != nil {
		t.Errorf("Update should return nil command, got %v", cmd)
	}

	// Should return new instance
	if newAdapter == adapter {
		t.Error("Update should return new instance")
	}

	// Type assert to get the new adapter
	newAdap, ok := newAdapter.(*BubbleteaToRenderModel)
	if !ok {
		t.Fatal("newAdapter should be *BubbleteaToRenderModel")
	}

	// Inner model should be updated
	if newAdap.inner.(*mockBubbleteaModel).value != 42 {
		t.Errorf("Inner model should be updated to 42, got %d",
			newAdap.inner.(*mockBubbleteaModel).value)
	}

	// Original adapter should not be modified
	if adapter.inner.(*mockBubbleteaModel).value != 0 {
		t.Errorf("Original adapter should not be modified, got %d", adapter.inner.(*mockBubbleteaModel).value)
	}
}

// TestBubbleteaToRenderModel_View 验证 View 方法
func TestBubbleteaToRenderModel_View(t *testing.T) {
	inner := &mockBubbleteaModel{value: 65} // ASCII 'A'
	adapter := NewBubbleteaToRenderModel(inner)

	view := adapter.View()
	if view != "value: A" {
		t.Errorf("View should return 'value: A', got %s", view)
	}
}

// TestBubbleteaToRenderModel_GetInner 验证 GetInner 方法
func TestBubbleteaToRenderModel_GetInner(t *testing.T) {
	inner := &mockBubbleteaModel{value: 42}
	adapter := NewBubbleteaToRenderModel(inner)

	retrievedInner := adapter.GetInner()
	if retrievedInner != inner {
		t.Error("GetInner should return the original inner model")
	}

	// Verify it's the same instance
	if retrievedInner.(*mockBubbleteaModel).value != 42 {
		t.Errorf("Retrieved inner model should have value 42, got %d", retrievedInner.(*mockBubbleteaModel).value)
	}
}

// TestBubbleteaToRenderModel_WithInner 验证 WithInner 方法
func TestBubbleteaToRenderModel_WithInner(t *testing.T) {
	inner1 := &mockBubbleteaModel{value: 42}
	inner2 := &mockBubbleteaModel{value: 100}
	adapter := NewBubbleteaToRenderModel(inner1)

	newAdapter := adapter.WithInner(inner2)

	// Should return new instance
	if newAdapter == adapter {
		t.Error("WithInner should return new instance")
	}

	// New adapter should have new inner model
	if newAdapter.inner.(*mockBubbleteaModel).value != 100 {
		t.Errorf("New adapter should have inner model with value 100, got %d", newAdapter.inner.(*mockBubbleteaModel).value)
	}

	// Original adapter should not be modified
	if adapter.inner.(*mockBubbleteaModel).value != 42 {
		t.Errorf("Original adapter should not be modified, got %d", adapter.inner.(*mockBubbleteaModel).value)
	}
}

// mockViewable is a simple type that implements View() string
type mockViewable struct {
	content string
}

func (m mockViewable) View() string {
	return m.content
}

// TestRenderToBubbleteaModel_Init 验证 Init 方法
func TestRenderToBubbleteaModel_Init(t *testing.T) {
	inner := mockViewable{content: "test"}
	adapter := NewRenderToBubbleteaModel(inner)

	cmd := adapter.Init()
	if cmd != nil {
		t.Errorf("Init should return nil command, got %v", cmd)
	}
}

// TestRenderToBubbleteaModel_Update 验证 Update 方法
func TestRenderToBubbleteaModel_Update(t *testing.T) {
	inner := mockViewable{content: "test"}
	adapter := NewRenderToBubbleteaModel(inner)

	newModel, cmd := adapter.Update(tea.KeyMsg{})

	if cmd != nil {
		t.Errorf("Update should return nil command, got %v", cmd)
	}

	// Should return same instance (no state changes)
	if newModel != adapter {
		t.Error("Update should return same instance")
	}
}

// TestRenderToBubbleteaModel_View 验证 View 方法
func TestRenderToBubbleteaModel_View(t *testing.T) {
	inner := mockViewable{content: "test view"}
	adapter := NewRenderToBubbleteaModel(inner)

	view := adapter.View()
	if view != "test view" {
		t.Errorf("View should return 'test view', got %s", view)
	}
}
