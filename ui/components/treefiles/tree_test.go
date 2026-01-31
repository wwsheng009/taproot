package treefiles

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestTree(t *testing.T) (string, func()) {
	// Create test directory structure
	baseDir, err := os.MkdirTemp("", "treefiles_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create subdirectories
	dirs := []string{
		"dir1",
		"dir2",
		filepath.Join("dir1", "subdir1"),
		filepath.Join("dir1", "subdir2"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			os.RemoveAll(baseDir)
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	// Create files
	files := map[string]string{
		"file1.txt":         "content1",
		"file2.go":          "package main\n",
		"dir1/subfile1.txt": "subcontent1",
		"dir1/subdir1/file.py": "import os\n",
		".hidden":         "hidden",
		"dir2/file.json":  `{"key": "value"}`,
	}
	for file, content := range files {
		if err := os.WriteFile(filepath.Join(baseDir, file), []byte(content), 0644); err != nil {
			os.RemoveAll(baseDir)
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	cleanup := func() {
		os.RemoveAll(baseDir)
	}

	return baseDir, cleanup
}

func TestNewFileTree(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	if tree.Root() == nil {
		t.Fatal("Root node is nil")
	}

	if tree.Root().Path() != baseDir {
		t.Errorf("Root path mismatch: got %s, want %s", tree.Root().Path(), baseDir)
	}

	if !tree.Root().IsDir() {
		t.Error("Root should be a directory")
	}
}

func TestNewFileTreeWithOptions(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir,
		WithSort(SortBySize, SortDescending),
		WithHidden(true),
		WithMaxDepth(3),
	)
	if err != nil {
		t.Fatalf("NewFileTree with options failed: %v", err)
	}

	if tree.sortBy != SortBySize {
		t.Errorf("sortBy not set correctly: got %v, want %v", tree.sortBy, SortBySize)
	}

	if tree.sortOrder != SortDescending {
		t.Errorf("sortOrder not set correctly: got %v, want %v", tree.sortOrder, SortDescending)
	}

	if !tree.IncludesHidden() {
		t.Error("hidden not set to true")
	}

	if tree.maxDepth != 3 {
		t.Errorf("maxDepth not set correctly: got %d, want 3", tree.maxDepth)
	}
}

func TestFlatten(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	// Default: root expanded, children collapsed
	flat := tree.Flatten()
	if len(flat) < 1 {
		t.Fatal("Flatten returned empty list")
	}

	// First node should be root
	if flat[0].Path() != baseDir {
		t.Errorf("First node should be root: got %s", flat[0].Path())
	}

	// Should have root files and directories (not deep yet)
	if len(flat) < 3 {
		t.Error("Flatten should return at least root + dir1 + dir2 + files")
	}
}

func TestToggleNode(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	dir1 := tree.Root().FindChild("dir1")
	if dir1 == nil {
		t.Fatal("dir1 not found in root")
	}

	// Initially collapsed (except root)
	if dir1.Expanded() {
		t.Error("dir1 should initially be collapsed")
	}

	// Toggle to expand
	tree.ToggleNode(dir1)
	if !dir1.Expanded() {
		t.Error("dir1 should be expanded after toggle")
	}

	// Toggle to collapse
	tree.ToggleNode(dir1)
	if dir1.Expanded() {
		t.Error("dir1 should be collapsed after second toggle")
	}
}

func TestExpandAll(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	// Expand all directories
	tree.ExpandAll()

	flat := tree.Flatten()

	// Should include all nested directories and files
	if len(flat) < 8 {
		t.Errorf("ExpandAll should show all nodes: got %d nodes, want at least 8", len(flat))
	}

	// Check that directories are expanded
	for _, node := range flat {
		if node.IsDir() {
			if !node.Expanded() && node != tree.Root() {
				t.Errorf("Directory %s should be expanded", node.Name())
			}
		}
	}
}

func TestCollapseAll(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	// Expand first
	tree.ExpandAll()

	// Then collapse
	tree.CollapseAll()

	flat := tree.Flatten()

	// Should be minimal (root + 4 children = 5)
	if len(flat) > 5 {
		t.Errorf("After collapse, should be minimal: got %d nodes", len(flat))
	}

	// Root should remain expanded
	if !tree.Root().Expanded() {
		t.Error("Root should remain expanded after CollapseAll")
	}
}

func TestFindNode(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	tests := []struct {
		path   string
		want   string
		exists bool
	}{
		{baseDir, filepath.Base(baseDir), true},
		{filepath.Join(baseDir, "dir1"), "dir1", true},
		{filepath.Join(baseDir, "dir1", "subdir1"), "subdir1", true},
		{filepath.Join(baseDir, "file1.txt"), "file1.txt", true},
		{filepath.Join(baseDir, "nonexistent.txt"), "", false},
	}

	for _, tt := range tests {
		node := tree.FindNode(tt.path)
		if tt.exists {
			if node == nil {
				t.Errorf("FindNode(%s) returned nil, want exists", tt.path)
			} else if node.Name() != tt.want {
				t.Errorf("FindNode(%s) = %s, want %s", tt.path, node.Name(), tt.want)
			}
		} else {
			if node != nil {
				t.Errorf("FindNode(%s) returned node, want nil", tt.path)
			}
		}
	}
}

func TestStats(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	stats := tree.Stats()

	if stats.TotalNodes == 0 {
		t.Error("Stats should have total nodes")
	}

	if stats.TotalDirs == 0 {
		t.Error("Stats should have at least root directory")
	}

	if stats.TotalFiles == 0 {
		t.Error("Stats should have files")
	}

	// All files exist except hidden
	// File structure:
	// - file1.txt, file2.go
	// - dir1/subfile1.txt
	// - dir1/subdir1/file.py
	// - dir2/file.json
	// = 5 non-hidden files
	if stats.TotalFiles != 5 {
		t.Errorf("Stats.TotalFiles = %d, want 5 (excluding hidden)", stats.TotalFiles)
	}

	if stats.ExpandedDirs != 1 {
		t.Errorf("Stats.ExpandedDirs = %d, want 1 (only root expanded)", stats.ExpandedDirs)
	}
}

func TestStatsWithHidden(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir, WithHidden(true))
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	stats := tree.Stats()

	// With hidden files included, should have 6 files (5 + .hidden)
	if stats.TotalFiles != 6 {
		t.Errorf("Stats.TotalFiles with hidden = %d, want 6", stats.TotalFiles)
	}
}

func TestToggleHidden(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	stats1 := tree.Stats()
	tree.ToggleHidden()
	stats2 := tree.Stats()

	// After toggling hidden, should have more files
	if stats2.TotalFiles <= stats1.TotalFiles {
		t.Errorf("ToggleHidden should increase file count: before %d, after %d", stats1.TotalFiles, stats2.TotalFiles)
	}

	if !tree.IncludesHidden() {
		t.Error("IncludesHidden should be true after toggle")
	}
}

func TestSortByName(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir, WithSort(SortByName, SortAscending))
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	// Check that siblings are sorted by name (directories first within each parent)
	checkSorted := func(parent *FileNode) {
		if !parent.IsDir() || len(parent.Children()) == 0 {
			return
		}
		// Within directories, children should be sorted with directories first
		// Find split between dirs and files
		lastDirIdx := -1
		for i, child := range parent.Children() {
			if child.IsDir() {
				lastDirIdx = i
			}
		}
		// Check directories are sorted
		if lastDirIdx >= 1 {
			for i := 1; i <= lastDirIdx; i++ {
				if strings.Compare(parent.Children()[i-1].Name(), parent.Children()[i].Name()) > 0 {
					t.Errorf("Under %s: directories not sorted: %s before %s", parent.Name(), parent.Children()[i-1].Name(), parent.Children()[i].Name())
				}
			}
		}
		// Check files are sorted
		for i := lastDirIdx + 2; i < len(parent.Children()); i++ {
			if strings.Compare(parent.Children()[i-1].Name(), parent.Children()[i].Name()) > 0 {
				t.Errorf("Under %s: files not sorted: %s before %s", parent.Name(), parent.Children()[i-1].Name(), parent.Children()[i].Name())
			}
		}
	}

	checkSorted(tree.Root())
	for _, child := range tree.Root().Children() {
		if child.IsDir() {
			checkSorted(child)
		}
	}
}

func TestSortBySize(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir, WithSort(SortBySize, SortAscending))
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	tree.ExpandAll()

	// Just verify it doesn't crash - file sizes may vary
	flat := tree.Flatten()
	if len(flat) == 0 {
		t.Fatal("Flatten returned empty after sort")
	}
}

func TestSortByType(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir, WithSort(SortByType, SortAscending))
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	// Check that siblings have directories before files (within each parent)
	checkTypeSorted := func(parent *FileNode) {
		if !parent.IsDir() || len(parent.Children()) == 0 {
			return
		}
		children := parent.Children()
		seenFile := false
		for _, child := range children {
			if !child.IsDir() {
				seenFile = true
			} else if child.IsDir() && seenFile {
				t.Errorf("Under %s: directory %s appears after a file", parent.Name(), child.Name())
			}
		}
	}

	checkTypeSorted(tree.Root())
	for _, child := range tree.Root().Children() {
		if child.IsDir() {
			checkTypeSorted(child)
		}
	}
}

func TestSetSort(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir, WithSort(SortByName, SortAscending))
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	// Change sort to size descending
	tree.SetSort(SortBySize, SortDescending)

	if tree.sortBy != SortBySize {
		t.Errorf("SetSort didn't update sortBy: got %v, want %v", tree.sortBy, SortBySize)
	}

	if tree.sortOrder != SortDescending {
		t.Errorf("SetSort didn't update sortOrder: got %v, want %v", tree.sortOrder, SortDescending)
	}
}

func TestRescan(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	// Create a new file after tree creation
	newFile := filepath.Join(baseDir, "newfile.txt")
	if err := os.WriteFile(newFile, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Rescan should pick up the new file
	if err := tree.Rescan(); err != nil {
		t.Fatalf("Rescan failed: %v", err)
	}

	if tree.FindNode(newFile) == nil {
		t.Error("Rescan didn't find newly created file")
	}
}

func TestMaxDepth(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	// Create deeper directory structure
	deepDir := filepath.Join(baseDir, "dir1", "subdir1", "deep")
	if err := os.MkdirAll(deepDir, 0755); err != nil {
		cleanup()
		t.Fatalf("Failed to create deep dir: %v", err)
	}
	os.WriteFile(filepath.Join(deepDir, "deep.txt"), []byte("deep"), 0644)

	tree, err := NewFileTree(baseDir, WithMaxDepth(2))
	if err != nil {
		cleanup()
		t.Fatalf("NewFileTree failed: %v", err)
	}

	tree.ExpandAll()
	flat := tree.Flatten()

	// Should include root (depth 0), dir1 (depth 1), subdirs (depth 2)
	// But not the deep directory (depth 3)
	for _, node := range flat {
		if node.Depth() > 2 {
			t.Errorf("Node %s at depth %d exceeds max depth 2", node.Path(), node.Depth())
		}
	}
}

func TestNodeMethods(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	root := tree.Root()
	_ = root.FindChild("dir1")

	// Test AddChild and RemoveChild
	newNode := &FileNode{
		path:   filepath.Join(baseDir, "test"),
		name:   "test",
		isDir:  false,
		size:   100,
		modTime: "Jan 01 12:00",
		depth:  1,
	}

	root.AddChild(newNode)

	if root.FindChild("test") == nil {
		t.Error("AddChild didn't add the node")
	}

	if newNode.Parent() != root {
		t.Error("Parent not set correctly in AddChild")
	}

	if newNode.Depth() != 1 {
		t.Errorf("Depth not set correctly: got %d, want 1", newNode.Depth())
	}

	// Test HasChildren
	if !root.HasChildren() {
		t.Error("Root should have children")
	}

	file1 := tree.FindNode(filepath.Join(baseDir, "file1.txt"))
	if file1 != nil && file1.HasChildren() {
		t.Error("File should not have children")
	}

	// Test RemoveChild
	root.RemoveChild(newNode)
	if root.FindChild("test") != nil {
		t.Error("RemoveChild didn't remove the node")
	}
}

func TestFileNodeGetters(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	root := tree.Root()

	// Test getters
	if root.Path() != baseDir {
		t.Errorf("Path() = %s, want %s", root.Path(), baseDir)
	}

	if root.Name() != filepath.Base(baseDir) {
		t.Errorf("Name() = %s, want %s", root.Name(), filepath.Base(baseDir))
	}

	if !root.IsDir() {
		t.Error("Root should be a directory")
	}

	// Size should be 0 for directories
	if root.Size() != 0 {
		t.Errorf("Size() for directory = %d, want 0", root.Size())
	}

	// Find a file
	file1 := tree.FindNode(filepath.Join(baseDir, "file1.txt"))
	if file1 == nil {
		t.Fatal("file1.txt not found")
	}

	if file1.IsDir() {
		t.Error("file1.txt should not be a directory")
	}

	if file1.Size() <= 0 {
		t.Errorf("Size() for file = %d, want > 0", file1.Size())
	}

	if file1.ModTime() == "" {
		t.Error("ModTime() should not be empty")
	}
}

func TestGetFileIcon(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"test.go", "Go"},
		{"test.py", "Python"},
		{"test.js", "JavaScript"},
		{"test.txt", "Text"},
		{"test.json", "JSON"},
		{"test.md", "Markdown"},
		{"test.png", "Image"},
		{"test.pdf", "PDF"},
		{"test.zip", "Archive"},
		{"test.yaml", "YAML"},
		{"test.yml", "YAML"},
		{"test.toml", "Config"},
		{"test.sh", "Shell"},
		{"test", "File"},
	}

	for _, tt := range tests {
		result := GetFileIcon(tt.name)
		if result != tt.expected {
			t.Errorf("GetFileIcon(%s) = %s, want %s", tt.name, result, tt.expected)
		}
	}
}

func TestGetTreeIcon(t *testing.T) {
	baseDir, cleanup := setupTestTree(t)
	defer cleanup()

	tree, err := NewFileTree(baseDir)
	if err != nil {
		t.Fatalf("NewFileTree failed: %v", err)
	}

	dir1 := tree.Root().FindChild("dir1")
	if dir1 == nil {
		t.Fatal("dir1 not found")
	}

	// Collapsed directory
	icon := GetTreeIcon(dir1)
	if icon != IconCollapsed {
		t.Errorf("GetTreeIcon(collapsed dir) = %s, want %s", icon, IconCollapsed)
	}

	// Expanded directory
	dir1.SetExpanded(true)
	icon = GetTreeIcon(dir1)
	if icon != IconExpanded {
		t.Errorf("GetTreeIcon(expanded dir) = %s, want %s", icon, IconExpanded)
	}

	// File
	file1 := tree.FindNode(filepath.Join(baseDir, "file1.txt"))
	if file1 == nil {
		t.Fatal("file1.txt not found")
	}

	icon = GetTreeIcon(file1)
	if icon != "Text" {
		t.Errorf("GetTreeIcon(file) = %s, want Text", icon)
	}
}
