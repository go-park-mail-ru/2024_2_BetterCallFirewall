package chat

import (
	"context"
	"net/http"
)

type Controller interface {
	SetConnection(w http.ResponseWriter, r *http.Request)
	GetChats(w http.ResponseWriter, r *http.Request)
	GetChatById(w http.ResponseWriter, r *http.Request)
	SendChatMsg(ctx context.Context, reqID string)
}
