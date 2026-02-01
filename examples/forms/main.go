package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/forms"
	"github.com/wwsheng009/taproot/ui/render"
)

type model struct {
	form *forms.Form
}

func initialModel() model {
	// 1. Email
	email := forms.NewTextInput("Enter email")
	email.SetShowBorder(true)
	email.AddValidator(forms.Email)
	email.AddValidator(forms.Required)

	// 2. Password
	password := forms.NewTextInput("Enter password")
	password.SetShowBorder(true)
	password.SetHidden(true)
	password.AddValidator(forms.Required)
	password.AddValidator(forms.MinLength(8))

	// 3. Role (Select)
	role := forms.NewSelect("Role", []string{"Developer", "Designer", "Product Manager", "DevOps"})
	role.SetPlaceholder("Select a role")
	role.AddValidator(forms.Required)

	// 4. Age
	age := forms.NewNumberInput("Enter age")
	age.SetShowBorder(true)
	age.SetRange(0, 150)
	age.SetStep(1)

	// 5. Bio
	bio := forms.NewTextArea("Enter bio")
	bio.SetShowBorder(true)
	bio.AddValidator(forms.Required)
	bio.AddValidator(forms.MinLength(10))

	// 6. Notification preferences
	notify := forms.NewRadioGroup("Notifications:", []string{"Email", "SMS", "Push", "None"})
	notify.AddValidator(forms.Required)

	// 7. Terms
	terms := forms.NewCheckbox("I accept the Terms and Conditions")

	form := forms.NewForm(email, password, role, age, bio, notify, terms)

	return model{
		form: form,
	}
}

func (m model) Init() tea.Cmd {
	return adaptCmd(m.form.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// Optional: Handle global validation trigger?
		// For now, let's rely on individual field validation state
	}

	newForm, cmd := m.form.Update(msg)
	m.form = newForm.(*forms.Form)

	return m, adaptCmd(cmd)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("Form Example (Tab/Enter to navigate, Ctrl+C to quit)\n\n")
	b.WriteString(m.form.View())
	b.WriteString("\n\n")

	// Show validation status summary?
	if err := m.form.Validate(); err == nil {
		b.WriteString("Status: Valid\n")
	} else {
		fmt.Fprintf(&b, "Status: Invalid (%s)\n", err.Error())
	}

	return b.String()
}

// adaptCmd converts a render.Cmd (engine-agnostic) to a tea.Cmd (Bubbletea specific).
// It recursively handles Batches and wrapping of function types.
func adaptCmd(cmd render.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}

	switch c := cmd.(type) {
	case render.BatchCmd:
		var cmds []tea.Cmd
		for _, nested := range c {
			if ac := adaptCmd(nested); ac != nil {
				cmds = append(cmds, ac)
			}
		}
		return tea.Batch(cmds...)
	
	case []render.Cmd: // Just in case it's a raw slice
		var cmds []tea.Cmd
		for _, nested := range c {
			if ac := adaptCmd(nested); ac != nil {
				cmds = append(cmds, ac)
			}
		}
		return tea.Batch(cmds...)

	case func() render.Msg:
		return func() tea.Msg {
			return c()
		}

	case func() error:
		return func() tea.Msg {
			if err := c(); err != nil {
				return render.ErrorMsg{Error: err}
			}
			return nil
		}
	
	case render.Command: // Alias for func() error
		return func() tea.Msg {
			if err := c(); err != nil {
				return render.ErrorMsg{Error: err}
			}
			return nil
		}

	case render.Cmd:
		// Check for Quit command interface
		if render.IsQuit(c) {
			return tea.Quit
		}
		// If it's a type we don't know but implements the empty interface (everything does),
		// check if it's a batch implicitly? No, Go types are strict.
	}

	return nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
