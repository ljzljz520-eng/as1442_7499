package model

import "strings"

type RegisterCommand struct {
	ID          string
	Community   string
	Title       string
	Body        string
	Description string
	Tags        []string
	Actor       string
}

type ReviewCommand struct {
	RecordID string
	Actor    string
	Approved bool
	Note     string
}

type ChangeCommand struct {
	RecordID    string
	Title       string
	Body        string
	Description string
	Tags        []string
	Actor       string
}

type PublishCommand struct {
	RecordID string
	Actor    string
}

type ArchiveCommand struct {
	RecordID string
	Actor    string
	Reason   string
}

type SearchQuery struct {
	Community string
	Text      string
	Status    RecordStatus
	Page      int
	PageSize  int
}

func (q SearchQuery) Normalized() SearchQuery {
	q.Community = strings.TrimSpace(strings.ToLower(q.Community))
	q.Text = strings.TrimSpace(strings.ToLower(q.Text))
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 10
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	return q
}

type SearchResult struct {
	Items      []Record `json:"items"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	Total      int      `json:"total"`
	TotalPages int      `json:"total_pages"`
}

func (r SearchResult) HasNext() bool {
	return r.Page < r.TotalPages
}

type ImportRow struct {
	ID          string
	Community   string
	Title       string
	Body        string
	Description string
	Tags        []string
}

type ImportReport struct {
	Accepted int      `json:"accepted"`
	Rejected int      `json:"rejected"`
	Errors   []string `json:"errors"`
	IDs      []string `json:"ids"`
}

func (r *ImportReport) AddError(message string) {
	r.Rejected++
	r.Errors = append(r.Errors, message)
}
