package files

import (
	"os"
	"time"
)

// FileItem represents a file or directory in the file list.
type FileItem interface {
	// ID returns a unique identifier for the file.
	ID() string
	// Name returns the display name of the file.
	Name() string
	// Path returns the full path to the file.
	Path() string
	// Size returns the file size in bytes.
	Size() int64
	// Mode returns the file mode (permissions and type).
	Mode() os.FileMode
	// ModTime returns the modification time.
	ModTime() time.Time
	// IsDir returns true if this is a directory.
	IsDir() bool
	// Icon returns an icon representing this file type.
	Icon() string
	// Extension returns the file extension (without the dot).
	Extension() string
	// Description returns a brief description (e.g., file size, type).
	Description() string
}

// FileInfo represents basic file information.
type FileInfo struct {
	id       string
	name     string
	path     string
	size     int64
	mode     os.FileMode
	modTime  time.Time
	isDir    bool
	extension string
}

// NewFileInfo creates a new FileInfo from os.FileInfo.
func NewFileInfo(path string, info os.FileInfo) *FileInfo {
	name := info.Name()
	ext := ""

	// Extract extension for files
	if !info.IsDir() {
		for i := len(name) - 1; i >= 0; i-- {
			if name[i] == '.' {
				ext = name[i+1:]
				break
			}
		}
	}

	return &FileInfo{
		id:        path,
		name:      name,
		path:      path,
		size:      info.Size(),
		mode:      info.Mode(),
		modTime:   info.ModTime(),
		isDir:     info.IsDir(),
		extension: ext,
	}
}

// NewFileInfoSimple creates a FileInfo with provided values.
func NewFileInfoSimple(id, name, path string, size int64, mode os.FileMode, modTime time.Time, isDir bool) *FileInfo {
	ext := ""

	// Extract extension for files
	if !isDir {
		for i := len(name) - 1; i >= 0; i-- {
			if name[i] == '.' {
				ext = name[i+1:]
				break
			}
		}
	}

	return &FileInfo{
		id:        id,
		name:      name,
		path:      path,
		size:      size,
		mode:      mode,
		modTime:   modTime,
		isDir:     isDir,
		extension: ext,
	}
}

// ID returns the unique identifier for the file.
func (f *FileInfo) ID() string {
	return f.id
}

// Name returns the display name of the file.
func (f *FileInfo) Name() string {
	return f.name
}

// Path returns the full path to the file.
func (f *FileInfo) Path() string {
	return f.path
}

// Size returns the file size in bytes.
func (f *FileInfo) Size() int64 {
	return f.size
}

// Mode returns the file mode (permissions and type).
func (f *FileInfo) Mode() os.FileMode {
	return f.mode
}

// ModTime returns the modification time.
func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

// IsDir returns true if this is a directory.
func (f *FileInfo) IsDir() bool {
	return f.isDir
}

// Extension returns the file extension.
func (f *FileInfo) Extension() string {
	return f.extension
}

// Icon returns an icon representing this file type.
func (f *FileInfo) Icon() string {
	return GetIcon(f.isDir, f.extension)
}

// Description returns a brief description.
func (f *FileInfo) Description() string {
	if f.isDir {
		return "folder"
	}
	return formatSize(f.size)
}

// SortBy defines sorting criteria for file lists.
type SortBy int

const (
	// SortByName sorts files alphabetically by name.
	SortByName SortBy = iota
	// SortBySize sorts files by size (largest first).
	SortBySize
	// SortByTime sorts files by modification time (newest first).
	SortByTime
	// SortByExtension sorts files by extension.
	SortByExtension
)

// SortOrder defines the sort direction.
type SortOrder int

const (
	// SortAscending sorts in ascending order.
	SortAscending SortOrder = iota
	// SortDescending sorts in descending order.
	SortDescending
)

// FilterOptions defines filtering options for file lists.
type FilterOptions struct {
	// Pattern is the filter pattern (supports wildcards).
	Pattern string
	// CaseSensitive enables case-sensitive matching.
	CaseSensitive bool
	// IncludeDirs includes directories in results.
	IncludeDirs bool
	// IncludeFiles includes files in results.
	IncludeFiles bool
	// Extensions filters by specific extensions (e.g., ["go", "md"]).
	Extensions []string
	// MinSize filters files by minimum size in bytes.
	MinSize int64
	// MaxSize filters files by maximum size in bytes.
	MaxSize int64
	// HiddenFiles includes hidden files (starting with .).
	HiddenFiles bool
}

// DefaultFilterOptions returns default filter options.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{
		Pattern:       "",
		CaseSensitive: false,
		IncludeDirs:   true,
		IncludeFiles:  true,
		Extensions:    nil,
		MinSize:       0,
		MaxSize:       0,
		HiddenFiles:   false,
	}
}
