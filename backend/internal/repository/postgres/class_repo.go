package postgres

import (
	"school-schedule/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ClassRepo struct{ db *sqlx.DB }

func NewClassRepo(db *sqlx.DB) *ClassRepo { return &ClassRepo{db} }

func (r *ClassRepo) GetAll() ([]domain.Class, error) {
	classes := make([]domain.Class, 0)
	err := r.db.Select(&classes, `SELECT id, name FROM classes ORDER BY name`)
	return classes, err
}

func (r *ClassRepo) GetByID(id int) (*domain.Class, error) {
	var c domain.Class
	err := r.db.Get(&c, `SELECT id, name FROM classes WHERE id=$1`, id)
	return &c, err
}

func (r *ClassRepo) Create(c *domain.Class) error {
	return r.db.QueryRow(
		`INSERT INTO classes(name) VALUES($1) RETURNING id`, c.Name,
	).Scan(&c.ID)
}

func (r *ClassRepo) Update(c *domain.Class) error {
	_, err := r.db.Exec(`UPDATE classes SET name=$1 WHERE id=$2`, c.Name, c.ID)
	return err
}

func (r *ClassRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM classes WHERE id=$1`, id)
	return err
}
