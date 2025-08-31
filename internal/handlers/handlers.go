package handlers

import (
	"PartTrack/internal"
	"PartTrack/internal/db"
	"PartTrack/internal/handlers/auth"
	"PartTrack/internal/resource/components"
	"PartTrack/internal/resource/users"
	"PartTrack/internal/templates"
	"net/http"

	"github.com/labstack/echo/v4"
)

func indexPage(c echo.Context) error {
	err := users.ValidateSession(c)
	if err == nil {
		c.Response().Header().Set("HX-Redirect", "/protected/dashboard")
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

func Setup(e *echo.Echo) {
	db.Init()

	auth.Setup(e)

	e.GET("/", indexPage)
	g := e.Group("/protected")
	g.GET("/dashboard", auth.Middleware(dashboardPage, notAuthorizedPage))

	componentsHandler := components.NewHandler()
	g.GET("/components/:id", auth.Middleware(componentsHandler.ViewOne, notAuthorizedPage))
	g.GET("/components/page/:id", auth.Middleware(componentsHandler.ViewComponents, notAuthorizedPage))
	g.POST("/components/page/:id", auth.Middleware(componentsHandler.ViewComponents, notAuthorizedPage))
}
