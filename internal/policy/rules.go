package policy

import (
	"strings"

	"noticeword/internal/model"
)

type Decision struct {
	Allowed bool
	Code    string
	Reasons []string
}

func EvaluateRegistration(command model.RegisterCommand) Decision {
	reasons := make([]string, 0, 5)
	if strings.TrimSpace(command.Actor) == "" {
		reasons = append(reasons, "actor_missing")
	}
	if strings.TrimSpace(command.ID) == "" {
		reasons = append(reasons, "id_missing")
	}
	if strings.TrimSpace(command.Community) == "" {
		reasons = append(reasons, "community_missing")
	}
	if strings.TrimSpace(command.Title) == "" {
		reasons = append(reasons, "title_missing")
	}
	if strings.TrimSpace(command.Body) == "" {
		reasons = append(reasons, "body_missing")
	}
	return Decision{Allowed: len(reasons) == 0, Code: decisionCode(reasons), Reasons: reasons}
}

func EvaluateReview(record model.Record, actor string, approved bool) Decision {
	reasons := make([]string, 0, 4)
	if strings.TrimSpace(actor) == "" {
		reasons = append(reasons, "actor_missing")
	}
	if record.Status != model.StatusInReview {
		reasons = append(reasons, "not_in_review")
	}
	if approved && record.CharacterCount == 0 {
		reasons = append(reasons, "empty_body")
	}
	if len(record.Title) > 120 {
		reasons = append(reasons, "title_too_long")
	}
	return Decision{Allowed: len(reasons) == 0, Code: decisionCode(reasons), Reasons: reasons}
}

func EvaluatePublication(record model.Record, actor string) Decision {
	reasons := make([]string, 0, 3)
	if strings.TrimSpace(actor) == "" {
		reasons = append(reasons, "actor_missing")
	}
	if record.Status != model.StatusApproved {
		reasons = append(reasons, "not_approved")
	}
	if strings.TrimSpace(record.Description) == "" {
		reasons = append(reasons, "description_missing")
	}
	return Decision{Allowed: len(reasons) == 0, Code: decisionCode(reasons), Reasons: reasons}
}

func EvaluateArchive(record model.Record, actor, reason string) Decision {
	reasons := make([]string, 0, 3)
	if strings.TrimSpace(actor) == "" {
		reasons = append(reasons, "actor_missing")
	}
	if record.Status != model.StatusPublished {
		reasons = append(reasons, "not_published")
	}
	if strings.TrimSpace(reason) == "" {
		reasons = append(reasons, "reason_missing")
	}
	return Decision{Allowed: len(reasons) == 0, Code: decisionCode(reasons), Reasons: reasons}
}

func decisionCode(reasons []string) string {
	if len(reasons) == 0 {
		return "allow"
	}
	return "deny"
}

func NormalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		normalized := strings.TrimSpace(strings.ToLower(tag))
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func ContainsBlockedTerm(value string, blocked []string) bool {
	value = strings.ToLower(value)
	for _, term := range blocked {
		if strings.TrimSpace(term) != "" && strings.Contains(value, strings.ToLower(strings.TrimSpace(term))) {
			return true
		}
	}
	return false
}
