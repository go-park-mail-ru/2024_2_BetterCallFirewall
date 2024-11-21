package main

import (
	"flag"
	"log"

	"github.com/2024_2_BetterCallFirewall/internal/app/auth"
	"github.com/2024_2_BetterCallFirewall/internal/config"
)

func main() {
	confPath := flag.String("c", ".env", "path to config file")
	flag.Parse()

	cfg, err := config.GetConfig(*confPath)
	if err != nil {
		panic(err)
	}

	httpServer, err := auth.GetHTTPServer(cfg)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on port %s", cfg.AUTH.Port)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
