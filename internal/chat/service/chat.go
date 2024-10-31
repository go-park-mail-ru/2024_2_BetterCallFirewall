package service

import (
	"context"
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ChatRepo interface {
	GetAllChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]models.Chat, error)
}

type ProfileService interface {
	GetHeader(ctx context.Context, userID uint32) (models.Header, error)
}

type ChatService struct {
	repo           ChatRepo
	profileService ProfileService
}

func NewChatService(repo ChatRepo, service ProfileService) *ChatService {
	return &ChatService{
		repo:           repo,
		profileService: service,
	}
}

func (cs *ChatService) GetAllChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]models.Chat, error) {
	chats, err := cs.repo.GetAllChats(ctx, userID, lastUpdateTime)

	if err != nil {
		return nil, fmt.Errorf("get all chats: %w", err)
	}

	for _, chat := range chats {
		chat.Receiver, err = cs.profileService.GetHeader(ctx, chat.Receiver.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get all chats: %w", err)
		}
	}

	return chats, nil
}
