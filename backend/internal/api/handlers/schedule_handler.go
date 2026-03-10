package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"school-schedule/internal/domain"
	"school-schedule/internal/service"
)

type ScheduleHandler struct{ svc *service.ScheduleService }

func NewScheduleHandler(s *service.ScheduleService) *ScheduleHandler { return &ScheduleHandler{s} }

func (h *ScheduleHandler) Generate(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.Generate()
	if err != nil {
		writeError(w, 422, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{
		"message": "Расписание успешно сгенерировано",
		"count":   len(entries),
	})
}

func (h *ScheduleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.GetAll()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, list)
}

func (h *ScheduleHandler) GetByClass(w http.ResponseWriter, r *http.Request) {
	classID, err := strconv.Atoi(r.URL.Query().Get("class_id"))
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
}

func (h *ScheduleHandler) GetByTeacher(w http.ResponseWriter, r *http.Request) {
	teacherID, err := strconv.Atoi(r.URL.Query().Get("teacher_id"))
	if err != nil {
		writeError(w, 400, "invalid teacher_id")
		return
	}
	list, err := h.svc.GetByTeacher(teacherID)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, list)
}

func (h *ScheduleHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	var e domain.ScheduleEntry
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	if err := h.svc.UpdateEntry(&e); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, e)
}

func (h *ScheduleHandler) Clear(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Clear(); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(204)
}

func (h *ScheduleHandler) ExportByClass(w http.ResponseWriter, r *http.Request) {
	f, err := h.svc.ExportByClass()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="schedule_by_class.xlsx"`)
	f.Write(w)
}

func (h *ScheduleHandler) ExportByTeacher(w http.ResponseWriter, r *http.Request) {
	f, err := h.svc.ExportByTeacher()
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="schedule_by_teacher.xlsx"`)
	f.Write(w)
}
