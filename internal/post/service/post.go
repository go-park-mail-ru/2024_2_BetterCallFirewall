package service

import (
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/post/entities"
)

type DB interface {
	Create(post *entities.PostDB) (uint32, error)
	Get(postID uint32) (*models.Post, error)
	Update(post *entities.PostDB) error
	Delete(postID uint32) error
	CheckAccess(profileID uint32, postID uint32) (bool, error)
	GetPosts(lastID uint32, newRequest bool) ([]*models.Post, error)
	GetFriendsPosts(friendsID []uint32, lastID uint32, newRequest bool) ([]*models.Post, error)
}

type ProfileRepo interface {
	GetHeader(userID uint32) (models.Header, error)
	GetFriendsID(userID uint32) ([]uint32, error)
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

func (s *PostServiceImpl) Create(post *models.Post) (uint32, error) {
	createTime := time.Now()

	id, err := s.db.Create(&entities.PostDB{AuthorID: post.AuthorID, Content: post.PostContent.Text, Created: createTime, Updated: createTime})
	if err != nil {
		return 0, fmt.Errorf("create post: %w", err)
	}

	return id, nil
}

func (s *PostServiceImpl) Get(postID uint32) (*models.Post, error) {
	post, err := s.db.Get(postID)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}

	header, err := s.profileRepo.GetHeader(post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("get author: %w", err)
	}
	post.Header = header

	return post, nil
}

func (s *PostServiceImpl) Delete(postID uint32) error {
	err := s.db.Delete(postID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	return nil
}

func (s *PostServiceImpl) Update(post *models.Post) error {
	post.PostContent.UpdatedAt = time.Now()

	err := s.db.Update(&entities.PostDB{ID: post.ID, Content: post.PostContent.Text, Updated: post.PostContent.UpdatedAt})
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}

	return nil
}

func (s *PostServiceImpl) CheckUserAccess(userID uint32, postID uint32) (bool, error) {
	ok, err := s.db.CheckAccess(userID, postID)
	if err != nil {
		return false, fmt.Errorf("check access: %w", err)
	}

	return ok, nil
}

func (s *PostServiceImpl) GetBatch(lastID uint32, newRequest bool) ([]*models.Post, error) {
	posts, err := s.db.GetPosts(lastID, newRequest)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	for _, post := range posts {
		header, err := s.profileRepo.GetHeader(post.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get author: %w", err)
		}
		post.Header = header
	}

	return posts, nil
}

func (s *PostServiceImpl) GetBatchFromFriend(userID uint32, lastID uint32, newRequest bool) ([]*models.Post, error) {
	friends, err := s.profileRepo.GetFriendsID(userID)
	if err != nil {
		return nil, fmt.Errorf("get friends: %w", err)
	}

	posts, err := s.db.GetFriendsPosts(friends, lastID, newRequest)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	for _, post := range posts {
		header, err := s.profileRepo.GetHeader(post.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get author: %w", err)
		}
		post.Header = header
	}

	return posts, nil
}
