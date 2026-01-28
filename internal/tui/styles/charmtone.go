package styles

import (
	"github.com/charmbracelet/lipgloss"
)

func NewCharmtoneTheme() *Theme {
	t := &Theme{
		Name:   "charmtone",
		IsDark: true,

		Primary:   lipgloss.Color("#1a1a2e"),
		Secondary: lipgloss.Color("#16213e"),
		Tertiary:  lipgloss.Color("#0f3460"),
		Accent:    lipgloss.Color("#e94560"),

		// Backgrounds
		BgBase:        lipgloss.Color("#0f0f1a"),
		BgBaseLighter: lipgloss.Color("#1a1a2e"),
		BgSubtle:      lipgloss.Color("#252540"),
		BgOverlay:     lipgloss.Color("#303050"),

		// Foregrounds
		FgBase:      lipgloss.Color("#e0e0e0"),
		FgMuted:     lipgloss.Color("#b0b0c0"),
		FgHalfMuted: lipgloss.Color("#9090a0"),
		FgSubtle:    lipgloss.Color("#707080"),
		FgSelected:  lipgloss.Color("#ffffff"),

		// Borders
		Border:      lipgloss.Color("#404060"),
		BorderFocus: lipgloss.Color("#e94560"),

		// Status
		Success: lipgloss.Color("#4ade80"),
		Error:   lipgloss.Color("#f87171"),
		Warning: lipgloss.Color("#fbbf24"),
		Info:    lipgloss.Color("#60a5fa"),

		// Colors
		White: lipgloss.Color("#ffffff"),

		BlueLight: lipgloss.Color("#60a5fa"),
		BlueDark:  lipgloss.Color("#1e3a8a"),
		Blue:      lipgloss.Color("#3b82f6"),

		Yellow: lipgloss.Color("#f59e0b"),
		Citron: lipgloss.Color("#84cc16"),

		Green:      lipgloss.Color("#22c55e"),
		GreenDark:  lipgloss.Color("#16a34a"),
		GreenLight: lipgloss.Color("#4ade80"),

		Red:      lipgloss.Color("#ef4444"),
		RedDark:  lipgloss.Color("#b91c1c"),
		RedLight: lipgloss.Color("#f87171"),
		Cherry:   lipgloss.Color("#dc2626"),
	}

	return t
}
