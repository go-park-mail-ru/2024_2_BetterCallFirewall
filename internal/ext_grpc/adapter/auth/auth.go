package auth

import (
	"context"

	"google.golang.org/grpc"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/auth_api"
	"github.com/2024_2_BetterCallFirewall/internal/ext_grpc/port/auth"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type GrpcSender struct {
	client auth_api.AuthServiceClient
}

func New(conn grpc.ClientConnInterface) *GrpcSender {
	client := auth_api.NewAuthServiceClient(conn)

	return &GrpcSender{client: client}
}

func (s *GrpcSender) Create(userID uint32) (*models.Session, error) {
	req := auth.NewSearchRequest(userID)
	resp, err := s.client.Create(context.Background(), req)
	if err != nil {
		return nil, err
	}

	res := auth.UnmarshalCreateResponse(resp)
	return res, nil
}

func (s *GrpcSender) Check(cookie string) (*models.Session, error) {
	req := auth.NewCheckRequest(cookie)
	resp, err := s.client.Check(context.Background(), req)
	if err != nil {
		return nil, err
	}

	res := auth.UnmarshalCheckResponse(resp)
	return res, nil
}

func (s *GrpcSender) Destroy(session *models.Session) error {
	req := auth.NewDestroyRequest(session)
	_, err := s.client.Destroy(context.Background(), req)

	return err
}
