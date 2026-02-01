package forms

import "github.com/wwsheng009/taproot/ui/render"

// Input is the interface for all form input components.
type Input interface {
	render.Model

	// Value returns the current value of the input.
	Value() string

	// SetValue sets the value of the input.
	SetValue(string)

	// Focus focuses the input.
	Focus() render.Cmd

	// Blur blurs the input.
	Blur()

	// Focused returns whether the input is focused.
	Focused() bool

	// Validate validates the input and returns an error if invalid.
	Validate() error

	// Error returns the current validation error, if any.
	Error() error
}
