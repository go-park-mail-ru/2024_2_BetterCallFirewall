package main

import (
	"flag"
	"fmt"
	"log"
	"net"

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

	grpcServer, err := community.GetGRPCServer(cfg)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.POSTGRPC.Port))
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on :%s with protocol gRPC", cfg.COMMUNITYGRPC.Port)
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}
