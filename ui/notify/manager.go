package notify

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

type tickMsg struct {
	time time.Time
}

func tickCmd() render.Cmd {
	return func() render.Msg {
		time.Sleep(time.Millisecond * 100)
		return tickMsg{time: time.Now()}
	}
}

// Manager handles the display and lifecycle of notifications
type Manager struct {
	config        Config
	notifications []Notification
	styles        styles.Styles
	width         int
	height        int
}

// NewManager creates a new notification manager
func NewManager(config Config) *Manager {
	return &Manager{
		config: config,
		styles: styles.DefaultStyles(),
	}
}

// Init initializes the manager
func (m *Manager) Init() error {
	return nil
}

// Update handles messages
func (m *Manager) Update(msg any) (render.Model, render.Cmd) {
	var cmd render.Cmd

	switch msg := msg.(type) {
	case ShowNotificationMsg:
		// Add new notification
		n := msg.Notification
		if n.Duration == 0 {
			n.Duration = m.config.DefaultDuration
		}
		n.CreatedAt = time.Now()
		
		// Prepend or append based on position? usually new on top or bottom.
		// Let's prepend for now (stack grows from top)
		m.notifications = append([]Notification{n}, m.notifications...)
		
		// Trim if exceeding max
		if len(m.notifications) > m.config.MaxVisible {
			m.notifications = m.notifications[:m.config.MaxVisible]
		}
		
		// Ensure we are ticking
		return m, tickCmd()

	case tickMsg:
		now := msg.time
		active := make([]Notification, 0, len(m.notifications))
		hasActive := false
		
		for _, n := range m.notifications {
			if now.Sub(n.CreatedAt) < n.Duration {
				active = append(active, n)
				hasActive = true
			}
		}
		m.notifications = active
		
		if hasActive {
			return m, tickCmd()
		}

	case render.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

// View renders the notifications
func (m *Manager) View() string {
	if len(m.notifications) == 0 {
		return ""
	}

	var views []string
	for _, n := range m.notifications {
		views = append(views, m.renderNotification(n))
	}

	// Join them vertically
	content := lipgloss.JoinVertical(lipgloss.Top, views...)

	// Position the content
	// If width/height are set, we can position it.
	// Otherwise return the content itself.
	if m.width > 0 && m.height > 0 {
		// Map Config.Position to lipgloss positions
		var hPos, vPos lipgloss.Position
		
		switch m.config.Position {
		case TopRight:
			hPos, vPos = lipgloss.Right, lipgloss.Top
		case BottomRight:
			hPos, vPos = lipgloss.Right, lipgloss.Bottom
		case TopLeft:
			hPos, vPos = lipgloss.Left, lipgloss.Top
		case BottomLeft:
			hPos, vPos = lipgloss.Left, lipgloss.Bottom
		case TopCenter:
			hPos, vPos = lipgloss.Center, lipgloss.Top
		case BottomCenter:
			hPos, vPos = lipgloss.Center, lipgloss.Bottom
		default:
			hPos, vPos = lipgloss.Right, lipgloss.Top
		}

		return lipgloss.Place(m.width, m.height, hPos, vPos, content)
	}

	return content
}

func (m *Manager) renderNotification(n Notification) string {
	var (
		icon  string
		color lipgloss.Color
		titleStyle = m.styles.Base.Bold(true)
		msgStyle   = m.styles.Muted
	)

	switch n.Type {
	case TypeInfo:
		icon = styles.InfoIcon
		color = m.styles.Info
	case TypeSuccess:
		icon = styles.CheckIcon
		color = m.styles.Green
	case TypeWarning:
		icon = styles.WarningIcon
		color = m.styles.Warning
	case TypeError:
		icon = styles.ErrorIcon
		color = m.styles.Error
	default:
		icon = styles.InfoIcon
		color = m.styles.Info
	}

	// Apply color to icon and border
	icon = lipgloss.NewStyle().Foreground(color).SetString(icon).String()
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(0, 1).
		Width(30).
		MarginBottom(1) // Add margin to prevent overlap

	// Title
	title := titleStyle.Render(n.Title)
	
	// Layout
	// Icon Title
	// Message
	
	// Or:
	// Icon Title
	//      Message
	
	header := icon + " " + title
	
	content := header
	if n.Message != "" {
		content += "\n" + msgStyle.Render(n.Message)
	}

	return border.Render(content)
}
