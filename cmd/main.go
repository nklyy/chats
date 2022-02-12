package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"noname-realtime-support-chat/internal/health"
	"noname-realtime-support-chat/pkg/logger"
)

func main() {
	// Init logger
	newLogger, err := logger.NewLogger("development")
	if err != nil {
		log.Fatalf("can't create logger: %v", err)
	}

	zapLogger, err := newLogger.SetupZapLogger()
	if err != nil {
		log.Fatalf("can't setup zap logger: %v", err)
	}
	defer zapLogger.Sync()

	// Set-up Route
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Handlers
	healthHandler := health.NewHandler()

	router.Route("/api/v1", func(r chi.Router) {
		healthHandler.SetupRoutes(r)
	})

	// Start App
	zapLogger.Infof("Starting HTTP server on port: %v", 5000)
	err = http.ListenAndServe(":5000", router)
	if err != nil {
		fmt.Println(err)
		zapLogger.Fatalf("Failed to start HTTP server: %v", err)
		return
	}
}
