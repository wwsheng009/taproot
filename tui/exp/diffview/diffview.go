package diffview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/highlight"
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
	filename    string
	useSyntaxHighlighting bool
	width       int
	height      int
	xOffset     int
	yOffset     int

	lines      []DiffLine
	totalLines int
	
	// Cache for highlighted content
	highlightCache map[string]string
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

// SetFilename sets the filename for syntax detection
func (dv *DiffView) SetFilename(name string) *DiffView {
	dv.filename = name
	return dv
}

// SetSyntaxHighlighting enables/disables syntax highlighting in split view
func (dv *DiffView) SetSyntaxHighlighting(enabled bool) *DiffView {
	dv.useSyntaxHighlighting = enabled
	dv.highlightCache = make(map[string]string)
	return dv
}

// SetSize sets the width and height
func (dv *DiffView) SetSize(width, height int) {
	dv.width = width
	dv.height = height
}

// SetXOffset sets the horizontal scroll offset
func (dv *DiffView) SetXOffset(offset int) {
	dv.xOffset = max(0, offset)
}

// SetYOffset sets the vertical scroll offset
func (dv *DiffView) SetYOffset(offset int) {
	dv.yOffset = max(0, min(offset, dv.totalLines-dv.height))
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

	// If syntax highlighting is enabled, use highlighted version
	if dv.useSyntaxHighlighting && dv.filename != "" {
		return dv.renderUnifiedWithHighlight()
	}

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

// renderUnifiedWithHighlight renders unified diff view with syntax highlighting
func (dv *DiffView) renderUnifiedWithHighlight() string {
	s := dv.styles
	var result strings.Builder

	// Get syntax-highlighted lines
	beforeLines, afterLines := dv.getHighlightedLines()

	// Reconstruct line indices to map diff lines to highlighted content
	beforeIdx := 0
	afterIdx := 0

	// Calculate visible range
	start := dv.yOffset
	renderedRows := 0
	for i := 0; i < len(dv.lines) && renderedRows < dv.height; i++ {
		line := dv.lines[i]

		// Skip lines before the visible range
		if i < start {
			dv.updateLineIndices(line, &beforeIdx, &afterIdx)
			continue
		}

		// Check if we should stop rendering
		if renderedRows >= dv.height {
			break
		}

		// Get highlighted content for this line
		content := ""

		switch line.Type {
		case LineAdded:
			if afterIdx < len(afterLines) {
				content = afterLines[afterIdx]
				afterIdx++
			}
		case LineDeleted:
			if beforeIdx < len(beforeLines) {
				content = beforeLines[beforeIdx]
				beforeIdx++
			}
		case LineContext:
			if beforeIdx < len(beforeLines) && afterIdx < len(afterLines) {
				content = beforeLines[beforeIdx]
				beforeIdx++
				afterIdx++
			}
		case LineHeader:
			content = line.Content
		}

		// Render the line with diff styling
		result.WriteString(dv.renderHighlightedLine(line, content))
		renderedRows++
	}

	// Footer
	footer := s.Base.Foreground(s.FgMuted).
		Render(fmt.Sprintf("Lines %d-%d of %d | Scroll: ←/→ ↑/↓ (Syntax: %s)",
			start+1, min(start+dv.height, len(dv.lines)), dv.totalLines, dv.filename))
	result.WriteString("\n" + footer)

	return result.String()
}

// updateLineIndices updates before/after line indices for a diff line
func (dv *DiffView) updateLineIndices(line DiffLine, beforeIdx, afterIdx *int) {
	switch line.Type {
	case LineAdded:
		(*afterIdx)++
	case LineDeleted:
		(*beforeIdx)++
	case LineContext:
		(*beforeIdx)++
		(*afterIdx)++
	case LineHeader:
		// Headers don't affect line indices
	}
}

// renderHighlightedLine renders a diff line with syntax highlighting and diff styling
func (dv *DiffView) renderHighlightedLine(line DiffLine, highlightedContent string) string {
	s := dv.styles
	var lineBuilder strings.Builder

	// Line marker
	var marker string
	switch line.Type {
	case LineAdded:
		marker = "+ "
	case LineDeleted:
		marker = "- "
	case LineContext:
		marker = "  "
	case LineHeader:
		marker = "@@"
	}

	// Get the diff base style for this line type
	diffStyle := dv.getStyleForLineType(line.Type)

	// Render marker with diff style
	markerStyled := diffStyle.Bold(true).Render(marker)

	// Render highlighted content (already has syntax highlighting colors)
	// Apply diff background color while preserving syntax colors
	var contentStyled string
	if highlightedContent != "" {
		contentStyled = diffStyle.Render(highlightedContent)
	} else if line.Content != "" {
		contentStyled = diffStyle.Render(line.Content)
	} else {
		contentStyled = ""
	}

	// Line number
	lineNumStr := ""
	if dv.lineNumbers && line.LineNum > 0 {
		lineNumStr = fmt.Sprintf("%5d ", line.LineNum)
		lineNumStyled := s.Base.Foreground(s.FgMuted).Render(lineNumStr)
		lineBuilder.WriteString(lineNumStyled)
	}

	lineBuilder.WriteString(markerStyled)
	lineBuilder.WriteString(contentStyled)
	lineBuilder.WriteString("\n")

	return lineBuilder.String()
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
	dividerWidth := 3

	// Calculate available code width for each column
	availableWidth := totalWidth - lineNumWidth*2 - dividerWidth
	leftColWidth := availableWidth / 2
	rightColWidth := availableWidth - leftColWidth

	// Calculate visible range (divide by 2 since each diff line takes 2 rows in split view)
	start := dv.yOffset
	visibleLines := dv.height

	// Get syntax-highlighted lines if enabled
	beforeLines, afterLines := dv.getHighlightedLines()

	// Track line indices for both sides
	beforeLineIdx := 0
	afterLineIdx := 0

	// Render each line pair
	i := start
	renderedRows := 0
	for renderedRows < visibleLines && i < len(dv.lines) {
		if i >= len(dv.lines) {
			break
		}

		line := dv.lines[i]
		i++

		// Get left and right content with syntax highlighting
		leftContent, rightContent := "", ""
		if dv.useSyntaxHighlighting && dv.filename != "" {
			leftContent, rightContent = dv.splitLineContentWithHighlight(line, &beforeLineIdx, &afterLineIdx, beforeLines, afterLines)
		} else {
			leftContent, rightContent = dv.splitLineContent(line)
		}

		// Get line numbers for both sides
		leftLineNum := ""
		rightLineNum := ""

		if dv.lineNumbers {
			if line.Type == LineDeleted || line.Type == LineContext {
				// For deleted or context lines, use the line number from the diff
				leftLineNum = fmt.Sprintf("%5d ", line.LineNum)
			} else {
				leftLineNum = "      "
			}
			if line.Type == LineAdded || line.Type == LineContext {
				// For added or context lines, use the line number from the diff
				rightLineNum = fmt.Sprintf("%5d ", line.LineNum)
			} else {
				rightLineNum = "      "
			}
		}

		leftStyle := dv.getStyleForLineType(line.Type)
		rightStyle := dv.getStyleForLineType(line.Type)

		leftStyled := leftStyle.Render(leftContent)
		rightStyled := rightStyle.Render(rightContent)

		// Handle horizontal scrolling by applying offset
		if dv.xOffset > 0 {
			leftRunes := []rune(leftStyled)
			rightRunes := []rune(rightStyled)
			if dv.xOffset < len(leftRunes) {
				leftStyled = string(leftRunes[dv.xOffset:])
			} else {
				leftStyled = ""
			}
			if dv.xOffset < len(rightRunes) {
				rightStyled = string(rightRunes[dv.xOffset:])
			} else {
				rightStyled = ""
			}
		}

		if lipgloss.Width(leftStyled) > leftColWidth {
			leftStyled = lipgloss.NewStyle().MaxWidth(leftColWidth).Render(leftStyled)
		}
		if lipgloss.Width(rightStyled) > rightColWidth {
			rightStyled = lipgloss.NewStyle().MaxWidth(rightColWidth).Render(rightStyled)
		}

		leftPad := strings.Repeat(" ", leftColWidth-lipgloss.Width(leftStyled))
		rightPad := strings.Repeat(" ", rightColWidth-lipgloss.Width(rightStyled))

		divider := s.Base.Foreground(s.Border).Render(" │ ")

		lineBuilder := strings.Builder{}
		lineBuilder.WriteString(s.Base.Foreground(s.FgMuted).Render(leftLineNum))
		lineBuilder.WriteString(leftStyled + leftPad)
		lineBuilder.WriteString(divider)
		lineBuilder.WriteString(s.Base.Foreground(s.FgMuted).Render(rightLineNum))
		lineBuilder.WriteString(rightStyled + rightPad)

		result.WriteString(lineBuilder.String() + "\n")
		renderedRows++
	}

	// Footer
	end := start + renderedRows
	syntaxStatus := ""
	if dv.useSyntaxHighlighting && dv.filename != "" {
		syntaxStatus = fmt.Sprintf(" | Syntax: %s", dv.filename)
	}
	footer := s.Base.Foreground(s.FgMuted).
		Render(fmt.Sprintf("Lines %d-%d of %d%s | Scroll: ←/→ ↑/↓",
			start+1, end, dv.totalLines, syntaxStatus))
	result.WriteString("\n" + footer)

	return result.String()
}

// getHighlightedLines returns syntax-highlighted lines for the before/after content
func (dv *DiffView) getHighlightedLines() (beforeLines, afterLines []string) {
	if !dv.useSyntaxHighlighting || dv.filename == "" {
		return strings.Split(dv.before, "\n"), strings.Split(dv.after, "\n")
	}

	// Check cache
	beforeKey := "before:" + dv.filename
	afterKey := "after:" + dv.filename

	if cached, ok := dv.highlightCache[beforeKey]; ok {
		beforeLines = strings.Split(cached, "\n")
	} else {
		highlighted, err := highlight.SyntaxHighlight(dv.styles, dv.before, dv.filename, dv.styles.BgSubtle)
		if err == nil {
			dv.highlightCache[beforeKey] = highlighted
			beforeLines = strings.Split(highlighted, "\n")
		} else {
			beforeLines = strings.Split(dv.before, "\n")
		}
	}

	if cached, ok := dv.highlightCache[afterKey]; ok {
		afterLines = strings.Split(cached, "\n")
	} else {
		highlighted, err := highlight.SyntaxHighlight(dv.styles, dv.after, dv.filename, dv.styles.BgSubtle)
		if err == nil {
			dv.highlightCache[afterKey] = highlighted
			afterLines = strings.Split(highlighted, "\n")
		} else {
			afterLines = strings.Split(dv.after, "\n")
		}
	}

	return beforeLines, afterLines
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

// splitLineContentWithHighlight splits a diff line using syntax-highlighted content
func (dv *DiffView) splitLineContentWithHighlight(line DiffLine, beforeLineIndex, afterLineIndex *int, beforeLines, afterLines []string) (string, string) {
	// If syntax highlighting is disabled, fall back to original method
	if !dv.useSyntaxHighlighting || dv.filename == "" {
		return dv.splitLineContent(line)
	}

	var leftContent, rightContent string

	switch line.Type {
	case LineAdded:
		// Only right side, increment after line index
		if *afterLineIndex < len(afterLines) {
			rightContent = afterLines[*afterLineIndex]
			(*afterLineIndex)++
		}
		leftContent = ""
	case LineDeleted:
		// Only left side, increment before line index
		if *beforeLineIndex < len(beforeLines) {
			leftContent = beforeLines[*beforeLineIndex]
			(*beforeLineIndex)++
		}
		rightContent = ""
	case LineContext:
		// Both sides, increment both indices
		if *beforeLineIndex < len(beforeLines) {
			leftContent = beforeLines[*beforeLineIndex]
			(*beforeLineIndex)++
		}
		if *afterLineIndex < len(afterLines) {
			rightContent = afterLines[*afterLineIndex]
			(*afterLineIndex)++
		}
	case LineHeader:
		// Headers have no line indices, use content as-is
		leftContent = line.Content
		rightContent = line.Content
	default:
		leftContent = ""
		rightContent = ""
	}

	return leftContent, rightContent
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
		fmt.Fprintf(&content, "%4d", line.LineNum)
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
