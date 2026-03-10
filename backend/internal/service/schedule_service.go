package service

import (
	"fmt"

	"school-schedule/internal/domain"
	"school-schedule/internal/repository"
	"school-schedule/internal/scheduler"

	"github.com/xuri/excelize/v2"
)

type ScheduleService struct {
	schedRepo   repository.ScheduleRepository
	currRepo    repository.CurriculumRepository
	teacherRepo repository.TeacherRepository
}

func NewScheduleService(
	sr repository.ScheduleRepository,
	cr repository.CurriculumRepository,
	tr repository.TeacherRepository,
) *ScheduleService {
	return &ScheduleService{schedRepo: sr, currRepo: cr, teacherRepo: tr}
}

func (s *ScheduleService) Generate() ([]domain.ScheduleEntry, error) {
	curricula, err := s.currRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("curriculum load: %w", err)
	}

	teachers, err := s.teacherRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("teachers load: %w", err)
	}

	teacherSubjects := make(map[int][]int) // teacherID → []subjectID
	for _, t := range teachers {
		subjects, err := s.teacherRepo.GetSubjects(t.ID)
		if err != nil {
			return nil, err
		}
		ids := make([]int, len(subjects))
		for i, sub := range subjects {
			ids[i] = sub.ID
		}
		teacherSubjects[t.ID] = ids
	}

	gen := scheduler.NewGenerator(curricula, teachers, teacherSubjects)
	entries, err := gen.Generate()
	if err != nil {
		return nil, err
	}

	if err = s.schedRepo.SaveAll(entries); err != nil {
		return nil, fmt.Errorf("save schedule: %w", err)
	}
	return entries, nil
}

func (s *ScheduleService) GetAll() ([]domain.ScheduleEntryRich, error) {
	return s.schedRepo.GetAll()
}

func (s *ScheduleService) GetByClass(classID int) ([]domain.ScheduleEntryRich, error) {
	return s.schedRepo.GetByClass(classID)
}

func (s *ScheduleService) GetByTeacher(teacherID int) ([]domain.ScheduleEntryRich, error) {
	return s.schedRepo.GetByTeacher(teacherID)
}

func (s *ScheduleService) UpdateEntry(e *domain.ScheduleEntry) error {
	return s.schedRepo.UpdateEntry(e)
}

func (s *ScheduleService) Clear() error {
	return s.schedRepo.Clear()
}

// ─── Excel Export ─────────────────────────────────────────────────────────────

var dayNames = []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница"}

// ExportByClass генерирует Excel-файл с расписанием по классам.
func (s *ScheduleService) ExportByClass() (*excelize.File, error) {
	entries, err := s.schedRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// класс → день → урок → запись
	grid := make(map[string]map[int]map[int]domain.ScheduleEntryRich)
	for _, e := range entries {
		if grid[e.ClassName] == nil {
			grid[e.ClassName] = make(map[int]map[int]domain.ScheduleEntryRich)
		}
		if grid[e.ClassName][e.Day] == nil {
			grid[e.ClassName][e.Day] = make(map[int]domain.ScheduleEntryRich)
		}
		grid[e.ClassName][e.Day][e.LessonNumber] = e
	}

	f := excelize.NewFile()
	for className, days := range grid {
		f.NewSheet(className)
		f.SetCellValue(className, "A1", "Урок / День")
		for d, name := range dayNames {
			col := string(rune('B' + d))
			f.SetCellValue(className, col+"1", name)
		}
		for lesson := 0; lesson <= 6; lesson++ {
			row := fmt.Sprintf("A%d", lesson+2)
			f.SetCellValue(className, row, fmt.Sprintf("Урок %d", lesson+1))
			for day := 0; day <= 4; day++ {
				col := string(rune('B' + day))
				cell := fmt.Sprintf("%s%d", col, lesson+2)
				if e, ok := days[day][lesson]; ok {
					f.SetCellValue(className, cell, fmt.Sprintf("%s\n(%s)", e.SubjectName, e.TeacherName))
				}
			}
		}
	}
	f.DeleteSheet("Sheet1")
	return f, nil
}

// ExportByTeacher генерирует Excel-файл с расписанием по учителям.
func (s *ScheduleService) ExportByTeacher() (*excelize.File, error) {
	entries, err := s.schedRepo.GetAll()
	if err != nil {
		return nil, err
	}

	grid := make(map[string]map[int]map[int]domain.ScheduleEntryRich)
	for _, e := range entries {
		if grid[e.TeacherName] == nil {
			grid[e.TeacherName] = make(map[int]map[int]domain.ScheduleEntryRich)
		}
		if grid[e.TeacherName][e.Day] == nil {
			grid[e.TeacherName][e.Day] = make(map[int]domain.ScheduleEntryRich)
		}
		grid[e.TeacherName][e.Day][e.LessonNumber] = e
	}

	f := excelize.NewFile()
	for teacherName, days := range grid {
		f.NewSheet(teacherName)
		f.SetCellValue(teacherName, "A1", "Урок / День")
		for d, name := range dayNames {
			col := string(rune('B' + d))
			f.SetCellValue(teacherName, col+"1", name)
		}
		for lesson := 0; lesson <= 6; lesson++ {
			row := fmt.Sprintf("A%d", lesson+2)
			f.SetCellValue(teacherName, row, fmt.Sprintf("Урок %d", lesson+1))
			for day := 0; day <= 4; day++ {
				col := string(rune('B' + day))
				cell := fmt.Sprintf("%s%d", col, lesson+2)
				if e, ok := days[day][lesson]; ok {
					f.SetCellValue(teacherName, cell,
						fmt.Sprintf("%s (%s)", e.SubjectName, e.ClassName))
				}
			}
		}
	}
	f.DeleteSheet("Sheet1")
	return f, nil
}
