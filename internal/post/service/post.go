package service

import (
	"context"
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type DB interface {
	Create(ctx context.Context, post *models.Post) (uint32, error)
	Get(ctx context.Context, postID uint32) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, postID uint32) error
	GetPosts(ctx context.Context, lastID uint32) ([]*models.Post, error)
	GetFriendsPosts(ctx context.Context, friendsID []uint32, lastID uint32) ([]*models.Post, error)
	GetPostAuthor(ctx context.Context, postID uint32) (uint32, error)

	CreateCommunityPost(ctx context.Context, post *models.Post, communityID uint32) (uint32, error)
	GetCommunityPosts(ctx context.Context, communityID uint32, lastID uint32) ([]*models.Post, error)

	SetLikeToPost(ctx context.Context, postID uint32, userID uint32) error
	DeleteLikeFromPost(ctx context.Context, postID uint32, userID uint32) error
	GetLikesOnPost(ctx context.Context, postID uint32) (uint32, error)
	CheckLikes(ctx context.Context, postID, userID uint32) (bool, error)
}

type ProfileRepo interface {
	GetHeader(ctx context.Context, userID uint32) (*models.Header, error)
	GetFriendsID(ctx context.Context, userID uint32) ([]uint32, error)
}

type CommunityRepo interface {
	CheckAccess(ctx context.Context, communityID, userID uint32) bool
	GetHeader(ctx context.Context, communityID uint32) (*models.Header, error)
}

type PostServiceImpl struct {
	db            DB
	profileRepo   ProfileRepo
	communityRepo CommunityRepo
}

func NewPostServiceImpl(db DB, profileRepo ProfileRepo, repo CommunityRepo) *PostServiceImpl {
	return &PostServiceImpl{
		db:            db,
		profileRepo:   profileRepo,
		communityRepo: repo,
	}
}

func (s *PostServiceImpl) Create(ctx context.Context, post *models.Post) (uint32, error) {
	id, err := s.db.Create(ctx, post)
	if err != nil {
		return 0, fmt.Errorf("create post: %w", err)
	}

	return id, nil
}

func (s *PostServiceImpl) Get(ctx context.Context, postID, userID uint32) (*models.Post, error) {
	post, err := s.db.Get(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}

	if err := s.setPostFields(ctx, post, userID); err != nil {
		return nil, fmt.Errorf("set post fields: %w", err)
	}

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

func (s *PostServiceImpl) GetBatch(ctx context.Context, lastID, userID uint32) ([]*models.Post, error) {
	posts, err := s.db.GetPosts(ctx, lastID)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	for _, post := range posts {
		if err := s.setPostFields(ctx, post, userID); err != nil {
			return nil, fmt.Errorf("set post fields: %w", err)
		}
	}

	return posts, nil
}

func (s *PostServiceImpl) GetBatchFromFriend(ctx context.Context, userID uint32, lastID uint32) ([]*models.Post, error) {
	friends, err := s.profileRepo.GetFriendsID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get friends: %w", err)
	}

	if len(friends) == 0 {
		return nil, my_err.ErrNoMoreContent
	}

	posts, err := s.db.GetFriendsPosts(ctx, friends, lastID)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	for _, post := range posts {
		if err := s.setPostFields(ctx, post, userID); err != nil {
			return nil, fmt.Errorf("set post fields: %w", err)
		}
	}

	return posts, err
}

func (s *PostServiceImpl) GetPostAuthorID(ctx context.Context, postID uint32) (uint32, error) {
	return s.db.GetPostAuthor(ctx, postID)
}

func (s *PostServiceImpl) CreateCommunityPost(ctx context.Context, post *models.Post) (uint32, error) {
	id, err := s.db.CreateCommunityPost(ctx, post, post.Header.CommunityID)
	if err != nil {
		return 0, fmt.Errorf("create post: %w", err)
	}

	return id, nil
}

func (s *PostServiceImpl) GetCommunityPost(ctx context.Context, communityID, lastID uint32) ([]*models.Post, error) {
	posts, err := s.db.GetCommunityPosts(ctx, communityID, lastID)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	return posts, nil
}

func (s *PostServiceImpl) CheckAccessToCommunity(ctx context.Context, userID uint32, communityID uint32) bool {
	return s.communityRepo.CheckAccess(ctx, userID, communityID)
}

func convertTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}

func (s *PostServiceImpl) SetLikeToPost(ctx context.Context, postID uint32, userID uint32) error {
	err := s.db.SetLikeToPost(ctx, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostServiceImpl) DeleteLikeFromPost(ctx context.Context, postID uint32, userID uint32) error {
	err := s.db.DeleteLikeFromPost(ctx, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostServiceImpl) CheckLikes(ctx context.Context, postID, userID uint32) (bool, error) {
	res, err := s.db.CheckLikes(ctx, postID, userID)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (s *PostServiceImpl) setPostFields(ctx context.Context, post *models.Post, userID uint32) error {
	var (
		header *models.Header
		err    error
	)
	if post.Header.CommunityID == 0 {
		header, err = s.profileRepo.GetHeader(ctx, post.Header.AuthorID)
		if err != nil {
			return fmt.Errorf("get header: %w", err)
		}
	} else {
		header, err = s.communityRepo.GetHeader(ctx, post.Header.CommunityID)
		if err != nil {
			return fmt.Errorf("get community header: %w", err)
		}
	}
	post.Header = *header

	likes, err := s.db.GetLikesOnPost(ctx, post.ID)
	if err != nil {
		return fmt.Errorf("get likes: %w", err)
	}
	post.LikesCount = likes

	liked, err := s.db.CheckLikes(ctx, post.ID, userID)
	if err != nil {
		return fmt.Errorf("check likes: %w", err)
	}
	post.IsLiked = liked

	post.PostContent.CreatedAt = convertTime(post.PostContent.CreatedAt)

	return nil
}
