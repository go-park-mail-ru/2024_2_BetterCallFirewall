package service

import (
	"context"
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type DB interface {
	Create(ctx context.Context, post *models.Post) (uint32, error)
	Get(ctx context.Context, postID uint32) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, postID uint32) error
	GetPosts(ctx context.Context, lastID uint32) ([]*models.Post, error)
	GetFriendsPosts(ctx context.Context, friendsID []uint32, lastID uint32) ([]*models.Post, error)
	GetPostAuthor(ctx context.Context, postID uint32) (uint32, error)
}

type ProfileRepo interface {
	GetHeader(ctx context.Context, userID uint32) (models.Header, error)
	GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error)
}

type PostServiceImpl struct {
	db          DB
	profileRepo ProfileRepo
}

func NewPostServiceImpl(db DB, profileRepo ProfileRepo) *PostServiceImpl {
	return &PostServiceImpl{
		db:          db,
		profileRepo: profileRepo,
	}
}

func (s *PostServiceImpl) Create(ctx context.Context, post *models.Post) (uint32, error) {
	id, err := s.db.Create(ctx, post)
	if err != nil {
		return 0, fmt.Errorf("create post: %w", err)
	}

	return id, nil
}

func (s *PostServiceImpl) Get(ctx context.Context, postID uint32) (*models.Post, error) {
	post, err := s.db.Get(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}

	header, err := s.profileRepo.GetHeader(ctx, post.Header.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("get header:%w", err)
	}
	post.Header = header

	return post, nil
}

func (s *PostServiceImpl) Delete(ctx context.Context, postID uint32) error {
	err := s.db.Delete(ctx, postID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	return nil
}

func (s *PostServiceImpl) Update(ctx context.Context, post *models.Post) error {
	post.PostContent.UpdatedAt = time.Now()

	err := s.db.Update(ctx, post)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}

	return nil
}

func (s *PostServiceImpl) GetBatch(ctx context.Context, lastID uint32) ([]*models.Post, error) {
	var (
		err    error
		header models.Header
	)

	posts, err := s.db.GetPosts(ctx, lastID)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	for _, post := range posts {
		header, err = s.profileRepo.GetHeader(ctx, post.Header.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get header: %w", err)
		}
		post.Header = header
	}

	return posts, nil
}

func (s *PostServiceImpl) GetBatchFromFriend(ctx context.Context, userID uint32, lastID uint32) ([]*models.Post, error) {
	var (
		err    error
		header models.Header
	)

	friends, err := s.profileRepo.GetFriendsID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get friends: %w", err)
	}

	if len(friends) == 0 {
		return nil, myErr.ErrNoMoreContent
	}

	posts, err := s.db.GetFriendsPosts(ctx, friends, lastID)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	for _, post := range posts {
		header, err = s.profileRepo.GetHeader(ctx, post.Header.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get header: %w", err)
		}
		post.Header = header
	}

	return posts, err
}

func (s *PostServiceImpl) GetPostAuthorID(ctx context.Context, postID uint32) (uint32, error) {
	return s.db.GetPostAuthor(ctx, postID)
}
