package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwsheng009/taproot/ui/markdown"
	"github.com/wwsheng009/taproot/ui/styles"
)

type model struct {
	content    string
	rendered   string
	width      int
	showSource bool
	styles     *styles.Styles
}

func initialModel() *model {
	sty := styles.DefaultStyles()

	// Build comprehensive markdown content
	var md strings.Builder

	// Title and intro
	md.WriteString("# Markdown Enhancement Demo\n\n")
	md.WriteString("This demo showcases the **table**, **task list**, and **link** rendering capabilities.\n\n")

	// Tables section
	md.WriteString("## Tables\n\n")
	tb := markdown.NewTableBuilder([]string{"Feature", "Status", "Priority"})
	tb.SetAlignment(1, "center")
	tb.SetAlignment(2, "left")
	tb.AddRow([]string{"Table Rendering", "✅ Done", "High"})
	tb.AddRow([]string{"Task Lists", "✅ Done", "High"})
	tb.AddRow([]string{"Link Handling", "✅ Done", "Medium"})
	tb.AddRow([]string{"Helper Functions", "✅ Done", "Medium"})
	tb.AddRow([]string{"Code Blocks", "✅ Done", "High"})
	md.WriteString(tb.String())
	md.WriteString("\n")

	// Task list section
	md.WriteString("## Task Lists\n\n")
	tlb := markdown.NewTaskListBuilder()
	tlb.AddItem("Implement table rendering", true)
	tlb.AddItem("Implement task list rendering", true)
	tlb.AddItem("Add link handling", true)
	tlb.AddSubtask("Parse markdown links", false)
	tlb.AddItem("Create code block helpers", true)
	tlb.AddItem("Add comprehensive tests", true)
	tlb.AddItem("Create interactive demo", false)
	tlb.AddSubtask("Add keyboard controls", false)
	tlb.AddSubtask("Add responsive layout", false)
	md.WriteString(tlb.String())
	md.WriteString("\n")

	// Links section
	md.WriteString("## Links\n\n")
	lb := markdown.NewLinkBuilder("Taproot GitHub", "https://github.com/wwsheng009/taproot")
	lb.SetTitle("Taproot TUI Framework")
	md.WriteString("- " + lb.String() + "\n")
	lb2 := markdown.NewLinkBuilder("Bubbletea", "https://github.com/charmbracelet/bubbletea")
	md.WriteString("- " + lb2.String() + "\n")
	lb3 := markdown.NewLinkBuilder("Glamour", "https://github.com/charmbracelet/glamour")
	md.WriteString("- " + lb3.String() + "\n\n")

	// Code block section
	md.WriteString("## Code Blocks\n\n")
	cbb := markdown.NewCodeBlockBuilder("go", `// Render markdown with custom options
opts := markdown.DefaultRenderOptions()
opts.Width = 80
result, err := markdown.RenderWithStyles(content, sty, opts)`)
	md.WriteString(cbb.String())
	md.WriteString("\n\n")

	// Block quote section
	md.WriteString("## Block Quote\n\n")
	md.WriteString(markdown.BlockQuote("This is a block quote demonstrating\nhow text can be quoted across multiple lines\nwhile maintaining proper formatting."))
	md.WriteString("\n\n")

	// Horizontal rule
	md.WriteString("## Separator\n\n")
	md.WriteString(markdown.HorizontalRule())
	md.WriteString("\n\n")

	// Escaping section
	md.WriteString("## Escaping Special Characters\n\n")
	md.WriteString("To display literal special characters, escape them:\n\n")
	md.WriteString("- " + markdown.InlineCode("*not bold*") + "\n")
	md.WriteString("- " + markdown.InlineCode("_not italic_") + "\n")
	md.WriteString("- " + markdown.InlineCode("[not a link]") + "\n\n")

	return &model{
		content:    md.String(),
		rendered:   "",
		width:      80,
		showSource: false,
		styles:     &sty,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q":
			return m, tea.Quit
		case "s", "S":
			m.showSource = !m.showSource
		case "r", "R":
			m.renderMarkdown()
		case "+", "=":
			m.width += 10
			m.renderMarkdown()
		case "-", "_":
			m.width -= 10
			if m.width < 20 {
				m.width = 20
			}
			m.renderMarkdown()
		}
	}

	return m, nil
}

func (m *model) renderMarkdown() {
	opts := markdown.DefaultRenderOptions()
	opts.Width = m.width

	result, err := markdown.RenderWithStyles(m.content, m.styles, opts)
	if err != nil {
		m.rendered = fmt.Sprintf("Error: %v", err)
	} else {
		m.rendered = result
	}
}

func (m *model) View() string {
	sty := m.styles

	var b strings.Builder

	// Title
	title := sty.Base.Bold(true).Render("Markdown Enhancement Demo")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Help text
	help := sty.Subtle.Render("Controls: s=toggle source | r=render | +/- width | q=quit")
	b.WriteString(help)
	b.WriteString("\n\n")

	// Show either source or rendered
	if m.showSource {
		b.WriteString(sty.Subtle.Render("Source Markdown:"))
		b.WriteString("\n\n")
		b.WriteString(sty.Base.Render(m.content))
	} else {
		b.WriteString(sty.Subtle.Render("Rendered Output:"))
		b.WriteString("\n\n")
		if m.rendered == "" {
			m.renderMarkdown()
		}
		b.WriteString(sty.Base.Render(m.rendered))
	}

	// Width indicator
	widthInfo := sty.Subtle.Render(fmt.Sprintf("Width: %d", m.width))
	b.WriteString("\n\n")
	b.WriteString(widthInfo)

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
