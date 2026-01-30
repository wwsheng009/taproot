package diffview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/styles"
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

	// Header
	header := s.Base.Bold(true).Foreground(s.Primary).
		Render("─ Diff View (Unified) ─")
	result.WriteString(header + "\n\n")

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
