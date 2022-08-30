package main

import (
	"errors"
	"log"
	"net/http"
	"one-time-session-chat/config"
	"one-time-session-chat/internal/chat"
	"one-time-session-chat/internal/chat/room"
	"one-time-session-chat/internal/health"

	"one-time-session-chat/pkg/logger"
	"one-time-session-chat/pkg/redis"
	"os"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init logger
	newLogger, err := logger.NewLogger(cfg.Environment)
	if err != nil {
		log.Fatalf("can't create logger: %v", err)
	}

	zapLogger, err := newLogger.SetupZapLogger()
	if err != nil {
		log.Fatalf("can't setup zap logger: %v", err)
	}
	defer func(zapLogger *zap.SugaredLogger) {
		err := zapLogger.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Fatalf("can't setup zap logger: %v", err)
		}
	}(zapLogger)

	// Create roomKey folder if bot exist
	_, err = os.Stat("keys")
	if os.IsNotExist(err) {
		err = os.Mkdir("keys", 0755)
		if err != nil {
			zapLogger.Fatalf("failed to create roomKeys folder %v", err)
		}
	}

	// Redis
	redisChatClient, err := redis.NewClient(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		zapLogger.Fatalf("failed to connect to chat redis: %v", err)
	}
	zapLogger.Info("Redis(chat) connected successfully")

	// Services
	roomService, err := room.NewService(zapLogger)
	if err != nil {
		zapLogger.Fatalf("failed to set up room service %v", err)
	}

	chatService, err := chat.NewService(redisChatClient, roomService, zapLogger)
	if err != nil {
		zapLogger.Fatalf("failed to set up chat service %v", err)
	}

	// Set-up Route
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Content-Type", "JWT-Token"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Handlers
	healthHandler := health.NewHandler()

	chatHandler, err := chat.NewHandler(chatService)
	if err != nil {
		zapLogger.Fatalf("failed to set up chat handler %v", err)
	}

	// Routes
	router.Route("/api/v1", func(r chi.Router) {
		healthHandler.SetupRoutes(r)
	})

	router.Route("/", func(r chi.Router) {
		chatHandler.SetupRoutes(r)
	})

	// Start App
	zapLogger.Infof("Starting HTTP server on port: %v", cfg.PORT)
	err = http.ListenAndServe(":"+cfg.PORT, router)
	if err != nil {
		zapLogger.Fatalf("Failed to start HTTP server: %v", err)
		return
	}
}
