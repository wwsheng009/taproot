package util

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInfoType(t *testing.T) {
	t.Run("InfoType values are correct", func(t *testing.T) {
		tests := []struct {
			name  string
			value InfoType
			want  int
		}{
			{"InfoTypeInfo", InfoTypeInfo, 0},
			{"InfoTypeSuccess", InfoTypeSuccess, 1},
			{"InfoTypeWarn", InfoTypeWarn, 2},
			{"InfoTypeError", InfoTypeError, 3},
			{"InfoTypeUpdate", InfoTypeUpdate, 4},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if int(tt.value) != tt.want {
					t.Errorf("%s = %d, want %d", tt.name, tt.value, tt.want)
				}
			})
		}
	})
}

func TestNewInfoMsg(t *testing.T) {
	t.Run("NewInfoMsg creates InfoMsg with correct type", func(t *testing.T) {
		msg := NewInfoMsg("test message")

		if msg.Type != InfoTypeInfo {
			t.Errorf("NewInfoMsg() Type = %d, want InfoTypeInfo (%d)", msg.Type, InfoTypeInfo)
		}
		if msg.Msg != "test message" {
			t.Errorf("NewInfoMsg() Msg = %s, want 'test message'", msg.Msg)
		}
	})
}

func TestNewSuccessMsg(t *testing.T) {
	t.Run("NewSuccessMsg creates InfoMsg with success type", func(t *testing.T) {
		msg := NewSuccessMsg("success")

		if msg.Type != InfoTypeSuccess {
			t.Errorf("NewSuccessMsg() Type = %d, want InfoTypeSuccess (%d)", msg.Type, InfoTypeSuccess)
		}
		if msg.Msg != "success" {
			t.Errorf("NewSuccessMsg() Msg = %s, want 'success'", msg.Msg)
		}
	})
}

func TestNewWarnMsg(t *testing.T) {
	t.Run("NewWarnMsg creates InfoMsg with warn type", func(t *testing.T) {
		msg := NewWarnMsg("warning")

		if msg.Type != InfoTypeWarn {
			t.Errorf("NewWarnMsg() Type = %d, want InfoTypeWarn (%d)", msg.Type, InfoTypeWarn)
		}
		if msg.Msg != "warning" {
			t.Errorf("NewWarnMsg() Msg = %s, want 'warning'", msg.Msg)
		}
	})
}

func TestNewErrorMsg(t *testing.T) {
	t.Run("NewErrorMsg creates InfoMsg with error type", func(t *testing.T) {
		msg := NewErrorMsg(errors.New("error"))

		if msg.Type != InfoTypeError {
			t.Errorf("NewErrorMsg() Type = %d, want InfoTypeError (%d)", msg.Type, InfoTypeError)
		}
		if msg.Msg != "error" {
			t.Errorf("NewErrorMsg() Msg = %s, want 'error'", msg.Msg)
		}
	})
}

func TestReportError(t *testing.T) {
	t.Run("ReportError returns a command", func(t *testing.T) {
		err := errors.New("test error")
		cmd := ReportError(err)

		if cmd == nil {
			t.Fatal("ReportError() should return a non-nil command")
		}

		// Execute the command to get the message
		msg := cmd()
		if msg == nil {
			t.Fatal("ReportError() command should return a message")
		}

		errorMsg, ok := msg.(InfoMsg)
		if !ok {
			t.Fatal("ReportError() command should return InfoMsg")
		}

		if errorMsg.Type != InfoTypeError {
			t.Errorf("ReportError() Type = %d, want InfoTypeError (%d)", errorMsg.Type, InfoTypeError)
		}
	})
}

func TestReportInfo(t *testing.T) {
	t.Run("ReportInfo returns a command with info type", func(t *testing.T) {
		cmd := ReportInfo("info message")

		msg := cmd()
		infoMsg, ok := msg.(InfoMsg)
		if !ok {
			t.Fatal("ReportInfo() command should return InfoMsg")
		}

		if infoMsg.Type != InfoTypeInfo {
			t.Errorf("ReportInfo() Type = %d, want InfoTypeInfo (%d)", infoMsg.Type, InfoTypeInfo)
		}
		if infoMsg.Msg != "info message" {
			t.Errorf("ReportInfo() Msg = %s, want 'info message'", infoMsg.Msg)
		}
	})
}

func TestReportSuccess(t *testing.T) {
	t.Run("ReportSuccess returns a command with success type", func(t *testing.T) {
		cmd := ReportSuccess("success message")

		msg := cmd()
		successMsg, ok := msg.(InfoMsg)
		if !ok {
			t.Fatal("ReportSuccess() command should return InfoMsg")
		}

		if successMsg.Type != InfoTypeSuccess {
			t.Errorf("ReportSuccess() Type = %d, want InfoTypeSuccess (%d)", successMsg.Type, InfoTypeSuccess)
		}
	})
}

func TestReportWarn(t *testing.T) {
	t.Run("ReportWarn returns a command with warn type", func(t *testing.T) {
		cmd := ReportWarn("warn message")

		msg := cmd()
		warnMsg, ok := msg.(InfoMsg)
		if !ok {
			t.Fatal("ReportWarn() command should return InfoMsg")
		}

		if warnMsg.Type != InfoTypeWarn {
			t.Errorf("ReportWarn() Type = %d, want InfoTypeWarn (%d)", warnMsg.Type, InfoTypeWarn)
		}
	})
}

func TestCmdHandler(t *testing.T) {
	t.Run("CmdHandler returns command that produces the message", func(t *testing.T) {
		testMsg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd := CmdHandler(testMsg)

		msg := cmd()
		keyMsg, ok := msg.(tea.KeyMsg)
		if !ok {
			t.Fatal("CmdHandler() should produce tea.KeyMsg")
		}
		if keyMsg.Type != tea.KeyEnter {
			t.Errorf("CmdHandler() produced %v, want KeyEnter", keyMsg.Type)
		}
	})
}

func TestModelInterface(t *testing.T) {
	t.Run("Model interface is satisfiable", func(t *testing.T) {
		var _ Model = mockModel{}

		m := mockModel{}
		cmd := m.Init()
		if cmd != nil {
			t.Error("mockModel.Init() should return nil")
		}

		model, _ := m.Update(nil)
		if model == nil {
			t.Error("mockModel.Update() should return non-nil Model")
		}

		view := m.View()
		if view != "test" {
			t.Errorf("mockModel.View() = %s, want 'test'", view)
		}
	})
}

// mockModel implements Model interface for testing
type mockModel struct{}

func (m mockModel) Init() tea.Cmd { return nil }
func (m mockModel) Update(msg tea.Msg) (Model, tea.Cmd) { return m, nil }
func (m mockModel) View() string { return "test" }
