package list

// Action represents a keyboard action that can be performed on a list.
type Action int

const (
	// ActionNone represents no action.
	ActionNone Action = iota
	// ActionMoveUp moves the cursor up.
	ActionMoveUp
	// ActionMoveDown moves the cursor down.
	ActionMoveDown
	// ActionMoveLeft moves the cursor left (for nested lists).
	ActionMoveLeft
	// ActionMoveRight moves the cursor right (for nested lists).
	ActionMoveRight
	// ActionPageUp moves up by a page.
	ActionPageUp
	// ActionPageDown moves down by a page.
	ActionPageDown
	// ActionMoveToTop jumps to the first item.
	ActionMoveToTop
	// ActionMoveToBottom jumps to the last item.
	ActionMoveToBottom
	// ActionSelect selects the current item.
	ActionSelect
	// ActionDeselect deselects the current item.
	ActionDeselect
	// ActionToggleSelection toggles selection of the current item.
	ActionToggleSelection
	// ActionSelectAll selects all items.
	ActionSelectAll
	// ActionDeselectAll deselects all items.
	ActionDeselectAll
	// ActionInvertSelection inverts the selection.
	ActionInvertSelection
	// ActionConfirm confirms the current selection.
	ActionConfirm
	// ActionCancel cancels the current operation.
	ActionCancel
	// ActionDelete deletes the current item.
	ActionDelete
	// ActionEdit edits the current item.
	ActionEdit
	// ActionNew creates a new item.
	ActionNew
	// ActionFilter enters filter mode.
	ActionFilter
	// ActionFilterClear clears the filter.
	ActionFilterClear
	// ActionToggleGroup toggles expansion of current group.
	ActionToggleGroup
	// ActionExpandAll expands all groups.
	ActionExpandAll
	// ActionCollapseAll collapses all groups.
	ActionCollapseAll
	// ActionHelp shows help.
	ActionHelp
	// ActionQuit quits the application.
	ActionQuit
)

// KeyBinding maps a key string to an Action.
type KeyBinding struct {
	Key    string
	Action Action
}

// KeyMap defines keyboard shortcuts for list actions.
type KeyMap struct {
	// Navigation
	Up       []string
	Down     []string
	Left     []string
	Right    []string
	PageUp   []string
	PageDown []string
	Home     []string
	End      []string

	// Selection
	Select       []string
	Deselect     []string
	ToggleSelect []string
	SelectAll    []string
	DeselectAll  []string
	InvertSelect []string

	// Actions
	Confirm []string
	Cancel  []string
	Delete  []string
	Edit    []string
	New     []string

	// Filter
	Filter       []string
	FilterClear  []string

	// Groups
	ToggleGroup  []string
	ExpandAll    []string
	CollapseAll  []string

	// System
	Help  []string
	Quit  []string
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		Up:       []string{"up", "k"},
		Down:     []string{"down", "j"},
		Left:     []string{"left", "h"},
		Right:    []string{"right", "l"},
		PageUp:   []string{"pgup", "ctrl+u"},
		PageDown: []string{"pgdown", "ctrl+d"},
		Home:     []string{"home", "g"},
		End:      []string{"end", "G"},

		ToggleSelect: []string{" ", "enter"},
		SelectAll:    []string{"ctrl+a"},
		DeselectAll:  []string{"ctrl+shift+a"},
		InvertSelect: []string{"ctrl+i"},

		Confirm: []string{"enter"},
		Cancel:  []string{"esc"},
		Delete:  []string{"d", "delete"},
		Edit:    []string{"e"},
		New:     []string{"n", "ctrl+n"},

		Filter:      []string{"/"},
		FilterClear: []string{"esc"},

		ToggleGroup:  []string{"enter", " "},
		ExpandAll:    []string{"ctrl+e"},
		CollapseAll:  []string{"ctrl+w"},

		Help: []string{"?", "ctrl+g"},
		Quit: []string{"q", "ctrl+c"},
	}
}

// MatchAction returns the action for a given key string.
func (km *KeyMap) MatchAction(key string) Action {
	// Navigation
	for _, k := range km.Up {
		if k == key {
			return ActionMoveUp
		}
	}
	for _, k := range km.Down {
		if k == key {
			return ActionMoveDown
		}
	}
	for _, k := range km.Left {
		if k == key {
			return ActionMoveLeft
		}
	}
	for _, k := range km.Right {
		if k == key {
			return ActionMoveRight
		}
	}
	for _, k := range km.PageUp {
		if k == key {
			return ActionPageUp
		}
	}
	for _, k := range km.PageDown {
		if k == key {
			return ActionPageDown
		}
	}
	for _, k := range km.Home {
		if k == key {
			return ActionMoveToTop
		}
	}
	for _, k := range km.End {
		if k == key {
			return ActionMoveToBottom
		}
	}

	// Selection
	for _, k := range km.ToggleSelect {
		if k == key {
			return ActionToggleSelection
		}
	}
	for _, k := range km.SelectAll {
		if k == key {
			return ActionSelectAll
		}
	}
	for _, k := range km.DeselectAll {
		if k == key {
			return ActionDeselectAll
		}
	}
	for _, k := range km.InvertSelect {
		if k == key {
			return ActionInvertSelection
		}
	}

	// Actions
	for _, k := range km.Confirm {
		if k == key {
			return ActionConfirm
		}
	}
	for _, k := range km.Cancel {
		if k == key {
			return ActionCancel
		}
	}
	for _, k := range km.Delete {
		if k == key {
			return ActionDelete
		}
	}
	for _, k := range km.Edit {
		if k == key {
			return ActionEdit
		}
	}
	for _, k := range km.New {
		if k == key {
			return ActionNew
		}
	}

	// Filter
	for _, k := range km.Filter {
		if k == key {
			return ActionFilter
		}
	}
	for _, k := range km.FilterClear {
		if k == key {
			return ActionFilterClear
		}
	}

	// Groups
	for _, k := range km.ToggleGroup {
		if k == key {
			return ActionToggleGroup
		}
	}
	for _, k := range km.ExpandAll {
		if k == key {
			return ActionExpandAll
		}
	}
	for _, k := range km.CollapseAll {
		if k == key {
			return ActionCollapseAll
		}
	}

	// System
	for _, k := range km.Help {
		if k == key {
			return ActionHelp
		}
	}
	for _, k := range km.Quit {
		if k == key {
			return ActionQuit
		}
	}

	return ActionNone
}

// BaseList contains the core state shared by all list types.
type BaseList struct {
	// Dimensions
	width  int
	height int

	// Focus state
	focused bool

	// Key bindings
	keyMap *KeyMap

	// State flags
	initialized bool
}

// NewBaseList creates a new base list.
func NewBaseList() *BaseList {
	return &BaseList{
		width:   80,
		height:  24,
		focused: true,
		keyMap:  DefaultKeyMap(),
	}
}

// Size returns the current dimensions.
func (l *BaseList) Size() (width, height int) {
	return l.width, l.height
}

// SetSize updates the dimensions.
func (l *BaseList) SetSize(width, height int) {
	l.width = width
	l.height = height
}

// Width returns the current width.
func (l *BaseList) Width() int {
	return l.width
}

// Height returns the current height.
func (l *BaseList) Height() int {
	return l.height
}

// Focused returns true if the list has focus.
func (l *BaseList) Focused() bool {
	return l.focused
}

// Focus sets the list to receive input.
func (l *BaseList) Focus() {
	l.focused = true
}

// Blur removes focus from the list.
func (l *BaseList) Blur() {
	l.focused = false
}

// SetKeyMap sets the key bindings.
func (l *BaseList) SetKeyMap(km *KeyMap) {
	l.keyMap = km
}

// KeyMap returns the current key bindings.
func (l *BaseList) KeyMap() *KeyMap {
	return l.keyMap
}

// MatchAction returns the action for a given key string.
func (l *BaseList) MatchAction(key string) Action {
	if l.keyMap == nil {
		return ActionNone
	}
	return l.keyMap.MatchAction(key)
}

// Initialized returns true if the list has been initialized.
func (l *BaseList) Initialized() bool {
	return l.initialized
}

// SetInitialized sets the initialization flag.
func (l *BaseList) SetInitialized(init bool) {
	l.initialized = init
}
