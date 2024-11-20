package service

import (
	"context"
)

type repoHelper interface {
	CheckAccess(ctx context.Context, communityID, userID uint32) bool
}

type ServiceHelper struct {
	repo repoHelper
}

func NewServiceHelper(repo repoHelper) *ServiceHelper {
	return &ServiceHelper{
		repo: repo,
	}
}

func (s *ServiceHelper) CheckAccess(ctx context.Context, communityID, userID uint32) bool {
	return s.repo.CheckAccess(ctx, communityID, userID)
}