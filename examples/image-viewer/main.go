// Buffer-based image viewer test
// This demonstrates the buffer rendering system with actual image

package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"strings"
	"golang.org/x/term"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-sixel"
	"github.com/nfnt/resize"
	"github.com/wwsheng009/taproot/ui/render/buffer"
)

// Model for the image viewer
type Model struct {
	width       int
	height      int
	termWidth   int
	termHeight  int
	imageData   []byte
	sixelData   []byte
	displayHeight int
}

func initialModel() Model {
	// Load the image file
	imagePath := "test-image.jpeg"
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Printf("Error reading image: %v\n", err)
		os.Exit(1)
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		os.Exit(1)
	}

	// Get image dimensions
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	m := Model{
		termWidth:   80,
		termHeight:  30,
		imageData:   imageData,
		width:       imgWidth,
		height:      imgHeight,
	}

	// Get terminal size
	termW, termH, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		m.termWidth = termW
		m.termHeight = termH
	}

	// Calculate display dimensions (resize to fit terminal while maintaining aspect ratio)
	availableWidth := m.termWidth
	availableHeight := m.termHeight - 6 // Leave room for header and footer

	scaleW := float64(availableWidth) / float64(imgWidth)
	scaleH := float64(availableHeight*6) / float64(imgHeight) // *6 for Sixel height conversion

	scale := min(scaleW, scaleH)

	if scale > 1.0 {
		scale = 1.0 // Don't upscale
	}

	newWidth := int(float64(imgWidth) * scale)
	newHeight := int(float64(imgHeight) * scale)

	// Resize image
	resized := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	// Convert to Sixel
	var sixelBuf bytes.Buffer
	encoder := sixel.NewEncoder(&sixelBuf)
	if err := encoder.Encode(resized); err != nil {
		fmt.Printf("Error encoding to sixel: %v\n", err)
		os.Exit(1)
	}

	m.sixelData = sixelBuf.Bytes()
	m.displayHeight = (newHeight + 5) / 6 // Approximate Sixel display height

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				return m, tea.Quit
			case "+":
				// Zoom in (not implemented for demo)
			case "-":
				// Zoom out (not implemented for demo)
			}
		}
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	// Create layout manager using buffer system
	lm := buffer.NewLayoutManager(m.termWidth, m.termHeight)

	// Calculate layout with display height hint
	lm.ImageLayout(m.displayHeight)

	// Create header component
	header := buffer.NewTextComponent(
		fmt.Sprintf("Image Viewer - %dx%d px | Display: ~%d lines",
			m.width, m.height, m.displayHeight),
		buffer.Style{
			Bold:       true,
			Foreground: "86", // Cyan
			Background: "236",
		},
	).SetCenterH(true)

	// Create footer component
	footer := buffer.NewTextComponent(
		"Controls: Q/ctrl+c quit | + zoom in | - zoom out",
		buffer.Style{
			Foreground: "244",
		},
	).SetCenterH(true)

	// Create image placeholder component (will be replaced by actual Sixel output)
	imagePlaceholder := buffer.NewImageComponent(m.termWidth-10, m.displayHeight*3)

	// Add components to layout
	lm.AddComponent("header", header)
	lm.AddComponent("content", imagePlaceholder)
	lm.AddComponent("footer", footer)

	// Render layout using buffer system
	output := lm.Render()

	// Split output by lines
	lines := strings.Split(output, "\n")

	// Find content area and insert Sixel image
	headerHeight := 5
	contentStart := headerHeight

	// Build final output
	var result strings.Builder
	for i, line := range lines {
		if i == contentStart {
			result.WriteString(line + "\n")
			// Insert Sixel image here
			result.Write(m.sixelData)
			result.WriteByte('\n')
		} else {
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
