package sidebar

import (
	"github.com/wwsheng009/taproot/layout"
)

// Sidebar represents a sidebar component for navigation and information display.
type Sidebar interface {
	layout.Sizeable

	// Init initializes the sidebar.
	Init() error

	// Update handles messages and updates state.
	Update(msg any) (Sidebar, any)

	// View renders the sidebar.
	View() string

	// SetModelInfo sets the model information to display.
	SetModelInfo(info ModelInfo)

	// SetSession sets the session information.
	SetSession(info SessionInfo)

	// SetCompactMode toggles compact mode.
	SetCompactMode(compact bool)

	// AddFile adds a file to the modified files list.
	AddFile(file FileInfo)

	// ClearFiles clears the modified files list.
	ClearFiles()

	// SetLSPStatus sets the LSP service status.
	SetLSPStatus(services []LSPService)

	// SetMCPStatus sets the MCP service status.
	SetMCPStatus(services []MCPService)
}

// ModelInfo represents information about the current AI model.
type ModelInfo struct {
	Name         string  // Model name
	Icon         string  // Icon to display (e.g., "M")
	Provider     string  // Provider name (e.g., "openai", "anthropic")
	CanReason    bool    // Whether model supports reasoning
	ReasoningOn  bool    // Whether reasoning is enabled
	ReasoningEffort string // Reasoning effort level (e.g., "low", "medium", "high")
	ContextWindow int64   // Context window size in tokens
}

// SessionInfo represents session-level information.
type SessionInfo struct {
	ID              string // Session ID
	Title           string // Session title
	PromptTokens    int64  // Prompt tokens used
	CompletionTokens int64 // Completion tokens used
	Cost            float64 // Total cost in dollars
	WorkingDir      string // Current working directory
}

// FileInfo represents a modified file.
type FileInfo struct {
	Path      string // File path (relative to working dir)
	Additions int    // Number of lines added
	Deletions int    // Number of lines deleted
}

// LSPService represents an LSP service status.
type LSPService struct {
	Name       string // Service name
	Language   string // Language (e.g., "go", "python")
	Connected  bool   // Whether connected
	ErrorCount int    // Number of errors from this LSP
}

// MCPService represents an MCP service status.
type MCPService struct {
	Name      string // Service name
	Connected bool   // Whether connected
}

// Config represents sidebar configuration options.
type Config struct {
	// Width of the sidebar (used if not using layout.Area)
	Width int

	// Height of the sidebar (used if not using layout.Area)
	Height int

	// Logo display options
	ShowLogo       bool
	LogoHeight     int  // Minimum height to show full logo
	LogoProvider   LogoProvider

	// Section display options
	MaxFiles       int  // Maximum files to show
	MaxLSPs        int  // Maximum LSP services to show
	MaxMCPs        int  // Maximum MCP services to show

	// Mode
	CompactMode    bool // Enable compact mode
}

// DefaultConfig returns default sidebar configuration.
func DefaultConfig() Config {
	return Config{
		Width:         30,
		Height:        50,
		ShowLogo:      true,
		LogoHeight:    30,
		LogoProvider:  nil, // Will use default ASCII logo
		MaxFiles:      10,
		MaxLSPs:       8,
		MaxMCPs:       8,
		CompactMode:   false,
	}
}

// LogoProvider is a function that renders a logo with the given width.
type LogoProvider func(width int) string

// DefaultLogoProvider returns a simple ASCII logo.
func DefaultLogoProvider(width int) string {
	logo := `
 ┌─────────────────┐
 │  TR  │  TAPROOT │
 └─────────────────┘
`
	// Truncate to fit width
	if width > 0 && width < len(logo) {
		return logo[:width]
	}
	return logo
}
