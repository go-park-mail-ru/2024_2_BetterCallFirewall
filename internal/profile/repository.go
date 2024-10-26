package profile

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Repository interface {
	GetProfileById(context.Context, uint32) (*models.FullProfile, error)
	GetAll(ctx context.Context, self uint32) ([]*models.ShortProfile, error)
	UpdateProfile(*models.FullProfile) error
	DeleteProfile(uint32) error

	AddFriendsReq(reciever uint32, sender uint32) error
	AcceptFriendsReq(who uint32, whose uint32) error
	MoveToSubs(who uint32, whom uint32) error
	RemoveSub(who uint32, whom uint32) error
	GetAllFriends(context.Context, uint32) ([]*models.ShortProfile, error)
	GetAllSubs(context.Context, uint32) ([]*models.ShortProfile, error)
	GetAllSubscriptions(context.Context, uint32) ([]*models.ShortProfile, error)

	GetFriendsID(context.Context, uint32) ([]uint32, error)
	GetHeader(context.Context, uint32) (*models.Header, error)
}

type PostGetter interface {
	GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error)
}
