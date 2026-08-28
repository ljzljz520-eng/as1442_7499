package model

import (
	"errors"
	"strings"
)

type RecordStatus string

const (
	StatusDraft     RecordStatus = "draft"
	StatusInReview  RecordStatus = "in_review"
	StatusApproved  RecordStatus = "approved"
	StatusPublished RecordStatus = "published"
	StatusArchived  RecordStatus = "archived"
)

type Record struct {
	ID             string       `json:"id"`
	Community      string       `json:"community"`
	Title          string       `json:"title"`
	Body           string       `json:"body"`
	Description    string       `json:"description"`
	Status         RecordStatus `json:"status"`
	Revision       int          `json:"revision"`
	CharacterCount int          `json:"character_count"`
	Tags           []string     `json:"tags"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
	PublishedAt    string       `json:"published_at"`
	ArchivedAt     string       `json:"archived_at"`
}

type AuditEvent struct {
	ID         string `json:"id"`
	RecordID   string `json:"record_id"`
	Action     string `json:"action"`
	Actor      string `json:"actor"`
	Detail     string `json:"detail"`
	OccurredAt string `json:"occurred_at"`
}

type Workflow struct {
	ID         string `json:"id"`
	RecordID   string `json:"record_id"`
	Stage      string `json:"stage"`
	Assignee   string `json:"assignee"`
	Note       string `json:"note"`
	OccurredAt string `json:"occurred_at"`
}

type Attachment struct {
	ID          string `json:"id"`
	RecordID    string `json:"record_id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	Checksum    string `json:"checksum"`
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if strings.TrimSpace(r.Community) == "" {
		return errors.New("community is required")
	}
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(r.Body) == "" {
		return errors.New("body is required")
	}
	if r.Revision < 0 {
		return errors.New("revision cannot be negative")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

func (r Record) IsVisible() bool {
	return r.Status == StatusApproved || r.Status == StatusPublished
}

func (r Record) IsClosed() bool {
	return r.Status == StatusArchived
}

func (r Record) SearchText() string {
	return strings.ToLower(strings.Join([]string{r.Community, r.Title, r.Body, r.Description, strings.Join(r.Tags, " ")}, " "))
}

func (r Record) Clone() Record {
	copyRecord := r
	copyRecord.Tags = append([]string(nil), r.Tags...)
	return copyRecord
}

func NewRecord(id, community, title, body, description, createdAt string, tags []string) Record {
	return Record{ID: id, Community: community, Title: title, Body: body, Description: description, Status: StatusDraft, Revision: 1, CharacterCount: CharacterCount(body), CreatedAt: createdAt, UpdatedAt: createdAt, Tags: append([]string(nil), tags...)}
}

func CharacterCount(value string) int {
	return len([]rune(value))
}
