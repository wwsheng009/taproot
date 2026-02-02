package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/components/image"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Get image path from command line or use a placeholder
	imgPath := "placeholder.png"
	if len(os.Args) > 1 {
		imgPath = os.Args[1]
	}

	model := NewModel(imgPath)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

type Model struct {
	image     *image.Image
	quitting  bool
	imgPath   string
	renderer  image.RendererType
	width     int
	height    int
}

func NewModel(imgPath string) Model {
	img := image.New(imgPath)

	return Model{
		image:    img,
		quitting: false,
		imgPath:  imgPath,
		renderer: image.RendererAuto,
	}
}

func (m Model) Init() tea.Cmd {
	return m.image.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// First, update the image component with the message
	var imgCmd tea.Cmd
	updatedModel, imgCmd := m.image.Update(msg)
	m.image = updatedModel.(*image.Image)

	// Then handle our own messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "r":
			// Reload image
			return m, m.image.Reload()
		case "1":
			// Auto detect renderer
			m.renderer = image.RendererAuto
			return m, m.image.SetRenderer(image.RendererAuto)
		case "2":
			// Kitty renderer
			m.renderer = image.RendererKitty
			return m, m.image.SetRenderer(image.RendererKitty)
		case "3":
			// iTerm2 renderer
			m.renderer = image.RendereriTerm2
			return m, m.image.SetRenderer(image.RendereriTerm2)
		case "4":
			// Block renderer
			m.renderer = image.RendererBlocks
			return m, m.image.SetRenderer(image.RendererBlocks)
		case "s":
			// Switch to a different path (demo)
			m.imgPath = "demo-" + m.imgPath
			return m, m.image.SetPath(m.imgPath)
		case "+", "=":
			// Zoom in
			return m, m.image.ZoomIn()
		case "-", "_":
			// Zoom out
			return m, m.image.ZoomOut()
		case "0", " ":
			// Reset zoom to fit screen
			return m, m.image.ResetZoom()
		case "m":
			// Cycle zoom mode
			return m, m.image.CycleZoomMode()
		}
	case tea.WindowSizeMsg:
		// Check if this is the first time we get window size (or if it changed)
		oldWidth := m.width
		m.width = msg.Width
		m.height = msg.Height
		
		// If we just got window size for the first time, or it changed significantly
		if oldWidth == 0 || oldWidth != msg.Width {
			newW := msg.Width - 4
			newH := msg.Height - 4
			m.image.SetSize(newW, newH)
		}
	}

	// Return any command from the image component
	return m, imgCmd
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	theme := lipgloss.NewStyle()
	title := theme.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Image Component Demo")

	var b strings.Builder

	// Header
	b.WriteString(title + "\n\n")

	// Image info
	rendererName := "Auto"
	switch m.renderer {
	case image.RendererKitty:
		rendererName = "Kitty"
	case image.RendereriTerm2:
		rendererName = "iTerm2"
	case image.RendererBlocks:
		rendererName = "Blocks (Unicode)"
	}

	fmt.Fprintf(&b, "Path: %s\n", m.imgPath)
	fmt.Fprintf(&b, "Renderer: %s\n", rendererName)
	fmt.Fprintf(&b, "Loaded: %v\n", m.image.IsLoaded())

	// Add debugging info
	imgW, imgH := m.image.Size()
	scaledW, scaledH := m.image.ScaledSize()
	fmt.Fprintf(&b, "Terminal: %dx%d | Image: %dx%d | Scaled: %dx%d\n\n", m.width, m.height, imgW, imgH, scaledW, scaledH)

	if m.image.Error() != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		b.WriteString(errorStyle.Render("Error: " + m.image.Error()))
		b.WriteString("\n\n")
		b.WriteString("Tip: Provide an image path as argument:\n")
		b.WriteString("  go run examples/image/main.go /path/to/image.png\n\n")
	}

	// Image view
	b.WriteString(m.image.View())

	// Add separator line
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	b.WriteString("\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", min(m.width, 80))))
	b.WriteString("\n")

	// Footer hints
	b.WriteString("\n")
	
	// Show zoom level and mode
	scalePercent := int(m.image.GetScale() * 100)
	zoomMode := m.image.GetZoomModeName()
	zoomText := fmt.Sprintf("Zoom: %s %d%%", zoomMode, scalePercent)
	
	hints := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(
		"+/-: Zoom | 0: Reset | m: Mode | 1-4: Renderer | r: Reload | s: Path | q: Quit  [" + zoomText + "]",
	)
	b.WriteString(hints)

	return b.String()
}
