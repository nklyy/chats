package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"nn-realtime-support-chat/internal/health"
)

func main() {
	// Init logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// Set-up Route
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Handlers
	healthHandler := health.NewHandler()

	router.Route("/api/v1", func(r chi.Router) {
		healthHandler.SetupRoutes(r)
	})

	// Start App
	err = http.ListenAndServe(":5000", router)
	if err != nil {
		fmt.Println(err)
		logger.Fatal("Failed to start HTTP server!", zap.Error(err))
		return
	}
}
