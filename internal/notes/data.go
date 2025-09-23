package notes

import (
	"database/sql"
	"errors"
)

type NotesRepository interface {
	GetNote(id int) (*Note, error)
	GetUserNotes(userid int) (*[]Note, error)
	CreateNote(user_id int, title, body string) error
	UpdateNote(id int, title, body string) error
	DeleteNote(id int) error
}

type NotesDbRepository struct {
	db *sql.DB
}

func NewNotesDbRepository(db *sql.DB) *NotesDbRepository {
	return &NotesDbRepository{db: db}
}

func (r *NotesDbRepository) GetNote(id int) (*Note, error) {
	var note Note
	err := r.db.QueryRow(`SELECT * FROM notes WHERE id = $1`, id).
		Scan(&note.Id, &note.UserId, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *NotesDbRepository) GetUserNotes(userid int) (*[]Note, error) {
	var notes []Note
	rows, err := r.db.Query(`SELECT * FROM notes WHERE user_id = $1`, userid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.Id, &note.UserId, &note.Title, &note.Body, &note.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return &notes, nil
}

func (r *NotesDbRepository) CreateNote(user_id int, title, body string) error {
	_, err := r.db.Exec(
		`INSERT INTO notes (user_id, title, body, created_at) VALUES ($1, $2, $3, NOW())`,
		user_id,
		title,
		body)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotesDbRepository) UpdateNote(id int, title, body string) error {
	res, err := r.db.Exec(
		`UPDATE notes SET title = $1, body = $2 WHERE id = $3`,
		title,
		body,
		id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("NoteNotFound")
	}

	return nil
}

func (r *NotesDbRepository) DeleteNote(id int) error {
	_, err := r.db.Exec(`DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
