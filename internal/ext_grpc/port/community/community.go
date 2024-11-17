package community

import (
	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/community_api"
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
