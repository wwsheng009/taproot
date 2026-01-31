package styles

import (
	"testing"

	"github.com/alecthomas/chroma/v2"
	"github.com/charmbracelet/glamour/ansi"
)

func TestChromaStyle(t *testing.T) {
	tests := []struct {
		name     string
		style    ansi.StylePrimitive
		expected string
	}{
		{
			name:     "empty style",
			style:    ansi.StylePrimitive{},
			expected: "",
		},
		{
			name: "color only",
			style: ansi.StylePrimitive{
				Color: stringPtr("#FF0000"),
			},
			expected: "#FF0000",
		},
		{
			name: "background only",
			style: ansi.StylePrimitive{
				BackgroundColor: stringPtr("#00FF00"),
			},
			expected: "bg:#00FF00",
		},
		{
			name: "color and background",
			style: ansi.StylePrimitive{
				Color:           stringPtr("#FF0000"),
				BackgroundColor: stringPtr("#00FF00"),
			},
			expected: "#FF0000 bg:#00FF00",
		},
		{
			name: "bold",
			style: ansi.StylePrimitive{
				Bold: boolPtr(true),
			},
			expected: "bold",
		},
		{
			name: "italic",
			style: ansi.StylePrimitive{
				Italic: boolPtr(true),
			},
			expected: "italic",
		},
		{
			name: "underline",
			style: ansi.StylePrimitive{
				Underline: boolPtr(true),
			},
			expected: "underline",
		},
		{
			name: "all styles",
			style: ansi.StylePrimitive{
				Color:           stringPtr("#FF0000"),
				BackgroundColor: stringPtr("#00FF00"),
				Bold:            boolPtr(true),
				Italic:          boolPtr(true),
				Underline:       boolPtr(true),
			},
			expected: "#FF0000 bg:#00FF00 italic bold underline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chromaStyle(tt.style)
			if result != tt.expected {
				t.Errorf("chromaStyle() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetChromaTheme(t *testing.T) {
	s := DefaultStyles()
	theme := s.ChromaTheme()

	if theme == nil {
		t.Fatal("ChromaTheme() returned nil")
	}

	// Check that we have entries for common token types
	expectedTokens := []struct {
		name  string
		value chroma.TokenType
	}{
		{"Text", chroma.Text},
		{"Error", chroma.Error},
		{"Comment", chroma.Comment},
		{"Keyword", chroma.Keyword},
		{"Name", chroma.Name},
		{"NameFunction", chroma.NameFunction},
		{"LiteralString", chroma.LiteralString},
		{"LiteralNumber", chroma.LiteralNumber},
	}

	for _, token := range expectedTokens {
		if theme[token.value] == "" {
			t.Errorf("ChromaTheme missing entry for %s", token.name)
		}
	}
}

func TestPlainMarkdownStyle(t *testing.T) {
	s := DefaultStyles()
	style := s.PlainMarkdown

	// Check that Document has colors set
	if style.Document.Color == nil {
		t.Error("PlainMarkdown Document.Color is nil")
	}
	if style.Document.BackgroundColor == nil {
		t.Error("PlainMarkdown Document.BackgroundColor is nil")
	}

	// Check that headings exist
	if style.Heading.Color == nil {
		t.Error("PlainMarkdown Heading.Color is nil")
	}

	// CodeBlock.Chroma is nil in PlainMarkdownStyle (no syntax highlighting)
	// This is expected for plain text mode
}

func TestHelperFunctions(t *testing.T) {
	t.Run("stringPtr", func(t *testing.T) {
		s := "test"
		ptr := stringPtr(s)
		if ptr == nil {
			t.Fatal("stringPtr() returned nil")
		}
		if *ptr != s {
			t.Errorf("stringPtr() = %p, want %q", ptr, s)
		}
	})

	t.Run("boolPtr", func(t *testing.T) {
		b := true
		ptr := boolPtr(b)
		if ptr == nil {
			t.Fatal("boolPtr() returned nil")
		}
		if *ptr != b {
			t.Errorf("boolPtr() = %v, want %v", *ptr, b)
		}
	})

	t.Run("uintPtr", func(t *testing.T) {
		u := uint(42)
		ptr := uintPtr(u)
		if ptr == nil {
			t.Fatal("uintPtr() returned nil")
		}
		if *ptr != u {
			t.Errorf("uintPtr() = %d, want %d", *ptr, u)
		}
	})
}
