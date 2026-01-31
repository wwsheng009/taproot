package notify

import (
	"testing"
	"time"

	"github.com/wwsheng009/taproot/ui/render"
)

func TestManager_AddNotification(t *testing.T) {
	cfg := DefaultConfig()
	m := NewManager(cfg)

	n := Notification{
		ID:       "test-1",
		Title:    "Test",
		Message:  "Message",
		Duration: time.Second,
	}

	msg := ShowNotificationMsg{Notification: n}
	model, cmd := m.Update(msg)

	// Check type assertion
	nm, ok := model.(*Manager)
	if !ok {
		t.Fatalf("Expected *Manager, got %T", model)
	}

	// Check notification added
	if len(nm.notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(nm.notifications))
	}

	// Check command (should be tick)
	if cmd == nil {
		t.Error("Expected tick command, got nil")
	}
}

func TestManager_MaxVisible(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxVisible = 2
	m := NewManager(cfg)

	// Add 3 notifications
	for i := 0; i < 3; i++ {
		n := Notification{
			ID:       "test",
			Title:    "Test",
			Duration: time.Second,
		}
		m.Update(ShowNotificationMsg{Notification: n})
	}

	// Should have 2 (MaxVisible)
	if len(m.notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(m.notifications))
	}
}

func TestManager_Expiry(t *testing.T) {
	cfg := DefaultConfig()
	m := NewManager(cfg)

	// Add notification
	n := Notification{
		ID:       "test-expiry",
		Title:    "Test",
		Duration: 100 * time.Millisecond,
	}
	m.Update(ShowNotificationMsg{Notification: n})

	// Verify added
	if len(m.notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(m.notifications))
	}

	// Capture CreatedAt from the added notification
	addedN := m.notifications[0]

	// Simulate tick after expiry
	// tickMsg is private, so we can't create it directly from outside package?
	// Wait, we are in package `notify` (same package), so we can access it.
	
	future := addedN.CreatedAt.Add(200 * time.Millisecond)
	tick := tickMsg{time: future}

	model, _ := m.Update(tick)
	nm := model.(*Manager)

	// Should be removed
	if len(nm.notifications) != 0 {
		t.Errorf("Expected 0 notifications after expiry, got %d", len(nm.notifications))
	}
}

func TestManager_View_Empty(t *testing.T) {
	cfg := DefaultConfig()
	m := NewManager(cfg)

	if view := m.View(); view != "" {
		t.Errorf("Expected empty view for no notifications, got %q", view)
	}
}

func TestManager_WindowSize(t *testing.T) {
	cfg := DefaultConfig()
	m := NewManager(cfg)

	m.Update(render.WindowSizeMsg{Width: 100, Height: 100})

	if m.width != 100 || m.height != 100 {
		t.Errorf("Expected size 100x100, got %dx%d", m.width, m.height)
	}
}
