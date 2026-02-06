// Buffer Table Example - Table layout using buffer system
//
// This example demonstrates:
// - Table with headers and rows
// - Column alignment (left, center, right)
// - Cell borders and styling
// - Dynamic content rendering
// - Sorting indicators
//
// Usage: go run main.go

package main

import (
	"fmt"
	"strconv"

	"github.com/wwsheng009/taproot/ui/render/buffer"
)

// Column defines a table column
type Column struct {
	title    string
	width    int
	align    string // "left", "center", "right"
	sortable bool
	sortDir  int // 0 = none, 1 = asc, -1 = desc
}

// Cell represents table cell content
type Cell struct {
	content string
	style   buffer.Style
	colSpan int
}

// Row represents a table row
type Row struct {
	cells   []*Cell
	alt     bool // Alternating row color
}

// Table component
type Table struct {
	columns []*Column
	rows    []*Row
	title   string
	footer  string
	border  bool
}

func (t *Table) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw title
	if t.title != "" {
		titleX := rect.X + (rect.Width - len(t.title)) / 2
		buf.WriteString(buffer.Point{X: titleX, Y: rect.Y}, t.title, buffer.Style{Foreground: "#86", Bold: true})
	}

	// Calculate column widths
	totalWidth := rect.Width
	colWidths := make([]int, len(t.columns))
	flexCols := 0
	fixedWidth := 0

	for i, col := range t.columns {
		if col.width > 0 {
			colWidths[i] = col.width
			fixedWidth += col.width
		} else {
			flexCols++
		}
	}

	remainingWidth := totalWidth - fixedWidth - len(t.columns) - 1 // Account for borders
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	flexWidth := remainingWidth
	if flexCols > 0 {
		flexWidth = remainingWidth / flexCols
	}

	for i, col := range t.columns {
		if col.width == 0 {
			colWidths[i] = flexWidth
		}
	}

	// Starting Y position
	startY := rect.Y
	if t.title != "" {
		startY += 2
	}

	// Draw header
	headerY := startY
	for x := rect.X; x < rect.X+rect.Width; x++ {
		buf.WriteString(buffer.Point{X: x, Y: headerY}, "─", buffer.Style{Foreground: "#244"})
	}
	buf.WriteString(buffer.Point{X: rect.X, Y: headerY}, "┬", buffer.Style{Foreground: "#244"})
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: headerY}, "┬", buffer.Style{Foreground: "#244"})

	headerTextY := headerY + 1
	x := rect.X + 1
	for i, col := range t.columns {
		title := col.title
		if col.sortable {
			if col.sortDir == 1 {
				title += " ▲"
			} else if col.sortDir == -1 {
				title += " ▼"
			}
		}

		// Align title
		textX := x
		if col.align == "center" {
			textX = x + (colWidths[i] - len(title)) / 2
		} else if col.align == "right" {
			textX = x + colWidths[i] - len(title) - 1
		}

		if textX < x {
			textX = x
		}

		buf.WriteString(buffer.Point{X: textX, Y: headerTextY}, title, buffer.Style{Foreground: "#86", Bold: true})

		// Draw column separator
		sepX := x + colWidths[i]
		if sepX < rect.X+rect.Width-1 {
			buf.WriteString(buffer.Point{X: sepX, Y: headerTextY}, "│", buffer.Style{Foreground: "#244"})
		}

		x += colWidths[i] + 1
	}

	// Draw header bottom border
	for x := rect.X; x < rect.X+rect.Width; x++ {
		buf.WriteString(buffer.Point{X: x, Y: headerTextY + 1}, "─", buffer.Style{Foreground: "#244"})
	}
	buf.WriteString(buffer.Point{X: rect.X, Y: headerTextY + 1}, "├", buffer.Style{Foreground: "#244"})
	x = rect.X + 1
	for _, colWidth := range colWidths {
		sepX := x + colWidth
		if sepX < rect.X+rect.Width-1 {
			buf.WriteString(buffer.Point{X: sepX, Y: headerTextY + 1}, "┼", buffer.Style{Foreground: "#244"})
		}
		x += colWidth + 1
	}
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: headerTextY + 1}, "┤", buffer.Style{Foreground: "#244"})

	// Draw rows
	rowY := headerTextY + 2
	maxRows := rect.Height - 3
	if t.title != "" {
		maxRows -= 2
	}

	for rowIndex, row := range t.rows {
		if rowY >= rect.Y+rect.Height-1 {
			break
		}

		// Draw row background for alt rows
		if row.alt {
			for bx := rect.X; bx < rect.X+rect.Width-1; bx++ {
				bgStyle := buffer.Style{Background: "#234"}
				buf.WriteString(buffer.Point{X: bx, Y: rowY}, " ", bgStyle)
			}
		}

		// Draw cells
		x = rect.X + 1
		for cellIndex, cell := range row.cells {
			if cellIndex >= len(colWidths) {
				break
			}
			colWidth := colWidths[cellIndex]

			content := cell.content
			if len(content) > colWidth-2 {
				content = content[:colWidth-3] + "…"
			}

			// Align content
			textX := x
			align := "left"
			if cellIndex < len(t.columns) {
				align = t.columns[cellIndex].align
			}
			if align == "center" {
				textX = x + (colWidth - len(content)) / 2
			} else if align == "right" {
				textX = x + colWidth - len(content) - 1
			}

			style := cell.style
			if row.alt {
				style.Background = "#234"
			}
			buf.WriteString(buffer.Point{X: textX, Y: rowY}, content, style)

			// Draw column separator
			sepX := x + colWidth
			if sepX < rect.X+rect.Width-1 {
				buf.WriteString(buffer.Point{X: sepX, Y: rowY}, "│", buffer.Style{Foreground: "#244"})
			}

			x += colWidth + 1
		}

		// Draw row separator
		if rowIndex < len(t.rows)-1 {
			for bx := rect.X; bx < rect.X+rect.Width; bx++ {
				buf.WriteString(buffer.Point{X: bx, Y: rowY + 1}, "─", buffer.Style{Foreground: "#238"})
			}
			buf.WriteString(buffer.Point{X: rect.X, Y: rowY + 1}, "├", buffer.Style{Foreground: "#244"})
			x = rect.X + 1
			for _, colWidth := range colWidths {
				sepX := x + colWidth
				if sepX < rect.X+rect.Width-1 {
					buf.WriteString(buffer.Point{X: sepX, Y: rowY + 1}, "┼", buffer.Style{Foreground: "#244"})
				}
				x += colWidth + 1
			}
			buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rowY + 1}, "┤", buffer.Style{Foreground: "#244"})
		}

		rowY += 2
	}

	// Draw bottom border
	for bx := rect.X; bx < rect.X+rect.Width; bx++ {
		buf.WriteString(buffer.Point{X: bx, Y: rect.Height - 2}, "─", buffer.Style{Foreground: "#244"})
	}
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Height - 2}, "└", buffer.Style{Foreground: "#244"})
	x = rect.X + 1
	for _, colWidth := range colWidths {
		sepX := x + colWidth
		if sepX < rect.X+rect.Width-1 {
			buf.WriteString(buffer.Point{X: sepX, Y: rect.Height - 2}, "┴", buffer.Style{Foreground: "#244"})
		}
		x += colWidth + 1
	}
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Height - 2}, "┘", buffer.Style{Foreground: "#244"})

	// Draw footer
	if t.footer != "" {
		footerY := rect.Height - 1
		buf.WriteString(buffer.Point{X: rect.X, Y: footerY}, t.footer, buffer.Style{Foreground: "#244"})
	}
}

func (t *Table) MinSize() (int, int)     { return 40, 10 }
func (t *Table) PreferredSize() (int, int) { return 80, 20 }

// Helper function to create a simple cell
func NewCell(content string, color string) *Cell {
	return &Cell{
		content: content,
		style:   buffer.Style{Foreground: color},
	}
}

// Number cell (right aligned)
func NewNumberCell(num int, color string) *Cell {
	return &Cell{
		content: strconv.Itoa(num),
		style:   buffer.Style{Foreground: color},
	}
}

// Status cell with color coding
func NewStatusCell(status string) *Cell {
	color := "#250" // Green
	switch status {
	case "Active":
		color = "#120"
	case "Pending":
		color = "#226"
	case "Failed":
		color = "#160"
	case "Offline":
		color = "#244"
	}
	return &Cell{
		content: status,
		style:   buffer.Style{Foreground: color, Bold: true},
	}
}

func main() {
	width := 100
	height := 30

	fmt.Println(repeat("=", width))
	fmt.Println("Buffer Table Example - Taproot")
	fmt.Println(repeat("=", width))
	fmt.Println()

	buf := buffer.NewBuffer(width, height)

	// Example 1: Simple data table
	simpleTable := &Table{
		title: "Server Status",
		columns: []*Column{
			{title: "Server", width: 20, align: "left", sortable: true, sortDir: 1},
			{title: "Status", width: 15, align: "center", sortable: true},
			{title: "CPU", width: 10, align: "right", sortable: true, sortDir: -1},
			{title: "Memory", width: 10, align: "right", sortable: true},
			{title: "Uptime", width: 15, align: "right", sortable: true},
		},
		rows: []*Row{
			{
				cells: []*Cell{
					NewCell("web-server-01", "#250"),
					NewStatusCell("Active"),
					NewNumberCell(45, "#86"),
					NewNumberCell(62, "#86"),
					NewCell("15d 4h 23m", "#244"),
				},
				alt: false,
			},
			{
				cells: []*Cell{
					NewCell("web-server-02", "#250"),
					NewStatusCell("Active"),
					NewNumberCell(32, "#86"),
					NewNumberCell(58, "#86"),
					NewCell("12d 1h 45m", "#244"),
				},
				alt: true,
			},
			{
				cells: []*Cell{
					NewCell("db-server-01", "#250"),
					NewStatusCell("Active"),
					NewNumberCell(78, "#226"),
					NewNumberCell(85, "#226"),
					NewCell("45d 12h 30m", "#244"),
				},
				alt: false,
			},
			{
				cells: []*Cell{
					NewCell("db-server-02", "#250"),
					NewStatusCell("Pending"),
					NewNumberCell(0, "#244"),
					NewNumberCell(0, "#244"),
					NewCell("0d 0h 0m", "#244"),
				},
				alt: true,
			},
			{
				cells: []*Cell{
					NewCell("cache-server-01", "#250"),
					NewStatusCell("Failed"),
					NewNumberCell(0, "#244"),
					NewNumberCell(0, "#244"),
					NewCell("0d 0h 5m", "#160"),
				},
				alt: false,
			},
		},
		footer: " Showing 5 of 12 servers | ↑↓: Navigate | Enter: Select | F1: Filter ",
	}

	simpleTable.Render(buf, buffer.Rect{X: 2, Y: 1, Width: width - 4, Height: 18})

	// Example 2: Compact metrics table below
	metricsTable := &Table{
		title: "Performance Metrics",
		columns: []*Column{
			{title: "Metric", width: 0, align: "left"},
			{title: "Value", width: 0, align: "right"},
			{title: "Change", width: 0, align: "right"},
		},
		rows: []*Row{
			{
				cells: []*Cell{
					NewCell("Requests/sec", "#250"),
					NewNumberCell(1234, "#86"),
					NewCell("+12.5%", "#120"),
				},
				alt: false,
			},
			{
				cells: []*Cell{
					NewCell("Avg Response", "#250"),
					NewCell("45ms", "#86"),
					NewCell("-5.2%", "#120"),
				},
				alt: true,
			},
			{
				cells: []*Cell{
					NewCell("Error Rate", "#250"),
					NewCell("0.12%", "#86"),
					NewCell("+0.01%", "#226"),
				},
				alt: false,
			},
			{
				cells: []*Cell{
					NewCell("Active Users", "#250"),
					NewNumberCell(8432, "#86"),
					NewCell("+234", "#120"),
				},
				alt: true,
			},
		},
		footer: " Updated: 5 seconds ago ",
	}

	metricsTable.Render(buf, buffer.Rect{X: 2, Y: 20, Width: width - 4, Height: height - 2})

	// Render output
	output := buf.Render()
	fmt.Print(output)

	// Print summary
	fmt.Printf("\n\nTable Features Demonstrated:\n")
	fmt.Printf("  • Columns with custom widths and alignment\n")
	fmt.Printf("  • Sortable columns with indicators (▲/▼)\n")
	fmt.Printf("  • Alternating row colors\n")
	fmt.Printf("  • Status cells with color coding\n")
	fmt.Printf("  • Row and column borders\n")
	fmt.Printf("  • Title and footer sections\n")
	fmt.Printf("  • Multiple tables on same screen\n")
}

func repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	result := s
	for i := 1; i < count; i++ {
		result += s
	}
	return result
}
