package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/components/status"
	"github.com/wwsheng009/taproot/ui/styles"
)

type model struct {
	services       []*status.ServiceCmp
	diagnostics    *status.DiagnosticStatusCmp
	serviceIndex   int
	screenWidth    int
	screenHeight   int}

func initialModel() *model {
	// Create LSP services
	lspGo := status.NewService("lsp-go", "Go LSP")
	lspGo.Init()

	lspTS := status.NewService("lsp-ts", "TypeScript LSP")
	lspTS.Init()

	// Create MCP services
	mcpFs := status.NewService("mcp-fs", "File System MCP")
	mcpFs.Init()

	mcpGit := status.NewService("mcp-git", "Git MCP")
	mcpGit.Init()

	// Create diagnostic status
	diag := status.NewDiagnosticStatus("workspace")
	diag.Init()

	return &model{
		services:     []*status.ServiceCmp{lspGo, lspTS, mcpFs, mcpGit},
		diagnostics:  diag,
		serviceIndex: 0,
		screenWidth:  80,
		screenHeight: 24,
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
		case "up", "k":
			m.serviceIndex--
			if m.serviceIndex < 0 {
				m.serviceIndex = len(m.services) - 1
			}
		case "down", "j":
			m.serviceIndex++
			if m.serviceIndex >= len(m.services) {
				m.serviceIndex = 0
			}
		case "1":
			m.services[0].SetStatus(status.ServiceStatusOffline)
			m.services[0].SetErrorCount(0)
		case "2":
			m.services[0].SetStatus(status.ServiceStatusStarting)
		case "3":
			m.services[0].SetStatus(status.ServiceStatusConnecting)
		case "4":
			m.services[0].SetStatus(status.ServiceStatusOnline)
			m.services[0].SetErrorCount(0)
		case "5":
			m.services[0].SetStatus(status.ServiceStatusBusy)
		case "6":
			m.services[0].SetStatus(status.ServiceStatusError)
			m.services[0].SetErrorCount(3)
		case "e", "E":
			// Add error to current service
			svc := m.services[m.serviceIndex]
			svc.SetErrorCount(svc.ErrorCount() + 1)
			m.diagnostics.AddDiagnostic(status.DiagnosticSeverityError)
		case "w", "W":
			// Add warning to diagnostics
			m.diagnostics.AddDiagnostic(status.DiagnosticSeverityWarning)
		case "i", "I":
			// Add info to diagnostics
			m.diagnostics.AddDiagnostic(status.DiagnosticSeverityInfo)
		case "h", "H":
			// Add hint to diagnostics
			m.diagnostics.AddDiagnostic(status.DiagnosticSeverityHint)
		case "c", "C":
			// Clear all diagnostics
			m.diagnostics.Clear()
			for _, svc := range m.services {
				svc.SetErrorCount(0)
			}
		case "m", "M":
			// Toggle compact mode for all services
			compact := !m.services[0].Compact()
			for _, svc := range m.services {
				svc.SetCompact(compact)
			}
		case "d", "D":
			// Toggle compact mode for diagnostics
			m.diagnostics.SetCompact(!m.diagnostics.Compact())
		}

	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
		// Update service max widths
		for _, svc := range m.services {
			svc.SetMaxWidth(m.screenWidth - 10)
		}
		m.diagnostics.SetMaxWidth(m.screenWidth - 10)
	}

	return m, nil
}

func (m *model) View() string {
	sty := styles.DefaultStyles()
	var b strings.Builder

	// Calculate padding
	padding := strings.Repeat(" ", 2)

	// Title
	b.WriteString(sty.Section.Title.Render(padding + "Status Component Demo"))
	b.WriteString("\n\n")

	// Services section
	b.WriteString(padding)
	b.WriteString(sty.Section.Title.Render("LSP/MCP Services"))
	b.WriteString("\n\n")

	// Render service list
	for i, svc := range m.services {
		// Focus/unfocus service
		if i == m.serviceIndex {
			svc.Focus()
		} else {
			svc.Blur()
		}

		// Render service with proper indentation
		serviceLine := padding + padding + svc.View()
		b.WriteString(serviceLine)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Diagnostics section
	b.WriteString(padding)
	b.WriteString(sty.Section.Title.Render("Diagnostics"))
	b.WriteString("\n\n")

	// Render diagnostics
	diagLine := padding + padding + m.diagnostics.View()
	b.WriteString(diagLine)
	b.WriteString("\n\n")

	// Instructions section
	b.WriteString(padding)
	b.WriteString(sty.Section.Title.Render("Keyboard Controls"))
	b.WriteString("\n\n")

	instructions := []struct {
		key     string
		desc    string
	}{
		{"Q / Ctrl+C", "Quit"},
		{"↑ / ↓ or K/J", "Navigate services"},
		{"1-6", "Set service status (1=Offline, 2=Starting, 3=Connecting, 4=Online, 5=Busy, 6=Error)"},
		{"E", "Add error to current service"},
		{"W", "Add warning to diagnostics"},
		{"I", "Add info to diagnostics"},
		{"H", "Add hint to diagnostics"},
		{"C", "Clear all diagnostics"},
		{"M", "Toggle compact mode for services"},
		{"D", "Toggle compact mode for diagnostics"},
	}

	for _, instr := range instructions {
		b.WriteString(padding + padding)
		b.WriteString(fmt.Sprintf("%-20s - %s\n", instr.key, instr.desc))
	}

	// Current state info
	b.WriteString("\n")
	b.WriteString(padding)
	b.WriteString(sty.Section.Title.Render("Current State"))
	b.WriteString("\n\n")

	b.WriteString(padding + padding)
	b.WriteString(fmt.Sprintf("Focused Service: %s\n", m.services[m.serviceIndex].Name()))
	b.WriteString(padding + padding)
	b.WriteString(fmt.Sprintf("Service Status: %s\n", m.services[m.serviceIndex].Status()))
	b.WriteString(padding + padding)
	b.WriteString(fmt.Sprintf("Service Errors: %d\n", m.services[m.serviceIndex].ErrorCount()))
	b.WriteString(padding + padding)
	b.WriteString(fmt.Sprintf("Total Diagnostics: %d\n", m.diagnostics.Total()))
	b.WriteString(padding + padding)
	b.WriteString(fmt.Sprintf("Has Problems: %t\n", m.diagnostics.HasProblems()))
	b.WriteString(padding + padding)
	b.WriteString(fmt.Sprintf("Services Compact: %t\n", m.services[0].Compact()))
	b.WriteString(padding + padding)
	b.WriteString(fmt.Sprintf("Diagnostics Compact: %t\n", m.diagnostics.Compact()))

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
