package clipboard

import (
	"strings"
	"testing"
	"time"
)

// TestClipboardTypes - Test clipboard type enum
func TestClipboardTypes(t *testing.T) {
	tests := []struct {
		name string
		ct   ClipboardType
		want string
	}{
		{"OSC 52", ClipboardOSC52, "OSC 52"},
		{"Native", ClipboardNative, "Native"},
		{"Platform", ClipboardPlatform, "Platform"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.want {
				t.Errorf("ClipboardType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestClipboardData - Test clipboard data
func TestClipboardData(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		data := NewClipboardData(FormatText, "test")
		if data.Format != FormatText {
			t.Errorf("Expected format %v, got %v", FormatText, data.Format)
		}
		if data.Text != "test" {
			t.Errorf("Expected text 'test', got %v", data.Text)
		}
		if data.Timestamp.IsZero() {
			t.Error("Expected non-zero timestamp")
		}
	})

	t.Run("NewBytes", func(t *testing.T) {
		bytes := []byte("test")
		data := NewClipboardDataBytes(FormatImagePNG, bytes)
		if data.Format != FormatImagePNG {
			t.Errorf("Expected format %v, got %v", FormatImagePNG, data.Format)
		}
		if string(data.Bytes) != "test" {
			t.Errorf("Expected bytes 'test', got %v", string(data.Bytes))
		}
	})

	t.Run("IsImage", func(t *testing.T) {
		tests := []struct {
			format Format
			want   bool
		}{
			{FormatText, false},
			{FormatImagePNG, true},
			{FormatImageJPEG, true},
			{FormatImageGIF, true},
		}

		for _, tt := range tests {
			data := NewClipboardData(tt.format, "test")
			if got := data.IsImage(); got != tt.want {
				t.Errorf("ClipboardData.IsImage() for %v = %v, want %v", tt.format, got, tt.want)
			}
		}
	})

	t.Run("Size", func(t *testing.T) {
		data := NewClipboardData(FormatText, "hello")
		if got := data.Size(); got != 5 {
			t.Errorf("ClipboardData.Size() = %v, want 5", got)
		}

		data2 := NewClipboardDataBytes(FormatText, []byte("world"))
		if got := data2.Size(); got != 5 {
			t.Errorf("ClipboardData.Size() bytes = %v, want 5", got)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		data := NewClipboardData(FormatText, "")
		if !data.Empty() {
			t.Error("Expected empty data")
		}

		data2 := NewClipboardData(FormatText, "test")
		if data2.Empty() {
			t.Error("Expected non-empty data")
		}
	})
}

// TestClipboardState - Test clipboard state enum
func TestClipboardState(t *testing.T) {
	tests := []struct {
		name string
		cs   ClipboardState
		want string
	}{
		{"Idle", StateIdle, "Idle"},
		{"Reading", StateReading, "Reading"},
		{"Writing", StateWriting, "Writing"},
		{"Error", StateError, "Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.String(); got != tt.want {
				t.Errorf("ClipboardState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDefaultConfigs - Test default configuration
func TestDefaultConfigs(t *testing.T) {
	t.Run("OSC52Config", func(t *testing.T) {
		config := DefaultOSC52Config()
		if config.Selection != "c" {
			t.Errorf("Expected selection 'c', got %v", config.Selection)
		}
		if config.MaxSize != 100*1024 {
			t.Errorf("Expected max_size 100KB, got %v", config.MaxSize)
		}
		if !config.EncodeBase64 {
			t.Error("Expected Base64 encoding enabled")
		}
	})

	t.Run("NativeConfig", func(t *testing.T) {
		config := DefaultNativeConfig()
		if len(config.Formats) != 1 {
			t.Errorf("Expected 1 format, got %v", len(config.Formats))
		}
		if config.Timeout != 5*time.Second {
			t.Errorf("Expected timeout 5s, got %v", config.Timeout)
		}
	})

	t.Run("HistoryConfig", func(t *testing.T) {
		config := DefaultHistoryConfig()
		if config.MaxItems != 100 {
			t.Errorf("Expected max_items 100, got %v", config.MaxItems)
		}
		if config.Expiration != 24*time.Hour {
			t.Errorf("Expected expiration 24h, got %v", config.Expiration)
		}
		if !config.Deduplicate {
			t.Error("Expected deduplication enabled")
		}
	})
}

// TestHistoryManager - Test history manager
func TestHistoryManager(t *testing.T) {
	mockProvider := &mockClipboard{
		available: true,
		data:      NewClipboardData(FormatText, "test"),
	}
	config := HistoryConfig{
		MaxItems:    10,
		Expiration:  time.Hour,
		Deduplicate: true,
	}

	hm := NewHistoryManager(config, mockProvider)

	t.Run("AddEntry", func(t *testing.T) {
		data := NewClipboardData(FormatText, "hello")
		hm.AddEntry(data)

		if got := hm.EntryCount(); got != 1 {
			t.Errorf("Expected 1 entry, got %v", got)
		}

		// Check duplicate detection - adding same data twice shouldn't create duplicate
		hm.AddEntry(data)
		if got := hm.EntryCount(); got != 1 {
			t.Errorf("Expected still 1 entry (dedup on consecutive), got %v", got)
		}

		// Adding different data should create new entry
		data2 := NewClipboardData(FormatText, "world")
		hm.AddEntry(data2)
		if got := hm.EntryCount(); got != 2 {
			t.Errorf("Expected 2 entries, got %v", got)
		}
	})

	t.Run("GetEntry", func(t *testing.T) {
		hm.Clear()
		data1 := NewClipboardData(FormatText, "first")
		data2 := NewClipboardData(FormatText, "second")
		hm.AddEntry(data1)
		hm.AddEntry(data2)

		entry, err := hm.GetEntry(0)
		if err != nil {
			t.Fatalf("GetEntry(0) failed: %v", err)
		}
		if entry.Data.Text != "first" {
			t.Errorf("Expected 'first', got %v", entry.Data.Text)
		}

		_, err = hm.GetEntry(10)
		if err == nil {
			t.Error("Expected error for out of range index")
		}
	})

	t.Run("RemoveEntry", func(t *testing.T) {
		hm.Clear()
		hm.AddEntry(NewClipboardData(FormatText, "one"))
		hm.AddEntry(NewClipboardData(FormatText, "two"))
		hm.AddEntry(NewClipboardData(FormatText, "three"))

		err := hm.RemoveEntry(1)
		if err != nil {
			t.Fatalf("RemoveEntry(1) failed: %v", err)
		}

		if got := hm.EntryCount(); got != 2 {
			t.Errorf("Expected 2 entries, got %v", got)
		}

		// Check remaining entries
		entry, _ := hm.GetEntry(0)
		if entry.Data.Text != "one" {
			t.Errorf("Expected 'one', got %v", entry.Data.Text)
		}

		entry, _ = hm.GetEntry(1)
		if entry.Data.Text != "three" {
			t.Errorf("Expected 'three', got %v", entry.Data.Text)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		hm.AddEntry(NewClipboardData(FormatText, "test"))
		if got := hm.EntryCount(); got == 0 {
			t.Error("Expected at least 1 entry")
		}

		hm.Clear()
		if got := hm.EntryCount(); got != 0 {
			t.Errorf("Expected 0 entries after clear, got %v", got)
		}
	})

	t.Run("MaxItems", func(t *testing.T) {
		config := HistoryConfig{
			MaxItems:    3,
			Deduplicate: false,
		}
		hm2 := NewHistoryManager(config, mockProvider)

		for i := 0; i < 10; i++ {
			hm2.AddEntry(NewClipboardData(FormatText, string(rune('a'+i))))
		}

		if got := hm2.EntryCount(); got != 3 {
			t.Errorf("Expected 3 entries (max limit), got %v", got)
		}
	})
}

// TestOSC52 - Test OSC 52 functions
func TestOSC52(t *testing.T) {
	t.Run("DetectSupport", func(t *testing.T) {
		// This will depend on the environment
		available := DetectOSC52Support()
		// Just verify it doesn't crash
		_ = available
	})

	t.Run("EncodeDecode", func(t *testing.T) {
		text := "hello world"
		encoded := EncodeOSC52(text)
		decoded, err := DecodeOSC52(encoded)
		if err != nil {
			t.Fatalf("DecodeOSC52 failed: %v", err)
		}
		if decoded != text {
			t.Errorf("Expected %v, got %v", text, decoded)
		}
	})

	t.Run("ParseSequence", func(t *testing.T) {
		sequence := "\x1b]52;c;SGVsbG8=\x1b\\"
		selection, encoded, err := ParseOSC52Sequence(sequence)
		if err != nil {
			t.Fatalf("ParseOSC52Sequence failed: %v", err)
		}
		if selection != "c" {
			t.Errorf("Expected selection 'c', got %v", selection)
		}
		if encoded != "SGVsbG8=" {
			t.Errorf("Expected encoded 'SGVsbG8=', got %v", encoded)
		}
	})

	t.Run("ParseSequenceInvalid", func(t *testing.T) {
		tests := []struct {
			name      string
			sequence  string
			wantError bool
		}{
			{"Missing prefix", "52;c;data\x1b\\", true},
			{"Missing semicolon", "\x1b]52;cdata\x1b\\", true},
			{"Missing ST", "\x1b]52;c;data", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, _, err := ParseOSC52Sequence(tt.sequence)
				if (err != nil) != tt.wantError {
					t.Errorf("ParseOSC52Sequence() error = %v, wantError %v", err, tt.wantError)
				}
			})
		}
	})

	t.Run("ValidateData", func(t *testing.T) {
		tests := []struct {
			name        string
			encoded     string
			base64      bool
			wantError   bool
		}{
			{"Valid base64", "SGVsbG8=", true, false},
			{"Invalid base64", "Invalid!!", true, true},
			{"Empty string", "", true, false},
			{"Plain text", "hello", false, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := ValidateOSC52Data(tt.encoded, tt.base64)
				if (err != nil) != tt.wantError {
					t.Errorf("ValidateOSC52Data() error = %v, wantError %v", err, tt.wantError)
				}
			})
		}
	})

	t.Run("EstimateEncodedSize", func(t *testing.T) {
		textSize := 100
		encodedSize := EstimateEncodedSize(textSize, true)
		// Base64 encoding increases size by ~33%
		if encodedSize < textSize || encodedSize > textSize*150/100 {
			t.Errorf("Estimated size %v out of reasonable range for %v", encodedSize, textSize)
		}

		noEncodeSize := EstimateEncodedSize(textSize, false)
		if noEncodeSize != textSize {
			t.Errorf("Estimated size without encoding should be %v, got %v", textSize, noEncodeSize)
		}
	})

	t.Run("SizeLimitExceeded", func(t *testing.T) {
		text := strings.Repeat("a", 1000)
		maxSize := 500

		if !SizeLimitExceeded(text, maxSize) {
			t.Error("Expected size limit exceeded")
		}

		if SizeLimitExceeded(text, 2000) {
			t.Error("Expected size limit not exceeded")
		}
	})

	t.Run("TruncateToLimit", func(t *testing.T) {
		text := strings.Repeat("a", 1000)
		maxSize := 500

		truncated := TruncateToLimit(text, maxSize)
		if len(truncated) != maxSize {
			t.Errorf("Expected length %v, got %v", maxSize, len(truncated))
		}
	})
}

// TestTerminalInfo - Test terminal info functions
func TestTerminalInfo(t *testing.T) {
	t.Run("GetTerminalName", func(t *testing.T) {
		name := GetTerminalName()
		if name == "" {
			t.Error("Expected non-empty terminal name")
		}
	})

	t.Run("GetTerminalVersion", func(t *testing.T) {
		version := GetTerminalVersion()
		// Can be empty if not available
		_ = version
	})

	t.Run("GetPlatformName", func(t *testing.T) {
		provider := NewNativeProvider(DefaultNativeConfig())
		name := provider.GetPlatformName()
		if name == "Unknown" {
			// Unknown is acceptable
		} else if name == "" {
			t.Error("Expected non-empty platform name")
		}
	})
}

// TestErrorTypes - Test error types
func TestErrorTypes(t *testing.T) {
	errors := []error{
		ErrClipboardUnavailable,
		ErrClipboardLocked,
		ErrInvalidFormat,
		ErrDataTooLarge,
		ErrEmptyData,
		ErrTerminalNotSupported,
	}

	for _, err := range errors {
		if err == nil {
			t.Error("Expected non-nil error")
		}
		if err.Error() == "" {
			t.Error("Expected non-empty error message")
		}
	}
}

// Mock clipboard for testing
type mockClipboard struct {
	available bool
	data      *ClipboardData
}

func (m *mockClipboard) Available() bool {
	return m.available
}

func (m *mockClipboard) Read(format Format) (*ClipboardData, error) {
	if !m.available {
		return nil, ErrClipboardUnavailable
	}
	return m.data, nil
}

func (m *mockClipboard) Write(data *ClipboardData) error {
	if !m.available {
		return ErrClipboardUnavailable
	}
	m.data = data
	return nil
}

func (m *mockClipboard) Clear() error {
	if !m.available {
		return ErrClipboardUnavailable
	}
	m.data = nil
	return nil
}

func (m *mockClipboard) Formats() []Format {
	return []Format{FormatText}
}

func (m *mockClipboard) Type() ClipboardType {
	return ClipboardNative
}
