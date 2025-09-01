package auth

import (
	"PartTrack/internal/db/views"

	"github.com/labstack/echo/v4"
)

func SessionMiddleware(next echo.HandlerFunc, stop echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := views.ValidateSession(c)
		if err != nil {
			return stop(c)
		}

		return next(c)
	}
}
