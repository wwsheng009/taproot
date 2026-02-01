package pills

import (
	"testing"
)

func TestPillStatusString(t *testing.T) {
	tests := []struct {
		status   PillStatus
		expected string
	}{
		{PillStatusPending, "pending"},
		{PillStatusInProgress, "in-progress"},
		{PillStatusCompleted, "completed"},
		{PillStatusError, "error"},
		{PillStatusWarning, "warning"},
		{PillStatusInfo, "info"},
		{PillStatusNeutral, "neutral"},
	}

	for _, tt := range tests {
		result := tt.status.String()
		if result != tt.expected {
			t.Errorf("For status %v, expected '%s', got '%s'", tt.status, tt.expected, result)
		}
	}
}

func TestDefaultPillConfig(t *testing.T) {
	config := DefaultPillConfig()

	if !config.ShowItems {
		t.Error("Expected ShowItems to be true")
	}
	if !config.ShowCount {
		t.Error("Expected ShowCount to be true")
	}
	if config.CompactMode {
		t.Error("Expected CompactMode to be false")
	}
	if config.MaxItemWidth != 60 {
		t.Errorf("Expected MaxItemWidth 60, got %d", config.MaxItemWidth)
	}
	if !config.ShowIcons {
		t.Error("Expected ShowIcons to be true")
	}
	if config.InlineMode {
		t.Error("Expected InlineMode to be false")
	}
}

func TestNewPillList(t *testing.T) {
	pills := []*Pill{
		{
			ID:     "1",
			Label:  "Tasks",
			Count:  5,
			Status: PillStatusPending,
		},
	}

	pl := NewPillList(pills)

	if len(pl.GetPills()) != 1 {
		t.Errorf("Expected 1 pill, got %d", len(pl.GetPills()))
	}

	if pl.GetPills()[0].Label != "Tasks" {
		t.Errorf("Expected pill label 'Tasks', got '%s'", pl.GetPills()[0].Label)
	}

	if pl.width != 80 {
		t.Errorf("Expected width 80, got %d", pl.width)
	}

	if pl.focused {
		t.Error("Expected focused to be false")
	}
}

func TestAddPill(t *testing.T) {
	pl := NewPillList([]*Pill{})

	newPill := &Pill{
		ID:     "1",
		Label:  "New Pill",
		Count:  3,
		Status: PillStatusInfo,
	}

	pl.AddPill(newPill)

	if len(pl.GetPills()) != 1 {
		t.Errorf("Expected 1 pill, got %d", len(pl.GetPills()))
	}
}

func TestRemovePill(t *testing.T) {
	pills := []*Pill{
		{ID: "1", Label: "Pill 1", Count: 1, Status: PillStatusPending},
		{ID: "2", Label: "Pill 2", Count: 2, Status: PillStatusCompleted},
	}

	pl := NewPillList(pills)

	// Remove existing pill
	if !pl.RemovePill("1") {
		t.Error("Expected RemovePill to return true")
	}

	if len(pl.GetPills()) != 1 {
		t.Errorf("Expected 1 pill after removal, got %d", len(pl.GetPills()))
	}

	// Try to remove non-existent pill
	if pl.RemovePill("3") {
		t.Error("Expected RemovePill to return false for non-existent ID")
	}
}

func TestGetPill(t *testing.T) {
	pills := []*Pill{
		{ID: "1", Label: "Pill 1", Count: 1, Status: PillStatusPending},
	}

	pl := NewPillList(pills)

	// Get existing pill
	pill := pl.GetPill("1")
	if pill == nil {
		t.Error("Expected to find pill with ID '1'")
	}

	if pill.Label != "Pill 1" {
		t.Errorf("Expected label 'Pill 1', got '%s'", pill.Label)
	}

	// Get non-existent pill
	pill = pl.GetPill("2")
	if pill != nil {
		t.Error("Expected nil for non-existent pill")
	}
}

func TestToggleExpanded(t *testing.T) {
	pills := []*Pill{
		{
			ID:       "1",
			Label:    "Tasks",
			Count:    5,
			Status:   PillStatusPending,
			Expanded: false,
			Items:    []string{"Task 1", "Task 2"},
		},
	}

	pl := NewPillList(pills)

	// Initially not expanded
	if pl.GetPill("1").Expanded {
		t.Error("Expected pill to not be expanded initially")
	}

	// Toggle to expand
	if !pl.ToggleExpanded("1") {
		t.Error("Expected ToggleExpanded to return true")
	}

	if !pl.GetPill("1").Expanded {
		t.Error("Expected pill to be expanded after toggle")
	}

	// Toggle to collapse
	pl.ToggleExpanded("1")
	if pl.GetPill("1").Expanded {
		t.Error("Expected pill to be collapsed after second toggle")
	}

	// Try to toggle non-existent pill
	if pl.ToggleExpanded("2") {
		t.Error("Expected ToggleExpanded to return false for non-existent ID")
	}
}

func TestGetTotalCount(t *testing.T) {
	pills := []*Pill{
		{ID: "1", Label: "Pill 1", Count: 5, Status: PillStatusPending},
		{ID: "2", Label: "Pill 2", Count: 3, Status: PillStatusCompleted},
		{ID: "3", Label: "Pill 3", Count: 7, Status: PillStatusInProgress},
	}

	pl := NewPillList(pills)

	totalCount := pl.GetTotalCount()
	if totalCount != 15 {
		t.Errorf("Expected total count 15, got %d", totalCount)
	}
}

func TestGetCountByStatus(t *testing.T) {
	pills := []*Pill{
		{ID: "1", Label: "Pill 1", Count: 1, Status: PillStatusPending},
		{ID: "2", Label: "Pill 2", Count: 2, Status: PillStatusPending},
		{ID: "3", Label: "Pill 3", Count: 3, Status: PillStatusCompleted},
		{ID: "4", Label: "Pill 4", Count: 4, Status: PillStatusCompleted},
		{ID: "5", Label: "Pill 5", Count: 5, Status: PillStatusError},
	}

	pl := NewPillList(pills)

	counts := pl.GetCountByStatus()

	if counts[PillStatusPending] != 2 {
		t.Errorf("Expected 2 pending pills, got %d", counts[PillStatusPending])
	}

	if counts[PillStatusCompleted] != 2 {
		t.Errorf("Expected 2 completed pills, got %d", counts[PillStatusCompleted])
	}

	if counts[PillStatusError] != 1 {
		t.Errorf("Expected 1 error pill, got %d", counts[PillStatusError])
	}
}

func TestExpandAll(t *testing.T) {
	pills := []*Pill{
		{ID: "1", Label: "Pill 1", Count: 1, Status: PillStatusPending, Expanded: false},
		{ID: "2", Label: "Pill 2", Count: 2, Status: PillStatusCompleted, Expanded: false},
	}

	pl := NewPillList(pills)

	if pl.GetPill("1").Expanded || pl.GetPill("2").Expanded {
		t.Error("Expected pills to not be expanded initially")
	}

	pl.ExpandAll()

	if !pl.GetPill("1").Expanded || !pl.GetPill("2").Expanded {
		t.Error("Expected all pills to be expanded after ExpandAll()")
	}
}

func TestCollapseAll(t *testing.T) {
	pills := []*Pill{
		{ID: "1", Label: "Pill 1", Count: 1, Status: PillStatusPending, Expanded: true},
		{ID: "2", Label: "Pill 2", Count: 2, Status: PillStatusCompleted, Expanded: true},
	}

	pl := NewPillList(pills)

	if !pl.GetPill("1").Expanded || !pl.GetPill("2").Expanded {
		t.Error("Expected pills to be expanded initially")
	}

	pl.CollapseAll()

	if pl.GetPill("1").Expanded || pl.GetPill("2").Expanded {
		t.Error("Expected all pills to be collapsed after CollapseAll()")
	}
}

func TestFocusAndBlur(t *testing.T) {
	pl := NewPillList([]*Pill{})

	if pl.Focused() {
		t.Error("Expected component to not be focused initially")
	}

	pl.Focus()
	if !pl.Focused() {
		t.Error("Expected component to be focused after Focus()")
	}

	pl.Blur()
	if pl.Focused() {
		t.Error("Expected component to not be focused after Blur()")
	}
}

func TestSetWidth(t *testing.T) {
	pl := NewPillList([]*Pill{})

	pl.SetWidth(120)
	if pl.width != 120 {
		t.Errorf("Expected width 120, got %d", pl.width)
	}
}

func TestSetConfig(t *testing.T) {
	pl := NewPillList([]*Pill{})

	newConfig := PillConfig{
		ShowItems:    false,
		ShowCount:    false,
		CompactMode:  true,
		MaxItemWidth: 80,
		ShowIcons:    false,
		InlineMode:   true,
	}

	pl.SetConfig(newConfig)

	if pl.config.ShowItems {
		t.Error("Expected ShowItems to be false")
	}
	if pl.config.CompactMode != true {
		t.Error("Expected CompactMode to be true")
	}
	if pl.config.InlineMode != true {
		t.Error("Expected InlineMode to be true")
	}
}

func TestEmptyPillList(t *testing.T) {
	pl := NewPillList([]*Pill{})

	view := pl.View()
	if view == "" {
		t.Error("Expected non-empty view for empty pill list")
	}

	if len(pl.GetPills()) != 0 {
		t.Errorf("Expected 0 pills, got %d", len(pl.GetPills()))
	}

	totalCount := pl.GetTotalCount()
	if totalCount != 0 {
		t.Errorf("Expected total count 0, got %d", totalCount)
	}
}

func TestView(t *testing.T) {
	pills := []*Pill{
		{
			ID:     "1",
			Label:  "Tasks",
			Count:  5,
			Status: PillStatusPending,
			Items:  []string{"Task 1", "Task 2"},
		},
	}

	pl := NewPillList(pills)
	view := pl.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Check that label appears in view
	// This is a basic check - more detailed rendering tests would verify exact format
}

func TestPillStatusIcons(t *testing.T) {
	pl := NewPillList([]*Pill{})

	tests := []struct {
		status PillStatus
		icon   string
	}{
		{PillStatusPending, "☐"},
		{PillStatusInProgress, "⟳"},
		{PillStatusCompleted, "✓"},
		{PillStatusError, "×"},
		{PillStatusWarning, "⚠"},
		{PillStatusInfo, "ℹ"},
		{PillStatusNeutral, "•"},
	}

	for _, tt := range tests {
		icon := pl.getPillIcon(tt.status)
		if icon != tt.icon {
			t.Errorf("For status %v, expected icon '%s', got '%s'", tt.status, tt.icon, icon)
		}
	}
}

func TestPillStatusText(t *testing.T) {
	tests := []struct {
		status   PillStatus
		expected string
	}{
		{PillStatusPending, "pending"},
		{PillStatusInProgress, "in-progress"},
		{PillStatusCompleted, "completed"},
		{PillStatusError, "error"},
		{PillStatusWarning, "warning"},
		{PillStatusInfo, "info"},
		{PillStatusNeutral, "neutral"},
	}

	for _, tt := range tests {
		result := tt.status.String()
		if result != tt.expected {
			t.Errorf("For status %v, expected '%s', got '%s'", tt.status, tt.expected, result)
		}
	}
}

func TestPillFields(t *testing.T) {
	pill := &Pill{
		ID:       "test-id",
		Label:    "Test Pill",
		Count:    42,
		Status:   PillStatusInfo,
		Expanded: false,
		Items:    []string{"Item 1", "Item 2", "Item 3"},
	}

	if pill.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", pill.ID)
	}

	if pill.Label != "Test Pill" {
		t.Errorf("Expected label 'Test Pill', got %s", pill.Label)
	}

	if pill.Count != 42 {
		t.Errorf("Expected count 42, got %d", pill.Count)
	}

	if pill.Status != PillStatusInfo {
		t.Errorf("Expected status PillStatusInfo, got %v", pill.Status)
	}

	if pill.Expanded {
		t.Error("Expected expanded to be false")
	}

	if len(pill.Items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(pill.Items))
	}
}

func TestInlineMode(t *testing.T) {
	pills := []*Pill{
		{ID: "1", Label: "Tasks", Count: 5, Status: PillStatusPending},
		{ID: "2", Label: "Issues", Count: 3, Status: PillStatusError},
	}

	pl := NewPillList(pills)
	config := DefaultPillConfig()
	config.InlineMode = true
	pl.SetConfig(config)

	view := pl.View()
	if view == "" {
		t.Error("Expected non-empty view in inline mode")
	}
}
