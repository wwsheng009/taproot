package pills

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// PillStatus represents the status of a pill.
type PillStatus int

const (
	PillStatusPending PillStatus = iota
	PillStatusInProgress
	PillStatusCompleted
	PillStatusError
	PillStatusWarning
	PillStatusInfo
	PillStatusNeutral
)

// String returns the string representation of the pill status.
func (ps PillStatus) String() string {
	switch ps {
	case PillStatusPending:
		return "pending"
	case PillStatusInProgress:
		return "in-progress"
	case PillStatusCompleted:
		return "completed"
	case PillStatusError:
		return "error"
	case PillStatusWarning:
		return "warning"
	case PillStatusInfo:
		return "info"
	case PillStatusNeutral:
		return "neutral"
	default:
		return "unknown"
	}
}

// Pill represents a single pill with count and status.
type Pill struct {
	ID       string
	Label    string
	Count    int
	Status   PillStatus
	Expanded bool
	Items    []string // List of items (optional)
}

// PillConfig holds configuration for pill display.
type PillConfig struct {
	ShowItems    bool // Show item list in expanded state
	ShowCount    bool // Show count badge
	CompactMode  bool // Compact display mode
	MaxItemWidth int  // Maximum width for item display
	ShowIcons    bool // Show status icons
	InlineMode   bool // Display pills in a single line
}

// DefaultPillConfig returns default pill configuration.
func DefaultPillConfig() PillConfig {
	return PillConfig{
		ShowItems:    true,
		ShowCount:    true,
		CompactMode:  false,
		MaxItemWidth: 60,
		ShowIcons:    true,
		InlineMode:   false,
	}
}

// PillList is a component for displaying a list of pills.
type PillList struct {
	pills   []*Pill
	config  PillConfig
	styles  *styles.Styles
	focused bool
	width   int

	// Render cache
	cached     string
	cacheValid bool
}

// NewPillList creates a new PillList component.
func NewPillList(pills []*Pill) *PillList {
	return &PillList{
		pills:      pills,
		width:      80,
		focused:    false,
		config:     DefaultPillConfig(),
		styles:     &styles.Styles{},
		cached:     "",
		cacheValid: false,
	}
}

// Init initializes the component. Implements render.Model.
func (pl *PillList) Init() render.Cmd {
	return nil
}

// Update handles incoming messages. Implements render.Model.
func (pl *PillList) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		pl.Focus()
	case *render.BlurMsg:
		pl.Blur()
	}
	return pl, nil
}

// View renders the component. Implements render.Model.
func (pl *PillList) View() string {
	if pl.cacheValid && pl.cached != "" {
		return pl.cached
	}

	if len(pl.pills) == 0 {
		sty := pl.styles
		pl.cached = sty.Subtle.Render("No pills")
		pl.cacheValid = true
		return pl.cached
	}

	if pl.config.InlineMode {
		pl.cached = pl.renderInline()
	} else {
		var b strings.Builder
		for i, pill := range pl.pills {
			if i > 0 && !pl.config.InlineMode {
				b.WriteString("\n")
			}
			b.WriteString(pl.renderPill(pill, 0))
		}
		pl.cached = b.String()
	}

	pl.cacheValid = true
	return pl.cached
}

// renderInline renders pills in a single inline row.
func (pl *PillList) renderInline() string {
	pillParts := []string{}

	for _, pill := range pl.pills {
		pillStyle := pl.getPillStyle(pill.Status)
		icon := ""
		if pl.config.ShowIcons {
			icon = pl.getPillIcon(pill.Status) + " "
		}

		label := pill.Label
		if pl.config.ShowCount && pill.Count > 0 {
			countBadge := pillStyle.Render(fmt.Sprintf("[%d]", pill.Count))
			label = fmt.Sprintf("%s %s", label, countBadge)
		}

		rendered := pillStyle.Render(icon + label)
		pillParts = append(pillParts, rendered)
	}

	return strings.Join(pillParts, " • ")
}

// renderPill renders a single pill with optional indentation.
func (pl *PillList) renderPill(pill *Pill, indent int) string {
	sty := pl.styles
	var b strings.Builder

	prefix := strings.Repeat("  ", indent)
	pillStyle := pl.getPillStyle(pill.Status)

	// Icon and label
	icon := ""
	if pl.config.ShowIcons {
		icon = pl.getPillIcon(pill.Status) + " "
	}

	label := pill.Label
	if pl.config.ShowCount && pill.Count > 0 {
		label = fmt.Sprintf("%s [%d]", label, pill.Count)
	}

	b.WriteString(prefix)
	b.WriteString(pillStyle.Render(icon + label))

	// Show items if expanded
	if pill.Expanded && pl.config.ShowItems && len(pill.Items) > 0 {
		b.WriteString("\n")
		for i, item := range pill.Items {
			b.WriteString(prefix)
			b.WriteString("  ")

			// Item prefix
			itemPrefix := "• "
			if i < len(pill.Items)-1 {
				itemPrefix = "├─ "
			} else {
				itemPrefix = "└─ "
			}

			b.WriteString(sty.Subtle.Render(itemPrefix))

			// Item content (truncate if too long)
			itemText := item
			if pl.config.MaxItemWidth > 0 && len(itemText) > pl.config.MaxItemWidth {
				itemText = itemText[:pl.config.MaxItemWidth] + "..."
			}

			b.WriteString(sty.Base.Render(itemText))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// getPillIcon returns the icon for a pill status.
func (pl *PillList) getPillIcon(status PillStatus) string {
	switch status {
	case PillStatusPending:
		return "☐"
	case PillStatusInProgress:
		return "⟳"
	case PillStatusCompleted:
		return "✓"
	case PillStatusError:
		return "×"
	case PillStatusWarning:
		return "⚠"
	case PillStatusInfo:
		return "ℹ"
	case PillStatusNeutral:
		return "•"
	default:
		return "?"
	}
}

// getPillStyle returns the style for a pill status.
func (pl *PillList) getPillStyle(status PillStatus) lipgloss.Style {
	sty := pl.styles

	switch status {
	case PillStatusPending:
		return sty.Base.Foreground(sty.FgMuted)
	case PillStatusInProgress:
		return sty.Base.Foreground(sty.Secondary)
	case PillStatusCompleted:
		return sty.Base.Foreground(sty.Info)
	case PillStatusError:
		return sty.Base.Foreground(sty.Error)
	case PillStatusWarning:
		return sty.Base.Foreground(sty.Warning)
	case PillStatusInfo:
		return sty.Base.Foreground(sty.Primary)
	case PillStatusNeutral:
		return sty.Base.Foreground(sty.FgMuted)
	default:
		return sty.Base
	}
}

// Focus focuses the component.
func (pl *PillList) Focus() {
	pl.focused = true
	pl.cacheValid = false
}

// Blur blurs the component.
func (pl *PillList) Blur() {
	pl.focused = false
	pl.cacheValid = false
}

// Focused returns true if the component is focused.
func (pl *PillList) Focused() bool {
	return pl.focused
}

// SetWidth sets the width for rendering.
func (pl *PillList) SetWidth(width int) {
	pl.width = width
	pl.cacheValid = false
}

// SetConfig sets the pill list configuration.
func (pl *PillList) SetConfig(config PillConfig) {
	pl.config = config
	pl.cacheValid = false
}

// AddPill adds a pill to the list.
func (pl *PillList) AddPill(pill *Pill) {
	pl.pills = append(pl.pills, pill)
	pl.cacheValid = false
}

// RemovePill removes a pill from the list by ID.
func (pl *PillList) RemovePill(id string) bool {
	for i, pill := range pl.pills {
		if pill.ID == id {
			pl.pills = append(pl.pills[:i], pl.pills[i+1:]...)
			pl.cacheValid = false
			return true
		}
	}
	return false
}

// GetPill retrieves a pill by ID.
func (pl *PillList) GetPill(id string) *Pill {
	for _, pill := range pl.pills {
		if pill.ID == id {
			return pill
		}
	}
	return nil
}

// GetPills returns all pills in the list.
func (pl *PillList) GetPills() []*Pill {
	return pl.pills
}

// ToggleExpanded toggles the expanded state of a pill.
func (pl *PillList) ToggleExpanded(id string) bool {
	pill := pl.GetPill(id)
	if pill != nil {
		pill.Expanded = !pill.Expanded
		pl.cacheValid = false
		return true
	}
	return false
}

// GetTotalCount returns the total count across all pills.
func (pl *PillList) GetTotalCount() int {
	total := 0
	for _, pill := range pl.pills {
		total += pill.Count
	}
	return total
}

// GetCountByStatus returns the count of pills by status.
func (pl *PillList) GetCountByStatus() map[PillStatus]int {
	counts := make(map[PillStatus]int)
	for _, pill := range pl.pills {
		counts[pill.Status]++
	}
	return counts
}

// ExpandAll expands all pills.
func (pl *PillList) ExpandAll() {
	for _, pill := range pl.pills {
		pill.Expanded = true
	}
	pl.cacheValid = false
}

// CollapseAll collapses all pills.
func (pl *PillList) CollapseAll() {
	for _, pill := range pl.pills {
		pill.Expanded = false
	}
	pl.cacheValid = false
}
