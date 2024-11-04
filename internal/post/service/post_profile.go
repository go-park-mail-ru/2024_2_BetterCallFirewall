package service

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type FileService interface {
	GetPostPicture(ctx context.Context, postID uint32) *models.Picture
}

type PostProfileDB interface {
	GetAuthorPosts(ctx context.Context, header *models.Header) ([]*models.Post, error)
}

type PostProfileImpl struct {
	fileService FileService
	db          PostProfileDB
}

func NewPostProfileImpl(fileService FileService, db PostProfileDB) *PostProfileImpl {
	return &PostProfileImpl{
		fileService: fileService,
		db:          db,
	}
}
func (p *PostProfileImpl) GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error) {
	posts, err := p.db.GetAuthorPosts(ctx, header)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		post.PostContent.File = p.fileService.GetPostPicture(ctx, post.ID)
		posts[i] = post
	}

	return posts, nil
}
