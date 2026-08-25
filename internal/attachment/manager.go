package attachment

import (
	"fmt"
	"strings"

	"noticeword/internal/model"
	"noticeword/internal/store"
)

type Manager struct {
	storage *store.Store
}

func NewManager(storage *store.Store) *Manager { return &Manager{storage: storage} }

func (m *Manager) Add(recordID, name, contentType string, size int, checksum string) (model.Attachment, error) {
	if strings.TrimSpace(recordID) == "" {
		return model.Attachment{}, fmt.Errorf("record id is required")
	}
	if strings.TrimSpace(name) == "" {
		return model.Attachment{}, fmt.Errorf("attachment name is required")
	}
	if size < 0 {
		return model.Attachment{}, fmt.Errorf("attachment size cannot be negative")
	}
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}
	item := model.Attachment{ID: store.AttachmentKey(recordID, name), RecordID: recordID, Name: name, ContentType: contentType, Size: size, Checksum: checksum}
	if err := m.storage.PutAttachment(item); err != nil {
		return model.Attachment{}, err
	}
	return item, nil
}

func (m *Manager) List(recordID string) ([]model.Attachment, error) {
	return m.storage.ListAttachments(recordID)
}

func (m *Manager) TotalSize(recordID string) (int, error) {
	items, err := m.List(recordID)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, item := range items {
		total += item.Size
	}
	return total, nil
}

func (m *Manager) Validate(item model.Attachment) error {
	if item.RecordID == "" || item.Name == "" {
		return fmt.Errorf("attachment identity is incomplete")
	}
	if item.Size > 5*1024*1024 {
		return fmt.Errorf("attachment exceeds 5MB")
	}
	if !strings.Contains(item.ContentType, "/") {
		return fmt.Errorf("attachment content type is invalid")
	}
	return nil
}
