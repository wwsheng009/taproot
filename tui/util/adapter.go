package util

import (
	tea "github.com/charmbracelet/bubbletea"
)

// BubbleteaToRenderModel creates a render.Model from a util.Model (Bubbletea-based).
// This adapter allows Bubbletea components to work with render-engine-agnostic interfaces.
//
// Example:
//
//	bubbleteaModel := &MyComponent{}
//	renderModel := BubbleteaToRenderModel(bubbleteaModel)
type BubbleteaToRenderModel struct {
	inner Model
}

// NewBubbleteaToRenderModel creates a new adapter.
func NewBubbleteaToRenderModel(inner Model) *BubbleteaToRenderModel {
	return &BubbleteaToRenderModel{inner: inner}
}

// Init initializes the model.
func (a *BubbleteaToRenderModel) Init() tea.Cmd {
	return a.inner.Init()
}

// Update handles messages. Converts any message to tea.Msg if possible.
func (a *BubbleteaToRenderModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	newInner, cmd := a.inner.Update(msg)
	return &BubbleteaToRenderModel{inner: newInner}, cmd
}

// View returns the rendered string.
func (a *BubbleteaToRenderModel) View() string {
	return a.inner.View()
}

// GetInner returns the underlying Bubbletea model.
func (a *BubbleteaToRenderModel) GetInner() Model {
	return a.inner
}

// WithInner updates the inner model and returns the adapter.
func (a *BubbleteaToRenderModel) WithInner(inner Model) *BubbleteaToRenderModel {
	return &BubbleteaToRenderModel{inner: inner}
}

// RenderToBubbleteaModel creates a util.Model from any component with View() string.
// This allows simple renderable components to work with Bubbletea.
type RenderToBubbleteaModel struct {
	inner interface {
		View() string
	}
}

// NewRenderToBubbleteaModel creates a new adapter.
func NewRenderToBubbleteaModel(inner interface {
	View() string
}) *RenderToBubbleteaModel {
	return &RenderToBubbleteaModel{inner: inner}
}

// Init returns nil command.
func (a *RenderToBubbleteaModel) Init() tea.Cmd {
	return nil
}

// Update returns the model itself with no command.
func (a *RenderToBubbleteaModel) Update(msg tea.Msg) (Model, tea.Cmd) {
	return a, nil
}

// View returns the rendered string.
func (a *RenderToBubbleteaModel) View() string {
	return a.inner.View()
}
