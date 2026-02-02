package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/wwsheng009/taproot/ui/components/status"
	"github.com/wwsheng009/taproot/ui/components/messages"
	"github.com/wwsheng009/taproot/ui/components/progress"
)

// Model - Main dashboard model
type Model struct {
	navItems    []string
	navIndex    int
	panels      []string
	activePanel int

	// Component focus states
	servicesSelected   int // 0 for LSP, 1 for MCP
	activitySelected   int // Index in activity list
	filesSelected      int // Index in files list
	todosSelected      int // Index in todo items
	statBoxSelected    int // 0-3 for stat boxes

	// Components
	lspList    *status.LSPList
	mcpList    *status.MCPList
	activity   []*activityItem
	todos      *messages.TodoMessage
	progresses []*progressTask
	stats      *statsData
	files      []string

	// State
	width  int
	height int
	ticker *time.Ticker
}

// activityItem - Activity feed item
type activityItem struct {
	message   string
	timestamp time.Time
	icon      string
	color     lipgloss.Color
}

// progressTask - Progress task
type progressTask struct {
	name     string
	current  int64
	total    int64
	status   string
	bar      *progress.ProgressBar
}

// statsData - Statistics data
type statsData struct {
	filesWatched    int
	eventsProcessed int
	tasksCompleted  int
	uptime          time.Duration
}

// TickMsg - Tick message for updates
type TickMsg time.Time

func tickCmd(ticker *time.Ticker) tea.Cmd {
	return func() tea.Msg {
		return TickMsg(<-ticker.C)
	}
}

// Init - Initialize model
func (m Model) Init() tea.Cmd {
	// Start ticker for periodic updates
	if m.ticker == nil {
		m.ticker = time.NewTicker(2 * time.Second)
	}
	return tickCmd(m.ticker)
}

// Update - Update model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			// Panel navigation
			if m.activePanel < len(m.panels)-1 {
				m.activePanel++
			} else {
				m.activePanel = 0 // Wrap to first
			}
		case tea.KeyShiftTab:
			if m.activePanel > 0 {
				m.activePanel--
			} else {
				m.activePanel = len(m.panels) - 1 // Wrap to last
			}
		case tea.KeyUp:
			m.handleUp()
		case tea.KeyDown:
			m.handleDown()
		case tea.KeyLeft:
			m.handleLeft()
		case tea.KeyRight:
			m.handleRight()
		case tea.KeyEnter:
			m.handleEnter()
		}
		// Handle character keys
		switch msg.String() {
		case "q", "Q":
			return m, tea.Quit
		case "r", "R":
			m.resetStats()
		case "1":
			m.activePanel = 0
		case "2":
			m.activePanel = 1
		case "3":
			m.activePanel = 2
		case "4":
			m.activePanel = 3
		case "5":
			m.activePanel = 4
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case TickMsg:
		m.updateSimulation()
		return m, tickCmd(m.ticker)
	}

	return m, nil
}

// handleUp - Handle up arrow key
func (m *Model) handleUp() {
	switch m.panels[m.activePanel] {
	case "Services":
		m.servicesSelected = (m.servicesSelected + 1) % 2
	case "Activity":
		if m.activitySelected > 0 {
			m.activitySelected--
		}
	case "Files":
		if m.filesSelected > 0 {
			m.filesSelected--
		}
	case "Tasks":
		if m.todosSelected > 0 {
			m.todosSelected--
		}
	case "Stats":
		if m.statBoxSelected > 0 {
			m.statBoxSelected--
		}
	}
}

// handleDown - Handle down arrow key
func (m *Model) handleDown() {
	switch m.panels[m.activePanel] {
	case "Services":
		m.servicesSelected = (m.servicesSelected + 1) % 2
	case "Activity":
		if m.activitySelected < len(m.activity)-1 {
			m.activitySelected++
		}
	case "Files":
		if m.filesSelected < len(m.files)-1 {
			m.filesSelected++
		}
	case "Tasks":
		maxTodo := m.todos.TodoCount() - 1
		if maxTodo > 0 && m.todosSelected < maxTodo {
			m.todosSelected++
		}
	case "Stats":
		if m.statBoxSelected < 3 {
			m.statBoxSelected++
		}
	}
}

// handleLeft - Handle left arrow key
func (m *Model) handleLeft() {
	switch m.panels[m.activePanel] {
	case "Services":
		m.servicesSelected = 0
	case "Stats":
		if m.statBoxSelected > 0 {
			m.statBoxSelected--
		}
	}
}

// handleRight - Handle right arrow key
func (m *Model) handleRight() {
	switch m.panels[m.activePanel] {
	case "Services":
		m.servicesSelected = 1
	case "Stats":
		if m.statBoxSelected < 3 {
			m.statBoxSelected++
		}
	}
}

// handleEnter - Handle enter key
func (m *Model) handleEnter() {
	switch m.panels[m.activePanel] {
	case "Services":
		m.servicesSelected = (m.servicesSelected + 1) % 2
	case "Tasks":
		m.todos.ToggleExpanded()
	case "Stats":
		// Toggle stat box selection
		m.statBoxSelected = (m.statBoxSelected + 1) % 4
	}
}

// View - Render dashboard
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Main content area
	contentHeight := m.height - 10 // header + footer
	panels := m.renderPanels(contentHeight)

	// Split into main sidebar and content
	if m.width > 100 {
		// Wide layout - sidebar on left
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left,
			m.renderSidebar(20),
			panels,
		))
	} else {
		// Narrow layout - stacked
		b.WriteString(m.renderSidebar(20))
		b.WriteString("\n")
		b.WriteString(panels)
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(m.renderFooter())

	return b.String()
}

// renderHeader - Render header
func (m Model) renderHeader() string {
	title := styleTitle.Render("🚀 Taproot v2.0 Dashboard")
	status := styleStatus.Render(fmt.Sprintf("● Running | Uptime: %s | Panels: %d/%d",
		formatDuration(m.stats.uptime),
		m.activePanel+1,
		len(m.panels)))

	header := lipgloss.JoinHorizontal(lipgloss.Left, title, status)
	return styleHeader.Width(m.width).Render(header)
}

// renderSidebar - Render navigation sidebar
func (m Model) renderSidebar(width int) string {
	content := styleSectionBox.Render("Navigation")
	content += "\n"

	for i, item := range m.navItems {
		var baseStyle lipgloss.Style
		if i == m.activePanel {
			baseStyle = styleNavActive
		} else {
			baseStyle = styleNavItem
		}

		prefix := " "
		if i == m.activePanel {
			prefix = "▸ "
		}

		content += baseStyle.Width(width-4).Render(prefix+item)
		content += "\n"
	}

	content += "\n"
	content += styleSectionBox.Render("Quick Actions")
	content += "\n"
	content += styleNavItem.Render(" [R] Reset Stats") + "\n"
	content += styleNavItem.Render(" [Q] Quit") + "\n"

	return styleSidebar.Width(width).Height(m.height-8).Render(content)
}

// renderPanels - Render content panels
func (m Model) renderPanels(height int) string {
	panel := m.panels[m.activePanel]
	contentWidth := m.width - 24 // sidebar - margins

	var panelContent string
	switch panel {
	case "Services":
		panelContent = m.renderServicesPanel(contentWidth, height)
	case "Activity":
		panelContent = m.renderActivityPanel(contentWidth, height)
	case "Files":
		panelContent = m.renderFilesPanel(contentWidth, height)
	case "Tasks":
		panelContent = m.renderTasksPanel(contentWidth, height)
	case "Stats":
		panelContent = m.renderStatsPanel(contentWidth, height)
	}

	panelBox := stylePanelBox.Width(contentWidth).Height(height).Render(panelContent)

	// Panel tabs
	tabs := m.renderTabs(contentWidth)

	return lipgloss.JoinVertical(lipgloss.Left, tabs, panelBox)
}

// renderTabs - Render panel tabs
func (m Model) renderTabs(width int) string {
	tabWidth := width / len(m.panels)
	tabs := make([]string, len(m.panels))

	for i, panel := range m.panels {
		var style lipgloss.Style
		if i == m.activePanel {
			style = styleTabActive
		} else {
			style = styleTabInactive
		}

		tabs[i] = style.Width(tabWidth-2).Render(fmt.Sprintf("%d. %s", i+1, panel))
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
	return styleTabBar.Width(width).Render(tabRow)
}

// renderServicesPanel - Services panel
func (m Model) renderServicesPanel(width, height int) string {
	var b strings.Builder

	// LSP Services
	b.WriteString(styleSectionHeader.Render("LSP Services"))
	b.WriteString("\n")
	b.WriteString(m.lspList.View())
	b.WriteString("\n\n")

	// MCP Services
	b.WriteString(styleSectionHeader.Render("MCP Services"))
	b.WriteString("\n")
	b.WriteString(m.mcpList.View())

	return b.String()
}

// renderActivityPanel - Activity panel
func (m Model) renderActivityPanel(width, height int) string {
	var b strings.Builder

	b.WriteString(styleSectionHeader.Render("Recent Activity"))
	b.WriteString("\n\n")

	// Show last 10 activities
	maxItems := min(10, len(m.activity))
	for i := len(m.activity) - maxItems; i < len(m.activity); i++ {
		item := m.activity[i]
		timeAgo := formatTimeSince(item.timestamp)

		activityLine := lipgloss.JoinHorizontal(lipgloss.Left,
			styleIcon.Foreground(item.color).Render(item.icon),
			styleMessage.Render(item.message),
			styleTime.Render(fmt.Sprintf(" %s", timeAgo)),
		)

		b.WriteString(activityLine)
		b.WriteString("\n")
	}

	if len(m.activity) == 0 {
		b.WriteString(styleMuted.Render("No recent activity"))
	}

	return b.String()
}

// renderFilesPanel - Files panel
func (m Model) renderFilesPanel(width, height int) string {
	var b strings.Builder

	b.WriteString(styleSectionHeader.Render("File Monitoring"))
	b.WriteString("\n\n")

	// File list with selection
	for i, file := range m.files {
		prefix := "  "
		if i == m.filesSelected {
			prefix = "▸ "
		}

		if i == m.filesSelected {
			b.WriteString(styleSelected.Render(prefix + file))
		} else {
			b.WriteString(styleMuted.Render(prefix + file))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderTasksPanel - Tasks panel
func (m Model) renderTasksPanel(width, height int) string {
	var b strings.Builder

	b.WriteString(styleSectionHeader.Render("Active Tasks"))
	b.WriteString("\n\n")

	// Progress bars
	for _, task := range m.progresses {
		taskName := styleLabel.Render(task.name)
		taskStatus := fmt.Sprintf(" [%s]", styleStatus.Render(task.status))

		taskLine := lipgloss.JoinHorizontal(lipgloss.Top,
			taskName,
			taskStatus,
		)

		percentage := float64(task.current) / float64(task.total) * 100
		if task.bar != nil {
			task.bar.SetCurrent(float64(task.current))
			task.bar.SetTotal(float64(task.total))
		}

		b.WriteString(taskLine)
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  %.1f%%", percentage))
		b.WriteString("\n")
		if task.bar != nil {
			b.WriteString("  " + task.bar.View())
		}
		b.WriteString("\n\n")
	}

	// Todos
	b.WriteString(styleSectionHeader.Render("Todo List"))
	b.WriteString("\n")
	b.WriteString(m.todos.View())

	return b.String()
}

// renderStatsPanel - Stats panel
func (m Model) renderStatsPanel(width, height int) string {
	var b strings.Builder

	b.WriteString(styleSectionHeader.Render("Statistics Overview"))
	b.WriteString("\n")

	// Grid of stats
	stats := []struct {
		label string
		value string
		icon  string
		color lipgloss.Color
	}{
		{"Files Watched", fmt.Sprintf("%d", m.stats.filesWatched), "📁", lipgloss.Color("#36A64F")},
		{"Events Processed", fmt.Sprintf("%d", m.stats.eventsProcessed), "⚡", lipgloss.Color("#FF6B6B")},
		{"Tasks Completed", fmt.Sprintf("%d", m.stats.tasksCompleted), "✅", lipgloss.Color("#4ECDC4")},
		{"Uptime", formatDuration(m.stats.uptime), "⏱️", lipgloss.Color("#FFE66D")},
	}

	// Render stats in 2x2 grid using lipgloss Join
	var boxes []string
	for i, stat := range stats {
		var boxStyle lipgloss.Style
		if i == m.statBoxSelected && m.activePanel == 4 { // Stats panel
			boxStyle = styleStatBoxSelected
		} else {
			boxStyle = styleStatBox
		}

		box := boxStyle.Render(
			styleIcon.Foreground(stat.color).Render(stat.icon) + " " +
				styleStatLabel.Render(stat.label) + "\n" +
				styleStatValue.Foreground(stat.color).Render(stat.value),
		)
		boxes = append(boxes, box)
	}

	// First row: Files Watched + Events Processed
	row1 := lipgloss.JoinHorizontal(lipgloss.Bottom, boxes[0], "  ", boxes[1])
	// Second row: Tasks Completed + Uptime
	row2 := lipgloss.JoinHorizontal(lipgloss.Bottom, boxes[2], "  ", boxes[3])
	// Combine rows vertically without extra spacing
	grid := lipgloss.JoinVertical(lipgloss.Top, row1, row2)

	b.WriteString(grid)

	b.WriteString("\n")
	b.WriteString(styleSectionHeader.Render("Components Status"))
	b.WriteString("\n")
	b.WriteString("✓ LSP Service List - Online\n")
	b.WriteString("✓ MCP Service List - Online\n")
	b.WriteString("✓ Activity Feed - Active\n")
	b.WriteString("✓ File Monitor - Active\n")

	return b.String()
}

// renderFooter - Render footer
func (m Model) renderFooter() string {
	help := styleHelp.Render(
		fmt.Sprintf("[Tab] Panel  [↑↓←→] Navigate  [Enter] Action  [R] Reset  [Q] Quit"),
	)
	return styleFooter.Width(m.width).Render(help)
}

// updateSimulation - Update simulated data
func (m *Model) updateSimulation() {
	m.stats.uptime += 2 * time.Second

	// Random activity
	if rand.Float32() < 0.3 {
		m.addRandomActivity()
	}

	// Update progress
	for _, task := range m.progresses {
		if task.status == "Running" {
			task.current += int64(rand.Intn(10) + 1)
			if task.current >= task.total {
				task.current = task.total
				task.status = "Complete"
				m.stats.tasksCompleted++
				m.addActivity("Task completed", "✅", lipgloss.Color("#4ECDC4"))
			}
		}
	}

	// Update stats
	m.stats.filesWatched += rand.Intn(5)
	m.stats.eventsProcessed += rand.Intn(20)
}

// addRandomActivity - Add random activity
func (m *Model) addRandomActivity() {
	activities := []struct {
		message string
		icon    string
		color   lipgloss.Color
	}{
		{"File changed: main.go", "📝", lipgloss.Color("#FFD93D")},
		{"LSP diagnostic: 3 errors", "❌", lipgloss.Color("#FF6B6B")},
		{"MCP tool: filesystem.read", "🔨", lipgloss.Color("#4ECDC4")},
		{"Git status: 2 modified", "🔀", lipgloss.Color("#95E1D3")},
		{"Build completed: success", "✅", lipgloss.Color("#6BCB77")},
		{"Test run: passing", "🧪", lipgloss.Color("#A8E6CF")},
		{"Buffer saved", "💾", lipgloss.Color("#FFD93D")},
		{"Format applied", "✨", lipgloss.Color("#DDA0DD")},
	}

	for i := len(m.activity) - 1; i >= 0; i-- {
		if time.Since(m.activity[i].timestamp) > 5*time.Minute {
			m.activity = append(m.activity[:i], m.activity[i+1:]...)
		}
	}

	act := activities[rand.Intn(len(activities))]
	m.addActivity(act.message, act.icon, act.color)
}

// addActivity - Add activity entry
func (m *Model) addActivity(message string, icon string, color lipgloss.Color) {
	m.activity = append(m.activity, &activityItem{
		message:   message,
		timestamp: time.Now(),
		icon:      icon,
		color:     color,
	})
}

// resetStats - Reset statistics
func (m *Model) resetStats() {
	m.stats = &statsData{
		filesWatched:    0,
		eventsProcessed: 0,
		tasksCompleted:  0,
		uptime:          0,
	}
	m.activity = make([]*activityItem, 0)

	// Reset progress tasks
	for _, task := range m.progresses {
		task.current = 0
		task.status = "Running"
		if task.bar != nil {
			task.bar.SetCurrent(0)
		}
	}
}

// formatDuration - Format duration
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// formatTimeSince - Format time since
func formatTimeSince(t time.Time) string {
	d := time.Since(t)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}

// Styles
var (
	styleHeader = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color("#1e1e2e")).
		Foreground(lipgloss.Color("#cdd6f4"))

	styleTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		PaddingRight(2)

	styleStatus = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1"))

	styleSidebar = lipgloss.NewStyle().
		Padding(1).
		Margin(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#45475a"))

	styleNavActive = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cba6f7")).
		Bold(true).
		Padding(0, 1)

	styleNavItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Padding(0, 1)

	styleSectionBox = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#89b4fa")).
		MarginTop(1)

	stylePanelBox = lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#45475a")).
		Margin(1, 0, 0, 1)

	styleTabBar = lipgloss.NewStyle().
		Margin(1, 0, 0, 1)

	styleTabActive = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cba6f7")).
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color("#313244"))

	styleTabInactive = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Padding(0, 1)

	styleSectionHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#89b4fa")).
		Underline(true)

	styleIcon = lipgloss.NewStyle()

	styleMessage = lipgloss.NewStyle().
		PaddingLeft(1)

	styleTime = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		MarginLeft(2)

	styleMuted = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086"))

	styleNumber = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#a6e3a1"))

	styleLabel = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4"))

	styleStatBox = lipgloss.NewStyle().
		Width(25).
		Height(4).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#45475a"))

	styleStatBoxSelected = lipgloss.NewStyle().
		Width(25).
		Height(4).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#cba6f7")).
		BorderStyle(lipgloss.Border{
			TopLeft:     ">",
			TopRight:    "<",
			BottomRight: ">",
			BottomLeft:  "<",
		})

	styleStatLabel = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086"))

	styleStatValue = lipgloss.NewStyle().
		Bold(true)

	styleFooter = lipgloss.NewStyle().
		Padding(1).
		Background(lipgloss.Color("#1e1e2e")).
		Foreground(lipgloss.Color("#a6adc8"))

	styleHelp = lipgloss.NewStyle().
		Align(lipgloss.Center)

	styleSelected = lipgloss.NewStyle().
		Background(lipgloss.Color("#313244")).
		Foreground(lipgloss.Color("#cba6f7"))
)

// min - minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// NewModel - Create new dashboard model
func NewModel() Model {
	m := Model{
		navItems:    []string{"Services", "Activity", "Files", "Tasks", "Stats"},
		panels:      []string{"Services", "Activity", "Files", "Tasks", "Stats"},
		navIndex:    0,
		activePanel: 0,
		width:       80,
		height:      24,
		ticker:      time.NewTicker(2 * time.Second),
		stats: &statsData{
			filesWatched:    127,
			eventsProcessed: 8542,
			tasksCompleted:  42,
			uptime:          0,
		},
		activity: make([]*activityItem, 0),
		files: []string{
			"main.go",
			"dashboard.go",
			"services.go",
			"utils.go",
			"config.yaml",
		},
		servicesSelected: 0,
		activitySelected: 0,
		filesSelected:    0,
		todosSelected:    0,
		statBoxSelected:  0,
	}

	// Initialize LSP list
	m.lspList = status.NewLSPList()
	m.lspList.SetWidth(40)
	m.lspList.SetShowTitle(false)
	m.lspList.AddService(status.LSPServiceInfo{
		Name:     "gopls",
		Language: "go",
		State:    status.StateReady,
		Diagnostics: status.DiagnosticSummary{
			Error:   0,
			Warning: 2,
			Hint:    5,
		},
	})
	m.lspList.AddService(status.LSPServiceInfo{
		Name:     "rust-analyzer",
		Language: "rust",
		State:    status.StateReady,
		Diagnostics: status.DiagnosticSummary{
			Error:   0,
			Warning: 1,
		},
	})
	m.lspList.AddService(status.LSPServiceInfo{
		Name:     "clangd",
		Language: "c++",
		State:    status.StateStarting,
	})

	// Initialize MCP list
	m.mcpList = status.NewMCPList()
	m.mcpList.SetWidth(40)
	m.mcpList.SetShowTitle(false)
	m.mcpList.AddService(status.MCPServiceInfo{
		Name:  "filesystem",
		State: status.StateReady,
		ToolCounts: status.ToolCounts{
			Tools:   8,
			Prompts: 2,
		},
	})
	m.mcpList.AddService(status.MCPServiceInfo{
		Name:  "git",
		State: status.StateReady,
		ToolCounts: status.ToolCounts{
			Tools:   5,
			Prompts: 1,
		},
	})

	// Initialize todos
	m.todos = messages.NewTodoMessage("todos-1", "Dashboard")
	m.todos.AddTodo(messages.Todo{
		ID:          "todo-1",
		Description: "Review file changes",
		Status:      messages.TodoStatusPending,
	})
	m.todos.AddTodo(messages.Todo{
		ID:          "todo-2",
		Description: "Fix LSP diagnostics",
		Status:      messages.TodoStatusPending,
	})
	m.todos.AddTodo(messages.Todo{
		ID:          "todo-3",
		Description: "Deploy to production",
		Status:      messages.TodoStatusCompleted,
	})

	// Initialize progress tasks
	m.progresses = []*progressTask{
		{
			name:    "Build Project",
			current: 650,
			total:   1000,
			status:  "Running",
			bar:     progress.NewProgressBar(1000),
		},
		{
			name:    "Run Tests",
			current: 100,
			total:   100,
			status:  "Complete",
			bar:     progress.NewProgressBar(100),
		},
		{
			name:    "Generate Docs",
			current: 30,
			total:   200,
			status:  "Running",
			bar:     progress.NewProgressBar(200),
		},
	}

	// Add initial activity
	m.addActivity("Dashboard initialized", "🚀", lipgloss.Color("#CBA6F7"))
	m.addActivity("Connected to LSP services", "🔌", lipgloss.Color("#4ECDC4"))
	m.addActivity("File monitoring active", "👁️", lipgloss.Color("#95E1D3"))

	return m
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
