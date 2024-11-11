package service

import (
	"context"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Repo interface {
	GetBatch(ctx context.Context, lastID uint32) ([]*models.Community, error)
	GetOne(ctx context.Context, id uint32) (*models.Community, error)
	Create(ctx context.Context, community *models.Community) (uint32, error)
	Update(ctx context.Context, community *models.Community) error
	Delete(ctx context.Context, id uint32) error
	CheckAccess(ctx context.Context, communityID, userID uint32) error
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, lastID uint32) ([]*models.Community, error) {
	coms, err := s.repo.GetBatch(ctx, lastID)
	if err != nil {
		return nil, fmt.Errorf("get community list: %w", err)
	}

	return coms, nil
}

func (s *Service) GetOne(ctx context.Context, id uint32) (*models.Community, error) {
	com, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get community: %w", err)
	}

	return com, nil
}

func (s *Service) Create(ctx context.Context, community *models.Community) error {
	id, err := s.repo.Create(ctx, community)
	if err != nil {
		return fmt.Errorf("create community: %w", err)
	}
	community.ID = id

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
	err := s.repo.CheckAccess(ctx, communityID, userID)
	if err != nil {
		return false
	}

	return true
}
