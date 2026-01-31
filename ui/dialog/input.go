package dialog

import (
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/list"
	"github.com/wwsheng009/taproot/ui/render"
)

const (
	DialogIDInput     DialogID = "input"
	DialogIDSelectList DialogID = "selectlist"
)

// InputDialog prompts the user for text input.
type InputDialog struct {
	id         DialogID
	title      string
	prompt     string
	input      *InputField
	width      int
	height     int
	callback   Callback
	quitting   bool
	hint       string
}

// NewInputDialog creates a new input dialog.
func NewInputDialog(title, prompt string, callback func(string)) *InputDialog {
	return &InputDialog{
		id:       DialogIDInput,
		title:    title,
		prompt:   prompt,
		input:    NewInputField(""),
		width:    60,
		height:   10,
		callback: func(result ActionResult, data any) {
			if result == ActionConfirm && data != nil {
				if val, ok := data.(string); ok && callback != nil {
					callback(val)
				}
			}
		},
		quitting: false,
		hint:     "",
	}
}

// SetHint sets a hint text below the input.
func (d *InputDialog) SetHint(hint string) *InputDialog {
	d.hint = hint
	return d
}

// SetPlaceholder sets the input placeholder.
func (d *InputDialog) SetPlaceholder(placeholder string) *InputDialog {
	d.input.SetPlaceholder(placeholder)
	return d
}

// SetHidden sets whether input is hidden (password).
func (d *InputDialog) SetHidden(hidden bool) *InputDialog {
	d.input.SetHidden(hidden)
	return d
}

// SetMaxLength sets the maximum input length.
func (d *InputDialog) SetMaxLength(max int) *InputDialog {
	d.input.SetMaxLength(max)
	return d
}

// SetID sets a custom dialog ID.
func (d *InputDialog) SetID(id DialogID) {
	d.id = id
}

// Init implements render.Model.
func (d *InputDialog) Init() error {
	d.input.SetFocused(true)
	return nil
}

// Update implements render.Model.
func (d *InputDialog) Update(msg any) (render.Model, render.Cmd) {
	if d.quitting {
		return d, nil
	}

	switch msg := msg.(type) {
	case render.KeyMsg:
		if !d.input.Focused() {
			switch msg.String() {
			case "enter":
				d.quitting = true
				if d.callback != nil {
					d.callback(ActionConfirm, d.input.Value())
				}
				return d, render.Batch()
			case "escape", "q":
				d.quitting = true
				if d.callback != nil {
					d.callback(ActionCancel, nil)
				}
				return d, render.Batch()
			}
			return d, nil
		}

		switch msg.String() {
		case "escape", "q":
			d.quitting = true
			if d.callback != nil {
				d.callback(ActionCancel, nil)
			}
			return d, render.Batch()
		case "enter":
			d.quitting = true
			if d.callback != nil {
				d.callback(ActionConfirm, d.input.Value())
			}
			return d, render.Batch()
		case "backspace", "ctrl+h":
			d.input.Delete()
		case "ctrl+u": // Clear line
			d.input.Clear()
		case "ctrl+a": // Start of line
			d.input.MoveCursorToStart()
		case "ctrl+e": // End of line
			d.input.MoveCursorToEnd()
		default:
			// Regular character input
			if len(msg.String()) == 1 {
				d.input.Insert([]rune(msg.String())[0])
			}
		}
	}

	return d, nil
}

// View implements render.Model.
func (d *InputDialog) View() string {
	if d.quitting {
		return ""
	}

	var b strings.Builder

	border := "╔════════════════════════════════════════════════════════════════╗"
	b.WriteString(border + "\n")

	// Title
	titleLine := centerText(d.title, 62)
	b.WriteString("║ " + titleLine + " ║\n")

	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	// Prompt
	promptLine := fmt.Sprintf("%-60s", d.prompt+":")
	b.WriteString("║ " + promptLine + " ║\n")

	// Input field
	inputValue := d.input.Value()
	if d.input.Hidden() {
		inputValue = strings.Repeat("•", len(inputValue))
	}
	if d.input.Focused() && d.input.Value() == "" && d.input.Placeholder() != "" {
		inputValue = d.input.Placeholder()
	}
	
	inputLine := fmt.Sprintf("  %s%-56s  ", ">", inputValue)
	b.WriteString("║ " + inputLine + " ║\n")

	// Hint
	if d.hint != "" {
		hintLine := fmt.Sprintf("  %s", d.hint)
		b.WriteString("║ " + fmt.Sprintf("%-60s", hintLine) + " ║\n")
	}

	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	// Footer
	footerLine := centerText("[ Enter to submit | ESC to cancel ]", 62)
	b.WriteString("║ " + footerLine + " ║\n")

	b.WriteString("╚════════════════════════════════════════════════════════════════╝")

	return b.String()
}

// ID returns the dialog ID.
func (d *InputDialog) ID() DialogID {
	return d.id
}

// Title returns the dialog title.
func (d *InputDialog) Title() string {
	return d.title
}

// SetTitle sets the dialog title.
func (d *InputDialog) SetTitle(title string) {
	d.title = title
}

// SetSize sets the dialog dimensions.
func (d *InputDialog) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// SelectListDialog allows selecting from a list of items.
type SelectListDialog struct {
	id        DialogID
	title     string
	selectList *SelectList
	width     int
	height    int
	callback  Callback
	quitting  bool
}

// NewSelectListDialog creates a new select list dialog.
func NewSelectListDialog(title string, items []string, callback func(int, string)) *SelectListDialog {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = list.NewListItem(fmt.Sprintf("%d", i), item, "")
	}

	return &SelectListDialog{
		id:        DialogIDSelectList,
		title:     title,
		selectList: NewSelectList(listItems),
		width:     60,
		height:    15,
		callback: func(result ActionResult, data any) {
			if result == ActionConfirm && data != nil {
				if idx, ok := data.(int); ok && callback != nil {
					item := items[idx]
					callback(idx, item)
				}
			}
		},
		quitting: false,
	}
}

// SetID sets a custom dialog ID.
func (d *SelectListDialog) SetID(id DialogID) {
	d.id = id
}

// Init implements render.Model.
func (d *SelectListDialog) Init() error {
	d.selectList.SetExpanded(true)
	d.selectList.SetFocused(true)
	return nil
}

// Update implements render.Model.
func (d *SelectListDialog) Update(msg any) (render.Model, render.Cmd) {
	if d.quitting {
		return d, nil
	}

	switch msg := msg.(type) {
	case render.KeyMsg:
		switch msg.String() {
		case "escape", "q":
			d.quitting = true
			if d.callback != nil {
				d.callback(ActionCancel, nil)
			}
			return d, render.Batch()
		case "enter":
			d.quitting = true
			if d.callback != nil {
				d.callback(ActionConfirm, d.selectList.Selected())
			}
			return d, render.Batch()
		case "up", "k":
			d.selectList.MoveUp()
		case "down", "j":
			d.selectList.MoveDown()
		case "home", "g":
			d.selectList.SetSelected(0)
		case "end", "G":
			d.selectList.SetSelected(len(d.selectList.Items()) - 1)
		}
	}

	return d, nil
}

// View implements render.Model.
func (d *SelectListDialog) View() string {
	if d.quitting {
		return ""
	}

	var b strings.Builder

	border := "╔════════════════════════════════════════════════════════════════╗"
	b.WriteString(border + "\n")

	// Title
	titleLine := centerText(d.title, 62)
	b.WriteString("║ " + titleLine + " ║\n")

	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	// List items
	items := d.selectList.Items()
	maxVisible := d.selectList.MaxVisible()
	start := max(0, d.selectList.Selected()-maxVisible+1)
	end := min(len(items), start+maxVisible)

	for i := start; i < end; i++ {
		item := items[i]
		cursor := " "
		if i == d.selectList.Selected() {
			cursor = "→"
		}

		if li, ok := item.(*list.ListItem); ok {
			line := fmt.Sprintf("%s %s", cursor, li.Title())
			b.WriteString("║ " + fmt.Sprintf("%-60s", line) + " ║\n")
		}
	}

	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	// Footer
	selectedItem := d.selectList.SelectedItem()
	var selectedText string
	if li, ok := selectedItem.(*list.ListItem); ok {
		selectedText = fmt.Sprintf("Selected: %s", li.Title())
	}
	footerLine := fmt.Sprintf("%-30s [ Enter to select | ESC to cancel ]", selectedText)
	b.WriteString("║ " + fmt.Sprintf("%-60s", footerLine) + " ║\n")

	b.WriteString("╚════════════════════════════════════════════════════════════════╝")

	return b.String()
}

// ID returns the dialog ID.
func (d *SelectListDialog) ID() DialogID {
	return d.id
}

// Title returns the dialog title.
func (d *SelectListDialog) Title() string {
	return d.title
}

// SetTitle sets the dialog title.
func (d *SelectListDialog) SetTitle(title string) {
	d.title = title
}

// SetSize sets the dialog dimensions.
func (d *SelectListDialog) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// SetItems sets the list items.
func (d *SelectListDialog) SetItems(items []string) {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = list.NewListItem(fmt.Sprintf("%d", i), item, "")
	}
	d.selectList.SetItems(listItems)
}

func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}
