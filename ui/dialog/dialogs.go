package dialog

import (
	"fmt"
	"strings"

	"github.com/wwsheng009/taproot/ui/render"
)

const (
	DialogIDInfo     DialogID = "info"
	DialogIDConfirm  DialogID = "confirm"
)

// InfoDialog displays an informational message to the user.
type InfoDialog struct {
	id       DialogID
	title    string
	message  string
	width    int
	height   int
	callback Callback
	quitting bool
}

// NewInfoDialog creates a new info dialog.
func NewInfoDialog(title, message string) *InfoDialog {
	return &InfoDialog{
		id:       DialogIDInfo,
		title:    title,
		message:  message,
		width:    DefaultWidth(),
		height:   DefaultHeight(),
		callback: nil,
		quitting: false,
	}
}

// SetCallback sets a callback for when the dialog is closed.
func (d *InfoDialog) SetCallback(cb Callback) *InfoDialog {
	d.callback = cb
	return d
}

// Init implements render.Model.
func (d *InfoDialog) Init() error {
	return nil
}

// Update implements render.Model.
func (d *InfoDialog) Update(msg any) (render.Model, render.Cmd) {
	if d.quitting {
		return d, nil
	}

	switch msg := msg.(type) {
	case render.KeyMsg:
		switch msg.String() {
		case "enter", "escape", "q":
			d.quitting = true
			if d.callback != nil {
				d.callback(ActionConfirm, nil)
			}
			return d, render.Batch()
		}
	}

	return d, nil
}

// View implements render.Model.
func (d *InfoDialog) View() string {
	if d.quitting {
		return ""
	}

	var b strings.Builder

	// Box frame
	border := "╔════════════════════════════════════════════════════════════════╗"
	b.WriteString(border + "\n")

	// Title centered
	titleLine := d.centerText(d.title, 62)
	b.WriteString("║ " + titleLine + " ║\n")

	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	// Message with word wrap
	msgLines := d.wordWrap(d.message, 60)
	for _, line := range msgLines {
		msgLine := fmt.Sprintf("%-60s", line)
		b.WriteString("║ " + msgLine + " ║\n")
	}

	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	// OK button centered
	buttonLine := d.centerText("[ Enter or ESC to close ]", 62)
	b.WriteString("║ " + buttonLine + " ║\n")

	b.WriteString("╚════════════════════════════════════════════════════════════════╝")

	return b.String()
}

// ID returns the dialog ID.
func (d *InfoDialog) ID() DialogID {
	return d.id
}

// Title returns the dialog title.
func (d *InfoDialog) Title() string {
	return d.title
}

// SetTitle sets the dialog title.
func (d *InfoDialog) SetTitle(title string) {
	d.title = title
}

// Message returns the dialog message.
func (d *InfoDialog) Message() string {
	return d.message
}

// SetMessage sets the dialog message.
func (d *InfoDialog) SetMessage(message string) {
	d.message = message
}

// SetSize sets the dialog dimensions.
func (d *InfoDialog) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// SetID sets a custom dialog ID.
func (d *InfoDialog) SetID(id DialogID) {
	d.id = id
}

func (d *InfoDialog) centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}

func (d *InfoDialog) wordWrap(text string, width int) []string {
	var lines []string
	words := strings.Fields(text)
	currentLine := ""
	currentLen := 0

	for _, word := range words {
		wordLen := len(word)
		if currentLen == 0 {
			currentLine = word
			currentLen = wordLen
		} else if currentLen+1+wordLen <= width {
			currentLine += " " + word
			currentLen += 1 + wordLen
		} else {
			lines = append(lines, currentLine)
			currentLine = word
			currentLen = wordLen
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// ConfirmDialog presents a yes/no confirmation to the user.
type ConfirmDialog struct {
	id        DialogID
	title     string
	message   string
	width     int
	height    int
	callback  Callback
	selected  int // 0 = cancel, 1 = confirm
	quitting  bool
}

// NewConfirmDialog creates a new confirm dialog.
func NewConfirmDialog(title, message string, callback Callback) *ConfirmDialog {
	return &ConfirmDialog{
		id:        DialogIDConfirm,
		title:     title,
		message:   message,
		width:     DefaultWidth(),
		height:    DefaultHeight(),
		callback:  callback,
		selected:  0, // Default to cancel
		quitting:  false,
	}
}

// SetCallback sets the callback for dialog result.
func (d *ConfirmDialog) SetCallback(cb Callback) *ConfirmDialog {
	d.callback = cb
	return d
}

// Init implements render.Model.
func (d *ConfirmDialog) Init() error {
	return nil
}

// Update implements render.Model.
func (d *ConfirmDialog) Update(msg any) (render.Model, render.Cmd) {
	if d.quitting {
		return d, nil
	}

	switch msg := msg.(type) {
	case render.KeyMsg:
		switch msg.String() {
		case "left", "h":
			d.selected = 0
		case "right", "l":
			d.selected = 1
		case "enter":
			d.quitting = true
			result := ActionCancel
			if d.selected == 1 {
				result = ActionConfirm
			}
			if d.callback != nil {
				d.callback(result, nil)
			}
			return d, render.Batch()
		case "escape", "q":
			d.quitting = true
			if d.callback != nil {
				d.callback(ActionCancel, nil)
			}
			return d, render.Batch()
		}
	}

	return d, nil
}

// View implements render.Model.
func (d *ConfirmDialog) View() string {
	if d.quitting {
		return ""
	}

	var b strings.Builder

	border := "╔════════════════════════════════════════════════════════════════╗"
	b.WriteString(border + "\n")

	// Title
	titleLine := d.centerText(d.title, 62)
	b.WriteString("║ " + titleLine + " ║\n")

	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	// Message
	msgLines := d.wordWrap(d.message, 60)
	for _, line := range msgLines {
		msgLine := fmt.Sprintf("%-60s", line)
		b.WriteString("║ " + msgLine + " ║\n")
	}

	b.WriteString("╠════════════════════════════════════════════════════════════════╣\n")

	// Buttons
	cancelStyle := "     [ Cancel ]     "
	confirmStyle := "     [ Confirm ]     "
	
	if d.selected == 0 {
		cancelStyle = " >>>  [ Cancel ]  <<< "
	} else {
		confirmStyle = " >>>  [ Confirm ]  <<< "
	}

	buttonLine := d.centerText(cancelStyle+confirmStyle, 62)
	b.WriteString("║ " + buttonLine + " ║\n")

	b.WriteString("╚════════════════════════════════════════════════════════════════╝")

	return b.String()
}

// ID returns the dialog ID.
func (d *ConfirmDialog) ID() DialogID {
	return d.id
}

// Title returns the dialog title.
func (d *ConfirmDialog) Title() string {
	return d.title
}

// SetTitle sets the dialog title.
func (d *ConfirmDialog) SetTitle(title string) {
	d.title = title
}

// Message returns the dialog message.
func (d *ConfirmDialog) Message() string {
	return d.message
}

// SetMessage sets the dialog message.
func (d *ConfirmDialog) SetMessage(message string) {
	d.message = message
}

// SetSize sets the dialog dimensions.
func (d *ConfirmDialog) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// SetID sets a custom dialog ID.
func (d *ConfirmDialog) SetID(id DialogID) {
	d.id = id
}

// Selected returns the selected button index.
func (d *ConfirmDialog) Selected() int {
	return d.selected
}

// SetSelected sets the selected button index.
func (d *ConfirmDialog) SetSelected(index int) {
	d.selected = index
}

func (d *ConfirmDialog) centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}

func (d *ConfirmDialog) wordWrap(text string, width int) []string {
	var lines []string
	words := strings.Fields(text)
	currentLine := ""
	currentLen := 0

	for _, word := range words {
		wordLen := len(word)
		if currentLen == 0 {
			currentLine = word
			currentLen = wordLen
		} else if currentLen+1+wordLen <= width {
			currentLine += " " + word
			currentLen += 1 + wordLen
		} else {
			lines = append(lines, currentLine)
			currentLine = word
			currentLen = wordLen
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
