package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

var errMock = errors.New("mock error")

type MockRepo struct{}

func (m MockRepo) GetChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]*models.Chat, error) {
	if userID == 0 {
		return nil, errMock
	}
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
	tests := []TestStructGetAllChats{
		{
			userID:    0,
			wantChats: nil,
			wantErr:   errMock,
		},
		{
			userID:    1,
			wantChats: []*models.Chat{},
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		res, err := chatServ.GetAllChats(context.Background(), tt.userID, tt.lastUpdateTime)
		assert.Equal(t, tt.wantChats, res)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("GetAllChats() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
