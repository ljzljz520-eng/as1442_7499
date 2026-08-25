package catalog

import (
	"testing"

	"noticeword/internal/model"
)

func sampleCatalogRecords() []model.Record {
	return []model.Record{{ID: "a", Community: "星河", Status: model.StatusPublished, Tags: []string{"活动", "报名"}, Title: "活动", Body: "周六"}, {ID: "b", Community: "星河", Status: model.StatusDraft, Tags: []string{"报名"}, Title: "招募", Body: "志愿者"}, {ID: "c", Community: "晨光", Status: model.StatusPublished, Tags: []string{"分享"}, Title: "读书", Body: "本月"}}
}

func TestIndexQueriesAndFacets(t *testing.T) {
	index := NewIndex()
	index.Rebuild(sampleCatalogRecords())
	items := index.Query("星河", model.StatusPublished, []string{"报名"})
	if len(items) != 1 || items[0].ID != "a" {
		t.Fatalf("items=%#v", items)
	}
	if len(index.Facets().Tags) != 3 {
		t.Fatal("facets missing")
	}
	if !index.Remove("c") || index.Remove("missing") {
		t.Fatal("remove result incorrect")
	}
	if _, ok := index.Get("c"); ok {
		t.Fatal("removed record remains")
	}
}

func TestScheduleAndQueue(t *testing.T) {
	schedule := NewSchedule()
	if err := schedule.Add(Window{ID: "w1", RecordID: "a", Start: "2025-01-02", End: "2025-01-03", Audience: "all"}); err != nil {
		t.Fatal(err)
	}
	if err := schedule.Add(Window{ID: "w2", RecordID: "a", Start: "2025-01-02T12", End: "2025-01-04", Audience: "all"}); err != nil {
		t.Fatal(err)
	}
	if len(schedule.Conflicts(Window{ID: "w3", RecordID: "a", Start: "2025-01-02T08", End: "2025-01-02T16"})) != 2 {
		t.Fatal("conflict detection failed")
	}
	if err := schedule.Publish("w1"); err != nil || len(schedule.PublishedFor("a")) != 1 {
		t.Fatal("publish failed")
	}
	queue := NewQueue()
	if err := queue.Enqueue(QueueItem{ID: "q1", RecordID: "a", Priority: 1, Enqueued: "1"}); err != nil {
		t.Fatal(err)
	}
	if err := queue.Enqueue(QueueItem{ID: "q2", RecordID: "b", Priority: 2, Enqueued: "2"}); err != nil {
		t.Fatal(err)
	}
	item, ok := queue.Next("")
	if !ok || item.ID != "q2" {
		t.Fatalf("next=%#v", item)
	}
	if err := queue.Complete(item.ID, "done"); err != nil {
		t.Fatal(err)
	}
	if counts := queue.Counts(); counts["completed"] != 1 || counts["pending"] != 1 {
		t.Fatalf("counts=%v", counts)
	}
}

func TestTokenRanking(t *testing.T) {
	if !MatchAll("公告说明与报名", []string{"说明", "报名"}) {
		t.Fatal("match all failed")
	}
	if Rank("公告 公告 报名", []string{"公告"}) != 3 {
		t.Fatal("rank failed")
	}
	terms := TopTerms([]string{"公告报名", "公告说明"}, 2)
	if len(terms) != 2 {
		t.Fatalf("terms=%v", terms)
	}
	if Highlight("公告说明", "说明", "[", "]") != "公告[说明]" {
		t.Fatal("highlight failed")
	}
}
