package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/tui/components/core/status"
	"github.com/yourorg/taproot/internal/tui/components/dialogs"
	"github.com/yourorg/taproot/internal/tui/page"
	"github.com/yourorg/taproot/internal/tui/util"
)

// AppModel represents the main application model that manages pages and dialogs.
type AppModel struct {
	width, height int

	currentPage  page.PageID
	previousPage page.PageID
	pages        map[page.PageID]util.Model
	pageStack    []page.PageID

	status     status.StatusCmp
	dialogs    dialogs.DialogCmp
	quitting   bool
	showingFullHelp bool
}

// NewApp creates a new application model.
func NewApp() AppModel {
	return AppModel{
		pages:     make(map[page.PageID]util.Model),
		pageStack: []page.PageID{},
		status:    status.NewStatusCmp(),
		dialogs:   dialogs.NewDialogCmp(),
	}
}

// RegisterPage registers a page with the given ID.
func (a *AppModel) RegisterPage(id page.PageID, model util.Model) {
	a.pages[id] = model
}

// SetPage sets the current page.
func (a *AppModel) SetPage(id page.PageID) tea.Cmd {
	if _, ok := a.pages[id]; !ok {
		return nil
	}
	if a.currentPage != "" {
		a.pageStack = append(a.pageStack, a.currentPage)
	}
	a.currentPage = id
	return a.initPage(id)
}

func (a *AppModel) initPage(id page.PageID) tea.Cmd {
	model, ok := a.pages[id]
	if !ok {
		return nil
	}
	return model.Init()
}

func (a AppModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	if a.currentPage != "" {
		if cmd := a.initPage(a.currentPage); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	cmds = append(cmds, a.status.Init())
	cmds = append(cmds, a.dialogs.Init())
	return tea.Batch(cmds...)
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle dialog messages first
	if a.dialogs.HasDialogs() {
		switch msg.(type) {
		case tea.KeyMsg:
			// Forward to dialog
			_, cmd := a.dialogs.Update(msg)
			return a, cmd
		default:
			// Forward other messages
			_, cmd := a.dialogs.Update(msg)
			return a, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.QuitMsg:
		a.quitting = true
		return a, tea.Quit

	case tea.KeyMsg:
		// Handle global keys
		switch msg.String() {
		case "ctrl+c", "q":
			a.quitting = true
			return a, tea.Quit
		case "ctrl+g":
			a.showingFullHelp = !a.showingFullHelp
			a.status.ToggleFullHelp()
		case "esc":
			// Go back to previous page
			if len(a.pageStack) > 0 {
				lastIdx := len(a.pageStack) - 1
				a.currentPage = a.pageStack[lastIdx]
				a.pageStack = a.pageStack[:lastIdx]
				return a, a.initPage(a.currentPage)
			}
		}

	case page.PageChangeMsg:
		a.SetPage(msg.ID)
		return a, nil

	case page.PageBackMsg:
		if len(a.pageStack) > 0 {
			lastIdx := len(a.pageStack) - 1
			a.currentPage = a.pageStack[lastIdx]
			a.pageStack = a.pageStack[:lastIdx]
			return a, a.initPage(a.currentPage)
		}
		return a, nil

	case dialogs.OpenDialogMsg:
		_, cmd := a.dialogs.Update(msg)
		return a, cmd

	case dialogs.CloseDialogMsg:
		_, cmd := a.dialogs.Update(msg)
		return a, cmd

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		// Forward to status and dialogs
		a.status.Update(msg)
		a.dialogs.Update(msg)

		// Forward to current page
		if currentPage, ok := a.pages[a.currentPage]; ok {
			updated, cmd := currentPage.Update(msg)
			a.pages[a.currentPage] = updated
			return a, cmd
		}
	}

	// Forward to current page
	if currentPage, ok := a.pages[a.currentPage]; ok {
		updated, cmd := currentPage.Update(msg)
		a.pages[a.currentPage] = updated

		// Also forward to status and dialogs
		_, statusCmd := a.status.Update(msg)
		if statusCmd != nil {
			return a, tea.Batch(cmd, statusCmd)
		}

		return a, cmd
	}

	// Forward to status
	_, cmd := a.status.Update(msg)
	return a, cmd
}

func (a AppModel) View() string {
	if a.quitting {
		return ""
	}

	// If there's an active dialog, show it
	if a.dialogs.HasDialogs() {
		dialogView := a.dialogs.View()
		currentPageView := a.currentPageView()

		// Simple overlay - dialog on top of page
		// In a real implementation, you'd use proper layering
		return lipgloss.JoinVertical(lipgloss.Top,
			currentPageView,
			a.status.View(),
			dialogView,
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		a.currentPageView(),
		a.status.View(),
	)
}

func (a AppModel) currentPageView() string {
	if currentPage, ok := a.pages[a.currentPage]; ok {
		return currentPage.View()
	}
	return ""
}

// CurrentPage returns the current page ID.
func (a AppModel) CurrentPage() page.PageID {
	return a.currentPage
}

// HasDialogs returns true if there are active dialogs.
func (a AppModel) HasDialogs() bool {
	return a.dialogs.HasDialogs()
}

// Status returns the status component.
func (a AppModel) Status() status.StatusCmp {
	return a.status
}

// Dialogs returns the dialog component.
func (a AppModel) Dialogs() dialogs.DialogCmp {
	return a.dialogs
}
