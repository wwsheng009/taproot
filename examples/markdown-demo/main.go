package main

import (
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/styles"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("          Taproot Markdown Rendering Enhancement Demo")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// 1. Table Rendering
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("1. Markdown Table Rendering")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	headers := []string{"Name", "Age", "City", "Occupation"}
	rows := [][]string{
		{"Alice Johnson", "30", "New York", "Software Engineer"},
		{"Bob Smith", "25", "San Francisco", "Designer"},
		{"Charlie Brown", "35", "Los Angeles", "Product Manager"},
		{"Diana Prince", "28", "Seattle", "Data Scientist"},
	}

	tableStyles := styles.DefaultMarkdownTable()
	table := styles.RenderTable(headers, rows, tableStyles, 80)
	fmt.Println(table)
	fmt.Println()

	// 2. Task List Rendering
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("2. Markdown Task List Rendering")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	tasks := []styles.TaskItem{
		{Checked: true, Text: "Implement table rendering"},
		{Checked: true, Text: "Implement task list rendering"},
		{Checked: true, Text: "Implement link rendering"},
		{Checked: true, Text: "Implement image rendering"},
		{Checked: false, Text: "Create comprehensive tests"},
		{Checked: false, Text: "Create interactive demo"},
	}

	taskStyles := styles.DefaultMarkdownTaskList()
	taskList := styles.RenderTaskList(tasks, taskStyles)
	fmt.Println(taskList)

	// 3. Link Rendering
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("3. Markdown Link Rendering")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	linkStyles := styles.DefaultMarkdownLink()

	links := []struct {
		text string
		url  string
	}{
		{"Taproot GitHub", "https://github.com/wwsheng009/taproot"},
		{"Bubbletea Framework", "https://github.com/charmbracelet/bubbletea"},
		{"Lipgloss Styling", "https://github.com/charmbracelet/lipgloss"},
		{"Glamour Markdown", "https://github.com/charmbracelet/glamour"},
	}

	for _, link := range links {
		renderedLink := styles.RenderLink(link.text, link.url, linkStyles)
		fmt.Printf("  • %s\n", renderedLink)
	}
	fmt.Println()

	// 4. Image Rendering
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("4. Markdown Image Rendering")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	imageStyles := styles.DefaultMarkdownImage()

	images := []struct {
		alt string
		url string
	}{
		{"Taproot Logo", "https://example.com/logo.png"},
		{"Architecture Diagram", "https://example.com/architecture.svg"},
		{"Component Hierarchy", "https://example.com/hierarchy.png"},
	}

	for _, img := range images {
		renderedImage := styles.RenderImage(img.alt, img.url, imageStyles)
		fmt.Printf("  • %s\n", renderedImage)
	}
	fmt.Println()

	// 5. Parse Table from Markdown
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("5. Parse Table from Markdown Text")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	markdownTable := `| Feature | Status | Priority |
| :--- | :---: | ----: |
| Table Rendering | ✅ Done | High |
| Task Lists | ✅ Done | High |
| Links | ✅ Done | Medium |
| Images | ✅ Done | Medium |
| Tests | ✅ Done | High |`

	fmt.Println("Input Markdown:")
	fmt.Println("  " + strings.ReplaceAll(markdownTable, "\n", "\n  "))
	fmt.Println()

	parsedHeaders, parsedRows, parsedAlignments := styles.ParseTable(markdownTable)

	fmt.Printf("Parsed Headers: %v\n", parsedHeaders)
	fmt.Printf("Parsed Rows: %d rows\n", len(parsedRows))
	fmt.Printf("Parsed Alignments: %v\n", parsedAlignments)
	fmt.Println()

	// Render the parsed table
	parsedTableStyles := styles.DefaultMarkdownTable()
	parsedTable := styles.RenderTable(parsedHeaders, parsedRows, parsedTableStyles, 80)
	fmt.Println("Rendered Output:")
	fmt.Println(parsedTable)
	fmt.Println()

	// 6. Custom Styling Example
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("6. Custom Styled Table (Compact)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	compactHeaders := []string{"Type", "Count", "%"}
	compactRows := [][]string{
		{"Users", "150", "60%"},
		{"Admins", "50", "20%"},
		{"Guests", "50", "20%"},
	}

	// Create compact styling
	baseStyle := styles.DefaultStyles()
	compactStyles := styles.MarkdownTable{
		Header:       baseStyle.Base.Bold(true).Foreground(baseStyle.Primary),
		Border:       baseStyle.Base.Foreground(baseStyle.FgMuted),
		BorderRowSep: baseStyle.Base.Foreground(baseStyle.FgMuted),
		Cell:         baseStyle.Base.Foreground(baseStyle.FgBase),
		EvenRow:      baseStyle.Base,
		OddRow:       baseStyle.Base.Background(baseStyle.BgSubtle),
	}

	compactTable := styles.RenderTable(compactHeaders, compactRows, compactStyles, 50)
	fmt.Println(compactTable)
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("                    Demo Complete")
	fmt.Println("═══════════════════════════════════════════════════════════════")
}
