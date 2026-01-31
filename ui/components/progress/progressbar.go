package progress

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
)

var _ render.Model = (*ProgressBar)(nil)

// ProgressBarStyle defines the visual style configuration for progress bars.
type ProgressBarStyle struct {
	FullBarStyle  lipgloss.Style
	EmptyBarStyle lipgloss.Style
	ShowPercent   bool
	ShowLabel     bool
	Width         int
}

// DefaultProgressBarStyle returns the default progress bar style.
func DefaultProgressBarStyle() *ProgressBarStyle {
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("#7c3aed"))
	return &ProgressBarStyle{
		FullBarStyle:  s,
		EmptyBarStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")),
		ShowPercent:   true,
		ShowLabel:     true,
		Width:         40,
	}
}

// ProgressBar is a progress bar component implementing render.Model.
type ProgressBar struct {
	current     float64
	total       float64
	label       string
	style       *ProgressBarStyle
	initialized bool
}

// NewProgressBar creates a new progress bar component with total value.
func NewProgressBar(total float64) *ProgressBar {
	return &ProgressBar{
		current:     0,
		total:       total,
		label:       "",
		style:       DefaultProgressBarStyle(),
		initialized: false,
	}
}

// Init initializes the progress bar component.
// Implements render.Model interface.
func (pb *ProgressBar) Init() error {
	pb.initialized = true
	return nil
}

// Update handles incoming messages and returns updated state and commands.
// Implements render.Model interface.
func (pb *ProgressBar) Update(msg any) (render.Model, render.Cmd) {
	return pb, render.None()
}

// View returns the string representation for rendering.
// Implements render.Model interface.
func (pb *ProgressBar) View() string {
	if pb.total <= 0 {
		return ""
	}

	width := pb.style.Width
	if width <= 0 {
		width = 40
	}

	percent := pb.current / pb.total
	if percent > 1 {
		percent = 1
	}
	if percent < 0 {
		percent = 0
	}

	fullChars := int(percent * float64(width))
	emptyChars := width - fullChars

	fullBar := strings.Repeat("█", fullChars)
	emptyBar := strings.Repeat("░", emptyChars)

	barStr := pb.style.FullBarStyle.Render(fullBar) + pb.style.EmptyBarStyle.Render(emptyBar)

	var result strings.Builder

	if pb.style.ShowLabel && pb.label != "" {
		result.WriteString(pb.label)
		result.WriteString(" ")
	}

	result.WriteString(barStr)

	if pb.style.ShowPercent {
		percentText := fmt.Sprintf(" %d%%", int(percent*100))
		result.WriteString(pb.style.FullBarStyle.Render(percentText))
	}

	return result.String()
}

// Current returns the current value.
func (pb *ProgressBar) Current() float64 {
	return pb.current
}

// SetCurrent sets the current value.
func (pb *ProgressBar) SetCurrent(current float64) {
	if current < 0 {
		current = 0
	}
	if current > pb.total {
		current = pb.total
	}
	pb.current = current
}

// Add increments the current value.
func (pb *ProgressBar) Add(delta float64) {
	pb.SetCurrent(pb.current + delta)
}

// Total returns the total value.
func (pb *ProgressBar) Total() float64 {
	return pb.total
}

// SetTotal sets the total value.
func (pb *ProgressBar) SetTotal(total float64) {
	if total <= 0 {
		total = 1
	}
	pb.total = total
	if pb.current > total {
		pb.current = total
	}
}

// Label returns the progress bar's label.
func (pb *ProgressBar) Label() string {
	return pb.label
}

// SetLabel sets the progress bar's label.
func (pb *ProgressBar) SetLabel(label string) {
	pb.label = label
}

// Style returns the progress bar's style configuration.
func (pb *ProgressBar) Style() *ProgressBarStyle {
	return pb.style
}

// SetStyle sets the progress bar's style configuration.
func (pb *ProgressBar) SetStyle(style *ProgressBarStyle) {
	pb.style = style
}

// Completed returns whether the progress is complete.
func (pb *ProgressBar) Completed() bool {
	return pb.current >= pb.total
}

// Reset resets the progress to 0.
func (pb *ProgressBar) Reset() {
	pb.current = 0
}

// Percent returns the completion percentage (0-100).
func (pb *ProgressBar) Percent() float64 {
	if pb.total <= 0 {
		return 0
	}
	return (pb.current / pb.total) * 100
}

// Increment increments the progress by 1.
func (pb *ProgressBar) Increment() {
	pb.Add(1)
}
