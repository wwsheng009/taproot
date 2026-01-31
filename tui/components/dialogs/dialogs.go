package dialogs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/tui/util"
)

type DialogID string

// DialogModel represents a dialog component that can be displayed.
type DialogModel interface {
	util.Model
	Position() (int, int)
	ID() DialogID
}

// CloseCallback allows dialogs to perform cleanup when closed.
type CloseCallback interface {
	Close() tea.Cmd
}

// OpenDialogMsg is sent to open a new dialog with specified dimensions.
type OpenDialogMsg struct {
	Model DialogModel
}

// CloseDialogMsg is sent to close the topmost dialog.
type CloseDialogMsg struct{}

// DialogCmp manages a stack of dialogs with keyboard navigation.
type DialogCmp interface {
	util.Model

	Dialogs() []DialogModel
	HasDialogs() bool
	ActiveModel() util.Model
	ActiveDialogID() DialogID
}

type dialogCmp struct {
	width, height int
	dialogs       []DialogModel
	idMap         map[DialogID]int
}

// NewDialogCmp creates a new dialog manager.
func NewDialogCmp() DialogCmp {
	return &dialogCmp{
		dialogs: []DialogModel{},
		idMap:   make(map[DialogID]int),
	}
}

func (d dialogCmp) Init() tea.Cmd {
	return nil
}

// Update handles dialog lifecycle and forwards messages to the active dialog.
func (d *dialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	newDialogs := *d  // Deep copy

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmds []tea.Cmd
		newDialogs.width = msg.Width
		newDialogs.height = msg.Height
		for i := range newDialogs.dialogs {
			u, cmd := newDialogs.dialogs[i].Update(msg)
			newDialogs.dialogs[i] = u.(DialogModel)
			cmds = append(cmds, cmd)
		}
		return &newDialogs, tea.Batch(cmds...)
	case OpenDialogMsg:
		return newDialogs.handleOpen(msg)
	case CloseDialogMsg:
		if len(newDialogs.dialogs) == 0 {
			return &newDialogs, nil
		}
		// Call Close callback if available
		if cc, ok := newDialogs.dialogs[len(newDialogs.dialogs)-1].(CloseCallback); ok {
			cmd := cc.Close()
			// Remove the dialog
			newDialogs.pop()
			return &newDialogs, cmd
		}
		newDialogs.pop()
		return &newDialogs, nil
	}

	// Forward messages to the active dialog
	if len(newDialogs.dialogs) > 0 {
		activeIdx := len(newDialogs.dialogs) - 1
		updated, cmd := newDialogs.dialogs[activeIdx].Update(msg)
		newDialogs.dialogs[activeIdx] = updated.(DialogModel)
		return &newDialogs, cmd
	}

	return &newDialogs, nil
}

func (d *dialogCmp) handleOpen(msg OpenDialogMsg) (util.Model, tea.Cmd) {
	dialog := msg.Model
	newDialogs := *d  // Deep copy
	newDialogs.dialogs = append(newDialogs.dialogs, dialog)
	newDialogs.idMap[dialog.ID()] = len(newDialogs.dialogs) - 1
	return &newDialogs, nil
}

func (d *dialogCmp) pop() {
	if len(d.dialogs) > 0 {
		id := d.dialogs[len(d.dialogs)-1].ID()
		delete(d.idMap, id)
		d.dialogs = d.dialogs[:len(d.dialogs)-1]
	}
}

func (d dialogCmp) View() string {
	if len(d.dialogs) == 0 {
		return ""
	}

	// Render the active dialog only
	activeDialog := d.dialogs[len(d.dialogs)-1]
	return activeDialog.View()
}

func (d dialogCmp) Dialogs() []DialogModel {
	return d.dialogs
}

func (d dialogCmp) HasDialogs() bool {
	return len(d.dialogs) > 0
}

func (d dialogCmp) ActiveModel() util.Model {
	if len(d.dialogs) == 0 {
		return nil
	}
	return d.dialogs[len(d.dialogs)-1]
}

func (d dialogCmp) ActiveDialogID() DialogID {
	if len(d.dialogs) == 0 {
		return ""
	}
	return d.dialogs[len(d.dialogs)-1].ID()
}
