package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type FavqsResponse struct {
	Quote struct {
		Body string `json:"body"`
	} `json:"quote"`
}

// localhost:8000/api/note/:id
func (s *Service) GetNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	notesRepository := s.notesRepository
	note, err := notesRepository.GetNote(id)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	s.logger.Infof("Note with id %s was given", id)
	return c.JSON(http.StatusOK, Response{Object: note})
}

// localhost:8000/api/notes
func (s *Service) GetUserNotes(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	email, err := user.Claims.GetSubject()
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	usersRepository := s.usersRepository
	dbUser, err := usersRepository.GetUserByEmail(email)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	notesRepository := s.notesRepository
	notes, err := notesRepository.GetUserNotes(dbUser.Id)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	s.logger.Infof("User %s took his notes", dbUser.Id)
	return c.JSON(http.StatusOK, Response{Object: notes})
}

// localhost:8000/api/note
func (s *Service) CreateNote(c echo.Context) error {
	var note Note
	err := c.Bind(&note)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	user := c.Get("user").(*jwt.Token)
	email, err := user.Claims.GetSubject()
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	resp, err := http.Get("https://favqs.com/api/qotd")
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}
	defer resp.Body.Close()

	var favqs FavqsResponse
	if err := json.NewDecoder(resp.Body).Decode(&favqs); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	note.Body = note.Body + "\nQuote of the day: " + favqs.Quote.Body

	notesRepository := s.notesRepository
	usersRepository := s.usersRepository
	dbUser, err := usersRepository.GetUserByEmail(email)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	err = notesRepository.CreateNote(dbUser.Id, note.Title, note.Body)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	s.logger.Infof("User %s created note", dbUser.Email)
	return c.String(http.StatusOK, "OK")
}

// localhost:8000/note/:id
func (s *Service) UpdateNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	var note Note
	err = c.Bind(&note)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	notesRepository := s.notesRepository
	err = notesRepository.UpdateNote(id, note.Title, note.Body)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	s.logger.Infof("Note with id %s was updated", id)
	return c.String(http.StatusOK, "OK")
}

// localhost:8000/note/:id
func (s *Service) DeleteNote(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	notesRepository := s.notesRepository
	err = notesRepository.DeleteNote(id)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	s.logger.Infof("Note with id %s was deleted", id)
	return c.String(http.StatusOK, "OK")
}
