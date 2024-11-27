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
	"github.com/2024_2_BetterCallFirewall/internal/metrics"
)

func main() {
	confPath := flag.String("c", ".env", "path to config file")
	flag.Parse()

	cfg, err := config.GetConfig(*confPath)
	if err != nil {
		panic(err)
	}

	grpcMetrics, err := metrics.NewGrpcMetrics("auth")
	if err != nil {
		panic(err)
	}
	defer grpcMetrics.ShutDown()

	grpcServer, err := auth.GetGRPCServer(cfg, grpcMetrics)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.AUTHGRPC.Port))
	if err != nil {
		panic(err)
	}
	go func() {
		http.Handle("/api/v1/metrics", promhttp.Handler())
		http.Handle(
			"/", http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			),
		)
		http.ListenAndServe(":6001", nil)
	}()

	log.Printf("Listening on :%s with protocol gRPC", cfg.AUTHGRPC.Port)
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}
