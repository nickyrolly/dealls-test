package main

import (
	"os"

	"github.com/nickyrolly/dealls-test/internal/config"
	"github.com/nickyrolly/dealls-test/internal/services/user"
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

	// Define models
	var userModel = &user.Entity{}
	// var memberModel = &member.Entity{}
	// var orderModel = &order.Entity{}

	if err := db.Migrator().AutoMigrate(userModel); err != nil {
		log.Errorf("Migration error : %+v", err)
	}

	// Seeder for development
	if env == "development" {
		var email string = "nickyrolly1@gmail.com"
		result := db.First(userModel, "email = ?", email)
		if result.Error != nil {
			log.Errorf("Create seed development error : %+v", result.Error)
		}

		if result.RowsAffected != 0 {
			log.Warnf("Development seed data already exist")
		}

		if result.RowsAffected < 1 {
			var newUser = user.Entity{
				Username: "administrator",
				Email:    email,
			}

			result := db.Create(&newUser)
			if result.Error != nil {
				log.Errorf("Create database error : %+v", result.Error)
				panic(result.Error)
			}
		}
	}
}
