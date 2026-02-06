// Buffer Forms Example - Interactive form using Bubbletea
//
// This example demonstrates:
// - Form fields with keyboard navigation
// - Input validation
// - Checkbox toggles
// - Radio selection
// - Submit/Cancel buttons
//
// Usage: go run main.go
// Keys: Tab/Shift+Tab: Navigate | Space: Toggle checkbox/radio | Enter: Submit

package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/render/buffer"
)

// FieldType represents different field types
type FieldType int

const (
	FieldInput FieldType = iota
	FieldCheckbox
	FieldRadio
	FieldButton
)

// FormField represents a form field
type FormField struct {
	id       string
	label    string
	fieldType FieldType
	value    string
	options  []string // For radio buttons
	checked  bool     // For checkbox
	selected int      // For radio selection
	focused  bool
	valid    bool
	errorMsg string
	primary  bool // For buttons
}

// FormModel represents the form state
type FormModel struct {
	width    int
	height   int
	quitting bool
	fields   []*FormField
	focused  int
	submitted bool
}

func NewFormModel() FormModel {
	return FormModel{
		width:  80,
		height: 30,
		fields: []*FormField{
			{
				id:        "name",
				label:     "Full Name *",
				fieldType: FieldInput,
				value:     "",
				valid:     false,
				errorMsg:  "Name is required",
			},
			{
				id:        "email",
				label:     "Email *",
				fieldType: FieldInput,
				value:     "",
				valid:     false,
				errorMsg:  "Email is required",
			},
			{
				id:        "phone",
				label:     "Phone",
				fieldType: FieldInput,
				value:     "",
				valid:     true,
			},
			{
				id:        "newsletter",
				label:     "Subscribe to newsletter",
				fieldType: FieldCheckbox,
				checked:   false,
			},
			{
				id:        "terms",
				label:     "Accept Terms of Service *",
				fieldType: FieldCheckbox,
				checked:   false,
				valid:     false,
				errorMsg:  "Must accept terms",
			},
			{
				id:        "accountType",
				label:     "Account Type",
				fieldType: FieldRadio,
				options:   []string{"Personal", "Professional", "Enterprise"},
				selected:  0,
			},
			{
				id:        "submit",
				label:     "Submit",
				fieldType: FieldButton,
				primary:   true,
			},
			{
				id:        "cancel",
				label:     "Cancel",
				fieldType: FieldButton,
				primary:   false,
			},
		},
		focused: 0,
	}
}

func (m FormModel) Init() tea.Cmd {
	return nil
}

type TickMsg struct{}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		if m.submitted {
			return m, tea.Quit
		}
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyStr := msg.String()
		keyRunes := msg.Runes

		switch keyStr {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.focused = (m.focused + 1) % len(m.fields)
			return m, nil
		case "shift+tab":
			m.focused = (m.focused - 1 + len(m.fields)) % len(m.fields)
			return m, nil
		case "up", "k":
			m.focused = (m.focused - 1 + len(m.fields)) % len(m.fields)
			return m, nil
		case "down", "j":
			m.focused = (m.focused + 1) % len(m.fields)
			return m, nil
		case "enter", " ":
			field := m.fields[m.focused]
			if field.fieldType == FieldButton {
				if field.id == "submit" {
					if m.validate() {
						m.submitted = true
						m.quitting = true
						return m, tea.Quit
					}
				} else {
					m.quitting = true
					return m, tea.Quit
				}
			} else if field.fieldType == FieldCheckbox {
				m.fields[m.focused].checked = !m.fields[m.focused].checked
			}
			return m, nil
		case "backspace":
			field := m.fields[m.focused]
			if field.fieldType == FieldInput && len(field.value) > 0 {
				m.fields[m.focused].value = field.value[:len(field.value)-1]
				m.fields[m.focused].valid = len(m.fields[m.focused].value) > 0
			}
			return m, nil
		default:
			// Handle character input for text fields
			field := m.fields[m.focused]
			if field.fieldType == FieldInput && len(keyRunes) == 1 {
				// Only allow printable characters
				r := keyRunes[0]
				if r >= 32 && r <= 126 {
					m.fields[m.focused].value += string(r)
					m.fields[m.focused].valid = len(m.fields[m.focused].value) > 0
					if len(m.fields[m.focused].value) > 0 {
						m.fields[m.focused].errorMsg = ""
					}
				}
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m FormModel) validate() bool {
	valid := true
	for _, field := range m.fields {
		field.focused = false
		if field.id == "name" || field.id == "email" {
			if len(field.value) == 0 {
				field.valid = false
				valid = false
			}
		}
		if field.id == "terms" && !field.checked {
			field.valid = false
			valid = false
		}
	}
	return valid
}

func (m FormModel) View() string {
	if m.submitted {
		return "\n\n\n" + strings.Repeat(" ", 35) + "✓ Form Submitted Successfully!\n\n"
	}

	if m.quitting {
		return "\n\n\n" + strings.Repeat(" ", 40) + "Cancelled\n\n"
	}

	buf := buffer.NewBuffer(m.width, m.height)

	// Title
	title := " User Registration "
	titleX := (m.width - len(title)) / 2
	buf.WriteString(buffer.Point{X: titleX, Y: 0}, title, buffer.Style{Foreground: "#86", Background: "#235", Bold: true})

	// Subtitle
	subtitle := "Fill in the form below (Fields marked with * are required)"
	subtitleX := (m.width - len(subtitle)) / 2
	buf.WriteString(buffer.Point{X: subtitleX, Y: 2}, subtitle, buffer.Style{Foreground: "#244"})

	// Draw fields
	y := 5
	buttonY := 0
	for i, field := range m.fields {
		field.focused = (i == m.focused)
		if field.fieldType == FieldButton {
			buttonY = y // Store button Y position
		} else {
			renderField(buf, field, i, y, m.width)
			if field.fieldType == FieldInput {
				y += 5
			} else if field.fieldType == FieldCheckbox {
				y += 2
			} else if field.fieldType == FieldRadio {
				y += len(field.options) + 2
			}
		}
	}

	// Draw buttons together on same row
	buttonWidth := 15
	submitX := m.width/2 - buttonWidth - 2
	cancelX := m.width/2 + 2

	// Find button fields
	for i, field := range m.fields {
		if field.fieldType == FieldButton {
			field.focused = (i == m.focused)
			if field.id == "submit" {
				renderButton(buf, field.label, submitX, buttonY, buttonWidth, field.primary, field.focused)
			} else if field.id == "cancel" {
				renderButton(buf, field.label, cancelX, buttonY, buttonWidth, field.primary, field.focused)
			}
		}
	}

	// Footer
	footer := " Tab: Navigate | Type: Input | Backspace: Delete | Space/Enter: Select | Ctrl+C: Cancel "
	buf.WriteString(buffer.Point{X: 0, Y: m.height - 1}, repeat(" ", m.width), buffer.Style{Background: "#235"})
	buf.WriteString(buffer.Point{X: 2, Y: m.height - 1}, footer, buffer.Style{Background: "#235", Foreground: "#250"})

	return buf.Render()
}

func renderField(buf *buffer.Buffer, field *FormField, index, y, width int) {
	x := 10
	fieldWidth := width - 20

	if field.fieldType == FieldInput {
		// Label
		label := field.label
		buf.WriteString(buffer.Point{X: x, Y: y}, label, buffer.Style{
			Foreground: "#86",
			Bold:       true,
		})

		// Error message
		if !field.valid && field.errorMsg != "" && !field.focused {
			buf.WriteString(buffer.Point{X: x + fieldWidth, Y: y}, " ✗", buffer.Style{Foreground: "#160", Bold: true})
		} else if field.valid && len(field.value) > 0 {
			buf.WriteString(buffer.Point{X: x + fieldWidth, Y: y}, " ✓", buffer.Style{Foreground: "#120", Bold: true})
		}

		// Input box
		inputY := y + 2
		buf.WriteString(buffer.Point{X: x, Y: inputY}, "┌", buffer.Style{Foreground: field.borderColor()})
		for i := 1; i < fieldWidth-1; i++ {
			buf.WriteString(buffer.Point{X: x + i, Y: inputY}, "─", buffer.Style{Foreground: field.borderColor()})
		}
		buf.WriteString(buffer.Point{X: x + fieldWidth - 1, Y: inputY}, "┐", buffer.Style{Foreground: field.borderColor()})

		// Input content
		content := field.value
		if content == "" && !field.focused {
			content = "Type here..."
		}
		if len(content) > fieldWidth-4 {
			content = content[:fieldWidth-4] + ".."
		}
		contentColor := "#250"
		if field.value == "" {
			contentColor = "#244"
		}
		if field.focused {
			contentColor = "#15"
			content = field.value + "_"
		}

		buf.WriteString(buffer.Point{X: x, Y: inputY + 1}, "│", buffer.Style{Foreground: field.borderColor()})
		buf.WriteString(buffer.Point{X: x + 1, Y: inputY + 1}, content, buffer.Style{Foreground: contentColor})
		buf.WriteString(buffer.Point{X: x + fieldWidth - 1, Y: inputY + 1}, "│", buffer.Style{Foreground: field.borderColor()})

		// Bottom border
		buf.WriteString(buffer.Point{X: x, Y: inputY + 2}, "└", buffer.Style{Foreground: field.borderColor()})
		for i := 1; i < fieldWidth-1; i++ {
			buf.WriteString(buffer.Point{X: x + i, Y: inputY + 2}, "─", buffer.Style{Foreground: field.borderColor()})
		}
		buf.WriteString(buffer.Point{X: x + fieldWidth - 1, Y: inputY + 2}, "┘", buffer.Style{Foreground: field.borderColor()})

		// Error message below input
		if !field.valid && field.errorMsg != "" {
			buf.WriteString(buffer.Point{X: x, Y: inputY + 3}, field.errorMsg, buffer.Style{Foreground: "#160"})
		}

	} else if field.fieldType == FieldCheckbox {
		// Checkbox
		indicator := "[ ]"
		if field.checked {
			indicator = "[✓]"
		}
		color := "#250"
		if field.focused {
			color = "#32"
		}
		if field.checked {
			color = "#120"
		}

		buf.WriteString(buffer.Point{X: x, Y: y}, indicator, buffer.Style{Foreground: color, Bold: field.focused})
		buf.WriteString(buffer.Point{X: x + 5, Y: y}, field.label, buffer.Style{Foreground: "#250", Bold: field.focused})

		// Required indicator
		if strings.Contains(field.label, "*") {
			buf.WriteString(buffer.Point{X: x + len(field.label) + 6, Y: y}, " ✗", buffer.Style{Foreground: "#160"})
		} else if field.checked {
			buf.WriteString(buffer.Point{X: x + len(field.label) + 6, Y: y}, " ✓", buffer.Style{Foreground: "#120"})
		}

	} else if field.fieldType == FieldRadio {
		// Label
		buf.WriteString(buffer.Point{X: x, Y: y}, field.label+":", buffer.Style{Foreground: "#86", Bold: field.focused})

		// Options
		for i, option := range field.options {
			optY := y + 2 + i
			indicator := "○"
			if field.selected == i {
				indicator = "●"
			}
			color := "#250"
			if field.focused && field.selected == i {
				color = "#32"
			} else if field.selected == i {
				color = "#120"
			}

			buf.WriteString(buffer.Point{X: x + 2, Y: optY}, indicator, buffer.Style{Foreground: color, Bold: field.selected == i})
			buf.WriteString(buffer.Point{X: x + 5, Y: optY}, option, buffer.Style{Foreground: color})
		}
	}
	// Buttons are handled separately in View()
}

func renderButton(buf *buffer.Buffer, label string, x, y, width int, primary, focused bool) {
	if primary {
		// Primary button - filled
		for i := 0; i < width; i++ {
			color := "#32"
			if focused {
				color = "#50"
			}
			buf.WriteString(buffer.Point{X: x + i, Y: y}, " ", buffer.Style{Background: color})
		}
		labelX := x + (width-len(label))/2
		labelColor := "#15"
		if focused {
			labelColor = "#230"
		}
		buf.WriteString(buffer.Point{X: labelX, Y: y}, label, buffer.Style{Background: "#32", Foreground: labelColor, Bold: true})
	} else {
		// Secondary button - outline
		buf.WriteString(buffer.Point{X: x, Y: y}, "┌", buffer.Style{Foreground: "#32"})
		for i := 1; i < width-1; i++ {
			buf.WriteString(buffer.Point{X: x + i, Y: y}, "─", buffer.Style{Foreground: "#32"})
		}
		buf.WriteString(buffer.Point{X: x + width - 1, Y: y}, "┐", buffer.Style{Foreground: "#32"})

		labelX := x + (width-len(label))/2
		color := "#32"
		if focused {
			color = "#50"
		}
		buf.WriteString(buffer.Point{X: labelX, Y: y}, label, buffer.Style{Foreground: color})
	}
}

func (f *FormField) borderColor() string {
	if f.focused {
		return "#32" // Green when focused
	}
	if !f.valid && f.errorMsg != "" {
		return "#160" // Red when invalid
	}
	return "#244" // Gray default
}

func repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	result := s
	for i := 1; i < count; i++ {
		result += s
	}
	return result
}

func main() {
	p := tea.NewProgram(NewFormModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
