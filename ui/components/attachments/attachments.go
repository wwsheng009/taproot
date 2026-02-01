package attachments

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// AttachmentList is a component for displaying a list of file attachments.
type AttachmentList struct {
	attachments []*Attachment
	config      AttachmentConfig
	styles      *styles.Styles
	focused     bool
	expanded    bool
	width       int

	// Render cache
	cached     string
	cacheValid bool
}

// NewAttachmentList creates a new AttachmentList component.
func NewAttachmentList(attachments []*Attachment) *AttachmentList {
	return &AttachmentList{
		attachments: attachments,
		width:       80,
		expanded:    true,
		focused:     false,
		config:      DefaultAttachmentConfig(),
		styles:      &styles.Styles{},
		cached:      "",
		cacheValid:  false,
	}
}

// Init initializes the component. Implements render.Model.
func (al *AttachmentList) Init() render.Cmd {
	return nil
}

// Update handles incoming messages. Implements render.Model.
func (al *AttachmentList) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		al.Focus()
	case *render.BlurMsg:
		al.Blur()
	}
	return al, nil
}

// View renders the component. Implements render.Model.
func (al *AttachmentList) View() string {
	if al.cacheValid && al.cached != "" {
		return al.cached
	}

	var b strings.Builder

	if !al.expanded {
		al.cached = b.String()
		al.cacheValid = true
		return al.cached
	}

	if len(al.attachments) == 0 {
		sty := al.styles
		b.WriteString(sty.Subtle.Render("No attachments"))
		al.cached = b.String()
		al.cacheValid = true
		return al.cached
	}

	for i, att := range al.attachments {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(al.renderAttachment(att, 0))
	}

	al.cached = b.String()
	al.cacheValid = true
	return al.cached
}

// renderAttachment renders a single attachment with optional indentation.
func (al *AttachmentList) renderAttachment(att *Attachment, indent int) string {
	sty := al.styles
	var b strings.Builder

	// Indentation
	prefix := strings.Repeat("  ", indent)

	// Get icon based on type
	icon := al.getAttachmentIcon(att.Type)
	iconStyle := al.getAttachmentIconStyle(att.Type)

	// Icon and name
	b.WriteString(prefix)
	b.WriteString(iconStyle.Render(icon))
	b.WriteString(" ")

	nameStyle := sty.Base
	if al.focused {
		nameStyle = nameStyle.Foreground(sty.Primary)
	}
	b.WriteString(nameStyle.Render(att.Name))

	// File size
	if al.config.ShowSize {
		b.WriteString(" ")
		b.WriteString(sty.Subtle.Render("(" + FormatSize(att.Size) + ")"))
	}

	// MIME type
	if !al.config.CompactMode {
		b.WriteString("\n")
		b.WriteString(prefix)
		b.WriteString("  ")
		b.WriteString(sty.Subtle.Render(att.MimeType))
	}

	// Preview for text files
	if al.config.ShowPreview && att.Preview != "" && !al.config.CompactMode {
		b.WriteString("\n")
		b.WriteString(prefix)
		b.WriteString("  ")
		preview := att.Preview
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		b.WriteString(sty.Muted.Render(preview))
	}

	// Attachment type label
	if !al.config.CompactMode {
		b.WriteString("\n")
		b.WriteString(prefix)
		b.WriteString("  ")
		b.WriteString(sty.Subtle.Render("Type: " + att.Type.String()))
	}

	return b.String()
}

// getAttachmentIcon returns the icon for an attachment type.
func (al *AttachmentList) getAttachmentIcon(attType AttachmentType) string {
	switch attType {
	case AttachmentTypeFile:
		return "📄"
	case AttachmentTypeImage:
		return "🖼️"
	case AttachmentTypeVideo:
		return "🎬"
	case AttachmentTypeAudio:
		return "🎵"
	case AttachmentTypeDocument:
		return "📝"
	case AttachmentTypeArchive:
		return "📦"
	default:
		return "•"
	}
}

// getAttachmentIconStyle returns the lipgloss Style for an attachment type.
func (al *AttachmentList) getAttachmentIconStyle(attType AttachmentType) lipgloss.Style {
	return al.styles.Base
}

// Focus focuses the component.
func (al *AttachmentList) Focus() {
	al.focused = true
	al.cacheValid = false
}

// Blur blurs the component.
func (al *AttachmentList) Blur() {
	al.focused = false
	al.cacheValid = false
}

// Focused returns true if the component is focused.
func (al *AttachmentList) Focused() bool {
	return al.focused
}

// SetWidth sets the width for rendering.
func (al *AttachmentList) SetWidth(width int) {
	al.width = width
	al.cacheValid = false
}

// SetConfig sets the attachment list configuration.
func (al *AttachmentList) SetConfig(config AttachmentConfig) {
	al.config = config
	al.cacheValid = false
}

// AddAttachment adds an attachment to the list.
func (al *AttachmentList) AddAttachment(attachment *Attachment) {
	al.attachments = append(al.attachments, attachment)
	al.cacheValid = false
}

// RemoveAttachment removes an attachment from the list by ID.
func (al *AttachmentList) RemoveAttachment(id string) bool {
	for i, att := range al.attachments {
		if att.ID == id {
			al.attachments = append(al.attachments[:i], al.attachments[i+1:]...)
			al.cacheValid = false
			return true
		}
	}
	return false
}

// GetAttachment retrieves an attachment by ID.
func (al *AttachmentList) GetAttachment(id string) *Attachment {
	for _, att := range al.attachments {
		if att.ID == id {
			return att
		}
	}
	return nil
}

// GetAttachments returns all attachments in the list.
func (al *AttachmentList) GetAttachments() []*Attachment {
	return al.attachments
}

// FilterByType filters attachments by type.
func (al *AttachmentList) FilterByType(attType AttachmentType) []*Attachment {
	var filtered []*Attachment
	for _, att := range al.attachments {
		if att.Type == attType {
			filtered = append(filtered, att)
		}
	}
	return filtered
}

// GetTotalSize returns the total size of all attachments.
func (al *AttachmentList) GetTotalSize() int64 {
	var total int64
	for _, att := range al.attachments {
		total += att.Size
	}
	return total
}

// GetCountByType returns the count of attachments by type.
func (al *AttachmentList) GetCountByType() map[AttachmentType]int {
	counts := make(map[AttachmentType]int)
	for _, att := range al.attachments {
		counts[att.Type]++
	}
	return counts
}
