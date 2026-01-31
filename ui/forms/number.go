package forms

import (
	"fmt"
	"strconv"
	
	"github.com/wwsheng009/taproot/ui/render"
)

// NumberInput is an input for numeric values.
type NumberInput struct {
	*TextInput
	min, max  float64
	step      float64
	precision int
}

// NewNumberInput creates a new number input.
func NewNumberInput(placeholder string) *NumberInput {
	n := &NumberInput{
		TextInput: NewTextInput(placeholder),
		step:      1,
		precision: 0,
		min:       -1e9, // Default large range
		max:       1e9,
	}
	// Add validator ensuring it's a number and within range
	n.AddValidator(func(v string) error {
		if v == "" {
			return nil
		}
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("must be a number")
		}
		if val < n.min || val > n.max {
			return fmt.Errorf("must be between %v and %v", n.min, n.max)
		}
		return nil
	})
	return n
}

// Update implements render.Model.
func (n *NumberInput) Update(msg any) (render.Model, render.Cmd) {
	// Handle specific keys for increment/decrement
	if n.TextInput.focused {
		if keyMsg, ok := msg.(render.KeyMsg); ok {
			switch keyMsg.String() {
			case "up":
				n.Increment()
				return n, nil
			case "down":
				n.Decrement()
				return n, nil
			}
		}
	}
	
	model, cmd := n.TextInput.Update(msg)
	if ti, ok := model.(*TextInput); ok {
		n.TextInput = ti
	}
	return n, cmd
}

// Increment increases the value by step.
func (n *NumberInput) Increment() {
	val := n.FloatValue()
	n.SetFloatValue(val + n.step)
}

// Decrement decreases the value by step.
func (n *NumberInput) Decrement() {
	val := n.FloatValue()
	n.SetFloatValue(val - n.step)
}

// FloatValue returns the value as float64.
func (n *NumberInput) FloatValue() float64 {
	if n.Value() == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(n.Value(), 64)
	return v
}

// SetFloatValue sets the value from float64.
func (n *NumberInput) SetFloatValue(v float64) {
	if v < n.min {
		v = n.min
	}
	if v > n.max {
		v = n.max
	}
	format := fmt.Sprintf("%%.%df", n.precision)
	n.SetValue(fmt.Sprintf(format, v))
}

// SetRange sets the min and max values.
func (n *NumberInput) SetRange(min, max float64) {
	n.min = min
	n.max = max
}

// SetStep sets the increment step.
func (n *NumberInput) SetStep(step float64) {
	n.step = step
}

// SetPrecision sets the decimal precision.
func (n *NumberInput) SetPrecision(p int) {
	n.precision = p
}
