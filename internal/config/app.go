package config

import (
	http "github.com/nickyrolly/dealls-test/internal/delivery/http"
	"gorm.io/gorm"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	Config *viper.Viper
	Router *chi.Mux
	Log    *logrus.Logger
	DB     *gorm.DB
}

func Bootstrap(config *BootstrapConfig) {
	route := &http.RouteConfig{
		Router: config.Router,
	}

	route.Setup()
}
