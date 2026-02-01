package forms

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// RadioGroup represents a single-select radio button group.
type RadioGroup struct {
	label    string
	options  []string
	selected int
	focused  bool

	// Styles
	focusedStyle  lipgloss.Style
	blurredStyle  lipgloss.Style
	selectedStyle lipgloss.Style
	labelStyle    lipgloss.Style
	optionStyle   lipgloss.Style

	// Validation
	validators []Validator
	err        error
	errorStyle lipgloss.Style
}

// NewRadioGroup creates a new radio group.
func NewRadioGroup(label string, options []string) *RadioGroup {
	s := styles.DefaultStyles()
	return &RadioGroup{
		label:    label,
		options:  options,
		selected: 0,
		focusedStyle: lipgloss.NewStyle().
			Foreground(s.Primary),
		blurredStyle: lipgloss.NewStyle().
			Foreground(s.FgMuted),
		selectedStyle: lipgloss.NewStyle().
			Foreground(s.Primary).
			Bold(true),
		labelStyle: lipgloss.NewStyle().
			Foreground(s.FgBase).
			Bold(true),
		optionStyle: lipgloss.NewStyle().
			PaddingLeft(2),
		errorStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
	}
}

// Init implements render.Model.
func (r *RadioGroup) Init() render.Cmd {
	return nil
}

// Update implements render.Model.
func (r *RadioGroup) Update(msg any) (render.Model, render.Cmd) {
	if !r.focused {
		return r, nil
	}

	var keyStr string
	if k, ok := msg.(tea.KeyMsg); ok {
		keyStr = k.String()
	} else if k, ok := msg.(render.KeyMsg); ok {
		keyStr = k.String()
	}

	switch keyStr {
	case "up", "k":
		if r.selected > 0 {
			r.selected--
		} else {
			r.selected = len(r.options) - 1 // Wrap around
		}
	case "down", "j":
		if r.selected < len(r.options)-1 {
			r.selected++
		} else {
			r.selected = 0 // Wrap around
		}
	case " ":
		// Space just confirms selection (already selected by navigation)
		// But if we supported "unselect", we'd do it here. 
		// Radio buttons usually can't be deselected once selected.
	}

	return r, nil
}

// View implements render.Model.
func (r *RadioGroup) View() string {
	var b strings.Builder

	// Render Group Label
	if r.label != "" {
		b.WriteString(r.labelStyle.Render(r.label))
		b.WriteString("\n")
	}

	// Render Options
	for i, option := range r.options {
		var indicator string
		var text string

		if i == r.selected {
			if r.focused {
				indicator = r.focusedStyle.Render("(o)")
				text = r.selectedStyle.Render(option)
			} else {
				indicator = r.blurredStyle.Render("(o)")
				text = r.blurredStyle.Render(option)
			}
		} else {
			indicator = r.blurredStyle.Render("( )")
			text = r.blurredStyle.Render(option)
		}

		b.WriteString(r.optionStyle.Render(indicator + " " + text))
		if i < len(r.options)-1 {
			b.WriteString("\n")
		}
	}

	// Render error
	if r.err != nil {
		b.WriteString("\n")
		b.WriteString(r.errorStyle.Render(r.err.Error()))
	}

	return b.String()
}

// Value returns the selected option string.
func (r *RadioGroup) Value() string {
	if len(r.options) == 0 {
		return ""
	}
	return r.options[r.selected]
}

// SetValue sets the selected option by value string.
// If value not found, it defaults to the first option.
func (r *RadioGroup) SetValue(v string) {
	for i, opt := range r.options {
		if opt == v {
			r.selected = i
			return
		}
	}
	r.selected = 0
}

// SelectedIndex returns the index of the selected option.
func (r *RadioGroup) SelectedIndex() int {
	return r.selected
}

// SetSelectedIndex sets the selected index.
func (r *RadioGroup) SetSelectedIndex(i int) {
	if i >= 0 && i < len(r.options) {
		r.selected = i
	}
}

// Focus focuses the radio group.
func (r *RadioGroup) Focus() render.Cmd {
	r.focused = true
	return nil
}

// Blur blurs the radio group.
func (r *RadioGroup) Blur() {
	r.focused = false
}

// Focused returns true if focused.
func (r *RadioGroup) Focused() bool {
	return r.focused
}

// Validate validates the input.
func (r *RadioGroup) Validate() error {
	for _, v := range r.validators {
		if err := v(r.Value()); err != nil {
			r.err = err
			return err
		}
	}
	r.err = nil
	return nil
}

// Error returns the last validation error.
func (r *RadioGroup) Error() error {
	return r.err
}

// AddValidator adds a validator.
func (r *RadioGroup) AddValidator(v Validator) {
	r.validators = append(r.validators, v)
}
