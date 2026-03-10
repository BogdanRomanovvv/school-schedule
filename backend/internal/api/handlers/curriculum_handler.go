package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"school-schedule/internal/domain"
	"school-schedule/internal/service"

	"github.com/go-chi/chi/v5"
)

type CurriculumHandler struct{ svc *service.CurriculumService }

func NewCurriculumHandler(s *service.CurriculumService) *CurriculumHandler {
	return &CurriculumHandler{s}
}

func (h *CurriculumHandler) List(w http.ResponseWriter, r *http.Request) {
	classIDStr := r.URL.Query().Get("class_id")
	if classIDStr != "" {
		classID, err := strconv.Atoi(classIDStr)
		if err != nil {
			writeError(w, 400, "invalid class_id")
			return
		}
		list, err := h.svc.GetByClass(classID)
		if err != nil {
			writeError(w, 500, err.Error())
			return
		}
		writeJSON(w, 200, list)
		return
	}
	list, err := h.svc.GetAll()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, list)
}

func (h *CurriculumHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var c domain.Curriculum
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err := h.svc.Upsert(&c); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, c)
}

func (h *CurriculumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	classID, err := strconv.Atoi(chi.URLParam(r, "classId"))
	if err != nil {
		writeError(w, 400, "invalid class id")
		return
	}
	subjectID, err := strconv.Atoi(chi.URLParam(r, "subjectId"))
	if err != nil {
		writeError(w, 400, "invalid subject id")
		return
	}
	if err := h.svc.Delete(classID, subjectID); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(204)
}
