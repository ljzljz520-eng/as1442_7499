package repository

import (
	"fmt"
	"noticeword/internal/model"
	"noticeword/internal/store"
)

type WorkflowRepository struct {
	storage *store.Store
}

func NewWorkflowRepository(storage *store.Store) *WorkflowRepository {
	return &WorkflowRepository{storage: storage}
}

func (r *WorkflowRepository) Record(stage, recordID, assignee, note, at string, sequence int) (model.Workflow, error) {
	if stage == "" || recordID == "" {
		return model.Workflow{}, fmt.Errorf("workflow stage and record are required")
	}
	workflow := model.Workflow{ID: store.WorkflowKey(recordID, sequence), RecordID: recordID, Stage: stage, Assignee: assignee, Note: note, OccurredAt: at}
	if err := r.storage.PutWorkflow(workflow); err != nil {
		return model.Workflow{}, err
	}
	return workflow, nil
}

func (r *WorkflowRepository) ForRecord(recordID string) ([]model.Workflow, error) {
	return r.storage.ListWorkflows(recordID)
}
