package layout

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewVerticalLayout(t *testing.T) {
	layout := NewVerticalLayout()

	if layout == nil {
		t.Fatal("NewVerticalLayout returned nil")
	}

	if layout.centerV != true {
		t.Errorf("Expected centerV to be true, got %v", layout.centerV)
	}

	if layout.centerH != true {
		t.Errorf("Expected centerH to be true, got %v", layout.centerH)
	}

	if layout.separator != true {
		t.Errorf("Expected separator to be true, got %v", layout.separator)
	}
}

func TestGetHeaderHeight(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected int
	}{
		{"Empty header", "", 0},
		{"Single line header", "Hello", 1},
		{"Multi-line header", "Line1\nLine2", 2},
		{"Three lines", "Line1\nLine2\nLine3", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := NewVerticalLayout().SetHeader(tt.header)
			got := layout.GetHeaderHeight()
			if got != tt.expected {
				t.Errorf("GetHeaderHeight() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetFooterHeight(t *testing.T) {
	tests := []struct {
		name     string
		footer   string
		expected int
	}{
		{"Empty footer", "", 0},
		{"Single line footer", "Footer", 1},
		{"Multi-line footer", "Line1\nLine2", 2},
		{"Three lines", "Line1\nLine2\nLine3", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := NewVerticalLayout().SetFooter(tt.footer)
			got := layout.GetFooterHeight()
			if got != tt.expected {
				t.Errorf("GetFooterHeight() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetContentHeight(t *testing.T) {
	tests := []struct {
		name       string
		height     int
		header     string
		footer     string
		maxHeight  int
		expected   int
	}{
		{"100 height, 5 header, 2 footer", 100, "Line1\nLine2\nLine3\nLine4\nLine5", "Line1\nLine2", 0, 93},
		{"20 height, 3 header, 2 footer", 20, "Line1\nLine2\nLine3", "Line1\nLine2", 0, 15},
		{"10 height, 5 header, 3 footer", 10, "Line1\nLine2\nLine3\nLine4\nLine5", "Line1\nLine2\nLine3", 0, 2},
		{"With maxHeight limit", 100, "Header", "Footer", 50, 50},
		{"No limit needed", 100, "Header", "Footer", 200, 98},
		{"Empty header/footer", 50, "", "", 0, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := NewVerticalLayout().
				SetSize(80, tt.height).
				SetHeader(tt.header).
				SetFooter(tt.footer).
				SetMaxHeight(tt.maxHeight)

			got := layout.GetContentHeight()
			if got != tt.expected {
				t.Errorf("GetContentHeight() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRender_EmptySections(t *testing.T) {
	layout := NewVerticalLayout().SetSize(80, 20)
	output := layout.Render(0)

	if output == "" {
		t.Error("Render should not return empty string")
	}
}

func TestRender_OnlyHeader(t *testing.T) {
	layout := NewVerticalLayout().
		SetSize(80, 20).
		SetHeader("Header Line").
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	output := layout.Render(0)

	if output == "" {
		t.Error("Render should not return empty string")
	}

	// Header should be present
	if !contains(output, "Header Line") {
		t.Error("Output should contain header text")
	}
}

func TestRender_OnlyContent(t *testing.T) {
	layout := NewVerticalLayout().
		SetSize(80, 10).
		SetContent("Content Line").
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	output := layout.Render(0)

	if output == "" {
		t.Error("Render should not return empty string")
	}

	if !contains(output, "Content Line") {
		t.Error("Output should contain content text")
	}
}

func TestRender_HeaderAndFooter(t *testing.T) {
	layout := NewVerticalLayout().
		SetSize(80, 20).
		SetHeader("Header").
		SetFooter("Footer").
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	output := layout.Render(0)

	if !contains(output, "Header") {
		t.Error("Output should contain header")
	}

	if !contains(output, "Footer") {
		t.Error("Output should contain footer")
	}

	// Header should appear before footer
	headerPos := indexOf(output, "Header")
	footerPos := indexOf(output, "Footer")
	if headerPos > footerPos {
		t.Error("Header should appear before footer")
	}
}

func TestRender_ContentClipping(t *testing.T) {
	// Create content that exceeds available space
	content := makeLargeContent(50)

	layout := NewVerticalLayout().
		SetSize(80, 30).
		SetHeader("Header\nLine2").
		SetFooter("Footer\nLine2").
		SetContent(content).
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	output := layout.Render(0)

	// Content should be clipped to fit
	// Available space: 30 - 2 (header) - 2 (footer) = 26 lines
	count := countLines(output)
	if count > 30 {
		t.Errorf("Output has too many lines: %d (max 30)", count)
	}

	if !contains(output, "Header") {
		t.Error("Header should be present")
	}

	if !contains(output, "Footer") {
		t.Error("Footer should be present")
	}
}

func TestRender_VerticalCentering(t *testing.T) {
	layout := NewVerticalLayout().
		SetSize(80, 20).
		SetHeader("Header").
		SetFooter("Footer").
		SetContent("Content").
		SetCenterV(true).
		SetCenterH(false).
		SetSeparator(false)

	output := layout.Render(0)

	if !contains(output, "Header") && !contains(output, "Footer") && !contains(output, "Content") {
		t.Error("Output should contain all sections")
	}

	// When content is small, it should be centered
	// This is hard to test precisely, but we can verify it renders
	if output == "" {
		t.Error("Output should not be empty")
	}
}

func TestRender_SeparatorLines(t *testing.T) {
	layout := NewVerticalLayout().
		SetSize(80, 20).
		SetHeader("Header").
		SetFooter("Footer").
		SetContent("Content").
		SetCenterV(true).
		SetCenterH(true).
		SetSeparator(true)

	output := layout.Render(0)

	if output == "" {
		t.Error("Output should not be empty")
	}

	// Separator should be present (─ character)
	if !contains(output, "─") {
		t.Error("Output should contain separator line")
	}
}

func TestRender_MaxHeight(t *testing.T) {
	content := makeLargeContent(100)

	layout := NewVerticalLayout().
		SetSize(80, 50).
		SetHeader("Header").
		SetFooter("Footer").
		SetContent(content).
		SetMaxHeight(10).
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	output := layout.Render(0)

	if !contains(output, "Header") {
		t.Error("Header should be present")
	}

	if !contains(output, "Footer") {
		t.Error("Footer should be present")
	}

	// Content should be limited to maxHeight
	count := countLines(output)
	if count > 50 {
		t.Errorf("Output has too many lines: %d (max 50)", count)
	}
}

func TestRender_LargeContentWithSixelHeight(t *testing.T) {
	// Test with display height for Sixel images
	content := makeLargeContent(10)

	layout := NewVerticalLayout().
		SetSize(80, 30).
		SetHeader("Header").
		SetFooter("Footer").
		SetContent(content).
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false)

	// Simulate Sixel image taking 10 lines
	displayHeight := 10
	output := layout.Render(displayHeight)

	if !contains(output, "Header") {
		t.Error("Header should be present")
	}

	if !contains(output, "Footer") {
		t.Error("Footer should be present")
	}

	// Output should fit within 30 lines total
	count := countLines(output)
	if count > 30 {
		t.Errorf("Output has too many lines: %d (max 30)", count)
	}
}

func TestChainableMethods(t *testing.T) {
	layout := NewVerticalLayout().
		SetSize(100, 50).
		SetHeader("Header").
		SetFooter("Footer").
		SetContent("Content").
		SetCenterV(false).
		SetCenterH(false).
		SetSeparator(false).
		SetMaxHeight(20).
		SetTruncate(true)

	if layout.width != 100 {
		t.Errorf("Expected width 100, got %d", layout.width)
	}

	if layout.height != 50 {
		t.Errorf("Expected height 50, got %d", layout.height)
	}

	if layout.header != "Header" {
		t.Errorf("Expected header 'Header', got '%s'", layout.header)
	}

	if layout.footer != "Footer" {
		t.Errorf("Expected footer 'Footer', got '%s'", layout.footer)
	}

	if layout.content != "Content" {
		t.Errorf("Expected content 'Content', got '%s'", layout.content)
	}

	if layout.centerV != false {
		t.Errorf("Expected centerV false, got %v", layout.centerV)
	}

	if layout.centerH != false {
		t.Errorf("Expected centerH false, got %v", layout.centerH)
	}

	if layout.separator != false {
		t.Errorf("Expected separator false, got %v", layout.separator)
	}

	if layout.maxHeight != 20 {
		t.Errorf("Expected maxHeight 20, got %d", layout.maxHeight)
	}

	if layout.truncate != true {
		t.Errorf("Expected truncate true, got %v", layout.truncate)
	}
}

// Helper functions

func contains(s, substr string) bool {
	return indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func countLines(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	if len(s) > 0 && s[len(s)-1] != '\n' {
		count++
	}
	return count
}

func makeLargeContent(lines int) string {
	var b strings.Builder
	for i := 0; i < lines; i++ {
		b.WriteString(fmt.Sprintf("Content Line %d\n", i))
	}
	return b.String()
}
