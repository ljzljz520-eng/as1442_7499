package catalog

import (
	"fmt"
	"sort"
	"strings"
)

type Window struct {
	ID        string
	RecordID  string
	Start     string
	End       string
	Audience  string
	Published bool
}

type Schedule struct {
	windows map[string]Window
}

func NewSchedule() *Schedule { return &Schedule{windows: map[string]Window{}} }

func (s *Schedule) Add(window Window) error {
	if strings.TrimSpace(window.ID) == "" {
		return fmt.Errorf("window id is required")
	}
	if strings.TrimSpace(window.RecordID) == "" {
		return fmt.Errorf("record id is required")
	}
	if strings.TrimSpace(window.Start) == "" || strings.TrimSpace(window.End) == "" {
		return fmt.Errorf("window bounds are required")
	}
	if window.Start >= window.End {
		return fmt.Errorf("window start must precede end")
	}
	s.windows[window.ID] = window
	return nil
}

func (s *Schedule) Remove(id string) bool {
	if _, exists := s.windows[id]; !exists {
		return false
	}
	delete(s.windows, id)
	return true
}

func (s *Schedule) Upcoming(at string, audience string) []Window {
	result := make([]Window, 0)
	for _, window := range s.windows {
		if window.Start < at {
			continue
		}
		if audience != "" && window.Audience != "" && window.Audience != audience {
			continue
		}
		result = append(result, window)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Start == result[j].Start {
			return result[i].ID < result[j].ID
		}
		return result[i].Start < result[j].Start
	})
	return result
}

func (s *Schedule) Conflicts(candidate Window) []Window {
	conflicts := make([]Window, 0)
	for _, existing := range s.windows {
		if existing.ID == candidate.ID || existing.RecordID != candidate.RecordID {
			continue
		}
		if existing.Start < candidate.End && candidate.Start < existing.End {
			conflicts = append(conflicts, existing)
		}
	}
	sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].ID < conflicts[j].ID })
	return conflicts
}

func (s *Schedule) Publish(id string) error {
	window, exists := s.windows[id]
	if !exists {
		return fmt.Errorf("window %s not found", id)
	}
	window.Published = true
	s.windows[id] = window
	return nil
}

func (s *Schedule) PublishedFor(recordID string) []Window {
	result := make([]Window, 0)
	for _, window := range s.windows {
		if window.RecordID == recordID && window.Published {
			result = append(result, window)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Start < result[j].Start })
	return result
}

func (s *Schedule) Count() int { return len(s.windows) }
