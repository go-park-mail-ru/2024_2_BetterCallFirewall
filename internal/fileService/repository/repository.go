package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type FileRepo struct {
	DB *sql.DB
}

func NewFileRepo(db *sql.DB) *FileRepo {
	repo := &FileRepo{
		DB: db,
	}
	return repo
}

func (fr FileRepo) InsertPostFilePath(ctx context.Context, filePath string, postId uint32) error {
	_, err := fr.DB.ExecContext(ctx, InsertPostFile, filePath, postId)
	if err != nil {
		return fmt.Errorf("insert file: %w", err)
	}
	return nil
}

func (fr FileRepo) InsertProfileFilePath(ctx context.Context, filePath string, profileId uint32) error {
	_, err := fr.DB.ExecContext(ctx, InsertProfileFile, filePath, profileId)
	if err != nil {
		return fmt.Errorf("insert file: %w", err)
	}
	return nil
}

func (fr FileRepo) GetProfileFiles(ctx context.Context, profileId uint32) ([]*string, error) {
	res := make([]*string, 0)
	rows, err := fr.DB.QueryContext(ctx, GetProfileFile, profileId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErr.ErrNoFile
		}
		return nil, fmt.Errorf("get file db: %w", err)

	}
	for rows.Next() {
		var file string
		err := rows.Scan(&file)
		if err != nil {
			return nil, fmt.Errorf("get file db: %w", err)
		}
		res = append(res, &file)
	}
	rows.Close()
	return res, nil
}

func (fr FileRepo) GetPostFiles(ctx context.Context, postId uint32) (string, error) {
	var res string

	if err := fr.DB.QueryRowContext(ctx, GetPostFile, postId).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", myErr.ErrNoFile
		}
		return "", fmt.Errorf("get file db: %w", err)
	}

	return res, nil
}

func (fr FileRepo) UpdatePostFile(ctx context.Context, filepath string, postId uint32) error {
	_, err := fr.DB.ExecContext(ctx, UpdatePostFile, filepath, postId)
	if err != nil {
		return fmt.Errorf("update file: %w", err)
	}

	return nil
}
