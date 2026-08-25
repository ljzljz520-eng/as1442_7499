package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"go.etcd.io/bbolt"
	"noticeword/internal/model"
)

var (
	BucketRecords     = []byte("records")
	BucketAuditEvents = []byte("audit_events")
	BucketWorkflows   = []byte("workflows")
	BucketAttachments = []byte("attachments")
)

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{NoSync: true})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	s := &Store{db: db}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{BucketRecords, BucketAuditEvents, BucketWorkflows, BucketAttachments} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return fmt.Errorf("create bucket %s: %w", bucket, err)
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) PutRecord(record model.Record) error {
	return putJSON(s.db, BucketRecords, record.ID, record)
}

func (s *Store) GetRecord(id string) (model.Record, error) {
	var record model.Record
	err := getJSON(s.db, BucketRecords, id, &record)
	return record, err
}

func (s *Store) DeleteRecord(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(BucketRecords).Delete([]byte(id)) })
}

func (s *Store) PutAudit(event model.AuditEvent) error {
	return putJSON(s.db, BucketAuditEvents, event.ID, event)
}

func (s *Store) ListAudit(recordID string) ([]model.AuditEvent, error) {
	var events []model.AuditEvent
	err := listJSON(s.db, BucketAuditEvents, func(event model.AuditEvent) bool { return recordID == "" || event.RecordID == recordID }, &events)
	return events, err
}

func (s *Store) PutWorkflow(workflow model.Workflow) error {
	return putJSON(s.db, BucketWorkflows, workflow.ID, workflow)
}

func (s *Store) ListWorkflows(recordID string) ([]model.Workflow, error) {
	var workflows []model.Workflow
	err := listJSON(s.db, BucketWorkflows, func(workflow model.Workflow) bool { return recordID == "" || workflow.RecordID == recordID }, &workflows)
	return workflows, err
}

func (s *Store) PutAttachment(attachment model.Attachment) error {
	return putJSON(s.db, BucketAttachments, attachment.ID, attachment)
}

func (s *Store) ListAttachments(recordID string) ([]model.Attachment, error) {
	var attachments []model.Attachment
	err := listJSON(s.db, BucketAttachments, func(attachment model.Attachment) bool { return recordID == "" || attachment.RecordID == recordID }, &attachments)
	return attachments, err
}

func putJSON[T any](db *bbolt.DB, bucket []byte, key string, value T) error {
	if key == "" {
		return errors.New("key is required")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", key, err)
	}
	return db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Put([]byte(key), data) })
}

func getJSON[T any](db *bbolt.DB, bucket []byte, key string, target *T) error {
	return db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucket).Get([]byte(key))
		if value == nil {
			return fmt.Errorf("%s %q not found", bucket, key)
		}
		return json.Unmarshal(value, target)
	})
}

func listJSON[T any](db *bbolt.DB, bucket []byte, keep func(T) bool, output *[]T) error {
	return db.View(func(tx *bbolt.Tx) error {
		items := tx.Bucket(bucket)
		return items.ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var item T
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}
			if keep(item) {
				*output = append(*output, item)
			}
			return nil
		})
	})
}

func (s *Store) ListRecords() ([]model.Record, error) {
	var records []model.Record
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(BucketRecords).ForEach(func(_, value []byte) error {
			if value == nil {
				return nil
			}
			var record model.Record
			if err := json.Unmarshal(value, &record); err != nil {
				return err
			}
			records = append(records, record)
			return nil
		})
	})
	return records, err
}
