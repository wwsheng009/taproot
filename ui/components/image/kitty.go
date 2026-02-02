package image

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

// KittyRenderer implements the Kitty graphics protocol
// Documentation: https://sw.kovidgoyal.net/kitty/graphics-protocol/
type KittyRenderer struct {
	data       *decoder.ImageData
	cellWidth  int
	cellHeight int
	sampledW   int // Sampled width (from zoom level)
	sampledH   int // Sampled height (from zoom level)
}

// NewKittyRenderer creates a new Kitty protocol renderer
func NewKittyRenderer(data *decoder.ImageData) *KittyRenderer {
	return &KittyRenderer{
		data:       data,
		cellWidth:  10, // Default cell width (will be updated)
		cellHeight: 20, // Default cell height (will be updated)
	}
}

// SetCellSize sets the terminal cell dimensions in pixels
func (k *KittyRenderer) SetCellSize(w, h int) {
	k.cellWidth = w
	k.cellHeight = h
}

// SetSampledSize sets the sampled dimensions (from zoom level)
func (k *KittyRenderer) SetSampledSize(w, h int) {
	k.sampledW = w
	k.sampledH = h
}

// Render returns the Kitty graphics protocol escape sequence
func (k *KittyRenderer) Render(width, height int) string {
	// Use sampled dimensions instead of calculating from display size
	// This allows zoom to control the scale independently of display size
	scaledW, scaledH := k.sampledW, k.sampledH

	// Fall back to original images if sampled size not set or invalid
	if scaledW <= 0 || scaledH <= 0 {
		// Calculate the display size in pixels
		displayWidth := width * k.cellWidth
		displayHeight := height * k.cellHeight
		// Scale the image to fit
		scaledW, scaledH = k.data.Scale(displayWidth, displayHeight)
	}

	// Generate PNG data and encode to base64
	pngData := k.encodePNG(scaledW, scaledH)
	encoded := base64.StdEncoding.EncodeToString(pngData)

	// Chunk the base64 data (Kitty has limits on escape sequence length)
	chunkSize := 4096
	var chunks []string
	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunks = append(chunks, encoded[i:end])
	}

	var result strings.Builder

	// Send the image with the first chunk
	result.WriteString(fmt.Sprintf("\033_Ga=T,f=100,t=f,m=1;%s\033\\", chunks[0]))

	// Send remaining chunks
	for i := 1; i < len(chunks); i++ {
		last := (i == len(chunks)-1)
		m := 0
		if last {
			m = 0 // Last chunk
		} else {
			m = 1 // More chunks coming
		}
		result.WriteString(fmt.Sprintf("\033_Gm=%d;%s\033\\", m, chunks[i]))
	}

	// Display the image at cursor position
	result.WriteString(fmt.Sprintf("\033_Ga=d,d=c,q=2,i=1,C=1,c=%d,r=%d\033\\", width, height))

	return result.String()
}

// scaleToFit calculates scaled dimensions that fit within the bounds
func (k *KittyRenderer) scaleToFit(maxW, maxH int) (int, int) {
	return k.data.Scale(maxW, maxH)
}

// encodePNG converts the image to PNG format
func (k *KittyRenderer) encodePNG(width, height int) []byte {
	encoder := NewEncoder()
	data, err := encoder.EncodePNGScaled(k.data, width, height)
	if err != nil {
		// Return 1x1 transparent PNG as fallback
		return []byte("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==")
	}
	return data
}

// Delete deletes the displayed image from memory
func (k *KittyRenderer) Delete(imageID int) string {
	return fmt.Sprintf("\033_Ga=d,i=%d\033\\", imageID)
}

// DeleteAll deletes all displayed images
func (k *KittyRenderer) DeleteAll() string {
	return "\033_Ga=d,d=A\033\\"
}

// SupportsTransparency returns whether the renderer supports transparency
func (k *KittyRenderer) SupportsTransparency() bool {
	return true
}

// SupportsAnimation returns whether the renderer supports animation
func (k *KittyRenderer) SupportsAnimation() bool {
	return true
}

// DetectKitty checks if the terminal supports Kitty graphics protocol
func DetectKitty() bool {
	// Check TERM environment variable
	term := termEnv()
	if strings.Contains(term, "kitty") {
		return true
	}

	// Check TERM_PROGRAM
	termProgram := env("TERM_PROGRAM")
	if strings.Contains(termProgram, "kitty") {
		return true
	}

	// Check KITTY_WINDOW_ID
	if env("KITTY_WINDOW_ID") != "" {
		return true
	}

	return false
}
