package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/components/image"
	"github.com/wwsheng009/taproot/ui/components/layout"
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
	image    *image.Image
	quitting bool
	imgPath  string
	renderer image.RendererType
	width    int
	height   int
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
	// Handle our own messages
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "r":
			// Reload image
			cmds = append(cmds, m.image.Reload())
		case "1":
			// Auto detect renderer
			m.renderer = image.RendererAuto
			cmds = append(cmds, m.image.SetRenderer(image.RendererAuto))
		case "2":
			// Kitty renderer
			m.renderer = image.RendererKitty
			cmds = append(cmds, m.image.SetRenderer(image.RendererKitty))
		case "3":
			// iTerm2 renderer
			m.renderer = image.RendereriTerm2
			cmds = append(cmds, m.image.SetRenderer(image.RendereriTerm2))
		case "4":
			// Block renderer
			m.renderer = image.RendererBlocks
			cmds = append(cmds, m.image.SetRenderer(image.RendererBlocks))
		case "5":
			// Sixel renderer (high-quality)
			m.renderer = image.RendererSixel
			cmds = append(cmds, m.image.SetRenderer(image.RendererSixel))
		case "6":
			// ASCII renderer (fallback)
			m.renderer = image.RendererASCII
			cmds = append(cmds, m.image.SetRenderer(image.RendererASCII))
		case "s":
			// Switch to a different path (demo)
			m.imgPath = "demo-" + m.imgPath
			cmds = append(cmds, m.image.SetPath(m.imgPath))
		case "+", "=":
			// Zoom in
			cmd := m.image.ZoomIn()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case "-", "_":
			// Zoom out
			cmd := m.image.ZoomOut()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case "0", " ":
			// Reset zoom to fit screen
			cmd := m.image.ResetZoom()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case "m":
			// Cycle zoom mode
			cmd := m.image.CycleZoomMode()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case tea.WindowSizeMsg:
		// Check if this is the first time we get window size (or if it changed)
		oldWidth := m.width
		oldHeight := m.height
		m.width = msg.Width
		m.height = msg.Height

		// Update image size if window size changed (or first time getting size)
		if oldWidth == 0 || oldHeight == 0 || oldWidth != msg.Width || oldHeight != msg.Height {
			newW := msg.Width - 4
			newH := msg.Height - 4
			m.image.SetSize(newW, newH)
		}
	}

	// Always update the image component with the message
	updatedImage, imgCmd := m.image.Update(msg)
	m.image = updatedImage.(*image.Image)

	// Add image component's command if any
	if imgCmd != nil {
		cmds = append(cmds, imgCmd)
	}

	// Return all commands batched together
	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	theme := lipgloss.NewStyle()
	title := theme.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Image Component Demo")

	// Build header
	var header strings.Builder
	header.WriteString(title + "\n\n")

	// Image info
	rendererName := "Auto"
	switch m.renderer {
	case image.RendererSixel:
		rendererName = "Sixel (High-Quality)"
	case image.RendererKitty:
		rendererName = "Kitty"
	case image.RendereriTerm2:
		rendererName = "iTerm2"
	case image.RendererBlocks:
		rendererName = "Blocks (Unicode)"
	case image.RendererASCII:
		rendererName = "ASCII (Fallback)"
	}

	fmt.Fprintf(&header, "Path: %s\n", m.imgPath)
	fmt.Fprintf(&header, "Renderer: %s\n", rendererName)
	fmt.Fprintf(&header, "Loaded: %v\n", m.image.IsLoaded())

	// Add debugging info only if we have window size
	if m.width > 0 && m.height > 0 {
		imgW, imgH := m.image.Size()
		scaledW, scaledH := m.image.ScaledSize()
		fmt.Fprintf(&header, "Terminal: %dx%d | Image: %dx%d | Scaled: %dx%d\n", m.width, m.height, imgW, imgH, scaledW, scaledH)
	}

	if m.image.Error() != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		header.WriteString(errorStyle.Render("Error: " + m.image.Error()))
		header.WriteString("\n\n")
		header.WriteString("Tip: Provide an image path as argument:\n")
		header.WriteString("  go run examples/image/main.go /path/to/image.png\n")
	}

	header.WriteString("\n")

	// Build footer
	// Show zoom level and mode
	scalePercent := int(m.image.GetScale() * 100)
	zoomMode := m.image.GetZoomModeName()
	zoomText := fmt.Sprintf("Zoom: %s %d%%", zoomMode, scalePercent)

	hints := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(
		"+/-: Zoom | 0: Reset | m: Mode | 1-6: Renderer | r: Reload | s: Path | q: Quit  [" + zoomText + "]",
	)
	var footer strings.Builder
	footer.WriteString(hints)
	footer.WriteString("\n")

	// Get image view
	imageView := m.image.View()

	// For Sixel renderer, calculate actual display height
	displayHeight := 0
	if m.renderer == image.RendererSixel || m.renderer == image.RendererAuto {
		_, scaledH := m.image.ScaledSize()
		// Sixel renders 6 pixels per line (standard Sixel resolution)
		sixelHeight := scaledH / 6
		if sixelHeight > 1 {
			displayHeight = sixelHeight
		}
	}

	// Use vertical layout component
	vbox := layout.NewVerticalLayout().
		SetSize(m.width, m.height).
		SetHeader(header.String()).
		SetContent(imageView).
		SetFooter(footer.String()).
		SetCenterV(true).
		SetCenterH(false).
		SetSeparator(true)

	// Render the layout
	return vbox.Render(displayHeight)
}
