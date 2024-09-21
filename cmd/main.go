package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/controller"
	"github.com/2024_2_BetterCallFirewall/internal/repository"
	"github.com/2024_2_BetterCallFirewall/internal/service"
)

func main() {
	repo := repository.NewSampleDB()
	authServ := service.NewAuthServiceImpl(repo)
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	responder := controller.NewResponder(logger)
	control := controller.NewAuthController(responder, authServ)
	router := controller.NewAuthRouter(control)

	server := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
