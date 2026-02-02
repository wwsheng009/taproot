package image

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/styles"
	"github.com/wwsheng009/taproot/tui/util"
)

// ImageLoadedMsg is sent when image loading completes
type ImageLoadedMsg struct {
	Image  image.Image
	Loaded bool
	Error  string
}

// RendererType represents the image rendering protocol
type RendererType int

const (
	RendererAuto RendererType = iota
	RendererKitty
	RendereriTerm2
	RendererBlocks // Unicode block characters (fallback)
)

// ZoomMode represents different scaling behaviors
type ZoomMode int

const (
	ZoomFitScreen ZoomMode = iota // Fit within available space (maintain aspect ratio)
	ZoomFitWidth                   // Fit width to available, allow overflow in height
	ZoomFitHeight                  // Fit height to available, allow overflow in width
	ZoomFill                       // Fill available space (may distort aspect ratio)
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
	img         image.Image
	scaled      image.Image
	cachedView  string        // Cached rendered output
	cacheValid  bool          // Whether cache is valid
	blockSize   int           // Size of each rendered block (default 1)
	scale       float64       // Zoom scale factor (1.0 = fit to screen)
	zoomMode    ZoomMode      // Current zoom mode
}

const (
	maxWidth = 100
)

// New creates a new image component
func New(path string) *Image {
	s := styles.DefaultStyles()
	return &Image{
		path:      path,
		renderer:  RendererAuto,
		loaded:    false,
		styles:    &s,
		blockSize: 1, // Default to 1x1 blocks
		scale:     1.0, // Default zoom level (fit to screen)
		zoomMode:  ZoomFitScreen, // Default mode
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
			return ImageLoadedMsg{
				Error: fmt.Sprintf("File not found: %s", img.path),
				Loaded: false,
			}
		}

		// Open and decode the image
		file, err := os.Open(img.path)
		if err != nil {
			return ImageLoadedMsg{
				Error: fmt.Sprintf("Failed to open file: %v", err),
				Loaded: false,
			}
		}
		defer file.Close()

		decoded, _, err := image.Decode(file)
		if err != nil {
			return ImageLoadedMsg{
				Error: fmt.Sprintf("Failed to decode image: %v", err),
				Loaded: false,
			}
		}

		return ImageLoadedMsg{
			Image:  decoded,
			Loaded: true,
			Error:  "",
		}
	}
}

func (img *Image) Init() tea.Cmd {
	return img.Reload()
}

func (img *Image) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ImageLoadedMsg:
		// Handle image loaded message
		if msg.Error != "" {
			img.error = msg.Error
			img.loaded = false
		} else {
			img.img = msg.Image
			img.loaded = true
			img.error = ""
			// Scale image if size is set
			if img.width > 0 && img.height > 0 {
				img.scaleImage()
			} else {
				// Auto-scale based on image size (max 80x40 terminal cells)
				img.autoScale()
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return img, img.Reload()
		}
	case tea.WindowSizeMsg:
		// Invalidate cache on window resize to force re-render
		img.cacheValid = false
		img.cachedView = ""
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

// renderKitty uses colored blocks to render the image (simplified version)
func (img *Image) renderKitty() string {
	return img.renderBlocks()
}

// renderiTerm2 uses colored blocks to render the image (simplified version)
func (img *Image) renderiTerm2() string {
	return img.renderBlocks()
}

// renderBlocks uses colored blocks to render the image
func (img *Image) renderBlocks() string {
	// If image is not loaded or scaled, show placeholder
	if img.scaled == nil {
		return img.renderPlaceholder()
	}

	// Return cached view if available
	if img.cacheValid && img.cachedView != "" {
		return img.cachedView
	}

	s := img.styles
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Border).
		Padding(1)

	bounds := img.scaled.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	var sb strings.Builder
	// Pre-allocate capacity for better performance (estimated)
	sb.Grow(h * w * 20)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			pixel := img.scaled.At(x, y)
			r, g, b, _ := pixel.RGBA()
			// Convert to 8-bit color and write ANSI sequence directly
			sb.WriteString("\x1b[48;2;")

			// Write red component
			rc := r >> 8
			if rc >= 100 {
				sb.WriteByte('0' + byte(rc/100))
			}
			if rc >= 10 {
				sb.WriteByte('0' + byte((rc%100)/10))
			}
			sb.WriteByte('0' + byte(rc%10))

			sb.WriteByte(';')

			// Write green component
			gc := g >> 8
			if gc >= 100 {
				sb.WriteByte('0' + byte(gc/100))
			}
			if gc >= 10 {
				sb.WriteByte('0' + byte((gc%100)/10))
			}
			sb.WriteByte('0' + byte(gc%10))

			sb.WriteByte(';')

			// Write blue component
			bc := b >> 8
			if bc >= 100 {
				sb.WriteByte('0' + byte(bc/100))
			}
			if bc >= 10 {
				sb.WriteByte('0' + byte((bc%100)/10))
			}
			sb.WriteByte('0' + byte(bc%10))

			sb.WriteString("m ")
		}
		sb.WriteString("\x1b[0m\n")
	}

	rendered := boxStyle.Render(sb.String())

	// Cache the result
	img.cachedView = lipgloss.NewStyle().
		Width(img.width).
		Align(lipgloss.Center).
		Render(rendered)
	img.cacheValid = true

	return img.cachedView
}

// renderPlaceholder shows a placeholder when image is not loaded
func (img *Image) renderPlaceholder() string {
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

	// Handle boundary case when width is not yet set
	hrizWidth := img.width - 4
	if hrizWidth < 2 {
		hrizWidth = 2
	}
	sb.WriteString("│" + strings.Repeat("─", hrizWidth) + "│\n")
	sb.WriteString(fmt.Sprintf("│ %dx%d  │\n", max(img.width, 20), min(img.height, 20)))
	sb.WriteString("└" + strings.Repeat("─", hrizWidth) + "┘")

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
	oldWidth := img.width
	oldHeight := img.height
	img.width = w
	img.height = h

	// Re-scale image if already loaded
	// Also re-scale if old width was 0 (first time getting window size)
	if img.img != nil && (oldWidth != w || oldHeight != h || oldWidth == 0) {
		// Invalidate cache
		img.cacheValid = false
		img.cachedView = ""
		img.scaleImage()
	}
}

// scaleImage scales the image to fit the current size
func (img *Image) scaleImage() {
	if img.img == nil || img.width <= 0 || img.height <= 0 {
		return
	}

	// Invalidate cache
	img.cacheValid = false
	img.cachedView = ""

	srcBounds := img.img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// Calculate aspect ratio
	aspectRatio := float64(srcW) / float64(srcH)

	// Determine display size
	// displayW and displayH are measured in character cells (not pixels)
	// Border adds 2 chars (left/right), Padding adds 2 chars (left/right) = 4 total
	availableW := img.width - 4
	availableH := img.height - 4

	// Step 1: Calculate base dimensions based on zoom mode
	var baseW, baseH int
	
	switch img.zoomMode {
	case ZoomFitScreen:
		// Fit within available space (maintain aspect ratio)
		baseW = availableW
		baseH = int(float64(baseW) / aspectRatio)
		if baseH > availableH {
			baseH = availableH
			baseW = int(float64(baseH) * aspectRatio)
		}
		
	case ZoomFitWidth:
		// Fit width to available, allow height overflow
		baseW = availableW
		baseH = int(float64(baseW) / aspectRatio)
		
	case ZoomFitHeight:
		// Fit height to available, allow width overflow
		baseH = availableH
		baseW = int(float64(baseH) * aspectRatio)
		
	case ZoomFill:
		// Fill available space (may distort aspect ratio)
		baseW = availableW
		baseH = availableH
	}

	// Step 2: Apply zoom scale factor
	displayW := int(float64(baseW) * img.scale)
	displayH := int(float64(baseH) * img.scale)

	// Step 3: Ensure minimum dimensions (don't go below 1)
	if displayW < 1 {
		displayW = 1
	}
	if displayH < 1 {
		displayH = 1
	}

	// If available space is too small, use minimum reasonable size
	if availableW < 5 || availableH < 3 {
		// Very small terminal, use minimal size
		displayW = max(1, availableW)
		displayH = max(1, availableH)
	}

	// Create scaled image
	scaled := image.NewRGBA(image.Rect(0, 0, displayW, displayH))
	for y := 0; y < displayH; y++ {
		for x := 0; x < displayW; x++ {
			// Map pixel from source
			srcX := x * srcW / displayW
			srcY := y * srcH / displayH
			scaled.Set(x, y, img.img.At(srcX, srcY))
		}
	}

	img.scaled = scaled
	img.aspectRatio = aspectRatio
}

// autoScale automatically scales the image to a reasonable size based on original dimensions
func (img *Image) autoScale() {
	if img.img == nil {
		return
	}

	// If window size is already set, use scaleImage instead
	if img.width > 0 && img.height > 0 {
		img.scaleImage()
		return
	}

	srcBounds := img.img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// Calculate aspect ratio
	aspectRatio := float64(srcW) / float64(srcH)

	// Set reasonable default size (max 80 columns x 40 rows for terminal)
	maxW := 80
	maxH := 40

	// Determine display size maintaining aspect ratio
	displayW := maxW
	displayH := int(float64(displayW) / aspectRatio / 0.5)

	if displayH > maxH {
		displayH = maxH
		displayW = int(float64(displayH) * aspectRatio * 0.5)
	}

	// Set the size (will be overridden if WindowSizeMsg arrives)
	img.width = displayW + 4 // Account for padding/border (2 + 2)
	img.height = displayH + 4

	// Create scaled image
	scaled := image.NewRGBA(image.Rect(0, 0, displayW, displayH))
	for y := 0; y < displayH; y++ {
		for x := 0; x < displayW; x++ {
			// Map pixel from source
			srcX := x * srcW / displayW
			srcY := y * srcH / displayH
			scaled.Set(x, y, img.img.At(srcX, srcY))
		}
	}

	img.scaled = scaled
	img.aspectRatio = aspectRatio
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

// ScaledSize returns the size of the scaled image (actual pixel data, not including border/padding)
func (img *Image) ScaledSize() (width, height int) {
	if img.scaled == nil {
		return 0, 0
	}
	bounds := img.scaled.Bounds()
	return bounds.Dx(), bounds.Dy()
}

// ZoomIn increases the zoom scale factor
func (img *Image) ZoomIn() tea.Cmd {
	if img.scale >= 4.0 {
		// Maximum zoom level reached
		return nil
	}
	img.scale *= 1.25 // Increase by 25%
	if img.scale > 4.0 {
		img.scale = 4.0
	}
	
	// Re-scale the image with new zoom factor
	if img.img != nil {
		img.cacheValid = false
		img.cachedView = ""
		img.scaleImage()
	}
	
	return nil
}

// ZoomOut decreases the zoom scale factor
func (img *Image) ZoomOut() tea.Cmd {
	if img.scale <= 0.25 {
		// Minimum zoom level reached
		return nil
	}
	img.scale /= 1.25 // Decrease by 25%
	if img.scale < 0.25 {
		img.scale = 0.25
	}
	
	// Re-scale the image with new zoom factor
	if img.img != nil {
		img.cacheValid = false
		img.cachedView = ""
		img.scaleImage()
	}
	
	return nil
}

// ResetZoom resets the zoom factor to 1.0 (fit to screen)
func (img *Image) ResetZoom() tea.Cmd {
	img.scale = 1.0
	
	// Re-scale the image with default zoom factor
	if img.img != nil {
		img.cacheValid = false
		img.cachedView = ""
		img.scaleImage()
	}
	
	return nil
}

// GetScale returns the current zoom scale factor
func (img *Image) GetScale() float64 {
	return img.scale
}

// GetZoomMode returns the current zoom mode
func (img *Image) GetZoomMode() ZoomMode {
	return img.zoomMode
}

// SetZoomMode sets the zoom mode and re-scales the image
func (img *Image) SetZoomMode(mode ZoomMode) tea.Cmd {
	img.zoomMode = mode
	
	// Re-scale the image with new zoom mode
	if img.img != nil {
		img.cacheValid = false
		img.cachedView = ""
		img.scaleImage()
	}
	
	return nil
}

// CycleZoomMode cycles through available zoom modes
func (img *Image) CycleZoomMode() tea.Cmd {
	// Cycle: FitScreen -> FitWidth -> FitHeight -> Fill -> FitScreen
	nextMode := img.zoomMode + 1
	if nextMode > ZoomFill {
		nextMode = ZoomFitScreen
	}
	
	img.zoomMode = nextMode
	
	// Re-scale the image with new zoom mode
	if img.img != nil {
		img.cacheValid = false
		img.cachedView = ""
		img.scaleImage()
	}
	
	return nil
}

// GetZoomModeName returns a human-readable name for the zoom mode
func (img *Image) GetZoomModeName() string {
	switch img.zoomMode {
	case ZoomFitScreen:
		return "Fit"
	case ZoomFitWidth:
		return "Width"
	case ZoomFitHeight:
		return "Height"
	case ZoomFill:
		return "Fill"
	default:
		return "Unknown"
	}
}
