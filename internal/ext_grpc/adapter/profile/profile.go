package profile

import (
	"context"

	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/profile_api"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/port/profile"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type GrpcSender struct {
	client profile_api.ProfileServiceClient
}

func New(conn grpc.ClientConnInterface) *GrpcSender {
	client := profile_api.NewProfileServiceClient(conn)

	return &GrpcSender{
		client: client,
	}
}

func (g *GrpcSender) GetHeader(ctx context.Context, userID uint32) (models.Header, error) {
	req := profile.NewGetHeaderRequest(userID)
	resp, err := g.client.GetHeader(ctx, req)
	if err != nil {
		return models.Header{}, err
	}

	res := profile.UnmarshallHeaderResponse(resp)
	return *res, nil
}

func (g *GrpcSender) GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error) {
	req := profile.NewGetFriendsIDRequest(userID)
	resp, err := g.client.GetFriendsID(ctx, req)
	if err != nil {
		return nil, err
	}

	res := profile.UnmarshallGetFriendsIDResponse(resp)
	return res, nil
}
