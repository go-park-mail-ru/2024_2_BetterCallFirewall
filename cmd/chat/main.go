package main

import (
	"flag"
	"log"

	"github.com/2024_2_BetterCallFirewall/internal/app/chat"
	"github.com/2024_2_BetterCallFirewall/internal/config"
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
)

func main() {
	confPath := flag.String("c", ".env", "path to config file")
	flag.Parse()

	cfg, err := config.GetConfig(*confPath)
	if err != nil {
		panic(err)
	}

	chatMetrics, err := metrics.NewHTTPMetrics("chat")
	if err != nil {
		panic(err)
	}
	defer chatMetrics.ShutDown()

	server, err := chat.GetServer(cfg, chatMetrics)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on port %s", cfg.CHAT.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
