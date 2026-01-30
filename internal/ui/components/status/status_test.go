package status

import (
	"strings"
	"testing"

	"github.com/wwsheng009/taproot/internal/ui/render"
	"github.com/wwsheng009/taproot/internal/ui/styles"
)

// Test ServiceStatus String method
func TestServiceStatus_String(t *testing.T) {
	tests := []struct {
		status  ServiceStatus
		want    string
	}{
		{ServiceStatusOffline, "offline"},
		{ServiceStatusStarting, "starting"},
		{ServiceStatusConnecting, "connecting"},
		{ServiceStatusOnline, "online"},
		{ServiceStatusBusy, "busy"},
		{ServiceStatusError, "error"},
		{ServiceStatus(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("ServiceStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test DiagnosticSeverity String method
func TestDiagnosticSeverity_String(t *testing.T) {
	tests := []struct {
		severity DiagnosticSeverity
		want     string
	}{
		{DiagnosticSeverityError, "error"},
		{DiagnosticSeverityWarning, "warning"},
		{DiagnosticSeverityInfo, "info"},
		{DiagnosticSeverityHint, "hint"},
		{DiagnosticSeverity(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.severity.String(); got != tt.want {
				t.Errorf("DiagnosticSeverity.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test DiagnosticSummary methods
func TestDiagnosticSummary_Methods(t *testing.T) {
	tests := []struct {
		name     string
		add      []DiagnosticSeverity
		wantTotal   int
		wantError   int
		wantWarning int
		wantInfo    int
		wantHint    int
	}{
		{
			name:     "empty",
			add:      nil,
			wantTotal: 0,
		},
		{
			name:     "single error",
			add:      []DiagnosticSeverity{DiagnosticSeverityError},
			wantTotal: 1,
			wantError: 1,
		},
		{
			name:     "multiple types",
			add:      []DiagnosticSeverity{
				DiagnosticSeverityError,
				DiagnosticSeverityError,
				DiagnosticSeverityWarning,
				DiagnosticSeverityInfo,
				DiagnosticSeverityHint,
			},
			wantTotal:   5,
			wantError:   2,
			wantWarning: 1,
			wantInfo:    1,
			wantHint:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DiagnosticSummary{}

			for _, severity := range tt.add {
				d.Add(severity)
			}

			if got := d.Total(); got != tt.wantTotal {
				t.Errorf("DiagnosticSummary.Total() = %v, want %v", got, tt.wantTotal)
			}
			if got := d.Error; got != tt.wantError {
				t.Errorf("DiagnosticSummary.Error = %v, want %v", got, tt.wantError)
			}
			if got := d.Warning; got != tt.wantWarning {
				t.Errorf("DiagnosticSummary.Warning = %v, want %v", got, tt.wantWarning)
			}
			if got := d.Info; got != tt.wantInfo {
				t.Errorf("DiagnosticSummary.Info = %v, want %v", got, tt.wantInfo)
			}
			if got := d.Hint; got != tt.wantHint {
				t.Errorf("DiagnosticSummary.Hint = %v, want %v", got, tt.wantHint)
			}
		})
	}
}

// Test DiagnosticSummary Has methods
func TestDiagnosticSummary_HasMethods(t *testing.T) {
	tests := []struct {
		name          string
		errors        int
		warnings      int
		wantHasErrors bool
		wantHasWarnings bool
		wantHasProblems bool
	}{
		{"none", 0, 0, false, false, false},
		{"errors only", 1, 0, true, false, true},
		{"warnings only", 0, 1, false, true, true},
		{"both", 1, 1, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DiagnosticSummary{
				Error:   tt.errors,
				Warning: tt.warnings,
			}

			if got := d.HasErrors(); got != tt.wantHasErrors {
				t.Errorf("DiagnosticSummary.HasErrors() = %v, want %v", got, tt.wantHasErrors)
			}
			if got := d.HasWarnings(); got != tt.wantHasWarnings {
				t.Errorf("DiagnosticSummary.HasWarnings() = %v, want %v", got, tt.wantHasWarnings)
			}
			if got := d.HasProblems(); got != tt.wantHasProblems {
				t.Errorf("DiagnosticSummary.HasProblems() = %v, want %v", got, tt.wantHasProblems)
			}
		})
	}
}

// Test DiagnosticSummary Clear
func TestDiagnosticSummary_Clear(t *testing.T) {
	d := &DiagnosticSummary{
		Error:   5,
		Warning: 3,
		Info:    2,
		Hint:    1,
	}

	d.Clear()

	if d.Error != 0 || d.Warning != 0 || d.Info != 0 || d.Hint != 0 {
		t.Errorf("Clear() did not reset all fields: Error=%d, Warning=%d, Info=%d, Hint=%d",
			d.Error, d.Warning, d.Info, d.Hint)
	}
}

// Test ServiceCmp creation and initialization
func TestServiceCmp_NewAndInit(t *testing.T) {
	service := NewService("lsp", "Language Server")

	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	if err := service.Init(); err != nil {
		t.Errorf("Init() error = %v, want nil", err)
	}

	if !service.initialized {
		t.Error("Init() did not set initialized flag")
	}
}

// Test ServiceCmp interface implementation
func TestServiceCmp_Interface(t *testing.T) {
	var _ Service = (*ServiceCmp)(nil)
	var _ render.Model = (*ServiceCmp)(nil)
}

// Test ServiceCmp getters and setters
func TestServiceCmp_GettersAndSetters(t *testing.T) {
	service := NewService("test-id", "Test Service")

	if service.ID() != "test-id" {
		t.Errorf("ID() = %v, want test-id", service.ID())
	}

	if service.Name() != "Test Service" {
		t.Errorf("Name() = %v, want Test Service", service.Name())
	}

	service.SetStatus(ServiceStatusOnline)
	if service.Status() != ServiceStatusOnline {
		t.Errorf("Status() after SetStatus(ServiceStatusOnline) = %v, want ServiceStatusOnline", service.Status())
	}

	service.SetErrorCount(5)
	if service.ErrorCount() != 5 {
		t.Errorf("ErrorCount() after SetErrorCount(5) = %v, want 5", service.ErrorCount())
	}

	if !service.IsOnline() {
		t.Error("IsOnline() returned false after setting to ServiceStatusOnline")
	}

	service.SetCompact(true)
	if !service.Compact() {
		t.Error("Compact() returned false after setting to true")
	}

	service.SetMaxWidth(100)
	if service.MaxWidth() != 100 {
		t.Errorf("MaxWidth() = %v, want 100", service.MaxWidth())
	}
}

// Test ServiceCmp focus handling
func TestServiceCmp_Focus(t *testing.T) {
	service := NewService("test", "Test")

	if service.Focused() {
		t.Fatal("Focused() returned true initially")
	}

	service.Focus()
	if !service.Focused() {
		t.Error("Focus() did not set focused flag")
	}

	service.Blur()
	if service.Focused() {
		t.Error("Blur() did not clear focused flag")
	}
}

// Test ServiceCmp Update message handling
func TestServiceCmp_Update(t *testing.T) {
	service := NewService("test", "Test")
	service.Init()

	// Test FocusGainMsg
	model, cmd := service.Update(&render.FocusGainMsg{})
	if service.Focused() != true {
		t.Error("Update(FocusGainMsg) did not set focused flag")
	}
	if cmd != nil {
		t.Errorf("Update(FocusGainMsg) returned non-nil command")
	}
	if model != service {
		t.Error("Update() did not return the same model")
	}

	// Test BlurMsg
	model, cmd = service.Update(&render.BlurMsg{})
	if service.Focused() != false {
		t.Error("Update(BlurMsg) did not clear focused flag")
	}
	if cmd != nil {
		t.Errorf("Update(BlurMsg) returned non-nil command")
	}

	// Test other message types
	model, cmd = service.Update("other message")
	if service.Focused() {
		t.Error("Update(other message) changed focused state")
	}
	if model != service {
		t.Error("Update(other message) did not return the same model")
	}
}

// Test ServiceCmp View rendering
func TestServiceCmp_View(t *testing.T) {
	service := NewService("lsp", "Language Server")
	service.Init()

	view := service.View()

	if view == "" {
		t.Error("View() returned empty string")
	}

	if !strings.Contains(view, "●") {
		t.Error("View() does not contain status icon")
	}

	if !strings.Contains(view, "Language Server") {
		t.Error("View() does not contain service name")
	}
}

// Test ServiceCmp View with error count
func TestServiceCmp_ViewWithErrorCount(t *testing.T) {
	service := NewService("lsp", "Language Server")
	service.Init()
	service.SetErrorCount(3)

	view := service.View()

	if !strings.Contains(view, "×") {
		t.Error("View() with error count does not contain error icon")
	}

	if !strings.Contains(view, "3") {
		t.Error("View() with error count does not contain error count value")
	}
}

// Test ServiceCmp View in compact mode
func TestServiceCmp_ViewCompact(t *testing.T) {
	service := NewService("lsp", "Language Server")
	service.Init()
	service.SetCompact(true)

	view := service.View()

	if view == "" {
		t.Error("View() in compact mode returned empty string")
	}

	if strings.Contains(view, "offline") {
		t.Error("View() in compact mode should not show status label")
	}
}

// Test ServiceCmp View with max width
func TestServiceCmp_ViewWithMaxWidth(t *testing.T) {
	service := NewService("very-long-service-name", "Very Long Service Name")
	service.Init()
	service.SetMaxWidth(20)

	view := service.View()

	if view == "" {
		t.Error("View() returned empty string")
	}

}

// Test DiagnosticStatusCmp creation and initialization
func TestDiagnosticStatusCmp_NewAndInit(t *testing.T) {
	diag := NewDiagnosticStatus("lsp")

	if diag == nil {
		t.Fatal("NewDiagnosticStatus() returned nil")
	}

	if err := diag.Init(); err != nil {
		t.Errorf("Init() error = %v, want nil", err)
	}

	if !diag.initialized {
		t.Error("Init() did not set initialized flag")
	}

	if diag.Source() != "lsp" {
		t.Errorf("Source() = %v, want lsp", diag.Source())
	}
}

// Test DiagnosticStatusCmp interface implementation
func TestDiagnosticStatusCmp_Interface(t *testing.T) {
	var _ render.Model = (*DiagnosticStatusCmp)(nil)
}

// Test DiagnosticStatusCmp AddDiagnostic
func TestDiagnosticStatusCmp_AddDiagnostic(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()

	diag.AddDiagnostic(DiagnosticSeverityError)
	diag.AddDiagnostic(DiagnosticSeverityError)
	diag.AddDiagnostic(DiagnosticSeverityWarning)

	if diag.ErrorCount() != 2 {
		t.Errorf("ErrorCount() = %v, want 2", diag.ErrorCount())
	}

	if diag.WarningCount() != 1 {
		t.Errorf("WarningCount() = %v, want 1", diag.WarningCount())
	}

	if diag.Total() != 3 {
		t.Errorf("Total() = %v, want 3", diag.Total())
	}
}

// Test DiagnosticStatusCmp SetSummary
func TestDiagnosticStatusCmp_SetSummary(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()

	summary := DiagnosticSummary{
		Error:   5,
		Warning: 3,
		Info:    2,
		Hint:    1,
	}

	diag.SetSummary(summary)

	if diag.ErrorCount() != 5 {
		t.Errorf("ErrorCount() after SetSummary() = %v, want 5", diag.ErrorCount())
	}

	if diag.Total() != 11 {
		t.Errorf("Total() after SetSummary() = %v, want 11", diag.Total())
	}
}

// Test DiagnosticStatusCmp Clear
func TestDiagnosticStatusCmp_Clear(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()

	diag.AddDiagnostic(DiagnosticSeverityError)
	diag.AddDiagnostic(DiagnosticSeverityWarning)

	diag.Clear()

	if diag.Total() != 0 {
		t.Errorf("Total() after Clear() = %v, want 0", diag.Total())
	}
}

// Test DiagnosticStatusCmp source-specific diagnostics
func TestDiagnosticStatusCmp_SourceDiagnostics(t *testing.T) {
	diag := NewDiagnosticStatus("main")
	diag.Init()

	diag.AddSourceDiagnostic("file1.go", DiagnosticSeverityError)
	diag.AddSourceDiagnostic("file1.go", DiagnosticSeverityWarning)
	diag.AddSourceDiagnostic("file2.go", DiagnosticSeverityError)

	summary1 := diag.GetSourceSummary("file1.go")
	if summary1.Error != 1 || summary1.Warning != 1 {
		t.Errorf("file1.go summary: Error=%d, Warning=%d, want Error=1, Warning=1",
			summary1.Error, summary1.Warning)
	}

	summary2 := diag.GetSourceSummary("file2.go")
	if summary2.Error != 1 {
		t.Errorf("file2.go summary: Error=%d, want 1", summary2.Error)
	}
}

// Test DiagnosticStatusCmp Has methods
func TestDiagnosticStatusCmp_HasMethods(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()

	if diag.HasProblems() {
		t.Error("HasProblems() returned true initially")
	}

	if diag.HasErrors() {
		t.Error("HasErrors() returned true initially")
	}

	diag.AddDiagnostic(DiagnosticSeverityError)

	if !diag.HasErrors() {
		t.Error("HasErrors() returned false after adding error")
	}

	if !diag.HasProblems() {
		t.Error("HasProblems() returned false after adding error")
	}
}

// Test DiagnosticStatusCmp focus handling
func TestDiagnosticStatusCmp_Focus(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()

	if diag.Focused() {
		t.Fatal("Focused() returned true initially")
	}

	diag.Focus()
	if !diag.Focused() {
		t.Error("Focus() did not set focused flag")
	}

	diag.Blur()
	if diag.Focused() {
		t.Error("Blur() did not clear focused flag")
	}
}

// Test DiagnosticStatusCmp Update message handling
func TestDiagnosticStatusCmp_Update(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()

	// Test FocusGainMsg
	model, cmd := diag.Update(&render.FocusGainMsg{})
	if diag.Focused() != true {
		t.Error("Update(FocusGainMsg) did not set focused flag")
	}
	if cmd != nil {
		t.Errorf("Update(FocusGainMsg) returned non-nil command")
	}
	if model != diag {
		t.Error("Update() did not return the same model")
	}

	// Test BlurMsg
	model, cmd = diag.Update(&render.BlurMsg{})
	if diag.Focused() != false {
		t.Error("Update(BlurMsg) did not clear focused flag")
	}

	// Test other message types
	model, cmd = diag.Update("other")
	if cmd != nil {
		t.Errorf("Update(other) returned non-nil command")
	}
}

// Test DiagnosticStatusCmp View rendering
func TestDiagnosticStatusCmp_View(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()

	view := diag.View()

	if view == "" {
		t.Error("View() returned empty string")
	}

	// Empty state should show "No diagnostics"
	if !strings.Contains(view, "No diagnostics") {
		t.Error("View() with empty summary does not show 'No diagnostics'")
	}
}

// Test DiagnosticStatusCmp View with errors
func TestDiagnosticStatusCmp_ViewWithErrors(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()
	diag.AddDiagnostic(DiagnosticSeverityError)

	view := diag.View()

	if !strings.Contains(view, string(styles.ErrorIcon)) {
		t.Error("View() with error does not contain error icon")
	}

	if !strings.Contains(view, "1") {
		t.Error("View() with error does not contain error count")
	}
}

// Test DiagnosticStatusCmp View in compact mode
func TestDiagnosticStatusCmp_ViewCompact(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()
	diag.SetCompact(true)
	diag.AddDiagnostic(DiagnosticSeverityError)

	view := diag.View()

	if view == "" {
		t.Error("View() in compact mode returned empty string")
	}

	if strings.Contains(view, "No diagnostics") {
		t.Error("View() in compact mode should not show 'No diagnostics' when there are errors")
	}
}

// Test DiagnosticStatusCmp View with multiple severities
func TestDiagnosticStatusCmp_ViewMultipleSeverities(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()
	diag.AddDiagnostic(DiagnosticSeverityError)
	diag.AddDiagnostic(DiagnosticSeverityWarning)
	diag.AddDiagnostic(DiagnosticSeverityInfo)

	view := diag.View()

	if !strings.Contains(view, string(styles.ErrorIcon)) {
		t.Error("View() does not contain error icon")
	}

	if !strings.Contains(view, string(styles.WarningIcon)) {
		t.Error("View() does not contain warning icon")
	}

	if !strings.Contains(view, string(styles.InfoIcon)) {
		t.Error("View() does not contain info icon")
	}
}

// Test DiagnosticStatusCmp ShowHints
func TestDiagnosticStatusCmp_ShowHints(t *testing.T) {
	diag := NewDiagnosticStatus("test")
	diag.Init()
	diag.AddDiagnostic(DiagnosticSeverityHint)

	view1 := diag.View()

	// Hints should be shown by default
	if diag.ShowHints() != true {
		t.Error("ShowHints() should return true by default")
	}

	if !strings.Contains(view1, string(styles.HintIcon)) {
		t.Error("View() should show hint icon when showHints is true")
	}

	diag.SetShowHints(false)

	if diag.ShowHints() {
		t.Error("SetShowHints(false) did not clear showHints flag")
	}
}

// Test DiagnosticStatusCmp SetSource
func TestDiagnosticStatusCmp_SetSource(t *testing.T) {
	diag := NewDiagnosticStatus("old-source")
	diag.Init()

	diag.SetSource("new-source")

	if diag.Source() != "new-source" {
		t.Errorf("Source() after SetSource() = %v, want new-source", diag.Source())
	}
}
