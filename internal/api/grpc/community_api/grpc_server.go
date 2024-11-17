package community_api

import (
	"context"
)

type CommunityService interface {
	CheckAccess(ctx context.Context, communityID, userID uint32) bool
}

type Adapter struct {
	UnimplementedCommunityServiceServer
	serv CommunityService
}

func New(s CommunityService) *Adapter {
	return &Adapter{
		serv: s,
	}
}

func (a *Adapter) CheckAccess(ctx context.Context, req *CheckAccessRequest) (*CheckAccessResponse, error) {
	res := a.serv.CheckAccess(ctx, req.CommunityID, req.UserID)

	return &CheckAccessResponse{Access: res}, nil
}
