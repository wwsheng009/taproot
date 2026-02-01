package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// TaskStatus represents the status of a task.
type TaskStatus int

const (
	// TaskStatusPending means the task is not started.
	TaskStatusPending TaskStatus = iota
	// TaskStatusInProgress means the task is in progress.
	TaskStatusInProgress
	// TaskStatusCompleted means the task is completed.
	TaskStatusCompleted
	// TaskStatusBlocked means the task is blocked.
	TaskStatusBlocked
	// TaskStatusCancelled means the task was cancelled.
	TaskStatusCancelled
)

// String returns the string representation of task status.
func (s TaskStatus) String() string {
	switch s {
	case TaskStatusPending:
		return "pending"
	case TaskStatusInProgress:
		return "in_progress"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusBlocked:
		return "blocked"
	case TaskStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// Task represents a single task in the task list.
type Task struct {
	ID          string
	Title       string
	Description string
	Status      TaskStatus
	Progress    float64 // 0.0 to 1.0
	Priority    int     // 1-5, higher is more important
	Assignee    string
	DueDate     *time.Time
	Tags        []string
	Subtasks    []*Task
	Expanded    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TaskList represents a list of tasks with rendering and interaction.
type TaskList struct {
	tasks       []*Task
	width       int
	expanded    bool
	focused     bool
	config      *TaskListConfig
	styles      *styles.Styles
	cached      string
	cacheValid  bool
}

// TaskListConfig contains configuration for task list rendering.
type TaskListConfig struct {
	ShowProgress      bool
	ShowPriority      bool
	ShowAssignee      bool
	ShowDueDate       bool
	ShowTags          bool
	ShowSubtasks      bool
	ShowStatusIcons   bool
	CompactMode       bool
	MaxWidth          int
	AnimationEnabled  bool
}

// DefaultTaskListConfig returns the default task list configuration.
func DefaultTaskListConfig() *TaskListConfig {
	return &TaskListConfig{
		ShowProgress:     true,
		ShowPriority:     false,
		ShowAssignee:     false,
		ShowDueDate:      false,
		ShowTags:         true,
		ShowSubtasks:     true,
		ShowStatusIcons:  true,
		CompactMode:      false,
		MaxWidth:         80,
		AnimationEnabled: true,
	}
}

// NewTaskList creates a new TaskList component.
func NewTaskList(tasks []*Task) *TaskList {
	return &TaskList{
		tasks:      tasks,
		width:      80,
		expanded:   true,
		focused:    false,
		config:     DefaultTaskListConfig(),
		styles:     &styles.Styles{},
		cached:     "",
		cacheValid: false,
	}
}

// Init initializes the component. Implements render.Model.
func (tl *TaskList) Init() render.Cmd {
	return nil
}

// Update handles incoming messages. Implements render.Model.
func (tl *TaskList) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		tl.Focus()
	case *render.BlurMsg:
		tl.Blur()
	}

	tl.cacheValid = false
	return tl, render.None()
}

// View renders the task list. Implements render.Model.
func (tl *TaskList) View() string {
	if tl.cacheValid && tl.cached != "" {
		return tl.cached
	}

	var b strings.Builder

	if !tl.expanded {
		tl.cached = b.String()
		tl.cacheValid = true
		return tl.cached
	}

	if len(tl.tasks) == 0 {
		sty := tl.styles
		b.WriteString(sty.Subtle.Render("No tasks"))
		tl.cached = b.String()
		tl.cacheValid = true
		return tl.cached
	}

	for i, task := range tl.tasks {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(tl.renderTask(task, 0))
	}

	tl.cached = b.String()
	tl.cacheValid = true
	return tl.cached
}

// renderTask renders a single task with optional indentation.
func (tl *TaskList) renderTask(task *Task, indent int) string {
	sty := tl.styles
	var b strings.Builder

	// Indentation
	prefix := strings.Repeat("  ", indent)

	// Status row
	statusIcon := tl.getStatusIcon(task.Status)
	statusText := tl.getStatusText(task.Status)
	statusStyle := tl.getStatusStyle(task.Status)

	if tl.config.ShowStatusIcons {
		b.WriteString(prefix)
		b.WriteString(statusStyle.Render(statusIcon))
		b.WriteString(" ")

		// Task title
		titleStyle := sty.Base
		if task.Status == TaskStatusCompleted {
			titleStyle = titleStyle.Foreground(sty.FgMuted)
		} else if tl.focused {
			titleStyle = titleStyle.Foreground(sty.Primary)
		}
		b.WriteString(titleStyle.Render(task.Title))
	} else {
		b.WriteString(prefix)
		b.WriteString(statusStyle.Render(statusText))
		b.WriteString(": ")
		b.WriteString(sty.Base.Render(task.Title))
	}

	// Progress bar
	if tl.config.ShowProgress && task.Progress > 0 && task.Status != TaskStatusCompleted {
		b.WriteString("\n")
		b.WriteString(tl.renderProgressBar(task.Progress, indent))
	}

	// Description
	if task.Description != "" && !tl.config.CompactMode {
		b.WriteString("\n")
		b.WriteString(prefix)
		b.WriteString("  ")
		b.WriteString(sty.Subtle.Render(task.Description))
	}

	// Metadata
	var metadata []string
	if tl.config.ShowPriority && task.Priority > 0 {
		priority := strings.Repeat("!", task.Priority)
		metadata = append(metadata, sty.Base.Foreground(sty.Secondary).Render(priority))
	}
	if tl.config.ShowAssignee && task.Assignee != "" {
		metadata = append(metadata, sty.Subtle.Render(fmt.Sprintf("@%s", task.Assignee)))
	}
	if tl.config.ShowDueDate && task.DueDate != nil {
		due := time.Until(*task.DueDate)
		dueText := ""
		if due < 0 {
			dueText = "overdue"
		} else if due < 24*time.Hour {
			dueText = "due today"
		} else if due < 7*24*time.Hour {
			dueText = "due this week"
		}
		if dueText != "" {
			metadata = append(metadata, sty.Base.Foreground(sty.Error).Render(dueText))
		}
	}

	if len(metadata) > 0 {
		b.WriteString("\n")
		b.WriteString(prefix)
		b.WriteString("  ")
		b.WriteString(strings.Join(metadata, " • "))
	}

	// Tags
	if tl.config.ShowTags && len(task.Tags) > 0 {
		b.WriteString("\n")
		b.WriteString(prefix)
		b.WriteString("  ")
		for j, tag := range task.Tags {
			if j > 0 {
				b.WriteString(" ")
			}
			tagStyle := lipgloss.NewStyle().
				Foreground(sty.Secondary).
				Background(sty.BgSubtle).
				Padding(0, 1)
			b.WriteString(tagStyle.Render("#"+tag))
		}
	}

	// Subtasks
	if tl.config.ShowSubtasks && len(task.Subtasks) > 0 {
		subtaskCount := 0
		completedSubtasks := 0
		for _, sub := range task.Subtasks {
			subtaskCount++
			if sub.Status == TaskStatusCompleted {
				completedSubtasks++
			}
		}

		b.WriteString("\n")
		subtaskInfo := fmt.Sprintf("%d/%d subtasks", completedSubtasks, subtaskCount)
		if task.Expanded {
			b.WriteString(prefix)
			b.WriteString("  ")
			b.WriteString(sty.Subtle.Render("▼ "))
			b.WriteString(sty.Subtle.Render(subtaskInfo))
			

			for _, subtask := range task.Subtasks {
				b.WriteString("\n")
				b.WriteString(tl.renderTask(subtask, indent+1))
			}
		} else {
			b.WriteString(prefix)
			b.WriteString("  ")
			b.WriteString(sty.Subtle.Render("▶ "))
			b.WriteString(sty.Subtle.Render(subtaskInfo))
		}
	}

	return b.String()
}

// renderProgressBar renders a progress bar.
func (tl *TaskList) renderProgressBar(progress float64, indent int) string {
	sty := tl.styles
	prefix := strings.Repeat("  ", indent)

	barWidth := 30
	if tl.width > 0 {
		barWidth = tl.width/3
		if barWidth > 40 {
			barWidth = 40
		}
	}

	filled := int(progress * float64(barWidth))
	empty := barWidth - filled
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}

	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", empty)
	percent := int(progress * 100)

	barStyle := lipgloss.NewStyle().Foreground(sty.Primary)
	emptyStyle := lipgloss.NewStyle().Foreground(sty.BgSubtle)

	builder := strings.Builder{}
	builder.WriteString(prefix)
	builder.WriteString("  ")
	builder.WriteString(barStyle.Render(filledBar))
	builder.WriteString(emptyStyle.Render(emptyBar))
	builder.WriteString(" ")
	builder.WriteString(sty.Subtle.Render(fmt.Sprintf("%d%%", percent)))

	return builder.String()
}

// getStatusIcon returns the icon for a task status.
func (tl *TaskList) getStatusIcon(status TaskStatus) string {
	switch status {
	case TaskStatusPending:
		return "☐"
	case TaskStatusInProgress:
		return "⟳"
	case TaskStatusCompleted:
		return "☑"
	case TaskStatusBlocked:
		return "⚠"
	case TaskStatusCancelled:
		return "✕"
	default:
		return "•"
	}
}

// getStatusText returns the text for a task status.
func (tl *TaskList) getStatusText(status TaskStatus) string {
	switch status {
	case TaskStatusPending:
		return "Pending"
	case TaskStatusInProgress:
		return "In Progress"
	case TaskStatusCompleted:
		return "Done"
	case TaskStatusBlocked:
		return "Blocked"
	case TaskStatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

// getStatusStyle returns the style for a task status.
func (tl *TaskList) getStatusStyle(status TaskStatus) lipgloss.Style {
	sty := tl.styles
	switch status {
	case TaskStatusPending:
		return sty.Base.Foreground(sty.FgMuted)
	case TaskStatusInProgress:
		return sty.Base.Foreground(sty.Secondary)
	case TaskStatusCompleted:
		return sty.Base.Foreground(sty.Info)
	case TaskStatusBlocked:
		return sty.Base.Foreground(sty.Warning)
	case TaskStatusCancelled:
		return sty.Base.Foreground(sty.Error)
	default:
		return sty.Base
	}
}

// Focus focuses the component.
func (tl *TaskList) Focus() {
	tl.focused = true
	tl.cacheValid = false
}

// Blur blurs the component.
func (tl *TaskList) Blur() {
	tl.focused = false
	tl.cacheValid = false
}

// Focused returns true if the component is focused.
func (tl *TaskList) Focused() bool {
	return tl.focused
}

// SetWidth sets the width for rendering.
func (tl *TaskList) SetWidth(width int) {
	tl.width = width
	tl.cacheValid = false
}

// SetConfig sets the task list configuration.
func (tl *TaskList) SetConfig(config *TaskListConfig) {
	tl.config = config
	tl.cacheValid = false
}

// AddTask adds a task to the list.
func (tl *TaskList) AddTask(task *Task) {
	tl.tasks = append(tl.tasks, task)
	tl.cacheValid = false
}

// RemoveTask removes a task from the list by ID.
func (tl *TaskList) RemoveTask(id string) bool {
	for i, task := range tl.tasks {
		if task.ID == id {
			tl.tasks = append(tl.tasks[:i], tl.tasks[i+1:]...)
			tl.cacheValid = false
			return true
		}
	}
	return false
}

// GetTask retrieves a task by ID.
func (tl *TaskList) GetTask(id string) *Task {
	for _, task := range tl.tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

// ToggleExpanded toggles the expanded state of a task.
func (tl *TaskList) ToggleExpanded(id string) bool {
	task := tl.GetTask(id)
	if task != nil {
		task.Expanded = !task.Expanded
		tl.cacheValid = false
		return true
	}
	return false
}

// GetTasks returns all tasks in the list.
func (tl *TaskList) GetTasks() []*Task {
	return tl.tasks
}

// GetProgress returns the overall progress of the task list.
func (tl *TaskList) GetProgress() float64 {
	if len(tl.tasks) == 0 {
		return 0.0
	}

	totalProgress := 0.0
	for _, task := range tl.tasks {
		if task.Status == TaskStatusCompleted {
			totalProgress += 1.0
		} else {
			totalProgress += task.Progress
		}

		if len(task.Subtasks) > 0 {
			subtaskProgress := 0.0
			for _, sub := range task.Subtasks {
				if sub.Status == TaskStatusCompleted {
					subtaskProgress += 1.0
				}
			}
			totalProgress += subtaskProgress / float64(len(task.Subtasks))
		}
	}

	return totalProgress / float64(len(tl.tasks))
}

// GetStatusCounts returns the count of tasks by status.
func (tl *TaskList) GetStatusCounts() map[TaskStatus]int {
	counts := make(map[TaskStatus]int)
	for _, task := range tl.tasks {
		counts[task.Status]++
	}
	return counts
}

// FilterByStatus filters the task list by status.
func (tl *TaskList) FilterByStatus(status TaskStatus) []*Task {
	var filtered []*Task
	for _, task := range tl.tasks {
		if task.Status == status {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// FilterByTag filters the task list by tag.
func (tl *TaskList) FilterByTag(tag string) []*Task {
	var filtered []*Task
	for _, task := range tl.tasks {
		for _, t := range task.Tags {
			if t == tag {
				filtered = append(filtered, task)
				break
			}
		}
	}
	return filtered
}
