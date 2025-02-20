package subscription

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

func NewService(db *gorm.DB, log *logrus.Logger) *Service {
	return &Service{
		DB:  db,
		Log: log,
	}
}

func (s *Service) CreateSubscription(userID uuid.UUID, isSubscribed bool) error {
	// Create a new subscription record
	subscription := Subscription{
		ID:           uuid.New(),
		UserID:       userID,
		IsSubscribed: isSubscribed,
		CreatedAt:    time.Now(),
	}
	return s.DB.Create(&subscription).Error
}
