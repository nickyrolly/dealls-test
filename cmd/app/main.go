package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	config "github.com/nickyrolly/dealls-test/internal/config"
)

func main() {

	fmt.Println("Hello, World!")

	cfg := config.NewConfig()
	log := config.NewLogger(cfg)
	router := chi.NewRouter()

	db := config.NewDatabase(config.DatabaseOption{
		Driver:   cfg.GetString("database.driver"),
		DBName:   cfg.GetString("database.name"),
		Username: cfg.GetString("database.username"),
		Password: cfg.GetString("database.password"),
	})

	redis := config.NewRedis(cfg)

	config.Bootstrap(&config.BootstrapConfig{
		Config:       cfg,
		Router:       router,
		Log:          log,
		DB:           db,
		RedisGeneral: redis,
	})

	fmt.Println("Service Running on : ", ":"+cfg.GetString("application.port"))
	log.Fatal(http.ListenAndServe(":"+cfg.GetString("application.port"), router))
}
