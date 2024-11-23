package service

import (
	"context"
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/chat"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type CSATStat interface {
	NewMessage(uint32)
}

type ChatService struct {
	repo chat.ChatRepository
	stat CSATStat
}

func NewChatService(repo chat.ChatRepository, csat CSATStat) *ChatService {
	return &ChatService{
		repo: repo,
		stat: csat,
	}
}

func (cs *ChatService) GetAllChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]*models.Chat, error) {
	chats, err := cs.repo.GetChats(ctx, userID, lastUpdateTime)

	if err != nil {
		return nil, fmt.Errorf("get all chats: %w", err)
	}

	return chats, nil
}

func (cs *ChatService) GetChat(ctx context.Context, userID uint32, chatID uint32, lastSent time.Time) ([]*models.Message, error) {
	messages, err := cs.repo.GetMessages(ctx, userID, chatID, lastSent)
	if err != nil {
		return nil, fmt.Errorf("get all messages: %w", err)
	}

	for i, m := range messages {
		messages[i].CreatedAt = convertTime(m.CreatedAt)
	}

	return messages, nil
}

func (cs *ChatService) SendNewMessage(ctx context.Context, receiver uint32, sender uint32, message string) error {
	err := cs.repo.SendNewMessage(ctx, receiver, sender, message)
	if err != nil {
		return fmt.Errorf("send new message: %w", err)
	}
	cs.stat.NewMessage(sender)
	return nil
}

func convertTime(t time.Time) time.Time {
	newTime, _ := time.Parse("2006-01-02T15:04:05.000000Z", t.Format("2006-01-02T15:04:05.000000Z"))
	return newTime
}
