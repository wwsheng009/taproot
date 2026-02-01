package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"github.com/wwsheng009/taproot/ui/styles"
)

// FetchType represents the type of fetch operation.
type FetchType int

const (
	// FetchTypeBasic is a basic fetch operation.
	FetchTypeBasic FetchType = iota
	// FetchTypeWebFetch is a web fetch operation.
	FetchTypeWebFetch
	// FetchTypeWebSearch is a web search operation.
	FetchTypeWebSearch
	// FetchTypeAgentic is an agentic fetch with nested tools.
	FetchTypeAgentic
)

// String returns the string representation of fetch type.
func (f FetchType) String() string {
	switch f {
	case FetchTypeBasic:
		return "fetch"
	case FetchTypeWebFetch:
		return "web_fetch"
	case FetchTypeWebSearch:
		return "web_search"
	case FetchTypeAgentic:
		return "agentic_fetch"
	default:
		return "unknown"
	}
}

// FetchRequest represents a fetch request.
type FetchRequest struct {
	Type     FetchType
	URL      string
	Query    string      // For web_search
	Prompt   string      // For agentic_fetch
	Format   string      // text, markdown, html
	Timeout  time.Duration
	Params   map[string]any
	Children []*FetchRequest // Nested fetches for agentic_fetch
}

// FetchResult represents the result of a fetch operation.
type FetchResult struct {
	Content   string
	Error     string
	SavedPath string // For large content saved to file
	Size      int64  // Content size in bytes
	Duration  time.Duration
}

// FetchMessage represents a fetch/agentic fetch/web search operation message.
//
// This component supports:
// - Multiple fetch types (basic, web_fetch, web_search, agentic_fetch)
// - Request parameters display (URL, query, prompt, format, timeout)
// - Result display with syntax highlighting
// - Nested tool calls for agentic fetch (tree structure)
// - Collapsible results
// - Error display
// - Focus state with different styling
// - Caching for performance
type FetchMessage struct {
	id         string
	fetchType  FetchType
	request    FetchRequest
	result     *FetchResult
	timestamp  time.Time
	focused    bool
	expanded   bool
	loading    bool
	nested     []MessageItem // Nested messages for agentic_fetch
	config     *MessageConfig
	maxWidth   int

	// Caching fields
	cachedRender string
	cachedWidth  int
	cachedHeight int
	cacheValid   bool
}

// NewFetchMessage creates a new FetchMessage component.
func NewFetchMessage(id string, fetchType FetchType) *FetchMessage {
	return &FetchMessage{
		id:        id,
		fetchType: fetchType,
		timestamp: time.Now(),
		focused:   false,
		expanded:  false, // Start collapsed
		loading:   false,
		config:    &MessageConfig{},
	}
}

// NewBasicFetchMessage creates a new basic fetch message.
func NewBasicFetchMessage(id, url string) *FetchMessage {
	msg := NewFetchMessage(id, FetchTypeBasic)
	msg.request.URL = url
	msg.request.Format = "markdown"
	return msg
}

// NewWebFetchMessage creates a new web fetch message.
func NewWebFetchMessage(id, url string) *FetchMessage {
	msg := NewFetchMessage(id, FetchTypeWebFetch)
	msg.request.URL = url
	return msg
}

// NewWebSearchMessage creates a new web search message.
func NewWebSearchMessage(id, query string) *FetchMessage {
	msg := NewFetchMessage(id, FetchTypeWebSearch)
	msg.request.Query = query
	return msg
}

// NewAgenticFetchMessage creates a new agentic fetch message.
func NewAgenticFetchMessage(id, prompt string) *FetchMessage {
	msg := NewFetchMessage(id, FetchTypeAgentic)
	msg.request.Prompt = prompt
	return msg
}

// Init initializes the component. Implements render.Model.
func (m *FetchMessage) Init() render.Cmd {
	return nil
}

// Update handles incoming messages. Implements render.Model.
func (m *FetchMessage) Update(msg any) (render.Model, render.Cmd) {
	switch msg.(type) {
	case *render.FocusGainMsg:
		m.focused = true
	case *render.BlurMsg:
		m.focused = false
	}

	// Invalidate cache on any update
	m.cacheValid = false

	return m, render.None()
}

// View renders the fetch message. Implements render.Model.
func (m *FetchMessage) View() string {
	width := m.config.MaxWidth
	// Apply max width limit
	if m.maxWidth > 0 && width > m.maxWidth {
		width = m.maxWidth
	}

	return m.render(width)
}

// ID returns the message ID. Implements Identifiable interface.
func (m *FetchMessage) ID() string {
	return m.id
}

// Role returns RoleTool (fetches are tool messages). Implements Message interface.
func (m *FetchMessage) Role() MessageRole {
	return RoleTool
}

// Content returns a summary of the fetch. Implements Message interface.
func (m *FetchMessage) Content() string {
	switch m.fetchType {
	case FetchTypeWebSearch:
		return fmt.Sprintf("Search: %s", m.request.Query)
	case FetchTypeAgentic:
		return fmt.Sprintf("Agentic Fetch: %s", truncateString(m.request.Prompt, 50))
	default:
		return fmt.Sprintf("Fetch: %s", m.request.URL)
	}
}

// SetContent sets the prompt for agentic fetch. Implements Message interface.
func (m *FetchMessage) SetContent(content string) {
	m.request.Prompt = content
	m.cacheValid = false
}

// Timestamp returns the message timestamp. Implements Message interface.
func (m *FetchMessage) Timestamp() time.Time {
	return m.timestamp
}

// SetTimestamp sets the message timestamp.
func (m *FetchMessage) SetTimestamp(ts time.Time) {
	m.timestamp = ts
	m.cacheValid = false
}

// FetchType returns the fetch type.
func (m *FetchMessage) FetchType() FetchType {
	return m.fetchType
}

// Request returns the fetch request.
func (m *FetchMessage) Request() FetchRequest {
	return m.request
}

// SetRequest sets the fetch request.
func (m *FetchMessage) SetRequest(req FetchRequest) {
	m.request = req
	m.cacheValid = false
}

// Result returns the fetch result.
func (m *FetchMessage) Result() *FetchResult {
	return m.result
}

// SetResult sets the fetch result.
func (m *FetchMessage) SetResult(result *FetchResult) {
	m.result = result
	m.loading = false
	m.cacheValid = false
}

// SetURL sets the URL for the fetch.
func (m *FetchMessage) SetURL(url string) {
	m.request.URL = url
	m.cacheValid = false
}

// SetQuery sets the query for web search.
func (m *FetchMessage) SetQuery(query string) {
	m.request.Query = query
	m.cacheValid = false
}

// Prompt returns the prompt for agentic fetch.
func (m *FetchMessage) Prompt() string {
	return m.request.Prompt
}

// SetPrompt sets the prompt for agentic fetch.
func (m *FetchMessage) SetPrompt(prompt string) {
	m.request.Prompt = prompt
	m.cacheValid = false
}

// Format returns the format for the fetch.
func (m *FetchMessage) Format() string {
	return m.request.Format
}

// SetFormat sets the format for the fetch.
func (m *FetchMessage) SetFormat(format string) {
	m.request.Format = format
	m.cacheValid = false
}

// Timeout returns the timeout for the fetch.
func (m *FetchMessage) Timeout() time.Duration {
	return m.request.Timeout
}

// SetTimeout sets the timeout for the fetch.
func (m *FetchMessage) SetTimeout(timeout time.Duration) {
	m.request.Timeout = timeout
	m.cacheValid = false
}

// Loading returns true if the fetch is in progress.
func (m *FetchMessage) Loading() bool {
	return m.loading
}

// SetLoading sets the loading state.
func (m *FetchMessage) SetLoading(loading bool) {
	m.loading = loading
	m.cacheValid = false
}

// Nested returns the nested messages for agentic fetch.
func (m *FetchMessage) Nested() []MessageItem {
	return m.nested
}

// SetNested sets the nested messages for agentic fetch.
func (m *FetchMessage) SetNested(nested []MessageItem) {
	m.nested = nested
	m.cacheValid = false
}

// AddNested adds a nested message for agentic fetch.
func (m *FetchMessage) AddNested(item MessageItem) {
	m.nested = append(m.nested, item)
	m.cacheValid = false
}

// Expanded returns true if the message is expanded.
func (m *FetchMessage) Expanded() bool {
	return m.expanded
}

// SetExpanded sets the expansion state.
func (m *FetchMessage) SetExpanded(expanded bool) {
	m.expanded = expanded
	m.cacheValid = false
}

// ToggleExpanded toggles the expansion state.
func (m *FetchMessage) ToggleExpanded() {
	m.expanded = !m.expanded
	m.cacheValid = false
}

// Focus focuses the component. Implements Focusable interface.
func (m *FetchMessage) Focus() {
	m.focused = true
	m.cacheValid = false
}

// Blur blurs the component. Implements Focusable interface.
func (m *FetchMessage) Blur() {
	m.focused = false
	m.cacheValid = false
}

// Focused returns true if the component is focused. Implements Focusable interface.
func (m *FetchMessage) Focused() bool {
	return m.focused
}

// SetMaxWidth sets the maximum width for rendering.
func (m *FetchMessage) SetMaxWidth(width int) {
	m.maxWidth = width
	m.cacheValid = false
}

// SetConfig sets the message configuration.
func (m *FetchMessage) SetConfig(config *MessageConfig) {
	m.config = config
	m.cacheValid = false
}

// render renders the fetch message with the given width.
func (m *FetchMessage) render(width int) string {
	// Check cache first
	if m.cacheValid && m.cachedWidth == width && m.cachedRender != "" {
		return m.cachedRender
	}

	sty := styles.DefaultStyles()

	var builder strings.Builder

	// Render header
	header := m.renderHeader(&sty, width)
	builder.WriteString(header)

	// Render content if expanded
	if m.expanded {
		// Render request details
		details := m.renderRequestDetails(&sty, width)
		if details != "" {
			if header != "" {
				builder.WriteString("\n")
			}
			builder.WriteString(details)
		}

		// Render loading state
		if m.loading {
			loadingText := m.renderLoading(&sty)
			if loadingText != "" {
				builder.WriteString("\n")
				builder.WriteString(loadingText)
			}
		}

		// Render result if available
		if m.result != nil && !m.loading {
			result := m.renderResult(&sty, width)
			if result != "" {
				builder.WriteString("\n")
				builder.WriteString(result)
			}
		}

		// Render nested tools for agentic fetch
		if m.fetchType == FetchTypeAgentic && len(m.nested) > 0 {
			nested := m.renderNested(&sty, width)
			if nested != "" {
				builder.WriteString("\n")
				builder.WriteString(nested)
			}
		}
	}

	// Add expand hint if collapsed
	if !m.expanded && (m.result != nil || len(m.nested) > 0) {
		builder.WriteString("\n")
		expandHint := sty.Base.Foreground(sty.FgMuted).Italic(true).Render(
			"[press enter to expand]",
		)
		builder.WriteString(expandHint)
	}

	result := builder.String()

	// Apply focus styling
	if m.focused {
		result = lipgloss.NewStyle().Foreground(sty.Secondary).Render(result)
	} else {
		result = lipgloss.NewStyle().Foreground(sty.FgMuted).Render(result)
	}

	// Cache the result
	m.cachedRender = result
	m.cachedWidth = width
	m.cachedHeight = lipgloss.Height(result)
	m.cacheValid = true

	return result
}

// renderHeader renders the fetch message header.
func (m *FetchMessage) renderHeader(sty *styles.Styles, width int) string {
	// Determine icon and name based on fetch type
	icon, name := m.getFetchTypeDisplay()

	headerText := fmt.Sprintf("%s %s", icon, name)

	// Add key parameter
	switch m.fetchType {
	case FetchTypeBasic, FetchTypeWebFetch:
		if m.request.URL != "" {
			headerText += fmt.Sprintf(" %s", truncateString(m.request.URL, 40))
		}
	case FetchTypeWebSearch:
		if m.request.Query != "" {
			headerText += fmt.Sprintf(" %s", truncateString(m.request.Query, 40))
		}
	case FetchTypeAgentic:
		if m.request.Prompt != "" {
			headerText += fmt.Sprintf(" %s", truncateString(m.request.Prompt, 30))
		}
	}

	// Add format for basic fetch
	if m.fetchType == FetchTypeBasic && m.request.Format != "" {
		headerText += fmt.Sprintf(" (format: %s)", m.request.Format)
	}

	headerStyle := sty.Base.Bold(true).Foreground(sty.Primary)

	// Add status icon
	if m.result != nil {
		if m.result.Error != "" {
			headerStyle = headerStyle.Foreground(sty.Error)
		} else {
			headerStyle = headerStyle.Foreground(sty.Green) // Success green
		}
	} else if m.loading {
		headerStyle = headerStyle.Foreground(sty.Warning)
	}

	return headerStyle.Render(headerText)
}

// renderRequestDetails renders the request details.
func (m *FetchMessage) renderRequestDetails(sty *styles.Styles, width int) string {
	var details []string

	switch m.fetchType {
	case FetchTypeBasic, FetchTypeWebFetch:
		if m.request.URL != "" {
			urlStyle := sty.Base.Foreground(sty.FgMuted)
			details = append(details, urlStyle.Render("URL: "+m.request.URL))
		}
		if m.request.Timeout > 0 {
			timeoutStyle := sty.Base.Foreground(sty.FgMuted)
			details = append(details, timeoutStyle.Render("Timeout: "+m.request.Timeout.String()))
		}
	case FetchTypeWebSearch:
		if m.request.Query != "" {
			queryStyle := sty.Base.Foreground(sty.FgMuted)
			details = append(details, queryStyle.Render("Query: "+m.request.Query))
		}
	case FetchTypeAgentic:
		if m.request.Prompt != "" {
			promptLabel := sty.Base.Foreground(sty.Warning).Bold(true).Render("Prompt:")
			promptStyle := sty.Base.Foreground(sty.FgBase)
			// Wrap prompt to fit width
			wrappedPrompt := wrapText(m.request.Prompt, width-4)
			details = append(details, promptLabel)
			for _, line := range wrappedPrompt {
				details = append(details, promptStyle.Render("  "+line))
			}
		}
		if m.request.URL != "" {
			urlStyle := sty.Base.Foreground(sty.FgMuted)
			details = append(details, urlStyle.Render("URL: "+m.request.URL))
		}
	}

	if len(details) == 0 {
		return ""
	}

	return strings.Join(details, "\n")
}

// renderLoading renders the loading indicator.
func (m *FetchMessage) renderLoading(sty *styles.Styles) string {
	loadingStyle := sty.Base.Foreground(sty.Warning).Italic(true)
	return loadingStyle.Render("  Fetching...")
}

// renderResult renders the fetch result.
func (m *FetchMessage) renderResult(sty *styles.Styles, width int) string {
	if m.result == nil {
		return ""
	}

	// Check for error
	if m.result.Error != "" {
		errorTag := sty.Base.
			Foreground(sty.Error).
			Bold(true).
			Render("ERROR")
		errorText := sty.Base.Foreground(sty.Error).Width(width-4).Render(m.result.Error)
		return fmt.Sprintf("  %s\n  %s", errorTag, errorText)
	}

	// For large content saved to file
	if m.result.SavedPath != "" {
		infoStyle := sty.Base.Foreground(sty.Info).Italic(true)
		return infoStyle.Render(fmt.Sprintf("  ✓ Content saved to: %s (%d bytes)", m.result.SavedPath, m.result.Size))
	}

	// Render content
	content := m.result.Content
	if content == "" {
		return ""
	}

	// Add duration info
	var header string
	if m.result.Duration > 0 {
		durationStyle := sty.Base.Foreground(sty.FgMuted).Italic(true)
		header = durationStyle.Render(fmt.Sprintf("  ✓ Fetched in %s", m.result.Duration))
	}

	// Format content based on fetch type
	var contentBody string
	if m.fetchType == FetchTypeWebSearch || (m.fetchType == FetchTypeBasic && m.request.Format == "html") {
		// Render as markdown
		// Simplified markdown rendering
		contentBody = sty.Base.Foreground(sty.FgBase).Width(width-4).Render(content)
	} else {
		// Render as code block
		contentBody = renderCodeBlock(sty, content, m.request.Format, width-4)
	}

	if header != "" {
		return header + "\n" + contentBody
	}
	return contentBody
}

// renderNested renders nested tool calls for agentic fetch.
func (m *FetchMessage) renderNested(sty *styles.Styles, width int) string {
	if len(m.nested) == 0 {
		return ""
	}

	var lines []string

	// Header
	nestedLabel := sty.Base.Foreground(sty.Secondary).Bold(true).Render("Nested Tools:")
	lines = append(lines, nestedLabel)

	// Render each nested item
	for i, item := range m.nested {
		prefix := "  "
		if i < len(m.nested)-1 {
			prefix += "├─ "
		} else {
			prefix += "└─ "
		}

		// Get nested item view
		itemView := item.View()

		// Add prefix to each line
		itemLines := strings.Split(itemView, "\n")
		for j, line := range itemLines {
			if j == 0 {
				lines = append(lines, prefix+line)
			} else {
				// Continue lines with appropriate prefix
				if i < len(m.nested)-1 {
					lines = append(lines, "│  "+line)
				} else {
					lines = append(lines, "   "+line)
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}

// getFetchTypeDisplay returns the icon and name for the fetch type.
func (m *FetchMessage) getFetchTypeDisplay() (icon, name string) {
	switch m.fetchType {
	case FetchTypeBasic:
		return "🌐", "Fetch"
	case FetchTypeWebFetch:
		return "🌐", "Fetch"
	case FetchTypeWebSearch:
		return "🔍", "Search"
	case FetchTypeAgentic:
		return "🤖", "Agentic Fetch"
	default:
		return "❓", "Fetch"
	}
}

// HandleClick handles click events for expanding/collapsing.
// Implements ClickHandler interface.
func (m *FetchMessage) HandleClick(line, col int) render.Cmd {
	if m.result != nil || len(m.nested) > 0 {
		m.ToggleExpanded()
		return render.None()
	}
	return render.None()
}

// HandleKeyEvent handles keyboard events.
// Implements KeyEventHandler interface.
func (m *FetchMessage) HandleKeyEvent(msg any) (bool, render.Cmd) {
	keyMsg, ok := msg.(*render.KeyMsg)
	if !ok {
		return false, nil
	}

	switch keyMsg.String() {
	case " ", "enter":
		if m.result != nil || len(m.nested) > 0 {
			m.ToggleExpanded()
			return true, render.None()
		}
	}

	return false, nil
}

// Helper functions

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	// This is a simplified word wrap
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		testLine := currentLine.String()
		if testLine == "" {
			testLine = word
		} else {
			testLine = testLine + " " + word
		}

		if lipgloss.Width(testLine) <= width {
			currentLine.WriteString(word + " ")
		} else {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
			}
			currentLine.Reset()
			currentLine.WriteString(word + " ")
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

func renderCodeBlock(sty *styles.Styles, content, format string, width int) string {
	// Create code block style
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(sty.Border).
		Foreground(sty.FgBase)

	// In a real implementation, this would use syntax highlighting
	// For now, just wrap the content
	contentStyle := sty.Base.Width(width - 4).Render(content)
	return borderStyle.Render(contentStyle)
}
