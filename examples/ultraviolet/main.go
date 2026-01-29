package main

import (
	"fmt"
	"strings"

	"github.com/yourorg/taproot/internal/ui/render"
)

// UVCounterModel demonstrates the Ultraviolet engine with a simple counter
type UVCounterModel struct {
	count int
	paused bool
}

// NewUVCounterModel creates a new counter model
func NewUVCounterModel() *UVCounterModel {
	return &UVCounterModel{
		count:  0,
		paused: false,
	}
}

// Init initializes the model
func (m *UVCounterModel) Init() error {
	return nil
}

// Update handles incoming messages
func (m *UVCounterModel) Update(msg any) (render.Model, render.Cmd) {
	if msg == nil {
		return m, render.None()
	}

	switch msg := msg.(type) {
	case render.KeyMsg:
		switch msg.Key {
		case "q", "ctrl+c":
			// Quit
			return m, render.Quit()
		case "up", "right", "+", "=":
			if !m.paused {
				m.count++
			}
		case "down", "left", "-", "_":
			if m.count > 0 && !m.paused {
				m.count--
			}
		case " ", "enter":
			m.paused = !m.paused
		case "r":
			m.count = 0
		}
	}

	return m, render.None()
}

// View returns the string representation for rendering
func (m *UVCounterModel) View() string {
	var b strings.Builder

	b.WriteString("Taproot v2.0.0 - Ultraviolet Engine Demo\n")
	b.WriteString("==========================================\n\n")
	b.WriteString("This demo uses the Ultraviolet rendering engine\n")
	b.WriteString("for high-performance terminal UI rendering.\n\n")
	b.WriteString("Keys:\n")
	b.WriteString("  ↑/→ or +/- : Increment counter (if not paused)\n")
	b.WriteString("  ↓/← or -_  : Decrement counter (if not paused)\n")
	b.WriteString("  space/enter : Toggle pause\n")
	b.WriteString("  r           : Reset counter\n")
	b.WriteString("  q/ctrl+c    : Quit\n\n")

	b.WriteString("──────────────────────────────────────────\n\n")

	b.WriteString(fmt.Sprintf("Count: %d\n\n", m.count))

	// Progress bar using block characters
	bars := 50
	filled := int(float64(m.count)/100.0 * float64(bars))
	if filled > bars {
		filled = bars
	}

	b.WriteString("[")
	for i := 0; i < bars; i++ {
		if i < filled {
			b.WriteString("█")
		} else {
			b.WriteString("░")
		}
	}
	b.WriteString(fmt.Sprintf("] %d%%\n\n", m.count))

	// Status
	if m.paused {
		b.WriteString("⏸️  PAUSED")
	} else {
		b.WriteString("▶️  RUNNING")
	}

	return b.String()
}

func main() {
	// Create Ultraviolet engine
	config := render.DefaultConfig()
	config.EnableAltScreen = true

	engine, err := render.CreateEngine(render.EngineUltraviolet, config)
	if err != nil {
		fmt.Printf("Failed to create engine: %v\n", err)
		return
	}

	model := NewUVCounterModel()

	fmt.Println("Starting Ultraviolet engine...")
	fmt.Println("Press q or ctrl+c to quit")

	if err := engine.Start(model); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
