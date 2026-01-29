package list

// Viewport manages the visible window into a larger list.
type Viewport struct {
	offset   int
	visible  int
	total    int
	cursor   int
}

// NewViewport creates a new viewport.
func NewViewport(visible, total int) *Viewport {
	return &Viewport{
		offset:   0,
		visible:  visible,
		total:    total,
		cursor:   0,
	}
}

// SetTotal updates the total number of items.
func (v *Viewport) SetTotal(total int) {
	v.total = total
	v.clamp()
}

// SetVisible updates the number of visible items.
func (v *Viewport) SetVisible(visible int) {
	v.visible = visible
	v.clamp()
}

// SetCursor updates the cursor position and adjusts offset if needed.
func (v *Viewport) SetCursor(cursor int) {
	v.cursor = cursor
	v.adjustOffset()
	v.clamp()
}

// Cursor returns the current cursor position.
func (v *Viewport) Cursor() int {
	return v.cursor
}

// Offset returns the current scroll offset.
func (v *Viewport) Offset() int {
	return v.offset
}

// Visible returns the number of visible items.
func (v *Viewport) Visible() int {
	return v.visible
}

// Total returns the total number of items.
func (v *Viewport) Total() int {
	return v.total
}

// Range returns the start and end indices of the visible range.
func (v *Viewport) Range() (start, end int) {
	start = v.offset
	end = min(start+v.visible, v.total)
	return start, end
}

// MoveUp moves the cursor up by one item.
func (v *Viewport) MoveUp() {
	if v.cursor > 0 {
		v.cursor--
		if v.cursor < v.offset {
			v.offset--
		}
	}
}

// MoveDown moves the cursor down by one item.
func (v *Viewport) MoveDown() {
	if v.cursor < v.total-1 {
		v.cursor++
		if v.cursor >= v.offset+v.visible {
			v.offset++
		}
	}
}

// MoveTo moves the cursor to a specific position.
func (v *Viewport) MoveTo(pos int) {
	v.cursor = pos
	v.adjustOffset()
	v.clamp()
}

// MoveToTop moves the cursor to the first item.
func (v *Viewport) MoveToTop() {
	v.cursor = 0
	v.offset = 0
}

// MoveToBottom moves the cursor to the last item.
func (v *Viewport) MoveToBottom() {
	v.cursor = max(0, v.total-1)
	v.offset = max(0, v.total-v.visible)
}

// PageUp moves the cursor up by one page.
func (v *Viewport) PageUp() {
	v.cursor = max(0, v.cursor-v.visible)
	v.offset = max(0, v.cursor-v.visible+1)
}

// PageDown moves the cursor down by one page.
func (v *Viewport) PageDown() {
	v.cursor = min(v.total-1, v.cursor+v.visible)
	v.adjustOffset()
}

// IsFirst returns true if the cursor is at the first item.
func (v *Viewport) IsFirst() bool {
	return v.cursor == 0
}

// IsLast returns true if the cursor is at the last item.
func (v *Viewport) IsLast() bool {
	return v.cursor == v.total-1 || v.total == 0
}

// IsVisible returns true if the given index is currently visible.
func (v *Viewport) IsVisible(index int) bool {
	return index >= v.offset && index < v.offset+v.visible
}

// ScrollTo makes a specific index visible by adjusting the offset.
func (v *Viewport) ScrollTo(index int) {
	if index < 0 || index >= v.total {
		return
	}
	if index < v.offset {
		v.offset = index
	} else if index >= v.offset+v.visible {
		v.offset = index - v.visible + 1
	}
}

// clamp ensures cursor and offset are within valid bounds.
func (v *Viewport) clamp() {
	if v.total == 0 {
		v.cursor = 0
		v.offset = 0
		return
	}
	if v.cursor < 0 {
		v.cursor = 0
	} else if v.cursor >= v.total {
		v.cursor = v.total - 1
	}
	if v.offset < 0 {
		v.offset = 0
	} else if v.offset > v.total-v.visible {
		v.offset = max(0, v.total-v.visible)
	}
}

// adjustOffset adjusts offset to keep cursor in viewport.
func (v *Viewport) adjustOffset() {
	if v.cursor < v.offset {
		v.offset = v.cursor
	} else if v.cursor >= v.offset+v.visible {
		v.offset = v.cursor - v.visible + 1
	}
}

// HasScroll returns true if scrolling is available (content > viewport).
func (v *Viewport) HasScroll() bool {
	return v.total > v.visible
}

// CanScrollUp returns true if we can scroll up.
func (v *Viewport) CanScrollUp() bool {
	return v.offset > 0
}

// CanScrollDown returns true if we can scroll down.
func (v *Viewport) CanScrollDown() bool {
	return v.offset+v.visible < v.total
}

// ScrollIndicator returns a string representation of scroll state.
func (v *Viewport) ScrollIndicator() string {
	if !v.HasScroll() {
		return "All"
	}
	if v.CanScrollUp() && v.CanScrollDown() {
		return "▲ ▼"
	}
	if v.CanScrollUp() {
		return "▲"
	}
	if v.CanScrollDown() {
		return "▼"
	}
	return "All"
}
