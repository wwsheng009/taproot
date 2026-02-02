package image

import (
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

// SixelRenderer implements the Sixel graphics protocol
// Documentation: https://en.wikipedia.org/wiki/Sixel
// Supported by: Windows Terminal, xterm, mlterm, etc.
type SixelRenderer struct {
	data       *decoder.ImageData
	cellWidth  int
	cellHeight int
	colorReg   int // Number of color registers (max 256 for Sixel)
	sampledW   int // Sampled width (from zoom level)
	sampledH   int // Sampled height (from zoom level)
}

// NewSixelRenderer creates a new Sixel renderer
func NewSixelRenderer(data *decoder.ImageData) *SixelRenderer {
	return &SixelRenderer{
		data:     data,
		colorReg: 256, // Maximum for Sixel
	}
}

// SetCellSize sets the terminal cell dimensions in pixels
func (s *SixelRenderer) SetCellSize(w, h int) {
	s.cellWidth = w
	s.cellHeight = h
}

// SetSampledSize sets the sampled dimensions (from zoom level)
func (s *SixelRenderer) SetSampledSize(w, h int) {
	s.sampledW = w
	s.sampledH = h
}

// SetColorRegisters sets the number of color registers (1-256)
func (s *SixelRenderer) SetColorRegisters(n int) {
	if n < 1 {
		n = 1
	}
	if n > 256 {
		n = 256
	}
	s.colorReg = n
}

// Render returns the Sixel escape sequence
func (s *SixelRenderer) Render(width, height int) string {
	// Use sampled dimensions instead of calculating from display size
	// This allows zoom to control the scale independently of display size
	scaledW, scaledH := s.sampledW, s.sampledH

	// Fall back to original scale if sampled size not set or invalid
	if scaledW <= 0 || scaledH <= 0 {
		// Sixel displays 6 pixels vertically per character
		// So we need to adjust our calculations
		displayWidth := width * s.cellWidth
		displayHeight := height * s.cellHeight * 6
		// Scale image
		scaledW, scaledH = s.data.Scale(displayWidth, displayHeight)
	}

	// Generate Sixel data
	sixelData := s.generateSixel(scaledW, scaledH)

	return sixelData
}

// scaleToFit calculates scaled dimensions
func (s *SixelRenderer) scaleToFit(maxW, maxH int) (int, int) {
	return s.data.Scale(maxW, maxH)
}

// generateSixel generates Sixel graphics data from image
func (s *SixelRenderer) generateSixel(width, height int) string {
	var sb strings.Builder

	// Start Sixel mode
	// \033Pq starts Sixel mode with device control string
	// Parameters: <width>;<height>
	sb.WriteString(fmt.Sprintf("\033Pq%d;%d", width, height))

	// Build color palette using simple quantization
	// Limit to 64 colors for compatibility
	maxColors := 64
	if s.colorReg < maxColors {
		maxColors = s.colorReg
	}

	// Collect unique colors and build palette
	palette := make(map[struct{ R, G, B uint8 }]int)
	colors := []struct{ R, G, B, A uint8 }{}
	paletteSize := 0

	// Sample pixels to build palette (step sampling for speed)
	step := 1
	if width > 100 || height > 100 {
		step = 2
	}
	if width > 200 || height > 200 {
		step = 4
	}

	for y := 0; y < height; y += step {
		for x := 0; x < width; x += step {
			r, g, b, a := s.data.GetPixelColor(
				(x*s.data.Width)/width,
				(y*s.data.Height)/height,
			)

			if a < 128 { // Skip transparent pixels
				continue
			}

			key := struct{ R, G, B uint8 }{R: r, G: g, B: b}
			if _, exists := palette[key]; !exists {
				if paletteSize >= maxColors {
					break
				}
				palette[key] = paletteSize
				colors = append(colors, struct{ R, G, B, A uint8 }{R: r, G: g, B: b, A: a})
				paletteSize++
			}
		}
	}

	// Add black as color 0 if not present
	if _, hasBlack := palette[struct{ R, G, B uint8 }{R: 0, G: 0, B: 0}]; !hasBlack {
		palette[struct{ R, G, B uint8 }{R: 0, G: 0, B: 0}] = paletteSize
		colors = append(colors, struct{ R, G, B, A uint8 }{R: 0, G: 0, B: 0, A: 255})
		paletteSize++
	}

	// Define color palette with # Pc ; Pu ; Px ; Py ; Pz format
	// Pc: color number, Pu: pixel type (always 2 for RGB), Px,Py,Pz: RGB (0-100)
	for i, color := range colors {
		if i >= maxColors {
			break
		}
		// Convert 0-255 to 0-100
		r := int(color.R) * 100 / 255
		g := int(color.G) * 100 / 255
		b := int(color.B) * 100 / 255
		sb.WriteString(fmt.Sprintf(";2;%d;%d;%d#%d", r, g, b, i))
	}

	// Default to color 0
	sb.WriteString("#0")

	// Encode image in 6-pixel vertical strips (Sixel rows)
	for bandY := 0; bandY < height; bandY += 6 {
		for x := 0; x < width; x++ {
			// Process 6 pixels vertically at this x position
			var sixelBits uint8 = 0
			var pixelR, pixelG, pixelB uint8
			hasPixel := false

			for bit := 0; bit < 6 && (bandY+bit) < height; bit++ {
				y := bandY + bit
				imgX := (x * s.data.Width) / width
				imgY := (y * s.data.Height) / height

				r, g, b, a := s.data.GetPixelColor(imgX, imgY)

				if a >= 128 { // Non-transparent pixel
					sixelBits |= 1 << bit
					pixelR, pixelG, pixelB = r, g, b
					hasPixel = true
				}
			}

			if hasPixel {
				// Find the color that best matches this pixel
				key := struct{ R, G, B uint8 }{R: pixelR, G: pixelG, B: pixelB}
				if colorIndex, exists := palette[key]; exists {
					// Switch to this color if different from current
					sb.WriteString(fmt.Sprintf("#%d", colorIndex))
				}
				// Output the sixel character (63-126 range)
				sb.WriteByte(63 + sixelBits)
			} else {
				// Empty pixel, output space (63, which is '?')
				sb.WriteByte(63)
			}
		}
		// End line with carriage return (-) and optionally move to next line
		sb.WriteString("$")
	}

	// End Sixel mode with string terminator
	sb.WriteString("\033\\")

	return sb.String()
}

// SupportsTransparency returns whether Sixel supports transparency
func (s *SixelRenderer) SupportsTransparency() bool {
	return false // Sixel doesn't support transparency
}

// SupportsAnimation returns whether Sixel supports animation
func (s *SixelRenderer) SupportsAnimation() bool {
	return false
}

// DetectSixel checks if the terminal supports Sixel
func DetectSixel() bool {
	// Check TERM variable
	term := termEnv()

	// Terminals known to support Sixel
	sixelTerms := []string{
		"xterm",
		"vt340",
		"vt382",
		"vt330",
		"dtterm",
		"ms-terminal", // Windows Terminal
		"windows Terminal",
	}

	for _, t := range sixelTerms {
		if strings.Contains(strings.ToLower(term), t) {
			return true
		}
	}

	// Check WT_SESSION (Windows Terminal)
	if env("WT_SESSION") != "" {
		return true
	}

	// Check VTE version (vte >= 0.52 supports Sixel)
	vteVersion := env("VTE_VERSION")
	if vteVersion != "" {
		// VTE_VERSION is in format: major * 10000 + minor * 100 + patch
		// 0.52 = 5200
		return true // Simplified check
	}

	return false
}
