package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/tui/exp/list"
)

// Model holds the application state
type model struct {
	list *list.FilterableList
}

// InitialModel creates the initial model
func initialModel() model {
	items := []list.ListItem{
		{ID: "1", Title: "Go", Desc: "Programming language"},
		{ID: "2", Title: "Python", Desc: "Programming language"},
		{ID: "3", Title: "JavaScript", Desc: "Web programming"},
		{ID: "4", Title: "Rust", Desc: "Systems programming"},
		{ID: "5", Title: "TypeScript", Desc: "Typed JavaScript"},
		{ID: "6", Title: "Java", Desc: "Enterprise language"},
		{ID: "7", Title: "C++", Desc: "Systems language"},
		{ID: "8", Title: "Ruby", Desc: "Web framework language"},
		{ID: "9", Title: "PHP", Desc: "Server-side scripting"},
		{ID: "10", Title: "Swift", Desc: "Apple platform language"},
		{ID: "11", Title: "Kotlin", Desc: "Android development"},
		{ID: "12", Title: "Dart", Desc: "Flutter language"},
		{ID: "13", Title: "Haskell", Desc: "Functional language"},
		{ID: "14", Title: "Elixir", Desc: "Functional concurrent"},
		{ID: "15", Title: "Lua", Desc: "Embedded scripting"},
	}

	return model{
		list: list.NewFilterableList(items),
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
	m.list = newList.(*list.FilterableList)
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
