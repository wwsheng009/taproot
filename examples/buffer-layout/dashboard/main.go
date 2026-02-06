// Buffer Dashboard - Multi-panel dashboard example
//
// This example demonstrates advanced buffer layout features:
// - Multiple panels in different positions
// - Custom Renderable components
// - Dynamic content rendering
// - Status indicators, progress bars, and metrics
//
// Usage: go run main.go

package main

import (
	"fmt"
	"time"

	"github.com/wwsheng009/taproot/ui/render/buffer"
)

// Custom component types

// StatusPanel shows system status with an indicator
type StatusPanel struct {
	status  string
	message string
	color   string
}

func (s *StatusPanel) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw status indicator
	indicator := "●"
	buf.WriteString(buffer.Point{X: rect.X + 1, Y: rect.Y}, indicator, buffer.Style{Foreground: s.color, Bold: true})
	buf.WriteString(buffer.Point{X: rect.X + 3, Y: rect.Y}, s.status, buffer.Style{Foreground: s.color, Bold: true})
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - len(s.message) - 2, Y: rect.Y}, s.message, buffer.Style{Foreground: "#244"})
}

func (s *StatusPanel) MinSize() (int, int) { return 20, 1 }
func (s *StatusPanel) PreferredSize() (int, int) { return 40, 1 }

// MetricPanel shows a metric with label and value
type MetricPanel struct {
	label string
	value string
	delta string
}

func (m *MetricPanel) Render(buf *buffer.Buffer, rect buffer.Rect) {
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, m.label, buffer.Style{Foreground: "#244"})
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + 1}, m.value, buffer.Style{Foreground: "#86", Bold: true})
	if m.delta != "" {
		color := "#120" // Green
		if m.delta[0] == '-' {
			color = "#160" // Red
		}
		buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + 2}, m.delta, buffer.Style{Foreground: color})
	}
}

func (m *MetricPanel) MinSize() (int, int)     { return 15, 3 }
func (m *MetricPanel) PreferredSize() (int, int) { return 20, 3 }

// ProgressPanel shows a horizontal progress bar
type ProgressPanel struct {
	label  string
	percent int
	total   int
	current int
}

func (p *ProgressPanel) Render(buf *buffer.Buffer, rect buffer.Rect) {
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, p.label, buffer.Style{Foreground: "#244"})

	// Draw progress bar background
	barWidth := rect.Width - 2
	barY := rect.Y + 1

	for i := 0; i < barWidth; i++ {
		buf.WriteString(buffer.Point{X: rect.X + i, Y: barY}, "░", buffer.Style{Foreground: "#238"})
	}

	// Draw filled portion
	filledWidth := (barWidth * p.percent) / 100
	for i := 0; i < filledWidth; i++ {
		buf.WriteString(buffer.Point{X: rect.X + i, Y: barY}, "█", buffer.Style{Foreground: "#120", Bold: true})
	}

	// Draw percentage text
	percentText := fmt.Sprintf(" %d%% (%d/%d)", p.percent, p.current, p.total)
	textX := rect.X + barWidth/2 - len(percentText)/2
	if textX < rect.X {
		textX = rect.X
	}
	buf.WriteString(buffer.Point{X: textX, Y: barY + 1}, percentText, buffer.Style{Foreground: "#250"})
}

func (p *ProgressPanel) MinSize() (int, int)     { return 20, 3 }
func (p *ProgressPanel) PreferredSize() (int, int) { return 40, 3 }

// LogPanel shows scrolling log entries
type LogPanel struct {
	entries []string
}

func (l *LogPanel) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw border using box drawing characters
	topLeft := "┌"
	topRight := "┐"
	bottomLeft := "└"
	bottomRight := "┘"
	horizontal := "─"
	vertical := "│"

	// Top border
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, topLeft, buffer.Style{Foreground: "#238"})
	for i := 1; i < rect.Width-1; i++ {
		buf.WriteString(buffer.Point{X: rect.X + i, Y: rect.Y}, horizontal, buffer.Style{Foreground: "#238"})
	}
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y}, topRight, buffer.Style{Foreground: "#238"})

	// Bottom border
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + rect.Height - 1}, bottomLeft, buffer.Style{Foreground: "#238"})
	for i := 1; i < rect.Width-1; i++ {
		buf.WriteString(buffer.Point{X: rect.X + i, Y: rect.Y + rect.Height - 1}, horizontal, buffer.Style{Foreground: "#238"})
	}
	buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y + rect.Height - 1}, bottomRight, buffer.Style{Foreground: "#238"})

	// Side borders
	for i := 1; i < rect.Height-1; i++ {
		buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y + i}, vertical, buffer.Style{Foreground: "#238"})
		buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y + i}, vertical, buffer.Style{Foreground: "#238"})
	}

	// Draw title
	title := " System Logs "
	buf.WriteString(buffer.Point{X: rect.X + 2, Y: rect.Y}, title, buffer.Style{Foreground: "#86", Bold: true})

	// Draw log entries
	y := rect.Y + 2
	maxY := rect.Y + rect.Height - 1
	for _, entry := range l.entries {
		if y >= maxY {
			break
		}
		// Truncate if too long
		maxLen := rect.Width - 3
		if len(entry) > maxLen {
			entry = entry[:maxLen] + "..."
		}
		buf.WriteString(buffer.Point{X: rect.X + 1, Y: y}, entry, buffer.Style{Foreground: "#250"})
		y++
	}
}

func (l *LogPanel) MinSize() (int, int)     { return 30, 8 }
func (l *LogPanel) PreferredSize() (int, int) { return 50, 12 }

// ChartPanel shows a simple vertical bar chart
type ChartPanel struct {
	title  string
	values []int
	labels []string
}

func (c *ChartPanel) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw title
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, c.title, buffer.Style{Foreground: "#86", Bold: true})

	// Find max value for scaling
	maxVal := 1
	for _, v := range c.values {
		if v > maxVal {
			maxVal = v
		}
	}

	// Draw chart
	chartHeight := rect.Height - 3
	chartWidth := rect.Width

	barWidth := (chartWidth / len(c.values)) - 1
	if barWidth < 1 {
		barWidth = 1
	}

	for i, val := range c.values {
		x := rect.X + i*(barWidth+1)
		barHeight := (val * chartHeight) / maxVal

		// Draw bar from bottom
		for h := 0; h < barHeight; h++ {
			y := rect.Y + rect.Height - 2 - h
			for w := 0; w < barWidth; w++ {
				buf.WriteString(buffer.Point{X: x + w, Y: y}, "█", buffer.Style{Foreground: "#33"})
			}
		}

		// Draw label
		if i < len(c.labels) {
			labelX := x + barWidth/2 - len(c.labels[i])/2
			buf.WriteString(buffer.Point{X: labelX, Y: rect.Y + rect.Height - 1}, c.labels[i], buffer.Style{Foreground: "#244"})
		}

		// Draw value on top
		valStr := fmt.Sprintf("%d", val)
		valX := x + barWidth/2 - len(valStr)/2
		buf.WriteString(buffer.Point{X: valX, Y: rect.Y + rect.Height - 2 - barHeight}, valStr, buffer.Style{Foreground: "#86", Bold: true})
	}
}

func (c *ChartPanel) MinSize() (int, int)     { return 30, 8 }
func (c *ChartPanel) PreferredSize() (int, int) { return 50, 15 }

func main() {
	width := 100
	height := 30

	fmt.Println(repeat("=", width))
	fmt.Println("Buffer Dashboard Example - Taproot")
	fmt.Println(repeat("=", width))
	fmt.Println()

	// Create layout manager
	lm := buffer.NewLayoutManager(width, height)

	// Header
	header := buffer.NewTextComponent(
		"  📊 System Dashboard  |  Taproot Buffer Layout System  |  ",
		buffer.Style{Bold: true, Foreground: "#86", Background: "#235"},
	)

	// Status row - three status panels
	statusOK := &StatusPanel{status: "RUNNING", message: "All systems operational", color: "#120"}
	statusDB := &StatusPanel{status: "CONNECTED", message: " PostgreSQL", color: "#120"}
	statusCache := &StatusPanel{status: "WARNING", message: "High memory", color: "#226"}

	// Metric panels - showing various metrics
	metricCPU := &MetricPanel{label: "CPU Usage", value: "23%", delta: "+2%"}
	metricMem := &MetricPanel{label: "Memory", value: "4.2 GB", delta: "-128 MB"}
	metricDisk := &MetricPanel{label: "Disk I/O", value: "45 MB/s", delta: "+12 MB/s"}
	metricNet := &MetricPanel{label: "Network", value: "120 Mbps", delta: "-5 Mbps"}

	// Progress panels
	progressBuild := &ProgressPanel{label: "Build Progress", percent: 75, total: 100, current: 75}
	progressDeploy := &ProgressPanel{label: "Deployment", percent: 30, total: 10, current: 3}

	// Log panel
	logPanel := &LogPanel{
		entries: []string{
			"[14:32:01] INFO  Build started (commit: abc123)",
			"[14:32:05] INFO  Tests passed: 245/245",
			"[14:32:08] WARN  High memory usage detected",
			"[14:32:12] INFO  Docker image built",
			"[14:32:15] INFO  Pushing to registry...",
			"[14:32:20] INFO  Deployment started",
			"[14:32:25] INFO  Rolling update: 3/10 pods",
		},
	}

	// Chart panel
	chartPanel := &ChartPanel{
		title:  "Request Rate (req/min)",
		values: []int{120, 245, 180, 320, 290, 410, 380},
		labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
	}

	// Add header
	lm.AddComponent("header", header)

	// Add status indicators in a row
	lm.AddComponent("status1", statusOK)
	lm.AddComponent("status2", statusDB)
	lm.AddComponent("status3", statusCache)

	// Add metrics
	lm.AddComponent("cpu", metricCPU)
	lm.AddComponent("mem", metricMem)
	lm.AddComponent("disk", metricDisk)
	lm.AddComponent("net", metricNet)

	// Add progress bars
	lm.AddComponent("progress1", progressBuild)
	lm.AddComponent("progress2", progressDeploy)

	// Add log panel and chart
	lm.AddComponent("logs", logPanel)
	lm.AddComponent("chart", chartPanel)

	// Calculate layout with custom positioning
	lm.CalculateLayout()

	// Get the layout rectangles and manually position components
	// This demonstrates fine-grained control over component placement
	buf := buffer.NewBuffer(width, height)

	// Header at top
	headerRect := buffer.Rect{X: 0, Y: 0, Width: width, Height: 1}
	header.Render(buf, headerRect)

	// Status row
	statusY := 2
	statusWidth := width / 3
	statusOK.Render(buf, buffer.Rect{X: 1, Y: statusY, Width: statusWidth - 2, Height: 1})
	statusDB.Render(buf, buffer.Rect{X: statusWidth, Y: statusY, Width: statusWidth - 2, Height: 1})
	statusCache.Render(buf, buffer.Rect{X: statusWidth * 2, Y: statusY, Width: statusWidth - 2, Height: 1})

	// Metrics row (4 metrics in a row)
	metricY := statusY + 2
	metricWidth := width / 4
	metricCPU.Render(buf, buffer.Rect{X: 2, Y: metricY, Width: metricWidth - 4, Height: 3})
	metricMem.Render(buf, buffer.Rect{X: metricWidth, Y: metricY, Width: metricWidth - 4, Height: 3})
	metricDisk.Render(buf, buffer.Rect{X: metricWidth * 2, Y: metricY, Width: metricWidth - 4, Height: 3})
	metricNet.Render(buf, buffer.Rect{X: metricWidth * 3, Y: metricY, Width: metricWidth - 4, Height: 3})

	// Progress bars row
	progressY := metricY + 4
	progressWidth := width / 2
	progressBuild.Render(buf, buffer.Rect{X: 2, Y: progressY, Width: progressWidth - 4, Height: 3})
	progressDeploy.Render(buf, buffer.Rect{X: progressWidth, Y: progressY, Width: progressWidth - 4, Height: 3})

	// Split bottom section: logs on left, chart on right
	bottomY := progressY + 4
	bottomHeight := height - bottomY - 1
	logPanel.Render(buf, buffer.Rect{X: 1, Y: bottomY, Width: width/2 - 2, Height: bottomHeight})
	chartPanel.Render(buf, buffer.Rect{X: width/2 + 1, Y: bottomY, Width: width/2 - 2, Height: bottomHeight})

	// Footer
	footerText := fmt.Sprintf(" Updated: %s | Components: 10 | Layout: Manual Positioning ", time.Now().Format("15:04:05"))
	buf.WriteString(buffer.Point{X: 0, Y: height - 1}, repeat(" ", width), buffer.Style{Background: "#235", Foreground: "#250"})
	buf.WriteString(buffer.Point{X: 2, Y: height - 1}, footerText, buffer.Style{Background: "#235", Foreground: "#250"})

	// Render and display
	output := buf.Render()
	fmt.Print(output)

	// Print summary
	fmt.Printf("\n\nDashboard Components:\n")
	fmt.Printf("  Header: 1 row\n")
	fmt.Printf("  Status Row: 1 row (3 panels)\n")
	fmt.Printf("  Metrics: 3 rows (4 panels)\n")
	fmt.Printf("  Progress: 3 rows (2 panels)\n")
	fmt.Printf("  Bottom Panel: %d rows (logs + chart)\n", bottomHeight)
	fmt.Printf("  Footer: 1 row\n")
	fmt.Printf("\nTotal: %d components in %d×%d grid (%d cells)\n", 11, width, height, width*height)
	fmt.Printf("\nKey Takeaway: Buffer layout enables pixel-perfect positioning\n")
	fmt.Printf("of multiple components with complex layouts.\n")
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
