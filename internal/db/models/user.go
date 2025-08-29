package models

import (
	"PartTrack/internal/db"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Store *db.Store
}

type UserRole string

const (
	GUEST    UserRole = "guest"
	CUSTOMER UserRole = "customer"
	EMPLOYEE UserRole = "employee"
	ADMIN    UserRole = "admin"
)

type User struct {
	Id           int        `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password_hash"`
	Role         UserRole   `json:"role"`
	Created      *time.Time `json:"created"`
	Deleted      *time.Time `json:"deleted"`
}

func getUser(db *sql.DB, id int64) (User, error) {
	user := User{}
	row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role, &user.Created, &user.Deleted)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)

	return user, err
}

func getUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, username, password_hash, role, created, deleted FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role, &user.Created, &user.Deleted)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (h *UserHandler) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user, err := getUser(h.Store.DB, id)
	if err != nil {
		panic(err)
	}

	log.Println(user)

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := getUsers(h.Store.DB)
	if err != nil {
		panic(err)
	}

	log.Println(users)

	return c.NoContent(http.StatusOK)
}
