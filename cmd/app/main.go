package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// fmt.Println("Service Running on : ", ":"+cfg.GetString("application.port"))
	// log.Fatal(http.ListenAndServe(":"+cfg.GetString("application.port"), router))

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.GetString("application.port"),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting server on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exiting")
}
