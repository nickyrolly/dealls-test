package subscription

import (
	"errors"
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

func (s *Service) CreateSubscription(userID uuid.UUID) error {
	var subscription Subscription
	err := s.DB.Where("user_id = ?", userID).First(&subscription).Error
	if err == nil {
		// Subscription already exists, you can return an error or update the existing record
		return errors.New("subscription already exists for this user")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other error occurred
		return err
	}

	// Create a new subscription record since it doesn't exist
	data := Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	return s.DB.Create(&data).Error
}

func (s *Service) CheckSubscription(userID uuid.UUID) (bool, error) {
	var subscription Subscription
	err := s.DB.Where("user_id = ?", userID).First(&subscription).Error
	if err == nil {
		// Subscription exists
		return true, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Subscription does not exist
		return false, nil
	} else {
		// Some other error occurred
		return false, err
	}
}
