package forms

import (
	"sync/atomic"
	"time"

	"github.com/wwsheng009/taproot/ui/render"
)

// Global counter for unique blink IDs to prevent collisions between components
var blinkCounter int64

// NextBlinkID returns a globally unique blink ID.
func NextBlinkID() int {
	return int(atomic.AddInt64(&blinkCounter, 1))
}

// BlinkMsg is sent when the cursor should blink.
type BlinkMsg struct {
	id int
}

// BlinkCmd returns a command that waits and sends a BlinkMsg.
// Duration is set to 500ms for standard cursor blink rate.
func BlinkCmd(id int) render.Cmd {
	return func() render.Msg {
		time.Sleep(500 * time.Millisecond)
		return BlinkMsg{id: id}
	}
}
