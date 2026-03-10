package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"school-schedule/internal/domain"
	"school-schedule/internal/service"

	"github.com/go-chi/chi/v5"
)

// ─── helpers ─────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseID(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}

// ─── Classes ──────────────────────────────────────────────────────────────────

type ClassHandler struct{ svc *service.ClassService }

func NewClassHandler(s *service.ClassService) *ClassHandler { return &ClassHandler{s} }

func (h *ClassHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.GetAll()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, list)
}

func (h *ClassHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	c, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, c)
}

func (h *ClassHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c domain.Class
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err := h.svc.Create(&c); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, c)
}

func (h *ClassHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	var c domain.Class
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	c.ID = id
	if err := h.svc.Update(&c); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, c)
}

func (h *ClassHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(204)
}

// ─── Subjects ────────────────────────────────────────────────────────────────

type SubjectHandler struct{ svc *service.SubjectService }

func NewSubjectHandler(s *service.SubjectService) *SubjectHandler { return &SubjectHandler{s} }

func (h *SubjectHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.GetAll()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, list)
}

func (h *SubjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	s, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, s)
}

func (h *SubjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var s domain.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err := h.svc.Create(&s); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, s)
}

func (h *SubjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	var s domain.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	s.ID = id
	if err := h.svc.Update(&s); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s)
}

func (h *SubjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(204)
}
