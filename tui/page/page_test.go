package page

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPageID(t *testing.T) {
	t.Run("PageID can be created and compared", func(t *testing.T) {
		id1 := PageID("home")
		id2 := PageID("home")
		id3 := PageID("about")

		if id1 != id2 {
			t.Error("PageID 'home' should equal 'home'")
		}
		if id1 == id3 {
			t.Error("PageID 'home' should not equal 'about'")
		}
	})

	t.Run("PageID can be used as string", func(t *testing.T) {
		id := PageID("settings")
		if string(id) != "settings" {
			t.Errorf("PageID string = %s, want 'settings'", string(id))
		}
	})
}

func TestPageChangeMsg(t *testing.T) {
	t.Run("PageChangeMsg contains page ID", func(t *testing.T) {
		msg := PageChangeMsg{ID: "test"}

		if msg.ID != "test" {
			t.Errorf("PageChangeMsg.ID = %s, want 'test'", msg.ID)
		}
	})

	t.Run("PageChangeMsg is a valid tea.Msg", func(t *testing.T) {
		msg := PageChangeMsg{ID: "home"}
		var _ tea.Msg = msg
	})
}

func TestPageCloseMsg(t *testing.T) {
	t.Run("PageCloseMsg is a valid tea.Msg", func(t *testing.T) {
		msg := PageCloseMsg{}
		var _ tea.Msg = msg
	})
}

func TestPageBackMsg(t *testing.T) {
	t.Run("PageBackMsg is a valid tea.Msg", func(t *testing.T) {
		msg := PageBackMsg{}
		var _ tea.Msg = msg
	})
}
