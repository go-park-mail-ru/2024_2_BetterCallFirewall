package post_api

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type PostService interface {
	GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error)
}

type Adapter struct {
	UnimplementedPostServiceServer
	service PostService
}

func New(s PostService) *Adapter {
	return &Adapter{service: s}
}

func (a *Adapter) GetAuthorsPosts(ctx context.Context, req *Request) (*Response, error) {
	request := &models.Header{
		AuthorID:    req.Head.AuthorID,
		Author:      req.Head.Author,
		CommunityID: req.Head.CommunityID,
		Avatar:      models.Picture(req.Head.Avatar),
	}

	res, err := a.service.GetAuthorsPosts(ctx, request)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Posts: make([]*Post, 0, len(res)),
	}
	for _, post := range res {
		resp.Posts = append(resp.Posts, &Post{
			ID: post.ID,
			Head: &Header{
				AuthorID:    post.Header.AuthorID,
				CommunityID: post.Header.CommunityID,
				Author:      post.Header.Author,
				Avatar:      string(post.Header.Avatar),
			},
			PostContent: &Content{
				Text:      post.PostContent.Text,
				File:      string(post.PostContent.File),
				CreatedAt: post.PostContent.CreatedAt.Unix(),
				UpdatedAt: post.PostContent.UpdatedAt.Unix(),
			},
		})
	}

	return resp, nil
}
