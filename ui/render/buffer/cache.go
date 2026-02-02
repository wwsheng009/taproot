package buffer

import (
	"strings"
	"sync"
)

// StyleCache caches ANSI escape sequences for styles
type StyleCache struct {
	cache map[uint32]string
	mu    sync.RWMutex
}

// NewStyleCache creates a new style cache
func NewStyleCache() *StyleCache {
	return &StyleCache{
		cache: make(map[uint32]string, 64),
	}
}

// Get retrieves cached style string
func (sc *StyleCache) Get(s Style) string {
	key := sc.styleToKey(s)

	sc.mu.RLock()
	cached, exists := sc.cache[key]
	sc.mu.RUnlock()

	if exists {
		return cached
	}

	// Build and cache
	styleStr := sc.buildStyle(s)
	sc.mu.Lock()
	sc.cache[key] = styleStr
	sc.mu.Unlock()

	return styleStr
}

// styleToKey converts style to a unique uint32 key
func (sc *StyleCache) styleToKey(s Style) uint32 {
	var key uint32

	if s.Bold {
		key |= 1 << 0
	}
	if s.Italic {
		key |= 1 << 1
	}
	if s.Underline {
		key |= 1 << 2
	}
	if s.Reverse {
		key |= 1 << 3
	}

	// Simple hash of foreground/background (since they're strings)
	if s.Foreground != "" {
		key = key*31 + uint32(len(s.Foreground))
		for i, c := range s.Foreground {
			key = key*31 + uint32(c)*uint32(i+1)
		}
	}
	if s.Background != "" {
		key = key*31 + uint32(len(s.Background))
		for i, c := range s.Background {
			key = key*31 + uint32(c)*uint32(i+1)
		}
	}

	return key
}

// buildStyle builds ANSI style string from style
func (sc *StyleCache) buildStyle(s Style) string {
	if !s.Bold && !s.Italic && !s.Underline && !s.Reverse &&
		s.Foreground == "" && s.Background == "" {
		return ""
	}

	var styles []string

	if s.Bold {
		styles = append(styles, "1")
	}
	if s.Italic {
		styles = append(styles, "3")
	}
	if s.Underline {
		styles = append(styles, "4")
	}
	if s.Reverse {
		styles = append(styles, "7")
	}

	if s.Foreground != "" {
		styles = append(styles, s.Foreground)
	}
	if s.Background != "" {
		styles = append(styles, s.Background)
	}

	if len(styles) > 0 {
		return "\x1b[" + strings.Join(styles, ";") + "m"
	}

	return ""
}

// Global style cache instance
var globalStyleCache = NewStyleCache()
