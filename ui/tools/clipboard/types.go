package clipboard

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrClipboardUnavailable - Clipboard not available
	ErrClipboardUnavailable = errors.New("clipboard: not available")
	// ErrClipboardLocked - Clipboard locked by another process
	ErrClipboardLocked = errors.New("clipboard: locked")
	// ErrInvalidFormat - Invalid clipboard format
	ErrInvalidFormat = errors.New("clipboard: invalid format")
	// ErrDataTooLarge - Data exceeds size limit
	ErrDataTooLarge = errors.New("clipboard: data too large")
	// ErrEmptyData - Empty data
	ErrEmptyData = errors.New("clipboard: empty data")
	// ErrTerminalNotSupported - Terminal doesn't support OSC 52
	ErrTerminalNotSupported = errors.New("clipboard: terminal doesn't support OSC 52")
)

// ClipboardType - Clipboard type
type ClipboardType int

const (
	// ClipboardOSC52 - OSC 52 terminal clipboard
	ClipboardOSC52 ClipboardType = iota
	// ClipboardNative - Native OS clipboard
	ClipboardNative
	// ClipboardPlatform - Platform-specific (auto-detect)
	ClipboardPlatform
)

func (ct ClipboardType) String() string {
	switch ct {
	case ClipboardOSC52:
		return "OSC 52"
	case ClipboardNative:
		return "Native"
	case ClipboardPlatform:
		return "Platform"
	default:
		return "Unknown"
	}
}

// Format - Clipboard data format
type Format string

const (
	// FormatText - Plain text format
	FormatText Format = "text/plain"
	// FormatHTML - HTML format
	FormatHTML Format = "text/html"
	// FormatRTF - Rich Text Format
	FormatRTF Format = "text/rtf"
	// FormatImagePNG - PNG image format
	FormatImagePNG Format = "image/png"
	// FormatImageJPEG - JPEG image format
	FormatImageJPEG Format = "image/jpeg"
	// FormatImageGIF - GIF image format
	FormatImageGIF Format = "image/gif"
)

// ClipboardData - Clipboard data container
type ClipboardData struct {
	Format Format
	Text   string
	Bytes  []byte
	Timestamp time.Time
}

// NewClipboardData - Create clipboard data
func NewClipboardData(format Format, text string) *ClipboardData {
	return &ClipboardData{
		Format:    format,
		Text:      text,
		Timestamp: time.Now(),
	}
}

// NewClipboardDataBytes - Create clipboard data from bytes
func NewClipboardDataBytes(format Format, data []byte) *ClipboardData {
	return &ClipboardData{
		Format:    format,
		Bytes:     data,
		Timestamp: time.Now(),
	}
}

// IsImage - Check if data is image
func (cd *ClipboardData) IsImage() bool {
	switch cd.Format {
	case FormatImagePNG, FormatImageJPEG, FormatImageGIF:
		return true
	default:
		return false
	}
}

// Size - Get data size in bytes
func (cd *ClipboardData) Size() int {
	if len(cd.Bytes) > 0 {
		return len(cd.Bytes)
	}
	return len(cd.Text)
}

// Empty - Check if data is empty
func (cd *ClipboardData) Empty() bool {
	return cd.Size() == 0
}

// ClipboardState - Clipboard state
type ClipboardState int

const (
	// StateIdle - Clipboard idle
	StateIdle ClipboardState = iota
	// StateReading - Reading from clipboard
	StateReading
	// StateWriting - Writing to clipboard
	StateWriting
	// StateError - Error state
	StateError
)

func (cs ClipboardState) String() string {
	switch cs {
	case StateIdle:
		return "Idle"
	case StateReading:
		return "Reading"
	case StateWriting:
		return "Writing"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// OSC52Config - OSC 52 configuration
type OSC52Config struct {
	// Selection - Clipboard selection (c=clipboard, p=primary, q=secondary)
	Selection string
	// MaxSize - Maximum size for OSC 52 (default 100KB)
	MaxSize int
	// Truncate - Truncate if exceeds max size
	Truncate bool
	// EncodeBase64 - Use base64 encoding
	EncodeBase64 bool
}

// DefaultOSC52Config - Default OSC 52 config
func DefaultOSC52Config() OSC52Config {
	return OSC52Config{
		Selection:    "c",
		MaxSize:      100 * 1024, // 100KB
		Truncate:     false,
		EncodeBase64: true,
	}
}

// NativeConfig - Native clipboard configuration
type NativeConfig struct {
	// Formats - Supported formats
	Formats []Format
	// Timeout - Operation timeout
	Timeout time.Duration
	// RetryCount - Retry count on failure
	RetryCount int
	// RetryDelay - Delay between retries
	RetryDelay time.Duration
}

// DefaultNativeConfig - Default native clipboard config
func DefaultNativeConfig() NativeConfig {
	return NativeConfig{
		Formats:    []Format{FormatText},
		Timeout:    5 * time.Second,
		RetryCount: 3,
		RetryDelay: 100 * time.Millisecond,
	}
}

// HistoryConfig - Clipboard history configuration
type HistoryConfig struct {
	// MaxItems - Maximum history items
	MaxItems int
	// Expiration - History entry expiration time (0 = no expiration)
	Expiration time.Duration
	// PersistPath - Path to persist history (empty = in-memory only)
	PersistPath string
	// Deduplicate - Remove duplicate consecutive entries
	Deduplicate bool
}

// DefaultHistoryConfig - Default history config
func DefaultHistoryConfig() HistoryConfig {
	return HistoryConfig{
		MaxItems:    100,
		Expiration:  24 * time.Hour,
		PersistPath: "",
		Deduplicate: true,
	}
}

// Clipboard - Clipboard interface
type Clipboard interface {
	// Available - Check if clipboard is available
	Available() bool
	// Read - Read from clipboard
	Read(format Format) (*ClipboardData, error)
	// Write - Write to clipboard
	Write(data *ClipboardData) error
	// Clear - Clear clipboard
	Clear() error
	// Formats - Get supported formats
	Formats() []Format
	// Type - Get clipboard type
	Type() ClipboardType
}

// ClipboardWithHistory - Clipboard with history tracking
type ClipboardWithHistory interface {
	Clipboard
	// History - Get clipboard history
	History() []*ClipboardData
	// AddToHistory - Add entry to history
	AddToHistory(data *ClipboardData)
	// ClearHistory - Clear history
	ClearHistory()
	// RestoreHistory - Restore from history
	RestoreFromHistory(index int) error
	// PersistHistory - Persist history to disk
	PersistHistory() error
	// LoadHistory - Load history from disk
	LoadHistory() error
}

// Provider - Provider interface for clipboard operations
type Provider interface {
	// Read - Read from clipboard
	Read(ctx context.Context) (*ClipboardData, error)
	// Write - Write to clipboard
	Write(ctx context.Context, data *ClipboardData) error
	// Clear - Clear clipboard
	Clear(ctx context.Context) error
	// Available - Check if available
	Available() bool
}

// HistoryEntry - History entry
type HistoryEntry struct {
	Data      *ClipboardData
	Added     time.Time
	ExpiresAt time.Time
	Index     int
}

// HistoryManager - History manager
type HistoryManager struct {
	config   HistoryConfig
	entries  []*HistoryEntry
	provider Clipboard
}

// NewHistoryManager - Create history manager
func NewHistoryManager(config HistoryConfig, provider Clipboard) *HistoryManager {
	return &HistoryManager{
		config:   config,
		provider: provider,
		entries:  make([]*HistoryEntry, 0, config.MaxItems),
	}
}

// AddEntry - Add entry to history
func (hm *HistoryManager) AddEntry(data *ClipboardData) {
	entry := &HistoryEntry{
		Data:  data,
		Added: time.Now(),
		Index: len(hm.entries),
	}

	if hm.config.Expiration > 0 {
		entry.ExpiresAt = entry.Added.Add(hm.config.Expiration)
	}

	// Deduplicate consecutive entries
	if hm.config.Deduplicate && len(hm.entries) > 0 {
		lastEntry := hm.entries[len(hm.entries)-1]
		if lastEntry.Data.Format == data.Format && lastEntry.Data.Text == data.Text {
			return
		}
	}

	// Add entry
	hm.entries = append(hm.entries, entry)

	// Trim if exceeds max items
	if len(hm.entries) > hm.config.MaxItems {
		hm.entries = hm.entries[len(hm.entries)-hm.config.MaxItems:]
	}

	// Re-index entries
	for i, e := range hm.entries {
		e.Index = i
	}
}

// GetEntry - Get history entry by index
func (hm *HistoryManager) GetEntry(index int) (*HistoryEntry, error) {
	if index < 0 || index >= len(hm.entries) {
		return nil, errors.New("history: index out of range")
	}
	return hm.entries[index], nil
}

// GetAllEntries - Get all entries
func (hm *HistoryManager) GetAllEntries() []*HistoryEntry {
	// Clean expired entries
	now := time.Now()
	clean := make([]*HistoryEntry, 0, len(hm.entries))
	for _, entry := range hm.entries {
		if entry.ExpiresAt.IsZero() || entry.ExpiresAt.After(now) {
			clean = append(clean, entry)
		}
	}
	hm.entries = clean
	return hm.entries
}

// Clear - Clear all history entries
func (hm *HistoryManager) Clear() {
	hm.entries = make([]*HistoryEntry, 0, hm.config.MaxItems)
}

// RemoveEntry - Remove entry by index
func (hm *HistoryManager) RemoveEntry(index int) error {
	if index < 0 || index >= len(hm.entries) {
		return errors.New("history: index out of range")
	}
	hm.entries = append(hm.entries[:index], hm.entries[index+1:]...)

	// Re-index entries
	for i := index; i < len(hm.entries); i++ {
		hm.entries[i].Index = i
	}

	return nil
}

// EntryCount - Get number of entries
func (hm *HistoryManager) EntryCount() int {
	return len(hm.entries)
}

// RestoreEntry - Restore entry to clipboard
func (hm *HistoryManager) RestoreEntry(ctx context.Context, index int) error {
	entry, err := hm.GetEntry(index)
	if err != nil {
		return err
	}
	return hm.provider.Write(entry.Data)
}
