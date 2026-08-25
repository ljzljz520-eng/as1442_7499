package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"noticeword/internal/service"
	"noticeword/internal/store"
)

func testServer(t *testing.T) (*Server, func()) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	app := service.New(db, store.NewSequenceClock("2025-01-01T00:00:00Z"))
	return NewServer(app), func() { db.Close() }
}

func TestHealthAndRegister(t *testing.T) {
	server, closeStore := testServer(t)
	defer closeStore()
	health := httptest.NewRecorder()
	server.Handler().ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("health=%d", health.Code)
	}
	body := `{"id":"api-1","community":"API社群","title":"公告","body":"API正文","description":"API说明","actor":"editor"}`
	created := httptest.NewRecorder()
	server.Handler().ServeHTTP(created, httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(body)))
	if created.Code != http.StatusCreated {
		t.Fatalf("create=%d body=%s", created.Code, created.Body.String())
	}
	var response map[string]any
	if err := json.Unmarshal(created.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["id"] != "api-1" {
		t.Fatalf("response=%v", response)
	}
}

func TestSearchRouteRejectsWrongMethod(t *testing.T) {
	server, closeStore := testServer(t)
	defer closeStore()
	result := httptest.NewRecorder()
	server.Handler().ServeHTTP(result, httptest.NewRequest(http.MethodPost, "/search", nil))
	if result.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d", result.Code)
	}
}
