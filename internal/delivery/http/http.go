package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gomodule/redigo/redis"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/authentication"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/healthcheck"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/middleware"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/profile"
	authService "github.com/nickyrolly/dealls-test/internal/services/authentication"
	profileService "github.com/nickyrolly/dealls-test/internal/services/profile"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	Router                   *chi.Mux
	DB                       *gorm.DB
	RedisPool                *redis.Pool
	AuthenticationController *authentication.Controller
	ProfileController        *profile.Controller
}

func NewRouteConfig(router *chi.Mux, redisPool *redis.Pool, db *gorm.DB) *Config {
	// Get standard logger
	log := logrus.StandardLogger()

	// Initialize services
	authSvc := authService.NewService(db, log)
	profileSvc := profileService.NewService(db, log)

	// Initialize controllers
	authController := authentication.NewController(authSvc, log)
	profileController := profile.NewController(log, profileSvc)

	// Create route config
	config := &Config{
		Router:                   router,
		DB:                       db,
		RedisPool:                redisPool,
		AuthenticationController: authController,
		ProfileController:        profileController,
	}

	return config
}

func Setup(c *Config) error {
	fmt.Println("========== Dealls Setup API ==========")

	// Initialize middleware
	mw := middleware.NewMiddleware(logrus.StandardLogger(), c.DB, c.RedisPool, "your-jwt-secret-here")

	// Apply global middleware
	c.Router.Use(mw.Logger)
	c.Router.Use(mw.Recover)
	c.Router.Use(mw.CORS)

	// Health check endpoint
	c.Router.Get("/", healthcheck.HandleHealthCheck)

	// API routes
	c.Router.Route("/api", func(r chi.Router) {
		// Public routes (no auth required)
		r.Group(func(r chi.Router) {
			r.Use(mw.Logger)
			r.Use(mw.Recover)
			r.Use(mw.CORS)
			r.Post("/signup", c.AuthenticationController.SignUp)
		})

		// Authentication routes
		r.Group(func(r chi.Router) {
			r.Use(mw.JWT)
			r.Post("/login", c.AuthenticationController.Login)
		})

		// API v1 routes
		r.Route("/v1", func(r chi.Router) {
			// Apply authentication check for all /v1 routes
			r.Use(mw.JWT)

			// Handle preflight requests for CORS
			r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Profile routes
			r.Route("/profile", func(r chi.Router) {
				r.Get("/", c.ProfileController.HandleGetProfile)
				r.Put("/", c.ProfileController.HandleUpdateProfile)
				r.Get("/matches", c.ProfileController.HandleGetMatches)
				r.Post("/like/{id}", c.ProfileController.HandleLikeProfile)
				r.Post("/pass/{id}", c.ProfileController.HandlePassProfile)
				r.Get("/potential", c.ProfileController.HandleGetPotentialMatches)
			})
		})
	})

	return nil
}
