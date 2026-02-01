package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/forms"
	"github.com/wwsheng009/taproot/ui/render"
)

type model struct {
	email    *forms.TextInput
	password *forms.TextInput
	age      *forms.NumberInput
	bio      *forms.TextArea
	focus    int // 0: email, 1: password, 2: age, 3: bio
}

func initialModel() model {
	email := forms.NewTextInput("Enter email")
	email.SetShowBorder(true)
	email.AddValidator(forms.Email)
	email.AddValidator(forms.Required)
	email.Focus()

	password := forms.NewTextInput("Enter password")
	password.SetShowBorder(true)
	password.SetHidden(true)
	password.AddValidator(forms.Required)
	password.AddValidator(forms.MinLength(8))

	age := forms.NewNumberInput("Enter age")
	age.SetShowBorder(true)
	age.SetRange(0, 150)
	age.SetStep(1)

	bio := forms.NewTextArea("Enter bio")
	bio.SetShowBorder(true)
	bio.AddValidator(forms.Required)
	bio.AddValidator(forms.MinLength(10))

	return model{
		email:    email,
		password: password,
		age:      age,
		bio:      bio,
		focus:    0,
	}
}

func (m model) Init() tea.Cmd {
	return m.updateFocus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focus = (m.focus + 1) % 4
			return m, m.updateFocus()
		case "shift+tab":
			m.focus = (m.focus - 1 + 4) % 4
			return m, m.updateFocus()
		case "enter":
			// Validate all
			m.email.Validate()
			m.password.Validate()
			m.age.Validate()
			m.bio.Validate()
		}
	}

	// Update components
	// Email
	newEmail, cmd := m.email.Update(msg)
	m.email = newEmail.(*forms.TextInput)
	if cmd != nil {
		cmds = append(cmds, adaptCmd(cmd))
	}

	// Password
	newPass, cmd := m.password.Update(msg)
	m.password = newPass.(*forms.TextInput)
	if cmd != nil {
		cmds = append(cmds, adaptCmd(cmd))
	}

	// Age
	newAge, cmd := m.age.Update(msg)
	m.age = newAge.(*forms.NumberInput)
	if cmd != nil {
		cmds = append(cmds, adaptCmd(cmd))
	}

	// Bio
	newBio, cmd := m.bio.Update(msg)
	m.bio = newBio.(*forms.TextArea)
	if cmd != nil {
		cmds = append(cmds, adaptCmd(cmd))
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateFocus() tea.Cmd {
	focuses := []func() render.Cmd{
		m.email.Focus,
		m.password.Focus,
		m.age.Focus,
		m.bio.Focus,
	}
	blurs := []func(){
		m.email.Blur,
		m.password.Blur,
		m.age.Blur,
		m.bio.Blur,
	}

	var cmd tea.Cmd
	for i := range 4 {
		if i == m.focus {
			if c := focuses[i](); c != nil {
				cmd = adaptCmd(c)
			}
		} else {
			blurs[i]()
		}
	}
	return cmd
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("Form Example (Tab/Shift+Tab to switch, Enter to validate, Ctrl+C to quit)\n\n")

	b.WriteString("Email:\n")
	b.WriteString(m.email.View())
	b.WriteString("\n\n")

	b.WriteString("Password:\n")
	b.WriteString(m.password.View())
	b.WriteString("\n\n")

	b.WriteString("Age:\n")
	b.WriteString(m.age.View())
	b.WriteString("\n\n")

	b.WriteString("Bio:\n")
	b.WriteString(m.bio.View())
	b.WriteString("\n\n")

	return b.String()
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
