// Buffer Forms Example - Form layout using buffer system
//
// This example demonstrates:
// - Form field layout with labels and inputs
// - Validation indicators
// - Button layout
// - Multi-section forms
//
// Usage: go run main.go

package main

import (
	"fmt"

	"github.com/wwsheng009/taproot/ui/render/buffer"
)

// FormField represents a single form field
type FormField struct {
	label    string
	value    string
	placeholder string
	required bool
	valid    bool
	errorMsg string
}

func (f *FormField) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw label
	labelText := f.label
	if f.required {
		labelText += " *"
	}
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, labelText, buffer.Style{Foreground: "#86", Bold: true})

	// Draw input field box
	inputY := rect.Y + 1
	boxWidth := rect.Width

	// Top border
	buf.WriteString(buffer.Point{X: rect.X, Y: inputY}, "┌", buffer.Style{Foreground: "#244"})
	for i := 1; i < boxWidth-1; i++ {
		buf.WriteString(buffer.Point{X: rect.X + i, Y: inputY}, "─", buffer.Style{Foreground: "#244"})
	}
	buf.WriteString(buffer.Point{X: rect.X + boxWidth - 1, Y: inputY}, "┐", buffer.Style{Foreground: "#244"})

	// Middle with value or placeholder
	value := f.value
	if value == "" {
		value = f.placeholder
	}
	if len(value) > boxWidth-4 {
		value = value[:boxWidth-4] + ".."
	}
	valueColor := "#250"
	if f.value == "" {
		valueColor = "#244" // Dim for placeholder
	}

	buf.WriteString(buffer.Point{X: rect.X, Y: inputY + 1}, "│", buffer.Style{Foreground: "#244"})
	buf.WriteString(buffer.Point{X: rect.X + 1, Y: inputY + 1}, value, buffer.Style{Foreground: valueColor})
	buf.WriteString(buffer.Point{X: rect.X + boxWidth - 1, Y: inputY + 1}, "│", buffer.Style{Foreground: "#244"})

	// Bottom border
	buf.WriteString(buffer.Point{X: rect.X, Y: inputY + 2}, "└", buffer.Style{Foreground: "#244"})
	for i := 1; i < boxWidth-1; i++ {
		buf.WriteString(buffer.Point{X: rect.X + i, Y: inputY + 2}, "─", buffer.Style{Foreground: "#244"})
	}
	buf.WriteString(buffer.Point{X: rect.X + boxWidth - 1, Y: inputY + 2}, "┘", buffer.Style{Foreground: "#244"})

	// Draw validation indicator
	indicatorX := rect.X + boxWidth
	if f.valid {
		buf.WriteString(buffer.Point{X: indicatorX, Y: inputY + 1}, "✓", buffer.Style{Foreground: "#120", Bold: true})
	} else if f.errorMsg != "" {
		buf.WriteString(buffer.Point{X: indicatorX, Y: inputY + 1}, "✗", buffer.Style{Foreground: "#160", Bold: true})
	}

	// Draw error message
	if f.errorMsg != "" {
		buf.WriteString(buffer.Point{X: rect.X, Y: inputY + 3}, f.errorMsg, buffer.Style{Foreground: "#160"})
	}
}

func (f *FormField) MinSize() (int, int)     { return 30, 2 }
func (f *FormField) PreferredSize() (int, int) { return 40, 4 }

// Button represents a clickable button
type Button struct {
	label    string
	primary  bool
	disabled bool
}

func (b *Button) Render(buf *buffer.Buffer, rect buffer.Rect) {
	label := b.label

	// Draw button
	if b.primary {
		// Primary button - filled background
		for x := rect.X; x < rect.X+rect.Width; x++ {
			buf.WriteString(buffer.Point{X: x, Y: rect.Y}, " ", buffer.Style{Background: "#32", Foreground: "#15"})
		}
		labelX := rect.X + (rect.Width-len(label))/2
		buf.WriteString(buffer.Point{X: labelX, Y: rect.Y}, label, buffer.Style{Background: "#32", Foreground: "#15", Bold: true})
	} else if b.disabled {
		// Disabled button
		for x := rect.X; x < rect.X+rect.Width; x++ {
			buf.WriteString(buffer.Point{X: x, Y: rect.Y}, " ", buffer.Style{Background: "#236"})
		}
		labelX := rect.X + (rect.Width-len(label))/2
		buf.WriteString(buffer.Point{X: labelX, Y: rect.Y}, label, buffer.Style{Background: "#236", Foreground: "#244"})
	} else {
		// Secondary button - outline
		buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, "┌", buffer.Style{Foreground: "#32"})
		for i := 1; i < rect.Width-1; i++ {
			buf.WriteString(buffer.Point{X: rect.X + i, Y: rect.Y}, "─", buffer.Style{Foreground: "#32"})
		}
		buf.WriteString(buffer.Point{X: rect.X + rect.Width - 1, Y: rect.Y}, "┐", buffer.Style{Foreground: "#32"})

		labelX := rect.X + (rect.Width-len(label))/2
		buf.WriteString(buffer.Point{X: labelX, Y: rect.Y}, label, buffer.Style{Foreground: "#32"})
	}
}

func (b *Button) MinSize() (int, int)     { return len(b.label) + 4, 1 }
func (b *Button) PreferredSize() (int, int) { return len(b.label) + 8, 1 }

// Section represents a form section with title
type Section struct {
	title  string
	fields []*FormField
}

func (s *Section) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw section title
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, "═ "+s.title+" ", buffer.Style{Foreground: "#86", Bold: true})
	for i := len(s.title) + 3; i < rect.Width; i++ {
		buf.WriteString(buffer.Point{X: rect.X + i, Y: rect.Y}, "═", buffer.Style{Foreground: "#244"})
	}

	// Render fields
	fieldY := rect.Y + 2
	fieldWidth := rect.Width - 10 // Leave space for validation indicator
	for _, field := range s.fields {
		field.Render(buf, buffer.Rect{X: rect.X, Y: fieldY, Width: fieldWidth, Height: 4})
		fieldY += 5 // Each field takes 4 lines + 1 line spacing
	}
}

func (s *Section) MinSize() (int, int)     { return 40, 10 }
func (s *Section) PreferredSize() (int, int) { return 60, 20 }

// Checkbox component
type Checkbox struct {
	label    string
	checked  bool
}

func (c *Checkbox) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw checkbox
	if c.checked {
		buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, "[✓]", buffer.Style{Foreground: "#120", Bold: true})
	} else {
		buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, "[ ]", buffer.Style{Foreground: "#244"})
	}
	// Draw label
	buf.WriteString(buffer.Point{X: rect.X + 5, Y: rect.Y}, c.label, buffer.Style{Foreground: "#250"})
}

func (c *Checkbox) MinSize() (int, int)     { return len(c.label) + 6, 1 }
func (c *Checkbox) PreferredSize() (int, int) { return len(c.label) + 10, 1 }

// RadioGroup component
type RadioGroup struct {
	label   string
	options []string
	selected int
}

func (r *RadioGroup) Render(buf *buffer.Buffer, rect buffer.Rect) {
	// Draw label
	buf.WriteString(buffer.Point{X: rect.X, Y: rect.Y}, r.label+":", buffer.Style{Foreground: "#86", Bold: true})

	// Draw options
	y := rect.Y + 2
	for i, option := range r.options {
		if y >= rect.Y+rect.Height {
			break
		}
		if i == r.selected {
			buf.WriteString(buffer.Point{X: rect.X + 2, Y: y}, "●", buffer.Style{Foreground: "#32", Bold: true})
		} else {
			buf.WriteString(buffer.Point{X: rect.X + 2, Y: y}, "○", buffer.Style{Foreground: "#244"})
		}
		buf.WriteString(buffer.Point{X: rect.X + 5, Y: y}, option, buffer.Style{Foreground: "#250"})
		y++
	}
}

func (r *RadioGroup) MinSize() (int, int)     { return 20, len(r.options) + 2 }
func (r *RadioGroup) PreferredSize() (int, int) { return 30, len(r.options) + 4 }

func main() {
	width := 80
	height := 35

	fmt.Println(repeat("=", width))
	fmt.Println("Buffer Forms Example - Taproot")
	fmt.Println(repeat("=", width))
	fmt.Println()

	// Create buffer
	buf := buffer.NewBuffer(width, height)

	// Title
	title := " User Registration Form "
	buf.WriteString(buffer.Point{X: (width-len(title))/2, Y: 0}, title, buffer.Style{Foreground: "#86", Background: "#235", Bold: true})
	buf.WriteString(buffer.Point{X: 0, Y: 1}, repeat("─", width), buffer.Style{Foreground: "#244"})

	// Create form sections

	// Personal Information Section
	personalSection := &Section{
		title: "Personal Information",
		fields: []*FormField{
			{
				label:       "Full Name",
				value:       "John Doe",
				placeholder: "Enter your full name",
				required:    true,
				valid:       true,
			},
			{
				label:       "Email Address",
				value:       "john@example.com",
				placeholder: "your@email.com",
				required:    true,
				valid:       true,
			},
			{
				label:       "Phone Number",
				value:       "",
				placeholder: "+1 (555) 123-4567",
				required:    false,
				valid:       true,
			},
		},
	}

	// Account Details Section
	accountSection := &Section{
		title: "Account Details",
		fields: []*FormField{
			{
				label:       "Username",
				value:       "johndoe",
				placeholder: "Choose a username",
				required:    true,
				valid:       true,
			},
			{
				label:       "Password",
				value:       "",
				placeholder: "••••••••",
				required:    true,
				valid:       false,
				errorMsg:    "Password must be at least 8 characters",
			},
		},
	}

	// Render sections side by side
	sectionWidth := width/2 - 4
	personalSection.Render(buf, buffer.Rect{X: 2, Y: 3, Width: sectionWidth, Height: 18})
	accountSection.Render(buf, buffer.Rect{X: sectionWidth + 4, Y: 3, Width: sectionWidth, Height: 18})

	// Preferences section (bottom)
	prefsY := 23
	buf.WriteString(buffer.Point{X: 2, Y: prefsY}, "═ Preferences ", buffer.Style{Foreground: "#86", Bold: true})
	for i := 15; i < width-2; i++ {
		buf.WriteString(buffer.Point{X: i, Y: prefsY}, "═", buffer.Style{Foreground: "#244"})
	}

	// Checkboxes
	checkY := prefsY + 2
	check1 := &Checkbox{label: "Subscribe to newsletter", checked: true}
	check1.Render(buf, buffer.Rect{X: 4, Y: checkY, Width: 30, Height: 1})
	check2 := &Checkbox{label: "Accept Terms of Service", checked: true}
	check2.Render(buf, buffer.Rect{X: 4, Y: checkY + 1, Width: 30, Height: 1})
	check3 := &Checkbox{label: "Remember me", checked: false}
	check3.Render(buf, buffer.Rect{X: 4, Y: checkY + 2, Width: 30, Height: 1})

	// Radio group for user type
	radioGroup := &RadioGroup{
		label:    "Account Type",
		options:  []string{"Personal", "Professional", "Enterprise"},
		selected: 1,
	}
	radioGroup.Render(buf, buffer.Rect{X: 45, Y: checkY, Width: 30, Height: 6})

	// Buttons at bottom
	buttonY := height - 3
	btnSubmit := &Button{label: "Submit", primary: true}
	btnSubmit.Render(buf, buffer.Rect{X: width - 30, Y: buttonY, Width: 12, Height: 1})
	btnCancel := &Button{label: "Cancel", primary: false, disabled: false}
	btnCancel.Render(buf, buffer.Rect{X: width - 15, Y: buttonY, Width: 12, Height: 1})

	// Footer hint
	buf.WriteString(buffer.Point{X: 2, Y: height - 1}, "* Required fields | Tab: Navigate | Enter: Submit", buffer.Style{Foreground: "#244"})

	// Render output
	output := buf.Render()
	fmt.Print(output)

	// Print summary
	fmt.Printf("\n\nForm Components Demonstrated:\n")
	fmt.Printf("  • FormField: Text input with label, placeholder, validation\n")
	fmt.Printf("  • Button: Primary, secondary, and disabled states\n")
	fmt.Printf("  • Section: Grouped form fields with title\n")
	fmt.Printf("  • Checkbox: Toggle selection\n")
	fmt.Printf("  • RadioGroup: Single selection from options\n")
	fmt.Printf("\nLayout: Two-column sections + preferences + buttons\n")
	fmt.Printf("Total: %d components in %d×%d grid\n", 12, width, height)
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
