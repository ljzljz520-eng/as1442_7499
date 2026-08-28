package service

import (
	"fmt"
	"strings"

	"noticeword/internal/model"
)

func ValidateRegister(command model.RegisterCommand) []string {
	issues := make([]string, 0, 4)
	if strings.TrimSpace(command.ID) == "" {
		issues = append(issues, "id is required")
	}
	if strings.TrimSpace(command.Community) == "" {
		issues = append(issues, "community is required")
	}
	if strings.TrimSpace(command.Title) == "" {
		issues = append(issues, "title is required")
	}
	if strings.TrimSpace(command.Body) == "" {
		issues = append(issues, "body is required")
	}
	if len([]rune(command.Body)) > 10000 {
		issues = append(issues, "body exceeds 10000 characters")
	}
	return issues
}

func ValidateChange(command model.ChangeCommand) error {
	if strings.TrimSpace(command.RecordID) == "" {
		return fmt.Errorf("record id is required")
	}
	if command.Body != "" && len([]rune(command.Body)) > 10000 {
		return fmt.Errorf("body exceeds 10000 characters")
	}
	return nil
}

func StatusLabel(status model.RecordStatus) string {
	switch status {
	case model.StatusDraft:
		return "草稿"
	case model.StatusInReview:
		return "审核中"
	case model.StatusApproved:
		return "已批准"
	case model.StatusPublished:
		return "已发布"
	case model.StatusArchived:
		return "已归档"
	default:
		return "未知"
	}
}
