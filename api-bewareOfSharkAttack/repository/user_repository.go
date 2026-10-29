package repository

import (
	"api-bewareOfSharkAttack/entity"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	query := `INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.DB.QueryRow(query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	query := `SELECT id, username, email, password, role, created_at FROM users WHERE email = $1`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
