package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width      int
	height     int
	quitting   bool
	detecting  bool
	sixel      bool
	kitty      bool
	iterm2     bool
	tmux       bool
	term       string
	program    string
	error      string
}

func initialModel() model {
	return model{
		detecting: true,
	}
}

type detectionCompleteMsg struct {
	sixel   bool
	kitty   bool
	iterm2  bool
	tmux    bool
	term    string
	program string
	error   string
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(300 * time.Millisecond) // Simulate detection time

		term := os.Getenv("TERM")
		program := os.Getenv("TERM_PROGRAM")
		tmux := os.Getenv("TMUX") != ""

		// Simple detection for demonstration
		sixel := false
		kitty := false
		iterm2 := false

		// Environment-based detection
		if strings.Contains(program, "kitty") || strings.Contains(term, "kitty") {
			kitty = true
		}
		if program == "iTerm.app" {
			iterm2 = true
		}
		if program == "WezTerm" || strings.Contains(program, "mintty") || strings.Contains(term, "xterm") {
			sixel = true
		}

		return detectionCompleteMsg{
			sixel:   sixel,
			kitty:   kitty,
			iterm2:  iterm2,
			tmux:    tmux,
			term:    term,
			program: program,
			error:   "",
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case detectionCompleteMsg:
		m.detecting = false
		m.sixel = msg.sixel
		m.kitty = msg.kitty
		m.iterm2 = msg.iterm2
		m.tmux = msg.tmux
		m.term = msg.term
		m.program = msg.program
		m.error = msg.error
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Align(lipgloss.Center)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("43")).
		Bold(true)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226"))

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("61")).
		Padding(1).
		Width(min(80, m.width-4))

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Width(m.width).Render("🔍 Terminal Capability Checker"))
	b.WriteString("\n\n")

	// Detection status
	if m.detecting {
		b.WriteString(mutedStyle.Render("Detecting terminal capabilities..."))
		b.WriteString("\n")
		b.WriteString(strings.Repeat("░", min(60, m.width-4)))
		return b.String()
	}

	// Environment info
	b.WriteString(sectionStyle.Render("Environment Information"))
	b.WriteString("\n\n")

	fmt.Fprintf(&b, "  TERM:        %s\n", m.term)
	fmt.Fprintf(&b, "  TERM_PROGRAM: %s\n", m.program)
	fmt.Fprintf(&b, "  TMUX:        ")
	if m.tmux {
		b.WriteString(successStyle.Render("✓ Yes") + "\n")
	} else {
		b.WriteString(mutedStyle.Render("✗ No") + "\n")
	}

	b.WriteString("\n")

	// Protocol support
	b.WriteString(sectionStyle.Render("Image Rendering Protocol Support"))
	b.WriteString("\n\n")

	// Sixel
	fmt.Fprintf(&b, "  Sixel Protocol:   ")
	if m.sixel {
		icon := successStyle.Render("✓ SUPPORTED")
		if m.tmux {
			icon += warningStyle.Render(" ⚠ (tmux)")
		}
		b.WriteString(icon + "\n")
	} else {
		b.WriteString(errorStyle.Render("✗ NOT SUPPORTED") + "\n")
	}

	// Kitty
	fmt.Fprintf(&b, "  Kitty Protocol:   ")
	if m.kitty {
		b.WriteString(successStyle.Render("✓ SUPPORTED") + "\n")
	} else {
		b.WriteString(mutedStyle.Render("✗ NOT AVAILABLE") + "\n")
	}

	// iTerm2
	fmt.Fprintf(&b, "  iTerm2 Protocol:  ")
	if m.iterm2 {
		b.WriteString(successStyle.Render("✓ SUPPORTED") + "\n")
	} else {
		b.WriteString(mutedStyle.Render("✗ NOT AVAILABLE") + "\n")
	}

	// Unicode Block (always available)
	b.WriteString(successStyle.Render("  Unicode Blocks:   ✓ SUPPORTED (fallback)") + "\n")

	b.WriteString("\n")

	// Recommended renderer
	b.WriteString(sectionStyle.Render("Recommended Renderer"))
	b.WriteString("\n\n")

	var recommendation string
	switch {
	case m.sixel:
		recommendation = "Sixel (High-Quality)"
	case m.kitty:
		recommendation = "Kitty (Native)"
	case m.iterm2:
		recommendation = "iTerm2 (Native)"
	default:
		recommendation = "Unicode Blocks (Compatible)"
	}

	recommendationBox := boxStyle.Render(
		fmt.Sprintf("Best: %s", recommendation),
	)
	b.WriteString(recommendationBox)

	b.WriteString("\n\n")

	// Notes
	if m.tmux && m.sixel {
		notes := warningStyle.Render(
			"⚠️  Note: Sixel in tmux requires 'set -g allow-passthrough on'\n" +
			"   in ~/.tmux.conf for best results.",
		)
		b.WriteString(notes)
		b.WriteString("\n\n")
	}

	// Footer
	hints := mutedStyle.Render(
		"Press q or ctrl+c to quit",
	)
	b.WriteString(hints)

	return b.String()
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
