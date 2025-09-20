package main

import (
	"NotesService/cmd/config"
	"NotesService/pkg/logs"

	"github.com/labstack/echo/v4"
)

func main() {
	logger := logs.NewLogger(false)

	_, err := PostgresConnection()
	if err != nil {
		logger.Fatal(err)
	}

	router := echo.New()
	appConf, err := config.GetConfig()
	if err != nil {
		logger.Fatal(err)
	}
	port := appConf.App.Port

	router.Logger.Fatal(router.Start(":" + port))
}
