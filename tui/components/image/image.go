package image

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-sixel"
	"github.com/nfnt/resize"
	"github.com/qeesung/image2ascii/convert"
	"github.com/wwsheng009/taproot/ui/styles"
	"github.com/wwsheng009/taproot/tui/util"
	"golang.org/x/term"
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
	RendererSixel   // Sixel protocol (DEC graphics)
	RendererKitty   // Kitty graphics protocol
	RendereriTerm2  // iTerm2 inline images protocol
	RendererBlocks  // Unicode block characters with ANSI colors
	RendererASCII   // ASCII/ANSI art (fallback)
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
	width            int
	height           int
	path             string
	renderer         RendererType
	loaded           bool
	error            string
	aspectRatio      float64
	styles           *styles.Styles
	img              image.Image
	scaled           image.Image
	cachedView       string     // Cached rendered output
	cacheValid       bool       // Whether cache is valid
	sixelCacheValid  bool       // Whether Sixel cache is valid
	cacheValidASCII  bool       // Whether ASCII cache is valid
	blockSize        int        // Size of each rendered block (default 1)
	scale            float64    // Zoom scale factor (1.0 = fit to screen)
	zoomMode         ZoomMode   // Current zoom mode
	sixelCache       string     // Cached Sixel output
	asciiCache       string     // Cached ASCII output
	detectedRenderer RendererType // Cached detected renderer
	sixelSupported   bool       // Whether Sixel is supported (cached)
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
	img.detectedRenderer = RendererAuto // Reset detected renderer
	// Invalidate all caches when renderer changes
	img.cacheValid = false
	img.sixelCacheValid = false  // Also invalidate Sixel cache
	img.cacheValidASCII = false  // Invalidate ASCII cache
	img.cachedView = ""  // Clear cached content to prevent cross-protocol contamination
	img.sixelCacheValid = false
	img.sixelCache = ""   // Clear Sixel cache
	img.cacheValidASCII = false
	img.asciiCache = ""   // Clear ASCII cache

	// We need to clear screen when switching between graphic and text renderers
	// especially when leaving Sixel mode to prevent graphical residue
	return tea.Batch(
		tea.ClearScreen,
		img.Reload(),
	)
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
	// We don't handle tea.KeyMsg here - the parent component handles it
	// This prevents double-processing of keyboard events
	switch msg := msg.(type) {
	case ImageLoadedMsg:
		// Handle image loaded message
		if msg.Error != "" {
			img.error = msg.Error
			img.loaded = false
			// Clear old image on error
			img.img = nil
			img.scaled = nil
		} else {
			// Release old scaled image to free memory
			img.scaled = nil
			// Set new image
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
			// Invalidate all caches when image is reloaded
			img.cacheValid = false
			img.sixelCacheValid = false
			img.cacheValidASCII = false
		}
	case tea.WindowSizeMsg:
		// Invalidate all caches on window resize to force re-render
		img.cacheValid = false
		img.cachedView = ""
		img.sixelCacheValid = false
		img.sixelCache = ""
		img.cacheValidASCII = false
		img.asciiCache = ""
	}

	return img, nil
}

func (img *Image) View() string {
	s := img.styles

	if img.error != "" {
		// Show error message
		errorStyle := lipgloss.NewStyle().
			Foreground(s.Error).
			Bold(true)
		
		if img.width > 0 {
			errorStyle.Width(img.width).Align(lipgloss.Center)
		}

		return errorStyle.Render("⚠️  " + img.error)
	}

	if !img.loaded {
		// Show loading state
		loadingStyle := lipgloss.NewStyle().
			Foreground(s.FgMuted)
		
		if img.width > 0 {
			loadingStyle.Width(img.width).Align(lipgloss.Center)
		}

		return loadingStyle.Render("Loading image...")
	}

	// Detect renderer type
	renderer := img.renderer
	if renderer == RendererAuto {
		renderer = detectRenderer()
	}

	// Render based on type
	switch renderer {
	case RendererSixel:
		return img.renderSixel()
	case RendererKitty:
		return img.renderKitty()
	case RendereriTerm2:
		return img.renderiTerm2()
	case RendererASCII:
		return img.renderASCII()
	default:
		return img.renderBlocks()
	}
}

// renderSixel uses the Sixel protocol to output high-quality images
func (img *Image) renderSixel() string {
	if img.img == nil {
		return img.renderPlaceholder()
	}

	// Return cached Sixel output if available
	if img.sixelCacheValid && img.sixelCache != "" {
		return img.sixelCache
	}

	// Prepare image: resize to fit terminal width
	var displayImg image.Image

	if img.scaled != nil {
		// Use pre-scaled image from zoom system
		// img.scaled is already in the correct character cell dimensions
		// For Sixel output, convert char cells to pixels (1 char ≈ 2 pixels)
		scaledBounds := img.scaled.Bounds()
		scaledW := scaledBounds.Dx()
		scaledH := scaledBounds.Dy()

		// Convert char dimensions to pixel dimensions
		// Width: 1 char cell ≈ 2 pixels
		// Height: 1 char cell ≈ 1 pixel (Sixel is pixel-based)
		targetWidth := uint(scaledW * 2)
		targetHeight := uint(scaledH)

		// Resize the pre-scaled image to pixel dimensions
		displayImg = resize.Resize(targetWidth, targetHeight, img.scaled, resize.Lanczos3)
	} else {
		// Apply reasonable default scaling
		img.autoScale()
		displayImg = img.img

		// Calculate target width in pixels
		targetWidth := uint(img.width * 2)
		if targetWidth == 0 {
			targetWidth = 800 // Default width
		}

		// Resize to fit terminal width
		displayImg = resize.Resize(targetWidth, 0, displayImg, resize.Lanczos3)
	}

	// encode to Sixel
	var buf bytes.Buffer
	enc := sixel.NewEncoder(&buf)
	enc.Dither = true // Enable dithering for better quality
	if err := enc.Encode(displayImg); err != nil {
		// If Sixel encoding fails, fall back to block rendering
		return img.renderBlocks()
	}

	output := buf.String()

	// Add cursor positioning after Sixel to ensure text appears on the next line
	// Sixel doesn't automatically move the cursor, so we need to do it explicitly
	output += "\n"

	// Tmux passthrough support
	if os.Getenv("TMUX") != "" {
		// Wrap Sixel data in tmux passthrough sequence
		// This tells tmux to pass through these escape sequences directly to the terminal
		output = "\x1bPtmux;\x1b" + output + "\x1b\\"
	}

	img.sixelCache = output
	img.sixelCacheValid = true

	return output
}

// renderKitty uses colored blocks to render the image (simplified version)
func (img *Image) renderKitty() string {
	return img.renderBlocks()
}

// renderiTerm2 uses colored blocks to render the image (simplified version)
func (img *Image) renderiTerm2() string {
	return img.renderBlocks()
}

// renderASCII uses ASCII/ANSI art fallback rendering
func (img *Image) renderASCII() string {
	if !img.loaded || img.img == nil {
		return img.renderPlaceholder()
	}

	// Return cached ASCII output if available
	if img.cacheValidASCII && img.asciiCache != "" {
		return img.asciiCache
	}

	// Use image2ascii library for ASCII conversion
	converter := convert.NewImageConverter()

	// Configure conversion options
	convertOptions := convert.DefaultOptions

	// Use scaled image if available (supports zoom and zoom modes)
	var sourceImage image.Image
	if img.scaled != nil {
		sourceImage = img.scaled
	} else {
		sourceImage = img.img
	}

	// Calculate width based on zoom mode and scale
	// For ASCII, we work in character units directly
	baseWidth := img.width - 4  // Account for padding/border

	// Apply zoom scale to width
	width := int(float64(baseWidth) * img.scale)

	// Ensure minimum width
	if width < 20 {
		width = 20
	}

	convertOptions.FixedWidth = width  // Limit to terminal width considering zoom
	convertOptions.Colored = true       // Use ANSI colors
	convertOptions.Reversed = false     // Don't reverse colors
	convertOptions.FitScreen = false    // Use custom sizing based on zoom mode and scale

	// Convert image to ASCII string
	// Note: image2ascii library automatically adjusts height to maintain aspect ratio
	result := converter.Image2ASCIIString(sourceImage, &convertOptions)

	// Apply styling
	s := img.styles
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Border).
		Padding(1).
		Width(img.width).
		Align(lipgloss.Center)

	rendered := boxStyle.Render(result)

	img.asciiCache = rendered
	img.cacheValidASCII = true

	return rendered
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

	// Sixel detection - if terminal supports Sixel, try to use it
	if supportsSixel() {
		return RendererSixel
	}

	// Kitty detection
	if strings.Contains(term, "kitty") || strings.Contains(program, "kitty") {
		return RendererKitty
	}

	// iTerm2 detection
	if program == "iTerm.app" {
		return RendereriTerm2
	}

	// Default to blocks (works everywhere with UTF-8)
	return RendererBlocks
}

// supportsSixel checks if the current terminal supports Sixel protocol
func supportsSixel() bool {
	// Get Stdin file descriptor
	fd := int(os.Stdin.Fd())

	// If not a terminal, return false
	if !term.IsTerminal(fd) {
		return false
	}

	// Check for known Sixel-supporting terminals via environment variables
	termProg := os.Getenv("TERM_PROGRAM")
	termEnv := os.Getenv("TERM")

	// Mintty (Git Bash, Windows) supports Sixel
	if strings.Contains(termProg, "mintty") {
		return true
	}

	// WezTerm supports Sixel
	if termProg == "WezTerm" {
		return true
	}

	// Check TERM variable for Sixel support
	sixelTerms := []string{"xterm", "vt340", "vt330"}
	for _, t := range sixelTerms {
		if strings.Contains(termEnv, t) {
			// Could also check Sixel-specific versions like xterm-sixel
			return true
		}
	}

	// Direct DA (Device Attributes) query
	// This is the most reliable method but temporarily sets terminal to raw mode
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return false
	}
	defer term.Restore(fd, oldState)

	// Send DA query: ESC [ c
	_, err = os.Stdout.Write([]byte("\x1b[c"))
	if err != nil {
		return false
	}

	// Read response with timeout
	result := make([]byte, 0, 100)
	buffer := make([]byte, 1)

	done := make(chan bool)
	go func() {
		for {
			n, err := os.Stdin.Read(buffer)
			if err != nil || n == 0 {
				break
			}
			result = append(result, buffer[0])
			if buffer[0] == 'c' {
				break
			}
		}
		done <- true
	}()

	select {
	case <-done:
		// Read completed, check for Sixel capability (code 4)
		response := string(result)
		return strings.Contains(response, ";4") || strings.Contains(response, "?4")
	case <-time.After(100 * time.Millisecond):
		// Timeout - terminal doesn't respond or doesn't support DA query
		return false
	}
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
		// Release old scaled image to free memory
		img.scaled = nil
		// Invalidate all caches
		img.cacheValid = false
		img.cachedView = ""
		img.sixelCacheValid = false
		img.sixelCache = ""
		img.cacheValidASCII = false
		img.asciiCache = ""
		img.scaleImage()
	}
}

// scaleImage scales the image to fit the current size
func (img *Image) scaleImage() {
	if img.img == nil || img.width <= 0 || img.height <= 0 {
		return
	}

	// Release old scaled image to free memory
	img.scaled = nil

	// Invalidate all caches
	img.cacheValid = false
	img.cachedView = ""
	img.sixelCacheValid = false
	img.sixelCache = ""
	img.cacheValidASCII = false
	img.asciiCache = ""

	srcBounds := img.img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// Calculate aspect ratio (width / height)
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
		// First try fitting to width
		baseW = availableW
		baseH = int(float64(baseW) / aspectRatio)
		// If height exceeds, fit to height instead
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

// GetImageDimensions returns the original image dimensions
func (img *Image) GetImageDimensions() (width, height int) {
	if img.img == nil {
		return 0, 0
	}
	bounds := img.img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

// SetPath changes the image path and reloads
func (img *Image) SetPath(path string) tea.Cmd {
	img.path = path
	// Clear screen when changing image to prevent visual artifacts
	return tea.Batch(
		tea.ClearScreen,
		img.Reload(),
	)
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
		// Release old scaled image to free memory
		img.scaled = nil
		// Invalidate all caches when zoom changes
		img.cacheValid = false
		img.cachedView = ""
		img.sixelCacheValid = false
		img.sixelCache = ""
		img.cacheValidASCII = false
		img.asciiCache = ""
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
		// Release old scaled image to free memory
		img.scaled = nil
		// Invalidate all caches when zoom changes
		img.cacheValid = false
		img.cachedView = ""
		img.sixelCacheValid = false
		img.sixelCache = ""
		img.cacheValidASCII = false
		img.asciiCache = ""
		img.scaleImage()
	}

	return nil
}

// ResetZoom resets the zoom factor to 1.0 (fit to screen)
func (img *Image) ResetZoom() tea.Cmd {
	img.scale = 1.0

	// Re-scale the image with default zoom factor
	if img.img != nil {
		// Release old scaled image to free memory
		img.scaled = nil
		// Invalidate all caches when zoom changes
		img.cacheValid = false
		img.cachedView = ""
		img.sixelCacheValid = false
		img.sixelCache = ""
		img.cacheValidASCII = false
		img.asciiCache = ""
		img.scaleImage()
	}

	return nil
}

// GetScale returns the current zoom scale factor
func (img *Image) GetScale() float64 {
	return img.scale
}

// SetScale sets the zoom scale factor directly
func (img *Image) SetScale(scale float64) tea.Cmd {
	// Clamp scale to valid range
	if scale < 0.1 {
		scale = 0.1
	}
	if scale > 5.0 {
		scale = 5.0
	}
	img.scale = scale

	// Re-scale the image with new zoom factor
	if img.img != nil {
		// Release old scaled image to free memory
		img.scaled = nil
		// Invalidate all caches when zoom changes
		img.cacheValid = false
		img.cachedView = ""
		img.sixelCacheValid = false
		img.sixelCache = ""
		img.cacheValidASCII = false
		img.asciiCache = ""
		img.scaleImage()
	}

	return nil
}

// GetRenderer returns the current renderer type
func (img *Image) GetRenderer() RendererType {
	return img.renderer
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
		// Invalidate all caches when zoom mode changes
		img.cacheValid = false
		img.cachedView = ""
		img.sixelCacheValid = false
		img.sixelCache = ""
		img.cacheValidASCII = false
		img.asciiCache = ""
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

// GetRendererDescription returns a human-readable description for the renderer type
func GetRendererDescription(renderer RendererType) string {
	switch renderer {
	case RendererAuto:
		return "Auto - Automatically detect the best available renderer"
	case RendererSixel:
		return "Sixel - DEC graphics protocol for high-quality terminal images"
	case RendererKitty:
		return "Kitty - Kitty graphics protocol for modern terminals"
	case RendereriTerm2:
		return "iTerm2 - iTerm2 inline images protocol"
	case RendererBlocks:
		return "Blocks - Unicode block characters with ANSI colors"
	case RendererASCII:
		return "ASCII - ASCII/ANSI art fallback renderer"
	default:
		return "Unknown - Unrecognized renderer type"
	}
}
