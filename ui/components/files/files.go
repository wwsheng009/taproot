package files

import (
	"os"
	"path/filepath"
)

// FileList manages a list of files with sorting, filtering, and navigation.
type FileList struct {
	path      string          // Current directory path
	items     []FileItem      // All items in the directory
	filtered  []FileItem      // Items after filtering
	sortBy    SortBy          // Current sort criteria
	sortOrder SortOrder       // Current sort order
	filter    FilterOptions   // Current filter options
}

// FileOption is a function that configures a FileList.
type FileOption func(*FileList)

// WithSort sets the sort criteria and order.
func WithSort(sortBy SortBy, sortOrder SortOrder) FileOption {
	return func(fl *FileList) {
		fl.sortBy = sortBy
		fl.sortOrder = sortOrder
	}
}

// WithFilter sets the filter options.
func WithFilter(filter FilterOptions) FileOption {
	return func(fl *FileList) {
		fl.filter = filter
	}
}

// WithIncludeHidden includes hidden files in the list.
func WithIncludeHidden(include bool) FileOption {
	return func(fl *FileList) {
		fl.filter.HiddenFiles = include
	}
}

// WithExtensions sets the allowed extensions.
func WithExtensions(extensions []string) FileOption {
	return func(fl *FileList) {
		fl.filter.Extensions = extensions
	}
}

// NewFileList creates a new FileList.
func NewFileList(path string, opts ...FileOption) (*FileList, error) {
	fl := &FileList{
		path:      path,
		sortBy:    SortByName,
		sortOrder: SortAscending,
		filter:    DefaultFilterOptions(),
	}

	// Apply options
	for _, opt := range opts {
		opt(fl)
	}

	// Load directory
	if err := fl.LoadDirectory(path); err != nil {
		return nil, err
	}

	return fl, nil
}

// LoadDirectory loads files from a directory.
func (fl *FileList) LoadDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Clear current items
	fl.items = make([]FileItem, 0, len(entries))
	fl.path = path

	// Create FileItem for each entry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // Skip files that can't be read
		}

		fullPath := filepath.Join(path, entry.Name())
		fileItem := NewFileInfo(fullPath, info)
		fl.items = append(fl.items, fileItem)
	}

	// Sort and filter items
	fl.refresh()

	return nil
}

// refresh applies sorting and filtering to items.
func (fl *FileList) refresh() {
	// Sort items
	SortFiles(fl.items, fl.sortBy, fl.sortOrder)

	// Filter items
	fl.filtered = FilterFiles(fl.items, fl.filter)
}

// Path returns the current directory path.
func (fl *FileList) Path() string {
	return fl.path
}

// Items returns all items (before filtering).
func (fl *FileList) Items() []FileItem {
	return fl.items
}

// Filtered returns items after filtering.
func (fl *FileList) Filtered() []FileItem {
	return fl.filtered
}

// Count returns the number of items before filtering.
func (fl *FileList) Count() int {
	return len(fl.items)
}

// FilteredCount returns the number of items after filtering.
func (fl *FileList) FilteredCount() int {
	return len(fl.filtered)
}

// GetItem returns an item by index from the filtered list.
func (fl *FileList) GetItem(index int) (FileItem, bool) {
	if index < 0 || index >= len(fl.filtered) {
		return nil, false
	}
	return fl.filtered[index], true
}

// SetSort sets the sort criteria and order.
func (fl *FileList) SetSort(sortBy SortBy, sortOrder SortOrder) {
	fl.sortBy = sortBy
	fl.sortOrder = sortOrder
	fl.refresh()
}

// SortBy returns the current sort criteria.
func (fl *FileList) SortBy() SortBy {
	return fl.sortBy
}

// SortOrder returns the current sort order.
func (fl *FileList) SortOrder() SortOrder {
	return fl.sortOrder
}

// SetFilter sets the filter options.
func (fl *FileList) SetFilter(filter FilterOptions) {
	fl.filter = filter
	fl.refresh()
}

// Filter returns the current filter options.
func (fl *FileList) Filter() FilterOptions {
	return fl.filter
}

// SetPattern sets the filter pattern.
func (fl *FileList) SetPattern(pattern string) {
	fl.filter.Pattern = pattern
	fl.refresh()
}

// Pattern returns the current filter pattern.
func (fl *FileList) Pattern() string {
	return fl.filter.Pattern
}

// ToggleSortOrder toggles between ascending and descending sort order.
func (fl *FileList) ToggleSortOrder() {
	if fl.sortOrder == SortAscending {
		fl.sortOrder = SortDescending
	} else {
		fl.sortOrder = SortAscending
	}
	fl.refresh()
}

// ClearFilter clears the filter pattern.
func (fl *FileList) ClearFilter() {
	fl.filter.Pattern = ""
	fl.refresh()
}

// ToggleHiddenFiles toggles visibility of hidden files.
func (fl *FileList) ToggleHiddenFiles() {
	fl.filter.HiddenFiles = !fl.filter.HiddenFiles
	fl.refresh()
}

// IncludesHidden returns true if hidden files are included.
func (fl *FileList) IncludesHidden() bool {
	return fl.filter.HiddenFiles
}

// ListDirectories returns only directories from the filtered list.
func (fl *FileList) ListDirectories() []FileItem {
	var dirs []FileItem
	for _, item := range fl.filtered {
		if item.IsDir() {
			dirs = append(dirs, item)
		}
	}
	return dirs
}

// ListFiles returns only files from the filtered list.
func (fl *FileList) ListFiles() []FileItem {
	var files []FileItem
	for _, item := range fl.filtered {
		if !item.IsDir() {
			files = append(files, item)
		}
	}
	return files
}

// FindItem finds an item by its path in the filtered list.
func (fl *FileList) FindItem(path string) (FileItem, bool) {
	for _, item := range fl.filtered {
		if item.Path() == path {
			return item, true
		}
	}
	return nil, false
}

// FindItemByName finds an item by name in the filtered list.
func (fl *FileList) FindItemByName(name string) (FileItem, bool) {
	for _, item := range fl.filtered {
		if item.Name() == name {
			return item, true
		}
	}
	return nil, false
}

// GetStats returns statistics about the file list.
func (fl *FileList) GetStats() FileListStats {
	totalSize := int64(0)
	dirCount := 0
	fileCount := 0

	for _, item := range fl.filtered {
		if item.IsDir() {
			dirCount++
		} else {
			fileCount++
			totalSize += item.Size()
		}
	}

	return FileListStats{
		TotalItems:  len(fl.filtered),
		TotalSize:   totalSize,
		DirectoryCount: dirCount,
		FileCount:    fileCount,
		AverageSize:  averageSize(fileCount, totalSize),
	}
}

// FileListStats contains statistics about a file list.
type FileListStats struct {
	TotalItems     int
	TotalSize      int64
	DirectoryCount int
	FileCount      int
	AverageSize    int64
}

// averageSize calculates the average file size.
func averageSize(count int, total int64) int64 {
	if count == 0 {
		return 0
	}
	return total / int64(count)
}
