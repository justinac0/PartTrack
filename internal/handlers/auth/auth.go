package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Setup(e *echo.Echo) {
	e.POST("/signin", func(c echo.Context) error {
		user := c.FormValue("username")
		pass := c.FormValue("password")

		passBytes := []byte(pass)

		passHash, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		fmt.Printf("username: %s, pass_hash: %s\n", user, string(passHash))

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
