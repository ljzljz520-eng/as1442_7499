package repository

import (
	"path/filepath"
	"testing"

	"noticeword/internal/model"
	"noticeword/internal/store"
)

func TestSearchFiltersAndPaginates(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "repo.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewRecordRepository(db)
	for i := 0; i < 5; i++ {
		record := model.NewRecord(string(rune('a'+i)), "社群", "公告", "搜索内容", "说明", "2025", nil)
		record.Status = model.StatusPublished
		if err := repo.Save(record); err != nil {
			t.Fatal(err)
		}
	}
	result, err := repo.Search(model.SearchQuery{Text: "内容", Page: 2, PageSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 5 || len(result.Items) != 2 || !result.HasNext() {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestSupportingRepositories(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "support.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	audit := NewAuditRepository(db)
	if _, err := audit.Append("r", "registered", "a", "d", "t", 1); err != nil {
		t.Fatal(err)
	}
	if len(mustAudit(t, audit, "r")) != 1 {
		t.Fatal("audit missing")
	}
	workflow := NewWorkflowRepository(db)
	if _, err := workflow.Record("registration", "r", "a", "n", "t", 1); err != nil {
		t.Fatal(err)
	}
	if len(mustWorkflow(t, workflow, "r")) != 1 {
		t.Fatal("workflow missing")
	}
}

func mustAudit(t *testing.T, repo *AuditRepository, id string) []model.AuditEvent {
	t.Helper()
	items, err := repo.ForRecord(id)
	if err != nil {
		t.Fatal(err)
	}
	return items
}
func mustWorkflow(t *testing.T, repo *WorkflowRepository, id string) []model.Workflow {
	t.Helper()
	items, err := repo.ForRecord(id)
	if err != nil {
		t.Fatal(err)
	}
	return items
}
