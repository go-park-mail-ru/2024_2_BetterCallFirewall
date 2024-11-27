package service

import (
	"context"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock_helper.go -source=$GOFILE -package=${GOPACKAGE}
type repository interface {
	Create(user *models.User, ctx context.Context) (uint32, error)
	GetByEmail(email string, ctx context.Context) (*models.User, error)
	GetFriendsID(context.Context, uint32) ([]uint32, error)
	GetHeader(context.Context, uint32) (*models.Header, error)
}

type ProfileHelper struct {
	repo repository
}

func NewProfileHelper(repo repository) *ProfileHelper {
	return &ProfileHelper{repo}
}

func (p ProfileHelper) Create(ctx context.Context, user *models.User) (uint32, error) {
	id, err := p.repo.Create(user, ctx)
	if err != nil {
		return 0, fmt.Errorf("get profile: %w", err)
	}

	return id, nil
}

func (p ProfileHelper) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := p.repo.GetByEmail(email, ctx)
	if err != nil {
		return nil, fmt.Errorf("get user by email usecase: %w", err)
	}

	return user, nil
}

func (p ProfileHelper) GetHeader(ctx context.Context, userID uint32) (*models.Header, error) {
	header, err := p.repo.GetHeader(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get header usecase: %w", err)
	}

	return header, nil
}

func (p ProfileHelper) GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error) {
	res, err := p.repo.GetFriendsID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get friends id usecase: %w", err)
	}

	return res, nil
}
