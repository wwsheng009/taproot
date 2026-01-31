package progress

import (
	"testing"

	"github.com/wwsheng009/taproot/ui/render"
)

func TestNewProgressBar(t *testing.T) {
	pb := NewProgressBar(100)

	if pb == nil {
		t.Fatal("NewProgressBar returned nil")
	}

	if pb.total != 100 {
		t.Errorf("Expected total 100, got %f", pb.total)
	}

	if pb.current != 0 {
		t.Errorf("Expected current 0, got %f", pb.current)
	}
}

func TestProgressBar_Init(t *testing.T) {
	pb := NewProgressBar(100)

	err := pb.Init()
	if err != nil {
		t.Fatalf("Init() returned error: %v", err)
	}

	if !pb.initialized {
		t.Error("ProgressBar should be initialized after Init()")
	}
}

func TestProgressBar_ImplementsModel(t *testing.T) {
	pb := NewProgressBar(100)
	var _ render.Model = pb
}

func TestProgressBar_SetCurrent(t *testing.T) {
	pb := NewProgressBar(100)

	pb.SetCurrent(50)

	if pb.current != 50 {
		t.Errorf("Expected current 50, got %f", pb.current)
	}
}

func TestProgressBar_Add(t *testing.T) {
	pb := NewProgressBar(100)
	pb.SetCurrent(25)

	pb.Add(25)

	if pb.current != 50 {
		t.Errorf("Expected current 50, got %f", pb.current)
	}
}

func TestProgressBar_Increment(t *testing.T) {
	pb := NewProgressBar(100)
	pb.SetCurrent(49)

	pb.Increment()

	if pb.current != 50 {
		t.Errorf("Expected current 50, got %f", pb.current)
	}
}

func TestProgressBar_SetTotal(t *testing.T) {
	pb := NewProgressBar(100)

	pb.SetTotal(200)

	if pb.total != 200 {
		t.Errorf("Expected total 200, got %f", pb.total)
	}
}

func TestProgressBar_SetCurrentDoesNotExceedTotal(t *testing.T) {
	pb := NewProgressBar(100)

	pb.SetCurrent(150)

	if pb.current != 100 {
		t.Errorf("Expected current to be clamped to 100, got %f", pb.current)
	}
}

func TestProgressBar_SetCurrentDoesNotGoNegative(t *testing.T) {
	pb := NewProgressBar(100)
	pb.SetCurrent(50)

	pb.SetCurrent(-10)

	if pb.current < 0 {
		t.Errorf("Expected current to not be negative, got %f", pb.current)
	}
}

func TestProgressBar_SetLabel(t *testing.T) {
	pb := NewProgressBar(100)

	pb.SetLabel("Loading...")

	if pb.label != "Loading..." {
		t.Errorf("Expected label 'Loading...', got '%s'", pb.label)
	}
}

func TestProgressBar_Completed(t *testing.T) {
	pb := NewProgressBar(100)

	if pb.Completed() {
		t.Error("Progress should not be complete initially")
	}

	pb.SetCurrent(100)

	if !pb.Completed() {
		t.Error("Progress should be complete at 100%")
	}
}

func TestProgressBar_Percent(t *testing.T) {
	pb := NewProgressBar(100)
	pb.SetCurrent(50)

	percent := pb.Percent()

	if percent != 50 {
		t.Errorf("Expected percent 50, got %f", percent)
	}
}

func TestProgressBar_Reset(t *testing.T) {
	pb := NewProgressBar(100)
	pb.SetCurrent(75)

	pb.Reset()

	if pb.current != 0 {
		t.Errorf("Expected current to be 0 after reset, got %f", pb.current)
	}
}

func TestProgressBar_View(t *testing.T) {
	pb := NewProgressBar(100)
	pb.SetCurrent(50)
	pb.Init()

	view := pb.View()

	if len(view) == 0 {
		t.Error("View() returned empty string")
	}
}

func TestProgressBar_ViewWithLabel(t *testing.T) {
	pb := NewProgressBar(100)
	pb.SetCurrent(50)
	pb.SetLabel("Progress:")
	pb.Init()

	view := pb.View()

	if len(view) == 0 {
		t.Error("View() returned empty string")
	}
	if !contains(view, "Progress:") {
		t.Error("View should contain label 'Progress:'")
	}
}

func TestProgressBar_Style(t *testing.T) {
	pb := NewProgressBar(100)
	style := DefaultProgressBarStyle()

	pb.SetStyle(style)

	if pb.style != style {
		t.Error("Style should be set")
	}
}

func TestNewSpinner(t *testing.T) {
	sp := NewSpinner()

	if sp == nil {
		t.Fatal("NewSpinner returned nil")
	}

	if sp.state != 0 {
		t.Errorf("Expected state 0, got %d", sp.state)
	}
}

func TestSpinner_Init(t *testing.T) {
	sp := NewSpinner()

	err := sp.Init()
	if err != nil {
		t.Fatalf("Init() returned error: %v", err)
	}

	if !sp.initialized {
		t.Error("Spinner should be initialized after Init()")
	}

	if !sp.started {
		t.Error("Spinner should be started after Init()")
	}
}

func TestSpinner_ImplementsModel(t *testing.T) {
	sp := NewSpinner()
	var _ render.Model = sp
}

func TestSpinner_Update_TickMsg(t *testing.T) {
	sp := NewSpinner()
	sp.Init()

	initialState := sp.state
	msg := &TickMsg{id: sp.id}
	sp.Update(msg)

	if sp.state <= initialState {
		t.Error("State should be incremented after TickMsg")
	}
}

func TestSpinner_Update_InvalidTickMsg(t *testing.T) {
	sp := NewSpinner()
	sp.Init()

	initialState := sp.state
	msg := &TickMsg{id: sp.id + 1}
	sp.Update(msg)

	if sp.state != initialState {
		t.Error("State should not change for invalid TickMsg ID")
	}
}

func TestSpinner_Update_RenderTickMsg(t *testing.T) {
	sp := NewSpinner()
	sp.Init()

	initialState := sp.state
	msg := &render.TickMsg{}
	sp.Update(msg)

	if sp.state <= initialState {
		t.Error("State should be incremented after render.TickMsg")
	}
}

func TestSpinner_View(t *testing.T) {
	sp := NewSpinner()
	sp.Init()

	view := sp.View()

	if len(view) == 0 {
		t.Error("View() returned empty string")
	}
}

func TestSpinner_ViewWithLabel(t *testing.T) {
	sp := NewSpinner()
	sp.SetLabel("Loading...")
	sp.Init()

	view := sp.View()

	if !contains(view, "Loading...") {
		t.Error("View should contain label 'Loading...'")
	}
}

func TestSpinner_Stop(t *testing.T) {
	sp := NewSpinner()
	sp.Init()

	sp.Stop()

	if !sp.stopped {
		t.Error("Spinner should be stopped after Stop()")
	}

	if sp.Running() {
		t.Error("Spinner should not be running when stopped")
	}
}

func TestSpinner_Reset(t *testing.T) {
	sp := NewSpinner()
	sp.Init()
	sp.Update(&TickMsg{id: sp.id})

	sp.Reset()

	if sp.state != 0 {
		t.Errorf("Expected state 0 after reset, got %d", sp.state)
	}

	if sp.stopped {
		t.Error("Spinner should not be stopped after reset")
	}
}

func TestSpinner_Running(t *testing.T) {
	sp := NewSpinner()

	if sp.Running() {
		t.Error("Spinner should not be running before Init()")
	}

	sp.Init()

	if !sp.Running() {
		t.Error("Spinner should be running after Init()")
	}
}

func TestSpinner_SetLabel(t *testing.T) {
	sp := NewSpinner()

	sp.SetLabel("Please wait")

	if sp.Label() != "Please wait" {
		t.Errorf("Expected label 'Please wait', got '%s'", sp.Label())
	}
}

func TestSpinner_SetType(t *testing.T) {
	sp := NewSpinner()

	sp.SetType(SpinnerTypeLine)

	if sp.Type() != SpinnerTypeLine {
		t.Errorf("Expected type SpinnerTypeLine, got %v", sp.Type())
	}
}

func TestSpinner_SetFPS(t *testing.T) {
	sp := NewSpinner()

	sp.SetFPS(20)

	if sp.FPS() != 20 {
		t.Errorf("Expected FPS 20, got %d", sp.FPS())
	}
}

func TestSpinner_SetFPS_ClampsToMin(t *testing.T) {
	sp := NewSpinner()

	sp.SetFPS(-5)

	if sp.FPS() <= 0 {
		t.Error("FPS should be clamped to minimum of 10")
	}
}

func TestSpinner_frames(t *testing.T) {
	sp := NewSpinner()

	frames := sp.frames()

	if len(frames) == 0 {
		t.Error("frames() should return non-empty slice")
	}
}

func TestSpinner_frames_Lines(t *testing.T) {
	sp := NewSpinner()
	sp.SetType(SpinnerTypeLine)

	frames := sp.frames()

	if len(frames) == 0 {
		t.Error("frames() should return non-empty slice for SpinnerTypeLine")
	}
}

func TestSpinner_StateWraps(t *testing.T) {
	sp := NewSpinner()
	sp.Init()

	// Advance beyond frame count
	for i := 0; i < 100; i++ {
		sp.Update(&TickMsg{id: sp.id})
	}

	// State should wrap around
	if sp.state >= len(sp.frames()) {
		t.Errorf("State should be less than frame count, got %d", sp.state)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
