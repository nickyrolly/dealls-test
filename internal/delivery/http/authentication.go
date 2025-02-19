package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
	gorilla_context "github.com/gorilla/context"
	"github.com/nickyrolly/dealls-test/common"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/middleware"
	"github.com/nickyrolly/dealls-test/internal/services/user"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationController struct {
	log       *logrus.Logger
	redisPool *redis.Pool
	db        *gorm.DB
}

func NewAuthenticationController(log *logrus.Logger, redisPool *redis.Pool, db *gorm.DB) *AuthenticationController {
	return &AuthenticationController{
		log:       log,
		redisPool: redisPool,
		db:        db,
	}
}

func (c *AuthenticationController) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.WithError(err).Error("Failed to decode signup request")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success":       false,
			"error_message": "Invalid request format",
		})
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.Gender == "" {
		c.log.Error("Missing required fields in signup request")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success":       false,
			"error_message": "Missing required fields",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.log.WithError(err).Error("Failed to hash password")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success":       false,
			"error_message": "Internal server error",
		})
		return
	}

	// Create user entity
	newUser := user.Entity{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		PhoneNumber:  req.PhoneNumber,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		DateOfBirth:  req.DateOfBirth,
		Gender:       req.Gender,
		LocationLat:  req.LocationLat,
		LocationLng:  req.LocationLng,
	}

	// Check if email exists
	var existingUser user.Entity
	if err := c.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.log.Error("Email already exists")
		common.CustomResponseAPI(w, r, http.StatusConflict, map[string]interface{}{
			"success":       false,
			"error_message": "Email already exists",
		})
		return
	}

	// Check if phone number exists (if provided)
	if req.PhoneNumber != nil {
		if err := c.db.Where("phone_number = ?", req.PhoneNumber).First(&existingUser).Error; err == nil {
			c.log.Error("Phone number already exists")
			common.CustomResponseAPI(w, r, http.StatusConflict, map[string]interface{}{
				"success":       false,
				"error_message": "Phone number already exists",
			})
			return
		}
	}

	// Create user
	result := c.db.Create(&newUser)
	if result.Error != nil {
		c.log.WithError(result.Error).Error("Failed to create user")

		// Check for unique constraint violations
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
			common.CustomResponseAPI(w, r, http.StatusConflict, map[string]interface{}{
				"success":       false,
				"error_message": "Email already exists",
			})
			return
		} else if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.phone_number") {
			common.CustomResponseAPI(w, r, http.StatusConflict, map[string]interface{}{
				"success":       false,
				"error_message": "Phone number already exists",
			})
			return
		}

		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success":       false,
			"error_message": "Failed to create user",
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User created successfully",
	})
}

func (c *AuthenticationController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// The actual authentication is handled by the middleware
	// This endpoint just returns success since if we got here, authentication was successful

	user, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success":       false,
			"error_message": "User session not found",
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data": map[string]interface{}{
			"token":            user.SessionToken,
			"token_expires_at": user.SessionExpirationTime,
		},
	})
}

func (c *AuthenticationController) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user session from context (set by authentication middleware)
	userSession, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success":       false,
			"error_message": "User session not found",
		})
		return
	}

	// Get fresh user data from database
	var user user.Entity
	if err := c.db.First(&user, "id = ?", userSession.ID).Error; err != nil {
		c.log.WithError(err).Error("Failed to get user from database")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success":       false,
			"error_message": "Failed to get user profile",
		})
		return
	}

	// Don't send password hash in response
	user.PasswordHash = ""

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

func (c *AuthenticationController) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user session from context
	userSession, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success":       false,
			"error_message": "User session not found",
		})
		return
	}

	// Get current user from database
	var user user.Entity
	if err := c.db.First(&user, "id = ?", userSession.ID).Error; err != nil {
		c.log.WithError(err).Error("Failed to get user from database")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success":       false,
			"error_message": "Failed to update profile",
		})
		return
	}

	// Decode update request
	var req basicProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.WithError(err).Error("Failed to decode update request")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success":       false,
			"error_message": "Invalid request format",
		})
		return
	}

	// Update only provided fields
	if req.PhoneNumber != nil {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.DateOfBirth != nil {
		user.DateOfBirth = *req.DateOfBirth
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.LocationLat != nil {
		user.LocationLat = req.LocationLat
	}
	if req.LocationLng != nil {
		user.LocationLng = req.LocationLng
	}

	// Update user in database
	if err := c.db.Save(&user).Error; err != nil {
		c.log.WithError(err).Error("Failed to update user")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success":       false,
			"error_message": "Failed to update profile",
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
		"data":    user,
	})
}
