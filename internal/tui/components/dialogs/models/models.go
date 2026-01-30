package models

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs"
	"github.com/wwsheng009/taproot/internal/ui/styles"
	"github.com/wwsheng009/taproot/internal/tui/util"
)

const (
	ModelDialogID dialogs.DialogID = "models"
	defaultWidth  int               = 70
)

// Model represents an AI model
type Model struct {
	ID          string
	Name        string
	Provider    string
	ContextSize int
	Description string
	// Recently used timestamp
	LastUsed time.Time
}

// ModelProvider provides available models
type ModelProvider interface {
	// Models returns a list of available models
	Models() []Model
	// RecentModels returns recently used models
	RecentModels() []Model
	// SetModel sets the active model
	SetModel(modelID string) tea.Cmd
}

// ModelsDialog represents the model selection dialog
type ModelsDialog interface {
	dialogs.DialogModel
}

type modelsDialogCmp struct {
	styles       *styles.Styles
	width        int
	height       int
	x, y         int
	provider     ModelProvider
	allModels    []Model
	recentModels []Model
	filtered     []Model
	selectedIdx  int
	showRecent   bool
	query        string
}

// NewModelsDialog creates a new model selection dialog
func NewModelsDialog(provider ModelProvider) ModelsDialog {
	allModels := provider.Models()
	recentModels := provider.RecentModels()
	s := styles.DefaultStyles()

	return &modelsDialogCmp{
		styles:       &s,
		width:        defaultWidth,
		height:       20,
		x:            10,
		y:            5,
		provider:     provider,
		allModels:    allModels,
		recentModels: recentModels,
		filtered:     allModels,
		selectedIdx:  0,
		showRecent:   len(recentModels) > 0,
		query:        "",
	}
}

func (d *modelsDialogCmp) Init() tea.Cmd {
	return nil
}

func (d *modelsDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return d, func() tea.Msg {
				return dialogs.CloseDialogMsg{}
			}
		case "tab":
			// Toggle between all models and recent
			d.showRecent = !d.showRecent
			d.filterModels()
			d.selectedIdx = 0
			return d, nil
		case "enter":
			if d.selectedIdx < len(d.filtered) {
				model := d.filtered[d.selectedIdx]
				// Update last used time
				for i, m := range d.allModels {
					if m.ID == model.ID {
						d.allModels[i].LastUsed = time.Now()
						break
					}
				}
				// Call provider to set model
				cmd := d.provider.SetModel(model.ID)
				// Close dialog
				return d, tea.Batch(
					func() tea.Msg { return dialogs.CloseDialogMsg{} },
					cmd,
					util.ReportInfo(fmt.Sprintf("Model changed to %s", model.Name)),
				)
			}
		case "up", "k":
			if d.selectedIdx > 0 {
				d.selectedIdx--
			}
		case "down", "j":
			if d.selectedIdx < len(d.filtered)-1 {
				d.selectedIdx++
			}
		default:
			// Handle filtering
			if len(msg.String()) == 1 {
				d.query += msg.String()
				d.filterModels()
				d.selectedIdx = 0
				return d, nil
			}
			if msg.String() == "backspace" || msg.String() == "ctrl+h" {
				if len(d.query) > 0 {
					d.query = d.query[:len(d.query)-1]
					d.filterModels()
					d.selectedIdx = 0
				}
				return d, nil
			}
		}
	}

	return d, nil
}

func (d *modelsDialogCmp) filterModels() {
	source := d.allModels
	if d.showRecent {
		source = d.recentModels
	}

	if d.query == "" {
		d.filtered = source
		return
	}

	query := strings.ToLower(d.query)
	d.filtered = []Model{}
	for _, model := range source {
		searchText := fmt.Sprintf("%s %s %s", model.Name, model.Provider, model.Description)
		if strings.Contains(strings.ToLower(searchText), query) {
			d.filtered = append(d.filtered, model)
		}
	}
}

func (d *modelsDialogCmp) View() string {
	s := d.styles

	// Dialog box style
	boxStyle := s.Dialog.View.Padding(0, 1)

	var content strings.Builder

	// Header
	title := "Model Selection"
	if d.showRecent {
		title += " (Recent)"
	}
	header := s.Base.Bold(true).Foreground(s.Primary).Render(title)
	content.WriteString(header + "\n\n")

	// Filter input
	filterLabel := s.Base.Foreground(s.FgMuted).Render("Search: ")
	content.WriteString(filterLabel + d.query + "_\n\n")

	// Toggle hint
	toggleHint := s.Base.Foreground(s.FgSubtle).Render("Tab: Toggle Recent/All | Enter: Select | ESC: Close")
	content.WriteString(toggleHint + "\n\n")

	// Model list
	visible := d.height - 6 // Account for header, filter, hints, footer
	start := (d.selectedIdx / visible) * visible
	end := min(start+visible, len(d.filtered))

	for i := start; i < end; i++ {
		if i >= len(d.filtered) {
			break
		}

		model := d.filtered[i]
		prefix := " "
		if i == d.selectedIdx {
			prefix = ">"
		}

		itemStyle := s.Dialog.NormalItem
		if i == d.selectedIdx {
			itemStyle = s.Dialog.SelectedItem
		}

		// Model info
		line := fmt.Sprintf("%s %s", prefix, model.Name)
		if model.Provider != "" {
			line += fmt.Sprintf(" [%s]", model.Provider)
		}
		if model.ContextSize > 0 {
			line += fmt.Sprintf(" (%dK)", model.ContextSize/1024)
		}

		content.WriteString(itemStyle.Render(line) + "\n")

		// Description on next line if selected
		if i == d.selectedIdx && model.Description != "" {
			descStyle := s.Muted.Italic(true)
			content.WriteString(descStyle.Render("  " + model.Description) + "\n")
		}
	}

	// Footer
	source := d.allModels
	if d.showRecent {
		source = d.recentModels
	}
	footerCount := fmt.Sprintf("%d/%d models", len(d.filtered), len(source))
	footer := s.Muted.Render(footerCount)
	content.WriteString("\n" + footer)

	return boxStyle.Render(content.String())
}

func (d *modelsDialogCmp) Position() (int, int) {
	return d.x, d.y
}

func (d *modelsDialogCmp) ID() dialogs.DialogID {
	return ModelDialogID
}

