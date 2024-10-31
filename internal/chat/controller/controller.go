package controller

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type ChatService interface {
	GetAllChats(ctx context.Context, userID uint32, lastUpdateTime time.Time) ([]models.Chat, error)
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestId string)
	OutputNoMoreContentJSON(w http.ResponseWriter, requestId string)

	ErrorInternal(w http.ResponseWriter, err error, requestId string)
	ErrorBadRequest(w http.ResponseWriter, err error, requestId string)
	LogError(err error, requestId string)
}

type ChatController struct {
	chatService ChatService
	responder   Responder
}

func NewChatController(service ChatService, responder Responder) *ChatController {
	return &ChatController{
		chatService: service,
		responder:   responder,
	}
}

func (cc *ChatController) GetAllChats(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok     = r.Context().Value("requestID").(string)
		lastTimeQuery = r.URL.Query().Get("lastIme")
		lastTime      time.Time
		err           error
	)

	if !ok {
		cc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	if lastTimeQuery == "" {
		lastTime = time.Now()
	} else {
		lastTime, err = time.Parse("2006-01-02 15:04:05.000", lastTimeQuery)
		if err != nil {
			cc.responder.ErrorBadRequest(w, myErr.ErrWrongDateFormat, reqID)
			return
		}
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		cc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	chats, err := cc.chatService.GetAllChats(r.Context(), sess.UserID, lastTime)
	if errors.Is(err, myErr.ErrNoMoreContent) {
		cc.responder.OutputNoMoreContentJSON(w, reqID)
		return
	}

	if err != nil {
		cc.responder.ErrorInternal(w, err, reqID)
		return
	}

	cc.responder.OutputJSON(w, chats, reqID)
}

// TODO open WS? skrol with batch?
func (cc *ChatController) GetOneChat(w http.ResponseWriter, r *http.Request) {}
