package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"school-schedule/internal/api"
	"school-schedule/internal/api/handlers"
	"school-schedule/internal/repository/postgres"
	"school-schedule/internal/service"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "school_schedule"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	// Repositories
	classRepo := postgres.NewClassRepo(db)
	subjectRepo := postgres.NewSubjectRepo(db)
	teacherRepo := postgres.NewTeacherRepo(db)
	curriculumRepo := postgres.NewCurriculumRepo(db)
	scheduleRepo := postgres.NewScheduleRepo(db)

	// Services
	classSvc := service.NewClassService(classRepo)
	subjectSvc := service.NewSubjectService(subjectRepo)
	teacherSvc := service.NewTeacherService(teacherRepo)
	curriculumSvc := service.NewCurriculumService(curriculumRepo)
	scheduleSvc := service.NewScheduleService(scheduleRepo, curriculumRepo, teacherRepo)

	// Handlers
	classH := handlers.NewClassHandler(classSvc)
	subjectH := handlers.NewSubjectHandler(subjectSvc)
	teacherH := handlers.NewTeacherHandler(teacherSvc)
	curriculumH := handlers.NewCurriculumHandler(curriculumSvc)
	scheduleH := handlers.NewScheduleHandler(scheduleSvc)

	router := api.NewRouter(classH, subjectH, teacherH, curriculumH, scheduleH)

	addr := ":" + getEnv("PORT", "8080")
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
