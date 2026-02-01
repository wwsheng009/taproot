package status

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

var _ Service = (*ServiceCmp)(nil)

// ServiceCmp is a component that displays the status of a single service.
type ServiceCmp struct {
	id         string
	name       string
	status     ServiceStatus
	errorCount int
	focused    bool
	compact    bool
	maxWidth   int
	initialized bool
}

// NewService creates a new service status component.
func NewService(id, name string) *ServiceCmp {
	return &ServiceCmp{
		id:         id,
		name:       name,
		status:     ServiceStatusOffline,
		errorCount: 0,
		focused:    false,
		compact:    false,
		maxWidth:   0,
		initialized: false,
	}
}

// Init initializes the component.
// Implements render.Model interface.
func (s *ServiceCmp) Init() render.Cmd {
	s.initialized = true
	return nil
}

// Update handles incoming messages and returns updated state and commands.
// Implements render.Model interface.
func (s *ServiceCmp) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		s.focused = true
	case *render.BlurMsg:
		s.focused = false
	}
	return s, render.None()
}

// View returns the string representation for rendering.
// Implements render.Model interface.
func (s *ServiceCmp) View() string {
	sty := styles.DefaultStyles()

	var statusIcon string
	var statusStyle lipgloss.Style

	switch s.status {
	case ServiceStatusOnline:
		statusIcon = "●"
		statusStyle = sty.ItemOnlineIcon
	case ServiceStatusOffline:
		statusIcon = "●"
		statusStyle = sty.ItemOfflineIcon
	case ServiceStatusBusy:
		statusIcon = "●"
		statusStyle = sty.ItemBusyIcon
	case ServiceStatusError:
		statusIcon = "●"
		statusStyle = sty.ItemErrorIcon
	case ServiceStatusStarting, ServiceStatusConnecting:
		statusIcon = "●"
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")) // Orange for transient states
	default:
		statusIcon = "●"
		statusStyle = sty.ItemOfflineIcon
	}

	iconStr := statusStyle.SetString(statusIcon).Render(" ")
	nameStr := sty.Base.Render(s.name)

	if s.compact {
		result := iconStr + nameStr
		if s.maxWidth > 0 && lipgloss.Width(result) > s.maxWidth {
			result = lipgloss.NewStyle().MaxWidth(s.maxWidth).Render(result)
		}
		return result
	}

	var result string
	result = iconStr + nameStr

	// Add error count if exists
	if s.errorCount > 0 {
		errorIcon := sty.LSP.ErrorDiagnostic.Foreground(sty.Error).Render(string(styles.ErrorIcon))
		errorCount := sty.LSP.ErrorDiagnostic.Render(fmt.Sprintf("%d", s.errorCount))
		result += " " + errorIcon + errorCount
	}

	// Add status label if focused
	if s.focused {
		statusColor := lipgloss.Color("#666666")
		if s.status == ServiceStatusError {
			statusColor = sty.Error
		} else if s.status == ServiceStatusOnline {
			statusColor = sty.GreenLight
		}
		statusLabel := lipgloss.NewStyle().Foreground(statusColor).Render(" " + s.status.String())
		result += statusLabel
	}

	if s.maxWidth > 0 && lipgloss.Width(result) > s.maxWidth {
		result = lipgloss.NewStyle().MaxWidth(s.maxWidth).Render(result)
	}

	return result
}

// ID returns the service identifier.
// Implements Service interface.
func (s *ServiceCmp) ID() string {
	return s.id
}

// Name returns the service name.
// Implements Service interface.
func (s *ServiceCmp) Name() string {
	return s.name
}

// Status returns the current service status.
// Implements Service interface.
func (s *ServiceCmp) Status() ServiceStatus {
	return s.status
}

// SetStatus sets the service status.
// Implements Service interface.
func (s *ServiceCmp) SetStatus(status ServiceStatus) {
	s.status = status
}

// IsOnline returns whether the service is online.
// Implements Service interface.
func (s *ServiceCmp) IsOnline() bool {
	return s.status == ServiceStatusOnline
}

// ErrorCount returns the error count.
// Implements Service interface.
func (s *ServiceCmp) ErrorCount() int {
	return s.errorCount
}

// SetErrorCount sets the error count.
// Implements Service interface.
func (s *ServiceCmp) SetErrorCount(count int) {
	s.errorCount = count
}

// Focus focuses the component.
func (s *ServiceCmp) Focus() {
	s.focused = true
}

// Blur unfocuses the component.
func (s *ServiceCmp) Blur() {
	s.focused = false
}

// SetCompact sets whether to use compact mode.
func (s *ServiceCmp) SetCompact(compact bool) {
	s.compact = compact
}

// Compact returns whether compact mode is enabled.
func (s *ServiceCmp) Compact() bool {
	return s.compact
}

// SetMaxWidth sets the maximum width for rendering.
func (s *ServiceCmp) SetMaxWidth(width int) {
	s.maxWidth = width
}

// MaxWidth returns the maximum width for rendering.
func (s *ServiceCmp) MaxWidth() int {
	return s.maxWidth
}

// Focused returns whether the component is focused.
func (s *ServiceCmp) Focused() bool {
	return s.focused
}
