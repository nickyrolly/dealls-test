package subscription

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	gorilla_context "github.com/gorilla/context"
	"github.com/nickyrolly/dealls-test/common"
	"github.com/nickyrolly/dealls-test/internal/services/subscription"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	log     *logrus.Logger
	service *subscription.Service
}

func NewController(log *logrus.Logger, service *subscription.Service) *Controller {
	return &Controller{
		log:     log,
		service: service,
	}
}

func (c *Controller) HandleSubscription(w http.ResponseWriter, r *http.Request) {
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
		IsSubscribed bool `json:"is_subscribed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// Handle error
		c.log.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.CreateSubscription(userID, request.IsSubscribed); err != nil {
		// Handle error
		c.log.WithError(err).Error("Failed to create subscription")
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
