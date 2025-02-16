package authentication

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nickyrolly/dealls-test/internal/services/authentication"
	"github.com/sirupsen/logrus"
)

type SignupRequest struct {
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required,min=8"`
	FirstName   string    `json:"first_name" validate:"required"`
	LastName    *string   `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
	Gender      string    `json:"gender" validate:"required,oneof=male female other"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Controller struct {
	service *authentication.Service
	log     *logrus.Logger
}

func NewController(service *authentication.Service, log *logrus.Logger) *Controller {
	return &Controller{
		service: service,
		log:     log,
	}
}

func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.WithError(err).Error("Failed to decode signup request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, token, err := c.service.SignUp(
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		req.DateOfBirth,
		req.Gender,
	)
	if err != nil {
		c.log.WithError(err).Error("Failed to sign up user")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.WithError(err).Error("Failed to decode login request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, token, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		c.log.WithError(err).Error("Failed to login user")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}
