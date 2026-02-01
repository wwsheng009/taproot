package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MarkdownTable holds styles for rendering Markdown tables.
type MarkdownTable struct {
	// Header styles
	Header lipgloss.Style

	// Border styles
	Border       lipgloss.Style
	BorderTop    lipgloss.Style
	BorderBottom lipgloss.Style
	BorderLeft   lipgloss.Style
	BorderRight  lipgloss.Style
	BorderRowSep lipgloss.Style // Row separator

	// Cell styles
	Cell lipgloss.Style

	// Alignment
	AlignLeft   lipgloss.Style
	AlignCenter lipgloss.Style
	AlignRight  lipgloss.Style

	// Alternating row colors
	EvenRow  lipgloss.Style
	OddRow   lipgloss.Style
	FirstRow lipgloss.Style // Header row
}

// MarkdownTaskList holds styles for rendering Markdown task lists.
type MarkdownTaskList struct {
	// Checkbox styles
	Checkbox      lipgloss.Style
	CheckboxChecked lipgloss.Style
	CheckboxUnchecked lipgloss.Style

	// Text styles
	TextChecked   lipgloss.Style
	TextUnchecked lipgloss.Style

	// Indicators
	CheckedIndicator   string // Default: "☑"
	UncheckedIndicator string // Default: "☐"
}

// MarkdownLink holds styles for rendering Markdown links.
type MarkdownLink struct {
	// Link styles
	URL       lipgloss.Style
	URLHover  lipgloss.Style
	URLActive lipgloss.Style

	// Link text
	Text lipgloss.Style

	// Reference links
	Reference lipgloss.Style

	// Auto links
	AutoLink lipgloss.Style

	// Email links
	Email lipgloss.Style

	// Show URL in parentheses after text
	ShowURL bool

	// Underline style
	Underline bool
}

// MarkdownImage holds styles for rendering Markdown images.
type MarkdownImage struct {
	// Image placeholder
	Placeholder lipgloss.Style

	// Alt text
	AltText lipgloss.Style

	// Image URL
	URL lipgloss.Style

	// Size indicators
	Width  lipgloss.Style
	Height lipgloss.Style

	// Border
	Border lipgloss.Style

	// Show image info (dimensions, format)
	ShowInfo bool

	// Placeholder icon
	Icon string // Default: "🖼"
}

// RenderTable renders a Markdown table with the given styles.
func RenderTable(headers []string, rows [][]string, styles MarkdownTable, maxWidth int) string {
	if len(headers) == 0 {
		return ""
	}

	var result strings.Builder

	// Calculate column widths
	colCount := len(headers)
	colWidths := make([]int, colCount)

	// Initialize with header widths
	for i, header := range headers {
		colWidths[i] = lipgloss.Width(header)
	}

	// Find maximum width for each column
	for _, row := range rows {
		for i, cell := range row {
			if i < colCount {
				width := lipgloss.Width(cell)
				if width > colWidths[i] {
					colWidths[i] = width
				}
			}
		}
	}

	// Calculate total table width
	totalWidth := 0
	for _, w := range colWidths {
		totalWidth += w
	}
	// Add padding and borders
	totalWidth += (colCount - 1) * 3 // Space between cells: " | "
	totalWidth += 4                  // Side borders and padding

	// Adjust column widths if too wide
	if maxWidth > 0 && totalWidth > maxWidth {
		availableWidth := maxWidth - 4 - (colCount-1)*3
		if availableWidth > 0 {
			avgWidth := availableWidth / colCount
			for i := range colWidths {
				colWidths[i] = min(colWidths[i], avgWidth)
			}
		}
	}

	// Render top border
	if lipgloss.Width(styles.BorderTop.String()) > 0 {
		result.WriteString(styles.BorderTop.Render(renderTableBorder(colWidths)))
	}

	// Render header
	for i, header := range headers {
		cell := renderTableCell(header, colWidths[i], styles.AlignLeft)
		result.WriteString(styles.Header.Render(cell))

		if i < colCount-1 {
			result.WriteString(styles.Border.Render(" │ "))
		}
	}
	result.WriteString("\n")

	// Render header separator
	separator := renderTableSeparator(colWidths)
	result.WriteString(styles.BorderRowSep.Render(separator))
	result.WriteString("\n")

	// Render rows
	for rowIdx, row := range rows {
		rowStyle := styles.OddRow
		if rowIdx%2 == 0 {
			rowStyle = styles.EvenRow
		}

		for colIdx, cell := range row {
			if colIdx >= colCount {
				break
			}

			// Apply row style to cell
			cell = rowStyle.Render(cell)

			renderedCell := renderTableCell(cell, colWidths[colIdx], styles.AlignLeft)
			result.WriteString(styles.Cell.Render(renderedCell))

			if colIdx < colCount-1 {
				result.WriteString(styles.Border.Render(" │ "))
			}
		}
		result.WriteString("\n")
	}

	// Render bottom border
	if lipgloss.Width(styles.BorderBottom.String()) > 0 {
		result.WriteString(styles.BorderBottom.Render(renderTableBorder(colWidths)))
	}

	return result.String()
}

// renderTableCell renders a single table cell with proper padding and truncation.
func renderTableCell(content string, width int, _ lipgloss.Style) string {
	contentWidth := lipgloss.Width(content)

	if contentWidth > width {
		// Truncate
		content = lipgloss.NewStyle().Width(width).MaxWidth(width).Render(content)
		content = content[:width-1] + "…"
	} else if contentWidth < width {
		// Pad
		padding := width - contentWidth
		content = content + strings.Repeat(" ", padding)
	}

	// Add side padding
	return " " + content + " "
}

// renderTableBorder renders a table border line.
func renderTableBorder(colWidths []int) string {
	var parts []string
	for _, w := range colWidths {
		parts = append(parts, strings.Repeat("─", w+2)) // +2 for padding
	}
	return "┌" + strings.Join(parts, "┬") + "┐"
}

// renderTableSeparator renders a table row separator line.
func renderTableSeparator(colWidths []int) string {
	var parts []string
	for _, w := range colWidths {
		parts = append(parts, strings.Repeat("─", w+2)) // +2 for padding
	}
	return "├" + strings.Join(parts, "┼") + "┤"
}

// RenderTaskList renders a Markdown task list with the given styles.
func RenderTaskList(items []TaskItem, styles MarkdownTaskList) string {
	var result strings.Builder

	for _, item := range items {
		// Render checkbox
		var checkbox string
		var textStyle lipgloss.Style

		if item.Checked {
			checkbox = styles.CheckboxChecked.Render(styles.CheckedIndicator)
			textStyle = styles.TextChecked
		} else {
			checkbox = styles.CheckboxUnchecked.Render(styles.UncheckedIndicator)
			textStyle = styles.TextUnchecked
		}

		// Render text
		text := textStyle.Render(item.Text)

		result.WriteString(checkbox)
		result.WriteString(" ")
		result.WriteString(text)
		result.WriteString("\n")
	}

	return result.String()
}

// TaskItem represents a single task list item.
type TaskItem struct {
	Checked bool
	Text    string
}

// RenderLink renders a Markdown link with the given styles.
func RenderLink(text, url string, styles MarkdownLink) string {
	var result strings.Builder

	// Render link text
	linkText := styles.Text.Render(text)
	result.WriteString(linkText)

	// Render URL if requested
	if styles.ShowURL && url != "" {
		urlText := styles.URL.Render("(" + url + ")")
		result.WriteString(" ")
		result.WriteString(urlText)
	}

	return result.String()
}

// RenderImage renders a Markdown image with the given styles.
func RenderImage(alt, url string, styles MarkdownImage) string {
	var result strings.Builder

	icon := styles.Icon
	if icon == "" {
		icon = "🖼"
	}

	// Render icon and alt text
	placeholder := icon + " " + alt
	result.WriteString(styles.Placeholder.Render(placeholder))

	// Render URL info if requested
	if styles.ShowInfo && url != "" {
		info := " [" + url + "]"
		result.WriteString(styles.URL.Render(info))
	}

	return result.String()
}

// ParseTable parses a Markdown table from text.
// Returns headers, rows, and alignment indicators.
func ParseTable(markdown string) ([]string, [][]string, []string) {
	lines := strings.Split(markdown, "\n")
	if len(lines) < 2 {
		return nil, nil, nil
	}

	// Parse header
	headers := parseTableRow(lines[0])

	// Parse alignment separator
	alignments := parseTableAlignment(lines[1])

	// Parse data rows
	var rows [][]string
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "|") {
			break // End of table
		}
		row := parseTableRow(line)
		if len(row) > 0 {
			rows = append(rows, row)
		}
	}

	return headers, rows, alignments
}

// parseTableRow parses a single table row.
func parseTableRow(line string) []string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return nil
	}

	// Remove leading and trailing pipes
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")

	// Split by pipe
	cells := strings.Split(line, "|")

	// Trim whitespace from each cell
	for i := range cells {
		cells[i] = strings.TrimSpace(cells[i])
	}

	return cells
}

// parseTableAlignment parses the alignment separator row.
func parseTableAlignment(line string) []string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return nil
	}

	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")

	cells := strings.Split(line, "|")

	var alignments []string
	for _, cell := range cells {
		cell = strings.TrimSpace(cell)
		if strings.HasPrefix(cell, ":") && strings.HasSuffix(cell, ":") {
			alignments = append(alignments, "center")
		} else if strings.HasSuffix(cell, ":") {
			alignments = append(alignments, "right")
		} else if strings.HasPrefix(cell, ":") {
			alignments = append(alignments, "left")
		} else {
			alignments = append(alignments, "left")
		}
	}

	return alignments
}

// DefaultMarkdownTable returns default table styles.
func DefaultMarkdownTable() MarkdownTable {
	sty := DefaultStyles()

	return MarkdownTable{
		Header:       sty.Base.Bold(true).Foreground(sty.Primary),
		Border:       sty.Base.Foreground(sty.Border),
		BorderTop:    sty.Base.Foreground(sty.Border),
		BorderBottom: sty.Base.Foreground(sty.Border),
		BorderRowSep: sty.Base.Foreground(sty.Border),
		Cell:         sty.Base.Foreground(sty.FgBase),
		AlignLeft:    sty.Base,
		AlignCenter:  sty.Base,
		AlignRight:   sty.Base,
		EvenRow:      sty.Base.Background(lipgloss.Color("#1e1e1e")),
		OddRow:       sty.Base.Background(lipgloss.Color("#252525")),
		FirstRow:     sty.Base.Bold(true),
	}
}

// DefaultMarkdownTaskList returns default task list styles.
func DefaultMarkdownTaskList() MarkdownTaskList {
	sty := DefaultStyles()

	return MarkdownTaskList{
		Checkbox:           sty.Base.Foreground(sty.Primary),
		CheckboxChecked:    sty.Base.Foreground(lipgloss.Color("#4ade80")),
		CheckboxUnchecked:  sty.Base.Foreground(sty.FgMuted),
		TextChecked:        sty.Base.Strikethrough(true).Foreground(sty.FgMuted),
		TextUnchecked:      sty.Base.Foreground(sty.FgBase),
		CheckedIndicator:   "☑",
		UncheckedIndicator: "☐",
	}
}

// DefaultMarkdownLink returns default link styles.
func DefaultMarkdownLink() MarkdownLink {
	sty := DefaultStyles()

	return MarkdownLink{
		URL:        sty.Base.Foreground(sty.Secondary).Italic(true),
		URLHover:   sty.Base.Foreground(sty.Primary).Underline(true),
		URLActive:  sty.Base.Foreground(sty.Primary).Bold(true),
		Text:       sty.Base.Foreground(sty.Blue).Underline(true),
		Reference:  sty.Base.Foreground(sty.FgMuted),
		AutoLink:   sty.Base.Foreground(sty.Blue),
		Email:      sty.Base.Foreground(sty.Secondary),
		ShowURL:    true,
		Underline:  true,
	}
}

// DefaultMarkdownImage returns default image styles.
func DefaultMarkdownImage() MarkdownImage {
	sty := DefaultStyles()

	return MarkdownImage{
		Placeholder: sty.Base.Foreground(sty.FgMuted).Background(sty.BgSubtle).Padding(0, 1),
		AltText:     sty.Base.Foreground(sty.FgBase).Italic(true),
		URL:         sty.Base.Foreground(sty.Secondary).Italic(true),
		ShowInfo:    true,
		Icon:        "🖼",
	}
}
