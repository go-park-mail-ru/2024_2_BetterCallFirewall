package models

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type Session struct {
	ID        string
	UserID    uint32
	CreatedAt int64
}

func NewSession(userID uint32) (*Session, error) {
	randID := make([]byte, 16)
	_, err := rand.Read(randID)
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}
	return &Session{
		ID:        fmt.Sprintf("%x", randID),
		UserID:    userID,
		CreatedAt: time.Now().Unix(),
	}, nil
}

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, my_err.ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
