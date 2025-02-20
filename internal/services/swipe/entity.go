package swipe

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Swipe struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	ProfileID uuid.UUID `gorm:"type:uuid;not null;column:profile_id"`
	Action    string    `gorm:"type:varchar(10);not null;column:action"` // "like" or "pass"
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:created_at"`
}

func (e *Swipe) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

func (e *Swipe) TableName() string {
	return "swipe"
}
