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
	"github.com/wwsheng009/taproot/ui/markdown"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Application state
type AppState int

const (
	StateFileBrowser AppState = iota
	StatePreviewPanel
)

// PanelFocus - Which panel has focus
type PanelFocus int

const (
	FocusFileList PanelFocus = iota
	FocusPreview
	FocusCommandOutput
)

// Model - Complete application model
type Model struct {
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
	fileListWidth int
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
		currentPath:      wd,
		fileList:         fileList,
		fileListWidth:    40,
		keyMap:           DefaultKeyMap,
		panelFocus:       FocusFileList,
		lspList:          lspList,
		mcpList:          mcpList,
		todos:            todos,
		commandInput:     commandInput,
		searchInput:      searchInput,
		viewport:         vp,
		commandOutput:    cmdOutput,
		previewFile:      "",
		previewData:      "",
		previewType:      PreviewNone,
		previewLoading:   false,
		outputLines:      []string{"Ready. Type ! to enter command mode."},
		contentCache:     make(map[string]string),
		commandHistory:   make([]string, 0),
		historyIndex:     -1,
		width:            80,
		height:           24,
		quitting:         false,
		previewThrottler: nil,
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

		// Update file list size
		listHeight := m.height - 10
		m.fileList.SetSize(m.fileListWidth-2, listHeight)

		// Update viewport size
		viewportWidth := m.width - m.fileListWidth - 4
		viewportHeight := m.height - 10
		m.viewport.Width = viewportWidth
		m.viewport.Height = viewportHeight
		m.commandOutput.Width = viewportWidth
		m.commandOutput.Height = viewportHeight

	case []list.Item:
		// File list loaded
		m.fileList.SetItems(msg)
	}

	return m, tea.Batch(cmds...)
}

// updateCommandMode - Handle command mode input
func (m Model) updateCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.commandInput, cmd = m.commandInput.Update(msg)

	switch msg.Type {
	case tea.KeyEnter:
		command := m.commandInput.Value()
		if command != "" {
			// Add to history
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
		// Navigate command history
		if m.historyIndex > 0 {
			m.historyIndex--
			m.commandInput.SetValue(m.commandHistory[m.historyIndex])
		}

	case tea.KeyDown:
		// Navigate command history
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

	// Real-time filtering
	query := m.searchInput.Value()
	if query != "" {
		m.fileList.FilterInput.SetValue(query)
	} else {
		m.fileList.ResetFilter()
	}

	switch msg.Type {
	case tea.KeyEnter:
		// Accept search
		m.showSearch = false
		m.searchInput.Blur()
		return m, nil

	case tea.KeyEsc:
		// Cancel search
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
		// Tab: Toggle between file list, preview, and output
		m.panelFocus = (m.panelFocus + 1) % 3
		return m, nil

	case key.Matches(msg, m.keyMap.ResizePanel):
		// Resize file list using [ and ]
		if msg.Type == tea.KeyRunes && len(msg.Runes) > 0 {
			rune := msg.Runes[0]
			if rune == '[' {
				// Decrease width
				if m.fileListWidth > 30 {
					m.fileListWidth -= 5
				}
			} else if rune == ']' {
				// Increase width
				if m.fileListWidth < m.width-30 {
					m.fileListWidth += 5
				}
			}

			// Update sizes
			listHeight := m.height - 10
			m.fileList.SetSize(m.fileListWidth-2, listHeight)
			viewportWidth := m.width - m.fileListWidth - 4
			viewportHeight := m.height - 10
			m.viewport.Width = viewportWidth
			m.viewport.Height = viewportHeight
			m.commandOutput.Width = viewportWidth
			m.commandOutput.Height = viewportHeight
		}
		return m, nil

	case key.Matches(msg, m.keyMap.Left):
		// Left: switch to file list
		if m.panelFocus != FocusFileList {
			m.panelFocus = FocusFileList
			return m, nil
		}
		// Let list handle it if focused
		var listCmd tea.Cmd
		m.fileList, listCmd = m.fileList.Update(msg)
		return m, listCmd

	case key.Matches(msg, m.keyMap.Right):
		// Right: switch to preview/output
		if m.panelFocus == FocusFileList {
			m.panelFocus = FocusPreview
			return m, nil
		} else if m.panelFocus == FocusPreview {
			m.panelFocus = FocusCommandOutput
			return m, nil
		}
		// Let viewport handle it if focused
		var viewCmd tea.Cmd
		if m.panelFocus == FocusPreview {
			m.viewport, viewCmd = m.viewport.Update(msg)
		} else if m.panelFocus == FocusCommandOutput {
			m.commandOutput, viewCmd = m.commandOutput.Update(msg)
		}
		return m, viewCmd

	case key.Matches(msg, m.keyMap.Up), key.Matches(msg, m.keyMap.Down),
	     key.Matches(msg, m.keyMap.PageUp), key.Matches(msg, m.keyMap.PageDown),
	     key.Matches(msg, m.keyMap.Home), key.Matches(msg, m.keyMap.End):
		// Handle based on focus
		var cmd tea.Cmd

		switch m.panelFocus {
		case FocusFileList:
			m.fileList, cmd = m.fileList.Update(msg)

			// Update preview when selection changes (with throttling)
			if m.fileList.SelectedItem() != nil {
				selected := m.fileList.SelectedItem().(FileItem)
				if !selected.isDir {
					// Throttle preview updates to avoid lag
					if time.Since(m.lastPreview) > 300*time.Millisecond {
						m.loadPreview(selected.path)
						m.lastPreview = time.Now()
					}
				}
			}

		case FocusPreview:
			m.viewport, cmd = m.viewport.Update(msg)

		case FocusCommandOutput:
			m.commandOutput, cmd = m.commandOutput.Update(msg)
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

		// Add parent directory entry if not at root
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

		// Add directory entries first
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

		// Add file entries
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

	// File selected - show preview
	m.loadPreview(fileItem.path)
	m.panelFocus = FocusPreview
	return m, nil
}

// loadPreview - Load file preview asynchronously
func (m *Model) loadPreview(filePath string) {
	// Skip if already previewing this file
	if m.previewFile == filePath && !m.previewLoading {
		return
	}

	m.previewFile = filePath
	m.previewLoading = true

	// Check cache first
	m.cacheMutex.RLock()
	if cached, ok := m.contentCache[filePath]; ok {
		m.cacheMutex.RUnlock()
		m.previewLoading = false
		m.renderPreview(cached, filePath)
		return
	}
	m.cacheMutex.RUnlock()

	// Read file content asynchronously
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

		// Cache the content
		m.cacheMutex.Lock()
		m.contentCache[filePath] = contentStr
		m.cacheMutex.Unlock()

		// Determine file type and render
		m.renderPreview(contentStr, filePath)
		m.previewLoading = false
	}()
}

// renderPreview - Render preview based on file type
func (m *Model) renderPreview(content, filePath string) {
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check for markdown
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

	// Check for text files based on extension
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
		// Binary file
		m.previewType = PreviewBinary
		m.viewport.SetContent("[Binary file]")
		m.previewData = ""
	}

	// Reset viewport position
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

	// Log command to output
	m.outputLines = append(m.outputLines, fmt.Sprintf(">>> %s", cmd))

	switch command {
	case "cd":
		return m.cmdChangeDir(args)

	case "ls", "dir":
		return m.cmdList(args)

	case "help", "?":
		m.panelFocus = FocusCommandOutput
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

// cmdChangeDir - Change directory command
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

// cmdList - List directory command
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

// cmdHelp - Show help command
func (m Model) cmdHelp(args []string) (tea.Model, tea.Cmd) {
	helpText := `=== Command Help ===

Built-in Commands:
  cd <path>      - Change directory (supports .. and ~)
  ls [path]      - List directory contents
  preview <file> - Preview a file in the right panel
  cat <file>     - Display file contents
  echo <text>    - Echo text to output
  task <desc>    - Add a new task
  clear          - Clear search filter
  refresh        - Refresh current directory
  help           - Show this help message
  exit           - Exit the program

Shell Commands:
  Any command will be executed in the shell
  Examples:
    !ls -la
    !git status
    !pwd

Keyboard Shortcuts:
  Ctrl+C / q     - Quit
  /              - Enter search mode
  !              - Enter command mode
  r              - Refresh directory
  ↑/k, ↓/j       - Navigate
  ←/h, →/l       - Switch between left/right panels
  Tab            - Toggle between Files/Preview
  u              - Go to parent directory
  [ and ]        - Resize file browser panel
`

	m.outputLines = append(m.outputLines, helpText)
	m.showCommand = false
	m.commandInput.Reset()
	m.showCommandOutput()
	return m, nil
}

// cmdTask - Add todo task command
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

// cmdPreview - Preview file command
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

// cmdCat - Display file contents
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

// cmdEcho - Echo text to output
func (m Model) cmdEcho(args []string) (tea.Model, tea.Cmd) {
	text := strings.Join(args, " ")
	m.outputLines = append(m.outputLines, text)
	m.showCommand = false
	m.commandInput.Reset()
	m.showCommandOutput()
	return m, nil
}

// cmdShell - Execute shell command
func (m Model) cmdShell(cmd string) (tea.Model, tea.Cmd) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = "cmd.exe"
		} else {
			shell = "/bin/sh"
		}
	}

	// Create command with timeout
	var execCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		execCmd = exec.Command("cmd", "/C", cmd)
	} else {
		execCmd = exec.Command(shell, "-c", cmd)
	}

	// Add timeout to prevent blocking on continuous commands like ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	execCmd = exec.CommandContext(ctx, execCmd.Path, execCmd.Args[1:]...)

	// Start command and capture output
	out, err := execCmd.CombinedOutput()

	if err != nil {
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			m.outputLines = append(m.outputLines, "Note: Command timed out after 5 seconds")
			m.outputLines = append(m.outputLines, "Output captured up to timeout...")
		} else {
			m.outputLines = append(m.outputLines, fmt.Sprintf("Error: %v", err))
		}
	}

	// Convert encoding for Windows (GBK -> UTF-8)
	var outputStr string
	if runtime.GOOS == "windows" {
		// Try to decode GBK
		outputStr = decodeGBK(out)
	} else {
		outputStr = string(out)
	}
	
	if outputStr != "" {
		lines := strings.Split(outputStr, "\n")
		// Limit output to prevent memory issues
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
	m.panelFocus = FocusCommandOutput
	return m, nil
}

// decodeGBK - Convert GBK encoded bytes to UTF-8 string
func decodeGBK(data []byte) string {
	// Try GBK decoding
	reader := transform.NewReader(strings.NewReader(string(data)), simplifiedchinese.GBK.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		// If GBK decoding fails, fall back to UTF-8
		return string(data)
	}
	return string(decoded)
}

// showError - Show error message
func (m Model) showError(err string) tea.Cmd {
	m.outputLines = append(m.outputLines, fmt.Sprintf("Error: %s", err))
	m.commandOutput.SetContent(strings.Join(m.outputLines, "\n"))
	m.commandOutput.GotoBottom()
	return nil
}

// showCommandOutput - Set output content, scroll to bottom, and switch to output panel
func (m *Model) showCommandOutput() {
	m.commandOutput.SetContent(strings.Join(m.outputLines, "\n"))
	m.commandOutput.GotoBottom()
	// Don't auto-switch to output for all commands
	// Shell commands will set panelFocus directly
}

// showHelp - Show help
func (m Model) showHelp() Model {
	return m
}

// View - Render application
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader(m.width))
	b.WriteString("\n")

	// Main content area
	mainHeight := m.height - 8
	if mainHeight < 10 {
		mainHeight = 10
	}

	leftContent := m.renderFileBrowser(m.fileListWidth, mainHeight)
	// Calculate right width: total width minus left panel width
	// Ensure rightWidth includes space for right border
	rightWidth := m.width - m.fileListWidth
	if rightWidth < 20 {
		rightWidth = 20
	}
	rightContent := m.renderRightPanel(rightWidth, mainHeight)

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Left,
		leftContent,
		rightContent,
	))

	// Command/Search/Footer
	b.WriteString("\n")
	if m.showCommand {
		b.WriteString(m.renderCommandPanel(m.width))
	} else if m.showSearch {
		b.WriteString(m.renderSearchPanel(m.width))
	} else {
		b.WriteString(m.renderFooter(m.width))
	}

	return b.String()
}

// renderHeader - Render header
func (m Model) renderHeader(width int) string {
	return styleHeader.Render(fmt.Sprintf(
		" 📁 Taproot File Browser - %s ",
		m.currentPath,
	))
}

// renderFileBrowser - Render file browser with focus indicator
func (m Model) renderFileBrowser(width, height int) string {
	focused := m.panelFocus == FocusFileList
	
	var panelStyle lipgloss.Style
	if focused {
		panelStyle = stylePanelFocused
	} else {
		panelStyle = stylePanel
	}

	m.fileList.Title = "Files"
	return panelStyle.
		Height(height).
		Width(width).
		Render(m.fileList.View())
}

// renderRightPanel - Render right panel (Preview or Output)
func (m Model) renderRightPanel(width, height int) string {
	// Check if width is sufficient for border and content
	minWidth := 10 // minimum width for border + padding
	if width < minWidth {
		return ""
	}

	panelFocused := m.panelFocus == FocusPreview || m.panelFocus == FocusCommandOutput
	
	// Render content based on current panel
	var content string
	contentWidth := width - 4 // border (2) + padding (2)
	contentHeight := height - 3 // tabs (1) + border/padding (2)
	
	if contentWidth < 1 {
		contentWidth = 1
	}
	if contentHeight < 1 {
		contentHeight = 1
	}

	switch m.panelFocus {
	case FocusPreview:
		content = m.renderPreviewContent(contentWidth, contentHeight)
	case FocusCommandOutput:
		content = m.renderOutputContent(contentWidth, contentHeight)
	default:
		content = m.renderPreviewContent(contentWidth, contentHeight)
	}
	
	// Build panel with content and border
	var panelStyle lipgloss.Style
	if panelFocused {
		panelStyle = stylePanelFocused
	} else {
		panelStyle = stylePanel
	}
	
	panel := panelStyle.
		Height(height-1).
		Width(width).
		Render(content)
	
	// Simple tabs at top (outside the panel border)
	previewTab := " Preview "
	outputTab := " Output "
	
	if m.panelFocus == FocusPreview {
		previewTab = "* Preview "
	} else if m.panelFocus == FocusCommandOutput {
		outputTab = "* Output "
	}
	
	tabsStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1e1e2e")).
		Foreground(lipgloss.Color("#cba6f7")).
		Width(width)
	
	tabsText := previewTab + outputTab
	tabs := tabsStyle.Render(tabsText)
	
	// Tabs on top, panel below
	return tabs + "\n" + panel
}

// renderPreviewContent - Render preview content without header
func (m Model) renderPreviewContent(width, height int) string {
	content := m.viewport.View()
	
	if content == "" {
		content = "Select a file to preview"
	}
	
	// Ensure content has enough lines to render full height
	lines := strings.Split(content, "\n")
	missingLines := height - len(lines)
	if missingLines > 0 {
		for i := 0; i < missingLines; i++ {
			content += "\n"
		}
	}
	
	return content
}

// renderOutputContent - Render command output content without header
func (m Model) renderOutputContent(width, height int) string {
	content := m.commandOutput.View()
	
	if content == "" {
		content = `No output yet

Type ! to enter command mode

Available commands:
  ls         - List directory
  cd <path>  - Change directory
  cat <file> - Display file contents
  echo <text>- Echo text
  any cmd    - Execute shell command

Note: Long-running commands timeout after 5 seconds
Examples:
  ping -n 4 baidu.com    (Windows)
  ping -c 4 baidu.com    (Linux)`
	}
	
	// Ensure content has enough lines to render full height
	lines := strings.Split(content, "\n")
	missingLines := height - len(lines)
	if missingLines > 0 {
		for i := 0; i < missingLines; i++ {
			content += "\n"
		}
	}
	
	return content
}

// renderPreviewPanel - Render file preview panel with focus indicator
func (m Model) renderPreviewPanel(width, height int) string {
	focused := m.panelFocus == FocusPreview
	
	var panelStyle, headerStyle lipgloss.Style
	if focused {
		panelStyle = stylePanelFocused
		headerStyle = styleHeaderFocused
	} else {
		panelStyle = stylePanel
		headerStyle = styleHeader
	}

	var title string
	if m.previewLoading {
		title = "⏳ Loading..."
	} else {
		switch m.previewType {
		case PreviewMarkdown:
			title = "📖 Preview"
		case PreviewText:
			title = "📄 Preview"
		case PreviewBinary:
			title = "🔒 Preview"
		default:
			title = "👁️ Preview"
		}
	}

	header := headerStyle.Render(title)
	content := m.viewport.View()
	
	if content == "" {
		content = "Select a file to preview"
	}

	return header + "\n" + panelStyle.
		Height(height-2).
		Width(width).
		Render(content)
}

// renderCommandOutputPanel - Render command output panel with focus indicator
func (m Model) renderCommandOutputPanel(width, height int) string {
	focused := m.panelFocus == FocusCommandOutput
	
	var panelStyle, headerStyle lipgloss.Style
	if focused {
		panelStyle = stylePanelFocused
		headerStyle = styleHeaderFocused
	} else {
		panelStyle = stylePanel
		headerStyle = styleHeader
	}

	header := headerStyle.Render("⚡ Output")
	content := m.commandOutput.View()
	
	if content == "" {
		content = "No output yet"
	}

	return header + "\n" + panelStyle.
		Height(height-2).
		Width(width).
		Render(content)
}

// renderCommandPanel - Render command input panel
func (m Model) renderCommandPanel(width int) string {
	return styleCommand.Render(m.commandInput.View())
}

// renderCommandInput - Render command input line for Output panel
func (m Model) renderCommandInput(width int) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#313244")).
		Foreground(lipgloss.Color("#cdd6f4")).
		Padding(0, 1).
		Width(width).
		Render(m.commandInput.View())
}

// renderSearchPanel - Render search input panel
func (m Model) renderSearchPanel(width int) string {
	return styleCommand.Render(m.searchInput.View())
}

// renderFooter - Render footer
func (m Model) renderFooter(width int) string {
	help := fmt.Sprintf(
		"[%s] Quit  [%s] Search  [%s] Cmd  [%s] Refresh  [%s/%s] Panel  [%s] Parent  [%s] Resize",
		m.keyMap.Quit.Help().Key,
		m.keyMap.Search.Help().Key,
		m.keyMap.Command.Help().Key,
		m.keyMap.Refresh.Help().Key,
		m.keyMap.Left.Help().Key,
		m.keyMap.Right.Help().Key,
		m.keyMap.GoToParent.Help().Key,
		"[",
	)
	return styleFooter.Render(help)
}

// Styles
var (
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("#1e1e2e")).
			Foreground(lipgloss.Color("#cba6f7"))

	styleHeaderFocused = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("#cba6f7")).
			Foreground(lipgloss.Color("#1e1e2e"))

	stylePanel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#45475a")).
		Padding(1, 1)

	stylePanelFocused = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#cba6f7")).
		Padding(1, 1)

	styleCommand = lipgloss.NewStyle().
		Background(lipgloss.Color("#313244")).
		Foreground(lipgloss.Color("#cdd6f4")).
		Padding(0, 1)

	styleFooter = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("#1e1e2e")).
		Foreground(lipgloss.Color("#6c7086"))
)

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
