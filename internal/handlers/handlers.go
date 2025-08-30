package handlers

import (
	"PartTrack/internal/db"
	"PartTrack/internal/handlers/auth"
	"PartTrack/internal/resource/sessions"
	"PartTrack/internal/resource/users"
	"PartTrack/internal/templates"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, status int, t templ.Component) error {
	c.Response().Writer.WriteHeader(status)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(c.Request().Context(), c.Response().Writer)
}

func indexPage(c echo.Context) error {
	err := auth.CheckSession(c)
	if err == nil {
		c.Request().Header.Add("HX-Redirect", "/dashboard")
		return render(c, http.StatusOK, templates.DashboardPage())
	}

	return render(c, http.StatusOK, templates.IndexPage())
}

func dashboardPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.DashboardPage())
}

func Setup(e *echo.Echo) {
	db.Init()

	userHandler := users.NewHandler()
	sessionHandler := sessions.NewHandler()

	auth.Setup(e, userHandler, sessionHandler)

	e.GET("/", indexPage)
	e.GET("/dashboard", auth.Middleware(dashboardPage))
}
