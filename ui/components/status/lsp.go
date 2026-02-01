package status

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// LSPList represents a list of LSP services with their status and diagnostics.
type LSPList struct {
	services    []LSPServiceInfo
	maxItems    int
	width       int
	showTitle   bool
	title       string
	initialized bool
	cached      string
	cacheValid  bool
}

// LSPServiceInfo represents information about an LSP service.
type LSPServiceInfo struct {
	Name         string
	Language     string
	State        State
	Error        string
	Diagnostics  DiagnosticSummary
	ConnectedAt  string
}

// NewLSPList creates a new LSP list component.
func NewLSPList() *LSPList {
	return &LSPList{
		services:    []LSPServiceInfo{},
		maxItems:    5,
		width:       40,
		showTitle:   false,
		title:       "LSPs",
		initialized: false,
		cached:      "",
		cacheValid:  false,
	}
}

// Init initializes the component.
// Implements render.Model interface.
func (l *LSPList) Init() render.Cmd {
	l.initialized = true
	return nil
}

// Update handles incoming messages and returns updated state and commands.
// Implements render.Model interface.
func (l *LSPList) Update(msg any) (render.Model, render.Cmd) {
	l.cacheValid = false
	return l, render.None()
}

// View returns the string representation for rendering.
// Implements render.Model interface.
func (l *LSPList) View() string {
	if l.cacheValid && l.cached != "" {
		return l.cached
	}

	sty := styles.DefaultStyles()
	var result strings.Builder

	// Add title if enabled
	if l.showTitle {
		title := sty.Subtle.Render(l.title)
		result.WriteString(title)
		if len(l.services) > 0 {
			result.WriteString("\n\n")
		}
	}

	// Handle empty state
	if len(l.services) == 0 {
		result.WriteString(sty.Subtle.Render("None"))
		l.cached = result.String()
		l.cacheValid = true
		return l.cached
	}

	// Render services
	rendered := l.renderServices(sty)
	result.WriteString(rendered)

	l.cached = result.String()
	l.cacheValid = true
	return l.cached
}

// renderServices renders the LSP services with truncation if needed.
func (l *LSPList) renderServices(sty styles.Styles) string {
	var renderedServices []string

	for _, service := range l.services {
		renderedServices = append(renderedServices, l.renderService(sty, service))
	}

	// Handle truncation
	if len(renderedServices) > l.maxItems {
		visibleItems := renderedServices[:l.maxItems-1]
		remaining := len(renderedServices) - l.maxItems + 1
		moreText := sty.Subtle.Render(fmt.Sprintf("…and %d more", remaining))
		visibleItems = append(visibleItems, moreText)
		return lipgloss.JoinVertical(lipgloss.Left, visibleItems...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedServices...)
}

// renderService renders a single LSP service.
func (l *LSPList) renderService(sty styles.Styles, service LSPServiceInfo) string {
	var icon string
	var description string
	var diagnostics string

	switch service.State {
	case StateStarting:
		icon = sty.ItemBusyIcon.String()
		description = sty.Subtle.Render("starting...")
	case StateReady:
		icon = sty.ItemOnlineIcon.String()
		diagnostics = l.renderDiagnostics(sty, service.Diagnostics)
	case StateError:
		icon = sty.ItemErrorIcon.String()
		description = sty.Subtle.Render("error")
		if service.Error != "" {
			description = sty.Subtle.Render(fmt.Sprintf("error: %s", service.Error))
		}
	case StateDisabled:
		icon = sty.ItemOfflineIcon.String()
		description = sty.Subtle.Render("inactive")
	default:
		icon = sty.ItemOfflineIcon.String()
	}

	return renderStatus(sty, renderStatusOpts{
		Icon:         icon,
		Title:        service.Name,
		Description:  description,
		ExtraContent: diagnostics,
	}, l.width)
}

// renderDiagnostics formats diagnostic counts with icons and colors.
func (l *LSPList) renderDiagnostics(sty styles.Styles, diags DiagnosticSummary) string {
	var parts []string

	if diags.Error > 0 {
		parts = append(parts, sty.LSP.ErrorDiagnostic.Render(
			fmt.Sprintf("%s %d", styles.ErrorIcon, diags.Error)))
	}
	if diags.Warning > 0 {
		parts = append(parts, sty.LSP.WarningDiagnostic.Render(
			fmt.Sprintf("%s %d", styles.WarningIcon, diags.Warning)))
	}
	if diags.Hint > 0 {
		parts = append(parts, sty.LSP.HintDiagnostic.Render(
			fmt.Sprintf("%s %d", styles.HintIcon, diags.Hint)))
	}
	if diags.Information > 0 {
		parts = append(parts, sty.LSP.InfoDiagnostic.Render(
			fmt.Sprintf("%s %d", styles.InfoIcon, diags.Information)))
	}

	return strings.Join(parts, " ")
}

// SetServices sets the LSP services to display.
func (l *LSPList) SetServices(services []LSPServiceInfo) {
	l.services = services
	l.cacheValid = false
}

// Services returns the current LSP services.
func (l *LSPList) Services() []LSPServiceInfo {
	return l.services
}

// AddService adds an LSP service to the list.
func (l *LSPList) AddService(service LSPServiceInfo) {
	l.services = append(l.services, service)
	l.cacheValid = false
}

// ClearServices clears all LSP services.
func (l *LSPList) ClearServices() {
	l.services = []LSPServiceInfo{}
	l.cacheValid = false
}

// SetMaxItems sets the maximum number of services to display.
func (l *LSPList) SetMaxItems(max int) {
	l.maxItems = max
	l.cacheValid = false
}

// MaxItems returns the maximum number of services to display.
func (l *LSPList) MaxItems() int {
	return l.maxItems
}

// SetWidth sets the display width.
func (l *LSPList) SetWidth(width int) {
	l.width = width
	l.cacheValid = false
}

// Width returns the display width.
func (l *LSPList) Width() int {
	return l.width
}

// SetShowTitle sets whether to show the title.
func (l *LSPList) SetShowTitle(show bool) {
	l.showTitle = show
	l.cacheValid = false
}

// ShowTitle returns whether the title is shown.
func (l *LSPList) ShowTitle() bool {
	return l.showTitle
}

// SetTitle sets the title text.
func (l *LSPList) SetTitle(title string) {
	l.title = title
	l.cacheValid = false
}

// Title returns the title text.
func (l *LSPList) Title() string {
	return l.title
}

// HasErrors returns true if any LSP service has errors.
func (l *LSPList) HasErrors() bool {
	for _, s := range l.services {
		if s.State == StateError || s.Diagnostics.Error > 0 {
			return true
		}
	}
	return false
}

// TotalErrors returns the total number of errors across all LSP services.
func (l *LSPList) TotalErrors() int {
	total := 0
	for _, s := range l.services {
		total += s.Diagnostics.Error
	}
	return total
}

// TotalWarnings returns the total number of warnings across all LSP services.
func (l *LSPList) TotalWarnings() int {
	total := 0
	for _, s := range l.services {
		total += s.Diagnostics.Warning
	}
	return total
}

// TotalDiagnostics returns the total number of diagnostics across all LSP services.
func (l *LSPList) TotalDiagnostics() int {
	total := 0
	for _, s := range l.services {
		total += s.Diagnostics.Total()
	}
	return total
}

// OnlineCount returns the number of online LSP services.
func (l *LSPList) OnlineCount() int {
	count := 0
	for _, s := range l.services {
		if s.State == StateReady {
			count++
		}
	}
	return count
}

// renderStatusOpts contains options for rendering a status line.
type renderStatusOpts struct {
	Icon         string
	Title        string
	Description  string
	ExtraContent string
}

// renderStatus renders a status line with icon, title, description, and extra content.
func renderStatus(sty styles.Styles, opts renderStatusOpts, width int) string {
	content := []string{}

	if opts.Icon != "" {
		content = append(content, opts.Icon)
	}

	title := sty.Base.Render(opts.Title)
	content = append(content, title)

	if opts.Description != "" {
		content = append(content, opts.Description)
	}

	if opts.ExtraContent != "" {
		content = append(content, opts.ExtraContent)
	}

	result := strings.Join(content, " ")

	// Truncate if needed
	if lipgloss.Width(result) > width {
		result = lipgloss.NewStyle().MaxWidth(width).Render(result)
	}

	return result
}
