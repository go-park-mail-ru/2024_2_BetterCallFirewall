package stickers

import (
	"context"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Usecase interface {
	AddNewSticker(ctx context.Context, filepath string, userID uint32) error
	GetAllStickers(ctx context.Context) ([]*models.Picture, error)
	GetMineStickers(ctx context.Context, userID uint32) ([]*models.Picture, error)
}
