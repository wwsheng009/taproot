package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/yourorg/taproot/internal/ui/render"
)

// DualEngineModel demonstrates the same model working with both engines
type DualEngineModel struct {
	count    int
	history  []int
	engine   string
	fps      int
	lastTime time.Time
}

// NewDualEngineModel creates a new dual-engine model
func NewDualEngineModel(engine string) *DualEngineModel {
	return &DualEngineModel{
		count:    0,
		history:  make([]int, 0, 50),
		engine:   engine,
		fps:      60,
		lastTime: time.Now(),
	}
}

// Init initializes the model
func (m *DualEngineModel) Init() error {
	return nil
}

// Update handles incoming messages
func (m *DualEngineModel) Update(msg any) (render.Model, render.Cmd) {
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
			m.count++
		case "down", "left", "-", "_":
			if m.count > 0 {
				m.count--
			}
		case " ":
			// Add to history
			m.history = append(m.history, m.count)
			if len(m.history) > 50 {
				m.history = m.history[1:]
			}
		}
	case render.TickMsg:
		// Calculate FPS
		now := time.Now()
		elapsed := now.Sub(m.lastTime)
		if elapsed > 0 {
			m.fps = int(time.Second / elapsed)
		}
		m.lastTime = now
	}

	return m, render.None()
}

// View returns the string representation for rendering
func (m *DualEngineModel) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Taproot v2.0.0 - Engine Comparison (%s)\n", m.engine))
	b.WriteString("===============================================\n\n")
	b.WriteString("Same Model • Different Engines\n\n")

	b.WriteString("Keys:\n")
	b.WriteString("  ↑/→ or +/- : Increment counter\n")
	b.WriteString("  ↓/← or -_  : Decrement counter\n")
	b.WriteString("  space       : Save to history\n")
	b.WriteString("  q/ctrl+c    : Quit\n\n")

	b.WriteString("─────────────────────────────────────────────\n\n")

	// Display stats
	engineType := render.EngineBubbletea
	if m.engine == "ultraviolet" {
		engineType = render.EngineUltraviolet
	}

	b.WriteString(fmt.Sprintf("Current Engine : %s\n", engineType.String()))
	b.WriteString(fmt.Sprintf("Current Count  : %d\n", m.count))
	b.WriteString(fmt.Sprintf("History Size   : %d\n", len(m.history)))
	b.WriteString(fmt.Sprintf("Estimated FPS  : %d\n\n", m.fps))

	// Progress bar visualization
	bars := 60
	filled := int(float64(m.count)/100.0 * float64(bars))
	if filled > bars {
		filled = bars
	}

	b.WriteString(fmt.Sprintf("Progress Bar:\n  [%s", strings.Repeat("█", filled)))
	b.WriteString(fmt.Sprintf("%s] %d%%\n\n", strings.Repeat("░", bars-filled), m.count))

	// History graph
	if len(m.history) > 0 {
		b.WriteString("History Graph (last 50 values):\n")

		maxVal := 100
		for _, val := range m.history {
			if val > maxVal {
				maxVal = val
			}
		}

		graphWidth := 40
		visibleHistory := m.history
		if len(visibleHistory) > 20 {
			visibleHistory = visibleHistory[len(visibleHistory)-20:]
		}

		for i, val := range visibleHistory {
			barsInRow := int(float64(val) / float64(maxVal) * float64(graphWidth))
			b.WriteString("  ")
			b.WriteString(fmt.Sprintf("%3d │", val))
			b.WriteString(strings.Repeat("█", barsInRow))
			b.WriteString(strings.Repeat("░", graphWidth-barsInRow))
			if i == len(visibleHistory)-1 {
				b.WriteString(" ← Current")
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func main() {
	engineType := flag.String("engine", "bubbletea", "Rendering engine: bubbletea or ultraviolet")
	flag.Parse()

	var engine render.Engine
	var err error

	config := render.DefaultConfig()
	config.EnableAltScreen = true

	switch *engineType {
	case "ultraviolet":
		engine, err = render.CreateEngine(render.EngineUltraviolet, config)
	case "bubbletea", "":
		engine, err = render.CreateEngine(render.EngineBubbletea, config)
	default:
		fmt.Printf("Unknown engine: %s\nUse: -engine=[bubbletea|ultraviolet]\n", *engineType)
		return
	}

	if err != nil {
		fmt.Printf("Failed to create engine: %v\n", err)
		return
	}

	model := NewDualEngineModel(*engineType)

	fmt.Printf("Starting %s engine...\n", *engineType)
	fmt.Println("Press q or ctrl+c to quit")

	if err := engine.Start(model); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
