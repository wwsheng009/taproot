package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// TodoMessage represents a TODO list message containing task items with status tracking.
//
// This component supports:
// - Multiple todo items with different statuses (Pending, InProgress, Completed)
// - Progress tracking with percentage
// - Tags/labels for each todo
// - Collapsible todo list
// - Status icons and colors
// - Focus state with different styling
// - Caching for performance
type TodoMessage struct {
	id         string
	title      string
	todos      []Todo
	timestamp  time.Time
	focused    bool
	expanded   bool
	config     *MessageConfig
	maxWidth   int

	// Caching fields
	cachedRender string
	cachedWidth  int
	cachedHeight int
	cacheValid   bool
}

// NewTodoMessage creates a new TodoMessage component.
func NewTodoMessage(id, title string) *TodoMessage {
	return &TodoMessage{
		id:        id,
		title:     title,
		todos:     []Todo{},
		timestamp: time.Now(),
		focused:   false,
		expanded:  false,
		config:    &MessageConfig{},
	}
}

// NewTodoMessageFromTodos creates a new TodoMessage from existing todos.
func NewTodoMessageFromTodos(id, title string, todos []Todo) *TodoMessage {
	return &TodoMessage{
		id:        id,
		title:     title,
		todos:     todos,
		timestamp: time.Now(),
		focused:   false,
		expanded:  false,
		config:    &MessageConfig{},
	}
}

// Init initializes the component. Implements render.Model.
func (m *TodoMessage) Init() render.Cmd {
	return nil
}

// Update handles incoming messages. Implements render.Model.
func (m *TodoMessage) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		m.focused = true
	case *render.BlurMsg:
		m.focused = false
	}

	// Invalidate cache on any update
	m.cacheValid = false

	return m, render.None()
}

// View renders the todo message. Implements render.Model.
func (m *TodoMessage) View() string {
	width := m.config.MaxWidth
	// Apply max width limit
	if m.maxWidth > 0 && width > m.maxWidth {
		width = m.maxWidth
	}

	return m.render(width)
}

// ID returns the message ID. Implements Identifiable interface.
func (m *TodoMessage) ID() string {
	return m.id
}

// Role returns RoleTool (todos are system/tool messages). Implements Message interface.
func (m *TodoMessage) Role() MessageRole {
	return RoleTool
}

// Content returns a summary of todos. Implements Message interface.
func (m *TodoMessage) Content() string {
	completed := m.CompletedCount()
	return fmt.Sprintf("%s (%d/%d completed)", m.title, completed, len(m.todos))
}

// SetContent sets the title. Implements Message interface.
func (m *TodoMessage) SetContent(content string) {
	m.title = content
	m.cacheValid = false
}

// Timestamp returns the message timestamp. Implements Message interface.
func (m *TodoMessage) Timestamp() time.Time {
	return m.timestamp
}

// SetTimestamp sets the message timestamp.
func (m *TodoMessage) SetTimestamp(ts time.Time) {
	m.timestamp = ts
	m.cacheValid = false
}

// Title returns the todo message title.
func (m *TodoMessage) Title() string {
	return m.title
}

// SetTitle sets the title.
func (m *TodoMessage) SetTitle(title string) {
	m.title = title
	m.cacheValid = false
}

// Todos returns the todos.
func (m *TodoMessage) Todos() []Todo {
	return m.todos
}

// SetTodos sets the todos.
func (m *TodoMessage) SetTodos(todos []Todo) {
	m.todos = todos
	m.cacheValid = false
}

// AddTodo adds a todo item.
func (m *TodoMessage) AddTodo(todo Todo) {
	m.todos = append(m.todos, todo)
	m.cacheValid = false
}

// UpdateTodoStatus updates the status of a todo by ID.
func (m *TodoMessage) UpdateTodoStatus(todoID string, status TodoStatus) {
	for i := range m.todos {
		if m.todos[i].ID == todoID {
			m.todos[i].Status = status
			m.todos[i].Timestamp = time.Now()
			m.cacheValid = false
			return
		}
	}
}

// UpdateTodoProgress updates the progress of a todo by ID.
func (m *TodoMessage) UpdateTodoProgress(todoID string, progress float64) {
	for i := range m.todos {
		if m.todos[i].ID == todoID {
			m.todos[i].Progress = progress
			m.todos[i].Timestamp = time.Now()
			// Auto-update status based on progress
			if progress >= 1.0 {
				m.todos[i].Status = TodoStatusCompleted
			} else if progress > 0 {
				m.todos[i].Status = TodoStatusInProgress
			}
			m.cacheValid = false
			return
		}
	}
}

// Expanded returns true if the message is expanded.
func (m *TodoMessage) Expanded() bool {
	return m.expanded
}

// SetExpanded sets the expansion state.
func (m *TodoMessage) SetExpanded(expanded bool) {
	m.expanded = expanded
	m.cacheValid = false
}

// ToggleExpanded toggles the expansion state.
func (m *TodoMessage) ToggleExpanded() {
	m.expanded = !m.expanded
	m.cacheValid = false
}

// Focus focuses the component. Implements Focusable interface.
func (m *TodoMessage) Focus() {
	m.focused = true
	m.cacheValid = false
}

// Blur blurs the component. Implements Focusable interface.
func (m *TodoMessage) Blur() {
	m.focused = false
	m.cacheValid = false
}

// Focused returns true if the component is focused. Implements Focusable interface.
func (m *TodoMessage) Focused() bool {
	return m.focused
}

// SetMaxWidth sets the maximum width for rendering.
func (m *TodoMessage) SetMaxWidth(width int) {
	m.maxWidth = width
	m.cacheValid = false
}

// SetConfig sets the message configuration.
func (m *TodoMessage) SetConfig(config *MessageConfig) {
	m.config = config
	m.cacheValid = false
}

// render renders the todo message with the given width.
func (m *TodoMessage) render(width int) string {
	// Check cache first
	if m.cacheValid && m.cachedWidth == width && m.cachedRender != "" {
		return m.cachedRender
	}

	sty := styles.DefaultStyles()

	var builder strings.Builder

	// Render header
	header := m.renderHeader(&sty, width)
	builder.WriteString(header)

	// Render todos if expanded
	if m.expanded && len(m.todos) > 0 {
		for i, todo := range m.todos {
			if i > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(m.renderTodo(&sty, todo, width))
		}
	}

	// Add expand hint if collapsed and there are todos
	if !m.expanded && len(m.todos) > 0 {
		builder.WriteString("\n")
		expandHint := sty.Base.Foreground(sty.FgMuted).Italic(true).Render(
			fmt.Sprintf("[%d todo(s) - press enter to expand]", len(m.todos)),
		)
		builder.WriteString(expandHint)
	}

	result := builder.String()

	// Apply focus styling
	if m.focused {
		result = lipgloss.NewStyle().Foreground(sty.Secondary).Render(result)
	} else {
		result = lipgloss.NewStyle().Foreground(sty.FgMuted).Render(result)
	}

	// Cache the result
	m.cachedRender = result
	m.cachedWidth = width
	m.cachedHeight = lipgloss.Height(result)
	m.cacheValid = true

	return result
}

// renderHeader renders the todo message header with progress bar.
func (m *TodoMessage) renderHeader(sty *styles.Styles, width int) string {
	completed := m.CompletedCount()
	total := len(m.todos)
	progress := 0.0
	if total > 0 {
		progress = float64(completed) / float64(total)
	}

	// Todo icon
	todoIcon := "📋"
	headerText := fmt.Sprintf("%s %s", todoIcon, m.title)
	if total > 0 {
		headerText += fmt.Sprintf(" (%d/%d)", completed, total)
	}

	headerStyle := sty.Base.Bold(true).Foreground(sty.Primary)

	var parts []string
	parts = append(parts, headerStyle.Render(headerText))

	// Progress bar if there are todos
	if total > 0 {
		progressBar := m.renderProgressBar(sty, progress, width)
		parts = append(parts, progressBar)
	}

	return strings.Join(parts, "\n")
}

// renderProgressBar renders a progress bar.
func (m *TodoMessage) renderProgressBar(sty *styles.Styles, progress float64, maxContentWidth int) string {
	if maxContentWidth <= 0 {
		maxContentWidth = 40
	}

	barWidth := maxContentWidth - 2 // Leave room for []
	if barWidth < 10 {
		barWidth = 10
	}

	filled := int(progress * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	empty := barWidth - filled

	filledChar := "="
	emptyChar := " "

	bar := strings.Repeat(filledChar, filled) + strings.Repeat(emptyChar, empty)
	if filled == barWidth {
		bar += "✓"
	}

	progressStyle := sty.Base.Foreground(lipgloss.Color("#4ade80"))
	barStyle := progressStyle.Render(fmt.Sprintf("[%s]", bar))

	percentageText := fmt.Sprintf("%.0f%%", progress*100)
	percentageStyle := sty.Base.Foreground(sty.FgMuted).Italic(true)
	totalBar := barStyle + " " + percentageStyle.Render(percentageText)

	return totalBar
}

// renderTodo renders a single todo item.
func (m *TodoMessage) renderTodo(sty *styles.Styles, todo Todo, width int) string {
	var parts []string

	// Todo status icon and description
	statusIcon := getTodoStatusIcon(todo.Status)
	todoHeader := fmt.Sprintf("  %s %s", statusIcon, todo.Description)

	headerStyle := sty.Base.Bold(true)
	switch todo.Status {
	case TodoStatusPending:
		headerStyle = headerStyle.Foreground(sty.FgMuted)
	case TodoStatusInProgress:
		headerStyle = headerStyle.Foreground(sty.Secondary)
	case TodoStatusCompleted:
		headerStyle = headerStyle.Foreground(lipgloss.Color("#4ade80")).Strikethrough(true)
	}

	parts = append(parts, headerStyle.Render(todoHeader))

	// Progress bar if progress is set and between 0 and 1
	if todo.Progress > 0 && todo.Progress < 1.0 {
		progressBar := m.renderTodoProgressBar(sty, todo.Progress, width-4)
		if progressBar != "" {
			parts = append(parts, progressBar)
		}
	}

	// Tags if present
	if len(todo.Tags) > 0 {
		tags := m.renderTags(sty, todo.Tags)
		parts = append(parts, tags)
	}

	return strings.Join(parts, "\n")
}

// renderTodoProgressBar renders a progress bar for a single todo.
func (m *TodoMessage) renderTodoProgressBar(sty *styles.Styles, progress float64, maxContentWidth int) string {
	if maxContentWidth <= 0 {
		maxContentWidth = 20
	}

	barWidth := maxContentWidth - 2
	if barWidth < 5 {
		barWidth = 5
	}

	filled := int(progress * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	empty := barWidth - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	progressStyle := sty.Base.Foreground(sty.Secondary)
	return progressStyle.Render(fmt.Sprintf("    [%s]", bar))
}

// renderTags renders todo tags.
func (m *TodoMessage) renderTags(sty *styles.Styles, tags []string) string {
	tagStrings := make([]string, len(tags))
	for i, tag := range tags {
		tagStyle := sty.Base.
			Foreground(sty.FgMuted).
			Background(lipgloss.Color("#374151")).
			Padding(0, 1)
		tagStrings[i] = tagStyle.Render(tag)
	}

	return strings.Join(tagStrings, " ")
}

// getTodoStatusIcon returns the icon for a todo status.
func getTodoStatusIcon(status TodoStatus) string {
	switch status {
	case TodoStatusPending:
		return "⬜"
	case TodoStatusInProgress:
		return "🔄"
	case TodoStatusCompleted:
		return "✅"
	default:
		return "❓"
	}
}

// HandleClick handles click events for expanding/collapsing.
// Implements ClickHandler interface.
func (m *TodoMessage) HandleClick(line, col int) render.Cmd {
	if len(m.todos) > 0 {
		m.ToggleExpanded()
		return render.None()
	}
	return render.None()
}

// HandleKeyEvent handles keyboard events.
// Implements KeyEventHandler interface.
func (m *TodoMessage) HandleKeyEvent(msg any) (bool, render.Cmd) {
	keyMsg, ok := msg.(*render.KeyMsg)
	if !ok {
		return false, nil
	}

	switch keyMsg.String() {
	case " ", "enter":
		// Toggle expansion if there are todos
		if len(m.todos) > 0 {
			m.ToggleExpanded()
			return true, render.None()
		}
	}

	return false, nil
}

// TodoCount returns the total number of todos.
func (m *TodoMessage) TodoCount() int {
	return len(m.todos)
}

// CompletedCount returns the number of completed todos.
func (m *TodoMessage) CompletedCount() int {
	count := 0
	for _, todo := range m.todos {
		if todo.Status == TodoStatusCompleted {
			count++
		}
	}
	return count
}

// InProgressCount returns the number of in-progress todos.
func (m *TodoMessage) InProgressCount() int {
	count := 0
	for _, todo := range m.todos {
		if todo.Status == TodoStatusInProgress {
			count++
		}
	}
	return count
}

// PendingCount returns the number of pending todos.
func (m *TodoMessage) PendingCount() int {
	count := 0
	for _, todo := range m.todos {
		if todo.Status == TodoStatusPending {
			count++
		}
	}
	return count
}

// OverallProgress returns the overall progress (0.0 to 1.0).
func (m *TodoMessage) OverallProgress() float64 {
	total := len(m.todos)
	if total == 0 {
		return 0.0
	}

	sumProgress := 0.0
	for _, todo := range m.todos {
		sumProgress += todo.Progress
	}

	return sumProgress / float64(total)
}
