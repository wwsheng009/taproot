package dialogs

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourorg/taproot/internal/tui/util"
)

// mockDialogModel implements DialogModel for testing
type mockDialogModel struct {
	id      DialogID
	initCnt int
	updateCnt int
	viewCnt  int
}

func (m *mockDialogModel) Init() tea.Cmd {
	m.initCnt++
	return nil
}

func (m *mockDialogModel) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	m.updateCnt++
	return m, nil
}

func (m *mockDialogModel) View() string {
	m.viewCnt++
	return "mock dialog"
}

func (m *mockDialogModel) Position() (int, int) {
	return 0, 0
}

func (m *mockDialogModel) ID() DialogID {
	return m.id
}

func TestOpenDialogMsg(t *testing.T) {
	t.Run("OpenDialogMsg contains model", func(t *testing.T) {
		mock := &mockDialogModel{id: "test"}
		msg := OpenDialogMsg{Model: mock}

		if msg.Model == nil {
			t.Error("OpenDialogMsg.Model should not be nil")
		}
		if msg.Model.ID() != "test" {
			t.Errorf("OpenDialogMsg.Model.ID() = %s, want 'test'", msg.Model.ID())
		}
	})
}

func TestCloseDialogMsg(t *testing.T) {
	t.Run("CloseDialogMsg is a valid message type", func(t *testing.T) {
		msg := CloseDialogMsg{}
		// Just verify it can be created and passed as tea.Msg
		var _ tea.Msg = msg
	})
}

func TestDialogID(t *testing.T) {
	t.Run("DialogID can be created and compared", func(t *testing.T) {
		id1 := DialogID("dialog1")
		id2 := DialogID("dialog1")
		id3 := DialogID("dialog2")

		if id1 != id2 {
			t.Error("DialogID 'dialog1' should equal 'dialog1'")
		}
		if id1 == id3 {
			t.Error("DialogID 'dialog1' should not equal 'dialog2'")
		}
	})
}

func TestDialogModelInterface(t *testing.T) {
	t.Run("DialogModel interface is satisfied", func(t *testing.T) {
		var _ DialogModel = &mockDialogModel{}
		mock := &mockDialogModel{id: "test"}

		// Test Init
		mock.Init()
		if mock.initCnt != 1 {
			t.Errorf("Init() called count = %d, want 1", mock.initCnt)
		}

		// Test Update
		mock.Update(tea.KeyMsg{})
		if mock.updateCnt != 1 {
			t.Errorf("Update() called count = %d, want 1", mock.updateCnt)
		}

		// Test View
		view := mock.View()
		if mock.viewCnt != 1 {
			t.Errorf("View() called count = %d, want 1", mock.viewCnt)
		}
		if view != "mock dialog" {
			t.Errorf("View() = %s, want 'mock dialog'", view)
		}

		// Test ID
		if mock.ID() != "test" {
			t.Errorf("ID() = %s, want 'test'", mock.ID())
		}

		// Test Position
		x, y := mock.Position()
		if x != 0 || y != 0 {
			t.Errorf("Position() = (%d, %d), want (0, 0)", x, y)
		}
	})
}
