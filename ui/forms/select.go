package forms

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// Select represents a dropdown selection component.
type Select struct {
	label       string
	options     []string
	placeholder string
	selected    int // Index of the confirmed selection (-1 if none)
	highlighted int // Index of the currently highlighted option in the dropdown
	expanded    bool
	focused     bool
	width       int

	// Styles
	focusedStyle      lipgloss.Style
	blurredStyle      lipgloss.Style
	labelStyle        lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	errorStyle        lipgloss.Style

	// Validation
	validators []Validator
	err        error
}

// NewSelect creates a new select component.
func NewSelect(label string, options []string) *Select {
	s := styles.DefaultStyles()
	return &Select{
		label:       label,
		options:     options,
		selected:    -1, // No selection by default
		highlighted: 0,
		width:       40,
		focusedStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(s.BorderColor).
			Padding(0, 1),
		blurredStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(s.Border).
			Padding(0, 1),
		labelStyle: lipgloss.NewStyle().
			Foreground(s.FgBase).
			Bold(true),
		itemStyle: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(s.FgMuted),
		selectedItemStyle: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(s.Primary).
			Bold(true),
		errorStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
	}
}

// Init implements render.Model.
func (s *Select) Init() error {
	return nil
}

// Update implements render.Model.
func (s *Select) Update(msg any) (render.Model, render.Cmd) {
	if !s.focused {
		return s, nil
	}

	var keyStr string
	if k, ok := msg.(tea.KeyMsg); ok {
		keyStr = k.String()
	} else if k, ok := msg.(render.KeyMsg); ok {
		keyStr = k.String()
	}

	switch keyStr {
	case "enter", " ":
		if s.expanded {
			// Confirm selection
			s.selected = s.highlighted
			s.expanded = false
		} else {
			// Expand
			s.expanded = true
			// If we have a selection, highlight it
			s.highlighted = max(0, s.selected)

		}
	case "esc":
		if s.expanded {
			s.expanded = false
		}
	case "up", "k":
		if s.expanded {
			if s.highlighted > 0 {
				s.highlighted--
			} else if len(s.options) > 0 {
				s.highlighted = len(s.options) - 1 // Wrap
			}
		}
	case "down", "j":
		if s.expanded {
			if s.highlighted < len(s.options)-1 {
				s.highlighted++
			} else if len(s.options) > 0 {
				s.highlighted = 0 // Wrap
			}
		}
	}

	return s, nil
}

// View implements render.Model.
func (s *Select) View() string {
	var b strings.Builder

	// Render Label
	if s.label != "" {
		b.WriteString(s.labelStyle.Render(s.label))
		b.WriteString("\n")
	}

	// Render Input Box
	var displayText string
	if s.selected >= 0 && s.selected < len(s.options) {
		displayText = s.options[s.selected]
	} else if s.placeholder != "" {
		displayText = s.placeholder
	}

	// Add dropdown indicator
	// Calculate padding to align the arrow to the right
	// We want the text on left, arrow on right, total width s.width
	// We need to account for visual width.
	// For simplicity, let's just use space padding.
	
	// Truncate text if too long (leave room for " v")
	availWidth := s.width - 2 
	if len(displayText) > availWidth {
		displayText = displayText[:availWidth]
	}
	
	padding := s.width - len(displayText) - 2 // -2 for " v"
	padding = max(0, padding)

	
	displayContent := displayText + strings.Repeat(" ", padding) + " ▼"
	
	var box string
	if s.focused {
		box = s.focusedStyle.Render(displayContent)
	} else {
		box = s.blurredStyle.Render(displayContent)
	}
	b.WriteString(box)

	// Render Dropdown List if expanded
	if s.expanded && len(s.options) > 0 {
		b.WriteString("\n")
		// Draw a box around options? Or just list them?
		// For now, just list them with indentation/box
		// Using a simple list style
		
		for i, opt := range s.options {
			var item string
			prefix := "  "
			if i == s.highlighted {
				prefix = "> "
				item = s.selectedItemStyle.Render(prefix + opt)
			} else {
				item = s.itemStyle.Render(prefix + opt)
			}
			b.WriteString(item)
			if i < len(s.options)-1 {
				b.WriteString("\n")
			}
		}
	}

	// Render error
	if s.err != nil {
		b.WriteString("\n")
		b.WriteString(s.errorStyle.Render(s.err.Error()))
	}

	return b.String()
}

// Value returns the selected option string.
func (s *Select) Value() string {
	if s.selected >= 0 && s.selected < len(s.options) {
		return s.options[s.selected]
	}
	return ""
}

// SetValue sets the selected option by value string.
func (s *Select) SetValue(v string) {
	for i, opt := range s.options {
		if opt == v {
			s.selected = i
			return
		}
	}
	s.selected = -1
}

// SelectedIndex returns the index of the selected option.
func (s *Select) SelectedIndex() int {
	return s.selected
}

// SetSelectedIndex sets the selected index.
func (s *Select) SetSelectedIndex(i int) {
	if i >= -1 && i < len(s.options) {
		s.selected = i
	}
}

// SetPlaceholder sets the placeholder text.
func (s *Select) SetPlaceholder(p string) {
	s.placeholder = p
}

// Focus focuses the component.
func (s *Select) Focus() render.Cmd {
	s.focused = true
	return nil
}

// Blur blurs the component.
func (s *Select) Blur() {
	s.focused = false
	s.expanded = false // Auto-collapse on blur
}

// Focused returns true if focused.
func (s *Select) Focused() bool {
	return s.focused
}

// Validate validates the input.
func (s *Select) Validate() error {
	for _, v := range s.validators {
		if err := v(s.Value()); err != nil {
			s.err = err
			return err
		}
	}
	s.err = nil
	return nil
}

// Error returns the last validation error.
func (s *Select) Error() error {
	return s.err
}

// AddValidator adds a validator.
func (s *Select) AddValidator(v Validator) {
	s.validators = append(s.validators, v)
}

// SetWidth sets the width of the input box.
func (s *Select) SetWidth(w int) {
	s.width = w
}
