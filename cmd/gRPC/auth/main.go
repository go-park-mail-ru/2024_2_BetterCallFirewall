package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	grpcServer, err := auth.GetGRPCServer(cfg)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.AUTHGRPC.Port))
	if err != nil {
		panic(err)
	}
	go func() {
		http.Handle("/api/v1/metrics", promhttp.Handler())
		http.ListenAndServe(":6001", nil)
	}()

	log.Printf("Listening on :%s with protocol gRPC", cfg.AUTHGRPC.Port)
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}
