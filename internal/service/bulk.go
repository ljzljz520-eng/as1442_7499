package service

import (
	"fmt"
	"strings"

	"noticeword/internal/attachment"
	"noticeword/internal/export"
	"noticeword/internal/model"
	"noticeword/internal/policy"
	"noticeword/internal/store"
)

func (s *Service) ImportCSV(input, actor string) model.ImportReport {
	rows, err := export.Decode(input)
	if err != nil {
		return model.ImportReport{Rejected: 1, Errors: []string{err.Error()}}
	}
	return s.Import(rows, actor)
}

func (s *Service) ExportCSV(query model.SearchQuery) (string, error) {
	records, err := s.Export(query)
	if err != nil {
		return "", err
	}
	return export.Encode(records)
}

func (s *Service) AddAttachment(recordID, name, contentType string, size int, checksum string) (model.Attachment, error) {
	record, err := s.Get(recordID)
	if err != nil {
		return model.Attachment{}, err
	}
	if record.IsClosed() {
		return model.Attachment{}, fmt.Errorf("archived record cannot accept attachment")
	}
	manager := attachment.NewManager(s.db())
	item, err := manager.Add(recordID, name, contentType, size, checksum)
	if err != nil {
		return model.Attachment{}, err
	}
	if err := manager.Validate(item); err != nil {
		return model.Attachment{}, err
	}
	return item, nil
}

func (s *Service) AttachmentSize(recordID string) (int, error) {
	return attachment.NewManager(s.db()).TotalSize(recordID)
}

func (s *Service) db() *store.Store {
	return s.storage
}

func normalizeImportActor(actor string) string {
	if strings.TrimSpace(actor) == "" {
		return "system-import"
	}
	return strings.TrimSpace(actor)
}

func normalizeTagsForRecord(tags []string) []string {
	return policy.NormalizeTags(tags)
}
