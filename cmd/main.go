package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"noname-realtime-support-chat/config"
	"noname-realtime-support-chat/internal/chat"
	"noname-realtime-support-chat/internal/health"
	"noname-realtime-support-chat/internal/support"
	"noname-realtime-support-chat/internal/support/auth"
	"noname-realtime-support-chat/internal/support/jwt"
	"noname-realtime-support-chat/pkg/logger"
	"noname-realtime-support-chat/pkg/mongodb"
	"noname-realtime-support-chat/pkg/redis"
	"syscall"
)

func main() {
	cfg, err := config.Get(".")
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

	// Connect to database
	db, ctx, cancel, err := mongodb.NewConnection(cfg)
	if err != nil {
		zapLogger.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer mongodb.Close(db, ctx, cancel)

	// Ping db
	err = mongodb.Ping(db, ctx)
	if err != nil {
		log.Fatal(err)
	}
	zapLogger.Info("DB connected successfully")

	// Redis
	redisClient, err := redis.NewClient(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		zapLogger.Fatalf("failed to connect to redis: %v", err)
	}
	zapLogger.Info("Redis connected successfully")

	// RabbitMq
	//rabbitmqConnection, err := rabbitmq.NewConnection(cfg.RabbitMqUrl)
	//if err != nil {
	//	zapLogger.Fatalf("can't connect to amqp host: %v", err)
	//}
	//
	//rabbitmqChannel, err := rabbitmq.NewChanel(rabbitmqConnection)
	//if err != nil {
	//	zapLogger.Fatalf("can't create amqp channel: %v", err)
	//}
	//defer rabbitmq.Close(rabbitmqConnection, rabbitmqChannel)
	//zapLogger.Info("RabbitMq connected successfully")

	// Repositories
	supportRepository, err := support.NewRepository(db, cfg.MongoDbName, zapLogger)
	if err != nil {
		zapLogger.Fatalf("failde to create support repository: %v", err)
	}

	// Services
	jwtSvc, err := jwt.NewJwtService(
		cfg.JwtSecretAccess,
		&cfg.JwtExpiryAccess,
		cfg.JwtSecretRefresh,
		&cfg.JwtExpiryRefresh,
		&cfg.AutoLogout,
		redisClient)
	if err != nil {
		zapLogger.Fatalf("failde to jwt service: %v", err)
	}

	supportService, err := support.NewService(supportRepository, zapLogger, &cfg.Salt)
	if err != nil {
		zapLogger.Fatalf("failde to create support service: %v", err)
	}

	supportAuthService, err := auth.NewService(supportService, zapLogger, jwtSvc)
	if err != nil {
		zapLogger.Fatalf("failde to create support service: %v", err)
	}

	//Middleware
	supportMiddleware, err := support.NewMiddleware(jwtSvc, zapLogger)
	if err != nil {
		zapLogger.Fatalf("failed to set up support middleware %v", err)
	}

	// Set-up Route
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Handlers
	healthHandler := health.NewHandler()

	supportHandler, err := support.NewHandler(supportService)
	if err != nil {
		zapLogger.Fatalf("failde to create support handler: %v", err)
	}

	supportAuthHandler, err := auth.NewHandler(supportAuthService)
	if err != nil {
		zapLogger.Fatalf("failde to create support auth handler: %v", err)
	}

	chatHandler, err := chat.NewHandler()
	if err != nil {
		zapLogger.Fatalf("failde to create chat handler: %v", err)
	}

	router.Route("/api/v1/auth", func(r chi.Router) {
		supportAuthHandler.SetupRoutes(r)
	})

	router.Route("/api/v1", func(r chi.Router) {
		supportRoute := r.With(supportMiddleware.JwtMiddleware)

		healthHandler.SetupRoutes(r)
		supportHandler.SetupRoutes(supportRoute)
		chatHandler.SetupRoutes(r)
	})

	// Start App
	zapLogger.Infof("Starting HTTP server on port: %v", cfg.PORT)
	err = http.ListenAndServe(cfg.PORT, router)
	if err != nil {
		fmt.Println(err)
		zapLogger.Fatalf("Failed to start HTTP server: %v", err)
		return
	}
}
