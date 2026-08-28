package service

import (
	"errors"
	"fmt"
	"strings"

	"noticeword/internal/model"
	"noticeword/internal/policy"
	"noticeword/internal/repository"
	"noticeword/internal/store"
)

type Service struct {
	storage   *store.Store
	records   *repository.RecordRepository
	audit     *repository.AuditRepository
	workflows *repository.WorkflowRepository
	clock     store.Clock
	sequence  int
}

func New(storage *store.Store, clock store.Clock) *Service {
	return &Service{storage: storage, records: repository.NewRecordRepository(storage), audit: repository.NewAuditRepository(storage), workflows: repository.NewWorkflowRepository(storage), clock: clock}
}

func (s *Service) nextSequence() int {
	s.sequence++
	return s.sequence
}

func (s *Service) Register(command model.RegisterCommand) (model.Record, error) {
	decision := policy.EvaluateRegistration(command)
	if !decision.Allowed {
		return model.Record{}, fmt.Errorf("registration denied: %s", policy.JoinReasons(decision.Reasons))
	}
	if strings.TrimSpace(command.Actor) == "" {
		return model.Record{}, errors.New("actor is required")
	}
	command.Tags = policy.NormalizeTags(command.Tags)
	created := model.NewRecord(command.ID, command.Community, command.Title, command.Body, command.Description, s.clock.Now(), command.Tags)
	if err := created.Validate(); err != nil {
		return model.Record{}, err
	}
	if _, err := s.records.Find(created.ID); err == nil {
		return model.Record{}, fmt.Errorf("record %q already exists", created.ID)
	}
	if err := s.records.Save(created); err != nil {
		return model.Record{}, err
	}
	if _, err := s.audit.Append(created.ID, model.ActionRegistered, command.Actor, "registered notice", created.CreatedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	if _, err := s.workflows.Record(model.WorkflowFor(created.Status), created.ID, command.Actor, "registration created", created.CreatedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	return created, nil
}

func (s *Service) SubmitForReview(recordID, actor string) (model.Record, error) {
	record, err := s.records.Find(recordID)
	if err != nil {
		return model.Record{}, err
	}
	if record.Status != model.StatusDraft {
		return model.Record{}, fmt.Errorf("record %s is not a draft", recordID)
	}
	record.Status = model.StatusInReview
	record.UpdatedAt = s.clock.Now()
	if err := s.records.Save(record); err != nil {
		return model.Record{}, err
	}
	if _, err := s.audit.Append(recordID, model.ActionSubmitted, actor, "sent to review", record.UpdatedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	if _, err := s.workflows.Record(model.WorkflowFor(record.Status), recordID, actor, "review requested", record.UpdatedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Review(command model.ReviewCommand) (model.Record, error) {
	record, err := s.records.Find(command.RecordID)
	if err != nil {
		return model.Record{}, err
	}
	decision := policy.EvaluateReview(record, command.Actor, command.Approved)
	if !decision.Allowed {
		return model.Record{}, fmt.Errorf("review denied: %s", policy.JoinReasons(decision.Reasons))
	}
	if record.Status != model.StatusInReview {
		return model.Record{}, fmt.Errorf("record %s is not awaiting review", command.RecordID)
	}
	if command.Approved {
		record.Status = model.StatusApproved
	} else {
		record.Status = model.StatusDraft
	}
	record.UpdatedAt = s.clock.Now()
	if err := s.records.Save(record); err != nil {
		return model.Record{}, err
	}
	action := model.ActionRejected
	if command.Approved {
		action = model.ActionApproved
	}
	if _, err := s.audit.Append(record.ID, action, command.Actor, command.Note, record.UpdatedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	if _, err := s.workflows.Record(model.WorkflowFor(record.Status), record.ID, command.Actor, command.Note, record.UpdatedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Change(command model.ChangeCommand) (model.Record, error) {
	record, err := s.records.Find(command.RecordID)
	if err != nil {
		return model.Record{}, err
	}
	if record.IsClosed() || record.Status == model.StatusPublished {
		return model.Record{}, errors.New("published or archived record must be copied before changing")
	}
	if strings.TrimSpace(command.Title) != "" {
		record.Title = strings.TrimSpace(command.Title)
	}
	if strings.TrimSpace(command.Body) != "" {
		record.Body = command.Body
		record.CharacterCount = model.CharacterCount(command.Body)
	}
	if command.Description != "" {
		record.Description = command.Description
	}
	if command.Tags != nil {
		record.Tags = append([]string(nil), command.Tags...)
	}
	record.Revision++
	record.UpdatedAt = s.clock.Now()
	if err := record.Validate(); err != nil {
		return model.Record{}, err
	}
	if err := s.records.Save(record); err != nil {
		return model.Record{}, err
	}
	if _, err := s.audit.Append(record.ID, model.ActionChanged, command.Actor, "notice content changed", record.UpdatedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Publish(command model.PublishCommand) (model.Record, error) {
	record, err := s.records.Find(command.RecordID)
	if err != nil {
		return model.Record{}, err
	}
	decision := policy.EvaluatePublication(record, command.Actor)
	if !decision.Allowed {
		return model.Record{}, fmt.Errorf("publication denied: %s", policy.JoinReasons(decision.Reasons))
	}
	if record.Status != model.StatusApproved {
		return model.Record{}, errors.New("only approved records can be published")
	}
	record.Status = model.StatusPublished
	record.PublishedAt = s.clock.Now()
	record.UpdatedAt = record.PublishedAt
	if err := s.records.Save(record); err != nil {
		return model.Record{}, err
	}
	if _, err := s.audit.Append(record.ID, model.ActionPublished, command.Actor, "published to community", record.PublishedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	if _, err := s.workflows.Record(model.WorkflowFor(record.Status), record.ID, command.Actor, "publication complete", record.PublishedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Archive(command model.ArchiveCommand) (model.Record, error) {
	record, err := s.records.Find(command.RecordID)
	if err != nil {
		return model.Record{}, err
	}
	decision := policy.EvaluateArchive(record, command.Actor, command.Reason)
	if !decision.Allowed {
		return model.Record{}, fmt.Errorf("archive denied: %s", policy.JoinReasons(decision.Reasons))
	}
	if record.Status != model.StatusPublished {
		return model.Record{}, errors.New("only published records can be archived")
	}
	record.Status = model.StatusArchived
	record.ArchivedAt = s.clock.Now()
	record.UpdatedAt = record.ArchivedAt
	if err := s.records.Save(record); err != nil {
		return model.Record{}, err
	}
	if _, err := s.audit.Append(record.ID, model.ActionArchived, command.Actor, command.Reason, record.ArchivedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	if _, err := s.workflows.Record(model.WorkflowFor(record.Status), record.ID, command.Actor, "archive complete", record.ArchivedAt, s.nextSequence()); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Search(query model.SearchQuery) (model.SearchResult, error) {
	return s.records.Search(query)
}

func (s *Service) Get(recordID string) (model.Record, error) {
	return s.records.Find(recordID)
}

func (s *Service) Audit(recordID string) ([]model.AuditEvent, error) {
	return s.audit.ForRecord(recordID)
}

func (s *Service) Workflow(recordID string) ([]model.Workflow, error) {
	return s.workflows.ForRecord(recordID)
}
