package mock

import (
	"time"

	"github.com/google/uuid"
	"github.com/nickyrolly/dealls-test/internal/services/user"
	"golang.org/x/crypto/bcrypt"
)

func CreateMockUser() *user.Entity {
	id := uuid.New()
	lastName := "Doe"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return &user.Entity{
		ID:           id,
		Email:        "john.doe@example.com",
		PasswordHash: string(hashedPassword),
		FirstName:    "John",
		LastName:     &lastName,
		DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:       "male",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
