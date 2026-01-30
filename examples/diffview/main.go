package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/internal/tui/exp/diffview"
)

// Model holds the application state
type model struct {
	diffView *diffview.DiffView
	width    int
	height   int
}

// InitialModel creates the initial model
func initialModel() model {
	before := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	
	x := 42
	y := 100
	
	sum := x + y
	fmt.Printf("The sum is: %d\\n", sum)
	
	items := []string{"apple", "banana", "cherry"}
	for _, item := range items {
		fmt.Println(item)
	}
	
	return nil
}`

	after := `package main

import "fmt"

func main() {
	fmt.Println("Hello, Taproot!")
	
	x := 42
	y := 200  // Changed from 100
	
	sum := x + y
	product := x * y  // New line
	fmt.Printf("The sum is: %d\\n", sum)
	fmt.Printf("The product is: %d\\n", product)
	
	items := []string{"apple", "blueberry", "cherry", "date"}
	for _, item := range items {
		fmt.Println(item)
	}
	
	// New function
	greet := func(name string) {
		fmt.Printf("Greetings, %s!\\n", name)
	}
	greet("User")
	
	return nil
}`

	dv := diffview.New()
	dv.Before(before)
	dv.After(after)
	dv.SetLayout(diffview.LayoutUnified)
	dv.SetLineNumbers(true)

	return model{
		diffView: dv,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.diffView.SetSize(msg.Width, msg.Height-3) // Leave space for header/footer

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.diffView.ScrollUp()
		case "down", "j":
			m.diffView.ScrollDown()
		case "left", "h":
			m.diffView.ScrollLeft()
		case "right", "l":
			m.diffView.ScrollRight()
		case "g", "home":
			m.diffView.ScrollToTop()
		case "G", "end":
			m.diffView.ScrollToBottom()
		case "ctrl+u":
			for i := 0; i < 5; i++ {
				m.diffView.ScrollUp()
			}
		case "ctrl+d":
			for i := 0; i < 5; i++ {
				m.diffView.ScrollDown()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	return m.diffView.Render()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
