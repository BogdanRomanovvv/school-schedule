package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"school-schedule/internal/domain"
	"school-schedule/internal/service"

	"github.com/go-chi/chi/v5"
)

type TeacherHandler struct{ svc *service.TeacherService }

func NewTeacherHandler(s *service.TeacherService) *TeacherHandler { return &TeacherHandler{s} }

func (h *TeacherHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.GetAll()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, list)
}

func (h *TeacherHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	t, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, t)
}

func (h *TeacherHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t domain.Teacher
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err := h.svc.Create(&t); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, t)
}

func (h *TeacherHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	var t domain.Teacher
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	t.ID = id
	if err := h.svc.Update(&t); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, t)
}

func (h *TeacherHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *TeacherHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	subjects, err := h.svc.GetSubjects(id)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, subjects)
}

func (h *TeacherHandler) AssignSubject(w http.ResponseWriter, r *http.Request) {
	tID, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid teacher id")
		return
	}
	sID, err := strconv.Atoi(chi.URLParam(r, "subjectId"))
	if err != nil {
		writeError(w, 400, "invalid subject id")
		return
	}
	if err := h.svc.AssignSubject(tID, sID); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(204)
}

func (h *TeacherHandler) RemoveSubject(w http.ResponseWriter, r *http.Request) {
	tID, err := parseID(r)
	if err != nil {
		writeError(w, 400, "invalid teacher id")
		return
	}
	sID, err := strconv.Atoi(chi.URLParam(r, "subjectId"))
	if err != nil {
		writeError(w, 400, "invalid subject id")
		return
	}
	if err := h.svc.RemoveSubject(tID, sID); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(204)
}
