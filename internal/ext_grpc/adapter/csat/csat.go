package csat

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/csat_api"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/port/csat"
)

type GRPCSender struct {
	client csat_api.CsatServiceClient
}

func New(conn grpc.ClientConnInterface) *GRPCSender {
	client := csat_api.NewCsatServiceClient(conn)

	return &GRPCSender{client: client}
}

func (g *GRPCSender) NewLike(id uint32) {
	req := csat.NewRequest(id)
	_, err := g.client.NewLike(context.Background(), req)
	if err != nil {
		log.Println(err)
	}
}

func (g *GRPCSender) NewMessage(id uint32) {
	req := csat.NewRequest(id)
	_, err := g.client.NewMessage(context.Background(), req)
	if err != nil {
		log.Println(err)
	}
}

func (g *GRPCSender) NewFriend(id uint32) {
	req := csat.NewRequest(id)
	_, err := g.client.NewFriend(context.Background(), req)
	if err != nil {
		log.Println(err)
	}
}

func GetCSATProvider(host, port string) (grpc.ClientConnInterface, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return conn, nil
}
