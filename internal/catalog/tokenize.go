package catalog

import (
	"sort"
	"strings"
	"unicode"
)

func Tokenize(value string) []string {
	words := make([]string, 0)
	current := make([]rune, 0)
	flush := func() {
		if len(current) > 0 {
			words = append(words, strings.ToLower(string(current)))
			current = current[:0]
		}
	}
	for _, r := range []rune(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.In(r, unicode.Han) {
			current = append(current, r)
			continue
		}
		flush()
	}
	flush()
	return words
}

func MatchAll(value string, terms []string) bool {
	lower := strings.ToLower(value)
	for _, term := range terms {
		if strings.TrimSpace(term) != "" && !strings.Contains(lower, strings.ToLower(strings.TrimSpace(term))) {
			return false
		}
	}
	return true
}

func Rank(value string, terms []string) int {
	lower := strings.ToLower(value)
	score := 0
	for _, term := range terms {
		term = strings.ToLower(strings.TrimSpace(term))
		if term == "" {
			continue
		}
		if strings.Contains(lower, term) {
			score += 1 + strings.Count(lower, term)
		}
	}
	return score
}

func TopTerms(values []string, limit int) []string {
	counts := map[string]int{}
	for _, value := range values {
		for _, token := range Tokenize(value) {
			counts[token]++
		}
	}
	type pair struct {
		term  string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for term, count := range counts {
		pairs = append(pairs, pair{term, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].term < pairs[j].term
		}
		return pairs[i].count > pairs[j].count
	})
	if limit < 0 {
		limit = 0
	}
	if len(pairs) < limit {
		limit = len(pairs)
	}
	result := make([]string, 0, limit)
	for _, item := range pairs[:limit] {
		result = append(result, item.term)
	}
	return result
}

func Highlight(value, term, left, right string) string {
	if strings.TrimSpace(term) == "" {
		return value
	}
	index := strings.Index(strings.ToLower(value), strings.ToLower(term))
	if index < 0 {
		return value
	}
	end := index + len(term)
	return value[:index] + left + value[index:end] + right + value[end:]
}
