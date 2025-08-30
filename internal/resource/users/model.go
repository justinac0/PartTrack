package users

import (
	"time"
)

type UserRole string

const (
	RoleGuest    UserRole = "guest"
	RoleCustomer UserRole = "customer"
	RoleEmployee UserRole = "employee"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	Id           uint64     `json:"id"`
	Email        string     `json:"email"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password_hash"`
	Role         UserRole   `json:"role"`
	CreatedAt    *time.Time `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
