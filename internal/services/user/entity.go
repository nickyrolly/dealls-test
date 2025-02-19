package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Entity struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;column:id"`
	Email        string     `gorm:"type:varchar(255);unique;not null;column:email"`
	PhoneNumber  *string    `gorm:"type:varchar(20);unique;column:phone_number"`
	PasswordHash string     `gorm:"type:varchar(255);not null;column:password_hash"`
	FirstName    string     `gorm:"type:varchar(50);not null;column:first_name"`
	LastName     *string    `gorm:"type:varchar(50);column:last_name"`
	DateOfBirth  time.Time  `gorm:"type:date;not null;column:date_of_birth"`
	Gender       string     `gorm:"type:varchar(20);not null;column:gender"`
	LocationLat  *float64   `gorm:"type:decimal(10,8);column:location_lat"`
	LocationLng  *float64   `gorm:"type:decimal(11,8);column:location_lng"`
	LastActiveAt *time.Time `gorm:"type:timestamp;column:last_active_at"`
	CreatedAt    time.Time  `gorm:"type:timestamp;not null;default:current_timestamp;column:created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamp;not null;default:current_timestamp;column:updated_at"`
	DeletedAt    *time.Time `gorm:"type:timestamp;column:deleted_at"`
}

func (e *Entity) TableName() string {
	return "users"
}

// BeforeCreate will set a UUID rather than numeric ID.
func (e *Entity) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

// VerifyPassword checks if the provided password matches the hashed password
func (e *Entity) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(e.PasswordHash), []byte(password))
	return err == nil
}
