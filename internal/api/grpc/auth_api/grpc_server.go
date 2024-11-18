package auth_api

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

type Adapter struct {
	UnimplementedAuthServiceServer
	authServer SessionManager
}

func New(authServer SessionManager) *Adapter {
	return &Adapter{
		authServer: authServer,
	}
}

func (a *Adapter) Check(ctx context.Context, reqGRPC *CheckRequest) (*CheckResponse, error) {
	req := reqGRPC.Cookie
	sess, err := a.authServer.Check(req)
	if err != nil {
		return nil, err
	}

	res := &CheckResponse{
		Sess: &Session{
			ID:        sess.ID,
			UserID:    sess.UserID,
			CreatedAt: sess.CreatedAt,
		},
	}

	return res, nil
}

func (a *Adapter) Create(ctx context.Context, reqGRPC *CreateRequest) (*CreateResponse, error) {
	req := reqGRPC.UserID
	sess, err := a.authServer.Create(req)
	if err != nil {
		return nil, err
	}
	res := &CreateResponse{
		Sess: &Session{
			ID:        sess.ID,
			UserID:    sess.UserID,
			CreatedAt: sess.CreatedAt,
		},
	}

	return res, nil
}

func (a *Adapter) Destroy(ctx context.Context, reqGRPC *DestroyRequest) (*EmptyResponse, error) {
	req := &models.Session{
		ID:        reqGRPC.Sess.ID,
		UserID:    reqGRPC.Sess.UserID,
		CreatedAt: reqGRPC.Sess.CreatedAt,
	}
	err := a.authServer.Destroy(req)

	return nil, err
}
