package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/tui/components/completions"
	"github.com/wwsheng009/taproot/tui/components/dialogs"
	"github.com/wwsheng009/taproot/tui/util"
	"github.com/wwsheng009/taproot/ui/styles"
)

const (
	CommandDialogID dialogs.DialogID = "commands"
	defaultWidth    int              = 60
)

// ArgDef defines an argument for a command
type ArgDef struct {
	Name        string
	Description string
	Placeholder string
}

// Command represents an executable command
type Command struct {
	ID          string
	Title       string
	Description string
	Args        []ArgDef
	// Callback is called when the command is executed
	// Returns a command to send back to the parent
	Callback func(args map[string]string) tea.Cmd
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
	styles      *styles.Styles
	width       int
	height      int
	x, y        int
	commands    []Command
	completions completions.CompletionsCmp
	filtered    []Command
	selectedIdx int
	visible     int

	// Argument collection state
	collectingArgs bool
	currentArg     int
	argValues      map[string]string
	textInput      textinput.Model
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

	ti := textinput.New()
	ti.Focus()
	ti.Width = defaultWidth - 4

	s := styles.DefaultStyles()

	return &commandDialogCmp{
		styles:      &s,
		width:       defaultWidth,
		height:      20,
		x:           10,
		y:           5,
		commands:    commands,
		filtered:    commands,
		completions: completions.New(),
		visible:     10,
		textInput:   ti,
		argValues:   make(map[string]string),
	}
}

func (d *commandDialogCmp) Init() tea.Cmd {
	return textinput.Blink
}

func (d *commandDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	if d.collectingArgs {
		return d.updateArgs(msg)
	}
	return d.updateList(msg)
}

func (d *commandDialogCmp) updateArgs(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel argument collection, return to list
			d.collectingArgs = false
			d.currentArg = 0
			d.argValues = make(map[string]string)
			d.textInput.Reset()
			return d, nil
		case "enter":
			// Save current argument
			cmd := d.filtered[d.selectedIdx]
			argDef := cmd.Args[d.currentArg]
			d.argValues[argDef.Name] = d.textInput.Value()

			// Move to next argument
			d.currentArg++
			d.textInput.Reset()

			// Check if we have collected all arguments
			if d.currentArg >= len(cmd.Args) {
				// Execute command
				if cmd.Callback != nil {
					return d, tea.Batch(
						func() tea.Msg { return dialogs.CloseDialogMsg{} },
						cmd.Callback(d.argValues),
					)
				}
				return d, func() tea.Msg { return dialogs.CloseDialogMsg{} }
			}

			// Setup prompt for next argument
			return d, nil
		}
	}

	var cmd tea.Cmd
	d.textInput, cmd = d.textInput.Update(msg)
	return d, cmd
}

func (d *commandDialogCmp) updateList(msg tea.Msg) (util.Model, tea.Cmd) {
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
				
				// Check if command needs arguments
				if len(cmd.Args) > 0 {
					d.collectingArgs = true
					d.currentArg = 0
					d.argValues = make(map[string]string)
					d.textInput.Reset()
					return d, nil
				}

				if cmd.Callback != nil {
					// Close dialog and execute command
					return d, tea.Batch(
						func() tea.Msg { return dialogs.CloseDialogMsg{} },
						cmd.Callback(nil),
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
	return len(s) >= len(substr) && (strings.Contains(strings.ToLower(s), strings.ToLower(substr)))
}

func (d *commandDialogCmp) View() string {
	s := d.styles

	// Dialog box style
	boxStyle := s.Base.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Primary).
		Padding(0, 1).
		Width(d.width)

	var content strings.Builder

	// Header
	header := s.Base.Bold(true).Foreground(s.Primary).Render("Commands")
	content.WriteString(header + "\n\n")

	if d.collectingArgs {
		// Input view
		cmd := d.filtered[d.selectedIdx]
		arg := cmd.Args[d.currentArg]

		// Show command title
		content.WriteString(s.Base.Foreground(s.Secondary).Render(cmd.Title) + "\n\n")

		// Show prompt
		label := arg.Name
		if arg.Description != "" {
			label = fmt.Sprintf("%s (%s)", arg.Name, arg.Description)
		}
		content.WriteString(s.Base.Bold(true).Render(label + ":") + "\n")
		
		// Set placeholder
		d.textInput.Placeholder = arg.Placeholder
		content.WriteString(d.textInput.View() + "\n")
		
		// Help
		content.WriteString("\n" + s.Base.Foreground(s.FgMuted).Render("Enter to confirm, Esc to cancel"))

	} else {
		// List view
		// Filter input
		filterLabel := s.Base.Foreground(s.FgMuted).Render("Filter: ")
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

			itemStyle := s.Base
			if i == d.selectedIdx {
				itemStyle = s.TextSelection
			}

			line := fmt.Sprintf("%s %s", prefix, cmd.Title)
			if cmd.Description != "" {
				line += fmt.Sprintf(": %s", cmd.Description)
			}

			content.WriteString(itemStyle.Render(line) + "\n")
		}
		
		// Fill remaining lines to maintain height
		linesRendered := end - start
		if linesRendered < d.visible {
			content.WriteString(strings.Repeat("\n", d.visible-linesRendered))
		}

		// Footer
		footerCount := fmt.Sprintf("%d/%d commands", len(d.filtered), len(d.commands))
		footerHelp := "↑↓: Navigate | Enter: Execute | Type: Filter | ESC: Close"
		footer := s.Base.Foreground(s.FgMuted).Render(footerCount + " | " + footerHelp)
		content.WriteString("\n" + footer)
	}

	return boxStyle.Render(content.String())
}

func (d *commandDialogCmp) Position() (int, int) {
	return d.x, d.y
}

func (d *commandDialogCmp) ID() dialogs.DialogID {
	return CommandDialogID
}
