package profile

import (
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Repository interface {
	GetProfileById(uint32) (*models.FullProfile, error)
	GetAll(self uint32) ([]*models.ShortProfile, error)
	UpdateProfile(*models.FullProfile) error
	DeleteProfile(uint32) error

	AddFriendsReq(reciever uint32, sender uint32) error
	AcceptFriendsReq(who uint32, whose uint32) error
	MoveToSubs(who uint32, whom uint32) error
	RemoveSub(who uint32, whom uint32) error
	CheckStatus(uint32, uint32) (int, error)
	GetAllFriends(uint32) ([]*models.ShortProfile, error)
}

type PostGetter interface {
	GetAuthorsPosts(uint32) ([]*models.Post, error)
}