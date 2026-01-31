package core

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/wwsheng009/taproot/ui/styles"
)

type KeyMapHelp interface {
	Help() help.KeyMap
}

type simpleHelp struct {
	shortList []key.Binding
	fullList  [][]key.Binding
}

func NewSimpleHelp(shortList []key.Binding, fullList [][]key.Binding) help.KeyMap {
	return &simpleHelp{
		shortList: shortList,
		fullList:  fullList,
	}
}

// FullHelp implements help.KeyMap.
func (s *simpleHelp) FullHelp() [][]key.Binding {
	return s.fullList
}

// ShortHelp implements help.KeyMap.
func (s *simpleHelp) ShortHelp() []key.Binding {
	return s.shortList
}

func Section(s *styles.Styles, text string, width int) string {
	char := "─"
	length := lipgloss.Width(text) + 1
	remainingWidth := width - length
	lineStyle := s.Base.Foreground(s.Border)
	if remainingWidth > 0 {
		text = text + " " + lineStyle.Render(strings.Repeat(char, remainingWidth))
	}
	return text
}

func SectionWithInfo(s *styles.Styles, text string, width int, info string) string {
	char := "─"
	length := lipgloss.Width(text) + 1
	remainingWidth := width - length

	if info != "" {
		remainingWidth -= lipgloss.Width(info) + 1 // 1 for the space before info
	}
	lineStyle := s.Base.Foreground(s.Border)
	if remainingWidth > 0 {
		text = text + " " + lineStyle.Render(strings.Repeat(char, remainingWidth)) + " " + info
	}
	return text
}

func Title(s *styles.Styles, title string, width int) string {
	char := "╱"
	length := lipgloss.Width(title) + 1
	remainingWidth := width - length
	titleStyle := s.Base.Foreground(s.Primary)
	if remainingWidth > 0 {
		lines := strings.Repeat(char, remainingWidth)
		lines = styles.ApplyForegroundGrad(s, lines, s.Primary, s.Secondary)
		title = titleStyle.Render(title) + " " + lines
	}
	return title
}

type StatusOpts struct {
	Icon             string // if empty no icon will be shown
	Title            string
	TitleColor       lipgloss.Color
	Description      string
	DescriptionColor lipgloss.Color
	ExtraContent     string // additional content to append after the description
}

func Status(s *styles.Styles, opts StatusOpts, width int) string {
	icon := opts.Icon
	title := opts.Title
	titleColor := s.FgMuted
	if opts.TitleColor != "" {
		titleColor = opts.TitleColor
	}
	description := opts.Description
	descriptionColor := s.FgSubtle
	if opts.DescriptionColor != "" {
		descriptionColor = opts.DescriptionColor
	}
	title = s.Base.Foreground(titleColor).Render(title)
	if description != "" {
		extraContentWidth := lipgloss.Width(opts.ExtraContent)
		if extraContentWidth > 0 {
			extraContentWidth += 1
		}
		description = ansi.Truncate(description, width-lipgloss.Width(icon)-lipgloss.Width(title)-2-extraContentWidth, "…")
		description = s.Base.Foreground(descriptionColor).Render(description)
	}

	content := []string{}
	if icon != "" {
		content = append(content, icon)
	}
	content = append(content, title)
	if description != "" {
		content = append(content, description)
	}
	if opts.ExtraContent != "" {
		content = append(content, opts.ExtraContent)
	}

	return strings.Join(content, " ")
}

type ButtonOpts struct {
	Text           string
	UnderlineIndex int  // Index of character to underline (0-based)
	Selected       bool // Whether this button is selected
}

// SelectableButton creates a button with an underlined character and selection state
func SelectableButton(s *styles.Styles, opts ButtonOpts) string {
	// Base style for the button
	buttonStyle := s.Base

	// Apply selection styling
	if opts.Selected {
		buttonStyle = buttonStyle.Foreground(s.White).Background(s.Secondary)
	} else {
		buttonStyle = buttonStyle.Background(s.BgSubtle)
	}

	// Create the button text with underlined character
	text := opts.Text
	if opts.UnderlineIndex >= 0 && opts.UnderlineIndex < len(text) {
		before := text[:opts.UnderlineIndex]
		underlined := text[opts.UnderlineIndex : opts.UnderlineIndex+1]
		after := text[opts.UnderlineIndex+1:]

		message := buttonStyle.Render(before) +
			buttonStyle.Underline(true).Render(underlined) +
			buttonStyle.Render(after)

		return buttonStyle.Padding(0, 2).Render(message)
	}

	// Fallback if no underline index specified
	return buttonStyle.Padding(0, 2).Render(text)
}

// SelectableButtons creates a horizontal row of selectable buttons
func SelectableButtons(s *styles.Styles, buttons []ButtonOpts, spacing string) string {
	if spacing == "" {
		spacing = "  "
	}

	var parts []string
	for i, button := range buttons {
		parts = append(parts, SelectableButton(s, button))
		if i < len(buttons)-1 {
			parts = append(parts, spacing)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

// SelectableButtonsVertical creates a vertical row of selectable buttons
func SelectableButtonsVertical(s *styles.Styles, buttons []ButtonOpts, spacing int) string {
	var parts []string
	for i, button := range buttons {
		parts = append(parts, SelectableButton(s, button))
		if i < len(buttons)-1 {
			for range spacing {
				parts = append(parts, "")
			}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}
