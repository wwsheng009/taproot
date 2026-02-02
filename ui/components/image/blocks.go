package image

import (
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

// BlocksRenderer uses Unicode block characters for image display
// This works in any terminal that supports Unicode
// Each character cell can show 2 pixels (upper and lower half)
type BlocksRenderer struct {
	data       *decoder.ImageData
	cellWidth  int
	cellHeight int
	useColor   bool
	useASCII   bool // Fallback to ASCII only
	sampledW   int  // Pre-calculated sampled dimensions from original image
	sampledH   int
}

// NewBlocksRenderer creates a new blocks renderer
func NewBlocksRenderer(data *decoder.ImageData) *BlocksRenderer {
	return &BlocksRenderer{
		data:     data,
		useColor: true,
		useASCII: false,
	}
}

// SetCellSize sets the terminal cell dimensions in pixels
func (b *BlocksRenderer) SetCellSize(w, h int) {
	b.cellWidth = w
	b.cellHeight = h
}

// SetColorEnabled enables or disables colored output
func (b *BlocksRenderer) SetColorEnabled(enabled bool) {
	b.useColor = enabled
}

// SetASCIIMode sets whether to use ASCII-only mode
func (b *BlocksRenderer) SetASCIIMode(ascii bool) {
	b.useASCII = ascii
}

// Render returns the image rendered with Unicode blocks or ASCII
func (b *BlocksRenderer) Render(width, height int) string {
	if b.useASCII {
		return b.renderASCII(width, height)
	}
	return b.renderBlocks(width, height)
}

// SetSampledSize sets the source image dimensions for sampling
func (b *BlocksRenderer) SetSampledSize(width, height int) {
	b.sampledW = width
	b.sampledH = height
}

// renderBlocks renders using Unicode block elements (2 pixels per cell)
func (b *BlocksRenderer) renderBlocks(width, height int) string {
	// Use pre-calculated sampled size if available, otherwise calculate
	sampledW, sampledH := b.sampledW, b.sampledH
	if sampledW == 0 || sampledH == 0 {
		// Fallback: calculate from image data
		scaledW, scaledH := b.data.Scale(width, height*2)
		sampledW = scaledW
		sampledH = scaledH
	}

	var lines []string
	for y := 0; y < height; y++ {
		var line strings.Builder
		for x := 0; x < width; x++ {
			// Map grid position to image position
			imgX := (x * sampledW) / width
			imgY := ((y * 2) * sampledH) / (height * 2)

			// Get colors for upper and lower pixels
			upperR, upperG, upperB, _ := b.data.GetPixelColor(imgX, imgY)
			lowerR, lowerG, lowerB, _ := b.data.GetPixelColor(imgX, imgY+1)

			// Always use upper half block (▀) for two-color rendering
			// Upper pixel = foreground, Lower pixel = background
			line.WriteString(b.formatTwoColorCell(upperR, upperG, upperB, lowerR, lowerG, lowerB))
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

// formatTwoColorCell formats a cell with foreground and background colors using direct ANSI codes
func (b *BlocksRenderer) formatTwoColorCell(fgR, fgG, fgB, bgR, bgG, bgB uint8) string {
	if !b.useColor {
		return "▀"
	}

	// Use upper half block character with foreground and background colors
	// Format: ESC[38;2;R;G;Bm (foreground RGB)
	//         ESC[48;2;R;G;Bm (background RGB)
	//         ▀ (upper half block)
	//         ESC[0m (reset)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀\033[0m",
		fgR, fgG, fgB, bgR, bgG, bgB)
}

// renderASCII renders a simple ASCII approximation
func (b *BlocksRenderer) renderASCII(width, height int) string {
	// Use pre-calculated sampled size if available, otherwise calculate
	scaledW, scaledH := b.sampledW, b.sampledH
	if scaledW == 0 || scaledH == 0 {
		// Fallback: calculate from image data
		scaledW, scaledH = b.data.Scale(width, height)
	}

	// Create ASCII art using brightness levels
	asciiChars := " .:-=+*#%@" // From darkest to lightest

	var lines []string
	for y := 0; y < height; y++ {
		var line strings.Builder
		for x := 0; x < width; x++ {
			imgX := (x * scaledW) / width
			imgY := (y * scaledH) / height

			r, g, bVal, _ := b.data.GetPixelColor(imgX, imgY)
			brightness := (int(r) + int(g) + int(bVal)) / 3

			charIdx := (brightness * (len(asciiChars) - 1)) / 255
			if charIdx >= len(asciiChars) {
				charIdx = len(asciiChars) - 1
			}

			line.WriteByte(asciiChars[charIdx])
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

// SupportsTransparency returns whether blocks support transparency
func (b *BlocksRenderer) SupportsTransparency() bool {
	return false
}

// SupportsAnimation returns whether blocks support animation
func (b *BlocksRenderer) SupportsAnimation() bool {
	return false
}

// DetectUnicodeSupport checks if terminal supports Unicode block characters
func DetectUnicodeSupport() bool {
	// Most modern terminals support Unicode
	// This is a safe default
	return true
}

// DetectColorSupport checks if terminal supports 24-bit color
func DetectColorSupport() bool {
	// Check COLORTERM
	colorTerm := env("COLORTERM")
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return true
	}

	// Check WT_SESSION (Windows Terminal)
	if env("WT_SESSION") != "" {
		return true
	}

	// Check TERM
	term := termEnv()
	truecolorTerms := []string{
		"xterm-256color",
		"screen-256color",
		"tmux-256color",
		"kitty",
		"alacritty",
		"wezterm",
	}

	for _, t := range truecolorTerms {
		if strings.Contains(term, t) {
			return true
		}
	}

	return false
}
