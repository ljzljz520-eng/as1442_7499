package policy

import (
	"fmt"
	"strings"

	"noticeword/internal/model"
)

type Label struct {
	Key      string
	Text     string
	Tone     string
	Terminal bool
}

func StatusLabel(status model.RecordStatus) Label {
	switch status {
	case model.StatusDraft:
		return Label{Key: string(status), Text: "草稿", Tone: "neutral"}
	case model.StatusInReview:
		return Label{Key: string(status), Text: "审核中", Tone: "warning"}
	case model.StatusApproved:
		return Label{Key: string(status), Text: "已批准", Tone: "positive"}
	case model.StatusPublished:
		return Label{Key: string(status), Text: "已发布", Tone: "positive"}
	case model.StatusArchived:
		return Label{Key: string(status), Text: "已归档", Tone: "muted", Terminal: true}
	default:
		return Label{Key: string(status), Text: "未知", Tone: "danger"}
	}
}

func TransitionMessage(from, to model.RecordStatus) string {
	if from == to {
		return "状态未变化"
	}
	if !model.AllowedTransition(from, to) {
		return fmt.Sprintf("不允许从%s变更为%s", StatusLabel(from).Text, StatusLabel(to).Text)
	}
	return fmt.Sprintf("%s -> %s", StatusLabel(from).Text, StatusLabel(to).Text)
}

func JoinReasons(reasons []string) string {
	clean := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		if strings.TrimSpace(reason) != "" {
			clean = append(clean, strings.TrimSpace(reason))
		}
	}
	return strings.Join(clean, ", ")
}
