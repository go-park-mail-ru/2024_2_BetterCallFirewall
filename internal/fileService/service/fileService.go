package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) Download(ctx context.Context, file multipart.File) (string, error) {
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
