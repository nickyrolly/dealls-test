package profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	gorilla_context "github.com/gorilla/context"
	"github.com/nickyrolly/dealls-test/common"
	"github.com/nickyrolly/dealls-test/internal/services/profile"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	log     *logrus.Logger
	service *profile.Service
}

func NewController(log *logrus.Logger, service *profile.Service) *Controller {
	return &Controller{
		log:     log,
		service: service,
	}
}

func (c *Controller) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
	fmt.Println("--- HandleGetProfile")
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from context
	userIDStr, ok := gorilla_context.Get(r, "user").(string)
	fmt.Println("userIDStr : ", userIDStr)
	if !ok {
		c.log.Error("User session not found in context")

		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		fmt.Println("context userID : ", err)
		// c.log.WithError(err).Error("Failed to parse user ID")
		// http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := c.service.GetUser(userID)
	if err != nil {
		fmt.Println("get user : ", err)
		c.log.WithError(err).Error("Failed to get user")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to get user",
		})
		return
	}

	// Get profile
	userProfile, err := c.service.GetUserProfile(userID)
	if err != nil {
		fmt.Println("get user profile : ", err)
		c.log.WithError(err).Error("Failed to get user profile")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to get user profile",
		})
		return
	}

	// Get photos
	photos, err := c.service.GetUserPhotos(userID)
	if err != nil {
		c.log.WithError(err).Error("Failed to get user photos")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to get user photos",
		})
		return
	}

	// Get preferences
	preferences, err := c.service.GetUserPreference(userID)
	if err != nil {
		c.log.WithError(err).Error("Failed to get user preferences")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to get user preferences",
		})
		return
	}

	userData := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"phone":      user.PhoneNumber,
		"date_of_birth": map[string]interface{}{
			"day":   user.DateOfBirth.Day(),
			"month": user.DateOfBirth.Month(),
			"year":  user.DateOfBirth.Year(),
		},
		"gender": user.Gender,
	}

	response := map[string]interface{}{
		"success":     true,
		"user":        userData,
		"profile":     userProfile,
		"preferences": preferences,
		"photos":      photos,
	}

	common.CustomResponseAPI(w, r, http.StatusOK, response)
}

func (c *Controller) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from context
	userIDStr, ok := gorilla_context.Get(r, "user").(string)
	if !ok {
		c.log.Error("User session not found in context")

		return
	}

	fmt.Println("userIDStr : ", userIDStr)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse user ID")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Decode request body
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		c.log.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("updates : %v\n", updates)

	// Update profile
	if err := c.service.UpdateUserProfile(userID, updates); err != nil {
		c.log.WithError(err).Error("Failed to update user profile")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get updated profile
	profile, err := c.service.GetUserProfile(userID)
	if err != nil {
		c.log.WithError(err).Error("Failed to get updated user profile")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func (c *Controller) HandleUpdatePreferences(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
	fmt.Println("--- HandleUpdatePreferences")
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from context
	userIDStr, ok := gorilla_context.Get(r, "user").(string)
	if !ok {
		c.log.Error("User session not found in context")

		return
	}

	fmt.Println("userIDStr : ", userIDStr)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse user ID")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Decode request body
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		c.log.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("updates : %v\n", updates)

	// Update preferences
	if err := c.service.UpdateUserPreference(userID, updates); err != nil {
		c.log.WithError(err).Error("Failed to update user preferences")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to update preferences",
		})
		return
	}

	// Get updated preferences
	preference, err := c.service.GetUserPreference(userID)
	if err != nil {
		c.log.WithError(err).Error("Failed to get updated user preferences")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to get updated preferences",
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success":     true,
		"preferences": preference,
	})
}

func (c *Controller) HandleGetMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "Not implemented"})
}

func (c *Controller) HandleLikeProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "Not implemented"})
}

func (c *Controller) HandlePassProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "Not implemented"})
}

func (c *Controller) HandleGetPotentialMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "Not implemented"})
}
