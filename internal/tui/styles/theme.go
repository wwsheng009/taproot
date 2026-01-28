package styles

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/uniseg"
)

type Theme struct {
	Name   string
	IsDark bool

	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Tertiary  lipgloss.Color
	Accent    lipgloss.Color

	BgBase        lipgloss.Color
	BgBaseLighter lipgloss.Color
	BgSubtle      lipgloss.Color
	BgOverlay     lipgloss.Color

	FgBase      lipgloss.Color
	FgMuted     lipgloss.Color
	FgHalfMuted lipgloss.Color
	FgSubtle    lipgloss.Color
	FgSelected  lipgloss.Color

	Border      lipgloss.Color
	BorderFocus lipgloss.Color

	Success lipgloss.Color
	Error   lipgloss.Color
	Warning lipgloss.Color
	Info    lipgloss.Color

	// Colors
	// White
	White lipgloss.Color

	// Blues
	BlueLight lipgloss.Color
	BlueDark  lipgloss.Color
	Blue      lipgloss.Color

	// Yellows
	Yellow lipgloss.Color
	Citron lipgloss.Color

	// Greens
	Green      lipgloss.Color
	GreenDark  lipgloss.Color
	GreenLight lipgloss.Color

	// Reds
	Red      lipgloss.Color
	RedDark  lipgloss.Color
	RedLight lipgloss.Color
	Cherry   lipgloss.Color

	styles     *Styles
	stylesOnce sync.Once
}

type Styles struct {
	Base         lipgloss.Style
	SelectedBase lipgloss.Style

	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Text     lipgloss.Style
	TextSelected lipgloss.Style
	Muted    lipgloss.Style
	Subtle   lipgloss.Style

	Success lipgloss.Style
	Error   lipgloss.Style
	Warning lipgloss.Style
	Info    lipgloss.Style

	// Help
	Help help.Styles

	// Markdown & Chroma
	Markdown ansi.StyleConfig
}

func (t *Theme) S() *Styles {
	t.stylesOnce.Do(func() {
		t.styles = t.buildStyles()
	})
	return t.styles
}

func (t *Theme) buildStyles() *Styles {
	base := lipgloss.NewStyle().
		Foreground(t.FgBase)
	return &Styles{
		Base: base,

		SelectedBase: base.Background(t.Primary),

		Title: base.
			Foreground(t.Accent).
			Bold(true),

		Subtitle: base.
			Foreground(t.Secondary).
			Bold(true),

		Text:         base,
		TextSelected: base.Background(t.Primary).Foreground(t.FgSelected),

		Muted: base.Foreground(t.FgMuted),

		Subtle: base.Foreground(t.FgSubtle),

		Success: base.Foreground(t.Success),

		Error: base.Foreground(t.Error),

		Warning: base.Foreground(t.Warning),

		Info: base.Foreground(t.Info),

		Markdown: ansi.StyleConfig{
			Document: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					// BlockPrefix: "\n",
					// BlockSuffix: "\n",
					Color: stringPtr(ColorSmoke),
				},
				// Margin: uintPtr(defaultMargin),
			},
			BlockQuote: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
				Indent:         uintPtr(1),
				IndentToken:    stringPtr("│ "),
			},
			List: ansi.StyleList{
				LevelIndent: 2, // defaultListIndent hardcoded or define constant if available
			},
			Heading: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					BlockSuffix: "\n",
					Color:       stringPtr(ColorMalibu),
					Bold:        boolPtr(true),
				},
			},
			H1: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix:          " ",
					Suffix:          " ",
					Color:           stringPtr(ColorZest),
					BackgroundColor: stringPtr(ColorCharple),
					Bold:            boolPtr(true),
				},
			},
			H2: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "## ",
				},
			},
			H3: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "### ",
				},
			},
			H4: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "#### ",
				},
			},
			H5: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "##### ",
				},
			},
			H6: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix: "###### ",
					Color:  stringPtr(ColorGuac),
					Bold:   boolPtr(false),
				},
			},
			Strikethrough: ansi.StylePrimitive{
				CrossedOut: boolPtr(true),
			},
			Emph: ansi.StylePrimitive{
				Italic: boolPtr(true),
			},
			Strong: ansi.StylePrimitive{
				Bold: boolPtr(true),
			},
			HorizontalRule: ansi.StylePrimitive{
				Color:  stringPtr(ColorCharcoal),
				Format: "\n--------\n",
			},
			Item: ansi.StylePrimitive{
				BlockPrefix: "• ",
			},
			Enumeration: ansi.StylePrimitive{
				BlockPrefix: ". ",
			},
			Task: ansi.StyleTask{
				StylePrimitive: ansi.StylePrimitive{},
				Ticked:         "[✓] ",
				Unticked:       "[ ] ",
			},
			Link: ansi.StylePrimitive{
				Color:     stringPtr(ColorZinc),
				Underline: boolPtr(true),
			},
			LinkText: ansi.StylePrimitive{
				Color: stringPtr(ColorGuac),
				Bold:  boolPtr(true),
			},
			Image: ansi.StylePrimitive{
				Color:     stringPtr(ColorCheeky),
				Underline: boolPtr(true),
			},
			ImageText: ansi.StylePrimitive{
				Color:  stringPtr(ColorSquid),
				Format: "Image: {{.text}} →",
			},
			Code: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Prefix:          " ",
					Suffix:          " ",
					Color:           stringPtr(ColorCoral),
					BackgroundColor: stringPtr(ColorCharcoal),
				},
			},
			CodeBlock: ansi.StyleCodeBlock{
				StyleBlock: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Color: stringPtr(ColorCharcoal),
					},
					Margin: uintPtr(2), // defaultMargin
				},
				Chroma: &ansi.Chroma{
					Text: ansi.StylePrimitive{
						Color: stringPtr(ColorSmoke),
					},
					Error: ansi.StylePrimitive{
						Color:           stringPtr(ColorButter),
						BackgroundColor: stringPtr(ColorSriracha),
					},
					Comment: ansi.StylePrimitive{
						Color: stringPtr(ColorOyster),
					},
					CommentPreproc: ansi.StylePrimitive{
						Color: stringPtr(ColorBengal),
					},
					Keyword: ansi.StylePrimitive{
						Color: stringPtr(ColorMalibu),
					},
					KeywordReserved: ansi.StylePrimitive{
						Color: stringPtr(ColorPony),
					},
					KeywordNamespace: ansi.StylePrimitive{
						Color: stringPtr(ColorPony),
					},
					KeywordType: ansi.StylePrimitive{
						Color: stringPtr(ColorGuppy),
					},
					Operator: ansi.StylePrimitive{
						Color: stringPtr(ColorSalmon),
					},
					Punctuation: ansi.StylePrimitive{
						Color: stringPtr(ColorZest),
					},
					Name: ansi.StylePrimitive{
						Color: stringPtr(ColorSmoke),
					},
					NameBuiltin: ansi.StylePrimitive{
						Color: stringPtr(ColorCheeky),
					},
					NameTag: ansi.StylePrimitive{
						Color: stringPtr(ColorMauve),
					},
					NameAttribute: ansi.StylePrimitive{
						Color: stringPtr(ColorHazy),
					},
					NameClass: ansi.StylePrimitive{
						Color:     stringPtr(ColorSalt),
						Underline: boolPtr(true),
						Bold:      boolPtr(true),
					},
					NameDecorator: ansi.StylePrimitive{
						Color: stringPtr(ColorCitron),
					},
					NameFunction: ansi.StylePrimitive{
						Color: stringPtr(ColorGuac),
					},
					LiteralNumber: ansi.StylePrimitive{
						Color: stringPtr(ColorJulep),
					},
					LiteralString: ansi.StylePrimitive{
						Color: stringPtr(ColorCumin),
					},
					LiteralStringEscape: ansi.StylePrimitive{
						Color: stringPtr(ColorBok),
					},
					GenericDeleted: ansi.StylePrimitive{
						Color: stringPtr(ColorCoral),
					},
					GenericEmph: ansi.StylePrimitive{
						Italic: boolPtr(true),
					},
					GenericInserted: ansi.StylePrimitive{
						Color: stringPtr(ColorGuac),
					},
					GenericStrong: ansi.StylePrimitive{
						Bold: boolPtr(true),
					},
					GenericSubheading: ansi.StylePrimitive{
						Color: stringPtr(ColorSquid),
					},
					Background: ansi.StylePrimitive{
						BackgroundColor: stringPtr(ColorCharcoal),
					},
				},
			},
			Table: ansi.StyleTable{
				StyleBlock: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{},
				},
			},
			DefinitionDescription: ansi.StylePrimitive{
				BlockPrefix: "\n ",
			},
		},

		Help: help.Styles{
			ShortKey:       base.Foreground(t.FgMuted),
			ShortDesc:      base.Foreground(t.FgSubtle),
			ShortSeparator: base.Foreground(t.Border),
			Ellipsis:       base.Foreground(t.Border),
			FullKey:        base.Foreground(t.FgMuted),
			FullDesc:       base.Foreground(t.FgSubtle),
			FullSeparator:  base.Foreground(t.Border),
		},
	}
}

type Manager struct {
	themes  map[string]*Theme
	current *Theme
}

var (
	defaultManager     *Manager
	defaultManagerOnce sync.Once
)

func initDefaultManager() *Manager {
	defaultManagerOnce.Do(func() {
		defaultManager = newManager()
	})
	return defaultManager
}

func SetDefaultManager(m *Manager) {
	defaultManager = m
}

func DefaultManager() *Manager {
	return initDefaultManager()
}

func CurrentTheme() *Theme {
	return initDefaultManager().Current()
}

func newManager() *Manager {
	m := &Manager{
		themes: make(map[string]*Theme),
	}

	t := NewCharmtoneTheme() // default theme
	m.Register(t)
	m.current = m.themes[t.Name]

	return m
}

func (m *Manager) Register(theme *Theme) {
	m.themes[theme.Name] = theme
}

func (m *Manager) Current() *Theme {
	return m.current
}

func (m *Manager) SetTheme(name string) error {
	if theme, ok := m.themes[name]; ok {
		m.current = theme
		return nil
	}
	return fmt.Errorf("theme %s not found", name)
}

func (m *Manager) List() []string {
	names := make([]string, 0, len(m.themes))
	for name := range m.themes {
		names = append(names, name)
	}
	return names
}

func ForegroundGrad(input string, bold bool, color1, color2 lipgloss.Color) []string {
	if input == "" {
		return []string{""}
	}
	t := CurrentTheme()
	if len(input) == 1 {
		style := t.S().Base.Foreground(color1)
		if bold {
			style.Bold(true)
		}
		return []string{style.Render(input)}
	}
	var clusters []string
	gr := uniseg.NewGraphemes(input)
	for gr.Next() {
		clusters = append(clusters, string(gr.Runes()))
	}

	// Parse colors for blending
	c1, _ := colorful.MakeColor(lipgloss.Color(color1))
	c2, _ := colorful.MakeColor(lipgloss.Color(color2))

	ramp := blendLipglossColors(len(clusters), c1, c2)
	for i, c := range ramp {
		style := t.S().Base.Foreground(lipgloss.Color(c.Hex()))
		if bold {
			style.Bold(true)
		}
		clusters[i] = style.Render(clusters[i])
	}
	return clusters
}

// ApplyForegroundGrad renders a given string with a horizontal gradient
// foreground.
func ApplyForegroundGrad(input string, color1, color2 lipgloss.Color) string {
	if input == "" {
		return ""
	}
	var o strings.Builder
	clusters := ForegroundGrad(input, false, color1, color2)
	for _, c := range clusters {
		fmt.Fprint(&o, c)
	}
	return o.String()
}

// ApplyBoldForegroundGrad renders a given string with a horizontal gradient
// foreground.
func ApplyBoldForegroundGrad(input string, color1, color2 lipgloss.Color) string {
	if input == "" {
		return ""
	}
	var o strings.Builder
	clusters := ForegroundGrad(input, true, color1, color2)
	for _, c := range clusters {
		fmt.Fprint(&o, c)
	}
	return o.String()
}

// blendLipglossColors returns a slice of colorful.Color blended between the given keys.
// Blending is done in Hcl to stay in gamut.
func blendLipglossColors(size int, c1, c2 colorful.Color) []colorful.Color {
	if size < 1 {
		return nil
	}

	blended := make([]colorful.Color, size)

	for j := range size {
		var t float64
		if size > 1 {
			t = float64(j) / float64(size-1)
		}
		c := c1.BlendHcl(c2, t)
		blended[j] = c
	}

	return blended
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func uintPtr(u uint) *uint {
	return &u
}
