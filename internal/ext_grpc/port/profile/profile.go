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
	res = append(res, users.UserID...)

	return res
}

func NewCreateRequest(user *models.User) *profile_api.CreateRequest {
	return &profile_api.CreateRequest{
		User: &profile_api.User{
			ID:        user.ID,
			Email:     user.Email,
			Password:  user.Password,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    string(user.Avatar),
		},
	}
}

func UnmarshallCreateResponse(user *profile_api.CreateResponse) uint32 {
	return uint32(user.ID)
}

func NewGetUserByEmailRequest(email string) *profile_api.GetByEmailRequest {
	return &profile_api.GetByEmailRequest{
		Email: email,
	}
}

func UnmarshallGetUserByEmailRequest(response *profile_api.GetByEmailResponse) *models.User {
	return &models.User{
		ID:        response.User.ID,
		Email:     response.User.Email,
		FirstName: response.User.FirstName,
		LastName:  response.User.LastName,
		Password:  response.User.Password,
		Avatar:    models.Picture(response.User.Avatar),
	}
}
