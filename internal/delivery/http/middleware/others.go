package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/nickyrolly/dealls-test/internal/services/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func credentialCheck(r *http.Request, email string, password string, redisPool *redis.Pool, db *gorm.DB) (UserSession, error) {
	// First, try to get from Redis
	conn := redisPool.Get()
	defer conn.Close()

	// Redis key format: "user:credentials:email"
	redisKey := fmt.Sprintf("user:credentials:%s", email)

	// Try to get user data from Redis
	userJSON, redisErr := redis.String(conn.Do("GET", redisKey))
	if redisErr == nil {
		// Found in Redis, unmarshal and verify password
		var cachedUser UserSession
		if unmarshalErr := json.Unmarshal([]byte(userJSON), &cachedUser); unmarshalErr == nil {
			// Verify password
			if verifyErr := bcrypt.CompareHashAndPassword([]byte(cachedUser.PasswordHash), []byte(password)); verifyErr == nil {
				// Update session token expiration
				cachedUser.SessionExpirationTime = time.Now().Add(360 * time.Minute)
				return cachedUser, nil
			}
		}
	}

	// Not found in Redis or password incorrect, check database
	var dbUser user.Entity
	result := db.Where("email = ?", email).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return UserSession{}, fmt.Errorf("user not found")
		}
		return UserSession{}, result.Error
	}

	// Verify password
	if verifyErr := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(password)); verifyErr != nil {
		return UserSession{}, fmt.Errorf("invalid password")
	}

	// Update last active timestamp
	now := time.Now()
	if updateErr := db.Model(&dbUser).Update("last_active_at", now).Error; updateErr != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to update last_active_at: %v\n", updateErr)
	}

	// Convert to UserSession
	userSession := UserSession{
		ID:                    dbUser.ID.String(),
		Email:                 dbUser.Email,
		PasswordHash:          dbUser.PasswordHash,
		PhoneNumber:           dbUser.PhoneNumber,
		FirstName:             dbUser.FirstName,
		LastName:              dbUser.LastName,
		DateOfBirth:           dbUser.DateOfBirth,
		Gender:                dbUser.Gender,
		LocationLat:           dbUser.LocationLat,
		LocationLng:           dbUser.LocationLng,
		LastActiveAt:          dbUser.LastActiveAt,
		SessionExpirationTime: now.Add(360 * time.Minute),
	}

	// Cache in Redis for future requests (expire in 15 minutes)
	cachedJSON, marshalErr := json.Marshal(userSession)
	if marshalErr == nil {
		_, redisErr = conn.Do("SETEX", redisKey, 900, cachedJSON) // 900 seconds = 15 minutes
		if redisErr != nil {
			// Log Redis error but don't fail the request
			fmt.Printf("Redis cache error: %v\n", redisErr)
		}
	}

	return userSession, nil
}

func profileCheck(r *http.Request, xUserID string) (user UserSession, err error) {

	// participantID, err := strconv.ParseInt(xUserID, 10, 64)
	// if err != nil {
	// 	return user, err
	// }

	// user, err = authRepo.ValidateProfile(r.Context(), participantID)
	// if err != nil {
	// 	log.Println("[credential check] error", err)
	// }
	return user, err
}
