package render

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"
)

// BubbleteaEngine implements Engine using the Bubbletea framework.
type BubbleteaEngine struct {
	program *tea.Program
	running bool
	config  *EngineConfig
}

// NewBubbleteaEngine creates a new Bubbletea engine.
func NewBubbleteaEngine(config *EngineConfig) Engine {
	return &BubbleteaEngine{
		config: config,
	}
}

// Type returns the engine type.
func (e *BubbleteaEngine) Type() EngineType {
	return EngineBubbletea
}

// Start initializes the engine and begins the event loop.
func (e *BubbleteaEngine) Start(model Model) error {
	opts := []tea.ProgramOption{}

	if e.config.EnableMouse {
		opts = append(opts, tea.WithMouseCellMotion())
	}
	if e.config.EnableAltScreen {
		opts = append(opts, tea.WithAltScreen())
	}
	// Note: tea.WithoutCursor() is not available in the vendored version
	// if !e.config.EnableCursor {
	// 	opts = append(opts, tea.WithoutCursor())
	// }

	// Input/Output
	if e.config.Input != nil {
		if r, ok := e.config.Input.(io.Reader); ok {
			opts = append(opts, tea.WithInput(r))
		}
	}
	if e.config.Output != nil {
		if w, ok := e.config.Output.(io.Writer); ok {
			opts = append(opts, tea.WithOutput(w))
		}
	}

	// Wrap model
	teaModel := &teaAdapter{internal: model}

	e.program = tea.NewProgram(teaModel, opts...)
	e.running = true

	_, err := e.program.Run()
	e.running = false
	return err
}

// Stop gracefully shuts down the engine.
func (e *BubbleteaEngine) Stop() error {
	if e.program != nil {
		e.program.Quit()
	}
	return nil
}

// Send sends a message to the model.
func (e *BubbleteaEngine) Send(msg Msg) error {
	if e.program != nil {
		// If msg is already a tea.Msg, send it directly.
		// Otherwise, wrap it or send as is (tea.Msg is interface{}).
		e.program.Send(msg)
	}
	return nil
}

// Resize notifies the engine of a terminal size change.
func (e *BubbleteaEngine) Resize(width, height int) error {
	if e.program != nil {
		e.program.Send(tea.WindowSizeMsg{Width: width, Height: height})
	}
	return nil
}

// Running returns true if the engine is active.
func (e *BubbleteaEngine) Running() bool {
	return e.running
}

// teaAdapter wraps our Model to satisfy tea.Model.
type teaAdapter struct {
	internal Model
}

func (m *teaAdapter) Init() tea.Cmd {
	cmd := m.internal.Init()
	return adaptCmd(cmd)
}

func (m *teaAdapter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var internalMsg Msg = msg

	// Convert tea.KeyMsg to render.KeyMsg
	if key, ok := msg.(tea.KeyMsg); ok {
		internalMsg = KeyMsg{
			Key: key.String(),
			// Note: We use the string representation for matching,
			// similar to how internal/ui/list works.
			// Detailed modifier bools are not populated here to avoid
			// imperfect parsing, as Key.String() is the source of truth.
		}
	}
	// Convert tea.WindowSizeMsg to render.WindowSizeMsg
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		internalMsg = WindowSizeMsg{
			Width:  size.Width,
			Height: size.Height,
		}
	}

	newModel, cmd := m.internal.Update(internalMsg)
	m.internal = newModel
	return m, adaptCmd(cmd)
}

func (m *teaAdapter) View() string {
	return m.internal.View()
}

// adaptCmd converts our Cmd interface to tea.Cmd.
func adaptCmd(cmd Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}

	// If it's a Command func (func() error), wrap it
	if c, ok := cmd.(Command); ok {
		return func() tea.Msg {
			err := c.Execute()
			if err != nil {
				return ErrorMsg{Error: err}
			}
			return nil
		}
	}

	// If it's already a tea.Cmd (func() tea.Msg), cast it
	if c, ok := cmd.(tea.Cmd); ok {
		return c
	}
	
	// Check if it matches func() tea.Msg signature but not typed as tea.Cmd
	if c, ok := cmd.(func() tea.Msg); ok {
		return c
	}

	// If it's a func() Msg, wrap it
	if c, ok := cmd.(func() Msg); ok {
		return func() tea.Msg {
			return c()
		}
	}

	return nil
}
