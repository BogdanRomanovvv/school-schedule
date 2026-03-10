package postgres

import (
	"school-schedule/internal/domain"

	"github.com/jmoiron/sqlx"
)

type CurriculumRepo struct{ db *sqlx.DB }

func NewCurriculumRepo(db *sqlx.DB) *CurriculumRepo { return &CurriculumRepo{db} }

func (r *CurriculumRepo) GetAll() ([]domain.Curriculum, error) {
	list := make([]domain.Curriculum, 0)
	err := r.db.Select(&list, `SELECT id, class_id, subject_id, hours_per_week FROM curriculum ORDER BY class_id, subject_id`)
	return list, err
}

func (r *CurriculumRepo) GetByClass(classID int) ([]domain.Curriculum, error) {
	list := make([]domain.Curriculum, 0)
	err := r.db.Select(&list,
		`SELECT id, class_id, subject_id, hours_per_week FROM curriculum WHERE class_id=$1`,
		classID)
	return list, err
}

func (r *CurriculumRepo) Upsert(c *domain.Curriculum) error {
	return r.db.QueryRow(`
		INSERT INTO curriculum(class_id, subject_id, hours_per_week)
		VALUES($1,$2,$3)
		ON CONFLICT (class_id, subject_id) DO UPDATE SET hours_per_week=EXCLUDED.hours_per_week
		RETURNING id`, c.ClassID, c.SubjectID, c.HoursPerWeek).Scan(&c.ID)
}

func (r *CurriculumRepo) Delete(classID, subjectID int) error {
	_, err := r.db.Exec(
		`DELETE FROM curriculum WHERE class_id=$1 AND subject_id=$2`,
		classID, subjectID,
	)
	return err
}
