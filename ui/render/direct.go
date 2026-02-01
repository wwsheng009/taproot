package render

import (
	"fmt"
	"strings"
	"sync"
)

// DirectEngine is a simple rendering engine that writes directly to output.
// It's primarily useful for testing and simple applications.
type DirectEngine struct {
	config    *EngineConfig
	model     Model
	running   bool
	mu        sync.RWMutex
	output    *stringWriter
	initFunc  func() error
	cleanupFn func() error
}

// stringWriter is a simple writer interface.
type stringWriter struct {
	buf strings.Builder
}

func newStringWriter() *stringWriter {
	return &stringWriter{}
}

func (w *stringWriter) Write(s string) {
	w.buf.WriteString(s)
}

func (w *stringWriter) String() string {
	return w.buf.String()
}

func (w *stringWriter) Clear() {
	w.buf.Reset()
}

// NewDirectEngine creates a new direct engine.
func NewDirectEngine(config *EngineConfig) Engine {
	if config == nil {
		config = DefaultConfig()
	}
	return &DirectEngine{
		config:  config,
		running: false,
		output:  newStringWriter(),
	}
}

// Type returns EngineDirect.
func (e *DirectEngine) Type() EngineType {
	return EngineDirect
}

// Start initializes and runs the engine.
func (e *DirectEngine) Start(model Model) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return fmt.Errorf("engine already running")
	}

	e.model = model

	// Initialize the model
	if err := model.Init(); err != nil {
		return fmt.Errorf("model init failed: %w", err)
	}

	// Call custom init function if provided
	if e.initFunc != nil {
		if err := e.initFunc(); err != nil {
			return err
		}
	}

	e.running = true

	// Initial render
	e.render()

	return nil
}

// Stop shuts down the engine.
func (e *DirectEngine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.running {
		return nil
	}

	// Call custom cleanup function if provided
	if e.cleanupFn != nil {
		if err := e.cleanupFn(); err != nil {
			return err
		}
	}

	e.running = false
	e.model = nil

	return nil
}

// Send sends a message to the model.
func (e *DirectEngine) Send(msg Msg) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.running || e.model == nil {
		return fmt.Errorf("engine not running")
	}

	// Update the model with the message
	newModel, cmd := e.model.Update(msg)
	if newModel != nil {
		e.model = newModel
	}

	// Execute command if returned
	if cmd != nil {
		if c, ok := cmd.(Command); ok {
			if err := c.Execute(); err != nil {
				return err
			}
		}
	}

	// Re-render after update
	e.render()

	return nil
}

// Resize notifies the engine of a size change.
func (e *DirectEngine) Resize(width, height int) error {
	return e.Send(WindowSizeMsg{Width: width, Height: height})
}

// Running returns true if the engine is active.
func (e *DirectEngine) Running() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

// render renders the current model view.
func (e *DirectEngine) render() {
	if e.model == nil {
		return
	}

	view := e.model.View()
	e.output.Clear()
	e.output.Write(view)
}

// Output returns the rendered output.
func (e *DirectEngine) Output() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.output.String()
}

// SetInitFunc sets a custom initialization function.
func (e *DirectEngine) SetInitFunc(fn func() error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.initFunc = fn
}

// SetCleanupFunc sets a custom cleanup function.
func (e *DirectEngine) SetCleanupFunc(fn func() error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cleanupFn = fn
}

// Model returns the current model.
func (e *DirectEngine) Model() Model {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.model
}

// TestModel is a simple test model implementation.
type TestModel struct {
	content string
	width   int
	height  int
}

// NewTestModel creates a new test model.
func NewTestModel(content string) *TestModel {
	return &TestModel{
		content: content,
		width:   80,
		height:  24,
	}
}

// Init initializes the model.
func (m *TestModel) Init() error {
	return nil
}

// Update handles messages.
func (m *TestModel) Update(msg any) (Model, Cmd) {
	switch msg := msg.(type) {
	case WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case KeyMsg:
		if msg.Key == "q" {
			m.content = "quit"
		}
	}
	return m, None()
}

// View renders the model.
func (m *TestModel) View() string {
	return fmt.Sprintf("%s [%dx%d]", m.content, m.width, m.height)
}

// SetContent sets the content.
func (m *TestModel) SetContent(content string) {
	m.content = content
}

// Content returns the content.
func (m *TestModel) Content() string {
	return m.content
}

// Size returns the size.
func (m *TestModel) Size() (int, int) {
	return m.width, m.height
}
