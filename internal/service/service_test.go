package service_test

import (
	"NotesService/internal/notes"
	"NotesService/internal/service"
	"NotesService/internal/users"
	"NotesService/pkg/logs"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotesRepository struct {
	mock.Mock
}

func (m *MockNotesRepository) GetNote(id int) (*notes.Note, error) {
	args := m.Called(id)
	return args.Get(0).(*notes.Note), args.Error(1)
}
func (m *MockNotesRepository) GetUserNotes(userId int) (*[]notes.Note, error) {
	args := m.Called(userId)
	return args.Get(0).(*[]notes.Note), args.Error(1)
}
func (m *MockNotesRepository) CreateNote(userId int, title, body string) error {
	args := m.Called(userId, title, body)
	return args.Error(0)
}
func (m *MockNotesRepository) UpdateNote(id int, title, body string) error {
	args := m.Called(id, title, body)
	return args.Error(0)
}
func (m *MockNotesRepository) DeleteNote(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) GetUserById(id int) (*users.User, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*users.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUsersRepository) GetUserByEmail(email string) (*users.User, error) {
	args := m.Called(email)
	if user, ok := args.Get(0).(*users.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUsersRepository) CreateUser(email, hashedPassword string) error {
	args := m.Called(email, hashedPassword)
	return args.Error(0)
}

func (m *MockUsersRepository) UpdateUser(id int, email, hashedPassword string) error {
	args := m.Called(id, email, hashedPassword)
	return args.Error(0)
}

func (m *MockUsersRepository) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetNote_Success(t *testing.T) {
	//Arrange
	c, rec := newEchoContext(http.MethodGet, "/api/note/1", nil)
	c.SetPath("/api/note/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockNotes := new(MockNotesRepository)
	mockUsers := new(MockUsersRepository)
	mockNotes.On("GetNote", 1).Return(&notes.Note{Id: 1, Title: "test", Body: "body"}, nil)
	s := service.NewService(logs.NewLogger(false), mockNotes, mockUsers)

	//Act
	err := s.GetNote(c)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetNote_NotFound(t *testing.T) {
	//Arrange
	c, rec := newEchoContext(http.MethodGet, "/api/note/-1", nil)
	c.SetPath("/api/note/:id")
	c.SetParamNames("id")
	c.SetParamValues("-1")

	mockNotes := new(MockNotesRepository)
	mockUsers := new(MockUsersRepository)

	mockNotes.On("GetNote", -1).Return((*notes.Note)(nil), errors.New(service.InvalidParams))

	s := service.NewService(logs.NewLogger(false), mockNotes, mockUsers)

	//Act
	err := s.GetNote(c)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockNotes.AssertExpectations(t)
}

func TestGetUserNotes_Success(t *testing.T) {
	//Arrange
	c, rec := newEchoContext(http.MethodGet, "/api/notes", nil)

	claims := jwt.RegisteredClaims{Subject: "user@test.com"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)

	mockNotes := new(MockNotesRepository)
	mockUsers := new(MockUsersRepository)
	mockUsers.On("GetUserByEmail", "user@test.com").Return(&users.User{Id: 1, Email: "user@test.com"}, nil)
	mockNotes.On("GetUserNotes", 1).Return(&[]notes.Note{{Id: 1, Title: "t", Body: "b"}}, nil)

	s := service.NewService(logs.NewLogger(false), mockNotes, mockUsers)

	//Act
	err := s.GetUserNotes(c)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateNote_Success(t *testing.T) {
	// Arrange
	body := []byte(`{"title":"Updated","body":"Changed"}`)
	c, rec := newEchoContext(http.MethodPut, "/note/5", body)
	c.SetPath("/note/:id")
	c.SetParamNames("id")
	c.SetParamValues("5")

	mockNotes := new(MockNotesRepository)
	mockUsers := new(MockUsersRepository)

	mockNotes.On("UpdateNote", 5, "Updated", "Changed").Return(nil)

	s := service.NewService(logs.NewLogger(false), mockNotes, mockUsers)

	// Act
	err := s.UpdateNote(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())

	mockNotes.AssertExpectations(t)
}

func TestUpdateNote_InvalidID(t *testing.T) {
	// Arrange
	body := []byte(`{"title":"T","body":"B"}`)
	c, rec := newEchoContext(http.MethodPut, "/note/abc", body)
	c.SetPath("/note/:id")
	c.SetParamNames("id")
	c.SetParamValues("abc")

	mockNotes := new(MockNotesRepository)
	mockUsers := new(MockUsersRepository)

	s := service.NewService(logs.NewLogger(false), mockNotes, mockUsers)

	// Act
	err := s.UpdateNote(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateNote_RepoError(t *testing.T) {
	// Arrange
	body := []byte(`{"title":"T","body":"B"}`)
	c, rec := newEchoContext(http.MethodPut, "/note/5", body)
	c.SetPath("/note/:id")
	c.SetParamNames("id")
	c.SetParamValues("5")

	mockNotes := new(MockNotesRepository)
	mockUsers := new(MockUsersRepository)

	mockNotes.On("UpdateNote", 5, "T", "B").
		Return(errors.New("db error"))

	s := service.NewService(logs.NewLogger(false), mockNotes, mockUsers)

	// Act
	err := s.UpdateNote(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteNote_Success(t *testing.T) {
	// Arrange
	c, rec := newEchoContext(http.MethodDelete, "/note/10", nil)
	c.SetPath("/note/:id")
	c.SetParamNames("id")
	c.SetParamValues("10")

	mockNotes := new(MockNotesRepository)
	mockUsers := new(MockUsersRepository)

	mockNotes.On("DeleteNote", 10).Return(nil)

	s := service.NewService(logs.NewLogger(false), mockNotes, mockUsers)

	// Act
	err := s.DeleteNote(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())

	mockNotes.AssertExpectations(t)
}

func TestGetNote_Integration(t *testing.T) {
	// Arrange
	e := echo.New()

	mockNotes := new(MockNotesRepository)
	mockUsers := new(MockUsersRepository)

	expectedNote := &notes.Note{Id: 1, Title: "Test title", Body: "Test body"}
	mockNotes.On("GetNote", 1).Return(expectedNote, nil)

	s := service.NewService(logs.NewLogger(false), mockNotes, mockUsers)

	e.GET("/api/note/:id", s.GetNote)

	req := httptest.NewRequest(http.MethodGet, "/api/note/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp service.Response
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	noteMap, ok := resp.Object.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(1), noteMap["id"])
	assert.Equal(t, "Test title", noteMap["title"])
	assert.Equal(t, "Test body", noteMap["body"])

	mockNotes.AssertExpectations(t)
}

func newEchoContext(method, path string, body []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
