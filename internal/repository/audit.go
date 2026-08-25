package repository

import (
	"fmt"
	"noticeword/internal/model"
	"noticeword/internal/store"
)

type AuditRepository struct {
	storage *store.Store
}

func NewAuditRepository(storage *store.Store) *AuditRepository {
	return &AuditRepository{storage: storage}
}

func (r *AuditRepository) Append(recordID, action, actor, detail, at string, sequence int) (model.AuditEvent, error) {
	if recordID == "" || action == "" {
		return model.AuditEvent{}, fmt.Errorf("record and action are required")
	}
	event := model.AuditEvent{ID: store.AuditKey(recordID, sequence), RecordID: recordID, Action: action, Actor: actor, Detail: detail, OccurredAt: at}
	if err := r.storage.PutAudit(event); err != nil {
		return model.AuditEvent{}, err
	}
	return event, nil
}

func (r *AuditRepository) ForRecord(recordID string) ([]model.AuditEvent, error) {
	return r.storage.ListAudit(recordID)
}

func Summarize(events []model.AuditEvent) string {
	if len(events) == 0 {
		return "no activity"
	}
	last := events[len(events)-1]
	return fmt.Sprintf("%s by %s at %s", last.Action, last.Actor, last.OccurredAt)
}
