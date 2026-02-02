package decoder

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

// ImageData holds the decoded image information
type ImageData struct {
	Img        image.Image
	Width      int
	Height     int
	Path       string
	Format     string
	Background uint32 // RGBA background color for transparency
}

// Decoder handles image decoding operations
type Decoder struct {
	background uint32 // Default background color (black)
}

// NewDecoder creates a new image decoder
func NewDecoder() *Decoder {
	return &Decoder{
		background: 0xFF000000, // Opaque black
	}
}

// SetBackground sets the background color for transparent images (RGBA format)
func (d *Decoder) SetBackground(r, g, b, a uint8) {
	d.background = uint32(a)<<24 | uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}

// DecodeFile decodes an image from a file path
func (d *Decoder) DecodeFile(path string) (*ImageData, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return nil, errors.New("invalid image dimensions")
	}

	return &ImageData{
		Img:        img,
		Width:      width,
		Height:     height,
		Path:       path,
		Format:     format,
		Background: d.background,
	}, nil
}

// DetectFormat detects the image format from the file extension
func DetectFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	case ".bmp":
		return "bmp"
	case ".webp":
		return "webp"
	default:
		return "unknown"
	}
}

// ValidatePath checks if the path points to a supported image file
func ValidatePath(path string) error {
	if path == "" {
		return errors.New("empty path")
	}

	// Check file extension
	format := DetectFormat(path)
	if format == "unknown" {
		return fmt.Errorf("unsupported image format: %s", filepath.Ext(path))
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", path)
	}

	return nil
}

// GetRGBA returns the image as RGBA format
func (d *ImageData) GetRGBA() *image.RGBA {
	if rgba, ok := d.Img.(*image.RGBA); ok {
		return rgba
	}
	rgba := image.NewRGBA(d.Img.Bounds())
	// Draw the original image onto the RGBA image
	draw.Draw(rgba, rgba.Bounds(), d.Img, d.Img.Bounds().Min, draw.Src)
	return rgba
}

// Scale calculates scaled dimensions while maintaining aspect ratio
func (d *ImageData) Scale(maxWidth, maxHeight int) (int, int) {
	return d.ScaleWithOptions(maxWidth, maxHeight, false)
}

// ScaleWithOptions calculates scaled dimensions with option to allow upscaling
func (d *ImageData) ScaleWithOptions(maxWidth, maxHeight int, allowUpscale bool) (int, int) {
	if maxWidth <= 0 || maxHeight <= 0 {
		return d.Width, d.Height
	}

	// Calculate aspect ratios
	widthRatio := float64(maxWidth) / float64(d.Width)
	heightRatio := float64(maxHeight) / float64(d.Height)

	// Use the smaller ratio to fit within bounds
	ratio := widthRatio
	if heightRatio < widthRatio {
		ratio = heightRatio
	}

	// Don't upscale unless explicitly allowed
	if !allowUpscale && ratio > 1.0 {
		ratio = 1.0
	}

	newWidth := int(float64(d.Width) * ratio)
	newHeight := int(float64(d.Height) * ratio)

	return newWidth, newHeight
}

// GetPixelColor returns the RGBA color of a pixel at the given position
func (d *ImageData) GetPixelColor(x, y int) (uint8, uint8, uint8, uint8) {
	if x < 0 || x >= d.Width || y < 0 || y >= d.Height {
		return 0, 0, 0, 255
	}

	// Get color from image
	r, g, b, a := d.Img.At(x, y).RGBA()

	// Convert from 16-bit to 8-bit
	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)
}

// GetAverageColor returns the average color of a region
func (d *ImageData) GetAverageColor(x, y, width, height int) (uint8, uint8, uint8, uint8) {
	if width <= 0 || height <= 0 {
		return 0, 0, 0, 255
	}

	var rSum, gSum, bSum, aSum uint64
	count := 0

	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			px := x + dx
			py := y + dy

			if px >= 0 && px < d.Width && py >= 0 && py < d.Height {
				r, g, b, a := d.GetPixelColor(px, py)
				rSum += uint64(r)
				gSum += uint64(g)
				bSum += uint64(b)
				aSum += uint64(a)
				count++
			}
		}
	}

	if count == 0 {
		return 0, 0, 0, 255
	}

	return uint8(rSum / uint64(count)),
		uint8(gSum / uint64(count)),
		uint8(bSum / uint64(count)),
		uint8(aSum / uint64(count))
}

// Bytes returns the approximate memory size in bytes
func (d *ImageData) Bytes() int64 {
	return int64(d.Width * d.Height * 4) // 4 bytes per pixel (RGBA)
}
