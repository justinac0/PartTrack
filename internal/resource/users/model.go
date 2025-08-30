package users

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
	Id           uint64     `json:"id"`
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password_hash"`
	Role         UserRole   `json:"role"`
	Created      *time.Time `json:"created_at"`
	Deleted      *time.Time `json:"deleted_at"`
}
