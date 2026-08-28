package analytics

import (
	"sort"
	"strings"

	"noticeword/internal/model"
)

type Bucket struct {
	Key         string `json:"key"`
	Count       int    `json:"count"`
	Characters  int    `json:"characters"`
	Published   int    `json:"published"`
	Communities int    `json:"communities"`
}

type Trend struct {
	Buckets []Bucket `json:"buckets"`
	Total   int      `json:"total"`
}

func Build(records []model.Record) Trend {
	groups := map[string]*Bucket{}
	communitySets := map[string]map[string]struct{}{}
	for _, record := range records {
		key := period(record.UpdatedAt)
		bucket := groups[key]
		if bucket == nil {
			bucket = &Bucket{Key: key}
			groups[key] = bucket
			communitySets[key] = map[string]struct{}{}
		}
		bucket.Count++
		bucket.Characters += record.CharacterCount
		if record.Status == model.StatusPublished {
			bucket.Published++
		}
		communitySets[key][record.Community] = struct{}{}
	}
	result := Trend{Buckets: make([]Bucket, 0, len(groups)), Total: len(records)}
	for key, bucket := range groups {
		bucket.Communities = len(communitySets[key])
		result.Buckets = append(result.Buckets, *bucket)
	}
	sort.Slice(result.Buckets, func(i, j int) bool { return result.Buckets[i].Key < result.Buckets[j].Key })
	return result
}

func period(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 7 {
		return value[:7]
	}
	if value == "" {
		return "unknown"
	}
	return value
}

func PublishedRate(bucket Bucket) int {
	if bucket.Count == 0 {
		return 0
	}
	return (bucket.Published*100 + bucket.Count/2) / bucket.Count
}

func CharacterAverage(bucket Bucket) int {
	if bucket.Count == 0 {
		return 0
	}
	return bucket.Characters / bucket.Count
}

func Merge(left, right Trend) Trend {
	merged := map[string]Bucket{}
	for _, bucket := range left.Buckets {
		merged[bucket.Key] = bucket
	}
	for _, bucket := range right.Buckets {
		current := merged[bucket.Key]
		current.Key = bucket.Key
		current.Count += bucket.Count
		current.Characters += bucket.Characters
		current.Published += bucket.Published
		current.Communities += bucket.Communities
		merged[bucket.Key] = current
	}
	result := Trend{Total: left.Total + right.Total, Buckets: make([]Bucket, 0, len(merged))}
	for _, bucket := range merged {
		result.Buckets = append(result.Buckets, bucket)
	}
	sort.Slice(result.Buckets, func(i, j int) bool { return result.Buckets[i].Key < result.Buckets[j].Key })
	return result
}

func Filter(trend Trend, minimumCount int) Trend {
	result := Trend{Buckets: make([]Bucket, 0, len(trend.Buckets))}
	for _, bucket := range trend.Buckets {
		if bucket.Count >= minimumCount {
			result.Buckets = append(result.Buckets, bucket)
			result.Total += bucket.Count
		}
	}
	return result
}

func Growth(left, right Bucket) int {
	if left.Count == 0 {
		if right.Count > 0 {
			return 100
		}
		return 0
	}
	return ((right.Count-left.Count)*100 + left.Count/2) / left.Count
}

func BucketKeys(trend Trend) []string {
	keys := make([]string, 0, len(trend.Buckets))
	for _, bucket := range trend.Buckets {
		keys = append(keys, bucket.Key)
	}
	return keys
}
