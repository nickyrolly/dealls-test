package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func NewConfig() *viper.Viper {
	config := viper.New()
	config.SetConfigName(getConfigFileName())
	config.SetConfigType("yaml")
	config.AddConfigPath(getConfigFolder())
	err := config.ReadInConfig()

	if err != nil {
		// cannot up if config not exist
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return config
}

func getConfigFileName() string {
	var env = os.Getenv("ENV")

	if env == "" {
		env = "development"
	}
	var filename = fmt.Sprintf("%s.yaml", env)

	return filename
}

func getConfigFolder() string {
	var (
		env = os.Getenv("ENV")
	)
	if env == "" {
		env = "development"
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// use local files in dev
	return filepath.Join(path, "files/config")
}
