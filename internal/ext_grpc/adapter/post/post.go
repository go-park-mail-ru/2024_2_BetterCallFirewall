package post

import (
	"context"

	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/post_api"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/port/post"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type GrpcSender struct {
	client post_api.PostServiceClient
}

func New(conn grpc.ClientConnInterface) *GrpcSender {
	client := post_api.NewPostServiceClient(conn)

	return &GrpcSender{client: client}
}

func (g *GrpcSender) GetAuthorsPosts(ctx context.Context, header *models.Header, userID uint32) ([]*models.Post, error) {
	req := post.NewRequest(header, userID)
	resp, err := g.client.GetAuthorsPosts(ctx, req)
	if err != nil {
		return nil, err
	}

	res := post.UnmarshalResponse(resp)
	return res, nil
}