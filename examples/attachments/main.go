package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/components/attachments"
)

const (
	maxWidth  = 80
	maxHeight = 24
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7c3aed")).
			Bold(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280")).
			Italic(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d1d5db"))
)

type model struct {
	attachmentList *attachments.AttachmentList
	width          int
	height         int
	cursor         int
	config         attachments.AttachmentConfig
}

func initialModel() model {
	// Create sample attachments
	sampleAttachments := []*attachments.Attachment{
		{
			ID:       "1",
			Name:     "report.pdf",
			Type:     attachments.AttachmentTypeDocument,
			Size:     1024 * 150,
			MimeType: "application/pdf",
			Preview:  "Monthly sales report with charts and graphs...",
		},
		{
			ID:       "2",
			Name:     "photo.jpg",
			Type:     attachments.AttachmentTypeImage,
			Size:     1024 * 2500,
			MimeType: "image/jpeg",
			Preview:  "Photograph of the team event...",
		},
		{
			ID:       "3",
			Name:     "code.go",
			Type:     attachments.AttachmentTypeDocument,
			Size:     1024 * 5,
			MimeType: "text/plain",
			Preview:  "package main\n\nimport \"fmt\"\n\nfunc main() { ... }",
		},
		{
			ID:       "4",
			Name:     "backup.zip",
			Type:     attachments.AttachmentTypeArchive,
			Size:     1024 * 1024 * 50,
			MimeType: "application/zip",
			Preview:  "Compressed backup of project files...",
		},
		{
			ID:       "5",
			Name:     "tutorial.mp4",
			Type:     attachments.AttachmentTypeVideo,
			Size:     1024 * 1024 * 100,
			MimeType: "video/mp4",
			Preview:  "Screencast tutorial on using the system...",
		},
		{
			ID:       "6",
			Name:     "podcast.mp3",
			Type:     attachments.AttachmentTypeAudio,
			Size:     1024 * 1024 * 5,
			MimeType: "audio/mpeg",
			Preview:  "Episode 45: Advanced Techniques...",
		},
		{
			ID:       "7",
			Name:     "instructions.txt",
			Type:     attachments.AttachmentTypeDocument,
			Size:     1024 * 2,
			MimeType: "text/plain",
			Preview:  "Step 1: Download the installer\nStep 2: Run the setup wizard...",
		},
	}

	al := attachments.NewAttachmentList(sampleAttachments)
	config := attachments.DefaultAttachmentConfig()

	return model{
		attachmentList: al,
		width:          maxWidth,
		height:         maxHeight,
		cursor:         0,
		config:         config,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		switch key {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.attachmentList.GetAttachments())-1 {
				m.cursor++
			}

		case "a":
			// Add new attachment
			newAtt := &attachments.Attachment{
				ID:       fmt.Sprintf("%d", len(m.attachmentList.GetAttachments())+1),
				Name:     fmt.Sprintf("newfile%d.txt", len(m.attachmentList.GetAttachments())),
				Type:     attachments.AttachmentTypeDocument,
				Size:     1024 * 10,
				MimeType: "text/plain",
				Preview:  "This is a newly added file...",
			}
			m.attachmentList.AddAttachment(newAtt)

		case "r":
			// Remove selected attachment
			if len(m.attachmentList.GetAttachments()) > 0 {
				att := m.attachmentList.GetAttachments()[m.cursor]
				m.attachmentList.RemoveAttachment(att.ID)
				if m.cursor >= len(m.attachmentList.GetAttachments()) {
					m.cursor = len(m.attachmentList.GetAttachments()) - 1
				}
			}

		case "c":
			// Toggle compact mode
			m.config.CompactMode = !m.config.CompactMode
			m.attachmentList.SetConfig(m.config)

		case "p":
			// Toggle preview
			m.config.ShowPreview = !m.config.ShowPreview
			m.attachmentList.SetConfig(m.config)

		case "s":
			// Toggle size display
			m.config.ShowSize = !m.config.ShowSize
			m.attachmentList.SetConfig(m.config)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width > maxWidth {
			m.width = maxWidth
		}
		m.attachmentList.SetWidth(m.width - 4)
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Attachment List Demo"))
	b.WriteString("\n\n")

	// Attachment list view
	b.WriteString(m.attachmentList.View())
	b.WriteString("\n\n")

	// Selected attachment details
	if len(m.attachmentList.GetAttachments()) > 0 {
		att := m.attachmentList.GetAttachments()[m.cursor]

		// Separator
		b.WriteString(separatorStyle.Render(strings.Repeat("─", m.width)))
		b.WriteString("\n\n")

		// Details
		details := fmt.Sprintf(
			"Selected: %s\nType: %s | Size: %s | MIME: %s\nPreview: %s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#7c3aed")).Render(att.Name),
			att.Type.String(),
			attachments.FormatSize(att.Size),
			att.MimeType,
			att.Preview,
		)
		b.WriteString(details)
	}

	// Separator
	b.WriteString("\n\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", m.width)))
	b.WriteString("\n")

	// Help text
	help := "Controls: ↑/k • ↓/j • aadd • rremove • ccompact • ppreview • ssize • qquit"
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
