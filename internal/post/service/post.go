package service

import (
	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type DB interface {
	GetAll() []*models.Post
}

type PostServiceImpl struct {
	db DB
}

func NewPostServiceImpl(db DB) *PostServiceImpl {
	return &PostServiceImpl{
		db: db,
	}
}

func (s *PostServiceImpl) GetAll() []*models.Post {
	posts := s.db.GetAll()
	return posts
}
