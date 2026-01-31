package lifecycle

import (
	"testing"
	"time"
)

func TestLifecycleManager_Register(t *testing.T) {
	lm := NewLifecycleManager()

	ctx1 := lm.Register("page1")
	ctx2 := lm.Register("page2")

	// Check both contexts exist
	if ctx1 == nil {
		t.Error("First context is nil")
	}
	if ctx2 == nil {
		t.Error("Second context is nil")
	}

	// Check they are different
	if ctx1 == ctx2 {
		t.Error("Contexts should be different")
	}

	// Check active count
	if lm.ActiveCount() != 2 {
		t.Errorf("Expected 2 active contexts, got %d", lm.ActiveCount())
	}
}

func TestLifecycleManager_CancelContext(t *testing.T) {
	lm := NewLifecycleManager()

	ctx := lm.Register("page1")

	// Cancel context
	lm.CancelContext("page1")

	// Check context is cancelled
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be cancelled")
	}

	// Check active count
	if lm.ActiveCount() != 0 {
		t.Errorf("Expected 0 active contexts, got %d", lm.ActiveCount())
	}
}

func TestLifecycleManager_GetContext(t *testing.T) {
	lm := NewLifecycleManager()

	ctx, exists := lm.GetContext("nonexistent")
	if exists {
		t.Error("Non-existent context should not exist")
	}
	if ctx != nil {
		t.Error("Non-existent context should be nil")
	}

	lm.Register("page1")
	ctx, exists = lm.GetContext("page1")
	if !exists {
		t.Error("Context should exist")
	}
	if ctx == nil {
		t.Error("Context should not be nil")
	}
}

func TestLifecycleManager_CancelAll(t *testing.T) {
	lm := NewLifecycleManager()

	ctx1 := lm.Register("page1")
	ctx2 := lm.Register("page2")

	// Cancel all
	lm.CancelAll()

	// Check all contexts are cancelled
	select {
	case <-ctx1.Done():
		// Expected
	default:
		t.Error("Context 1 should be cancelled")
	}

	select {
	case <-ctx2.Done():
		// Expected
	default:
		t.Error("Context 2 should be cancelled")
	}

	// Check active count
	if lm.ActiveCount() != 0 {
		t.Errorf("Expected 0 active contexts, got %d", lm.ActiveCount())
	}
}

func TestLifecycleManager_ReRegister(t *testing.T) {
	lm := NewLifecycleManager()

	// Register first context
	ctx1 := lm.Register("page1")
	ctx1Done := ctx1.Done()

	// Re-register same ID
	ctx2 := lm.Register("page1")

	// First context should be cancelled
	select {
	case <-ctx1Done:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("First context should be cancelled immediately")
	}

	// Second context should exist
	if ctx2 == nil {
		t.Error("Second context should not be nil")
	}

	// Active count should be 1
	if lm.ActiveCount() != 1 {
		t.Errorf("Expected 1 active context, got %d", lm.ActiveCount())
	}
}

func TestLifecycleManager_ContextUsage(t *testing.T) {
	lm := NewLifecycleManager()

	ctx, exists := lm.GetContext("test")
	if ctx != nil {
		t.Error("Context should not exist initially")
	}
	if exists {
		t.Error("Context should not exist initially")
	}

	// Register context
	lm.Register("test")
	ctx, exists = lm.GetContext("test")

	if ctx == nil {
		t.Error("Context should exist after registration")
	}
	if !exists {
		t.Error("Context should exist after registration")
	}

	// Use context to create a channel
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(done)
	}()

	// Cancel context
	lm.CancelContext("test")

	select {
	case <-done:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled")
	}
}
