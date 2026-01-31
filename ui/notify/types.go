package notify

import (
	"time"
)

// Type represents the type of notification
type Type int

const (
	TypeInfo Type = iota
	TypeSuccess
	TypeWarning
	TypeError
)

// Notification represents a single notification message
type Notification struct {
	ID       string
	Type     Type
	Title    string
	Message  string
	Duration time.Duration
	// CreatedAt is used to calculate remaining time
	CreatedAt time.Time
	// Animating indicates if the notification is currently animating (in or out)
	Animating bool
}

// Config holds global configuration for the notification system
type Config struct {
	// DefaultDuration is the default duration for notifications
	DefaultDuration time.Duration
	// MaxVisible is the maximum number of visible notifications
	MaxVisible int
	// Position determines where notifications appear (TODO: implement positioning)
	Position Position
}

// Position determines where notifications appear on screen
type Position int

const (
	TopRight Position = iota
	BottomRight
	TopLeft
	BottomLeft
	TopCenter
	BottomCenter
)

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		DefaultDuration: 5 * time.Second,
		MaxVisible:      5,
		Position:        TopRight,
	}
}

// Msg is the message type for notification events
type Msg struct {
	Notification Notification
}

// ShowNotificationMsg is a specific message to show a notification
type ShowNotificationMsg struct {
	Notification Notification
}
