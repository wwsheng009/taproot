package clipboard

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// Manager - Clipboard manager
type Manager struct {
	osc52      *OSC52Provider
	native     *NativeProvider
	primary    Provider // Active provider
	history    *HistoryManager
	config     ManagerConfig
	state      ClipboardState
	stateMutex sync.RWMutex
}

// ManagerConfig - Manager configuration
type ManagerConfig struct {
	// PreferredType - Preferred clipboard type
	PreferredType ClipboardType
	// OSC52Config - OSC 52 configuration
	OSC52Config OSC52Config
	// NativeConfig - Native clipboard configuration
	NativeConfig NativeConfig
	// HistoryConfig - History configuration
	HistoryConfig HistoryConfig
	// EnableHistory - Enable history tracking
	EnableHistory bool
	// AutoSwitch - Auto-switch between providers
	AutoSwitch bool
}

// DefaultManagerConfig - Default manager config
func DefaultManagerConfig() ManagerConfig {
	return ManagerConfig{
		PreferredType: ClipboardPlatform,
		OSC52Config:   DefaultOSC52Config(),
		NativeConfig:  DefaultNativeConfig(),
		HistoryConfig: DefaultHistoryConfig(),
		EnableHistory: true,
		AutoSwitch:    true,
	}
}

// NewManager - Create clipboard manager
func NewManager(config ManagerConfig) *Manager {
	m := &Manager{
		config: config,
		state:  StateIdle,
	}

	// Initialize providers
	m.osc52 = NewOSC52Provider(config.OSC52Config)
	m.native = NewNativeProvider(config.NativeConfig)

	// Select primary provider
	m.selectProvider()

	// Initialize history
	if config.EnableHistory {
		m.history = NewHistoryManager(config.HistoryConfig, m)
	}

	return m
}

// NewDefaultManager - Create clipboard manager with default config
func NewDefaultManager() *Manager {
	return NewManager(DefaultManagerConfig())
}

// selectProvider - Select primary provider based on config
func (m *Manager) selectProvider() {
	switch m.config.PreferredType {
	case ClipboardOSC52:
		if m.osc52.Available() {
			m.primary = m.osc52
			return
		}
	case ClipboardNative:
		if m.native.Available() {
			m.primary = m.native
			return
		}
	case ClipboardPlatform:
		// Auto-detect best provider
		if m.native.Available() && m.native.IsWriteSupported() {
			m.primary = m.native
			return
		}
		if m.osc52.Available() {
			m.primary = m.osc52
			return
		}
	}

	// Fallback to whatever is available
	if m.native.Available() {
		m.primary = m.native
	} else if m.osc52.Available() {
		m.primary = m.osc52
	}
}

// Available - Check if any clipboard is available
func (m *Manager) Available() bool {
	return m.osc52.Available() || m.native.Available()
}

// Read - Read from clipboard
func (m *Manager) Read(format Format) (*ClipboardData, error) {
	if !m.Available() {
		return nil, ErrClipboardUnavailable
	}

	m.setState(StateReading)
	defer m.setState(StateIdle)

	ctx := context.Background()
	data, err := m.primary.Read(ctx)
	if err != nil {
		// Try fallback provider
		if m.config.AutoSwitch {
			if m.primary == m.native && m.osc52.Available() {
				data, err = m.osc52.Read(ctx)
			} else if m.primary == m.osc52 && m.native.Available() {
				data, err = m.native.Read(ctx)
			}
		}
	}

	return data, err
}

// Write - Write to clipboard
func (m *Manager) Write(data *ClipboardData) error {
	if !m.Available() {
		return ErrClipboardUnavailable
	}

	if data.Empty() {
		return ErrEmptyData
	}

	m.setState(StateWriting)
	defer m.setState(StateIdle)

	ctx := context.Background()
	err := m.primary.Write(ctx, data)
	if err != nil {
		// Try fallback provider
		if m.config.AutoSwitch {
			if m.primary == m.native && m.osc52.Available() {
				err = m.osc52.Write(ctx, data)
			} else if m.primary == m.osc52 && m.native.Available() {
				err = m.native.Write(ctx, data)
			}
		}
	}

	// Add to history on success
	if err == nil && m.config.EnableHistory && m.history != nil {
		m.history.AddEntry(data)
	}

	return err
}

// Clear - Clear clipboard
func (m *Manager) Clear() error {
	ctx := context.Background()
	return m.primary.Clear(ctx)
}

// Formats - Get supported formats
func (m *Manager) Formats() []Format {
	return []Format{FormatText}
}

// Type - Get clipboard type
func (m *Manager) Type() ClipboardType {
	if m.primary == m.osc52 {
		return ClipboardOSC52
	}
	if m.primary == m.native {
		return ClipboardNative
	}
	return ClipboardPlatform
}

// History - Get clipboard history
func (m *Manager) History() []*ClipboardData {
	if !m.config.EnableHistory || m.history == nil {
		return nil
	}

	entries := m.history.GetAllEntries()
	result := make([]*ClipboardData, len(entries))
	for i, entry := range entries {
		result[i] = entry.Data
	}
	return result
}

// AddToHistory - Add to history manually
func (m *Manager) AddToHistory(data *ClipboardData) {
	if m.config.EnableHistory && m.history != nil {
		m.history.AddEntry(data)
	}
}

// ClearHistory - Clear history
func (m *Manager) ClearHistory() {
	if m.config.EnableHistory && m.history != nil {
		m.history.Clear()
	}
}

// RestoreFromHistory - Restore from history
func (m *Manager) RestoreFromHistory(index int) error {
	if !m.config.EnableHistory || m.history == nil {
		return fmt.Errorf("history not enabled")
	}

	ctx := context.Background()
	return m.history.RestoreEntry(ctx, index)
}

// PersistHistory - Persist history to disk
func (m *Manager) PersistHistory() error {
	if !m.config.EnableHistory || m.history == nil {
		return nil
	}

	if m.config.HistoryConfig.PersistPath == "" {
		return nil // No persistence configured
	}

	// TODO: Implement persistence
	return nil
}

// LoadHistory - Load history from disk
func (m *Manager) LoadHistory() error {
	if !m.config.EnableHistory || m.history == nil {
		return nil
	}

	if m.config.HistoryConfig.PersistPath == "" {
		return nil // No persistence configured
	}

	// TODO: Implement loading
	return nil
}

// SetProvider - Set provider type
func (m *Manager) SetProvider(clipboardType ClipboardType) error {
	switch clipboardType {
	case ClipboardOSC52:
		if !m.osc52.Available() {
			return ErrClipboardUnavailable
		}
		m.primary = m.osc52
	case ClipboardNative:
		if !m.native.Available() {
			return ErrClipboardUnavailable
		}
		m.primary = m.native
	case ClipboardPlatform:
		m.selectProvider()
	default:
		return ErrInvalidFormat
	}
	return nil
}

// GetProvider - Get current provider
func (m *Manager) GetProvider() Provider {
	return m.primary
}

// GetState - Get state
func (m *Manager) GetState() ClipboardState {
	m.stateMutex.RLock()
	defer m.stateMutex.RUnlock()
	return m.state
}

// setState - Set state
func (m *Manager) setState(state ClipboardState) {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()
	m.state = state
}

// GetHistoryCount - Get history count
func (m *Manager) GetHistoryCount() int {
	if m.config.EnableHistory && m.history != nil {
		return m.history.EntryCount()
	}
	return 0
}

// GetHistoryInfo - Get history info
func (m *Manager) GetHistoryInfo() ([]*ClipboardData, int, int) {
	if !m.config.EnableHistory || m.history == nil {
		return nil, 0, 0
	}

	entries := m.history.GetAllEntries()
	data := make([]*ClipboardData, len(entries))
	for i, entry := range entries {
		data[i] = entry.Data
	}

	maxItems := m.config.HistoryConfig.MaxItems
	currentItems := len(entries)

	return data, currentItems, maxItems
}

// Copy - Copy text to clipboard
func (m *Manager) Copy(text string) error {
	return m.Write(NewClipboardData(FormatText, text))
}

// Paste - Paste from clipboard
func (m *Manager) Paste() (string, error) {
	data, err := m.Read(FormatText)
	if err != nil {
		return "", err
	}
	return data.Text, nil
}

// TryCopy - Try to copy with fallback
func (m *Manager) TryCopy(text string) error {
	// Try OSC 52 first
	if m.osc52.Available() {
		ctx := context.Background()
		data := NewClipboardData(FormatText, text)
		err := m.osc52.Write(ctx, data)
		if err == nil && m.config.EnableHistory && m.history != nil {
			m.history.AddEntry(data)
		}
		if err == nil {
			return nil
		}
	}

	// Try native clipboard
	if m.native.Available() {
		ctx := context.Background()
		data := NewClipboardData(FormatText, text)
		err := m.native.Write(ctx, data)
		if err == nil && m.config.EnableHistory && m.history != nil {
			m.history.AddEntry(data)
		}
		return err
	}

	return ErrClipboardUnavailable
}

// TryRead - Try to read with fallback
func (m *Manager) TryRead() (string, error) {
	ctx := context.Background()

	// Try native clipboard first
	if m.native.Available() && m.native.IsReadSupported() {
		data, err := m.native.Read(ctx)
		if err == nil {
			return data.Text, nil
		}
	}

	// Try OSC 52 (doesn't support read)
	// OSC 52 doesn't support reading, so we can't try it

	return "", ErrClipboardUnavailable
}

// GetPlatformInfo - Get platform information
func (m *Manager) GetPlatformInfo() map[string]interface{} {
	info := make(map[string]interface{})

	info["OS"] = runtime.GOOS
	info["Arch"] = runtime.GOARCH
	info["OSC52Available"] = m.osc52.Available()
	info["NativeAvailable"] = m.native.Available()
	info["NativePlatform"] = m.native.GetPlatformName()
	info["NativeWriteSupported"] = m.native.IsWriteSupported()
	info["NativeReadSupported"] = m.native.IsReadSupported()
	info["PrimaryProvider"] = m.Type().String()

	return info
}

// GetHistoryEntries - Get detailed history entries
func (m *Manager) GetHistoryEntries() []*HistoryEntry {
	if !m.config.EnableHistory || m.history == nil {
		return nil
	}
	return m.history.GetAllEntries()
}

// RemoveFromHistory - Remove from history by index
func (m *Manager) RemoveFromHistory(index int) error {
	if !m.config.EnableHistory || m.history == nil {
		return fmt.Errorf("history not enabled")
	}
	return m.history.RemoveEntry(index)
}

// GetHistoryEntry - Get specific history entry
func (m *Manager) GetHistoryEntry(index int) (*HistoryEntry, error) {
	if !m.config.EnableHistory || m.history == nil {
		return nil, fmt.Errorf("history not enabled")
	}
	return m.history.GetEntry(index)
}

// UpdateConfig - Update manager config
func (m *Manager) UpdateConfig(config ManagerConfig) {
	m.config = config

	// Update provider configs
	m.osc52.SetConfig(config.OSC52Config)
	m.native.SetConfig(config.NativeConfig)

	// Reinitialize history if needed
	if config.EnableHistory {
		if m.history == nil {
			m.history = NewHistoryManager(config.HistoryConfig, m)
		}
	} else {
		m.history = nil
	}

	// Re-select provider
	m.selectProvider()
}

// GetConfig - Get current config
func (m *Manager) GetConfig() ManagerConfig {
	return m.config
}

// EnableHistory - Enable/disable history
func (m *Manager) EnableHistory(enable bool) {
	m.config.EnableHistory = enable
	if enable && m.history == nil {
		m.history = NewHistoryManager(m.config.HistoryConfig, m)
	} else if !enable {
		m.history = nil
	}
}

// IsHistoryEnabled - Check if history is enabled
func (m *Manager) IsHistoryEnabled() bool {
	return m.config.EnableHistory && m.history != nil
}
