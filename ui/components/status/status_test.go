package status

import (
	"testing"

	"github.com/wwsheng009/taproot/ui/render"
)

func TestStateString(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		expected string
	}{
		{"Disabled", StateDisabled, "disabled"},
		{"Starting", StateStarting, "starting"},
		{"Ready", StateReady, "ready"},
		{"Error", StateError, "error"},
		{"Unknown", State(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("State.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStateIcon(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		expected string
	}{
		{"Disabled", StateDisabled, "○"},
		{"Starting", StateStarting, "⟳"},
		{"Ready", StateReady, "●"},
		{"Error", StateError, "×"},
		{"Unknown", State(99), "?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.Icon(); got != tt.expected {
				t.Errorf("State.Icon() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDiagnosticCounts(t *testing.T) {
	t.Run("Total", func(t *testing.T) {
		d := DiagnosticCounts{
			Error:       3,
			Warning:     2,
			Information: 5,
			Hint:        1,
		}
		expected := 11
		if got := d.Total(); got != expected {
			t.Errorf("DiagnosticCounts.Total() = %v, want %v", got, expected)
		}
	})

	t.Run("HasAny", func(t *testing.T) {
		d := DiagnosticCounts{Error: 1}
		if !d.HasAny() {
			t.Error("Expected HasAny() to return true when Error > 0")
		}

		d2 := DiagnosticCounts{}
		if d2.HasAny() {
			t.Error("Expected HasAny() to return false when all counts are 0")
		}
	})

	t.Run("HasErrors", func(t *testing.T) {
		d := &DiagnosticCounts{Error: 1}
		if !d.HasErrors() {
			t.Error("Expected HasErrors() to return true")
		}

		d.Error = 0
		if d.HasErrors() {
			t.Error("Expected HasErrors() to return false")
		}
	})

	t.Run("HasWarnings", func(t *testing.T) {
		d := &DiagnosticCounts{Warning: 1}
		if !d.HasWarnings() {
			t.Error("Expected HasWarnings() to return true")
		}

		d.Warning = 0
		if d.HasWarnings() {
			t.Error("Expected HasWarnings() to return false")
		}
	})

	t.Run("HasProblems", func(t *testing.T) {
		d := &DiagnosticCounts{Error: 1}
		if !d.HasProblems() {
			t.Error("Expected HasProblems() to return true when Error > 0")
		}

		d.Error = 0
		d.Warning = 1
		if !d.HasProblems() {
			t.Error("Expected HasProblems() to return true when Warning > 0")
		}

		d.Warning = 0
		if d.HasProblems() {
			t.Error("Expected HasProblems() to return false when no problems")
		}
	})

	t.Run("Add", func(t *testing.T) {
		d := &DiagnosticCounts{}
		d.Add(DiagnosticSeverityError)
		if d.Error != 1 {
			t.Errorf("Expected Error = 1, got %d", d.Error)
		}

		d.Add(DiagnosticSeverityWarning)
		if d.Warning != 1 {
			t.Errorf("Expected Warning = 1, got %d", d.Warning)
		}

		d.Add(DiagnosticSeverityInfo)
		if d.Information != 1 {
			t.Errorf("Expected Information = 1, got %d", d.Information)
		}

		d.Add(DiagnosticSeverityHint)
		if d.Hint != 1 {
			t.Errorf("Expected Hint = 1, got %d", d.Hint)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		d := &DiagnosticCounts{
			Error:       5,
			Warning:     3,
			Information: 2,
			Hint:        1,
		}
		d.Clear()
		if d.Error != 0 || d.Warning != 0 || d.Information != 0 || d.Hint != 0 {
			t.Error("Clear() did not reset all counts to 0")
		}
	})
}

func TestToolCounts(t *testing.T) {
	t.Run("Total", func(t *testing.T) {
		tc := ToolCounts{Tools: 5, Prompts: 3}
		expected := 8
		if got := tc.Total(); got != expected {
			t.Errorf("ToolCounts.Total() = %v, want %v", got, expected)
		}
	})

	t.Run("HasAny", func(t *testing.T) {
		tc := ToolCounts{Tools: 1}
		if !tc.HasAny() {
			t.Error("Expected HasAny() to return true when Tools > 0")
		}

		tc2 := ToolCounts{}
		if tc2.HasAny() {
			t.Error("Expected HasAny() to return false when both counts are 0")
		}
	})
}

func TestDiagnosticSeverityString(t *testing.T) {
	tests := []struct {
		name     string
		severity DiagnosticSeverity
		expected string
	}{
		{"Error", DiagnosticSeverityError, "error"},
		{"Warning", DiagnosticSeverityWarning, "warning"},
		{"Info", DiagnosticSeverityInfo, "info"},
		{"Hint", DiagnosticSeverityHint, "hint"},
		{"Unknown", DiagnosticSeverity(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.severity.String(); got != tt.expected {
				t.Errorf("DiagnosticSeverity.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	s := NewService("test-id", "Test Service")
	if s.ID() != "test-id" {
		t.Errorf("Expected ID = 'test-id', got %v", s.ID())
	}
	if s.Name() != "Test Service" {
		t.Errorf("Expected Name = 'Test Service', got %v", s.Name())
	}
	if s.Status() != StateDisabled {
		t.Errorf("Expected default Status = StateDisabled, got %v", s.Status())
	}
	if s.ErrorCount() != 0 {
		t.Errorf("Expected default ErrorCount = 0, got %v", s.ErrorCount())
	}
}

func TestServiceCmp(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		s := NewService("id", "name")
		cmd := s.Init()
		if cmd != nil {
			t.Error("Init() should return nil")
		}
	})

	t.Run("Update", func(t *testing.T) {
		s := NewService("id", "name")
		newModel, cmd := s.Update(&render.FocusGainMsg{})
		if cmd != render.None() {
			t.Error("Update() should return None() command")
		}
		if !newModel.(*ServiceCmp).Focused() {
			t.Error("FocusGainMsg should set focused to true")
		}

		newModel, cmd = s.Update(&render.BlurMsg{})
		if newModel.(*ServiceCmp).Focused() {
			t.Error("BlurMsg should set focused to false")
		}
	})

	t.Run("SetStatus", func(t *testing.T) {
		s := NewService("id", "name")
		s.SetStatus(StateReady)
		if s.Status() != StateReady {
			t.Errorf("Expected Status = StateReady, got %v", s.Status())
		}
		if !s.IsOnline() {
			t.Error("IsOnline() should return true when Status = StateReady")
		}
	})

	t.Run("SetErrorCount", func(t *testing.T) {
		s := NewService("id", "name")
		s.SetErrorCount(5)
		if s.ErrorCount() != 5 {
			t.Errorf("Expected ErrorCount = 5, got %v", s.ErrorCount())
		}
	})

	t.Run("FocusBlur", func(t *testing.T) {
		s := NewService("id", "name")
		s.Focus()
		if !s.Focused() {
			t.Error("Focus() should set focused to true")
		}

		s.Blur()
		if s.Focused() {
			t.Error("Blur() should set focused to false")
		}
	})

	t.Run("Compact", func(t *testing.T) {
		s := NewService("id", "name")
		s.SetCompact(true)
		if !s.Compact() {
			t.Error("Compact() should return true after SetCompact(true)")
		}
	})

	t.Run("MaxWidth", func(t *testing.T) {
		s := NewService("id", "name")
		s.SetMaxWidth(50)
		if s.MaxWidth() != 50 {
			t.Errorf("Expected MaxWidth = 50, got %v", s.MaxWidth())
		}
	})
}

func TestNewDiagnosticStatus(t *testing.T) {
	d := NewDiagnosticStatus("test-source")
	if d.Source() != "test-source" {
		t.Errorf("Expected Source = 'test-source', got %v", d.Source())
	}
	if d.Total() != 0 {
		t.Errorf("Expected default Total = 0, got %v", d.Total())
	}
}

func TestDiagnosticStatusCmp(t *testing.T) {
	t.Run("AddDiagnostic", func(t *testing.T) {
		d := NewDiagnosticStatus("test")
		d.AddDiagnostic(DiagnosticSeverityError)
		if d.ErrorCount() != 1 {
			t.Errorf("Expected ErrorCount = 1, got %v", d.ErrorCount())
		}

		d.AddDiagnostic(DiagnosticSeverityWarning)
		if d.WarningCount() != 1 {
			t.Errorf("Expected WarningCount = 1, got %v", d.WarningCount())
		}
	})

	t.Run("Clear", func(t *testing.T) {
		d := NewDiagnosticStatus("test")
		d.AddDiagnostic(DiagnosticSeverityError)
		d.Clear()
		if d.Total() != 0 {
			t.Errorf("Expected Total = 0 after Clear(), got %v", d.Total())
		}
	})

	t.Run("SetSummary", func(t *testing.T) {
		d := NewDiagnosticStatus("test")
		summary := DiagnosticSummary{Error: 5, Warning: 3}
		d.SetSummary(summary)
		if d.ErrorCount() != 5 {
			t.Errorf("Expected ErrorCount = 5, got %v", d.ErrorCount())
		}
	})

	t.Run("HasProblems", func(t *testing.T) {
		d := NewDiagnosticStatus("test")
		if d.HasProblems() {
			t.Error("Expected HasProblems() = false with empty summary")
		}

		d.AddDiagnostic(DiagnosticSeverityError)
		if !d.HasProblems() {
			t.Error("Expected HasProblems() = true with error")
		}
	})

	t.Run("FocusBlur", func(t *testing.T) {
		d := NewDiagnosticStatus("test")
		d.Focus()
		if !d.Focused() {
			t.Error("Focus() should set focused to true")
		}

		d.Blur()
		if d.Focused() {
			t.Error("Blur() should set focused to false")
		}
	})

	t.Run("Compact", func(t *testing.T) {
		d := NewDiagnosticStatus("test")
		d.SetCompact(true)
		if !d.Compact() {
			t.Error("Compact() should return true after SetCompact(true)")
		}
	})

	t.Run("ShowHints", func(t *testing.T) {
		d := NewDiagnosticStatus("test")
		d.SetShowHints(false)
		if d.ShowHints() {
			t.Error("ShowHints() should return false after SetShowHints(false)")
		}
	})
}

func TestNewLSPList(t *testing.T) {
	l := NewLSPList()
	if l.Title() != "LSPs" {
		t.Errorf("Expected default Title = 'LSPs', got %v", l.Title())
	}
	if l.MaxItems() != 5 {
		t.Errorf("Expected default MaxItems = 5, got %v", l.MaxItems())
	}
	if len(l.Services()) != 0 {
		t.Errorf("Expected empty Services slice, got %v", len(l.Services()))
	}
}

func TestLSPList(t *testing.T) {
	t.Run("AddService", func(t *testing.T) {
		l := NewLSPList()
		service := LSPServiceInfo{
			Name:     "gopls",
			Language: "go",
			State:    StateReady,
		}
		l.AddService(service)
		if len(l.Services()) != 1 {
			t.Errorf("Expected 1 service, got %v", len(l.Services()))
		}
	})

	t.Run("SetServices", func(t *testing.T) {
		l := NewLSPList()
		services := []LSPServiceInfo{
			{Name: "gopls", State: StateReady},
			{Name: "rust-analyzer", State: StateReady},
		}
		l.SetServices(services)
		if len(l.Services()) != 2 {
			t.Errorf("Expected 2 services, got %v", len(l.Services()))
		}
	})

	t.Run("ClearServices", func(t *testing.T) {
		l := NewLSPList()
		l.AddService(LSPServiceInfo{Name: "gopls", State: StateReady})
		l.ClearServices()
		if len(l.Services()) != 0 {
			t.Errorf("Expected 0 services after ClearServices(), got %v", len(l.Services()))
		}
	})

	t.Run("TotalErrors", func(t *testing.T) {
		l := NewLSPList()
		l.AddService(LSPServiceInfo{
			Name:        "gopls",
			State:       StateReady,
			Diagnostics: DiagnosticSummary{Error: 3},
		})
		l.AddService(LSPServiceInfo{
			Name:        "rust-analyzer",
			State:       StateReady,
			Diagnostics: DiagnosticSummary{Error: 2},
		})
		if l.TotalErrors() != 5 {
			t.Errorf("Expected TotalErrors = 5, got %v", l.TotalErrors())
		}
	})

	t.Run("OnlineCount", func(t *testing.T) {
		l := NewLSPList()
		l.AddService(LSPServiceInfo{Name: "gopls", State: StateReady})
		l.AddService(LSPServiceInfo{Name: "rust-analyzer", State: StateStarting})
		l.AddService(LSPServiceInfo{Name: "pylsp", State: StateError})
		if l.OnlineCount() != 1 {
			t.Errorf("Expected OnlineCount = 1, got %v", l.OnlineCount())
		}
	})
}

func TestNewMCPList(t *testing.T) {
	m := NewMCPList()
	if m.Title() != "MCPs" {
		t.Errorf("Expected default Title = 'MCPs', got %v", m.Title())
	}
	if m.MaxItems() != 5 {
		t.Errorf("Expected default MaxItems = 5, got %v", m.MaxItems())
	}
	if len(m.Services()) != 0 {
		t.Errorf("Expected empty Services slice, got %v", len(m.Services()))
	}
}

func TestMCPList(t *testing.T) {
	t.Run("AddService", func(t *testing.T) {
		m := NewMCPList()
		service := MCPServiceInfo{
			Name:  "filesystem",
			State: StateReady,
		}
		m.AddService(service)
		if len(m.Services()) != 1 {
			t.Errorf("Expected 1 service, got %v", len(m.Services()))
		}
	})

	t.Run("SetServices", func(t *testing.T) {
		m := NewMCPList()
		services := []MCPServiceInfo{
			{Name: "filesystem", State: StateReady},
			{Name: "git", State: StateReady},
		}
		m.SetServices(services)
		if len(m.Services()) != 2 {
			t.Errorf("Expected 2 services, got %v", len(m.Services()))
		}
	})

	t.Run("TotalTools", func(t *testing.T) {
		m := NewMCPList()
		m.AddService(MCPServiceInfo{
			Name:       "filesystem",
			State:      StateReady,
			ToolCounts: ToolCounts{Tools: 5},
		})
		m.AddService(MCPServiceInfo{
			Name:       "git",
			State:      StateReady,
			ToolCounts: ToolCounts{Tools: 3},
		})
		if m.TotalTools() != 8 {
			t.Errorf("Expected TotalTools = 8, got %v", m.TotalTools())
		}
	})

	t.Run("ConnectedCount", func(t *testing.T) {
		m := NewMCPList()
		m.AddService(MCPServiceInfo{Name: "filesystem", State: StateReady})
		m.AddService(MCPServiceInfo{Name: "git", State: StateStarting})
		if m.ConnectedCount() != 1 {
			t.Errorf("Expected ConnectedCount = 1, got %v", m.ConnectedCount())
		}
	})
}

func TestBackwardCompatibility(t *testing.T) {
	t.Run("ServiceStatusAliases", func(t *testing.T) {
		// Test that ServiceStatus is an alias for State
		var status ServiceStatus = StateReady
		if status != StateReady {
			t.Error("ServiceStatus should be an alias for State")
		}

		// Test constant aliases
		if ServiceStatusOffline != StateDisabled {
			t.Error("ServiceStatusOffline should equal StateDisabled")
		}
		if ServiceStatusStarting != StateStarting {
			t.Error("ServiceStatusStarting should equal StateStarting")
		}
		if ServiceStatusOnline != StateReady {
			t.Error("ServiceStatusOnline should equal StateReady")
		}
		if ServiceStatusError != StateError {
			t.Error("ServiceStatusError should equal StateError")
		}
	})
}
