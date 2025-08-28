package main

import (
	"PartTrack/cmd/internal/db"
	"PartTrack/cmd/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "internal/static",
		Browse: true,
	}))
	e.Use(middleware.CORS())

	return e
}

func main() {
	// setup echo + handlers
	e := InitEcho()

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
