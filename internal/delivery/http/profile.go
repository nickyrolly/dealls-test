package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	gorilla_context "github.com/gorilla/context"
	"github.com/nickyrolly/dealls-test/common"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/middleware"
	"github.com/nickyrolly/dealls-test/internal/services/profile"
	"github.com/sirupsen/logrus"
)

type ProfileController struct {
	log     *logrus.Logger
	service *profile.Service
}

func NewProfileController(log *logrus.Logger, service *profile.Service) *ProfileController {
	return &ProfileController{
		log:     log,
		service: service,
	}
}

func (c *ProfileController) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userSession, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "User session not found",
		})
		return
	}

	userID, err := uuid.Parse(userSession.ID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse user ID")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	// Get profile data
	userProfile, err := c.service.GetUserProfile(userID)
	if err != nil {
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

	response := map[string]interface{}{
		"success": true,
		"profile": userProfile,
		"photos":  photos,
	}
	if preferences != nil {
		response["preferences"] = preferences
	}

	common.CustomResponseAPI(w, r, http.StatusOK, response)
}

func (c *ProfileController) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	userSession, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "User session not found",
		})
		return
	}

	userID, err := uuid.Parse(userSession.ID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse user ID")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		c.log.WithError(err).Error("Failed to decode request body")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Validate updates
	allowedFields := map[string]bool{
		"height":     true,
		"weight":     true,
		"occupation": true,
		"education":  true,
		"religion":   true,
		"ethnicity":  true,
		"interests":  true,
		"about_me":   true,
	}

	for field := range updates {
		if !allowedFields[field] {
			c.log.WithField("field", field).Error("Invalid field in update request")
			common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   "Invalid field: " + field,
			})
			return
		}
	}

	// Update profile
	if err := c.service.UpdateUserProfile(userID, updates); err != nil {
		c.log.WithError(err).Error("Failed to update user profile")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to update profile",
		})
		return
	}

	// Get updated profile
	profile, err := c.service.GetUserProfile(userID)
	if err != nil {
		c.log.WithError(err).Error("Failed to get updated user profile")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to get updated profile",
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
		"profile": profile,
	})
}

func (c *ProfileController) HandleUpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userSession, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "User session not found",
		})
		return
	}

	userID, err := uuid.Parse(userSession.ID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse user ID")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	var pref profile.UserPreference
	if err := json.NewDecoder(r.Body).Decode(&pref); err != nil {
		c.log.WithError(err).Error("Failed to decode request body")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Set user ID
	pref.UserID = userID

	// Update preferences
	if err := c.service.UpdateUserPreference(&pref); err != nil {
		c.log.WithError(err).Error("Failed to update user preferences")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to update preferences",
		})
		return
	}

	// Get updated preferences
	updatedPref, err := c.service.GetUserPreference(userID)
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
		"preferences": updatedPref,
	})
}
