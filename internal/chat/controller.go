package chat

import (
	"net/http"
)

type Controller interface {
	SetConnection(w http.ResponseWriter, r *http.Request)
	GetChats(w http.ResponseWriter, r *http.Request)
	GetChatById(w http.ResponseWriter, r *http.Request)
	SendChatMsg(w http.ResponseWriter, r *http.Request, userID uint32)
}
