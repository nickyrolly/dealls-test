package authentication

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nickyrolly/dealls-test/internal/services/user"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewService(db *gorm.DB, log *logrus.Logger) *Service {
	return &Service{
		db:  db,
		log: log,
	}
}

func (s *Service) SignUp(email, password, firstName string, lastName *string, dateOfBirth time.Time, gender string) (*user.Entity, string, error) {
	// Check if user already exists
	var existingUser user.Entity
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, "", errors.New("user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	newUser := &user.Entity{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirth:  dateOfBirth,
		Gender:       gender,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(newUser).Error; err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": newUser.ID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return nil, "", err
	}

	return newUser, tokenString, nil
}

func (s *Service) Login(email, password string) (*user.Entity, string, error) {
	var user user.Entity
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return nil, "", err
	}

	return &user, tokenString, nil
}
