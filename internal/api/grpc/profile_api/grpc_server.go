package profile_api

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type profileService interface {
	GetHeader(ctx context.Context, userID uint32) (models.Header, error)
	GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error)
}

type Adapter struct {
	UnimplementedProfileServiceServer
	service profileService
}

func New(s profileService) *Adapter {
	return &Adapter{
		service: s,
	}
}

func (a *Adapter) GetHeader(ctx context.Context, req *HeaderRequest) (*HeaderResponse, error) {
	userID := req.UserID
	res, err := a.service.GetHeader(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &HeaderResponse{
		Head: &Header{
			AuthorID:    res.AuthorID,
			CommunityID: res.CommunityID,
			Author:      res.Author,
			Avatar:      string(res.Avatar),
		},
	}

	return resp, nil
}

func (a *Adapter) GetFriendsID(ctx context.Context, req *FriendsRequest) (*FriendsResponse, error) {
	userID := req.UserID
	res, err := a.service.GetFriendsID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &FriendsResponse{
		UserID: make([]uint32, 0, len(res)),
	}

	for _, id := range res {
		resp.UserID = append(resp.UserID, id)
	}

	return resp, nil
}
