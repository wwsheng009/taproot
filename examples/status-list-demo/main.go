package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/components/status"
	"github.com/wwsheng009/taproot/ui/styles"
)

type model struct {
	lspList      *status.LSPList
	mcpList      *status.MCPList
	screenWidth  int
	screenHeight int
	showTitles   bool
	initialized  bool
}

func initialModel() *model {
	// Create LSP list
	lspList := status.NewLSPList()
	lspList.Init()
	lspList.SetWidth(50)
	lspList.SetShowTitle(true)

	// Add some initial LSP services
	lspList.AddService(status.LSPServiceInfo{
		Name:     "gopls",
		Language: "go",
		State:    status.StateReady,
		Diagnostics: status.DiagnosticSummary{
			Error:   0,
			Warning: 2,
			Hint:    5,
		},
	})
	lspList.AddService(status.LSPServiceInfo{
		Name:     "rust-analyzer",
		Language: "rust",
		State:    status.StateReady,
		Diagnostics: status.DiagnosticSummary{
			Error:   1,
			Warning: 0,
			Hint:    2,
		},
	})
	lspList.AddService(status.LSPServiceInfo{
		Name:     "pylsp",
		Language: "python",
		State:    status.StateStarting,
	})

	// Create MCP list
	mcpList := status.NewMCPList()
	mcpList.Init()
	mcpList.SetWidth(50)
	mcpList.SetShowTitle(true)

	// Add some initial MCP services
	mcpList.AddService(status.MCPServiceInfo{
		Name:  "filesystem",
		State: status.StateReady,
		ToolCounts: status.ToolCounts{
			Tools:   5,
			Prompts: 0,
		},
	})
	mcpList.AddService(status.MCPServiceInfo{
		Name:  "git",
		State: status.StateReady,
		ToolCounts: status.ToolCounts{
			Tools:   3,
			Prompts: 1,
		},
	})

	return &model{
		lspList:      lspList,
		mcpList:      mcpList,
		screenWidth:  80,
		screenHeight: 24,
		showTitles:   true,
		initialized:  false,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q":
			return m, tea.Quit
		case "t", "T":
			m.showTitles = !m.showTitles
			m.lspList.SetShowTitle(m.showTitles)
			m.mcpList.SetShowTitle(m.showTitles)
		case "1":
			// Add a new LSP service
			m.lspList.AddService(status.LSPServiceInfo{
				Name:     "typescript-language-server",
				Language: "typescript",
				State:    status.StateReady,
				Diagnostics: status.DiagnosticSummary{
					Error:   0,
					Warning: 1,
					Hint:    3,
				},
			})
		case "2":
			// Add an LSP service with errors
			m.lspList.AddService(status.LSPServiceInfo{
				Name:     "clangd",
				Language: "c++",
				State:    status.StateError,
				Error:    "failed to start",
				Diagnostics: status.DiagnosticSummary{
					Error: 3,
				},
			})
		case "3":
			// Add a disabled LSP service
			m.lspList.AddService(status.LSPServiceInfo{
				Name:     "hls",
				Language: "haskell",
				State:    status.StateDisabled,
			})
		case "4":
			// Add a starting LSP service
			m.lspList.AddService(status.LSPServiceInfo{
				Name:     "jdtls",
				Language: "java",
				State:    status.StateStarting,
			})
		case "5":
			// Clear all LSP services
			m.lspList.ClearServices()
		case "a", "A":
			// Add a new MCP service
			m.mcpList.AddService(status.MCPServiceInfo{
				Name:  "database",
				State: status.StateReady,
				ToolCounts: status.ToolCounts{
					Tools:   8,
					Prompts: 2,
				},
			})
		case "s", "S":
			// Add a starting MCP service
			m.mcpList.AddService(status.MCPServiceInfo{
				Name:  "http",
				State: status.StateStarting,
			})
		case "e", "E":
			// Add an error MCP service
			m.mcpList.AddService(status.MCPServiceInfo{
				Name:  "ssh",
				State: status.StateError,
				Error: "connection refused",
			})
		case "d", "D":
			// Add a disabled MCP service
			m.mcpList.AddService(status.MCPServiceInfo{
				Name:  "s3",
				State: status.StateDisabled,
			})
		case "c", "C":
			// Clear all MCP services
			m.mcpList.ClearServices()
		case "+", "=":
			// Increase max items
			m.lspList.SetMaxItems(m.lspList.MaxItems() + 1)
			m.mcpList.SetMaxItems(m.mcpList.MaxItems() + 1)
		case "-", "_":
			// Decrease max items
			if m.lspList.MaxItems() > 1 {
				m.lspList.SetMaxItems(m.lspList.MaxItems() - 1)
				m.mcpList.SetMaxItems(m.mcpList.MaxItems() - 1)
			}
		case "w", "W":
			// Increase width
			m.lspList.SetWidth(m.lspList.Width() + 5)
			m.mcpList.SetWidth(m.mcpList.Width() + 5)
		case "n", "N":
			// Decrease width
			if m.lspList.Width() > 20 {
				m.lspList.SetWidth(m.lspList.Width() - 5)
				m.mcpList.SetWidth(m.mcpList.Width() - 5)
			}
		}

	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
	}

	// Update lists
	m.lspList.Update(msg)
	m.mcpList.Update(msg)

	return m, nil
}

func (m *model) View() string {
	sty := styles.DefaultStyles()

	var b strings.Builder

	// Title
	title := sty.Base.Bold(true).Render("Status Components Demo")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Help text
	help := sty.Subtle.Render("Controls: " +
		"1-4 add LSP | 5 clear LSP | a/e/d/s add MCP | c clear MCP | " +
		"+/- max items | w/n width | t toggle titles | q quit")
	b.WriteString(help)
	b.WriteString("\n\n")

	// LSP and MCP lists side by side
	lspView := m.lspList.View()
	mcpView := m.mcpList.View()

	// Combine views
	b.WriteString(lspView)
	b.WriteString("\n\n")
	b.WriteString(mcpView)

	// Status info at bottom
	b.WriteString("\n\n")
	statusLine := sty.Subtle.Render(fmt.Sprintf(
		"LSPs: %d online | MCPs: %d connected | Width: %d | Max items: %d",
		m.lspList.OnlineCount(),
		m.mcpList.ConnectedCount(),
		m.lspList.Width(),
		m.lspList.MaxItems(),
	))
	b.WriteString(statusLine)

	return sty.Base.Render(b.String())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
