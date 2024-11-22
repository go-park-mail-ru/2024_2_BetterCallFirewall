package profile_api

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
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
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &FriendsResponse{
		UserID: make([]uint32, 0, len(res)),
	}

	resp.UserID = append(resp.UserID, res...)

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
	if errors.Is(err, my_err.ErrUserAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &CreateResponse{
		ID: res,
	}

	return resp, nil
}

func (a *Adapter) GetUserByEmail(ctx context.Context, req *GetByEmailRequest) (*GetByEmailResponse, error) {
	email := req.Email

	user, err := a.service.GetByEmail(ctx, email)

	if errors.Is(err, my_err.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
