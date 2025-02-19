package config

import (
	"github.com/go-chi/chi"
	"github.com/gomodule/redigo/redis"
	"github.com/nickyrolly/dealls-test/internal/delivery/http"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/authentication"
	"github.com/nickyrolly/dealls-test/internal/delivery/http/profile"
	authService "github.com/nickyrolly/dealls-test/internal/services/authentication"
	profileService "github.com/nickyrolly/dealls-test/internal/services/profile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Config       *viper.Viper
	Router       *chi.Mux
	Log          *logrus.Logger
	DB           *gorm.DB
	RedisGeneral *redis.Pool
}

func Bootstrap(config *BootstrapConfig) {
	// Initialize services
	authSvc := authService.NewService(config.DB, config.Log)
	profileSvc := profileService.NewService(config.DB, config.Log)

	// Initialize controllers
	authController := authentication.NewController(
		authSvc,
		config.Log,
	)

	profileController := profile.NewController(
		config.Log,
		profileSvc,
	)

	matchController := http.NewMatchController(
		config.Log,
		profileSvc,
	)

	// Initialize route configuration with all dependencies
	route := &http.Config{
		Router:                   config.Router,
		AuthenticationController: authController,
		ProfileController:        profileController,
		MatchesController:        matchController,
		RedisPool:                config.RedisGeneral,
		DB:                       config.DB,
	}

	// Setup routes
	if err := http.Setup(route); err != nil {
		config.Log.Fatalf("Failed to setup routes: %v", err)
	}
}
