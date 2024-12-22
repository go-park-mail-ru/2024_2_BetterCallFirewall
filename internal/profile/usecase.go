package profile

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ProfileUsecase interface {
	GetProfileById(context.Context, uint32) (*models.FullProfile, error)
	GetAll(ctx context.Context, self uint32, lastId uint32) ([]*models.ShortProfile, error)
	UpdateProfile(context.Context, *models.FullProfile) error
	DeleteProfile(uint32) error
	Search(ctx context.Context, subStr string, lastId uint32) ([]*models.ShortProfile, error)
	ChangePassword(ctx context.Context, userID uint32, oldPassword, newPassword string) error

	SendFriendReq(receiver uint32, sender uint32) error
	AcceptFriendReq(who uint32, whose uint32) error
	RemoveFromFriends(who uint32, whose uint32) error
	Unsubscribe(who uint32, whose uint32) error
	GetAllFriends(ctx context.Context, id uint32, lastId uint32) ([]*models.ShortProfile, error)
	GetAllSubs(ctx context.Context, id uint32, lastId uint32) ([]*models.ShortProfile, error)
	GetAllSubscriptions(ctx context.Context, id uint32, lastId uint32) ([]*models.ShortProfile, error)
	GetHeader(ctx context.Context, userID uint32) (*models.Header, error)

	GetCommunitySubs(ctx context.Context, communityID, lastID uint32) ([]*models.ShortProfile, error)
}
