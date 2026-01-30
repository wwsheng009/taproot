// Package anim provides an animated spinner.
package anim

import (
	"math/rand/v2"
	"strings"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	fps           = 20
	initialChar   = '.'
	labelGap      = " "
	labelGapWidth = 1

	// Periods of ellipsis animation speed in steps.
	//
	// If the FPS is 20 (50 milliseconds) this means that the ellipsis will
	// change every 8 frames (400 milliseconds).
	ellipsisAnimSpeed = 8

	// The maximum amount of time that can pass before a character appears.
	// This is used to create a staggered entrance effect.
	maxBirthOffset = time.Second

	// Number of frames to prerender for the animation. After this number
	// of frames, the animation will loop. This only applies when color
	// cycling is disabled.
	prerenderedFrames = 10

	// Default number of cycling chars.
	defaultNumCyclingChars = 10
)

// Default colors for gradient.
var (
	defaultGradColorA = lipgloss.Color("#ff0000")
	defaultGradColorB = lipgloss.Color("#0000ff")
	defaultLabelColor = lipgloss.Color("#cccccc")
)

var (
	availableRunes = []rune("0123456789abcdefABCDEF~!@#$£€%^&*()+=_")
	ellipsisFrames = []string{".", "..", "...", ""}
)

// Internal ID management. Used during animating to ensure that frame messages
// are received only by spinner components that sent them.
var lastID int64

func nextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

// StepMsg is a message type used to trigger the next step in the animation.
type StepMsg struct{ id int }

// Settings defines settings for the animation.
type Settings struct {
	Size        int
	Label       string
	LabelColor  lipgloss.Color
	GradColorA  lipgloss.Color
	GradColorB  lipgloss.Color
	CycleColors bool
}

// Default settings.
var (
	DefaultSettings = Settings{
		Size:        defaultNumCyclingChars,
		LabelColor:  defaultLabelColor,
		GradColorA:  defaultGradColorA,
		GradColorB:  defaultGradColorB,
		CycleColors: true,
	}
)

// Anim is a Bubbletea component for an animated spinner.
type Anim struct {
	width            int
	cyclingCharWidth int
	label            []string
	labelWidth       int
	labelColor       lipgloss.Color
	startTime        time.Time
	birthOffsets     []time.Duration
	initialFrames    [][]string // frames for the initial characters
	initialized      atomic.Bool
	cyclingFrames    [][]string // frames for the cycling characters
	step             atomic.Int64 // current main frame step
	ellipsisStep     atomic.Int64 // current ellipsis frame step
	ellipsisFrames   []string     // ellipsis animation frames
	id               int
}

// New creates a new Anim instance with the specified width and label.
func New(opts Settings) *Anim {
	a := &Anim{}
	// Validate settings.
	if opts.Size < 1 {
		opts.Size = defaultNumCyclingChars
	}
	if opts.GradColorA == "" {
		opts.GradColorA = defaultGradColorA
	}
	if opts.GradColorB == "" {
		opts.GradColorB = defaultGradColorB
	}
	if opts.LabelColor == "" {
		opts.LabelColor = defaultLabelColor
	}

	a.id = nextID()
	a.startTime = time.Now()
	a.cyclingCharWidth = opts.Size
	a.labelColor = opts.LabelColor

	a.labelWidth = lipgloss.Width(opts.Label)

	// Total width of anim, in cells.
	a.width = opts.Size
	if opts.Label != "" {
		a.width += labelGapWidth + lipgloss.Width(opts.Label)
	}

	// Render the label
	a.renderLabel(opts.Label)

	// Pre-generate gradient.
	var ramp []colorful.Color
	numFrames := prerenderedFrames
	if opts.CycleColors {
		ramp = makeGradientRamp(a.width*3, opts.GradColorA, opts.GradColorB, opts.GradColorA, opts.GradColorB)
		numFrames = a.width * 2
	} else {
		ramp = makeGradientRamp(a.width, opts.GradColorA, opts.GradColorB)
	}

	// Pre-render initial characters.
	a.initialFrames = make([][]string, numFrames)
	offset := 0
	for i := range a.initialFrames {
		a.initialFrames[i] = make([]string, a.width+labelGapWidth+a.labelWidth)
		for j := range a.initialFrames[i] {
			if j+offset >= len(ramp) {
				continue // skip if we run out of colors
			}

			var c lipgloss.Color
			if j <= a.cyclingCharWidth {
				c = lipgloss.Color(ramp[j+offset].Hex())
			} else {
				c = opts.LabelColor
			}

			// Also prerender the initial character with Lip Gloss to avoid
			// processing in the render loop.
			a.initialFrames[i][j] = lipgloss.NewStyle().
				Foreground(c).
				Render(string(initialChar))
		}
		if opts.CycleColors {
			offset++
		}
	}

	// Prerender scrambled rune frames for the animation.
	a.cyclingFrames = make([][]string, numFrames)
	offset = 0
	for i := range a.cyclingFrames {
		a.cyclingFrames[i] = make([]string, a.width)
		for j := range a.cyclingFrames[i] {
			if j+offset >= len(ramp) {
				continue // skip if we run out of colors
			}

			// Also prerender the color with Lip Gloss here to avoid processing
			// in the render loop.
			r := availableRunes[rand.IntN(len(availableRunes))]
			a.cyclingFrames[i][j] = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ramp[j+offset].Hex())).
				Render(string(r))
		}
		if opts.CycleColors {
			offset++
		}
	}

	a.ellipsisFrames = ellipsisFrames

	// Random assign a birth to each character for a staggered entrance effect.
	a.birthOffsets = make([]time.Duration, a.width)
	for i := range a.birthOffsets {
		a.birthOffsets[i] = time.Duration(rand.IntN(int(maxBirthOffset)))
	}

	return a
}

// renderLabel renders the label with the given color.
func (m *Anim) renderLabel(label string) {
	if label == "" {
		m.label = nil
		return
	}
	m.label = make([]string, 0, len(label))
	for _, r := range label {
		m.label = append(m.label, lipgloss.NewStyle().
			Foreground(m.labelColor).
			Render(string(r)))
	}
}

// Init implements tea.Model.
func (m *Anim) Init() tea.Cmd {
	return m.tickCmd()
}

// Update implements tea.Model.
func (m *Anim) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case time.Time:
		m.step.Add(1)
		m.ellipsisStep.Add(1)
		return m, m.tickCmd()
	case StepMsg:
		m.step.Add(1)
		m.ellipsisStep.Add(1)
		return m, m.tickCmd()
	}
	return m, nil
}

// View implements tea.Model.
func (m *Anim) View() string {
	// Determine which frame to show based on time since start
	elapsed := time.Since(m.startTime)
	frameIndex := int(elapsed.Seconds() * float64(fps))

	// Show initial frames during birth period
	cyclingIndex := frameIndex % len(m.cyclingFrames)
	ellipsisIndex := (frameIndex / ellipsisAnimSpeed) % len(m.ellipsisFrames)

	// Pre-allocate with estimated size (width + label + ellipsis + padding)
	estimatedSize := m.width + len(m.label) + 10 + 20 // extra for ANSI codes
	var b strings.Builder
	b.Grow(estimatedSize)

	// Draw cycling characters
	for i := 0; i < m.cyclingCharWidth; i++ {
		if elapsed < m.birthOffsets[i] {
			// Not yet born, draw nothing
			b.WriteString(" ")
		} else {
			// Born, draw cycling character
			if i < len(m.cyclingFrames[cyclingIndex]) {
				b.WriteString(m.cyclingFrames[cyclingIndex][i])
			}
		}
	}

	// Draw label gap and label
	if len(m.label) > 0 {
		b.WriteString(labelGap)
		for _, r := range m.label {
			b.WriteString(r)
		}
	}

	// Draw ellipsis
	b.WriteString(m.ellipsisFrames[ellipsisIndex])

	return b.String()
}

func (m *Anim) tickCmd() tea.Cmd {
	return tea.Tick(time.Second/time.Duration(fps), func(t time.Time) tea.Msg {
		return StepMsg{id: m.id}
	})
}

// makeGradientRamp creates a gradient ramp with the given colors.
func makeGradientRamp(size int, colors ...lipgloss.Color) []colorful.Color {
	if len(colors) < 2 {
		return nil
	}

	stops := make([]colorful.Color, len(colors))
	for i, c := range colors {
		cc, _ := colorful.MakeColor(lipgloss.Color(c))
		stops[i] = cc
	}

	numSegments := len(stops) - 1
	blended := make([]colorful.Color, 0, size)

	// Calculate how many colors each segment should have.
	segmentSizes := make([]int, numSegments)
	baseSize := size / numSegments
	remainder := size % numSegments

	// Distribute the remainder across segments.
	for i := range numSegments {
		segmentSizes[i] = baseSize
		if i < remainder {
			segmentSizes[i]++
		}
	}

	// Generate colors for each segment.
	for i := range numSegments {
		c1 := stops[i]
		c2 := stops[i+1]
		segmentSize := segmentSizes[i]

		for j := range segmentSize {
			var t float64
			if segmentSize > 1 {
				t = float64(j) / float64(segmentSize-1)
			}
			c := c1.BlendHcl(c2, t)
			blended = append(blended, c)
		}
	}

	return blended
}
