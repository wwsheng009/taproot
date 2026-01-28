package commands

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/tui/components/completions"
	"github.com/yourorg/taproot/internal/tui/components/dialogs"
	"github.com/yourorg/taproot/internal/tui/styles"
	"github.com/yourorg/taproot/internal/tui/util"
)

const (
	CommandDialogID dialogs.DialogID = "commands"
	defaultWidth    int               = 60
)

// Command represents an executable command
type Command struct {
	ID          string
	Title       string
	Description string
	// Callback is called when the command is executed
	// Returns a command to send back to the parent
	Callback func() tea.Cmd
}

// CommandProvider provides commands to the dialog
type CommandProvider interface {
	// Commands returns a list of available commands
	Commands() []Command
}

// CommandsDialog represents the command palette dialog
type CommandsDialog interface {
	dialogs.DialogModel
}

type commandDialogCmp struct {
	width       int
	height      int
	x, y        int
	commands    []Command
	completions completions.CompletionsCmp
	filtered    []Command
	selectedIdx int
	visible     int
}

// NewCommandsDialog creates a new command palette dialog
func NewCommandsDialog(provider CommandProvider) CommandsDialog {
	commands := provider.Commands()
	
	// Create completions for command search
	items := make([]completions.CompletionItem, 0, len(commands))
	for _, cmd := range commands {
		items = append(items, completions.NewCompletionItem(
			cmd.ID,
			cmd.Title,
			cmd,
		))
	}

	return &commandDialogCmp{
		width:       defaultWidth,
		height:      20,
		x:           10,
		y:           5,
		commands:    commands,
		filtered:    commands,
		completions: completions.New(),
		visible:     10,
	}
}

func (d *commandDialogCmp) Init() tea.Cmd {
	return nil
}

func (d *commandDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return d, func() tea.Msg {
				return dialogs.CloseDialogMsg{}
			}
		case "enter":
			if d.selectedIdx < len(d.filtered) {
				cmd := d.filtered[d.selectedIdx]
				if cmd.Callback != nil {
					// Close dialog and execute command
					return d, tea.Batch(
						func() tea.Msg { return dialogs.CloseDialogMsg{} },
						cmd.Callback(),
					)
				}
				return d, func() tea.Msg {
					return dialogs.CloseDialogMsg{}
				}
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
			// Forward to completions for filtering
			if len(msg.String()) == 1 || msg.String() == "backspace" || msg.String() == "ctrl+h" {
				newComps, cmd := d.completions.Update(msg)
				d.completions = newComps.(completions.CompletionsCmp)
				
				// Update filtered commands based on completions query
				d.filterCommands()
				
				return d, cmd
			}
		}
	}

	return d, nil
}

func (d *commandDialogCmp) filterCommands() {
	query := d.completions.Query()
	if query == "" {
		d.filtered = d.commands
		d.selectedIdx = 0
		return
	}

	d.filtered = []Command{}
	for _, cmd := range d.commands {
		// Simple substring match on title and description
		if contains(cmd.Title, query) || contains(cmd.Description, query) {
			d.filtered = append(d.filtered, cmd)
		}
	}
	d.selectedIdx = 0
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		len(substr) == 0 || 
		indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func (d *commandDialogCmp) View() string {
	t := styles.CurrentTheme()

	// Dialog box style
	boxStyle := t.S().Base.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Padding(0, 1)

	var content strings.Builder

	// Header
	header := t.S().Base.Bold(true).Foreground(t.Primary).Render("Commands")
	content.WriteString(header + "\n\n")

	// Filter input
	filterLabel := t.S().Base.Foreground(t.FgMuted).Render("Filter: ")
	content.WriteString(filterLabel + d.completions.Query() + "_\n\n")

	// Command list
	start := d.selectedIdx / d.visible * d.visible
	end := min(start+d.visible, len(d.filtered))

	for i := start; i < end; i++ {
		if i >= len(d.filtered) {
			break
		}

		cmd := d.filtered[i]
		prefix := " "
		if i == d.selectedIdx {
			prefix = ">"
		}

		itemStyle := t.S().Base
		if i == d.selectedIdx {
			itemStyle = t.S().TextSelected
		}

		line := fmt.Sprintf("%s %s", prefix, cmd.Title)
		if cmd.Description != "" {
			line += fmt.Sprintf(": %s", cmd.Description)
		}

		content.WriteString(itemStyle.Render(line) + "\n")
	}

	// Footer
	footerCount := fmt.Sprintf("%d/%d commands", len(d.filtered), len(d.commands))
	footerHelp := "↑↓: Navigate | Enter: Execute | Type: Filter | ESC: Close"
	footer := t.S().Base.Foreground(t.FgMuted).Render(footerCount + " | " + footerHelp)
	content.WriteString("\n" + footer)

	return boxStyle.Render(content.String())
}

func (d *commandDialogCmp) Position() (int, int) {
	return d.x, d.y
}

func (d *commandDialogCmp) ID() dialogs.DialogID {
	return CommandDialogID
}
