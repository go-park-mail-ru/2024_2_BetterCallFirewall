package auth

import (
	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/auth_api"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

func NewSearchRequest(userID uint32) *auth_api.CreateRequest {
	return &auth_api.CreateRequest{
		UserID: userID,
	}
}

func UnmarshalCreateResponse(response *auth_api.CreateResponse) *models.Session {
	return &models.Session{
		ID:        response.Sess.ID,
		UserID:    response.Sess.UserID,
		CreatedAt: response.Sess.CreatedAt,
	}
}

func NewCheckRequest(cookie string) *auth_api.CheckRequest {
	return &auth_api.CheckRequest{
		Cookie: cookie,
	}
}

func UnmarshalCheckResponse(response *auth_api.CheckResponse) *models.Session {
	return &models.Session{
		ID:        response.Sess.ID,
		UserID:    response.Sess.UserID,
		CreatedAt: response.Sess.CreatedAt,
	}
}

func NewDestroyRequest(session *models.Session) *auth_api.DestroyRequest {
	return &auth_api.DestroyRequest{
		Sess: &auth_api.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			CreatedAt: session.CreatedAt,
		},
	}
}
