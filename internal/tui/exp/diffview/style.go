package diffview

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// LineStyle contains styling for different parts of a diff line.
type LineStyle struct {
	LineNumber lipgloss.Style
	Symbol     lipgloss.Style
	Code       lipgloss.Style
}

// Style defines the styling for diff view elements.
type Style struct {
	// Divider line between before/after in split view
	DividerLine lipgloss.Style

	// Line number styling
	LineNumber lipgloss.Style

	// Line types
	MissingLine lipgloss.Style // Empty line placeholder
	EqualLine   LineStyle
	InsertLine  LineStyle
	DeleteLine  LineStyle
	HeaderLine  LineStyle

	// Syntax highlighting
	SyntaxHighlight bool
	SyntaxTheme     string
}

// DefaultLightStyle returns the default light theme style.
func DefaultLightStyle() Style {
	sty := styles.DefaultStyles()
	base := sty.Base

	return Style{
		DividerLine: base.Foreground(sty.Border),
		LineNumber: base.Foreground(sty.FgMuted),
		MissingLine: base.Foreground(lipgloss.Color("#e0e0e0")),
		EqualLine: LineStyle{
			LineNumber: base.Foreground(sty.FgMuted),
			Symbol:     base.Foreground(sty.FgMuted),
			Code:       base.Foreground(sty.FgBase),
		},
		InsertLine: LineStyle{
			LineNumber: base.Foreground(sty.FgMuted),
			Symbol:     base.Foreground(lipgloss.Color("#2da44e")),
			Code:       base.Foreground(sty.FgBase).Background(lipgloss.Color("#d8f8dd")),
		},
		DeleteLine: LineStyle{
			LineNumber: base.Foreground(sty.FgMuted),
			Symbol:     base.Foreground(lipgloss.Color("#f85149")),
			Code:       base.Foreground(sty.FgBase).Background(lipgloss.Color("#ffebe9")),
		},
		HeaderLine: LineStyle{
			LineNumber: base.Foreground(sty.FgMuted),
			Symbol:     base.Foreground(sty.FgMuted),
			Code:       base.Bold(true).Foreground(sty.FgMuted),
		},
		SyntaxHighlight: true,
		SyntaxTheme:     "github",
	}
}

// DefaultDarkStyle returns the default dark theme style.
func DefaultDarkStyle() Style {
	sty := styles.DefaultStyles()
	base := sty.Base

	return Style{
		DividerLine: base.Foreground(sty.Border),
		LineNumber: base.Foreground(sty.FgMuted),
		MissingLine: base.Foreground(lipgloss.Color("#3b3b3b")),
		EqualLine: LineStyle{
			LineNumber: base.Foreground(sty.FgMuted),
			Symbol:     base.Foreground(sty.FgMuted),
			Code:       base.Foreground(sty.FgBase),
		},
		InsertLine: LineStyle{
			LineNumber: base.Foreground(sty.FgMuted),
			Symbol:     base.Foreground(lipgloss.Color("#3fb950")),
			Code:       base.Foreground(sty.FgBase).Background(lipgloss.Color("#1a3f29")),
		},
		DeleteLine: LineStyle{
			LineNumber: base.Foreground(sty.FgMuted),
			Symbol:     base.Foreground(lipgloss.Color("#f85149")),
			Code:       base.Foreground(sty.FgBase).Background(lipgloss.Color("#3b1d1d")),
		},
		HeaderLine: LineStyle{
			LineNumber: base.Foreground(sty.FgMuted),
			Symbol:     base.Foreground(sty.FgMuted),
			Code:       base.Bold(true).Foreground(sty.FgMuted),
		},
		SyntaxHighlight: true,
		SyntaxTheme:     "monokai",
	}
}
