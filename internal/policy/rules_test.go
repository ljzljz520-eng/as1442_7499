package policy

import (
	"testing"

	"noticeword/internal/model"
)

func TestPolicyDecisions(t *testing.T) {
	allowed := EvaluateRegistration(model.RegisterCommand{ID: "r", Community: "c", Title: "t", Body: "b", Actor: "a"})
	if !allowed.Allowed || allowed.Code != "allow" {
		t.Fatalf("allowed=%#v", allowed)
	}
	denied := EvaluateReview(model.Record{Status: model.StatusDraft}, "", true)
	if denied.Allowed || len(denied.Reasons) < 2 {
		t.Fatalf("denied=%#v", denied)
	}
	if !EvaluatePublication(model.Record{Status: model.StatusApproved, Description: "d"}, "p").Allowed {
		t.Fatal("publication should be allowed")
	}
	if EvaluateArchive(model.Record{Status: model.StatusPublished}, "a", "").Allowed {
		t.Fatal("archive needs reason")
	}
}

func TestTagAndBlockedTermRules(t *testing.T) {
	tags := NormalizeTags([]string{" 活动 ", "活动", "", "报名"})
	if len(tags) != 2 || tags[0] != "活动" {
		t.Fatalf("tags=%v", tags)
	}
	if !ContainsBlockedTerm("含有禁词的说明", []string{"禁词"}) {
		t.Fatal("blocked term not found")
	}
	if TransitionMessage(model.StatusDraft, model.StatusPublished) == "" || JoinReasons([]string{"a", ""}) != "a" {
		t.Fatal("policy labels missing")
	}
}
