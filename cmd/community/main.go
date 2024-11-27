package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/2024_2_BetterCallFirewall/internal/app/community"
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
	communityMetrics, err := metrics.NewHTTPMetrics("community")
	if err != nil {
		panic(err)
	}
	defer communityMetrics.ShutDown()

	grpcMetrics, err := metrics.NewGrpcMetrics("community")
	if err != nil {
		panic(err)
	}

	httpServer, grpcServer, err := community.GetServers(cfg, grpcMetrics, communityMetrics)
	if err != nil {
		panic(err)
	}

	go func() {
		defer grpcMetrics.ShutDown()
		l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.COMMUNITYGRPC.Port))
		if err != nil {
			panic(err)
		}

		log.Printf("Listening on :%s with protocol gRPC", cfg.COMMUNITYGRPC.Port)
		if err := grpcServer.Serve(l); err != nil {
			panic(err)
		}
	}()

	log.Printf("Starting server on posrt %s", cfg.COMMUNITY.Port)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
