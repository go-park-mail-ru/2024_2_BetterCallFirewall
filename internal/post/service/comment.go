package service

import (
	"context"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

//go:generate mockgen -destination=comment_mock.go -source=$GOFILE -package=${GOPACKAGE}
type dbI interface {
	CreateComment(ctx context.Context, comment *models.Content, userID, postID uint32) (uint32, error)
	DeleteComment(ctx context.Context, commentID uint32) error
	UpdateComment(ctx context.Context, comment *models.Content, commentID uint32) error
	GetComments(ctx context.Context, postID, lastID uint32) ([]*models.Comment, error)
	GetCommentAuthor(ctx context.Context, commentID uint32) (uint32, error)
}

type profileRepoI interface {
	GetHeader(ctx context.Context, userID uint32) (*models.Header, error)
}

type CommentService struct {
	db          dbI
	profileRepo profileRepoI
}

func NewCommentService(db dbI, profileRepo profileRepoI) *CommentService {
	return &CommentService{
		db:          db,
		profileRepo: profileRepo,
	}
}

func (s *CommentService) Comment(
	ctx context.Context, userID, postID uint32, comment *models.Content,
) (*models.Comment, error) {
	id, err := s.db.CreateComment(ctx, comment, userID, postID)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	header, err := s.profileRepo.GetHeader(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get header: %w", err)
	}

	newComment := &models.Comment{
		ID:      id,
		Content: *comment,
		Header:  *header,
	}

	return newComment, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, commentID, userID uint32) error {
	authorID, err := s.db.GetCommentAuthor(ctx, commentID)
	if err != nil {
		return fmt.Errorf("get comment author: %w", err)
	}

	if authorID != userID {
		return my_err.ErrAccessDenied
	}

	err = s.db.DeleteComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	return nil
}

func (s *CommentService) EditComment(ctx context.Context, commentID, userID uint32, comment *models.Content) error {
	authorID, err := s.db.GetCommentAuthor(ctx, commentID)
	if err != nil {
		return fmt.Errorf("get comment author: %w", err)
	}

	if authorID != userID {
		return my_err.ErrAccessDenied
	}

	err = s.db.UpdateComment(ctx, comment, commentID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	return nil
}

func (s *CommentService) GetComments(ctx context.Context, postID, lastID uint32) ([]*models.Comment, error) {
	comments, err := s.db.GetComments(ctx, postID, lastID)
	if err != nil {
		return nil, fmt.Errorf("get comments: %w", err)
	}

	for i, c := range comments {
		header, err := s.profileRepo.GetHeader(ctx, c.Header.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get header %d: %w", i, err)
		}

		comments[i].Header = *header
	}

	return comments, nil
}
