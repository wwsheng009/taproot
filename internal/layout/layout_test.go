package layout

import (
	"testing"
)

// mockComponent implements all interfaces for testing
type mockComponent struct {
	focused  bool
	width    int
	height   int
	x, y     int
	helpText []string
}

func (m *mockComponent) Focus()        { m.focused = true }
func (m *mockComponent) Blur()         { m.focused = false }
func (m *mockComponent) Focused() bool { return m.focused }
func (m *mockComponent) Size() (int, int) { return m.width, m.height }
func (m *mockComponent) SetSize(w, h int) { m.width, m.height = w, h }
func (m *mockComponent) Help() []string { return m.helpText }
func (m *mockComponent) Position() (int, int) { return m.x, m.y }
func (m *mockComponent) SetPosition(x, y int) { m.x, m.y = x, y }

func TestInterfacesExist(t *testing.T) {
	t.Run("Focusable interface exists", func(t *testing.T) {
		var _ Focusable = &mockComponent{}
		m := &mockComponent{}
		m.Focus()
		if !m.Focused() {
			t.Error("Focus() should set focused to true")
		}
		m.Blur()
		if m.Focused() {
			t.Error("Blur() should set focused to false")
		}
	})

	t.Run("Sizeable interface exists", func(t *testing.T) {
		var _ Sizeable = &mockComponent{}
		m := &mockComponent{}
		m.SetSize(100, 50)
		w, h := m.Size()
		if w != 100 || h != 50 {
			t.Errorf("Size() = (%d, %d), want (100, 50)", w, h)
		}
	})

	t.Run("Help interface exists", func(t *testing.T) {
		var _ Help = &mockComponent{}
		m := &mockComponent{helpText: []string{"help1", "help2"}}
		help := m.Help()
		if len(help) != 2 {
			t.Errorf("Help() returned %d items, want 2", len(help))
		}
	})

	t.Run("Positional interface exists", func(t *testing.T) {
		var _ Positional = &mockComponent{}
		m := &mockComponent{}
		m.SetPosition(10, 20)
		x, y := m.Position()
		if x != 10 || y != 20 {
			t.Errorf("Position() = (%d, %d), want (10, 20)", x, y)
		}
	})
}
