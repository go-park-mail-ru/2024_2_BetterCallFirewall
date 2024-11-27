package main

import (
	"flag"
	"log"

	"github.com/2024_2_BetterCallFirewall/internal/app/auth"
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

	authMetrics, err := metrics.NewHTTPMetrics("auth")
	if err != nil {
		panic(err)
	}
	defer authMetrics.ShutDown()
	httpServer, err := auth.GetHTTPServer(cfg, authMetrics)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on port %s", cfg.AUTH.Port)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
