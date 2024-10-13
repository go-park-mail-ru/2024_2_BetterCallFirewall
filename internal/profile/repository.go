package profile

import (
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Repository interface {
	GetProfileById(uint32) (*models.FullProfile, error)
	GetAll(self uint32) ([]*models.ShortProfile, error)
	UpdateProfile(*models.FullProfile) error
	DeleteProfile(uint32) error
}

type PostGetter interface {
	GetAuthorsPosts(uint32) ([]*models.Post, error)
}
