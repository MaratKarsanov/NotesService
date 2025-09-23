package main

import (
	"NotesService/cmd/config"
	"NotesService/internal/service"
	"NotesService/pkg/logs"

	"github.com/golang-jwt/jwt/v5"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	logger := logs.NewLogger(false)

	db, err := PostgresConnection()
	if err != nil {
		logger.Fatal(err)
	}

	router := echo.New()
	appConf, err := config.GetConfig()
	if err != nil {
		logger.Fatal(err)
	}

	svc := service.NewService(db, logger)

	router.POST("/login", svc.Login)
	router.POST("/register", svc.Register)

	api := router.Group("api")
	jwtKey := []byte(appConf.App.JWTKey)
	api.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  jwtKey,
		TokenLookup: "header:Authorization",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(service.Claims)
		},
	}))

	api.GET("/notes", svc.GetUserNotes)
	api.GET("/note/:id", svc.GetNote)
	api.POST("/note", svc.CreateNote)
	api.PUT("/note/:id", svc.UpdateNote)
	api.DELETE("/note/:id", svc.DeleteNote)

	port := appConf.App.Port

	logger.Info("Starting application...")
	router.Logger.Fatal(router.Start(":" + port))
}
