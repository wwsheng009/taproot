package image

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/components/image/decoder"
	"github.com/wwsheng009/taproot/ui/styles"
)

// RendererType represents the image rendering protocol
type RendererType int

const (
	RendererAuto RendererType = iota
	RendererKitty
	RendereriTerm2
	RendererBlocks // Unicode block characters (fallback)
	RendererSixel
	RendererASCII  // Pure ASCII fallback
)

// ZoomMode represents how the image is scaled
type ZoomMode int

const (
	ZoomFit      ZoomMode = iota // Fit to screen (maintain aspect ratio)
	ZoomFill                     // Fill screen (may crop)
	ZoomStretch                  // Stretch to fill (ignore aspect ratio)
	ZoomOriginal                 // Original size (may scroll)
)

// String returns the zoom mode name
func (z ZoomMode) String() string {
	switch z {
	case ZoomFit:
		return "Fit"
	case ZoomFill:
		return "Fill"
	case ZoomStretch:
		return "Stretch"
	case ZoomOriginal:
		return "Original"
	default:
		return "Unknown"
	}
}

// String returns the renderer name
func (r RendererType) String() string {
	switch r {
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
		return "Auto"
	}
}

// Image displays an image in the terminal
type Image struct {
	width       int
	height      int
	path        string
	renderer    RendererType
	loaded      bool
	err         string
	aspectRatio float64
	styles      *styles.Styles

	// Image data
	imgData    *decoder.ImageData
	decoder    *decoder.Decoder

	// Renderers
	kitty      *KittyRenderer
	iterm      *ITerm2Renderer
	sixel      *SixelRenderer
	blocks     *BlocksRenderer

	// Zoom state
	zoomLevel  float64
	zoomMode   ZoomMode
}

const (
	maxWidth = 100
)

// New creates a new image component
func New(path string) *Image {
	s := styles.DefaultStyles()
	dec := decoder.NewDecoder()

	img := &Image{
		path:      path,
		renderer:  RendererAuto,
		loaded:    false,
		styles:    &s,
		decoder:   dec,
		zoomLevel: 1.0,
		zoomMode:  ZoomFit,
		err:       "Loading...",
	}

	// Load image synchronously
	img.loadImage()

	return img
}

// loadImage loads the image synchronously
func (img *Image) loadImage() {
	// Clear previous error
	img.err = ""
	img.loaded = false

	// Validate path
	if err := decoder.ValidatePath(img.path); err != nil {
		img.err = err.Error()
		return
	}

	// Decode image
	data, err := img.decoder.DecodeFile(img.path)
	if err != nil {
		img.err = err.Error()
		return
	}

	img.imgData = data
	img.loaded = true
	img.aspectRatio = float64(data.Width) / float64(data.Height)

	// Create renderers
	img.kitty = NewKittyRenderer(data)
	img.iterm = NewITerm2Renderer(data)
	img.sixel = NewSixelRenderer(data)
	img.blocks = NewBlocksRenderer(data)

	// Detect color support for blocks renderer
	img.blocks.SetColorEnabled(DetectColorSupport())
}

// SetRenderer sets the rendering type
func (img *Image) SetRenderer(renderer RendererType) render.Cmd {
	img.renderer = renderer
	return nil
}

// Reload reloads the image (synchronous)
func (img *Image) Reload() render.Cmd {
	img.loadImage()
	return nil
}

// Init implements render.Model
func (img *Image) Init() render.Cmd {
	return img.Reload()
}

// Update implements render.Model
func (img *Image) Update(msg any) (render.Model, render.Cmd) {
	switch msg := msg.(type) {
	case render.KeyMsg:
		switch msg.String() {
		// Reload
		case "r":
			return img, img.Reload()

		// Renderer selection
		case "1":
			return img, img.SetRenderer(RendererAuto)
		case "2":
			return img, img.SetRenderer(RendererKitty)
		case "3":
			return img, img.SetRenderer(RendereriTerm2)
		case "4":
			return img, img.SetRenderer(RendererBlocks)
		case "5":
			return img, img.SetRenderer(RendererSixel)
		case "6":
			return img, img.SetRenderer(RendererASCII)

		// Zoom controls
		case "+", "=":
			img.ZoomIn()
			return img, nil
		case "-", "_":
			img.ZoomOut()
			return img, nil
		case "0":
			img.ResetZoom()
			return img, nil
		case "*":
			img.zoomLevel = 2.0
			return img, nil
		case "%":
			img.zoomLevel = 0.5
			return img, nil

		// Zoom mode switching
		case "m":
			img.CycleZoomMode()
			return img, nil
		case "f":
			img.SetZoomMode(ZoomFit)
			return img, nil
		case "F":
			img.SetZoomMode(ZoomFill)
			return img, nil
		case "s":
			img.SetZoomMode(ZoomStretch)
			return img, nil
		case "o":
			img.SetZoomMode(ZoomOriginal)
			return img, nil

		// Fine zoom control
		case "[":
			img.zoomLevel -= 0.01
			if img.zoomLevel < 0.1 {
				img.zoomLevel = 0.1
			}
			return img, nil
		case "]":
			img.zoomLevel += 0.01
			if img.zoomLevel > 5.0 {
				img.zoomLevel = 5.0
			}
			return img, nil
		}
	case render.WindowSizeMsg:
		img.width = msg.Width
		img.height = msg.Height
	}

	return img, nil
}

// View implements render.Model
func (img *Image) View() string {
	if img.err != "" {
		return img.renderError(img.err)
	}

	if !img.loaded || img.imgData == nil {
		return img.renderLoading()
	}

	// Detect renderer type
	renderer := img.renderer
	if renderer == RendererAuto {
		renderer = img.detectRenderer()
	}

	// Get display size (where to render)
	displayW, displayH := img.calculateDisplaySize()

	// Render based on type
	switch renderer {
	case RendererKitty:
		return img.renderKitty(displayW, displayH)
	case RendereriTerm2:
		return img.renderITerm2(displayW, displayH)
	case RendererSixel:
		return img.renderSixel(displayW, displayH)
	case RendererASCII:
		return img.renderASCII(displayW, displayH)
	case RendererBlocks:
		return img.renderBlocks(displayW, displayH)
	default:
		return img.renderBlocks(displayW, displayH)
	}
}

// renderError displays an error message
func (img *Image) renderError(errMsg string) string {
	s := img.styles
	errorStyle := lipgloss.NewStyle().
		Foreground(s.Error).
		Bold(true).
		Width(img.width).
		Align(lipgloss.Center).
		Padding(1)

	return errorStyle.Render("⚠️  " + errMsg)
}

// renderLoading displays a loading state
func (img *Image) renderLoading() string {
	s := img.styles
	loadingStyle := lipgloss.NewStyle().
		Foreground(s.FgMuted).
		Width(img.width).
		Align(lipgloss.Center).
		Padding(1)

	return loadingStyle.Render("Loading image...")
}

// renderKitty renders using Kitty graphics protocol
func (img *Image) renderKitty(width, height int) string {
	if img.kitty == nil {
		return img.renderFallback("Kitty renderer not available")
	}

	// Set the sampled dimensions based on zoom level
	sampledW, sampledH := img.ScaledSize()
	img.kitty.SetSampledSize(sampledW, sampledH)

	// Set cell size (default)
	img.kitty.SetCellSize(10, 20)

	output := img.kitty.Render(width, height)
	return output
}

// renderITerm2 renders using iTerm2 inline images
func (img *Image) renderITerm2(width, height int) string {
	if img.iterm == nil {
		return img.renderFallback("iTerm2 renderer not available")
	}

	// Set the sampled dimensions based on zoom level
	sampledW, sampledH := img.ScaledSize()
	img.iterm.SetSampledSize(sampledW, sampledH)

	img.iterm.SetCellSize(10, 20)
	return img.iterm.Render(width, height)
}

// renderSixel renders using Sixel graphics
func (img *Image) renderSixel(width, height int) string {
	if img.sixel == nil {
		return img.renderFallback("Sixel renderer not available")
	}

	// Windows Terminal has limited Sixel support
	// Always fallback to Blocks for Windows Terminal users
	if env("WT_SESSION") != "" {
		// Create fallback message
		output := img.renderBlocks(width, height)
		msgStyle := img.styles.Base.
			Foreground(img.styles.FgMuted).
			Italic(true)
		return output + "\n\n" + msgStyle.Render("  Note: Sixel is not fully supported in Windows Terminal. Using Blocks renderer instead.")
	}

	// Check if terminal supports Sixel
	if !DetectSixel() {
		output := img.renderBlocks(width, height)
		return output + "\n\n" + img.styles.Base.Render("  Note: Sixel is not supported by this terminal. Using Blocks renderer instead.")
	}

	// Set the sampled dimensions based on zoom level
	sampledW, sampledH := img.ScaledSize()

	// Sixel works better with larger dimensions
	// Use at least 40x20 pixels
	if sampledW < 40 {
		sampledW = 40
	}
	if sampledH < 20 {
		sampledH = 20
	}

	img.sixel.SetSampledSize(sampledW, sampledH)

	return img.sixel.Render(width, height)
}

// renderBlocks renders using Unicode block characters
func (img *Image) renderBlocks(width, height int) string {
	if img.blocks == nil {
		return img.renderFallback("Blocks renderer not available")
	}

	// Set the sampled dimensions based on zoom level
	sampledW, sampledH := img.ScaledSize()
	img.blocks.SetSampledSize(sampledW, sampledH)
	
	// Ensure we use blocks mode (not ASCII)
	img.blocks.SetASCIIMode(false)

	return img.blocks.Render(width, height)
}

// renderASCII renders using ASCII characters
func (img *Image) renderASCII(width, height int) string {
	if img.blocks == nil {
		return img.renderFallback("ASCII renderer not available")
	}

	// Set the sampled dimensions based on zoom level
	sampledW, sampledH := img.ScaledSize()
	img.blocks.SetSampledSize(sampledW, sampledH)
	img.blocks.SetASCIIMode(true)

	return img.blocks.Render(width, height)
}

// renderFallback displays a fallback message
func (img *Image) renderFallback(msg string) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(img.styles.Border).
		Padding(1).
		Width(img.width).
		Align(lipgloss.Center)

	return boxStyle.Render(msg)
}

// detectRenderer attempts to detect the best available renderer
func (img *Image) detectRenderer() RendererType {
	// Check for Kitty support
	if DetectKitty() {
		return RendererKitty
	}

	// Check for iTerm2 support
	if DetectITerm2() {
		return RendereriTerm2
	}

	// Use Blocks as default (works everywhere with true color)
	// Sixel is detected but not fully implemented yet
	return RendererBlocks
}

// calculateDisplaySize calculates the display size for the image
func (img *Image) calculateDisplaySize() (int, int) {
	if img.width <= 0 || img.height <= 0 {
		return 80, 24
	}

	// Reserve space for UI elements
	displayH := img.height - 6 // Leave room for header/footer
	if displayH < 1 {
		displayH = 1
	}

	displayW := img.width - 4 // Leave margin
	if displayW < 1 {
		displayW = 1
	}

	return displayW, displayH
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
func (img *Image) SetPath(path string) render.Cmd {
	img.path = path
	img.loadImage()
	return nil
}

// IsLoaded returns whether the image is loaded
func (img *Image) IsLoaded() bool {
	return img.loaded
}

// Error returns any loading error
func (img *Image) Error() string {
	return img.err
}

// GetRenderer returns the current renderer type
func (img *Image) GetRenderer() RendererType {
	return img.renderer
}

// GetImageData returns the decoded image data
func (img *Image) GetImageData() *decoder.ImageData {
	return img.imgData
}

// GetImageDimensions returns the original image dimensions
func (img *Image) GetImageDimensions() (int, int) {
	if img.imgData == nil {
		return 0, 0
	}
	return img.imgData.Width, img.imgData.Height
}

// GetSupportedRenderers returns a list of renderers supported by the current terminal
func GetSupportedRenderers() []RendererType {
	var renderers []RendererType

	// Check each renderer type
	if DetectKitty() {
		renderers = append(renderers, RendererKitty)
	}
	if DetectITerm2() {
		renderers = append(renderers, RendereriTerm2)
	}
	if DetectSixel() {
		renderers = append(renderers, RendererSixel)
	}

	// Blocks and ASCII are always supported
	renderers = append(renderers, RendererBlocks)
	renderers = append(renderers, RendererASCII)

	return renderers
}

// GetRendererDescription returns a description of the renderer
func GetRendererDescription(r RendererType) string {
	switch r {
	case RendererAuto:
		return "Auto-detect best renderer"
	case RendererKitty:
		return "Kitty graphics protocol (high quality)"
	case RendereriTerm2:
		return "iTerm2 inline images (high quality)"
	case RendererSixel:
		return "Sixel graphics (medium quality, wide support)"
	case RendererBlocks:
		return "Unicode blocks (2 pixels per cell, universal support)"
	case RendererASCII:
		return "ASCII art (lowest quality, maximum compatibility)"
	default:
		return "Unknown renderer"
	}
}

// truncatePath truncates a path to fit within maxLen
func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CreateInfoBox creates an info box with image details
func (img *Image) CreateInfoBox() string {
	if img.imgData == nil {
		return ""
	}

	s := img.styles
	infoStyle := lipgloss.NewStyle().
		Foreground(s.FgMuted).
		Width(img.width)

	var info []string

	info = append(info, fmt.Sprintf("Path: %s", truncatePath(img.path, 40)))
	info = append(info, fmt.Sprintf("Size: %dx%d", img.imgData.Width, img.imgData.Height))
	info = append(info, fmt.Sprintf("Format: %s", img.imgData.Format))
	info = append(info, fmt.Sprintf("Renderer: %s", img.renderer.String()))

	// Get platform info
	platform := GetPlatformInfo()
	var support []string
	if platform.SupportsKitty {
		support = append(support, "Kitty")
	}
	if platform.SupportsITerm2 {
		support = append(support, "iTerm2")
	}
	if platform.SupportsSixel {
		support = append(support, "Sixel")
	}
	if len(support) > 0 {
		info = append(info, fmt.Sprintf("Supports: %s", strings.Join(support, ", ")))
	}

	return infoStyle.Render(strings.Join(info, " • "))
}

// ========== Zoom Methods ==========

// ZoomIn increases the zoom level
func (img *Image) ZoomIn() {
	img.zoomLevel += 0.1
	if img.zoomLevel > 5.0 {
		img.zoomLevel = 5.0
	}
}

// ZoomOut decreases the zoom level
func (img *Image) ZoomOut() {
	img.zoomLevel -= 0.1
	if img.zoomLevel < 0.1 {
		img.zoomLevel = 0.1
	}
}

// ResetZoom resets the zoom level to 1.0
func (img *Image) ResetZoom() {
	img.zoomLevel = 1.0
}

// CycleZoomMode cycles through zoom modes
func (img *Image) CycleZoomMode() {
	switch img.zoomMode {
	case ZoomFit:
		img.zoomMode = ZoomFill
	case ZoomFill:
		img.zoomMode = ZoomStretch
	case ZoomStretch:
		img.zoomMode = ZoomOriginal
	case ZoomOriginal:
		img.zoomMode = ZoomFit
	}
}

// GetScale returns the current zoom level
func (img *Image) GetScale() float64 {
	return img.zoomLevel
}

// GetZoomModeName returns the current zoom mode name
func (img *Image) GetZoomModeName() string {
	return img.zoomMode.String()
}

// SetZoomMode sets the zoom mode
func (img *Image) SetZoomMode(mode ZoomMode) {
	img.zoomMode = mode
}

// GetZoomMode returns the current zoom mode
func (img *Image) GetZoomMode() ZoomMode {
	return img.zoomMode
}

// SetScale sets the zoom level directly
func (img *Image) SetScale(scale float64) {
	if scale < 0.1 {
		scale = 0.1
	}
	if scale > 5.0 {
		scale = 5.0
	}
	img.zoomLevel = scale
}

// ScaledSize returns the scaled dimensions of the image
// This calculates what size the image should be displayed at, considering zoom level
func (img *Image) ScaledSize() (int, int) {
	if img.imgData == nil {
		return 0, 0
	}

	// Get display bounds (terminal size minus margins)
	displayW, displayH := img.calculateDisplaySize()
	
	// Get original image dimensions
	origW, origH := img.imgData.Width, img.imgData.Height
	
	// Calculate aspect ratio
	aspectRatio := float64(origW) / float64(origH)
	
	// Step 1: Calculate base dimensions that fit in display area
	var baseW, baseH int
	
	switch img.zoomMode {
	case ZoomFit:
		// Fit within display area (maintain aspect ratio)
		baseW = displayW
		baseH = int(float64(baseW) / aspectRatio)
		
		// If height exceeds, fit to height instead
		if baseH > displayH {
			baseH = displayH
			baseW = int(float64(baseH) * aspectRatio)
		}
		
	case ZoomFill:
		// Fill display area (may crop/distort)
		// Use the larger scale to ensure coverage
		scaleX := float64(displayW) / float64(origW)
		scaleY := float64(displayH) / float64(origH)
		baseScale := scaleX
		if scaleY > baseScale {
			baseScale = scaleY
		}
		baseW = int(float64(origW) * baseScale)
		baseH = int(float64(origH) * baseScale)
		
	case ZoomStretch:
		// Stretch to fill display (ignore aspect ratio)
		baseW = displayW
		baseH = displayH
		
	case ZoomOriginal:
		// Use original size (may be larger/smaller than display)
		baseW = origW
		baseH = origH
	}
	
	// Step 2: Apply zoom level by INVERTING the sampling area
	// This changes resolution, not display size
	// Higher zoom = smaller sampled area = higher resolution/detail
	// zoomLevel 2.0 = sample half the area (2x zoom in)
	// zoomLevel 0.5 = sample double the area (2x zoom out)
	sampledW := int(float64(baseW) / img.zoomLevel)
	sampledH := int(float64(baseH) / img.zoomLevel)
	
	// Ensure minimum size (at least 1 pixel from original)
	if sampledW < 1 {
		sampledW = 1
	}
	if sampledH < 1 {
		sampledH = 1
	}
	
	// Don't exceed original image size (can't sample more than we have)
	if sampledW > origW {
		sampledW = origW
	}
	if sampledH > origH {
		sampledH = origH
	}
	
	return sampledW, sampledH
}
