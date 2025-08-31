package internal

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderTempl(c echo.Context, status int, t templ.Component) error {
	c.Response().Writer.WriteHeader(status)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(c.Request().Context(), c.Response().Writer)
}
