package dialog

// OpenDialogMsg is sent to open a new dialog.
type OpenDialogMsg struct {
	Dialog Dialog
}

// CloseDialogMsg is sent to close the active dialog.
type CloseDialogMsg struct{}

// CloseAllDialogsMsg is sent to close all dialogs.
type CloseAllDialogsMsg struct{}

// DialogActionMsg is sent when a dialog action occurs.
type DialogActionMsg struct {
	DialogID DialogID
	Result  ActionResult
	Data    any
}

// DialogResultMsg is sent with the final dialog result.
type DialogResultMsg struct {
	DialogID DialogID
	Result   ActionResult
	Data     any
}

// Overlay manages a stack of dialogs with proper message routing.
type Overlay struct {
	dialogs []Dialog
	width   int
	height  int
}

// NewOverlay creates a new dialog overlay.
func NewOverlay() *Overlay {
	return &Overlay{
		dialogs: []Dialog{},
		width:   80,
		height:  24,
	}
}

// Push adds a dialog to the top of the stack.
func (o *Overlay) Push(d Dialog) {
	o.dialogs = append(o.dialogs, d)
}

// Pop removes the top dialog from the stack.
func (o *Overlay) Pop() Dialog {
	if len(o.dialogs) == 0 {
		return nil
	}
	last := o.dialogs[len(o.dialogs)-1]
	o.dialogs = o.dialogs[:len(o.dialogs)-1]
	return last
}

// Peek returns the top dialog without removing it.
func (o *Overlay) Peek() Dialog {
	if len(o.dialogs) == 0 {
		return nil
	}
	return o.dialogs[len(o.dialogs)-1]
}

// HasDialogs returns true if there are any dialogs.
func (o *Overlay) HasDialogs() bool {
	return len(o.dialogs) > 0
}

// Count returns the number of dialogs.
func (o *Overlay) Count() int {
	return len(o.dialogs)
}

// Dialogs returns all dialogs in the stack.
func (o *Overlay) Dialogs() []Dialog {
	return o.dialogs
}

// Clear removes all dialogs.
func (o *Overlay) Clear() {
	o.dialogs = []Dialog{}
}

// SetSize sets the overlay dimensions.
func (o *Overlay) SetSize(width, height int) {
	o.width = width
	o.height = height
	for _, d := range o.dialogs {
		d.SetSize(width, height)
	}
}

// Size returns the overlay dimensions.
func (o *Overlay) Size() (width, height int) {
	return o.width, o.height
}

// FindByID finds a dialog by its ID.
func (o *Overlay) FindByID(id DialogID) Dialog {
	for i := len(o.dialogs) - 1; i >= 0; i-- {
		if o.dialogs[i].ID() == id {
			return o.dialogs[i]
		}
	}
	return nil
}

// RemoveByID removes a dialog by its ID.
func (o *Overlay) RemoveByID(id DialogID) Dialog {
	for i := len(o.dialogs) - 1; i >= 0; i-- {
		if o.dialogs[i].ID() == id {
			d := o.dialogs[i]
			o.dialogs = append(o.dialogs[:i], o.dialogs[i+1:]...)
			return d
		}
	}
	return nil
}

// ActiveDialog returns the topmost (active) dialog.
func (o *Overlay) ActiveDialog() Dialog {
	return o.Peek()
}

// IsActive returns true if the given dialog is the active one.
func (o *Overlay) IsActive(d Dialog) bool {
	active := o.ActiveDialog()
	return active != nil && active.ID() == d.ID()
}

// UpdatePosition updates the position of all dialogs.
func (o *Overlay) UpdatePosition(width, height int) {
	o.width = width
	o.height = height
}

// DialogBounds calculates the bounds for a dialog.
type DialogBounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

// CalculateBounds calculates default dialog bounds centered on screen.
func CalculateBounds(dialogWidth, dialogHeight, screenWidth, screenHeight int) DialogBounds {
	x := (screenWidth - dialogWidth) / 2
	y := (screenHeight - dialogHeight) / 2
	
	// Ensure bounds are within screen
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x + dialogWidth > screenWidth {
		x = screenWidth - dialogWidth
	}
	if y + dialogHeight > screenHeight {
		y = screenHeight - dialogHeight
	}
	
	return DialogBounds{
		X:      x,
		Y:      y,
		Width:  dialogWidth,
		Height: dialogHeight,
	}
}

// DefaultWidth returns a reasonable default width for dialogs.
func DefaultWidth() int {
	return 60
}

// DefaultHeight returns a reasonable default height for dialogs.
func DefaultHeight() int {
	return 15
}

// MaxWidth returns the maximum width for dialogs.
func MaxWidth(screenWidth int) int {
	max := screenWidth - 8 // Leave padding
	if max < 40 {
		return 40
	}
	if max > 100 {
		return 100
	}
	return max
}

// MaxHeight returns the maximum height for dialogs.
func MaxHeight(screenHeight int) int {
	max := screenHeight - 6 // Leave padding
	if max < 10 {
		return 10
	}
	if max > 30 {
		return 30
	}
	return max
}
