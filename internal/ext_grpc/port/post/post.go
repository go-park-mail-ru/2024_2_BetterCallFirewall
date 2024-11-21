package post

import (
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/post_api"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

func NewRequest(header *models.Header) *post_api.Request {
	return &post_api.Request{
		Head: &post_api.Header{
			AuthorID:    header.AuthorID,
			CommunityID: header.CommunityID,
			Avatar:      string(header.Avatar),
			Author:      header.Author,
		},
	}
}

func UnmarshalResponse(response *post_api.Response) []*models.Post {
	res := make([]*models.Post, 0, len(response.Posts))
	for _, post := range response.Posts {
		res = append(res, &models.Post{
			ID: post.ID,
			Header: models.Header{
				AuthorID:    post.Head.AuthorID,
				CommunityID: post.Head.CommunityID,
				Avatar:      models.Picture(post.Head.Avatar),
				Author:      post.Head.Author,
			},
			PostContent: models.Content{
				Text:      post.PostContent.Text,
				File:      models.Picture(post.PostContent.File),
				CreatedAt: time.Unix(post.PostContent.CreatedAt, 0),
				UpdatedAt: time.Unix(post.PostContent.UpdatedAt, 0),
			},
		})
	}

	return res
}
