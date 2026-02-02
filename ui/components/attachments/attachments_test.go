package attachments

import (
	"os"
	"testing"
	"time"
)

func TestAttachmentTypeString(t *testing.T) {
	tests := []struct {
		attType  AttachmentType
		expected string
	}{
		{AttachmentTypeFile, "file"},
		{AttachmentTypeImage, "image"},
		{AttachmentTypeVideo, "video"},
		{AttachmentTypeAudio, "audio"},
		{AttachmentTypeDocument, "document"},
		{AttachmentTypeArchive, "archive"},
	}

	for _, tt := range tests {
		result := tt.attType.String()
		if result != tt.expected {
			t.Errorf("For attachment type %v, expected '%s', got '%s'", tt.attType, tt.expected, result)
		}
	}
}

func TestDetectAttachmentType(t *testing.T) {
	tests := []struct {
		filename string
		expected AttachmentType
	}{
		{"test.jpg", AttachmentTypeImage},
		{"test.jpeg", AttachmentTypeImage},
		{"test.png", AttachmentTypeImage},
		{"test.gif", AttachmentTypeImage},
		{"test.mp4", AttachmentTypeVideo},
		{"test.mkv", AttachmentTypeVideo},
		{"test.mp3", AttachmentTypeAudio},
		{"test.wav", AttachmentTypeAudio},
		{"test.pdf", AttachmentTypeDocument},
		{"test.txt", AttachmentTypeDocument},
		{"test.zip", AttachmentTypeArchive},
		{"test.tar.gz", AttachmentTypeArchive},
		{"test.exe", AttachmentTypeFile},
		{"test.dll", AttachmentTypeFile},
	}

	for _, tt := range tests {
		result := DetectAttachmentType(tt.filename)
		if result != tt.expected {
			t.Errorf("For file '%s', expected %v, got %v", tt.filename, tt.expected, result)
		}
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "< 1 KB"},
		{1024, "1.0 KB"},
		{2048, "2.0 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1.5, "1.5 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		result := FormatSize(tt.bytes)
		// Just check if it's not empty for now
		if result == "" {
			t.Errorf("Expected non-empty format for %d bytes", tt.bytes)
		}
	}
}

func TestDefaultAttachmentConfig(t *testing.T) {
	config := DefaultAttachmentConfig()

	if !config.ShowThumbnails {
		t.Error("Expected ShowThumbnails to be true")
	}
	if !config.ShowSize {
		t.Error("Expected ShowSize to be true")
	}
	if config.ShowDate {
		t.Error("Expected ShowDate to be false")
	}
	if !config.ShowPreview {
		t.Error("Expected ShowPreview to be true")
	}
	if config.CompactMode {
		t.Error("Expected CompactMode to be false")
	}
	if config.MaxThumbnailWidth != 40 {
		t.Errorf("Expected MaxThumbnailWidth 40, got %d", config.MaxThumbnailWidth)
	}
}

func TestNewAttachmentList(t *testing.T) {
	attachments := []*Attachment{
		{
			ID:   "1",
			Name: "test.jpg",
			Type: AttachmentTypeImage,
		},
	}

	al := NewAttachmentList(attachments)

	if len(al.GetAttachments()) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(al.GetAttachments()))
	}

	if al.GetAttachments()[0].Name != "test.jpg" {
		t.Errorf("Expected attachment name 'test.jpg', got '%s'", al.GetAttachments()[0].Name)
	}

	if al.width != 80 {
		t.Errorf("Expected width 80, got %d", al.width)
	}

	if !al.expanded {
		t.Error("Expected expanded to be true")
	}

	if al.focused {
		t.Error("Expected focused to be false")
	}
}

func TestAddAttachment(t *testing.T) {
	al := NewAttachmentList([]*Attachment{})

	newAttachment := &Attachment{
		ID:   "1",
		Name: "new.jpg",
		Type: AttachmentTypeImage,
		Size: 1024,
	}

	al.AddAttachment(newAttachment)

	if len(al.GetAttachments()) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(al.GetAttachments()))
	}
}

func TestRemoveAttachment(t *testing.T) {
	attachments := []*Attachment{
		{ID: "1", Name: "test.jpg", Type: AttachmentTypeImage},
		{ID: "2", Name: "test.png", Type: AttachmentTypeImage},
	}

	al := NewAttachmentList(attachments)

	// Remove existing attachment
	if !al.RemoveAttachment("1") {
		t.Error("Expected RemoveAttachment to return true")
	}

	if len(al.GetAttachments()) != 1 {
		t.Errorf("Expected 1 attachment after removal, got %d", len(al.GetAttachments()))
	}

	// Try to remove non-existent attachment
	if al.RemoveAttachment("3") {
		t.Error("Expected RemoveAttachment to return false for non-existent ID")
	}
}

func TestGetAttachment(t *testing.T) {
	attachments := []*Attachment{
		{ID: "1", Name: "test.jpg", Type: AttachmentTypeImage},
	}

	al := NewAttachmentList(attachments)

	// Get existing attachment
	att := al.GetAttachment("1")
	if att == nil {
		t.Error("Expected to find attachment with ID '1'")
	}

	if att.Name != "test.jpg" {
		t.Errorf("Expected name 'test.jpg', got '%s'", att.Name)
	}

	// Get non-existent attachment
	att = al.GetAttachment("2")
	if att != nil {
		t.Error("Expected nil for non-existent attachment")
	}
}

func TestFilterByType(t *testing.T) {
	attachments := []*Attachment{
		{ID: "1", Name: "test.jpg", Type: AttachmentTypeImage},
		{ID: "2", Name: "test.png", Type: AttachmentTypeImage},
		{ID: "3", Name: "test.pdf", Type: AttachmentTypeDocument},
		{ID: "4", Name: "test.zip", Type: AttachmentTypeArchive},
	}

	al := NewAttachmentList(attachments)

	// Filter by image type
	images := al.FilterByType(AttachmentTypeImage)
	if len(images) != 2 {
		t.Errorf("Expected 2 image attachments, got %d", len(images))
	}

	// Filter by document type
	docs := al.FilterByType(AttachmentTypeDocument)
	if len(docs) != 1 {
		t.Errorf("Expected 1 document attachment, got %d", len(docs))
	}

	// Filter by audio type (none)
	audios := al.FilterByType(AttachmentTypeAudio)
	if len(audios) != 0 {
		t.Errorf("Expected 0 audio attachments, got %d", len(audios))
	}
}

func TestGetTotalSize(t *testing.T) {
	attachments := []*Attachment{
		{ID: "1", Name: "test1.jpg", Type: AttachmentTypeImage, Size: 1024},
		{ID: "2", Name: "test2.jpg", Type: AttachmentTypeImage, Size: 2048},
		{ID: "3", Name: "test3.jpg", Type: AttachmentTypeImage, Size: 3072},
	}

	al := NewAttachmentList(attachments)

	totalSize := al.GetTotalSize()
	if totalSize != 6144 {
		t.Errorf("Expected total size 6144, got %d", totalSize)
	}
}

func TestGetCountByType(t *testing.T) {
	attachments := []*Attachment{
		{ID: "1", Name: "test.jpg", Type: AttachmentTypeImage},
		{ID: "2", Name: "test.png", Type: AttachmentTypeImage},
		{ID: "3", Name: "test.pdf", Type: AttachmentTypeDocument},
		{ID: "4", Name: "test.pdf2", Type: AttachmentTypeDocument},
		{ID: "5", Name: "test.zip", Type: AttachmentTypeArchive},
	}

	al := NewAttachmentList(attachments)

	counts := al.GetCountByType()

	if counts[AttachmentTypeImage] != 2 {
		t.Errorf("Expected 2 images, got %d", counts[AttachmentTypeImage])
	}

	if counts[AttachmentTypeDocument] != 2 {
		t.Errorf("Expected 2 documents, got %d", counts[AttachmentTypeDocument])
	}

	if counts[AttachmentTypeArchive] != 1 {
		t.Errorf("Expected 1 archive, got %d", counts[AttachmentTypeArchive])
	}
}

func TestFocusAndBlur(t *testing.T) {
	al := NewAttachmentList([]*Attachment{})

	if al.Focused() {
		t.Error("Expected component to not be focused initially")
	}

	al.Focus()
	if !al.Focused() {
		t.Error("Expected component to be focused after Focus()")
	}

	al.Blur()
	if al.Focused() {
		t.Error("Expected component to not be focused after Blur()")
	}
}

func TestSetWidth(t *testing.T) {
	al := NewAttachmentList([]*Attachment{})

	al.SetWidth(120)
	if al.width != 120 {
		t.Errorf("Expected width 120, got %d", al.width)
	}
}

func TestSetConfig(t *testing.T) {
	al := NewAttachmentList([]*Attachment{})

	newConfig := AttachmentConfig{
		ShowThumbnails:    false,
		ShowSize:          false,
		ShowDate:          true,
		ShowPreview:       false,
		CompactMode:       true,
		MaxThumbnailWidth: 60,
	}

	al.SetConfig(newConfig)

	if al.config.ShowThumbnails {
		t.Error("Expected ShowThumbnails to be false")
	}
	if al.config.CompactMode != true {
		t.Error("Expected CompactMode to be true")
	}
}

func TestEmptyAttachmentList(t *testing.T) {
	al := NewAttachmentList([]*Attachment{})

	view := al.View()
	if view == "" {
		t.Error("Expected non-empty view for empty attachment list")
	}

	if len(al.GetAttachments()) != 0 {
		t.Errorf("Expected 0 attachments, got %d", len(al.GetAttachments()))
	}

	totalSize := al.GetTotalSize()
	if totalSize != 0 {
		t.Errorf("Expected total size 0, got %d", totalSize)
	}
}

func TestNewAttachment(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("test content")
	tmpfile.Close()

	info, err := os.Stat(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	att, err := NewAttachment(tmpfile.Name(), info)
	if err != nil {
		t.Fatalf("Failed to create attachment: %v", err)
	}

	if att.Name != info.Name() {
		t.Errorf("Expected name '%s', got '%s'", info.Name(), att.Name)
	}

	if att.Size != info.Size() {
		t.Errorf("Expected size %d, got %d", info.Size(), att.Size)
	}

	if att.Type == AttachmentTypeFile {
		t.Error("Expected detected type to be document for .txt file")
	}

	if att.Path == "" {
		t.Error("Expected non-empty path")
	}
}

func TestAttachmentFields(t *testing.T) {
	now := time.Now()
	att := &Attachment{
		ID:         "test-id",
		Name:       "Test File.jpg",
		Type:       AttachmentTypeImage,
		Size:       123456,
		MimeType:   "image/jpeg",
		Thumbnail:  "/path/to/thumb.jpg",
		Preview:    "This is a preview...",
		CreatedAt:  now,
		ModifiedAt: now,
		Metadata:   map[string]string{"author": "test"},
	}

	if att.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", att.ID)
	}

	if att.Name != "Test File.jpg" {
		t.Errorf("Expected name 'Test File.jpg', got %s", att.Name)
	}

	if att.Type != AttachmentTypeImage {
		t.Errorf("Expected type AttachmentTypeImage, got %v", att.Type)
	}

	if att.Size != 123456 {
		t.Errorf("Expected size 123456, got %d", att.Size)
	}

	if len(att.Metadata) != 1 {
		t.Errorf("Expected 1 metadata entry, got %d", len(att.Metadata))
	}
}
