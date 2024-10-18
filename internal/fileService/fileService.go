package fileService

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) Upload(file multipart.File) (*models.Picture, error) {
	var firstData [1024]byte
	n, err := file.Read(firstData[:])
	if err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	name := string(firstData[:n])
	hash := md5.Sum([]byte(name))
	hashName := hex.EncodeToString(hash[:])
	dst, err := os.Create(hashName)
	if err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	dir, err := filepath.Abs(hashName)
	if err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	pic := models.Picture(dir)

	return &pic, nil
}

// TODO realize
func (f *FileService) GetPostPicture(postID uint32) (*models.Picture, error) {
	return nil, nil
}
