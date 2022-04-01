package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"log"
	"net/http"
	"noname-realtime-support-chat/config"
	"noname-realtime-support-chat/internal/chat/old-chat"
	"noname-realtime-support-chat/internal/health"
	"noname-realtime-support-chat/internal/user"
	"noname-realtime-support-chat/internal/user/auth"
	"noname-realtime-support-chat/pkg/jwt"
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
	redisAuthClient, err := redis.NewClient(cfg.RedisHost, cfg.RedisPortAuth)
	if err != nil {
		zapLogger.Fatalf("failed to connect to auth redis: %v", err)
	}
	zapLogger.Info("Redis(auth) connected successfully")

	redisChatClient, err := redis.NewClient(cfg.RedisHost, cfg.RedisPortChat)
	if err != nil {
		zapLogger.Fatalf("failed to connect to chat redis: %v", err)
	}
	zapLogger.Info("Redis(chat) connected successfully")

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
	userRepository, err := user.NewRepository(db, cfg.MongoDbName, zapLogger)
	if err != nil {
		zapLogger.Fatalf("failde to create user repository: %v", err)
	}

	// Services
	jwtSvc, err := jwt.NewJwtService(
		cfg.JwtSecretAccess,
		&cfg.JwtExpiryAccess,
		cfg.JwtSecretRefresh,
		&cfg.JwtExpiryRefresh,
		&cfg.AutoLogout,
		redisAuthClient)
	if err != nil {
		zapLogger.Fatalf("failde to jwt service: %v", err)
	}

	userService, err := user.NewService(userRepository, zapLogger, &cfg.Salt)
	if err != nil {
		zapLogger.Fatalf("failde to create user service: %v", err)
	}

	userAuthService, err := auth.NewService(userService, zapLogger, jwtSvc)
	if err != nil {
		zapLogger.Fatalf("failde to create user service: %v", err)
	}

	chatService, err := old_chat.NewService(zapLogger, redisChatClient, make(map[*old_chat.Client]bool), make(map[*old_chat.Room]bool), jwtSvc)
	if err != nil {
		zapLogger.Fatalf("failde to create chat service: %v", err)
	}

	//Middleware
	supportMiddleware, err := user.NewMiddleware(jwtSvc, userService, zapLogger)
	if err != nil {
		zapLogger.Fatalf("failed to set up user middleware %v", err)
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

	userHandler, err := user.NewHandler(userService)
	if err != nil {
		zapLogger.Fatalf("failde to create user handler: %v", err)
	}

	userAuthHandler, err := auth.NewHandler(userAuthService)
	if err != nil {
		zapLogger.Fatalf("failde to create user auth handler: %v", err)
	}

	chatHandler, err := old_chat.NewHandler(chatService)
	if err != nil {
		zapLogger.Fatalf("failde to create chat handler: %v", err)
	}

	router.Route("/api/v1/auth", func(r chi.Router) {
		userAuthHandler.SetupRoutes(r)
	})

	router.Route("/api/v1", func(r chi.Router) {
		supportRoute := r.With(supportMiddleware.JwtMiddleware)

		healthHandler.SetupRoutes(r)
		userHandler.SetupRoutes(supportRoute)
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
