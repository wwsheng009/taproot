package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/components/pills"
)

const (
	maxWidth  = 80
	maxHeight = 24
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7c3aed")).
			Bold(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280")).
			Italic(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d1d5db"))
)

type model struct {
	pillList *pills.PillList
	width    int
	height   int
	config   pills.PillConfig
}

func initialModel() model {
	// Create sample pills
	samplePills := []*pills.Pill{
		{
			ID:       "1",
			Label:    "Tasks",
			Count:    5,
			Status:   pills.PillStatusPending,
			Expanded: false,
			Items:    []string{"Implement user authentication", "Design database schema", "Write API documentation", "Set up CI/CD pipeline", "Create unit tests"},
		},
		{
			ID:       "2",
			Label:    "In Progress",
			Count:    3,
			Status:   pills.PillStatusInProgress,
			Expanded: false,
			Items:    []string{"Build frontend components", "Integrate payment gateway", "Optimize database queries"},
		},
		{
			ID:       "3",
			Label:    "Completed",
			Count:    12,
			Status:   pills.PillStatusCompleted,
			Expanded: false,
			Items:    []string{"Setup project structure", "Define API endpoints", "Create database models", "Implement user registration", "Build dashboard UI", "Add dark mode support", "Implement search functionality", "Add file upload feature", "Create admin panel", "Setup logging system", "Add unit tests", "Deploy to staging"},
		},
		{
			ID:       "4",
			Label:    "Errors",
			Count:    2,
			Status:   pills.PillStatusError,
			Expanded: false,
			Items:    []string{"Fix memory leak in caching module", "Resolve race condition in concurrent updates"},
		},
		{
			ID:       "5",
			Label:    "Warnings",
			Count:    4,
			Status:   pills.PillStatusWarning,
			Expanded: false,
			Items:    []string{"Update deprecated dependencies", "Refactor legacy code", "Improve error handling", "Add input validation"},
		},
		{
			ID:       "6",
			Label:    "Info",
			Count:    1,
			Status:   pills.PillStatusInfo,
			Expanded: false,
			Items:    []string{"Check system logs for performance metrics"},
		},
	}

	pl := pills.NewPillList(samplePills)
	config := pills.DefaultPillConfig()

	return model{
		pillList: pl,
		width:    maxWidth,
		height:   maxHeight,
		config:   config,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "e":
			// Expand all pills
			m.pillList.ExpandAll()

		case "c":
			// Collapse all pills
			m.pillList.CollapseAll()

		case "1":
			// Toggle pill 1
			m.pillList.ToggleExpanded("1")

		case "2":
			// Toggle pill 2
			m.pillList.ToggleExpanded("2")

		case "3":
			// Toggle pill 3
			m.pillList.ToggleExpanded("3")

		case "4":
			// Toggle pill 4
			m.pillList.ToggleExpanded("4")

		case "5":
			// Toggle pill 5
			m.pillList.ToggleExpanded("5")

		case "6":
			// Toggle pill 6
			m.pillList.ToggleExpanded("6")

		case "a":
			// Add new pill
			newPill := &pills.Pill{
				ID:       fmt.Sprintf("%d", len(m.pillList.GetPills())+1),
				Label:    "New Category",
				Count:    1,
				Status:   pills.PillStatusInfo,
				Expanded: false,
				Items:    []string{fmt.Sprintf("New item %d", len(m.pillList.GetPills())+1)},
			}
			m.pillList.AddPill(newPill)

		case "r":
			// Remove last pill
			if len(m.pillList.GetPills()) > 0 {
				lastPill := m.pillList.GetPills()[len(m.pillList.GetPills())-1]
				m.pillList.RemovePill(lastPill.ID)
			}

		case "m":
			// Toggle inline mode
			m.config.InlineMode = !m.config.InlineMode
			m.pillList.SetConfig(m.config)

		case "i":
			// Toggle icons
			m.config.ShowIcons = !m.config.ShowIcons
			m.pillList.SetConfig(m.config)

		case "s":
			// Toggle count display
			m.config.ShowCount = !m.config.ShowCount
			m.pillList.SetConfig(m.config)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width > maxWidth {
			m.width = maxWidth
		}
		m.pillList.SetWidth(m.width - 4)
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Pill List Demo"))
	b.WriteString("\n\n")

	// Stats
	stats := fmt.Sprintf(
		"Total Pills: %d | Total Items: %d",
		len(m.pillList.GetPills()),
		m.pillList.GetTotalCount(),
	)
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Render(stats))
	b.WriteString("\n\n")

	// Pill list view
	b.WriteString(m.pillList.View())
	b.WriteString("\n\n")

	// Separator
	b.WriteString(separatorStyle.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")

	// Help text
	help := "Controls: 1-6toggle • eexpand all • ccollapse all • mmode • aadd • rremove • iicons • ssize • qquit"
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
