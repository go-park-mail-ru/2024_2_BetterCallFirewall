package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

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
	profile, err := p.repo.GetProfileById(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("get profile by id usecase: %w", err)
	}

	sess, err := models.SessionFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get session usecase: %w", myErr.ErrSessionNotFound)
	}

	self := sess.UserID
	if u == self {
		profile.IsAuthor = true
	} else {
		status, err := p.repo.GetStatus(ctx, self, u)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("get status usecase: %w", err)
			}
		} else {
			profile.IsFriend = status == 0
			profile.IsSubscription = status == 1
			profile.IsSubscriber = status == -1

		}
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

func (p ProfileUsecaseImplementation) GetAll(ctx context.Context, self uint32, lastId uint32) ([]*models.ShortProfile, error) {
	profiles, err := p.repo.GetAll(ctx, self, lastId)
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

func (p ProfileUsecaseImplementation) SendFriendReq(receiver uint32, sender uint32) error {
	if receiver == sender {
		return myErr.ErrSameUser
	}
	err := p.repo.AddFriendsReq(receiver, sender)
	if err != nil {
		return fmt.Errorf("add friend req usecase: %w", err)
	}

	return nil
}

func (p ProfileUsecaseImplementation) AcceptFriendReq(who uint32, whose uint32) error {
	if who == whose {
		return myErr.ErrSameUser
	}
	err := p.repo.AcceptFriendsReq(who, whose)
	if err != nil {
		return fmt.Errorf("accept friend req usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecaseImplementation) RemoveFromFriends(who uint32, whom uint32) error {
	if who == whom {
		return myErr.ErrSameUser
	}
	err := p.repo.RemoveSub(who, whom)
	if err != nil {
		return fmt.Errorf("remove sub usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecaseImplementation) Unsubscribe(who uint32, whom uint32) error {
	if who == whom {
		return myErr.ErrSameUser
	}

	err := p.repo.MoveToSubs(who, whom)
	if err != nil {
		return fmt.Errorf("unsub usecase: %w", err)
	}
	return nil
}

func (p ProfileUsecaseImplementation) setStatuses(ctx context.Context, profiles []*models.ShortProfile) error {
	sess, err := models.SessionFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get self session usecase: %w", myErr.ErrSessionNotFound)
	}
	selfId := sess.UserID

	friends, err := p.repo.GetFriendsID(ctx, selfId)
	if err != nil {
		return fmt.Errorf("get self friends usecase: %w", err)
	}

	subs, err := p.repo.GetSubscribersID(ctx, selfId)
	if err != nil {
		return fmt.Errorf("get self friends usecase: %w", err)
	}

	subscriptions, err := p.repo.GetSubscriptionsID(ctx, selfId)
	if err != nil {
		return fmt.Errorf("get self friends usecase: %w", err)
	}

	for _, profile := range profiles {
		profile.IsFriend = slices.Contains(friends, profile.ID)
		profile.IsSubscriber = slices.Contains(subs, profile.ID)
		profile.IsSubscription = slices.Contains(subscriptions, profile.ID)
		profile.IsAuthor = profile.ID == selfId
	}
	return nil
}

func (p ProfileUsecaseImplementation) GetAllFriends(ctx context.Context, id uint32, lastId uint32) ([]*models.ShortProfile, error) {
	res, err := p.repo.GetAllFriends(ctx, id, lastId)
	if err != nil {
		return nil, fmt.Errorf("get all friends usecase: %w", err)
	}

	err = p.setStatuses(ctx, res)
	if err != nil {
		return nil, fmt.Errorf("get all friends usecase: %w", err)
	}

	return res, nil
}

func (p ProfileUsecaseImplementation) GetAllSubs(ctx context.Context, id uint32, lastId uint32) ([]*models.ShortProfile, error) {
	res, err := p.repo.GetAllSubs(ctx, id, lastId)
	if err != nil {
		return nil, fmt.Errorf("get all subs usecase: %w", err)
	}

	err = p.setStatuses(ctx, res)
	if err != nil {
		return nil, fmt.Errorf("get all friends usecase: %w", err)
	}

	return res, nil
}

func (p ProfileUsecaseImplementation) GetAllSubscriptions(ctx context.Context, id uint32, lastId uint32) ([]*models.ShortProfile, error) {
	res, err := p.repo.GetAllSubscriptions(ctx, id, lastId)
	if err != nil {
		return nil, fmt.Errorf("get all subscriptions usecase: %w", err)
	}

	err = p.setStatuses(ctx, res)
	if err != nil {
		return nil, fmt.Errorf("get all friends usecase: %w", err)
	}

	return res, nil
}

func (p ProfileUsecaseImplementation) GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error) {
	res, err := p.repo.GetFriendsID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get friends id usecase: %w", err)
	}
	return res, nil
}

func (p ProfileUsecaseImplementation) GetHeader(ctx context.Context, userID uint32) (models.Header, error) {
	header, err := p.repo.GetHeader(ctx, userID)
	if err != nil {
		return models.Header{}, fmt.Errorf("get header usecase: %w", err)
	}
	return *header, nil
}
