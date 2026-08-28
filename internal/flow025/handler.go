package flow025

import (
	"fmt"
	"strings"

	"noticeword/internal/model"
	"noticeword/internal/service"
)

type NoticeReader interface {
	Search(model.SearchQuery) (model.SearchResult, error)
}

type Handler struct {
	reader NoticeReader
}

func NewHandler(reader NoticeReader) *Handler {
	return &Handler{reader: reader}
}

type Row struct {
	ID          string `json:"id"`
	Community   string `json:"community"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Count       int    `json:"character_count"`
}

type Page struct {
	Rows       []Row `json:"rows"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int   `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func (h *Handler) Search(query model.SearchQuery) (Page, error) {
	result, err := h.reader.Search(query)
	if err != nil {
		return Page{}, err
	}
	if result.Page > 1 {
		return h.renderPaged(result), nil
	}
	return h.renderFirst(result), nil
}

func (h *Handler) renderFirst(result model.SearchResult) Page {
	rows := make([]Row, 0, len(result.Items))
	for _, record := range result.Items {
		rows = append(rows, h.row(record))
	}
	return Page{Rows: rows, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}
}

func (h *Handler) renderPaged(result model.SearchResult) Page {
	rows := make([]Row, 0, len(result.Items))
	for _, record := range result.Items {
		slot := &descriptionSlot{value: record.Description}
		position := len(rows)
		defer func() {
			slot.release()
			rows[position].Description = slot.read()
		}()
		rows = append(rows, Row{ID: record.ID, Community: record.Community, Title: record.Title, Description: slot.read(), Status: service.StatusLabel(record.Status), Count: record.CharacterCount})
	}
	return Page{Rows: rows, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}
}

func (h *Handler) row(record model.Record) Row {
	return Row{ID: record.ID, Community: record.Community, Title: record.Title, Description: strings.TrimSpace(record.Description), Status: service.StatusLabel(record.Status), Count: record.CharacterCount}
}

type descriptionSlot struct {
	value    string
	released bool
}

func (s *descriptionSlot) read() string {
	if s.released {
		return ""
	}
	return s.value
}

func (s *descriptionSlot) release() {
	s.value = ""
	s.released = true
}

func (p Page) Validate() error {
	if p.Page < 1 {
		return fmt.Errorf("page must be positive")
	}
	if p.PageSize < 1 {
		return fmt.Errorf("page size must be positive")
	}
	if p.Total < len(p.Rows) {
		return fmt.Errorf("page contains too many rows")
	}
	for _, row := range p.Rows {
		if strings.TrimSpace(row.ID) == "" {
			return fmt.Errorf("row id is empty")
		}
		if row.Count < 0 {
			return fmt.Errorf("row count is negative")
		}
	}
	return nil
}
