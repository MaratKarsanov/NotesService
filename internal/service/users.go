package service

import (
	"NotesService/cmd/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"email"`
	jwt.RegisteredClaims
}

// localhost:8000/api/login
func (s *Service) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	usersRepository := s.usersRepository
	user, err := usersRepository.GetUserByEmail(email)

	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidCredentials))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidCredentials))
	}

	token, err := GenerateJWT(email)
	if err != nil {
		return err
	}

	s.logger.Infof("User %s authorized successfully", email)
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// localhost:8000/api/login
func (s *Service) Register(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	usersRepository := s.usersRepository
	user, err := usersRepository.GetUserByEmail(email)

	if user != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(UserAlreadyExists))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)

	err = usersRepository.CreateUser(email, string(hashedPassword))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	s.logger.Infof("User %s registered successfully", email)
	return c.JSON(http.StatusOK, "OK")
}

func GenerateJWT(username string) (string, error) {
	appConf, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	jwtKey := []byte(appConf.App.JWTKey)

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
