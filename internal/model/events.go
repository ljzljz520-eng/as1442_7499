package model

const (
	ActionRegistered = "registered"
	ActionSubmitted  = "submitted_for_review"
	ActionApproved   = "approved"
	ActionRejected   = "rejected"
	ActionChanged    = "changed"
	ActionPublished  = "published"
	ActionArchived   = "archived"
	ActionImported   = "imported"
)

func WorkflowFor(status RecordStatus) string {
	switch status {
	case StatusDraft:
		return "registration"
	case StatusInReview:
		return "review"
	case StatusApproved:
		return "approval"
	case StatusPublished:
		return "publication"
	case StatusArchived:
		return "archive"
	default:
		return "unknown"
	}
}

func AllowedTransition(from, to RecordStatus) bool {
	if from == StatusDraft && to == StatusInReview {
		return true
	}
	if from == StatusInReview && (to == StatusApproved || to == StatusDraft) {
		return true
	}
	if from == StatusApproved && (to == StatusPublished || to == StatusDraft) {
		return true
	}
	if from == StatusPublished && to == StatusArchived {
		return true
	}
	return false
}
