package community

import (
	"context"

	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/community_api"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/port/community"
)

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
