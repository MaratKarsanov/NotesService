package users

import (
	"database/sql"
	"errors"
)

type UsersRepository interface {
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(email, hashed_password string) error
	UpdateUser(id int, email, hashedPassword string) error
	DeleteUser(id int) error
}

type UsersDbRepository struct {
	db *sql.DB
}

func NewUsersDbRepository(db *sql.DB) *UsersDbRepository {
	return &UsersDbRepository{db: db}
}

func (r *UsersDbRepository) GetUserById(id int) (*User, error) {
	var user User
	err := r.db.QueryRow(`SELECT * FROM users WHERE id = $1`, id).
		Scan(&user.Id, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersDbRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(`SELECT * FROM users WHERE email = $1`, email).
		Scan(&user.Id, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersDbRepository) CreateUser(email, hashed_password string) error {
	_, err := r.db.Exec(
		`INSERT INTO users (email, hashed_password, created_at) VALUES ($1, $2, NOW())`,
		email,
		hashed_password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersDbRepository) UpdateUser(id int, email, hashedPassword string) error {
	res, err := r.db.Exec(
		`UPDATE users SET email = $1, hashed_password = $2 WHERE id = $3`,
		email,
		hashedPassword,
		id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("UserNotFound")
	}

	return nil
}

func (r *UsersDbRepository) DeleteUser(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
