package profile

import (
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ProfileUsecase interface {
	GetProfileById(uint32) (*models.Profile, error)
	GetAll(self uint32) ([]*models.Profile, error)
	CreateProfile(models.Profile) (uint32, error)
	UpdateProfile(*models.Profile) error
	DeleteProfile(uint32) error
}
