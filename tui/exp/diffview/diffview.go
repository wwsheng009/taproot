package diffview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/styles"
)

// Layout type for diff view
type Layout int

const (
	LayoutUnified Layout = iota
	LayoutSplit
)

// DiffLine represents a single line in the diff
type DiffLine struct {
	Type    LineType
	LineNum int
	Content string
}

// LineType represents the type of change
type LineType int

const (
	LineContext LineType = iota
	LineAdded
	LineDeleted
	LineHeader
)

// DiffView represents a simplified diff viewer
type DiffView struct {
	styles      *styles.Styles
	before      string
	after       string
	layout      Layout
	lineNumbers bool
	width       int
	height      int
	xOffset     int
	yOffset     int

	lines      []DiffLine
	totalLines int
}

// New creates a new DiffView
func New() *DiffView {
	s := styles.DefaultStyles()
	return &DiffView{
		styles:      &s,
		layout:      LayoutUnified,
		lineNumbers: true,
	}
}

// Before sets the before content
func (dv *DiffView) Before(content string) *DiffView {
	dv.before = content
	return dv
}

// After sets the after content
func (dv *DiffView) After(content string) *DiffView {
	dv.after = content
	return dv
}

// SetLayout sets the layout type
func (dv *DiffView) SetLayout(layout Layout) *DiffView {
	dv.layout = layout
	return dv
}

// SetLineNumbers enables/disables line numbers
func (dv *DiffView) SetLineNumbers(show bool) *DiffView {
	dv.lineNumbers = show
	return dv
}

// SetSize sets the width and height
func (dv *DiffView) SetSize(width, height int) {
	dv.width = width
	dv.height = height
}

// SetXOffset sets the horizontal scroll offset
func (dv *DiffView) SetXOffset(offset int) {
	dv.xOffset = offset
	if dv.xOffset < 0 {
		dv.xOffset = 0
	}
}

// SetYOffset sets the vertical scroll offset
func (dv *DiffView) SetYOffset(offset int) {
	dv.yOffset = offset
	if dv.yOffset < 0 {
		dv.yOffset = 0
	}
	if dv.yOffset > dv.totalLines-dv.height {
		dv.yOffset = max(0, dv.totalLines-dv.height)
	}
}

// Compute computes the diff
func (dv *DiffView) Compute() {
	dv.lines = dv.computeDiff()
	dv.totalLines = len(dv.lines)
}

// computeDiff computes a simple line-by-line diff
func (dv *DiffView) computeDiff() []DiffLine {
	beforeLines := strings.Split(dv.before, "\n")
	afterLines := strings.Split(dv.after, "\n")

	var result []DiffLine

	// Simple unified diff algorithm
	beforeIdx := 0
	afterIdx := 0

	for beforeIdx < len(beforeLines) || afterIdx < len(afterLines) {
		if beforeIdx >= len(beforeLines) {
			// Remaining lines are additions
			result = append(result, DiffLine{
				Type:    LineAdded,
				LineNum: afterIdx + 1,
				Content: afterLines[afterIdx],
			})
			afterIdx++
			continue
		}

		if afterIdx >= len(afterLines) {
			// Remaining lines are deletions
			result = append(result, DiffLine{
				Type:    LineDeleted,
				LineNum: beforeIdx + 1,
				Content: beforeLines[beforeIdx],
			})
			beforeIdx++
			continue
		}

		beforeLine := strings.TrimRight(beforeLines[beforeIdx], "\r")
		afterLine := strings.TrimRight(afterLines[afterIdx], "\r")

		if beforeLine == afterLine {
			// Lines are equal
			result = append(result, DiffLine{
				Type:    LineContext,
				LineNum: afterIdx + 1,
				Content: afterLine,
			})
			beforeIdx++
			afterIdx++
		} else {
			// Lines differ - this is a simplified approach
			// Mark as deletion and addition
			result = append(result, DiffLine{
				Type:    LineDeleted,
				LineNum: beforeIdx + 1,
				Content: beforeLine,
			})
			result = append(result, DiffLine{
				Type:    LineAdded,
				LineNum: afterIdx + 1,
				Content: afterLine,
			})
			beforeIdx++
			afterIdx++
		}
	}

	return result
}

// Render renders the diff view
func (dv *DiffView) Render() string {
	if dv.lines == nil {
		dv.Compute()
	}

	s := dv.styles

	var result strings.Builder

	// Header - show layout type
	layoutName := "Unified"
	if dv.layout == LayoutSplit {
		layoutName = "Split"
	}
	header := s.Base.Bold(true).Foreground(s.Primary).
		Render(fmt.Sprintf("─ Diff View (%s) ─", layoutName))
	result.WriteString(header + "\n\n")

	// Render based on layout type
	if dv.layout == LayoutSplit {
		result.WriteString(dv.renderSplit())
	} else {
		result.WriteString(dv.renderUnified())
	}

	return result.String()
}

// renderUnified renders unified diff view
func (dv *DiffView) renderUnified() string {
	s := dv.styles
	var result strings.Builder

	// Calculate visible range
	start := dv.yOffset
	end := min(start+dv.height, len(dv.lines))

	for i := start; i < end; i++ {
		line := dv.lines[i]
		result.WriteString(dv.renderLine(line))
	}

	// Footer
	footer := s.Base.Foreground(s.FgMuted).
		Render(fmt.Sprintf("Lines %d-%d of %d | Scroll: ←/→ ↑/↓",
			start+1, end, dv.totalLines))
	result.WriteString("\n" + footer)

	return result.String()
}

// renderSplit renders side-by-side split diff view
func (dv *DiffView) renderSplit() string {
	s := dv.styles
	var result strings.Builder

	// Calculate column widths for split view
	totalWidth := dv.width
	if totalWidth <= 0 {
		totalWidth = 80
	}

	// Reserve space for line numbers and divider
	lineNumWidth := 6
	dividerWidth := 3 // " │ "
	
	// Calculate available code width for each column
	availableWidth := totalWidth - lineNumWidth*2 - dividerWidth
	leftColWidth := availableWidth / 2
	rightColWidth := availableWidth - leftColWidth

	// Calculate visible range
	start := dv.yOffset
	end := min(start+dv.height, len(dv.lines))

	// Render each line
	for i := start; i < end; i++ {
		if i >= len(dv.lines) {
			break
		}

		line := dv.lines[i]
		
		// Split the content into left (before) and right (after) parts
		leftContent, rightContent := dv.splitLineContent(line)
		
		// Build the line
		var lineBuilder strings.Builder

		// Left side (before)
		if dv.lineNumbers {
			lineNum := "      " // Empty for insertions
			if line.Type == LineDeleted || line.Type == LineContext {
				lineNum = fmt.Sprintf("%6d", line.LineNum)
			}
			lineBuilder.WriteString(s.Base.Foreground(s.FgMuted).Render(lineNum))
		}

		leftStyle := dv.getStyleForLineType(line.Type)
		leftContent = leftStyle.Render(leftContent)
		if lipgloss.Width(leftContent) > leftColWidth {
			leftContent = lipgloss.NewStyle().MaxWidth(leftColWidth).Render(leftContent)
		}
		lineBuilder.WriteString(leftContent + "\n")

		// Right side (after)
		if dv.lineNumbers {
			lineNum := "      " // Empty for deletions
			if line.Type == LineAdded || line.Type == LineContext {
				lineNum = fmt.Sprintf("%6d", line.LineNum)
			}
			lineBuilder.WriteString(s.Base.Foreground(s.FgMuted).Render(lineNum))
		}

		rightStyle := dv.getStyleForLineType(line.Type)
		rightContent = rightStyle.Render(rightContent)
		if lipgloss.Width(rightContent) > rightColWidth {
			rightContent = lipgloss.NewStyle().MaxWidth(rightColWidth).Render(rightContent)
		}
		lineBuilder.WriteString(rightContent + "\n")

		result.WriteString(lineBuilder.String())
	}

	// Footer
	footer := s.Base.Foreground(s.FgMuted).
		Render(fmt.Sprintf("Lines %d-%d of %d | Scroll: ←/→ ↑/↓",
			start+1, end, dv.totalLines))
	result.WriteString("\n" + footer)

	return result.String()
}

// splitLineContent splits a diff line into before and after content
func (dv *DiffView) splitLineContent(line DiffLine) (string, string) {
	switch line.Type {
	case LineAdded:
		return "", line.Content
	case LineDeleted:
		return line.Content, ""
	case LineContext:
		return line.Content, line.Content
	case LineHeader:
		return line.Content, line.Content
	default:
		return "", ""
	}
}

// getStyleForLineType returns the style for a given line type
func (dv *DiffView) getStyleForLineType(lineType LineType) lipgloss.Style {
	s := dv.styles
	
	switch lineType {
	case LineAdded:
		return s.Base.Foreground(s.Green).Background(lipgloss.Color("#d8f8dd"))
	case LineDeleted:
		return s.Base.Foreground(s.Error).Background(lipgloss.Color("#ffebe9"))
	case LineContext:
		return s.Base.Foreground(s.FgBase)
	case LineHeader:
		return s.Base.Bold(true).Foreground(s.FgMuted)
	default:
		return s.Base
	}
}

// renderLine renders a single diff line
func (dv *DiffView) renderLine(line DiffLine) string {
	s := dv.styles

	var prefix string
	var style lipgloss.Style

	switch line.Type {
	case LineAdded:
		prefix = "+"
		style = s.Base.Foreground(s.Green)
	case LineDeleted:
		prefix = "-"
		style = s.Base.Foreground(s.Error)
	case LineContext:
		prefix = " "
		style = s.Base.Foreground(s.FgBase)
	case LineHeader:
		prefix = "@"
		style = s.Base.Bold(true).Foreground(s.Warning)
	}

	// Build line content
	var content strings.Builder
	if dv.lineNumbers && line.Type != LineHeader {
		content.WriteString(fmt.Sprintf("%4d", line.LineNum))
		content.WriteString(" | ")
	}

	// Handle horizontal scrolling
	displayContent := line.Content
	if dv.xOffset > 0 {
		runes := []rune(displayContent)
		if dv.xOffset < len(runes) {
			displayContent = string(runes[dv.xOffset:])
		} else {
			displayContent = ""
		}
	}

	content.WriteString(displayContent)

	// Truncate to width
	lineStr := style.Render(prefix + " " + content.String())
	if dv.width > 0 && len(lineStr) > dv.width {
		lineStr = lineStr[:dv.width]
	}

	return lineStr + "\n"
}

// CanScrollDown returns true if can scroll down
func (dv *DiffView) CanScrollDown() bool {
	return dv.yOffset < dv.totalLines-dv.height
}

// CanScrollUp returns true if can scroll up
func (dv *DiffView) CanScrollUp() bool {
	return dv.yOffset > 0
}

// ScrollDown scrolls down by one line
func (dv *DiffView) ScrollDown() {
	dv.SetYOffset(dv.yOffset + 1)
}

// ScrollUp scrolls up by one line
func (dv *DiffView) ScrollUp() {
	dv.SetYOffset(dv.yOffset - 1)
}

// ScrollLeft scrolls left by one column
func (dv *DiffView) ScrollLeft() {
	dv.SetXOffset(dv.xOffset - 4)
}

// ScrollRight scrolls right by one column
func (dv *DiffView) ScrollRight() {
	dv.SetXOffset(dv.xOffset + 4)
}

// ScrollToTop scrolls to the top
func (dv *DiffView) ScrollToTop() {
	dv.yOffset = 0
}

// ScrollToBottom scrolls to the bottom
func (dv *DiffView) ScrollToBottom() {
	dv.yOffset = max(0, dv.totalLines-dv.height)
}
