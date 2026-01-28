package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/tui/app"
	"github.com/yourorg/taproot/internal/tui/components/dialogs"
	"github.com/yourorg/taproot/internal/tui/components/dialogs/sessions"
	"github.com/yourorg/taproot/internal/tui/page"
	"github.com/yourorg/taproot/internal/tui/util"
)

const (
	pageHome page.PageID = "home"
)

func main() {
	application := app.NewApp()
	provider := NewMySessionProvider()
	homePage := NewHomePage(provider)
	application.RegisterPage(pageHome, homePage)
	application.SetPage(pageHome)

	p := tea.NewProgram(application, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

// MySessionProvider implements the SessionProvider interface
type MySessionProvider struct {
	sessions      []sessions.Session
	currentID     string
	nextID        int
}

func NewMySessionProvider() *MySessionProvider {
	now := time.Now()
	return &MySessionProvider{
		sessions: []sessions.Session{
			{
				ID:        "sess-1",
				Title:     "Project Planning",
				CreatedAt: now.Add(-24 * time.Hour),
				UpdatedAt: now.Add(-2 * time.Hour),
			},
			{
				ID:        "sess-2",
				Title:     "Code Review Discussion",
				CreatedAt: now.Add(-12 * time.Hour),
				UpdatedAt: now.Add(-1 * time.Hour),
			},
			{
				ID:        "sess-3",
				Title:     "Algorithm Design",
				CreatedAt: now.Add(-6 * time.Hour),
				UpdatedAt: now.Add(-30 * time.Minute),
			},
			{
				ID:        "sess-4",
				Title:     "Debugging Session",
				CreatedAt: now.Add(-3 * time.Hour),
				UpdatedAt: now,
			},
		},
		currentID: "sess-4",
		nextID:    5,
	}
}

func (p *MySessionProvider) Sessions() []sessions.Session {
	return p.sessions
}

func (p *MySessionProvider) CreateSession(title string) tea.Cmd {
	p.nextID++
	newSession := sessions.Session{
		ID:        fmt.Sprintf("sess-%d", p.nextID),
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	p.sessions = append([]sessions.Session{newSession}, p.sessions...)
	p.currentID = newSession.ID
	return util.ReportSuccess(fmt.Sprintf("Created session: %s", title))
}

func (p *MySessionProvider) DeleteSession(id string) tea.Cmd {
	for i, s := range p.sessions {
		if s.ID == id {
			p.sessions = append(p.sessions[:i], p.sessions[i+1:]...)
			return util.ReportInfo(fmt.Sprintf("Deleted session: %s", s.Title))
		}
	}
	return nil
}

func (p *MySessionProvider) SelectSession(id string) tea.Cmd {
	p.currentID = id
	for _, s := range p.sessions {
		if s.ID == id {
			return util.ReportInfo(fmt.Sprintf("Switched to: %s", s.Title))
		}
	}
	return nil
}

// HomePage is the home page
type HomePage struct {
	provider *MySessionProvider
}

func NewHomePage(provider *MySessionProvider) HomePage {
	return HomePage{provider: provider}
}

func (h HomePage) Init() tea.Cmd {
	return nil
}

func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			// Open sessions dialog
			dialog := sessions.New(h.provider)
			return h, func() tea.Msg {
				return dialogs.OpenDialogMsg{Model: dialog}
			}
		case "ctrl+c", "q":
			return h, tea.Quit
		}
	}

	return h, nil
}

func (h HomePage) View() string {
	t := lipgloss.NewStyle()
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Sessions Dialog Demo")

	var b strings.Builder
	b.WriteString(title + "\n\n")

	// Find current session
	var currentTitle string
	for _, s := range h.provider.Sessions() {
		if s.ID == h.provider.currentID {
			currentTitle = s.Title
			break
		}
	}

	b.WriteString(fmt.Sprintf("Current Session: %s\n\n", currentTitle))
	b.WriteString("Press ctrl+s to open sessions dialog\n\n")
	b.WriteString("Features:\n")
	b.WriteString("  - Browse all sessions\n")
	b.WriteString("  - Filter by typing\n")
	b.WriteString("  - Create new session (n)\n")
	b.WriteString("  - Delete session (d)\n")
	b.WriteString("  - View timestamps\n\n")
	b.WriteString("Available sessions:\n")
	for _, s := range h.provider.Sessions() {
		cursor := " "
		if s.ID == h.provider.currentID {
			cursor = ">"
		}
		fmt.Fprintf(&b, "  %s %s (%s)\n", cursor, s.Title, formatTimeAgo(s.UpdatedAt))
	}
	b.WriteString("\nPress q or ctrl+c to quit")

	return b.String()
}

func formatTimeAgo(t time.Time) string {
	diff := time.Since(t)
	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
}
