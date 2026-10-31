package repository

import (
	"api-bewareOfSharkAttack/entity"
	"database/sql"
	"errors"
)

type ActivityRepository struct {
	DB *sql.DB
}

func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{DB: db}
}

func (r *ActivityRepository) GetAll() ([]entity.Activity, error) {
	rows, err := r.DB.Query("SELECT id, name, category, created_at FROM activities ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []entity.Activity
	for rows.Next() {
		var act entity.Activity
		if err := rows.Scan(&act.ID, &act.Name, &act.Category, &act.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, act)
	}
	return activities, nil
}

func (r *ActivityRepository) GetByID(id int) (*entity.Activity, error) {
	var act entity.Activity
	query := "SELECT id, name, category, created_at FROM activities WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&act.ID, &act.Name, &act.Category, &act.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("data aktivitas tidak ditemukan")
	}
	return &act, err
}

func (r *ActivityRepository) Create(act *entity.Activity) error {
	query := "INSERT INTO activities (name, category) VALUES ($1, $2) RETURNING id, created_at"
	return r.DB.QueryRow(query, act.Name, act.Category).Scan(&act.ID, &act.CreatedAt)
}

func (r *ActivityRepository) Update(id int, act *entity.Activity) error {
	query := "UPDATE activities SET name = $1, category = $2 WHERE id = $3"
	res, err := r.DB.Exec(query, act.Name, act.Category, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data aktivitas tidak ditemukan")
	}
	act.ID = id
	return nil
}

func (r *ActivityRepository) Delete(id int) error {
	query := "DELETE FROM activities WHERE id = $1"
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data aktivitas tidak ditemukan")
	}
	return nil
}
