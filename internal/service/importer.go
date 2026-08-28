package service

import (
	"fmt"
	"strings"

	"noticeword/internal/model"
)

func (s *Service) Import(rows []model.ImportRow, actor string) model.ImportReport {
	report := model.ImportReport{IDs: make([]string, 0, len(rows))}
	for index, row := range rows {
		if strings.TrimSpace(row.ID) == "" {
			report.AddError(fmt.Sprintf("row %d has no id", index+1))
			continue
		}
		if strings.TrimSpace(row.Body) == "" {
			report.AddError(fmt.Sprintf("row %d has no body", index+1))
			continue
		}
		_, err := s.Register(model.RegisterCommand{ID: row.ID, Community: row.Community, Title: row.Title, Body: row.Body, Description: row.Description, Tags: row.Tags, Actor: actor})
		if err != nil {
			report.AddError(fmt.Sprintf("row %d: %v", index+1, err))
			continue
		}
		report.Accepted++
		report.IDs = append(report.IDs, row.ID)
		s.markImported(row.ID, actor)
	}
	return report
}

func (s *Service) markImported(recordID, actor string) {
	record, err := s.records.Find(recordID)
	if err != nil {
		return
	}
	_, _ = s.audit.Append(recordID, model.ActionImported, actor, "loaded from deterministic import", s.clock.Now(), s.nextSequence())
	_ = s.records.Save(record)
}

func (s *Service) Export(query model.SearchQuery) ([]model.Record, error) {
	query.Page = 1
	query.PageSize = 100
	result, err := s.Search(query)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
