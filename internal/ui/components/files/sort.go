package files

import (
	"sort"
	"strings"
)

// SortFiles sorts a slice of FileItems by the specified criteria.
func SortFiles(files []FileItem, sortBy SortBy, order SortOrder) {
	switch sortBy {
	case SortByName:
		sortFilesByName(files, order)
	case SortBySize:
		sortFilesBySize(files, order)
	case SortByTime:
		sortFilesByTime(files, order)
	case SortByExtension:
		sortFilesByExtension(files, order)
	}
}

// sortFilesByName sorts files alphabetically by name.
// Directories always come first.
func sortFilesByName(files []FileItem, order SortOrder) {
	sort.Slice(files, func(i, j int) bool {
		di, dj := files[i].IsDir(), files[j].IsDir()
		if di != dj {
			return di // directories first
		}

		if order == SortDescending {
			i, j = j, i
		}

		return strings.Compare(files[i].Name(), files[j].Name()) < 0
	})
}

// sortFilesBySize sorts files by size.
// Directories always come first (with size 0).
func sortFilesBySize(files []FileItem, order SortOrder) {
	sort.Slice(files, func(i, j int) bool {
		di, dj := files[i].IsDir(), files[j].IsDir()
		if di != dj {
			return di // directories first
		}
		if di && dj {
			return strings.Compare(files[i].Name(), files[j].Name()) < 0
		}

		if order == SortDescending {
			return files[i].Size() > files[j].Size()
		}
		return files[i].Size() < files[j].Size()
	})
}

// sortFilesByTime sorts files by modification time.
// Directories always come first.
func sortFilesByTime(files []FileItem, order SortOrder) {
	sort.Slice(files, func(i, j int) bool {
		di, dj := files[i].IsDir(), files[j].IsDir()
		if di != dj {
			return di // directories first
		}
		if di && dj {
			return strings.Compare(files[i].Name(), files[j].Name()) < 0
		}

		ti, tj := files[i].ModTime(), files[j].ModTime()
		if order == SortDescending {
			return ti.After(tj)
		}
		return ti.Before(tj)
	})
}

// sortFilesByExtension sorts files by extension.
// Directories always come first.
func sortFilesByExtension(files []FileItem, order SortOrder) {
	sort.Slice(files, func(i, j int) bool {
		di, dj := files[i].IsDir(), files[j].IsDir()
		if di != dj {
			return di // directories first
		}
		if di && dj {
			return strings.Compare(files[i].Name(), files[j].Name()) < 0
		}

		ei, ej := files[i].Extension(), files[j].Extension()
		if ei == ej {
			return strings.Compare(files[i].Name(), files[j].Name()) < 0
		}

		if order == SortDescending {
			return strings.Compare(ej, ei) < 0
		}
		return strings.Compare(ei, ej) < 0
	})
}

// DirectoryFirstSorter is a sort helper that puts directories first.
type DirectoryFirstSorter struct {
	files  []FileItem
	less   func(i, j int) bool
}

// NewDirectoryFirstSorter creates a new sorter that always sorts directories first.
func NewDirectoryFirstSorter(files []FileItem, less func(i, j int) bool) *DirectoryFirstSorter {
	return &DirectoryFirstSorter{
		files: files,
		less:  less,
	}
}

func (s *DirectoryFirstSorter) Len() int      { return len(s.files) }
func (s *DirectoryFirstSorter) Swap(i, j int) { s.files[i], s.files[j] = s.files[j], s.files[i] }
func (s *DirectoryFirstSorter) Less(i, j int) bool {
	di, dj := s.files[i].IsDir(), s.files[j].IsDir()
	if di != dj {
		return di // directories first
	}
	return s.less(i, j)
}

// Sort implements the sort.Interface.
func (s *DirectoryFirstSorter) Sort() {
	sort.Sort(s)
}
