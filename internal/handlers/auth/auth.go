package auth

import (
	"PartTrack/internal/resources/users"

	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc, stop echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := users.ValidateSession(c)
		if err != nil {
			return stop(c)
		}

		return next(c)
	}
}

func Setup(e *echo.Echo) {
	userHandler := users.NewHandler()
	g := e.Group("/auth")
	g.POST("/signin", userHandler.SignIn)
	g.POST("/signout", userHandler.SignOut)
	g.POST("/register", userHandler.Register)
	g.GET("/who-am-i", userHandler.WhoAmI)
}
