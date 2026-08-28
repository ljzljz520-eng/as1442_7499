package service

import (
	"testing"

	"noticeword/internal/model"
)

func TestValidationMessages(t *testing.T) {
	issues := ValidateRegister(model.RegisterCommand{})
	if len(issues) != 4 {
		t.Fatalf("issues=%v", issues)
	}
	if err := ValidateChange(model.ChangeCommand{}); err == nil {
		t.Fatal("missing record should fail")
	}
	if StatusLabel(model.StatusPublished) != "已发布" {
		t.Fatal("published label mismatch")
	}
}

func TestImportRejectsLongBody(t *testing.T) {
	body := make([]rune, 10001)
	for i := range body {
		body[i] = 'x'
	}
	if len(ValidateRegister(model.RegisterCommand{ID: "x", Community: "c", Title: "t", Body: string(body)})) == 0 {
		t.Fatal("long body should be rejected")
	}
}
