package models

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type Session struct {
	ID     string
	UserID uint32
}

func NewSession(userID uint32) (*Session, error) {
	randID := make([]byte, 16)
	_, err := rand.Read(randID)
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}
	return &Session{
		ID:     fmt.Sprintf("%x", randID),
		UserID: userID,
	}, nil
}

// SessionKey TODO сделать тип
var (
	SessionKey string = "sessionKey"
)

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, myErr.ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
