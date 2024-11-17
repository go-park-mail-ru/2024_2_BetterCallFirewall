package profile

import (
	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/profile_api"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

func NewGetHeaderRequest(userID uint32) *profile_api.HeaderRequest {
	return &profile_api.HeaderRequest{
		UserID: userID,
	}
}

func UnmarshallHeaderResponse(header *profile_api.HeaderResponse) *models.Header {
	return &models.Header{
		AuthorID:    header.Head.AuthorID,
		CommunityID: header.Head.CommunityID,
		Avatar:      models.Picture(header.Head.Avatar),
		Author:      header.Head.Author,
	}
}

func NewGetFriendsIDRequest(userID uint32) *profile_api.FriendsRequest {
	return &profile_api.FriendsRequest{
		UserID: userID,
	}
}

func UnmarshallGetFriendsIDResponse(users *profile_api.FriendsResponse) []uint32 {
	res := make([]uint32, 0, len(users.UserID))
	for _, id := range users.UserID {
		res = append(res, id)
	}

	return res
}
