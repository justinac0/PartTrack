package handlers

import (
	"PartTrack/internal/db"
	"PartTrack/internal/handlers/auth"
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
	err := users.ValidateSession(c)
	if err == nil {
		c.Response().Header().Set("HX-Redirect", "/protected/dashboard")
		return render(c, http.StatusOK, templates.DashboardPage())
	}

	return render(c, http.StatusOK, templates.IndexPage())
}

func dashboardPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.DashboardPage())
}

func componentsPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.ComponentsPage())
}

func Setup(e *echo.Echo) {
	db.Init()

	auth.Setup(e)

	e.GET("/", indexPage)
	g := e.Group("/protected")
	g.GET("/dashboard", auth.Middleware(dashboardPage, func(c echo.Context) error {
		return c.String(http.StatusUnauthorized, "you are not authorized to view this content")
	}))
	g.GET("/components", auth.Middleware(componentsPage, func(c echo.Context) error {
		return c.String(http.StatusUnauthorized, "you are not authorized to view this content")
	}))
}
