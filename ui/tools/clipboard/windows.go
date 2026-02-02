//go:build windows
// +build windows

package clipboard

import (
	"context"
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procOpenClipboard   = user32.NewProc("OpenClipboard")
	procCloseClipboard  = user32.NewProc("CloseClipboard")
	procEmptyClipboard  = user32.NewProc("EmptyClipboard")
	procGetClipboardData = user32.NewProc("GetClipboardData")
	procSetClipboardData = user32.NewProc("SetClipboardData")
	procIsClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable")
	procGlobalAlloc      = kernel32.NewProc("GlobalAlloc")
	procGlobalFree       = kernel32.NewProc("GlobalFree")
	procGlobalLock       = kernel32.NewProc("GlobalLock")
	procGlobalUnlock     = kernel32.NewProc("GlobalUnlock")
	procGlobalSize       = kernel32.NewProc("GlobalSize")
	procGetClipboardOwner = user32.NewProc("GetClipboardOwner")
)

const (
	CF_TEXT     = 1
	CF_UNICODETEXT = 13
	GMEM_MOVEABLE = 0x0002
	GMEM_ZEROINIT = 0x0040
)

// WindowsClipboard - Windows clipboard wrapper
type WindowsClipboard struct{}

// NewWindowsClipboard - Create Windows clipboard
func NewWindowsClipboard() *WindowsClipboard {
	return &WindowsClipboard{}
}

// Open - Open clipboard
func (w *WindowsClipboard) Open() error {
	ret, _, err := procOpenClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("failed to open clipboard: %v", err)
	}
	return nil
}

// Close - Close clipboard
func (w *WindowsClipboard) Close() error {
	ret, _, err := procCloseClipboard.Call()
	if ret == 0 {
		return fmt.Errorf("failed to close clipboard: %v", err)
	}
	return nil
}

// Empty - Empty clipboard
func (w *WindowsClipboard) Empty() error {
	ret, _, err := procEmptyClipboard.Call()
	if ret == 0 {
		return fmt.Errorf("failed to empty clipboard: %v", err)
	}
	return nil
}

// SetText - Set text to clipboard
func (w *WindowsClipboard) SetText(text string) error {
	// Convert to UTF-16
	utf16Text := utf16.Encode([]rune(text))
	size := len(utf16Text) * 2 // 2 bytes per char

	// Allocate global memory
	hMem, _, err := procGlobalAlloc.Call(GMEM_MOVEABLE|GMEM_ZEROINIT, uintptr(size+2))
	if hMem == 0 {
		return fmt.Errorf("failed to allocate memory: %v", err)
	}

	// Lock memory
	ptr, _, err := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return fmt.Errorf("failed to lock memory: %v", err)
	}

	// Copy data
	for i, c := range utf16Text {
		*(*uint16)(unsafe.Pointer(ptr + uintptr(i*2))) = c
	}

	// Unlock memory
	procGlobalUnlock.Call(hMem)

	// Set clipboard data
	ret, _, err := procSetClipboardData.Call(CF_UNICODETEXT, hMem)
	if ret == 0 {
		procGlobalFree.Call(hMem)
		return fmt.Errorf("failed to set clipboard data: %v", err)
	}

	return nil
}

// GetText - Get text from clipboard
func (w *WindowsClipboard) GetText() (string, error) {
	// Check if Unicode text is available
	ret, _, _ := procIsClipboardFormatAvailable.Call(CF_UNICODETEXT)
	if ret == 0 {
		// Try ANSI text
		ret, _, _ = procIsClipboardFormatAvailable.Call(CF_TEXT)
		if ret == 0 {
			return "", fmt.Errorf("no text in clipboard")
		}
	}

	// Get clipboard data
	hMem, _, err := procGetClipboardData.Call(CF_UNICODETEXT)
	if hMem == 0 {
		return "", fmt.Errorf("failed to get clipboard data: %v", err)
	}

	// Lock memory
	ptr, _, err := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return "", fmt.Errorf("failed to lock memory: %v", err)
	}

	// Get size
	size, _, err := procGlobalSize.Call(hMem)
	if size == 0 {
		return "", fmt.Errorf("failed to get memory size: %v", err)
	}

	// Read UTF-16 data
	utf16Text := make([]uint16, size/2)
	for i := uintptr(0); i < uintptr(size)/2; i++ {
		utf16Text[i] = *(*uint16)(unsafe.Pointer(ptr + i*2))
	}

	// Unlock memory
	procGlobalUnlock.Call(hMem)

	// Convert UTF-16 to string
	text := string(utf16.Decode(utf16Text))
	
	// Remove trailing null
	if len(text) > 0 && text[len(text)-1] == 0 {
		text = text[:len(text)-1]
	}

	return text, nil
}

// IsOpen - Check if clipboard is open
func (w *WindowsClipboard) IsOpen() bool {
	// Try to get clipboard owner
	owner, _, _ := procGetClipboardOwner.Call()
	return owner != 0
}

// readWindows - Read from Windows clipboard
func (p *NativeProvider) readWindows(ctx context.Context) (string, error) {
	win := NewWindowsClipboard()

	// Open clipboard
	err := win.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open clipboard: %w", err)
	}
	defer win.Close()

	// Get text
	text, err := win.GetText()
	if err != nil {
		return "", err
	}

	return text, nil
}

// writeWindows - Write to Windows clipboard
func (p *NativeProvider) writeWindows(ctx context.Context, text string) error {
	win := NewWindowsClipboard()

	// Open clipboard
	err := win.Open()
	if err != nil {
		return fmt.Errorf("cannot open clipboard: %w", err)
	}
	defer win.Close()

	// Empty clipboard
	err = win.Empty()
	if err != nil {
		return fmt.Errorf("cannot empty clipboard: %w", err)
	}

	// Set text
	err = win.SetText(text)
	if err != nil {
		return fmt.Errorf("cannot set clipboard text: %w", err)
	}

	return nil
}
