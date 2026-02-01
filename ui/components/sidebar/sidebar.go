package sidebar

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/wwsheng009/taproot/ui/render"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// defaultStyles provides default sidebar styling.
type defaultStyles struct {
	base       lipgloss.Style
	title      lipgloss.Style
	muted      lipgloss.Style
	subtle     lipgloss.Style
	success    lipgloss.Style
	error      lipgloss.Style
	primary    lipgloss.Color
	secondary  lipgloss.Color
	warning    lipgloss.Color
}

func newDefaultStyles() defaultStyles {
	return defaultStyles{
		base:      lipgloss.NewStyle(),
		title:     lipgloss.NewStyle().Bold(true),
		muted:     lipgloss.NewStyle().Faint(true),
		subtle:    lipgloss.NewStyle().Foreground(lipgloss.Color("239")),
		success:   lipgloss.NewStyle().Foreground(lipgloss.Color("46")),
		error:     lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		primary:   lipgloss.Color("212"),
		secondary: lipgloss.Color("57"),
		warning:   lipgloss.Color("226"),
	}
}

// sidebarImpl implements the Sidebar interface.
type sidebarImpl struct {
	conf        Config
	styles      defaultStyles
	width       int
	height      int
	modelInfo   ModelInfo
	sessionInfo SessionInfo
	files       []FileInfo
	lsps        []LSPService
	mcps        []MCPService
}

// New creates a new sidebar component with the given configuration.
func New(conf Config) Sidebar {
	if conf.LogoProvider == nil {
		conf.LogoProvider = DefaultLogoProvider
	}
	return &sidebarImpl{
		conf:   conf,
		styles: newDefaultStyles(),
		files:  make([]FileInfo, 0),
		lsps:   make([]LSPService, 0),
		mcps:   make([]MCPService, 0),
	}
}

// NewDefault creates a sidebar with default configuration.
func NewDefault() Sidebar {
	return New(DefaultConfig())
}

// Init initializes the sidebar.
func (s *sidebarImpl) Init() render.Cmd {
	if s.conf.Width > 0 {
		s.width = s.conf.Width
	}
	if s.conf.Height > 0 {
		s.height = s.conf.Height
	}
	return nil
}

// Update handles messages and updates state.
func (s *sidebarImpl) Update(msg any) (Sidebar, any) {
	// For engine-agnostic design, Update just returns itself
	// External components are responsible for calling SetXXX methods
	return s, nil
}

// View renders the sidebar.
func (s *sidebarImpl) View() string {
	var parts []string

	// Build container style
	containerStyle := s.styles.base.
		Width(s.width).
		Height(s.height).
		Padding(1)
	if s.conf.CompactMode {
		containerStyle = containerStyle.PaddingTop(0)
	}

	// Add logo if not in compact mode
	if !s.conf.CompactMode && s.conf.ShowLogo {
		if s.height > s.conf.LogoHeight {
			parts = append(parts, s.renderLogo())
		} else {
			parts = append(parts, s.renderSmallLogo())
		}
	}

	// Add session title if present
	if s.sessionInfo.ID != "" {
		if s.conf.CompactMode {
			parts = append(parts, s.styles.base.Render(s.sessionInfo.Title), "")
		} else {
			parts = append(parts, s.styles.muted.Render(s.sessionInfo.Title), "")
		}
	}

	// Add working directory if not in compact mode
	if !s.conf.CompactMode && s.sessionInfo.WorkingDir != "" {
		parts = append(parts, s.renderWorkingDir(), "")
	}

	// Add model info
	parts = append(parts, s.renderModelInfo())

	// Add sections
	if s.conf.CompactMode && s.width > s.height {
		// Horizontal layout for compact mode when width > height
		sections := s.renderSectionsHorizontal()
		if sections != "" {
			parts = append(parts, "", sections)
		}
	} else {
		// Vertical layout
		if s.sessionInfo.ID != "" && len(s.files) > 0 {
			parts = append(parts, "", s.renderFilesSection())
		}
		parts = append(parts, "", s.renderLSPSection(), "", s.renderMCPSection())
	}

	return containerStyle.Render(lipgloss.JoinVertical(lipgloss.Left, parts...))
}

// renderLogo renders the full logo.
func (s *sidebarImpl) renderLogo() string {
	return s.conf.LogoProvider(s.width - 2)
}

// renderSmallLogo renders a simplified logo for small screens.
func (s *sidebarImpl) renderSmallLogo() string {
	return s.styles.title.Render("TR | TAPROOT")
}

// renderWorkingDir renders the working directory.
func (s *sidebarImpl) renderWorkingDir() string {
	return s.styles.muted.Render(s.sessionInfo.WorkingDir)
}

// renderModelInfo renders model information and token usage.
func (s *sidebarImpl) renderModelInfo() string {
	parts := []string{}

	// Model name with icon
	if s.modelInfo.Name != "" {
		icon := "M"
		if s.modelInfo.Icon != "" {
			icon = s.modelInfo.Icon
		}
		modelText := s.styles.base.Render(fmt.Sprintf("%s %s", icon, s.modelInfo.Name))
		parts = append(parts, modelText)
	}

	// Reasoning info
	if s.modelInfo.CanReason {
		reasoningStyle := s.styles.subtle.PaddingLeft(2)
		var reasoningText string
		switch s.modelInfo.Provider {
		case "anthropic":
			if s.modelInfo.ReasoningOn {
				reasoningText = "Thinking ON"
			} else {
				reasoningText = "Thinking OFF"
			}
		default:
			effort := "medium"
			if s.modelInfo.ReasoningEffort != "" {
				effort = s.modelInfo.ReasoningEffort
			}
			reasoningText = fmt.Sprintf("Reasoning %s", cases.Title(language.English).String(effort))
		}
		parts = append(parts, reasoningStyle.Render(reasoningText))
	}

	// Token usage and cost (only if session is active)
	if s.sessionInfo.ID != "" {
		tokenText := s.renderTokenUsage()
		if tokenText != "" {
			parts = append(parts, "  "+tokenText)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderTokenUsage renders token usage and cost.
func (s *sidebarImpl) renderTokenUsage() string {
	if s.modelInfo.ContextWindow == 0 {
		return ""
	}

	totalTokens := s.sessionInfo.PromptTokens + s.sessionInfo.CompletionTokens
	if totalTokens == 0 {
		return ""
	}

	// Format tokens
	var formattedTokens string
	switch {
	case totalTokens >= 1_000_000:
		formattedTokens = fmt.Sprintf("%.1fM", float64(totalTokens)/1_000_000)
	case totalTokens >= 1_000:
		formattedTokens = fmt.Sprintf("%.1fK", float64(totalTokens)/1_000)
	default:
		formattedTokens = strconv.FormatInt(totalTokens, 10)
	}

	// Remove .0 suffix
	formattedTokens = strings.Replace(formattedTokens, ".0K", "K", 1)
	formattedTokens = strings.Replace(formattedTokens, ".0M", "M", 1)

	// Calculate percentage
	percentage := float64(totalTokens) / float64(s.modelInfo.ContextWindow) * 100
	percentageText := s.styles.muted.Render(fmt.Sprintf("%d%%", int(percentage)))

	// Format cost
	costText := s.styles.muted.Render(fmt.Sprintf("$%.2f", s.sessionInfo.Cost))

	// Token count
	tokenCountText := s.styles.subtle.Render(fmt.Sprintf("(%s)", formattedTokens))

	// Build full text
	result := fmt.Sprintf("%s %s %s", percentageText, tokenCountText, costText)

	// Add warning if > 80%
	if percentage > 80 {
		result = fmt.Sprintf("⚠ %s", result)
	}

	return s.styles.base.Render(result)
}

// renderFilesSection renders the modified files section.
func (s *sidebarImpl) renderFilesSection() string {
	maxWidth := min(s.width-2, 60)

	maxFiles := s.conf.MaxFiles
	if maxFiles <= 0 {
		maxFiles = 10
	}

	return s.filesSection("Modified Files", s.files, maxWidth, maxFiles)
}

// renderLSPSection renders the LSP services section.
func (s *sidebarImpl) renderLSPSection() string {
	maxWidth := min(s.width-2, 60)

	maxLSPs := s.conf.MaxLSPs
	if maxLSPs <= 0 {
		maxLSPs = 8
	}

	return s.lspSection("LSPs", s.lsps, maxWidth, maxLSPs)
}

// renderMCPSection renders the MCP services section.
func (s *sidebarImpl) renderMCPSection() string {
	maxWidth := min(s.width-2, 60)

	maxMCPs := s.conf.MaxMCPs
	if maxMCPs <= 0 {
		maxMCPs = 8
	}

	return s.mcpSection("MCPs", s.mcps, maxWidth, maxMCPs)
}

// renderSectionsHorizontal renders all sections horizontally for compact mode.
func (s *sidebarImpl) renderSectionsHorizontal() string {
	totalWidth := s.width - 4
	sectionWidth := 50
	sectionWidth = min(sectionWidth, totalWidth/3)

	filesContent := s.filesSectionCompact("Modified Files", s.files, sectionWidth)
	lspContent := s.lspSectionCompact("LSPs", s.lsps, sectionWidth)
	mcpContent := s.mcpSectionCompact("MCPs", s.mcps, sectionWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top, filesContent, " ", lspContent, " ", mcpContent)
}

// filesSection renders the files section.
func (s *sidebarImpl) filesSection(title string, files []FileInfo, maxWidth, maxItems int) string {
	if len(files) == 0 {
		return ""
	}

	// Section header
	header := s.styles.title.Render(title)

	// Limit items
	maxItems = min(maxItems, len(files))

	var lines []string
	for i := 0; i < maxItems; i++ {
		file := files[i]
		path := file.Path
		if maxWidth > 0 && len(path) > maxWidth {
			path = path[:maxWidth-3] + "..."
		}

		// Add additions/deletions indicators
		var addons []string
		if file.Additions > 0 {
			addons = append(addons, s.styles.success.Render("+"+strconv.Itoa(file.Additions)))
		}
		if file.Deletions > 0 {
			addons = append(addons, s.styles.error.Render("-"+strconv.Itoa(file.Deletions)))
		}

		line := s.styles.base.Render(path)
		if len(addons) > 0 {
			line += " " + strings.Join(addons, " ")
		}
		lines = append(lines, line)
	}

	if len(files) > maxItems {
		remaining := len(files) - maxItems
		lines = append(lines, s.styles.muted.Render(fmt.Sprintf("...and %d more", remaining)))
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, strings.Join(lines, "\n"))
}

// filesSectionCompact renders files section for horizontal layout.
func (s *sidebarImpl) filesSectionCompact(title string, files []FileInfo, maxWidth int) string {
	maxItems := max(2, min(5, len(files)))

	// Reserve height for header and other content
	availableHeight := s.height - 8
	if availableHeight > 0 && availableHeight < maxItems {
		maxItems = availableHeight
	}

	return s.filesSection(title, files, maxWidth, maxItems)
}

// lspSection renders the LSP section.
func (s *sidebarImpl) lspSection(title string, lsps []LSPService, maxWidth, maxItems int) string {
	if len(lsps) == 0 {
		return ""
	}

	header := s.styles.title.Render(title)
	errorCount := 0
	for _, lsp := range lsps {
		errorCount += lsp.ErrorCount
	}

	var lines []string
	if errorCount > 0 {
		lines = append(lines, s.styles.error.Render(fmt.Sprintf("✗ %d errors", errorCount)))
	}

	maxItems = min(maxItems, len(lsps))

	for i := 0; i < maxItems; i++ {
		lsp := lsps[i]
		name := lsp.Name
		if maxWidth > 0 && len(name) > maxWidth {
			name = name[:maxWidth-3] + "..."
		}

		status := "✓"
		if !lsp.Connected {
			status = "✗"
		}

		line := fmt.Sprintf("%s %s", status, name)
		lines = append(lines, s.styles.base.Render(line))
	}

	if len(lsps) > maxItems {
		remaining := len(lsps) - maxItems
		lines = append(lines, s.styles.muted.Render(fmt.Sprintf("...and %d more", remaining)))
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, strings.Join(lines, "\n"))
}

// lspSectionCompact renders LSP section for horizontal layout.
func (s *sidebarImpl) lspSectionCompact(title string, lsps []LSPService, maxWidth int) string {
	maxItems := max(2, min(5, len(lsps)))
	availableHeight := s.height - 8
	if availableHeight > 0 && availableHeight < maxItems {
		maxItems = availableHeight
	}

	return s.lspSection(title, lsps, maxWidth, maxItems)
}

// mcpSection renders the MCP section.
func (s *sidebarImpl) mcpSection(title string, mcps []MCPService, maxWidth, maxItems int) string {
	if len(mcps) == 0 {
		return ""
	}

	header := s.styles.title.Render(title)

	maxItems = min(maxItems, len(mcps))

	var lines []string
	for i := 0; i < maxItems; i++ {
		mcp := mcps[i]
		name := mcp.Name
		if maxWidth > 0 && len(name) > maxWidth {
			name = name[:maxWidth-3] + "..."
		}

		status := "✓"
		if !mcp.Connected {
			status = "✗"
		}

		line := fmt.Sprintf("%s %s", status, name)
		lines = append(lines, s.styles.base.Render(line))
	}

	if len(mcps) > maxItems {
		remaining := len(mcps) - maxItems
		lines = append(lines, s.styles.muted.Render(fmt.Sprintf("...and %d more", remaining)))
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, strings.Join(lines, "\n"))
}

// mcpSectionCompact renders MCP section for horizontal layout.
func (s *sidebarImpl) mcpSectionCompact(title string, mcps []MCPService, maxWidth int) string {
	maxItems := max(2, min(5, len(mcps)))
	availableHeight := s.height - 8
	if availableHeight > 0 && availableHeight < maxItems {
		maxItems = availableHeight
	}

	return s.mcpSection(title, mcps, maxWidth, maxItems)
}

// Size returns the current dimensions.
func (s *sidebarImpl) Size() (int, int) {
	return s.width, s.height
}

// SetSize updates the dimensions.
func (s *sidebarImpl) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// SetModelInfo implements Sidebar.
func (s *sidebarImpl) SetModelInfo(info ModelInfo) {
	s.modelInfo = info
}

// SetSession implements Sidebar.
func (s *sidebarImpl) SetSession(info SessionInfo) {
	s.sessionInfo = info
}

// SetCompactMode implements Sidebar.
func (s *sidebarImpl) SetCompactMode(compact bool) {
	s.conf.CompactMode = compact
}

// AddFile implements Sidebar.
func (s *sidebarImpl) AddFile(file FileInfo) {
	s.files = append(s.files, file)
}

// ClearFiles implements Sidebar.
func (s *sidebarImpl) ClearFiles() {
	s.files = make([]FileInfo, 0)
}

// SetLSPStatus implements Sidebar.
func (s *sidebarImpl) SetLSPStatus(services []LSPService) {
	s.lsps = services
}

// SetMCPStatus implements Sidebar.
func (s *sidebarImpl) SetMCPStatus(services []MCPService) {
	s.mcps = services
}
