package repository

import "school-schedule/internal/domain"

type ClassRepository interface {
	GetAll() ([]domain.Class, error)
	GetByID(id int) (*domain.Class, error)
	Create(c *domain.Class) error
	Update(c *domain.Class) error
	Delete(id int) error
}

type SubjectRepository interface {
	GetAll() ([]domain.Subject, error)
	GetByID(id int) (*domain.Subject, error)
	Create(s *domain.Subject) error
	Update(s *domain.Subject) error
	Delete(id int) error
}

type TeacherRepository interface {
	GetAll() ([]domain.Teacher, error)
	GetByID(id int) (*domain.Teacher, error)
	Create(t *domain.Teacher) error
	Update(t *domain.Teacher) error
	Delete(id int) error
	GetSubjects(teacherID int) ([]domain.Subject, error)
	AssignSubject(teacherID, subjectID int) error
	RemoveSubject(teacherID, subjectID int) error
	GetTeachersBySubject(subjectID int) ([]domain.Teacher, error)
}

type CurriculumRepository interface {
	GetAll() ([]domain.Curriculum, error)
	GetByClass(classID int) ([]domain.Curriculum, error)
	Upsert(c *domain.Curriculum) error
	Delete(classID, subjectID int) error
}

type ScheduleRepository interface {
	GetAll() ([]domain.ScheduleEntryRich, error)
	GetByClass(classID int) ([]domain.ScheduleEntryRich, error)
	GetByTeacher(teacherID int) ([]domain.ScheduleEntryRich, error)
	SaveAll(entries []domain.ScheduleEntry) error
	Clear() error
	UpdateEntry(e *domain.ScheduleEntry) error
}
