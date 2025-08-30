package auth

import (
	"PartTrack/internal/resource/sessions"
	"PartTrack/internal/resource/users"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := CheckSession(c)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return next(c)
	}
}

func CheckSession(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err == nil {
		return errors.New("session does not exist")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sessionStore := sessions.NewStore()
	session, err := sessionStore.GetBySessionID(ctx, cookie.Value)
	if err != nil {
		return err
	}

	log.Println(session)

	return nil
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
