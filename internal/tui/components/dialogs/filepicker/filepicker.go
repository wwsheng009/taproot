package filepicker

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs"
	"github.com/wwsheng009/taproot/internal/tui/util"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

const (
	ID dialogs.DialogID = "filepicker"
)

type Callback func(path string) tea.Cmd

type FilePicker struct {
	styles      *styles.Styles
	id          dialogs.DialogID
	currentDir  string
	files       []os.DirEntry
	cursor      int
	callback    Callback
	showHidden  bool
	width       int
	height      int
	err         error
	
	// Scrolling state
	scrollOffset int
}

func New(startPath string, callback Callback) *FilePicker {
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		absPath = "."
	}
	s := styles.DefaultStyles()
	
	fp := &FilePicker{
		styles:     &s,
		id:         ID,
		currentDir: absPath,
		callback:   callback,
		cursor:     0,
	}
	
	fp.readDir()
	return fp
}

func (m *FilePicker) Init() tea.Cmd {
	return nil
}

func (m *FilePicker) readDir() {
	entries, err := os.ReadDir(m.currentDir)
	if err != nil {
		m.err = err
		return
	}
	
	m.files = []os.DirEntry{}
	
	for _, e := range entries {
		if !m.showHidden && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		m.files = append(m.files, e)
	}
	
	sort.Slice(m.files, func(i, j int) bool {
		// Dirs first
		if m.files[i].IsDir() && !m.files[j].IsDir() {
			return true
		}
		if !m.files[i].IsDir() && m.files[j].IsDir() {
			return false
		}
		return m.files[i].Name() < m.files[j].Name()
	})
	
	m.err = nil
	m.cursor = 0
	m.scrollOffset = 0
}

func (m *FilePicker) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return dialogs.CloseDialogMsg{} }
		case "up", "k":
			if m.cursor > -1 {
				m.cursor--
			}
			m.updateScroll()
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
			m.updateScroll()
		case "enter":
			return m.handleSelect()
		case "left", "backspace":
			// Go up a directory
			parent := filepath.Dir(m.currentDir)
			if parent != m.currentDir {
				m.currentDir = parent
				m.readDir()
			}
		case ".":
			m.showHidden = !m.showHidden
			m.readDir()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *FilePicker) updateScroll() {
	maxItems := 10 // Visible items excluding parent
	
	// Adjust scrollOffset to keep cursor in view
	if m.cursor == -1 {
		// Parent directory is always at top, but technically outside the list index
		// We can just set scroll to 0
		m.scrollOffset = 0
		return
	}
	
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+maxItems {
		m.scrollOffset = m.cursor - maxItems + 1
	}
}

func (m *FilePicker) handleSelect() (util.Model, tea.Cmd) {
	if m.cursor == -1 {
		// Parent directory
		parent := filepath.Dir(m.currentDir)
		if parent != m.currentDir {
			m.currentDir = parent
			m.readDir()
		}
		return m, nil
	}
	
	if len(m.files) == 0 {
		return m, nil
	}
	
	selected := m.files[m.cursor]
	path := filepath.Join(m.currentDir, selected.Name())
	
	if selected.IsDir() {
		m.currentDir = path
		m.readDir()
		return m, nil
	}
	
	// It's a file
	if m.callback != nil {
		cmd := m.callback(path)
		return m, tea.Batch(
			cmd,
			func() tea.Msg { return dialogs.CloseDialogMsg{} },
		)
	}
	
	return m, func() tea.Msg { return dialogs.CloseDialogMsg{} }
}

func (m *FilePicker) View() string {
	s := m.styles
	
	// Simple Box
	width := 60
	height := 20
	
	boxStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Border).
		Padding(1, 2)
		
	var sb strings.Builder
	
	// Title
	sb.WriteString(styles.ApplyForegroundGrad(s, fmt.Sprintf("📂 %s", m.currentDir), s.Primary, s.Secondary))
	sb.WriteString("\n\n")
	
	if m.err != nil {
		sb.WriteString(lipgloss.NewStyle().Foreground(s.Error).Render(fmt.Sprintf("Error: %v", m.err)))
		return boxStyle.Render(sb.String())
	}
	
	// Parent dir option
	cursor := "  "
	style := lipgloss.NewStyle().Foreground(s.FgBase)
	if m.cursor == -1 {
		cursor = "> "
		style = style.Foreground(s.Primary).Bold(true)
	}
	sb.WriteString(style.Render(cursor + ".. (Parent Directory)"))
	sb.WriteString("\n")
	
	// Files
	maxItems := 10
	
	for i := range maxItems {
		idx := m.scrollOffset + i
		if idx >= len(m.files) {
			break
		}
		
		f := m.files[idx]
		isDir := f.IsDir()
		
		cursor := "  "
		style := lipgloss.NewStyle().Foreground(s.FgBase)
		
		if idx == m.cursor {
			cursor = "> "
			style = style.Foreground(s.Primary).Bold(true)
		} else if isDir {
			style = style.Foreground(s.Secondary)
		}
		
		icon := "📄"
		if isDir {
			icon = "📁"
		}
		
		name := f.Name()
		if isDir {
			name += "/"
		}
		
		// Truncate name if too long
		if len(name) > 45 {
			name = name[:42] + "..."
		}
		
		sb.WriteString(style.Render(fmt.Sprintf("%s%s %s", cursor, icon, name)))
		sb.WriteString("\n")
	}
	
	// Simple footer hint
	sb.WriteString(lipgloss.NewStyle().Foreground(s.FgMuted).Render("\n[esc] cancel  [.] toggle hidden"))
	
	return boxStyle.Render(sb.String())
}

func (m *FilePicker) Position() (int, int) {
	return 0, 0
}

func (m *FilePicker) ID() dialogs.DialogID {
	return m.id
}
