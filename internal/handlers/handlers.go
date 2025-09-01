package handlers

import (
	"PartTrack/internal"
	"PartTrack/internal/db"
	"PartTrack/internal/db/views"
	"PartTrack/internal/templates"
	"net/http"

	"github.com/labstack/echo/v4"
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

func SessionMiddleware(next echo.HandlerFunc, stop echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := views.ValidateSession(c)
		if err != nil {
			return stop(c)
		}

		return next(c)
	}
}

// TODO: re-render on DB changes: https://readmedium.com/creating-a-custom-change-data-capture-cdc-tool-in-golang-5a580ba7ac98
func Setup(e *echo.Echo) {
	db.Init()

	e.GET("/", indexPage)

	userHandler := views.NewUsersHandler()
	componentsHandler := views.NewComponentsHandler()

	auth := e.Group("/auth")
	auth.POST("/signin", userHandler.SignIn)
	auth.POST("/signout", userHandler.SignOut)
	auth.POST("/register", userHandler.Register)
	auth.GET("/who-am-i", userHandler.WhoAmI)

	protected := e.Group("/protected")
	protected.GET("/dashboard", SessionMiddleware(dashboardPage, notAuthorizedPage))
	protected.GET("/components/:id", SessionMiddleware(componentsHandler.SingleComponentView, notAuthorizedPage))
	protected.GET("/components/page/:id", SessionMiddleware(componentsHandler.ComponentsTableView, notAuthorizedPage))
}
