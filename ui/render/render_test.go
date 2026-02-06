package render

import (
	"slices"
	"testing"
)

func TestEngineRegistry(t *testing.T) {
	// Test unregistered engine first
	t.Run("UnregisteredEngine", func(t *testing.T) {
		_, err := CreateEngine(EngineType(999), nil)
		if err == nil {
			t.Error("expected error for unregistered engine")
		}
	})

	// Create a dummy factory that returns the correct type
	dummyFactory := func(config *EngineConfig) Engine {
		// Return a direct engine but we'll check the type differently
		// The key is that factory returns something, not the specific type
		return NewDirectEngine(config)
	}

	// Test registration
	t.Run("RegisterEngine", func(t *testing.T) {
		RegisterEngine(EngineUltraviolet, dummyFactory)
		if !IsEngineRegistered(EngineUltraviolet) {
			t.Error("expected ultraviolet engine to be registered")
		}
	})

	// Test available engines
	t.Run("AvailableEngines", func(t *testing.T) {
		engines := AvailableEngines()
		if !slices.Contains(engines, EngineUltraviolet) {
			t.Error("expected ultraviolet in available engines")
		}
	})

	// Test creation
	t.Run("CreateEngine", func(t *testing.T) {
		engine, err := CreateEngine(EngineUltraviolet, nil)
		if err != nil {
			t.Errorf("failed to create engine: %v", err)
		}
		if engine == nil {
			t.Error("expected non-nil engine")
		}
		// Just verify it was created successfully
	})
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("expected non-nil config")
	}
	if config.EnableMouse {
		t.Error("expected mouse disabled by default")
	}
	if !config.EnableAltScreen {
		t.Error("expected alt screen enabled by default")
	}
	if config.EnableCursor {
		t.Error("expected cursor disabled by default")
	}
}

func TestDirectEngine(t *testing.T) {
	t.Run("NewDirectEngine", func(t *testing.T) {
		engine := NewDirectEngine(nil)
		if engine.Type() != EngineDirect {
			t.Errorf("expected direct type, got %v", engine.Type())
		}
	})

	t.Run("StartStop", func(t *testing.T) {
		engine := NewDirectEngine(nil).(*DirectEngine)
		model := NewTestModel("hello")

		if err := engine.Start(model); err != nil {
			t.Fatalf("failed to start: %v", err)
		}

		if !engine.Running() {
			t.Error("expected engine to be running")
		}

		if err := engine.Stop(); err != nil {
			t.Fatalf("failed to stop: %v", err)
		}

		if engine.Running() {
			t.Error("expected engine to be stopped")
		}
	})

	t.Run("Send", func(t *testing.T) {
		engine := NewDirectEngine(nil).(*DirectEngine)
		model := NewTestModel("hello")

		if err := engine.Start(model); err != nil {
			t.Fatalf("failed to start: %v", err)
		}
		defer engine.Stop()

		// Send a key message
		msg := KeyMsg{Key: "q"}
		if err := engine.Send(msg); err != nil {
			t.Errorf("failed to send message: %v", err)
		}

		// Check model was updated
		testModel := engine.Model().(*TestModel)
		if testModel.Content() != "quit" {
			t.Errorf("expected content 'quit', got %q", testModel.Content())
		}
	})

	t.Run("Resize", func(t *testing.T) {
		engine := NewDirectEngine(nil).(*DirectEngine)
		model := NewTestModel("hello")

		if err := engine.Start(model); err != nil {
			t.Fatalf("failed to start: %v", err)
		}
		defer engine.Stop()

		if err := engine.Resize(100, 50); err != nil {
			t.Errorf("failed to resize: %v", err)
		}

		testModel := engine.Model().(*TestModel)
		w, h := testModel.Size()
		if w != 100 || h != 50 {
			t.Errorf("expected 100x50, got %dx%d", w, h)
		}
	})

	t.Run("Output", func(t *testing.T) {
		engine := NewDirectEngine(nil).(*DirectEngine)
		model := NewTestModel("test")

		if err := engine.Start(model); err != nil {
			t.Fatalf("failed to start: %v", err)
		}
		defer engine.Stop()

		output := engine.Output()
		if output == "" {
			t.Error("expected non-empty output")
		}
	})

	t.Run("DoubleStart", func(t *testing.T) {
		engine := NewDirectEngine(nil).(*DirectEngine)
		model := NewTestModel("hello")

		if err := engine.Start(model); err != nil {
			t.Fatalf("failed to start: %v", err)
		}
		defer engine.Stop()

		if err := engine.Start(model); err == nil {
			t.Error("expected error when starting already running engine")
		}
	})
}

func TestKeyMsg(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		tests := []struct {
			key      KeyMsg
			expected string
		}{
			{KeyMsg{Key: "a"}, "a"},
			{KeyMsg{Key: "a", Ctrl: true}, "ctrl+a"},
			{KeyMsg{Key: "a", Alt: true}, "alt+a"},
			{KeyMsg{Key: "a", Shift: true}, "shift+a"},
			{KeyMsg{Key: "a", Ctrl: true, Alt: true}, "alt+ctrl+a"},
		}

		for _, tt := range tests {
			if got := tt.key.String(); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		}
	})

	t.Run("IsMouse", func(t *testing.T) {
		msg := KeyMsg{Key: "a"}
		if msg.IsMouse() {
			t.Error("expected IsMouse to return false")
		}
	})
}

func TestTestModel(t *testing.T) {
	t.Run("NewTestModel", func(t *testing.T) {
		model := NewTestModel("content")
		if model.Content() != "content" {
			t.Errorf("expected 'content', got %q", model.Content())
		}
		w, h := model.Size()
		if w != 80 || h != 24 {
			t.Errorf("expected 80x24, got %dx%d", w, h)
		}
	})

	t.Run("Init", func(t *testing.T) {
		model := NewTestModel("")
		if err := model.Init(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		model := NewTestModel("test")

		// Test WindowSizeMsg
		newModel, _ := model.Update(WindowSizeMsg{Width: 100, Height: 50})
		model = newModel.(*TestModel)
		w, h := model.Size()
		if w != 100 || h != 50 {
			t.Errorf("expected 100x50, got %dx%d", w, h)
		}

		// Test KeyMsg
		newModel, _ = model.Update(KeyMsg{Key: "q"})
		model = newModel.(*TestModel)
		if model.Content() != "quit" {
			t.Errorf("expected 'quit', got %q", model.Content())
		}
	})

	t.Run("View", func(t *testing.T) {
		model := NewTestModel("test")
		view := model.View()
		expected := "test [80x24]"
		if view != expected {
			t.Errorf("expected %q, got %q", expected, view)
		}
	})

	t.Run("SetContent", func(t *testing.T) {
		model := NewTestModel("")
		model.SetContent("new")
		if model.Content() != "new" {
			t.Errorf("expected 'new', got %q", model.Content())
		}
	})
}

func TestCommand(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		cmd := None()
		if cmd == nil {
			t.Error("expected non-nil none command")
		}
		if !IsNone(cmd) {
			t.Error("expected none command")
		}
	})

	t.Run("Batch", func(t *testing.T) {
		executed := 0
		cmd1 := Command(func() error {
			executed++
			return nil
		})
		cmd2 := Command(func() error {
			executed += 10
			return nil
		})

		batch := Batch(cmd1, cmd2)
		if batch == nil {
			t.Fatal("expected non-nil batch")
		}

		// Batch returns a BatchCmd (slice of Cmd)
		if cmds, ok := batch.(BatchCmd); ok {
			if len(cmds) != 2 {
				t.Errorf("expected 2 commands, got %d", len(cmds))
			}

			// We need to execute the commands manually since they are wrapped
			for _, c := range cmds {
				if cmd, ok := c.(Command); ok {
					if err := cmd(); err != nil {
						t.Errorf("unexpected error: %v", err)
					}
				} else {
					t.Errorf("expected Command type in batch, got %T", c)
				}
			}

			if executed != 11 {
				t.Errorf("expected executed 11, got %d", executed)
			}
		} else {
			t.Errorf("expected BatchCmd type from Batch, got %T", batch)
		}
	})

	t.Run("BatchEmpty", func(t *testing.T) {
		batch := Batch()
		if batch == nil {
			t.Error("expected non-nil batch command")
		}
		if !IsNone(batch) {
			t.Error("expected none command for empty batch")
		}
	})

	t.Run("Quit", func(t *testing.T) {
		cmd := Quit()
		if cmd == nil {
			t.Error("expected non-nil quit command")
		}
		if !IsQuit(cmd) {
			t.Error("expected IsQuit to return true")
		}

		// Test that None is not quit
		if IsQuit(None()) {
			t.Error("expected IsQuit to return false for None()")
		}

		// Test that regular command is not quit
		regularCmd := Command(func() error { return nil })
		if IsQuit(regularCmd) {
			t.Error("expected IsQuit to return false for regular command")
		}
	})
}
