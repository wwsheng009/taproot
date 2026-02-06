package logo

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/styles"
)

var (
	randCaches   = make(map[int]int)
	randCachesMu sync.Mutex
)

func cachedRandN(n int) int {
	randCachesMu.Lock()
	defer randCachesMu.Unlock()

	if n, ok := randCaches[n]; ok {
		return n
	}

	r := rand.IntN(n)
	randCaches[n] = r
	return r
}

// Opts are the options for rendering the title art.
type Opts struct {
	FieldColor   lipgloss.Color // diagonal lines
	TitleColorA  lipgloss.Color // left gradient ramp point
	TitleColorB  lipgloss.Color // right gradient ramp point
	CharmColor   lipgloss.Color // Charm™ text color
	VersionColor lipgloss.Color // Version text color
	Width        int            // width of the rendered logo, used for truncation
}

const diag = `╱`

// Render renders the logo. Set compact to true for a narrow version.
func Render(s *styles.Styles, version string, compact bool, o Opts) string {
	const charm = " Taproot"

	fg := func(c lipgloss.Color, s string) string {
		return lipgloss.NewStyle().Foreground(c).Render(s)
	}

	// Title.
	const spacing = 1
	letterforms := []letterform{
		letterT,
		letterA,
		letterP,
		letterR,
		letterO1,
		letterO2,
		letterT,
	}
	stretchIndex := -1 // -1 means no stretching.
	if !compact {
		stretchIndex = cachedRandN(len(letterforms))
	}

	taproot := renderWord(spacing, stretchIndex, letterforms...)
	taprootWidth := lipgloss.Width(taproot)
	b := new(strings.Builder)
	for _, r := range strings.Split(taproot, "\n") {
		fmt.Fprintln(b, styles.ApplyForegroundGrad(s, r, o.TitleColorA, o.TitleColorB))
	}
	taproot = b.String()

	// Charm and version.
	metaRowGap := 1
	maxVersionWidth := taprootWidth - lipgloss.Width(charm) - metaRowGap
	if maxVersionWidth < 0 {
		maxVersionWidth = 0
	}
	// Don't truncate version - if it doesn't fit, just let it overflow
	// This ensures the full version string is always visible
	gap := max(0, taprootWidth-lipgloss.Width(charm)-lipgloss.Width(version))
	metaRow := fg(o.CharmColor, charm) + strings.Repeat(" ", gap) + fg(o.VersionColor, version)

	// Join the meta row and big Taproot title.
	taproot = strings.TrimSpace(metaRow + "\n" + taproot)

	// Narrow version.
	if compact {
		field := fg(o.FieldColor, strings.Repeat(diag, taprootWidth))
		return strings.Join([]string{field, field, taproot, field, ""}, "\n")
	}

	fieldHeight := lipgloss.Height(taproot)

	// Left field.
	const leftWidth = 6
	leftFieldRow := fg(o.FieldColor, strings.Repeat(diag, leftWidth))
	leftField := new(strings.Builder)
	for range fieldHeight {
		fmt.Fprintln(leftField, leftFieldRow)
	}

	// Right field.
	rightWidth := max(15, o.Width-taprootWidth-leftWidth-2) // 2 for the gap.
	const stepDownAt = 0
	rightField := new(strings.Builder)
	for i := range fieldHeight {
		width := rightWidth
		if i >= stepDownAt {
			width = rightWidth - (i - stepDownAt)
		}
		if width < 0 {
			width = 0
		}
		fmt.Fprint(rightField, fg(o.FieldColor, strings.Repeat(diag, width)), "\n")
	}

	// Return the wide version.
	const hGap = " "
	logo := lipgloss.JoinHorizontal(lipgloss.Top, leftField.String(), hGap, taproot, hGap, rightField.String())
	if o.Width > 0 {
		// Truncate the logo to the specified width.
		lines := strings.Split(logo, "\n")
		for i, line := range lines {
			if len(line) > o.Width {
				lines[i] = line[:o.Width]
			}
		}
		logo = strings.Join(lines, "\n")
	}
	return logo
}

// SmallRender renders a smaller version of the logo.
func SmallRender(s *styles.Styles, width int) string {
	title := s.Base.Foreground(s.Secondary).Render("Taproot")
	title = fmt.Sprintf("%s %s", title, styles.ApplyBoldForegroundGrad(s, "TUI", s.Secondary, s.Primary))
	remainingWidth := width - lipgloss.Width(title) - 1
	if remainingWidth > 0 {
		lines := strings.Repeat("╱", remainingWidth)
		title = fmt.Sprintf("%s %s", title, s.Base.Foreground(s.Primary).Render(lines))
	}
	return title
}

// letterform represents a letterform. It can be stretched horizontally by
// a given amount via the boolean argument.
type letterform func(bool) string

// renderWord renders letterforms to form a word.
func renderWord(spacing int, stretchIndex int, letterforms ...letterform) string {
	if spacing < 0 {
		spacing = 0
	}

	renderedLetterforms := make([]string, len(letterforms))

	// pick one letter randomly to stretch
	for i, letter := range letterforms {
		renderedLetterforms[i] = letter(i == stretchIndex)
	}

	if spacing > 0 {
		// Add spaces between the letters.
		spacer := strings.Repeat(" ", spacing)
		var result []string
		for i, lf := range renderedLetterforms {
			result = append(result, lf)
			if i < len(renderedLetterforms)-1 {
				result = append(result, spacer)
			}
		}
		return strings.TrimSpace(strings.Join(result, "\n"))
	}

	return strings.TrimSpace(
		lipgloss.JoinHorizontal(lipgloss.Top, renderedLetterforms...),
	)
}

func joinLetterform(letters ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, letters...)
}

// letterformProps defines letterform stretching properties.
type letterformProps struct {
	width      int
	minStretch int
	maxStretch int
	stretch    bool
}

// stretchLetterformPart is a helper function for letter stretching.
func stretchLetterformPart(s string, p letterformProps) string {
	if p.maxStretch < p.minStretch {
		p.minStretch, p.maxStretch = p.maxStretch, p.minStretch
	}
	n := p.width
	if p.stretch {
		n = cachedRandN(p.maxStretch-p.minStretch) + p.minStretch
	}
	parts := make([]string, n)
	for i := range parts {
		parts[i] = s
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

// Letter implementations.

func letterA(stretch bool) string {
	top := "▄▀\n"
	mid := "█ \n"
	bottom := "▀▀\n"
	return joinLetterform(
		stretchLetterformPart(top, letterformProps{
			stretch:    stretch,
			width:      2,
			minStretch: 6,
			maxStretch: 10,
		}),
		mid,
		bottom,
	)
}

func letterP(stretch bool) string {
	left := "█\n█\n"
	right := "▀\n▀\n"
	return joinLetterform(
		left,
		stretchLetterformPart(right, letterformProps{
			stretch:    stretch,
			width:      2,
			minStretch: 5,
			maxStretch: 10,
		}),
	)
}

func letterR(stretch bool) string {
	left := "█\n█\n▀\n"
	center := "▀\n▀\n"
	right := "▄\n▄\n▀\n"
	return joinLetterform(
		left,
		stretchLetterformPart(center, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 7,
			maxStretch: 12,
		}),
		right,
	)
}

func letterO1(stretch bool) string {
	top := "▄▀▀▄\n"
	mid := "█  \n"
	bottom := "▀▀▀\n"
	return joinLetterform(
		stretchLetterformPart(top, letterformProps{
			stretch:    stretch,
			width:      2,
			minStretch: 5,
			maxStretch: 9,
		}),
		mid,
		bottom,
	)
}

func letterO2(stretch bool) string {
	top := " ▄\n"
	mid := "█ \n"
	bottom := "▀\n"
	return joinLetterform(
		top,
		mid,
		bottom,
	)
}

func letterT(stretch bool) string {
	top := "▀▀▀\n"
	mid := "\n█\n"
	return joinLetterform(
		stretchLetterformPart(top, letterformProps{
			stretch:    stretch,
			width:      3,
			minStretch: 7,
			maxStretch: 12,
		}),
		mid,
	)
}
