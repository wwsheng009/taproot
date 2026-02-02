// Package watcher provides file system monitoring utilities.
package watcher

import (
	"time"
)

// EventType represents the type of file system event.
type EventType int

const (
	// EventCreate is triggered when a file is created
	EventCreate EventType = iota
	// EventWrite is triggered when a file is written to
	EventWrite
	// EventRemove is triggered when a file is removed
	EventRemove
	// EventRename is triggered when a file is renamed
	EventRename
	// EventChmod is triggered when file permissions change
	EventChmod
)

// String returns the string representation of the event type.
func (e EventType) String() string {
	switch e {
	case EventCreate:
		return "CREATE"
	case EventWrite:
		return "WRITE"
	case EventRemove:
		return "REMOVE"
	case EventRename:
		return "RENAME"
	case EventChmod:
		return "CHMOD"
	default:
		return "UNKNOWN"
	}
}

// Event represents a file system event.
type Event struct {
	Type      EventType // Type of event
	Path      string    // Path of the file/directory
	OldPath   string    // Previous path (for rename events)
	Timestamp time.Time // When the event occurred
}

// Filter defines criteria for filtering events.
type Filter struct {
	// IncludePatterns is a list of glob patterns for files to include
	IncludePatterns []string

	// ExcludePatterns is a list of glob patterns for files to exclude
	ExcludePatterns []string

	// IncludeDirs specifies whether to include directory events
	IncludeDirs bool

	// IgnoreHidden specifies whether to ignore hidden files/directories
	IgnoreHidden bool

	// EventTypes specifies which event types to include
	// If empty, all event types are included
	EventTypes []EventType

	// MinSize specifies minimum file size to track (in bytes)
	// 0 means no minimum size
	MinSize int64

	// MaxSize specifies maximum file size to track (in bytes)
	// 0 means no maximum size
	MaxSize int64

	// Extensions specifies file extensions to include
	// If empty, all extensions are included
	Extensions []string

	// Recurse specifies whether to watch subdirectories
	Recurse bool
}

// WatcherState represents the state of a file watcher.
type WatcherState int

const (
	// WatcherIdle is the initial state before watching starts
	WatcherIdle WatcherState = iota
	// WatcherRunning indicates the watcher is actively monitoring
	WatcherRunning
	// WatcherPaused indicates the watcher is paused but not stopped
	WatcherPaused
	// WatcherStopped indicates the watcher has been stopped
	WatcherStopped
	// WatcherError indicates an error occurred
	WatcherError
)

// String returns the string representation of the watcher state.
func (s WatcherState) String() string {
	switch s {
	case WatcherIdle:
		return "idle"
	case WatcherRunning:
		return "running"
	case WatcherPaused:
		return "paused"
	case WatcherStopped:
		return "stopped"
	case WatcherError:
		return "error"
	default:
		return "unknown"
	}
}

// DebounceConfig configures debouncing behavior.
type DebounceConfig struct {
	// Enabled specifies whether debouncing is enabled
	Enabled bool

	// Delay is the time to wait before collecting events
	// Events within this window are debounced
	Delay time.Duration

	// MaxWait is the maximum time to wait before emitting events
	// This ensures events are eventually emitted even if they keep coming
	MaxWait time.Duration

	// MergeEvents specifies whether to merge consecutive events of the same type for the same file
	MergeEvents bool

	// MergeWindow is the time window within which events are eligible for merging
	MergeWindow time.Duration
}

// DefaultDebounceConfig returns default debounce configuration.
func DefaultDebounceConfig() DebounceConfig {
	return DebounceConfig{
		Enabled:      true,
		Delay:        100 * time.Millisecond,
		MaxWait:      1 * time.Second,
		MergeEvents:  true,
		MergeWindow:  500 * time.Millisecond,
	}
}

// BatchConfig configures batch processing of events.
type BatchConfig struct {
	// Enabled specifies whether batching is enabled
	Enabled bool

	// MaxSize is the maximum number of events in a batch
	MaxSize int

	// MaxWait is the maximum time to wait before emitting a batch
	MaxWait time.Duration

	// MinSize is the minimum number of events required to emit a batch
	// If 0, any number of events triggers a batch
	MinSize int
}

// DefaultBatchConfig returns default batch configuration.
func DefaultBatchConfig() BatchConfig {
	return BatchConfig{
		Enabled:   false,
		MaxSize:   100,
		MaxWait:   500 * time.Millisecond,
		MinSize:   1,
	}
}

// Handler is a callback function that receives events.
type Handler func(events []Event)

// ErrorHandler handles watcher errors.
type ErrorHandler func(err error)

// Config configures a file watcher.
type Config struct {
	// Filter specifies which events to watch for
	Filter Filter

	// Debounce configures debouncing behavior
	Debounce DebounceConfig

	// Batch configures batch processing behavior
	Batch BatchConfig

	// Handler is called when events occur
	Handler Handler

	// ErrorHandler is called when errors occur
	ErrorHandler ErrorHandler

	// Reopen specifies whether to reopen files when they are renamed/removed
	Reopen bool

	// BufferSize is the size of the event buffer
	BufferSize int
}

// DefaultConfig returns default watcher configuration.
func DefaultConfig() Config {
	return Config{
		Filter:        Filter{},
		Debounce:      DefaultDebounceConfig(),
		Batch:         DefaultBatchConfig(),
		Handler:       nil,
		ErrorHandler:  nil,
		Reopen:        false,
		BufferSize:    1000,
	}
}

// Stats holds statistics about the watcher.
type Stats struct {
	TotalEvents      int64
	DroppedEvents    int64
	DebouncedEvents  int64
	BatchedEvents    int64
	LastEventTime    time.Time
	StartTime        time.Time
	FilesWatched     int
	DirsWatched      int
}

// DebouncedEvent represents a debounced collection of events for a file.
type DebouncedEvent struct {
	LastEvent    Event       // Last event in the collection
	FirstEvent   Event       // First event in the collection
	EventCount   int        // Number of events in the collection
	EventType    EventType  // Most common event type
	WindowStart  time.Time  // Start of the debouncing window
	WindowEnd    time.Time  // End of the debouncing window
}

// Errors
var (
	ErrWatcherStopped   = &Error{Code: "WATCHER_STOPPED", Message: "watcher is stopped"}
	ErrWatcherRunning   = &Error{Code: "WATCHER_RUNNING", Message: "watcher is already running"}
	ErrInvalidPath      = &Error{Code: "INVALID_PATH", Message: "invalid path"}
	ErrPathNotFound     = &Error{Code: "PATH_NOT_FOUND", Message: "path does not exist"}
	ErrNotDirectory     = &Error{Code: "NOT_DIRECTORY", Message: "path is not a directory"}
	ErrAlreadyWatching  = &Error{Code: "ALREADY_WATCHING", Message: "already watching path"}
	ErrNotWatching      = &Error{Code: "NOT_WATCHING", Message: "not watching path"}
	ErrNoHandler        = &Error{Code: "NO_HANDLER", Message: "no event handler configured"}
	ErrBufferOverflow   = &Error{Code: "BUFFER_OVERFLOW", Message: "event buffer overflow, events dropped"}
	ErrInvalidFilter    = &Error{Code: "INVALID_FILTER", Message: "invalid filter configuration"}
)

// Error represents a watcher error.
type Error struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	return e.Cause
}
