package analytics

import (
	"testing"

	"noticeword/internal/model"
)

func TestTrendBuildAndMerge(t *testing.T) {
	left := Build([]model.Record{{UpdatedAt: "2025-01-02T00", CharacterCount: 10, Status: model.StatusPublished, Community: "A"}, {UpdatedAt: "2025-01-02T01", CharacterCount: 6, Status: model.StatusDraft, Community: "B"}})
	right := Build([]model.Record{{UpdatedAt: "2025-02-01", CharacterCount: 8, Status: model.StatusPublished, Community: "A"}})
	if len(left.Buckets) != 1 || left.Buckets[0].Count != 2 {
		t.Fatalf("left=%#v", left)
	}
	if PublishedRate(left.Buckets[0]) != 50 || CharacterAverage(left.Buckets[0]) != 8 {
		t.Fatal("bucket calculations failed")
	}
	merged := Merge(left, right)
	if merged.Total != 3 || len(merged.Buckets) != 2 {
		t.Fatalf("merged=%#v", merged)
	}
	if Filter(merged, 2).Total != 2 {
		t.Fatal("filter failed")
	}
}

func TestCommunityForecast(t *testing.T) {
	records := []model.Record{{ID: "a", Community: "A", Status: model.StatusPublished, CharacterCount: 20}, {ID: "b", Community: "A", Status: model.StatusDraft, CharacterCount: 3}}
	if CompletionRate(records) != 50 {
		t.Fatal("completion rate failed")
	}
	ranked := RankCommunities(records)
	if len(ranked) != 1 || ranked[0].Visible != 1 {
		t.Fatalf("ranked=%#v", ranked)
	}
	if Longest(records, 1)[0].ID != "a" {
		t.Fatal("longest failed")
	}
}
