package status

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// MCPList represents a list of MCP services with their status and tool counts.
type MCPList struct {
	services    []MCPServiceInfo
	maxItems    int
	width       int
	showTitle   bool
	title       string
	initialized bool
	cached      string
	cacheValid  bool
}

// MCPServiceInfo represents information about an MCP service.
type MCPServiceInfo struct {
	Name        string
	State       State
	Error       string
	ToolCounts  ToolCounts
	ConnectedAt string
}

// NewMCPList creates a new MCP list component.
func NewMCPList() *MCPList {
	return &MCPList{
		services:    []MCPServiceInfo{},
		maxItems:    5,
		width:       40,
		showTitle:   false,
		title:       "MCPs",
		initialized: false,
		cached:      "",
		cacheValid:  false,
	}
}

// Init initializes the component.
// Implements render.Model interface.
func (m *MCPList) Init() render.Cmd {
	m.initialized = true
	return nil
}

// Update handles incoming messages and returns updated state and commands.
// Implements render.Model interface.
func (m *MCPList) Update(msg any) (render.Model, render.Cmd) {
	m.cacheValid = false
	return m, render.None()
}

// View returns the string representation for rendering.
// Implements render.Model interface.
func (m *MCPList) View() string {
	if m.cacheValid && m.cached != "" {
		return m.cached
	}

	sty := styles.DefaultStyles()
	var result strings.Builder

	// Add title if enabled
	if m.showTitle {
		title := sty.Subtle.Render(m.title)
		result.WriteString(title)
		if len(m.services) > 0 {
			result.WriteString("\n\n")
		}
	}

	// Handle empty state
	if len(m.services) == 0 {
		result.WriteString(sty.Subtle.Render("None"))
		m.cached = result.String()
		m.cacheValid = true
		return m.cached
	}

	// Render services
	rendered := m.renderServices(sty)
	result.WriteString(rendered)

	m.cached = result.String()
	m.cacheValid = true
	return m.cached
}

// renderServices renders the MCP services with truncation if needed.
func (m *MCPList) renderServices(sty styles.Styles) string {
	var renderedServices []string

	for _, service := range m.services {
		renderedServices = append(renderedServices, m.renderService(sty, service))
	}

	// Handle truncation
	if len(renderedServices) > m.maxItems {
		visibleItems := renderedServices[:m.maxItems-1]
		remaining := len(renderedServices) - m.maxItems + 1
		moreText := sty.Subtle.Render(fmt.Sprintf("…and %d more", remaining))
		visibleItems = append(visibleItems, moreText)
		return lipgloss.JoinVertical(lipgloss.Left, visibleItems...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedServices...)
}

// renderService renders a single MCP service.
func (m *MCPList) renderService(sty styles.Styles, service MCPServiceInfo) string {
	var icon string
	var description string
	var extraContent string

	switch service.State {
	case StateStarting:
		icon = sty.ItemBusyIcon.String()
		description = sty.Subtle.Render("starting...")
	case StateReady:
		icon = sty.ItemOnlineIcon.String()
		extraContent = m.renderToolCounts(sty, service.ToolCounts)
	case StateError:
		icon = sty.ItemErrorIcon.String()
		description = sty.Subtle.Render("error")
		if service.Error != "" {
			description = sty.Subtle.Render(fmt.Sprintf("error: %s", service.Error))
		}
	case StateDisabled:
		icon = sty.ItemOfflineIcon.String()
		description = sty.Subtle.Render("disabled")
	default:
		icon = sty.ItemOfflineIcon.String()
	}

	return renderStatus(sty, renderStatusOpts{
		Icon:         icon,
		Title:        service.Name,
		Description:  description,
		ExtraContent: extraContent,
	}, m.width)
}

// renderToolCounts formats tool and prompt counts for display.
func (m *MCPList) renderToolCounts(sty styles.Styles, counts ToolCounts) string {
	var parts []string

	if counts.Tools > 0 {
		label := "tools"
		if counts.Tools == 1 {
			label = "tool"
		}
		parts = append(parts, sty.Subtle.Render(fmt.Sprintf("%d %s", counts.Tools, label)))
	}

	if counts.Prompts > 0 {
		label := "prompts"
		if counts.Prompts == 1 {
			label = "prompt"
		}
		parts = append(parts, sty.Subtle.Render(fmt.Sprintf("%d %s", counts.Prompts, label)))
	}

	return strings.Join(parts, " ")
}

// SetServices sets the MCP services to display.
func (m *MCPList) SetServices(services []MCPServiceInfo) {
	m.services = services
	m.cacheValid = false
}

// Services returns the current MCP services.
func (m *MCPList) Services() []MCPServiceInfo {
	return m.services
}

// AddService adds an MCP service to the list.
func (m *MCPList) AddService(service MCPServiceInfo) {
	m.services = append(m.services, service)
	m.cacheValid = false
}

// ClearServices clears all MCP services.
func (m *MCPList) ClearServices() {
	m.services = []MCPServiceInfo{}
	m.cacheValid = false
}

// SetMaxItems sets the maximum number of services to display.
func (m *MCPList) SetMaxItems(max int) {
	m.maxItems = max
	m.cacheValid = false
}

// MaxItems returns the maximum number of services to display.
func (m *MCPList) MaxItems() int {
	return m.maxItems
}

// SetWidth sets the display width.
func (m *MCPList) SetWidth(width int) {
	m.width = width
	m.cacheValid = false
}

// Width returns the display width.
func (m *MCPList) Width() int {
	return m.width
}

// SetShowTitle sets whether to show the title.
func (m *MCPList) SetShowTitle(show bool) {
	m.showTitle = show
	m.cacheValid = false
}

// ShowTitle returns whether the title is shown.
func (m *MCPList) ShowTitle() bool {
	return m.showTitle
}

// SetTitle sets the title text.
func (m *MCPList) SetTitle(title string) {
	m.title = title
	m.cacheValid = false
}

// Title returns the title text.
func (m *MCPList) Title() string {
	return m.title
}

// HasErrors returns true if any MCP service has errors.
func (m *MCPList) HasErrors() bool {
	for _, s := range m.services {
		if s.State == StateError {
			return true
		}
	}
	return false
}

// ConnectedCount returns the number of connected MCP services.
func (m *MCPList) ConnectedCount() int {
	count := 0
	for _, s := range m.services {
		if s.State == StateReady {
			count++
		}
	}
	return count
}

// TotalTools returns the total number of tools across all MCP services.
func (m *MCPList) TotalTools() int {
	total := 0
	for _, s := range m.services {
		total += s.ToolCounts.Tools
	}
	return total
}

// TotalPrompts returns the total number of prompts across all MCP services.
func (m *MCPList) TotalPrompts() int {
	total := 0
	for _, s := range m.services {
		total += s.ToolCounts.Prompts
	}
	return total
}
