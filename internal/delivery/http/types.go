package http

import "time"

// Authentication requests
type signupRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	FirstName   string    `json:"first_name"`
	LastName    *string   `json:"last_name,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Bio         *string   `json:"bio,omitempty"`
	LocationLat *float64  `json:"location_lat,omitempty"`
	LocationLng *float64  `json:"location_lng,omitempty"`
}

// Basic profile requests
type basicProfileRequest struct {
	PhoneNumber *string    `json:"phone_number,omitempty"`
	FirstName   string     `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
	LocationLat *float64   `json:"location_lat,omitempty"`
	LocationLng *float64   `json:"location_lng,omitempty"`
}

// Extended profile requests
type extendedProfileRequest struct {
	Height     *int      `json:"height,omitempty"`
	Weight     *int      `json:"weight,omitempty"`
	Occupation *string   `json:"occupation,omitempty"`
	Education  *string   `json:"education,omitempty"`
	Religion   *string   `json:"religion,omitempty"`
	Ethnicity  *string   `json:"ethnicity,omitempty"`
	Interests  []string  `json:"interests,omitempty"`
	AboutMe    *string   `json:"about_me,omitempty"`
}

// Preference requests
type preferenceRequest struct {
	MinAge    *int    `json:"min_age,omitempty"`
	MaxAge    *int    `json:"max_age,omitempty"`
	MinHeight *int    `json:"min_height,omitempty"`
	MaxHeight *int    `json:"max_height,omitempty"`
	Religion  *string `json:"religion,omitempty"`
	Ethnicity *string `json:"ethnicity,omitempty"`
	Distance  *int    `json:"distance,omitempty"`
}
