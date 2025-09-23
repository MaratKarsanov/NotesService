package notes

import (
	"database/sql"
	"errors"
)

type NotesRepository struct {
	db *sql.DB
}

func NewNotesRepository(db *sql.DB) *NotesRepository {
	return &NotesRepository{db: db}
}

func (r *NotesRepository) GetNote(id int) (*Note, error) {
	var note Note
	err := r.db.QueryRow(`SELECT * FROM notes WHERE id = $1`, id).
		Scan(&note.Id, &note.UserId, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *NotesRepository) GetUserNotes(userid int) (*[]Note, error) {
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

func (r *NotesRepository) CreateNote(user_id int, title, body string) error {
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

func (r *NotesRepository) UpdateNote(id int, title, body string) error {
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

func (r *NotesRepository) DeleteNote(id int) error {
	_, err := r.db.Exec(`DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
