package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/tui/app"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs/reasoning"
	"github.com/wwsheng009/taproot/internal/tui/page"
	"github.com/wwsheng009/taproot/internal/tui/util"
)

const (
	pageHome page.PageID = "home"
)

func main() {
	application := app.NewApp()
	homePage := NewHomePage()
	application.RegisterPage(pageHome, homePage)
	application.SetPage(pageHome)

	p := tea.NewProgram(application, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

// HomePage is the home page
type HomePage struct {
	dialogOpen bool
}

func NewHomePage() HomePage {
	return HomePage{dialogOpen: false}
}

func (h HomePage) Init() tea.Cmd {
	return nil
}

func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if !h.dialogOpen {
				h.dialogOpen = true
				// Open reasoning dialog with initial content
				initialContent := "Thinking about the problem...\n\nAnalyzing requirements..."
				dialog := reasoning.New(initialContent)
				return h, tea.Sequence(
					func() tea.Msg {
						return dialogs.OpenDialogMsg{Model: dialog}
					},
					// Start streaming updates
					h.startStreamingReasoning(),
				)
			}
		case "ctrl+c", "q":
			return h, tea.Quit
		}
	case dialogs.CloseDialogMsg:
		h.dialogOpen = false
		return h, nil
	case reasoning.ReasoningUpdateMsg:
		// Pass through to the dialog
		return h, func() tea.Msg { return msg }
	}

	return h, nil
}

func (h HomePage) startStreamingReasoning() tea.Cmd {
	// Simulate streaming reasoning updates
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		// In a real app, this would come from an AI response stream
		reasoningText := `Step 1: Understanding the Problem
The user is asking about implementing a TUI framework.
Key considerations:
- Terminal compatibility
- Performance for large datasets
- Ease of use for developers

Step 2: Architecture Design
We need a modular approach:
- Component-based UI system
- Event-driven message passing
- Theme system for styling

Step 3: Implementation Strategy
Start with core components:
1. Layout management
2. Input handling
3. Rendering pipeline

Step 4: Refinement
Add advanced features:
- Dialog system
- Virtual scrolling
- Accessibility support

Conclusion: Taproot provides a solid foundation
for building terminal applications in Go.`
		return reasoning.UpdateReasoning(reasoningText)
	})
}

func (h HomePage) View() string {
	t := lipgloss.NewStyle()
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Reasoning Dialog Demo")

	var b strings.Builder
	b.WriteString(title + "\n\n")
	b.WriteString("The reasoning dialog displays collapsible thought processes.\n\n")
	b.WriteString("Features:\n")
	b.WriteString("  - Collapsible content (Enter to toggle)\n")
	b.WriteString("  - Streaming content updates\n")
	b.WriteString("  - Scrollable when expanded\n")
	b.WriteString("  - Auto-truncation for long text\n\n")
	b.WriteString("Press r to open reasoning dialog\n")
	b.WriteString("Press q or ctrl+c to quit")

	return b.String()
}
