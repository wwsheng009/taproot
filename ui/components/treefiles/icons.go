package treefiles

// Tree icons for displaying expand/collapse state
const (
	// Expanded directory icons
	IconExpanded   = "📂"
	IconExpandedAlt = "▼"

	// Collapsed directory icons
	IconCollapsed   = "📁"
	IconCollapsedAlt = "▶"

	// Vertical and horizontal tree lines
	IconVertical   = "│"
	IconHorizontal = "─"
	IconCorner     = "└"
	IconTee        = "├"
	IconSpace      = " "

	// File icons (from files/icon.go)
	IconGo      = "Go"
	IconPy      = "Python"
	IconJS      = "JavaScript"
	IconTS      = "TypeScript"
	IconRS      = "Rust"
	IconC       = "C"
	IconJava    = "Java"
	IconHTML    = "HTML"
	IconCSS     = "CSS"
	IconJSON    = "JSON"
	IconMD      = "Markdown"
	IconTXT     = "Text"
	IconPNG     = "Image"
	IconJPG     = "Image"
	IconPDF     = "PDF"
	IconZIP     = "Archive"
	IconMP3     = "Audio"
	IconMP4     = "Video"
	IconSQL     = "SQL"
	IconConf    = "Config"
	IconSH      = "Shell"
	IconUnknown = "File"
)

// GetTreeIcon returns the appropriate icon for a tree node.
func GetTreeIcon(node *FileNode) string {
	if !node.IsDir() {
		return GetFileIcon(node.name)
	}

	if node.Expanded() {
		return IconExpanded
	}
	return IconCollapsed
}

// GetTreePrefix returns the appropriate tree prefix based on node position.
func GetTreePrefix(node *FileNode, isLast bool) string {
	if node.Parent() == nil {
		return "" // Root node has no prefix
	}

	var prefix string

	// Build prefix from parent chain
	current := node.Parent()
	depth := node.Depth()

	for i := 1; i < depth; i++ {
		if current.Parent() != nil {
			siblings := current.Parent().Children()
			length := len(siblings)
			if length > 0 && siblings[length-1] == current {
				prefix += IconSpace + IconSpace + IconSpace
			} else {
				prefix += IconVertical + IconSpace + IconSpace
			}
			current = current.Parent()
		} else {
			prefix += IconSpace + IconSpace + IconSpace
		}
	}

	// Add connector for current level
	_, isLastChild := isLastNode(node)
	if isLastChild {
		prefix += IconCorner + IconHorizontal
	} else {
		prefix += IconTee + IconHorizontal
	}

	return prefix
}

// isLastNode returns true if this node is the last child of its parent.
func isLastNode(node *FileNode) (bool, bool) {
	if node.Parent() == nil {
		return true, true
	}

	siblings := node.Parent().Children()
	if len(siblings) == 0 {
		return false, true
	}

	return siblings[len(siblings)-1] == node, true
}

// GetFileIcon returns an icon based on file extension.
func GetFileIcon(name string) string {
	ext := getExtension(name)

	switch ext {
	case "go":
		return "Go"
	case "py":
		return "Python"
	case "js":
		return "JavaScript"
	case "ts":
		return "TypeScript"
	case "rs":
		return "Rust"
	case "c":
		return "C"
	case "cpp", "cc", "cxx":
		return "C++"
	case "h":
		return "C Header"
	case "java":
		return "Java"
	case "class":
		return "Java Class"
	case "jar":
		return "Java Archive"
	case "html", "htm":
		return "HTML"
	case "css":
		return "CSS"
	case "json":
		return "JSON"
	case "xml":
		return "XML"
	case "yaml", "yml":
		return "YAML"
	case "md", "markdown":
		return "Markdown"
	case "txt":
		return "Text"
	case "log":
		return "Log"
	case "png", "jpg", "jpeg", "gif", "svg", "webp", "ico", "bmp":
		return "Image"
	case "pdf":
		return "PDF"
	case "doc", "docx":
		return "Word"
	case "xls", "xlsx":
		return "Excel"
	case "ppt", "pptx":
		return "PowerPoint"
	case "zip", "tar", "gz", "7z", "rar", "bz2":
		return "Archive"
	case "mp3", "wav", "ogg", "flac", "aac":
		return "Audio"
	case "mp4", "avi", "mkv", "mov", "wmv", "flv":
		return "Video"
	case "sql":
		return "SQL"
	case "db", "sqlite", "sqlite3":
		return "Database"
	case "conf", "cfg", "ini", "toml":
		return "Config"
	case "sh", "bash", "zsh", "fish", "ps1":
		return "Shell"
	case "bat", "cmd":
		return "Batch"
	case "exe", "msi", "dmg", "app":
		return "Executable"
	case "dll", "so", "dylib", "lib":
		return "Library"
	case "dockerfile":
		return "Docker"
	case "makefile", "cmake":
		return "Build"
	default:
		return "File"
	}
}

// getExtension extracts file extension without the dot.
func getExtension(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i+1:]
		}
	}
	return ""
}
