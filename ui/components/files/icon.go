package files

// GetIcon returns an icon for a file/directory based on its type.
func GetIcon(isDir bool, extension string) string {
	if isDir {
		return "📁"
	}

	// File extension icons
	icons := map[string]string{
		// Code files
		"go":     "🐹",
		"py":     "🐍",
		"js":     "📜",
		"ts":     "📘",
		"tsx":    "⚛️",
		"jsx":    "⚛️",
		"rs":     "🦀",
		"c":      "⚙️",
		"cpp":    "⚙️",
		"h":      "📋",
		"hpp":    "📋",
		"java":   "☕",
		"kt":     "🎯",
		"swift":  "🍎",
		"rb":     "💎",
		"php":    "🐘",
		"scala":  "🔷",
		"dart":   "🎯",
		"lua":    "🌙",
		"r":      "📊",

		// Web files
		"html":   "🌐",
		"css":    "🎨",
		"scss":   "🎨",
		"sass":   "🎨",
		"less":   "🎨",
		"json":   "📋",
		"xml":    "📋",
		"yaml":   "📋",
		"yml":    "📋",
		"toml":   "📋",

		// Documentation
		"md":     "📝",
		"txt":    "📄",
		"rst":    "📝",
		"adoc":   "📝",

		// Images
		"png":    "🖼️",
		"jpg":    "🖼️",
		"jpeg":   "🖼️",
		"gif":    "🖼️",
		"svg":    "🎨",
		"ico":    "🖼️",
		"bmp":    "🖼️",
		"webp":   "🖼️",

		// Audio
		"mp3":    "🎵",
		"wav":    "🎵",
		"flac":   "🎵",
		"ogg":    "🎵",
		"aac":    "🎵",
		"m4a":    "🎵",

		// Video
		"mp4":    "🎬",
		"avi":    "🎬",
		"mkv":    "🎬",
		"mov":    "🎬",
		"webm":   "🎬",
		"flv":    "🎬",

		// Archives
		"zip":    "📦",
		"tar":    "📦",
		"gz":     "📦",
		"rar":    "📦",
		"7z":     "📦",
		"bz2":    "📦",
		"xz":     "📦",

		// Documents
		"pdf":    "📕",
		"doc":    "📘",
		"docx":   "📘",
		"xls":    "📗",
		"xlsx":   "📗",
		"ppt":    "📙",
		"pptx":   "📙",

		// Databases
		"sql":    "🗄️",
		"db":     "🗄️",
		"sqlite": "🗄️",
		"sqlite3": "🗄️",

		// Config files
		"conf":   "⚙️",
		"config": "⚙️",
		"ini":    "⚙️",
		"cfg":    "⚙️",
		"env":    "🔐",
		"dotenv": "🔐",

		// Shell/Scripts
		"sh":     "🐚",
		"bash":   "🐚",
		"zsh":    "🐚",
		"fish":   "🐚",
		"bat":    "🦇",
		"cmd":    "🦇",
		"ps1":    "💠",

		// Build files
		"makefile": "🛠️",
		"dockerfile": "🐳",
		"docker-compose": "🐳",
		"mk":     "🛠️",
		"gradle": "🐘",
		"pom.xml": "🐘",

		// Version control
		"git":   "🔀",
		"gitignore": "🔀",
		"gitattributes": "🔀",
		"gitmodules": "🔀",

		// Lock files
		"lock":   "🔒",
	}

	if icon, ok := icons[extension]; ok {
		return icon
	}

	// Default file icon
	return "📄"
}
