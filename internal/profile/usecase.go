package profile

import (
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ProfileUsecase interface {
	GetProfileById(uint32) (*models.FullProfile, error)
	GetAll(self uint32) ([]*models.ShortProfile, error)
	UpdateProfile(uint32, *models.FullProfile) error
	DeleteProfile(uint32) error

	SendFriendReq(reciever uint32, sender uint32) error
	AcceptFriendReq(who uint32, whose uint32) error
	RemoveFromFriends(who uint32, whose uint32) error
	GetAllFriends(self uint32) ([]*models.ShortProfile, error)
}