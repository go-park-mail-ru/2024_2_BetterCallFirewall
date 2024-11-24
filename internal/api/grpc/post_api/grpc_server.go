package post_api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type PostService interface {
	GetAuthorsPosts(ctx context.Context, header *models.Header, userID uint32) ([]*models.Post, error)
}

type Adapter struct {
	UnimplementedPostServiceServer
	service PostService
}

func New(s PostService) *Adapter {
	return &Adapter{service: s}
}

func (a *Adapter) GetAuthorsPosts(ctx context.Context, req *Request) (*Response, error) {
	header := &models.Header{
		AuthorID:    req.Head.AuthorID,
		Author:      req.Head.Author,
		CommunityID: req.Head.CommunityID,
		Avatar:      models.Picture(req.Head.Avatar),
	}
	userID := req.UserID

	res, err := a.service.GetAuthorsPosts(ctx, header, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
