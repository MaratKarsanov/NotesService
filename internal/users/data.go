package users

import (
	"database/sql"
	"errors"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetUserById(id int) (*User, error) {
	var user User
	err := r.db.QueryRow(`SELECT * FROM users WHERE id = $1`, id).
		Scan(&user.Id, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(`SELECT * FROM users WHERE email = $1`, email).
		Scan(&user.Id, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) CreateUser(email, hashed_password string) error {
	_, err := r.db.Exec(
		`INSERT INTO users (email, hashed_password, created_at) VALUES ($1, $2, NOW())`,
		email,
		hashed_password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) UpdateUser(id int, word, translate string) error {
	res, err := r.db.Exec(
		`UPDATE ru_en SET title = $1, translation = $2 WHERE id = $3`,
		word,
		translate,
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

func (r *UsersRepository) DeleteUser(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
