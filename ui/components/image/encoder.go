package image

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/wwsheng009/taproot/ui/components/image/decoder"
)

// Encoder handles image encoding operations
type Encoder struct{}

// NewEncoder creates a new image encoder
func NewEncoder() *Encoder {
	return &Encoder{}
}

// EncodePNG encodes an image to PNG format
func (e *Encoder) EncodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodePNGScaled encodes and scales an image to PNG format
func (e *Encoder) EncodePNGScaled(data *decoder.ImageData, width, height int) ([]byte, error) {
	// Create scaled image
	scaled := e.scaleImage(data, width, height)
	return e.EncodePNG(scaled)
}

// scaleImage scales an image to the specified dimensions
func (e *Encoder) scaleImage(data *decoder.ImageData, width, height int) image.Image {
	// Create a new RGBA image for the scaled result
	scaled := image.NewRGBA(image.Rect(0, 0, width, height))

	// Simple nearest-neighbor scaling
	srcW := data.Width
	srcH := data.Height

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Map destination pixel to source pixel
			srcX := (x * srcW) / width
			srcY := (y * srcH) / height

			// Get pixel color from source
			r, g, b, a := data.GetPixelColor(srcX, srcY)

			// Set pixel in destination
			scaled.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return scaled
}

// ImageRect represents a rectangle
type ImageRect struct {
	X, Y, Width, Height int
}

// image.Rect is a helper to create a rectangle
func imageRect(x, y, w, h int) ImageRect {
	return ImageRect{X: x, Y: y, Width: w, Height: h}
}

// ImageToRGBA converts any image.Image to *image.RGBA
func ImageToRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba
}

// EncodePNGReader encodes an image to PNG format and returns an io.Reader
func (e *Encoder) EncodePNGReader(img image.Image) (io.Reader, error) {
	data, err := e.EncodePNG(img)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
