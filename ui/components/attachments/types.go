package attachments

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// AttachmentType represents the type of attachment.
type AttachmentType int

const (
	AttachmentTypeFile AttachmentType = iota
	AttachmentTypeImage
	AttachmentTypeVideo
	AttachmentTypeAudio
	AttachmentTypeDocument
	AttachmentTypeArchive
)

// String returns the string representation of the attachment type.
func (at AttachmentType) String() string {
	switch at {
	case AttachmentTypeFile:
		return "file"
	case AttachmentTypeImage:
		return "image"
	case AttachmentTypeVideo:
		return "video"
	case AttachmentTypeAudio:
		return "audio"
	case AttachmentTypeDocument:
		return "document"
	case AttachmentTypeArchive:
		return "archive"
	default:
		return "unknown"
	}
}

// Attachment represents a file attachment.
type Attachment struct {
	ID         string
	Type       AttachmentType
	Name       string
	Path       string
	Size       int64
	MimeType   string
	Thumbnail  string // Path to thumbnail for images
	Preview    string // Text preview or description
	CreatedAt  time.Time
	ModifiedAt time.Time
	Metadata   map[string]string // Additional metadata
}

// AttachmentConfig holds configuration for attachment list display.
type AttachmentConfig struct {
	ShowThumbnails    bool // Show thumbnails for images
	ShowSize          bool // Show file size
	ShowDate          bool // Show modification date
	ShowPreview       bool // Show text preview
	CompactMode       bool // Compact display mode
	MaxThumbnailWidth int  // Maximum thumbnail width
}

// DefaultAttachmentConfig returns default attachment configuration.
func DefaultAttachmentConfig() AttachmentConfig {
	return AttachmentConfig{
		ShowThumbnails:    true,
		ShowSize:          true,
		ShowDate:          false,
		ShowPreview:       true,
		CompactMode:       false,
		MaxThumbnailWidth: 40,
	}
}

// DetectAttachmentType detects the attachment type from file extension.
func DetectAttachmentType(filename string) AttachmentType {
	ext := strings.ToLower(filepath.Ext(filename))

	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".avif"}
	videoExts := []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm"}
	audioExts := []string{".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a"}
	documentExts := []string{".pdf", ".doc", ".docx", ".txt", ".md", ".rtf"}
	archiveExts := []string{".zip", ".rar", ".7z", ".tar", ".gz", ".bz2"}

	for _, e := range imageExts {
		if ext == e {
			return AttachmentTypeImage
		}
	}
	for _, e := range videoExts {
		if ext == e {
			return AttachmentTypeVideo
		}
	}
	for _, e := range audioExts {
		if ext == e {
			return AttachmentTypeAudio
		}
	}
	for _, e := range documentExts {
		if ext == e {
			return AttachmentTypeDocument
		}
	}
	for _, e := range archiveExts {
		if ext == e {
			return AttachmentTypeArchive
		}
	}

	return AttachmentTypeFile
}

// FormatSize formats a file size in human-readable format.
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return "< 1 KB"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), []string{"KB", "MB", "GB", "TB"}[exp])
}

// NewAttachment creates a new attachment from file information.
func NewAttachment(path string, info fs.FileInfo) (*Attachment, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	attType := DetectAttachmentType(info.Name())

	return &Attachment{
		ID:         filepath.Base(path),
		Type:       attType,
		Name:       info.Name(),
		Path:       absPath,
		Size:       info.Size(),
		MimeType:   detectMimeType(path),
		CreatedAt:  time.Now(),
		ModifiedAt: info.ModTime(),
		Metadata:   make(map[string]string),
	}, nil
}

// detectMimeType detects the MIME type of a file.
// This is a simplified implementation. For production, use the "mime" package.
func detectMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".md":   "text/markdown",
		".html": "text/html",
		".css":  "text/css",
		".js":   "text/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".zip":  "application/zip",
		".tar":  "application/x-tar",
		".gz":   "application/gzip",
	}

	if mt, ok := mimeTypes[ext]; ok {
		return mt
	}

	return "application/octet-stream"
}
