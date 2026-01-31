package files

import "path/filepath"

// FilterFiles filters a slice of FileItem based on the provided options.
func FilterFiles(files []FileItem, opts FilterOptions) []FileItem {
	var result []FileItem

	for _, f := range files {
		if matchesFilter(f, opts) {
			result = append(result, f)
		}
	}

	return result
}

// matchesFilter checks if a file item matches the filter options.
func matchesFilter(f FileItem, opts FilterOptions) bool {
	// Check type filters
	if !matchesType(f, opts) {
		return false
	}

	// Check hidden files
	if !matchesHidden(f, opts) {
		return false
	}

	// Check extension filter
	if !matchesExtension(f, opts) {
		return false
	}

	// Check size filter (only applies to files)
	if !matchesSize(f, opts) {
		return false
	}

	// Check pattern filter
	if !matchesPattern(f, opts) {
		return false
	}

	return true
}

// matchesType checks if the file matches type filters (dirs/files).
func matchesType(f FileItem, opts FilterOptions) bool {
	if f.IsDir() {
		return opts.IncludeDirs
	}
	return opts.IncludeFiles
}

// matchesHidden checks if hidden files should be included.
func matchesHidden(f FileItem, opts FilterOptions) bool {
	// Hidden files start with a dot
	if !opts.HiddenFiles {
		name := f.Name()
		if len(name) > 0 && name[0] == '.' {
			return false
		}
	}
	return true
}

// matchesExtension checks if the file matches the extension filter.
func matchesExtension(f FileItem, opts FilterOptions) bool {
	// No extension filter specified
	if len(opts.Extensions) == 0 {
		return true
	}

	// Skip directories for extension filtering
	if f.IsDir() {
		return true
	}

	ext := f.Extension()
	for _, allowedExt := range opts.Extensions {
		if opts.CaseSensitive {
			if ext == allowedExt {
				return true
			}
		} else {
			if ContainsCaseInsensitive(ext, allowedExt) {
				return true
			}
		}
	}

	return false
}

// matchesSize checks if the file matches size constraints.
func matchesSize(f FileItem, opts FilterOptions) bool {
	// Size filters don't apply to directories
	if f.IsDir() {
		return true
	}

	// No size filters specified
	if opts.MinSize == 0 && opts.MaxSize == 0 {
		return true
	}

	size := f.Size()

	// Check minimum size
	if opts.MinSize > 0 && size < opts.MinSize {
		return false
	}

	// Check maximum size
	if opts.MaxSize > 0 && size > opts.MaxSize {
		return false
	}

	return true
}

// matchesPattern checks if the file matches the pattern filter.
func matchesPattern(f FileItem, opts FilterOptions) bool {
	// No pattern specified
	if opts.Pattern == "" {
		return true
	}

	// Check if pattern is a full wildcard (matches everything)
	if opts.Pattern == "*" {
		return true
	}

	name := f.Name()

	// Try exact match first
	if opts.CaseSensitive {
		if name == opts.Pattern {
			return true
		}
	} else {
		if ContainsCaseInsensitive(name, opts.Pattern) {
			return true
		}
	}

	// Try wildcard matching
	matched := WildcardMatch(opts.Pattern, name, opts.CaseSensitive)

	// If pattern matching failed, try as prefix
	if !matched {
		if opts.CaseSensitive {
			if HasPrefix(name, opts.Pattern) {
				return true
			}
		} else {
			if HasPrefixCaseInsensitive(name, opts.Pattern) {
				return true
			}
		}
	}

	return matched
}

// HasPrefix checks if string starts with prefix (case-sensitive).
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// FilterByExtension returns files matching the specified extensions.
func FilterByExtension(files []FileItem, extensions []string, caseSensitive bool) []FileItem {
	opts := FilterOptions{
		Pattern:       "",
		CaseSensitive: caseSensitive,
		IncludeDirs:   false,
		IncludeFiles:  true,
		Extensions:    extensions,
		MinSize:       0,
		MaxSize:       0,
		HiddenFiles:   true,
	}
	return FilterFiles(files, opts)
}

// FilterByName returns files matching the specified pattern.
func FilterByName(files []FileItem, pattern string, caseSensitive bool) []FileItem {
	opts := DefaultFilterOptions()
	opts.Pattern = pattern
	opts.CaseSensitive = caseSensitive
	return FilterFiles(files, opts)
}

// FilterBySize returns files within the specified size range.
// minSize and maxSize are in bytes. Use 0 for no limit.
func FilterBySize(files []FileItem, minSize, maxSize int64) []FileItem {
	opts := FilterOptions{
		Pattern:       "",
		CaseSensitive: false,
		IncludeDirs:   false,
		IncludeFiles:  true,
		Extensions:    nil,
		MinSize:       minSize,
		MaxSize:       maxSize,
		HiddenFiles:   true,
	}
	return FilterFiles(files, opts)
}

// FilterHidden returns files excluding hidden files (names starting with .).
func FilterHidden(files []FileItem) []FileItem {
	opts := DefaultFilterOptions()
	opts.HiddenFiles = false
	return FilterFiles(files, opts)
}

// FilterDirectories returns only directories.
func FilterDirectories(files []FileItem) []FileItem {
	opts := FilterOptions{
		IncludeDirs:  true,
		IncludeFiles: false,
	}
	return FilterFiles(files, opts)
}

// FilterFilesOnly returns only files (no directories).
func FilterFilesOnly(files []FileItem) []FileItem {
	opts := FilterOptions{
		IncludeDirs:  false,
		IncludeFiles: true,
	}
	return FilterFiles(files, opts)
}

// FilterByPath returns files whose path matches the pattern.
func FilterByPath(files []FileItem, pattern string, caseSensitive bool) []FileItem {
	var result []FileItem

	for _, f := range files {
		path := f.Path()
		match := false

		if caseSensitive {
			if path == pattern {
				match = true
			} else if WildcardMatch(pattern, path, true) {
				match = true
			}
		} else {
			if ContainsCaseInsensitive(path, pattern) {
				match = true
			} else if WildcardMatch(pattern, path, false) {
				match = true
			}
		}

		if match {
			result = append(result, f)
		}
	}

	return result
}

// FilterByGlobPath returns files matching a glob pattern.
// This uses filepath.Match for glob matching.
func FilterByGlobPath(files []FileItem, globPattern string) []FileItem {
	var result []FileItem

	for _, f := range files {
		matched, err := filepath.Match(globPattern, f.Path())
		if err == nil && matched {
			result = append(result, f)
		}
	}

	return result
}
