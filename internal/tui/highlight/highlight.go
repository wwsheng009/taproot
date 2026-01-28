package highlight

import (
	"bytes"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	chromaStyles "github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/tui/styles"
)

// SyntaxHighlight performs syntax highlighting on source code
func SyntaxHighlight(source, fileName string, bg lipgloss.Color) (string, error) {
	// Determine the language lexer to use
	l := lexers.Match(fileName)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Get the formatter
	f := formatters.Get("terminal16m")
	if f == nil {
		f = formatters.Fallback
	}

	style := chroma.MustNewStyle("taproot", styles.GetChromaTheme())

	// Convert lipgloss.Color to RGB values
	r, g, b, _ := lipgloss.Color(bg).RGBA()
	bgColor := chroma.NewColour(uint8(r>>8), uint8(g>>8), uint8(b>>8))

	// Modify the style to use the provided background
	s, err := style.Builder().Transform(
		func(t chroma.StyleEntry) chroma.StyleEntry {
			t.Background = bgColor
			return t
		},
	).Build()
	if err != nil {
		s = chromaStyles.Fallback
	}

	// Tokenize and format
	it, err := l.Tokenise(nil, source)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = f.Format(&buf, s, it)
	return buf.String(), err
}
