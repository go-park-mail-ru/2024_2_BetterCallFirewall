package service

import (
	"context"
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/chat"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ProfileService interface {
	GetHeader(ctx context.Context, userID uint32) (models.Header, error)
}

type ChatService struct {
	repo           chat.ChatRepository
	profileService ProfileService
}

func NewChatService(repo chat.ChatRepository, service ProfileService) *ChatService {
	return &ChatService{
		repo:           repo,
		profileService: service,
	}
}

func (cs *ChatService) GetAllChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]*models.Chat, error) {
	chats, err := cs.repo.GetChats(ctx, userID, lastUpdateTime)

	if err != nil {
		return nil, fmt.Errorf("get all chats: %w", err)
	}

	/*for _, chat := range chats {
		chat.Receiver, err = cs.profileService.GetHeader(ctx, chat.Receiver.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("get all chats: %w", err)
		}
	}*/

	return chats, nil
}

func (cs *ChatService) GetChat(ctx context.Context, userID uint32, chatID uint32, lastSent time.Time) ([]*models.Message, error) {
	messages, err := cs.repo.GetMessages(ctx, userID, chatID, lastSent)
	if err != nil {
		return nil, fmt.Errorf("get all messages: %w", err)
	}

	return messages, nil
}

func (cs *ChatService) SendNewMessage(ctx context.Context, receiver uint32, sender uint32, message string) error {
	err := cs.repo.SendNewMessage(ctx, receiver, sender, message)
	if err != nil {
		return fmt.Errorf("send new message: %w", err)
	}

	return nil
}
