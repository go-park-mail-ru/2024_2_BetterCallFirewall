package usecase

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/like"
)

type Usecase struct {
	Repo like.Repository
}

func NewLikeUsecase(repo like.Repository) *Usecase {
	return &Usecase{
		Repo: repo,
	}
}

func (u Usecase) SetLikeToPost(ctx context.Context, postID uint32, userID uint32) error {
	err := u.Repo.SetLikeToPost(ctx, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) SetLikeToComment(ctx context.Context, commentID uint32, userID uint32) error {
	err := u.Repo.SetLikeToComment(ctx, commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) SetLikeToFile(ctx context.Context, fileID uint32, userID uint32) error {
	err := u.Repo.SetLikeToFile(ctx, fileID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) DeleteLikeFromPost(ctx context.Context, postID uint32, userID uint32) error {
	err := u.Repo.DeleteLikeFromPost(ctx, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) DeleteLikeFromComment(ctx context.Context, commentID uint32, userID uint32) error {
	err := u.Repo.DeleteLikeFromComment(ctx, commentID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) DeleteLikeFromFile(ctx context.Context, fileID uint32, userID uint32) error {
	err := u.Repo.DeleteLikeFromFile(ctx, fileID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) GetLikesOnPost(ctx context.Context, postID uint32) (uint32, error) {
	likes, err := u.Repo.GetLikesOnPost(ctx, postID)
	if err != nil {
		return 0, err
	}
	return likes, nil
}
