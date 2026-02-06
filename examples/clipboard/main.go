package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/wwsheng009/taproot/ui/tools/clipboard"
)

func main() {
	fmt.Println("=== Taproot Clipboard Tool Demo ===")
	fmt.Println()

	// Get platform information
	fmt.Printf("Platform: %s (%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()

	fmt.Println("Creating clipboard manager...")
	manager := clipboard.NewDefaultManager()
	if !manager.Available() {
		fmt.Fprintln(os.Stderr, "❌ Clipboard not available on this system")
		fmt.Fprintln(os.Stderr, "   - OSC 52 requires terminal support")
		fmt.Fprintln(os.Stderr, "   - Native clipboard may require external tools (xclip, pbcopy, etc.)")
		os.Exit(1)
	}

	fmt.Println("✓ Clipboard initialized")
	fmt.Printf("  Provider: %s\n", manager.Type().String())
	fmt.Println()

	// Show platform info
	showPlatformInfo(manager)

	// Main menu loop
	for {
		fmt.Println()
		fmt.Println("--- Main Menu ---")
		fmt.Println("1. Copy text to clipboard")
		fmt.Println("2. Paste from clipboard")
		fmt.Println("3. View clipboard history")
		fmt.Println("4. Clear clipboard")
		fmt.Println("5. Restore from history")
		fmt.Println("6. Test OSC 52 encoding")
		fmt.Println("7. Switch provider")
		fmt.Println("8. Show platform info")
		fmt.Println("0. Exit")
		fmt.Print("\nSelect option: ")

		var choice int
		fmt.Scanln(&choice)
		fmt.Println()

		switch choice {
		case 0:
			fmt.Println("Exiting...")
			return
		case 1:
			copyToClipboard(manager)
		case 2:
			pasteFromClipboard(manager)
		case 3:
			viewHistory(manager)
		case 4:
			clearClipboard(manager)
		case 5:
			restoreFromHistory(manager)
		case 6:
			testOSC52()
		case 7:
			switchProvider(manager)
		case 8:
			showPlatformInfo(manager)
		default:
			fmt.Println("Invalid option")
		}

		fmt.Println("\nPress Enter to continue...")
		fmt.Scanln()
	}
}

func copyToClipboard(manager *clipboard.Manager) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text to copy: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	if text == "" {
		fmt.Println("❌ Empty text")
		return
	}

	fmt.Println("\nCopying text...")
	err := manager.Copy(text)
	if err != nil {
		fmt.Printf("❌ Failed to copy: %v\n", err)
		return
	}

	fmt.Printf("✓ Copied %d characters to clipboard\n", len(text))
	fmt.Printf("  Provider: %s\n", manager.Type().String())
}

func pasteFromClipboard(manager *clipboard.Manager) {
	fmt.Println("Pasting from clipboard...")
	text, err := manager.Paste()
	if err != nil {
		fmt.Printf("❌ Failed to paste: %v\n", err)
		fmt.Println("   Note: OSC 52 doesn't support reading from clipboard")
		fmt.Println("   Native clipboard read may require external tools")
		return
	}

	if text == "" {
		fmt.Println("Clipboard is empty")
		return
	}

	fmt.Printf("✓ Pasted %d characters:\n", len(text))
	fmt.Printf("  %s\n", text)
}

func viewHistory(manager *clipboard.Manager) {
	if !manager.IsHistoryEnabled() {
		fmt.Println("❌ History is not enabled")
		return
	}

	count := manager.GetHistoryCount()
	if count == 0 {
		fmt.Println("History is empty")
		return
	}

	fmt.Printf("History entries (%d items):\n", count)
	entries := manager.GetHistoryEntries()
	for i, entry := range entries {
		preview := entry.Data.Text
		if len(preview) > 40 {
			preview = preview[:40] + "..."
		}
		fmt.Printf("  [%d] %s (%s ago)\n", i,
			preview,
			time.Since(entry.Added).Round(time.Second))
	}
}

func clearClipboard(manager *clipboard.Manager) {
	fmt.Println("Clearing clipboard...")
	err := manager.Clear()
	if err != nil {
		fmt.Printf("❌ Failed to clear: %v\n", err)
		return
	}

	fmt.Println("✓ Clipboard cleared")
}

func restoreFromHistory(manager *clipboard.Manager) {
	if !manager.IsHistoryEnabled() {
		fmt.Println("❌ History is not enabled")
		return
	}

	count := manager.GetHistoryCount()
	if count == 0 {
		fmt.Println("History is empty")
		return
	}

	viewHistory(manager)

	fmt.Print("\nEnter index to restore: ")
	var index int
	fmt.Scanln(&index)

	if index < 0 || index >= count {
		fmt.Printf("❌ Invalid index (0-%d)\n", count-1)
		return
	}

	fmt.Printf("Restoring from history index %d...\n", index)
	err := manager.RestoreFromHistory(index)
	if err != nil {
		fmt.Printf("❌ Failed to restore: %v\n", err)
		return
	}

	fmt.Println("✓ Restored from history")
}

func testOSC52() {
	fmt.Println("Testing OSC 52 encoding/decoding...")

	testStrings := []string{
		"Hello, World!",
		"Special chars: @#$%^&*()",
		"Unicode: 你好世界 🌍",
		"Lorem ipsum dolor sit amet",
	}

	for _, text := range testStrings {
		fmt.Printf("Original: %s\n", text)

		// Encode
		encoded := clipboard.EncodeOSC52(text)
		fmt.Printf("Encoded:  %s\n", encoded[:min(50, len(encoded))] + "...")

		// Decode
		decoded, err := clipboard.DecodeOSC52(encoded)
		if err != nil {
			fmt.Printf("❌ Decode failed: %v\n", err)
			continue
		}

		if decoded != text {
			fmt.Printf("❌ Mismatch: %s\n", decoded)
		} else {
			fmt.Printf("✓ Decoded successfully\n")
		}
		fmt.Printf("Size: %d -> %d (%.1f%% increase)\n",
			len(text), len(encoded),
			float64(len(encoded)-len(text))/float64(len(text))*100)
		fmt.Println()
	}

	// Test size limit
	fmt.Println("\nTesting size limits...")
	longText := strings.Repeat("a", 200000)
	fmt.Printf("Long text size: %d characters\n", len(longText))

	config := clipboard.DefaultOSC52Config()
	maxSize := config.MaxSize
	fmt.Printf("Max size: %d bytes\n", maxSize)

	if clipboard.SizeLimitExceeded(longText, maxSize) {
		fmt.Printf("✓ Size limit exceeded\n")
		truncated := clipboard.TruncateToLimit(longText, maxSize)
		fmt.Printf("  Truncated to: %d bytes\n", len(truncated))

		// Encode truncated
		encoded := clipboard.EncodeOSC52(truncated)
		estimatedSize := clipboard.EstimateEncodedSize(len(truncated), true)
		fmt.Printf("  Estimated encoded size: %d\n", estimatedSize)
		fmt.Printf("  Actual encoded size: %d\n", len(encoded))
	}
}

func switchProvider(manager *clipboard.Manager) {
	fmt.Println("Current provider:", manager.Type().String())
	fmt.Println()

	fmt.Println("Available providers:")
	fmt.Println("1. OSC 52 (Terminal)")
	fmt.Println("2. Native (OS)")
	fmt.Println("3. Platform (Auto-detect)")
	fmt.Println()

	fmt.Print("Select provider: ")
	var choice int
	fmt.Scanln(&choice)

	var providerType clipboard.ClipboardType
	switch choice {
	case 1:
		providerType = clipboard.ClipboardOSC52
	case 2:
		providerType = clipboard.ClipboardNative
	case 3:
		providerType = clipboard.ClipboardPlatform
	default:
		fmt.Println("Invalid option")
		return
	}

	err := manager.SetProvider(providerType)
	if err != nil {
		fmt.Printf("❌ Failed to set provider: %v\n", err)
		return
	}

	fmt.Printf("✓ Switched to %s provider\n", manager.Type().String())
}

func showPlatformInfo(manager *clipboard.Manager) {
	info := manager.GetPlatformInfo()

	fmt.Println("=== Platform Information ===")
	for key, value := range info {
		switch v := value.(type) {
		case string:
			fmt.Printf("%s: %s\n", key, v)
		case bool:
			if v {
				fmt.Printf("%s: Yes\n", key)
			} else {
				fmt.Printf("%s: No\n", key)
			}
		}
	}

	// Terminal info
	terminalName := clipboard.GetTerminalName()
	fmt.Printf("Terminal: %s\n", terminalName)

	terminalVersion := clipboard.GetTerminalVersion()
	if terminalVersion != "" {
		fmt.Printf("Terminal Version: %s\n", terminalVersion)
	}

	// History info
	if manager.IsHistoryEnabled() {
		entries, count, max := manager.GetHistoryInfo()
		fmt.Printf("\nHistory: %d/%d entries\n", count, max)
		if count > 0 {
			oldest := entries[len(entries)-1]
			newest := entries[0]
			fmt.Printf("  Newest: %s ago\n", time.Since(newest.Timestamp).Round(time.Second))
			fmt.Printf("  Oldest: %s ago\n", time.Since(oldest.Timestamp).Round(time.Second))
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
