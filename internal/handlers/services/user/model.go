package user

import (
	"time"
)

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
