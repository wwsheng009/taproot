package logo

import (
	"strings"
	"testing"

	"github.com/wwsheng009/taproot/ui/styles"
)

func TestRender(t *testing.T) {
	s := styles.DefaultStyles()

	tests := []struct {
		name     string
		version  string
		compact  bool
		width    int
	}{
		{"Full version", "v2.0.0", false, 80},
		{"Compact version", "v2.0.0", true, 60},
		{"Full with wider", "v1.5.0", false, 120},
		{"Compact narrow", "v1.0.0", true, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Opts{
				FieldColor:   styles.DefaultStyles().Primary,
				TitleColorA:  styles.DefaultStyles().Primary,
				TitleColorB:  styles.DefaultStyles().Secondary,
				CharmColor:   styles.DefaultStyles().Primary,
				VersionColor: styles.DefaultStyles().FgMuted,
				Width:        tt.width,
			}

			result := Render(&s, tt.version, tt.compact, opts)

			if result == "" {
				t.Error("Render returned empty string")
			}

			// Should contain version number
			if !strings.Contains(result, tt.version) {
				t.Errorf("Expected version '%s' in output", tt.version)
			}

			// Should contain "Taproot"
			if !strings.Contains(result, "Taproot") {
				t.Error("Expected 'Taproot' in output")
			}
		})
	}
}

func TestSmallRender(t *testing.T) {
	s := styles.DefaultStyles()

	tests := []struct {
		name  string
		width int
	}{
		{"Narrow", 20},
		{"Medium", 60},
		{"Wide", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SmallRender(&s, tt.width)

			if result == "" {
				t.Error("SmallRender returned empty string")
			}

			// Should contain "Taproot"
			if !strings.Contains(result, "Taproot") {
				t.Error("Expected 'Taproot' in output")
			}

			// Should contain "TUI"
			if !strings.Contains(result, "TUI") {
				t.Error("Expected 'TUI' in output")
			}
		})
	}
}

func TestRenderEmptyVersion(t *testing.T) {
	s := styles.DefaultStyles()
	opts := Opts{
		Width: 80,
	}

	result := Render(&s, "", false, opts)

	if result == "" {
		t.Error("Render returned empty string for empty version")
	}

	// Should still contain "Taproot"
	if !strings.Contains(result, "Taproot") {
		t.Error("Expected 'Taproot' even with empty version")
	}
}

func TestRenderZeroWidth(t *testing.T) {
	s := styles.DefaultStyles()
	opts := Opts{
		Width: 0,
	}

	result := Render(&s, "v1.0.0", false, opts)

	// Should still produce output
	if result == "" {
		t.Error("Render returned empty string for zero width")
	}
}

func TestRenderNegativeWidth(t *testing.T) {
	s := styles.DefaultStyles()
	opts := Opts{
		Width: -10,
	}

	result := Render(&s, "v1.0.0", false, opts)

	// Should still produce output
	if result == "" {
		t.Error("Render returned empty string for negative width")
	}
}

func TestSmallRenderZeroWidth(t *testing.T) {
	s := styles.DefaultStyles()

	result := SmallRender(&s, 0)

	if result == "" {
		t.Error("SmallRender returned empty string for zero width")
	}
}

func TestLetterforms(t *testing.T) {
	letterFuncs := []letterform{
		letterT,
		letterA,
		letterP,
		letterR,
		letterO1,
		letterO2,
	}

	for _, fn := range letterFuncs {
		result := fn(false)
		if result == "" {
			t.Error("Letterform returned empty string")
		}

		resultStretched := fn(true)
		if resultStretched == "" {
			t.Error("Stretched letterform returned empty string")
		}

		// Stretched version should typically be wider (or equal)
		// This is a loose check since stretching is random
		lines := strings.Split(result, "\n")
		linesStretched := strings.Split(resultStretched, "\n")
		if len(lines) != len(linesStretched) {
			t.Error("Letterform should have same number of lines regardless of stretching")
		}
	}
}

func TestRenderWord(t *testing.T) {
	letterFuncs := []letterform{
		letterT,
		letterA,
		letterP,
		letterR,
		letterO1,
		letterO2,
		letterT,
	}

	tests := []struct {
		name         string
		spacing      int
		stretchIndex int
	}{
		{"No spacing", 0, -1},
		{"With spacing", 1, -1},
		{"Stretch first", 1, 0},
		{"Stretch middle", 1, 3},
		{"Stretch last", 1, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderWord(tt.spacing, tt.stretchIndex, letterFuncs...)

			if result == "" {
				t.Error("renderWord returned empty string")
			}

			// Result should be multi-line
			lines := strings.Split(result, "\n")
			if len(lines) < 2 {
				t.Errorf("Expected multi-line output, got %d lines", len(lines))
			}
		})
	}
}

func TestRenderWordNegativeSpacing(t *testing.T) {
	letterFuncs := []letterform{letterT, letterA, letterP}
	result := renderWord(-1, -1, letterFuncs...)

	if result == "" {
		t.Error("renderWord should handle negative spacing")
	}
}

func TestRenderWordEmpty(t *testing.T) {
	result := renderWord(1, -1)

	if result != "" {
		t.Error("renderWord with no letterforms should return empty string")
	}
}

func TestJoinLetterform(t *testing.T) {
	result := joinLetterform("A", "B", "C")

	if result == "" {
		t.Error("joinLetterform returned empty string")
	}

	if !strings.Contains(result, "A") || !strings.Contains(result, "B") || !strings.Contains(result, "C") {
		t.Error("joinLetterform should contain all parts")
	}
}

func TestJoinLetterformEmpty(t *testing.T) {
	result := joinLetterform()

	if result != "" {
		t.Error("joinLetterform with no arguments should return empty string")
	}
}



func TestMin(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{1, 2, 1},
		{5, 3, 3},
		{-1, 1, -1},
		{0, 0, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := min(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("min(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCachedRandN(t *testing.T) {
	// Test that caching works
	result1 := cachedRandN(10)
	result2 := cachedRandN(10)

	if result1 != result2 {
		t.Error("cachedRandN should return same value for same input")
	}

	// Different inputs should (likely) have different results
	// (though technically they could be the same)
	result3 := cachedRandN(20)
	_ = result3
}

func TestStretchLetterformPart(t *testing.T) {
	propsNoStretch := letterformProps{
		stretch:    false,
		width:      5,
		minStretch: 7,
		maxStretch: 10,
	}

	result1 := stretchLetterformPart("x", propsNoStretch)
	expectedLen := 5 // width times 1 character
	lines1 := strings.Split(result1, "\n")
	if len(lines1) != 1 || len(lines1[0]) != expectedLen {
		t.Errorf("Expected %d characters, got %d", expectedLen, len(lines1[0]))
	}

	// Test stretch
	propsStretch := letterformProps{
		stretch:    true,
		width:      1,
		minStretch: 5,
		maxStretch: 10,
	}

	result2 := stretchLetterformPart("x", propsStretch)
	lines2 := strings.Split(result2, "\n")
	if len(lines2) != 1 {
		t.Error("Expected single line output")
	}

	// Stretched version should be longer or equal to minStretch
	if len(lines2[0]) < propsStretch.minStretch {
		t.Errorf("Expected at least %d characters when stretched, got %d", propsStretch.minStretch, len(lines2[0]))
	}
}

func TestLetterT(t *testing.T) {
	normal := letterT(false)
	stretched := letterT(true)

	if normal == "" || stretched == "" {
		t.Error("letterT should return non-empty strings")
	}

	// Both should be multi-line
	linesN := strings.Split(normal, "\n")
	linesS := strings.Split(stretched, "\n")

	if len(linesN) < 2 || len(linesS) < 2 {
		t.Error("letterT should return multi-line output")
	}
}

func TestLetterA(t *testing.T) {
	normal := letterA(false)
	stretched := letterA(true)

	if normal == "" || stretched == "" {
		t.Error("letterA should return non-empty strings")
	}

	linesN := strings.Split(normal, "\n")
	linesS := strings.Split(stretched, "\n")

	if len(linesN) != len(linesS) {
		t.Error("letterA should have same number of lines regardless of stretch")
	}
}

func TestLetterP(t *testing.T) {
	normal := letterP(false)
	stretched := letterP(true)

	if normal == "" || stretched == "" {
		t.Error("letterP should return non-empty strings")
	}
}

func TestLetterR(t *testing.T) {
	normal := letterR(false)
	stretched := letterR(true)

	if normal == "" || stretched == "" {
		t.Error("letterR should return non-empty strings")
	}

	linesN := strings.Split(normal, "\n")
	linesS := strings.Split(stretched, "\n")

	if len(linesN) != len(linesS) {
		t.Error("letterR should have same number of lines regardless of stretch")
	}
}

func TestLetterO1(t *testing.T) {
	normal := letterO1(false)
	stretched := letterO1(true)

	if normal == "" || stretched == "" {
		t.Error("letterO1 should return non-empty strings")
	}
}

func TestLetterO2(t *testing.T) {
	normal := letterO2(false)
	stretched := letterO2(true)

	if normal == "" || stretched == "" {
		t.Error("letterO2 should return non-empty strings")
	}

	// O2 is special - no stretching
	if normal != stretched {
		t.Error("letterO2 should not stretch (no stretchIndex logic)")
	}
}

func TestOptsDefaults(t *testing.T) {
	var opts Opts

	if opts.Width != 0 {
		t.Error("Opts should have zero default width")
	}
}

func TestRenderWithColors(t *testing.T) {
	s := styles.DefaultStyles()

	opts := Opts{
		FieldColor:   s.Primary,
		TitleColorA:  s.Secondary,
		TitleColorB:  s.Tertiary,
		CharmColor:   s.Tertiary,
		VersionColor: s.FgMuted,
		Width:        80,
	}

	result := Render(&s, "v2.0.0", false, opts)

	if result == "" {
		t.Error("Render with full color options should work")
	}
}

func TestMultipleFormats(t *testing.T) {
	s := styles.DefaultStyles()
	opts := Opts{Width: 80}

	version := "v2.0.0"

	full := Render(&s, version, false, opts)
	compact := Render(&s, version, true, opts)
	small := SmallRender(&s, 80)

	// All should be non-empty
	if full == "" || compact == "" || small == "" {
		t.Error("All render formats should return non-empty strings")
	}

	// All should contain version
	if !strings.Contains(full, version) || !strings.Contains(compact, version) || !strings.Contains(small, version) {
		t.Error("All render formats should contain version")
	}
}

func BenchmarkRender(b *testing.B) {
	s := styles.DefaultStyles()
	opts := Opts{Width: 80}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Render(&s, "v2.0.0", false, opts)
	}
}

func BenchmarkSmallRender(b *testing.B) {
	s := styles.DefaultStyles()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SmallRender(&s, 80)
	}
}

func BenchmarkRenderWord(b *testing.B) {
	letterFuncs := []letterform{letterT, letterA, letterP, letterR, letterO1, letterO2, letterT}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = renderWord(1, -1, letterFuncs...)
	}
}
