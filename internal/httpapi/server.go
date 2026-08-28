package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"noticeword/internal/flow025"
	"noticeword/internal/model"
	"noticeword/internal/service"
)

type Server struct {
	service *service.Service
	flow    *flow025.Handler
}

func NewServer(app *service.Service) *Server {
	return &Server{service: app, flow: flow025.NewHandler(app)}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/records", s.records)
	mux.HandleFunc("/records/", s.record)
	mux.HandleFunc("/search", s.search)
	return requestLogger(mux)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "noticeword"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.register(w, r)
		return
	}
	if r.Method == http.MethodGet {
		s.list(w, r)
		return
	}
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/records/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "record id is required")
		return
	}
	if r.Method == http.MethodGet {
		s.get(w, r, id)
		return
	}
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/review") {
		s.review(w, r, strings.TrimSuffix(id, "/review"))
		return
	}
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/publish") {
		s.publish(w, r, strings.TrimSuffix(id, "/publish"))
		return
	}
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/archive") {
		s.archive(w, r, strings.TrimSuffix(id, "/archive"))
		return
	}
	if r.Method == http.MethodPatch {
		s.change(w, r, id)
		return
	}
	writeError(w, http.StatusNotFound, "route not found")
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var command model.RegisterCommand
	if !decode(r, &command) {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if issues := service.ValidateRegister(command); len(issues) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"errors": issues})
		return
	}
	record, err := s.service.Register(command)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	query := model.SearchQuery{Community: r.URL.Query().Get("community"), Text: r.URL.Query().Get("q"), Status: model.RecordStatus(r.URL.Query().Get("status"))}
	query.Page = parseInt(r.URL.Query().Get("page"), 1)
	query.PageSize = parseInt(r.URL.Query().Get("page_size"), 10)
	page, err := s.flow.Search(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	s.list(w, r)
}

func (s *Server) get(w http.ResponseWriter, r *http.Request, id string) {
	record, err := s.service.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) review(w http.ResponseWriter, r *http.Request, id string) {
	var command model.ReviewCommand
	if !decode(r, &command) {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	command.RecordID = id
	record, err := s.service.Review(command)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) publish(w http.ResponseWriter, r *http.Request, id string) {
	var command model.PublishCommand
	if !decode(r, &command) {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	command.RecordID = id
	record, err := s.service.Publish(command)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) archive(w http.ResponseWriter, r *http.Request, id string) {
	var command model.ArchiveCommand
	if !decode(r, &command) {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	command.RecordID = id
	record, err := s.service.Archive(command)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) change(w http.ResponseWriter, r *http.Request, id string) {
	var command model.ChangeCommand
	if !decode(r, &command) {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	command.RecordID = id
	if err := service.ValidateChange(command); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	record, err := s.service.Change(command)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func decode(r *http.Request, target any) bool {
	return json.NewDecoder(r.Body).Decode(target) == nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
