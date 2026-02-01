package markdown

import (
	"strings"
	"testing"

	"github.com/wwsheng009/taproot/ui/styles"
)

func TestRender(t *testing.T) {
	sty := styles.DefaultStyles()

	tests := []struct {
		name     string
		content  string
		plain    bool
		wantErr  bool
	}{
		{
			name:    "simple text",
			content: "Hello, world!",
			plain:   false,
			wantErr: false,
		},
		{
			name:    "bold text",
			content: "**Bold** text",
			plain:   false,
			wantErr: false,
		},
		{
			name:    "code block",
			content: "```go\nfunc main() {}\n```",
			plain:   false,
			wantErr: false,
		},
		{
			name:    "plain mode",
			content: "# Header\n\nContent",
			plain:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultRenderOptions()
			opts.Plain = tt.plain

			result, err := RenderWithStyles(tt.content, &sty, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == "" {
				t.Error("Render() returned empty result")
			}
		})
	}
}

func TestTableBuilder(t *testing.T) {
	t.Run("NewTableBuilder", func(t *testing.T) {
		tb := NewTableBuilder([]string{"Name", "Age", "City"})
		if tb == nil {
			t.Error("NewTableBuilder() returned nil")
		}
		if len(tb.headers) != 3 {
			t.Errorf("Expected 3 headers, got %d", len(tb.headers))
		}
	})

	t.Run("SetAlignment", func(t *testing.T) {
		tb := NewTableBuilder([]string{"A", "B", "C"})
		tb.SetAlignment(0, "left")
		tb.SetAlignment(1, "center")
		tb.SetAlignment(2, "right")

		if tb.alignment[0] != "left" {
			t.Errorf("Expected alignment[0] = 'left', got %s", tb.alignment[0])
		}
		if tb.alignment[1] != "center" {
			t.Errorf("Expected alignment[1] = 'center', got %s", tb.alignment[1])
		}
		if tb.alignment[2] != "right" {
			t.Errorf("Expected alignment[2] = 'right', got %s", tb.alignment[2])
		}
	})

	t.Run("AddRow", func(t *testing.T) {
		tb := NewTableBuilder([]string{"Name", "Age"})
		tb.AddRow([]string{"Alice", "30"})

		if len(tb.rows) != 1 {
			t.Errorf("Expected 1 row, got %d", len(tb.rows))
		}
	})

	t.Run("AddRow wrong length", func(t *testing.T) {
		tb := NewTableBuilder([]string{"Name", "Age"})
		tb.AddRow([]string{"Alice"})

		// Should not add rows with wrong length
		if len(tb.rows) != 0 {
			t.Errorf("Expected 0 rows, got %d", len(tb.rows))
		}
	})

	t.Run("String", func(t *testing.T) {
		tb := NewTableBuilder([]string{"Name", "Age"})
		tb.AddRow([]string{"Alice", "30"})
		tb.AddRow([]string{"Bob", "25"})

		result := tb.String()
		expected := "| Name | Age |\n| --- | --- |\n| Alice | 30 |\n| Bob | 25 |\n"
		if result != expected {
			t.Errorf("TableBuilder.String() mismatch\nGot: %q\nWant: %q", result, expected)
		}
	})
}

func TestTaskListBuilder(t *testing.T) {
	t.Run("NewTaskListBuilder", func(t *testing.T) {
		tlb := NewTaskListBuilder()
		if tlb == nil {
			t.Error("NewTaskListBuilder() returned nil")
		}
	})

	t.Run("AddItem", func(t *testing.T) {
		tlb := NewTaskListBuilder()
		tlb.AddItem("Task 1", true)
		tlb.AddItem("Task 2", false)

		if len(tlb.items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(tlb.items))
		}
		if tlb.items[0].Checked != true {
			t.Error("First item should be checked")
		}
		if tlb.items[1].Checked != false {
			t.Error("Second item should be unchecked")
		}
	})

	t.Run("AddSubtask", func(t *testing.T) {
		tlb := NewTaskListBuilder()
		tlb.AddItem("Parent task", false)
		tlb.AddSubtask("Subtask 1", true)
		tlb.AddSubtask("Subtask 2", false)

		if len(tlb.items[0].Subtasks) != 2 {
			t.Errorf("Expected 2 subtasks, got %d", len(tlb.items[0].Subtasks))
		}
	})

	t.Run("String", func(t *testing.T) {
		tlb := NewTaskListBuilder()
		tlb.AddItem("Task 1", true)
		tlb.AddItem("Task 2", false)

		result := tlb.String()
		expected := "- [x] Task 1\n- [ ] Task 2\n"
		if result != expected {
			t.Errorf("TaskListBuilder.String() mismatch\nGot: %q\nWant: %q", result, expected)
		}
	})
}

func TestLinkBuilder(t *testing.T) {
	t.Run("NewLinkBuilder", func(t *testing.T) {
		lb := NewLinkBuilder("Google", "https://google.com")
		if lb == nil {
			t.Error("NewLinkBuilder() returned nil")
		}
	})

	t.Run("String", func(t *testing.T) {
		lb := NewLinkBuilder("Google", "https://google.com")
		result := lb.String()
		expected := "[Google](https://google.com)"
		if result != expected {
			t.Errorf("LinkBuilder.String() = %q, want %q", result, expected)
		}
	})

	t.Run("SetTitle", func(t *testing.T) {
		lb := NewLinkBuilder("Google", "https://google.com")
		lb.SetTitle("Search Engine")
		result := lb.String()
		expected := "[Google](https://google.com \"Search Engine\")"
		if result != expected {
			t.Errorf("LinkBuilder.String() with title = %q, want %q", result, expected)
		}
	})
}

func TestImageBuilder(t *testing.T) {
	t.Run("NewImageBuilder", func(t *testing.T) {
		ib := NewImageBuilder("Logo", "logo.png")
		if ib == nil {
			t.Error("NewImageBuilder() returned nil")
		}
	})

	t.Run("String", func(t *testing.T) {
		ib := NewImageBuilder("Logo", "logo.png")
		result := ib.String()
		expected := "![Logo](logo.png)"
		if result != expected {
			t.Errorf("ImageBuilder.String() = %q, want %q", result, expected)
		}
	})

	t.Run("SetTitle", func(t *testing.T) {
		ib := NewImageBuilder("Logo", "logo.png")
		ib.SetTitle("Company Logo")
		result := ib.String()
		expected := "![Logo](logo.png \"Company Logo\")"
		if result != expected {
			t.Errorf("ImageBuilder.String() with title = %q, want %q", result, expected)
		}
	})
}

func TestCodeBlockBuilder(t *testing.T) {
	t.Run("NewCodeBlockBuilder", func(t *testing.T) {
		cbb := NewCodeBlockBuilder("go", "func main() {}")
		if cbb == nil {
			t.Error("NewCodeBlockBuilder() returned nil")
		}
	})

	t.Run("String", func(t *testing.T) {
		cbb := NewCodeBlockBuilder("go", "func main() {}")
		result := cbb.String()
		expected := "```go\nfunc main() {}\n```"
		if result != expected {
			t.Errorf("CodeBlockBuilder.String() = %q, want %q", result, expected)
		}
	})
}

func TestInlineCode(t *testing.T) {
	result := InlineCode("code")
	expected := "`code`"
	if result != expected {
		t.Errorf("InlineCode() = %q, want %q", result, expected)
	}

	// Test escaping
	result2 := InlineCode("code with `backtick`")
	expected2 := "`code with \\`backtick\\``"
	if result2 != expected2 {
		t.Errorf("InlineCode() with backtick = %q, want %q", result2, expected2)
	}
}

func TestBlockQuote(t *testing.T) {
	result := BlockQuote("This is a quote")
	expected := "> This is a quote\n"
	if result != expected {
		t.Errorf("BlockQuote() = %q, want %q", result, expected)
	}

	// Test multiline
	result2 := BlockQuote("Line 1\nLine 2")
	expected2 := "> Line 1\n> Line 2\n"
	if result2 != expected2 {
		t.Errorf("BlockQuote() multiline = %q, want %q", result2, expected2)
	}
}

func TestHorizontalRule(t *testing.T) {
	result := HorizontalRule()
	expected := "---"
	if result != expected {
		t.Errorf("HorizontalRule() = %q, want %q", result, expected)
	}
}

func TestEscapeText(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		{"*bold*", "\\*"},
		{"_italic_", "\\_"},
		{"[link]", "\\["},
		{"`code`", "\\`"},
		{"---", "\\-"},
		{"\\n", "\\\\n"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := EscapeText(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("EscapeText(%q) should contain %q, got %q", tt.input, tt.contains, result)
			}
		})
	}
}

func TestIsMarkdownTable(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{"table", "| A | B |\n|---|---|\n| 1 | 2 |", true},
		{"not table", "Some text", false},
		{"incomplete table", "| A | B |\n| 1 | 2 |", false},
		{"no header", "|---|---|\n| 1 | 2 |", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMarkdownTable(tt.content)
			if result != tt.expected {
				t.Errorf("IsMarkdownTable() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsTaskList(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{"task list checked", "- [x] Done", true},
		{"task list unchecked", "- [ ] Todo", true},
		{"task list uppercase", "- [X] Done", true},
		{"not task list", "- Item", false},
		{"no dash", "[x] Item", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTaskList(tt.content)
			if result != tt.expected {
				t.Errorf("IsTaskList() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractLinks(t *testing.T) {
	content := "[Google](https://google.com) and [GitHub](https://github.com)"
	links := ExtractLinks(content)

	if len(links) != 2 {
		t.Fatalf("Expected 2 links, got %d", len(links))
	}

	if links[0].Text != "Google" {
		t.Errorf("Link[0].Text = %q, want %q", links[0].Text, "Google")
	}
	if links[0].URL != "https://google.com" {
		t.Errorf("Link[0].URL = %q, want %q", links[0].URL, "https://google.com")
	}

	if links[1].Text != "GitHub" {
		t.Errorf("Link[1].Text = %q, want %q", links[1].Text, "GitHub")
	}
	if links[1].URL != "https://github.com" {
		t.Errorf("Link[1].URL = %q, want %q", links[1].URL, "https://github.com")
	}
}

func TestDefaultRenderOptions(t *testing.T) {
	opts := DefaultRenderOptions()

	if opts.Width != 80 {
		t.Errorf("Default Width = 80, got %d", opts.Width)
	}
	if opts.Plain != false {
		t.Error("Default Plain = false")
	}
	if opts.PreserveNewlines != false {
		t.Error("Default PreserveNewlines = false")
	}
	if opts.TrimNewlines != true {
		t.Error("Default TrimNewlines = true")
	}
}
