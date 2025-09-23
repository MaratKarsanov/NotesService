package service

import (
	"NotesService/internal/notes"
	"NotesService/internal/users"
	"database/sql"

	"github.com/labstack/echo/v4"
)

const (
	InvalidParams       = "invalid params"
	InvalidCredentials  = "invalid credentials"
	InternalServerError = "internal error"
	UserAlreadyExists   = "user already exists"
)

type Service struct {
	db     *sql.DB
	logger echo.Logger

	usersRepository *users.UsersRepository
	notesRepository *notes.NotesRepository
}

func NewService(db *sql.DB, logger echo.Logger) *Service {
	svc := &Service{
		db:     db,
		logger: logger,
	}
	svc.initRepositories(db)

	return svc
}

func (s *Service) initRepositories(db *sql.DB) {
	s.usersRepository = users.NewUsersRepository(db)
	s.notesRepository = notes.NewNotesRepository(db)
}

type Response struct {
	Object       any    `json:"object,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func (r *Response) Error() string {
	return r.ErrorMessage
}

func (s *Service) NewError(err string) (int, *Response) {
	return 400, &Response{ErrorMessage: err}
}
