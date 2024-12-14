package service

import (
	"context"
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/chat"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ChatService struct {
	repo chat.ChatRepository
}

func NewChatService(repo chat.ChatRepository) *ChatService {
	return &ChatService{
		repo: repo,
	}
}

func (cs *ChatService) GetAllChats(
	ctx context.Context, userID uint32, lastUpdateTime time.Time,
) ([]*models.Chat, error) {
	chats, err := cs.repo.GetChats(ctx, userID, lastUpdateTime)

	if err != nil {
		return nil, fmt.Errorf("get all chats: %w", err)
	}

	return chats, nil
}

func (cs *ChatService) GetChat(
	ctx context.Context, userID uint32, chatID uint32, lastSent time.Time,
) ([]*models.Message, error) {
	messages, err := cs.repo.GetMessages(ctx, userID, chatID, lastSent)
	if err != nil {
		return nil, fmt.Errorf("get all messages: %w", err)
	}

	for i, m := range messages {
		messages[i].CreatedAt = convertTime(m.CreatedAt)
	}

	res := make([]*models.Message, 0, len(messages))
	for _, m := range messages {
		mes := m.FromDto()
		res = append(res, &mes)
	}

	return res, nil
}

func (cs *ChatService) SendNewMessage(
	ctx context.Context, receiver uint32, sender uint32, message *models.MessageContentDto,
) error {
	err := cs.repo.SendNewMessage(ctx, receiver, sender, message)
	if err != nil {
		return fmt.Errorf("send new message: %w", err)
	}

	return nil
}

func convertTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}
