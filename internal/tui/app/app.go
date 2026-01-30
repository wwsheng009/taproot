package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/tui/components/core/status"
	"github.com/wwsheng009/taproot/internal/tui/components/dialogs"
	"github.com/wwsheng009/taproot/internal/tui/lifecycle"
	"github.com/wwsheng009/taproot/internal/tui/page"
	"github.com/wwsheng009/taproot/internal/tui/util"
)

// AppModel represents the main application model that manages pages and dialogs.
type AppModel struct {
	width, height int

	currentPage  page.PageID
	previousPage page.PageID
	pages        map[page.PageID]util.Model
	pageStack    []page.PageID

	status        status.StatusCmp
	dialogs       dialogs.DialogCmp
	lifecycleMgr  *lifecycle.LifecycleManager
	quitting      bool
	showingFullHelp bool
}

// NewApp creates a new application model.
func NewApp() AppModel {
	return AppModel{
		pages:        make(map[page.PageID]util.Model),
		pageStack:    []page.PageID{},
		status:       status.NewStatusCmp(),
		dialogs:      dialogs.NewDialogCmp(),
		lifecycleMgr: lifecycle.NewLifecycleManager(),
	}
}

// RegisterPage registers a page with the given ID.
func (a *AppModel) RegisterPage(id page.PageID, model util.Model) {
	a.pages[id] = model
}

// SetPage sets the current page and returns new model.
func (a *AppModel) SetPage(id page.PageID) *AppModel {
	if _, ok := a.pages[id]; !ok {
		return a
	}
	newApp := *a  // Deep copy
	if newApp.currentPage != "" {
		newApp.pageStack = append(newApp.pageStack, newApp.currentPage)
	}
	newApp.currentPage = id
	return &newApp
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
	newApp := a  // Create copy

	// Handle dialog messages first
	if newApp.dialogs.HasDialogs() {
		switch msg.(type) {
		case tea.KeyMsg:
			// Forward to dialog
			updatedDialogs, cmd := newApp.dialogs.Update(msg)
			newApp.dialogs = updatedDialogs.(dialogs.DialogCmp)
			return newApp, cmd
		default:
			// Forward other messages
			updatedDialogs, cmd := newApp.dialogs.Update(msg)
			newApp.dialogs = updatedDialogs.(dialogs.DialogCmp)
			return newApp, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.QuitMsg:
		newApp.quitting = true
		return newApp, tea.Quit

	case tea.KeyMsg:
		// Handle global keys
		switch msg.String() {
		case "ctrl+c", "q":
			newApp.quitting = true
			return newApp, tea.Quit
		case "ctrl+g":
			newApp.showingFullHelp = !newApp.showingFullHelp
			newApp.status.ToggleFullHelp()
		case "ctrl+m":
			// Forward to current page to handle
			if currentPage, ok := newApp.pages[newApp.currentPage]; ok {
				updated, cmd := currentPage.Update(msg)
				newApp.pages[newApp.currentPage] = updated
				return newApp, cmd
			}
		case "ctrl+p":
			// Forward to current page to handle
			if currentPage, ok := newApp.pages[newApp.currentPage]; ok {
				updated, cmd := currentPage.Update(msg)
				newApp.pages[newApp.currentPage] = updated
				return newApp, cmd
			}
		case "esc":
			// Go back to previous page
			if len(newApp.pageStack) > 0 {
				lastIdx := len(newApp.pageStack) - 1
				newApp.currentPage = newApp.pageStack[lastIdx]
				newApp.pageStack = newApp.pageStack[:lastIdx]
				return newApp, newApp.initPage(newApp.currentPage)
			}
		}

	case page.PageChangeMsg:
		if _, ok := newApp.pages[msg.ID]; ok {
			// Cancel old page lifecycle if exists
			if newApp.currentPage != "" {
				newApp.lifecycleMgr.CancelContext(string(newApp.currentPage))
				newApp.pageStack = append(newApp.pageStack, newApp.currentPage)
			}
			// Register new page lifecycle
			newApp.currentPage = msg.ID
			newApp.lifecycleMgr.Register(string(msg.ID))
			cmd := newApp.initPage(msg.ID)
			return newApp, cmd
		}
		return newApp, nil

	case page.PageBackMsg:
		if len(newApp.pageStack) > 0 {
			// Cancel current page lifecycle
			if newApp.currentPage != "" {
				newApp.lifecycleMgr.CancelContext(string(newApp.currentPage))
			}
			lastIdx := len(newApp.pageStack) - 1
			newApp.currentPage = newApp.pageStack[lastIdx]
			newApp.pageStack = newApp.pageStack[:lastIdx]
			// Register new page lifecycle
			newApp.lifecycleMgr.Register(string(newApp.currentPage))
			return newApp, newApp.initPage(newApp.currentPage)
		}
		return newApp, nil

	case dialogs.OpenDialogMsg:
		updatedDialogs, cmd := newApp.dialogs.Update(msg)
		newApp.dialogs = updatedDialogs.(dialogs.DialogCmp)
		return newApp, cmd

	case dialogs.CloseDialogMsg:
		updatedDialogs, cmd := newApp.dialogs.Update(msg)
		newApp.dialogs = updatedDialogs.(dialogs.DialogCmp)
		return newApp, cmd

	case tea.WindowSizeMsg:
		newApp.width = msg.Width
		newApp.height = msg.Height
		// Forward to status and dialogs
		newApp.status.Update(msg)
		newApp.dialogs.Update(msg)

		// Forward to current page
		if currentPage, ok := newApp.pages[newApp.currentPage]; ok {
			updated, cmd := currentPage.Update(msg)
			newApp.pages[newApp.currentPage] = updated
			return newApp, cmd
		}
	}

	// Forward to current page
	if currentPage, ok := newApp.pages[newApp.currentPage]; ok {
		updated, cmd := currentPage.Update(msg)
		newApp.pages[newApp.currentPage] = updated

		// Also forward to status and dialogs
		_, statusCmd := newApp.status.Update(msg)
		if statusCmd != nil {
			return newApp, tea.Batch(cmd, statusCmd)
		}

		return newApp, cmd
	}

	// Forward to status
	_, cmd := newApp.status.Update(msg)
	return newApp, cmd
}

func (a AppModel) View() string {
	if a.quitting {
		return ""
	}

	// If there's an active dialog, show it
	if a.dialogs.HasDialogs() {
		dialogView := a.dialogs.View()
		
		// Center the dialog using lipgloss.Place
		// This will replace the current view with the centered dialog
		return lipgloss.Place(a.width, a.height,
			lipgloss.Center, lipgloss.Center,
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

// GetPageContext returns the lifecycle context for a page.
func (a AppModel) GetPageContext(pageID page.PageID) (context.Context, bool) {
	return a.lifecycleMgr.GetContext(string(pageID))
}

// CancelPageContext cancels the lifecycle context for a page.
func (a AppModel) CancelPageContext(pageID page.PageID) {
	a.lifecycleMgr.CancelContext(string(pageID))
}

