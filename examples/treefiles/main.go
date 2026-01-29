package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/ui/components/treefiles"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
)

type model struct {
	tree      *treefiles.FileTree
	cursor    int
	visible   []*treefiles.FileNode
	scroll    int
	sortBy    treefiles.SortBy
	sortOrder treefiles.SortOrder
	sortNames []string
	width     int
	height    int
	ready     bool
	statusMsg string
}

func initialModel(path string) *model {
	m := &model{
		cursor:    0,
		scroll:    0,
		sortBy:    treefiles.SortByName,
		sortOrder: treefiles.SortAscending,
		sortNames: []string{"Name", "Size", "Time", "Type"},
		width:     defaultWidth,
		height:    defaultHeight,
	}

	tree, err := treefiles.NewFileTree(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file tree: %v\n", err)
		os.Exit(1)
	}
	m.tree = tree
	m.visible = tree.Flatten()

	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}

		case "down", "j":
			if m.cursor < len(m.visible)-1 {
				m.cursor++
				visibleHeight := m.height - 4
				if m.cursor >= m.scroll+visibleHeight {
					m.scroll = m.cursor - visibleHeight + 1
				}
			}

		case "pgup":
			visibleHeight := m.height - 4
			m.cursor -= visibleHeight / 2
			if m.cursor < 0 {
				m.cursor = 0
			}
			if m.cursor < m.scroll {
				m.scroll = m.cursor
			}

		case "pgdown":
			visibleHeight := m.height - 4
			m.cursor += visibleHeight / 2
			if m.cursor >= len(m.visible) {
				m.cursor = len(m.visible) - 1
			}
			visibleHeight = m.height - 4
			if m.cursor >= m.scroll+visibleHeight {
				m.scroll = m.cursor - visibleHeight + 1
			}

		case "home":
			m.cursor = 0
			m.scroll = 0

		case "end":
			m.cursor = len(m.visible) - 1
			visibleHeight := m.height - 4
			m.scroll = m.cursor - visibleHeight + 1
			if m.scroll < 0 {
				m.scroll = 0
			}

		case "enter":
			if len(m.visible) > 0 {
				selected := m.visible[m.cursor]
				if selected.IsDir() {
					m.tree.ToggleNode(selected)
					m.visible = m.tree.Flatten()
					m.statusMsg = fmt.Sprintf("Toggled: %s", selected.Name())
				} else {
					m.statusMsg = fmt.Sprintf("Selected: %s (%d bytes)", selected.Name(), selected.Size())
				}
			}

		case "+", "=":
			m.tree.ExpandAll()
			m.visible = m.tree.Flatten()
			m.statusMsg = "Expanded all directories"

		case "-", "_":
			m.tree.CollapseAll()
			m.visible = m.tree.Flatten()
			m.statusMsg = "Collapsed all directories"

		case "r":
			err := m.tree.Rescan()
			if err != nil {
				m.statusMsg = fmt.Sprintf("Rescan failed: %v", err)
			} else {
				m.visible = m.tree.Flatten()
				m.statusMsg = "Rescanned file tree"
			}

		case "h":
			m.tree.ToggleHidden()
			m.visible = m.tree.Flatten()
			if m.tree.IncludesHidden() {
				m.statusMsg = "Showing hidden files"
			} else {
				m.statusMsg = "Hiding hidden files"
			}

		case "s":
			m.cycleSort()
			m.tree.SetSort(m.sortBy, m.sortOrder)
			m.visible = m.tree.Flatten()
			sortOrder := "Asc"
			if m.sortOrder == treefiles.SortDescending {
				sortOrder = "Desc"
			}
			m.statusMsg = fmt.Sprintf("Sorted by %s (%s)", m.sortNames[m.sortBy], sortOrder)

		case "o":
			if m.sortOrder == treefiles.SortDescending {
				m.sortOrder = treefiles.SortAscending
			} else {
				m.sortOrder = treefiles.SortDescending
			}
			m.tree.SetSort(m.sortBy, m.sortOrder)
			m.visible = m.tree.Flatten()
			sortOrder := "Asc"
			if m.sortOrder == treefiles.SortDescending {
				sortOrder = "Desc"
			}
			m.statusMsg = fmt.Sprintf("Sort order: %s", sortOrder)

		case "esc":
			m.tree.CollapseAll()
			m.visible = m.tree.Flatten()
			m.statusMsg = "Reset view"
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	}

	return m, nil
}

func (m *model) cycleSort() {
	m.sortBy = (m.sortBy + 1) % 4
	m.sortOrder = treefiles.SortAscending
}

func (m *model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var sb strings.Builder

	// Header
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7DCFFF"))
	sb.WriteString(titleStyle.Render("🌲 File Tree Explorer"))
	sb.WriteString("\n\n")

	// Tree content
	visibleHeight := m.height - 8
	if len(m.visible) > 0 {
		start := m.scroll
		end := start + visibleHeight
		if end > len(m.visible) {
			end = len(m.visible)
		}

		selectedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#585B70"))

		dirStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7DCFFF")).Bold(true)
		fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#CDD6F4"))
		hiddenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6C7086")).Italic(true)

		for i := start; i < end; i++ {
			node := m.visible[i]
			line := ""
			isSelected := (i == m.cursor)

			prefix := ""
			if i > 0 {
				isLast := (i == len(m.visible)-1) || (i < len(m.visible)-1 && m.visible[i+1].Depth() <= node.Depth())
				prefix = treefiles.GetTreePrefix(node, isLast)
				line += prefix
			}

			name := node.Name()
			var style lipgloss.Style
			if isSelected {
				style = selectedStyle
			} else if node.IsDir() {
				style = dirStyle
			} else {
				style = fileStyle
			}

			if strings.HasPrefix(name, ".") {
				line += hiddenStyle.Render(name)
			} else {
				line += style.Render(name)
			}

			// Add file type icon on the right, aligned
			icon := treefiles.GetTreeIcon(node)
			iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6C7086"))
			iconText := iconStyle.Render(icon)

			// Calculate padding to align icon to the right
			// Current line has prefix + name, we need to add padding so icon appears at right
			availableWidth := m.width
			currentLineWidth := lipgloss.Width(line)
			iconWidth := lipgloss.Width(iconText)

			// Calculate padding to align icon to right edge
			padding := availableWidth - currentLineWidth - iconWidth
			if padding < 2 {
				padding = 2 // Minimum padding
			}

			for j := 0; j < padding; j++ {
				line += " "
			}

			line += iconText

			sb.WriteString(line + "\n")
		}
	} else {
		sb.WriteString("(empty)\n")
	}

	sb.WriteString("\n")

	// Status bar
	if m.statusMsg != "" {
		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1"))
		sb.WriteString(statusStyle.Render(m.statusMsg + "\n"))
	}

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6C7086"))
	help := []string{
		"↑/k: up    ↓/j: down    Enter: toggle   +/-: expand/collapse all",
		"s: sort    h: hidden    r: rescan      o: toggle order   q/ctrl+c: quit",
	}
	for _, line := range help {
		sb.WriteString(helpStyle.Render(line) + "\n")
	}

	// Selected file info
	if len(m.visible) > 0 && m.cursor >= 0 && m.cursor < len(m.visible) {
		selected := m.visible[m.cursor]
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1"))
		info := ""
		if selected.IsDir() {
			info = fmt.Sprintf("📁 %s", selected.Path())
		} else {
			info = fmt.Sprintf("%s %s (%d bytes)", treefiles.GetFileIcon(selected.Name()), selected.Path(), selected.Size())
		}
		sb.WriteString(infoStyle.Render(info))
	}

	return sb.String()
}

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// Resolve path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing path: %v\n", err)
		os.Exit(1)
	}

	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Path must be a directory: %s\n", absPath)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(absPath), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
