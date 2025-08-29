package handlers

import (
	"PartTrack/internal/handlers/auth"
	"PartTrack/internal/templates"
	"fmt"
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

	cookies := c.Cookies()
	fmt.Println(cookies)

	return render(c, http.StatusOK, templates.IndexPage())
}

func Setup(e *echo.Echo) {
	auth.Setup(e)

	e.GET("/", indexPage)
	e.GET("/dashboard", func(c echo.Context) error {
		return c.String(http.StatusOK, "dashboard")
	})
}
