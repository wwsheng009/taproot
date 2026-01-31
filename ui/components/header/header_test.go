package header

import (
	"testing"
)

func TestNew(t *testing.T) {
	h := New()

	if h == nil {
		t.Fatal("New() returned nil")
	}

	if h.width != 0 {
		t.Errorf("Expected width 0, got %d", h.width)
	}

	if h.height != 0 {
		t.Errorf("Expected height 0, got %d", h.height)
	}

	if h.brand != "Charm™" {
		t.Errorf("Expected brand 'Charm™', got '%s'", h.brand)
	}

	if h.title != "CRUSH" {
		t.Errorf("Expected title 'CRUSH', got '%s'", h.title)
	}
}

func TestSize(t *testing.T) {
	h := New()

	w, ht := h.Size()
	if w != 0 || ht != 0 {
		t.Errorf("Expected (0, 0), got (%d, %d)", w, ht)
	}

	h.SetSize(100, 2)
	w, ht = h.Size()
	if w != 100 || ht != 2 {
		t.Errorf("Expected (100, 2), got (%d, %d)", w, ht)
	}
}

func TestSetSize(t *testing.T) {
	h := New()

	h.SetSize(200, 3)
	if h.width != 200 || h.height != 3 {
		t.Errorf("Expected (200, 3), got (%d, %d)", h.width, h.height)
	}
}

func TestSetBrand(t *testing.T) {
	h := New()

	h.SetBrand("MyBrand", "TEST")
	if h.brand != "MyBrand" {
		t.Errorf("Expected brand 'MyBrand', got '%s'", h.brand)
	}

	if h.title != "TEST" {
		t.Errorf("Expected title 'TEST', got '%s'", h.title)
	}
}

func TestSetSessionTitle(t *testing.T) {
	h := New()

	h.SetSessionTitle("My Session")
	if h.sessionTitle != "My Session" {
		t.Errorf("Expected 'My Session', got '%s'", h.sessionTitle)
	}
}

func TestSetWorkingDirectory(t *testing.T) {
	h := New()

	h.SetWorkingDirectory("/path/to/project")
	if h.workingDir != "/path/to/project" {
		t.Errorf("Expected '/path/to/project', got '%s'", h.workingDir)
	}
}

func TestSetTokenUsage(t *testing.T) {
	h := New()

	h.SetTokenUsage(64000, 128000, 1.50)
	if h.tokenUsed != 64000 {
		t.Errorf("Expected 64000, got %d", h.tokenUsed)
	}

	if h.tokenMax != 128000 {
		t.Errorf("Expected 128000, got %d", h.tokenMax)
	}

	if h.cost != 1.50 {
		t.Errorf("Expected 1.50, got %f", h.cost)
	}
}

func TestSetErrorCount(t *testing.T) {
	h := New()

	h.SetErrorCount(5)
	if h.errorCount != 5 {
		t.Errorf("Expected 5, got %d", h.errorCount)
	}
}

func TestSetDetailsOpen(t *testing.T) {
	h := New()

	if h.ShowingDetails() {
		t.Error("Expected ShowingDetails() to be false initially")
	}

	h.SetDetailsOpen(true)
	if !h.ShowingDetails() {
		t.Error("Expected ShowingDetails() to be true after SetDetailsOpen(true)")
	}

	h.SetDetailsOpen(false)
	if h.ShowingDetails() {
		t.Error("Expected ShowingDetails() to be false after SetDetailsOpen(false)")
	}
}

func TestSetCompactMode(t *testing.T) {
	h := New()

	h.SetCompactMode(true)
	if !h.compactMode {
		t.Error("Expected compactMode to be true")
	}

	h.SetCompactMode(false)
	if h.compactMode {
		t.Error("Expected compactMode to be false")
	}
}

func TestView(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*HeaderComponent)
		contains []string
		notContains []string
	}{
		{
			name:  "empty brand",
			setup: func(h *HeaderComponent) { h.SetBrand("", "") },
			contains: []string{""},
		},
		{
			name:  "basic view",
			setup: func(h *HeaderComponent) {},
			contains: []string{"Charm™", "CRUSH", "ctrl+d"},
		},
		{
			name: "with errors",
			setup: func(h *HeaderComponent) {
				h.SetErrorCount(3)
			},
			contains: []string{"×3"},
		},
		{
			name: "with session title",
			setup: func(h *HeaderComponent) {
				h.SetSessionTitle("Test Session")
			},
			contains: []string{},
		},
		{
			name: "with working directory",
			setup: func(h *HeaderComponent) {
				h.SetWorkingDirectory("/path/to/project")
			},
			contains: []string{},
		},
		{
			name: "with token usage",
			setup: func(h *HeaderComponent) {
				h.SetTokenUsage(64000, 128000, 1.50)
			},
			contains: []string{"50%"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New()
			h.SetSize(100, 2)
			tt.setup(h)

			view := h.View()

			for _, s := range tt.contains {
				if !containsString(view, s) {
					t.Errorf("View does not contain expected string: %s\nView: %s", s, view)
				}
			}

			for _, s := range tt.notContains {
				if containsString(view, s) {
					t.Errorf("View should not contain string: %s\nView: %s", s, view)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	h := New()

	newH, _ := h.Update(nil)
	if newH != h {
		t.Error("Update should return the same instance")
	}
}

func TestFormatTokenUsage(t *testing.T) {
	tests := []struct {
		name     string
		used     int
		max      int
		cost     float64
		expected string
	}{
		{
			name:     "basic usage",
			used:     64000,
			max:      128000,
			cost:     1.50,
			expected: "50% (64.0K) $1.50",
		},
		{
		 name:     "no max",
		 used:     1000,
		 max:      0,
		 cost:     0.50,
		 expected: "1000 $0.50",
		},
		{
			name:     "large usage",
			used:     100000,
			max:      128000,
			cost:     2.00,
			expected: "78% (100.0K) $2.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTokenUsage(tt.used, tt.max, tt.cost)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestFormatErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{
			name:     "no errors",
			count:    0,
			expected: "",
		},
		{
			name:     "one error",
			count:    1,
			expected: "1",
		},
		{
			name:     "multiple errors",
			count:    5,
			expected: "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatErrorMessage(tt.count)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
