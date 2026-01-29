package files

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// formatSize formats a file size in bytes to a human-readable string.
func formatSize(size int64) string {
	if size < 0 {
		return "0 B"
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

// TrimFirstRune removes the first rune from a string.
func TrimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// ContainsRune checks if a string contains a specific rune.
func ContainsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

// HasPrefixCaseInsensitive checks if a string starts with a prefix, ignoring case.
func HasPrefixCaseInsensitive(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	sRunes := []rune(s)
	prefixRunes := []rune(prefix)
	for i, r := range prefixRunes {
		if unicode.ToLower(r) != unicode.ToLower(sRunes[i]) {
			return false
		}
	}
	return true
}

// ContainsCaseInsensitive checks if a string contains a substring, ignoring case.
func ContainsCaseInsensitive(s, substr string) bool {
	sRunes := []rune(s)
	substrRunes := []rune(substr)
	if len(substrRunes) > len(sRunes) {
		return false
	}

	for i := 0; i <= len(sRunes)-len(substrRunes); i++ {
		match := true
		for j, r := range substrRunes {
			if unicode.ToLower(r) != unicode.ToLower(sRunes[i+j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// WildcardMatch checks if a string matches a pattern with wildcards.
// Supports * (matches any sequence) and ? (matches any single character).
func WildcardMatch(pattern, text string, caseSensitive bool) bool {
	if !caseSensitive {
		pattern = string(toLowerCase([]rune(pattern)))
		text = string(toLowerCase([]rune(text)))
	}
	return wildcardMatchRune([]rune(pattern), []rune(text))
}

func wildcardMatchRune(pattern, text []rune) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			// Try to match zero or more characters
			for len(text) >= 0 {
				if wildcardMatchRune(pattern[1:], text) {
					return true
				}
				if len(text) == 0 {
					break
				}
				text = text[1:]
			}
			return false
		case '?':
			if len(text) == 0 {
				return false
			}
			text = text[1:]
		default:
			if len(text) == 0 || pattern[0] != text[0] {
				return false
			}
			text = text[1:]
		}
		pattern = pattern[1:]
	}
	return len(text) == 0
}

func toLowerCase(runes []rune) []rune {
	for i, r := range runes {
		runes[i] = unicode.ToLower(r)
	}
	return runes
}
