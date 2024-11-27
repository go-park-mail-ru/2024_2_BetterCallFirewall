package main

import (
	"flag"
	"log"

	"github.com/2024_2_BetterCallFirewall/internal/app/file"
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
	fileMetrics, err := metrics.NewFileMetrics("file")
	if err != nil {
		panic(err)
	}

	server, err := file.GetServer(cfg, fileMetrics)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on port %s", cfg.FILE.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
