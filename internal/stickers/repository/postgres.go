package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type StickerRepo struct {
	DB *sql.DB
}

func NewStickerRepo(db *sql.DB) *StickerRepo {
	repo := &StickerRepo{
		DB: db,
	}
	return repo
}

func (s StickerRepo) AddNewSticker(ctx context.Context, filepath string, userID uint32) error {
	_, err := s.DB.ExecContext(ctx, InsertNewSticker, filepath, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s StickerRepo) GetAllStickers(ctx context.Context) ([]*models.Picture, error) {
	rows, err := s.DB.QueryContext(ctx, GetAllSticker)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoStickers
		}
		return nil, err
	}
	defer rows.Close()
	var res []*models.Picture
	for rows.Next() {
		var pic models.Picture
		err = rows.Scan(&pic)
		if err != nil {
			return nil, err
		}
		res = append(res, &pic)
	}

	return res, nil
}

func (s StickerRepo) GetMineStickers(ctx context.Context, userID uint32) ([]*models.Picture, error) {
	rows, err := s.DB.QueryContext(ctx, GetUserStickers, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoStickers
		}
		return nil, err
	}
	defer rows.Close()
	var res []*models.Picture
	for rows.Next() {
		var pic models.Picture
		err = rows.Scan(&pic)
		if err != nil {
			return nil, err
		}
		res = append(res, &pic)
	}

	return res, nil
}
