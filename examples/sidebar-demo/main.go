package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/components/sidebar"
	"github.com/wwsheng009/taproot/internal/ui/render"
)

var (
	labelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	dimStyle = lipgloss.NewStyle().Faint(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("239"))

	mainAreaStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("239")).
		Padding(1)
)

type model struct {
	sidebar      sidebar.Sidebar
	screenWidth  int
	screenHeight int
}

func (m *model) Init() error {
	if err := m.sidebar.Init(); err != nil {
		return err
	}

	// Initialize with sample data
	m.loadSampleData()

	return nil
}

func (m *model) loadSampleData() {
	// Set model info
	m.sidebar.SetModelInfo(sidebar.ModelInfo{
		Name:         "gpt-4-turbo",
		Icon:         "M",
		Provider:     "openai",
		CanReason:    false,
		ContextWindow: 128000,
	})

	// Set session info
	m.sidebar.SetSession(sidebar.SessionInfo{
		ID:              "demo-session-123",
		Title:           "Sidebar Demo",
		PromptTokens:    64000,
		CompletionTokens: 32000,
		Cost:            0.75,
		WorkingDir:      "/home/user/taproot",
	})

	// Add some files
	m.sidebar.AddFile(sidebar.FileInfo{
		Path:      "internal/ui/layout/area.go",
		Additions: 15,
		Deletions: 3,
	})
	m.sidebar.AddFile(sidebar.FileInfo{
		Path:      "internal/ui/layout/split.go",
		Additions: 8,
		Deletions: 0,
	})
	m.sidebar.AddFile(sidebar.FileInfo{
		Path:      "internal/ui/layout/flex.go",
		Additions: 20,
		Deletions: 5,
	})
	m.sidebar.AddFile(sidebar.FileInfo{
		Path:      "internal/ui/layout/grid.go",
		Additions: 12,
		Deletions: 2,
	})

	// Set LSP status
	m.sidebar.SetLSPStatus([]sidebar.LSPService{
		{Name: "gopls", Language: "go", Connected: true, ErrorCount: 0},
		{Name: "pyright", Language: "python", Connected: true, ErrorCount: 0},
		{Name: "ts-language-server", Language: "typescript", Connected: false, ErrorCount: 0},
	})

	// Set MCP status
	m.sidebar.SetMCPStatus([]sidebar.MCPService{
		{Name: "filesystem", Connected: true},
		{Name: "github", Connected: true},
		{Name: "postgres", Connected: false},
	})
}

func (m *model) Update(msg any) (render.Model, render.Cmd) {
	switch msg := msg.(type) {
	case render.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, render.Quit()
		case "c":
			// Toggle compact mode
			m.sidebar.SetCompactMode(!m.isCompact())
		case "f":
			// Add more files
			m.sidebar.AddFile(sidebar.FileInfo{
				Path:      "examples/layout-demo/main.go",
				Additions: 30,
				Deletions: 5,
			})
			m.sidebar.AddFile(sidebar.FileInfo{
				Path:      "internal/ui/render/adapter_tea.go",
				Additions: 10,
				Deletions: 8,
			})
			m.sidebar.AddFile(sidebar.FileInfo{
				Path:      "examples/sidebar-demo/main.go",
				Additions: 50,
				Deletions: 0,
			})
		case "s":
			// Clear files
			m.sidebar.ClearFiles()
		case "u":
			// Update session (simulate progress)
			m.sidebar.SetSession(sidebar.SessionInfo{
				ID:              "demo-session-123",
				Title:           "Sidebar Demo (Updated)",
				PromptTokens:    75000,
				CompletionTokens: 45000,
				Cost:            1.25,
				WorkingDir:      "/home/user/taproot",
			})
		case "l":
			// Reset to sample data
			m.sidebar.ClearFiles()
			m.loadSampleData()
		}
	case render.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height

		// Calculate sidebar width (30 columns or 1/3 of screen, whichever is smaller)
		sidebarWidth := 30
		if msg.Width < sidebarWidth+20 {
			sidebarWidth = msg.Width - 10
		}

		m.sidebar.SetSize(sidebarWidth, msg.Height-3) // Reserve 3 lines for header
	}

	newSidebar, _ := m.sidebar.Update(msg)
	m.sidebar = newSidebar

	return m, nil
}

func (m *model) isCompact() bool {
	// Try to detect compact mode from the sidebar's current view
	// This is a simple heuristic - the actual state might be stored internally
	return m.screenWidth > 0 && m.screenWidth < 40
}

func (m *model) View() string {
	var b strings.Builder

	// Render title
	b.WriteString(labelStyle.Render("Taproot Sidebar Demo"))
	b.WriteString("\n\n")

	// Render help
	help := []string{
		"c: Toggle Compact Mode",
		"f: Add More Files",
		"s: Clear Files",
		"u: Update Session",
		"l: Load Sample Data",
		"q: Quit",
	}

	helpLine := strings.Join(help, " | ")
	b.WriteString(helpStyle.Render(helpLine))
	b.WriteString("\n\n")

	// Calculate layout
	sidebarWidth := 30
	if m.screenWidth > 0 && m.screenWidth < sidebarWidth+20 {
		sidebarWidth = m.screenWidth - 10
	}

	mainWidth := m.screenWidth - sidebarWidth - 2
	if mainWidth < 10 {
		mainWidth = 10
	}

	// Render sidebar
	sidebarView := m.sidebar.View()

	// Render main area
	mainContent := m.renderMainArea()
	mainView := mainAreaStyle.
		Width(mainWidth).
		Height(m.screenHeight - 5).
		Render(mainContent)

	// Join sidebar and main area horizontally
	fullView := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, mainView)

	b.WriteString(fullView)

	return b.String()
}

func (m *model) renderMainArea() string {
	var b strings.Builder

	b.WriteString(labelStyle.Render("Main Content Area"))
	b.WriteString("\n\n")

	b.WriteString("This area represents the main application interface.")
	b.WriteString("\n\n")

	w, h := m.sidebar.Size()
	info := fmt.Sprintf("Sidebar Size: %dx%d\n", w, h)
	b.WriteString(info)

	info2 := "Demo Instructions:\n"
	b.WriteString(info2)

	instructions := []string{
		"1. Use 'c' to toggle compact mode",
		"2. Use 'f' to add more files",
		"3. Use 's' to clear all files",
		"4. Use 'u' to update session info",
		"5. Use 'l' to reload sample data",
		"6. Try resizing the terminal",
	}

	for _, instr := range instructions {
		b.WriteString(instr)
		b.WriteString("\n")
	}

	return b.String()
}

func main() {
	engine, err := render.CreateEngine(render.EngineBubbletea, render.DefaultConfig())
	if err != nil {
		panic(err)
	}

	conf := sidebar.DefaultConfig()
	s := sidebar.New(conf)

	m := &model{
		sidebar:      s,
		screenWidth:  100,
		screenHeight: 40,
	}

	if err := engine.Start(m); err != nil {
		panic(err)
	}
}
