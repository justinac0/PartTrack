package handlers

import (
	"PartTrack/cmd/internal/templates"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// move to helper
func render(c echo.Context, status int, t templ.Component) error {
	c.Response().Writer.WriteHeader(status)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(c.Request().Context(), c.Response().Writer)
}

func indexPage(c echo.Context) error {
	return render(c, http.StatusOK, templates.IndexPage())
}

func Setup(e *echo.Echo) {
	e.GET("/", indexPage)
}
