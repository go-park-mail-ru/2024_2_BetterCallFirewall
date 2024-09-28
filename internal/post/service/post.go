package service

import (
	"errors"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type DB interface {
	GetAll() ([]*models.Post, error)
}

type PostServiceImpl struct {
	db DB
}

func NewPostServiceImpl(db DB) *PostServiceImpl {
	return &PostServiceImpl{
		db: db,
	}
}

func (s *PostServiceImpl) GetAll() ([]*models.Post, error) {
	posts, err := s.db.GetAll()
	if errors.Is(err, myErr.ErrPostEnd) {
		return nil, fmt.Errorf("get all posts: %w", err)
	}

	return posts, nil
}
