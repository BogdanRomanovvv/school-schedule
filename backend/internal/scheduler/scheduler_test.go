package scheduler_test

import (
	"testing"

	"school-schedule/internal/domain"
	"school-schedule/internal/scheduler"
)

// Минимальный сценарий: 1 класс, 2 предмета, 2 учителя, по 2 часа каждый
func TestGenerate_Simple(t *testing.T) {
	curricula := []domain.Curriculum{
		{ClassID: 1, SubjectID: 1, HoursPerWeek: 2},
		{ClassID: 1, SubjectID: 2, HoursPerWeek: 2},
	}
	teachers := []domain.Teacher{
		{ID: 1, Name: "Учитель1", MaxHoursPerWeek: 20},
		{ID: 2, Name: "Учитель2", MaxHoursPerWeek: 20},
	}
	teacherSubjects := map[int][]int{
		1: {1},
		2: {2},
	}

	gen := scheduler.NewGenerator(curricula, teachers, teacherSubjects)
	entries, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if len(entries) != 4 {
		t.Errorf("expected 4 lessons, got %d", len(entries))
	}
	assertNoConflicts(t, entries)
}

// Несколько классов, один предмет, один учитель
func TestGenerate_MultiClass(t *testing.T) {
	curricula := []domain.Curriculum{
		{ClassID: 1, SubjectID: 1, HoursPerWeek: 3},
		{ClassID: 2, SubjectID: 1, HoursPerWeek: 3},
		{ClassID: 3, SubjectID: 1, HoursPerWeek: 3},
	}
	teachers := []domain.Teacher{
		{ID: 10, Name: "Математик", MaxHoursPerWeek: 30},
	}
	teacherSubjects := map[int][]int{10: {1}}

	gen := scheduler.NewGenerator(curricula, teachers, teacherSubjects)
	entries, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if len(entries) != 9 {
		t.Errorf("expected 9 lessons, got %d", len(entries))
	}
	assertNoConflicts(t, entries)
}

// Нагрузка учителя превышена
func TestGenerate_OverloadedTeacher(t *testing.T) {
	curricula := []domain.Curriculum{
		{ClassID: 1, SubjectID: 1, HoursPerWeek: 5},
		{ClassID: 2, SubjectID: 1, HoursPerWeek: 5},
		{ClassID: 3, SubjectID: 1, HoursPerWeek: 5},
		{ClassID: 4, SubjectID: 1, HoursPerWeek: 5},
	}
	teachers := []domain.Teacher{
		{ID: 1, Name: "Один", MaxHoursPerWeek: 3}, // явно меньше нужного
	}
	teacherSubjects := map[int][]int{1: {1}}

	gen := scheduler.NewGenerator(curricula, teachers, teacherSubjects)
	_, err := gen.Generate()
	if err == nil {
		t.Fatal("expected error for overloaded teacher, got nil")
	}
}

// Проверка отсутствия конфликтов
func assertNoConflicts(t *testing.T, entries []domain.ScheduleEntry) {
	t.Helper()
	classSlot := make(map[[3]int]bool)
	teacherSlot := make(map[[3]int]bool)

	for _, e := range entries {
		ck := [3]int{e.ClassID, e.Day, e.LessonNumber}
		if classSlot[ck] {
			t.Errorf("conflict: class %d has two lessons on day %d lesson %d", e.ClassID, e.Day, e.LessonNumber)
		}
		classSlot[ck] = true

		tk := [3]int{e.TeacherID, e.Day, e.LessonNumber}
		if teacherSlot[tk] {
			t.Errorf("conflict: teacher %d has two lessons on day %d lesson %d", e.TeacherID, e.Day, e.LessonNumber)
		}
		teacherSlot[tk] = true
	}
}
