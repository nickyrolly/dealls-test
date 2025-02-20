package swipe

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

func (s *Service) CreateSwipe(userID, profileID uuid.UUID, action string) error {
	// Check if the user has already swiped on this profile today
	var count int64
	today := time.Now().Format("2006-01-02")
	s.DB.Model(&Swipe{}).Where("user_id = ? AND profile_id = ? AND DATE(created_at) = ?", userID, profileID, today).Count(&count)

	if count > 0 {
		return errors.New("you have already swiped on this profile today")
	}

	// Check total swipes today
	var totalSwipes int64
	s.DB.Model(&Swipe{}).Where("user_id = ? AND DATE(created_at) = ?", userID, today).Count(&totalSwipes)

	if totalSwipes >= 10 {
		return errors.New("you have reached the daily limit of swipes")
	}

	// Create a new swipe record
	swipe := Swipe{
		ID:        uuid.New(),
		UserID:    userID,
		ProfileID: profileID,
		Action:    action,
		CreatedAt: time.Now(),
	}
	return s.DB.Create(&swipe).Error
}
