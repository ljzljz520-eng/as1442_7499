package attachment

import (
	"path/filepath"
	"testing"

	"noticeword/internal/store"
)

func TestManagerPersistsAndTotals(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "attachment.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	manager := NewManager(db)
	first, err := manager.Add("record", "brief.txt", ContentType("brief.txt"), 4, Checksum([]byte("data")))
	if err != nil {
		t.Fatal(err)
	}
	if !Verify([]byte("data"), first.Checksum) {
		t.Fatal("checksum mismatch")
	}
	if size, err := manager.TotalSize("record"); err != nil || size != 4 {
		t.Fatalf("size=%d err=%v", size, err)
	}
	if err := manager.Validate(first); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Add("", "bad", "text/plain", 1, ""); err == nil {
		t.Fatal("missing record should fail")
	}
}

func TestContentTypes(t *testing.T) {
	if ContentType("a.csv") != "text/csv" || ContentType("a.bin") != "application/octet-stream" {
		t.Fatal("content type detection failed")
	}
}
