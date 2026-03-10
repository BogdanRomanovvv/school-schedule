package service

import (
	"school-schedule/internal/domain"
	"school-schedule/internal/repository"
)

type ClassService struct{ repo repository.ClassRepository }

func NewClassService(r repository.ClassRepository) *ClassService { return &ClassService{r} }

func (s *ClassService) GetAll() ([]domain.Class, error)       { return s.repo.GetAll() }
func (s *ClassService) GetByID(id int) (*domain.Class, error) { return s.repo.GetByID(id) }
func (s *ClassService) Create(c *domain.Class) error          { return s.repo.Create(c) }
func (s *ClassService) Update(c *domain.Class) error          { return s.repo.Update(c) }
func (s *ClassService) Delete(id int) error                   { return s.repo.Delete(id) }

type SubjectService struct{ repo repository.SubjectRepository }

func NewSubjectService(r repository.SubjectRepository) *SubjectService { return &SubjectService{r} }

func (s *SubjectService) GetAll() ([]domain.Subject, error)       { return s.repo.GetAll() }
func (s *SubjectService) GetByID(id int) (*domain.Subject, error) { return s.repo.GetByID(id) }
func (s *SubjectService) Create(sub *domain.Subject) error        { return s.repo.Create(sub) }
func (s *SubjectService) Update(sub *domain.Subject) error        { return s.repo.Update(sub) }
func (s *SubjectService) Delete(id int) error                     { return s.repo.Delete(id) }

type TeacherService struct{ repo repository.TeacherRepository }

func NewTeacherService(r repository.TeacherRepository) *TeacherService { return &TeacherService{r} }

func (s *TeacherService) GetAll() ([]domain.Teacher, error)       { return s.repo.GetAll() }
func (s *TeacherService) GetByID(id int) (*domain.Teacher, error) { return s.repo.GetByID(id) }
func (s *TeacherService) Create(t *domain.Teacher) error          { return s.repo.Create(t) }
func (s *TeacherService) Update(t *domain.Teacher) error          { return s.repo.Update(t) }
func (s *TeacherService) Delete(id int) error                     { return s.repo.Delete(id) }
func (s *TeacherService) GetSubjects(teacherID int) ([]domain.Subject, error) {
	return s.repo.GetSubjects(teacherID)
}
func (s *TeacherService) AssignSubject(tID, sID int) error { return s.repo.AssignSubject(tID, sID) }
func (s *TeacherService) RemoveSubject(tID, sID int) error { return s.repo.RemoveSubject(tID, sID) }

type CurriculumService struct {
	repo repository.CurriculumRepository
}

func NewCurriculumService(r repository.CurriculumRepository) *CurriculumService {
	return &CurriculumService{r}
}

func (s *CurriculumService) GetAll() ([]domain.Curriculum, error) { return s.repo.GetAll() }
func (s *CurriculumService) GetByClass(classID int) ([]domain.Curriculum, error) {
	return s.repo.GetByClass(classID)
}
func (s *CurriculumService) Upsert(c *domain.Curriculum) error { return s.repo.Upsert(c) }
func (s *CurriculumService) Delete(classID, subjectID int) error {
	return s.repo.Delete(classID, subjectID)
}
