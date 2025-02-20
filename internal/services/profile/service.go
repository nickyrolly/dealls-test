package profile

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nickyrolly/dealls-test/internal/services/user"
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

func (s *Service) GetUser(userID uuid.UUID) (*user.Entity, error) {
	var user user.Entity
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUserProfile(userID uuid.UUID) (*UserProfile, error) {
	var profile UserProfile
	if err := s.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (s *Service) GetUserPhotos(userID uuid.UUID) ([]UserPhoto, error) {
	var photos []UserPhoto
	if err := s.DB.Where("user_id = ?", userID).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (s *Service) GetUserPreference(userID uuid.UUID) (*UserPreference, error) {
	var pref UserPreference
	if err := s.DB.Where("user_id = ?", userID).First(&pref).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pref, nil
}

func (s *Service) UpdateUserProfile(userID uuid.UUID, data map[string]interface{}) error {
	var profile UserProfile
	if err := s.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new profile
			profile = UserProfile{
				UserID: userID,
			}
			if err := s.DB.Create(&profile).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if interests, ok := data["interests"].([]interface{}); ok {
		var interestsStr []string
		for _, interest := range interests {
			if str, ok := interest.(string); ok {
				interestsStr = append(interestsStr, str)
			} else {
				return fmt.Errorf("interests must be a string array") // Handle error
			}

		}
		data["interests"] = strings.Join(interestsStr, ",") // Join string dengan koma
	}

	if err := s.DB.Model(&profile).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (s *Service) AddUserPhoto(photo *UserPhoto) error {
	if photo.IsPrimary {
		// Set all other photos to non-primary
		if err := s.DB.Model(&UserPhoto{}).Where("user_id = ?", photo.UserID).Update("is_primary", false).Error; err != nil {
			return err
		}
	}

	return s.DB.Create(photo).Error
}

func (s *Service) UpdateUserPreference(userID uuid.UUID, data map[string]interface{}) error {
	// return s.DB.Save(pref).Error

	var preference UserPreference
	if err := s.DB.Where("user_id = ?", userID).First(&preference).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new preference
			preference = UserPreference{
				UserID: userID,
			}
			if err := s.DB.Create(&preference).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := s.DB.Model(&preference).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (s *Service) GetDiscovery(query string) ([]UserProfile, error) {
	var profiles []UserProfile
	if err := s.DB.Where("first_name LIKE ? OR last_name LIKE ?", "%"+query+"%", "%"+query+"%").Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

func (s *Service) CreateLike(like *UserLike) (*UserMatch, error) {
	// Check if users exist
	var likerProfile, likedProfile UserProfile
	if err := s.DB.Where("user_id = ?", like.LikerID).First(&likerProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create liker profile
			likerProfile = UserProfile{
				UserID: like.LikerID,
			}
			if err := s.DB.Create(&likerProfile).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	if err := s.DB.Where("user_id = ?", like.LikedID).First(&likedProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create liked profile
			likedProfile = UserProfile{
				UserID: like.LikedID,
			}
			if err := s.DB.Create(&likedProfile).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Check if like already exists
	var existingLike UserLike
	if err := s.DB.Where("liker_id = ? AND liked_id = ?", like.LikerID, like.LikedID).First(&existingLike).Error; err == nil {
		return nil, errors.New("like already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create like
	if err := s.DB.Create(like).Error; err != nil {
		return nil, err
	}

	// Check if mutual like exists
	var mutualLike UserLike
	if err := s.DB.Where("liker_id = ? AND liked_id = ?", like.LikedID, like.LikerID).First(&mutualLike).Error; err == nil {
		// Create match
		newUUID := uuid.New()
		match := &UserMatch{
			ID:        newUUID,
			User1ID:   like.LikerID,
			User2ID:   like.LikedID,
			CreatedAt: time.Now(),
		}
		if err := s.DB.Create(match).Error; err != nil {
			return nil, err
		}
		return match, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return nil, nil
}

func (s *Service) GetMatches(userID uuid.UUID) ([]UserMatch, error) {
	var matches []UserMatch
	if err := s.DB.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (s *Service) WithdrawLike(likerID, likedID uuid.UUID) error {
	return s.DB.Where("liker_id = ? AND liked_id = ?", likerID, likedID).Delete(&UserLike{}).Error
}

func (s *Service) UnmatchUsers(user1ID, user2ID uuid.UUID) error {
	return s.DB.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID).Delete(&UserMatch{}).Error
}

func (s *Service) SearchProfiles(query string) ([]UserProfile, error) {
	var profiles []UserProfile
	if err := s.DB.Where("first_name LIKE ? OR last_name LIKE ?", "%"+query+"%", "%"+query+"%").Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}
