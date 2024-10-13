package profile

import (
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ProfileUsecase interface {
	GetProfileById(uint32) (*models.FullProfile, error)
	GetAll(self uint32) ([]*models.ShortProfile, error)
	UpdateProfile(uint32, *models.FullProfile) error
	DeleteProfile(uint32) error
}
