package chat

import (
	"context"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ChatService interface {
	GetAllChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]*models.Chat, error)
	GetChat(ctx context.Context, userID uint32, chatID uint32, lastSentTime time.Time) ([]*models.Message, error)
	SendNewMessage(ctx context.Context, receiver uint32, sender uint32, message string) error
}
