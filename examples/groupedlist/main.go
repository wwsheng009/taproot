package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/internal/tui/exp/list"
)

// Model holds the application state
type model struct {
	list *list.GroupedList
}

// InitialModel creates the initial model
func initialModel() model {
	groups := []list.Group{
		{
			Title: "Programming Languages",
			Expanded: true,
			Items: []list.ListItem{
				{ID: "go", Title: "Go", Desc: "Fast, simple language"},
				{ID: "rust", Title: "Rust", Desc: "Safe systems language"},
				{ID: "python", Title: "Python", Desc: "Easy to learn"},
			},
		},
		{
			Title: "Web Technologies",
			Expanded: true,
			Items: []list.ListItem{
				{ID: "js", Title: "JavaScript", Desc: "Browser language"},
				{ID: "ts", Title: "TypeScript", Desc: "Typed JS"},
				{ID: "html", Title: "HTML", Desc: "Structure"},
				{ID: "css", Title: "CSS", Desc: "Styling"},
			},
		},
		{
			Title: "Databases",
			Expanded: true,
			Items: []list.ListItem{
				{ID: "postgres", Title: "PostgreSQL", Desc: "Relational DB"},
				{ID: "mysql", Title: "MySQL", Desc: "Popular relational DB"},
				{ID: "mongo", Title: "MongoDB", Desc: "NoSQL document DB"},
			},
		},
		{
			Title: "Tools",
			Expanded: true,
			Items: []list.ListItem{
				{ID: "git", Title: "Git", Desc: "Version control"},
				{ID: "docker", Title: "Docker", Desc: "Containers"},
				{ID: "vim", Title: "Vim", Desc: "Text editor"},
			},
		},
	}

	return model{
		list: list.NewGroupedList(groups),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList.(*list.GroupedList)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
