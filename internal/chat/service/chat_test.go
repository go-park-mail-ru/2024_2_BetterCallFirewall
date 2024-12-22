package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

var (
	errMock    = errors.New("mock error")
	createTime = time.Now()
)

type MockRepo struct{}

func (m MockRepo) GetChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]*models.Chat, error) {
	if userID == 0 {
		return nil, errMock
	}
	return []*models.Chat{}, nil
}

func (m MockRepo) GetMessages(
	ctx context.Context, userID uint32, chatID uint32, lastSentTime time.Time,
) ([]*models.MessageDto, error) {
	if userID == 0 || chatID == 0 {
		return nil, errMock
	}
	return []*models.MessageDto{
		{CreatedAt: createTime},
	}, nil
}

func (m MockRepo) SendNewMessage(
	ctx context.Context, receiver uint32, sender uint32, message *models.MessageContentDto,
) error {
	if receiver == 0 || sender == 0 || message.Text == "" {
		return errMock
	}
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

type TestStructGetChat struct {
	userID         uint32
	chatID         uint32
	lastUpdateTime time.Time
	wantMessage    []*models.Message
	wantErr        error
}

func TestGetChat(t *testing.T) {
	chatServ := NewChatService(MockRepo{})
	tests := []TestStructGetChat{
		{
			userID:      0,
			chatID:      10,
			wantMessage: nil,
			wantErr:     errMock,
		},
		{
			userID:      100,
			chatID:      0,
			wantMessage: nil,
			wantErr:     errMock,
		},
		{
			userID:      1,
			chatID:      100,
			wantMessage: []*models.Message{{CreatedAt: convertTime(createTime)}},
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		res, err := chatServ.GetChat(context.Background(), tt.userID, tt.chatID, tt.lastUpdateTime)
		assert.Equal(t, tt.wantMessage, res)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("GetAllChats() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

type TestStructSendNewMessage struct {
	sender   uint32
	receiver uint32
	message  *models.MessageContentDto
	wantErr  error
}

func TestSendNewMessage(t *testing.T) {
	chatServ := NewChatService(MockRepo{})
	tests := []TestStructSendNewMessage{
		{
			sender:   0,
			receiver: 10,
			message:  &models.MessageContentDto{Text: "hello"},
			wantErr:  errMock,
		},
		{
			sender:   10,
			receiver: 0,
			message:  &models.MessageContentDto{Text: "hello"},
			wantErr:  errMock,
		},
		{
			sender:   1,
			receiver: 10,
			message:  &models.MessageContentDto{Text: ""},
			wantErr:  errMock,
		},
		{
			sender:   1,
			receiver: 10,
			message:  &models.MessageContentDto{Text: "hello"},
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		err := chatServ.SendNewMessage(context.Background(), tt.receiver, tt.sender, tt.message)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("GetAllChats() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
