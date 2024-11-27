package profile

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Repository interface {
	GetProfileById(context.Context, uint32) (*models.FullProfile, error)
	GetStatus(context.Context, uint32, uint32) (int, error)
	GetAll(ctx context.Context, self uint32, lastId uint32) ([]*models.ShortProfile, error)
	UpdateProfile(context.Context, *models.FullProfile) error
	UpdateWithAvatar(context.Context, *models.FullProfile) error
	DeleteProfile(uint32) error
	Search(ctx context.Context, subStr string, lastId uint32) ([]*models.ShortProfile, error)

	CheckFriendship(context.Context, uint32, uint32) (bool, error)
	AddFriendsReq(receiver uint32, sender uint32) error
	AcceptFriendsReq(who uint32, whose uint32) error
	MoveToSubs(who uint32, whom uint32) error
	RemoveSub(who uint32, whom uint32) error
	GetAllFriends(ctx context.Context, u uint32, lastId uint32) ([]*models.ShortProfile, error)
	GetAllSubs(ctx context.Context, u uint32, lastId uint32) ([]*models.ShortProfile, error)
	GetAllSubscriptions(context.Context, uint32, uint32) ([]*models.ShortProfile, error)

	GetSubscriptionsID(context.Context, uint32) ([]uint32, error)
	GetSubscribersID(context.Context, uint32) ([]uint32, error)
	GetStatuses(context.Context, uint32) ([]uint32, []uint32, []uint32, error)
	GetHeader(context.Context, uint32) (*models.Header, error)

	GetCommunitySubs(ctx context.Context, communityID uint32, lastInsertId uint32) ([]*models.ShortProfile, error)
}

type PostGetter interface {
	GetAuthorsPosts(ctx context.Context, header *models.Header, userID uint32) ([]*models.Post, error)
}
