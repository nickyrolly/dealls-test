package subscription

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;column:id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	IsSubscribed bool      `gorm:"type:boolean;not null;default:false;column:is_subscribed"`
	CreatedAt    time.Time `gorm:"type:timestamp;not null;default:current_timestamp;column:created_at"`
}

func (e *Subscription) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}

func (e *Subscription) TableName() string {
	return "subscription"
}
