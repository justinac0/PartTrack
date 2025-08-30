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
			return c.HTML(http.StatusUnauthorized, fmt.Sprintf("<h1>Unauthorized Access: %v</h1><a href='/'>goto signin page</a>", err))
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
