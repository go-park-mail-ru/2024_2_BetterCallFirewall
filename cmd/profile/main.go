package main

import (
	"flag"
	"log"

	"github.com/2024_2_BetterCallFirewall/internal/app/profile"
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

	metric, err := metrics.NewHTTPMetrics("profile")
	if err != nil {
		panic(err)
	}

	server, err := profile.GetHTTPServer(cfg, metric)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on posrt %s", cfg.PROFILE.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
