package export

import (
	"testing"

	"noticeword/internal/model"
)

func TestEncodeDecode(t *testing.T) {
	input := []model.Record{{ID: "r", Community: "社群", Title: "标题", Body: "正文,带逗号", Description: "说明", Status: model.StatusDraft, CharacterCount: 6, Revision: 2}}
	encoded, err := Encode(input)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Body != input[0].Body {
		t.Fatalf("rows=%#v", rows)
	}
	if _, err := Decode("wrong\n"); err == nil {
		t.Fatal("bad header should fail")
	}
}
