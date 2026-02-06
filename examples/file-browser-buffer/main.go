package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/components/status"
	"github.com/wwsheng009/taproot/ui/components/messages"
	"github.com/wwsheng009/taproot/ui/render/buffer"
	"github.com/wwsheng009/taproot/ui/markdown"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// PanelFocus - Which panel has focus
type PanelFocus int

const (
	FocusFileList PanelFocus = iota
	FocusPreview
)

// Model - File browser with Buffer Layout system
type Model struct {
	// Buffer layout system
	layoutManager *buffer.LayoutManager

	// Layout rectangles
	mainRect     buffer.Rect
	fileListRect buffer.Rect
	previewRect  buffer.Rect

	// Navigation
	currentPath string
	fileList    list.Model
	keyMap      keyMap

	// Panels
	panelFocus  PanelFocus
	lspList     *status.LSPList
	mcpList     *status.MCPList
	diagnostics []*messages.DiagnosticMessage
	todos       *messages.TodoMessage

	// File preview
	viewport       viewport.Model
	previewFile    string
	previewData    string
	previewType    PreviewType
	previewLoading bool

	// Command palette - using ! for command mode
	commandInput   textinput.Model
	showCommand    bool
	commandHistory []string
	historyIndex   int

	// Command output
	commandOutput viewport.Model
	outputLines   []string

	// Search
	searchInput textinput.Model
	showSearch  bool

	// UI
	width    int
	height   int
	quitting bool

	// Content cache
	contentCache  map[string]string
	cacheMutex    sync.RWMutex
	lastPreview   time.Time
	previewThrottler *time.Timer

	// Resizable panels
	fileListWidthRatio float64
}

// PreviewType - Type of file preview
type PreviewType int

const (
	PreviewNone PreviewType = iota
	PreviewText
	PreviewMarkdown
	PreviewBinary
)

// keyMap - Application key bindings
type keyMap struct {
	Quit        key.Binding
	Search      key.Binding
	Command     key.Binding
	Refresh     key.Binding
	Up          key.Binding
	Down        key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	Home        key.Binding
	End         key.Binding
	Left        key.Binding
	Right       key.Binding
	Help        key.Binding
	GoToParent  key.Binding
	TogglePanel key.Binding
	ResizePanel key.Binding
}

var DefaultKeyMap = keyMap{
	Quit:        key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("ctrl+c/q", "quit")),
	Search:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	Command:     key.NewBinding(key.WithKeys("!"), key.WithHelp("!", "command")),
	Refresh:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	PageUp:      key.NewBinding(key.WithKeys("pgup", "shift+up", "ctrl+u"), key.WithHelp("PgUp", "page up")),
	PageDown:    key.NewBinding(key.WithKeys("pgdown", "shift+down", "ctrl+d"), key.WithHelp("PgDn", "page down")),
	Home:        key.NewBinding(key.WithKeys("home", "ctrl+a"), key.WithHelp("Home", "top")),
	End:         key.NewBinding(key.WithKeys("end", "ctrl+e"), key.WithHelp("End", "bottom")),
	Left:        key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "focus left")),
	Right:       key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "focus right")),
	Help:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	GoToParent:  key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "parent dir")),
	TogglePanel: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "toggle panel")),
	ResizePanel: key.NewBinding(key.WithKeys("[", "]"), key.WithHelp("[/]", "resize")),
}

// FileItem - File list item
type FileItem struct {
	name      string
	path      string
	isDir     bool
	size      int64
	modTime   string
	extension string
}

// Implement list.Item interface
func (i FileItem) Title() string       { return i.name }
func (i FileItem) Description() string { return i.description() }
func (i FileItem) FilterValue() string { return i.name }

func (i FileItem) description() string {
	if i.isDir {
		return "DIR"
	}
	return fmt.Sprintf("%s %s", formatSize(i.size), i.modTime)
}

// formatSize - Format file size
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Buffer-based component for header
type HeaderComponent struct {
	currentPath string
	color       lipgloss.Color
}

func NewHeaderComponent(path string) *HeaderComponent {
	return &HeaderComponent{
		currentPath: path,
		color:       lipgloss.Color("#cba6f7"),
	}
}

func (h *HeaderComponent) Render() (*buffer.Buffer, error) {
	text := fmt.Sprintf(" 📁 File Browser (Buffer Layout) - %s ", h.currentPath)
	buf := buffer.NewBuffer(len(text), 1)

	// Create text buffer
	for i, char := range text {
		buf.SetCell(buffer.Point{X: i, Y: 0}, buffer.Cell{
			Char: char,
			Style: buffer.Style{Foreground: "#cba6f7", Background: "#1e1e2e", Bold: true},
		})
	}

	return buf, nil
}

func (h *HeaderComponent) SetPath(path string) {
	h.currentPath = path
}

// Buffer-based component for panels
type PanelComponent struct {
	id      string
	title   string
	content string
	focused bool
}

func NewPanelComponent(id, title, content string, focused bool) *PanelComponent {
	return &PanelComponent{
		id:      id,
		title:   title,
		content: content,
		focused: focused,
	}
}

func (p *PanelComponent) Render() (*buffer.Buffer, error) {
	lines := strings.Split(p.content, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}

	minWidth := len(p.title)
	for _, line := range lines {
		if len(line) > minWidth {
			minWidth = len(line)
		}
	}

	contentWidth := minWidth + 2
	contentHeight := len(lines) + 2

	buf := buffer.NewBuffer(contentWidth+2, contentHeight+2)

	// Render border
	borderColor := "#45475a"
	if p.focused {
		borderColor = "#cba6f7"
	}

	// Draw border corners and edges
	for y := 0; y < buf.Height(); y++ {
		for x := 0; x < buf.Width(); x++ {
			if y == 0 || y == buf.Height()-1 {
				buf.SetCell(buffer.Point{X: x, Y: y}, buffer.Cell{
					Char:  '─',
					Style: buffer.Style{Foreground: borderColor},
				})
			} else if x == 0 || x == buf.Width()-1 {
				buf.SetCell(buffer.Point{X: x, Y: y}, buffer.Cell{
					Char:  '│',
					Style: buffer.Style{Foreground: borderColor},
				})
			}
		}
	}

	// Draw corners
	buf.SetCell(buffer.Point{X: 0, Y: 0}, buffer.Cell{Char: '╭', Style: buffer.Style{Foreground: borderColor}})
	buf.SetCell(buffer.Point{X: buf.Width() - 1, Y: 0}, buffer.Cell{Char: '╮', Style: buffer.Style{Foreground: borderColor}})
	buf.SetCell(buffer.Point{X: 0, Y: buf.Height() - 1}, buffer.Cell{Char: '╰', Style: buffer.Style{Foreground: borderColor}})
	buf.SetCell(buffer.Point{X: buf.Width() - 1, Y: buf.Height() - 1}, buffer.Cell{Char: '╯', Style: buffer.Style{Foreground: borderColor}})

	// Draw title at center top
	titleX := (buf.Width() - len(p.title)) / 2
	if titleX < 0 {
		titleX = 0
	}
	titleY := 0

	for i, char := range p.title {
		buf.SetCell(buffer.Point{X: titleX + i, Y: titleY}, buffer.Cell{
			Char: char,
			Style: buffer.Style{Foreground: "#cba6f7", Background: "#1e1e2e"},
		})
	}

	// Draw content with padding
	for i, line := range lines {
		if len(line) > contentWidth {
			line = line[:contentWidth]
		}
		buf.WriteString(buffer.Point{X: 1, Y: i + 1}, line, buffer.Style{})
	}

	return buf, nil
}

func (p *PanelComponent) SetContent(content string) {
	p.content = content
}

func (p *PanelComponent) SetFocused(focused bool) {
	p.focused = focused
}

func (p *PanelComponent) TitleHeight() int {
	return 1
}

// Buffer-based component for footer
type FooterComponent struct {
	help    string
	color   lipgloss.Color
}

func NewFooterComponent(help string) *FooterComponent {
	return &FooterComponent{
		help:  help,
		color: lipgloss.Color("#6c7086"),
	}
}

func (f *FooterComponent) Render() (*buffer.Buffer, error) {
	buf := buffer.NewBuffer(len(f.help), 1)

	for i, char := range f.help {
		buf.SetCell(buffer.Point{X: i, Y: 0}, buffer.Cell{
			Char: char,
			Style: buffer.Style{Foreground: "#6c7086", Background: "#1e1e2e"},
		})
	}

	return buf, nil
}

// NewModel - Create new application model
func NewModel() Model {
	// Initialize file list
	wd, _ := os.Getwd()
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#cba6f7"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#a6adc8"))

	fileList := list.New(nil, delegate, 0, 0)
	fileList.Title = "Files"

	// Initialize viewports
	vp := viewport.New(0, 0)
	cmdOutput := viewport.New(0, 0)

	// Initialize command input (! for command mode)
	commandInput := textinput.New()
	commandInput.Placeholder = "Command..."
	commandInput.Prompt = "! "

	// Initialize search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search..."
	searchInput.Prompt = "/ "

	// Initialize LSP list
	lspList := status.NewLSPList()
	lspList.SetWidth(40)
	lspList.AddService(status.LSPServiceInfo{
		Name:     "gopls",
		Language: "go",
		State:    status.StateReady,
		Diagnostics: status.DiagnosticSummary{
			Error:   0,
			Warning: 3,
			Hint:    8,
		},
	})
	lspList.AddService(status.LSPServiceInfo{
		Name:     "rust-analyzer",
		Language: "rust",
		State:    status.StateReady,
		Diagnostics: status.DiagnosticSummary{
			Error:   0,
			Warning: 1,
		},
	})

	// Initialize MCP list
	mcpList := status.NewMCPList()
	mcpList.SetWidth(40)
	mcpList.AddService(status.MCPServiceInfo{
		Name:  "filesystem",
		State: status.StateReady,
		ToolCounts: status.ToolCounts{
			Tools:   12,
			Prompts: 3,
		},
	})
	mcpList.AddService(status.MCPServiceInfo{
		Name:  "git",
		State: status.StateReady,
		ToolCounts: status.ToolCounts{
			Tools:   8,
			Prompts: 1,
		},
	})

	// Initialize todos
	todos := messages.NewTodoMessage("main-todos", "Active Tasks")
	todos.AddTodo(messages.Todo{
		ID:          "todo-1",
		Description: "Review PR #123",
		Status:      messages.TodoStatusPending,
	})
	todos.AddTodo(messages.Todo{
		ID:          "todo-2",
		Description: "Update documentation",
		Status:      messages.TodoStatusInProgress,
		Progress:    0.6,
	})
	todos.AddTodo(messages.Todo{
		ID:          "todo-3",
		Description: "Fix navigation bug",
		Status:      messages.TodoStatusCompleted,
	})

	return Model{
		currentPath:        wd,
		fileList:           fileList,
		fileListWidthRatio: 0.3,
		keyMap:             DefaultKeyMap,
		panelFocus:         FocusFileList,
		lspList:            lspList,
		mcpList:            mcpList,
		todos:              todos,
		commandInput:       commandInput,
		searchInput:        searchInput,
		viewport:           vp,
		commandOutput:      cmdOutput,
		previewFile:        "",
		previewData:        "",
		previewType:        PreviewNone,
		previewLoading:     false,
		outputLines:        []string{"Ready. Type ! to enter command mode."},
		contentCache:       make(map[string]string),
		commandHistory:     make([]string, 0),
		historyIndex:       -1,
		width:              80,
		height:             24,
		quitting:           false,
		previewThrottler:   nil,
	}
}

// Init - Initialize model
func (m Model) Init() tea.Cmd {
	// Load initial directory
	return m.loadDirectory(m.currentPath)
}

// Update - Update model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle command mode
		if m.showCommand {
			return m.updateCommandMode(msg)
		}

		// Handle search mode
		if m.showSearch {
			return m.updateSearchMode(msg)
		}

		// Normal mode
		return m.updateNormalMode(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Force recalculate layout on resize
		m.recalculateLayout()
		m.updateComponentSizes()

	case []list.Item:
		// File list loaded
		m.fileList.SetItems(msg)
	}

	return m, tea.Batch(cmds...)
}

// recalculateLayout - Simple clean layout: header + main area + footer
func (m *Model) recalculateLayout() {
	headerHeight := 1
	footerHeight := 1
	mainHeight := m.height - headerHeight - footerHeight
	
	if mainHeight < 10 {
		mainHeight = 10
	}

	// Split main area: 30% for file list, 70% for preview
	fileListWidth := int(float64(m.width) * 0.3)
	previewWidth := m.width - fileListWidth

	// Store areas - mainHeight is the height for panels, NOT the full screen
	m.mainRect = buffer.Rect{X: 0, Y: 0, Width: m.width, Height: mainHeight}
	m.fileListRect = buffer.Rect{X: 0, Y: 0, Width: fileListWidth, Height: mainHeight}
	m.previewRect = buffer.Rect{X: fileListWidth, Y: 0, Width: previewWidth, Height: mainHeight}
}

// updateComponentSizes - Update component sizes
func (m *Model) updateComponentSizes() {
	// File list: account for border and padding
	listWidth := m.fileListRect.Width - 4
	listHeight := m.fileListRect.Height - 4
	if listWidth > 0 {
		m.fileList.SetSize(listWidth, listHeight)
	}

	// Preview/output: account for border (2) + title bar (1)
	previewWidth := m.previewRect.Width - 2
	previewHeight := m.previewRect.Height - 3  // border(2) + title(1)
	if previewWidth > 0 && previewHeight > 0 {
		m.viewport.Width = previewWidth
		m.viewport.Height = previewHeight
		m.commandOutput.Width = previewWidth
		m.commandOutput.Height = previewHeight
	}
}

// updateCommandMode - Handle command mode input
func (m Model) updateCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.commandInput, cmd = m.commandInput.Update(msg)

	switch msg.Type {
	case tea.KeyEnter:
		command := m.commandInput.Value()
		if command != "" {
			m.commandHistory = append(m.commandHistory, command)
			m.historyIndex = len(m.commandHistory) - 1
			return m.executeCommand(command)
		}
		return m, nil

	case tea.KeyEsc:
		m.showCommand = false
		m.commandInput.Reset()
		m.commandInput.Blur()
		return m, nil

	case tea.KeyUp:
		if m.historyIndex > 0 {
			m.historyIndex--
			m.commandInput.SetValue(m.commandHistory[m.historyIndex])
		}

	case tea.KeyDown:
		if m.historyIndex < len(m.commandHistory)-1 {
			m.historyIndex++
			m.commandInput.SetValue(m.commandHistory[m.historyIndex])
		} else if m.historyIndex == len(m.commandHistory)-1 {
			m.historyIndex++
			m.commandInput.SetValue("")
		}
	}

	return m, cmd
}

// updateSearchMode - Handle search mode input
func (m Model) updateSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	query := m.searchInput.Value()
	if query != "" {
		m.fileList.FilterInput.SetValue(query)
	} else {
		m.fileList.ResetFilter()
	}

	switch msg.Type {
	case tea.KeyEnter:
		m.showSearch = false
		m.searchInput.Blur()
		return m, nil

	case tea.KeyEsc:
		m.showSearch = false
		m.searchInput.Reset()
		m.fileList.ResetFilter()
		m.searchInput.Blur()
		return m, nil
	}

	return m, cmd
}

// updateNormalMode - Handle normal mode input
func (m Model) updateNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch {
	case key.Matches(msg, m.keyMap.Quit):
		m.quitting = true
		return m, tea.Quit

	case key.Matches(msg, m.keyMap.Search):
		m.showSearch = true
		m.showCommand = false
		m.searchInput.Focus()
		m.searchInput.Reset()
		m.commandInput.Blur()
		m.fileList.ResetFilter()
		return m, nil

	case key.Matches(msg, m.keyMap.Command):
		m.showCommand = true
		m.showSearch = false
		m.commandInput.Focus()
		m.commandInput.Reset()
		m.commandInput.CursorEnd()
		m.searchInput.Blur()
		return m, nil

	case key.Matches(msg, m.keyMap.Refresh):
		return m, m.loadDirectory(m.currentPath)

	case key.Matches(msg, m.keyMap.Help):
		m = m.showHelp()
		return m, nil

	case key.Matches(msg, m.keyMap.GoToParent):
		parent := filepath.Dir(m.currentPath)
		if parent != m.currentPath {
			m.currentPath = parent
			return m, m.loadDirectory(m.currentPath)
		}
		return m, nil

	case key.Matches(msg, m.keyMap.TogglePanel):
		m.panelFocus = (m.panelFocus + 1) % 2
		return m, nil

	case key.Matches(msg, m.keyMap.ResizePanel):
		if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
			r := msg.Runes[0]
			if r == '[' {
				if m.fileListWidthRatio > 0.15 {
					m.fileListWidthRatio -= 0.05
					m.recalculateLayout()
					m.updateComponentSizes()
				}
			} else if r == ']' {
				if m.fileListWidthRatio < 0.5 {
					m.fileListWidthRatio += 0.05
					m.recalculateLayout()
					m.updateComponentSizes()
				}
			}
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Left):
		if m.panelFocus != FocusFileList {
			m.panelFocus = FocusFileList
			return m, nil
		}
		var listCmd tea.Cmd
		m.fileList, listCmd = m.fileList.Update(msg)
		return m, listCmd

	case key.Matches(msg, m.keyMap.Right):
		if m.panelFocus == FocusFileList {
			m.panelFocus = FocusPreview
			return m, nil
		}
		var viewCmd tea.Cmd
		m.viewport, viewCmd = m.viewport.Update(msg)
		return m, viewCmd

	case key.Matches(msg, m.keyMap.Up), key.Matches(msg, m.keyMap.Down),
	     key.Matches(msg, m.keyMap.PageUp), key.Matches(msg, m.keyMap.PageDown),
	     key.Matches(msg, m.keyMap.Home), key.Matches(msg, m.keyMap.End):
		var cmd tea.Cmd

		switch m.panelFocus {
		case FocusFileList:
			m.fileList, cmd = m.fileList.Update(msg)

			if m.fileList.SelectedItem() != nil {
				selected := m.fileList.SelectedItem().(FileItem)
				if !selected.isDir {
					if time.Since(m.lastPreview) > 300*time.Millisecond {
						m.loadPreview(selected.path)
						m.lastPreview = time.Now()
					}
				}
			}

		case FocusPreview:
			m.viewport, cmd = m.viewport.Update(msg)
		}

		return m, cmd

	case msg.Type == tea.KeyEnter:
		return m.handleFileSelection()
	}

	return m, tea.Batch(cmds...)
}

// loadDirectory - Load directory contents
func (m Model) loadDirectory(path string) tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(path)
		if err != nil {
			return m.showError(err.Error())
		}

		items := make([]list.Item, 0, len(entries)+1)

		parentDir := filepath.Dir(path)
		if parentDir != path {
			items = append(items, FileItem{
				name:      "..",
				path:      parentDir,
				isDir:     true,
				size:      0,
				modTime:   "",
				extension: "",
			})
		}

		for _, entry := range entries {
			if entry.IsDir() {
				entryInfo, _ := entry.Info()
				items = append(items, FileItem{
					name:      entry.Name(),
					path:      filepath.Join(path, entry.Name()),
					isDir:     true,
					size:      0,
					modTime:   entryInfo.ModTime().Format("2006-01-02"),
					extension: "",
				})
			}
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				entryInfo, _ := entry.Info()
				ext := strings.ToLower(filepath.Ext(entry.Name()))
				items = append(items, FileItem{
					name:      entry.Name(),
					path:      filepath.Join(path, entry.Name()),
					isDir:     false,
					size:      entryInfo.Size(),
					modTime:   entryInfo.ModTime().Format("2006-01-02"),
					extension: ext,
				})
			}
		}

		return items
	}
}

// handleFileSelection - Handle enter on file selection
func (m Model) handleFileSelection() (tea.Model, tea.Cmd) {
	selected := m.fileList.SelectedItem()
	if selected == nil {
		return m, nil
	}

	fileItem := selected.(FileItem)
	if fileItem.isDir {
		m.currentPath = fileItem.path
		return m, m.loadDirectory(m.currentPath)
	}

	m.loadPreview(fileItem.path)
	m.panelFocus = FocusPreview
	return m, nil
}

// loadPreview - Load file preview asynchronously
func (m *Model) loadPreview(filePath string) {
	if m.previewFile == filePath && !m.previewLoading {
		return
	}

	m.previewFile = filePath
	m.previewLoading = true

	m.cacheMutex.RLock()
	if cached, ok := m.contentCache[filePath]; ok {
		m.cacheMutex.RUnlock()
		m.previewLoading = false
		m.renderPreview(cached, filePath)
		return
	}
	m.cacheMutex.RUnlock()

	go func() {
		content, err := os.ReadFile(filePath)
		if err != nil {
			m.previewData = ""
			m.previewType = PreviewNone
			m.previewLoading = false
			m.viewport.SetContent("Error: " + err.Error())
			return
		}

		contentStr := string(content)

		m.cacheMutex.Lock()
		m.contentCache[filePath] = contentStr
		m.cacheMutex.Unlock()

		m.renderPreview(contentStr, filePath)
		m.previewLoading = false
	}()
}

// renderPreview - Render preview based on file type
func (m *Model) renderPreview(content, filePath string) {
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == ".md" {
		m.previewType = PreviewMarkdown
		width := m.viewport.Width
		if width <= 0 {
			width = 60
		}
		rendered, err := markdown.Render(content, markdown.RenderOptions{
			Width:            width,
			Plain:            false,
			PreserveNewlines: false,
			TrimNewlines:     true,
			EnableTables:     true,
			EnableTaskLists:  true,
		})
		if err != nil {
			m.viewport.SetContent(content)
		} else {
			m.viewport.SetContent(rendered)
		}
		m.previewData = content
		return
	}

	textExtensions := []string{
		".txt", ".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".h", ".hpp",
		".rs", ".rb", ".php", ".html", ".css", ".json", ".xml", ".yaml", ".yml",
		".toml", ".ini", ".conf", ".log", ".sh", ".bat", ".ps1", ".sql",
		".md", ".rst", ".cfg", ".env",
	}

	isText := false
	for _, texExt := range textExtensions {
		if ext == texExt {
			isText = true
			break
		}
	}

	if isText {
		m.previewType = PreviewText
		m.viewport.SetContent(content)
		m.previewData = content
	} else {
		m.previewType = PreviewBinary
		m.viewport.SetContent("[Binary file]")
		m.previewData = ""
	}

	m.viewport.GotoTop()
}

// executeCommand - Execute command
func (m Model) executeCommand(cmd string) (tea.Model, tea.Cmd) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		m.showCommand = false
		m.commandInput.Reset()
		return m, nil
	}

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		m.showCommand = false
		m.commandInput.Reset()
		return m, nil
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	m.outputLines = append(m.outputLines, fmt.Sprintf(">>> %s", cmd))

	switch command {
	case "cd":
		return m.cmdChangeDir(args)
	case "ls", "dir":
		return m.cmdList(args)
	case "help", "?":
		m.panelFocus = FocusPreview
		return m.cmdHelp(args)
	case "clear":
		m.fileList.ResetFilter()
		m.showCommand = false
		m.commandInput.Reset()
		m.outputLines = append(m.outputLines, "Search filter cleared.")
		m.commandOutput.SetContent(strings.Join(m.outputLines, "\n"))
		return m, nil
	case "task":
		return m.cmdTask(args)
	case "refresh", "reload":
		m.showCommand = false
		m.commandInput.Reset()
		return m, m.loadDirectory(m.currentPath)
	case "preview":
		return m.cmdPreview(args)
	case "exit", "quit":
		m.quitting = true
		return m, tea.Quit
	case "cat":
		return m.cmdCat(args)
	case "echo":
		return m.cmdEcho(args)
	default:
		return m.cmdShell(cmd)
	}
}

// cmdChangeDir
func (m Model) cmdChangeDir(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err == nil {
			m.currentPath = home
		}
	} else if args[0] == ".." {
		m.currentPath = filepath.Dir(m.currentPath)
	} else if args[0] == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			m.currentPath = home
		}
	} else {
		newPath := args[0]
		if !filepath.IsAbs(newPath) {
			newPath = filepath.Join(m.currentPath, newPath)
		}
		if info, err := os.Stat(newPath); err == nil && info.IsDir() {
			m.currentPath = newPath
		} else {
			m.outputLines = append(m.outputLines, fmt.Sprintf("Error: Directory not found: %s", newPath))
		}
	}

	m.showCommand = false
	m.commandInput.Reset()
	m.outputLines = append(m.outputLines, fmt.Sprintf("Changed to: %s", m.currentPath))
	m.showCommandOutput()
	return m, m.loadDirectory(m.currentPath)
}

// cmdList
func (m Model) cmdList(args []string) (tea.Model, tea.Cmd) {
	m.showCommand = false
	m.commandInput.Reset()

	var listPath string
	if len(args) == 0 {
		listPath = m.currentPath
	} else {
		listPath = args[0]
		if !filepath.IsAbs(listPath) {
			listPath = filepath.Join(m.currentPath, listPath)
		}
	}

	entries, err := os.ReadDir(listPath)
	if err != nil {
		m.outputLines = append(m.outputLines, fmt.Sprintf("Error: %v", err))
		m.showCommandOutput()
		return m, nil
	}

	m.outputLines = append(m.outputLines, fmt.Sprintf("Contents of %s:", listPath))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		m.outputLines = append(m.outputLines, fmt.Sprintf("  %s", name))
	}

	m.showCommand = false
	m.commandInput.Reset()
	m.showCommandOutput()
	return m, nil
}

// cmdHelp
func (m Model) cmdHelp(args []string) (tea.Model, tea.Cmd) {
	helpText := `=== Command Help ===

Built-in Commands:
  cd <path>      - Change directory
  ls [path]      - List directory
  preview <file> - Preview a file
  cat <file>     - Display file
  echo <text>    - Echo text
  task <desc>    - Add task
  clear          - Clear filter
  refresh        - Refresh dir
  help           - Show help
  exit           - Exit

Keyboard:
  Ctrl+C/q - Quit
  / - Search
  ! - Command
  r - Refresh
  ↑/k, ↓/j - Navigate
  ←/h, →/l - Panels
  Tab - Toggle
  u - Parent
  [/] - Resize
`

	m.outputLines = append(m.outputLines, helpText)
	m.showCommand = false
	m.commandInput.Reset()
	m.showCommandOutput()
	return m, nil
}

// cmdTask
func (m Model) cmdTask(args []string) (tea.Model, tea.Cmd) {
	if len(args) > 0 {
		task := strings.Join(args, " ")
		m.todos.AddTodo(messages.Todo{
			ID:          fmt.Sprintf("todo-%d", len(m.todos.Todos())+1),
			Description: task,
			Status:      messages.TodoStatusPending,
		})
		m.outputLines = append(m.outputLines, fmt.Sprintf("Added task: %s", task))
	}

	m.showCommand = false
	m.commandInput.Reset()
	m.showCommandOutput()
	return m, nil
}

// cmdPreview
func (m Model) cmdPreview(args []string) (tea.Model, tea.Cmd) {
	if len(args) > 0 {
		path := args[0]
		if !filepath.IsAbs(path) {
			path = filepath.Join(m.currentPath, path)
		}

		info, err := os.Stat(path)
		if err != nil {
			m.outputLines = append(m.outputLines, fmt.Sprintf("Error: %v", err))
			m.showCommand = false
			m.commandInput.Reset()
			m.commandOutput.SetContent(strings.Join(m.outputLines, "\n"))
			m.commandOutput.GotoBottom()
			return m, nil
		}

		if !info.IsDir() {
			m.loadPreview(path)
			m.panelFocus = FocusPreview
			m.outputLines = append(m.outputLines, fmt.Sprintf("Previewing: %s", path))
		}
	}

	m.showCommand = false
	m.commandInput.Reset()
	m.commandOutput.SetContent(strings.Join(m.outputLines, "\n"))
	m.commandOutput.GotoBottom()
	return m, nil
}

// cmdCat
func (m Model) cmdCat(args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		m.outputLines = append(m.outputLines, "Error: No file specified")
		m.showCommand = false
		m.commandInput.Reset()
		m.showCommandOutput()
		return m, nil
	}

	path := args[0]
	if !filepath.IsAbs(path) {
		path = filepath.Join(m.currentPath, path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		m.outputLines = append(m.outputLines, fmt.Sprintf("Error: %v", err))
		m.showCommand = false
		m.commandInput.Reset()
		m.showCommandOutput()
		return m, nil
	}

	m.outputLines = append(m.outputLines, fmt.Sprintf("--- %s ---", path))
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		m.outputLines = append(m.outputLines, line)
	}

	m.showCommand = false
	m.commandInput.Reset()
	m.showCommandOutput()
	return m, nil
}

// cmdEcho
func (m Model) cmdEcho(args []string) (tea.Model, tea.Cmd) {
	text := strings.Join(args, " ")
	m.outputLines = append(m.outputLines, text)
	m.showCommand = false
	m.commandInput.Reset()
	m.showCommandOutput()
	return m, nil
}

// cmdShell
func (m Model) cmdShell(cmd string) (tea.Model, tea.Cmd) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = "cmd.exe"
		} else {
			shell = "/bin/sh"
		}
	}

	var execCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		execCmd = exec.Command("cmd", "/C", cmd)
	} else {
		execCmd = exec.Command(shell, "-c", cmd)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	execCmd = exec.CommandContext(ctx, execCmd.Path, execCmd.Args[1:]...)

	out, err := execCmd.CombinedOutput()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			m.outputLines = append(m.outputLines, "Note: Command timed out after 5 seconds")
			m.outputLines = append(m.outputLines, "Output captured up to timeout...")
		} else {
			m.outputLines = append(m.outputLines, fmt.Sprintf("Error: %v", err))
		}
	}

	var outputStr string
	if runtime.GOOS == "windows" {
		outputStr = decodeGBK(out)
	} else {
		outputStr = string(out)
	}

	if outputStr != "" {
		lines := strings.Split(outputStr, "\n")
		maxLines := 100
		if len(lines) > maxLines {
			lines = lines[:maxLines]
			lines = append(lines, fmt.Sprintf("... (truncated, showing first %d lines)", maxLines))
		}
		for _, line := range lines {
			if line != "" {
				m.outputLines = append(m.outputLines, line)
			}
		}
	}

	m.showCommand = false
	m.commandInput.Reset()
	m.commandOutput.SetContent(strings.Join(m.outputLines, "\n"))
	m.commandOutput.GotoBottom()
	m.panelFocus = FocusPreview
	return m, nil
}

// decodeGBK
func decodeGBK(data []byte) string {
	reader := transform.NewReader(strings.NewReader(string(data)), simplifiedchinese.GBK.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return string(data)
	}
	return string(decoded)
}

// showError
func (m Model) showError(err string) tea.Cmd {
	m.outputLines = append(m.outputLines, fmt.Sprintf("Error: %s", err))
	m.commandOutput.SetContent(strings.Join(m.outputLines, "\n"))
	m.commandOutput.GotoBottom()
	return nil
}

// showCommandOutput
func (m *Model) showCommandOutput() {
	m.commandOutput.SetContent(strings.Join(m.outputLines, "\n"))
	m.commandOutput.GotoBottom()
}

// showHelp
func (m Model) showHelp() Model {
	return m
}

// View - Using image-viewer's layout technique
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Calculate layout if needed
	if m.layoutManager == nil {
		m.recalculateLayout()
		m.updateComponentSizes()
	}

	// Fixed header and footer lines
	headerLines := 1
	footerLines := 1
	availableHeight := m.height - headerLines - footerLines

	var b strings.Builder

	// === Header (fixed 1 line) ===
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color("#313244")).
		Foreground(lipgloss.Color("#cdd6f4"))
	
	headerText := fmt.Sprintf("📁 File Browser - %s", m.currentPath)
	b.WriteString(headerStyle.Render(headerText))
	b.WriteString("\n")

	// === Main Area ===
	// Get panel contents
	fileListPanel := m.renderFileListPanel(m.fileListRect.Width, availableHeight)
	previewPanel := m.renderPreviewPanel(m.previewRect.Width, availableHeight)
	
	// Join horizontally
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, fileListPanel, previewPanel)
	
	// Split into lines and limit to availableHeight
	mainLines := strings.Split(mainContent, "\n")
	
	// Remove trailing empty line if exists (JoinHorizontal might add one)
	if len(mainLines) > 0 && mainLines[len(mainLines)-1] == "" {
		mainLines = mainLines[:len(mainLines)-1]
	}
	
	displayedLines := 0
	for _, line := range mainLines {
		if displayedLines >= availableHeight {
			break
		}
		b.WriteString(line)
		b.WriteString("\n")
		displayedLines++
	}

	// Fill remaining space with empty lines to push footer to bottom
	remainingPadding := availableHeight - displayedLines
	for i := 0; i < remainingPadding; i++ {
		b.WriteString("\n")
	}

	// === Footer (fixed at bottom) ===
	footerText := "[/] Search  [!] Command  [Tab] Focus  [q] Quit  [?] Help"
	footerStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("#313244")).
		Foreground(lipgloss.Color("#a6adc8"))
	b.WriteString(footerStyle.Render(footerText))

	return b.String()
}

// renderFileListPanel - Render file list panel
func (m Model) renderFileListPanel(width, height int) string {
	focused := m.panelFocus == FocusFileList

	// Use consistent border style to prevent visual jumping during focus changes
	var panelStyle lipgloss.Style
	if focused {
		panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#cba6f7")).
			Padding(1, 1).
			Width(width)
	} else {
		panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#45475a")).
			Padding(1, 1).
			Width(width)
	}

	return panelStyle.Render(m.fileList.View())
}

// renderPreviewPanel - Render preview/output panel (right side)
func (m Model) renderPreviewPanel(width, height int) string {
	panelFocused := m.panelFocus == FocusPreview

	// DEBUG: Show width calculation
	debugWidth := fmt.Sprintf("DEBUG: panelWidth=%d, previewRect.Width=%d, m.width=%d\n\n", 
		width, m.previewRect.Width, m.width)

	// Determine what to show
	var content string
	var title string
	if m.previewFile != "" && m.previewData != "" {
		content = m.viewport.View()
		title = fmt.Sprintf(" Preview: %s ", m.previewFile)
	} else {
		content = debugWidth + m.commandOutput.View()
		title = " Command Output "
	}

	// Add focus indicator
	if panelFocused {
		title = "*" + title
	}

	// Create title bar
	titleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1e1e2e")).
		Foreground(lipgloss.Color("#cba6f7")).
		Padding(0, 1)
	
	titleBar := titleStyle.Render(title)

	// Create panel style - NO fixed height, let it grow with content
	var panelStyle lipgloss.Style
	if panelFocused {
		panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#cba6f7")).
			Width(width)
	} else {
		panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#45475a")).
			Width(width)
	}

	// Combine title and content
	fullContent := titleBar + "\n" + content
	return panelStyle.Render(fullContent)
}

// renderBottomPanel - Legacy, no longer used
func (m Model) renderBottomPanel(width, height int) string {
	return m.renderPreviewPanel(width, height)
}

// renderPreviewContent - Legacy function, no longer used
func (m Model) renderPreviewContent(width, height int) string {
	return m.viewport.View()
}

// renderOutputContent - Legacy function, no longer used
func (m Model) renderOutputContent(width, height int) string {
	return m.commandOutput.View()
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
