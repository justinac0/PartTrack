// TODO testing https://blog.jetbrains.com/go/2022/11/22/comprehensive-guide-to-testing-in-go/
// TODO login and route authentication RBAC
package main

import (
	"PartTrack/internal/db"
	"PartTrack/internal/handlers"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// setup db
	err := db.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.CloseDB()

	handlers.Setup(e)

	log.Println("db connection established...")

	// ready to listen
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}
