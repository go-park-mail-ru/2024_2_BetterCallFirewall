package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/2024_2_BetterCallFirewall/internal/auth/controller"
	"github.com/2024_2_BetterCallFirewall/internal/auth/repository/postgres"
	"github.com/2024_2_BetterCallFirewall/internal/auth/service"
	postController "github.com/2024_2_BetterCallFirewall/internal/post/controller"
	"github.com/2024_2_BetterCallFirewall/internal/post/repository"
	postServ "github.com/2024_2_BetterCallFirewall/internal/post/service"
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

	potgresDB, err := postgres.StartPostgres(connStr)
	if err != nil {
		log.Fatalf("Error starting postgres: %v", err)
	}

	repo := postgres.NewAdapter(potgresDB)
	err = repo.CreateNewUserTable()
	if err != nil {
		log.Fatalf("Error creating user table: %v", err)
	}
	err = repo.CreateNewSessionTable()
	if err != nil {
		log.Fatalf("Error creating session table: %v", err)
	}
	authServ := service.NewAuthServiceImpl(repo)
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	responder := router.NewResponder(logger)
	sessionManager := service.NewSessionManager(repo)
	control := controller.NewAuthController(responder, authServ, sessionManager)

	postRepo := repository.NewRepository()
	postService := postServ.NewPostServiceImpl(postRepo)
	postControl := postController.NewPostController(postService, responder)

	rout := router.NewAuthRouter(control, postControl, sessionManager)
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
