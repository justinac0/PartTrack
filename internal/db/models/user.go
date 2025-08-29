package models

type UserRole int

const (
	GUEST UserRole = iota
	CUSTOMER
	EMPLOYEE
	ADMIN
)

type User struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}
