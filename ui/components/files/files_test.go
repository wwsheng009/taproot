package files

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewFileInfo tests creating a new FileInfo.
func TestNewFileInfo(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Get file info
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	// Create FileInfo
	fileInfo := NewFileInfo(testFile, info)

	// Verify fields
	if fileInfo.ID() != testFile {
		t.Errorf("ID() = %s, want %s", fileInfo.ID(), testFile)
	}
	if fileInfo.Name() != "test.txt" {
		t.Errorf("Name() = %s, want test.txt", fileInfo.Name())
	}
	if fileInfo.Path() != testFile {
		t.Errorf("Path() = %s, want %s", fileInfo.Path(), testFile)
	}
	if fileInfo.IsDir() {
		t.Error("IsDir() = true, want false")
	}
	if fileInfo.Extension() != "txt" {
		t.Errorf("Extension() = %s, want txt", fileInfo.Extension())
	}
	if fileInfo.Size() != 11 {
		t.Errorf("Size() = %d, want 11", fileInfo.Size())
	}
}

// TestNewFileInfoWithDirectory tests FileInfo for a directory.
func TestNewFileInfoWithDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Get directory info
	info, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get dir info: %v", err)
	}

	// Create FileInfo
	fileInfo := NewFileInfo(tmpDir, info)

	// Verify it's a directory
	if !fileInfo.IsDir() {
		t.Error("IsDir() = false, want true")
	}
	if fileInfo.Extension() != "" {
		t.Errorf("Extension() = %s, want empty string", fileInfo.Extension())
	}
}

// TestNewFileInfoSimple tests creating FileInfo with custom values.
func TestNewFileInfoSimple(t *testing.T) {
	now := time.Now()
	fileInfo := NewFileInfoSimple(
		"test-id",
		"simple.txt",
		"/path/to/simple.txt",
		1024,
		0644,
		now,
		false,
	)

	if fileInfo.ID() != "test-id" {
		t.Errorf("ID() = %s, want test-id", fileInfo.ID())
	}
	if fileInfo.Name() != "simple.txt" {
		t.Errorf("Name() = %s, want simple.txt", fileInfo.Name())
	}
	if fileInfo.Size() != 1024 {
		t.Errorf("Size() = %d, want 1024", fileInfo.Size())
	}
	if fileInfo.IsDir() {
		t.Error("IsDir() = true, want false")
	}
	if fileInfo.Extension() != "txt" {
		t.Errorf("Extension() = %s, want txt", fileInfo.Extension())
	}
}

// TestGetIcon tests file icon mapping.
func TestGetIcon(t *testing.T) {
	tests := []struct {
		isDir     bool
		extension string
		want      string
	}{
		{true, "", "📁"},
		{false, "go", "🐹"},
		{false, "py", "🐍"},
		{false, "js", "📜"},
		{false, "ts", "📘"},
		{false, "rs", "🦀"},
		{false, "c", "⚙️"},
		{false, "java", "☕"},
		{false, "html", "🌐"},
		{false, "css", "🎨"},
		{false, "json", "📋"},
		{false, "md", "📝"},
		{false, "txt", "📄"},
		{false, "png", "🖼️"},
		{false, "jpg", "🖼️"},
		{false, "pdf", "📕"},
		{false, "zip", "📦"},
		{false, "mp3", "🎵"},
		{false, "mp4", "🎬"},
		{false, "sql", "🗄️"},
		{false, "conf", "⚙️"},
		{false, "sh", "🐚"},
		{false, "unknown", "📄"},
	}

	for _, tt := range tests {
		t.Run(tt.extension, func(t *testing.T) {
			got := GetIcon(tt.isDir, tt.extension)
			if got != tt.want {
				t.Errorf("GetIcon(%v, %s) = %s, want %s", tt.isDir, tt.extension, got, tt.want)
			}
		})
	}
}

// TestFormatSize tests file size formatting.
func TestFormatSize(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1024 * 1024, "1.0 MiB"},
		{1024 * 1024 * 1024, "1.0 GiB"},
		{1024 * 1024 * 1024 * 1024, "1.0 TiB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatSize(tt.size)
			if got != tt.want {
				t.Errorf("formatSize(%d) = %s, want %s", tt.size, got, tt.want)
			}
		})
	}
}

// TestWildcardMatch tests wildcard pattern matching.
func TestWildcardMatch(t *testing.T) {
	tests := []struct {
		pattern       string
		text          string
		caseSensitive bool
		want          bool
	}{
		{"*", "anything", true, true},
		{"*", "anything", false, true},
		{"test", "test", true, true},
		{"test", "test", false, true},
		{"test", "TEST", true, false},
		{"test", "TEST", false, true},
		{"test*", "testing", true, true},
		{"test*", "testing", false, true},
		{"*test", "pretest", true, true},
		{"*test", "pretest", false, true},
		{"*est", "test", true, true},
		{"*est", "test", false, true},
		{"t?st", "test", true, true},
		{"t?st", "toast", true, false},
		{"*.txt", "file.txt", true, true},
		{"*.txt", "file.go", true, false},
		{"*.txt", "file.TXT", true, false},
		{"*.txt", "file.TXT", false, true},
		{"file.*.txt", "file.1.txt", true, true},
		{"file.*.txt", "file.1.log", true, false},
		{"*file*", "prefix_file_suffix", true, true},
		{"*file*", "prefix_file_suffix", false, true},
		{"", "", true, true},
		{"", "non-empty", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"/"+tt.text, func(t *testing.T) {
			got := WildcardMatch(tt.pattern, tt.text, tt.caseSensitive)
			if got != tt.want {
				t.Errorf("WildcardMatch(%q, %q, %v) = %v, want %v",
					tt.pattern, tt.text, tt.caseSensitive, got, tt.want)
			}
		})
	}
}

// TestContainsCaseInsensitive tests case-insensitive substring check.
func TestContainsCaseInsensitive(t *testing.T) {
	tests := []struct {
		s     string
		substr string
		want  bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "HELLO", true},
		{"Hello World", "world", true},
		{"Hello World", "WORLD", true},
		{"Hello World", "xyz", false},
		{"Hello World", "", true},
		{"test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"/"+tt.substr, func(t *testing.T) {
			got := ContainsCaseInsensitive(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("ContainsCaseInsensitive(%q, %q) = %v, want %v",
					tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

// TestHasPrefixCaseInsensitive tests case-insensitive prefix check.
func TestHasPrefixCaseInsensitive(t *testing.T) {
	tests := []struct {
		s      string
		prefix string
		want   bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "HELLO", true},
		{"Hello World", "world", false},
		{"Hello World", "H", true},
		{"Hello World", "Hello World", true},
		{"Hello World", "Hello World ", false},
	}

	for _, tt := range tests {
		t.Run(tt.s+"/"+tt.prefix, func(t *testing.T) {
			got := HasPrefixCaseInsensitive(tt.s, tt.prefix)
			if got != tt.want {
				t.Errorf("HasPrefixCaseInsensitive(%q, %q) = %v, want %v",
					tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}

// TestSortFiles tests sorting functionality.
func TestSortFiles(t *testing.T) {
	now := time.Now()

	files := []FileItem{
		NewFileInfoSimple("3", "three.txt", "/three.txt", 300, 0644, now.Add(2*time.Hour), false),
		NewFileInfoSimple("1", "one.txt", "/one.txt", 100, 0644, now.Add(0*time.Hour), false),
		NewFileInfoSimple("2", "dir", "/dir", 0, 0755, now.Add(1*time.Hour), true),
	}

	// Sort by name ascending
	t.Run("SortByName ascending", func(t *testing.T) {
		sorted := make([]FileItem, len(files))
		copy(sorted, files)
		SortFiles(sorted, SortByName, SortAscending)
		if sorted[0].Name() != "dir" || sorted[1].Name() != "one.txt" || sorted[2].Name() != "three.txt" {
			t.Errorf("SortByName ascending failed: %v", getItemNames(sorted))
		}
	})

	// Sort by size ascending
	t.Run("SortBySize ascending", func(t *testing.T) {
		sorted := make([]FileItem, len(files))
		copy(sorted, files)
		SortFiles(sorted, SortBySize, SortAscending)
		if sorted[0].Name() != "dir" || sorted[1].Size() != 100 || sorted[2].Size() != 300 {
			t.Errorf("SortBySize ascending failed: %v", getItemNames(sorted))
		}
	})

	// Sort by time ascending (newest first)
	t.Run("SortByTime ascending", func(t *testing.T) {
		sorted := make([]FileItem, len(files))
		copy(sorted, files)
		SortFiles(sorted, SortByTime, SortAscending)
		if sorted[0].Name() != "dir" || sorted[1].ModTime().After(sorted[2].ModTime()) {
			t.Errorf("SortByTime ascending failed: %v", getItemNames(sorted))
		}
	})
}

func getItemNames(files []FileItem) []string {
	names := make([]string, len(files))
	for i, f := range files {
		names[i] = f.Name()
	}
	return names
}

// TestFilterFiles tests filtering functionality.
func TestFilterFiles(t *testing.T) {
	now := time.Now()

	files := []FileItem{
		NewFileInfoSimple("1", "test.txt", "/test.txt", 100, 0644, now, false),
		NewFileInfoSimple("2", "test.go", "/test.go", 200, 0644, now, false),
		NewFileInfoSimple("3", "main.go", "/main.go", 300, 0644, now, false),
		NewFileInfoSimple("4", ".hidden", "/.hidden", 10, 0644, now, false),
		NewFileInfoSimple("5", "dir", "/dir", 0, 0755, now, true),
	}

	// Filter by extension
	t.Run("FilterByExtension", func(t *testing.T) {
		result := FilterByExtension(files, []string{"go"}, false)
		if len(result) != 2 {
			t.Errorf("FilterByExtension: got %d files, want 2", len(result))
		}
	})

	// Filter by name pattern
	t.Run("FilterByName", func(t *testing.T) {
		result := FilterByName(files, "test", false)
		if len(result) != 2 {
			t.Errorf("FilterByName: got %d files, want 2", len(result))
		}
	})

	// Filter hidden files
	t.Run("FilterHidden", func(t *testing.T) {
		result := FilterHidden(files)
		for _, f := range result {
			if f.Name() == ".hidden" {
				t.Error("FilterHidden: hidden file found")
			}
		}
	})

	// Filter directories only
	t.Run("FilterDirectories", func(t *testing.T) {
		result := FilterDirectories(files)
		if len(result) != 1 {
			t.Errorf("FilterDirectories: got %d files, want 1", len(result))
		}
		if !result[0].IsDir() {
			t.Error("FilterDirectories: non-directory found")
		}
	})

	// Filter files only
	t.Run("FilterFilesOnly", func(t *testing.T) {
		result := FilterFilesOnly(files)
		for _, f := range result {
			if f.IsDir() {
				t.Error("FilterFilesOnly: directory found")
			}
		}
	})

	// Filter by size
	t.Run("FilterBySize", func(t *testing.T) {
		result := FilterBySize(files, 150, 250)
		if len(result) != 1 {
			t.Errorf("FilterBySize: got %d files, want 1", len(result))
		}
		if result[0].Size() != 200 {
			t.Error("FilterBySize: wrong file selected")
		}
	})
}

// TestFileList tests the FileList component.
func TestFileList(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	filesToCreate := []string{"test.txt", "test.go", "main.go", "README.md"}
	for _, name := range filesToCreate {
		path := filepath.Join(tmpDir, name)
		err := os.WriteFile(path, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", name, err)
		}
	}

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Create FileList
	fileList, err := NewFileList(tmpDir)
	if err != nil {
		t.Fatalf("NewFileList failed: %v", err)
	}

	// Check path
	if fileList.Path() != tmpDir {
		t.Errorf("Path() = %s, want %s", fileList.Path(), tmpDir)
	}

	// Check counts (should include all files + subdir)
	if fileList.Count() != len(filesToCreate)+1 {
		t.Errorf("Count() = %d, want %d", fileList.Count(), len(filesToCreate)+1)
	}

	// Test sorting by name
	fileList.SetSort(SortByName, SortAscending)
	items := fileList.Filtered()
	if len(items) == 0 {
		t.Error("Filtered() returned empty list")
	}

	// Test filtering
	fileList.SetPattern("test")
	if fileList.FilteredCount() == 0 || fileList.FilteredCount() > 2 {
		t.Errorf("After filtering by 'test', got %d items", fileList.FilteredCount())
	}

	// Clear filter
	fileList.ClearFilter()
	if fileList.FilteredCount() != fileList.Count() {
		t.Error("ClearFilter() didn't restore all items")
	}

	// Test toggle sort order
	originalOrder := fileList.SortOrder()
	fileList.ToggleSortOrder()
	if fileList.SortOrder() == originalOrder {
		t.Error("ToggleSortOrder() didn't change order")
	}

	// Test stats
	stats := fileList.GetStats()
	if stats.TotalItems != fileList.Count() {
		t.Errorf("Stats.TotalItems = %d, want %d", stats.TotalItems, fileList.Count())
	}
	if stats.FileCount == 0 {
		t.Error("Stats.FileCount is 0")
	}
}

// TestDefaultFilterOptions tests default filter options.
func TestDefaultFilterOptions(t *testing.T) {
	opts := DefaultFilterOptions()
	if opts.IncludeDirs != true {
		t.Error("DefaultFilterOptions: IncludeDirs should be true")
	}
	if opts.IncludeFiles != true {
		t.Error("DefaultFilterOptions: IncludeFiles should be true")
	}
	if opts.HiddenFiles != false {
		t.Error("DefaultFilterOptions: HiddenFiles should be false")
	}
}
