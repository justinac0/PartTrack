package auth

import (
	"PartTrack/internal/resource/sessions"
	"PartTrack/internal/resource/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := sessions.ValidateSession(c)
		if err != nil {
			return c.String(http.StatusUnauthorized, "unauthorized access")
		}

		return next(c)
	}
}

func Setup(e *echo.Echo, userHandler *users.Handler, sessionHandler *sessions.Handler) {
	e.POST("/signin", userHandler.SignIn)
	e.GET("/signout", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.GET("/register", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
