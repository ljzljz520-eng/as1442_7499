package repository

import (
	"fmt"
	"sort"
	"strings"

	"noticeword/internal/model"
	"noticeword/internal/store"
)

type RecordRepository struct {
	storage *store.Store
}

func NewRecordRepository(storage *store.Store) *RecordRepository {
	return &RecordRepository{storage: storage}
}

func (r *RecordRepository) Save(record model.Record) error {
	if r.storage == nil {
		return fmt.Errorf("record repository is not configured")
	}
	return r.storage.PutRecord(record)
}

func (r *RecordRepository) Find(id string) (model.Record, error) {
	return r.storage.GetRecord(id)
}

func (r *RecordRepository) All() ([]model.Record, error) {
	records, err := r.storage.ListRecords()
	if err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].UpdatedAt == records[j].UpdatedAt {
			return records[i].ID < records[j].ID
		}
		return records[i].UpdatedAt < records[j].UpdatedAt
	})
	return records, nil
}

func (r *RecordRepository) Search(query model.SearchQuery) (model.SearchResult, error) {
	query = query.Normalized()
	records, err := r.All()
	if err != nil {
		return model.SearchResult{}, err
	}
	filtered := make([]model.Record, 0, len(records))
	for _, record := range records {
		if query.Community != "" && !strings.Contains(strings.ToLower(record.Community), query.Community) {
			continue
		}
		if query.Text != "" && !strings.Contains(record.SearchText(), query.Text) {
			continue
		}
		if query.Status != "" && record.Status != query.Status {
			continue
		}
		filtered = append(filtered, record.Clone())
	}
	total := len(filtered)
	start := (query.Page - 1) * query.PageSize
	if start > total {
		start = total
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	items := append([]model.Record(nil), filtered[start:end]...)
	pages := 0
	if total > 0 {
		pages = (total + query.PageSize - 1) / query.PageSize
	}
	return model.SearchResult{Items: items, Page: query.Page, PageSize: query.PageSize, Total: total, TotalPages: pages}, nil
}

func (r *RecordRepository) SaveAudit(event model.AuditEvent) error {
	return r.storage.PutAudit(event)
}

func (r *RecordRepository) Audit(recordID string) ([]model.AuditEvent, error) {
	return r.storage.ListAudit(recordID)
}

func (r *RecordRepository) SaveWorkflow(workflow model.Workflow) error {
	return r.storage.PutWorkflow(workflow)
}

func (r *RecordRepository) Workflows(recordID string) ([]model.Workflow, error) {
	return r.storage.ListWorkflows(recordID)
}

func (r *RecordRepository) SaveAttachment(attachment model.Attachment) error {
	return r.storage.PutAttachment(attachment)
}

func (r *RecordRepository) Attachments(recordID string) ([]model.Attachment, error) {
	return r.storage.ListAttachments(recordID)
}
