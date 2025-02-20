package swipe

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	gorilla_context "github.com/gorilla/context"
	"github.com/nickyrolly/dealls-test/common"
	"github.com/nickyrolly/dealls-test/internal/services/subscription"
	"github.com/nickyrolly/dealls-test/internal/services/swipe"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	log                 *logrus.Logger
	service             *swipe.Service
	serviceSubscription *subscription.Service
}

func NewController(log *logrus.Logger, service *swipe.Service, serviceSubscription *subscription.Service) *Controller {
	return &Controller{
		log:                 log,
		service:             service,
		serviceSubscription: serviceSubscription,
	}
}

func (c *Controller) HandleSwipe(w http.ResponseWriter, r *http.Request) {
	userIDStr := gorilla_context.Get(r, "user").(string) // Get user ID from context
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		// Handle error
		c.log.WithError(err).Error("Failed to parse user ID")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Assuming you get profileID and action from the request body
	var request struct {
		ProfileID uuid.UUID `json:"profile_id"`
		Action    string    `json:"action"` // "like" or "pass"
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// Handle error
		c.log.WithError(err).Error("Failed to decode request body")
		//http.Error(w, "Invalid request body", http.StatusBadRequest)

		common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
			"success": false,
		})
		return
	}

	var isSubscribed bool
	if isSubscribedSubscription, err := c.serviceSubscription.CheckSubscription(userID); err != nil {
		// Handle error
		c.log.WithError(err).Error("Failed to get subscription")

		common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
			"success": false,
		})
		return
	} else {
		isSubscribed = isSubscribedSubscription
	}

	// Call the service to create a swipe
	if err := c.service.CreateSwipe(userID, request.ProfileID, request.Action, isSubscribed); err != nil {
		// Handle error
		c.log.WithError(err).Error("Failed to swipe user profile")

		common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
			"success": false,
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
