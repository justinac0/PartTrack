package handlers

import (
	"PartTrack/internal"
	"PartTrack/internal/db"
	"PartTrack/internal/handlers/auth"
	"PartTrack/internal/resources/components"
	"PartTrack/internal/resources/users"
	"PartTrack/internal/templates"
	"net/http"

	"github.com/labstack/echo/v4"
)

func indexPage(c echo.Context) error {
	err := users.ValidateSession(c)
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
func Setup(e *echo.Echo) {
	db.Init()

	auth.Setup(e)

	e.GET("/", indexPage)
	g := e.Group("/protected")
	g.GET("/dashboard", auth.Middleware(dashboardPage, notAuthorizedPage))

	componentsHandler := components.NewHandler()
	g.GET("/components/:id", auth.Middleware(componentsHandler.ViewOne, notAuthorizedPage))
	g.GET("/components/page/:id", auth.Middleware(componentsHandler.ViewComponents, notAuthorizedPage))
}
