package service

import "github.com/2024_2_BetterCallFirewall/internal/profile/models"

type ProfileUsecase interface {
	GetProfileById(uint64 uint32) (*models.Profile, error)
	GetAll() ([]*models.Profile, error)
	CreateProfile(models.Profile) (uint32, error)
	UpdateProfile(*models.Profile) (bool, error)
	DeleteProfile(uint32) error
}
