package service

import (
	"context"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type PostProfileDB interface {
	GetAuthorPosts(ctx context.Context, header *models.Header) ([]*models.Post, error)
	GetLikesOnPost(ctx context.Context, postID uint32) (uint32, error)
	CheckLikes(ctx context.Context, postID, userID uint32) (bool, error)
}

type PostProfileImpl struct {
	db PostProfileDB
}

func NewPostProfileImpl(db PostProfileDB) *PostProfileImpl {
	return &PostProfileImpl{
		db: db,
	}
}

func (p *PostProfileImpl) GetAuthorsPosts(ctx context.Context, header *models.Header, userId uint32) ([]*models.Post, error) {
	posts, err := p.db.GetAuthorPosts(ctx, header)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		likes, err := p.db.GetLikesOnPost(ctx, post.ID)
		if err != nil {
			return nil, fmt.Errorf("get likes: %w", err)
		}
		posts[i].LikesCount = likes

		liked, err := p.db.CheckLikes(ctx, post.ID, userId)
		if err != nil {
			return nil, fmt.Errorf("check likes: %w", err)
		}
		posts[i].IsLiked = liked
	}

	return posts, nil
}
