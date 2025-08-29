package main

import (
	"PartTrack/internal/db"
	"PartTrack/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func InitEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status}\n",
	}))
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("static/", "static/")

	return e
}

func main() {
	// setup echo + handlers
	e := InitEcho()
	e.Logger.SetLevel(log.DEBUG)

	// setup db
	err := db.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.CloseDB()

	handlers.Setup(e)

	// ready to listen
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}
