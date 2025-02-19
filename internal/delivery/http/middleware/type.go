package middleware

import "time"

type UserSession struct {
	ID                    string     `json:"id"`
	Email                 string     `json:"email"`
	PasswordHash          string     `json:"password_hash,omitempty"`
	PhoneNumber           *string    `json:"phone_number,omitempty"`
	FirstName             string     `json:"first_name"`
	LastName              *string    `json:"last_name,omitempty"`
	DateOfBirth           time.Time  `json:"date_of_birth"`
	Gender                string     `json:"gender"`
	LocationLat           *float64   `json:"location_lat,omitempty"`
	LocationLng           *float64   `json:"location_lng,omitempty"`
	LastActiveAt          *time.Time `json:"last_active_at,omitempty"`
	SessionToken          string     `json:"session_token"`
	SessionExpirationTime time.Time  `json:"session_expiration_time"`
}
