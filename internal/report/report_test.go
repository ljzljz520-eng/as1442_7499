package report

import (
	"testing"

	"noticeword/internal/model"
)

func TestBuildAndRender(t *testing.T) {
	summary := Build([]model.Record{{Community: "A", Status: model.StatusPublished, CharacterCount: 4}, {Community: "A", Status: model.StatusDraft, CharacterCount: 2}, {Community: "B", Status: model.StatusPublished, CharacterCount: 5}})
	if summary.Total != 3 || summary.Characters != 11 || summary.ByCommunity["A"] != 2 {
		t.Fatalf("summary=%#v", summary)
	}
	if Render(summary) == "" || Percent(1, 3) != 33 {
		t.Fatal("render or percent mismatch")
	}
}

func TestStatusOrder(t *testing.T) {
	if len(StatusOrder()) != 5 {
		t.Fatal("status order incomplete")
	}
}
