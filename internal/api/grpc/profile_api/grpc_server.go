package profile_api

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type profileService interface {
	GetHeader(ctx context.Context, userID uint32) (*models.Header, error)
	GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error)
	Create(ctx context.Context, user *models.User) (uint32, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
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

func (a *Adapter) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	user := &models.User{
		ID:        req.User.ID,
		Email:     req.User.Email,
		Password:  req.User.Password,
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Avatar:    models.Picture(req.User.Avatar),
	}

	res, err := a.service.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	resp := &CreateResponse{
		ID: res,
	}

	return resp, nil
}

func (a *Adapter) GetUserByEmail(ctx context.Context, req *GetByEmailRequest) (*GetByEmailResponse, error) {
	email := req.Email

	user, err := a.service.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	resp := &GetByEmailResponse{
		User: &User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     email,
			Password:  user.Password,
			Avatar:    string(user.Avatar),
		},
	}

	return resp, nil
}
