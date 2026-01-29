package treefiles

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileNode represents a node in the file tree.
type FileNode struct {
	path     string
	name     string
	isDir    bool
	size     int64
	modTime  string
	expanded bool
	children []*FileNode
	parent   *FileNode
	depth    int
}

// Path returns the full path of the node.
func (n *FileNode) Path() string {
	return n.path
}

// Name returns the display name of the node.
func (n *FileNode) Name() string {
	return n.name
}

// IsDir returns true if this is a directory.
func (n *FileNode) IsDir() bool {
	return n.isDir
}

// Size returns the file size (0 for directories).
func (n *FileNode) Size() int64 {
	if !n.isDir {
		return n.size
	}
	return 0
}

// ModTime returns the modification time string.
func (n *FileNode) ModTime() string {
	return n.modTime
}

// Expanded returns true if the directory is expanded.
func (n *FileNode) Expanded() bool {
	return n.expanded
}

// SetExpanded sets the expanded state.
func (n *FileNode) SetExpanded(expanded bool) {
	n.expanded = expanded
}

// Toggle expands or collapses the directory.
func (n *FileNode) Toggle() {
	n.expanded = !n.expanded
}

// Children returns the child nodes.
func (n *FileNode) Children() []*FileNode {
	return n.children
}

// Parent returns the parent node.
func (n *FileNode) Parent() *FileNode {
	return n.parent
}

// Depth returns the depth level in the tree.
func (n *FileNode) Depth() int {
	return n.depth
}

// HasChildren returns true if the node has children.
func (n *FileNode) HasChildren() bool {
	return n.isDir && len(n.children) > 0
}

// AddChild adds a child node.
func (n *FileNode) AddChild(child *FileNode) {
	child.parent = n
	child.depth = n.depth + 1
	n.children = append(n.children, child)
}

// RemoveChild removes a child node.
func (n *FileNode) RemoveChild(child *FileNode) {
	for i, c := range n.children {
		if c == child {
			n.children = append(n.children[:i], n.children[i+1:]...)
			break
		}
	}
}

// FindChild finds a child by name.
func (n *FileNode) FindChild(name string) *FileNode {
	for _, child := range n.children {
		if child.name == name {
			return child
		}
	}
	return nil
}

// FileTree represents a tree structure of files and directories.
type FileTree struct {
	root      *FileNode
	expanded  map[string]bool
	sortBy    SortBy
	sortOrder SortOrder
	hidden    bool
	maxDepth  int
}

// FileOption is a function that configures a FileTree.
type FileOption func(*FileTree)

// WithSort sets the sort criteria and order.
func WithSort(sortBy SortBy, sortOrder SortOrder) FileOption {
	return func(ft *FileTree) {
		ft.sortBy = sortBy
		ft.sortOrder = sortOrder
	}
}

// WithHidden sets whether to show hidden files.
func WithHidden(show bool) FileOption {
	return func(ft *FileTree) {
		ft.hidden = show
	}
}

// WithMaxDepth sets the maximum depth to scan.
func WithMaxDepth(depth int) FileOption {
	return func(ft *FileTree) {
		ft.maxDepth = depth
	}
}

// SortBy defines sorting criteria.
type SortBy int

const (
	SortByName SortBy = iota
	SortBySize
	SortByTime
	SortByType
)

// SortOrder defines the sort direction.
type SortOrder int

const (
	SortAscending SortOrder = iota
	SortDescending
)

// NewFileTree creates a new file tree rooted at the given path.
func NewFileTree(rootPath string, opts ...FileOption) (*FileTree, error) {
	rootPath = filepath.Clean(rootPath)

	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	root := &FileNode{
		path:    rootPath,
		name:    filepath.Base(rootPath),
		isDir:   info.IsDir(),
		size:    info.Size(),
		modTime: info.ModTime().Format("Jan 02 15:04"),
		depth:   0,
	}

	tree := &FileTree{
		root:      root,
		expanded:  make(map[string]bool),
		sortBy:    SortByName,
		sortOrder: SortAscending,
		hidden:    false,
		maxDepth:  100, // Default max depth
	}

	// Apply options
	for _, opt := range opts {
		opt(tree)
	}

	// Set root as expanded by default
	tree.expanded[rootPath] = true
	root.expanded = true

	// Build the tree
	if root.isDir {
		if err := tree.buildTree(root); err != nil {
			return nil, err
		}
	}

	return tree, nil
}

// buildTree recursively builds the file tree.
func (ft *FileTree) buildTree(node *FileNode) error {
	if !node.isDir {
		return nil
	}

	if node.depth >= ft.maxDepth {
		return nil
	}

	entries, err := os.ReadDir(node.path)
	if err != nil {
		return err
	}

	var children []*FileNode
	for _, entry := range entries {
		// Skip hidden files if not showing hidden
		if !ft.hidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(node.path, entry.Name())

		child := &FileNode{
			path:    fullPath,
			name:    entry.Name(),
			isDir:   info.IsDir(),
			size:    info.Size(),
			modTime: info.ModTime().Format("Jan 02 15:04"),
			depth:   node.depth + 1,
		}

		// Only recurse into directories
		if child.isDir && node.depth+1 < ft.maxDepth {
			ft.buildTree(child)
		}

		children = append(children, child)
		node.AddChild(child)
	}

	// Sort children
	ft.sortChildren(children)

	return nil
}

// sortChildren sorts a slice of child nodes.
func (ft *FileTree) sortChildren(children []*FileNode) {
	switch ft.sortBy {
	case SortByName:
		sort.Slice(children, func(i, j int) bool {
			di, dj := children[i].isDir, children[j].isDir
			if di != dj {
				return di // directories first
			}
			if ft.sortOrder == SortDescending {
				i, j = j, i
			}
			return strings.Compare(children[i].name, children[j].name) < 0
		})

	case SortByType:
		sort.Slice(children, func(i, j int) bool {
			di, dj := children[i].isDir, children[j].isDir
			if di != dj {
				return di // directories first
			}
			if di && dj {
				if ft.sortOrder == SortDescending {
					i, j = j, i
				}
				return strings.Compare(children[i].name, children[j].name) < 0
			}

			// Extract extensions
			getExt := func(n *FileNode) string {
				if n.isDir {
					return ""
				}
				ext := filepath.Ext(n.name)
				if len(ext) > 0 && ext[0] == '.' {
					ext = ext[1:]
				}
				return ext
			}

			ei, ej := getExt(children[i]), getExt(children[j])
			if ei == ej {
				if ft.sortOrder == SortDescending {
					i, j = j, i
				}
				return strings.Compare(children[i].name, children[j].name) < 0
			}

			if ft.sortOrder == SortDescending {
				return strings.Compare(ej, ei) < 0
			}
			return strings.Compare(ei, ej) < 0
		})

	case SortBySize, SortByTime:
		// For files, sort by size/time
		sort.Slice(children, func(i, j int) bool {
			di, dj := children[i].isDir, children[j].isDir
			if di != dj {
				return di // directories first
			}
			if di && dj {
				if ft.sortOrder == SortDescending {
					i, j = j, i
				}
				return strings.Compare(children[i].name, children[j].name) < 0
			}

			var compare bool
			if ft.sortBy == SortBySize {
				if ft.sortOrder == SortDescending {
					compare = children[i].size > children[j].size
				} else {
					compare = children[i].size < children[j].size
				}
			} else { // SortByTime - use string comparison
				ti, tj := children[i].modTime, children[j].modTime
				if ft.sortOrder == SortDescending {
					compare = strings.Compare(tj, ti) < 0
				} else {
					compare = strings.Compare(ti, tj) < 0
				}
			}
			return compare
		})
	}
}

// Flatten returns a flattened list of visible nodes.
// Only includes expanded directories' children.
func (ft *FileTree) Flatten() []*FileNode {
	return ft.flattenNode(ft.root)
}

// flattenNode recursively flattens a node.
func (ft *FileTree) flattenNode(node *FileNode) []*FileNode {
	var result []*FileNode

	result = append(result, node)

	if node.isDir && node.expanded {
		for _, child := range node.children {
			result = append(result, ft.flattenNode(child)...)
		}
	}

	return result
}

// Root returns the root node.
func (ft *FileTree) Root() *FileNode {
	return ft.root
}

// FindNode finds a node by its path.
func (ft *FileTree) FindNode(path string) *FileNode {
	return ft.findNode(ft.root, filepath.Clean(path))
}

// findNode recursively finds a node by path.
func (ft *FileTree) findNode(node *FileNode, path string) *FileNode {
	if node.path == path {
		return node
	}

	if node.isDir {
		for _, child := range node.children {
			if found := ft.findNode(child, path); found != nil {
				return found
			}
		}
	}

	return nil
}

// ToggleNode toggles the expanded state of a node.
func (ft *FileTree) ToggleNode(node *FileNode) {
	if node.isDir {
		node.Toggle()
		ft.expanded[node.path] = node.expanded
	}
}

// ExpandAll expands all directories.
func (ft *FileTree) ExpandAll() {
	ft.expandNode(ft.root)
}

// expandNode recursively expands a node.
func (ft *FileTree) expandNode(node *FileNode) {
	if node.isDir {
		node.expanded = true
		ft.expanded[node.path] = true
		for _, child := range node.children {
			ft.expandNode(child)
		}
	}
}

// CollapseAll collapses all directories.
func (ft *FileTree) CollapseAll() {
	ft.collapseNode(ft.root)
	// Keep root expanded
	ft.root.expanded = true
	ft.expanded[ft.root.path] = true
}

// collapseNode recursively collapses a node.
func (ft *FileTree) collapseNode(node *FileNode) {
	if node.isDir {
		node.expanded = false
		ft.expanded[node.path] = false
		for _, child := range node.children {
			ft.collapseNode(child)
		}
	}
}

// Rescan rescans the tree from the root.
func (ft *FileTree) Rescan() error {
	// Clear expanded state except root
	ft.expanded = make(map[string]bool)
	ft.expanded[ft.root.path] = true
	ft.root.expanded = true
	ft.root.children = nil

	return ft.buildTree(ft.root)
}

// SetSort sets the sort criteria and order.
func (ft *FileTree) SetSort(sortBy SortBy, sortOrder SortOrder) {
	ft.sortBy = sortBy
	ft.sortOrder = sortOrder
	ft.resortAll()
}

// resortAll recursively sorts all nodes.
func (ft *FileTree) resortAll() {
	ft.sortChildren(ft.root.children)
	for _, child := range ft.root.children {
		ft.resortNode(child)
	}
}

// resortNode recursively sorts a node's children.
func (ft *FileTree) resortNode(node *FileNode) {
	if node.isDir {
		ft.sortChildren(node.children)
		for _, child := range node.children {
			ft.resortNode(child)
		}
	}
}

// ToggleHidden toggles hidden file visibility.
func (ft *FileTree) ToggleHidden() {
	ft.hidden = !ft.hidden
	ft.Rescan()
}

// IncludesHidden returns true if hidden files are shown.
func (ft *FileTree) IncludesHidden() bool {
	return ft.hidden
}

// GetStats returns statistics about the tree.
type TreeStats struct {
	TotalNodes   int
	TotalFiles   int
	TotalDirs    int
	TotalSize    int64
	ExpandedDirs int
}

// Stats returns statistics about the tree.
func (ft *FileTree) Stats() TreeStats {
	stats := TreeStats{}
	ft.countNodes(ft.root, &stats)
	return stats
}

// countNodes recursively counts nodes.
func (ft *FileTree) countNodes(node *FileNode, stats *TreeStats) {
	stats.TotalNodes++
	if node.isDir {
		stats.TotalDirs++
		if node.expanded {
			stats.ExpandedDirs++
		}
	} else {
		stats.TotalFiles++
		stats.TotalSize += node.size
	}

	for _, child := range node.children {
		ft.countNodes(child, stats)
	}
}
