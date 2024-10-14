package service

import (
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
)

type ProfileUsecase struct {
	repo        profile.Repository
	postManager profile.PostGetter
}

func NewProfileUsecase(repo profile.Repository) *ProfileUsecase {
	return &ProfileUsecase{repo: repo}
}

func (p ProfileUsecase) GetProfileById(u uint32) (*models.FullProfile, error) {
	profile, err := p.repo.GetProfileById(u)
	if err != nil {
		return nil, fmt.Errorf("get profile by id usecase: %w", err)
	}

	posts, err := p.postManager.GetAuthorsPosts(u)
	if err != nil {
		return nil, fmt.Errorf("get authors posts usecase: %w", err)
	}

	profile.Posts = posts
	return profile, nil
}

func (p ProfileUsecase) GetAll(self uint32) ([]*models.ShortProfile, error) {
	profiles, err := p.repo.GetAll(self)
	if err != nil {
		return nil, fmt.Errorf("get all profiles usecase: %w", err)
	}
	return profiles, nil
}

func validateOwner(ownerId uint32, profile *models.FullProfile) bool {
	return ownerId == profile.ID
}

func (p ProfileUsecase) UpdateProfile(owner uint32, newProfile *models.FullProfile) error {
	if !validateOwner(owner, newProfile) {
		return myErr.ErrWrongOwner
	}
	err := p.repo.UpdateProfile(newProfile)
	if err != nil {
		return fmt.Errorf("update profile usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecase) DeleteProfile(u uint32) error {
	err := p.repo.DeleteProfile(u)
	if err != nil {
		return fmt.Errorf("delete profile usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecase) SendFriendReq(reciever uint32, sender uint32) error {
	if reciever == sender {
		return myErr.ErrSameUser
	}

	panic("implement me")
}

func (p ProfileUsecase) AcceptFriendReq(who uint32, whose uint32) error {
	//TODO implement me
	panic("implement me")
}

func (p ProfileUsecase) RemoveFromFriends(who uint32, whose uint32) error {
	//TODO implement me
	panic("implement me")
}

func (p ProfileUsecase) GetAllFriends(self uint32) ([]*models.ShortProfile, error) {
	//TODO implement me
	panic("implement me")
}
