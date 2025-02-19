package profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	gorilla_context "github.com/gorilla/context"
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

	// Get profile
	userProfile, err := c.service.GetUserProfile(userID)
	if err != nil {
		fmt.Println("get user profile : ", err)
		// c.log.WithError(err).Error("Failed to get user profile")
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if userProfile == nil {
		fmt.Println("userProfile is nil")
		// http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(userProfile)
}

func (c *Controller) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from context
	userIDStr, ok := gorilla_context.Get(r, "user").(string)
	fmt.Println("userIDStr : ", userIDStr)
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
