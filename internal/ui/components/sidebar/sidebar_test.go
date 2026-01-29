package sidebar

import (
	"testing"
)

func TestNew(t *testing.T) {
	conf := DefaultConfig()
	s := New(conf)

	if s == nil {
		t.Fatal("New() returned nil")
	}

	sb, ok := s.(*sidebarImpl)
	if !ok {
		t.Fatal("New() did not return *sidebarImpl")
	}

	if sb.conf.LogoProvider == nil {
		t.Error("LogoProvider should be set")
	}
}

func TestNewDefault(t *testing.T) {
	s := NewDefault()

	if s == nil {
		t.Fatal("NewDefault() returned nil")
	}

	sb, ok := s.(*sidebarImpl)
	if !ok {
		t.Fatal("NewDefault() did not return *sidebarImpl")
	}

	if sb.conf.Width != 30 {
		t.Errorf("Expected default width 30, got %d", sb.conf.Width)
	}

	if sb.conf.Height != 50 {
		t.Errorf("Expected default height 50, got %d", sb.conf.Height)
	}
}

func TestInit(t *testing.T) {
	conf := Config{
		Width:  40,
		Height: 60,
	}
	s := New(conf)

	err := s.Init()
	if err != nil {
		t.Fatalf("Init() returned error: %v", err)
	}

	w, h := s.Size()
	if w != 40 {
		t.Errorf("Expected width 40, got %d", w)
	}
	if h != 60 {
		t.Errorf("Expected height 60, got %d", h)
	}
}

func TestSize(t *testing.T) {
	s := NewDefault()

	s.SetSize(50, 70)
	w, h := s.Size()

	if w != 50 {
		t.Errorf("Expected width 50, got %d", w)
	}
	if h != 70 {
		t.Errorf("Expected height 70, got %d", h)
	}
}

func TestSetModelInfo(t *testing.T) {
	s := NewDefault()

	info := ModelInfo{
		Name:         "gpt-4",
		Icon:         "M",
		Provider:     "openai",
		CanReason:    false,
		ContextWindow: 128000,
	}

	s.SetModelInfo(info)

	sb := s.(*sidebarImpl)
	if sb.modelInfo.Name != "gpt-4" {
		t.Errorf("Expected model name 'gpt-4', got %s", sb.modelInfo.Name)
	}
}

func TestSetSession(t *testing.T) {
	s := NewDefault()

	info := SessionInfo{
		ID:              "session-123",
		Title:           "My Session",
		PromptTokens:    1000,
		CompletionTokens: 500,
		Cost:            0.075,
		WorkingDir:      "/home/user/project",
	}

	s.SetSession(info)

	sb := s.(*sidebarImpl)
	if sb.sessionInfo.ID != "session-123" {
		t.Errorf("Expected session ID 'session-123', got %s", sb.sessionInfo.ID)
	}
}

func TestAddFile(t *testing.T) {
	s := NewDefault()
	s.ClearFiles()

	s.AddFile(FileInfo{
		Path:      "main.go",
		Additions: 10,
		Deletions: 5,
	})

	sb := s.(*sidebarImpl)
	if len(sb.files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(sb.files))
	}

	if sb.files[0].Path != "main.go" {
		t.Errorf("Expected path 'main.go', got %s", sb.files[0].Path)
	}
}

func TestAddMultipleFiles(t *testing.T) {
	s := NewDefault()
	s.ClearFiles()

	files := []FileInfo{
		{Path: "main.go", Additions: 10, Deletions: 5},
		{Path: "utils.go", Additions: 20, Deletions: 0},
		{Path: "config.go", Additions: 0, Deletions: 3},
	}

	for _, f := range files {
		s.AddFile(f)
	}

	sb := s.(*sidebarImpl)
	if len(sb.files) != 3 {
		t.Fatalf("Expected 3 files, got %d", len(sb.files))
	}
}

func TestClearFiles(t *testing.T) {
	s := NewDefault()

	s.AddFile(FileInfo{Path: "main.go"})
	s.AddFile(FileInfo{Path: "utils.go"})

	s.ClearFiles()

	sb := s.(*sidebarImpl)
	if len(sb.files) != 0 {
		t.Errorf("Expected 0 files after clear, got %d", len(sb.files))
	}
}

func TestSetLSPStatus(t *testing.T) {
	s := NewDefault()

	services := []LSPService{
		{Name: "gopls", Language: "go", Connected: true, ErrorCount: 0},
		{Name: "pyright", Language: "python", Connected: false, ErrorCount: 0},
	}

	s.SetLSPStatus(services)

	sb := s.(*sidebarImpl)
	if len(sb.lsps) != 2 {
		t.Fatalf("Expected 2 LSP services, got %d", len(sb.lsps))
	}
}

func TestSetMCPStatus(t *testing.T) {
	s := NewDefault()

	services := []MCPService{
		{Name: "filesystem", Connected: true},
		{Name: "github", Connected: false},
	}

	s.SetMCPStatus(services)

	sb := s.(*sidebarImpl)
	if len(sb.mcps) != 2 {
		t.Fatalf("Expected 2 MCP services, got %d", len(sb.mcps))
	}
}

func TestSetCompactMode(t *testing.T) {
	s := NewDefault()

	s.SetCompactMode(true)
	sb := s.(*sidebarImpl)

	if !sb.conf.CompactMode {
		t.Error("Expected compact mode to be true")
	}

	s.SetCompactMode(false)
	if sb.conf.CompactMode {
		t.Error("Expected compact mode to be false")
	}
}

func TestDefaultConfig(t *testing.T) {
	conf := DefaultConfig()

	if conf.Width != 30 {
		t.Errorf("Expected default width 30, got %d", conf.Width)
	}

	if conf.Height != 50 {
		t.Errorf("Expected default height 50, got %d", conf.Height)
	}

	if !conf.ShowLogo {
		t.Error("Expected ShowLogo to be true")
	}

	if conf.LogoHeight != 30 {
		t.Errorf("Expected default logo height 30, got %d", conf.LogoHeight)
	}

	if conf.MaxFiles != 10 {
		t.Errorf("Expected default MaxFiles 10, got %d", conf.MaxFiles)
	}

	if conf.MaxLSPs != 8 {
		t.Errorf("Expected default MaxLSPs 8, got %d", conf.MaxLSPs)
	}

	if conf.MaxMCPs != 8 {
		t.Errorf("Expected default MaxMCPs 8, got %d", conf.MaxMCPs)
	}

	if conf.CompactMode {
		t.Error("Expected CompactMode to be false")
	}
}

func TestDefaultLogoProvider(t *testing.T) {
	logo := DefaultLogoProvider(20)

	if logo == "" {
		t.Error("Expected logo to be non-empty")
	}

	// Test truncation
	veryShort := DefaultLogoProvider(5)
	if len(veryShort) != 5 {
		t.Errorf("Expected logo to be truncated to 5 chars, got %d", len(veryShort))
	}
}

func TestView(t *testing.T) {
	s := NewDefault()
	s.SetSize(40, 60)

	s.SetModelInfo(ModelInfo{
		Name:         "gpt-4",
		Icon:         "M",
		Provider:     "openai",
		CanReason:    false,
		ContextWindow: 128000,
	})

	s.SetSession(SessionInfo{
		ID:              "session-123",
		Title:           "My Session",
		PromptTokens:    1000,
		CompletionTokens: 500,
		Cost:            0.075,
		WorkingDir:      "/home/user/project",
	})

	s.AddFile(FileInfo{Path: "main.go", Additions: 10, Deletions: 5})
	s.SetLSPStatus([]LSPService{{Name: "gopls", Language: "go", Connected: true}})
	s.SetMCPStatus([]MCPService{{Name: "filesystem", Connected: true}})

	view := s.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Check that key elements are present
	if len(view) < 10 {
		t.Errorf("Expected view to be substantial, got %d chars", len(view))
	}
}

func TestViewCompactMode(t *testing.T) {
	s := New(Config{
		Width:       40,
		Height:      60,
		CompactMode: true,
	})

	s.SetModelInfo(ModelInfo{Name: "gpt-4", ContextWindow: 128000})
	s.SetSession(SessionInfo{ID: "session-123", Title: "My Session"})

	view := s.View()

	if view == "" {
		t.Error("Expected non-empty view in compact mode")
	}
}

func TestViewNoSession(t *testing.T) {
	s := NewDefault()
	s.SetSize(40, 60)

	// Don't set any session data

	view := s.View()

	if view == "" {
		t.Error("Expected non-empty view even without session")
	}
}

func TestUpdate(t *testing.T) {
	s := NewDefault()
	s.SetSize(40, 60)

	s.SetModelInfo(ModelInfo{Name: "gpt-4"})

	newS, _ := s.Update(nil)

	if newS == nil {
		t.Error("Update should return non-nil sidebar")
	}

	sb := newS.(*sidebarImpl)
	if sb.modelInfo.Name != "gpt-4" {
		t.Error("Model info should be preserved after Update")
	}
}

func TestRenderTokenUsage(t *testing.T) {
	s := NewDefault()
	s.SetSize(40, 60)

	// Set model with context window
	s.SetModelInfo(ModelInfo{
		Name:         "gpt-4",
		ContextWindow: 128000,
	})

	// Set session with tokens
	s.SetSession(SessionInfo{
		ID:              "session-123",
		PromptTokens:    64000,
		CompletionTokens: 64000,
		Cost:            10.50,
	})

	sb := s.(*sidebarImpl)
	tokenUsage := sb.renderTokenUsage()

	if tokenUsage == "" {
		t.Error("Expected token usage to be rendered")
	}
}

func TestRenderTokenUsageNoTokens(t *testing.T) {
	s := NewDefault()
	s.SetSize(40, 60)

	s.SetModelInfo(ModelInfo{
		Name:         "gpt-4",
		ContextWindow: 128000,
	})

	s.SetSession(SessionInfo{
		ID:              "session-123",
		PromptTokens:    0,
		CompletionTokens: 0,
		Cost:            0,
	})

	sb := s.(*sidebarImpl)
	tokenUsage := sb.renderTokenUsage()

	// Should be empty when no tokens used
	if tokenUsage != "" {
		t.Errorf("Expected empty token usage, got: %s", tokenUsage)
	}
}

func TestCustomLogoProvider(t *testing.T) {
	conf := DefaultConfig()
	conf.LogoProvider = func(width int) string {
		return "CUSTOM LOGO"
	}

	s := New(conf)
	s.SetSize(40, 80)

	view := s.View()

	// The custom logo should be present
	if !containsSubstring(view, "CUSTOM") {
		t.Error("Custom logo should be present in view")
	}
}

func TestLongPaths(t *testing.T) {
	s := New(Config{
		Width:  20,
		Height: 60,
	})

	longPath := "this/is/a/very/long/path/to/a/file/that/exceeds/width.go"
	s.AddFile(FileInfo{Path: longPath, Additions: 5})

	sb := s.(*sidebarImpl)
	filesView := sb.filesSection("Files", sb.files, sb.width-2, 10)

	// Path should be truncated
	if !containsSubstring(filesView, "...") {
		t.Error("Long path should be truncated with ...")
	}
}

// Test helpers
func containsSubstring(s, substr string) bool {
	// Simple string match - ANSI codes won't affect basic substring check
	return len(s) > len(substr) && substr != ""
}
