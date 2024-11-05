package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/auth/controller"
	"github.com/2024_2_BetterCallFirewall/internal/auth/repository/postgres"
	redismy "github.com/2024_2_BetterCallFirewall/internal/auth/repository/redis"
	"github.com/2024_2_BetterCallFirewall/internal/auth/service"
	ChatController "github.com/2024_2_BetterCallFirewall/internal/chat/controller"
	chatRepository "github.com/2024_2_BetterCallFirewall/internal/chat/repository/postgres"
	chatService "github.com/2024_2_BetterCallFirewall/internal/chat/service"
	filecontrol "github.com/2024_2_BetterCallFirewall/internal/fileService/controller"
	fileRepo "github.com/2024_2_BetterCallFirewall/internal/fileService/repository"
	fileservis "github.com/2024_2_BetterCallFirewall/internal/fileService/service"
	postController "github.com/2024_2_BetterCallFirewall/internal/post/controller"
	postgresProfile "github.com/2024_2_BetterCallFirewall/internal/post/repository/postgres"
	postServ "github.com/2024_2_BetterCallFirewall/internal/post/service"
	profileController "github.com/2024_2_BetterCallFirewall/internal/profile/controller"
	profileRepository "github.com/2024_2_BetterCallFirewall/internal/profile/repository"
	profileService "github.com/2024_2_BetterCallFirewall/internal/profile/service"
	"github.com/2024_2_BetterCallFirewall/internal/router"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	redisPool := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "redis:6379")
		},
	}

	postgresDB, err := StartPostgres(connStr, logger)
	if err != nil {
		logger.Fatalf("Error starting postgres: %v", err)
	}

	repo := postgres.NewAdapter(postgresDB)
	profileRepo := profileRepository.NewProfileRepo(postgresDB)
	authServ := service.NewAuthServiceImpl(repo)

	responder := router.NewResponder(logger)
	sessionRepo := redismy.NewSessionRedisRepository(redisPool)
	sessionManager := service.NewSessionManager(sessionRepo)
	control := controller.NewAuthController(responder, authServ, sessionManager)

	postRepo := postgresProfile.NewAdapter(postgresDB)
	chatRepo := chatRepository.NewChatRepository(postgresDB)

	fileRepository := fileRepo.NewFileRepo(postgresDB)
	fileServ := fileservis.NewFileService(fileRepository)
	fileController := filecontrol.NewFileController(fileServ, responder)

	postsHelper := postServ.NewPostProfileImpl(fileServ, postRepo)
	profileUsecase := profileService.NewProfileUsecase(profileRepo, postsHelper)
	profileControl := profileController.NewProfileController(profileUsecase, fileServ, responder)

	chatService := chatService.NewChatService(chatRepo, profileUsecase)
	chatControl := ChatController.NewChatController(chatService, responder)

	postService := postServ.NewPostServiceImpl(postRepo, profileUsecase)
	postControl := postController.NewPostController(postService, responder, fileServ)

	rout := router.NewRouter(control, profileControl, postControl, fileController, sessionManager, logger, chatControl)
	server := http.Server{
		Addr:         ":8080",
		Handler:      rout,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("listen: %s\n", err)
	}
}

func StartPostgres(connStr string, logger *logrus.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	db.SetMaxOpenConns(10)

	retrying := 10
	i := 1
	logger.Infof("try ping postgresql:%v", i)
	for err = db.Ping(); err != nil; err = db.Ping() {
		if i >= retrying {
			return nil, fmt.Errorf("postgres connect: %w", err)
		}
		i++
		time.Sleep(1 * time.Second)
		logger.Infof("try ping postgresql: %v", i)
	}

	return db, nil
}
