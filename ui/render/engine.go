package render

// EngineType identifies the rendering engine.
type EngineType int

const (
	// EngineBubbletea uses the Bubbletea TUI framework.
	EngineBubbletea EngineType = iota
	// EngineUltraviolet uses the Ultraviolet rendering engine (future).
	EngineUltraviolet
	// EngineDirect uses direct terminal output (for testing/simple cases).
	EngineDirect
)

// String returns the engine type name.
func (e EngineType) String() string {
	switch e {
	case EngineBubbletea:
		return "bubbletea"
	case EngineUltraviolet:
		return "ultraviolet"
	case EngineDirect:
		return "direct"
	default:
		return "unknown"
	}
}

// Engine is the abstraction layer for different rendering engines.
type Engine interface {
	// Type returns the engine type.
	Type() EngineType
	// Start initializes the engine and begins the event loop.
	Start(model Model) error
	// Stop gracefully shuts down the engine.
	Stop() error
	// Send sends a message to the model.
	Send(msg Msg) error
	// Resize notifies the engine of a terminal size change.
	Resize(width, height int) error
	// Running returns true if the engine is active.
	Running() bool
}

// EngineConfig holds configuration for engine initialization.
type EngineConfig struct {
	// EnableMouse enables mouse event support.
	EnableMouse bool
	// EnableAltScreen enables alternate screen mode.
	EnableAltScreen bool
	// EnableCursor enables cursor visibility.
	EnableCursor bool
	// Input is the input reader (nil uses stdin).
	Input any
	// Output is the output writer (nil uses stdout).
	Output any
}

// DefaultConfig returns the default engine configuration.
func DefaultConfig() *EngineConfig {
	return &EngineConfig{
		EnableMouse:     false,
		EnableAltScreen: true,
		EnableCursor:    false,
		Input:           nil,
		Output:          nil,
	}
}

// EngineFactory creates a new engine instance.
type EngineFactory func(config *EngineConfig) Engine

// registry holds registered engine factories.
var registry = make(map[EngineType]EngineFactory)

// RegisterEngine registers an engine factory for a given type.
func RegisterEngine(engineType EngineType, factory EngineFactory) {
	registry[engineType] = factory
}

// init registers the built-in engines
func init() {
	RegisterEngine(EngineBubbletea, NewBubbleteaEngine)
	RegisterEngine(EngineUltraviolet, NewUltravioletEngine)
	RegisterEngine(EngineDirect, NewDirectEngine)
}

// CreateEngine creates a new engine of the specified type.
func CreateEngine(engineType EngineType, config *EngineConfig) (Engine, error) {
	factory, ok := registry[engineType]
	if !ok {
		return nil, ErrEngineNotRegistered{Type: engineType}
	}
	if config == nil {
		config = DefaultConfig()
	}
	return factory(config), nil
}

// ErrEngineNotRegistered is returned when an engine type is not registered.
type ErrEngineNotRegistered struct {
	Type EngineType
}

func (e ErrEngineNotRegistered) Error() string {
	return "engine not registered: " + e.Type.String()
}

// AvailableEngines returns a list of registered engine types.
func AvailableEngines() []EngineType {
	types := make([]EngineType, 0, len(registry))
	for t := range registry {
		types = append(types, t)
	}
	return types
}

// IsEngineRegistered returns true if the engine type is registered.
func IsEngineRegistered(engineType EngineType) bool {
	_, ok := registry[engineType]
	return ok
}
