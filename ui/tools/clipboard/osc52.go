package clipboard

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// OSC52Provider - OSC 52 clipboard provider
type OSC52Provider struct {
	config    OSC52Config
	available bool
	stdout    *os.File
}

// NewOSC52Provider - Create OSC 52 provider
func NewOSC52Provider(config OSC52Config) *OSC52Provider {
	return &OSC52Provider{
		config:    config,
		available: DetectOSC52Support(),
		stdout:    os.Stdout,
	}
}

// DetectOSC52Support - Detect if terminal supports OSC 52
func DetectOSC52Support() bool {
	// Check for environment variable override
	if os.Getenv("FORCE_OSC52") == "1" {
		return true
	}

	// On Windows, assume modern terminals support OSC 52
	if runtime.GOOS == "windows" {
		// Check for Windows Terminal, ConEmu, etc.
		wtSession := os.Getenv("WT_SESSION")
		if wtSession != "" {
			return true // Windows Terminal supports OSC 52
		}
		// Assume support for other modern terminals on Windows
		// We can't reliably detect, so return true and let the user override
		return true
	}

	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}

	// Check for known supporting terminals
	supportedTerms := []string{
		"xterm", "xterm-256color", "xterm-color",
		"screen", "screen-256color",
		"tmux", "tmux-256color",
		"rxvt", "rxvt-unicode",
		"urxvt",
		"vscode",
		"iTerm.app",
		"WezTerm",
		"Alacritty",
		"mintty",
		"msys",
		"cygwin",
	}

	for _, t := range supportedTerms {
		if strings.HasPrefix(term, t) {
			return true
		}
	}

	// Check for feature detection environment variables
	if os.Getenv("VTE_VERSION") != "" {
		return true
	}

	// Default to false for unknown terminals
	return false
}

// Available - Check if OSC 52 is available
func (p *OSC52Provider) Available() bool {
	return p.available
}

// Read - Read from clipboard (OSC 52 doesn't support reading)
func (p *OSC52Provider) Read(ctx context.Context) (*ClipboardData, error) {
	return nil, ErrTerminalNotSupported
}

// Write - Write to clipboard using OSC 52
func (p *OSC52Provider) Write(ctx context.Context, data *ClipboardData) error {
	if !p.available {
		return ErrTerminalNotSupported
	}

	if data.Empty() {
		return ErrEmptyData
	}

	// Only support text format
	if data.Format != FormatText {
		return ErrInvalidFormat
	}

	text := data.Text

	// Check size
	if len(text) > p.config.MaxSize {
		if !p.config.Truncate {
			return ErrDataTooLarge
		}
		text = text[:p.config.MaxSize]
	}

	var encoded string
	if p.config.EncodeBase64 {
		// Base64 encode
		encoded = base64.StdEncoding.EncodeToString([]byte(text))
	} else {
		// Plain text (need to escape special characters)
		encoded = escapeOSCString(text)
	}

	// Build OSC 52 sequence
	// Format: OSC 52 ; <selection> ; <data> ST
	// OSC = \x1b]
	// ST = \x1b\
	selection := p.config.Selection
	if selection == "" {
		selection = "c" // default to clipboard
	}

	sequence := fmt.Sprintf("\x1b]52;%s;%s\x1b\\", selection, encoded)

	// Write to terminal
	_, err := p.stdout.WriteString(sequence)
	if err != nil {
		return fmt.Errorf("failed to write OSC 52 sequence: %w", err)
	}

	return p.stdout.Sync()
}

// Clear - Clear clipboard (OSC 52 doesn't support clearing)
func (p *OSC52Provider) Clear(ctx context.Context) error {
	// Can write empty string to "clear" clipboard
	return p.Write(ctx, NewClipboardData(FormatText, ""))
}

// escapeOSCString - Escape OSC string
func escapeOSCString(s string) string {
	// Replace special characters with ST (String Terminator) escapes
	// \x07 (BEL), \x1b (ESC) need to be escaped

	// Simple approach: remove or replace problematic characters
	s = strings.ReplaceAll(s, "\x07", "?")
	s = strings.ReplaceAll(s, "\x1b", "?")

	return s
}

// EncodeOSC52 - Encode text for OSC 52
func EncodeOSC52(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

// DecodeOSC52 - Decode OSC 52 text
func DecodeOSC52(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode OSC 52: %w", err)
	}
	return string(data), nil
}

// IsOSC52Supported - Check if OSC 52 is supported
func IsOSC52Supported() bool {
	return DetectOSC52Support()
}

// QueryClipboard - Query clipboard using bracketed paste mode
// Note: This is not widely supported and works mainly in xterm
func QueryClipboard(selection string) (string, error) {
	if selection == "" {
		selection = "c"
	}

	// This writes the query sequence but doesn't read the response
	// Reading requires terminal input handling which is outside scope
	sequence := fmt.Sprintf("\x1b]52;%s;?\x1b\\", selection)
	_, err := os.Stdout.WriteString(sequence)
	if err != nil {
		return "", fmt.Errorf("failed to query clipboard: %w", err)
	}

	return "", fmt.Errorf("OSC 52 read not implemented (requires terminal input handling)")
}

// SetClipboard - Set clipboard using OSC 52 (convenience function)
func SetClipboard(text string, selection string, maxSize int) error {
	config := OSC52Config{
		Selection:    selection,
		MaxSize:      maxSize,
		Truncate:     true,
		EncodeBase64: true,
	}
	provider := NewOSC52Provider(config)
	if !provider.Available() {
		return ErrTerminalNotSupported
	}

	ctx := context.Background()
	data := NewClipboardData(FormatText, text)
	return provider.Write(ctx, data)
}

// CloneClipboardProvider - Clone provider with new config
func (p *OSC52Provider) Clone(config OSC52Config) *OSC52Provider {
	return &OSC52Provider{
		config:    config,
		available: p.available,
		stdout:    p.stdout,
	}
}

// GetConfig - Get current config
func (p *OSC52Provider) GetConfig() OSC52Config {
	return p.config
}

// SetConfig - Set new config
func (p *OSC52Provider) SetConfig(config OSC52Config) {
	p.config = config
}

// ParseOSC52Sequence - Parse OSC 52 sequence
// Format: \x1b]52;<selection>;<data>\x1b\
func ParseOSC52Sequence(sequence string) (selection string, encoded string, err error) {
	// Check for prefix
	if !strings.HasPrefix(sequence, "\x1b]52;") {
		return "", "", fmt.Errorf("invalid OSC 52 sequence: missing prefix")
	}

	// Remove prefix
	seq := sequence[len("\x1b]52;"):]

	// Find semicolon
	semicolon := strings.Index(seq, ";")
	if semicolon == -1 {
		return "", "", fmt.Errorf("invalid OSC 52 sequence: missing semicolon")
	}

	selection = seq[:semicolon]

	// Find ST (String Terminator)
	st := strings.Index(seq[semicolon+1:], "\x1b\\")
	if st == -1 {
		return "", "", fmt.Errorf("invalid OSC 52 sequence: missing ST")
	}

	encoded = seq[semicolon+1 : semicolon+1+st]

	return selection, encoded, nil
}

// ValidateOSC52Data - Validate OSC 52 data
func ValidateOSC52Data(encoded string, base64Encoded bool) error {
	if encoded == "" {
		return nil // empty is valid
	}

	if base64Encoded {
		// Validate base64
		_, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return fmt.Errorf("invalid base64: %w", err)
		}
	}

	return nil
}

// EstimateEncodedSize - Estimate encoded size
func EstimateEncodedSize(textSize int, base64Encoded bool) int {
	if base64Encoded {
		// Base64 increases size by ~33%
		return (textSize*4 + 2) / 3
	}
	return textSize
}

// SizeLimitExceeded - Check if size exceeds limit
func SizeLimitExceeded(text string, maxSize int) bool {
	return len(text) > maxSize
}

// TruncateToLimit - Truncate text to fit limit
func TruncateToLimit(text string, maxSize int) string {
	if len(text) <= maxSize {
		return text
	}
	return text[:maxSize]
}

// GetTerminalName - Get terminal name from env
func GetTerminalName() string {
	term := os.Getenv("TERM")
	program := os.Getenv("TERM_PROGRAM")
	wtSession := os.Getenv("WT_SESSION")

	if wtSession != "" {
		return "Windows Terminal"
	}

	if program != "" {
		return program
	}

	if term != "" {
		return term
	}

	return "unknown"
}

// GetTerminalVersion - Get terminal version
func GetTerminalVersion() string {
	program := os.Getenv("TERM_PROGRAM")
	if program != "" {
		version := os.Getenv("TERM_PROGRAM_VERSION")
		if version != "" {
			return program + " " + version
		}
		return program
	}

	// VTE version
	vteVersion := os.Getenv("VTE_VERSION")
	if vteVersion != "" {
		version, err := strconv.Atoi(vteVersion)
		if err == nil {
			major := version >> 16
			minor := (version >> 8) & 0xff
			patch := version & 0xff
			return fmt.Sprintf("VTE %d.%d.%d", major, minor, patch)
		}
	}

	return ""
}
