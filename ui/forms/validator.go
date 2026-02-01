package forms

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Validator is a function that validates a string value.
// It returns an error if validation fails, or nil if successful.
type Validator func(value string) error

// Required checks that the value is not empty.
func Required(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("this field is required")
	}
	return nil
}

// MinLength returns a validator that checks the minimum length.
func MinLength(min int) Validator {
	return func(value string) error {
		if len(value) < min {
			return fmt.Errorf("must be at least %d characters", min)
		}
		return nil
	}
}

// MaxLength returns a validator that checks the maximum length.
func MaxLength(max int) Validator {
	return func(value string) error {
		if len(value) > max {
			return fmt.Errorf("must be at most %d characters", max)
		}
		return nil
	}
}

// EmailRegex is a simple regex for email validation.
var EmailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// Email checks that the value is a valid email address.
func Email(value string) error {
	if value == "" {
		return nil // Use Required() if empty is invalid
	}
	if !EmailRegex.MatchString(strings.ToLower(value)) {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

// Regex returns a validator that checks against a regular expression.
func Regex(re *regexp.Regexp, errorMsg string) Validator {
	return func(value string) error {
		if value == "" {
			return nil
		}
		if !re.MatchString(value) {
			return errors.New(errorMsg)
		}
		return nil
	}
}

// Range returns a validator that checks numeric range (for string inputs representing numbers).
// It assumes the value can be parsed as a float.
func Range(min, max float64) Validator {
	return func(value string) error {
		if value == "" {
			return nil
		}
		var v float64
		_, err := fmt.Sscanf(value, "%f", &v)
		if err != nil {
			return fmt.Errorf("invalid number")
		}
		if v < min || v > max {
			return fmt.Errorf("must be between %v and %v", min, max)
		}
		return nil
	}
}

// Compose combines multiple validators.
func Compose(validators ...Validator) Validator {
	return func(value string) error {
		for _, v := range validators {
			if err := v(value); err != nil {
				return err
			}
		}
		return nil
	}
}
