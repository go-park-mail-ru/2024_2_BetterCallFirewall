package post

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/post_api"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/port/post"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type GrpcSender struct {
	client post_api.PostServiceClient
}

func New(conn grpc.ClientConnInterface) *GrpcSender {
	client := post_api.NewPostServiceClient(conn)

	return &GrpcSender{client: client}
}

func (g *GrpcSender) GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error) {
	req := post.NewRequest(header)
	resp, err := g.client.GetAuthorsPosts(ctx, req)
	if err != nil {
		return nil, err
	}

	res := post.UnmarshalResponse(resp)
	return res, nil
}

func GetPostProvider(host, port string) (grpc.ClientConnInterface, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return conn, nil
}
