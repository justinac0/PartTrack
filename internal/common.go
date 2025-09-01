package internal

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func ErrorPopup(c echo.Context, status int, msg string) error {
	c.Response().Header().Add("HX-Retarget", "#errors")
	c.Response().Header().Add("HX-Reswap", "outerHTML")
	html := fmt.Sprintf("<div id='errors'>%d: %s</div>", status, msg)
	return c.HTML(status, html)
}
