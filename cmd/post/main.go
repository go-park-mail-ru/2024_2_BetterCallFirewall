package main

import (
	"flag"
	"log"

	"github.com/2024_2_BetterCallFirewall/internal/app/post"
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
	postMetric, err := metrics.NewHTTPMetrics("post")
	if err != nil {
		panic(err)
	}

	server, err := post.GetHTTPServer(cfg, postMetric)
	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on posrt %s", cfg.POST.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
