package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (cr *ChatRepository) GetAllChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]models.Chat, error) {
	var chats []models.Chat

	rows, err := cr.db.QueryContext(ctx, getAllChatBatch, lastUpdateTime, userID)
	defer rows.Close()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErr.ErrNoMoreContent
		}
		return nil, fmt.Errorf("postgres get chats: %w", err)
	}

	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.Receiver.AuthorID, &chat.LastMessage.Content, &chat.LastMessage.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres get chats: %w", err)
		}
		chats = append(chats, chat)
	}

	if len(chats) == 0 {
		return nil, myErr.ErrNoMoreContent
	}

	return chats, nil
}
