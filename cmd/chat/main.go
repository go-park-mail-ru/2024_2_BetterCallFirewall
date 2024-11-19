package main

import (
	"flag"
	"log"

	"github.com/2024_2_BetterCallFirewall/internal/app/chat"
	"github.com/2024_2_BetterCallFirewall/internal/config"
)

func main() {
	confPath := flag.String("c", ".env", "path to config file")
	flag.Parse()

	cfg, err := config.GetConfig(*confPath)
	if err != nil {
		panic(err)
	}

	server, err := chat.GetServer(cfg)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on port %s", cfg.CHAT.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
