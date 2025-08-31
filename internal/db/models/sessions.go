package models

import (
	"time"
)

type Session struct {
	SessionId string     `json:"session_id"`
	UserId    uint64     `json:"user_id"`
	ExpiresAt *time.Time `json:"expiry"`
	CreatedAt *time.Time `json:"created"`
}
