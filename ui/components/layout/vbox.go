package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// countLinesInBuilder is a helper to count lines in a string
func countLinesInBuilder(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	if len(s) > 0 && s[len(s)-1] != '\n' {
		count++
	}
	return count
}

// VerticalLayout represents a vertical layout with fixed header and footer,
// and flexible content area in the middle.
type VerticalLayout struct {
	header       string
	footer       string
	content      string
	height       int
	width        int
	centerV      bool
	centerH      bool
	separator    bool
	truncate     bool     // Whether to truncate overflowing content
	maxHeight    int      // Maximum content height (0 = use available space)
}

// NewVerticalLayout creates a new vertical layout container.
func NewVerticalLayout() *VerticalLayout {
	return &VerticalLayout{
		centerV:   true,
		centerH:   true,
		separator: true,
	}
}

// SetHeader sets the header content.
func (l *VerticalLayout) SetHeader(header string) *VerticalLayout {
	l.header = header
	return l
}

// SetFooter sets the footer content.
func (l *VerticalLayout) SetFooter(footer string) *VerticalLayout {
	l.footer = footer
	return l
}

// SetContent sets the content (middle section).
func (l *VerticalLayout) SetContent(content string) *VerticalLayout {
	l.content = content
	return l
}

// SetSize sets the total size of the layout.
func (l *VerticalLayout) SetSize(width, height int) *VerticalLayout {
	l.width = width
	l.height = height
	return l
}

// SetCenterV sets whether to vertically center the content.
func (l *VerticalLayout) SetCenterV(center bool) *VerticalLayout {
	l.centerV = center
	return l
}

// SetCenterH sets whether to horizontally center the content.
func (l *VerticalLayout) SetCenterH(center bool) *VerticalLayout {
	l.centerH = center
	return l
}

// SetSeparator sets whether to draw separator lines between sections.
func (l *VerticalLayout) SetSeparator(show bool) *VerticalLayout {
	l.separator = show
	return l
}

// SetMaxHeight sets the maximum height for content.
// If content exceeds this, it will be truncated (if truncate is enabled).
// Use 0 to use all available space.
func (l *VerticalLayout) SetMaxHeight(height int) *VerticalLayout {
	l.maxHeight = height
	return l
}

// SetTruncate sets whether to truncate content that exceeds available space.
func (l *VerticalLayout) SetTruncate(enable bool) *VerticalLayout {
	l.truncate = enable
	return l
}

// GetHeaderHeight returns the number of lines in the header.
func (l *VerticalLayout) GetHeaderHeight() int {
	if l.header == "" {
		return 0
	}
	return strings.Count(l.header, "\n") + 1
}

// GetFooterHeight returns the number of lines in the footer.
func (l *VerticalLayout) GetFooterHeight() int {
	if l.footer == "" {
		return 0
	}
	return strings.Count(l.footer, "\n") + 1
}

// GetContentHeight returns the calculated height available for content.
// If maxHeight is set, it returns the smaller of available space and maxHeight.
func (l *VerticalLayout) GetContentHeight() int {
	headerHeight := l.GetHeaderHeight()
	footerHeight := l.GetFooterHeight()
	contentHeight := l.height - headerHeight - footerHeight
	if contentHeight < 0 {
		return 0
	}

	// Respect maxHeight if set
	if l.maxHeight > 0 && l.maxHeight < contentHeight {
		return l.maxHeight
	}

	return contentHeight
}

// Render renders the complete layout.
func (l *VerticalLayout) Render(displayHeight int) string {
	var builder strings.Builder

	// Calculate fixed heights
	headerHeight := l.GetHeaderHeight()
	footerHeight := l.GetFooterHeight()

	// Calculate available space for content (fixed, not affected by content size)
	contentSpace := l.height - headerHeight - footerHeight
	if contentSpace < 0 {
		contentSpace = 0
	}

	// Render header (fixed at top)
	if l.header != "" {
		builder.WriteString(l.header)

		// Add separator line if enabled and there's space for content
		if l.separator && contentSpace > 0 {
			sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
			sepWidth := 80
			if l.width > 0 && l.width < 80 {
				sepWidth = l.width
			}
			builder.WriteString("\n")
			builder.WriteString(sepStyle.Render(strings.Repeat("─", sepWidth)))
			builder.WriteString("\n")
		}
	}

	// Prepare content - truncate to fit available space
	contentLines := strings.Split(l.content, "\n")
	renderContent := l.content

	// Calculate content height based on renderer
	actualContentHeight := len(contentLines)
	if displayHeight > 0 && displayHeight > actualContentHeight {
		actualContentHeight = displayHeight
	}

	// Respect maxHeight if set
	maxAllowedHeight := contentSpace
	if l.maxHeight > 0 && l.maxHeight < maxAllowedHeight {
		maxAllowedHeight = l.maxHeight
	}

	// Clip content to available space (this is critical!)
	if actualContentHeight > maxAllowedHeight && maxAllowedHeight > 0 {
		// For text content, truncate lines
		if displayHeight == 0 || l.truncate {
			if maxAllowedHeight < len(contentLines) {
				renderContent = strings.Join(contentLines[:maxAllowedHeight], "\n")
			}
		}
		actualContentHeight = maxAllowedHeight
	}

	// Calculate space for padding to center content
	var paddingTop int
	if l.centerV && actualContentHeight < contentSpace {
		paddingTop = (contentSpace - actualContentHeight) / 2
	}

	// Calculate space needed after content for footer
	// We want: header + paddingTop + content + paddingBottom + footer = height
	spaceUsedForCentering := paddingTop + (contentSpace - actualContentHeight - paddingTop)
	totalContentSpace := headerHeight + actualContentHeight + spaceUsedForCentering
	remainingLines := l.height - totalContentSpace - footerHeight
	if remainingLines < 0 {
		remainingLines = 0
	}

	// Render padding before content
	for i := 0; i < paddingTop; i++ {
		builder.WriteString("\n")
	}

	// Render content
	if l.centerH {
		// For displayHeight > 0 (Sixel images), avoid extending content with width/centering
		if displayHeight > 0 {
			// Don't use centered rendering for graphics - raw content only
			builder.WriteString(renderContent)
		} else {
			// Use centered rendering for text content
			contentStyle := lipgloss.NewStyle().Width(l.width).Align(lipgloss.Center)
			builder.WriteString(contentStyle.Render(renderContent))
		}
	} else {
		builder.WriteString(renderContent)
	}

	// Fill remaining space to push footer to bottom
	for i := 0; i < remainingLines; i++ {
		builder.WriteString("\n")
	}

	// Render footer (fixed at bottom)
	if l.footer != "" {
		if l.separator {
			sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
			sepWidth := 80
			if l.width > 0 && l.width < 80 {
				sepWidth = l.width
			}
			builder.WriteString(sepStyle.Render(strings.Repeat("─", sepWidth)))
			builder.WriteString("\n\n")
		} else {
			builder.WriteString("\n")
		}
		builder.WriteString(l.footer)
	}

	return builder.String()
}
