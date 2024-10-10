package postgres

import (
	"database/sql"

	"github.com/taufiksty/be-socket/pkg/domain/entity"
	"github.com/taufiksty/be-socket/pkg/domain/repository"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	query := `INSERT INTO users (id, email, username, password) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.Password)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE username = $1`
	row := r.db.QueryRow(query, username)

	var user entity.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
