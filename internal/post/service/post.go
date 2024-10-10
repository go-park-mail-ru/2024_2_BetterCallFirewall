package service

import (
	"fmt"
	"time"

	modelsCon "github.com/2024_2_BetterCallFirewall/internal/content/models"
	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type DB interface {
	Create(post *models.PostDB) (uint32, error)
	Get(postID uint32) (*models.PostDB, error)
	Delete(postID uint32) error
}

type ContentRepo interface {
	Create(content *modelsCon.Content) (uint32, error)
	Get(contentID uint32) (*modelsCon.Content, error)
	Update(content *modelsCon.Content) error
}

type ProfileRepo interface {
	GetName(userID uint32) (string, error)
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
	content := &modelsCon.Content{
		Text:      post.Body,
		FilesPath: post.FilesPath,
		CreatedAt: createTime,
		UpdatedAt: createTime,
	}

	contentID, err := s.contentRepo.Create(content)
	if err != nil {
		return 0, fmt.Errorf("create content failed: %w", err)
	}

	id, err := s.db.Create(&models.PostDB{AuthorID: post.UserID, ContentID: contentID})
	if err != nil {
		return 0, fmt.Errorf("create post failed: %w", err)
	}

	return id, nil
}

func (s *PostServiceImpl) Get(postID uint32) (*models.Post, error) {
	postDB, err := s.db.Get(postID)
	if err != nil {
		return nil, fmt.Errorf("get post failed: %w", err)
	}

	content, err := s.contentRepo.Get(postDB.ContentID)
	if err != nil {
		return nil, fmt.Errorf("get content failed: %w", err)
	}

	author, err := s.profileRepo.GetName(postDB.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("get author failed: %w", err)
	}

	post := &models.Post{
		Header:    author,
		Body:      content.Text,
		FilesPath: content.FilesPath,
		CreatedAt: content.CreatedAt,
	}

	return post, nil
}

func (s *PostServiceImpl) Delete(postID uint32) error {
	err := s.db.Delete(postID)
	if err != nil {
		return fmt.Errorf("delete post failed: %w", err)
	}

	return nil
}

func (s *PostServiceImpl) Update(post *models.Post) error {
	newContent := &modelsCon.Content{
		Text:      post.Body,
		FilesPath: post.FilesPath,
		UpdatedAt: time.Now(),
	}

	err := s.contentRepo.Update(newContent)
	if err != nil {
		return fmt.Errorf("update post failed: %w", err)
	}

	return nil
}
