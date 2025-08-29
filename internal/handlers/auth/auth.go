package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// func verifyPassword(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

// TODO: add role based auth paired with session auth
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func Setup(e *echo.Echo) {
	e.POST("/signin", func(c echo.Context) error {
		user := c.FormValue("username")
		pass := c.FormValue("password")

		passHash, err := hashPassword(pass)
		if err != nil {
			panic(err)
		}

		fmt.Println(user, passHash)

		c.Response().Header().Add("HX-Redirect", "/dashboard")
		return c.NoContent(http.StatusOK)
	})
	e.GET("/signout", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.GET("/register", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
