package service

import (
	"sort"

	"noticeword/internal/model"
	"noticeword/internal/report"
)

type Metrics struct {
	Summary       report.Summary `json:"summary"`
	Visible       int            `json:"visible"`
	NeedsReview   int            `json:"needs_review"`
	AverageLength int            `json:"average_length"`
	TopTags       []string       `json:"top_tags"`
}

func BuildMetrics(records []model.Record) Metrics {
	metrics := Metrics{Summary: report.Build(records), TopTags: topTags(records)}
	for _, record := range records {
		if record.IsVisible() {
			metrics.Visible++
		}
		if record.Status == model.StatusInReview {
			metrics.NeedsReview++
		}
	}
	if metrics.Summary.Total > 0 {
		metrics.AverageLength = metrics.Summary.Characters / metrics.Summary.Total
	}
	return metrics
}

func topTags(records []model.Record) []string {
	counts := map[string]int{}
	for _, record := range records {
		for _, tag := range record.Tags {
			counts[tag]++
		}
	}
	type pair struct {
		name  string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for name, count := range counts {
		pairs = append(pairs, pair{name, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].name < pairs[j].name
		}
		return pairs[i].count > pairs[j].count
	})
	limit := 5
	if len(pairs) < limit {
		limit = len(pairs)
	}
	result := make([]string, 0, limit)
	for _, item := range pairs[:limit] {
		result = append(result, item.name)
	}
	return result
}

func (s *Service) Metrics(query model.SearchQuery) (Metrics, error) {
	records, err := s.Export(query)
	if err != nil {
		return Metrics{}, err
	}
	return BuildMetrics(records), nil
}
