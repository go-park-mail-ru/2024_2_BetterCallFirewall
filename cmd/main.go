package main

import (
	"fmt"
	"log"
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
	"github.com/2024_2_BetterCallFirewall/internal/fileService"
	postController "github.com/2024_2_BetterCallFirewall/internal/post/controller"
	postgresProfile "github.com/2024_2_BetterCallFirewall/internal/post/repository/postgres"
	postServ "github.com/2024_2_BetterCallFirewall/internal/post/service"
	profileController "github.com/2024_2_BetterCallFirewall/internal/profile/controller"
	profileRepository "github.com/2024_2_BetterCallFirewall/internal/profile/repository"
	profileService "github.com/2024_2_BetterCallFirewall/internal/profile/service"
	"github.com/2024_2_BetterCallFirewall/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	redisConn, err := redis.Dial("tcp", "redis:6379")
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, err := postgres.StartPostgres(connStr)
	if err != nil {
		log.Fatalf("Error starting postgres: %v", err)
	}

	repo := postgres.NewAdapter(postgresDB)
	profileRepo := profileRepository.NewProfileRepo(postgresDB)
	authServ := service.NewAuthServiceImpl(repo)

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	}

	responder := router.NewResponder(logger)
	sessionRepo := redismy.NewSessionRedisRepository(redisConn)
	sessionManager := service.NewSessionManager(sessionRepo)
	control := controller.NewAuthController(responder, authServ, sessionManager)

	postRepo := postgresProfile.NewAdapter(postgresDB)

	profileUsecase := profileService.NewProfileUsecase(profileRepo, postRepo)
	profileControl := profileController.NewProfileController(profileUsecase, responder)

	fileServ := fileService.NewFileService()
	postService := postServ.NewPostServiceImpl(postRepo, profileUsecase)
	postControl := postController.NewPostController(postService, responder, fileServ)

	rout := router.NewRouter(control, profileControl, postControl, sessionManager, logger)
	server := http.Server{
		Addr:         ":8080",
		Handler:      rout,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
