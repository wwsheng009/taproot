package notify

import (
	"fmt"
	"time"

	"github.com/wwsheng009/taproot/ui/render"
)

// New creates a command to show a notification
func New(title, message string, typ Type, duration time.Duration) render.Cmd {
	return func() render.Msg {
		return ShowNotificationMsg{
			Notification: Notification{
				ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
				Type:     typ,
				Title:    title,
				Message:  message,
				Duration: duration,
			},
		}
	}
}

// Info creates a command to show an info notification
func Info(title, message string) render.Cmd {
	return New(title, message, TypeInfo, 5*time.Second)
}

// Success creates a command to show a success notification
func Success(title, message string) render.Cmd {
	return New(title, message, TypeSuccess, 5*time.Second)
}

// Warn creates a command to show a warning notification
func Warn(title, message string) render.Cmd {
	return New(title, message, TypeWarning, 7*time.Second)
}

// Error creates a command to show an error notification
func Error(title, message string) render.Cmd {
	return New(title, message, TypeError, 10*time.Second)
}
