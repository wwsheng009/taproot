package clipboard

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// NativeProvider - Native clipboard provider
type NativeProvider struct {
	config     NativeConfig
	platform   Platform
	available  bool
	state      ClipboardState
	stateMutex sync.RWMutex
}

// Platform - Platform type
type Platform int

const (
	// PlatformUnknown - Unknown platform
	PlatformUnknown Platform = iota
	// PlatformWindows - Windows
	PlatformWindows
	// PlatformLinux - Linux
	PlatformLinux
	// PlatformDarwin - macOS
	PlatformDarwin
)

// NewNativeProvider - Create native clipboard provider
func NewNativeProvider(config NativeConfig) *NativeProvider {
	provider := &NativeProvider{
		config:   config,
		platform: detectPlatform(),
		available: detectNativeClipboard(),
		state:    StateIdle,
	}

	// Initialize platform-specific clipboard
	if runtime.GOOS != "windows" {
		// Only Windows is fully supported
		// For Linux/Darwin, we'll detect at runtime
	}

	return provider
}

// detectPlatform - Detect platform
func detectPlatform() Platform {
	switch runtime.GOOS {
	case "windows":
		return PlatformWindows
	case "linux":
		return PlatformLinux
	case "darwin":
		return PlatformDarwin
	default:
		return PlatformUnknown
	}
}

// detectNativeClipboard - Detect native clipboard support
func detectNativeClipboard() bool {
	// Windows always has clipboard support
	if runtime.GOOS == "windows" {
		return true
	}

	// Linux: Check for xclip or xsel
	// Darwin: Check for pbcopy/pbpaste
	return true // We'll try runtime detection
}

// Available - Check if native clipboard is available
func (p *NativeProvider) Available() bool {
	return p.available
}

// Read - Read from native clipboard
func (p *NativeProvider) Read(ctx context.Context) (*ClipboardData, error) {
	if !p.available {
		return nil, ErrClipboardUnavailable
	}

	p.setState(StateReading)
	defer p.setState(StateIdle)

	// Create context with timeout
	if p.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.config.Timeout)
		defer cancel()
	}

	var text string
	var err error

	// Platform-specific implementation
	switch p.platform {
	case PlatformWindows:
		text, err = p.readWindows(ctx)
	case PlatformLinux:
		text, err = p.readLinux(ctx)
	case PlatformDarwin:
		text, err = p.readDarwin(ctx)
	default:
		return nil, ErrClipboardUnavailable
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read clipboard: %w", err)
	}

	return NewClipboardData(FormatText, text), nil
}

// Write - Write to native clipboard
func (p *NativeProvider) Write(ctx context.Context, data *ClipboardData) error {
	if !p.available {
		return ErrClipboardUnavailable
	}

	if data.Empty() {
		return ErrEmptyData
	}

	// Only support text format for now
	if data.Format != FormatText {
		return ErrInvalidFormat
	}

	p.setState(StateWriting)
	defer p.setState(StateIdle)

	// Create context with timeout
	if p.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.config.Timeout)
		defer cancel()
	}

	// Retry logic
	var lastErr error
	for i := 0; i < p.config.RetryCount; i++ {
		err := p.writeWithRetry(ctx, data.Text)
		if err == nil {
			return nil
		}

		lastErr = err

		// Wait before retry
		if i < p.config.RetryCount-1 && p.config.RetryDelay > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(p.config.RetryDelay):
			}
		}
	}

	return fmt.Errorf("failed to write clipboard after %d attempts: %w", p.config.RetryCount, lastErr)
}

// writeWithRetry - Write with retry logic
func (p *NativeProvider) writeWithRetry(ctx context.Context, text string) error {
	var err error

	switch p.platform {
	case PlatformWindows:
		err = p.writeWindows(ctx, text)
	case PlatformLinux:
		err = p.writeLinux(ctx, text)
	case PlatformDarwin:
		err = p.writeDarwin(ctx, text)
	default:
		return ErrClipboardUnavailable
	}

	return err
}

// Clear - Clear native clipboard
func (p *NativeProvider) Clear(ctx context.Context) error {
	return p.Write(ctx, NewClipboardData(FormatText, ""))
}

// GetState - Get current state
func (p *NativeProvider) GetState() ClipboardState {
	p.stateMutex.RLock()
	defer p.stateMutex.RUnlock()
	return p.state
}

// setState - Set state
func (p *NativeProvider) setState(state ClipboardState) {
	p.stateMutex.Lock()
	defer p.stateMutex.Unlock()
	p.state = state
}

// GetPlatform - Get platform
func (p *NativeProvider) GetPlatform() Platform {
	return p.platform
}

// SetConfig - Set config
func (p *NativeProvider) SetConfig(config NativeConfig) {
	p.config = config
}

// GetConfig - Get config
func (p *NativeProvider) GetConfig() NativeConfig {
	return p.config
}

// Linux implementation (using xclip/xsel)
func (p *NativeProvider) readLinux(ctx context.Context) (string, error) {
	// Try xclip first
	if output, err := execCommand(ctx, "xclip", "-o", "-selection", "clipboard"); err == nil {
		return output, nil
	}

	// Try xsel
	if output, err := execCommand(ctx, "xsel", "--clipboard", "--output"); err == nil {
		return output, nil
	}

	return "", fmt.Errorf("no clipboard tool found (xclip or xsel required)")
}

func (p *NativeProvider) writeLinux(ctx context.Context, text string) error {
	// Try xclip first
	if err := execStdin(ctx, text, "xclip", "-selection", "clipboard"); err == nil {
		return nil
	}

	// Try xsel
	if err := execStdin(ctx, text, "xsel", "--clipboard", "--input"); err == nil {
		return nil
	}

	return fmt.Errorf("no clipboard tool found (xclip or xsel required)")
}

// Darwin (macOS) implementation (using pbcopy/pbpaste)
func (p *NativeProvider) readDarwin(ctx context.Context) (string, error) {
	return execCommand(ctx, "pbpaste")
}

func (p *NativeProvider) writeDarwin(ctx context.Context, text string) error {
	return execStdin(ctx, text, "pbcopy")
}

// execCommand - Execute command and return stdout
func execCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// execStdin - Execute command with stdin
func execStdin(ctx context.Context, stdin string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// GetPlatformName - Get platform name
func (p *NativeProvider) GetPlatformName() string {
	switch p.platform {
	case PlatformWindows:
		return "Windows"
	case PlatformLinux:
		return "Linux"
	case PlatformDarwin:
		return "Darwin (macOS)"
	default:
		return "Unknown"
	}
}

// IsReadSupported - Check if read is supported
func (p *NativeProvider) IsReadSupported() bool {
	switch p.platform {
	case PlatformWindows:
		return true // Windows API implemented
	case PlatformLinux:
		return true  // xclip/xsel
	case PlatformDarwin:
		return true  // pbpaste
	default:
		return false
	}
}

// IsWriteSupported - Check if write is supported
func (p *NativeProvider) IsWriteSupported() bool {
	switch p.platform {
	case PlatformWindows:
		return true // Windows API implemented
	case PlatformLinux:
		return true  // xclip/xsel
	case PlatformDarwin:
		return true  // pbcopy
	default:
		return false
	}
}

// GetSupportedFormats - Get supported formats
func (p *NativeProvider) GetSupportedFormats() []Format {
	if len(p.config.Formats) > 0 {
		return p.config.Formats
	}
	return []Format{FormatText}
}

// CheckToolAvailability - Check if clipboard tools are available
func CheckToolAvailability() map[string]bool {
	tools := map[string]bool{
		"xclip":   false,
		"xsel":    false,
		"pbcopy":  false,
		"pbpaste": false,
	}

	// This would require actual command execution
	// For now, set based on platform
	switch runtime.GOOS {
	case "linux":
		tools["xclip"] = true
		tools["xsel"] = true
	case "darwin":
		tools["pbcopy"] = true
		tools["pbpaste"] = true
	}

	return tools
}

// GetDefaultTool - Get default tool for platform
func (p *NativeProvider) GetDefaultTool() string {
	switch p.platform {
	case PlatformLinux:
		// Check if xclip is available
		if tools := CheckToolAvailability(); tools["xclip"] {
			return "xclip"
		}
		// Fall back to xsel
		if tools := CheckToolAvailability(); tools["xsel"] {
			return "xsel"
		}
	case PlatformDarwin:
		return "pbcopy/pbpaste"
	case PlatformWindows:
		return "Windows Clipboard API"
	}
	return "none"
}

// SetMaxSize - Set max size for clipboard operations
func (p *NativeProvider) SetMaxSize(int) {
	// Not applicable for native clipboard directly
	// Size limits are platform-specific
}

// GetMaxSize - Get max size for clipboard operations
func (p *NativeProvider) GetMaxSize() int {
	// Platform-specific limits
	switch p.platform {
	case PlatformWindows:
		// Windows clipboard limit is large but not infinite
		return 10 * 1024 * 1024 // 10MB
	case PlatformLinux:
		// Depends on X11/Xorg limits
		return 4 * 1024 * 1024 // 4MB
	case PlatformDarwin:
		// macOS has no fixed limit
		return 10 * 1024 * 1024 // 10MB
	default:
		return 1 * 1024 * 1024 // 1MB
	}
}

// FormatData - Format data for clipboard
func (p *NativeProvider) FormatData(data *ClipboardData) (string, error) {
	if data == nil {
		return "", ErrEmptyData
	}

	if data.IsImage() {
		return "", ErrInvalidFormat
	}

	return data.Text, nil
}

// ValidateData - Validate clipboard data
func (p *NativeProvider) ValidateData(data *ClipboardData) error {
	if data.Empty() {
		return ErrEmptyData
	}

	if data.Size() > p.GetMaxSize() {
		return ErrDataTooLarge
	}

	supported := p.GetSupportedFormats()
	for _, format := range supported {
		if data.Format == format {
			return nil
		}
	}

	return ErrInvalidFormat
}
