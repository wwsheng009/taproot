package status

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/tui/util"
)

// TestStatusCmpImmutability 验证 Update 返回新实例
func TestStatusCmpImmutability(t *testing.T) {
	status := NewStatusCmp()
	if status == nil {
		t.Fatal("NewStatusCmp returned nil")
	}

	original := status

	// Test WindowSizeMsg
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updated, _ := status.Update(msg)

	if updated == status {
		t.Error("Update should return new instance for WindowSizeMsg")
	}

	// Original should not be modified
	originalStatus := original.(*statusCmp)
	updatedStatus := updated.(*statusCmp)
	if originalStatus.width == updatedStatus.width && updatedStatus.width == 100 {
		t.Error("Original model should not be modified")
	}

	// Test InfoMsg
	infoMsg := util.InfoMsg{
		Type: util.InfoTypeError,
		Msg:  "Test error message",
		TTL:  time.Second,
	}
	updated, _ = status.Update(infoMsg)

	if updated == status {
		t.Error("Update should return new instance for InfoMsg")
	}

	originalStatus = original.(*statusCmp)
	updatedStatus = updated.(*statusCmp)
	if originalStatus.info.Msg == "Test error message" {
		t.Error("Original model should not be modified by InfoMsg")
	}

	// Test ClearStatusMsg
	updated, _ = status.Update(util.ClearStatusMsg{})

	if updated == status {
		t.Error("Update should return new instance for ClearStatusMsg")
	}
}

// TestStatusCmpToggleFullHelp 验证 ToggleFullHelp 不修改原始模型
func TestStatusCmpToggleFullHelp(t *testing.T) {
	status := NewStatusCmp()

	originalShowAll := status.(*statusCmp).help.ShowAll

	updated := status.ToggleFullHelp()

	if updated == status {
		t.Error("ToggleFullHelp should return new instance")
	}

	// Original should not be modified
	currentShowAll := status.(*statusCmp).help.ShowAll
	if currentShowAll != originalShowAll {
		t.Error("Original model should not be modified by ToggleFullHelp")
	}

	// Updated should have toggled state
	updatedShowAll := updated.(*statusCmp).help.ShowAll
	if updatedShowAll == originalShowAll {
		t.Error("Updated model should have toggled state")
	}
}

// TestStatusCmpSetKeyMap 验证 SetKeyMap 不修改原始模型
func TestStatusCmpSetKeyMap(t *testing.T) {
	status := NewStatusCmp()

	// Use nil as a simple test (real keyMap would implement help.KeyMap interface)
	updated := status.SetKeyMap(nil)

	if updated == status {
		t.Error("SetKeyMap should return new instance")
	}

	// Original should not have keyMap changed (it was nil before)
	if status.(*statusCmp).keyMap != nil {
		t.Error("Original model should not be modified by SetKeyMap")
	}

	// Updated should have nil keyMap
	if updated.(*statusCmp).keyMap != nil {
		t.Error("Updated model should have nil keyMap")
	}
}

// TestStatusCmpWindowSizeMsg 验证窗口大小变化不修改原始模型
func TestStatusCmpWindowSizeMsg(t *testing.T) {
	status := NewStatusCmp()

	originalWidth := status.(*statusCmp).width

	msg := tea.WindowSizeMsg{Width: 120, Height: 60}
	updated, _ := status.Update(msg)

	// Original should not be modified
	if status.(*statusCmp).width != originalWidth {
		t.Error("Original model width should not be modified")
	}

	// Updated should have new width
	if updated.(*statusCmp).width != 120 {
		t.Error("Updated model should have new width")
	}
}

// TestStatusCmpClearStatusMsg 验证清除消息不修改原始模型
func TestStatusCmpClearStatusMsg(t *testing.T) {
	status := NewStatusCmp()

	// First add a message
	infoMsg := util.InfoMsg{
		Type: util.InfoTypeSuccess,
		Msg:  "Operation completed",
		TTL:  time.Second,
	}
	withInfo, _ := status.Update(infoMsg)

	// Now clear it
	updated, _ := withInfo.Update(util.ClearStatusMsg{})

	if updated == withInfo {
		t.Error("Update should return new instance for ClearStatusMsg")
	}

	// Original (withInfo) should still have the message
	if withInfo.(*statusCmp).info.Msg == "" {
		t.Error("Original model should still have message after ClearStatusMsg")
	}

	// Updated should have cleared the message
	if updated.(*statusCmp).info.Msg != "" {
		t.Error("Updated model should have cleared the message")
	}
}

// TestStatusCmpMultipleUpdates 验证多次Update调用
func TestStatusCmpMultipleUpdates(t *testing.T) {
	status := NewStatusCmp()

	// First update
	msg1 := tea.WindowSizeMsg{Width: 100, Height: 50}
	status1, _ := status.Update(msg1)

	// Second update
	msg2 := util.InfoMsg{
		Type: util.InfoTypeInfo,
		Msg:  "Info message",
		TTL:  time.Second,
	}
	status2, _ := status1.Update(msg2)

	// Third update
	status3, _ := status2.Update(util.ClearStatusMsg{})

	// All should be different instances
	if status == status1 || status1 == status2 || status2 == status3 {
		t.Error("Each Update should return a new instance")
	}

	// Original should be unchanged
	if status.(*statusCmp).width == 100 || status.(*statusCmp).info.Msg != "" {
		t.Error("Original model should not be modified by any updates")
	}

	// Each update should reflect its changes
	if status1.(*statusCmp).width != 100 {
		t.Error("status1 should have width 100")
	}
	if status2.(*statusCmp).info.Msg != "Info message" {
		t.Error("status2 should have info message")
	}
	if status3.(*statusCmp).info.Msg != "" {
		t.Error("status3 should have cleared message")
	}
}
