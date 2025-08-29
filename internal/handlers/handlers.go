package handlers

import (
	"PartTrack/internal/db"
	"PartTrack/internal/db/models"
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
	db.Init()

	userStore := models.UserHandler{Store: db.GetStore()}

	auth.Setup(e)

	e.GET("/", indexPage)
	e.GET("/user", auth.Middleware(userStore.GetUsers))
	e.GET("/dashboard", func(c echo.Context) error {
		return c.String(http.StatusOK, "dashboard")
	})
}
