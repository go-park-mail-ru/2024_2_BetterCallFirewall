package service

import (
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
)

type ProfileUsecase struct {
	repo profile.Repository
}

func NewProfileUsecase(repo profile.Repository) *ProfileUsecase {
	return &ProfileUsecase{repo: repo}
}

func (p ProfileUsecase) GetProfileById(u uint32) (*models.Profile, error) {
	profile, err := p.repo.GetProfileById(u)
	if err != nil {
		return nil, fmt.Errorf("get profile by id usecase: %w", err)
	}

	panic("implement me")
}

func (p ProfileUsecase) GetAll(self uint32) ([]*models.Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileUsecase) CreateProfile(m models.Profile) (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileUsecase) UpdateProfile(m *models.Profile) error {
	//TODO implement me
	panic("implement me")
}

func (p ProfileUsecase) DeleteProfile(u uint32) error {
	//TODO implement me
	panic("implement me")
}
