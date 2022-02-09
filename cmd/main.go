package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"nn-realtime-support-chat/internal/health"
)

func main() {
	// Init logger
	logger := logrus.New()

	// Set-up Route
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Handlers
	healthHandler := health.NewHandler()

	router.Route("/api/v1", func(r chi.Router) {
		healthHandler.SetupRoutes(r)
	})

	// Start App
	err := http.ListenAndServe(":5000", router)
	if err != nil {
		fmt.Println(err)
		logger.Fatalln("Failed to start HTTP server!")
		return
	}
}
