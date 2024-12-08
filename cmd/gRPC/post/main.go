package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	grpcMetrics, err := metrics.NewGrpcMetrics("post")
	if err != nil {
		panic(err)
	}

	grpcServer, err := post.GetGRPCServer(cfg, grpcMetrics)
	if err != nil {
		panic(err)
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.POSTGRPC.Port))
	if err != nil {
		panic(err)
	}
	go func() {
		http.Handle("/api/v1/metrics", promhttp.Handler())
		if err = http.ListenAndServe(":6002", nil); err != nil {
			panic(err)
		}
	}()

	log.Printf("Listening on :%s with protocol gRPC", cfg.POSTGRPC.Port)
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}
