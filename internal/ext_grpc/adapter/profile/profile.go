package profile

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/profile_api"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/port/profile"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type GrpcSender struct {
	client profile_api.ProfileServiceClient
}

func New(conn grpc.ClientConnInterface) *GrpcSender {
	client := profile_api.NewProfileServiceClient(conn)

	return &GrpcSender{
		client: client,
	}
}

func (g *GrpcSender) GetHeader(ctx context.Context, userID uint32) (*models.Header, error) {
	req := profile.NewGetHeaderRequest(userID)
	resp, err := g.client.GetHeader(ctx, req)
	if err != nil {
		return nil, err
	}

	res := profile.UnmarshallHeaderResponse(resp)
	return res, nil
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

func (g *GrpcSender) Create(ctx context.Context, user *models.User) (uint32, error) {
	req := profile.NewCreateRequest(user)
	resp, err := g.client.Create(ctx, req)
	if err != nil {
		return 0, err
	}

	res := profile.UnmarshallCreateResponse(resp)
	return res, nil
}

func (g *GrpcSender) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	req := profile.NewGetUserByEmailRequest(email)
	resp, err := g.client.GetUserByEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	res := profile.UnmarshallGetUserByEmailRequest(resp)
	return res, nil
}

func GetProfileProvider(host, port string) (grpc.ClientConnInterface, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return conn, nil
}
