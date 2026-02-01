package markdown

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/wwsheng009/taproot/ui/styles"
)

// RenderOptions contains options for markdown rendering.
type RenderOptions struct {
	// Width is the maximum width for rendering.
	Width int

	// Plain enables plain/minimal styling (for thinking content, etc).
	Plain bool

	// PreserveNewlines preserves newlines in the output.
	PreserveNewlines bool

	// TrimNewlines trims trailing newlines from the output.
	TrimNewlines bool

	// EnableTables enables custom table rendering (custom styles, borders, etc).
	EnableTables bool

	// EnableTaskLists enables custom task list rendering (custom icons, colors, etc).
	EnableTaskLists bool
}

// DefaultRenderOptions returns the default render options.
func DefaultRenderOptions() RenderOptions {
	return RenderOptions{
		Width:            80,
		Plain:            false,
		PreserveNewlines: false,
		TrimNewlines:     true,
		EnableTables:     true,
		EnableTaskLists:  true,
	}
}

// Render renders markdown content with styling.
func Render(content string, opts RenderOptions) (string, error) {
	sty := styles.DefaultStyles()
	return RenderWithStyles(content, &sty, opts)
}

// RenderWithStyles renders markdown content with custom styles.
func RenderWithStyles(content string, sty *styles.Styles, opts RenderOptions) (string, error) {
	// Preprocess content if custom rendering is enabled
	processedContent := content
	if opts.EnableTables {
		applyCustomTableMarkers(&processedContent)
	}
	if opts.EnableTaskLists {
		applyCustomTaskListMarkers(&processedContent)
	}

	// Get the appropriate style config
	var styleConfig glamour.TermRendererOption
	if opts.Plain {
		styleConfig = glamour.WithStyles(sty.PlainMarkdown)
	} else {
		styleConfig = glamour.WithStyles(sty.Markdown)
	}

	// Create renderer
	renderer, err := glamour.NewTermRenderer(
		styleConfig,
		glamour.WithWordWrap(opts.Width),
	)
	if err != nil {
		return "", err
	}

	// Render content
	result, err := renderer.Render(processedContent)
	if err != nil {
		return "", err
	}

	// Post-process
	if opts.TrimNewlines {
		result = strings.TrimSuffix(result, "\n")
	}
	if !opts.PreserveNewlines {
		// Glamour adds extra newlines, clean them up
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	return result, nil
}

// applyCustomTableMarkers adds special markers for custom table rendering.
func applyCustomTableMarkers(content *string) {
	// Parse tables and enhance alignment beyond glamour's defaults
	// This handles custom alignment that glamour doesn't support well
	*content = enhanceTableAlignment(*content)
}

// enhanceTableAlignment improves table alignment with custom column widths.
func enhanceTableAlignment(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check if this is a table header
		if !isTableHeader(line) {
			result = append(result, line)
			i++
			continue
		}

		// Extract column widths and alignment from the table
		tableLines := extractTable(lines[i:])
		if len(tableLines) == 0 {
			result = append(result, line)
			i++
			continue
		}

		// Rebuild table with enhanced alignment
		enhancedTable := rebuildTable(tableLines)
		result = append(result, enhancedTable...)

		i += len(tableLines)
	}

	return strings.Join(result, "\n")
}

// isTableHeader checks if a line is a table header.
func isTableHeader(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "|") &&
		strings.HasSuffix(strings.TrimSpace(line), "|")
}

// extractTable extracts a complete table from lines starting at index 0.
func extractTable(lines []string) []string {
	if len(lines) == 0 || !isTableHeader(lines[0]) {
		return nil
	}

	// Get header and alignment rows
	if len(lines) < 2 {
		return lines[:1]
	}

	// Check alignment row
	alignLine := lines[1]
	if !strings.Contains(alignLine, "---") {
		return lines[:1]
	}

	// Extract all consecutive table rows
	tableLines := []string{lines[0], lines[1]}
	for i := 2; i < len(lines); i++ {
		if isTableHeader(lines[i]) || strings.TrimSpace(lines[i]) == "" {
			break
		}
		tableLines = append(tableLines, lines[i])
	}

	return tableLines
}

// rebuildTable rebuilds a table with proper column alignment and cell wrapping.
func rebuildTable(tableLines []string) []string {
	if len(tableLines) < 2 {
		return tableLines
	}

	// Parse column widths and alignments
	columns := parseTableColumns(tableLines)

	// Rebuild each row
	for i := range tableLines {
		cells := parseTableRow(tableLines[i], len(columns))
		if len(cells) == len(columns) {
			tableLines[i] = formatTableRow(cells, columns)
		}
	}

	return tableLines
}

// TableColumn represents a column's width and alignment.
type TableColumn struct {
	Width    int
	Alignment string // "left", "center", "right"
}

// parseTableColumns parses column info from header and alignment rows.
func parseTableColumns(tableLines []string) []TableColumn {
	headerCells := parseTableRow(tableLines[0], 0)
	alignCells := parseTableRow(tableLines[1], len(headerCells))

	columns := make([]TableColumn, len(headerCells))
	for i, header := range headerCells {
		width := len(header)
		if i < len(alignCells) {
			align := strings.TrimSpace(alignCells[i])
			if strings.Contains(align, ":") {
				if strings.HasPrefix(align, ":") && strings.HasSuffix(align, ":") {
					columns[i].Alignment = "center"
				} else if strings.HasSuffix(align, ":") {
					columns[i].Alignment = "right"
				}
			}
		}

		// Scan data rows for max width
		for j := 2; j < len(tableLines); j++ {
			cells := parseTableRow(tableLines[j], len(headerCells))
			if i < len(cells) && len(cells[i]) > width {
				width = len(cells[i])
			}
		}

		columns[i].Width = width
	}

	return columns
}

// parseTableRow parses a table row into cells.
func parseTableRow(line string, expectedCols int) []string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return []string{}
	}

	// Remove outer pipes
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")

	// Split by pipe, but handle escaped pipes
	cells := strings.Split(line, "|")
	result := make([]string, len(cells))

	for i, cell := range cells {
		result[i] = strings.TrimSpace(cell)
	}

	return result
}

// formatTableRow formats a table row with proper alignment.
func formatTableRow(cells []string, columns []TableColumn) string {
	var sb strings.Builder

	sb.WriteString("|")
	for i, cell := range cells {
		if i < len(columns) {
			col := columns[i]
			padded := padCell(cell, col.Width, col.Alignment)
			sb.WriteString(" " + padded + " |")
		}
	}

	return sb.String()
}

// padCell pads a cell to the specified width with the given alignment.
func padCell(cell string, width int, alignment string) string {
	cellLen := len(cell)

	switch alignment {
	case "center":
		if cellLen >= width {
			return cell[:width]
		}
		leftPad := (width - cellLen) / 2
		rightPad := width - cellLen - leftPad
		return strings.Repeat(" ", leftPad) + cell + strings.Repeat(" ", rightPad)
	case "right":
		if cellLen >= width {
			return cell[:width]
		}
		return strings.Repeat(" ", width-cellLen) + cell
	default: // left
		if cellLen >= width {
			return cell[:width]
		}
		return cell + strings.Repeat(" ", width-cellLen)
	}
}

// applyCustomTaskListMarkers adds custom checkbox markers for task lists.
func applyCustomTaskListMarkers(content *string) {
	// Replace [x] with ☑ for checked tasks
	*content = strings.ReplaceAll(*content, "- [x] ", "- ☑ ")
	*content = strings.ReplaceAll(*content, "- [X] ", "- ☑ ")
	// Replace [ ] with ☐ for unchecked tasks
	*content = strings.ReplaceAll(*content, "- [ ] ", "- ☐ ")
}

// TableBuilder helps build markdown tables programmatically.
type TableBuilder struct {
	headers   []string
	rows      [][]string
	alignment []string // "left", "center", "right"
}

// NewTableBuilder creates a new table builder.
func NewTableBuilder(headers []string) *TableBuilder {
	return &TableBuilder{
		headers:   headers,
		rows:      make([][]string, 0),
		alignment: make([]string, len(headers)),
	}
}

// SetAlignment sets the alignment for a specific column.
// col is 0-indexed. align can be "left", "center", or "right".
func (t *TableBuilder) SetAlignment(col int, align string) {
	if col >= 0 && col < len(t.alignment) {
		t.alignment[col] = align
	}
}

// AddRow adds a row to the table.
func (t *TableBuilder) AddRow(cells []string) {
	if len(cells) == len(t.headers) {
		t.rows = append(t.rows, cells)
	}
}

// String returns the markdown table as a string.
func (t *TableBuilder) String() string {
	var sb strings.Builder

	// Header row
	sb.WriteString("| ")
	for i, h := range t.headers {
		sb.WriteString(h)
		if i < len(t.headers)-1 {
			sb.WriteString(" | ")
		}
	}
	sb.WriteString(" |\n")

	// Alignment row
	sb.WriteString("| ")
	for i, align := range t.alignment {
		switch align {
		case "center":
			sb.WriteString(":---:")
		case "right":
			sb.WriteString("---:")
		default: // left
			sb.WriteString("---")
		}
		if i < len(t.alignment)-1 {
			sb.WriteString(" | ")
		}
	}
	sb.WriteString(" |\n")

	// Data rows
	for _, row := range t.rows {
		sb.WriteString("| ")
		for i, cell := range row {
			// Escape pipes in cell content
			escaped := strings.ReplaceAll(cell, "|", "\\|")
			sb.WriteString(escaped)
			if i < len(row)-1 {
				sb.WriteString(" | ")
			}
		}
		sb.WriteString(" |\n")
	}

	return sb.String()
}

// TaskListBuilder helps build markdown task lists programmatically.
type TaskListBuilder struct {
	items []TaskItem
}

// TaskItem represents a single task item.
type TaskItem struct {
	Text     string
	Checked  bool
	Indent   int
	Subtasks []TaskItem
}

// NewTaskListBuilder creates a new task list builder.
func NewTaskListBuilder() *TaskListBuilder {
	return &TaskListBuilder{
		items: make([]TaskItem, 0),
	}
}

// AddItem adds a task item to the list.
func (t *TaskListBuilder) AddItem(text string, checked bool) {
	t.items = append(t.items, TaskItem{
		Text:    text,
		Checked: checked,
		Indent:  0,
	})
}

// AddSubtask adds a subtask to the last item.
func (t *TaskListBuilder) AddSubtask(text string, checked bool) {
	if len(t.items) > 0 {
		lastIdx := len(t.items) - 1
		t.items[lastIdx].Subtasks = append(t.items[lastIdx].Subtasks, TaskItem{
			Text:    text,
			Checked: checked,
			Indent:  1,
		})
	}
}

// String returns the task list as a markdown string.
func (t *TaskListBuilder) String() string {
	var sb strings.Builder

	for _, item := range t.items {
		t.renderTaskItem(&sb, item)
	}

	return sb.String()
}

// renderTaskItem renders a single task item recursively.
func (t *TaskListBuilder) renderTaskItem(sb *strings.Builder, item TaskItem) {
	indent := strings.Repeat("  ", item.Indent)

	// Checkbox
	if item.Checked {
		sb.WriteString(indent + "- [x] " + item.Text + "\n")
	} else {
		sb.WriteString(indent + "- [ ] " + item.Text + "\n")
	}

	// Subtasks
	for _, subtask := range item.Subtasks {
		t.renderTaskItem(sb, subtask)
	}
}

// LinkBuilder helps build markdown links programmatically.
type LinkBuilder struct {
	text string
	url  string
	title string
}

// NewLinkBuilder creates a new link builder.
func NewLinkBuilder(text, url string) *LinkBuilder {
	return &LinkBuilder{
		text: text,
		url:  url,
	}
}

// SetTitle sets the link title.
func (l *LinkBuilder) SetTitle(title string) *LinkBuilder {
	l.title = title
	return l
}

// String returns the link as markdown.
func (l *LinkBuilder) String() string {
	if l.title != "" {
		return fmt.Sprintf("[%s](%s %q)", l.text, l.url, l.title)
	}
	return fmt.Sprintf("[%s](%s)", l.text, l.url)
}

// ImageBuilder helps build markdown images programmatically.
type ImageBuilder struct {
	alt  string
	url  string
	title string
}

// NewImageBuilder creates a new image builder.
func NewImageBuilder(alt, url string) *ImageBuilder {
	return &ImageBuilder{
		alt: alt,
		url: url,
	}
}

// SetTitle sets the image title.
func (i *ImageBuilder) SetTitle(title string) *ImageBuilder {
	i.title = title
	return i
}

// String returns the image as markdown.
func (i *ImageBuilder) String() string {
	if i.title != "" {
		return fmt.Sprintf("![%s](%s %q)", i.alt, i.url, i.title)
	}
	return fmt.Sprintf("![%s](%s)", i.alt, i.url)
}

// CodeBlockBuilder helps build markdown code blocks.
type CodeBlockBuilder struct {
	lang     string
	content  string
	language string
}

// NewCodeBlockBuilder creates a new code block builder.
func NewCodeBlockBuilder(language, content string) *CodeBlockBuilder {
	return &CodeBlockBuilder{
		language: language,
		content:  content,
	}
}

// String returns the code block as markdown.
func (c *CodeBlockBuilder) String() string {
	return fmt.Sprintf("```%s\n%s\n```", c.language, c.content)
}

// InlineCode returns inline code markdown.
func InlineCode(code string) string {
	return fmt.Sprintf("`%s`", strings.ReplaceAll(code, "`", "\\`"))
}

// BlockQuote returns a block quote markdown.
func BlockQuote(text string) string {
	lines := strings.Split(text, "\n")
	var sb strings.Builder

	for _, line := range lines {
		sb.WriteString("> ")
		if line != "" {
			sb.WriteString(line)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// HorizontalRule returns a horizontal rule markdown.
func HorizontalRule() string {
	return "---"
}

// EscapeText escapes special markdown characters.
func EscapeText(text string) string {
	// Order matters - escape backslashes first, then other special chars
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"`", "\\`",
		"*", "\\*",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		".", "\\.",
		"!", "\\!",
		"|", "\\|",
	)

	return replacer.Replace(text)
}

// IsMarkdownTable checks if content contains a markdown table.
func IsMarkdownTable(content string) bool {
	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return false
	}

	// Check for table header row
	headerLine := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(headerLine, "|") || !strings.HasSuffix(headerLine, "|") {
		return false
	}

	// Check for alignment row (second line)
	if len(lines) < 2 {
		return false
	}
	alignLine := strings.TrimSpace(lines[1])
	return strings.HasPrefix(alignLine, "|") && strings.HasSuffix(alignLine, "|") &&
		strings.Contains(alignLine, "---")
}

// IsTaskList checks if content contains a task list.
func IsTaskList(content string) bool {
	return strings.Contains(content, "- [x]") ||
		strings.Contains(content, "- [X]") ||
		strings.Contains(content, "- [ ]")
}

// ExtractLinks extracts all markdown links from content.
// Returns []Link with text, url, and position information.
type Link struct {
	Text string
	URL  string
	Start int
	End  int
}

// ExtractLinks extracts all markdown links from content.
func ExtractLinks(content string) []Link {
	var links []Link

	// Pattern for [text](url) or [text](url "title")
	pattern := `\[([^\]]+)\]\(([^\)]+)`

	locs := indexOfPattern(content, pattern)
	for _, loc := range locs {
		matched := content[loc[0]:loc[1]]

		// Extract text between brackets
		textStart := strings.Index(matched, "[")
		textEnd := strings.Index(matched, "]")
		text := matched[textStart+1 : textEnd]

		// Extract URL from parentheses
		urlStart := strings.Index(matched, "(")
		urlEnd := strings.Index(matched, ")")
		urlPart := matched[urlStart+1 : urlEnd]

		// Remove title if present
		url := strings.Split(urlPart, " ")[0]

		links = append(links, Link{
			Text: text,
			URL:  url,
			Start: loc[0],
			End:   loc[1],
		})
	}

	return links
}

// indexOfPattern finds all occurrences of a regex pattern in content.
func indexOfPattern(content, pattern string) [][2]int {
	// Simple implementation for common markdown patterns
	var results [][2]int

	// For markdown links [text](url)
	if strings.Contains(pattern, `[`) && strings.Contains(pattern, `]\(`) {
		start := 0
		for {
			openBracket := strings.Index(content[start:], "[")
			if openBracket == -1 {
				break
			}
			openBracket += start

			closeBracket := strings.Index(content[openBracket:], "]")
			if closeBracket == -1 {
				break
			}
			closeBracket += openBracket

			if strings.HasPrefix(content[closeBracket+1:], "(") {
				openParen := closeBracket + 1
				closeParen := strings.Index(content[openParen:], ")")
				if closeParen != -1 {
					closeParen += openParen
					results = append(results, [2]int{openBracket, closeParen + 1})
					start = closeParen + 1
					continue
				}
			}

			start = closeBracket + 1
		}
	}

	return results
}
