package highlight

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourorg/taproot/internal/ui/styles"
)

func TestSyntaxHighlight(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		fileName  string
		wantEmpty bool
	}{
		{
			name:      "empty source",
			source:    "",
			fileName:  "test.go",
			wantEmpty: true,
		},
		{
			name: "simple Go code",
			source: `package main

func main() {
	println("Hello, World!")
}`,
			fileName:  "test.go",
			wantEmpty: false,
		},
		{
			name: "JavaScript code",
			source: `function hello() {
	console.log("Hello");
}`,
			fileName:  "test.js",
			wantEmpty: false,
		},
		{
			name: "unknown file extension",
			source: `some text`,
			fileName:  "test.unknown",
			wantEmpty: false, // Should still produce output with fallback lexer
		},
		{
			name: "no filename - auto detect",
			source: `package main

func main() {
	println("detected")
}`,
			fileName:  "",
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg := lipgloss.Color("#1a1a2e")
			s := styles.DefaultStyles()
			result, err := SyntaxHighlight(&s, tt.source, tt.fileName, bg)

			if err != nil {
				t.Errorf("SyntaxHighlight() error = %v", err)
				return
			}

			if tt.wantEmpty && result != "" {
				t.Errorf("SyntaxHighlight() = %q, want empty", result)
			}
			if !tt.wantEmpty && result == "" {
				t.Error("SyntaxHighlight() returned empty, want non-empty")
			}
		})
	}
}

func TestSyntaxHighlightBackgroundColor(t *testing.T) {
	source := `package main
func main() {}`
	bg := lipgloss.Color("#FF0000")
	s := styles.DefaultStyles()

	result, err := SyntaxHighlight(&s, source, "test.go", bg)
	if err != nil {
		t.Fatalf("SyntaxHighlight() error = %v", err)
	}

	// Check that result contains ANSI codes (should have color codes)
	if result == "" {
		t.Error("SyntaxHighlight() returned empty result")
	}

	// The result should be different from the source (formatted with colors)
	if result == source {
		t.Error("SyntaxHighlight() returned unformatted source")
	}
}
