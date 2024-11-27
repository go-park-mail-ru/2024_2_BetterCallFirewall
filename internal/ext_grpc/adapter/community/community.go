package community

import (
	"context"

	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/community_api"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/port/community"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type GrpcSender struct {
	client community_api.CommunityServiceClient
}

func New(conn grpc.ClientConnInterface) *GrpcSender {
	client := community_api.NewCommunityServiceClient(conn)

	return &GrpcSender{
		client: client,
	}
}

func (g *GrpcSender) CheckAccess(ctx context.Context, communityID, userID uint32) bool {
	req := community.NewRequest(communityID, userID)
	resp, err := g.client.CheckAccess(ctx, req)
	if err != nil {
		return false
	}

	res := community.UnmarshallResponse(resp)
	return res
}

func (g *GrpcSender) GetHeader(ctx context.Context, communityID uint32) (*models.Header, error) {
	req := community.NewHeaderRequest(communityID)
	resp, err := g.client.GetHeader(ctx, req)
	if err != nil {
		return nil, err
	}

	res := community.UnmarshallHeaderResponse(resp)
	return res, nil
}
