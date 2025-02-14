package middleware

import "time"

type UserSession struct {
	UserID                int64     `json:"user_id"`
	Email                 string    `json:"email"`
	Password              string    `json:"password,omitempty"`
	PhoneNumber           string    `json:"phone_number"`
	SessionToken          string    `json:"session_token"`
	SessionExpirationTime time.Time `json:"session_expiration_time"`
}
