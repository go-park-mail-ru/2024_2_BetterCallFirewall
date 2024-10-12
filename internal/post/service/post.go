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
	Delete(postID uint32) error
	GetContentID(postID uint32) (uint32, error)
}

type ContentRepo interface {
	Create(content *models.Content) (uint32, error)
	Update(content *models.Content) error
}

type ProfileRepo interface {
	GetHeader(userID uint32) (models.Header, error)
}

type PostServiceImpl struct {
	db          DB
	profileRepo ProfileRepo
	contentRepo ContentRepo
}

func NewPostServiceImpl(db DB, profileRepo ProfileRepo, contentRepo ContentRepo) *PostServiceImpl {
	return &PostServiceImpl{
		db:          db,
		profileRepo: profileRepo,
		contentRepo: contentRepo,
	}
}

func (s *PostServiceImpl) Create(post *models.Post) (uint32, error) {
	createTime := time.Now()
	post.PostContent.CreatedAt, post.PostContent.UpdatedAt = createTime, createTime

	contentID, err := s.contentRepo.Create(&post.PostContent)
	if err != nil {
		return 0, fmt.Errorf("create content failed: %w", err)
	}

	id, err := s.db.Create(&entities.PostDB{AuthorID: post.AuthorID, ContentID: contentID})
	if err != nil {
		return 0, fmt.Errorf("create post failed: %w", err)
	}

	return id, nil
}

func (s *PostServiceImpl) Get(postID uint32) (*models.Post, error) {
	post, err := s.db.Get(postID)
	if err != nil {
		return nil, fmt.Errorf("get post failed: %w", err)
	}

	header, err := s.profileRepo.GetHeader(post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("get author failed: %w", err)
	}
	post.Header = header

	return post, nil
}

// TODO проверить доступ юзера
func (s *PostServiceImpl) Delete(postID uint32) error {
	err := s.db.Delete(postID)
	if err != nil {
		return fmt.Errorf("delete post failed: %w", err)
	}

	return nil
}

// TODO проверить доступ юзера
func (s *PostServiceImpl) Update(post *models.Post) error {
	post.PostContent.UpdatedAt = time.Now()
	contentID, err := s.db.GetContentID(post.ID)
	if err != nil {
		return fmt.Errorf("get content failed: %w", err)
	}
	post.PostContent.ID = contentID

	err = s.contentRepo.Update(&post.PostContent)
	if err != nil {
		return fmt.Errorf("update post failed: %w", err)
	}

	return nil
}
