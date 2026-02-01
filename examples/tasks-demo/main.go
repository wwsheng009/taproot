package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/components/tasks"
	"github.com/wwsheng009/taproot/ui/styles"
)

// Model for the tasks demo
type model struct {
	taskList    *tasks.TaskList
	quitting    bool
	width       int
	statusMsg   string
}

func initialModel() model {
	now := time.Now()
	taskList := tasks.NewTaskList([]*tasks.Task{
		{
			ID:          "1",
			Title:       "Set up development environment",
			Description: "Install Go, configure IDE, clone repository",
			Status:      tasks.TaskStatusCompleted,
			Progress:    1.0,
			Priority:    5,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "2",
			Title:       "Implement core functionality",
			Description: "Build main features and APIs",
			Status:      tasks.TaskStatusInProgress,
			Progress:    0.65,
			Priority:    5,
			Assignee:    "Alice",
			Expanded:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
			Subtasks: []*tasks.Task{
				{ID: "2-1", Title: "Data models", Status: tasks.TaskStatusCompleted},
				{ID: "2-2", Title: "API endpoints", Status: tasks.TaskStatusCompleted},
				{ID: "2-3", Title: "UI components", Status: tasks.TaskStatusInProgress, Progress: 0.5},
			},
		},
		{
			ID:          "3",
			Title:       "Write unit tests",
			Description: "Add comprehensive test coverage",
			Status:      tasks.TaskStatusPending,
			Progress:    0.0,
			Priority:    3,
			Tags:        []string{"testing", "quality"},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "4",
			Title:       "Code review",
			Description: "Review and approve pull requests",
			Status:      tasks.TaskStatusBlocked,
			Progress:    0.3,
			Priority:    2,
			Assignee:    "Bob",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "5",
			Title:       "Deploy to production",
			Description: "Release version 1.0.0",
			Status:      tasks.TaskStatusPending,
			Progress:    0.0,
			Priority:    5,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	})
	taskList.SetWidth(60)
	taskList.Focus()

	return model{
		taskList:  taskList,
		quitting:  false,
		width:     80,
		statusMsg: "Use arrow keys to navigate, Enter to start tasks",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			// Start first pending task
			for _, task := range m.taskList.GetTasks() {
				if task.Status == tasks.TaskStatusPending {
					task.Status = tasks.TaskStatusInProgress
					task.Progress = 0.1
					task.UpdatedAt = time.Now()
					m.statusMsg = fmt.Sprintf("Started: %s", task.Title)
					return m, nil
				}
			}
			m.statusMsg = "No pending tasks"

		case tea.KeySpace:
			// Toggle expansion of first task with subtasks
			for _, task := range m.taskList.GetTasks() {
				if len(task.Subtasks) > 0 {
					m.taskList.ToggleExpanded(task.ID)
					m.statusMsg = fmt.Sprintf("Toggled: %s", task.Title)
					return m, nil
				}
			}
			m.statusMsg = "No tasks with subtasks"

		case tea.KeyUp:
			// Cycle first task status to previous
			task := m.taskList.GetTask("2")
			if task != nil {
				statuses := []tasks.TaskStatus{
					tasks.TaskStatusPending,
					tasks.TaskStatusInProgress,
					tasks.TaskStatusCompleted,
					tasks.TaskStatusBlocked,
					tasks.TaskStatusCancelled,
				}
				currentIdx := -1
				for i, s := range statuses {
					if s == task.Status {
						currentIdx = i
						break
					}
				}
				if currentIdx > 0 {
					task.Status = statuses[currentIdx-1]
					task.UpdatedAt = time.Now()
					m.statusMsg = fmt.Sprintf("Status: %s", task.Status)
				}
			}

		case tea.KeyDown:
			// Cycle first task status to next
			task := m.taskList.GetTask("2")
			if task != nil {
				statuses := []tasks.TaskStatus{
					tasks.TaskStatusPending,
					tasks.TaskStatusInProgress,
					tasks.TaskStatusCompleted,
					tasks.TaskStatusBlocked,
					tasks.TaskStatusCancelled,
				}
				currentIdx := -1
				for i, s := range statuses {
					if s == task.Status {
						currentIdx = i
						break
					}
				}
				if currentIdx >= 0 && currentIdx < len(statuses)-1 {
					task.Status = statuses[currentIdx+1]
					task.UpdatedAt = time.Now()
					m.statusMsg = fmt.Sprintf("Status: %s", task.Status)
				}
			}

		default:
			// Handle single character keys
			if len(msg.Runes) > 0 {
				switch strings.ToLower(string(msg.Runes[0])) {
				case "q":
					m.quitting = true
					return m, tea.Quit

				case "r":
					// Reset all tasks to pending
					for _, task := range m.taskList.GetTasks() {
						task.Status = tasks.TaskStatusPending
						task.Progress = 0.0
						task.UpdatedAt = time.Now()
					}
					m.statusMsg = "Reset all tasks to pending"

				case "p":
					// Show progress info
					progress := m.taskList.GetProgress()
					counts := m.taskList.GetStatusCounts()
					m.statusMsg = fmt.Sprintf("Progress: %.0f%% | Counts: %d done, %d in progress, %d pending",
						progress*100, counts[tasks.TaskStatusCompleted],
						counts[tasks.TaskStatusInProgress], counts[tasks.TaskStatusPending])
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = min(100, max(60, msg.Width))
		m.taskList.SetWidth(m.width - 10)
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	sty := styles.DefaultStyles()

	// Build help text
	help := "\n" + sty.Subtle.Render("Controls: ") +
		sty.Base.Render("Enter") + " start pending | " +
		sty.Base.Render("Space") + " toggle expansion | " +
		sty.Base.Render("↑/↓") + " change status\n           " +
		sty.Base.Render("R") + " reset all | " +
		sty.Base.Render("P") + " show progress | " +
		sty.Base.Render("Q") + " quit"

	// Build status bar
	statusBar := sty.Base.
		Width(m.width).
		Background(sty.BgSubtle).
		Foreground(sty.FgMuted).
		Render(" " + m.statusMsg + " ")

	// Build main content
	content := strings.Builder{}
	content.WriteString(sty.Base.Bold(true).Render(sty.Base.Foreground(sty.Primary).Render("Task List Demo")))
	content.WriteString("\n\n")
	content.WriteString(m.taskList.View())
	content.WriteString(help)

	// Wrap in border
	wrapped := sty.Base.
		Width(m.width).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(sty.Primary).
		Padding(1, 2).
		Render(content.String())

	// Combine with status bar
	final := lipgloss.JoinVertical(lipgloss.Left, wrapped, statusBar)

	return final
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
