// Interactive dialog examples using the new engine-agnostic dialog components with Bubbletea
package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/dialog"
	"github.com/wwsheng009/taproot/ui/styles"
)

// DialogsModel demonstrates all dialog types
type DialogsModel struct {
	styles   *styles.Styles
	overlay   *dialog.Overlay
	menuIndex int
	quitting  bool
	width     int
	height    int
}

func NewDialogsModel() *DialogsModel {
	s := styles.DefaultStyles()
	return &DialogsModel{
		styles:   &s,
		overlay:   dialog.NewOverlay(),
		menuIndex: 0,
		quitting:  false,
	}
}

// Init implements tea.Model
func (m *DialogsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *DialogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If a dialog is active, forward messages to it
	if m.overlay.HasDialogs() {
		return m.updateDialog(msg)
	}

	// Otherwise handle menu navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.menuIndex > 0 {
				m.menuIndex--
			}

		case "down", "j":
			if m.menuIndex < 4 {
				m.menuIndex++
			}

		case "enter", " ":
			return m, m.openSelectedDialog()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.overlay.SetSize(msg.Width, msg.Height)
	}

	return m, nil
}

func (m *DialogsModel) updateDialog(msg tea.Msg) (tea.Model, tea.Cmd) {
	activeDialog := m.overlay.ActiveDialog()

	// Convert tea.KeyMsg to render.KeyMsg for dialog compatibility
	var renderMsg any
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		renderMsg = dialog.ConvertKeyMsg(keyMsg)
	} else {
		renderMsg = msg
	}

	// Update the active dialog
	updated, _ := activeDialog.Update(renderMsg)
	updatedDialog, ok := updated.(dialog.Dialog)
	if !ok {
		// Dialog was dismissed, remove from overlay
		m.overlay.Pop()
		return m, nil
	}

	// Check if dialog view is empty (indicates quitting state)
	if updatedDialog.View() == "" {
		m.overlay.Pop()
		return m, nil
	}

	// Replace the dialog in the overlay with updated state
	m.overlay.Pop()
	m.overlay.Push(updatedDialog)

	return m, nil
}

func (m *DialogsModel) openSelectedDialog() tea.Cmd {
	switch m.menuIndex {
	case 0:
		// Info Dialog
		d := dialog.NewInfoDialog(
			"Information",
			"This is an informational dialog message.\n\nIt can display important details to the user.",
		)
		d.Init()
		m.overlay.Push(d)

	case 1:
		// Confirm Dialog
		d := dialog.NewConfirmDialog(
			"Confirm Action",
			"Are you sure you want to proceed with this action?",
			func(result dialog.ActionResult, data any) {
				// Handle result
			},
		)
		d.Init()
		m.overlay.Push(d)

	case 2:
		// Input Dialog
		d := dialog.NewInputDialog(
			"Enter Your Name",
			"Name",
			func(value string) {
				// Handle input
			},
		)
		d.SetPlaceholder("John Doe")
		d.Init()
		m.overlay.Push(d)

	case 3:
		// Select List Dialog
		items := []string{
			"Apple - A red fruit",
			"Banana - A yellow fruit",
			"Cherry - A small red fruit",
			"Date - A sweet brown fruit",
			"Elderberry - A dark purple fruit",
		}
		d := dialog.NewSelectListDialog("Choose a Fruit", items, func(index int, value string) {
			// Handle selection
		})
		d.Init()
		m.overlay.Push(d)

	case 4:
		// Quit
		m.quitting = true
		return tea.Quit
	}

	return nil
}

// View implements tea.Model
func (m *DialogsModel) View() string {
	if m.quitting {
		return ""
	}

	// Render the background
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("235")).
		Padding(0, 2).
		Bold(true)

	b.WriteString(titleStyle.Render("╔════════════════════════════════════════════════════════════════╗"))
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("║              Dialog Examples (v2.0)                            ║"))
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("╠════════════════════════════════════════════════════════════════╣"))
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("║  Select a dialog type and press Enter to open it                    ║"))
	b.WriteString("\n")
	b.WriteString(titleStyle.Render("╚════════════════════════════════════════════════════════════════╝"))
	b.WriteString("\n\n")

	// Menu items
	menuItems := []struct {
		icon  string
		label string
		desc  string
	}{
		{"ℹ️ ", "Info Dialog", "Displays an informational message"},
		{"✓", "Confirm Dialog", "Asks for user confirmation"},
		{"✏️ ", "Input Dialog", "Prompts for text input"},
		{"📋", "Select List Dialog", "Select from a list of items"},
		{"", "Quit", "Exit the application"},
	}

	for i, item := range menuItems {
		cursor := " "
		if i == m.menuIndex {
			cursor = "→"
		}

		labelStyle := lipgloss.NewStyle()
		if i == m.menuIndex {
			labelStyle = labelStyle.Foreground(lipgloss.Color("226")).Bold(true)
		} else {
			labelStyle = labelStyle.Foreground(lipgloss.Color("252"))
		}

		line := fmt.Sprintf("%s %s %-20s", cursor, item.icon, item.label)
		if item.desc != "" {
			line += fmt.Sprintf(" — %s", item.desc)
		}

		b.WriteString(labelStyle.Render("  "+line) + "\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	footerText := "↑/k/j down: Navigate | Enter/Space: Open | q: Quit"
	b.WriteString("\n")
	b.WriteString(footerStyle.Render(footerText))

	backgroundView := b.String()

	if m.width > 0 && m.height > 0 {
		backgroundView = lipgloss.Place(
			m.width, m.height,
			lipgloss.Left, lipgloss.Top,
			backgroundView,
		)
	}

	// If dialog is active, overlay it on top of background
	if m.overlay.HasDialogs() {
		return m.overlay.Render(backgroundView)
	}

	return backgroundView
}

func main() {
	p := tea.NewProgram(
		NewDialogsModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
