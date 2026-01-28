package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/tui/app"
	"github.com/yourorg/taproot/internal/tui/components/dialogs"
	"github.com/yourorg/taproot/internal/tui/components/dialogs/models"
	"github.com/yourorg/taproot/internal/tui/page"
	"github.com/yourorg/taproot/internal/tui/util"
)

const (
	pageHome page.PageID = "home"
)

// MyModelProvider implements the ModelProvider interface
type MyModelProvider struct {
	currentModel string
	models      []models.Model
	recent      []models.Model
}

func (p *MyModelProvider) Models() []models.Model {
	return p.models
}

func (p *MyModelProvider) RecentModels() []models.Model {
	return p.recent
}

func (p *MyModelProvider) SetModel(modelID string) tea.Cmd {
	p.currentModel = modelID
	
	// Update recent models
	for i, m := range p.models {
		if m.ID == modelID {
			p.models[i].LastUsed = p.models[i].LastUsed
			
			// Add to recent if not already there
			found := false
			for _, r := range p.recent {
				if r.ID == modelID {
					found = true
					break
				}
			}
			if !found {
				p.recent = append([]models.Model{m}, p.recent...)
				if len(p.recent) > 5 {
					p.recent = p.recent[:5]
				}
			}
			break
		}
	}
	
	return util.ReportInfo(fmt.Sprintf("Model changed to: %s", modelID))
}

func NewMyModelProvider() *MyModelProvider {
	allModels := []models.Model{
		{
			ID:          "gpt-4",
			Name:        "GPT-4",
			Provider:    "OpenAI",
			ContextSize: 8192,
			Description: "Most capable model, great for complex tasks",
		},
		{
			ID:          "gpt-3.5-turbo",
			Name:        "GPT-3.5 Turbo",
			Provider:    "OpenAI",
			ContextSize: 4096,
			Description: "Fast and cost-effective for most tasks",
		},
		{
			ID:          "claude-3-opus",
			Name:        "Claude 3 Opus",
			Provider:    "Anthropic",
			ContextSize: 200000,
			Description: "High-end model with excellent reasoning",
		},
		{
			ID:          "claude-3-sonnet",
			Name:        "Claude 3 Sonnet",
			Provider:    "Anthropic",
			ContextSize: 200000,
			Description: "Balanced model for good performance",
		},
		{
			ID:          "gemini-pro",
			Name:        "Gemini Pro",
			Provider:    "Google",
			ContextSize: 32000,
			Description: "Multimodal capabilities with good performance",
		},
		{
			ID:          "llama-3-70b",
			Name:        "Llama 3 70B",
			Provider:    "Meta",
			ContextSize: 8192,
			Description: "Open-source model, can run locally",
		},
	}

	return &MyModelProvider{
		currentModel: "gpt-3.5-turbo",
		models:      allModels,
		recent:      []models.Model{allModels[0], allModels[2]}, // Some recent models
	}
}

func main() {
	application := app.NewApp()
	provider := NewMyModelProvider()

	homePage := NewHomePage(provider)
	application.RegisterPage(pageHome, homePage)
	application.SetPage(pageHome)

	p := tea.NewProgram(application, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

// HomePage is the home page
type HomePage struct {
	provider *MyModelProvider
}

func NewHomePage(provider *MyModelProvider) HomePage {
	return HomePage{provider: provider}
}

func (h HomePage) Init() tea.Cmd {
	return nil
}

func (h HomePage) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+m":
			// Open model selection dialog
			dialog := models.NewModelsDialog(h.provider)
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
	title := t.Bold(true).Foreground(lipgloss.Color("86")).Render("Taproot Model Selection Demo")

	var b strings.Builder
	b.WriteString(title + "\n\n")
	b.WriteString(fmt.Sprintf("Current Model: %s\n\n", h.provider.currentModel))
	b.WriteString("Press ctrl+m to open the model selection dialog\n\n")
	b.WriteString("Features:\n")
	b.WriteString("  - Search models by name, provider, or description\n")
	b.WriteString("  - View recent models with Tab\n")
	b.WriteString("  - See model details (context size, description)\n")
	b.WriteString("  - Quick select with Enter\n\n")
	b.WriteString("Available models:\n")
	for _, m := range h.provider.models {
		b.WriteString(fmt.Sprintf("  - %s (%s)\n", m.Name, m.Provider))
	}
	b.WriteString("\nPress q or ctrl+c to quit")

	return b.String()
}
