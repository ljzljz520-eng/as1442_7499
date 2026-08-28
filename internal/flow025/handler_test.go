package flow025

import (
	"testing"

	"noticeword/internal/model"
)

type fixtureReader struct {
	items []model.Record
}

func (r fixtureReader) Search(query model.SearchQuery) (model.SearchResult, error) {
	return model.SearchResult{Items: r.items, Page: query.Page, PageSize: query.PageSize, Total: 3, TotalPages: 2}, nil
}

func TestFirstPageDescriptions(t *testing.T) {
	handler := NewHandler(fixtureReader{items: FixtureRecords()[:1]})
	page, err := handler.Search(FixtureQuery(1, 2))
	if err != nil {
		t.Fatal(err)
	}
	if page.Rows[0].Description != "集合地点与报名说明" {
		t.Fatalf("description=%q", page.Rows[0].Description)
	}
}

func Test1442BusinessRegression(t *testing.T) {
	handler := NewHandler(fixtureReader{items: FixtureRecords()[1:]})
	page, err := handler.Search(FixtureQuery(2, 2))
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Rows) != 2 {
		t.Fatalf("rows=%d", len(page.Rows))
	}
	if page.Rows[0].Description != "报名截止日期和联系人" {
		t.Fatalf("description=%q", page.Rows[0].Description)
	}
	if page.Rows[1].Description != "分享会流程和准备事项" {
		t.Fatalf("description=%q", page.Rows[1].Description)
	}
}

func TestPageValidation(t *testing.T) {
	handler := NewHandler(fixtureReader{items: FixtureRecords()})
	page, err := handler.Search(FixtureQuery(1, 3))
	if err != nil {
		t.Fatal(err)
	}
	if err := page.Validate(); err != nil {
		t.Fatal(err)
	}
}
