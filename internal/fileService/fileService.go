package fileService

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) Save(file multipart.File, fileHeader *multipart.FileHeader) (*models.Picture, error) {
	var firstData [1024]byte
	n, err := file.Read(firstData[:])
	if err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	name := string(firstData[:n]) + fileHeader.Filename
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
	log.Println("Save file")

	return nil, nil
}
