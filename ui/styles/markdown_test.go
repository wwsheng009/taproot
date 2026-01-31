package styles

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestRenderTable(t *testing.T) {
	t.Run("BasicTable", func(t *testing.T) {
		headers := []string{"Name", "Age", "City"}
		rows := [][]string{
			{"Alice", "30", "New York"},
			{"Bob", "25", "San Francisco"},
			{"Charlie", "35", "Los Angeles"},
		}

		styles := DefaultMarkdownTable()
		result := RenderTable(headers, rows, styles, 80)

		if result == "" {
			t.Error("RenderTable() should not return empty string")
		}

		// Check that headers are present
		if !contains(result, "Name") {
			t.Error("Table should contain 'Name' header")
		}

		// Check that rows are present
		if !contains(result, "Alice") {
			t.Error("Table should contain 'Alice'")
		}
	})

	t.Run("EmptyTable", func(t *testing.T) {
		headers := []string{}
		rows := [][]string{}

		styles := DefaultMarkdownTable()
		result := RenderTable(headers, rows, styles, 80)

		if result != "" {
			t.Error("Empty table should return empty string")
		}
	})

	t.Run("TableWithMaxWidth", func(t *testing.T) {
		headers := []string{"Very Long Header Name", "Another Long Header"}
		rows := [][]string{
			{"Some long content here", "More content"},
		}

		styles := DefaultMarkdownTable()
		result := RenderTable(headers, rows, styles, 40)

		if result == "" {
			t.Error("RenderTable() should handle max width constraint")
		}
	})
}

func TestRenderTaskList(t *testing.T) {
	t.Run("BasicTaskList", func(t *testing.T) {
		items := []TaskItem{
			{Checked: true, Text: "Completed task"},
			{Checked: false, Text: "Pending task"},
			{Checked: true, Text: "Another completed task"},
		}

		styles := DefaultMarkdownTaskList()
		result := RenderTaskList(items, styles)

		if result == "" {
			t.Error("RenderTaskList() should not return empty string")
		}

		// Check for checked indicator
		if !contains(result, "☑") {
			t.Error("Task list should contain checked indicator")
		}

		// Check for unchecked indicator
		if !contains(result, "☐") {
			t.Error("Task list should contain unchecked indicator")
		}
	})

	t.Run("EmptyTaskList", func(t *testing.T) {
		items := []TaskItem{}

		styles := DefaultMarkdownTaskList()
		result := RenderTaskList(items, styles)

		// Empty list should return empty string
		if result != "\n" && result != "" {
			t.Errorf("Empty task list should return empty or newline, got: %q", result)
		}
	})
}

func TestRenderLink(t *testing.T) {
	t.Run("BasicLink", func(t *testing.T) {
		styles := DefaultMarkdownLink()
		result := RenderLink("Example", "https://example.com", styles)

		if result == "" {
			t.Error("RenderLink() should not return empty string")
		}

		if !contains(result, "Example") {
			t.Error("Link should contain text 'Example'")
		}

		if !contains(result, "https://example.com") {
			t.Error("Link should contain URL when ShowURL is true")
		}
	})

	t.Run("LinkWithoutURL", func(t *testing.T) {
		styles := DefaultMarkdownLink()
		styles.ShowURL = false

		result := RenderLink("Example", "https://example.com", styles)

		if result == "" {
			t.Error("RenderLink() should not return empty string")
		}

		if !contains(result, "Example") {
			t.Error("Link should contain text 'Example'")
		}
	})
}

func TestRenderImage(t *testing.T) {
	t.Run("BasicImage", func(t *testing.T) {
		styles := DefaultMarkdownImage()
		result := RenderImage("Alt text", "https://example.com/image.png", styles)

		if result == "" {
			t.Error("RenderImage() should not return empty string")
		}

		if !contains(result, "Alt text") {
			t.Error("Image should contain alt text 'Alt text'")
		}

		if !contains(result, "🖼") {
			t.Error("Image should contain default icon")
		}
	})

	t.Run("ImageWithoutInfo", func(t *testing.T) {
		styles := DefaultMarkdownImage()
		styles.ShowInfo = false

		result := RenderImage("Alt text", "https://example.com/image.png", styles)

		if result == "" {
			t.Error("RenderImage() should not return empty string")
		}

		// Should not contain URL when ShowInfo is false
		if contains(result, "https://example.com/image.png") {
			t.Error("Image should not contain URL when ShowInfo is false")
		}
	})
}

func TestParseTable(t *testing.T) {
	t.Run("ValidTable", func(t *testing.T) {
		markdown := `| Name | Age | City |
| :--- | :---: | ----: |
| Alice | 30 | New York |
| Bob | 25 | San Francisco |`

		headers, rows, alignments := ParseTable(markdown)

		if len(headers) != 3 {
			t.Errorf("Expected 3 headers, got %d", len(headers))
		}

		if headers[0] != "Name" {
			t.Errorf("Expected first header 'Name', got '%s'", headers[0])
		}

		if len(rows) != 2 {
			t.Errorf("Expected 2 rows, got %d", len(rows))
		}

		if len(alignments) != 3 {
			t.Errorf("Expected 3 alignments, got %d", len(alignments))
		}

		if alignments[0] != "left" {
			t.Errorf("Expected first alignment 'left', got '%s'", alignments[0])
		}

		if alignments[1] != "center" {
			t.Errorf("Expected second alignment 'center', got '%s'", alignments[1])
		}

		if alignments[2] != "right" {
			t.Errorf("Expected third alignment 'right', got '%s'", alignments[2])
		}
	})

	t.Run("EmptyTable", func(t *testing.T) {
		markdown := ""

		headers, rows, alignments := ParseTable(markdown)

		if headers != nil {
			t.Error("Empty markdown should return nil headers")
		}

		if rows != nil {
			t.Error("Empty markdown should return nil rows")
		}

		if alignments != nil {
			t.Error("Empty markdown should return nil alignments")
		}
	})

	t.Run("TableWithWhitespace", func(t *testing.T) {
		markdown := `|  Name  |  Age  |  City  |
| :--- | :---: | ----: |
|  Alice  |  30  |  New York  |`

		headers, rows, _ := ParseTable(markdown)

		if headers[0] != "Name" {
			t.Errorf("Expected trimmed header 'Name', got '%s'", headers[0])
		}

		if rows[0][0] != "Alice" {
			t.Errorf("Expected trimmed cell 'Alice', got '%s'", rows[0][0])
		}
	})
}

func TestParseTableRow(t *testing.T) {
	t.Run("ValidRow", func(t *testing.T) {
		line := "| Name | Age | City |"
		result := parseTableRow(line)

		if len(result) != 3 {
			t.Errorf("Expected 3 cells, got %d", len(result))
		}

		if result[0] != "Name" {
			t.Errorf("Expected first cell 'Name', got '%s'", result[0])
		}
	})

	t.Run("InvalidRow", func(t *testing.T) {
		line := "Not a table row"
		result := parseTableRow(line)

		if result != nil {
			t.Error("Invalid row should return nil")
		}
	})

	t.Run("RowWithExtraSpaces", func(t *testing.T) {
		line := "|  Name  |  Age  |  City  |"
		result := parseTableRow(line)

		if result[0] != "Name" {
			t.Errorf("Expected trimmed cell 'Name', got '%s'", result[0])
		}

		if result[1] != "Age" {
			t.Errorf("Expected trimmed cell 'Age', got '%s'", result[1])
		}
	})
}

func TestParseTableAlignment(t *testing.T) {
	t.Run("MixedAlignment", func(t *testing.T) {
		line := "| :--- | :---: | ----: |"
		result := parseTableAlignment(line)

		if len(result) != 3 {
			t.Errorf("Expected 3 alignments, got %d", len(result))
		}

		if result[0] != "left" {
			t.Errorf("Expected first alignment 'left', got '%s'", result[0])
		}

		if result[1] != "center" {
			t.Errorf("Expected second alignment 'center', got '%s'", result[1])
		}

		if result[2] != "right" {
			t.Errorf("Expected third alignment 'right', got '%s'", result[2])
		}
	})

	t.Run("DefaultAlignment", func(t *testing.T) {
		line := "| --- | --- | --- |"
		result := parseTableAlignment(line)

		for i, align := range result {
			if align != "left" {
				t.Errorf("Expected default alignment 'left' at index %d, got '%s'", i, align)
			}
		}
	})
}

func TestDefaultStyles(t *testing.T) {
	t.Run("DefaultMarkdownTable", func(t *testing.T) {
		styles := DefaultMarkdownTable()

		if lipgloss.Width(styles.Header.Render("test")) == 0 {
			t.Error("Default table header style should not be empty")
		}

		if lipgloss.Width(styles.Cell.Render("test")) == 0 {
			t.Error("Default table cell style should not be empty")
		}
	})

	t.Run("DefaultMarkdownTaskList", func(t *testing.T) {
		styles := DefaultMarkdownTaskList()

		if styles.CheckedIndicator == "" {
			t.Error("Default checked indicator should not be empty")
		}

		if styles.UncheckedIndicator == "" {
			t.Error("Default unchecked indicator should not be empty")
		}

		if styles.CheckedIndicator != "☑" {
			t.Errorf("Expected checked indicator '☑', got '%s'", styles.CheckedIndicator)
		}

		if styles.UncheckedIndicator != "☐" {
			t.Errorf("Expected unchecked indicator '☐', got '%s'", styles.UncheckedIndicator)
		}
	})

	t.Run("DefaultMarkdownLink", func(t *testing.T) {
		styles := DefaultMarkdownLink()

		if !styles.ShowURL {
			t.Error("Default link style should show URL")
		}

		if !styles.Underline {
			t.Error("Default link style should underline text")
		}
	})

	t.Run("DefaultMarkdownImage", func(t *testing.T) {
		styles := DefaultMarkdownImage()

		if styles.Icon == "" {
			t.Error("Default image icon should not be empty")
		}

		if styles.Icon != "🖼" {
			t.Errorf("Expected icon '🖼', got '%s'", styles.Icon)
		}

		if !styles.ShowInfo {
			t.Error("Default image style should show info")
		}
	})
}

func TestRenderTableCell(t *testing.T) {
	t.Run("CellWithoutTruncation", func(t *testing.T) {
		cell := renderTableCell("Hello", 10, lipgloss.NewStyle())

		if !contains(cell, "Hello") {
			t.Error("Cell should contain 'Hello'")
		}

		if lipgloss.Width(cell) != 12 { // 10 + 2 padding
			t.Errorf("Expected cell width 12, got %d", lipgloss.Width(cell))
		}
	})

	t.Run("CellWithTruncation", func(t *testing.T) {
		cell := renderTableCell("This is a very long text", 10, lipgloss.NewStyle())

		if lipgloss.Width(cell) > 12 { // 10 + 2 padding
			t.Errorf("Cell should be truncated to fit width")
		}
	})

	t.Run("CellWithPadding", func(t *testing.T) {
		cell := renderTableCell("Hi", 10, lipgloss.NewStyle())

		if !strings.HasPrefix(cell, " ") {
			t.Error("Cell should have left padding")
		}

		if !strings.HasSuffix(cell, " ") {
			t.Error("Cell should have right padding")
		}
	})
}

func TestRenderTableBorder(t *testing.T) {
	t.Run("BasicBorder", func(t *testing.T) {
		widths := []int{10, 15, 20}
		result := renderTableBorder(widths)

		if !strings.HasPrefix(result, "┌") {
			t.Error("Border should start with ┌")
		}

		if !strings.HasSuffix(result, "┐") {
			t.Error("Border should end with ┐")
		}

		if !contains(result, "┬") {
			t.Error("Border should contain ┬ separator")
		}
	})
}

func TestRenderTableSeparator(t *testing.T) {
	t.Run("BasicSeparator", func(t *testing.T) {
		widths := []int{10, 15, 20}
		result := renderTableSeparator(widths)

		if !strings.HasPrefix(result, "├") {
			t.Error("Separator should start with ├")
		}

		if !strings.HasSuffix(result, "┤") {
			t.Error("Separator should end with ┤")
		}

		if !contains(result, "┼") {
			t.Error("Separator should contain ┼ separator")
		}
	})
}

// Helper function
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
