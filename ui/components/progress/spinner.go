package progress

import (
	"sync/atomic"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
)

var _ render.Model = (*Spinner)(nil)

// Internal ID management for spinner frame messages.
var lastSpinnerID int64

func nextSpinnerID() int {
	return int(atomic.AddInt64(&lastSpinnerID, 1))
}

// TickMsg is a message type used to trigger the next frame in the spinner animation.
type TickMsg struct{ id int }

// SpinnerType defines the style of spinner.
type SpinnerType int

const (
	SpinnerTypeDots SpinnerType = iota
	SpinnerTypeLine
	SpinnerTypeArrow
	SpinnerTypeMoon
)

// SpinnerStyle defines the visual style configuration for spinners.
type SpinnerStyle struct {
	Color  lipgloss.Color
	Label  string
	Type   SpinnerType
	FPS    int
	Width  int
}

// DefaultSpinnerStyle returns the default spinner style.
func DefaultSpinnerStyle() *SpinnerStyle {
	return &SpinnerStyle{
		Color: lipgloss.Color("#7c3aed"),
		Type:  SpinnerTypeDots,
		FPS:   10,
		Width: 4,
	}
}

// Spinner is an animated spinner component implementing render.Model.
type Spinner struct {
	id           int
	state        int
	style        *SpinnerStyle
	started      bool
	stopped      bool
	initialized  bool
	tickInterval time.Duration
}

// NewSpinner creates a new spinner component.
func NewSpinner() *Spinner {
	return &Spinner{
		id:           nextSpinnerID(),
		state:        0,
		style:        DefaultSpinnerStyle(),
		started:      false,
		stopped:      false,
		initialized:  false,
		tickInterval: time.Second / time.Duration(DefaultSpinnerStyle().FPS),
	}
}

// Init initializes the spinner component.
// Implements render.Model interface.
func (s *Spinner) Init() error {
	s.initialized = true
	s.start()
	return nil
}

// Update handles incoming messages and returns updated state and commands.
// Implements render.Model interface.
func (s *Spinner) Update(msg any) (render.Model, render.Cmd) {
	if s.stopped {
		return s, render.None()
	}

	switch m := msg.(type) {
	case *TickMsg:
		if m.id == s.id {
			s.state++
			if s.state >= len(s.frames()) {
				s.state = 0
			}
		}
	case *render.TickMsg:
		// Handle generic tick messages from the render engine
		s.state++
		if s.state >= len(s.frames()) {
			s.state = 0
		}
	}

	return s, render.None()
}

// View returns the string representation for rendering.
// Implements render.Model interface.
func (s *Spinner) View() string {
	frames := s.frames()
	if s.state >= len(frames) {
		s.state = 0
	}

	frame := frames[s.state]
	style := lipgloss.NewStyle().Foreground(s.style.Color)
	rendered := style.Render(frame)

	if s.style.Label != "" {
		rendered = s.style.Label + " " + rendered
	}

	return rendered
}

// start starts the spinner animation.
func (s *Spinner) start() {
	if s.started {
		return
	}
	s.started = true
}

// Stop stops the spinner animation.
func (s *Spinner) Stop() {
	s.stopped = true
}

// Reset resets the spinner to its initial state.
func (s *Spinner) Reset() {
	s.state = 0
	s.stopped = false
}

// Running returns whether the spinner is running.
func (s *Spinner) Running() bool {
	return s.started && !s.stopped
}

// Style returns the spinner's style configuration.
func (s *Spinner) Style() *SpinnerStyle {
	return s.style
}

// SetStyle sets the spinner's style configuration.
func (s *Spinner) SetStyle(style *SpinnerStyle) {
	s.style = style
	s.tickInterval = time.Second / time.Duration(style.FPS)
}

// Type returns the spinner type.
func (s *Spinner) Type() SpinnerType {
	return s.style.Type
}

// SetType sets the spinner type.
func (s *Spinner) SetType(spinnerType SpinnerType) {
	s.style.Type = spinnerType
}

// Color returns the spinner color.
func (s *Spinner) Color() lipgloss.Color {
	return s.style.Color
}

// SetColor sets the spinner color.
func (s *Spinner) SetColor(color lipgloss.Color) {
	s.style.Color = color
}

// Label returns the spinner's label.
func (s *Spinner) Label() string {
	return s.style.Label
}

// SetLabel sets the spinner's label.
func (s *Spinner) SetLabel(label string) {
	s.style.Label = label
}

// FPS returns the frames per second for the animation.
func (s *Spinner) FPS() int {
	return s.style.FPS
}

// SetFPS sets the frames per second for the animation.
func (s *Spinner) SetFPS(fps int) {
	if fps <= 0 {
		fps = 10
	}
	s.style.FPS = fps
	s.tickInterval = time.Second / time.Duration(fps)
}

// frames returns the frames for the current spinner type.
func (s *Spinner) frames() []string {
	switch s.style.Type {
	case SpinnerTypeDots:
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	case SpinnerTypeLine:
		return []string{"-", "\\", "|", "/"}
	case SpinnerTypeArrow:
		return []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}
	case SpinnerTypeMoon:
		return []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	default:
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}
}
