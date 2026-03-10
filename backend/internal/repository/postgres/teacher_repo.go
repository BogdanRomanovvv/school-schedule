package postgres

import (
	"school-schedule/internal/domain"

	"github.com/jmoiron/sqlx"
)

type TeacherRepo struct{ db *sqlx.DB }

func NewTeacherRepo(db *sqlx.DB) *TeacherRepo { return &TeacherRepo{db} }

func (r *TeacherRepo) GetAll() ([]domain.Teacher, error) {
	list := make([]domain.Teacher, 0)
	err := r.db.Select(&list, `SELECT id, name, max_hours_per_week, homeroom_class_id FROM teachers ORDER BY name`)
	return list, err
}

func (r *TeacherRepo) GetByID(id int) (*domain.Teacher, error) {
	var t domain.Teacher
	err := r.db.Get(&t, `SELECT id, name, max_hours_per_week, homeroom_class_id FROM teachers WHERE id=$1`, id)
	return &t, err
}

func (r *TeacherRepo) Create(t *domain.Teacher) error {
	return r.db.QueryRow(
		`INSERT INTO teachers(name, max_hours_per_week) VALUES($1,$2) RETURNING id`,
		t.Name, t.MaxHoursPerWeek,
	).Scan(&t.ID)
}

func (r *TeacherRepo) Update(t *domain.Teacher) error {
	_, err := r.db.Exec(
		`UPDATE teachers SET name=$1, max_hours_per_week=$2 WHERE id=$3`,
		t.Name, t.MaxHoursPerWeek, t.ID,
	)
	return err
}

func (r *TeacherRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM teachers WHERE id=$1`, id)
	return err
}

func (r *TeacherRepo) GetSubjects(teacherID int) ([]domain.Subject, error) {
	list := make([]domain.Subject, 0)
	err := r.db.Select(&list, `
		SELECT s.id, s.name FROM subjects s
		JOIN teacher_subjects ts ON ts.subject_id = s.id
		WHERE ts.teacher_id = $1 ORDER BY s.name`, teacherID)
	return list, err
}

func (r *TeacherRepo) AssignSubject(teacherID, subjectID int) error {
	_, err := r.db.Exec(
		`INSERT INTO teacher_subjects(teacher_id, subject_id) VALUES($1,$2) ON CONFLICT DO NOTHING`,
		teacherID, subjectID,
	)
	return err
}

func (r *TeacherRepo) RemoveSubject(teacherID, subjectID int) error {
	_, err := r.db.Exec(
		`DELETE FROM teacher_subjects WHERE teacher_id=$1 AND subject_id=$2`,
		teacherID, subjectID,
	)
	return err
}

func (r *TeacherRepo) GetTeachersBySubject(subjectID int) ([]domain.Teacher, error) {
	list := make([]domain.Teacher, 0)
	err := r.db.Select(&list, `
		SELECT t.id, t.name, t.max_hours_per_week FROM teachers t
		JOIN teacher_subjects ts ON ts.teacher_id = t.id
		WHERE ts.subject_id = $1`, subjectID)
	return list, err
}
