package auth

import (
	"PartTrack/internal/resource/sessions"
	"PartTrack/internal/resource/users"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := users.ValidateSession(c)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}

func Setup(e *echo.Echo, userHandler *users.Handler, sessionHandler *sessions.Handler) {
	e.POST("/signin", userHandler.SignIn)
	e.POST("/signout", userHandler.SignOut)
	e.POST("/register", userHandler.Register)
	e.GET("/who-am-i", userHandler.WhoAmI)
}
