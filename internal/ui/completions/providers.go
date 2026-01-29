package completions

import (
	"os"
	"path/filepath"
	"strings"
)

// StringProvider is a simple provider based on string slices.
type StringProvider struct {
	items []CompletionItem
}

// NewStringProvider creates a new StringProvider.
func NewStringProvider(items []CompletionItem) *StringProvider {
	return &StringProvider{items: items}
}

// NewStringProviderFromStrings creates a StringProvider from simple strings.
func NewStringProviderFromStrings(strings []string) *StringProvider {
	items := make([]CompletionItem, len(strings))
	for i, s := range strings {
		items[i] = NewSimpleCompletionItem(s, s, s)
	}
	return &StringProvider{items: items}
}

// GetItems returns all items.
func (p *StringProvider) GetItems() []CompletionItem {
	return p.items
}

// GetFilterValue returns the display text for filtering.
func (p *StringProvider) GetFilterValue(item CompletionItem) string {
	return item.Display()
}

// FileProvider provides file system completions.
type FileProvider struct {
	baseDir string
	depth   int
	ignore  []string
}

// FileCompletionValue represents a file completion value.
type FileCompletionValue struct {
	Path     string
	IsDir    bool
	Absolute bool
}

// NewFileProvider creates a new FileProvider.
func NewFileProvider(baseDir string, depth int, ignore []string) *FileProvider {
	return &FileProvider{
		baseDir: baseDir,
		depth:   depth,
		ignore:  ignore,
	}
}

// GetItems returns all files in the base directory up to the specified depth.
func (p *FileProvider) GetItems() []CompletionItem {
	items := []CompletionItem{}

	if p.baseDir == "" {
		p.baseDir = "."
	}

	err := filepath.Walk(p.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Determine depth
		relPath, err := filepath.Rel(p.baseDir, path)
		if err != nil {
			return nil
		}

		if relPath == "." {
			return nil // Skip base directory
		}

		depth := len(strings.Split(relPath, string(filepath.Separator)))

		// Check depth limit
		if p.depth > 0 && depth > p.depth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check ignore patterns
		for _, ignorePattern := range p.ignore {
			if strings.Contains(relPath, ignorePattern) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Create item
		value := FileCompletionValue{
			Path:     path,
			IsDir:    info.IsDir(),
			Absolute: filepath.IsAbs(path),
		}

		display := relPath
		if info.IsDir() {
			display += "/"
		}

		items = append(items, NewSimpleCompletionItem(relPath, display, value))

		return nil
	})

	if err != nil {
		return items
	}

	return items
}

// GetFilterValue returns the relative path for filtering.
func (p *FileProvider) GetFilterValue(item CompletionItem) string {
	return item.Display()
}

// CommandProvider provides command completions.
type CommandProvider struct {
	commands []CommandItem
}

// CommandItem represents a command with description.
type CommandItem struct {
	Name        string
	Description string
	Handler     func(args ...string) any
}

// NewCommandProvider creates a new CommandProvider.
func NewCommandProvider(commands []CommandItem) *CommandProvider {
	items := make([]CompletionItem, len(commands))
	for i, cmd := range commands {
		display := cmd.Name
		if cmd.Description != "" {
			display += " - " + cmd.Description
		}
		items[i] = NewSimpleCompletionItem(cmd.Name, display, cmd.Handler)
	}
	return &CommandProvider{commands: commands}
}

// GetItems returns all commands.
func (p *CommandProvider) GetItems() []CompletionItem {
	items := make([]CompletionItem, len(p.commands))
	for i, cmd := range p.commands {
		display := cmd.Name
		if cmd.Description != "" {
			display += " - " + cmd.Description
		}
		items[i] = NewSimpleCompletionItem(cmd.Name, display, cmd.Handler)
	}
	return items
}

// GetFilterValue returns the command name for filtering.
func (p *CommandProvider) GetFilterValue(item CompletionItem) string {
	return item.Display()
}
