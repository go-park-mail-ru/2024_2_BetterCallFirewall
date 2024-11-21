package main

import (
	"flag"
	"log"

	"github.com/2024_2_BetterCallFirewall/internal/app/community"
	"github.com/2024_2_BetterCallFirewall/internal/config"
)

func main() {
	confPath := flag.String("c", ".env", "path to config file")
	flag.Parse()

	cfg, err := config.GetConfig(*confPath)
	if err != nil {
		panic(err)
	}

	httpServer, err := community.GetHTTPServer(cfg)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on posrt %s", cfg.COMMUNITY.Port)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
