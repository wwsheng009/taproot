package forms

import (
	"fmt"
	"regexp"
	"strconv"
)

// Validator functions return an error if the value is invalid.
type Validator func(value string) error

// Required checks if the value is not empty.
func Required(value string) error {
	if value == "" {
		return fmt.Errorf("this field is required")
	}
	return nil
}

// MinLength checks if the value has at least n characters.
func MinLength(n int) Validator {
	return func(value string) error {
		if len(value) < n {
			return fmt.Errorf("must be at least %d characters", n)
		}
		return nil
	}
}

// MaxLength checks if the value has at most n characters.
func MaxLength(n int) Validator {
	return func(value string) error {
		if len(value) > n {
			return fmt.Errorf("must be at most %d characters", n)
		}
		return nil
	}
}

// Email checks if the value is a valid email address.
func Email(value string) error {
	if value == "" {
		return nil // Use Required for empty check
	}
	// Simple regex for email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

// Range checks if the numeric value is within min and max (inclusive).
func Range(min, max float64) Validator {
	return func(value string) error {
		if value == "" {
			return nil
		}
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid number")
		}
		if v < min || v > max {
			return fmt.Errorf("must be between %g and %g", min, max)
		}
		return nil
	}
}
