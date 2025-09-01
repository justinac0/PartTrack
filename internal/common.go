package internal

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func ErrorPopup(c echo.Context, status int, msg string) error {
	c.Response().Header().Set("HX-Retarget", "#errors")
	c.Response().Header().Set("HX-Reswap", "outerHTML")
	return c.String(status, fmt.Sprintf("%d: %s", status, msg))
}
