package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/components/image"
)

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
		}
	case tea.WindowSizeMsg:
		m.image.SetSize(msg.Width-4, msg.Height-8)
	}

	return m, nil
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
	fmt.Fprintf(&b, "Loaded: %v\n\n", m.image.IsLoaded())

	if m.image.Error() != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		b.WriteString(errorStyle.Render("Error: " + m.image.Error()))
		b.WriteString("\n\n")
		b.WriteString("Tip: Provide an image path as argument:\n")
		b.WriteString("  go run examples/image/main.go /path/to/image.png\n\n")
	}

	// Image view
	b.WriteString(m.image.View())

	// Footer hints
	b.WriteString("\n\n")
	hints := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(
		"1-4: Switch renderer | r: Reload | s: Switch path | q: Quit",
	)
	b.WriteString(hints)

	return b.String()
}
