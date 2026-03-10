package api

import (
	"school-schedule/internal/api/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	classH *handlers.ClassHandler,
	subjectH *handlers.SubjectHandler,
	teacherH *handlers.TeacherHandler,
	curriculumH *handlers.CurriculumHandler,
	scheduleH *handlers.ScheduleHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
	}))

	r.Route("/api", func(r chi.Router) {
		// Classes
		r.Route("/classes", func(r chi.Router) {
			r.Get("/", classH.List)
			r.Post("/", classH.Create)
			r.Get("/{id}", classH.Get)
			r.Put("/{id}", classH.Update)
			r.Delete("/{id}", classH.Delete)
		})

		// Subjects
		r.Route("/subjects", func(r chi.Router) {
			r.Get("/", subjectH.List)
			r.Post("/", subjectH.Create)
			r.Get("/{id}", subjectH.Get)
			r.Put("/{id}", subjectH.Update)
			r.Delete("/{id}", subjectH.Delete)
		})

		// Teachers
		r.Route("/teachers", func(r chi.Router) {
			r.Get("/", teacherH.List)
			r.Post("/", teacherH.Create)
			r.Get("/{id}", teacherH.Get)
			r.Put("/{id}", teacherH.Update)
			r.Delete("/{id}", teacherH.Delete)
			r.Get("/{id}/subjects", teacherH.GetSubjects)
			r.Post("/{id}/subjects/{subjectId}", teacherH.AssignSubject)
			r.Delete("/{id}/subjects/{subjectId}", teacherH.RemoveSubject)
		})

		// Curriculum
		r.Route("/curriculum", func(r chi.Router) {
			r.Get("/", curriculumH.List)
			r.Post("/", curriculumH.Upsert)
			r.Delete("/{classId}/{subjectId}", curriculumH.Delete)
		})

		// Schedule
		r.Route("/schedule", func(r chi.Router) {
			r.Get("/", scheduleH.GetAll)
			r.Post("/generate", scheduleH.Generate)
			r.Delete("/", scheduleH.Clear)
			r.Put("/entry", scheduleH.UpdateEntry)
			r.Get("/by-class", scheduleH.GetByClass)
			r.Get("/by-teacher", scheduleH.GetByTeacher)
			r.Get("/export/class", scheduleH.ExportByClass)
			r.Get("/export/teacher", scheduleH.ExportByTeacher)
		})
	})

	return r
}
