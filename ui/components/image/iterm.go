package image

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

// ITerm2Renderer implements the iTerm2 inline images protocol
// Documentation: https://iterm2.com/documentation-images.html
type ITerm2Renderer struct {
	data       *decoder.ImageData
	cellWidth  int
	cellHeight int
	sampledW   int // Sampled width (from zoom level)
	sampledH   int // Sampled height (from zoom level)
}

// NewITerm2Renderer creates a new iTerm2 protocol renderer
func NewITerm2Renderer(data *decoder.ImageData) *ITerm2Renderer {
	return &ITerm2Renderer{
		data:       data,
		cellWidth:  10,
		cellHeight: 20,
	}
}

// SetCellSize sets the terminal cell dimensions in pixels
func (i *ITerm2Renderer) SetCellSize(w, h int) {
	i.cellWidth = w
	i.cellHeight = h
}

// SetSampledSize sets the sampled dimensions (from zoom level)
func (i *ITerm2Renderer) SetSampledSize(w, h int) {
	i.sampledW = w
	i.sampledH = h
}

// Render returns the iTerm2 inline image escape sequence
func (i *ITerm2Renderer) Render(width, height int) string {
	// Use sampled dimensions instead of calculating from display size
	// This allows zoom to control the scale independently of display size
	scaledW, scaledH := i.sampledW, i.sampledH

	// Fall back to original scale if sampled size not set or invalid
	if scaledW <= 0 || scaledH <= 0 {
		// Calculate display size
		displayWidth := width * i.cellWidth
		displayHeight := height * i.cellHeight
		// Scale image to fit
		scaledW, scaledH = i.data.Scale(displayWidth, displayHeight)
	}

	// Encode image to base64
	pngData := i.encodePNG(scaledW, scaledH)
	encoded := base64.StdEncoding.EncodeToString(pngData)

	// Build iTerm2 escape sequence
	// Format: ESC ] 1337 ; File = <key=value> : <base64 data> BEL
	var sb strings.Builder

	// Start of escape sequence
	sb.WriteString("\033]1337;File=")

	// Add parameters
	params := []string{
		fmt.Sprintf("width=%d", scaledW),
		fmt.Sprintf("height=%d", scaledH),
		"preserveAspectRatio=0", // We already handle aspect ratio
		"inline=1",              // Display inline
	}

	sb.WriteString(strings.Join(params, ","))

	// Add data separator
	sb.WriteString(":")

	// Add base64 data
	sb.WriteString(encoded)

	// End of escape sequence (BEL character)
	sb.WriteString("\a")

	return sb.String()
}

// scaleToFit calculates scaled dimensions that fit within the bounds
func (i *ITerm2Renderer) scaleToFit(maxW, maxH int) (int, int) {
	return i.data.Scale(maxW, maxH)
}

// encodePNG converts the image to PNG format
func (i *ITerm2Renderer) encodePNG(width, height int) []byte {
	encoder := NewEncoder()
	data, err := encoder.EncodePNGScaled(i.data, width, height)
	if err != nil {
		// Return 1x1 transparent PNG as fallback
		return []byte("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==")
	}
	return data
}

// SupportsTransparency returns whether the renderer supports transparency
func (i *ITerm2Renderer) SupportsTransparency() bool {
	return true
}

// SupportsAnimation returns whether the renderer supports animation
func (i *ITerm2Renderer) SupportsAnimation() bool {
	return false
}

// DetectITerm2 checks if the terminal supports iTerm2 inline images
func DetectITerm2() bool {
	// Check TERM_PROGRAM
	termProgram := env("TERM_PROGRAM")
	if termProgram == "iTerm.app" {
		return true
	}

	// Check TERM
	term := termEnv()
	if strings.Contains(term, "iterm") {
		return true
	}

	// Check ITERM_SESSION_ID
	if env("ITERM_SESSION_ID") != "" {
		return true
	}

	return false
}
