package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/components/progress"
	"github.com/wwsheng009/taproot/ui/render"
)

const (
	maxWidth  = 80
	maxHeight = 24
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7c3aed")).
			Bold(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280")).
			Italic(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d1d5db"))
)

type tickMsg time.Time

type model struct {
	progressBars []*	progress.ProgressBar
	spinners     []*progress.Spinner
	width        int
	height       int
	showBars     bool
	showSpinners bool
	spinnerIndex int
}

func initialModel() model {
	// Create sample progress bars
	pb1 := progress.NewProgressBar(100)
	pb1.SetLabel("Download:")
	pb1.SetCurrent(45)

	pb2 := progress.NewProgressBar(200)
	pb2.SetLabel("Processing:")
	pb2.SetCurrent(120)

	pb3 := progress.NewProgressBar(50)
	pb3.SetLabel("Upload:")
	pb3.SetCurrent(10)

	// Create sample spinners
	sp1 := progress.NewSpinner()
	sp1.SetLabel("Loading files...")

	sp2 := progress.NewSpinner()
	sp2.SetLabel("Searching...")
	sp2.SetType(progress.SpinnerTypeArrow)

	sp3 := progress.NewSpinner()
	sp3.SetLabel("Computing...")
	sp3.SetType(progress.SpinnerTypeMoon)

	sp4 := progress.NewSpinner()
	sp4.SetLabel("Syncing...")
	sp4.SetType(progress.SpinnerTypeLine)

	return model{
		progressBars: []*progress.ProgressBar{pb1, pb2, pb3},
		spinners:     []*progress.Spinner{sp1, sp2, sp3, sp4},
		width:        maxWidth,
		height:       maxHeight,
		showBars:     true,
		showSpinners: true,
		spinnerIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.tickCmd(),
	)
}

func (m model) tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "b":
			// Toggle progress bars
			m.showBars = !m.showBars

		case "s":
			// Toggle spinners
			m.showSpinners = !m.showSpinners

		case "up":
			// Increment all progress bars
			for _, pb := range m.progressBars {
				if !pb.Completed() {
					pb.Increment()
				}
			}

		case "down":
			// Decrement all progress bars
			for _, pb := range m.progressBars {
				if pb.Current() > 0 {
					pb.SetCurrent(pb.Current() - 1)
				}
			}

		case "r":
			// Reset all progress bars
			for _, pb := range m.progressBars {
				pb.Reset()
			}

		case "t":
			// Cycle spinner types
			for _, sp := range m.spinners {
				sp.SetType(sp.Type() + 1)
				if sp.Type() > progress.SpinnerTypeMoon {
					sp.SetType(progress.SpinnerTypeDots)
				}
			}

		case "1", "2", "3":
			// Increase specific progress bar
			idx := int(msg.String()[0]) - int('1')
			if idx >= 0 && idx < len(m.progressBars) {
				m.progressBars[idx].Increment()
			}
		}

	case tickMsg:
		// Update spinners with generic tick messages
		for _, sp := range m.spinners {
			if sp.Running() {
				sp.Update(&render.TickMsg{})
			}
		}
		return m, m.tickCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width > maxWidth {
			m.width = maxWidth
		}
		for _, pb := range m.progressBars {
			pb.Style().Width = m.width - 10
		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Progress Bar & Spinner Demo"))
	b.WriteString("\n\n")

	// Progress bars section
	if m.showBars {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Progress Bars"))
		b.WriteString("\n")

		for _, pb := range m.progressBars {
			b.WriteString(pb.View())
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Spinners section
	if m.showSpinners {
		b.WriteString(lipgloss.NewStyle().Bold(true).Render("Spinners"))
		b.WriteString("\n")

		for _, sp := range m.spinners {
			b.WriteString(sp.View())
			b.WriteString("  ")
		}
		b.WriteString("\n\n")
	}

	// Progress summary
	total := 0.0
	completed := 0
	for _, pb := range m.progressBars {
		total += pb.Percent()
		if pb.Completed() {
			completed++
		}
	}
	avg := total / float64(len(m.progressBars))

	summary := fmt.Sprintf(
		"Progress Summary: %.1f%% average | %d/%d completed",
		avg,
		completed,
		len(m.progressBars),
	)
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Render(summary))
	b.WriteString("\n\n")

	// Separator
	b.WriteString(separatorStyle.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")

	// Help text
	help := "Controls: ↑/down • rreset • btoggle bars • stoggle spinners • ttype • 1-3increment specific bar • qquit"
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
