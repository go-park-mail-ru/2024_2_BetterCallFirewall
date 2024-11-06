package service

import (
	"context"
	"testing"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type MockRepo struct{}

func (m MockRepo) GetChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]*models.Chat, error) {
	return []*models.Chat{}, nil
}

func (m MockRepo) GetMessages(ctx context.Context, userID uint32, chatID uint32, lastSentTime time.Time) ([]*models.Message, error) {
	return []*models.Message{}, nil
}

func (m MockRepo) SendNewMessage(ctx context.Context, receiver uint32, sender uint32, message string) error {
	return nil
}

type TestStructGetAllChats struct {
	userID         uint32
	lastUpdateTime time.Time
	wantChats      []*models.Chat
	wantErr        error
}

func TestGetAllChats(t *testing.T) {
	chatServ := NewChatService(MockRepo{})

}
