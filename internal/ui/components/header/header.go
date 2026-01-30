package header

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/layout"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

var _ headerImpl = (*HeaderComponent)(nil)
var _ layout.Sizeable = (*HeaderComponent)(nil)

type headerImpl interface {
	Size() (width, height int)
	SetSize(width, height int)
	SetBrand(brand, title string)
	SetSessionTitle(title string)
	SetWorkingDirectory(cwd string)
	SetTokenUsage(used, max int, cost float64)
	SetErrorCount(count int)
	SetDetailsOpen(open bool)
	ShowingDetails() bool
}

// HeaderComponent is the header component implementation.
type HeaderComponent struct {
	width        int
	height       int
	brand        string
	title        string
	sessionTitle string
	workingDir   string
	tokenUsed    int
	tokenMax     int
	cost         float64
	errorCount   int
	detailsOpen  bool
	compactMode  bool
}

// New creates a new header component.
func New() *HeaderComponent {
	return &HeaderComponent{
		brand:       "Charm™",
		title:       "CRUSH",
		tokenMax:    128000,
		compactMode: false,
	}
}

// Size returns the header dimensions.
func (h *HeaderComponent) Size() (width, height int) {
	return h.width, h.height
}

// SetSize sets the header dimensions.
func (h *HeaderComponent) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// SetBrand sets the brand and title text.
func (h *HeaderComponent) SetBrand(brand, title string) {
	h.brand = brand
	h.title = title
}

// SetSessionTitle sets the session title.
func (h *HeaderComponent) SetSessionTitle(title string) {
	h.sessionTitle = title
}

// SetWorkingDirectory sets the working directory.
func (h *HeaderComponent) SetWorkingDirectory(cwd string) {
	h.workingDir = cwd
}

// SetTokenUsage sets token usage statistics.
func (h *HeaderComponent) SetTokenUsage(used, max int, cost float64) {
	h.tokenUsed = used
	h.tokenMax = max
	h.cost = cost
}

// SetErrorCount sets the error count display.
func (h *HeaderComponent) SetErrorCount(count int) {
	h.errorCount = count
}

// SetDetailsOpen sets whether details panel is open.
func (h *HeaderComponent) SetDetailsOpen(open bool) {
	h.detailsOpen = open
}

// ShowingDetails returns whether details panel is shown.
func (h *HeaderComponent) ShowingDetails() bool {
	return h.detailsOpen
}

// SetCompactMode sets whether to use compact mode.
func (h *HeaderComponent) SetCompactMode(compact bool) {
	h.compactMode = compact
}

// View renders the header.
func (h *HeaderComponent) View() string {
	if h.brand == "" {
		return ""
	}

	s := styles.DefaultStyles()

	const (
		gap          = " "
		diag         = "╱"
		minDiags     = 3
		leftPadding  = 1
		rightPadding = 1
	)

	var b strings.Builder

	// Render brand and title
	b.WriteString(s.Base.Foreground(s.Secondary).Render(h.brand))
	b.WriteString(gap)
	b.WriteString(styles.ApplyBoldForegroundGrad(&s, h.title, s.Secondary, s.Primary))
	b.WriteString(gap)

	// Calculate available width for content (excluding padding)
	availableWidth := h.width - leftPadding - rightPadding
	progressBarWidth := int(float64(availableWidth) * 0.25)

	// Always render progress bar area with fixed width
	if progressBarWidth > minDiags {
		// Calculate percentage (0 if no token usage)
		var percentage float64
		if h.tokenUsed > 0 && h.tokenMax > 0 {
			percentage = float64(h.tokenUsed) / float64(h.tokenMax)
		}

		// Calculate number of diags based purely on token percentage
		// At 0%, we still show minDiags, not 0
		diagsCount := minDiags + int(float64(progressBarWidth-minDiags)*percentage)

		// Render progress bar with padding to fill fixed width
		diagsStr := strings.Repeat(diag, diagsCount)
		paddingCount := progressBarWidth - diagsCount
		if paddingCount > 0 {
			diagsStr += strings.Repeat(" ", paddingCount)
		}

		b.WriteString(s.Base.Foreground(s.Primary).Render(diagsStr))
		b.WriteString(gap)
	}

	// Calculate remaining width for details
	usedWidth := lipgloss.Width(b.String())
	detailsAvailWidth := availableWidth - usedWidth

	// Render details and pad to fill remaining width
	if detailsAvailWidth > minDiags {
		details := h.renderDetails(detailsAvailWidth)
		detailsWidth := lipgloss.Width(details)
		if detailsWidth < detailsAvailWidth {
			// Pad with space to fill remaining width
			details += strings.Repeat(" ", detailsAvailWidth-detailsWidth)
		}
		b.WriteString(details)
	}

	// Ensure content fits within available width to prevent multi-line rendering
	// We use string manipulation rather than lipgloss.MaxWidth for reliability
	// with ANSI sequences.
	content := b.String()
	contentWidth := lipgloss.Width(content)

	// If content exceeds available width, truncate it to available columns
	if contentWidth > availableWidth {
		// Manually truncate by iterating through runes to preserve ANSI sequences
		var truncated strings.Builder
		currentWidth := 0
		runes := []rune(content)
		i := 0

		for i < len(runes) && currentWidth < availableWidth {
			r := runes[i]

			// If we hit an ANSI escape sequence, include it entirely
			if r == '\x1b' {
				// Find the sequence end
				end := i + 1
				for end < len(runes) && runes[end] != 'm' {
					end++
				}
				if end < len(runes) {
					for j := i; j <= end; j++ {
						truncated.WriteRune(runes[j])
					}
					i = end + 1
					continue
				}
			}

			// Check if this rune adds visible width
			runeWidth := lipgloss.Width(string(r))
			if currentWidth+runeWidth > availableWidth {
				break
			}

			truncated.WriteRune(r)
			currentWidth += runeWidth
			i++
		}

		// If we still exceeded the width (due to ANSI sequences), use lipgloss as fallback
		currentResult := truncated.String()
		if lipgloss.Width(currentResult) > availableWidth {
			// Fallback: use lipgloss with both MaxWidth and MaxHeight
			truncatedStyle := lipgloss.NewStyle().
				MaxWidth(availableWidth).
				MaxHeight(1)
			currentResult = truncatedStyle.Render(content)
		}
		content = currentResult
		contentWidth = lipgloss.Width(content)
	}

	// Pad with spaces to fill exactly availableWidth
	if contentWidth < availableWidth {
		content += strings.Repeat(" ", availableWidth-contentWidth)
	}

	// Apply padding without any width constraints
	// The padding should not introduce any newlines
	result := s.Base.Padding(0, rightPadding, 0, leftPadding).Render(content)

	// CRITICAL: Final safety check - ensure absolutely no line breaks
	if strings.ContainsAny(result, "\n\r") {
		// If somehow line breaks were introduced, remove them as a last resort
		result = strings.ReplaceAll(result, "\n", "")
		result = strings.ReplaceAll(result, "\r", "")

		// Ensure we still have correct width after removing newlines
		resultWidth := lipgloss.Width(result)
		targetWidth := h.width
		if resultWidth < targetWidth {
			result += strings.Repeat(" ", targetWidth-resultWidth)
		}
	}

	return result
}

// renderDetails renders the details section.
func (h *HeaderComponent) renderDetails(availWidth int) string {
	s := styles.DefaultStyles()

	var parts []string

	// Error count
	if h.errorCount > 0 {
		errorStyle := s.Base.Foreground(s.Error)
		parts = append(parts, errorStyle.Render(fmt.Sprintf("%s%d", styles.ErrorIcon, h.errorCount)))
	}

	// Token usage
	var tokenStr string
	if h.tokenMax > 0 {
		percentage := int(float64(h.tokenUsed) / float64(h.tokenMax) * 100)
		tokenStr = fmt.Sprintf("%d%%", percentage)
	} else {
		tokenStr = fmt.Sprintf("%d", h.tokenUsed)
	}
	parts = append(parts, s.Muted.Render(tokenStr))

	// Details toggle
	const keystroke = "ctrl+d"
	if h.detailsOpen {
		parts = append(parts, s.Muted.Render(keystroke)+s.Subtle.Render(" close"))
	} else {
		parts = append(parts, s.Muted.Render(keystroke)+s.Subtle.Render(" open "))
	}

	dot := s.Subtle.Render(" • ")
	metadata := strings.Join(parts, dot)
	metadata = dot + metadata

	// Truncate working directory if necessary
	cwd := h.workingDir
	if cwd == "" {
		cwd = "~"
	}

	// Truncate directory to show max 4 components
	dirs := strings.Split(cwd, string('/'))
	if len(dirs) > 4 {
		cwd = strings.Join(dirs[len(dirs)-4:], "/")
		cwd = "…" + cwd
	}

	cwd = lipgloss.NewStyle().
		MaxWidth(max(0, availWidth-lipgloss.Width(metadata))).
		MaxHeight(1). // Ensure CWD rendering never exceeds 1 line
		Render(cwd)
	cwd = s.Muted.Render(cwd)

	return cwd + metadata
}

// Update updates the header state.
func (h *HeaderComponent) Update(msg any) (*HeaderComponent, any) {
	// Placeholder - engine-agnostic Update method
	return h, nil
}
