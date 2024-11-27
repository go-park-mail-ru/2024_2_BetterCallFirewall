package community_api

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type CommunityService interface {
	CheckAccess(ctx context.Context, communityID, userID uint32) bool
	GetHeader(ctx context.Context, communityID uint32) (*models.Header, error)
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

func (a *Adapter) GetHeader(ctx context.Context, req *GetHeaderRequest) (*GetHeaderResponse, error) {
	res, err := a.serv.GetHeader(ctx, req.CommunityID)
	if err != nil {
		if errors.Is(err, my_err.ErrWrongCommunity) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &GetHeaderResponse{
		Head: &Header{
			AuthorID:    res.AuthorID,
			CommunityID: res.CommunityID,
			Author:      res.Author,
			Avatar:      string(res.Avatar),
		},
	}
	return resp, nil
}
