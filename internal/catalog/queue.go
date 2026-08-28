package catalog

import (
	"fmt"
	"sort"
	"strings"
)

type QueueItem struct {
	ID       string
	RecordID string
	Priority int
	Owner    string
	State    string
	Enqueued string
	Note     string
}

type Queue struct {
	items map[string]QueueItem
}

func NewQueue() *Queue { return &Queue{items: map[string]QueueItem{}} }

func (q *Queue) Enqueue(item QueueItem) error {
	if strings.TrimSpace(item.ID) == "" {
		return fmt.Errorf("queue id is required")
	}
	if strings.TrimSpace(item.RecordID) == "" {
		return fmt.Errorf("record id is required")
	}
	if item.Priority < 0 {
		item.Priority = 0
	}
	item.State = "pending"
	q.items[item.ID] = item
	return nil
}

func (q *Queue) Next(owner string) (QueueItem, bool) {
	candidates := q.Pending()
	for _, item := range candidates {
		if owner != "" && item.Owner != "" && item.Owner != owner {
			continue
		}
		item.State = "claimed"
		item.Owner = owner
		q.items[item.ID] = item
		return item, true
	}
	return QueueItem{}, false
}

func (q *Queue) Complete(id, note string) error { return q.finish(id, "completed", note) }
func (q *Queue) Reject(id, note string) error   { return q.finish(id, "rejected", note) }

func (q *Queue) finish(id, state, note string) error {
	item, exists := q.items[id]
	if !exists {
		return fmt.Errorf("queue item %s not found", id)
	}
	if item.State != "claimed" && item.State != "pending" {
		return fmt.Errorf("queue item %s is already closed", id)
	}
	item.State = state
	item.Note = note
	q.items[id] = item
	return nil
}

func (q *Queue) Pending() []QueueItem {
	result := make([]QueueItem, 0)
	for _, item := range q.items {
		if item.State == "pending" {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Priority == result[j].Priority {
			return result[i].Enqueued < result[j].Enqueued
		}
		return result[i].Priority > result[j].Priority
	})
	return result
}

func (q *Queue) Find(id string) (QueueItem, bool) { item, exists := q.items[id]; return item, exists }

func (q *Queue) Counts() map[string]int {
	counts := map[string]int{}
	for _, item := range q.items {
		counts[item.State]++
	}
	return counts
}

func (q *Queue) ResetClaim(id string) error {
	item, exists := q.items[id]
	if !exists {
		return fmt.Errorf("queue item %s not found", id)
	}
	if item.State != "claimed" {
		return fmt.Errorf("queue item %s is not claimed", id)
	}
	item.State = "pending"
	item.Owner = ""
	q.items[id] = item
	return nil
}
