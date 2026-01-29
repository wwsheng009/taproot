package image

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/ui/styles"
	"github.com/yourorg/taproot/internal/tui/util"
)

// RendererType represents the image rendering protocol
type RendererType int

const (
	RendererAuto RendererType = iota
	RendererKitty
	RendereriTerm2
	RendererBlocks // Unicode block characters (fallback)
)

// Image displays an image in the terminal
type Image struct {
	width       int
	height      int
	path        string
	renderer    RendererType
	loaded      bool
	error       string
	aspectRatio float64
	styles      *styles.Styles
}

const (
	maxWidth = 100
)

// New creates a new image component
func New(path string) *Image {
	s := styles.DefaultStyles()
	return &Image{
		path:     path,
		renderer: RendererAuto,
		loaded:   false,
		styles:   &s,
	}
}

// SetRenderer sets the rendering type
func (img *Image) SetRenderer(renderer RendererType) tea.Cmd {
	img.renderer = renderer
	return img.Reload()
}

// Reload reloads the image
func (img *Image) Reload() tea.Cmd {
	return func() tea.Msg {
		// Check if file exists
		if _, err := os.Stat(img.path); os.IsNotExist(err) {
			img.error = fmt.Sprintf("File not found: %s", img.path)
			img.loaded = false
			return nil
		}

		// For now, just mark as loaded
		// In a full implementation, this would decode the image
		img.loaded = true
		img.error = ""
		return nil
	}
}

func (img *Image) Init() tea.Cmd {
	return img.Reload()
}

func (img *Image) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return img, img.Reload()
		}
	case tea.WindowSizeMsg:
		img.width = msg.Width
		img.height = msg.Height
	}

	return img, nil
}

func (img *Image) View() string {
	s := img.styles

	if img.error != "" {
		// Show error message
		errorStyle := lipgloss.NewStyle().
			Foreground(s.Error).
			Bold(true).
			Width(img.width).
			Align(lipgloss.Center)

		return errorStyle.Render("⚠️  " + img.error)
	}

	if !img.loaded {
		// Show loading state
		loadingStyle := lipgloss.NewStyle().
			Foreground(s.FgMuted).
			Width(img.width).
			Align(lipgloss.Center)

		return loadingStyle.Render("Loading image...")
	}

	// Detect renderer type
	renderer := img.renderer
	if renderer == RendererAuto {
		renderer = detectRenderer()
	}

	// Render based on type
	switch renderer {
	case RendererKitty:
		return img.renderKitty()
	case RendereriTerm2:
		return img.renderiTerm2()
	default:
		return img.renderBlocks()
	}
}

// renderKitty uses the Kitty graphics protocol
func (img *Image) renderKitty() string {
	// Kitty protocol escape sequence
	// This is a simplified version - full implementation would encode the image
	s := img.styles
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Border).
		Padding(1)

	placeholder := fmt.Sprintf("[Image: %s]\nKitty graphics protocol\n(%dx%d)",
		img.path, img.width, img.height)

	return boxStyle.Render(placeholder)
}

// renderiTerm2 uses the iTerm2 inline image protocol
func (img *Image) renderiTerm2() string {
	s := img.styles
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Border).
		Padding(1)

	placeholder := fmt.Sprintf("[Image: %s]\niTerm2 inline images\n(%dx%d)",
		img.path, img.width, img.height)

	return boxStyle.Render(placeholder)
}

// renderBlocks uses Unicode block characters (fallback)
func (img *Image) renderBlocks() string {
	s := img.styles
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Border).
		Padding(1)

	// Create a simple ASCII art placeholder
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("┌─ %s ─┐\n", truncatePath(img.path, 30)))
	sb.WriteString("│ Image  │\n")
	sb.WriteString("│ Viewer │\n")
	sb.WriteString("│        │\n")
	sb.WriteString("│  ░▒▓█ │\n")
	sb.WriteString("│  █▓▒░ │\n")
	sb.WriteString("│        │\n")
	sb.WriteString("│" + strings.Repeat("─", img.width-4) + "│\n")
	sb.WriteString(fmt.Sprintf("│ %dx%d  │\n", img.width, min(img.height, 20)))
	sb.WriteString("└" + strings.Repeat("─", img.width-4) + "┘")

	rendered := boxStyle.Render(sb.String())

	return lipgloss.NewStyle().
		Width(img.width).
		Align(lipgloss.Center).
		Render(rendered)
}

func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}

// detectRenderer attempts to detect the best available renderer
func detectRenderer() RendererType {
	// Check environment variables
	term := os.Getenv("TERM")
	program := os.Getenv("TERM_PROGRAM")

	// Kitty detection
	if strings.Contains(term, "kitty") || strings.Contains(program, "kitty") {
		return RendererKitty
	}

	// iTerm2 detection
	if program == "iTerm.app" {
		return RendereriTerm2
	}

	// Default to blocks (works everywhere)
	return RendererBlocks
}

// Size returns the current size
func (img *Image) Size() (int, int) {
	return img.width, img.height
}

// SetSize sets the dimensions
func (img *Image) SetSize(w, h int) {
	img.width = w
	img.height = h
}

// Path returns the image path
func (img *Image) Path() string {
	return img.path
}

// SetPath changes the image path and reloads
func (img *Image) SetPath(path string) tea.Cmd {
	img.path = path
	return img.Reload()
}

// IsLoaded returns whether the image is loaded
func (img *Image) IsLoaded() bool {
	return img.loaded
}

// Error returns any loading error
func (img *Image) Error() string {
	return img.error
}
