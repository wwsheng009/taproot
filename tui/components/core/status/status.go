package status

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/wwsheng009/taproot/tui/util"
	"github.com/wwsheng009/taproot/ui/styles"
)

type StatusCmp interface {
	util.Model
	ToggleFullHelp() StatusCmp
	SetKeyMap(keyMap help.KeyMap) StatusCmp
}

type statusCmp struct {
	styles     *styles.Styles
	info       util.InfoMsg
	width      int
	messageTTL time.Duration
	help       help.Model
	keyMap     help.KeyMap
}

// clearMessageCmd is a command that clears status messages after a timeout
func (m *statusCmp) clearMessageCmd(ttl time.Duration) tea.Cmd {
	return tea.Tick(ttl, func(time.Time) tea.Msg {
		return util.ClearStatusMsg{}
	})
}

func (m *statusCmp) Init() tea.Cmd {
	return nil
}

func (m *statusCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	newModel := *m  // Deep copy
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newModel.width = msg.Width
		return &newModel, nil

	// Handle status info
	case util.InfoMsg:
		newModel.info = msg
		ttl := msg.TTL
		if ttl == 0 {
			ttl = newModel.messageTTL
		}
		return &newModel, newModel.clearMessageCmd(ttl)
	case util.ClearStatusMsg:
		newModel.info = util.InfoMsg{}
		return &newModel, nil
	}
	return &newModel, nil
}

func (m *statusCmp) View() string {
	s := m.styles
	// Only render help if keyMap is set
	if m.keyMap != nil {
		status := s.Base.Padding(0, 1, 1, 1).Render(m.help.View(m.keyMap))
		if m.info.Msg != "" {
			status = m.infoMsg()
		}
		return status
	}
	// Fallback: just show info message if available
	if m.info.Msg != "" {
		return m.infoMsg()
	}
	return ""
}

func (m *statusCmp) infoMsg() string {
	s := m.styles
	message := ""
	infoType := ""
	switch m.info.Type {
	case util.InfoTypeError:
		infoType = s.Base.Background(s.Error).Padding(0, 1).Render("ERROR")
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		message = s.Base.Background(s.Error).Width(widthLeft+2).Foreground(s.White).Padding(0, 1).Render(info)
	case util.InfoTypeWarn:
		infoType = s.Base.Foreground(s.BgOverlay).Background(s.Warning).Padding(0, 1).Render("WARNING")
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		message = s.Base.Foreground(s.BgOverlay).Width(widthLeft+2).Background(s.Warning).Padding(0, 1).Render(info)
	default:
		note := "OKAY!"
		if m.info.Type == util.InfoTypeUpdate {
			note = "HEY!"
		}
		infoType = s.Base.Foreground(s.BgSubtle).Background(s.Green).Padding(0, 1).Bold(true).Render(note)
		widthLeft := m.width - (lipgloss.Width(infoType) + 2)
		info := ansi.Truncate(m.info.Msg, widthLeft, "…")
		message = s.Base.Background(s.Green).Width(widthLeft+2).Foreground(s.BgSubtle).Padding(0, 1).Render(info)
	}
	return ansi.Truncate(infoType+message, m.width, "…")
}

func (m *statusCmp) ToggleFullHelp() StatusCmp {
	newModel := *m  // Deep copy
	newModel.help.ShowAll = !newModel.help.ShowAll
	return &newModel
}

func (m *statusCmp) SetKeyMap(keyMap help.KeyMap) StatusCmp {
	newModel := *m  // Deep copy
	newModel.keyMap = keyMap
	return &newModel
}

func NewStatusCmp() StatusCmp {
	s := styles.DefaultStyles()
	h := help.New()
	h.Styles = help.Styles{
		ShortKey:       s.Dialog.Help.ShortKey,
		ShortDesc:      s.Dialog.Help.ShortDesc,
		ShortSeparator: s.Dialog.Help.ShortSeparator,
		Ellipsis:       s.Dialog.Help.Ellipsis,
		FullKey:        s.Dialog.Help.FullKey,
		FullDesc:       s.Dialog.Help.FullDesc,
		FullSeparator:  s.Dialog.Help.FullSeparator,
	}
	return &statusCmp{
		styles:     &s,
		messageTTL: 5 * time.Second,
		help:       h,
	}
}
