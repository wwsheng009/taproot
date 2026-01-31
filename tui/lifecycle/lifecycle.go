package lifecycle

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

// Lifecycle defines the lifecycle hooks for components that need cleanup.
type Lifecycle interface {
	// OnMount is called when the component is mounted/becomes active.
	// Returns a command that will be executed.
	OnMount(ctx context.Context) tea.Cmd

	// OnUnmount is called when the component is unmounted/becomes inactive.
	// This is the place to cancel any ongoing operations (goroutines, timers, etc.).
	// Returns a command for any final cleanup.
	OnUnmount() tea.Cmd

	// IsActive returns whether the component is currently active.
	IsActive() bool
}

// ContextProvider provides lifecycle context to components.
type ContextProvider interface {
	// GetContext returns the lifecycle context for a given component ID.
	GetContext(id string) (context.Context, context.CancelFunc)

	// CancelContext cancels the context for a given component ID.
	CancelContext(id string)
}

// LifecycleManager manages lifecycle contexts for components.
type LifecycleManager struct {
	contexts map[string]context.Context
	cancels  map[string]context.CancelFunc
}

// NewLifecycleManager creates a new lifecycle manager.
func NewLifecycleManager() *LifecycleManager {
	return &LifecycleManager{
		contexts: make(map[string]context.Context),
		cancels:  make(map[string]context.CancelFunc),
	}
}

// Register creates a new context for a component.
func (lm *LifecycleManager) Register(id string) context.Context {
	if cancel, exists := lm.cancels[id]; exists {
		cancel() // Cancel old context if exists
	}

	ctx, cancel := context.WithCancel(context.Background())
	lm.contexts[id] = ctx
	lm.cancels[id] = cancel
	return ctx
}

// GetContext returns the context for a component.
func (lm *LifecycleManager) GetContext(id string) (context.Context, bool) {
	ctx, exists := lm.contexts[id]
	return ctx, exists
}

// CancelContext cancels the context for a component.
func (lm *LifecycleManager) CancelContext(id string) {
	if cancel, exists := lm.cancels[id]; exists {
		cancel()
		delete(lm.contexts, id)
		delete(lm.cancels, id)
	}
}

// CancelAll cancels all contexts.
func (lm *LifecycleManager) CancelAll() {
	for id, cancel := range lm.cancels {
		cancel()
		delete(lm.contexts, id)
		delete(lm.cancels, id)
	}
}

// ActiveCount returns the number of active contexts.
func (lm *LifecycleManager) ActiveCount() int {
	return len(lm.contexts)
}
