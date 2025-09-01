package main

import (
	"PartTrack/internal"
	"PartTrack/internal/auth"
	"PartTrack/internal/db"
	"PartTrack/internal/db/views"
	"PartTrack/internal/templates"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"
)

func indexPage(c echo.Context) error {
	err := views.ValidateSession(c)
	if err == nil {
		c.Response().Header().Add("HX-Redirect", "/protected/dashboard")
		return internal.RenderTempl(c, http.StatusOK, templates.DashboardPage())
	}

	return internal.RenderTempl(c, http.StatusOK, templates.IndexPage())
}

func dashboardPage(c echo.Context) error {
	return internal.RenderTempl(c, http.StatusOK, templates.DashboardPage())
}

func notAuthorizedPage(c echo.Context) error {
	return c.String(http.StatusUnauthorized, "you are not authorized to view this content")
}

// TODO: re-render on DB changes: https://readmedium.com/creating-a-custom-change-data-capture-cdc-tool-in-golang-5a580ba7ac98
func setupHandlers(e *echo.Echo) {
	db.Init()

	e.GET("/", indexPage)

	userHandler := views.NewUsersHandler()
	componentsHandler := views.NewComponentsHandler()

	a := e.Group("/auth")
	a.POST("/signin", userHandler.SignIn)
	a.POST("/signout", userHandler.SignOut)
	a.POST("/register", userHandler.Register)
	a.GET("/who-am-i", userHandler.WhoAmI)

	p := e.Group("/protected")
	p.GET("/dashboard", auth.SessionMiddleware(dashboardPage, notAuthorizedPage))
	p.GET("/components/:id", auth.SessionMiddleware(componentsHandler.SingleComponentView, notAuthorizedPage))
	p.GET("/components/page/:id", auth.SessionMiddleware(componentsHandler.ComponentsTableView, notAuthorizedPage))
}

func initEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status}\n",
	}))
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(5))))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("static/", "static")

	return e
}

func main() {
	// setup echo
	e := initEcho()
	e.Logger.SetLevel(log.DEBUG)

	// setup handlers
	setupHandlers(e)

	// ready to listen
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
}
