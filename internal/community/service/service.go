package service

import (
	"context"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type Repo interface {
	GetBatch(ctx context.Context, lastID uint32) ([]*models.CommunityCard, error)
	GetOne(ctx context.Context, id uint32) (*models.Community, error)
	Create(ctx context.Context, community *models.Community, author uint32) (uint32, error)
	Update(ctx context.Context, community *models.Community) error
	Delete(ctx context.Context, id uint32) error
	CheckAccess(ctx context.Context, communityID, userID uint32) bool
	JoinCommunity(ctx context.Context, communityId, author uint32) error
	LeaveCommunity(ctx context.Context, communityId, author uint32) error
	NewAdmin(ctx context.Context, communityId uint32, author uint32) error
	Search(ctx context.Context, query string, lastID uint32) ([]*models.CommunityCard, error)
	IsFollowed(ctx context.Context, communityId, userID uint32) (bool, error)
}

type Service struct {
	repo Repo
}

func NewCommunityService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, userID, lastID uint32) ([]*models.CommunityCard, error) {
	coms, err := s.repo.GetBatch(ctx, lastID)
	if err != nil {
		return nil, fmt.Errorf("get community list: %w", err)
	}

	for i, com := range coms {
		follow, err := s.repo.IsFollowed(ctx, com.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("get community list: %w", err)
		}
		coms[i].IsFollowed = follow
	}

	return coms, nil
}

func (s *Service) GetOne(ctx context.Context, id, userID uint32) (*models.Community, error) {
	com, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get community: %w", err)
	}
	com.IsAdmin = s.CheckAccess(ctx, id, userID)

	follow, err := s.repo.IsFollowed(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("get community list: %w", err)
	}
	com.IsFollowed = follow

	return com, nil
}

func (s *Service) Create(ctx context.Context, community *models.Community, authorID uint32) error {
	id, err := s.repo.Create(ctx, community, authorID)
	if err != nil {
		return fmt.Errorf("create community: %w", err)
	}
	community.ID = id

	s.AddAdmin(ctx, id, authorID)

	return nil
}

func (s *Service) Update(ctx context.Context, id uint32, community *models.Community) error {
	community.ID = id
	err := s.repo.Update(ctx, community)
	if err != nil {
		return fmt.Errorf("update community: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uint32) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete community: %w", err)
	}

	return nil
}

func (s *Service) CheckAccess(ctx context.Context, communityID, userID uint32) bool {
	return s.repo.CheckAccess(ctx, communityID, userID)
}

func (s *Service) JoinCommunity(ctx context.Context, communityId, author uint32) error {
	err := s.repo.JoinCommunity(ctx, communityId, author)
	if err != nil {
		return fmt.Errorf("join community: %w", err)
	}

	return nil
}

func (s *Service) LeaveCommunity(ctx context.Context, communityId, author uint32) error {
	err := s.repo.LeaveCommunity(ctx, communityId, author)
	if err != nil {
		return fmt.Errorf("leave community: %w", err)
	}

	return nil
}

func (s *Service) AddAdmin(ctx context.Context, communityId, author uint32) error {
	err := s.repo.NewAdmin(ctx, communityId, author)
	if err != nil {
		return fmt.Errorf("add admin: %w", err)
	}

	return nil
}

func (s *Service) Search(ctx context.Context, query string, userID, lastID uint32) ([]*models.CommunityCard, error) {
	cards, err := s.repo.Search(ctx, query, lastID)
	if err != nil {
		return nil, fmt.Errorf("search community: %w", err)
	}

	for i, card := range cards {
		follow, err := s.repo.IsFollowed(ctx, card.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("get community list: %w", err)
		}
		cards[i].IsFollowed = follow
	}

	return cards, nil
}
