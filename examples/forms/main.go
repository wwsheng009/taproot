package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/forms"
	"github.com/wwsheng009/taproot/ui/render"
)

type model struct {
	textInput   *forms.TextInput
	numberInput *forms.NumberInput
	textArea    *forms.TextArea
	focusIndex  int
}

func initialModel() model {
	ti := forms.NewTextInput("Enter name")
	ti.Focus()

	ni := forms.NewNumberInput("Enter age")
	ni.SetRange(0, 150)

	ta := forms.NewTextArea("Enter bio")
	ta.SetHeight(5)

	return model{
		textInput:   ti,
		numberInput: ni,
		textArea:    ta,
		focusIndex:  0,
	}
}

func (m model) Init() tea.Cmd {
	return cmdToTea(m.textInput.Focus())
}

func cmdToTea(cmd render.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	if c, ok := cmd.(func() render.Msg); ok {
		return func() tea.Msg {
			return c()
		}
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd render.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			m.focusNext()
			return m, nil
		}
	}

	// Update focused component
	switch m.focusIndex {
	case 0:
		_, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmdToTea(cmd))
	case 1:
		_, cmd = m.numberInput.Update(msg)
		cmds = append(cmds, cmdToTea(cmd))
	case 2:
		_, cmd = m.textArea.Update(msg)
		cmds = append(cmds, cmdToTea(cmd))
	}

	return m, tea.Batch(cmds...)
}

func (m *model) focusNext() {
	m.textInput.Blur()
	m.numberInput.Blur()
	m.textArea.Blur()

	m.focusIndex = (m.focusIndex + 1) % 3

	switch m.focusIndex {
	case 0:
		m.textInput.Focus()
	case 1:
		m.numberInput.Focus()
	case 2:
		m.textArea.Focus()
	}
}

func (m model) View() string {
	return fmt.Sprintf(
		"Forms Example\n\nName:\n%s\n\nAge:\n%s\n\nBio:\n%s\n\n(Tab to switch focus, Esc to quit)",
		m.textInput.View(),
		m.numberInput.View(),
		m.textArea.View(),
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
