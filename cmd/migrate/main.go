package main

import (
	"errors"
	"os"
	"time"

	"github.com/nickyrolly/dealls-test/internal/config"
	"github.com/nickyrolly/dealls-test/internal/services/profile"
	"github.com/nickyrolly/dealls-test/internal/services/user"
	"gorm.io/gorm"
)

func main() {
	cfg := config.NewConfig()
	log := config.NewLogger(cfg)
	db := config.NewDatabase(config.DatabaseOption{
		Driver:   cfg.GetString("database.driver"),
		DBName:   cfg.GetString("database.name"),
		Username: cfg.GetString("database.username"),
		Password: cfg.GetString("database.password"),
	})

	// runtime env
	env := os.Getenv("ENV")

	// Define models to migrate
	models := []interface{}{
		&user.Entity{},            // Base user model
		&profile.UserProfile{},    // Extended user profile
		&profile.UserPhoto{},      // User photos
		&profile.UserPreference{}, // User preferences
		&profile.UserMatch{},      // User matches
		&profile.UserLike{},       // User likes
	}

	// Migrate all tables
	for _, model := range models {
		if err := db.Migrator().AutoMigrate(model); err != nil {
			log.Errorf("Migration error for %T: %+v", model, err)
			return
		}
		log.Infof("Successfully migrated %T", model)
	}

	// Seeder for development
	if env == "development" {
		var email string = "nickyrolly1@gmail.com"
		var existingUser user.Entity
		result := db.Where("email = ?", email).First(&existingUser)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Record not found is expected for new users
				log.Info("No existing user found, proceeding with creation")
			} else {
				log.Errorf("Error checking for existing user: %+v", result.Error)
				return
			}
		}

		if result.RowsAffected > 0 {
			log.Warnf("Development seed data already exists for email: %s", email)
			return
		}

		// Create new user if not found
		newUser := user.Entity{
			Email:        email,
			PasswordHash: "$2a$10$6jM7G7hXMBQGBOgCYR.tB.cV3IgzKGx4ZhkqG4GHion0/eqfGHhPi", // hashed value for "password123"
			FirstName:    "Admin",
			DateOfBirth:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:       "other",
		}

		result = db.Create(&newUser)
		if result.Error != nil {
			log.Errorf("Failed to create seed user: %+v", result.Error)
			return
		}

		log.Infof("Successfully created seed user with email: %s", email)
	}
}
