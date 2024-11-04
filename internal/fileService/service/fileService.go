package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Repo interface {
	InsertPostFilePath(ctx context.Context, filePath string, postId uint32) error
	InsertProfileFilePath(ctx context.Context, filePath string, profileId uint32) error
	GetProfileFiles(ctx context.Context, profileId uint32) ([]*string, error)
	GetPostFiles(ctx context.Context, postId uint32) (string, error)
	UpdatePostFile(ctx context.Context, filepath string, postId uint32) error
}

type FileService struct {
	repo Repo
}

func NewFileService(repo Repo) *FileService {
	return &FileService{
		repo: repo,
	}
}

func (f *FileService) CreateFile(file multipart.File) (string, error) {
	var (
		fileName = uuid.New().String()
		filePath = fmt.Sprintf("/image/%s", fileName)
		dst, err = os.Create(filePath)
	)

	defer dst.Close()
	if err != nil {
		return "", fmt.Errorf("save file: %w", err)
	}

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("save file: %w", err)
	}

	return filePath, nil
}

func (f *FileService) Download(ctx context.Context, file multipart.File, postId, profileId uint32) error {
	filePath, err := f.CreateFile(file)
	if err != nil {
		return err
	}

	if profileId == 0 {
		err = f.repo.InsertPostFilePath(ctx, filePath, postId)
	} else if postId == 0 {
		err = f.repo.InsertProfileFilePath(ctx, filePath, profileId)
	}

	if err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

func (f *FileService) UpdatePostFile(ctx context.Context, file multipart.File, postId uint32) error {
	filePath, err := f.CreateFile(file)
	if err != nil {
		return err
	}

	err = f.repo.UpdatePostFile(ctx, filePath, postId)
	if err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

func (f *FileService) Upload(ctx context.Context, name string) ([]byte, error) {
	var (
		file, err = os.Open(fmt.Sprintf("/image/%s", name))
		res       []byte
		sl        = make([]byte, 1024)
	)

	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	for n, err := file.Read(sl); err != io.EOF; n, err = file.Read(sl) {
		res = append(res, sl[:n]...)
	}

	return res, nil
}

func (f *FileService) GetPostPicture(ctx context.Context, postID uint32) *models.Picture {
	pic, err := f.repo.GetPostFiles(ctx, postID)
	if err != nil {
		return nil
	}

	res := models.Picture(pic)
	return &res
}

func (f *FileService) GetProfilePictures(ctx context.Context, postID uint32) []*models.Picture {
	pics, err := f.repo.GetProfileFiles(ctx, postID)
	if err != nil {
		return nil
	}

	res := make([]*models.Picture, len(pics))
	for i, pic := range pics {
		myPic := models.Picture(*pic)
		res[i] = &myPic
	}

	return res
}