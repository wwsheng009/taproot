package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/components/files"
)

const (
	maxWidth = 80
)

type model struct {
	fileList       *files.FileList
	cursor         int
	sortBy         files.SortBy
	sortOrder      files.SortOrder
	showHidden     bool
	filter         string
	width          int
	height         int
	selectedFile   string
}

func initialModel() model {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Create file list
	fl, err := files.NewFileList(cwd,
		files.WithIncludeHidden(false),
	)
	if err != nil {
		// If we can't create a file list, create an empty one as fallback
		fl, _ = files.NewFileList(".")
	}

	return model{
		fileList:   fl,
		cursor:     0,
		sortBy:     files.SortByName,
		sortOrder:  files.SortAscending,
		showHidden: false,
		filter:     "",
		width:      maxWidth,
		height:     20,
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

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < m.fileList.FilteredCount()-1 {
				m.cursor++
			}

		case "home":
			m.cursor = 0

		case "end":
			m.cursor = m.fileList.FilteredCount() - 1
			if m.cursor < 0 {
				m.cursor = 0
			}

		case "enter":
			if item, ok := m.fileList.GetItem(m.cursor); ok {
				if item.IsDir() {
					// Change directory
					path := item.Path()
					if err := m.fileList.LoadDirectory(path); err == nil {
						m.cursor = 0
						m.selectedFile = ""
					}
				} else {
					// Select file
					m.selectedFile = item.Path()
				}
			}

		case "backspace":
			// Go to parent directory
			parentDir := parentPath(m.fileList.Path())
			if parentDir != m.fileList.Path() {
				if err := m.fileList.LoadDirectory(parentDir); err == nil {
					m.cursor = 0
				}
			}

		case "s":
			// Cycle sort by
			switch m.sortBy {
			case files.SortByName:
				m.sortBy = files.SortBySize
			case files.SortBySize:
				m.sortBy = files.SortByTime
			case files.SortByTime:
				m.sortBy = files.SortByExtension
			case files.SortByExtension:
				m.sortBy = files.SortByName
			}
			m.fileList.SetSort(m.sortBy, m.sortOrder)

		case "r":
			// Toggle sort order
			m.fileList.ToggleSortOrder()
			m.sortOrder = m.fileList.SortOrder()

		case "h":
			// Toggle hidden files
			m.fileList.ToggleHiddenFiles()
			m.showHidden = m.fileList.IncludesHidden()
			m.cursor = 0


		}

	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth)
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Width(m.width)

	pathStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("228")).
		Width(m.width)

	statStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Width(m.width)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("208")).
		Width(m.width)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Width(m.width)

	var b strings.Builder

	// Title
	title := titleStyle.Render("Taproot File Browser Demo")
	b.WriteString(title + "\n\n")

	// Current path
	path := pathStyle.Render(m.fileList.Path() + "/")
	b.WriteString(path + "\n\n")

	// Stats
	stats := m.fileList.GetStats()
	statText := fmt.Sprintf("Total: %d items (%d dirs, %d files) | Size: %s",
		stats.TotalItems,
		stats.DirectoryCount,
		stats.FileCount,
		formatSize(stats.TotalSize),
	)
	b.WriteString(statStyle.Render(statText) + "\n\n")

	// Header
	header := headerStyle.Render(fmt.Sprintf("%-35s %10s %12s %-10s", "Name", "Size", "Modified", "Type"))
	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("─", m.width) + "\n")

	// File list (limit to visible items)
	maxLines := m.height - 12
	filteredItems := m.fileList.Filtered()
	start := 0
	end := len(filteredItems)

	if len(filteredItems) > maxLines {
		if m.cursor >= maxLines/2 {
			start = m.cursor - maxLines/2
		}
		end = start + maxLines
		if end > len(filteredItems) {
			end = len(filteredItems)
			start = end - maxLines
			if start < 0 {
				start = 0
			}
		}
	}

	for i := start; i < end && i < len(filteredItems); {
		item := filteredItems[i]
		cursor := " "
		if i == m.cursor {
			cursor = "→"
		}

		name := item.Icon() + " " + item.Name()
		if len(name) > 33 {
			name = name[:30] + "..."
		}

		size := formatSize(item.Size())
		if len(size) > 8 {
			size = size[:8]
		}

		modTime := item.ModTime().Format("Jan 02 15:04")

		fileType := item.Extension()
		if item.IsDir() {
			fileType = "DIR"
		} else if fileType == "" {
			fileType = "---"
		}

		line := fmt.Sprintf("%s %-33s %10s %12s %-10s",
			cursor, name, size, modTime, fileType)
		b.WriteString(line + "\n")

		i++
	}

	// Scroll indicator if needed
	if len(filteredItems) > maxLines {
		scroll := fmt.Sprintf(" %d/%d ", m.cursor+1, len(filteredItems))
		b.WriteString(statStyle.Render(scroll) + "\n")
	}

	b.WriteString("\n")

	// Sort info
	sortText := fmt.Sprintf("Sort: %s (%s)",
		sortByLabel(m.sortBy),
		sortOrderLabel(m.sortOrder),
	)
	if m.showHidden {
		sortText += " | Hidden: shown"
	} else {
		sortText += " | Hidden: hidden"
	}
	b.WriteString(statStyle.Render(sortText) + "\n")

	// Selected file
	if m.selectedFile != "" {
		selectedStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("154")).
			Width(m.width)
		selectedText := fmt.Sprintf("Selected: %s", m.selectedFile)
		b.WriteString(selectedStyle.Render(selectedText) + "\n\n")
	} else {
		b.WriteString("\n")
	}

	// Help
	help := helpStyle.Render(strings.Join([]string{
		"↑/k: Move up  ↓/j: Move down  Enter: Open dir/Select file",
		"Home: Start    End: Bottom       Backspace: Parent dir",
		"s: Cycle sort  r: Reverse order  h: Toggle hidden",
		"ESC/q/ctrl+c: Quit",
	}, "\n"))
	b.WriteString(help)

	return b.String()
}

// sortByLabel returns a readable label for sort criteria
func sortByLabel(sortBy files.SortBy) string {
	switch sortBy {
	case files.SortByName:
		return "Name"
	case files.SortBySize:
		return "Size"
	case files.SortByTime:
		return "Time"
	case files.SortByExtension:
		return "Ext"
	default:
		return "?"
	}
}

// sortOrderLabel returns a readable label for sort order
func sortOrderLabel(sortOrder files.SortOrder) string {
	if sortOrder == files.SortAscending {
		return "ASC"
	}
	return "DESC"
}

// formatSize formats a file size in bytes to a human-readable string
func formatSize(size int64) string {
	if size < 0 {
		return "0 B"
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

// parentPath returns the parent directory of a path
func parentPath(path string) string {
	if path == "/" || path == "" {
		return path
	}
	cleanPath := strings.TrimRight(path, "/\\")
	lastSep := strings.LastIndexAny(cleanPath, "/\\")
	if lastSep == -1 {
		return "."
	}
	if lastSep == 0 {
		return "/"
	}
	return cleanPath[:lastSep]
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
