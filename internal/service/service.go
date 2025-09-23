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

	usersRepository users.UsersRepository
	notesRepository notes.NotesRepository
}

func NewService(
	logger echo.Logger,
	notesRepository notes.NotesRepository,
	usersRepository users.UsersRepository) *Service {
	svc := &Service{
		logger:          logger,
		usersRepository: usersRepository,
		notesRepository: notesRepository,
	}

	return svc
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
