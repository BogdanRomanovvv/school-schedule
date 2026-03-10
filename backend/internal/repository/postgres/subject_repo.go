package postgres

import (
	"school-schedule/internal/domain"

	"github.com/jmoiron/sqlx"
)

type SubjectRepo struct{ db *sqlx.DB }

func NewSubjectRepo(db *sqlx.DB) *SubjectRepo { return &SubjectRepo{db} }

func (r *SubjectRepo) GetAll() ([]domain.Subject, error) {
	subjects := make([]domain.Subject, 0)
	err := r.db.Select(&subjects, `SELECT id, name FROM subjects ORDER BY name`)
	return subjects, err
}

func (r *SubjectRepo) GetByID(id int) (*domain.Subject, error) {
	var s domain.Subject
	err := r.db.Get(&s, `SELECT id, name FROM subjects WHERE id=$1`, id)
	return &s, err
}

func (r *SubjectRepo) Create(s *domain.Subject) error {
	return r.db.QueryRow(
		`INSERT INTO subjects(name) VALUES($1) RETURNING id`, s.Name,
	).Scan(&s.ID)
}

func (r *SubjectRepo) Update(s *domain.Subject) error {
	_, err := r.db.Exec(`UPDATE subjects SET name=$1 WHERE id=$2`, s.Name, s.ID)
	return err
}

func (r *SubjectRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM subjects WHERE id=$1`, id)
	return err
}
