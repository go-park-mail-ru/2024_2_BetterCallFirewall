package community

import (
	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/community_api"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

func NewRequest(communityID, userID uint32) *community_api.CheckAccessRequest {
	return &community_api.CheckAccessRequest{
		UserID:      userID,
		CommunityID: communityID,
	}
}

func UnmarshallResponse(response *community_api.CheckAccessResponse) bool {
	return response.Access
}

func NewHeaderRequest(communityID uint32) *community_api.GetHeaderRequest {
	return &community_api.GetHeaderRequest{
		CommunityID: communityID,
	}
}

func UnmarshallHeaderResponse(response *community_api.GetHeaderResponse) *models.Header {
	return &models.Header{
		AuthorID:    response.Head.AuthorID,
		CommunityID: response.Head.CommunityID,
		Author:      response.Head.Author,
		Avatar:      models.Picture(response.Head.Avatar),
	}
}
