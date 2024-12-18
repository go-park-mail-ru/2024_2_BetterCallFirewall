package service

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/stickers"
)

type StickerUsecaseImplementation struct {
	repo stickers.Repository
}

func NewStickerUsecase(stickerRepo stickers.Repository) *StickerUsecaseImplementation {
	return &StickerUsecaseImplementation{repo: stickerRepo}
}

func (s StickerUsecaseImplementation) AddNewSticker(ctx context.Context, filepath string, userID uint32) error {
	err := s.repo.AddNewSticker(ctx, filepath, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s StickerUsecaseImplementation) GetAllStickers(ctx context.Context) ([]*models.Picture, error) {
	res, err := s.repo.GetAllStickers(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s StickerUsecaseImplementation) GetMineStickers(ctx context.Context, userID uint32) ([]*models.Picture, error) {
	res, err := s.repo.GetMineStickers(ctx, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
