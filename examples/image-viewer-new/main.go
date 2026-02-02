package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/components/image"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

type model struct {
	img      *image.Image
	quitting bool
	width    int
	height   int
	showInfo bool
	showHelp bool // Toggle help display
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	imgPath := os.Args[1]

	// Create image component
	img := image.New(imgPath)

	m := model{
		img:      img,
		quitting: false,
		showInfo: true,
		showHelp: true,
	}

	// Use the adapter wrapper
	adapter := &imageAdapter{model: m}

	p := tea.NewProgram(adapter, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run main.go <image-path>")
	fmt.Println("\nSupported formats: JPEG, PNG, GIF, BMP")
	fmt.Println("\nKey Controls:")
	fmt.Println("  Zoom:")
	fmt.Println("    +/=     Zoom in")
	fmt.Println("    -/_     Zoom out")
	fmt.Println("    0       Reset zoom to 100%")
	fmt.Println("    *       Zoom to 200%")
	fmt.Println("    %       Zoom to 50%")
	fmt.Println("    [/]     Fine zoom (1% steps)")
	fmt.Println("\n  Display Mode:")
	fmt.Println("    m       Cycle zoom mode (Fit→Fill→Stretch→Original)")
	fmt.Println("    f       Fit mode (maintain aspect)")
	fmt.Println("    F       Fill mode (may crop)")
	fmt.Println("    s       Stretch mode (ignore aspect)")
	fmt.Println("    o       Original size (1:1 pixels)")
	fmt.Println("\n  Renderer:")
	fmt.Println("    1       Auto-detect best renderer")
	fmt.Println("    2       Kitty graphics protocol")
	fmt.Println("    3       iTerm2 inline images")
	fmt.Println("    4       Unicode blocks (works everywhere)")
	fmt.Println("    5       Sixel graphics")
	fmt.Println("    6       ASCII art (max compatibility)")
	fmt.Println("\n  Other:")
	fmt.Println("    r       Reload image")
	fmt.Println("    i       Toggle info")
	fmt.Println("    h       Toggle help")
	fmt.Println("    q/ctrl+c Quit")
	fmt.Println("\nExample:")
	fmt.Println("  go run main.go photo.jpg")
}

// imageAdapter wraps our model to implement tea.Model
type imageAdapter struct {
	model model
}

func (a *imageAdapter) Init() tea.Cmd {
	return adaptCmd(a.model.img.Init())
}

func (a *imageAdapter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle our model's quitting state
	if a.model.quitting {
		return a, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			a.model.quitting = true
			return a, tea.Quit

		case "h", "?":
			a.model.showHelp = !a.model.showHelp
			return a, nil

		case "i":
			a.model.showInfo = !a.model.showInfo
			return a, nil

		// Zoom controls
		case "+", "=":
			a.model.img.ZoomIn()
			return a, nil
		case "-", "_":
			a.model.img.ZoomOut()
			return a, nil
		case "0":
			a.model.img.ResetZoom()
			return a, nil
		case "*":
			a.model.img.SetScale(2.0)
			return a, nil
		case "%":
			a.model.img.SetScale(0.5)
			return a, nil

		// Zoom mode switching
		case "m":
			a.model.img.CycleZoomMode()
			return a, nil
		case "f":
			a.model.img.SetZoomMode(image.ZoomFit)
			return a, nil
		case "F":
			a.model.img.SetZoomMode(image.ZoomFill)
			return a, nil
		case "s":
			a.model.img.SetZoomMode(image.ZoomStretch)
			return a, nil
		case "o":
			a.model.img.SetZoomMode(image.ZoomOriginal)
			return a, nil

		// Fine zoom control
		case "[":
			scale := a.model.img.GetScale()
			scale -= 0.01
			if scale < 0.1 {
				scale = 0.1
			}
			a.model.img.SetScale(scale)
			return a, nil
		case "]":
			scale := a.model.img.GetScale()
			scale += 0.01
			if scale > 5.0 {
				scale = 5.0
			}
			a.model.img.SetScale(scale)
			return a, nil

		// Renderer selection
		case "1":
			cmd := a.model.img.SetRenderer(image.RendererAuto)
			return a, adaptCmd(cmd)
		case "2":
			cmd := a.model.img.SetRenderer(image.RendererKitty)
			return a, adaptCmd(cmd)
		case "3":
			cmd := a.model.img.SetRenderer(image.RendereriTerm2)
			return a, adaptCmd(cmd)
		case "4":
			cmd := a.model.img.SetRenderer(image.RendererBlocks)
			return a, adaptCmd(cmd)
		case "5":
			cmd := a.model.img.SetRenderer(image.RendererSixel)
			return a, adaptCmd(cmd)
		case "6":
			cmd := a.model.img.SetRenderer(image.RendererASCII)
			return a, adaptCmd(cmd)

		// Reload
		case "r":
			cmd := a.model.img.Reload()
			return a, adaptCmd(cmd)
		}

	case tea.WindowSizeMsg:
		a.model.width = msg.Width
		a.model.height = msg.Height
		a.model.img.SetSize(msg.Width, msg.Height)
	}

	// Update image component - convert tea.Msg to render.Msg
	newImg, cmd := a.model.img.Update(convertMsg(msg))
	a.model.img = newImg.(*image.Image)

	return a, adaptCmd(cmd)
}

func (a *imageAdapter) View() string {
	return a.model.View()
}

// adaptCmd converts render.Cmd to tea.Cmd
func adaptCmd(cmd render.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}

	// Check for quit command
	if q, ok := cmd.(interface{ IsQuit() bool }); ok && q.IsQuit() {
		return tea.Quit
	}

	// If it's a func() render.Msg, wrap it
	if fn, ok := cmd.(func() render.Msg); ok {
		return func() tea.Msg {
			return fn()
		}
	}

	// If it's a render.Command, execute it
	if c, ok := cmd.(render.Command); ok {
		return func() tea.Msg {
			err := c.Execute()
			if err != nil {
				return render.ErrorMsg{Error: err}
			}
			return nil
		}
	}

	return nil
}

// convertMsg converts tea.Msg to render.Msg
func convertMsg(msg tea.Msg) any {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return render.KeyMsg{Key: msg.String()}
	case tea.WindowSizeMsg:
		return render.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height,
		}
	default:
		return msg
	}
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	s := styles.DefaultStyles()

	// Fixed header lines
	headerLines := 1
	infoLines := 1
	footerLines := 2 // Help + Renderer info

	// Calculate available height for image
	availableHeight := m.height - headerLines - infoLines - footerLines
	if availableHeight < 1 {
		availableHeight = 1
	}

	// Create layout
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(s.Primary).
		Bold(true).
		Padding(0, 1)

	b.WriteString(headerStyle.Render("🖼️  Taproot Image Viewer"))
	b.WriteString("\n")

	// Info bar
	if m.showInfo {
		infoStyle := lipgloss.NewStyle().
			Foreground(s.FgMuted).
			Width(m.width).
			Padding(0, 1)

		var info []string

		// Get image info
		imgW, imgH := m.img.GetImageDimensions()
		if imgW > 0 && imgH > 0 {
			info = append(info, fmt.Sprintf("%dx%d", imgW, imgH))
		}

		// Get renderer info
		renderer := m.img.GetRenderer()
		info = append(info, fmt.Sprintf("Renderer: %s", renderer.String()))

		// Get zoom info
		scale := m.img.GetScale()
		zoomMode := m.img.GetZoomModeName()
		info = append(info, fmt.Sprintf("Zoom: %s %.0f%%", zoomMode, scale*100))

		// Get scaled size
		scaledW, scaledH := m.img.ScaledSize()
		if scaledW > 0 && scaledH > 0 {
			info = append(info, fmt.Sprintf("Display: %dx%d", scaledW, scaledH))
		}

		// Get error if any
		if err := m.img.Error(); err != "" {
			info = append(info, fmt.Sprintf("Error: %s", err))
		}

		b.WriteString(infoStyle.Render(strings.Join(info, " • ")))
		b.WriteString("\n")
	}

	// Image display area - limit height to keep footer at bottom
	imageView := m.img.View()
	imageLines := strings.Split(imageView, "\n")

	// Display image up to available height, track displayed count
	displayedLines := 0
	for i, line := range imageLines {
		if i >= availableHeight {
			break
		}
		b.WriteString(line)
		b.WriteString("\n")
		displayedLines++
	}

	// Add padding lines to push footer to bottom
	remainingPadding := availableHeight - displayedLines
	if remainingPadding > 0 {
		for i := 0; i < remainingPadding; i++ {
			b.WriteString("\n")
		}
	}

	// Help/Controls (always at bottom)
	footerStyle := lipgloss.NewStyle().
		Foreground(s.FgSubtle).
		Padding(0, 1)

	if m.showHelp {
		controls := `Zoom: +/-/0/* | Fine: []/[%] | Mode: m/f/F/s/o | r:Reload | i:Info | h:Help | q:Quit`
		b.WriteString(footerStyle.Render(controls))
	} else {
		controls := fmt.Sprintf("Zoom: %.0f%% | Mode: %s | h:Help | q:Quit",
			m.img.GetScale()*100, m.img.GetZoomModeName())
		b.WriteString(footerStyle.Render(controls))
	}
	b.WriteString("\n")

	// Renderer info (always at bottom)
	rendererInfo := fmt.Sprintf("Renderer: %s | %s",
		m.img.GetRenderer().String(),
		image.GetRendererDescription(m.img.GetRenderer()))
	b.WriteString(footerStyle.Render(rendererInfo))

	return b.String()
}
