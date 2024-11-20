package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type LikeRepository struct {
	DB *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{
		DB: db,
	}
}

func (lr *LikeRepository) SetLikeToPost(ctx context.Context, postID uint32, userID uint32) error {
	res, err := lr.DB.ExecContext(ctx, AddLikeToPost, postID, userID)
	if num, err := res.RowsAffected(); err == nil && num == 0 {
		return myErr.ErrLikeAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}

func (lr *LikeRepository) SetLikeToComment(ctx context.Context, commentID uint32, userID uint32) error {
	res, err := lr.DB.ExecContext(ctx, AddLikeToComment, commentID, userID)
	if num, err := res.RowsAffected(); err == nil && num == 0 {
		return myErr.ErrLikeAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}

func (lr *LikeRepository) SetLikeToFile(ctx context.Context, fileID uint32, userID uint32) error {
	res, err := lr.DB.ExecContext(ctx, AddLikeToFile, fileID, userID)
	if num, err := res.RowsAffected(); err == nil && num == 0 {
		return myErr.ErrLikeAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}

func (lr *LikeRepository) DeleteLikeFromPost(ctx context.Context, postID uint32, userID uint32) error {
	_, err := lr.DB.ExecContext(ctx, DeleteLikeFromPost, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (lr *LikeRepository) DeleteLikeFromComment(ctx context.Context, userID uint32, commentID uint32) error {
	_, err := lr.DB.ExecContext(ctx, DeleteLikeFromComment, commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (lr *LikeRepository) DeleteLikeFromFile(ctx context.Context, fileID uint32, userID uint32) error {
	_, err := lr.DB.ExecContext(ctx, DeleteLikeFromFile, fileID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (lr *LikeRepository) GetLikesOnPost(ctx context.Context, postID uint32) (uint32, error) {
	var likes uint32
	err := lr.DB.QueryRowContext(ctx, GetLikesOnPost, postID).Scan(&likes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, myErr.ErrWrongPost
		}
		return 0, err
	}
	return likes, nil
}
