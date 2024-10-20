package service

import (
	"context"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
)

type ProfileUsecaseImplementation struct {
	repo        profile.Repository
	postManager profile.PostGetter
}

func NewProfileUsecase(profileRepo profile.Repository, postRepo profile.PostGetter) *ProfileUsecaseImplementation {
	return &ProfileUsecaseImplementation{repo: profileRepo, postManager: postRepo}
}

func (p ProfileUsecaseImplementation) GetProfileById(ctx context.Context, u uint32) (*models.FullProfile, error) {
	profile, err := p.repo.GetProfileById(u, ctx)
	if err != nil {
		return nil, fmt.Errorf("get profile by id usecase: %w", err)
	}

	header := models.Header{
		AuthorID: profile.ID,
		Author:   profile.FirstName + " " + profile.LastName,
		Avatar:   profile.Avatar,
	}
	posts, err := p.postManager.GetAuthorsPosts(context.Background(), &header)
	if err != nil {
		return nil, fmt.Errorf("get authors posts usecase: %w", err)
	}

	profile.Posts = posts
	return profile, nil
}

func (p ProfileUsecaseImplementation) GetAll(ctx context.Context, self uint32) ([]*models.ShortProfile, error) {
	profiles, err := p.repo.GetAll(self, ctx)
	if err != nil {
		return nil, fmt.Errorf("get all profiles usecase: %w", err)
	}
	return profiles, nil
}

func validateOwner(ownerId uint32, profile *models.FullProfile) bool {
	return ownerId == profile.ID
}

func (p ProfileUsecaseImplementation) UpdateProfile(owner uint32, newProfile *models.FullProfile) error {
	if !validateOwner(owner, newProfile) {
		return myErr.ErrWrongOwner
	}
	err := p.repo.UpdateProfile(newProfile)
	if err != nil {
		return fmt.Errorf("update profile usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecaseImplementation) DeleteProfile(u uint32) error {
	err := p.repo.DeleteProfile(u)
	if err != nil {
		return fmt.Errorf("delete profile usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecaseImplementation) SendFriendReq(reciever uint32, sender uint32) error {
	if reciever == sender {
		return myErr.ErrSameUser
	}
	err := p.repo.AddFriendsReq(reciever, sender)
	if err != nil {
		return fmt.Errorf("add friend req usecase: %w", err)
	}

	return nil
}

func (p ProfileUsecaseImplementation) AcceptFriendReq(who uint32, whose uint32) error {
	err := p.repo.AcceptFriendsReq(who, whose)
	if err != nil {
		return fmt.Errorf("accept friend req usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecaseImplementation) RemoveFromFriends(who uint32, whom uint32) error {
	err := p.repo.RemoveSub(who, whom)
	if err != nil {
		return fmt.Errorf("remove sub usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecaseImplementation) Unsubscribe(who uint32, whom uint32) error {
	err := p.repo.MoveToSubs(who, whom)
	if err != nil {
		return fmt.Errorf("unsub usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecaseImplementation) GetAllFriends(ctx context.Context, self uint32) ([]*models.ShortProfile, error) {
	res, err := p.repo.GetAllFriends(self, ctx)
	if err != nil {
		return nil, fmt.Errorf("get all friends usecase: %w", err)
	}
	return res, nil
}

func (p ProfileUsecaseImplementation) GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error) {
	res, err := p.repo.GetFriendsID(userID, ctx)
	if err != nil {
		return nil, fmt.Errorf("get friends id usecase: %w", err)
	}
	return res, nil
}

func (p ProfileUsecaseImplementation) GetHeader(ctx context.Context, userID uint32) (models.Header, error) {
	header, err := p.repo.GetHeader(userID)
	if err != nil {
		return models.Header{}, fmt.Errorf("get header usecase: %w", err)
	}
	return *header, nil
}
