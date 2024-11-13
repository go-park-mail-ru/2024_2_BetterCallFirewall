package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type Repo struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (cr *Repo) GetChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]*models.Chat, error) {
	var chats []*models.Chat

	rows, err := cr.db.QueryContext(ctx, getAllChatBatch, userID, pq.FormatTimestamp(lastUpdateTime))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoMoreContent
		}
		return nil, fmt.Errorf("postgres get chats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		chat := &models.Chat{}
		if err := rows.Scan(&chat.Receiver.AuthorID, &chat.Receiver.Author, &chat.Receiver.Avatar, &chat.LastMessage, &chat.LastDate); err != nil {
			return nil, fmt.Errorf("postgres get chats: %w", err)
		}
		chats = append(chats, chat)
	}

	if len(chats) == 0 {
		return nil, my_err.ErrNoMoreContent
	}

	return chats, nil
}

func (cr *Repo) GetMessages(ctx context.Context, userID uint32, chatID uint32, lastSentTime time.Time) ([]*models.Message, error) {
	var messages []*models.Message

	rows, err := cr.db.QueryContext(ctx, getLatestMessagesBatch, userID, chatID, pq.FormatTimestamp(lastSentTime))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoMoreContent
		}
		return nil, fmt.Errorf("postgres get messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		msg := &models.Message{}
		if err := rows.Scan(&msg.Sender, &msg.Receiver, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres get messages: %w", err)
		}
		messages = append(messages, msg)
	}
	if len(messages) == 0 {
		return nil, my_err.ErrNoMoreContent
	}

	return messages, nil

}

func (cr *Repo) SendNewMessage(ctx context.Context, receiver uint32, sender uint32, message string) error {
	_, err := cr.db.ExecContext(ctx, sendNewMessage, receiver, sender, message)
	if err != nil {
		return fmt.Errorf("postgres send new message: %w", err)
	}
	return nil
}
