package service

import (
	"path/filepath"
	"testing"

	"noticeword/internal/model"
	"noticeword/internal/store"
)

func newTestService(t *testing.T) (*Service, func()) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	app := New(db, store.NewSequenceClock("2025-01-01T00:00:00Z"))
	return app, func() { db.Close() }
}

func registerForTest(t *testing.T, app *Service, id string) model.Record {
	t.Helper()
	record, err := app.Register(model.RegisterCommand{ID: id, Community: "星河社群", Title: "公告", Body: "公告正文", Description: "公告说明", Tags: []string{"活动"}, Actor: "editor"})
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	app, closeStore := newTestService(t)
	defer closeStore()
	record := registerForTest(t, app, "wf-archive")
	if _, err := app.SubmitForReview(record.ID, "editor"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.Review(model.ReviewCommand{RecordID: record.ID, Actor: "reviewer", Approved: true, Note: "内容合规"}); err != nil {
		t.Fatal(err)
	}
	if _, err := app.Publish(model.PublishCommand{RecordID: record.ID, Actor: "publisher"}); err != nil {
		t.Fatal(err)
	}
	archived, err := app.Archive(model.ArchiveCommand{RecordID: record.ID, Actor: "owner", Reason: "活动结束"})
	if err != nil {
		t.Fatal(err)
	}
	if archived.Status != model.StatusArchived {
		t.Fatalf("status=%s", archived.Status)
	}
	if len(mustAuditService(t, app, record.ID)) != 5 {
		t.Fatal("audit trail incomplete")
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	app, closeStore := newTestService(t)
	defer closeStore()
	record := registerForTest(t, app, "wf-publish")
	changed, err := app.Change(model.ChangeCommand{RecordID: record.ID, Body: "更新后的公告正文", Description: "最新说明", Actor: "editor"})
	if err != nil {
		t.Fatal(err)
	}
	if changed.CharacterCount != model.CharacterCount("更新后的公告正文") {
		t.Fatal("character count did not update")
	}
	if _, err := app.SubmitForReview(record.ID, "editor"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.Review(model.ReviewCommand{RecordID: record.ID, Actor: "reviewer", Approved: true}); err != nil {
		t.Fatal(err)
	}
	published, err := app.Publish(model.PublishCommand{RecordID: record.ID, Actor: "publisher"})
	if err != nil {
		t.Fatal(err)
	}
	if !published.IsVisible() || published.PublishedAt == "" {
		t.Fatal("publication not visible")
	}
	result, err := app.Search(model.SearchQuery{Text: "更新", Page: 1, PageSize: 10})
	if err != nil || len(result.Items) != 1 {
		t.Fatalf("search result: %v %#v", err, result)
	}
}

func TestWorkflowImportReport(t *testing.T) {
	app, closeStore := newTestService(t)
	defer closeStore()
	report := app.Import([]model.ImportRow{{ID: "import-1", Community: "导入社群", Title: "有效", Body: "有效正文", Description: "有效说明"}, {ID: "", Community: "导入社群", Title: "无效", Body: "缺少编号"}}, "importer")
	if report.Accepted != 1 || report.Rejected != 1 || len(report.IDs) != 1 {
		t.Fatalf("report=%#v", report)
	}
	if _, err := app.Get("import-1"); err != nil {
		t.Fatal(err)
	}
}

func mustAuditService(t *testing.T, app *Service, id string) []model.AuditEvent {
	t.Helper()
	items, err := app.Audit(id)
	if err != nil {
		t.Fatal(err)
	}
	return items
}
