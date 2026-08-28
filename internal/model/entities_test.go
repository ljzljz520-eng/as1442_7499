package model

import "testing"

func TestRecordValidationAndSearchText(t *testing.T) {
	record := NewRecord("r1", "星河", "活动", "三人报名", "周六说明", "2025", []string{"活动", "报名"})
	if err := record.Validate(); err != nil {
		t.Fatal(err)
	}
	if record.CharacterCount != 4 {
		t.Fatalf("count=%d", record.CharacterCount)
	}
	if !record.IsVisible() && record.Status == StatusDraft {
		t.Log("draft is correctly hidden")
	}
	if record.SearchText() == "" {
		t.Fatal("search text should not be empty")
	}
	copyRecord := record.Clone()
	copyRecord.Tags[0] = "changed"
	if record.Tags[0] == copyRecord.Tags[0] {
		t.Fatal("clone shares tags")
	}
}

func TestTransitions(t *testing.T) {
	if !AllowedTransition(StatusDraft, StatusInReview) {
		t.Fatal("draft should enter review")
	}
	if AllowedTransition(StatusDraft, StatusPublished) {
		t.Fatal("draft should not publish")
	}
	if WorkflowFor(StatusArchived) != "archive" {
		t.Fatal("archive workflow label missing")
	}
}
