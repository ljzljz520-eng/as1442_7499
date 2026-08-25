package analytics

import (
	"sort"

	"noticeword/internal/model"
)

type CommunityScore struct {
	Community  string `json:"community"`
	Visible    int    `json:"visible"`
	Drafts     int    `json:"drafts"`
	Characters int    `json:"characters"`
	Score      int    `json:"score"`
}

func RankCommunities(records []model.Record) []CommunityScore {
	groups := map[string]*CommunityScore{}
	for _, record := range records {
		score := groups[record.Community]
		if score == nil {
			score = &CommunityScore{Community: record.Community}
			groups[record.Community] = score
		}
		if record.IsVisible() {
			score.Visible++
		} else {
			score.Drafts++
		}
		score.Characters += record.CharacterCount
	}
	result := make([]CommunityScore, 0, len(groups))
	for _, score := range groups {
		score.Score = score.Visible*10 + score.Characters/10 - score.Drafts*2
		result = append(result, *score)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].Community < result[j].Community
		}
		return result[i].Score > result[j].Score
	})
	return result
}

func CompletionRate(records []model.Record) int {
	if len(records) == 0 {
		return 0
	}
	completed := 0
	for _, record := range records {
		if record.Status == model.StatusPublished || record.Status == model.StatusArchived {
			completed++
		}
	}
	return (completed*100 + len(records)/2) / len(records)
}

func Longest(records []model.Record, limit int) []model.Record {
	result := append([]model.Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].CharacterCount == result[j].CharacterCount {
			return result[i].ID < result[j].ID
		}
		return result[i].CharacterCount > result[j].CharacterCount
	})
	if limit < 0 {
		limit = 0
	}
	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

func StatusBreakdown(records []model.Record) map[model.RecordStatus]int {
	breakdown := map[model.RecordStatus]int{}
	for _, record := range records {
		breakdown[record.Status]++
	}
	return breakdown
}

func CommunityNames(records []model.Record) []string {
	seen := map[string]struct{}{}
	for _, record := range records {
		if record.Community != "" {
			seen[record.Community] = struct{}{}
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func RevisionAverage(records []model.Record) int {
	if len(records) == 0 {
		return 0
	}
	total := 0
	for _, record := range records {
		total += record.Revision
	}
	return total / len(records)
}

func VisibleCharacters(records []model.Record) int {
	total := 0
	for _, record := range records {
		if record.IsVisible() {
			total += record.CharacterCount
		}
	}
	return total
}
