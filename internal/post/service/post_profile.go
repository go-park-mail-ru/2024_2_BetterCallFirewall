package service

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type PostProfileDB interface {
	GetAuthorPosts(ctx context.Context, header *models.Header) ([]*models.Post, error)
}

type PostProfileImpl struct {
	db PostProfileDB
}

func NewPostProfileImpl(db PostProfileDB) *PostProfileImpl {
	return &PostProfileImpl{
		db: db,
	}
}
func (p *PostProfileImpl) GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error) {
	posts, err := p.db.GetAuthorPosts(ctx, header)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
