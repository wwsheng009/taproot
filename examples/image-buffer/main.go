package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render/buffer"
)

var (
	termSupportsColor = strings.Contains(os.Getenv("TERM"), "color") ||
		strings.Contains(os.Getenv("TERM"), "xterm") ||
		strings.Contains(os.Getenv("TERM"), "screen")
)

func removeANSICodes(s string) string {
	ansiCode := regexp.MustCompile("\x1b\\[[0-9;]*m")
	return ansiCode.ReplaceAllString(s, s)
}

func main() {
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

type Header struct {
	imgPath     string
	renderer    RendererType
	imageLoaded bool
	width       int
	height      int
	imageW      int
	imageH      int
	scaledW     int
	scaledH     int
	imageError  string
}

func NewHeader(imgPath string) *Header {
	return &Header{
		imgPath:  imgPath,
		renderer: RendererAuto,
	}
}

func (h *Header) SetRenderer(rt RendererType) {
	h.renderer = rt
}

func (h *Header) SetPath(path string) {
	h.imgPath = path
}

func (h *Header) SetLoaded(loaded bool) {
	h.imageLoaded = loaded
}

func (h *Header) SetError(err string) {
	h.imageError = err
}

func (h *Header) SetImageSize(imgW, imgH, scaledW, scaledH int) {
	h.imageW = imgW
	h.imageH = imgH
	h.scaledW = scaledW
	h.scaledH = scaledH
}

func (h *Header) SetWindowSize(w, height int) {
	h.width = w
	h.height = height
}

func (h *Header) MinSize() (int, int) {
	return 60, 5
}

func (h *Header) PreferredSize() (int, int) {
	if h.width > 0 {
		return h.width, 6
	}
	return 80, 6
}

func (h *Header) Render(buf *buffer.Buffer, rect buffer.Rect) {
	rendererName := "Auto"
	switch h.renderer {
	case RendererSixel:
		rendererName = "Sixel (High-Quality)"
	case RendererKitty:
		rendererName = "Kitty"
	case RendereriTerm2:
		rendererName = "iTerm2"
	case RendererBlocks:
		rendererName = "Blocks (Unicode)"
	case RendererASCII:
		rendererName = "ASCII (Fallback)"
	}

	buf.WriteString(buffer.Point{X: 2, Y: 0}, "Taproot Image Demo (Buffer Layout)", buffer.Style{Foreground: "#86", Bold: true})

	y := 1
	buf.WriteString(buffer.Point{X: 2, Y: y}, fmt.Sprintf("Path: %s", h.imgPath), buffer.Style{})
	y++
	buf.WriteString(buffer.Point{X: 2, Y: y}, fmt.Sprintf("Renderer: %s", rendererName), buffer.Style{})
	y++
	buf.WriteString(buffer.Point{X: 2, Y: y}, fmt.Sprintf("Loaded: %v", h.imageLoaded), buffer.Style{})
	y++

	if h.width > 0 && h.height > 0 {
		buf.WriteString(buffer.Point{X: 2, Y: y}, fmt.Sprintf("Terminal: %dx%d | Image: %dx%d | Scaled: %dx%d",
			h.width, h.height, h.imageW, h.imageH, h.scaledW, h.scaledH), buffer.Style{})
	}

	if h.imageError != "" {
		buf.WriteString(buffer.Point{X: 2, Y: rect.Height-2}, "Error: "+h.imageError, buffer.Style{Foreground: "#196", Bold: true})
		buf.WriteString(buffer.Point{X: 2, Y: rect.Height-1}, "Tip: Provide an image path as argument:", buffer.Style{Foreground: "#196"})
	}
}

type Footer struct {
	zoomPercent int
	zoomMode    string
	renderer    RendererType
}

func NewFooter() *Footer {
	return &Footer{
		zoomPercent: 100,
		zoomMode:    "Fit",
		renderer:    RendererAuto,
	}
}

func (f *Footer) SetZoom(percent int, mode string) {
	f.zoomPercent = percent
	f.zoomMode = mode
}

func (f *Footer) SetRenderer(rt RendererType) {
	f.renderer = rt
}

func (f *Footer) MinSize() (int, int) {
	return 80, 1
}

func (f *Footer) PreferredSize() (int, int) {
	hints := fmt.Sprintf("+/-: Zoom | 0: Reset | m: Mode | 1-6: Renderer | r: Reload | s: Path | q: Quit  [Zoom: %s %d%%]", f.zoomMode, f.zoomPercent)
	return len(hints) + 2, 1
}

func (f *Footer) Render(buf *buffer.Buffer, rect buffer.Rect) {
	hints := fmt.Sprintf("+/-: Zoom | 0: Reset | m: Mode | 1-6: Renderer | r: Reload | s: Path | q: Quit  [Zoom: %s %d%%]", f.zoomMode, f.zoomPercent)
	buf.WriteString(buffer.Point{X: 1, Y: 0}, hints, buffer.Style{Foreground: "#242"})
}

type Content struct {
	imageView string
	height    int
	sixelMode bool
}

func NewContent() *Content {
	return &Content{
		imageView: "",
		height:    0,
		sixelMode: false,
	}
}

func (c *Content) SetView(view string) {
	c.imageView = view
}

func (c *Content) SetHeight(height int) {
	c.height = height
}

func (c *Content) SetSixelMode(sixel bool) {
	c.sixelMode = sixel
}

func (c *Content) MinSize() (int, int) {
	if c.imageView == "" {
		return 40, 10
	}
	lines := strings.Split(c.imageView, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	return maxWidth, len(lines)
}

func (c *Content) PreferredSize() (int, int) {
	minW, minH := c.MinSize()
	return minW, minH
}

func (c *Content) Render(buf *buffer.Buffer, rect buffer.Rect) {
	if c.imageView == "" {
		centerY := rect.Y + rect.Height/2
		buf.WriteString(buffer.Point{X: rect.X + 2, Y: centerY}, "No image loaded", buffer.Style{Foreground: "#242"})
		return
	}

	if c.sixelMode {
		lines := strings.Split(c.imageView, "\n")
		y := rect.Y
		for _, line := range lines {
			if y >= rect.Y+rect.Height {
				break
			}
			buf.WriteString(buffer.Point{X: rect.X, Y: y}, line, buffer.Style{})
			y++
		}
		return
	}

	if termSupportsColor {
		lines := strings.Split(c.imageView, "\n")
		y := rect.Y
		for _, line := range lines {
			if y >= rect.Y+rect.Height {
				break
			}
			buf.WriteString(buffer.Point{X: rect.X, Y: y}, line, buffer.Style{})
			y++
		}
	} else {
		cleanText := removeANSICodes(c.imageView)
		lines := strings.Split(cleanText, "\n")
		y := rect.Y
		for _, line := range lines {
			if y >= rect.Y+rect.Height {
				break
			}
			buf.WriteString(buffer.Point{X: rect.X, Y: y}, line, buffer.Style{})
			y++
		}
	}
}

type RendererType int

const (
	RendererAuto RendererType = iota
	RendererKitty
	RendereriTerm2
	RendererBlocks
	RendererSixel
	RendererASCII
)

func (rt RendererType) String() string {
	switch rt {
	case RendererAuto:
		return "Auto"
	case RendererKitty:
		return "Kitty"
	case RendereriTerm2:
		return "iTerm2"
	case RendererBlocks:
		return "Blocks"
	case RendererSixel:
		return "Sixel"
	case RendererASCII:
		return "ASCII"
	default:
		return "Unknown"
	}
}

type Model struct {
	quitting   bool
	width      int
	height     int
	imgPath    string
	renderer   RendererType
	zoomLevel  float64
	zoomMode   string

	layoutMgr  *buffer.LayoutManager
	header     *Header
	footer     *Footer
	content    *Content

	imageLoaded bool
	imageError  string
	imageW      int
	imageH      int
	scaledW     int
	scaledH     int
	imageView   string
}

func NewModel(imgPath string) Model {
	m := Model{
		quitting:  false,
		imgPath:   imgPath,
		renderer:  RendererAuto,
		zoomLevel: 1.0,
		zoomMode:  "Fit",

		header:  NewHeader(imgPath),
		footer:  NewFooter(),
		content: NewContent(),

		imageLoaded: false,
		imageError:  "Image not initialized",
		imageView:   "",
	}

	m.layoutMgr = buffer.NewLayoutManager(80, 24)

	m.layoutMgr.AddComponent("header", m.header)
	m.layoutMgr.AddComponent("footer", m.footer)
	m.layoutMgr.AddComponent("content", m.content)

	// Auto-load image on startup
	m.imageLoaded = true
	m.imageError = ""
	m.width = 80
	m.height = 24

	// Initialize with default layout and content
	m.layoutMgr.ImageLayout(12)
	m.generateMockImage()
	m.content.SetView(m.imageView)
	m.header.SetLoaded(true)
	m.updateImageDimensions(80, 24)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "r":
			m.imageLoaded = !m.imageLoaded
			m.updateContent()
		case "1":
			m.renderer = RendererAuto
			m.header.SetRenderer(m.renderer)
			m.footer.SetRenderer(m.renderer)
			m.updateContent()
		case "2":
			m.renderer = RendererKitty
			m.header.SetRenderer(m.renderer)
			m.footer.SetRenderer(m.renderer)
			m.updateContent()
		case "3":
			m.renderer = RendereriTerm2
			m.header.SetRenderer(m.renderer)
			m.footer.SetRenderer(m.renderer)
			m.updateContent()
		case "4":
			m.renderer = RendererBlocks
			m.header.SetRenderer(m.renderer)
			m.footer.SetRenderer(m.renderer)
			m.updateContent()
		case "5":
			m.renderer = RendererSixel
			m.header.SetRenderer(m.renderer)
			m.footer.SetRenderer(m.renderer)
			m.updateContent()
		case "6":
			m.renderer = RendererASCII
			m.header.SetRenderer(m.renderer)
			m.footer.SetRenderer(m.renderer)
			m.updateContent()
		case "s":
			m.imgPath = "demo-" + m.imgPath
			m.header.SetPath(m.imgPath)
		case "+", "=":
			if m.zoomLevel < 5.0 {
				m.zoomLevel += 0.1
				m.zoomMode = "Manual"
				m.updateContent()
			}
		case "-", "_":
			if m.zoomLevel > 0.2 {
				m.zoomLevel -= 0.1
				m.zoomMode = "Manual"
				m.updateContent()
			}
		case "0", " ":
			m.zoomLevel = 1.0
			m.zoomMode = "Fit"
			m.updateContent()
		case "m":
			if m.zoomMode == "Fill" {
				m.zoomMode = "Fit"
			} else {
				m.zoomMode = "Fill"
			}
			m.updateContent()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layoutMgr.SetSize(msg.Width, msg.Height)
		m.header.SetWindowSize(msg.Width, msg.Height)
		m.updateImageDimensions(msg.Width, msg.Height)
		m.updateContent()
	}

	return m, nil
}

func (m *Model) updateImageDimensions(width, height int) {
	m.imageW = 640
	m.imageH = 480

	m.scaledW = int(float64(m.imageW) * m.zoomLevel)
	m.scaledH = int(float64(m.imageH) * m.zoomLevel)

	contentWidth := width
	contentHeight := height - 6

	if m.zoomMode == "Fit" {
		scaleX := float64(contentWidth) / float64(m.imageW)
		scaleY := float64(contentHeight) / float64(m.imageH)
		scale := scaleX
		if scaleY < scale {
			scale = scaleY
		}
		if scale > 1.0 {
			scale = 1.0
		}
		m.scaledW = int(float64(m.imageW) * scale)
		m.scaledH = int(float64(m.imageH) * scale)
	}

	m.header.SetImageSize(m.imageW, m.imageH, m.scaledW, m.scaledH)
}

func (m *Model) updateContent() {
	m.imageView = ""

	percent := int(m.zoomLevel * 100)
	m.footer.SetZoom(percent, m.zoomMode)

	if !m.imageLoaded {
		m.header.SetLoaded(false)
		m.header.SetError("Image not loaded")
		m.content.SetView("")
		return
	}

	m.header.SetLoaded(true)
	m.header.SetError("")

	m.generateMockImage()
	m.content.SetView(m.imageView)

	contentHeight := 0
	if m.renderer == RendererSixel {
		contentHeight = m.scaledH / 6
		m.content.SetSixelMode(true)
	} else {
		contentHeight = min(m.scaledH/2, m.height-6)
		m.content.SetSixelMode(false)
	}

	m.content.SetHeight(contentHeight)
	m.layoutMgr.ImageLayout(contentHeight)
}

func (m *Model) generateMockImage() {
	rendererName := ""
	switch m.renderer {
	case RendererSixel:
		rendererName = "Sixel"
	case RendererKitty:
		rendererName = "Kitty"
	case RendereriTerm2:
		rendererName = "iTerm2"
	case RendererBlocks:
		rendererName = "Blocks"
	case RendererASCII:
		rendererName = "ASCII"
	default:
		rendererName = "Auto"
	}

	width := min(m.scaledW, m.width-4)
	height := min(m.scaledH/2, m.height-6)

	switch m.renderer {
	case RendererBlocks:
		m.generateBlocksImage(rendererName, width, height)
	case RendererASCII:
		m.generateASCIIImage(rendererName, width, height)
	case RendererSixel:
		m.generateSixelImage(rendererName)
	default:
		m.generateFallbackImage(rendererName, width, height)
	}
}

func (m *Model) generateBlocksImage(title string, width, height int) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n    [ %s Renderer ] %dx%d\n", title, m.scaledW, m.scaledH))

	boxWidth := min(width, 40)
	boxHeight := min(height, 10)

	sb.WriteString("    ")
	for i := 0; i < boxWidth; i++ {
		sb.WriteString("─")
	}
	sb.WriteString("\n")

	for h := 0; h < boxHeight; h++ {
		sb.WriteString("    │")
		for i := 0; i < boxWidth-2; i++ {
			if h > 1 && h < boxHeight-2 {
				if h == boxHeight/2 && i == boxWidth/2 {
					sb.WriteString("🖼")
					i += 1
				} else {
					sb.WriteString(" ")
				}
			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("│\n")
	}

	sb.WriteString("    ")
	for i := 0; i < boxWidth; i++ {
		sb.WriteString("─")
	}
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("    %s %dx%d (Zoom: %s %.0f%%)\n",
		title, m.scaledW, m.scaledH, m.zoomMode, m.zoomLevel*100))

	m.imageView = sb.String()
}

func (m *Model) generateASCIIImage(title string, width, height int) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n    [%s] %dx%d\n", title, m.scaledW, m.scaledH))

	boxWidth := min(width, 30)
	boxHeight := min(height, 8)

	for h := 0; h < boxHeight; h++ {
		sb.WriteString("    ")
		for i := 0; i < boxWidth; i++ {
			if i == 0 || i == boxWidth-1 {
				sb.WriteString("|")
			} else if h == 0 || h == boxHeight-1 {
				sb.WriteString("-")
			} else if h == boxHeight/2 && i == boxWidth/2 {
				sb.WriteString("IMG")
				i += 2
			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("    %s %dx%d (Zoom: %s)\n",
		m.renderer, m.scaledW, m.scaledH, m.zoomMode))

	m.imageView = sb.String()
}

func (m *Model) generateSixelImage(title string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n[ Sixel Mode - %dx%d ]\n", m.scaledW, m.scaledH))
	sb.WriteString("    Simulated Sixel graphics output\n")
	sb.WriteString("    (Would render actual image data here)\n")
	sb.WriteString(fmt.Sprintf("    Zoom: %s %.0f%%\n", m.zoomMode, m.zoomLevel*100))

	lines := m.scaledH / 6
	if lines > 10 {
		lines = 10
	}
	for i := 0; i < lines; i++ {
		sb.WriteString(fmt.Sprintf("    [%s] Sixel graphics line %d/%d\n",
			title, i+1, lines))
	}

	m.imageView = sb.String()
}

func (m *Model) generateFallbackImage(title string, width, height int) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n%s Renderer: %dx%d\n",
		title, m.scaledW, m.scaledH))

	boxWidth := min(width, 35)
	boxHeight := min(height, 6)

	for h := 0; h < boxHeight; h++ {
		for i := 0; i < boxWidth; i++ {
			if h == 0 || h == boxHeight-1 {
				if i == 0 {
					sb.WriteString("+")
				} else if i == boxWidth-1 {
					sb.WriteString("+")
				} else {
					sb.WriteString("-")
				}
			} else {
				if i == 0 || i == boxWidth-1 {
					sb.WriteString("|")
				} else if h == boxHeight/2 && i == boxWidth/2 {
					sb.WriteString("X")
				} else {
					sb.WriteString(" ")
				}
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("  %s Mode - Size: %dx%d - Zoom: %s\n",
		m.renderer, m.scaledW, m.scaledH, m.zoomMode))

	m.imageView = sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	return m.layoutMgr.Render()
}
