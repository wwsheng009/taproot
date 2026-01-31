package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/notify"
	"github.com/wwsheng009/taproot/ui/render"
)

type model struct {
	notifyManager *notify.Manager
	quitting      bool
	width         int
	height        int
}

func initialModel() model {
	cfg := notify.DefaultConfig()
	cfg.Position = notify.TopRight
	return model{
		notifyManager: notify.NewManager(cfg),
	}
}

func (m model) Init() tea.Cmd {
	_ = m.notifyManager.Init()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "i":
			cmds = append(cmds, adaptCmd(notify.Info("Information", "This is an info notification")))
		case "s":
			cmds = append(cmds, adaptCmd(notify.Success("Success", "Operation completed successfully")))
		case "w":
			cmds = append(cmds, adaptCmd(notify.Warn("Warning", "Something might be wrong")))
		case "e":
			cmds = append(cmds, adaptCmd(notify.Error("Error", "Critical failure occurred")))
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// We DON'T pass WindowSizeMsg to manager here because we want to render it inline for this demo
		// If we passed it, Manager would try to absolute position it.
		// To test absolute positioning, we could pass it, but then we need to overlay.
		// For simplicity, let's keep it inline.
		
		// Actually, let's try to test the positioning feature.
		// If we return the manager view, it will be full screen.
		// We can render our help text as a notification? No.
		
		// Let's stick to inline for now.
	}

	// Pass everything to manager (including ticks and notification messages)
	// We need to handle ShowNotificationMsg specifically if we want to log it or something, but manager handles it.
	
	// Convert specific messages if needed? 
	// tea.WindowSizeMsg -> render.WindowSizeMsg is needed if we wanted to use it.
	
	// For now, just pass msg.
	// Note: tickMsg is internal to notify package, so we can't cast it here, 
	// but passing 'msg' (interface{}) works fine.
	
	newModel, nCmd := m.notifyManager.Update(msg)
	m.notifyManager = newModel.(*notify.Manager)
	cmds = append(cmds, adaptCmd(nCmd))

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	s := "Taproot Notification Demo\n\n"
	s += "Press i(nfo), s(uccess), w(arn), e(rror) to show notifications.\n"
	s += "Press q to quit.\n\n"
	
	// Render notifications
	return s + m.notifyManager.View()
}

func adaptCmd(cmd render.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	if fn, ok := cmd.(func() render.Msg); ok {
		return func() tea.Msg {
			return fn()
		}
	}
	return nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
