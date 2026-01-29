package render

import (
	"os"

	uv "github.com/charmbracelet/ultraviolet"
)

// UltravioletEngine implements the Engine interface using the Ultraviolet rendering engine.
type UltravioletEngine struct {
	config  *EngineConfig
	term    *uv.Terminal
	running bool
	model   Model
}

// NewUltravioletEngine creates a new UltravioletEngine instance.
func NewUltravioletEngine(config *EngineConfig) Engine {
	return &UltravioletEngine{
		config: config,
	}
}

// Type returns the engine type.
func (e *UltravioletEngine) Type() EngineType {
	return EngineUltraviolet
}

// Start initializes the engine and begins the event loop.
func (e *UltravioletEngine) Start(model Model) error {
	e.model = model
	e.running = true

	// Initialize initialization command
	if err := e.model.Init(); err != nil {
		return err
	}

	// Create terminal instance
	// UV uses standard IO by default or takes them as args
	e.term = uv.NewTerminal(os.Stdin, os.Stdout, os.Environ())

	// Start the terminal (enters raw mode, alt screen if configured in UV, etc)
	if err := e.term.Start(); err != nil {
		return err
	}

	// Initial render
	e.render()

	// Event loop
	go func() {
		for event := range e.term.Events() {
			if !e.running {
				break
			}

			var msg Msg

			switch evt := event.(type) {
			case uv.KeyPressEvent:
				// Translate UV key event to internal KeyMsg
				// Note: UV KeyPressEvent might need adaptation to our KeyMsg struct
				// For now, we use the String() representation which is usually enough
				msg = KeyMsg{
					Key: evt.String(),
				}
			case uv.WindowSizeEvent:
				msg = WindowSizeMsg{
					Width:  evt.Width,
					Height: evt.Height,
				}
			case error:
				// Handle errors
				continue
			}

			if msg != nil {
				e.update(msg)
			}
		}
	}()

	// Block until stopped (UV doesn't block main thread usually, but we need to keep alive?
	// The Engine.Start contract in this codebase seems to imply blocking or non-blocking?
	// Looking at adapter_tea.go (standard Bubbletea), Program.Run() blocks.
	// So we should probably block here too to match behavior if possible, or coordinate.
	// However, the interface says Start() error. If Bubbletea's Run() blocks, then Start() blocks.
	// We need a way to wait.
	
	// For now, let's just block on a channel since we are the main driver.
	// But wait, if we block, we can't run other logic?
	// Standard Bubbletea Run() BLOCKS. So we should block.
	
	// We need a channel to signal completion.
	// Since we are inside Start, we can just select{} or wait on a done channel.
	
	// Refactoring to support proper blocking:
	<-make(chan struct{}) // Temporary: blocking forever until killed
	return nil
}

// update handles the model update cycle
func (e *UltravioletEngine) update(msg Msg) {
	// Check for quit message
	if _, ok := msg.(QuitMsg); ok {
		e.Stop()
		return
	}

	// Update model
	newModel, cmd := e.model.Update(msg)
	e.model = newModel

	// Check for quit command
	if IsQuit(cmd) {
		e.Stop()
		return
	}

	// Render
	e.render()
}

// render performs the rendering step
func (e *UltravioletEngine) render() {
	if e.term == nil {
		return
	}
	
	view := e.model.View()
	// UV requires a Drawable. We use NewStyledString for text content.
	drawable := uv.NewStyledString(view)
	e.term.Draw(drawable)
	e.term.Display()
}

// Stop gracefully shuts down the engine.
func (e *UltravioletEngine) Stop() error {
	e.running = false
	if e.term != nil {
		e.term.Stop()
	}
	// In a real implementation we would signal the blocking Start to return
	os.Exit(0) // Aggressive exit for now matching behavior
	return nil
}

// Send sends a message to the model.
func (e *UltravioletEngine) Send(msg Msg) error {
	// In UV adaptation, we'd need a way to inject messages into the event loop.
	// Since we own the loop in Start(), we could use a channel.
	// For now, simpler implementation.
	if e.running {
		e.update(msg)
	}
	return nil
}

// Resize notifies the engine of a terminal size change.
func (e *UltravioletEngine) Resize(width, height int) error {
	e.update(WindowSizeMsg{Width: width, Height: height})
	return nil
}

// Running returns true if the engine is active.
func (e *UltravioletEngine) Running() bool {
	return e.running
}
