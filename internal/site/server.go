package site

import (
	"encoding/json"
	"errors"
	"net/http"

	"example.com/pottery-experience/internal/booking"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	service   *booking.UpdateService
	staticDir string
}

func NewServer(service *booking.UpdateService, staticDir string) *Server {
	return &Server{service: service, staticDir: staticDir}
}

func (s *Server) Handler() http.Handler {
	router := chi.NewRouter()
	router.Get("/api/records/{id}", s.getRecord)
	router.Post("/api/records/{id}/confirm", s.confirm)
	router.Get("/api/records/{id}/summary", s.summary)
	router.Handle("/*", s.staticHandler())
	return router
}

func (s *Server) getRecord(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, err := s.service.GetSummary(id)
	if errors.Is(err, booking.ErrRecordNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	s.getRecord(w, r)
}

func (s *Server) confirm(w http.ResponseWriter, r *http.Request) {
	var change booking.OperatorChange
	if err := json.NewDecoder(r.Body).Decode(&change); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	record, err := s.service.Confirm(chi.URLParam(r, "id"), change)
	if errors.Is(err, booking.ErrRecordNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if errors.Is(err, booking.ErrInvalidChange) {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record.Summary())
}

func (s *Server) staticHandler() http.Handler {
	return http.FileServer(http.Dir(s.staticDir))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
