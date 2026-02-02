// Package watcher provides file system monitoring utilities.
package watcher

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mattn/go-isatty"
	"github.com/ryanuber/go-glob"
)

// Watcher monitors file system events.
type Watcher struct {
	fsWatcher    *fsnotify.Watcher
	config       Config
	state        WatcherState
	stateLock    sync.RWMutex
	watchedPaths map[string]bool
	watchedLock  sync.Mutex
	stats        Stats
	statsLock    sync.RWMutex

	// Event handling
	eventChan    chan Event
	batchChan    chan []Event
	stopChan     chan struct{}
	pauseChan    chan struct{}
	paused       atomic.Bool

	// Debouncing
	debounceMap map[string]*DebouncedEvent
	debounceMu  sync.Mutex
	debounceT   *time.Timer

	// Batching
	batchMu     sync.Mutex
	batchBuffer []Event
	batchT      *time.Timer
}

// NewWatcher creates a new file watcher with default configuration.
func NewWatcher() (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	return &Watcher{
		fsWatcher:    fsWatcher,
		config:       DefaultConfig(),
		state:        WatcherIdle,
		watchedPaths: make(map[string]bool),
		eventChan:    make(chan Event, 1000),
		batchChan:    make(chan []Event, 100),
		stopChan:     make(chan struct{}),
		pauseChan:    make(chan struct{}),
		debounceMap:  make(map[string]*DebouncedEvent),
		batchBuffer:  make([]Event, 0, 100),
	}, nil
}

// NewWatcherWithConfig creates a new file watcher with custom configuration.
func NewWatcherWithConfig(config Config) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	return &Watcher{
		fsWatcher:    fsWatcher,
		config:       config,
		state:        WatcherIdle,
		watchedPaths: make(map[string]bool),
		eventChan:    make(chan Event, config.BufferSize),
		batchChan:    make(chan []Event, 100),
		stopChan:     make(chan struct{}),
		pauseChan:    make(chan struct{}),
		debounceMap:  make(map[string]*DebouncedEvent),
		batchBuffer:  make([]Event, 0, config.Batch.MaxSize),
	}, nil
}

// Add adds a path to be watched.
func (w *Watcher) Add(path string) error {
	w.stateLock.RLock()
	if w.state == WatcherRunning {
		w.stateLock.RUnlock()
		return ErrWatcherRunning
	}
	w.stateLock.RUnlock()

	// Validate path
	if path == "" {
		return ErrInvalidPath
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrPathNotFound
		}
		return fmt.Errorf("%w: %v", ErrInvalidPath, err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if already watching
	w.watchedLock.Lock()
	if w.watchedPaths[absPath] {
		w.watchedLock.Unlock()
		return ErrAlreadyWatching
	}
	w.watchedLock.Unlock()

	// Add to fsnotify watcher
	if err := w.fsWatcher.Add(absPath); err != nil {
		return fmt.Errorf("failed to add path to watcher: %w", err)
	}

	// Mark as watched
	w.watchedLock.Lock()
	w.watchedPaths[absPath] = true
	w.watchedLock.Unlock()

	// Update stats
	w.statsLock.Lock()
	if info.IsDir() {
		w.stats.DirsWatched++
	} else {
		w.stats.FilesWatched++
	}
	w.statsLock.Unlock()

	return nil
}

// AddRecursive adds a path and all subdirectories to be watched.
func (w *Watcher) AddRecursive(path string) error {
	w.stateLock.RLock()
	if w.state == WatcherRunning {
		w.stateLock.RUnlock()
		return ErrWatcherRunning
	}
	w.stateLock.RUnlock()

	// Add the root path
	if err := w.Add(path); err != nil {
		return err
	}

	// Walk through directory and add subdirectories
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if info.IsDir() && walkPath != path {
			absPath, err := filepath.Abs(walkPath)
			if err != nil {
				return nil
			}

			w.watchedLock.Lock()
			if !w.watchedPaths[absPath] {
				if err := w.fsWatcher.Add(absPath); err == nil {
					w.watchedPaths[absPath] = true
					w.statsLock.Lock()
					w.stats.DirsWatched++
					w.statsLock.Unlock()
				}
			}
			w.watchedLock.Unlock()
		}

		return nil
	})
}

// Remove removes a path from being watched.
func (w *Watcher) Remove(path string) error {
	w.stateLock.RLock()
	if w.state == WatcherRunning {
		w.stateLock.RUnlock()
		return ErrWatcherRunning
	}
	w.stateLock.RUnlock()

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if watching
	w.watchedLock.Lock()
	if !w.watchedPaths[absPath] {
		w.watchedLock.Unlock()
		return ErrNotWatching
	}
	w.watchedLock.Unlock()

	// Remove from fsnotify watcher
	if err := w.fsWatcher.Remove(absPath); err != nil {
		return fmt.Errorf("failed to remove path from watcher: %w", err)
	}

	// Mark as not watching
	w.watchedLock.Lock()
	delete(w.watchedPaths, absPath)
	w.watchedLock.Unlock()

	return nil
}

// SetConfig updates the watcher configuration.
func (w *Watcher) SetConfig(config Config) error {
	w.stateLock.RLock()
	if w.state == WatcherRunning {
		w.stateLock.RUnlock()
		return ErrWatcherRunning
	}
	w.stateLock.RUnlock()

	w.config = config
	return nil
}

// Start begins watching for events.
func (w *Watcher) Start() error {
	w.stateLock.Lock()
	defer w.stateLock.Unlock()

	if w.state != WatcherIdle && w.state != WatcherStopped {
		return ErrWatcherRunning
	}

	if w.config.Handler == nil && w.config.Batch.Enabled {
		return ErrNoHandler
	}

	w.state = WatcherRunning
	w.stats.StartTime = time.Now()

	// Start event processing goroutines
	go w.processEvents()
	if w.config.Batch.Enabled {
		go w.processBatches()
	}

	return nil
}

// Stop stops watching for events.
func (w *Watcher) Stop() error {
	w.stateLock.Lock()
	defer w.stateLock.Unlock()

	if w.state != WatcherRunning && w.state != WatcherPaused {
		return ErrWatcherStopped
	}

	w.state = WatcherStopped
	close(w.stopChan)

	if err := w.fsWatcher.Close(); err != nil {
		return fmt.Errorf("failed to close watcher: %w", err)
	}

	return nil
}

// Pause temporarily pauses event processing.
func (w *Watcher) Pause() {
	w.paused.Store(true)
}

// Resume resumes event processing after a pause.
func (w *Watcher) Resume() {
	w.paused.Store(false)
}

// Wait blocks until the watcher is stopped.
func (w *Watcher) Wait() {
	<-w.stopChan
}

// State returns the current watcher state.
func (w *Watcher) State() WatcherState {
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()
	return w.state
}

// Stats returns current statistics.
func (w *Watcher) Stats() Stats {
	w.statsLock.RLock()
	defer w.statsLock.RUnlock()
	return w.stats
}

// processEvents processes events from the fsnotify watcher.
func (w *Watcher) processEvents() {
	defer close(w.eventChan)

	for {
		select {
		case <-w.stopChan:
			return

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			if w.paused.Load() {
				continue // Skip events while paused
			}

			// Map fsnotify event type to our EventType
			var eventType EventType
			switch {
			case event.Op&fsnotify.Create == fsnotify.Create:
				eventType = EventCreate
			case event.Op&fsnotify.Write == fsnotify.Write:
				eventType = EventWrite
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				eventType = EventRemove
			case event.Op&fsnotify.Rename == fsnotify.Rename:
				eventType = EventRename
			case event.Op&fsnotify.Chmod == fsnotify.Chmod:
				eventType = EventChmod
			default:
				continue // Unknown event type
			}

			// Update stats
			w.statsLock.Lock()
			w.stats.TotalEvents++
			w.stats.LastEventTime = time.Now()
			w.statsLock.Unlock()

			// Apply filters
			if !w.shouldProcessEvent(eventType, event.Name) {
				continue
			}

			// Create our event
			wEvent := Event{
				Type:      eventType,
				Path:      event.Name,
				Timestamp: time.Now(),
			}

			// Handle event
			if w.config.Debounce.Enabled {
				w.handleDebouncedEvent(wEvent)
			} else if w.config.Batch.Enabled {
				w.handleBatchedEvent(wEvent)
			} else {
				// Direct handling
				if w.config.Handler != nil {
					w.config.Handler([]Event{wEvent})
				}
			}

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}

			// Update state to error
			w.stateLock.Lock()
			w.state = WatcherError
			w.stateLock.Unlock()

			// Call error handler if configured
			if w.config.ErrorHandler != nil {
				w.config.ErrorHandler(err)
			} else {
				slog.Error("Watcher error", "error", err)
			}
		}
	}
}

// shouldProcessEvent checks if an event should be processed based on filters.
func (w *Watcher) shouldProcessEvent(eventType EventType, path string) bool {
	filter := w.config.Filter

	// Check event type filter
	if len(filter.EventTypes) > 0 {
		typeIncluded := false
		for _, t := range filter.EventTypes {
			if t == eventType {
				typeIncluded = true
				break
			}
		}
		if !typeIncluded {
			return false
		}
	}

	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check directory filter
	if !filter.IncludeDirs && info.IsDir() {
		return false
	}

	// Check hidden files filter
	if filter.IgnoreHidden {
		base := filepath.Base(path)
		if strings.HasPrefix(base, ".") {
			return false
		}
	}

	// Check file size filter
	if !info.IsDir() {
		if filter.MinSize > 0 && info.Size() < filter.MinSize {
			return false
		}
		if filter.MaxSize > 0 && info.Size() > filter.MaxSize {
			return false
		}
	}

	// Check extension filter
	if len(filter.Extensions) > 0 {
		ext := strings.ToLower(filepath.Ext(path))
		extIncluded := false
		for _, e := range filter.Extensions {
			if strings.ToLower(e) == ext {
				extIncluded = true
				break
			}
		}
		if !extIncluded {
			return false
		}
	}

	// Check include patterns
	if len(filter.IncludePatterns) > 0 {
		included := false
		base := filepath.Base(path)
		for _, pattern := range filter.IncludePatterns {
			if glob.Glob(pattern, base) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	// Check exclude patterns
	if len(filter.ExcludePatterns) > 0 {
		base := filepath.Base(path)
		for _, pattern := range filter.ExcludePatterns {
			if glob.Glob(pattern, base) {
				return false
			}
		}
	}

	return true
}

// handleDebouncedEvent handles an event with debouncing.
func (w *Watcher) handleDebouncedEvent(event Event) {
	w.debounceMu.Lock()
	defer w.debounceMu.Unlock()

	// Reset existing timer
	if w.debounceT != nil {
		w.debounceT.Stop()
	}

	// Get or create debounced entry for this path
	entry, exists := w.debounceMap[event.Path]
	if !exists {
		entry = &DebouncedEvent{
			FirstEvent:  event,
			LastEvent:   event,
			EventCount:  1,
			EventType:   event.Type,
			WindowStart: time.Now(),
			WindowEnd:   time.Now(),
		}
		w.debounceMap[event.Path] = entry
	} else {
		entry.LastEvent = event
		entry.EventCount++
		entry.WindowEnd = time.Now()

		// Update most common event type if needed
		if w.config.Debounce.MergeEvents {
			entry.EventType = event.Type
		}
	}

	// Set timer to flush debounced events
	w.debounceT = time.AfterFunc(w.config.Debounce.MaxWait, w.flushDebouncedEvents)
}

// flushDebouncedEvents flushes all debounced events.
func (w *Watcher) flushDebouncedEvents() {
	w.debounceMu.Lock()
	defer w.debounceMu.Unlock()

	if len(w.debounceMap) == 0 {
		return
	}

	events := make([]Event, 0, len(w.debounceMap))

	for _, entry := range w.debounceMap {
		events = append(events, entry.LastEvent)
	}

	w.debounceMap = make(map[string]*DebouncedEvent)

	// Update stats
	w.statsLock.Lock()
	w.stats.DebouncedEvents += int64(len(events))
	w.statsLock.Unlock()

	// Handle events
	if w.config.Handler != nil {
		w.config.Handler(events)
	}
}

// handleBatchedEvent handles an event with batching.
func (w *Watcher) handleBatchedEvent(event Event) {
	w.batchMu.Lock()
	w.batchBuffer = append(w.batchBuffer, event)

	// Check if buffer is full
	if len(w.batchBuffer) >= w.config.Batch.MaxSize {
		w.flushBatch()
		w.batchMu.Unlock()
		return
	}

	// Set/reset timer
	if w.batchT != nil {
		w.batchT.Stop()
	}

	w.batchT = time.AfterFunc(w.config.Batch.MaxWait, w.flushBatch)
	w.batchMu.Unlock()
}

// flushBatch flushes the batch buffer.
func (w *Watcher) flushBatch() {
	w.batchMu.Lock()
	defer w.batchMu.Unlock()

	if len(w.batchBuffer) == 0 {
		return
	}

	batch := make([]Event, len(w.batchBuffer))
	copy(batch, w.batchBuffer)
	w.batchBuffer = w.batchBuffer[:0]

	// Update stats
	w.statsLock.Lock()
	w.stats.BatchedEvents++
	w.statsLock.Unlock()

	// Handle batch
	if w.config.Handler != nil {
		w.config.Handler(batch)
	}
}

// processBatches processes batches from the batch channel.
func (w *Watcher) processBatches() {
	defer close(w.batchChan)

	for {
		select {
		case <-w.stopChan:
			// Flush remaining batch before exiting
			w.batchMu.Lock()
			if len(w.batchBuffer) > 0 {
				w.flushBatch()
			}
			w.batchMu.Unlock()
			return

		case batch := <-w.batchChan:
			if w.config.Handler != nil {
				w.config.Handler(batch)
			}
		}
	}
}

// IsWatching returns true if a path is being watched.
func (w *Watcher) IsWatching(path string) bool {
	absPath, _ := filepath.Abs(path)
	w.watchedLock.Lock()
	defer w.watchedLock.Unlock()
	return w.watchedPaths[absPath]
}

// WatchedPaths returns a list of all watched paths.
func (w *Watcher) WatchedPaths() []string {
	w.watchedLock.Lock()
	defer w.watchedLock.Unlock()

	paths := make([]string, 0, len(w.watchedPaths))
	for path := range w.watchedPaths {
		paths = append(paths, path)
	}

	return paths
}

// WatchedCount returns the number of watched paths.
func (w *Watcher) WatchedCount() int {
	w.watchedLock.Lock()
	defer w.watchedLock.Unlock()
	return len(w.watchedPaths)
}

// Event returns the next event (blocking).
func (w *Watcher) Event() (Event, bool) {
	event, ok := <-w.eventChan
	return event, ok
}

// TryEvent attempts to get the next event without blocking.
func (w *Watcher) TryEvent() (Event, bool) {
	select {
	case event := <-w.eventChan:
		return event, true
	default:
		return Event{}, false
	}
}

// EventChan returns the event channel.
func (w *Watcher) EventChan() <-chan Event {
	return w.eventChan
}

// Batch returns the next batch (blocking).
func (w *Watcher) Batch() ([]Event, bool) {
	batch, ok := <-w.batchChan
	return batch, ok
}

// TryBatch attempts to get the next batch without blocking.
func (w *Watcher) TryBatch() ([]Event, bool) {
	select {
	case batch := <-w.batchChan:
		return batch, true
	default:
		return nil, false
	}
}

// BatchChan returns the batch channel.
func (w *Watcher) BatchChan() <-chan []Event {
	return w.batchChan
}

// Close closes the watcher and releases resources.
func (w *Watcher) Close() error {
	w.stateLock.Lock()
	defer w.stateLock.Unlock()

	if w.state == WatcherIdle {
		return nil
	}

	if w.state == WatcherRunning || w.state == WatcherPaused {
		close(w.stopChan)
		w.state = WatcherStopped
	}

	return w.fsWatcher.Close()
}

// IsTerminal returns true if standard output is a terminal (can use colors).
func IsTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// IsTerminalColor returns true if terminal supports colors.
func IsTerminalColor() bool {
	if !IsTerminal() {
		return false
	}

	// Check for common terminal color support indicators
	return os.Getenv("TERM") != "dumb" && os.Getenv("NO_COLOR") == ""
}

// Watch creates a simple watcher with a handler and starts it immediately.
func Watch(path string, handler Handler, errorHandler ErrorHandler) (*Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := w.Add(path); err != nil {
		w.Close()
		return nil, err
	}

	w.config.Handler = handler
	w.config.ErrorHandler = errorHandler

	if err := w.Start(); err != nil {
		w.Close()
		return nil, err
	}

	return w, nil
}

// WatchRecursive creates a watcher that watches path recursively.
func WatchRecursive(path string, handler Handler, errorHandler ErrorHandler) (*Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := w.AddRecursive(path); err != nil {
		w.Close()
		return nil, err
	}

	w.config.Handler = handler
	w.config.ErrorHandler = errorHandler

	if err := w.Start(); err != nil {
		w.Close()
		return nil, err
	}

	return w, nil
}

// WatchFiles watches specific files (not directories).
func WatchFiles(files []string, handler Handler, errorHandler ErrorHandler) (*Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if err := w.Add(file); err != nil {
			w.Close()
			return nil, fmt.Errorf("failed to watch %s: %w", file, err)
		}
	}

	w.config.Handler = handler
	w.config.ErrorHandler = errorHandler
	w.config.Filter.IncludeDirs = false

	if err := w.Start(); err != nil {
		w.Close()
		return nil, err
	}

	return w, nil
}

var _ error = (*Error)(nil)
var _ error = ErrWatcherStopped

// Is tests if error matches a specific error type.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As returns the target error if the error matches.
func As(err error, target any) bool {
	return errors.As(err, target)
}
