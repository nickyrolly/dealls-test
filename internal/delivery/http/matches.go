package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	gorilla_context "github.com/gorilla/context"
	"github.com/nickyrolly/dealls-test/common"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/middleware"
	"github.com/nickyrolly/dealls-test/internal/services/profile"
	"github.com/sirupsen/logrus"
)

type MatchController struct {
	log     *logrus.Logger
	service *profile.Service
}

func NewMatchController(log *logrus.Logger, service *profile.Service) *MatchController {
	return &MatchController{
		log:     log,
		service: service,
	}
}

type likeRequest struct {
	LikedUserID string `json:"liked_user_id"`
}

func (c *MatchController) HandleLike(w http.ResponseWriter, r *http.Request) {
	userSession, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "User session not found",
		})
		return
	}

	var req likeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.WithError(err).Error("Failed to decode request body")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	likerID, err := uuid.Parse(userSession.ID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse liker ID")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid liker ID",
		})
		return
	}

	likedID, err := uuid.Parse(req.LikedUserID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse liked user ID")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid liked user ID",
		})
		return
	}

	like := &profile.UserLike{
		LikerID: likerID,
		LikedID: likedID,
	}

	match, err := c.service.CreateLike(like)
	if err != nil {
		c.log.WithError(err).Error("Failed to create like")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to create like",
		})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"like": like,
		},
	}

	if match != nil {
		response["data"].(map[string]interface{})["match"] = match
		response["data"].(map[string]interface{})["is_match"] = true
	}

	common.CustomResponseAPI(w, r, http.StatusOK, response)
}

func (c *MatchController) HandleWithdrawLike(w http.ResponseWriter, r *http.Request) {
	userSession, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "User session not found",
		})
		return
	}

	likedUserID := chi.URLParam(r, "userId")
	if likedUserID == "" {
		c.log.Error("Missing user ID parameter")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Missing user ID",
		})
		return
	}

	likerID, err := uuid.Parse(userSession.ID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse liker ID")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	likedID, err := uuid.Parse(likedUserID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse liked user ID")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid liked user ID",
		})
		return
	}

	if err := c.service.WithdrawLike(likerID, likedID); err != nil {
		c.log.WithError(err).Error("Failed to withdraw like")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to withdraw like",
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Like withdrawn successfully",
	})
}

func (c *MatchController) HandleGetMatches(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := gorilla_context.Get(r, "user").(string)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "User session not found",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse user ID")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	matches, err := c.service.GetMatches(userID)
	if err != nil {
		c.log.WithError(err).Error("Failed to get matches")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to get matches",
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"matches": matches,
		},
	})
}

func (c *MatchController) HandleUnmatch(w http.ResponseWriter, r *http.Request) {
	userSession, ok := gorilla_context.Get(r, "user").(middleware.UserSession)
	if !ok {
		c.log.Error("User session not found in context")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "User session not found",
		})
		return
	}

	matchedUserID := chi.URLParam(r, "userId")
	if matchedUserID == "" {
		c.log.Error("Missing user ID parameter")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Missing user ID",
		})
		return
	}

	user1ID, err := uuid.Parse(userSession.ID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse user ID")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	user2ID, err := uuid.Parse(matchedUserID)
	if err != nil {
		c.log.WithError(err).Error("Failed to parse matched user ID")
		common.CustomResponseAPI(w, r, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid matched user ID",
		})
		return
	}

	if err := c.service.UnmatchUsers(user1ID, user2ID); err != nil {
		c.log.WithError(err).Error("Failed to unmatch users")
		common.CustomResponseAPI(w, r, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to unmatch users",
		})
		return
	}

	common.CustomResponseAPI(w, r, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Users unmatched successfully",
	})
}
