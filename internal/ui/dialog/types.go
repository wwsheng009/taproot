package dialog

import (
	"github.com/wwsheng009/taproot/internal/ui/list"
	"github.com/wwsheng009/taproot/internal/ui/render"
)

// DialogID uniquely identifies a dialog.
type DialogID string

// Dialog is the engine-agnostic interface for dialog components.
type Dialog interface {
	render.Model
	// ID returns the unique identifier for this dialog.
	ID() DialogID
	// Title returns the dialog title.
	Title() string
	// SetSize updates the dialog dimensions.
	SetSize(width, height int)
}

// ActionResult represents the result of a dialog action.
type ActionResult int

const (
	// ActionNone indicates no action was taken.
	ActionNone ActionResult = iota
	// ActionConfirm indicates the user confirmed.
	ActionConfirm
	// ActionCancel indicates the user cancelled.
	ActionCancel
	// ActionDismiss indicates the dialog was dismissed.
	ActionDismiss
)

// Callback is a function called when a dialog action occurs.
type Callback func(result ActionResult, data any)

// Button represents a dialog button with state.
type Button struct {
	label    string
	selected bool
	primary  bool
}

// NewButton creates a new button.
func NewButton(label string, primary bool) *Button {
	return &Button{
		label:    label,
		selected: false,
		primary:  primary,
	}
}

// Label returns the button label.
func (b *Button) Label() string {
	return b.label
}

// SetLabel sets the button label.
func (b *Button) SetLabel(label string) {
	b.label = label
}

// Selected returns true if the button is selected.
func (b *Button) Selected() bool {
	return b.selected
}

// SetSelected sets the selection state.
func (b *Button) SetSelected(selected bool) {
	b.selected = selected
}

// Primary returns true if this is the primary action button.
func (b *Button) Primary() bool {
	return b.primary
}

// SetPrimary sets whether this is the primary button.
func (b *Button) SetPrimary(primary bool) {
	b.primary = primary
}

// Toggle reverses the selection state.
func (b *Button) Toggle() {
	b.selected = !b.selected
}

// ButtonGroup manages a group of buttons with selection.
type ButtonGroup struct {
	buttons   []*Button
	selected  int
	horizontal bool
}

// NewButtonGroup creates a new button group.
func NewButtonGroup(horizontal bool) *ButtonGroup {
	return &ButtonGroup{
		buttons:    []*Button{},
		selected:   0,
		horizontal: horizontal,
	}
}

// Add adds a button to the group.
func (bg *ButtonGroup) Add(label string, primary bool) *ButtonGroup {
	bg.buttons = append(bg.buttons, NewButton(label, primary))
	return bg
}

// Buttons returns the buttons.
func (bg *ButtonGroup) Buttons() []*Button {
	return bg.buttons
}

// SetButtons sets the buttons.
func (bg *ButtonGroup) SetButtons(buttons []*Button) {
	bg.buttons = buttons
}

// Selected returns the index of the selected button.
func (bg *ButtonGroup) Selected() int {
	return bg.selected
}

// SetSelected sets the selected button index.
func (bg *ButtonGroup) SetSelected(index int) {
	// Deselect current
	if bg.selected >= 0 && bg.selected < len(bg.buttons) {
		bg.buttons[bg.selected].SetSelected(false)
	}
	// Set new selection
	bg.selected = index
	if bg.selected >= 0 && bg.selected < len(bg.buttons) {
		bg.buttons[bg.selected].SetSelected(true)
	}
}

// SelectNext moves selection to the next button.
func (bg *ButtonGroup) SelectNext() {
	if len(bg.buttons) == 0 {
		return
	}
	next := bg.selected + 1
	if next >= len(bg.buttons) {
		next = 0
	}
	bg.SetSelected(next)
}

// SelectPrev moves selection to the previous button.
func (bg *ButtonGroup) SelectPrev() {
	if len(bg.buttons) == 0 {
		return
	}
	prev := bg.selected - 1
	if prev < 0 {
		prev = len(bg.buttons) - 1
	}
	bg.SetSelected(prev)
}

// SelectedButton returns the currently selected button.
func (bg *ButtonGroup) SelectedButton() *Button {
	if bg.selected >= 0 && bg.selected < len(bg.buttons) {
		return bg.buttons[bg.selected]
	}
	return nil
}

// Count returns the number of buttons.
func (bg *ButtonGroup) Count() int {
	return len(bg.buttons)
}

// Horizontal returns true if buttons are arranged horizontally.
func (bg *ButtonGroup) Horizontal() bool {
	return bg.horizontal
}

// InputField represents a text input field.
type InputField struct {
	value       string
	placeholder string
	focused     bool
	hidden      bool // For passwords
	cursor      int
	width       int
	maxLength   int
}

// NewInputField creates a new input field.
func NewInputField(placeholder string) *InputField {
	return &InputField{
		placeholder: placeholder,
		focused:     false,
		hidden:      false,
		cursor:      0,
		width:       40,
		maxLength:   0, // 0 means no limit
	}
}

// Value returns the current input value.
func (i *InputField) Value() string {
	return i.value
}

// SetValue sets the input value.
func (i *InputField) SetValue(value string) {
	i.value = value
	i.setCursor(len(value))
}

// Placeholder returns the placeholder text.
func (i *InputField) Placeholder() string {
	return i.placeholder
}

// SetPlaceholder sets the placeholder text.
func (i *InputField) SetPlaceholder(placeholder string) {
	i.placeholder = placeholder
}

// Focused returns true if the input is focused.
func (i *InputField) Focused() bool {
	return i.focused
}

// SetFocused sets the focused state.
func (i *InputField) SetFocused(focused bool) {
	i.focused = focused
}

// Hidden returns true if input is hidden (password).
func (i *InputField) Hidden() bool {
	return i.hidden
}

// SetHidden sets whether input is hidden.
func (i *InputField) SetHidden(hidden bool) {
	i.hidden = hidden
}

// Cursor returns the cursor position.
func (i *InputField) Cursor() int {
	return i.cursor
}

// Width returns the input width.
func (i *InputField) Width() int {
	return i.width
}

// SetWidth sets the input width.
func (i *InputField) SetWidth(width int) {
	i.width = width
}

// MaxLength returns the max length limit.
func (i *InputField) MaxLength() int {
	return i.maxLength
}

// SetMaxLength sets the max length limit (0 for no limit).
func (i *InputField) SetMaxLength(max int) {
	i.maxLength = max
}

// Insert inserts a rune at the cursor position.
func (i *InputField) Insert(r rune) {
	if i.maxLength > 0 && len(i.value) >= i.maxLength {
		return
	}
	
	before := i.value[:i.cursor]
	after := i.value[i.cursor:]
	i.value = before + string(r) + after
	i.cursor++
}

// Delete deletes the rune before the cursor.
func (i *InputField) Delete() {
	if i.cursor > 0 {
		before := i.value[:i.cursor-1]
		after := i.value[i.cursor:]
		i.value = before + after
		i.cursor--
	}
}

// DeleteForward deletes the rune at the cursor.
func (i *InputField) DeleteForward() {
	if i.cursor < len(i.value) {
		before := i.value[:i.cursor]
		after := i.value[i.cursor+1:]
		i.value = before + after
	}
}

// MoveCursorLeft moves the cursor left.
func (i *InputField) MoveCursorLeft() {
	if i.cursor > 0 {
		i.cursor--
	}
}

// MoveCursorRight moves the cursor right.
func (i *InputField) MoveCursorRight() {
	if i.cursor < len(i.value) {
		i.cursor++
	}
}

// MoveCursorToStart moves cursor to start.
func (i *InputField) MoveCursorToStart() {
	i.cursor = 0
}

// MoveCursorToEnd moves cursor to end.
func (i *InputField) MoveCursorToEnd() {
	i.cursor = len(i.value)
}

// Clear clears the input value.
func (i *InputField) Clear() {
	i.value = ""
	i.cursor = 0
}

// setCursor sets the cursor position safely.
func (i *InputField) setCursor(pos int) {
	if pos < 0 {
		pos = 0
	} else if pos > len(i.value) {
		pos = len(i.value)
	}
	i.cursor = pos
}

// SelectList is a dropdown selection list.
type SelectList struct {
	items      []list.Item
	selected   int
	viewport   *list.Viewport
	expanded   bool
	focused    bool
	maxVisible int
}

// NewSelectList creates a new select list.
func NewSelectList(items []list.Item) *SelectList {
	return &SelectList{
		items:      items,
		selected:   0,
		viewport:   list.NewViewport(5, len(items)),
		expanded:   false,
		focused:    false,
		maxVisible: 5,
	}
}

// Items returns the list items.
func (s *SelectList) Items() []list.Item {
	return s.items
}

// SetItems sets the list items.
func (s *SelectList) SetItems(items []list.Item) {
	s.items = items
	s.viewport.SetTotal(len(items))
}

// Selected returns the index of the selected item.
func (s *SelectList) Selected() int {
	return s.selected
}

// SetSelected sets the selected index.
func (s *SelectList) SetSelected(index int) {
	if index >= 0 && index < len(s.items) {
		s.selected = index
		s.viewport.SetCursor(index)
	}
}

// SelectedItem returns the currently selected item.
func (s *SelectList) SelectedItem() list.Item {
	if s.selected >= 0 && s.selected < len(s.items) {
		return s.items[s.selected]
	}
	return nil
}

// Expanded returns true if the list is expanded.
func (s *SelectList) Expanded() bool {
	return s.expanded
}

// SetExpanded sets the expanded state.
func (s *SelectList) SetExpanded(expanded bool) {
	s.expanded = expanded
}

// Toggle toggles the expanded state.
func (s *SelectList) Toggle() {
	s.expanded = !s.expanded
}

// Focused returns true if the select is focused.
func (s *SelectList) Focused() bool {
	return s.focused
}

// SetFocused sets the focused state.
func (s *SelectList) SetFocused(focused bool) {
	s.focused = focused
}

// MoveUp moves selection up.
func (s *SelectList) MoveUp() {
	if s.selected > 0 {
		s.selected--
		s.viewport.MoveUp()
	}
}

// MoveDown moves selection down.
func (s *SelectList) MoveDown() {
	if s.selected < len(s.items)-1 {
		s.selected++
		s.viewport.MoveDown()
	}
}

// MaxVisible returns the maximum visible items.
func (s *SelectList) MaxVisible() int {
	return s.maxVisible
}

// SetMaxVisible sets the maximum visible items.
func (s *SelectList) SetMaxVisible(max int) {
	s.maxVisible = max
	s.viewport.SetVisible(max)
}
