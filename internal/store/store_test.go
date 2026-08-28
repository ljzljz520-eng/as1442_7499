package store

import (
	"path/filepath"
	"testing"

	"noticeword/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "noticeword.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := model.NewRecord("reopen-1", "重开社群", "重开公告", "重新打开后仍可读取", "持久化说明", "2025-01-01", nil)
	if err := first.PutRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := first.PutAudit(model.AuditEvent{ID: "a1", RecordID: record.ID, Action: model.ActionRegistered, Actor: "tester", OccurredAt: "2025-01-01"}); err != nil {
		t.Fatal(err)
	}
	if err := first.PutWorkflow(model.Workflow{ID: "w1", RecordID: record.ID, Stage: "registration", OccurredAt: "2025-01-01"}); err != nil {
		t.Fatal(err)
	}
	if err := first.PutAttachment(model.Attachment{ID: "f1", RecordID: record.ID, Name: "brief.txt", Size: 4}); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Body != record.Body || loaded.Description != record.Description {
		t.Fatalf("loaded record differs: %#v", loaded)
	}
	audit, err := second.ListAudit(record.ID)
	if err != nil || len(audit) != 1 {
		t.Fatalf("audit after reopen: %v %#v", err, audit)
	}
	workflows, err := second.ListWorkflows(record.ID)
	if err != nil || len(workflows) != 1 {
		t.Fatalf("workflows after reopen: %v %#v", err, workflows)
	}
	attachments, err := second.ListAttachments(record.ID)
	if err != nil || len(attachments) != 1 {
		t.Fatalf("attachments after reopen: %v %#v", err, attachments)
	}
}

func TestStoreMissingRecord(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "missing.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.GetRecord("missing"); err == nil {
		t.Fatal("expected missing record error")
	}
}
