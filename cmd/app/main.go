package main

import (
	"PartTrack/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"
)

func InitEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status}\n",
	}))
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(5))))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("static/", "static/")

	return e
}

func main() {
	// setup echo + handlers
	e := InitEcho()
	e.Logger.SetLevel(log.DEBUG)

	// setup handlers
	handlers.Setup(e)

	// ready to listen
	e.Logger.Fatal(e.Start("127.0.0.1:8000"))
}
