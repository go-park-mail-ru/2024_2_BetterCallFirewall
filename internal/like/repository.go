package like

import (
	"context"
)

type Repository interface {
	SetLikeToPost(ctx context.Context, postID uint32, userID uint32) error
	SetLikeToComment(ctx context.Context, commentID uint32, userID uint32) error
	SetLikeToFile(ctx context.Context, fileID uint32, userID uint32) error
	DeleteLikeFromPost(ctx context.Context, postID uint32, userID uint32) error
	DeleteLikeFromComment(ctx context.Context, commentID uint32, userID uint32) error
	DeleteLikeFromFile(ctx context.Context, fileID uint32, userID uint32) error
	GetLikesOnPost(ctx context.Context, postID uint32) (uint32, error)
}
