package sessions

import (
	"time"
)

type Session struct {
	SessionId string     `json:"session_id"`
	UserId    uint64     `json:"user_id"`
	Expiry    *time.Time `json:"expiry"`
	Created   *time.Time `json:"created"`
}
