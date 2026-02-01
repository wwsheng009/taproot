package tasks

import (
	"strings"
	"testing"
	"time"
)

func defaultTestTaskList() *TaskList {
	now := time.Now()
	tasks := []*Task{
		{
			ID:          "1",
			Title:       "Set up development environment",
			Description: "Install Go, configure IDE, clone repository",
			Status:      TaskStatusCompleted,
			Progress:    1.0,
			Priority:    5,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "2",
			Title:       "Implement core functionality",
			Description: "Build main features and APIs",
			Status:      TaskStatusInProgress,
			Progress:    0.65,
			Priority:    5,
			Assignee:    "Alice",
			Expanded:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
			Subtasks: []*Task{
				{ID: "2-1", Title: "Data models", Status: TaskStatusCompleted},
				{ID: "2-2", Title: "API endpoints", Status: TaskStatusCompleted},
				{ID: "2-3", Title: "UI components", Status: TaskStatusInProgress, Progress: 0.5},
			},
		},
		{
			ID:          "3",
			Title:       "Write unit tests",
			Description: "Add comprehensive test coverage",
			Status:      TaskStatusPending,
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
			Status:      TaskStatusBlocked,
			Progress:    0.3,
			Priority:    2,
			Assignee:    "Bob",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	return NewTaskList(tasks)
}

func TestNewTaskList(t *testing.T) {
	tl := NewTaskList([]*Task{})

	if tl == nil {
		t.Fatal("NewTaskList() returned nil")
	}

	if tl.focused {
		t.Error("Expected focused to be false initially")
	}

	if !tl.expanded {
		t.Error("Expected expanded to be true initially")
	}

	if tl.config == nil {
		t.Error("Expected config to be initialized")
	}
}

func TestTaskListDefaults(t *testing.T) {
	tasks := []*Task{
		{ID: "1", Title: "Test task", Status: TaskStatusPending},
	}
	tl := NewTaskList(tasks)

	if !tl.config.ShowProgress {
		t.Error("Expected ShowProgress to be true by default")
	}

	if !tl.config.ShowStatusIcons {
		t.Error("Expected ShowStatusIcons to be true by default")
	}

	if !tl.config.ShowTags {
		t.Error("Expected ShowTags to be true by default")
	}
}

func TestAddTask(t *testing.T) {
	tl := NewTaskList([]*Task{})

	task := &Task{
		ID:     "new-task",
		Title:  "New task",
		Status: TaskStatusPending,
	}

	tl.AddTask(task)

	if len(tl.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tl.tasks))
	}

	if tl.tasks[0].ID != "new-task" {
		t.Errorf("Expected task ID 'new-task', got %s", tl.tasks[0].ID)
	}
}

func TestRemoveTask(t *testing.T) {
	tl := defaultTestTaskList()

	initialCount := len(tl.tasks)
	if initialCount == 0 {
		t.Fatal("Expected tasks to be present")
	}

	removed := tl.RemoveTask("1")
	if !removed {
		t.Error("Expected RemoveTask to return true for existing task")
	}

	if len(tl.tasks) != initialCount-1 {
		t.Errorf("Expected %d tasks after removal, got %d", initialCount-1, len(tl.tasks))
	}

	// Try removing non-existent task
	removed = tl.RemoveTask("non-existent")
	if removed {
		t.Error("Expected RemoveTask to return false for non-existent task")
	}
}

func TestGetTask(t *testing.T) {
	tl := defaultTestTaskList()

	task := tl.GetTask("1")
	if task == nil {
		t.Error("Expected to find task with ID '1'")
	}

	task = tl.GetTask("non-existent")
	if task != nil {
		t.Error("Expected nil for non-existent task")
	}
}

func TestToggleExpanded(t *testing.T) {
	tl := defaultTestTaskList()

	task := tl.GetTask("2")
	if task == nil {
		t.Fatal("Expected task with ID '2' to exist")
	}

	initialExpanded := task.Expanded

	toggled := tl.ToggleExpanded("2")
	if !toggled {
		t.Error("Expected ToggleExpanded to return true")
	}

	if task.Expanded == initialExpanded {
		t.Error("Expected expanded state to change")
	}

	// Try toggling non-existent task
	toggled = tl.ToggleExpanded("non-existent")
	if toggled {
		t.Error("Expected ToggleExpanded to return false for non-existent task")
	}
}

func TestGetProgress(t *testing.T) {
	tl := defaultTestTaskList()

	progress := tl.GetProgress()
	if progress <= 0 {
		t.Errorf("Expected progress > 0, got %f", progress)
	}

	if progress > 1 {
		t.Errorf("Expected progress <= 1, got %f", progress)
	}
}

func TestGetStatusCounts(t *testing.T) {
	tl := defaultTestTaskList()

	counts := tl.GetStatusCounts()

	if counts == nil {
		t.Fatal("Expected counts to be initialized")
	}

	// We have 1 completed, 1 in progress, 1 pending, 1 blocked
	if counts[TaskStatusCompleted] != 1 {
		t.Errorf("Expected 1 completed task, got %d", counts[TaskStatusCompleted])
	}

	if counts[TaskStatusInProgress] != 1 {
		t.Errorf("Expected 1 in-progress task, got %d", counts[TaskStatusInProgress])
	}

	if counts[TaskStatusPending] != 1 {
		t.Errorf("Expected 1 pending task, got %d", counts[TaskStatusPending])
	}

	if counts[TaskStatusBlocked] != 1 {
		t.Errorf("Expected 1 blocked task, got %d", counts[TaskStatusBlocked])
	}
}

func TestFilterByStatus(t *testing.T) {
	tl := defaultTestTaskList()

	inProgressTasks := tl.FilterByStatus(TaskStatusInProgress)
	if len(inProgressTasks) != 1 {
		t.Errorf("Expected 1 in-progress task, got %d", len(inProgressTasks))
	}

	if inProgressTasks[0].Title != "Implement core functionality" {
		t.Errorf("Expected 'Implement core functionality', got %s", inProgressTasks[0].Title)
	}

	completedTasks := tl.FilterByStatus(TaskStatusCompleted)
	if len(completedTasks) != 1 {
		t.Errorf("Expected 1 completed task, got %d", len(completedTasks))
	}

	emptyTasks := tl.FilterByStatus(TaskStatusCancelled)
	if len(emptyTasks) != 0 {
		t.Errorf("Expected 0 cancelled tasks, got %d", len(emptyTasks))
	}
}

func TestFilterByTag(t *testing.T) {
	tl := defaultTestTaskList()

	testingTasks := tl.FilterByTag("testing")
	if len(testingTasks) != 1 {
		t.Errorf("Expected 1 task with 'testing' tag, got %d", len(testingTasks))
	}

	qualityTasks := tl.FilterByTag("quality")
	if len(qualityTasks) != 1 {
		t.Errorf("Expected 1 task with 'quality' tag, got %d", len(qualityTasks))
	}

	nonExistent := tl.FilterByTag("non-existent")
	if len(nonExistent) != 0 {
		t.Errorf("Expected 0 tasks with non-existent tag, got %d", len(nonExistent))
	}
}

func TestFocusAndBlur(t *testing.T) {
	tl := defaultTestTaskList()

	if tl.Focused() {
		t.Error("Expected to start unfocused")
	}

	tl.Focus()
	if !tl.Focused() {
		t.Error("Expected to be focused after Focus()")
	}

	tl.Blur()
	if tl.Focused() {
		t.Error("Expected to be unfocused after Blur()")
	}
}

func TestSetWidth(t *testing.T) {
	tl := defaultTestTaskList()

	tl.SetWidth(120)
	if tl.width != 120 {
		t.Errorf("Expected width 120, got %d", tl.width)
	}

	if tl.cacheValid {
		t.Error("Expected cache to be invalidated after SetWidth")
	}
}

func TestSetConfig(t *testing.T) {
	tl := defaultTestTaskList()

	newConfig := &TaskListConfig{
		ShowProgress:     false,
		ShowPriority:     true,
		ShowAssignee:     true,
		ShowDueDate:      true,
		ShowTags:         false,
		ShowSubtasks:     false,
		CompactMode:      true,
		AnimationEnabled: false,
	}

	tl.SetConfig(newConfig)

	if tl.config.ShowProgress {
		t.Error("Expected ShowProgress to be false")
	}

	if !tl.config.ShowPriority {
		t.Error("Expected ShowPriority to be true")
	}

	if tl.cacheValid {
		t.Error("Expected cache to be invalidated after SetConfig")
	}
}

func TestView(t *testing.T) {
	tl := defaultTestTaskList()

	view := tl.View()
	if view == "" {
		t.Error("Expected non-empty View output")
	}

	if !tl.cacheValid {
		t.Error("Expected cache to be valid after View()")
	}

	if tl.cached != view {
		t.Error("Expected cached view to match returned view")
	}
}

func TestStatusBarPendingIcon(t *testing.T) {
	tl := NewTaskList([]*Task{{ID: "1", Title: "Test"}})

	icon := tl.getStatusIcon(TaskStatusPending)
	if icon != "☐" {
		t.Errorf("Expected pending icon '☐', got %s", icon)
	}
}

func TestStatusBarInProgressIcon(t *testing.T) {
	tl := NewTaskList([]*Task{{ID: "1", Title: "Test"}})

	icon := tl.getStatusIcon(TaskStatusInProgress)
	if icon != "⟳" {
		t.Errorf("Expected in-progress icon '⟳', got %s", icon)
	}
}

func TestStatusBarCompletedIcon(t *testing.T) {
	tl := NewTaskList([]*Task{{ID: "1", Title: "Test"}})

	icon := tl.getStatusIcon(TaskStatusCompleted)
	if icon != "☑" {
		t.Errorf("Expected completed icon '☑', got %s", icon)
	}
}

func TestStatusBarBlockedIcon(t *testing.T) {
	tl := NewTaskList([]*Task{{ID: "1", Title: "Test"}})

	icon := tl.getStatusIcon(TaskStatusBlocked)
	if icon != "⚠" {
		t.Errorf("Expected blocked icon '⚠', got %s", icon)
	}
}

func TestStatusBarCancelledIcon(t *testing.T) {
	tl := NewTaskList([]*Task{{ID: "1", Title: "Test"}})

	icon := tl.getStatusIcon(TaskStatusCancelled)
	if icon != "✕" {
		t.Errorf("Expected cancelled icon '✕', got %s", icon)
	}
}

func TestGetStatusText(t *testing.T) {
	tl := NewTaskList([]*Task{{ID: "1", Title: "Test"}})

	tests := []struct {
		status TaskStatus
		text   string
	}{
		{TaskStatusPending, "Pending"},
		{TaskStatusInProgress, "In Progress"},
		{TaskStatusCompleted, "Done"},
		{TaskStatusBlocked, "Blocked"},
		{TaskStatusCancelled, "Cancelled"},
	}

	for _, tt := range tests {
		result := tl.getStatusText(tt.status)
		if result != tt.text {
			t.Errorf("For status %v, expected text '%s', got '%s'", tt.status, tt.text, result)
		}
	}
}

func TestTaskStatusString(t *testing.T) {
	tests := []struct {
		status TaskStatus
		text   string
	}{
		{TaskStatusPending, "pending"},
		{TaskStatusInProgress, "in_progress"},
		{TaskStatusCompleted, "completed"},
		{TaskStatusBlocked, "blocked"},
		{TaskStatusCancelled, "cancelled"},
	}

	for _, tt := range tests {
		result := tt.status.String()
		if result != tt.text {
			t.Errorf("For status %v, expected '%s', got '%s'", tt.status, tt.text, result)
		}
	}
}

func TestEmptyTaskList(t *testing.T) {
	tl := NewTaskList([]*Task{})

	view := tl.View()
	if view == "" {
		t.Error("Expected non-empty view for empty task list")
	}

	progress := tl.GetProgress()
	if progress != 0.0 {
		t.Errorf("Expected progress 0.0 for empty list, got %f", progress)
	}

	counts := tl.GetStatusCounts()
	if len(counts) != 0 {
		t.Errorf("Expected empty counts for empty list, got %d entries", len(counts))
	}
}

func TestRenderProgressBar(t *testing.T) {
	tl := NewTaskList([]*Task{{ID: "1", Title: "Test"}})

	// Test 50% progress
	bar := tl.renderProgressBar(0.5, 0)
	if bar == "" {
		t.Error("Expected non-empty progress bar")
	}

	// Test 0% progress
	bar = tl.renderProgressBar(0.0, 0)
	if !strings.Contains(bar, "0%") {
		t.Error("Expected '0%' in zero-progress bar")
	}

	// Test 100% progress
	bar = tl.renderProgressBar(1.0, 0)
	if !strings.Contains(bar, "100%") {
		t.Error("Expected '100%' in full-progress bar")
	}
}

func TestDefaultTaskListConfig(t *testing.T) {
	config := DefaultTaskListConfig()

	if !config.ShowProgress {
		t.Error("Expected ShowProgress to be true")
	}

	if config.ShowPriority {
		t.Error("Expected ShowPriority to be false")
	}

	if config.ShowAssignee {
		t.Error("Expected ShowAssignee to be false")
	}

	if !config.ShowTags {
		t.Error("Expected ShowTags to be true")
	}

	if !config.ShowSubtasks {
		t.Error("Expected ShowSubtasks to be true")
	}

	if config.CompactMode {
		t.Error("Expected CompactMode to be false")
	}

	if config.MaxWidth != 80 {
		t.Errorf("Expected MaxWidth 80, got %d", config.MaxWidth)
	}
}

func TestTaskFields(t *testing.T) {
	now := time.Now()
	task := &Task{
		ID:          "test-id",
		Title:       "Test Task",
		Description: "Test description",
		Status:      TaskStatusInProgress,
		Progress:    0.5,
		Priority:    3,
		Assignee:    "Test User",
		Tags:        []string{"tag1", "tag2"},
		Subtasks:    []*Task{{ID: "sub-1", Title: "Subtask"}},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if task.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", task.ID)
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", task.Title)
	}

	if task.Progress != 0.5 {
		t.Errorf("Expected progress 0.5, got %f", task.Progress)
	}

	if len(task.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(task.Tags))
	}
}

func TestGetTasks(t *testing.T) {
	tl := defaultTestTaskList()

	tasks := tl.GetTasks()
	if len(tasks) != 4 {
		t.Errorf("Expected 4 tasks, got %d", len(tasks))
	}

	if tasks[0].ID != "1" {
		t.Errorf("Expected first task ID '1', got '%s'", tasks[0].ID)
	}
}

