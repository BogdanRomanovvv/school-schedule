package postgres

import (
	"school-schedule/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ScheduleRepo struct{ db *sqlx.DB }

func NewScheduleRepo(db *sqlx.DB) *ScheduleRepo { return &ScheduleRepo{db} }

const richQuery = `
	SELECT
		s.id, s.class_id, s.subject_id, s.teacher_id, s.day, s.lesson_number,
		c.name AS class_name,
		sub.name AS subject_name,
		t.name AS teacher_name
	FROM schedule s
	JOIN classes c ON c.id = s.class_id
	JOIN subjects sub ON sub.id = s.subject_id
	JOIN teachers t ON t.id = s.teacher_id
`

func (r *ScheduleRepo) GetAll() ([]domain.ScheduleEntryRich, error) {
	list := make([]domain.ScheduleEntryRich, 0)
	err := r.db.Select(&list, richQuery+` ORDER BY s.class_id, s.day, s.lesson_number`)
	return list, err
}

func (r *ScheduleRepo) GetByClass(classID int) ([]domain.ScheduleEntryRich, error) {
	list := make([]domain.ScheduleEntryRich, 0)
	err := r.db.Select(&list, richQuery+` WHERE s.class_id=$1 ORDER BY s.day, s.lesson_number`, classID)
	return list, err
}

func (r *ScheduleRepo) GetByTeacher(teacherID int) ([]domain.ScheduleEntryRich, error) {
	list := make([]domain.ScheduleEntryRich, 0)
	err := r.db.Select(&list, richQuery+` WHERE s.teacher_id=$1 ORDER BY s.day, s.lesson_number`, teacherID)
	return list, err
}

func (r *ScheduleRepo) SaveAll(entries []domain.ScheduleEntry) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM schedule`); err != nil {
		tx.Rollback()
		return err
	}
	for _, e := range entries {
		if _, err = tx.Exec(
			`INSERT INTO schedule(class_id, subject_id, teacher_id, day, lesson_number)
			 VALUES($1,$2,$3,$4,$5)`,
			e.ClassID, e.SubjectID, e.TeacherID, e.Day, e.LessonNumber,
		); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *ScheduleRepo) Clear() error {
	_, err := r.db.Exec(`DELETE FROM schedule`)
	return err
}

func (r *ScheduleRepo) UpdateEntry(e *domain.ScheduleEntry) error {
	_, err := r.db.Exec(
		`UPDATE schedule SET class_id=$1, subject_id=$2, teacher_id=$3, day=$4, lesson_number=$5 WHERE id=$6`,
		e.ClassID, e.SubjectID, e.TeacherID, e.Day, e.LessonNumber, e.ID,
	)
	return err
}
