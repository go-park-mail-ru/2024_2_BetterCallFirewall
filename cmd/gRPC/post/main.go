package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/2024_2_BetterCallFirewall/internal/app/post"
	"github.com/2024_2_BetterCallFirewall/internal/config"
)

func main() {
	confPath := flag.String("c", ".env", "path to config file")
	flag.Parse()

	cfg, err := config.GetConfig(*confPath)
	if err != nil {
		panic(err)
	}

	grpcServer, err := post.GetGRPCServer(cfg)
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.POSTGRPC))
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on :%s with protocol gRPC", cfg.POSTGRPC)
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}
