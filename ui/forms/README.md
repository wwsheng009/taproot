# Forms Package

The `forms` package provides reusable TUI form components for the Taproot framework. These components support validation, focus management, and are engine-agnostic via the `render` package.

## Components

### TextInput

A single-line text input field with cursor management, input masking (for passwords), and validation support.

```go
import "github.com/wwsheng009/taproot/ui/forms"

input := forms.NewTextInput("Enter email")
input.AddValidator(forms.Email)
input.AddValidator(forms.Required)
input.Focus()
```

**Key Features:**
- Cursor movement (left/right, home/end)
- Character insertion and deletion
- Password masking with `SetHidden(true)`
- Max length limitation with `SetMaxLength()`
- Custom prompt labels with `SetPrompt()`

### NumberInput

A specialized input for numeric values with increment/decrement controls.

```go
ageInput := forms.NewNumberInput("Enter age")
ageInput.SetRange(0, 150)
ageInput.SetStep(1)
ageInput.SetPrecision(0) // Integer
```

**Key Features:**
- Range validation (min/max)
- Precision control (decimal places)
- Increment/decrement with Up/Down arrow keys
- Automatic clamping to range limits

### TextArea

A multi-line text input with scrolling support.

```go
bioArea := forms.NewTextArea("Enter bio")
bioArea.SetWidth(60)
bioArea.SetHeight(10)
bioArea.AddValidator(forms.MinLength(10))
```

**Key Features:**
- Multi-line editing with Enter key
- Line merging on backspace
- Vertical scrolling (viewport offset)
- Newline insertion and line navigation

## Validators

The `validator.go` file provides standard validators:

```go
// Built-in validators
forms.Required          // Field must not be empty
forms.MaxLength(100)    // Maximum length
forms.MinLength(8)      // Minimum length
forms.Email             // Email format validation
forms.Range(0, 100)     // Numeric range

// Custom validator
customValidator := func(value string) error {
    if invalid {
        return errors.New("custom error message")
    }
    return nil
}
input.AddValidator(customValidator)
```

### Using Validators

```go
// Single validator
input.AddValidator(forms.Required)

// Multiple validators
input.AddValidator(forms.Email)
input.AddValidator(forms.Required)

// Replace all validators
input.SetValidators(forms.Email, forms.Required)

// Validate and get error
if err := input.Validate(); err != nil {
    // Handle validation error
}
```

## Focus Management

All form components implement focus management:

```go
// Focus a component
input.Focus()

// Blur a component
input.Blur()

// Check focus state
if input.Focused() {
    // Component is focused
}
```

## Engine Integration

Form components are engine-agnostic and support both `tea.KeyMsg` (Bubbletea) and `render.KeyMsg`. They automatically detect and handle both message types.

```go
import tea "github.com/charmbracelet/bubbletea"

type model struct {
    input *forms.TextInput
}

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
    // Just pass tea.Msg directly - forms components handle it
    newInput, cmd := m.input.Update(msg)
    m.input = newInput.(*forms.TextInput)
    
    // Convert render.Cmd back to tea.Cmd
    if cmd != nil {
        return m, adaptCmd(cmd)
    }
    return m, nil
}
```

## Example

See `examples/forms/main.go` for a complete form with:
- Email input with validation
- Password input with masking
- Age number input with range
- Bio textarea

```bash
go run examples/forms/main.go
```

## API Reference

### TextInput

```go
func NewTextInput(placeholder string) *TextInput
func (t *TextInput) Value() string
func (t *TextInput) SetValue(v string)
func (t *TextInput) Focus() render.Cmd
func (t *TextInput) Blur()
func (t *TextInput) Focused() bool
func (t *TextInput) Validate() error
func (t *TextInput) AddValidator(v Validator)
func (t *TextInput) SetValidators(v ...Validator)
func (t *TextInput) SetPrompt(p string)
func (t *TextInput) SetHidden(h bool)
func (t *TextInput) SetWidth(w int)
func (t *TextInput) SetMaxLength(l int)
```

### NumberInput

```go
func NewNumberInput(placeholder string) *NumberInput
func (n *NumberInput) FloatValue() float64
func (n *NumberInput) SetFloatValue(v float64)
func (n *NumberInput) Increment()
func (n *NumberInput) Decrement()
func (n *NumberInput) SetRange(min, max float64)
func (n *NumberInput) SetStep(step float64)
func (n *NumberInput) SetPrecision(p int)
// Inherits all TextInput methods via embedding
```

### TextArea

```go
func NewTextArea(placeholder string) *TextArea
func (t *TextArea) Value() string
func (t *TextArea) SetValue(val string)
func (t *TextArea) Insert(r rune)
func (t *TextArea) InsertNewline()
func (t *TextArea) Delete()
func (t *TextArea) MoveUp()
func (t *TextArea) MoveDown()
func (t *TextArea) MoveLeft()
func (t *TextArea) MoveRight()
func (t *TextArea) Focus() render.Cmd
func (t *TextArea) Blur()
func (t *TextArea) Focused() bool
func (t *TextArea) Validate() error
func (t *TextArea) AddValidator(v Validator)
func (t *TextArea) SetValidators(v ...Validator)
func (t *TextArea) SetWidth(w int)
func (t *TextArea) SetHeight(h int)
```

## Running Tests

```bash
# Run tests with coverage
go test ./ui/forms/... -cover

# Run tests with verbose output
go test ./ui/forms/... -v
```

## Design Notes

- **Engine Agnostic**: Components work with any rendering engine via `render.Model` interface
- **Validation**: Per-component validation with error display in `View()`
- **Cursor Blinking**: Shared `BlinkMsg` and `BlinkCmd()` for all components
- **Keyboard Navigation**: Standard shortcuts (arrow keys, home/end, delete/backspace)
- **Extensible**: Easy to add new form components following existing patterns
