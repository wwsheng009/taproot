package forms

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// Checkbox represents a boolean input field.
type Checkbox struct {
	label   string
	checked bool
	focused bool

	// Styles
	focusedStyle lipgloss.Style
	blurredStyle lipgloss.Style
	checkStyle   lipgloss.Style
	labelStyle   lipgloss.Style

	// Validation (though rarely used for single checkboxes unless "required")
	validators []Validator
	err        error
	errorStyle lipgloss.Style
}

// NewCheckbox creates a new checkbox.
func NewCheckbox(label string) *Checkbox {
	s := styles.DefaultStyles()
	return &Checkbox{
		label: label,
		focusedStyle: lipgloss.NewStyle().
			Foreground(s.Primary),
		blurredStyle: lipgloss.NewStyle().
			Foreground(s.FgMuted),
		checkStyle: lipgloss.NewStyle().
			Foreground(s.Primary).
			Bold(true),
		labelStyle: lipgloss.NewStyle(), // Inherit default
		errorStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
	}
}

// Init implements render.Model.
func (c *Checkbox) Init() render.Cmd {
	return nil
}

// Update implements render.Model.
func (c *Checkbox) Update(msg any) (render.Model, render.Cmd) {
	if !c.focused {
		return c, nil
	}

	var keyStr string
	if k, ok := msg.(tea.KeyMsg); ok {
		keyStr = k.String()
	} else if k, ok := msg.(render.KeyMsg); ok {
		keyStr = k.String()
	}

	switch keyStr {
	case " ", "enter":
		c.checked = !c.checked
		// If "required" validation exists, re-validate immediately
		c.Validate()
	}

	return c, nil
}

// View implements render.Model.
func (c *Checkbox) View() string {
	var b strings.Builder

	var checkMark string
	if c.checked {
		checkMark = "[x]"
	} else {
		checkMark = "[ ]"
	}

	// Apply styles
	if c.focused {
		checkMark = c.focusedStyle.Render(checkMark)
		b.WriteString(checkMark)
		b.WriteString(" ")
		b.WriteString(c.focusedStyle.Render(c.label))
	} else {
		checkMark = c.blurredStyle.Render(checkMark)
		b.WriteString(checkMark)
		b.WriteString(" ")
		b.WriteString(c.blurredStyle.Render(c.label))
	}

	// Render error if present
	if c.err != nil {
		b.WriteString("\n")
		b.WriteString(c.errorStyle.Render(c.err.Error()))
	}

	return b.String()
}

// Value returns "true" or "false".
func (c *Checkbox) Value() string {
	return fmt.Sprintf("%v", c.checked)
}

// SetValue sets the value ("true" or "false").
func (c *Checkbox) SetValue(v string) {
	c.checked = v == "true"
}

// Checked returns the boolean state.
func (c *Checkbox) Checked() bool {
	return c.checked
}

// SetChecked sets the boolean state.
func (c *Checkbox) SetChecked(checked bool) {
	c.checked = checked
}

// Focus focuses the checkbox.
func (c *Checkbox) Focus() render.Cmd {
	c.focused = true
	return nil
}

// Blur blurs the checkbox.
func (c *Checkbox) Blur() {
	c.focused = false
}

// Focused returns true if focused.
func (c *Checkbox) Focused() bool {
	return c.focused
}

// Validate validates the input.
func (c *Checkbox) Validate() error {
	for _, v := range c.validators {
		if err := v(c.Value()); err != nil {
			c.err = err
			return err
		}
	}
	c.err = nil
	return nil
}

// Error returns the last validation error.
func (c *Checkbox) Error() error {
	return c.err
}

// AddValidator adds a validator.
func (c *Checkbox) AddValidator(v Validator) {
	c.validators = append(c.validators, v)
}
