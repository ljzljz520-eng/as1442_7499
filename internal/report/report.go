package report

import (
	"fmt"
	"sort"
	"strings"

	"noticeword/internal/model"
)

type Summary struct {
	Total       int            `json:"total"`
	ByStatus    map[string]int `json:"by_status"`
	ByCommunity map[string]int `json:"by_community"`
	Characters  int            `json:"characters"`
}

func Build(records []model.Record) Summary {
	summary := Summary{ByStatus: map[string]int{}, ByCommunity: map[string]int{}}
	for _, record := range records {
		summary.Total++
		summary.ByStatus[string(record.Status)]++
		summary.ByCommunity[record.Community]++
		summary.Characters += record.CharacterCount
	}
	return summary
}

func Render(summary Summary) string {
	communities := make([]string, 0, len(summary.ByCommunity))
	for community := range summary.ByCommunity {
		communities = append(communities, community)
	}
	sort.Strings(communities)
	parts := []string{fmt.Sprintf("total=%d", summary.Total), fmt.Sprintf("characters=%d", summary.Characters)}
	for _, community := range communities {
		parts = append(parts, fmt.Sprintf("%s=%d", community, summary.ByCommunity[community]))
	}
	return strings.Join(parts, " ")
}

func StatusOrder() []model.RecordStatus {
	return []model.RecordStatus{model.StatusDraft, model.StatusInReview, model.StatusApproved, model.StatusPublished, model.StatusArchived}
}

func Percent(value, total int) int {
	if total <= 0 || value <= 0 {
		return 0
	}
	return (value*100 + total/2) / total
}
